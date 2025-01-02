package datasources

import (
	"context"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	dbaasBackups "github.com/MagaluCloud/magalu/mgc/lib/products/dbaas/instances/backups"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/tfutil"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DataSourceDbBackups struct {
	sdkClient *mgcSdk.Client
	backups   dbaasBackups.Service
}

type dbBackupsModel struct {
	InstanceId types.String    `tfsdk:"instance_id"`
	Backups    []dbBackupModel `tfsdk:"backups"`
}

func NewDataSourceDbaasInstancesBackups() datasource.DataSource {
	return &DataSourceDbBackups{}
}

func (r *DataSourceDbBackups) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_instances_backups"
}

func (r *DataSourceDbBackups) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.backups = dbaasBackups.NewService(ctx, r.sdkClient)
}

func (r *DataSourceDbBackups) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all backups for a database instance.",
		Attributes: map[string]schema.Attribute{
			"instance_id": schema.StringAttribute{
				Description: "ID of the instance",
				Required:    true,
			},
			"backups": schema.ListNestedAttribute{
				Description: "List of backups",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the backup",
							Computed:    true,
						},
						"instance_id": schema.StringAttribute{
							Description: "ID of the instance",
							Required:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the backup",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Creation timestamp",
							Computed:    true,
						},
						"status": schema.StringAttribute{
							Description: "Status of the backup",
							Computed:    true,
						},
						"size": schema.Int64Attribute{
							Description: "Size of the backup in bytes",
							Computed:    true,
						},
						"mode": schema.StringAttribute{
							Description: "Backup mode",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *DataSourceDbBackups) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dbBackupsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	backups, err := r.backups.ListContext(ctx, dbaasBackups.ListParameters{
		InstanceId: data.InstanceId.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasBackups.ListConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("failed to list backups", err.Error())
		return
	}

	var backupModels []dbBackupModel
	for _, backup := range backups.Results {
		backupModels = append(backupModels, dbBackupModel{
			ID:        types.StringValue(backup.Id),
			Name:      types.StringPointerValue(backup.Name),
			CreatedAt: types.StringValue(backup.CreatedAt),
			Status:    types.StringValue(backup.Status),
			Size:      types.Int64PointerValue(tfutil.ConvertIntPointerToInt64Pointer(backup.Size)),
			Mode:      types.StringValue(backup.Mode),
		})
	}

	data.Backups = backupModels
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
