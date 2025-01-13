package resources

import (
	"context"
	"strings"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	networkPIP "github.com/MagaluCloud/magalu/mgc/lib/products/network/public_ips"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/tfutil"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkPublicIAttachPModel struct {
	PublicIpID  types.String `tfsdk:"public_ip_id"`
	InterfaceID types.String `tfsdk:"interface_id"`
}

type NetworkPublicIPAttachResource struct {
	sdkClient  *mgcSdk.Client
	networkPIP networkPIP.Service
}

func NewNetworkPublicIPAttachResource() resource.Resource {
	return &NetworkPublicIPAttachResource{}
}

func (r *NetworkPublicIPAttachResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_public_ips_attach"
}

func (r *NetworkPublicIPAttachResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	var err error
	var errDetail error
	r.sdkClient, err, errDetail = client.NewSDKClient(req, resp)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), errDetail.Error())
		return
	}

	r.networkPIP = networkPIP.NewService(ctx, r.sdkClient)
}

func (r *NetworkPublicIPAttachResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network Public IP Attach",
		Attributes: map[string]schema.Attribute{
			"interface_id": schema.StringAttribute{
				Description: "Interface ID",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_ip_id": schema.StringAttribute{
				Description: "Public IP ID",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *NetworkPublicIPAttachResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model NetworkPublicIAttachPModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.networkPIP.AttachContext(ctx, networkPIP.AttachParameters{
		PublicIpId: model.PublicIpID.ValueString(),
		PortId:     model.InterfaceID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkPIP.AttachConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Error attaching public IP", err.Error())
		return
	}

	resp.State.Set(ctx, &model)
}

func (r *NetworkPublicIPAttachResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model NetworkPublicIAttachPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.networkPIP.DetachContext(ctx, networkPIP.DetachParameters{
		PublicIpId: model.PublicIpID.ValueString(),
		PortId:     model.InterfaceID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkPIP.DetachConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Error detaching public IP", err.Error())
		return
	}
}

func (r *NetworkPublicIPAttachResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Update not supported")
}

func (r *NetworkPublicIPAttachResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model NetworkPublicIAttachPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pip, err := r.networkPIP.GetContext(ctx, networkPIP.GetParameters{
		PublicIpId: model.PublicIpID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkPIP.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Error reading public IP", err.Error())
		return
	}

	model.InterfaceID = types.StringPointerValue(pip.PortId)
	resp.State.Set(ctx, &model)
}

func (r *NetworkPublicIPAttachResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var model NetworkPublicIAttachPModel
	input := req.ID
	if input == "" {
		resp.Diagnostics.AddError("Invalid ID", "ID must be in the format public_ip_id,interface_id")
		return
	}
	parts := strings.Split(input, ",")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid ID", "ID must be in the format public_ip_id,interface_id")
		return
	}

	model.PublicIpID = types.StringValue(parts[0])
	model.InterfaceID = types.StringValue(parts[1])
	resp.State.Set(ctx, &model)
}
