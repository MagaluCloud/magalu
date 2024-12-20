package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	dbaasBackups "magalu.cloud/lib/products/dbaas/instances/backups"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type DataSourceDbBackup struct {
	sdkClient *mgcSdk.Client
	backups   dbaasBackups.Service
}

type dbBackupModel struct {
	ID         types.String `tfsdk:"id"`
	InstanceId types.String `tfsdk:"instance_id"`
	Name       types.String `tfsdk:"name"`
	CreatedAt  types.String `tfsdk:"created_at"`
	Status     types.String `tfsdk:"status"`
	Size       types.Int64  `tfsdk:"size"`
	Mode       types.String `tfsdk:"mode"`
}

func NewDataSourceDbaasInstancesBackup() datasource.DataSource {
	return &DataSourceDbBackup{}
}

func (r *DataSourceDbBackup) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_instances_backup"
}

func (r *DataSourceDbBackup) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *DataSourceDbBackup) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get a database backup by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "ID of the backup",
				Required:    true,
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
				Description: "Backup mode (FULL or INCREMENTAL)",
				Computed:    true,
			},
		},
	}
}

func (r *DataSourceDbBackup) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dbBackupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	backup, err := r.backups.GetContext(ctx, dbaasBackups.GetParameters{
		InstanceId: data.InstanceId.ValueString(),
		BackupId:   data.ID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasBackups.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("failed to get backup", err.Error())
		return
	}

	data.Name = types.StringPointerValue(backup.Name)
	data.CreatedAt = types.StringValue(backup.CreatedAt)
	data.Status = types.StringValue(backup.Status)
	data.Size = types.Int64PointerValue(tfutil.ConvertIntPointerToInt64Pointer(backup.Size))
	data.Mode = types.StringValue(backup.Mode)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
