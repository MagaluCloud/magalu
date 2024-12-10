package resources

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	dbaasSnapshots "magalu.cloud/lib/products/dbaas/instances/snapshots"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

const snapshotStatusTimeout = 30 * time.Minute

type DBaaSInstanceSnapshotStatus string

const (
	DBaaSInstanceSnapshotStatusPending       DBaaSInstanceSnapshotStatus = "PENDING"
	DBaaSInstanceSnapshotStatusCreating      DBaaSInstanceSnapshotStatus = "CREATING"
	DBaaSInstanceSnapshotStatusCreated       DBaaSInstanceSnapshotStatus = "CREATED"
	DBaaSInstanceSnapshotStatusError         DBaaSInstanceSnapshotStatus = "ERROR"
	DBaaSInstanceSnapshotStatusDeleting      DBaaSInstanceSnapshotStatus = "DELETING"
	DBaaSInstanceSnapshotStatusDeleted       DBaaSInstanceSnapshotStatus = "DELETED"
	DBaaSInstanceSnapshotStatusErrorDeleting DBaaSInstanceSnapshotStatus = "ERROR_DELETING"
)

func (s DBaaSInstanceSnapshotStatus) String() string {
	return string(s)
}

type DBaaSInstanceSnapshotModel struct {
	Id          types.String `tfsdk:"id"`
	InstanceId  types.String `tfsdk:"instance_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

type DBaaSInstanceSnapshotResource struct {
	sdkClient       *mgcSdk.Client
	snapshotService dbaasSnapshots.Service
}

func NewDBaaSInstanceSnapshotResource() resource.Resource {
	return &DBaaSInstanceSnapshotResource{}
}

func (r *DBaaSInstanceSnapshotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_instances_snapshots"
}

func (r *DBaaSInstanceSnapshotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.snapshotService = dbaasSnapshots.NewService(ctx, r.sdkClient)
}

func (r *DBaaSInstanceSnapshotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a DBaaS instance snapshot",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the snapshot",
				Computed:    true,
			},
			"instance_id": schema.StringAttribute{
				Description: "ID of the DBaaS instance to snapshot",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the snapshot",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the snapshot",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *DBaaSInstanceSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DBaaSInstanceSnapshotModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.snapshotService.CreateContext(ctx, dbaasSnapshots.CreateParameters{
		InstanceId:  data.InstanceId.ValueString(),
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasSnapshots.CreateConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create snapshot", err.Error())
		return
	}

	data.Id = types.StringValue(created.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	err = r.waitUntilSnapshotStatusMatches(ctx, data.InstanceId.ValueString(), created.Id, DBaaSInstanceSnapshotStatusCreated)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create snapshot", err.Error())
		return
	}
}

func (r *DBaaSInstanceSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DBaaSInstanceSnapshotModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	snapshot, err := r.snapshotService.GetContext(ctx, dbaasSnapshots.GetParameters{
		InstanceId: data.InstanceId.ValueString(),
		SnapshotId: data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasSnapshots.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read snapshot", err.Error())
		return
	}

	data.Name = types.StringValue(snapshot.Name)
	data.Description = types.StringValue(snapshot.Description)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DBaaSInstanceSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"DBaaS instance snapshots cannot be updated after creation",
	)
}

func (r *DBaaSInstanceSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DBaaSInstanceSnapshotModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.snapshotService.DeleteContext(ctx, dbaasSnapshots.DeleteParameters{
		InstanceId: data.InstanceId.ValueString(),
		SnapshotId: data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasSnapshots.DeleteConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete snapshot", err.Error())
		return
	}
}

func (r *DBaaSInstanceSnapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.Split(req.ID, ",")
	if len(ids) != 2 {
		resp.Diagnostics.AddError(
			"Invalid import format",
			"Format should be: instance_id,snapshot_id",
		)
		return
	}
	data := DBaaSInstanceSnapshotModel{
		InstanceId: types.StringValue(ids[0]),
		Id:         types.StringValue(ids[1]),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DBaaSInstanceSnapshotResource) waitUntilSnapshotStatusMatches(ctx context.Context, instanceID string, snapshotID string, status DBaaSInstanceSnapshotStatus) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, snapshotStatusTimeout)
	defer cancel()

	for {
		select {
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout waiting for snapshot %s to reach status %s", snapshotID, status)
		case <-time.After(10 * time.Second):
			snapshot, err := r.snapshotService.GetContext(ctx, dbaasSnapshots.GetParameters{
				InstanceId: instanceID,
				SnapshotId: snapshotID,
			}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasSnapshots.GetConfigs{}))
			if err != nil {
				return err
			}

			currentStatus := DBaaSInstanceSnapshotStatus(snapshot.Status)
			if currentStatus == status {
				return nil
			}
			if currentStatus == DBaaSInstanceSnapshotStatusError {
				return fmt.Errorf("snapshot %s is in error state", snapshotID)
			}
		}
	}
}
