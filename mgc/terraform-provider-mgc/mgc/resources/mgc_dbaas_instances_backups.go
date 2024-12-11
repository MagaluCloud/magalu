package resources

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	dbaasBackups "magalu.cloud/lib/products/dbaas/instances/backups"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

const backupStatusTimeout = 70 * time.Minute

type DBaaSInstanceBackupStatus string

const (
	DBaaSInstanceBackupStatusPending       DBaaSInstanceBackupStatus = "PENDING"
	DBaaSInstanceBackupStatusCreating      DBaaSInstanceBackupStatus = "CREATING"
	DBaaSInstanceBackupStatusCreated       DBaaSInstanceBackupStatus = "CREATED"
	DBaaSInstanceBackupStatusError         DBaaSInstanceBackupStatus = "ERROR"
	DBaaSInstanceBackupStatusDeleting      DBaaSInstanceBackupStatus = "DELETING"
	DBaaSInstanceBackupStatusDeleted       DBaaSInstanceBackupStatus = "DELETED"
	DBaaSInstanceBackupStatusErrorDeleting DBaaSInstanceBackupStatus = "ERROR_DELETING"
)

func (s DBaaSInstanceBackupStatus) String() string {
	return string(s)
}

type DBaaSInstanceBackupModel struct {
	Id         types.String `tfsdk:"id"`
	InstanceId types.String `tfsdk:"instance_id"`
	Mode       types.String `tfsdk:"mode"`
}

type DBaaSInstanceBackupResource struct {
	sdkClient     *mgcSdk.Client
	backupService dbaasBackups.Service
}

func NewDBaaSInstanceBackupResource() resource.Resource {
	return &DBaaSInstanceBackupResource{}
}

func (r *DBaaSInstanceBackupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_instances_backups"
}

func (r *DBaaSInstanceBackupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.backupService = dbaasBackups.NewService(ctx, r.sdkClient)
}

func (r *DBaaSInstanceBackupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a DBaaS instance backup",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the backup",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance_id": schema.StringAttribute{
				Description: "ID of the DBaaS instance to backup",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mode": schema.StringAttribute{
				Description: "Mode of the backup",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("FULL", "INCREMENTAL"),
				},
			},
		},
	}
}

func (r *DBaaSInstanceBackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DBaaSInstanceBackupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.backupService.CreateContext(ctx, dbaasBackups.CreateParameters{
		InstanceId: data.InstanceId.ValueString(),
		Mode:       data.Mode.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasBackups.CreateConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create backup", err.Error())
		return
	}

	data.Id = types.StringValue(created.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	err = r.waitUntilBackupStatusMatches(ctx, data.InstanceId.ValueString(), created.Id, DBaaSInstanceBackupStatusCreated)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create backup", err.Error())
		return
	}
}

func (r *DBaaSInstanceBackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DBaaSInstanceBackupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	backup, err := r.backupService.GetContext(ctx, dbaasBackups.GetParameters{
		InstanceId: data.InstanceId.ValueString(),
		BackupId:   data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasBackups.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read backup", err.Error())
		return
	}

	data.Mode = types.StringValue(backup.Mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DBaaSInstanceBackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"DBaaS instance backups cannot be updated after creation",
	)
}

func (r *DBaaSInstanceBackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DBaaSInstanceBackupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.backupService.DeleteContext(ctx, dbaasBackups.DeleteParameters{
		InstanceId: data.InstanceId.ValueString(),
		BackupId:   data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasBackups.DeleteConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete backup", err.Error())
		return
	}
}

func (r *DBaaSInstanceBackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError("Invalid import format", "Format should be: instance_id,backup_id")
		return
	}
	data := DBaaSInstanceBackupModel{
		InstanceId: types.StringValue(idParts[0]),
		Id:         types.StringValue(idParts[1]),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DBaaSInstanceBackupResource) waitUntilBackupStatusMatches(ctx context.Context, instanceID string, backupID string, status DBaaSInstanceBackupStatus) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, backupStatusTimeout)
	defer cancel()

	for {
		select {
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout waiting for backup %s to reach status %s", backupID, status)
		case <-time.After(10 * time.Second):
			backup, err := r.backupService.GetContext(ctx, dbaasBackups.GetParameters{
				InstanceId: instanceID,
				BackupId:   backupID,
			}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasBackups.GetConfigs{}))
			if err != nil {
				return err
			}

			currentStatus := DBaaSInstanceBackupStatus(backup.Status)
			if currentStatus == status {
				return nil
			}
			if currentStatus == DBaaSInstanceBackupStatusError || currentStatus == DBaaSInstanceBackupStatusErrorDeleting {
				return fmt.Errorf("backup %s is in error state", backupID)
			}
		}
	}
}
