package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/stoewer/go-strcase"
	"golang.org/x/exp/slices"
	"magalu.cloud/core"
	mgcSdk "magalu.cloud/sdk"
)

type mgcName string
type tfName string

type attribute struct {
	tfName     tfName
	mgcName    mgcName
	mgcSchema  *mgcSdk.Schema
	tfSchema   schema.Attribute
	attributes mgcAttributes
}

type mgcAttributes map[mgcName]*attribute

// Similar schemas are those with the same type and, depending on the type,
// similar properties or restrictions.
func checkSimilarJsonSchemas(a, b *mgcSdk.Schema) bool {
	if a == b {
		return true
	}

	tA, err := getJsonType(a)
	if err != nil {
		return false
	}

	tB, err := getJsonType(b)
	if err != nil {
		return false
	}

	if tA != tB {
		return false
	}

	switch tA {
	default:
		return true
	case "string":
		// Relax if one of them doesn't specify a format
		return a.Format == b.Format || a.Format == "" || b.Format == ""
	case "array":
		return checkSimilarJsonSchemas((*mgcSdk.Schema)(a.Items.Value), (*mgcSdk.Schema)(b.Items.Value))
	case "object":
		// TODO: should we compare Required? I don't think so, but it may be a problem
		if len(a.Properties) != len(b.Properties) {
			return false
		}
		for k, refA := range a.Properties {
			refB := b.Properties[k]
			if refB == nil {
				return false
			}
			if !checkSimilarJsonSchemas((*mgcSdk.Schema)(refA.Value), (*mgcSdk.Schema)(refB.Value)) {
				return false
			}
		}
		// TODO: handle additionalProperties and property patterns
		return true
	}
}

type attributeModifiers struct {
	isRequired                 bool
	isOptional                 bool
	isComputed                 bool
	useStateForUnknown         bool
	requiresReplaceWhenChanged bool
	getChildModifiers          func(mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers
}

func addMgcSchemaAttributes(
	attributes mgcAttributes,
	mgcSchema *mgcSdk.Schema,
	getModifiers func(mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers,
	resourceName string,
	ctx context.Context,
) error {
	for k, ref := range mgcSchema.Properties {
		mgcName := mgcName(k)
		mgcPropSchema := (*mgcSdk.Schema)(ref.Value)
		if ca, ok := attributes[mgcName]; ok {
			if !checkSimilarJsonSchemas(ca.mgcSchema, mgcPropSchema) {
				// Ignore update value in favor of create value (This is probably a bug with the API)
				tflog.Error(ctx, fmt.Sprintf("[resource] schema for `%s`: ignoring DIFFERENT attribute `%s`:\nOLD=%+v\nNEW=%+v", resourceName, k, ca.mgcSchema, mgcPropSchema))
				continue
			}
			tflog.Debug(ctx, fmt.Sprintf("[resource] schema for `%s`: ignoring already computed attribute `%s` ", resourceName, k))
			continue
		}

		tfSchema, childAttributes, err := mgcToTFSchema(mgcPropSchema, getModifiers(mgcSchema, mgcName), resourceName, ctx)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("[resource] schema for %q attribute %q schema: %+v; error: %s", resourceName, k, mgcPropSchema, err))
			return fmt.Errorf("attribute %q, error=%s", k, err)
		}

		attr := &attribute{
			tfName:     tfNameFromMgc(mgcName),
			mgcName:    mgcName,
			mgcSchema:  mgcPropSchema,
			tfSchema:   tfSchema,
			attributes: childAttributes,
		}
		attributes[mgcName] = attr
		tflog.Debug(ctx, fmt.Sprintf("[resource] schema for %q attribute %q: %+v", resourceName, k, attr))
	}

	return nil
}

func getInputChildModifiers(mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
	k := string(mgcName)
	isRequired := slices.Contains(mgcSchema.Required, k)
	return attributeModifiers{
		isRequired:                 isRequired,
		isOptional:                 !isRequired,
		isComputed:                 false, // This is being set to false because the parent may already be Computed, no further logic is needed here
		useStateForUnknown:         true,
		requiresReplaceWhenChanged: false,
		getChildModifiers:          getInputChildModifiers,
	}
}

