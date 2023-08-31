package provider

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

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
	mgcSdk "magalu.cloud/sdk"
)

type attribute struct {
	tfName                     string
	mgcSchema                  *mgcSdk.Schema
	attributes                 map[string]*attribute
	isID                       bool
	isRequired                 bool
	isOptional                 bool
	isComputed                 bool
	useStateForUnknown         bool
	requiresReplaceWhenChanged bool
}

var idRexp = regexp.MustCompile(`(^id$|_id$)`)

func getAttribute(
	name string,
	schema *mgcSdk.Schema,
	parent *mgcSdk.Schema,
	calcIsRequired func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool,
	calcIsOptional func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool,
	calcIsComputed func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool,
	calcRequiresReplaceWhenChanged func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool,
	path []string,
	useStateForUnknown bool,
) *attribute {
	// TODO: Better handle ID attributions
	// Consider other ID elements as ID if "id" doesn't exist
	isID := name == "id" || idRexp.MatchString(name)

	attr := &attribute{
		tfName:             kebabToSnakeCase(name),
		mgcSchema:          schema,
		isID:               isID,
		useStateForUnknown: useStateForUnknown,
		attributes:         map[string]*attribute{},
	}

	requiredByParent := slices.Contains(parent.Required, name)

	attr.isRequired = calcIsRequired(attr, name, path, isID, requiredByParent)
	attr.isOptional = calcIsOptional(attr, name, path, isID, requiredByParent)
	attr.isComputed = calcIsComputed(attr, name, path, isID, requiredByParent)
	attr.requiresReplaceWhenChanged = calcRequiresReplaceWhenChanged(attr, name, path, isID, requiredByParent)

	for propName, propRef := range schema.Properties {
		propAttr := getAttribute(
			propName,
			(*mgcSdk.Schema)(propRef.Value),
			schema,
			calcIsRequired,
			calcIsOptional,
			calcIsComputed,
			calcRequiresReplaceWhenChanged,
			append(path, name),
			useStateForUnknown,
		)
		attr.attributes[propName] = propAttr
	}

	if schema.Items != nil && schema.Items.Value != nil {
		itemsAttr := getAttribute(
			name,
			(*mgcSdk.Schema)(schema.Items.Value),
			parent, // Items isn't really a child, so we use the same parent
			calcIsRequired,
			calcIsOptional,
			calcIsComputed,
			calcRequiresReplaceWhenChanged,
			path,
			useStateForUnknown,
		)
		attr.attributes["0"] = itemsAttr
	}

	// TODO: Handle OneOf, AnyOf, AllOf...
	return attr
}

func lookupSubProperty(schema *mgcSdk.Schema, path []string, name string) (*mgcSdk.Schema, bool) {
	pathLen := len(path)
	strLookup := name
	if pathLen > 0 {
		strLookup = path[0]
	}

	propRef, ok := schema.Properties[strLookup]
	if !ok {
		if schema.Items != nil {
			propRef, ok = schema.Items.Value.Properties[strLookup]
			if !ok {
				return nil, false
			}
		} else {
			return nil, false
		}
	}

	if pathLen == 0 {
		return (*mgcSdk.Schema)(propRef.Value), true
	} else {
		return lookupSubProperty((*mgcSdk.Schema)(propRef.Value), path[1:], name)
	}
}

func fillCreateParamsAttributes(
	paramsSchema *mgcSdk.Schema,
	readResultSchema *mgcSdk.Schema,
	updateParamsSchema *mgcSdk.Schema,
	dst map[string]*attribute,
) {
	for name, ref := range paramsSchema.Properties {
		attr := getAttribute(
			name,
			(*mgcSdk.Schema)(ref.Value),
			paramsSchema,
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool {
				return requiredByParent
			}, // isRequired
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool {
				return !requiredByParent
			}, // isOptional
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { // isComputed
				if requiredByParent {
					return false
				}

				_, isInReadResult := lookupSubProperty(readResultSchema, path, name)
				return isInReadResult
			},
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { // requiresReplaceWhenChanged
				_, isInUpdateParams := lookupSubProperty(updateParamsSchema, path, name)
				return !isInUpdateParams
			},
			[]string{},
			false,
		)
		attr.isID = false // Force no ID on Create, even if, somehow, it came out as true
		dst[name] = attr
	}
}

