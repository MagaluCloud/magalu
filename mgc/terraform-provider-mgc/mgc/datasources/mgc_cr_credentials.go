package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	sdkCRCredentials "github.com/MagaluCloud/magalu/mgc/lib/products/container_registry/credentials"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DataSourceCRCredentials{}

type DataSourceCRCredentials struct {
	sdkClient     *mgcSdk.Client
	crCredentials sdkCRCredentials.Service
}

func NewDataSourceCRCredentials() datasource.DataSource {
	return &DataSourceCRCredentials{}
}

func (r *DataSourceCRCredentials) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_credentials"
}

type crCredentials struct {
	Email    types.String `tfsdk:"email"`
	Password types.String `tfsdk:"password"`
	Username types.String `tfsdk:"username"`
}

func (r *DataSourceCRCredentials) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.crCredentials = sdkCRCredentials.NewService(ctx, r.sdkClient)
}

func (r *DataSourceCRCredentials) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Credentials for Container Registry authentication",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Description: "Email address for the credentials",
				Computed:    true,
				Sensitive:   false,
			},
			"password": schema.StringAttribute{
				Description: "Password for authentication",
				Computed:    true,
				Sensitive:   true,
			},
			"username": schema.StringAttribute{
				Description: "Username for authentication",
				Computed:    true,
				Sensitive:   false,
			},
		},
	}
}

func (r *DataSourceCRCredentials) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data crCredentials

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutputList, err := r.crCredentials.ListContext(ctx, sdkCRCredentials.ListConfigs{})
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	data.Email = types.StringValue(sdkOutputList.Email)
	data.Password = types.StringValue(sdkOutputList.Password)
	data.Username = types.StringValue(sdkOutputList.Username)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
