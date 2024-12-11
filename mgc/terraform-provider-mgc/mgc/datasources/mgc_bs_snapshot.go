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

var _ datasource.DataSource = &DataSourceBsSnapshot{}

type DataSourceBsSnapshot struct {
	sdkClient   *mgcSdk.Client
	bsSnapshots sdkBlockStorageSnapshots.Service
}

func NewDataSourceBSSnapshot() datasource.DataSource {
	return &DataSourceBsSnapshot{}
}

func (r *DataSourceBsSnapshot) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_storage_snapshot"
}

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

func (r *DataSourceBsSnapshot) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func GetBsSnapshotAttributes(idRequired bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the volume snapshot.",
			Required:    idRequired,
			Computed:    !idRequired,
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
		"type": schema.StringAttribute{
			Description: "The type of the snapshot.",
			Computed:    true,
		},
		"availability_zones": schema.ListAttribute{
			Description: "The availability zones of the snapshot.",
			Computed:    true,
			ElementType: types.StringType,
		},
	}
}

func (r *DataSourceBsSnapshot) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	description := "Block storage snapshots"
	resp.Schema = schema.Schema{
		Description:         description,
		MarkdownDescription: description,
		Attributes:          GetBsSnapshotAttributes(true),
	}
}

func (r *DataSourceBsSnapshot) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data bsSnapshotsResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutput, err := r.bsSnapshots.GetContext(ctx, sdkBlockStorageSnapshots.GetParameters{Id: data.ID.ValueString()},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageSnapshots.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	list, diags := types.ListValueFrom(ctx, types.StringType, sdkOutput.AvailabilityZones)
	resp.Diagnostics.Append(diags...)

	data.ID = types.StringValue(sdkOutput.Id)
	data.Name = types.StringValue(sdkOutput.Name)
	data.Description = types.StringPointerValue(sdkOutput.Description)
	data.UpdatedAt = types.StringValue(sdkOutput.UpdatedAt)
	data.CreatedAt = types.StringValue(sdkOutput.CreatedAt)
	data.Volume = &bsSnapshotsVolumeIDModel{ID: types.StringPointerValue(sdkOutput.Volume.Id)}
	data.State = types.StringValue(sdkOutput.State)
	data.Status = types.StringValue(sdkOutput.Status)
	data.Size = types.Int64Value(int64(sdkOutput.Size))
	data.Type = types.StringValue(sdkOutput.Type)
	data.AvailabilityZones = list

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
