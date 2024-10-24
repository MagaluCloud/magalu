package resources

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	networkPIP "magalu.cloud/lib/products/network/publicIps"
	networkVPCPIP "magalu.cloud/lib/products/network/vpcs/publicIps"
)

type NetworkPublicIPModel struct {
	Id          types.String `tfsdk:"id"`
	PublicIP    types.String `tfsdk:"public_ip"`
	Description types.String `tfsdk:"description"`
	VPCId       types.String `tfsdk:"vpc_id"`

	//  "id": "fe140324-7795-4d8b-b607-ce6b9096ce4a",
	//  "description": "Created With Port: port-test-67474460-c42e-4a11-b1c7-0465a55002e1",
	//  "public_ip": "201.23.18.173",
	//  "vpc_id": "4ed41b5b-bf84-4a81-a16a-6193ce626f0c"

	//  "port_id": "945519cc-42d2-439c-a9a1-eb999c098204",
	//  "error": null,
	//  "external_id": null,
	//  "project_type": null,
	//  "status": "created",
	//  "tenant_id": null,
}

type NetworkPublicIPResource struct {
	sdkClient  *mgcSdk.Client
	networkPIP networkPIP.Service
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
	r.sdkClient, err, errDetail = client.NewSDKClient(req)
	if err != nil {
		resp.Diagnostics.AddError(err.Error(), errDetail.Error())
		return
	}

	r.networkPIP = networkPIP.NewService(ctx, r.sdkClient)
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

	createdPIP, err := r.networkPIP.CreateContext(ctx, convertCreateTFModelToSDKModel(data), tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkPIP.CreateConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to create Public IP", err.Error())
		return
	}

	pip, err := getPublicIPData(r.sdkClient, ctx, createdPIP.Id)
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch public IP address")
	}

	data.Id = types.StringPointerValue(pip.Id)
	data.PublicIP = types.StringPointerValue(pip.PublicIp)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func convertCreateTFModelToSDKModel(create NetworkPublicIPModel) networkPIP.CreateParameters {
	return networkPIP.CreateParameters{
		Description: create.Description.ValueStringPointer(),
		VpcId:       create.VPCId.ValueStringPointer(),
	}
}

func getPublicIPData(sdkClient mgcSdk.Client, ctx context.Context, pipId string) (
	result GetResult,
	err error,
) {
	return r.networkPIP.GetContext(ctx, networkPIP.GetParameters{
		PublicIpId: pipId,
	},
		tfutil.GetConfigsFromTags(sdkClient.Sdk().Config().Get, networkPIP.GetConfigs{}))
}

func (r *NetworkPublicIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetworkPublicIPModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pip, err := getPublicIPData(r.sdkClient, ctx, createdPIP.Id)
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
	}), tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkPIP.GetConfigs{})

	if err != nil {
		resp.Diagnostics.AddError("Failed to import public IP", err.Error())
		return
	}

	data.Id = types.StringPointerValue(pip.Id)
	data.PublicIP = types.StringPointerValue(pip.PublicIp)
	data.Description = types.StringPointerValue(pip.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