func fillUpdateParamsAttributes(
	paramsSchema *mgcSdk.Schema,
	ctx context.Context,
	resourceName string,
	dst map[string]*attribute,
) {
	for name, ref := range paramsSchema.Properties {
		if ca, ok := dst[name]; ok {
			us := ref.Value
			if !reflect.DeepEqual(ca.mgcSchema, (*mgcSdk.Schema)(us)) {
				// Ignore update value in favor of create value (This is probably a bug with the API)
				// TODO: Ignore default values when verifying equality
				// TODO: Don't forget to add the path when using recursion
				// err := fmt.Sprintf("[resource] schema for `%s`: input attribute `%s` is different between create and update - create: %+v - update: %+v ", r.name, attr, ca.schema, us)
				// d.AddError("Attribute schema is different between create and update schemas", err)
				continue
			}
			tflog.Debug(ctx, fmt.Sprintf("[resource] schema for `%s`: ignoring already computed attribute `%s` ", resourceName, name))
			continue
		}

		attr := getAttribute(
			name,
			(*mgcSdk.Schema)(ref.Value),
			paramsSchema,
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { // isRequired
				return requiredByParent && !isID
			},
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { // isOptional
				return !requiredByParent && !isID
			},
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { // isComputed
				return !requiredByParent || isID
			},
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { // requiresReplaceWhenChanged
				return false
			},
			[]string{},
			true,
		)
		dst[name] = attr
		tflog.Debug(ctx, fmt.Sprintf("[resource] schema for `%s`: attribute `%s` info - %+v", resourceName, name, attr))
	}
}

func fillCreateResultAttributes(
	resultSchema *mgcSdk.Schema,
	resourceName string,
	ctx context.Context,
	dst map[string]*attribute,
) {
	for name, ref := range resultSchema.Properties {
		attr := getAttribute(
			name,
			(*mgcSdk.Schema)(ref.Value),
			resultSchema,
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { return false }, // isRequired
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { return false }, // isOptional
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { return true },  // isComputed
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { // requiresReplaceWhenChanged
				return false // This one is useless in this case
			},
			[]string{},
			true,
		)
		dst[name] = attr
		tflog.Debug(ctx, fmt.Sprintf("[resource] schema for `%s`: attribute `%s` info - %+v", resourceName, name, attr))
	}
}

func fillReadResultAttributes(
	resultSchema *mgcSdk.Schema,
	resourceName string,
	ctx context.Context,
	dst map[string]*attribute,
) {
	for name, ref := range resultSchema.Properties {
		if ra, ok := dst[name]; ok {
			rs := ref.Value
			if !reflect.DeepEqual(ra.mgcSchema, (*mgcSdk.Schema)(rs)) {
				// Ignore read value in favor of create result value (This is probably a bug with the API)
				// err := fmt.Sprintf("[resource] schema for `%s`: output attribute `%s` is different between create result and read - create result: %+v - read: %+v ", r.name, attr, ra.schema, rs)
				// d.AddError("Attribute schema is different between create result and read schemas", err)
				continue
			}
			tflog.Debug(ctx, fmt.Sprintf("[resource] schema for `%s`: ignoring already computed attribute `%s` ", resourceName, name))
			continue
		}

		attr := getAttribute(
			name,
			(*mgcSdk.Schema)(ref.Value),
			resultSchema,
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { return false }, // isRequired
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { return false }, // isOptional
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { return true },  // isComputed
			func(attr *attribute, name string, path []string, isID bool, requiredByParent bool) bool { // requiresReplaceWhenChanged
				return false // This one is useless in this case
			},
			[]string{},
			true,
		)
		attr.isID = false
		dst[name] = attr
		tflog.Debug(ctx, fmt.Sprintf("[resource] schema for `%s`: attribute `%s` info - %+v", resourceName, name, attr))
	}
}

