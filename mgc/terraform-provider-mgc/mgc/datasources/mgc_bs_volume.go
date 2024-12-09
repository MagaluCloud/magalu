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

var _ datasource.DataSource = &DataSourceBsVolume{}

type DataSourceBsVolume struct {
	sdkClient *mgcSdk.Client
	bsVolumes sdkBlockStorageVolumes.Service
}

func NewDataSourceBsVolume() datasource.DataSource {
	return &DataSourceBsVolume{}
}

func (r *DataSourceBsVolume) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_storage_volume"
}

type bsVolumeResourceModel struct {
	ID               types.String  `tfsdk:"id"`
	Name             types.String  `tfsdk:"name"`
	AvailabilityZone types.String  `tfsdk:"availability_zones"`
	UpdatedAt        types.String  `tfsdk:"updated_at"`
	CreatedAt        types.String  `tfsdk:"created_at"`
	Size             types.Int64   `tfsdk:"size"`
	Type             *bsVolumeType `tfsdk:"type"`
	State            types.String  `tfsdk:"state"`
	Status           types.String  `tfsdk:"status"`
}

type bsVolumeType struct {
	DiskType types.String `tfsdk:"disk_type"`
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Status   types.String `tfsdk:"status"`
}

func (r *DataSourceBsVolume) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func GetBsVolumeAttributes(idRequired bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The unique identifier of the volume snapshot.",
			Required:    idRequired,
			Computed:    !idRequired,
		},
		"name": schema.StringAttribute{
			Description: "The name of the block storage.",
			Computed:    true,
		},
		"availability_zone": schema.StringAttribute{
			Description: "The availability zones where the block storage is available.",
			Computed:    true,
		},
		"size": schema.Int64Attribute{
			Description: "The size of the block storage in GB.",
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
		"type": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "The type of the block storage.",
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					Description: "The name of the block storage type.",
					Computed:    true,
				},
				"disk_type": schema.StringAttribute{
					Description: "The disk type of the block storage.",
					Computed:    true,
				},
				"id": schema.StringAttribute{
					Description: "The unique identifier of the block storage type.",
					Computed:    true,
				},
				"status": schema.StringAttribute{
					Description: "The status of the block storage type.",
					Computed:    true,
				},
			},
		},
	}
}

func (r *DataSourceBsVolume) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Block storage volumes",
		MarkdownDescription: "Block storage volumes",
		Attributes:          GetBsVolumeAttributes(true),
	}
}

func (r *DataSourceBsVolume) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data bsVolumeResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sdkOutput, err := r.bsVolumes.GetContext(ctx, sdkBlockStorageVolumes.GetParameters{Id: data.ID.ValueString()},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	data.ID = types.StringValue(sdkOutput.Id)
	data.Name = types.StringValue(sdkOutput.Name)
	data.AvailabilityZone = types.StringValue(sdkOutput.AvailabilityZone)
	data.UpdatedAt = types.StringValue(sdkOutput.UpdatedAt)
	data.CreatedAt = types.StringValue(sdkOutput.CreatedAt)
	data.Size = types.Int64Value(int64(sdkOutput.Size))
	data.Type = &bsVolumeType{
		DiskType: types.StringPointerValue(sdkOutput.Type.DiskType),
		Id:       types.StringValue(sdkOutput.Type.Id),
		Name:     types.StringPointerValue(sdkOutput.Type.Name),
		Status:   types.StringPointerValue(sdkOutput.Type.Status),
	}
	data.State = types.StringValue(sdkOutput.State)
	data.Status = types.StringValue(sdkOutput.Status)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
