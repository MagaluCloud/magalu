package resources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	networkSubenetpools "magalu.cloud/lib/products/network/subnetpools"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type NetworkSubnetPoolModel struct {
	ID          types.String `tfsdk:"id"`
	Cidr        types.String `tfsdk:"cidr"`
	Description types.String `tfsdk:"description"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
}

type mgcNetworkSubnetpoolsResource struct {
	sdkClient          *mgcSdk.Client
	subnetPoolsService networkSubenetpools.Service
}

func NewNetworkSubnetpoolsResource() resource.Resource {
	return &mgcNetworkSubnetpoolsResource{}
}

func (r *mgcNetworkSubnetpoolsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_subnetpools"
}

func (r *mgcNetworkSubnetpoolsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network Subnet Pool",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the subnet pool",
				Computed:    true,
			},
			"cidr": schema.StringAttribute{
				Description: "The CIDR block of the subnet pool",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("172.26.0.0/16"),
				Validators: []validator.String{
					tfutil.CidrValidator{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the subnet pool",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the subnet pool",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of the subnet pool. Possible values are 'pip' (Public IP) and 'default'",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("default"),
				Validators: []validator.String{
					stringvalidator.OneOf("pip", "default"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *mgcNetworkSubnetpoolsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.subnetPoolsService = networkSubenetpools.NewService(ctx, r.sdkClient)
}

func (r *mgcNetworkSubnetpoolsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	data := NetworkSubnetPoolModel{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	creatParam := networkSubenetpools.CreateParameters{
		Cidr:        data.Cidr.ValueStringPointer(),
		Description: data.Description.ValueString(),
		Name:        data.Name.ValueString(),
		Type:        data.Type.ValueStringPointer(),
	}

	subnetPool, err := r.subnetPoolsService.CreateContext(ctx, creatParam,
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSubenetpools.CreateConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to create subnet pool", err.Error())
		return
	}

	data.ID = types.StringValue(subnetPool.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mgcNetworkSubnetpoolsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetworkSubnetPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	subnetPool, err := r.subnetPoolsService.GetContext(ctx, networkSubenetpools.GetParameters{
		SubnetpoolId: data.ID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSubenetpools.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to read subnet pool", err.Error())
		return
	}

	data.Cidr = types.StringPointerValue(subnetPool.Cidr)
	data.Description = types.StringValue(subnetPool.Description)
	data.Name = types.StringValue(subnetPool.Name)
	data.Type = types.StringValue(isDefaultConverter(subnetPool.IsDefault))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *mgcNetworkSubnetpoolsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update is not supported for subnet pools", "")
}

func (r *mgcNetworkSubnetpoolsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetworkSubnetPoolModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.subnetPoolsService.DeleteContext(ctx, networkSubenetpools.DeleteParameters{
		SubnetpoolId: data.ID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSubenetpools.DeleteConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to delete subnet pool", err.Error())
		return
	}
}

func (r *mgcNetworkSubnetpoolsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	subnetPoolId := req.ID
	data := NetworkSubnetPoolModel{}

	subnetPool, err := r.subnetPoolsService.GetContext(ctx, networkSubenetpools.GetParameters{
		SubnetpoolId: subnetPoolId,
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSubenetpools.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to import subnet pool", err.Error())
		return
	}

	data.ID = types.StringValue(subnetPool.Id)
	data.Cidr = types.StringPointerValue(subnetPool.Cidr)
	data.Description = types.StringValue(subnetPool.Description)
	data.Name = types.StringValue(subnetPool.Name)
	data.Type = types.StringValue(isDefaultConverter(subnetPool.IsDefault))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func isDefaultConverter(isDefault bool) string {
	if isDefault {
		return "default"
	}
	return "pip"
}
