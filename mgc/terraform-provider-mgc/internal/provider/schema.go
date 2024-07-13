package provider

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"slices"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/stoewer/go-strcase"
	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
	mgcSdk "magalu.cloud/sdk"
)

type mgcName string
type tfName string

type resAttrInfo struct {
	tfName             tfName
	mgcName            mgcName
	mgcSchema          *mgcSdk.Schema
	tfSchema           schema.Attribute
	currentCounterpart *resAttrInfo
	childAttributes    resAttrInfoMap
	state              TerraformParams
}

type resAttrInfoMap map[mgcName]*resAttrInfo

func (m resAttrInfoMap) get(name tfName) (*resAttrInfo, bool) {
	for _, attr := range m {
		if attr.tfName == name {
			return attr, true
		}
	}
	return nil, false
}

type resAttrInfoTree struct {
	createInput  resAttrInfoMap
	createOutput resAttrInfoMap

	// No need for 'read' input, as it may have extra unneeded parameters (like 'expand', for some products)
	// which would taint the resource attributes. Parameters absolutely needed for Resource ID will be
	// present in the 'create' output and 'update' input
	readOutput resAttrInfoMap

	updateInput resAttrInfoMap
	// No need for 'update' output, some resources don't have them and they'd be the same
	// as the 'read' output anyway (or be empty)

	deleteInput resAttrInfoMap
	// No need for 'delete' output (basically always empty)

	propertySetterInputs  map[mgcName]map[*core.Schema]resAttrInfoMap
	propertySetterOutputs map[mgcName]map[*core.Schema]resAttrInfoMap

	// Input is the aggregate of all input attributes without duplicates
	input resAttrInfoMap
	// Output is the aggregate of all output attributes without duplicates
	output resAttrInfoMap
}

func (t resAttrInfoTree) getTFInputFirst(name tfName) (*resAttrInfo, bool) {
	if i, ok := t.input.get(name); ok {
		return i, ok
	}
	if o, ok := t.output.get(name); ok {
		return o, ok
	}
	return nil, false
}

