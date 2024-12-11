package resources

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	dbaasEngines "magalu.cloud/lib/products/dbaas/engines"
	dbaasInstanceTypes "magalu.cloud/lib/products/dbaas/instance_types"
	dbaasInstances "magalu.cloud/lib/products/dbaas/instances"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

const instanceStatusTimeout = 60 * time.Minute

type DBaaSInstanceStatus string

const (
	DBaaSInstanceStatusCreating      DBaaSInstanceStatus = "CREATING"
	DBaaSInstanceStatusError         DBaaSInstanceStatus = "ERROR"
	DBaaSInstanceStatusStopped       DBaaSInstanceStatus = "STOPPED"
	DBaaSInstanceStatusReboot        DBaaSInstanceStatus = "REBOOT"
	DBaaSInstanceStatusPending       DBaaSInstanceStatus = "PENDING"
	DBaaSInstanceStatusResizing      DBaaSInstanceStatus = "RESIZING"
	DBaaSInstanceStatusDeleted       DBaaSInstanceStatus = "DELETED"
	DBaaSInstanceStatusActive        DBaaSInstanceStatus = "ACTIVE"
	DBaaSInstanceStatusStarting      DBaaSInstanceStatus = "STARTING"
	DBaaSInstanceStatusStopping      DBaaSInstanceStatus = "STOPPING"
	DBaaSInstanceStatusBackingUp     DBaaSInstanceStatus = "BACKING_UP"
	DBaaSInstanceStatusDeleting      DBaaSInstanceStatus = "DELETING"
	DBaaSInstanceStatusRestoring     DBaaSInstanceStatus = "RESTORING"
	DBaaSInstanceStatusErrorDeleting DBaaSInstanceStatus = "ERROR_DELETING"
	DBaaSInstanceStatusMaintenance   DBaaSInstanceStatus = "MAINTENANCE"
)

func (s DBaaSInstanceStatus) String() string {
	return string(s)
}

func (s DBaaSInstanceStatus) IsError() bool {
	return strings.Contains(string(s), "ERROR")
}

type DBaaSInstanceModel struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	User                types.String `tfsdk:"user"`
	Password            types.String `tfsdk:"password"`
	EngineName          types.String `tfsdk:"engine_name"`
	EngineVersion       types.String `tfsdk:"engine_version"`
	InstanceType        types.String `tfsdk:"instance_type"`
	VolumeSize          types.Int64  `tfsdk:"volume_size"`
	BackupRetentionDays types.Int64  `tfsdk:"backup_retention_days"`
	BackupStartAt       types.String `tfsdk:"backup_start_at"`
}

type DBaaSInstanceResource struct {
	sdkClient          *mgcSdk.Client
	dbaasInstances     dbaasInstances.Service
	dbaasEngines       dbaasEngines.Service
	dbaasInstanceTypes dbaasInstanceTypes.Service
}

func NewDBaaSInstanceResource() resource.Resource {
	return &DBaaSInstanceResource{}
}

func (r *DBaaSInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dbaas_instances"
}

func (r *DBaaSInstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.dbaasInstances = dbaasInstances.NewService(ctx, r.sdkClient)
	r.dbaasEngines = dbaasEngines.NewService(ctx, r.sdkClient)
	r.dbaasInstanceTypes = dbaasInstanceTypes.NewService(ctx, r.sdkClient)
}

func (r *DBaaSInstanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a DBaaS (Database-as-a-Service) instance",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Unique identifier for the DBaaS instance. Generated automatically on creation.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the DBaaS instance. Must be unique within your account. Cannot be changed after creation.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(100),
				},
			},
			"user": schema.StringAttribute{
				Description: "Master username for the database. Must start with a letter and contain only alphanumeric characters.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(25),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`),
						"must start with a letter and contain only alphanumeric characters",
					),
				},
				Required: true,
			},
			"password": schema.StringAttribute{
				Description: "Master password for the database. Must be at least 8 characters long and contain letters, numbers and special characters.",
				Required:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(8),
					stringvalidator.LengthAtMost(50),
				},
			},
			"engine_name": schema.StringAttribute{
				Description: "Type of database engine to use (e.g., 'mysql', 'postgresql'). Cannot be changed after creation.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"engine_version": schema.StringAttribute{
				Description: "Version of the database engine (e.g., '8.0', '13.3'). Must be compatible with the selected engine_name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"instance_type": schema.StringAttribute{
				Description: "Compute and memory capacity of the instance (e.g., 'db.t3.micro'). Can be changed to scale the instance.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"volume_size": schema.Int64Attribute{
				Description: "Size of the storage volume in GB. Can be increased but not decreased after creation.",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.Between(10, 50000),
				},
			},
			"backup_retention_days": schema.Int64Attribute{
				Description: "Number of days to retain automated backups (1-35 days). Zero disables automated backups. Default is 7 days.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 35),
				},
			},
			"backup_start_at": schema.StringAttribute{
				Description: "Time to initiate the daily backup in UTC (format: 'HH:MM:SS'). Default is 04:00:00.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`),
						"must be in format HH:MM:SS (24-hour)",
					),
				},
			},
		},
	}
}

