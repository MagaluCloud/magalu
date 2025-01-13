package resources

import (
	"context"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	networkSubnetpools "github.com/MagaluCloud/magalu/mgc/lib/products/network/subnetpools"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/tfutil"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type NetworkSubnetPoolsBookCIDRModel struct {
	SubnetPoolID types.String `tfsdk:"subnet_pool_id"`
	CIDR         types.String `tfsdk:"cidr"`
}

type NetworkSubnetPoolsBookCIDRResource struct {
	sdkClient          *mgcSdk.Client
	subnetPoolsService networkSubnetpools.Service
}

func NewNetworkSubnetPoolsBookCIDRResource() resource.Resource {
	return &NetworkSubnetPoolsBookCIDRResource{}
}

func (r *NetworkSubnetPoolsBookCIDRResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_subnetpools_book_cidr"
}

func (r *NetworkSubnetPoolsBookCIDRResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.subnetPoolsService = networkSubnetpools.NewService(ctx, r.sdkClient)
}

func (r *NetworkSubnetPoolsBookCIDRResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network Subnet Pools Book CIDR",
		Attributes: map[string]schema.Attribute{
			"subnet_pool_id": schema.StringAttribute{
				Description: "Subnet Pool ID",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cidr": schema.StringAttribute{
				Description: "CIDR",
				Required:    true,
			},
		},
	}
}

func (r *NetworkSubnetPoolsBookCIDRResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := NetworkSubnetPoolsBookCIDRModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.Set(ctx, &data)
}

func (r *NetworkSubnetPoolsBookCIDRResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetworkSubnetPoolsBookCIDRModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.subnetPoolsService.CreateBookCidrContext(ctx, networkSubnetpools.CreateBookCidrParameters{
		SubnetpoolId: data.SubnetPoolID.ValueString(),
		Cidr:         data.CIDR.ValueStringPointer(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSubnetpools.CreateBookCidrConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("failed to book CIDR", err.Error())
		return
	}

	resp.State.Set(ctx, &data)
}

func (r *NetworkSubnetPoolsBookCIDRResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetworkSubnetPoolsBookCIDRModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.subnetPoolsService.CreateUnbookCidrContext(ctx, networkSubnetpools.CreateUnbookCidrParameters{
		SubnetpoolId: data.SubnetPoolID.ValueString(),
		Cidr:         data.CIDR.ValueStringPointer(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSubnetpools.CreateUnbookCidrConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("failed to delete CIDR", err.Error())
		return
	}
}

func (r *NetworkSubnetPoolsBookCIDRResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update is not supported for network subnet pools book CIDR", "")
}
