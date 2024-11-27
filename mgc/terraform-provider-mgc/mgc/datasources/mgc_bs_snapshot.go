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

type snapshotsModel struct {
	Snapshots []bsSnapshotsResourceModel `tfsdk:"snapshots"`
}

// bsSnapshotsResourceModel maps de resource schema data.
type bsSnapshotsResourceModel struct {
	ID                types.String              `tfsdk:"id"`
	Name              types.String              `tfsdk:"name"`
	Description       types.String              `tfsdk:"description"`
	UpdatedAt         types.String              `tfsdk:"updated_at"`
	CreatedAt         types.String              `tfsdk:"created_at"`
	Volume            *bsSnapshotsVolumeIDModel `tfsdk:"volume"`
	State             types.String              `tfsdk:"state"`
	Status            types.String              `tfsdk:"status"`
	Size              types.Int64               `tfsdk:"size"`
	Type              types.String              `tfsdk:"type"`
	AvailabilityZones types.List                `tfsdk:"availability_zones"`
}

type bsSnapshotsVolumeIDModel struct {
	ID types.String `tfsdk:"id"`
}

func (r *DataSourceBsSnapshots) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.bsSnapshots = sdkBlockStorageSnapshots.NewService(ctx, r.sdkClient)
}

func (r *DataSourceBsSnapshots) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	description := "Block storage snapshots"
	resp.Schema = schema.Schema{
		Description:         description,
		MarkdownDescription: description,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the volume snapshot.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the volume snapshot.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the volume snapshot.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "The timestamp when the block storage was last updated.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the block storage was created.",
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: "The current state of the virtual machine instance.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the virtual machine instance.",
				Computed:    true,
			},
			"size": schema.Int64Attribute{
				Description: "The size of the snapshot in GB.",
				Computed:    true,
			},
			"volume": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "ID of block storage volume",
						Computed:    true,
					},
				},
			},
			"snapshot_source_id": schema.StringAttribute{
				Description: "The ID of the snapshot source.",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the snapshot.",
				Computed:    true,
			},
			"availability_zones": schema.ListAttribute{
				Description: "The availability zones of the snapshot.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *DataSourceBsSnapshots) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data snapshotsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutput, err := r.bsSnapshots.ListContext(ctx, sdkBlockStorageSnapshots.ListParameters{},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageSnapshots.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	for _, snap := range sdkOutput.Snapshots {
		list, diags := types.ListValueFrom(ctx, types.StringType, snap.AvailabilityZones)
		resp.Diagnostics.Append(diags...)

		data.Snapshots = append(data.Snapshots, bsSnapshotsResourceModel{
			ID:                types.StringValue(snap.Id),
			Name:              types.StringValue(snap.Name),
			Description:       types.StringPointerValue(snap.Description),
			UpdatedAt:         types.StringValue(snap.UpdatedAt),
			CreatedAt:         types.StringValue(snap.CreatedAt),
			Volume:            &bsSnapshotsVolumeIDModel{ID: types.StringValue(snap.Volume.Id)},
			State:             types.StringValue(snap.State),
			Status:            types.StringValue(snap.Status),
			Size:              types.Int64Value(int64(snap.Size)),
			Type:              types.StringValue(snap.Type),
			AvailabilityZones: list,
		})

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
