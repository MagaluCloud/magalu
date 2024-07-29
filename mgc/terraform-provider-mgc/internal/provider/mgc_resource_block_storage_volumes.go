package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"

	sdkBlockStorageVolumes "magalu.cloud/lib/products/block_storage/volumes"
	"magalu.cloud/sdk"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &bsVolumes{}
	_ resource.ResourceWithConfigure = &bsVolumes{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewBlockStorageVolumesResource() resource.Resource {
	return &bsVolumes{}
}

// orderResource is the resource implementation.
type bsVolumes struct {
	sdkClient *mgcSdk.Client
	bsVolumes sdkBlockStorageVolumes.Service
}

// Metadata returns the resource type name.
func (r *bsVolumes) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_storage_volumes"
}

// Configure adds the provider configured client to the resource.
func (r *bsVolumes) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
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

	r.bsVolumes = sdkBlockStorageVolumes.NewService(ctx, r.sdkClient)
}

// vmSnapshotsResourceModel maps de resource schema data.
type bsVolumesResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Size      types.Int64  `tfsdk:"size"`
	Type      bsVolumeType `tfsdk:"type"`
	State     types.String `tfsdk:"state"`
	Status    types.String `tfsdk:"status"`
	UpdatedAt types.String `tfsdk:"updated_at"`
	CreatedAt types.String `tfsdk:"created_at"`
}

type bsVolumeType struct {
	DiskType types.String     `tfsdk:"disk_type"`
	Id       types.String     `tfsdk:"id"`
	Iops     bsVolumeTypeIops `tfsdk:"iops"`
	Name     types.String     `tfsdk:"name"`
	Status   types.String     `tfsdk:"status"`
}

type bsVolumeTypeIops struct {
	Read  types.Int64 `tfsdk:"read"`
	Write int         `tfsdk:"write"`
}

// Schema defines the schema for the resource.
func (r *bsVolumes) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	// YOUR CODE HERE
	// description := "change here"
}

func (r *bsVolumes) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	//do nothing
}

func (r *bsVolumes) setValuesFromServer(ctx context.Context, result sdkBlockStorageVolumes.GetResult, state *bsVolumesResourceModel) {
	state.ID = types.StringValue(result.Id)
	state.Name = types.StringValue(result.Name)
	state.Size = types.Int64Value(int64(result.Size))
	state.State = types.StringValue(result.State)
	state.Status = types.StringValue(result.Status)
	state.Type = bsVolumeType{
		DiskType: types.StringPointerValue(result.Type.DiskType),
		Id:       types.StringValue(result.Type.Id),
		Name:     types.StringPointerValue(result.Type.Name),
		Status:   types.StringPointerValue(result.Type.Status),
		Iops: bsVolumeTypeIops{
			Read:  types.Int64Value(int64(result.Type.Iops.Read)),
			Write: int(result.Type.Iops.Write),
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *bsVolumes) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	plan := &vmSnapshotsResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	_, _ = r.bsVolumes.Get(sdkBlockStorageVolumes.GetParameters{
		Id: plan.ID.ValueString(),
	}, sdkBlockStorageVolumes.GetConfigs{})

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Create creates the resource and sets the initial Terraform state.
func (r *bsVolumes) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &vmSnapshotsResourceModel{}
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// YOUR CODE HERE

	plan.CreatedAt = types.StringValue(time.Now().Format(time.RFC850))
	plan.UpdatedAt = types.StringValue(time.Now().Format(time.RFC850))
	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *bsVolumes) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// YOUR CODE HERE
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bsVolumes) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data vmSnapshotsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	// UNCOMMENT THE FOLLOWING LINE TO DELETE THE RESOURCE
	// err := r.vmSnapshots.Delete(
	// 	sdkVmSnapshots.DeleteParameters{
	// 		Id: data.ID.ValueString(),
	// 	},
	// 	sdkVmSnapshots.DeleteConfigs{})
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error deleting VM Snapshot",
	// 		"Could not delete VM Snapshot "+data.ID.ValueString()+": "+err.Error(),
	// 	)
	// 	return
	// }

}
