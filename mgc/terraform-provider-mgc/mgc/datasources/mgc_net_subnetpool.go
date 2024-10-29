package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	networkSubnetpools "magalu.cloud/lib/products/network/subnetpools"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type mgcNetworkSubnetpoolsModel struct {
	Cidr        types.String `tfsdk:"cidr"`
	CreatedAt   types.String `tfsdk:"created_at"`
	Description types.String `tfsdk:"description"`
	Id          types.String `tfsdk:"id"`
	IpVersion   types.Int64  `tfsdk:"ip_version"`
	Name        types.String `tfsdk:"name"`
	TenantId    types.String `tfsdk:"tenant_id"`
}

type mgcNetworkSubnetpoolsDatasource struct {
	sdkClient          *mgcSdk.Client
	networkSubnetpools networkSubnetpools.Service
}

func NewDataSourceNetworkSubnetpool() datasource.DataSource {
	return &mgcNetworkSubnetpoolsDatasource{}
}

func (r *mgcNetworkSubnetpoolsDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_subnetpool"
}

func (r *mgcNetworkSubnetpoolsDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network Subnet Pool",
		Attributes: map[string]schema.Attribute{
			"cidr": schema.StringAttribute{
				Description: "The CIDR range associated with the subnetpool",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the subnetpool was created",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the subnetpool",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The unique identifier of the subnetpool",
				Required:    true,
			},
			"ip_version": schema.Int64Attribute{
				Description: "The IP version of the subnetpool (4 or 6)",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the subnetpool",
				Computed:    true,
			},
			"tenant_id": schema.StringAttribute{
				Description: "The ID of the tenant that owns the subnetpool",
				Computed:    true,
			},
		},
	}
}

func (r *mgcNetworkSubnetpoolsDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &mgcNetworkSubnetpoolsModel{}
	resp.Diagnostics.Append(req.Config.Get(ctx, data)...)

	getParam := networkSubnetpools.GetParameters{
		SubnetpoolId: data.Id.ValueString(),
	}
	subnetPool, err := r.networkSubnetpools.GetContext(ctx, getParam,
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkSubnetpools.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("unable to get subnetpool", err.Error())
		return
	}

	data.Cidr = types.StringPointerValue(subnetPool.Cidr)
	data.CreatedAt = types.StringValue(subnetPool.CreatedAt)
	data.Description = types.StringValue(subnetPool.Description)
	data.IpVersion = types.Int64PointerValue(tfutil.ConvertIntPointerToInt64Pointer(&subnetPool.IpVersion))
	data.Name = types.StringValue(subnetPool.Name)
	data.TenantId = types.StringValue(subnetPool.TenantId)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *mgcNetworkSubnetpoolsDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.networkSubnetpools = networkSubnetpools.NewService(ctx, r.sdkClient)
}