func (r *MgcResource) getCreateParamsModifiers(mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
	k := string(mgcName)
	isRequired := slices.Contains(mgcSchema.Required, k)
	isComputed := !isRequired
	if isComputed {
		readSchema := r.read.ResultSchema().Properties[k]
		if readSchema == nil {
			isComputed = false
		} else {
			// If not required and present in read it can be compute
			isComputed = checkSimilarJsonSchemas((*core.Schema)(readSchema.Value), (*core.Schema)(mgcSchema.Properties[string(mgcName)].Value))
		}
	}

	return attributeModifiers{
		isRequired:                 isRequired,
		isOptional:                 !isRequired,
		isComputed:                 isComputed,
		useStateForUnknown:         false,
		requiresReplaceWhenChanged: r.update.ParametersSchema().Properties[k] == nil,
		getChildModifiers:          getInputChildModifiers,
	}
}

func (r *MgcResource) getUpdateParamsModifiers(mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
	k := string(mgcName)
	isCreated := r.create.ResultSchema().Properties[k] != nil
	required := slices.Contains(mgcSchema.Required, k)

	return attributeModifiers{
		isRequired:                 required && !isCreated,
		isOptional:                 !required && !isCreated,
		isComputed:                 !required || isCreated,
		useStateForUnknown:         true,
		requiresReplaceWhenChanged: false,
		getChildModifiers:          getInputChildModifiers,
	}
}

func (r *MgcResource) getDeleteParamsModifiers(mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
	// For now we consider all delete params as optionals, we need to think a way for the user to define
	// required delete params
	return attributeModifiers{
		isRequired:                 false,
		isOptional:                 true,
		isComputed:                 false,
		useStateForUnknown:         true,
		requiresReplaceWhenChanged: false,
		getChildModifiers:          getInputChildModifiers,
	}
}

func getResultModifiers(mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
	return attributeModifiers{
		isRequired:                 false,
		isOptional:                 false,
		isComputed:                 true,
		useStateForUnknown:         false,
		requiresReplaceWhenChanged: false,
		getChildModifiers:          getResultModifiers,
	}
}

func (r *MgcResource) readInputAttributes(ctx context.Context) diag.Diagnostics {
	d := diag.Diagnostics{}
	if len(r.inputAttr) != 0 {
		return d
	}
	tflog.Debug(ctx, fmt.Sprintf("[resource] schema for `%s`: reading input attributes", r.name))

	input := mgcAttributes{}
	err := addMgcSchemaAttributes(
		input,
		r.create.ParametersSchema(),
		r.getCreateParamsModifiers,
		r.name,
		ctx,
	)
	if err != nil {
		d.AddError("could not create TF input attributes", err.Error())
		return d
	}

	err = addMgcSchemaAttributes(
		input,
		r.update.ParametersSchema(),
		r.getUpdateParamsModifiers,
		r.name,
		ctx,
	)
	if err != nil {
		d.AddError("could not create TF input attributes", err.Error())
		return d
	}

	err = addMgcSchemaAttributes(
		input,
		r.delete.ParametersSchema(),
		r.getDeleteParamsModifiers,
		r.name,
		ctx,
	)
	if err != nil {
		d.AddError("could not create TF input attributes", err.Error())
		return d
	}

	r.inputAttr = input
	return d
}

func (r *MgcResource) readOutputAttributes(ctx context.Context) diag.Diagnostics {
	d := diag.Diagnostics{}
	if len(r.outputAttr) != 0 {
		return d
	}
	tflog.Debug(ctx, fmt.Sprintf("[resource] schema for `%s`: reading output attributes", r.name))

	output := mgcAttributes{}
	err := addMgcSchemaAttributes(
		output,
		r.create.ResultSchema(),
		getResultModifiers,
		r.name,
		ctx,
	)
	if err != nil {
		d.AddError("could not create TF output attributes", err.Error())
		return d
	}
	err = addMgcSchemaAttributes(
		output,
		r.read.ResultSchema(),
		getResultModifiers,
		r.name,
		ctx,
	)
	if err != nil {
		d.AddError("could not create TF output attributes", err.Error())
		return d
	}

	r.outputAttr = output
	return d
}

