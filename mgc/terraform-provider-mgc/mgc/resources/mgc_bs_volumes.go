package resources

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/tfutil"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdkBlockStorageVolumes "github.com/MagaluCloud/magalu/mgc/lib/products/block_storage/volumes"
)

const (
	completedBsSttus      = "completed"
	BsVolumeStatusTimeout = 60 * time.Minute
)

type VolumeStatus string

const (
	Provisioning       VolumeStatus = "provisioning"
	Creating           VolumeStatus = "creating"
	CreatingError      VolumeStatus = "creating_error"
	CreatingErrorQuota VolumeStatus = "creating_error_quota"
	Completed          VolumeStatus = "completed"
	ExtendPending      VolumeStatus = "extend_pending"
	Extending          VolumeStatus = "extending"
	ExtendError        VolumeStatus = "extend_error"
	ExtendErrorQuota   VolumeStatus = "extend_error_quota"
	AttachingPending   VolumeStatus = "attaching_pending"
	AttachingError     VolumeStatus = "attaching_error"
	Attaching          VolumeStatus = "attaching"
	DetachingPending   VolumeStatus = "detaching_pending"
	DetachingError     VolumeStatus = "detaching_error"
	Detaching          VolumeStatus = "detaching"
	RetypePending      VolumeStatus = "retype_pending"
	Retyping           VolumeStatus = "retyping"
	RetypeError        VolumeStatus = "retype_error"
	RetypeErrorQuota   VolumeStatus = "retype_error_quota"
	DeletingPending    VolumeStatus = "deleting_pending"
	Deleting           VolumeStatus = "deleting"
	Deleted            VolumeStatus = "deleted"
	DeletedError       VolumeStatus = "deleted_error"
)

func (s VolumeStatus) String() string {
	return string(s)
}

func (s VolumeStatus) isError() bool {
	return strings.HasSuffix(s.String(), "error")
}

func NewBlockStorageVolumesResource() resource.Resource {
	return &bsVolumes{}
}

type bsVolumes struct {
	sdkClient *mgcSdk.Client
	bsVolumes sdkBlockStorageVolumes.Service
}

func (r *bsVolumes) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_storage_volumes"
}

func (r *bsVolumes) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.bsVolumes = sdkBlockStorageVolumes.NewService(ctx, r.sdkClient)
}

type bsVolumesResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	SnapshotID       types.String `tfsdk:"snapshot_id"`
	AvailabilityZone types.String `tfsdk:"availability_zone"`
	CreatedAt        types.String `tfsdk:"created_at"`
	Size             types.Int64  `tfsdk:"size"`
	Type             types.String `tfsdk:"type"`
	Encrypted        types.Bool   `tfsdk:"encrypted"`
}

func (r *bsVolumes) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Block storage volumes are storage devices that can be attached to virtual machines. They are used to store data and can be detached and attached to other virtual machines.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the volume.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the volume.",
				Required:    true,
			},
			"snapshot_id": schema.StringAttribute{
				Description: "Create a volume from a snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
				Optional: true,
			},
			"availability_zone": schema.StringAttribute{
				Description: "The availability zones where the volume is available.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"size": schema.Int64Attribute{
				Description: "The size of the volume in GB.",
				Required:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the volume was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The name of the volume type.",
			},
			"encrypted": schema.BoolAttribute{
				Description: "Indicates if the volume is encrypted.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					tfutil.ReplaceIfChangeAndNotIsNotSetOnPlanBool{},
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}

}

func (r *bsVolumes) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	plan := &bsVolumesResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	getResult, err := r.bsVolumes.GetContext(ctx, sdkBlockStorageVolumes.GetParameters{
		Id:     plan.ID.ValueString(),
		Expand: &sdkBlockStorageVolumes.GetParametersExpand{"volume_type"},
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Error Reading block storage", err.Error())
		return
	}

	convertedResult := r.toTerraformModel(getResult, plan.SnapshotID.ValueStringPointer())
	resp.Diagnostics.Append(resp.State.Set(ctx, &convertedResult)...)
}

func (r *bsVolumes) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	state := &bsVolumesResourceModel{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createParam := sdkBlockStorageVolumes.CreateParameters{
		Name:      state.Name.ValueString(),
		Size:      int(state.Size.ValueInt64()),
		Encrypted: state.Encrypted.ValueBoolPointer(),
		Type: sdkBlockStorageVolumes.CreateParametersType{
			Name: state.Type.ValueStringPointer(),
		},
	}
	if !state.SnapshotID.IsNull() {
		createParam.Snapshot = &sdkBlockStorageVolumes.CreateParametersSnapshot{
			Id: state.SnapshotID.ValueStringPointer(),
		}
	}
	if !state.AvailabilityZone.IsUnknown() {
		createParam.AvailabilityZone = state.AvailabilityZone.ValueStringPointer()
	}

	createResult, err := r.bsVolumes.CreateContext(ctx, createParam,
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.CreateConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Error creating volume", err.Error())
		return
	}
	getResult, err := r.waitUntilVolumeStatusMatches(ctx, *createResult.Id, Completed)
	if err != nil {
		resp.Diagnostics.AddError("Error Reading block storage", err.Error())
		return
	}
	convertedResult := r.toTerraformModel(*getResult, state.SnapshotID.ValueStringPointer())
	resp.Diagnostics.Append(resp.State.Set(ctx, &convertedResult)...)
}

