package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkBlockStorageSnapshots "magalu.cloud/lib/products/block_storage/snapshots"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

var _ datasource.DataSource = &DataSourceBsSnapshots{}

type DataSourceBsSnapshots struct {
	sdkClient   *mgcSdk.Client
	bsSnapshots sdkBlockStorageSnapshots.Service
}

func NewDataSourceBSSnapshots() datasource.DataSource {
	return &DataSourceBsSnapshots{}
}

func (r *DataSourceBsSnapshots) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_storage_snapshots"
}

type bsSnapshotsResourceListModel struct {
	snapshots []bsSnapshotsResourceModel `tfsdk:"snapshots"`
}

func (r *DataSourceBsSnapshots) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.bsSnapshots = sdkBlockStorageSnapshots.NewService(ctx, r.sdkClient)
}

func (r *DataSourceBsSnapshots) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"snapshots": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of available Block Storage Snapshots.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: GetBsSnapshotAttributes(false),
				},
			},
		},
	}
}

func (r *DataSourceBsSnapshots) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data bsSnapshotsResourceListModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutputList, err := r.bsSnapshots.ListContext(ctx, sdkBlockStorageSnapshots.ListParameters{},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageSnapshots.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	for _, sdkOutput := range sdkOutputList.Snapshots {
		list, diags := types.ListValueFrom(ctx, types.StringType, sdkOutput.AvailabilityZones)
		resp.Diagnostics.Append(diags...)

		var item bsSnapshotsResourceModel

		item.ID = types.StringValue(sdkOutput.Id)
		item.Name = types.StringValue(sdkOutput.Name)
		item.Description = types.StringPointerValue(sdkOutput.Description)
		item.UpdatedAt = types.StringValue(sdkOutput.UpdatedAt)
		item.CreatedAt = types.StringValue(sdkOutput.CreatedAt)
		item.Volume = &bsSnapshotsVolumeIDModel{ID: types.StringPointerValue(sdkOutput.Volume.Id)}
		item.State = types.StringValue(sdkOutput.State)
		item.Status = types.StringValue(sdkOutput.Status)
		item.Size = types.Int64Value(int64(sdkOutput.Size))
		item.Type = types.StringValue(sdkOutput.Type)
		item.AvailabilityZones = list

		data.snapshots = append(data.snapshots, item)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
