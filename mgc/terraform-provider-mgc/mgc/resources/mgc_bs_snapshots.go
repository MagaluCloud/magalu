package resources

import (
	"context"
	"fmt"
	"time"

	bws "github.com/geffersonFerraz/brazilian-words-sorter"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"

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

	config, ok := req.ProviderData.(tfutil.ProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected provider config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	sdk := sdk.NewSdk()
	_ = sdk.Config().SetTempConfig("region", config.Region.ValueStringPointer())
	_ = sdk.Config().SetTempConfig("env", config.Env.ValueStringPointer())
	_ = sdk.Config().SetTempConfig("api_key", config.ApiKey.ValueStringPointer())
	r.sdkClient = mgcSdk.NewClient(sdk)

	r.bsSnapshots = sdkBlockStorageSnapshots.NewService(ctx, r.sdkClient)
}

// bsSnapshotsResourceModel maps de resource schema data.
type bsSnapshotsResourceModel struct {
	ID           types.String             `tfsdk:"id"`
	Name         types.String             `tfsdk:"name"`
	NameIsPrefix types.Bool               `tfsdk:"name_is_prefix"`
	Description  types.String             `tfsdk:"description"`
	FinalName    types.String             `tfsdk:"final_name"`
	UpdatedAt    types.String             `tfsdk:"updated_at"`
	CreatedAt    types.String             `tfsdk:"created_at"`
	Volume       bsSnapshotsVolumeIDModel `tfsdk:"volume"`
	State        types.String             `tfsdk:"state"`
	Status       types.String             `tfsdk:"status"`
	Size         types.Int64              `tfsdk:"size"`
}

type bsSnapshotsVolumeIDModel struct {
	ID types.String `tfsdk:"id"`
}

// Schema defines the schema for the resource.
func (r *bsSnapshots) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	description := "The block storage snapshots resource allows you to manage block storage snapshots in the Magalu Cloud."
	resp.Schema = schema.Schema{
		Description:         description,
		MarkdownDescription: description,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the volume snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"name_is_prefix": schema.BoolAttribute{
				Description: "Indicates whether the provided name is a prefix or the exact name of the volume snapshot.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"name": schema.StringAttribute{
				Description: "The name of the volume snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Required: true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the volume snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Required: true,
			},
			"final_name": schema.StringAttribute{
				Description: "The final name of the volume snapshot after applying any naming conventions or modifications.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Description: "The timestamp when the block storage was last updated.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "The timestamp when the block storage was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"state": schema.StringAttribute{
				Description: "The current state of the virtual machine instance.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "The status of the virtual machine instance.",
				Computed:    true,
			},
			"size": schema.Int64Attribute{
				Description: "The size of the snapshot in GB.",
				Computed:    true,
			},
			"volume": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "ID of block storage volume",
						Required:    true,
					},
				},
			},
		},
	}
}

func (r *bsSnapshots) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	//do nothing
}

func (r *bsSnapshots) setValuesFromServer(result sdkBlockStorageSnapshots.GetResult, state *bsSnapshotsResourceModel) {
	state.ID = types.StringValue(result.Id)
	state.FinalName = types.StringValue(result.Name)
	state.State = types.StringValue(result.State)
	state.Status = types.StringValue(result.Status)

}

// Read refreshes the Terraform state with the latest data.
func (r *bsSnapshots) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := &bsSnapshotsResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	result, err := r.bsSnapshots.GetContext(ctx, sdkBlockStorageSnapshots.GetParameters{
		Id: data.ID.ValueString()},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageSnapshots.GetConfigs{}),
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading block storage snapshot",
			"Could not read block storage snapshot "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	r.setValuesFromServer(result, data)

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

	state := plan
	state.FinalName = types.StringValue(state.Name.ValueString())

	if state.NameIsPrefix.ValueBool() {
		bwords := bws.BrazilianWords(3, "-")
		state.FinalName = types.StringValue(state.Name.ValueString() + "-" + bwords.Sort())
	}

	// Create the block storage
	createResult, err := r.bsSnapshots.CreateContext(ctx, sdkBlockStorageSnapshots.CreateParameters{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.FinalName.String(),
		Volume: sdkBlockStorageSnapshots.CreateParametersVolume{
			Id: plan.Volume.ID.ValueString(),
		},
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageSnapshots.CreateConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not create virtual-machine, unexpected error: "+err.Error(),
		)
		return
	}

	state.ID = types.StringValue(createResult.Id)
	state.CreatedAt = types.StringValue(time.Now().Format(time.RFC850))
	state.UpdatedAt = types.StringValue(time.Now().Format(time.RFC850))

	getCreatedResource, err := r.bsSnapshots.GetContext(ctx, sdkBlockStorageSnapshots.GetParameters{
		Id: state.ID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageSnapshots.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading BS",
			"Could not read BS ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	r.checkStatusIsCreating(ctx, state.ID.ValueString())

	r.setValuesFromServer(getCreatedResource, state)

	state.Size = types.Int64Value(int64(getCreatedResource.Size))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *bsSnapshots) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bsSnapshots) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bsSnapshotsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	err := r.bsSnapshots.DeleteContext(ctx, sdkBlockStorageSnapshots.DeleteParameters{
		Id: data.ID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageSnapshots.DeleteConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting VM Snapshot",
			"Could not delete VM Snapshot "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

}

func (r *bsSnapshots) checkStatusIsCreating(ctx context.Context, id string) {
	getResult := &sdkBlockStorageSnapshots.GetResult{}

	duration := 5 * time.Minute
	startTime := time.Now()
	getParam := sdkBlockStorageSnapshots.GetParameters{Id: id}
	var err error
	for {
		elapsed := time.Since(startTime)
		remaining := duration - elapsed
		if remaining <= 0 {
			if getResult.Status != "" {
				return
			}
			return
		}

		*getResult, err = r.bsSnapshots.GetContext(ctx, getParam, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageSnapshots.GetConfigs{}))
		if err != nil {
			return
		}
		if getResult.State == "available" {
			return
		}

		time.Sleep(3 * time.Second)
	}
}
