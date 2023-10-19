package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
	mgcSdk "magalu.cloud/sdk"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MgcActionResource{}
var _ resource.ResourceWithImportState = &MgcActionResource{}

// MgcResource defines the resource implementation.
type MgcActionResource struct {
	sdk        *mgcSdk.Sdk
	name       string
	readOwner  mgcSdk.Executor
	create     mgcSdk.Linker
	read       mgcSdk.Linker
	update     mgcSdk.Linker // TODO: Will action resources need/have updates?
	delete     mgcSdk.Linker
	inputAttr  mgcAttributes
	outputAttr mgcAttributes
	splitAttr  []splitMgcAttribute
	tfschema   *schema.Schema
}

// BEGIN: tfSchemaHandler implementation

func (r *MgcActionResource) Name() string {
	return r.name
}

func (r *MgcActionResource) getReadParamsModifiers(ctx context.Context, mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
	return attributeModifiers{
		isRequired:                 true,
		isOptional:                 false,
		isComputed:                 false,
		useStateForUnknown:         true,
		requiresReplaceWhenChanged: true,
		getChildModifiers:          getInputChildModifiers,
	}
}

func (r *MgcActionResource) getDeleteParamsModifiers(ctx context.Context, mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
	// For now we consider all delete params as optionals, we need to think a way for the user to define
	// required delete params
	return attributeModifiers{
		isRequired:                 false,
		isOptional:                 true,
		isComputed:                 false,
		useStateForUnknown:         true,
		requiresReplaceWhenChanged: true,
		getChildModifiers:          getInputChildModifiers,
	}
}

func (r *MgcActionResource) ReadInputAttributes(ctx context.Context) diag.Diagnostics {
	ctx = tflog.SubsystemSetField(ctx, string(schemaGenSubsystem), string(actionResourceNameField), r.name)
	d := diag.Diagnostics{}
	if len(r.inputAttr) != 0 {
		return d
	}
	tflog.SubsystemDebug(ctx, string(schemaGenSubsystem), fmt.Sprintf("[action-resource] schema for %q: reading input attributes", r.name))

	r.inputAttr = mgcAttributes{}

	err := addMgcSchemaAttributes(
		r.inputAttr,
		r.create.AdditionalParametersSchema(),
		r.getReadParamsModifiers,
		ctx,
	)
	if err != nil {
		d.AddError("could not create TF input attributes", err.Error())
		return d
	}

	err = addMgcSchemaAttributes(
		r.inputAttr,
		r.readOwner.ParametersSchema(),
		r.getReadParamsModifiers,
		ctx,
	)
	if err != nil {
		d.AddError("could not create TF input attributes", err.Error())
		return d
	}

	err = addMgcSchemaAttributes(
		r.inputAttr,
		r.read.AdditionalParametersSchema(),
		r.getReadParamsModifiers,
		ctx,
	)
	if err != nil {
		d.AddError("could not create TF input attributes", err.Error())
		return d
	}

	err = addMgcSchemaAttributes(
		r.inputAttr,
		r.delete.AdditionalParametersSchema(),
		r.getDeleteParamsModifiers,
		ctx,
	)
	if err != nil {
		d.AddError("could not create TF input attributes", err.Error())
		return d
	}

	return d
}

func (r *MgcActionResource) ReadOutputAttributes(ctx context.Context) diag.Diagnostics {
	ctx = tflog.SubsystemSetField(ctx, string(schemaGenSubsystem), string(actionResourceNameField), r.name)
	d := diag.Diagnostics{}
	if len(r.outputAttr) != 0 {
		return d
	}
	tflog.SubsystemDebug(ctx, string(schemaGenSubsystem), "reading output attributes")

	r.outputAttr = mgcAttributes{}
	err := addMgcSchemaAttributes(
		r.outputAttr,
		r.create.ResultSchema(),
		getResultModifiers,
		ctx,
	)
	if err != nil {
		d.AddError("could not create TF output attributes", err.Error())
		return d
	}

	err = addMgcSchemaAttributes(
		r.outputAttr,
		r.read.ResultSchema(),
		getResultModifiers,
		ctx,
	)
	if err != nil {
		d.AddError("could not create TF output attributes", err.Error())
		return d
	}

	return d
}

func (r *MgcActionResource) InputAttributes() mgcAttributes {
	return r.inputAttr
}

