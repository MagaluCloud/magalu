package resources

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	mgcSdk "magalu.cloud/lib"
	networkVpc "magalu.cloud/lib/products/network/vpcs"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

const NetworkPoolingTimeout = 5 * time.Minute

type NetworkVPCModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type NetworkVPCResource struct {
	sdkClient  *mgcSdk.Client
	networkVPC networkVpc.Service
}

func NewNetworkVPCResource() resource.Resource {
	return &NetworkVPCResource{}
}

func (r *NetworkVPCResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_vpcs"
}

func (r *NetworkVPCResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	var err error
	var errDetail error
	r.sdkClient, err, errDetail = client.NewSDKClient(req)
	if err != nil {
		resp.Diagnostics.AddError(
			err.Error(),
			errDetail.Error(),
		)
		return
	}

	r.networkVPC = networkVpc.NewService(ctx, r.sdkClient)
}

func (r *NetworkVPCResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network VPC",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the VPC",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the VPC",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the VPC",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *NetworkVPCResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetworkVPCModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createdVPC, err := r.networkVPC.CreateContext(ctx, convertCreateTFModelToSDKModel(data),
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkVpc.CreateConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create VPC", err.Error())
		return
	}

	for startTime := time.Now(); time.Since(startTime) < NetworkPoolingTimeout; {
		res, err := r.networkVPC.GetContext(ctx, networkVpc.GetParameters{
			VpcId: createdVPC.Id,
		},
			tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkVpc.GetConfigs{}))
		if err != nil {
			resp.Diagnostics.AddError("Failed to get VPC", err.Error())
			return
		}
		if res.Status == "created" {
			break
		}
		tflog.Info(ctx, "VPC is not yet created, waiting for 10 seconds",
			map[string]interface{}{"status": res.Status})
		time.Sleep(10 * time.Second)
	}

	data.Id = types.StringValue(createdVPC.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertCreateTFModelToSDKModel(create NetworkVPCModel) networkVpc.CreateParameters {
	return networkVpc.CreateParameters{
		Name:        create.Name.ValueString(),
		Description: create.Description.ValueStringPointer(),
	}
}

func (r *NetworkVPCResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetworkVPCModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vpc, err := r.networkVPC.GetContext(ctx, networkVpc.GetParameters{
		VpcId: data.Id.ValueString(),
	},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkVpc.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to read VPC", err.Error())
		return
	}

	if vpc.Description != nil && *vpc.Description == "" {
		vpc.Description = nil
	}

	data.Name = types.StringPointerValue(vpc.Name)
	data.Description = types.StringPointerValue(vpc.Description)
	data.Id = types.StringPointerValue(vpc.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *NetworkVPCResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetworkVPCModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.networkVPC.DeleteContext(ctx, networkVpc.DeleteParameters{
		VpcId: data.Id.ValueString(),
	},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkVpc.DeleteConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to delete VPC", err.Error())
		return
	}
}

func (r *NetworkVPCResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update is not supported for VPC", "")
}

func (r *NetworkVPCResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	vpcId := req.ID
	data := NetworkVPCModel{}

	vpc, err := r.networkVPC.GetContext(ctx, networkVpc.GetParameters{
		VpcId: vpcId,
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkVpc.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to import VPC", err.Error())
		return
	}

	data.Id = types.StringPointerValue(vpc.Id)
	data.Name = types.StringPointerValue(vpc.Name)
	data.Description = types.StringPointerValue(vpc.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
