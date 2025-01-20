package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	sdkProfileAvailabilityZones "github.com/MagaluCloud/magalu/mgc/lib/products/profile/availability_zones"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DataSourceAvailabilityZones{}

type DataSourceAvailabilityZones struct {
	sdkClient         *mgcSdk.Client
	profileAvailZones sdkProfileAvailabilityZones.Service
}

func NewDataSourceAvailabilityZones() datasource.DataSource {
	return &DataSourceAvailabilityZones{}
}

func (r *DataSourceAvailabilityZones) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_availability_zones"
}

type Regions struct {
	Regions []Region `tfsdk:"regions"`
}

type Region struct {
	Region            types.String        `tfsdk:"region"`
	AvailabilityZones []AvailabilityZones `tfsdk:"availability_zones"`
}

type AvailabilityZones struct {
	AvailabilityZone types.String `tfsdk:"availability_zone"`
	BlockType        types.String `tfsdk:"block_type"`
}

func (r *DataSourceAvailabilityZones) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.profileAvailZones = sdkProfileAvailabilityZones.NewService(ctx, r.sdkClient)
}

func (r *DataSourceAvailabilityZones) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of available regions and availability zones.",
		Attributes: map[string]schema.Attribute{
			"regions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Computed: true,
						},
						"availability_zones": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"availability_zone": schema.StringAttribute{
										Computed: true,
									},
									"block_type": schema.StringAttribute{
										Computed: true,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
func (r *DataSourceAvailabilityZones) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Regions

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutput, err := r.profileAvailZones.ListContext(ctx, sdkProfileAvailabilityZones.ListParameters{}, sdkProfileAvailabilityZones.ListConfigs{})
	if err != nil {
		resp.Diagnostics.AddError("Failed to get versions", err.Error())
		return
	}

	for _, region := range sdkOutput.Results {
		var regionData Region
		regionData.Region = types.StringValue(region.RegionId)
		for _, az := range region.AvailabilityZones {
			regionData.AvailabilityZones = append(regionData.AvailabilityZones, AvailabilityZones{
				AvailabilityZone: types.StringValue(az.AzId),
				BlockType:        types.StringPointerValue(removeNoneFromString(az.BlockType)),
			})
		}
		data.Regions = append(data.Regions, regionData)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func removeNoneFromString(str string) *string {
	if str == "none" {
		return nil
	}
	return &str
}
