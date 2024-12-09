package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkBlockStorageVolumes "magalu.cloud/lib/products/block_storage/volumes"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

var _ datasource.DataSource = &DataSourceBsVolumes{}

type DataSourceBsVolumes struct {
	sdkClient *mgcSdk.Client
	bsVolumes sdkBlockStorageVolumes.Service
}

func NewDataSourceBsVolumes() datasource.DataSource {
	return &DataSourceBsVolumes{}
}

func (r *DataSourceBsVolumes) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_storage_volumes"
}

type bsVolumesResourceModel struct {
	volumes []bsVolumeResourceModel `tfsdk:"volumes"`
}

func (r *DataSourceBsVolumes) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.bsVolumes = sdkBlockStorageVolumes.NewService(ctx, r.sdkClient)
}

func (r *DataSourceBsVolumes) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"volumes": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of available Block Storage Volumes.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: GetBsVolumeAttributes(false),
				},
			},
		},
	}
}

func (r *DataSourceBsVolumes) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data bsVolumesResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sdkOutputList, err := r.bsVolumes.ListContext(ctx, sdkBlockStorageVolumes.ListParameters{},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	for _, sdkOutput := range sdkOutputList.Volumes {

		var item bsVolumeResourceModel
		item.ID = types.StringValue(sdkOutput.Id)
		item.Name = types.StringValue(sdkOutput.Name)
		item.AvailabilityZone = types.StringValue(sdkOutput.AvailabilityZone)
		item.UpdatedAt = types.StringValue(sdkOutput.UpdatedAt)
		item.CreatedAt = types.StringValue(sdkOutput.CreatedAt)
		item.Size = types.Int64Value(int64(sdkOutput.Size))
		item.Type = &bsVolumeType{
			DiskType: types.StringPointerValue(sdkOutput.Type.DiskType),
			Id:       types.StringValue(sdkOutput.Type.Id),
			Name:     types.StringPointerValue(sdkOutput.Type.Name),
			Status:   types.StringPointerValue(sdkOutput.Type.Status),
		}
		item.State = types.StringValue(sdkOutput.State)
		item.Status = types.StringValue(sdkOutput.Status)

		data.volumes = append(data.volumes, item)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