func (r *DBaaSInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DBaaSInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	engineID, err := r.validateAndGetEngineID(ctx, data.EngineName.ValueString(), data.EngineVersion.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid engine name", err.Error())
		return
	}

	instanceTypeID, err := r.validateAndGetInstanceTypeID(ctx, data.InstanceType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid instance type", err.Error())
		return
	}

	params := dbaasInstances.CreateParameters{
		Name:           data.Name.ValueString(),
		User:           data.User.ValueString(),
		Password:       data.Password.ValueString(),
		EngineId:       &engineID,
		InstanceTypeId: &instanceTypeID,
		Volume: dbaasInstances.CreateParametersVolume{
			Size: *tfutil.ConvertInt64PointerToIntPointer(data.VolumeSize.ValueInt64Pointer()),
		},
		BackupRetentionDays: tfutil.ConvertInt64PointerToIntPointer(data.BackupRetentionDays.ValueInt64Pointer()),
		BackupStartAt:       data.BackupStartAt.ValueStringPointer(),
	}

	created, err := r.dbaasInstances.CreateContext(ctx, params,
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasInstances.CreateConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create DBaaS Instance", err.Error())
		return
	}

	data.Id = types.StringValue(created.Id)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	err = r.waitUntilInstanceStatusMatches(ctx, data.Id.ValueString(), DBaaSInstanceStatusActive.String())
	if err != nil {
		resp.Diagnostics.AddError("Failed to create DBaaS Instance", err.Error())
		return
	}
}