func (r *MgcActionResource) OutputAttributes() mgcAttributes {
	return r.outputAttr
}

func (r *MgcActionResource) AppendSplitAttribute(split splitMgcAttribute) {
	if r.splitAttr == nil {
		r.splitAttr = []splitMgcAttribute{}
	}
	r.splitAttr = append(r.splitAttr, split)
}

var _ tfSchemaHandler = (*MgcActionResource)(nil)

// END: tfSchemaHandler implementation

// BEGIN: tfStateHandler implementation

func (r *MgcActionResource) TFSchema() *schema.Schema {
	return r.tfschema
}

func (r *MgcActionResource) SplitAttributes() []splitMgcAttribute {
	return r.splitAttr
}

func (r *MgcActionResource) ReadResultSchema() *mgcSdk.Schema {
	return r.read.ResultSchema()
}

var _ tfStateHandler = (*MgcActionResource)(nil)

// END: tfStateHandler implementation

// BEGIN: Resource implementation

func (r *MgcActionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.name
}

func (r *MgcActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// TODO: Handle nullable values
	tflog.Debug(ctx, fmt.Sprintf("[action-resource] generating schema for %q", r.name))

	if r.tfschema == nil {
		tfs, d := generateTFSchema(r, ctx)
		resp.Diagnostics.Append(d...)
		if d.HasError() {
			return
		}

		tfs.MarkdownDescription = r.name
		r.tfschema = &tfs
	}

	attributes := []string{}
	for attrName := range (*r.tfschema).Attributes {
		attributes = append(attributes, attrName)
	}

	tflog.Debug(ctx, fmt.Sprintf("[action-resource] generated tf schema for %q", r.name), map[string]any{"attributes": attributes})
	resp.Schema = *r.tfschema
}

func (r *MgcActionResource) performOwnerRead(ctx context.Context, tfState tfsdk.State, d *diag.Diagnostics) core.ResultWithValue {
	configs := getConfigs(r.readOwner.ConfigsSchema())
	params := readMgcMap(r, r.readOwner.ParametersSchema(), ctx, tfState, d)
	if d.HasError() {
		return nil
	}

	return execute(r.name, ctx, r.readOwner, params, configs, d)
}

func (r *MgcActionResource) performLinkOperation(ctx context.Context, link core.Linker, inState tfsdk.State, outState *tfsdk.State, diag *diag.Diagnostics) {
	ctx = r.sdk.WrapContext(ctx)

	ownerResult := r.performOwnerRead(ctx, inState, diag)
	if diag.HasError() {
		return
	}

	configs := getConfigs(link.AdditionalConfigsSchema())
	params := readMgcMap(r, link.AdditionalParametersSchema(), ctx, inState, diag)
	if diag.HasError() {
		return
	}

	linkExec, err := link.CreateExecutor(ownerResult)
	if err != nil {
		diag.AddError("error when creating link executor", err.Error())
		return
	}
	result := execute(r.name, ctx, linkExec, params, configs, diag)
	if diag.HasError() {
		return
	}
	applyStateAfter(r, result, ctx, outState, diag)
}

func (r *MgcActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.performLinkOperation(ctx, r.create, tfsdk.State(req.Plan), &resp.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[resource] created a %q resource", r.name))
}

func (r *MgcActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.performLinkOperation(ctx, r.read, req.State, &resp.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		// When reading fails, that means that the resource was most likely altered outside of terraform.
		resp.Diagnostics.AddError("reading the resource failed", "was the resource altered outside of terraform?")
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[resource] read a %q resource", r.name))
}

// Update will most likely never be called, as we always require replace when changed
func (r *MgcActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.update == nil {
		resp.Diagnostics.AddError(
			"no 'update' operation was provided",
			fmt.Sprintf("action resource %q doesn't have an update operation to run", r.name),
		)
		return
	}
	r.performLinkOperation(ctx, r.update, tfsdk.State(req.Plan), &resp.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[resource] updated a %q resource", r.name))
}

func (r *MgcActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.performLinkOperation(ctx, r.delete, req.State, &resp.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[resource] deleted a %q resource", r.name))
}

func (r *MgcActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

var _ resource.Resource = (*MgcActionResource)(nil)

// END: Resource implementation
