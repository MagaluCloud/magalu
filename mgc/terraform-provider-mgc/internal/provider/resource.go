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
	"golang.org/x/exp/slices"
	"magalu.cloud/core"
	mgcSdk "magalu.cloud/sdk"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &MgcResource{}
var _ resource.ResourceWithImportState = &MgcResource{}

type splitMgcAttribute struct {
	current *attribute
	desired *attribute
}

// MgcResource defines the resource implementation.
type MgcResource struct {
	sdk        *mgcSdk.Sdk
	name       string
	group      mgcSdk.Grouper // TODO: is this needed?
	create     mgcSdk.Executor
	read       mgcSdk.Executor
	update     mgcSdk.Executor
	delete     mgcSdk.Executor
	inputAttr  mgcAttributes
	outputAttr mgcAttributes
	splitAttr  []splitMgcAttribute
	tfschema   *schema.Schema
}

// BEGIN: tfSchemaHandler implementation

func (r *MgcResource) Name() string {
	return r.name
}

func (r *MgcResource) getCreateParamsModifiers(ctx context.Context, mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
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

func (r *MgcResource) getUpdateParamsModifiers(ctx context.Context, mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
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

func (r *MgcResource) getDeleteParamsModifiers(ctx context.Context, mgcSchema *mgcSdk.Schema, mgcName mgcName) attributeModifiers {
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

func (r *MgcResource) ReadInputAttributes(ctx context.Context) diag.Diagnostics {
	d := diag.Diagnostics{}
	if len(r.inputAttr) != 0 {
		return d
	}
	tflog.Debug(ctx, fmt.Sprintf("[resource] schema for %q: reading input attributes", r.name))

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

func (r *MgcResource) ReadOutputAttributes(ctx context.Context) diag.Diagnostics {
	d := diag.Diagnostics{}
	if len(r.outputAttr) != 0 {
		return d
	}
	tflog.Debug(ctx, fmt.Sprintf("[resource] schema for %q: reading output attributes", r.name))

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

func (r *MgcResource) InputAttributes() mgcAttributes {
	return r.inputAttr
}

func (r *MgcResource) OutputAttributes() mgcAttributes {
	return r.outputAttr
}

func (r *MgcResource) AppendSplitAttribute(split splitMgcAttribute) {
	if r.splitAttr == nil {
		r.splitAttr = []splitMgcAttribute{}
	}
	r.splitAttr = append(r.splitAttr, split)
}

var _ tfSchemaHandler = (*MgcResource)(nil)

// END: tfSchemaHandler implementation

// BEGIN: tfStateHandler implementation

func (r *MgcResource) SplitAttributes() []splitMgcAttribute {
	return r.splitAttr
}

func (r *MgcResource) TFSchema() *schema.Schema {
	return r.tfschema
}

func (r *MgcResource) ReadResultSchema() *mgcSdk.Schema {
	return r.read.ResultSchema()
}

var _ tfStateHandler = (*MgcResource)(nil)

// END: tfStateHandler implementation

// BEGIN: Resource implementation

func (r *MgcResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.name
}

func (r *MgcResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// TODO: Handle nullable values
	tflog.Debug(ctx, fmt.Sprintf("[resource] generating schema for %q", r.name))

	if r.tfschema == nil {
		tfs, d := generateTFSchema(r, ctx)
		resp.Diagnostics.Append(d...)
		if d.HasError() {
			return
		}

		tfs.MarkdownDescription = r.name
		r.tfschema = &tfs
	}

	resp.Schema = *r.tfschema
}

func (r *MgcResource) performOperation(ctx context.Context, exec core.Executor, inState tfsdk.State, outState *tfsdk.State, diag *diag.Diagnostics) {
	ctx = r.sdk.WrapContext(ctx)

	configs := getConfigs(exec.ConfigsSchema())
	params := readMgcMap(r, exec.ParametersSchema(), ctx, inState, diag)
	if diag.HasError() {
		return
	}

	result := execute(r.name, ctx, exec, params, configs, diag)
	if diag.HasError() {
		return
	}

	applyStateAfter(r, result, ctx, outState, diag)
}

func (r *MgcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.performOperation(ctx, r.create, tfsdk.State(req.Plan), &resp.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[resource] created a %q resource", r.name))
}

func (r *MgcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	r.performOperation(ctx, r.read, req.State, &resp.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[resource] read a %q resource", r.name))
}

func (r *MgcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.performOperation(ctx, r.update, tfsdk.State(req.Plan), &resp.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[resource] updated a %q resource", r.name))
}

func (r *MgcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	r.performOperation(ctx, r.delete, req.State, &resp.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[resource] deleted a %q resource", r.name))
}

func (r *MgcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// END: Resource implemenation
