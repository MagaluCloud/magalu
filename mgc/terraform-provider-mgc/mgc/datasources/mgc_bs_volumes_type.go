package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkBlockStorageVolumeTypes "magalu.cloud/lib/products/block_storage/volume_types"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

var _ datasource.DataSource = &DataSourceBsVolumeTypes{}

type DataSourceBsVolumeTypes struct {
	sdkClient     *mgcSdk.Client
	bsVolumeTypes sdkBlockStorageVolumeTypes.Service
}

func NewDataSourceBsVolumeTypes() datasource.DataSource {
	return &DataSourceBsVolumeTypes{}
}

func (r *DataSourceBsVolumeTypes) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_storage_volume_types"
}

type volumeTypes struct {
	VolumeTypes []volumeType `tfsdk:"types"`
}

type volumeType struct {
	AvailabilityZones []string     `tfsdk:"availability_zones"`
	DiskType          types.String `tfsdk:"disk_type"`
	Id                types.String `tfsdk:"id"`
	Iops              types.Int64  `tfsdk:"iops"`
	Name              types.String `tfsdk:"name"`
	Status            types.String `tfsdk:"status"`
}

func (r *DataSourceBsVolumeTypes) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.bsVolumeTypes = sdkBlockStorageVolumeTypes.NewService(ctx, r.sdkClient)
}

func (r *DataSourceBsVolumeTypes) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Block-storage Volume Types",
		Attributes: map[string]schema.Attribute{
			"volume_types": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of available Block-storage Volume Types.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of image.",
						},
						"disk_type": schema.StringAttribute{
							Computed:    true,
							Description: "The disk type.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The volume type name.",
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "The volume type status.",
						},
						"iops": schema.Int64Attribute{
							Computed:    true,
							Description: "The volume type IOPS.",
						},
						"availability_zones": schema.ListAttribute{
							Computed:    true,
							Description: "The volume type availability zones.",
						},
					},
				},
			},
		},
	}
}

func (r *DataSourceBsVolumeTypes) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data volumeTypes

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutput, err := r.bsVolumeTypes.ListContext(ctx, sdkBlockStorageVolumeTypes.ListParameters{},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumeTypes.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	for _, image := range sdkOutput.Types {
		data.VolumeTypes = append(data.VolumeTypes, volumeType{
			AvailabilityZones: image.AvailabilityZones,
			DiskType:          types.StringValue(image.DiskType),
			Id:                types.StringValue(image.Id),
			Iops:              types.Int64Value(int64(image.Iops.Total)),
			Name:              types.StringValue(image.Name),
			Status:            types.StringValue(image.Status),
		})

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