func (r *MgcResource) generateTFSchema(ctx context.Context) (tfSchema schema.Schema, d diag.Diagnostics) {
	var tfsa map[tfName]schema.Attribute
	tfsa, d = r.generateTFAttributes(ctx)
	if d.HasError() {
		return
	}

	tfSchema = schema.Schema{Attributes: map[string]schema.Attribute{}}
	for tfName, tfAttr := range tfsa {
		tfSchema.Attributes[string(tfName)] = tfAttr
	}
	return
}

func (r *MgcResource) generateTFAttributes(ctx context.Context) (tfa map[tfName]schema.Attribute, d diag.Diagnostics) {
	d.Append(r.readInputAttributes(ctx)...)
	if d.HasError() {
		return
	}
	d.Append(r.readOutputAttributes(ctx)...)
	if d.HasError() {
		return
	}

	tfa = map[tfName]schema.Attribute{}
	for name, iattr := range r.inputAttr {
		// Split attributes that differ between input/output
		if oattr := r.outputAttr[name]; oattr != nil {
			if !checkSimilarJsonSchemas(oattr.mgcSchema, iattr.mgcSchema) {
				os, _ := oattr.mgcSchema.MarshalJSON()
				is, _ := iattr.mgcSchema.MarshalJSON()
				tflog.Debug(ctx, fmt.Sprintf("[resource] schema for `%s`: attribute `%s` differs between input and output. input: %s - output %s", r.name, name, is, os))
				iattr.tfName = iattr.tfName.asDesired()
				oattr.tfName = oattr.tfName.asCurrent()
			}
		}

		tfa[iattr.tfName] = iattr.tfSchema
	}

	for _, oattr := range r.outputAttr {
		// If they don't differ and it's already created skip
		if _, ok := tfa[oattr.tfName]; ok {
			continue
		}

		tfa[oattr.tfName] = oattr.tfSchema
	}

	return
}

