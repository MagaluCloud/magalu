package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"

	sdkVmSnapshots "magalu.cloud/lib/products/virtual_machine/snapshots"
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
	sdkClient   *mgcSdk.Client
	vmSnapshots sdkVmSnapshots.Service
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

	r.vmSnapshots = sdkVmSnapshots.NewService(ctx, r.sdkClient)
}

// vmSnapshotsResourceModel maps de resource schema data.
type vmSnapshotsResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	VirtualMachineID types.String `tfsdk:"virtual_machine_id"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
	CreatedAt        types.String `tfsdk:"created_at"`
}

// Schema defines the schema for the resource.
func (r *vmSnapshots) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	description := "Operations with snapshots for instances."
	resp.Schema = schema.Schema{
		Description:         description,
		MarkdownDescription: description,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"name": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Required: true,
			},
			"virtual_machine_id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The id of the virtual machine.",
				Required:            true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
		},
	}
}

func (r *vmSnapshots) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	//do nothing
}

func (r *vmSnapshots) getVmSnapshot(id string) (sdkVmSnapshots.GetResult, error) {
	getResult, err := r.vmSnapshots.Get(
		sdkVmSnapshots.GetParameters{
			Id: id,
		},
		sdkVmSnapshots.GetConfigs{})
	if err != nil {
		return sdkVmSnapshots.GetResult{}, err
	}
	return getResult, nil
}

// Read refreshes the Terraform state with the latest data.
func (r *vmSnapshots) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := &vmSnapshotsResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	getResult, err := r.getVmSnapshot(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading VM",
			"Could not read VM ID "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	data.ID = types.StringValue(getResult.Id)
	data.Name = types.StringValue(*getResult.Name)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Create creates the resource and sets the initial Terraform state.
func (r *vmSnapshots) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &vmSnapshotsResourceModel{}
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createParams := sdkVmSnapshots.CreateParameters{
		Name: plan.Name.ValueString(),
		VirtualMachine: sdkVmSnapshots.CreateParametersVirtualMachine{
			Id: plan.VirtualMachineID.ValueString(),
		},
	}

	result, err := r.vmSnapshots.Create(createParams, sdkVmSnapshots.CreateConfigs{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating VM Snapshot",
			"Could not create VM Snapshot: "+err.Error(),
		)
	}

	plan.Name = types.StringValue(plan.Name.ValueString())
	plan.ID = types.StringValue(result.Id)

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
func (r *vmSnapshots) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *vmSnapshots) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data vmSnapshotsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	err := r.vmSnapshots.Delete(
		sdkVmSnapshots.DeleteParameters{
			Id: data.ID.ValueString(),
		},
		sdkVmSnapshots.DeleteConfigs{})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting VM Snapshot",
			"Could not delete VM Snapshot "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

}
