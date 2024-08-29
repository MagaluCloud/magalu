package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkObjects "magalu.cloud/lib/products/object_storage/objects"
	"magalu.cloud/sdk"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

var _ datasource.DataSource = &DatasourceObject{}

type DatasourceObject struct {
	sdkClient *mgcSdk.Client
	objects   sdkObjects.Service
}

func NewDatasourceObject() datasource.DataSource {
	return &DatasourceObject{}
}

func (r *DatasourceObject) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_object_storage_object"
}

func (r *DatasourceObject) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	sdk, ok := req.ProviderData.(*sdk.Sdk)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sdk.Sdk, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.sdkClient = sdk.NewClient()
	r.objects = sdk.Objects()
}

func (r *DatasourceObject) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"key": schema.StringAttribute{
				Computed:    true,
				Description: "Object key",
			},
			"last_modified": schema.StringAttribute{
				Computed:    true,
				Description: "Object last modified",
			},
			"size": schema.Int64Attribute{
				Computed:    true,
				Description: "Object size",
			},
			"storage_class": schema.StringAttribute{
				Computed:    true,
				Description: "Object storage class",
			},
		},
	}
}

func (r *DatasourceObject) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ObjectModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	object, err := r.objects.Download(
		sdkObjects.DownloadParameters{},
		tfutil.GetConfigsFromTags(mgcSdk.DefaultSdk().Config().Get, sdkObjects.DownloadConfigs{}),
	)

	if err != nil {
		resp.Diagnostics.AddError("Error reading object", err.Error())
		return
	}

	data.Key = types.StringValue(object.Key)
	data.LastModified = types.StringValue(object.LastModified)
	data.Size = types.Int64Value(object.Size)
	data.StorageClass = types.StringValue(object.StorageClass)
}
