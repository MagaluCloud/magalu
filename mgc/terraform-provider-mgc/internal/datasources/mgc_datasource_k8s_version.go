package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkVersion "magalu.cloud/lib/products/kubernetes/version"
	"magalu.cloud/sdk"
)

type VersionsModel struct {
	Versions []VersionModel `tfsdk:"versions"`
}

type VersionModel struct {
	Deprecated types.Bool   `tfsdk:"deprecated"`
	Version    types.String `tfsdk:"version"`
}

type DataSourceKubernetesVersion struct {
	sdkClient *mgcSdk.Client
	nodepool  sdkVersion.Service
}

func NewDataSourceKubernetesVersion() datasource.DataSource {
	return &DataSourceKubernetesVersion{}
}

func (r *DataSourceKubernetesVersion) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_version"
}

func (r *DataSourceKubernetesVersion) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	sdk, ok := req.ProviderData.(*sdk.Sdk)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected provider config, got: %T. Please report this issue to the provider developers.",
		)
		return
	}

	r.sdkClient = mgcSdk.NewClient(sdk)
	r.nodepool = sdkVersion.NewService(ctx, r.sdkClient)
}

func (r *DataSourceKubernetesVersion) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"versions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of available Kubernetes versions.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"deprecated": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the version is deprecated.",
						},
						"version": schema.StringAttribute{
							Computed:    true,
							Description: "The Kubernetes version.",
						},
					},
				},
			},
		},
	}
	resp.Schema.Description = "Get the available versions of Kubernetes."
}

func (r *DataSourceKubernetesVersion) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VersionsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutput, err := r.nodepool.List(sdkVersion.ListConfigs{})
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	for _, version := range sdkOutput.Results {
		data.Versions = append(data.Versions, VersionModel{
			Deprecated: types.BoolValue(version.Deprecated),
			Version:    types.StringValue(version.Version),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