func (r *MgcResource) readInputAttributes(ctx context.Context) diag.Diagnostics {
	d := diag.Diagnostics{}
	if len(r.inputAttr) != 0 {
		return d
	}
	tflog.Debug(ctx, fmt.Sprintf("[resource] schema for `%s`: reading input attributes", r.name))

	inputAttributes := map[string]*attribute{}
	fillCreateParamsAttributes(
		r.create.ParametersSchema(),
		r.read.ResultSchema(),
		r.update.ParametersSchema(),
		inputAttributes,
	)
	fillUpdateParamsAttributes(
		r.update.ParametersSchema(),
		ctx,
		r.name,
		inputAttributes,
	)

	r.inputAttr = inputAttributes
	return d
}

func (r *MgcResource) readOutputAttributes(ctx context.Context) diag.Diagnostics {
	d := diag.Diagnostics{}
	if len(r.outputAttr) != 0 {
		return d
	}
	tflog.Debug(ctx, fmt.Sprintf("[resource] schema for `%s`: reading output attributes", r.name))

	outputAttributes := map[string]*attribute{}
	fillCreateResultAttributes(
		r.create.ResultSchema(),
		r.name,
		ctx,
		outputAttributes,
	)
	fillReadResultAttributes(
		r.read.ResultSchema(),
		r.name,
		ctx,
		outputAttributes,
	)
	r.outputAttr = outputAttributes
	return d
}

func (r *MgcResource) generateTFAttributes(ctx context.Context) (*map[string]schema.Attribute, diag.Diagnostics) {
	d := diag.Diagnostics{}
	d.Append(r.readInputAttributes(ctx)...)
	d.Append(r.readOutputAttributes(ctx)...)

	tfa := map[string]schema.Attribute{}
	for name, iattr := range r.inputAttr {
		// Split attributes that differ between input/output
		if oattr := r.outputAttr[name]; oattr != nil && !iattr.isID {
			if !reflect.DeepEqual(oattr.mgcSchema, iattr.mgcSchema) {
				os, _ := oattr.mgcSchema.MarshalJSON()
				is, _ := iattr.mgcSchema.MarshalJSON()
				tflog.Debug(ctx, fmt.Sprintf("[resource] schema for `%s`: attribute `%s` differs between input and output. input: %s - output %s", r.name, name, is, os))
				iattr.tfName = kebabToSnakeCase("desired_" + iattr.tfName)
				oattr.tfName = kebabToSnakeCase("current_" + oattr.tfName)
			}
		}

		at := sdkToTerraformAttribute(ctx, iattr, &d)
		// TODO: This shouldn't happen after we handle complex types like slices and objects
		// TODO: Remove debug log
		if at == nil {
			err := fmt.Sprintf("[resource] schema for `%s`: unable to create terraform attribute `%s` - data: %+v", r.name, iattr.tfName, iattr)
			tflog.Debug(ctx, err)
			// TODO: Uncomment the error
			// d.AddError("Unknown attribute type", err)
			continue
		}
		tflog.Debug(ctx, fmt.Sprintf("[resource] schema for `%s`: terraform input attribute `%s` created", r.name, iattr.tfName))
		tfa[iattr.tfName] = at
	}

	for _, oattr := range r.outputAttr {
		// If they don't differ and it's already created skip
		if _, ok := tfa[oattr.tfName]; ok {
			continue
		}

		at := sdkToTerraformAttribute(ctx, oattr, &d)
		if at == nil {
			// TODO: This shouldn't happen after we handle complex types like slices and objects
			// TODO: Remove debug log
			err := fmt.Sprintf("[resource] schema for `%s`: unable to create terraform attribute `%s` - data: %+v", r.name, oattr.tfName, oattr)
			tflog.Debug(ctx, err)
			// TODO: Uncomment the error
			// d.AddError("Unknown attribute type", err)
			continue
		}
		tfa[oattr.tfName] = at
	}

	return &tfa, d
}

