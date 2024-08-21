package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sdkSSHKeys "magalu.cloud/lib/products/ssh/public_keys"
	"magalu.cloud/sdk"
	tfutil "magalu.cloud/terraform-provider-mgc/internal/tfutil"

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
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Key_Type  types.String `tfsdk:"key_type"`
	Key       types.String `tfsdk:"key"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func (r *sshKeys) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ssh_keys": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of ssh-keys.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.BoolAttribute{
							Computed:    true,
							Description: "ID of ssh key.",
						},
						"key": schema.StringAttribute{
							Required:    true,
							Description: "The value of public ssh key.",
						},
						"key_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of ssh key.",
						},
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The name of ssh key",
						},
						"created_at": schema.StringAttribute{
							Description: "The timestamp when the ssh key was created.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
							Computed: true,
						},
					},
				},
			},
		},
	}
	resp.Schema.Description = "Get the available virtual-machine images."
}

func (r *sshKeys) setValuesFromServer(result sdkSSHKeys.GetResult, state *sshKeyModel) {
	state.ID = types.StringValue(result.Id)
	state.Key_Type = types.StringValue(result.KeyType)
}

func (r *sshKeys) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	plan := &sshKeyModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	getResult, err := r.sshKeys.Get(sdkSSHKeys.GetParameters{
		KeyId: plan.ID.ValueString(),
	},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkSSHKeys.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading block storage",
			"Could not read BS ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	r.setValuesFromServer(getResult, plan)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sshKeys) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &sshKeyModel{}
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := plan

	createResult, err := r.sshKeys.Create(sdkSSHKeys.CreateParameters{
		Key:  state.Key.ValueString(),
		Name: state.Name.ValueString(),
	},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkSSHKeys.CreateConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ssh key",
			"Could not create ssh-key, unexpected error: "+err.Error(),
		)
		return
	}

	getCreatedResource, err := r.sshKeys.Get(sdkSSHKeys.GetParameters{
		KeyId: createResult.Id,
	},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkSSHKeys.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading BS",
			"Could not read BS ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	r.setValuesFromServer(getCreatedResource, state)

	state.CreatedAt = types.StringValue(time.Now().Format(time.RFC850))
	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sshKeys) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// TODO // TODO // TODO
}

func (r *sshKeys) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data sshKeyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	err := r.sshKeys.Delete(
		sdkSSHKeys.DeleteParameters{
			KeyId: data.ID.ValueString(),
		},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkSSHKeys.DeleteConfigs{}),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting block storage volume",
			"Could not delete block storage volume "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}
}
