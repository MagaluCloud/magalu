package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"

	sdkVolumes "magalu.cloud/lib/products/block_storage/volumes"
	"magalu.cloud/sdk"
)

const (
	AttachVolumeTimeout         = 5 * time.Minute
	AttachVolumeCompletedStatus = "completed"
)

type VolumeAttach struct {
	sdkClient           *mgcSdk.Client
	blockStorageVolumes sdkVolumes.Service
}

type VolumeAttachResourceModel struct {
	BlockStorageID   types.String `tfsdk:"block_storage_id"`
	VirtualMachineID types.String `tfsdk:"virtual_machine_id"`
}

func NewVolumeAttachResource() resource.Resource {
	return &VolumeAttach{}
}

func (r *VolumeAttach) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_storage_volume_attachment"
}

func (r *VolumeAttach) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {

		return
	}

	sdk, ok := req.ProviderData.(*sdk.Sdk)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected provider config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.sdkClient = mgcSdk.NewClient(sdk)

	r.blockStorageVolumes = sdkVolumes.NewService(ctx, r.sdkClient)
}

func (r *VolumeAttach) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A block storage volume attachment.",
		Attributes: map[string]schema.Attribute{
			"block_storage_id": schema.StringAttribute{
				Description: "The ID of the block storage volume to attach.",
				Required:    true,
			},
			"virtual_machine_id": schema.StringAttribute{
				Description: "The ID of the virtual machine to attach the volume to.",
				Required:    true,
			},
		},
	}
}

func (r *VolumeAttach) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model VolumeAttachResourceModel
	diags := req.Plan.Get(ctx, &model)

	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	err := r.blockStorageVolumes.Attach(sdkVolumes.AttachParameters{
		Id:               model.BlockStorageID.ValueString(),
		VirtualMachineId: model.VirtualMachineID.ValueString(),
	}, GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVolumes.AttachConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to attach volume", err.Error())
		return
	}

	err = r.waitForVolumeAvailability(model.BlockStorageID.ValueString(), AttachVolumeCompletedStatus)

	if err != nil {
		resp.Diagnostics.AddError("Failed to attach volume in pooling", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

func (r *VolumeAttach) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var model VolumeAttachResourceModel
	diags := req.State.Get(ctx, &model)

	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	expand := sdkVolumes.GetParametersExpand{"attachment"}

	result, err := r.blockStorageVolumes.Get(sdkVolumes.GetParameters{
		Id:     model.BlockStorageID.ValueString(),
		Expand: &expand,
	}, GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVolumes.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to get volume", err.Error())
		return
	}

	if result.Attachment == nil {
		resp.Diagnostics.AddWarning("Volume is not attached to any virtual machine", "")
		model.VirtualMachineID = types.StringValue("")
	} else {
		model.VirtualMachineID = types.StringValue(result.Attachment.Instance.Id)
	}
	model.BlockStorageID = types.StringValue(result.Id)

	resp.State.Set(ctx, &model)
}

func (r *VolumeAttach) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var model VolumeAttachResourceModel
	diags := req.Plan.Get(ctx, &model)

	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	err := r.blockStorageVolumes.Attach(sdkVolumes.AttachParameters{
		Id:               model.BlockStorageID.ValueString(),
		VirtualMachineId: model.VirtualMachineID.ValueString(),
	}, GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVolumes.AttachConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to attach volume", err.Error())
		return
	}

	err = r.waitForVolumeAvailability(model.BlockStorageID.ValueString(), AttachVolumeCompletedStatus)

	if err != nil {
		resp.Diagnostics.AddError("Failed to attach volume in pooling", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

func (r *VolumeAttach) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model VolumeAttachResourceModel
	diags := req.State.Get(ctx, &model)

	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	err := r.blockStorageVolumes.Detach(sdkVolumes.DetachParameters{
		Id: model.BlockStorageID.ValueString(),
	}, GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVolumes.DetachConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to detach volume", err.Error())
		return
	}

	err = r.waitForVolumeAvailability(model.BlockStorageID.ValueString(), AttachVolumeCompletedStatus)

	if err != nil {
		resp.Diagnostics.AddError("Failed to detach volume in pooling", err.Error())
		return
	}
}

func (r *VolumeAttach) waitForVolumeAvailability(volumeID string, expetedStatus string) (err error) {
	for startTime := time.Now(); time.Since(startTime) < AttachVolumeTimeout; {
		time.Sleep(10 * time.Second)
		getResult, err := r.blockStorageVolumes.Get(sdkVolumes.GetParameters{
			Id: volumeID,
		}, GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkVolumes.GetConfigs{}))
		if err != nil {
			return err
		}
		if getResult.Status == expetedStatus {
			break
		}
	}
	return nil
}
