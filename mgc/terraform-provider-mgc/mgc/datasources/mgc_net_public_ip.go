package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	networkPIP "magalu.cloud/lib/products/network/public_ips"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type NetworkPublicIPModel struct {
	Id          types.String `tfsdk:"id"`
	Description types.String `tfsdk:"description"`
	PublicIP    types.String `tfsdk:"public_ip"`
	VPCId       types.String `tfsdk:"vpc_id"`
	PortId      types.String `tfsdk:"port_id"`
}

type NetworkPublicIPDataSource struct {
	sdkClient  *mgcSdk.Client
	networkPIP networkPIP.Service
}

func NewDataSourceNetworkPublicIP() datasource.DataSource {
	return &NetworkPublicIPDataSource{}
}

func (r *NetworkPublicIPDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_public_ip"
}

func (r *NetworkPublicIPDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.networkPIP = networkPIP.NewService(ctx, r.sdkClient)
}

func (r *NetworkPublicIPDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network Public IP",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the public IP",
				Required:    true,
			},
			"public_ip": schema.StringAttribute{
				Description: "The public IP address",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the public IP",
				Computed:    true,
			},
			"vpc_id": schema.StringAttribute{
				Description: "The related VPC ID",
				Computed:    true,
			},
			"port_id": schema.StringAttribute{
				Description: "The port ID it's attached to",
				Computed:    true,
			},
		},
	}
}

func (r *NetworkPublicIPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &NetworkPublicIPModel{}
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	pip, err := r.networkPIP.GetContext(ctx, networkPIP.GetParameters{
		PublicIpId: data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkPIP.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Unable to get public IP", err.Error())
		return
	}

	data.Id = types.StringPointerValue(pip.Id)
	data.Description = types.StringPointerValue(pip.Description)
	data.PublicIP = types.StringPointerValue(pip.PublicIp)
	data.VPCId = types.StringPointerValue(pip.VpcId)
	data.PortId = types.StringPointerValue(pip.PortId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
