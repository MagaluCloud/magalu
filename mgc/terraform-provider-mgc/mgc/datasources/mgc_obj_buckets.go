package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	sdkBuckets "github.com/MagaluCloud/magalu/mgc/lib/products/object_storage/buckets"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/tfutil"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
