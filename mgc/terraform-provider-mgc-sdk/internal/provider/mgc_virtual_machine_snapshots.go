package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"

	"magalu.cloud/sdk"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &vmSnapshots{}
	_ resource.ResourceWithConfigure = &vmSnapshots{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewVirtualMachineSnapshotsResource() resource.Resource {
	return &vmSnapshots{}
}

// orderResource is the resource implementation.
type vmSnapshots struct {
	sdkClient *mgcSdk.Client
}

// Metadata returns the resource type name.
func (r *vmSnapshots) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_virtual_machine_snapshots"
}

// Configure adds the provider configured client to the resource.
func (r *vmSnapshots) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
}

// vmSnapshotsResourceModel maps de resource schema data.
type vmSnapshotsResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Schema defines the schema for the resource.
func (r *vmSnapshots) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
}

func (r *vmSnapshots) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	//do nothing
}

// Read refreshes the Terraform state with the latest data.
func (r *vmSnapshots) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Create creates the resource and sets the initial Terraform state.
func (r *vmSnapshots) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *vmSnapshots) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vmSnapshots) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}