type resAttrInfoGenMetadata struct {
	schema    *mgcSdk.Schema
	modifiers func(ctx context.Context, parentSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers
}

type resAttrInfoTreeGenMetadata struct {
	createInput  resAttrInfoGenMetadata
	createOutput resAttrInfoGenMetadata

	readOutput resAttrInfoGenMetadata

	updateInput resAttrInfoGenMetadata

	deleteInput resAttrInfoGenMetadata

	propertySetterInputs  map[mgcName][]resAttrInfoGenMetadata
	propertySetterOutputs map[mgcName][]resAttrInfoGenMetadata
}

type attributeModifiers struct {
	isRequired                 bool
	isOptional                 bool
	isComputed                 bool
	useStateForUnknown         bool
	requiresReplaceWhenChanged bool
	nameOverride               tfName
	ignoreDefault              bool
	getChildModifiers          func(ctx context.Context, mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers
}

func addMgcSchemaAttributes(
	dst resAttrInfoMap,
	mgcSchema *mgcSdk.Schema,
	getModifiers func(ctx context.Context, mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers,
	ctx context.Context,
) error {
	mgcSchema, err := getXOfObjectSchemaTransformed(mgcSchema)
	if err != nil {
		return err
	}

	for propName, propSchemaRef := range mgcSchema.Properties {
		tflog.SubsystemDebug(ctx, schemaGenSubsystem, fmt.Sprintf("adding attribute %q", propName))
		propMgcName := mgcName(propName)
		propSchema := (*mgcSchemaPkg.Schema)(propSchemaRef.Value)

		modifiers := getModifiers(ctx, mgcSchema, propMgcName)
		if hasSchemaBeenPromoted(propSchema) {
			if modifiers.isRequired {
				modifiers.isRequired = false
			}
			if !modifiers.isComputed {
				modifiers.isOptional = true
			}
			tflog.SubsystemDebug(ctx, schemaGenSubsystem, fmt.Sprintf("computing %q as oneOf, %+v", propName, propSchema))
		}

		tfSchema, childAttributes, err := mgcSchemaToTFAttribute(propSchema, getModifiers(ctx, mgcSchema, propMgcName), ctx)
		tflog.SubsystemDebug(ctx, schemaGenSubsystem, fmt.Sprintf("attribute %q generated tfSchema %#v", propName, tfSchema))
		if err != nil {
			tflog.SubsystemError(ctx, schemaGenSubsystem, fmt.Sprintf("attribute %q schema: %+v; error: %s", propName, propSchema, err))
			return fmt.Errorf("attribute %q, error=%s", propName, err)
		}

		name := propMgcName.asTFName()
		if modifiers.nameOverride != "" {
			name = modifiers.nameOverride
		}

		attr := &resAttrInfo{
			tfName:          name,
			mgcName:         propMgcName,
			mgcSchema:       propSchema,
			tfSchema:        tfSchema,
			childAttributes: childAttributes,
		}
		dst[propMgcName] = attr
		tflog.SubsystemDebug(ctx, schemaGenSubsystem, fmt.Sprintf("attribute %q: %+v", propName, attr))
	}

	return nil
}

func getInputChildModifiers(ctx context.Context, mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
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

func getResultModifiers(ctx context.Context, mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
	return attributeModifiers{
		isRequired:                 false,
		isOptional:                 false,
		isComputed:                 true,
		useStateForUnknown:         false,
		requiresReplaceWhenChanged: false,
		ignoreDefault:              true,
		getChildModifiers:          getResultModifiers,
	}
}

func generateResAttrInfoTree(ctx context.Context, resName tfName, treeMetadata resAttrInfoTreeGenMetadata) (resAttrInfoTree, error) {
	diagnostics := Diagnostics{}
	tree := resAttrInfoTree{}
	var d Diagnostics

	tree.createInput, d = generateResAttrInfoMap(ctx, resName, treeMetadata.createInput)
	if diagnostics.AppendCheckError(d...) {
		return resAttrInfoTree{}, fmt.Errorf("errors when generating resource create input attributes: %#v", d.Errors())
	}

	tree.updateInput, d = generateResAttrInfoMap(ctx, resName, treeMetadata.updateInput)
	if diagnostics.AppendCheckError(d...) {
		return resAttrInfoTree{}, fmt.Errorf("errors when generating resource update input attributes: %#v", d.Errors())
	}

	tree.deleteInput, d = generateResAttrInfoMap(ctx, resName, treeMetadata.deleteInput)
	if diagnostics.AppendCheckError(d...) {
		return resAttrInfoTree{}, fmt.Errorf("errors when generating resource delete input attributes: %#v", d.Errors())
	}

	tree.createOutput, d = generateResAttrInfoMap(ctx, resName, treeMetadata.createOutput)
	if diagnostics.AppendCheckError(d...) {
		return resAttrInfoTree{}, fmt.Errorf("errors when generating resource create output attributes: %#v", d.Errors())
	}

	tree.readOutput, d = generateResAttrInfoMap(ctx, resName, treeMetadata.readOutput)
	if diagnostics.AppendCheckError(d...) {
		return resAttrInfoTree{}, fmt.Errorf("errors when generating resource read output attributes: %#v", d.Errors())
	}

	tree.propertySetterInputs, d = generatePropertySetterResAttrInfo(ctx, resName, treeMetadata.propertySetterInputs)
	if diagnostics.AppendCheckError(d...) {
		return resAttrInfoTree{}, fmt.Errorf("errors when generating resource property setter input attributes: %#v", d.Errors())
	}

	tree.propertySetterOutputs, d = generatePropertySetterResAttrInfo(ctx, resName, treeMetadata.propertySetterOutputs)
	if diagnostics.AppendCheckError(d...) {
		return resAttrInfoTree{}, fmt.Errorf("errors when generating resource property setter output attributes: %#v", d.Errors())
	}

	inputs := []resAttrInfoMap{tree.createInput, tree.updateInput, tree.deleteInput}
	for _, m := range tree.propertySetterInputs {
		for _, v := range m {
			inputs = append(inputs, v)
		}
	}
	tree.input, d = generateAggregateResAttrInfoMap(ctx, resName, "input", inputs)
	if diagnostics.AppendCheckError(d...) {
		return resAttrInfoTree{}, fmt.Errorf("errors when generating resource input attributes: %#v", d.Errors())
	}

	outputs := []resAttrInfoMap{tree.createOutput, tree.readOutput}
	for _, m := range tree.propertySetterOutputs {
		for _, v := range m {
			inputs = append(inputs, v)
		}
	}
	tree.output, d = generateAggregateResAttrInfoMap(ctx, resName, "output", outputs)
	if diagnostics.AppendCheckError(d...) {
		return resAttrInfoTree{}, fmt.Errorf("errors when generating resource output attributes: %#v", d.Errors())
	}

	return tree, nil
}

func generateAggregateResAttrInfoMap(ctx context.Context, resName tfName, attrType string, sources []resAttrInfoMap) (resAttrInfoMap, Diagnostics) {
	ctx = tflog.SubsystemSetField(ctx, schemaGenSubsystem, resourceNameField, resName)
	tflog.SubsystemDebug(ctx, schemaGenSubsystem, fmt.Sprintf("generating aggregate %s attributes", attrType))
	diagnostics := Diagnostics{}

	aggregateAttrInfoMap := resAttrInfoMap{}
	for _, attrInfoMap := range sources {
		for attrMgcName, attrInfo := range attrInfoMap {
			if current, ok := aggregateAttrInfoMap[attrMgcName]; ok {
				if err := mgcSchemaPkg.CompareJsonSchemas(current.mgcSchema, attrInfo.mgcSchema); err != nil && !isSubSchema(current.mgcSchema, attrInfo.mgcSchema) {
					return aggregateAttrInfoMap, diagnostics.AppendErrorReturn(
						fmt.Sprintf("Cannot generate aggregate CRUD attributes for resource %q", resName),
						fmt.Sprintf(
							"The same attribute in different CRUD operations has a different schema and is NOT a sub-schema. Diff: %v",
							err,
						),
					)
				}
				continue
			}

			aggregateAttrInfoMap[attrMgcName] = attrInfo

		}
	}

	return aggregateAttrInfoMap, diagnostics
}

func generateResAttrInfoMap(ctx context.Context, resName tfName, metadata resAttrInfoGenMetadata) (resAttrInfoMap, Diagnostics) {
	if metadata.schema == nil || metadata.modifiers == nil {
		return nil, nil
	}

	if metadata.schema.IsEmpty() {
		return nil, nil
	}

	ctx = tflog.SubsystemSetField(ctx, schemaGenSubsystem, resourceNameField, resName)
	tflog.SubsystemDebug(ctx, schemaGenSubsystem, "generating attr info map")
	diagnostics := Diagnostics{}

	attrInfoMap := resAttrInfoMap{}
	err := addMgcSchemaAttributes(attrInfoMap, metadata.schema, metadata.modifiers, ctx)
	if err != nil {
		return nil, diagnostics.AppendErrorReturn(
			"could not create TF attributes",
			err.Error(),
		)
	}

	return attrInfoMap, diagnostics
}

func generatePropertySetterResAttrInfo(
	ctx context.Context,
	resName tfName,
	metadatas map[mgcName][]resAttrInfoGenMetadata,
) (map[mgcName]map[*core.Schema]resAttrInfoMap, Diagnostics) {
	diagnostics := Diagnostics{}
	result := map[mgcName]map[*core.Schema]resAttrInfoMap{}

	for propName, propMetadatas := range metadatas {
		for _, metadata := range propMetadatas {
			attr, d := generateResAttrInfoMap(ctx, resName, metadata)
			if diagnostics.AppendCheckError(d...) {
				return nil, diagnostics
			}

			if attr == nil {
				continue
			}

			if result[propName] == nil {
				result[propName] = map[*core.Schema]resAttrInfoMap{}
			}

			result[propName][metadata.schema] = attr
		}
	}

	return result, diagnostics
}

func generateTFSchema(ctx context.Context, name tfName, description string, attrInfoTree resAttrInfoTree) schema.Schema {
	tflog.Debug(ctx, "generating schema")

	ctx = tflog.NewSubsystem(ctx, schemaGenSubsystem)
	ctx = tflog.SubsystemSetField(ctx, schemaGenSubsystem, resourceNameField, name)

	tfAttributes := generateTFAttributes(ctx, attrInfoTree)
	tfSchema := schema.Schema{Attributes: map[string]schema.Attribute{}}
	tfSchema.MarkdownDescription = description
	for tfName, tfAttr := range tfAttributes {
		tfSchema.Attributes[string(tfName)] = tfAttr
	}

	tfAttributeNames := []tfName{}
	for attrName := range tfAttributes {
		tfAttributeNames = append(tfAttributeNames, attrName)
	}

	tflog.Debug(ctx, "generated tf schema", map[string]any{"attributes": tfAttributeNames})

	return tfSchema
}

func generateTFAttributes(ctx context.Context, attrInfoTree resAttrInfoTree) map[tfName]schema.Attribute {
	tflog.SubsystemInfo(ctx, schemaGenSubsystem, "generating TF schema from input and output attributes in attr tree")

	tfAttributes := map[tfName]schema.Attribute{}
	tflog.SubsystemInfo(ctx, schemaGenSubsystem, "generating attributes using input")
	for name, iattr := range attrInfoTree.input {
		// Split attributes that differ between input/output
		for _, oattr := range attrInfoTree.output {
			if iattr.tfName != oattr.tfName {
				continue
			}
			if err := mgcSchemaPkg.CompareJsonSchemas(oattr.mgcSchema, iattr.mgcSchema); err != nil {
				os, _ := oattr.mgcSchema.MarshalJSON()
				is, _ := iattr.mgcSchema.MarshalJSON()
				tflog.SubsystemDebug(ctx, schemaGenSubsystem, fmt.Sprintf("attribute %q differs between input and output. input: %s - output %s\nerror=%s", name, is, os, err.Error()))
				iattr.tfName = iattr.tfName.asDesired()
				oattr.tfName = oattr.tfName.asCurrent()

				iattr.currentCounterpart = oattr
			}
		}

		tfAttributes[iattr.tfName] = iattr.tfSchema
	}

	tflog.SubsystemInfo(ctx, schemaGenSubsystem, "generating attributes using output")
	for _, oattr := range attrInfoTree.output {
		// If they don't differ and it's already created skip
		if _, ok := tfAttributes[oattr.tfName]; ok {
			continue
		}

		tfAttributes[oattr.tfName] = oattr.tfSchema
	}

	return tfAttributes
}

func mgcSchemaToTFAttribute(mgcSchema *mgcSdk.Schema, m attributeModifiers, ctx context.Context) (schema.Attribute, resAttrInfoMap, error) {
	description := mgcSchema.Description

	switch mgcSchema.Type {
	case "string":
		return mgcStringSchemaToTFAttribute(ctx, description, mgcSchema, m)
	case "number":
		return mgcNumberSchemaToTFAttribute(ctx, description, mgcSchema, m)
	case "integer":
		return mgcIntSchemaToTFAttribute(ctx, description, mgcSchema, m)
	case "boolean":
		return mgcBoolSchemaToTFAttribute(ctx, description, mgcSchema, m)
	case "array":
		return mgcArraySchemaToTFAttribute(ctx, description, mgcSchema, m)
	case "object":
		if mgcSchema.AdditionalProperties.Has != nil && *mgcSchema.AdditionalProperties.Has {
			return nil, nil, fmt.Errorf(
				"Unable to create Terraform Schema from MGC Schema. Schema with Additional Properties must have type information, and not just boolean: %#v",
				mgcSchema,
			)
		}
		if mgcSchema.AdditionalProperties.Schema != nil {
			if len(mgcSchema.Properties) > 0 {
				return nil, nil, fmt.Errorf(
					"Unable to create Terraform Schema from MGC Schema. Schema cannot have both Additional Properties and standard Properties: %#v",
					mgcSchema,
				)
			}

			return mgcMapSchemaToTFAttribute(ctx, description, mgcSchema, m)
		}
		return mgcObjectSchemaToTFAttribute(ctx, description, mgcSchema, m)
	default:
		return nil, nil, fmt.Errorf("type %q not supported", mgcSchema.Type)
	}
}

func mgcStringSchemaToTFAttribute(ctx context.Context, description string, mgcSchema *mgcSdk.Schema, m attributeModifiers) (schema.StringAttribute, resAttrInfoMap, error) {
	tflog.SubsystemDebug(ctx, schemaGenSubsystem, "generating attribute as string", map[string]any{"mgcSchema": mgcSchema})
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

	isComputed := m.isComputed
	var d defaults.String
	if v, ok := mgcSchema.Default.(string); ok && !m.isRequired && !m.ignoreDefault {
		d = stringdefault.StaticString(v)
		isComputed = true
	}

	return schema.StringAttribute{
		Description:   description,
		Required:      m.isRequired,
		Optional:      m.isOptional,
		Computed:      isComputed,
		PlanModifiers: mod,
		Default:       d,
	}, nil, nil
}

func mgcNumberSchemaToTFAttribute(ctx context.Context, description string, mgcSchema *mgcSdk.Schema, m attributeModifiers) (schema.NumberAttribute, resAttrInfoMap, error) {
	tflog.SubsystemDebug(ctx, schemaGenSubsystem, "generating attribute as number", map[string]any{"mgcSchema": mgcSchema})
	mod := []planmodifier.Number{}
	if m.useStateForUnknown {
		mod = append(mod, numberplanmodifier.UseStateForUnknown())
	}
	if m.requiresReplaceWhenChanged {
		mod = append(mod, numberplanmodifier.RequiresReplace())
	}

	isComputed := m.isComputed
	var d defaults.Number
	if v, ok := mgcSchema.Default.(float64); ok && !m.isRequired && !m.ignoreDefault {
		d = numberdefault.StaticBigFloat(big.NewFloat(v))
		isComputed = true
	}

	return schema.NumberAttribute{
		Description:   description,
		Required:      m.isRequired,
		Optional:      m.isOptional,
		Computed:      isComputed,
		PlanModifiers: mod,
		Default:       d,
	}, nil, nil
}

func mgcIntSchemaToTFAttribute(ctx context.Context, description string, mgcSchema *mgcSdk.Schema, m attributeModifiers) (schema.Int64Attribute, resAttrInfoMap, error) {
	tflog.SubsystemDebug(ctx, schemaGenSubsystem, "generating attribute as int", map[string]any{"mgcSchema": mgcSchema})
	mod := []planmodifier.Int64{}
	if m.useStateForUnknown {
		mod = append(mod, int64planmodifier.UseStateForUnknown())
	}
	if m.requiresReplaceWhenChanged {
		mod = append(mod, int64planmodifier.RequiresReplace())
	}

	isComputed := m.isComputed
	var d defaults.Int64
	if v, ok := mgcSchema.Default.(int64); ok && !m.isRequired && !m.ignoreDefault {
		d = int64default.StaticInt64(v)
		isComputed = true
	}

	return schema.Int64Attribute{
		Description:   description,
		Required:      m.isRequired,
		Optional:      m.isOptional,
		Computed:      isComputed,
		PlanModifiers: mod,
		Default:       d,
	}, nil, nil
}

func mgcBoolSchemaToTFAttribute(ctx context.Context, description string, mgcSchema *mgcSdk.Schema, m attributeModifiers) (schema.BoolAttribute, resAttrInfoMap, error) {
	tflog.SubsystemDebug(ctx, schemaGenSubsystem, "generating attribute as bool", map[string]any{"mgcSchema": mgcSchema})
	mod := []planmodifier.Bool{}
	if m.useStateForUnknown {
		mod = append(mod, boolplanmodifier.UseStateForUnknown())
	}
	if m.requiresReplaceWhenChanged {
		mod = append(mod, boolplanmodifier.RequiresReplace())
	}

	isComputed := m.isComputed
	var d defaults.Bool
	if v, ok := mgcSchema.Default.(bool); ok && !m.isRequired && !m.ignoreDefault {
		d = booldefault.StaticBool(v)
		isComputed = true
	}

	return schema.BoolAttribute{
		Description:   description,
		Required:      m.isRequired,
		Optional:      m.isOptional,
		Computed:      isComputed,
		PlanModifiers: mod,
		Default:       d,
	}, nil, nil
}

func mgcArraySchemaToTFAttribute(ctx context.Context, description string, mgcSchema *mgcSdk.Schema, m attributeModifiers) (schema.Attribute, resAttrInfoMap, error) {
	tflog.SubsystemDebug(ctx, schemaGenSubsystem, "generating attribute as array", map[string]any{"mgcSchema": mgcSchema})
	mgcItemSchema := (*core.Schema)(mgcSchema.Items.Value)
	elemAttr, elemAttrs, err := mgcSchemaToTFAttribute(mgcItemSchema, m.getChildModifiers(ctx, mgcItemSchema, "0"), ctx)
	if err != nil {
		return nil, nil, err
	}

	childAttrs := resAttrInfoMap{}
	childAttrs["0"] = &resAttrInfo{
		tfName:          "0",
		mgcName:         "0",
		mgcSchema:       mgcItemSchema,
		tfSchema:        elemAttr,
		childAttributes: elemAttrs,
	}

	mod := []planmodifier.List{}
	if m.requiresReplaceWhenChanged {
		mod = append(mod, listplanmodifier.RequiresReplace())
	}
	if m.useStateForUnknown {
		mod = append(mod, listplanmodifier.UseStateForUnknown())
	}

	isComputed := m.isComputed
	var d defaults.List
	if v, ok := mgcSchema.Default.([]any); ok && !m.isRequired && !m.ignoreDefault {
		lst, err := tfAttrListValueFromMgcSchema(ctx, mgcSchema, childAttrs["0"], v)
		if err != nil {
			return nil, nil, err
		}

		if l, ok := lst.(types.List); ok {
			d = listdefault.StaticValue(l)
			isComputed = true
		}
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
			Computed:      isComputed,
			PlanModifiers: mod,
			Default:       d,
		}, childAttrs, nil
	} else {
		return schema.ListAttribute{
			ElementType:   elemAttr.GetType(),
			Description:   description,
			Required:      m.isRequired,
			Optional:      m.isOptional,
			Computed:      isComputed,
			PlanModifiers: mod,
			Default:       d,
		}, childAttrs, nil
	}
}

func mgcMapSchemaToTFAttribute(ctx context.Context, description string, mgcSchema *mgcSdk.Schema, m attributeModifiers) (schema.Attribute, resAttrInfoMap, error) {
	tflog.SubsystemDebug(ctx, schemaGenSubsystem, "generating attribute as map", map[string]any{"mgcSchema": mgcSchema})
	mapValueMgcSchema := (*mgcSchemaPkg.Schema)(mgcSchema.AdditionalProperties.Schema.Value)

	mapValueTFSchema, mapValueChildAttributes, err := mgcSchemaToTFAttribute(mapValueMgcSchema, m, ctx)
	if err != nil {
		return nil, nil, err
	}

	childAttrs := resAttrInfoMap{}
	childAttrs["0"] = &resAttrInfo{
		tfName:          "0",
		mgcName:         "0",
		mgcSchema:       mapValueMgcSchema,
		tfSchema:        mapValueTFSchema,
		childAttributes: mapValueChildAttributes,
	}

	mod := []planmodifier.Map{}
	if m.requiresReplaceWhenChanged {
		mod = append(mod, mapplanmodifier.RequiresReplace())
	}
	if m.useStateForUnknown {
		mod = append(mod, mapplanmodifier.UseStateForUnknown())
	}

	isComputed := m.isComputed

	var d defaults.Map
	if v, ok := mgcSchema.Default.(map[string]any); ok && !m.isRequired && !m.ignoreDefault {
		m, err := tfAttrMapValueFromMgcSchema(ctx, mgcSchema, childAttrs["0"], v)
		if err != nil {
			return nil, nil, err
		}

		if m, ok := m.(types.Map); ok {
			d = mapdefault.StaticValue(m)
			isComputed = true
		}
	}

	if objAttr, ok := mapValueTFSchema.(schema.SingleNestedAttribute); ok {
		// This type assertion will/should NEVER fail, according to TF code
		nestedObj, ok := objAttr.GetNestedObject().(schema.NestedAttributeObject)
		if !ok {
			return nil, nil, fmt.Errorf("failed TF GetNestedObject")
		}
		return schema.MapNestedAttribute{
			NestedObject:        nestedObj,
			Description:         description,
			MarkdownDescription: description,
			Required:            m.isRequired,
			Optional:            m.isOptional,
			Computed:            isComputed,
			PlanModifiers:       mod,
			Default:             d,
		}, childAttrs, nil
	} else {
		return schema.MapAttribute{
			ElementType:         mapValueTFSchema.GetType(),
			Description:         description,
			MarkdownDescription: description,
			Required:            m.isRequired,
			Optional:            m.isOptional,
			Computed:            isComputed,
			PlanModifiers:       mod,
			Default:             d,
		}, childAttrs, nil
	}
}

func mgcObjectSchemaToTFAttribute(ctx context.Context, description string, mgcSchema *mgcSdk.Schema, m attributeModifiers) (schema.Attribute, resAttrInfoMap, error) {
	tflog.SubsystemDebug(ctx, schemaGenSubsystem, "generating attribute as object", map[string]any{"mgcSchema": mgcSchema})
	childAttrs := resAttrInfoMap{}
	err := addMgcSchemaAttributes(childAttrs, mgcSchema, m.getChildModifiers, ctx)
	if err != nil {
		return nil, nil, err
	}
	tfAttributes := map[string]schema.Attribute{}
	for _, attr := range childAttrs {
		tfAttributes[string(attr.tfName)] = attr.tfSchema
	}

	mod := []planmodifier.Object{}
	if m.requiresReplaceWhenChanged {
		mod = append(mod, objectplanmodifier.RequiresReplace())
	}
	if m.useStateForUnknown {
		mod = append(mod, objectplanmodifier.UseStateForUnknown())
	}

	isComputed := m.isComputed

	var d defaults.Object
	if v, ok := mgcSchema.Default.(map[string]any); ok && !m.isRequired && !m.ignoreDefault {
		obj, err := tfAttrObjectValueFromMgcSchema(ctx, mgcSchema, childAttrs, v)
		if err != nil {
			return nil, nil, err
		}

		if o, ok := obj.(types.Object); ok {
			d = objectdefault.StaticValue(o)
			isComputed = true
		}
	}

	return schema.SingleNestedAttribute{
		Attributes:    tfAttributes,
		Description:   description,
		Required:      m.isRequired,
		Optional:      m.isOptional,
		Computed:      isComputed,
		PlanModifiers: mod,
		Default:       d,
	}, childAttrs, nil
}

func tfAttrListValueFromMgcSchema(ctx context.Context, s *mgcSdk.Schema, listAttr *resAttrInfo, v []any) (attr.Value, error) {
	attrSchema := (*core.Schema)(s.Items.Value)
	attrType := listAttr.tfSchema.GetType()
	attrValues := []attr.Value{}
	for i := range v {
		v, ok, err := tfAttrValueFromMgcSchema(ctx, attrSchema, listAttr, v[i])
		if err != nil {
			return nil, err
		}

		if !ok {
			continue
		}

		attrValues = append(attrValues, v)
	}

	lst, diag := types.ListValue(attrType, attrValues)
	if diag.HasError() {
		return nil, fmt.Errorf("unable to create default list value")
	}
	return lst, nil
}

func tfAttrMapValueFromMgcSchema(ctx context.Context, s *mgcSdk.Schema, mapElemAttr *resAttrInfo, v map[string]any) (attr.Value, error) {
	mapType := mapElemAttr.tfSchema.GetType()
	mapElements := make(map[string]attr.Value, len(v))
	for k, sv := range v {
		propDefault, ok, err := tfAttrValueFromMgcSchema(ctx, mapElemAttr.mgcSchema, mapElemAttr, sv)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		mapElements[k] = propDefault
	}
	mapValue, d := types.MapValue(mapType, mapElements)
	if d.HasError() {
		return nil, fmt.Errorf("unable to create default map value: %v", d.Errors())
	}
	return mapValue, nil
}

func tfAttrObjectValueFromMgcSchema(ctx context.Context, s *mgcSdk.Schema, mapAttr map[mgcName]*resAttrInfo, v map[string]any) (attr.Value, error) {
	attrTypes := map[string]attr.Type{}
	attrValues := map[string]attr.Value{}
	for k := range v {
		attrSchema := (*core.Schema)(s.Properties[k].Value)

		val, ok, err := tfAttrValueFromMgcSchema(ctx, attrSchema, mapAttr[mgcName(k)], v[k])
		if err != nil {
			return nil, err
		}

		if !ok {
			continue
		}

		attrValues[k] = val
		attrTypes[k] = val.Type(ctx)
	}
	obj, diag := types.ObjectValue(attrTypes, attrValues)
	if diag.HasError() {
		return nil, fmt.Errorf("unable to create default object value")
	}
	return obj, nil
}

func tfAttrValueFromMgcSchema(ctx context.Context, s *mgcSdk.Schema, attrInfo *resAttrInfo, v any) (attr.Value, bool, error) {
	if v == nil {
		return nil, false, nil
	}

	switch s.Type {
	case "string":
		if dStr, ok := v.(string); ok {
			return types.StringValue(dStr), true, nil
		}
		return nil, false, fmt.Errorf("unable to create attr.Value of type string")
	case "number":
		if dFloat, ok := v.(float64); ok {
			return types.NumberValue(big.NewFloat(dFloat)), true, nil
		}
		return nil, false, fmt.Errorf("unable to create attr.Value of type number")
	case "integer":
		if dInt, ok := v.(int64); ok {
			return types.Int64Value(dInt), true, nil
		}
		return nil, false, fmt.Errorf("unable to create attr.Value of type integer")
	case "boolean":
		if b, ok := v.(bool); ok {
			return types.BoolValue(b), true, nil
		}
		return nil, false, fmt.Errorf("unable to create attr.Value of type boolean")
	case "array":
		listVal, ok := v.([]any)
		if !ok {
			return nil, false, fmt.Errorf("unable to create attr.Value of type list")
		}

		attrValue, err := tfAttrListValueFromMgcSchema(ctx, s, attrInfo, listVal)
		if err != nil {
			return nil, false, err
		}
		return attrValue, true, nil
	case "object":
		mapVal, ok := v.(map[string]any)
		if !ok {
			return nil, false, fmt.Errorf("unable to create attr.Value of type object")
		}

		if s.AdditionalProperties.Has != nil && *s.AdditionalProperties.Has {
			return nil, false, fmt.Errorf("unable to create attr.Value of type map when additional properties has no type information")
		}
		if s.AdditionalProperties.Schema != nil {
			if len(s.Properties) > 0 {
				return nil, false, fmt.Errorf("unable to create attr.Value of type map when MgcSchema has both additional and standard properties")
			}

			attrValue, err := tfAttrMapValueFromMgcSchema(ctx, s, attrInfo, mapVal)
			if err != nil {
				return nil, false, err
			}
			return attrValue, true, nil
		}

		attrValue, err := tfAttrObjectValueFromMgcSchema(ctx, s, attrInfo.childAttributes, mapVal)
		if err != nil {
			return nil, false, err
		}
		return attrValue, true, nil
	default:
		return nil, false, fmt.Errorf("type %q not supported", s.Type)
	}
}

func (n mgcName) asTFName() tfName {
	return tfName(strcase.SnakeCase(string(n)))
}

func (n mgcName) singular() mgcName {
	if len(n) == 0 {
		return n
	}

	if strings.HasSuffix(string(n), "ies") {
		return mgcName(strings.TrimSuffix(string(n), "ies") + "y")
	}

	return mgcName(strings.TrimSuffix(string(n), "s"))
}

func (n tfName) asDesired() tfName {
	return n
}

func (n tfName) asCurrent() tfName {
	if strings.HasPrefix(string(n), "current_") {
		return n
	}
	return "current_" + n
}