func (r *bsVolumes) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	planData := &bsVolumesResourceModel{}
	state := &bsVolumesResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if planData.Name.ValueString() != state.Name.ValueString() {
		err := r.bsVolumes.RenameContext(ctx, sdkBlockStorageVolumes.RenameParameters{
			Id:   state.ID.ValueString(),
			Name: planData.Name.ValueString(),
		}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.RenameConfigs{}))
		if err != nil {
			resp.Diagnostics.AddError("Error renaming block storage volume", err.Error())
			return
		}
		_, err = r.waitUntilVolumeStatusMatches(ctx, state.ID.ValueString(), Completed)
		if err != nil {
			resp.Diagnostics.AddError("Server side error renaming block storage volume", err.Error())
			return
		}
	}

	if planData.Type.ValueString() != state.Type.ValueString() {
		err := r.bsVolumes.RetypeContext(ctx, sdkBlockStorageVolumes.RetypeParameters{
			Id: planData.ID.ValueString(),
			NewType: sdkBlockStorageVolumes.RetypeParametersNewType{
				Name: planData.Type.ValueStringPointer(),
			},
		}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.RetypeConfigs{}))
		if err != nil {
			resp.Diagnostics.AddError("Error to retype the block storage volume", err.Error())
			return
		}
		_, err = r.waitUntilVolumeStatusMatches(ctx, state.ID.ValueString(), Completed)
		if err != nil {
			resp.Diagnostics.AddError("Server side error retyping block storage volume", err.Error())
			return
		}
	}

	if planData.Size.ValueInt64() > state.Size.ValueInt64() {
		err := r.bsVolumes.ExtendContext(ctx, sdkBlockStorageVolumes.ExtendParameters{
			Id:   planData.ID.ValueString(),
			Size: int(planData.Size.ValueInt64()),
		},
			tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.ExtendConfigs{}))
		if err != nil {
			resp.Diagnostics.AddError("Error to resize the block storage volume", err.Error())
			return
		}
		_, err = r.waitUntilVolumeStatusMatches(ctx, state.ID.ValueString(), Completed)
		if err != nil {
			resp.Diagnostics.AddError("Server side error resizing block storage volume", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *bsVolumes) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bsVolumesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.bsVolumes.DeleteContext(ctx,
		sdkBlockStorageVolumes.DeleteParameters{
			Id: data.ID.ValueString(),
		},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.DeleteConfigs{}),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting block storage volume", err.Error())
	}
}

func (r *bsVolumes) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	data := &bsVolumesResourceModel{ID: types.StringValue(req.ID)}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *bsVolumes) toTerraformModel(volume sdkBlockStorageVolumes.GetResult, snapshotId *string) bsVolumesResourceModel {
	return bsVolumesResourceModel{
		ID:               types.StringValue(volume.Id),
		Name:             types.StringValue(volume.Name),
		AvailabilityZone: types.StringValue(volume.AvailabilityZone),
		CreatedAt:        types.StringValue(volume.CreatedAt),
		Size:             types.Int64PointerValue(tfutil.ConvertIntPointerToInt64Pointer(&volume.Size)),
		Type:             types.StringPointerValue(volume.Type.Name),
		SnapshotID:       types.StringPointerValue(snapshotId),
		Encrypted:        types.BoolPointerValue(volume.Encrypted),
	}
}

func (r *bsVolumes) waitUntilVolumeStatusMatches(ctx context.Context, volumeID string, status VolumeStatus) (*sdkBlockStorageVolumes.GetResult, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, BsVolumeStatusTimeout)
	defer cancel()

	for {
		select {
		case <-timeoutCtx.Done():
			return nil, fmt.Errorf("timeout waiting for volume %s to reach status %s", volumeID, status)
		case <-time.After(10 * time.Second):
			volume, err := r.bsVolumes.GetContext(ctx, sdkBlockStorageVolumes.GetParameters{
				Expand: &sdkBlockStorageVolumes.GetParametersExpand{"volume_type"},
				Id:     volumeID,
			}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.GetConfigs{}))
			if err != nil {
				return nil, err
			}
			currentStatus := VolumeStatus(volume.Status)
			if currentStatus == status {
				return &volume, nil
			}
			if currentStatus.isError() {
				return nil, fmt.Errorf("volume %s is in error state: %s", volumeID, currentStatus)
			}
		}
	}
}
