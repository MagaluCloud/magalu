package provider

import (
	"context"
	"fmt"
	"strings"

	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
	mgcSdk "magalu.cloud/sdk"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MgcResource{}
var _ resource.ResourceWithImportState = &MgcResource{}

// MgcResource defines the resource implementation.
type MgcResource struct {
	sdk             *mgcSdk.Sdk
	resTfName       tfName
	resMgcName      mgcName
	description     string
	create          mgcSdk.Executor
	read            mgcSdk.Executor
	update          mgcSdk.Executor
	delete          mgcSdk.Executor
	attrTree        resAttrInfoTree
	tfschema        *schema.Schema
	propertySetters map[mgcName]propertySetter
}

func newMgcResource(
	ctx context.Context,
	sdk *mgcSdk.Sdk,
	resTfName tfName,
	resMgcName mgcName,
	description string,
	create, read, update, delete, list mgcSdk.Executor,
) (*MgcResource, error) {
	if create == nil {
		return nil, fmt.Errorf("resource %q misses create", resTfName)
	}
	if delete == nil {
		return nil, fmt.Errorf("resource %q misses delete", resTfName)
	}
	if read == nil {
		if list == nil {
			return nil, fmt.Errorf("resource %q misses read", resTfName)
		}

		readFromList, err := createReadFromList(list, create.ResultSchema(), delete.ParametersSchema())
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("unable to generate 'read' operation from 'list' for %q", resTfName), map[string]any{"error": err})
			return nil, fmt.Errorf("resource %q misses read", resTfName)
		}

		read = readFromList
		tflog.Debug(ctx, fmt.Sprintf("generated 'read' operation based on 'list' for %q", resTfName))
	}
	if update == nil {
		update = core.NoOpExecutor()
	}

	propertySetterContainer, err := collectPropertySetterContainers(read.Links())
	if err != nil {
		return nil, err
	}

	var propertySetters map[mgcName]propertySetter
	for key, container := range propertySetterContainer {
		if propertySetters == nil {
			propertySetters = make(map[mgcName]propertySetter, len(propertySetterContainer))
		}

		switch container.argCount {
		case 0:
			propertySetters[key] = newDefaultPropertySetter(key, container.entries[0].target)
		case 1:
			propertySetters[key] = newEnumPropertySetter(key, container.entries)
		case 2:
			propertySetters[key] = newStrTransitionPropertySetter(key, container.entries)
		default:
			// TODO: Handle more action types?
			continue
		}
	}

	r := &MgcResource{
		sdk:             sdk,
		resTfName:       resTfName,
		resMgcName:      resMgcName,
		description:     description,
		create:          create,
		read:            read,
		update:          update,
		delete:          delete,
		propertySetters: propertySetters,
	}

	propertySetterInputs := make(map[mgcName][]resAttrInfoGenMetadata, len(propertySetters))
	propertySetterOutputs := make(map[mgcName][]resAttrInfoGenMetadata, len(propertySetters))
	for propName, propSetter := range propertySetters {
		for _, paramsSchemas := range propSetter.allParametersSchemas() {
			propertySetterInputs[propName] = append(propertySetterInputs[propName], resAttrInfoGenMetadata{paramsSchemas, r.getUpdateParamsModifiers})
		}
		for _, resultSchemas := range propSetter.allResultSchemas() {
			propertySetterOutputs[propName] = append(propertySetterOutputs[propName], resAttrInfoGenMetadata{resultSchemas, r.getResultModifiers})
		}
	}

	attrTree, err := generateResAttrInfoTree(ctx, r.resTfName, resAttrInfoTreeGenMetadata{
		createInput: resAttrInfoGenMetadata{r.create.ParametersSchema(), r.getCreateParamsModifiers},
		updateInput: resAttrInfoGenMetadata{r.update.ParametersSchema(), r.getUpdateParamsModifiers},
		deleteInput: resAttrInfoGenMetadata{r.delete.ParametersSchema(), r.getDeleteParamsModifiers},

		createOutput: resAttrInfoGenMetadata{r.create.ResultSchema(), r.getResultModifiers},
		readOutput:   resAttrInfoGenMetadata{r.read.ResultSchema(), r.getResultModifiers},

		propertySetterInputs:  propertySetterInputs,
		propertySetterOutputs: propertySetterOutputs,
	})
	if err != nil {
		return nil, err
	}

	r.attrTree = attrTree
	return r, nil
}

