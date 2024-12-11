package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	networkPIP "magalu.cloud/lib/products/network/public_ips"
	networkVPCPIP "magalu.cloud/lib/products/network/vpcs/public_ips"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type NetworkPublicIPModel struct {
	Id          types.String `tfsdk:"id"`
	PublicIP    types.String `tfsdk:"public_ip"`
	Description types.String `tfsdk:"description"`
	VPCId       types.String `tfsdk:"vpc_id"`
}

type NetworkPublicIPResource struct {
	sdkClient     *mgcSdk.Client
	networkPIP    networkPIP.Service
	networkVPCPIP networkVPCPIP.Service
}

func NewNetworkPublicIPResource() resource.Resource {
	return &NetworkPublicIPResource{}
}

func (r *NetworkPublicIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_public_ips"
}

func (r *NetworkPublicIPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.networkVPCPIP = networkVPCPIP.NewService(ctx, r.sdkClient)
}

func (r *NetworkPublicIPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network Public IP",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the public IP",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_ip": schema.StringAttribute{
				Description: "The public IP address",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the public IP",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_id": schema.StringAttribute{
				Description: "The related VPC ID",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *NetworkPublicIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetworkPublicIPModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdPIP, err := r.networkVPCPIP.CreateContext(ctx, networkVPCPIP.CreateParameters{
		Description: data.Description.ValueStringPointer(),
		VpcId:       data.VPCId.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkVPCPIP.CreateConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to create Public IP", err.Error())
		return
	}

	pip, err := r.networkPIP.GetContext(ctx, networkPIP.GetParameters{
		PublicIpId: createdPIP.Id,
	},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkPIP.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch public IP address", err.Error())
	}

	data.Id = types.StringPointerValue(pip.Id)
	data.PublicIP = types.StringPointerValue(pip.PublicIp)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *NetworkPublicIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetworkPublicIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pip, err := r.networkPIP.GetContext(ctx, networkPIP.GetParameters{
		PublicIpId: data.Id.ValueString(),
	},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkPIP.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to read public IP", err.Error())
		return
	}

	data.Id = types.StringPointerValue(pip.Id)
	data.PublicIP = types.StringPointerValue(pip.PublicIp)
	data.Description = types.StringPointerValue(pip.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *NetworkPublicIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetworkPublicIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.networkPIP.DeleteContext(ctx, networkPIP.DeleteParameters{
		PublicIpId: data.Id.ValueString(),
	},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkPIP.DeleteConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to delete public IP", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *NetworkPublicIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update is not supported for public IP", "")
}

func (r *NetworkPublicIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	data := NetworkPublicIPModel{}

	pip, err := r.networkPIP.GetContext(ctx, networkPIP.GetParameters{
		PublicIpId: id,
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkPIP.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to import public IP", err.Error())
		return
	}

	data.Id = types.StringPointerValue(pip.Id)
	data.PublicIP = types.StringPointerValue(pip.PublicIp)
	data.Description = types.StringPointerValue(pip.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
