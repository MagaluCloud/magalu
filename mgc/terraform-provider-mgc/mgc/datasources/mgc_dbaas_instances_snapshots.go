package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	dbaasSnapshots "magalu.cloud/lib/products/dbaas/instances/snapshots"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type DataSourceDbSnapshots struct {
	sdkClient *mgcSdk.Client
	snapshots dbaasSnapshots.Service
}

type dbSnapshotsModel struct {
	InstanceId types.String      `tfsdk:"instance_id"`
	Snapshots  []dbSnapshotModel `tfsdk:"snapshots"`
}

func NewDataSourceDbaasInstancesSnapshots() datasource.DataSource {
	return &DataSourceDbSnapshots{}
}

func (r *DataSourceDbSnapshots) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_instances_snapshots"
}

func (r *DataSourceDbSnapshots) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.snapshots = dbaasSnapshots.NewService(ctx, r.sdkClient)
}

func (r *DataSourceDbSnapshots) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all snapshots for a database instance.",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Description: "ID of the instance",
				Required:    true,
			},
			"snapshots": schema.ListNestedAttribute{
				Description: "List of snapshots",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the snapshot",
							Computed:    true,
						},
						"instance_id": schema.StringAttribute{
							Description: "ID of the instance",
							Required:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the snapshot",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the snapshot",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Creation timestamp",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Status of the snapshot",
							Computed:    true,
						},
						"size": schema.Int64Attribute{
							Description: "Size of the snapshot in bytes",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *DataSourceDbSnapshots) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dbSnapshotsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	snapshots, err := r.snapshots.ListContext(ctx, dbaasSnapshots.ListParameters{
		InstanceId: data.InstanceId.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasSnapshots.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("failed to list snapshots", err.Error())
		return
	}

	var snapshotModels []dbSnapshotModel
	for _, snapshot := range snapshots.Results {
		snapshotModels = append(snapshotModels, dbSnapshotModel{
			ID:          types.StringValue(snapshot.Id),
			Name:        types.StringValue(snapshot.Name),
			Description: types.StringValue(snapshot.Description),
			CreatedAt:   types.StringValue(snapshot.CreatedAt),
			Status:      types.StringValue(snapshot.Status),
			Size:        types.Int64PointerValue(tfutil.ConvertIntPointerToInt64Pointer(&snapshot.AllocatedSize)),
		})
	}

	data.Snapshots = snapshotModels
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
