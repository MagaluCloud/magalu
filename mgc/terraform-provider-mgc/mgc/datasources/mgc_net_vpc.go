package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	networkVpc "magalu.cloud/lib/products/network/vpcs"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type NetworkVPCModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type NetworkVPCDatasource struct {
	sdkClient  *mgcSdk.Client
	networkVPC networkVpc.Service
}

func NewDataSourceNetworkVPC() datasource.DataSource {
	return &NetworkVPCDatasource{}
}

func (r *NetworkVPCDatasource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_vpc"
}

func (r *NetworkVPCDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.networkVPC = networkVpc.NewService(ctx, r.sdkClient)
}

func (r *NetworkVPCDatasource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Network VPC",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the VPC",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the VPC",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the VPC",
				Computed:    true,
			},
		},
	}
}

func (r *NetworkVPCDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := &NetworkVPCModel{}
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	vpc, err := r.networkVPC.GetContext(ctx, networkVpc.GetParameters{
		VpcId: data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, networkVpc.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("unable to get VPC", err.Error())
		return
	}

	data.Id = types.StringPointerValue(vpc.Id)
	data.Name = types.StringPointerValue(vpc.Name)
	data.Description = types.StringPointerValue(vpc.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