func sdkToTerraformAttribute(ctx context.Context, c *attribute, di *diag.Diagnostics) schema.Attribute {
	if c.mgcSchema == nil || c == nil {
		di.AddError("Invalid attribute pointer", fmt.Sprintf("ERROR invalid pointer, attribute pointer is nil %v %v", c.mgcSchema, c))
		return nil
	}

	conv := tfStateConverter{
		ctx:  ctx,
		diag: di,
	}

	// TODO: Handle default values

	value := c.mgcSchema
	t := conv.getAttributeType(value)
	if di.HasError() {
		return nil
	}

	switch t {
	case "string":
		// I wanted to use an interface to define the modifiers regardless of the attr type
		// but couldn't find the interface, it seems everything is redefined for each type
		// https://github.com/hashicorp/terraform-plugin-framework/blob/main/internal/fwschema/fwxschema/attribute_plan_modification.go
		mod := []planmodifier.String{}
		if c.useStateForUnknown {
			mod = append(mod, stringplanmodifier.UseStateForUnknown())
		}
		if c.requiresReplaceWhenChanged {
			mod = append(mod, stringplanmodifier.RequiresReplace())
		}
		return schema.StringAttribute{
			Description:   value.Description,
			Required:      c.isRequired,
			Optional:      c.isOptional,
			Computed:      c.isComputed,
			PlanModifiers: mod,
		}
	case "number":
		mod := []planmodifier.Number{}
		if c.useStateForUnknown {
			mod = append(mod, numberplanmodifier.UseStateForUnknown())
		}
		if c.requiresReplaceWhenChanged {
			mod = append(mod, numberplanmodifier.RequiresReplace())
		}
		return schema.NumberAttribute{
			Description:   value.Description,
			Required:      c.isRequired,
			Optional:      c.isOptional,
			Computed:      c.isComputed,
			PlanModifiers: mod,
		}
	case "integer":
		mod := []planmodifier.Int64{}
		if c.useStateForUnknown {
			mod = append(mod, int64planmodifier.UseStateForUnknown())
		}
		if c.requiresReplaceWhenChanged {
			mod = append(mod, int64planmodifier.RequiresReplace())
		}
		return schema.Int64Attribute{
			Description:   value.Description,
			Required:      c.isRequired,
			Optional:      c.isOptional,
			Computed:      c.isComputed,
			PlanModifiers: mod,
		}
	case "boolean":
		mod := []planmodifier.Bool{}
		if c.useStateForUnknown {
			mod = append(mod, boolplanmodifier.UseStateForUnknown())
		}
		if c.requiresReplaceWhenChanged {
			mod = append(mod, boolplanmodifier.RequiresReplace())
		}
		return schema.BoolAttribute{
			Description:   value.Description,
			Required:      c.isRequired,
			Optional:      c.isOptional,
			Computed:      c.isComputed,
			PlanModifiers: mod,
		}
	case "array":
		elemAttr := sdkToTerraformAttribute(ctx, c.attributes["0"], di)
		if elemAttr == nil {
			return nil
		}
		mod := []planmodifier.List{}
		if c.requiresReplaceWhenChanged {
			mod = append(mod, listplanmodifier.RequiresReplace())
		}

		// TODO: How will we handle List of Lists? Does it need to be handled at all? Does the
		// 'else' branch already cover that correctly?
		if objAttr, ok := elemAttr.(schema.SingleNestedAttribute); ok {
			// This type assertion will/should NEVER fail, according to TF code
			nestedObj, ok := objAttr.GetNestedObject().(schema.NestedAttributeObject)
			if !ok {
				return nil
			}
			return schema.ListNestedAttribute{
				NestedObject:  nestedObj,
				Description:   value.Description,
				Required:      c.isRequired,
				Optional:      c.isOptional,
				Computed:      c.isComputed,
				PlanModifiers: mod,
			}
		} else {
			return schema.ListAttribute{
				ElementType:   elemAttr.GetType(),
				Description:   value.Description,
				Required:      c.isRequired,
				Optional:      c.isOptional,
				Computed:      c.isComputed,
				PlanModifiers: mod,
			}
		}
	case "object":
		children := map[string]schema.Attribute{}
		for _, child := range c.attributes {
			childAttr := sdkToTerraformAttribute(ctx, child, di)
			if childAttr == nil {
				// Should not happen
				continue
			}
			children[child.tfName] = childAttr
		}
		mod := []planmodifier.Object{}
		if c.requiresReplaceWhenChanged {
			mod = append(mod, objectplanmodifier.RequiresReplace())
		}
		return schema.SingleNestedAttribute{
			Attributes:    children,
			Description:   value.Description,
			Required:      c.isRequired,
			Optional:      c.isOptional,
			Computed:      c.isComputed,
			PlanModifiers: mod,
		}
	default:
		return nil
	}
}

func kebabToSnakeCase(n string) string {
	return strcase.SnakeCase(n)
}