func (r *MgcResource) doesPropHaveSetter(name mgcName) bool {
	_, ok := r.propertySetters[name]
	return ok
}

func (r *MgcResource) getCreateParamsModifiers(ctx context.Context, mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
	var k = string(mgcName)

	propSchema := (*core.Schema)(mgcSchema.Properties[k].Value)
	isRequired := slices.Contains(mgcSchema.Required, k)
	isComputed := !isRequired

	if isComputed {
		readSchema := r.read.ResultSchema().Properties[k]
		if readSchema == nil {
			isComputed = false
		} else {
			// If not required and present in read it can be compute
			isComputed = mgcSchemaPkg.CheckSimilarJsonSchemas((*core.Schema)(readSchema.Value), propSchema)
		}
	}

	requiresReplace := r.update.ParametersSchema().Properties[k] == nil && !r.doesPropHaveSetter(mgcName) && propSchema.Default == nil

	return attributeModifiers{
		isRequired:                 isRequired,
		isOptional:                 !isRequired,
		isComputed:                 isComputed,
		useStateForUnknown:         true,
		requiresReplaceWhenChanged: requiresReplace,
		getChildModifiers:          getInputChildModifiers,
	}
}

func (r *MgcResource) getUpdateParamsModifiers(ctx context.Context, mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
	var k = string(mgcName)

	isComputed := r.create.ResultSchema().Properties[k] != nil
	required := slices.Contains(mgcSchema.Required, k)

	var nameOverride tfName
	if withoutResPrefix, hasResPrefix := strings.CutPrefix(k, string(r.resMgcName.singular()+"_")); hasResPrefix && r.create.ResultSchema().Properties[withoutResPrefix] != nil {
		nameOverride = tfName(withoutResPrefix)
		isComputed = true
	} else if withoutNewPrefix, hasNewPrefix := strings.CutPrefix(k, "new_"); hasNewPrefix && r.read.ResultSchema().Properties[withoutNewPrefix] != nil {
		nameOverride = tfName(withoutNewPrefix)
	}

	return attributeModifiers{
		isRequired:                 required && !isComputed,
		isOptional:                 !required && !isComputed,
		isComputed:                 !required || isComputed,
		useStateForUnknown:         true,
		nameOverride:               nameOverride,
		requiresReplaceWhenChanged: false,
		getChildModifiers:          getInputChildModifiers,
	}
}

func (r *MgcResource) getDeleteParamsModifiers(ctx context.Context, mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
	k := string(mgcName)
	if _, isInRead := r.read.ParametersSchema().Properties[k]; isInRead {
		return r.getUpdateParamsModifiers(ctx, mgcSchema, mgcName)
	}

	if _, isInUpdate := r.update.ParametersSchema().Properties[k]; isInUpdate {
		return r.getUpdateParamsModifiers(ctx, mgcSchema, mgcName)
	}

	isComputed := r.create.ResultSchema().Properties[k] != nil

	var nameOverride tfName
	if withoutResPrefix, hasResPrefix := strings.CutPrefix(k, string(r.resMgcName.singular()+"_")); hasResPrefix && r.create.ResultSchema().Properties[withoutResPrefix] != nil {
		nameOverride = tfName(withoutResPrefix)
		isComputed = true
	}

	// All Delete parameters need to be optional, since they're not returned by the server when creating the resources,
	// and thus would produce "inconsistent results" in Terraform. Since we calculate the 'Update' parameters first,
	// all strictly necessary parameters for deletion (like the resource ID) will already be computed correctly (in the Update
	// parameters)
	return attributeModifiers{
		isRequired:                 false,
		isOptional:                 true,
		isComputed:                 isComputed,
		useStateForUnknown:         true,
		nameOverride:               nameOverride,
		requiresReplaceWhenChanged: false,
		getChildModifiers:          getInputChildModifiers,
	}
}

var propsUpdatableByServer = []string{"error", "updated_at", "updated"}
var propsNotUpdatableByServer = []string{"created_at"}

