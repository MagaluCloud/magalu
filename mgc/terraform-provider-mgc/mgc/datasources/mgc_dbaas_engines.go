package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	dbaas "magalu.cloud/lib/products/dbaas/engines"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

var _ datasource.DataSource = &DataSourceDbEngines{}

type DataSourceDbEngines struct {
	sdkClient *mgcSdk.Client
	dbEngines dbaas.Service
}

type dbEngineModel struct {
	Status  types.String `tfsdk:"status"`
	Engines []DbEngines  `tfsdk:"engines"`
}

type DbEngines struct {
	ID      types.String `tfsdk:"id"`
	Engine  types.String `tfsdk:"engine"`
	Name    types.String `tfsdk:"name"`
	Status  types.String `tfsdk:"status"`
	Version types.String `tfsdk:"version"`
}

func NewDataSourceDbaasEngines() datasource.DataSource {
	return &DataSourceDbEngines{}
}

func (r *DataSourceDbEngines) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_engines"
}

func (r *DataSourceDbEngines) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.dbEngines = dbaas.NewService(ctx, r.sdkClient)
}

func (r *DataSourceDbEngines) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A list of available database engines.",
		Attributes: map[string]schema.Attribute{
			"engines": schema.ListNestedAttribute{
				Description: "List of available database engines",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The ID of the database engine",
							Computed:    true,
						},
						"engine": schema.StringAttribute{
							Description: "The type of the database engine",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the database engine",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "The status of the database engine",
							Computed:    true,
						},
						"version": schema.StringAttribute{
							Description: "The version of the database engine",
							Computed:    true,
						},
					},
				},
			},
			"status": schema.StringAttribute{
				Description: "Filter to return engines based on their status. If set to ACTIVE, returns only active engines. If set to DEPRECATED, returns only deprecated engines.",
				Optional:    true,
				Computed:    true,
				Validators:  []validator.String{stringvalidator.OneOf("DEPRECATED", "ACTIVE")},
			},
		},
	}
}

func (r *DataSourceDbEngines) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	data := dbEngineModel{}
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	limit := 50
	params := dbaas.ListParameters{
		Status: data.Status.ValueStringPointer(),
		Limit:  &limit,
	}
	engines, err := r.dbEngines.ListContext(ctx, params, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaas.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("failed to list db engines", err.Error())
		return
	}

	engineModels := make([]DbEngines, 0, len(engines.Results))
	for _, engine := range engines.Results {
		engineModels = append(engineModels, DbEngines{
			ID:      types.StringValue(engine.Id),
			Engine:  types.StringValue(engine.Engine),
			Name:    types.StringValue(engine.Name),
			Status:  types.StringValue(engine.Status),
			Version: types.StringValue(engine.Version),
		})
	}
	data.Engines = engineModels
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