func mgcToTFSchema(mgcSchema *mgcSdk.Schema, m attributeModifiers, resourceName string, ctx context.Context) (schema.Attribute, mgcAttributes, error) {
	// TODO: Handle default values

	t, err := getJsonType(mgcSchema)
	if err != nil {
		return nil, nil, err
	}
	description := mgcSchema.Description

	switch t {
	case "string":
		// I wanted to use an interface to define the modifiers regardless of the attr type
		// but couldn't find the interface, it seems everything is redefined for each type
		// https://github.com/hashicorp/terraform-plugin-framework/blob/main/internal/fwschema/fwxschema/attribute_plan_modification.go
		mod := []planmodifier.String{}
		if m.useStateForUnknown {
			mod = append(mod, stringplanmodifier.UseStateForUnknown())
		}
		if m.requiresReplaceWhenChanged {
			mod = append(mod, stringplanmodifier.RequiresReplace())
		}
		return schema.StringAttribute{
			Description:   description,
			Required:      m.isRequired,
			Optional:      m.isOptional,
			Computed:      m.isComputed,
			PlanModifiers: mod,
		}, nil, nil
	case "number":
		mod := []planmodifier.Number{}
		if m.useStateForUnknown {
			mod = append(mod, numberplanmodifier.UseStateForUnknown())
		}
		if m.requiresReplaceWhenChanged {
			mod = append(mod, numberplanmodifier.RequiresReplace())
		}
		return schema.NumberAttribute{
			Description:   description,
			Required:      m.isRequired,
			Optional:      m.isOptional,
			Computed:      m.isComputed,
			PlanModifiers: mod,
		}, nil, nil
	case "integer":
		mod := []planmodifier.Int64{}
		if m.useStateForUnknown {
			mod = append(mod, int64planmodifier.UseStateForUnknown())
		}
		if m.requiresReplaceWhenChanged {
			mod = append(mod, int64planmodifier.RequiresReplace())
		}
		return schema.Int64Attribute{
			Description:   description,
			Required:      m.isRequired,
			Optional:      m.isOptional,
			Computed:      m.isComputed,
			PlanModifiers: mod,
		}, nil, nil
	case "boolean":
		mod := []planmodifier.Bool{}
		if m.useStateForUnknown {
			mod = append(mod, boolplanmodifier.UseStateForUnknown())
		}
		if m.requiresReplaceWhenChanged {
			mod = append(mod, boolplanmodifier.RequiresReplace())
		}
		return schema.BoolAttribute{
			Description:   description,
			Required:      m.isRequired,
			Optional:      m.isOptional,
			Computed:      m.isComputed,
			PlanModifiers: mod,
		}, nil, nil
	case "array":
		mgcItemSchema := (*core.Schema)(mgcSchema.Items.Value)
		elemAttr, elemAttrs, err := mgcToTFSchema(mgcItemSchema, m.getChildModifiers(mgcItemSchema, "0"), resourceName, ctx)
		if err != nil {
			return nil, nil, err
		}

		childAttrs := mgcAttributes{"0": &attribute{
			tfName:     "0",
			mgcName:    "0",
			mgcSchema:  mgcItemSchema,
			tfSchema:   elemAttr,
			attributes: elemAttrs,
		}}

		mod := []planmodifier.List{}
		if m.requiresReplaceWhenChanged {
			mod = append(mod, listplanmodifier.RequiresReplace())
		}
		if m.useStateForUnknown {
			mod = append(mod, listplanmodifier.UseStateForUnknown())
		}

		// TODO: How will we handle List of Lists? Does it need to be handled at all? Does the
		// 'else' branch already cover that correctly?
		if objAttr, ok := elemAttr.(schema.SingleNestedAttribute); ok {
			// This type assertion will/should NEVER fail, according to TF code
			nestedObj, ok := objAttr.GetNestedObject().(schema.NestedAttributeObject)
			if !ok {
				return nil, nil, fmt.Errorf("failed TF GetNestedObject")
			}
			return schema.ListNestedAttribute{
				NestedObject:  nestedObj,
				Description:   description,
				Required:      m.isRequired,
				Optional:      m.isOptional,
				Computed:      m.isComputed,
				PlanModifiers: mod,
			}, childAttrs, nil
		} else {
			return schema.ListAttribute{
				ElementType:   elemAttr.GetType(),
				Description:   description,
				Required:      m.isRequired,
				Optional:      m.isOptional,
				Computed:      m.isComputed,
				PlanModifiers: mod,
			}, childAttrs, nil
		}
	case "object":
		mgcAttributes := mgcAttributes{}
		err := addMgcSchemaAttributes(mgcAttributes, mgcSchema, m.getChildModifiers, resourceName, ctx)
		if err != nil {
			return nil, nil, err
		}
		tfAttributes := map[string]schema.Attribute{}
		for _, attr := range mgcAttributes {
			tfAttributes[string(attr.tfName)] = attr.tfSchema
		}

		mod := []planmodifier.Object{}
		if m.requiresReplaceWhenChanged {
			mod = append(mod, objectplanmodifier.RequiresReplace())
		}
		if m.useStateForUnknown {
			mod = append(mod, objectplanmodifier.UseStateForUnknown())
		}
		return schema.SingleNestedAttribute{
			Attributes:    tfAttributes,
			Description:   description,
			Required:      m.isRequired,
			Optional:      m.isOptional,
			Computed:      m.isComputed,
			PlanModifiers: mod,
		}, mgcAttributes, nil
	default:
		return nil, nil, fmt.Errorf("type %q not supported", t)
	}
}

func tfNameFromMgc(n mgcName) tfName {
	return tfName(strcase.SnakeCase(string(n)))
}

func (n tfName) asDesired() tfName {
	return "desired_" + n
}

func (n tfName) asCurrent() tfName {
	return "current_" + n
}