func (r *MgcResource) getResultModifiers(ctx context.Context, mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
	propSchema := (*mgcSchemaPkg.Schema)(mgcSchema.Properties[string(mgcName)].Value)

	isOptional := r.doesPropHaveSetter(mgcName)
	useStateForUnknown := (isOptional || slices.Contains(propsNotUpdatableByServer, string(mgcName))) && !slices.Contains(propsUpdatableByServer, string(mgcName))
	// In the following cases, this Result Attr will have a different structure when compared to its
	// Input counterpart, and this it will be transformed into 'current_<name>'. These 'current_<name>'
	// attributes, with mismatches, should NEVER be optional, the User must never deal with them
	// directly. The Input counterpart will already have been computed, so this can safely be
	// 'optional == false' since the '<name>' version of this attribute will be 'optional == true'
	// (or 'required == true').
	//
	// Ideally we would make this check when generating the attributes and seeing if there are
	// mismatches, when we transform this into a 'current_<name>' attribute, but we can't modify
	// the 'isOptional' value to false in that step since it's after everything has already been
	// created and we only have an interface to deal with the TF Schema :/
	var inputCounterpart *mgcSchemaPkg.Schema
	if createProp, ok := r.create.ParametersSchema().Properties[string(mgcName)]; ok {
		inputCounterpart = (*mgcSchemaPkg.Schema)(createProp.Value)
	} else if updateProp, ok := r.update.ParametersSchema().Properties[string(mgcName)]; ok {
		inputCounterpart = (*mgcSchemaPkg.Schema)(updateProp.Value)
	}

	if inputCounterpart != nil {
		willSplit := !mgcSchemaPkg.CheckSimilarJsonSchemas(inputCounterpart, propSchema)
		isOptional = isOptional && !willSplit
		useStateForUnknown = !willSplit
	}

	return attributeModifiers{
		isRequired:                 false,
		isOptional:                 isOptional,
		isComputed:                 true,
		useStateForUnknown:         useStateForUnknown,
		requiresReplaceWhenChanged: false,
		ignoreDefault:              true,
		getChildModifiers:          getResultModifiers,
	}
}

// BEGIN: Resource implemenation

func (r *MgcResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = string(r.resTfName)
}

func (r *MgcResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	if r.tfschema == nil {
		ctx = tflog.SetField(ctx, resourceNameField, r.resTfName)
		tfs := generateTFSchema(ctx, r.resTfName, r.description, r.attrTree)
		r.tfschema = &tfs
	}
	resp.Schema = *r.tfschema
}

func (r *MgcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	diagnostics := Diagnostics{}
	defer func() {
		resp.Diagnostics = diag.Diagnostics(diagnostics)
	}()

	createOp := newMgcResourceCreate(r.resTfName, r.attrTree, r.create, r.read, r.propertySetters)
	operationRunner := newMgcOperationRunner(r.sdk, createOp, tfsdk.State(req.Plan), req.Plan, &resp.State)
	diagnostics = operationRunner.Run(ctx)
}

func (r *MgcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	diagnostics := Diagnostics{}
	defer func() {
		resp.Diagnostics = diag.Diagnostics(diagnostics)
	}()

	operation := newMgcResourceRead(r.resTfName, r.attrTree, r.read)
	runner := newMgcOperationRunner(r.sdk, operation, req.State, tfsdk.Plan(req.State), &resp.State)
	diagnostics = runner.Run(ctx)
}

func (r *MgcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	diagnostics := Diagnostics{}
	defer func() {
		resp.Diagnostics = diag.Diagnostics(diagnostics)
	}()

	operation := newMgcResourceUpdate(r.resTfName, r.attrTree, r.update, r.read, r.propertySetters)
	runner := newMgcOperationRunner(r.sdk, operation, req.State, req.Plan, &resp.State)
	diagnostics = runner.Run(ctx)
}

func (r *MgcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	diagnostics := Diagnostics{}
	defer func() {
		resp.Diagnostics = diag.Diagnostics(diagnostics)
	}()

	deleteOp := newMgcResourceDelete(r.resTfName, r.attrTree, r.delete)
	runner := newMgcOperationRunner(r.sdk, deleteOp, req.State, tfsdk.Plan(req.State), &resp.State)
	diagnostics = runner.Run(ctx)
}

func (r *MgcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ctx = tflog.SetField(ctx, rpcField, "import-state")
	ctx = tflog.SetField(ctx, resourceNameField, r.resTfName)
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// END: Resource implemenation
