package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	networkInterfaces "magalu.cloud/lib/products/network/ports"
	networkVpcInterfaces "magalu.cloud/lib/products/network/vpcs/ports"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type NetworkVPCInterfaceModel struct {
	Id    types.String `tfsdk:"id"`
	VpcId types.String `tfsdk:"vpc_id"`
	Name  types.String `tfsdk:"name"`
}

type NetworkVPCInterfaceResource struct {
	sdkClient           *mgcSdk.Client
	networkInterfaces   networkInterfaces.Service
	networkVPCInterface networkVpcInterfaces.Service
}

func NewNetworkVPCInterfaceResource() resource.Resource {
	return &NetworkVPCInterfaceResource{}
}

func (r *NetworkVPCInterfaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_vpcs_interfaces"
}

func (r *NetworkVPCInterfaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	var err error
	var errDetail error
	r.sdkClient, err, errDetail = client.NewSDKClient(req, resp)
	if err != nil {
		resp.Diagnostics.AddError(
			err.Error(),
			errDetail.Error(),
		)
		return
	}

	r.networkInterfaces = networkInterfaces.NewService(ctx, r.sdkClient)
	r.networkVPCInterface = networkVpcInterfaces.NewService(ctx, r.sdkClient)
}

func (r *NetworkVPCInterfaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network VPC Interface",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the VPC Interface",
				Computed:    true,
			},
			"vpc_id": schema.StringAttribute{
				Description: "The ID of the VPC",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the VPC Interface",
				Required:    true,
			},
		},
	}
}

func (r *NetworkVPCInterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model NetworkVPCInterfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vpcInterface, err := r.networkInterfaces.GetContext(ctx, networkInterfaces.GetParameters{
		PortId: model.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkInterfaces.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("unable to get VPC Interface", err.Error())
		return
	}

	model.Name = types.StringPointerValue(vpcInterface.Name)
	model.VpcId = types.StringPointerValue(vpcInterface.VpcId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *NetworkVPCInterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model NetworkVPCInterfaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	disableDefaultRules := false
	createdVPCInterface, err := r.networkVPCInterface.CreateContext(ctx, networkVpcInterfaces.CreateParameters{
		VpcId:  model.VpcId.ValueString(),
		Name:   model.Name.ValueString(),
		HasPip: &disableDefaultRules,
		HasSg:  &disableDefaultRules,
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkVpcInterfaces.CreateConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create VPC Interface", err.Error())
		return
	}

	model.Id = types.StringValue(createdVPCInterface.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func (r *NetworkVPCInterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update is not supported for VPC Interfaces", "")
}

func (r *NetworkVPCInterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model NetworkVPCInterfaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.networkInterfaces.DeleteContext(ctx, networkInterfaces.DeleteParameters{
		PortId: model.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkInterfaces.DeleteConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete VPC Interface", err.Error())
		return
	}
}

func (r *NetworkVPCInterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	vpcInterfaceId := req.ID
	model := NetworkVPCInterfaceModel{}

	vpcInterface, err := r.networkInterfaces.GetContext(ctx, networkInterfaces.GetParameters{
		PortId: vpcInterfaceId,
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkInterfaces.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to import VPC Interface", err.Error())
		return
	}

	model.Id = types.StringPointerValue(vpcInterface.Id)
	model.Name = types.StringPointerValue(vpcInterface.Name)
	model.VpcId = types.StringPointerValue(vpcInterface.VpcId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}
