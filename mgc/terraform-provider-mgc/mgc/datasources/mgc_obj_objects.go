package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkObjects "magalu.cloud/lib/products/object_storage/objects"
	"magalu.cloud/sdk"
	tfutil "magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

var _ datasource.DataSource = &DatasourceObjects{}

type DatasourceObjects struct {
	sdkClient *mgcSdk.Client
	objects   sdkObjects.Service
}

type ObjectModel struct {
	Key          types.String `tfsdk:"key"`
	LastModified types.String `tfsdk:"last_modified"`
	Size         types.Int64  `tfsdk:"size"`
	StorageClass types.String `tfsdk:"storage_class"`
}

type DatasourceObjectsModel struct {
	Objects []ObjectModel `tfsdk:"objects"`
}

func NewDatasourceObjects() datasource.DataSource {
	return &DatasourceObjects{}
}

func (r *DatasourceObjects) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage_objects"
}

func (r *DatasourceObjects) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	r.objects = sdkObjects.NewService(ctx, r.sdkClient)
}

func (r *DatasourceObjects) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get all objects.",
		Attributes: map[string]schema.Attribute{
			"objects": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of objects.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							Computed:    true,
							Description: "Object name.",
						},
						"last_modified": schema.StringAttribute{
							Computed:    true,
							Description: "Object last modified date.",
						},
						"size": schema.Int64Attribute{
							Computed:    true,
							Description: "Object size.",
						},
						"storage_class": schema.StringAttribute{
							Computed:    true,
							Description: "Object storage class.",
						},
					},
				},
			},
		},
	}
}

func (r *DatasourceObjects) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatasourceObjectsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sdkOutput, err := r.objects.List(sdkObjects.ListParameters{},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkObjects.ListConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to get objects", err.Error())
		return
	}

	for _, key := range sdkOutput.Contents {
		data.Objects = append(data.Objects, ObjectModel{
			Key:          types.StringValue(key.Key),
			LastModified: types.StringValue(key.LastModified),
			Size:         types.Int64Value(int64(key.ContentSize)),
			StorageClass: types.StringValue(key.StorageClass),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
