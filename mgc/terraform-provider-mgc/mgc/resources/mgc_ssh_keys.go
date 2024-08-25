package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkSSHKeys "magalu.cloud/lib/products/profile/ssh_keys"
	"magalu.cloud/sdk"
	tfutil "magalu.cloud/terraform-provider-mgc/mgc/tfutil"

	mgcSdk "magalu.cloud/lib"
)

var (
	_ resource.Resource              = &sshKeys{}
	_ resource.ResourceWithConfigure = &sshKeys{}
)

func NewSshKeysResource() resource.Resource {
	return &sshKeys{}
}

type sshKeys struct {
	sdkClient *mgcSdk.Client
	sshKeys   sdkSSHKeys.Service
}

func (r *sshKeys) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_keys"
}

func (r *sshKeys) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	sdk, ok := req.ProviderData.(*sdk.Sdk)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected provider config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.sdkClient = mgcSdk.NewClient(sdk)
	r.sshKeys = sdkSSHKeys.NewService(ctx, r.sdkClient)
}

type sshKeyModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Key  types.String `tfsdk:"key"`
}

func (r *sshKeys) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the SSH key",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"key": schema.StringAttribute{
				Description: "Public key",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Description: "ID of the SSH key",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	resp.Schema.Description = "Add a new SSH key to your account"
}

func setValuesFromCreate(result sdkSSHKeys.CreateResult) *sshKeyModel {
	return &sshKeyModel{
		ID:   types.StringValue(result.Id),
		Name: types.StringValue(result.Name),
		Key:  types.StringValue(result.Key),
	}
}

func setValuesFromGet(result sdkSSHKeys.GetResult) *sshKeyModel {
	return &sshKeyModel{
		ID:   types.StringValue(result.Id),
		Name: types.StringValue(result.Name),
		Key:  types.StringValue(result.Key),
	}
}

func (r *sshKeys) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	plan := &sshKeyModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	getResult, err := r.sshKeys.Get(sdkSSHKeys.GetParameters{
		KeyId: plan.ID.ValueString(),
	},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkSSHKeys.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Error Reading ssh key", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, setValuesFromGet(getResult))...)
}

func (r *sshKeys) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &sshKeyModel{}
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createResult, err := r.sshKeys.Create(sdkSSHKeys.CreateParameters{
		Key:  plan.Key.ValueString(),
		Name: plan.Name.ValueString(),
	},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkSSHKeys.CreateConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Error creating ssh key", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, setValuesFromCreate(createResult))...)
}

func (r *sshKeys) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *sshKeys) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data sshKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	_, err := r.sshKeys.Delete(
		sdkSSHKeys.DeleteParameters{
			KeyId: data.ID.ValueString(),
		},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkSSHKeys.DeleteConfigs{}),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting ssh key", err.Error())
	}
}
