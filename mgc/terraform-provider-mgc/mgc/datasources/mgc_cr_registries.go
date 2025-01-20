package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	sdkCRRegistries "github.com/MagaluCloud/magalu/mgc/lib/products/container_registry/registries"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DataSourceCRRegistries{}

type DataSourceCRRegistries struct {
	sdkClient    *mgcSdk.Client
	crRegistries sdkCRRegistries.Service
}

func NewDataSourceCRRegistries() datasource.DataSource {
	return &DataSourceCRRegistries{}
}

func (r *DataSourceCRRegistries) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_container_registries"
}

type crRegistries struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	StorageUsageBytes types.Int64  `tfsdk:"storage_usage_bytes"`
	CreatedAt         types.String `tfsdk:"created_at"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
}

type crRegistriesList struct {
	Registries []crRegistries `tfsdk:"registries"`
}

func (r *DataSourceCRRegistries) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.crRegistries = sdkCRRegistries.NewService(ctx, r.sdkClient)
}

func (r *DataSourceCRRegistries) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for Container Registry lists",
		Attributes: map[string]schema.Attribute{
			"registries": schema.ListNestedAttribute{
				Description: "List of container registries",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier of the registry",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the registry",
							Computed:    true,
						},
						"storage_usage_bytes": schema.Int64Attribute{
							Description: "The storage usage in bytes",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "The timestamp when the registry was created",
							Computed:    true,
						},
						"updated_at": schema.StringAttribute{
							Description: "The timestamp when the registry was last updated",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *DataSourceCRRegistries) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data crRegistriesList

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutputList, err := r.crRegistries.ListContext(ctx, sdkCRRegistries.ListParameters{}, sdkCRRegistries.ListConfigs{})
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	for _, sdkOutput := range sdkOutputList.Results {

		var item crRegistries

		item.ID = types.StringValue(sdkOutput.Id)
		item.Name = types.StringValue(sdkOutput.Name)
		item.UpdatedAt = types.StringValue(sdkOutput.UpdatedAt)
		item.CreatedAt = types.StringValue(sdkOutput.CreatedAt)

		data.Registries = append(data.Registries, item)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
