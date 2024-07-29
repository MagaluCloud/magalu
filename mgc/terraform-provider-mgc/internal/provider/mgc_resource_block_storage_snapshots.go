package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"

	sdkBlockStorageSnapshots "magalu.cloud/lib/products/block_storage/snapshots"
	"magalu.cloud/sdk"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &bsSnapshots{}
	_ resource.ResourceWithConfigure = &bsSnapshots{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewBlockStorageSnapshotsResource() resource.Resource {
	return &bsSnapshots{}
}

// orderResource is the resource implementation.
type bsSnapshots struct {
	sdkClient   *mgcSdk.Client
	bsSnapshots sdkBlockStorageSnapshots.Service
}

// Metadata returns the resource type name.
func (r *bsSnapshots) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_storage_snapshots"
}

// Configure adds the provider configured client to the resource.
func (r *bsSnapshots) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.bsSnapshots = sdkBlockStorageSnapshots.NewService(ctx, r.sdkClient)
}

// bsSnapshotsResourceModel maps de resource schema data.
type bsSnapshotsResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	VirtualMachineID types.String `tfsdk:"virtual_machine_id"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
	CreatedAt        types.String `tfsdk:"created_at"`
}

// Schema defines the schema for the resource.
func (r *bsSnapshots) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	// YOUR CODE HERE
	// description := "change here"
}

func (r *bsSnapshots) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	//do nothing
}

// Read refreshes the Terraform state with the latest data.
func (r *bsSnapshots) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := &bsSnapshotsResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	// YOUR CODE HERE

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Create creates the resource and sets the initial Terraform state.
func (r *bsSnapshots) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &bsSnapshotsResourceModel{}
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
func (r *bsSnapshots) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// YOUR CODE HERE
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bsSnapshots) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bsSnapshotsResourceModel
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
