package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkBuckets "magalu.cloud/lib/products/object_storage/buckets"
	"magalu.cloud/sdk"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

var _ datasource.DataSource = &DatasourceBuckets{}

type DatasourceBuckets struct {
	sdkClient *mgcSdk.Client
	buckets   sdkBuckets.Service
}

type BucketModel struct {
	Name         types.String `tfsdk:"name"`
	CreationDate types.String `tfsdk:"creation_date"`
}

type BucketsModel struct {
	Buckets []BucketModel `tfsdk:"ssh_keys"`
}

func NewDatasourceBuckets() datasource.DataSource {
	return &DatasourceBuckets{}
}

func (r *DatasourceBuckets) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage_buckets"
}

func (r *DatasourceBuckets) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	config, ok := req.ProviderData.(tfutil.ProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected provider config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	sdk := sdk.NewSdk()
	_ = sdk.Config().SetTempConfig("region", config.Region.ValueStringPointer())
	_ = sdk.Config().SetTempConfig("env", config.Env.ValueStringPointer())
	_ = sdk.Config().SetTempConfig("api_key", config.ApiKey.ValueStringPointer())
	if config.ObjectStorage != nil && config.ObjectStorage.ObjectKeyPair != nil {
		sdk.Config().AddTempKeyPair("apikey", config.ObjectStorage.ObjectKeyPair.KeyID.ValueString(), config.ObjectStorage.ObjectKeyPair.KeySecret.ValueString())
	}

	r.buckets = sdkBuckets.NewService(ctx, r.sdkClient)
}

func (r *DatasourceBuckets) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"buckets": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of ssh-keys.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Bucket name",
						},
						"creation_date": schema.StringAttribute{
							Computed:    true,
							Description: "Bucket creation date",
						},
					},
				},
			},
		},
	}
	resp.Schema.Description = "Get all buckets."
}

func (r *DatasourceBuckets) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BucketsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutput, err := r.buckets.ListContext(ctx, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBuckets.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	for _, key := range sdkOutput.Buckets {
		data.Buckets = append(data.Buckets, BucketModel{
			Name:         types.StringValue(key.Name),
			CreationDate: types.StringValue(key.CreationDate),
		})

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