func (r *DBaaSInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DBaaSInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instance, err := r.dbaasInstances.GetContext(ctx, dbaasInstances.GetParameters{
		InstanceId: data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasInstances.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to read DBaaS Instance", err.Error())
		return
	}

	engineName, engineVersion, err := r.getEngineNameAndVersionByID(ctx, instance.EngineId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get engine name and version", err.Error())
		return
	}

	instanceTypeName, err := r.getInstanceTypeNameByID(ctx, instance.InstanceTypeId)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get instance type name", err.Error())
		return
	}

	data.Name = types.StringValue(instance.Name)
	data.EngineName = types.StringValue(engineName)
	data.EngineVersion = types.StringValue(engineVersion)
	data.InstanceType = types.StringValue(instanceTypeName)
	data.VolumeSize = types.Int64Value(int64(instance.Volume.Size))
	data.BackupRetentionDays = types.Int64Value(int64(instance.BackupRetentionDays))
	data.BackupStartAt = types.StringValue(instance.BackupStartAt)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DBaaSInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData DBaaSInstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var currentData DBaaSInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &currentData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if planData.InstanceType.ValueString() != currentData.InstanceType.ValueString() {
		instanceTypeID, err := r.validateAndGetInstanceTypeID(ctx, planData.InstanceType.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid instance type", err.Error())
			return
		}

		_, err = r.dbaasInstances.ResizeContext(ctx, dbaasInstances.ResizeParameters{
			InstanceId:     currentData.Id.ValueString(),
			InstanceTypeId: &instanceTypeID,
		}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasInstances.ResizeConfigs{}))
		if err != nil {
			resp.Diagnostics.AddError("Failed to update DBaaS Instance", err.Error())
			return
		}
	}

	if planData.VolumeSize.ValueInt64() != currentData.VolumeSize.ValueInt64() {
		_, err := r.dbaasInstances.ResizeContext(ctx, dbaasInstances.ResizeParameters{
			InstanceId: currentData.Id.ValueString(),
			Volume: &dbaasInstances.ResizeParametersVolume{
				Size: *tfutil.ConvertInt64PointerToIntPointer(planData.VolumeSize.ValueInt64Pointer()),
			},
		}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasInstances.ResizeConfigs{}))
		if err != nil {
			resp.Diagnostics.AddError("Failed to update DBaaS Instance", err.Error())
			return
		}
	}

	if (planData.BackupRetentionDays.ValueInt64() != currentData.BackupRetentionDays.ValueInt64()) || (planData.BackupStartAt.ValueString() != currentData.BackupStartAt.ValueString()) {
		params := dbaasInstances.UpdateParameters{
			InstanceId:          planData.Id.ValueString(),
			BackupRetentionDays: tfutil.ConvertInt64PointerToIntPointer(planData.BackupRetentionDays.ValueInt64Pointer()),
			BackupStartAt:       planData.BackupStartAt.ValueStringPointer(),
		}
		_, err := r.dbaasInstances.UpdateContext(ctx, params,
			tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasInstances.UpdateConfigs{}))
		if err != nil {
			resp.Diagnostics.AddError("Failed to update DBaaS Instance", err.Error())
			return
		}
	}

	err := r.waitUntilInstanceStatusMatches(ctx, planData.Id.ValueString(), DBaaSInstanceStatusActive.String())
	if err != nil {
		resp.Diagnostics.AddError("Failed to update DBaaS Instance", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *DBaaSInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DBaaSInstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.dbaasInstances.DeleteContext(ctx, dbaasInstances.DeleteParameters{
		InstanceId: data.Id.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasInstances.DeleteConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete DBaaS Instance", err.Error())
		return
	}
}

func (r *DBaaSInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	data := DBaaSInstanceModel{}
	data.Id = types.StringValue(req.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DBaaSInstanceResource) validateAndGetEngineID(ctx context.Context, engineName string, engineVersion string) (string, error) {
	status := DBaaSInstanceStatusActive.String()
	engines, err := r.dbaasEngines.ListContext(ctx, dbaasEngines.ListParameters{
		Status: &status,
	},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasEngines.ListConfigs{}))
	if err != nil {
		return "", err
	}
	for _, engine := range engines.Results {
		if engine.Name == engineName && engine.Version == engineVersion {
			return engine.Id, nil
		}
	}
	return "", errors.New("engine not found")
}

func (r *DBaaSInstanceResource) validateAndGetInstanceTypeID(ctx context.Context, instanceType string) (string, error) {
	status := "ACTIVE"
	instanceTypes, err := r.dbaasInstanceTypes.ListContext(ctx, dbaasInstanceTypes.ListParameters{
		Status: &status,
	},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasInstanceTypes.ListConfigs{}))
	if err != nil {
		return "", err
	}
	for _, instance := range instanceTypes.Results {
		if instance.Label == instanceType {
			return instance.Id, nil
		}
	}
	return "", errors.New("instance type not found")
}

func (r *DBaaSInstanceResource) getEngineNameAndVersionByID(ctx context.Context, engineID string) (name string, version string, err error) {
	engine, err := r.dbaasEngines.GetContext(ctx, dbaasEngines.GetParameters{
		EngineId: engineID,
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasEngines.GetConfigs{}))
	if err != nil {
		return "", "", err
	}
	return engine.Name, engine.Version, nil
}

func (r *DBaaSInstanceResource) getInstanceTypeNameByID(ctx context.Context, instanceTypeID string) (string, error) {
	instanceType, err := r.dbaasInstanceTypes.GetContext(ctx, dbaasInstanceTypes.GetParameters{
		InstanceTypeId: instanceTypeID,
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasInstanceTypes.GetConfigs{}))
	if err != nil {
		return "", err
	}
	return instanceType.Label, nil
}

func (r *DBaaSInstanceResource) waitUntilInstanceStatusMatches(ctx context.Context, instanceID string, status string) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, instanceStatusTimeout)
	defer cancel()
	for {
		select {
		case <-timeoutCtx.Done():
			return fmt.Errorf("timeout waiting for instance %s to reach status %s", instanceID, status)
		case <-time.After(10 * time.Second):
			instance, err := r.dbaasInstances.GetContext(ctx, dbaasInstances.GetParameters{
				InstanceId: instanceID,
			}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, dbaasInstances.GetConfigs{}))
			if err != nil {
				return err
			}
			currentStatus := DBaaSInstanceStatus(instance.Status)
			if currentStatus.String() == status {
				return nil
			}
			if currentStatus.IsError() {
				return fmt.Errorf("instance %s is in error state", instanceID)
			}
		}
	}
}
