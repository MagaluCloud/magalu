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

type volumes struct {
	Volumes []bsVolumesResourceModel `tfsdk:"volumes"`
}

type bsVolumesResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	AvailabilityZones types.List   `tfsdk:"availability_zones"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
	CreatedAt         types.String `tfsdk:"created_at"`
	Size              types.Int64  `tfsdk:"size"`
	Type              bsVolumeType `tfsdk:"type"`
	State             types.String `tfsdk:"state"`
	Status            types.String `tfsdk:"status"`
}

type bsVolumeType struct {
	DiskType types.String `tfsdk:"disk_type"`
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Status   types.String `tfsdk:"status"`
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
		Description: "Block storage volumes",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the block storage.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the block storage.",
				Computed:    true,
			},
			"availability_zones": schema.ListAttribute{
				Description: "The availability zones where the block storage is available.",
				Computed:    true,
				ElementType: types.StringType,
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
		},
	}
}

func (r *DataSourceBsVolumes) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data volumes

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutput, err := r.bsVolumes.ListContext(ctx, sdkBlockStorageVolumes.ListParameters{},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	for _, volume := range sdkOutput.Volumes {
		list, diags := types.ListValueFrom(ctx, types.StringType, volume.AvailabilityZones)
		resp.Diagnostics.Append(diags...)

		data.Volumes = append(data.Volumes, bsVolumesResourceModel{
			ID:                types.StringValue(volume.Id),
			Name:              types.StringValue(volume.Name),
			AvailabilityZones: list,
			UpdatedAt:         types.StringValue(volume.UpdatedAt),
			CreatedAt:         types.StringValue(volume.CreatedAt),
			Size:              types.Int64Value(int64(volume.Size)),
			Type: bsVolumeType{
				DiskType: types.StringPointerValue(volume.Type.DiskType),
				Id:       types.StringValue(volume.Type.Id),
				Name:     types.StringPointerValue(volume.Type.Name),
				Status:   types.StringPointerValue(volume.Type.Status),
			},
			State:  types.StringValue(volume.State),
			Status: types.StringValue(volume.Status),
		})

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
