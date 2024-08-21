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
	"magalu.cloud/terraform-provider-mgc/internal/tfutil"

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
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	NameIsPrefix types.Bool   `tfsdk:"name_is_prefix"`
	FinalName    types.String `tfsdk:"final_name"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
	CreatedAt    types.String `tfsdk:"created_at"`
	Size         types.Int64  `tfsdk:"size"`
	Type         bsVolumeType `tfsdk:"type"`
	State        types.String `tfsdk:"state"`
	Status       types.String `tfsdk:"status"`
}

type bsVolumeType struct {
	DiskType types.String `tfsdk:"disk_type"`
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Status   types.String `tfsdk:"status"`
}

// Schema defines the schema for the resource.
func (r *bsVolumes) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	description := "Block storage volumes are storage devices that can be attached to virtual machines. They are used to store data and can be detached and attached to other virtual machines."
	resp.Schema = schema.Schema{
		Description:         description,
		MarkdownDescription: description,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the block storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"name_is_prefix": schema.BoolAttribute{
				Description: "Indicates whether the provided name is a prefix or the exact name of the block storage.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"name": schema.StringAttribute{
				Description: "The name of the block storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Required: true,
			},
			"final_name": schema.StringAttribute{
				Description: "The final name of the block storage after applying any naming conventions or modifications.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"size": schema.Int64Attribute{
				Description: "The size of the block storage in GB.",
				Required:    true,
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
			"type": schema.SingleNestedAttribute{
				Required:    true,
				Description: "The type of the block storage.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Description: "The name of the block storage type.",
						Required:    true,
					},
					"disk_type": schema.StringAttribute{
						Description: "The disk type of the block storage.",
						Computed:    true,
					},
					"id": schema.StringAttribute{
						Description: "The unique identifier of the block storage type.",
						Computed:    true,
					},
					"status": schema.StringAttribute{
						Description: "The status of the block storage type.",
						Computed:    true,
					},
				},
			},
		},
	}

}

func (r *bsVolumes) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	//do nothing
}

func (r *bsVolumes) setValuesFromServer(result sdkBlockStorageVolumes.GetResult, state *bsVolumesResourceModel) {
	state.ID = types.StringValue(result.Id)
	state.FinalName = types.StringValue(result.Name)
	state.Size = types.Int64Value(int64(result.Size))
	state.State = types.StringValue(result.State)
	state.Status = types.StringValue(result.Status)

	state.Type = bsVolumeType{
		DiskType: types.StringPointerValue(result.Type.DiskType),
		Id:       types.StringValue(result.Type.Id),
		Name:     types.StringPointerValue(result.Type.Name),
		Status:   types.StringPointerValue(result.Type.Status),
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *bsVolumes) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	plan := &bsVolumesResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	getResult, err := r.bsVolumes.Get(sdkBlockStorageVolumes.GetParameters{
		Id:     plan.ID.ValueString(),
		Expand: &sdkBlockStorageVolumes.GetParametersExpand{"volume_type"},
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading block storage",
			"Could not read BS ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}
	plan.ID = types.StringValue(getResult.Id)
	r.setValuesFromServer(getResult, plan)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Create creates the resource and sets the initial Terraform state.
func (r *bsVolumes) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &bsVolumesResourceModel{}
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

	createResult, err := r.bsVolumes.Create(sdkBlockStorageVolumes.CreateParameters{
		Name: state.FinalName.ValueString(),
		Size: int(state.Size.ValueInt64()),
		Type: sdkBlockStorageVolumes.CreateParametersType{
			Name: state.Type.Name.ValueStringPointer(),
		},
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.CreateConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating vm",
			"Could not create virtual-machine, unexpected error: "+err.Error(),
		)
		return
	}

	state.ID = types.StringValue(createResult.Id)

	getCreatedResource, err := r.bsVolumes.Get(sdkBlockStorageVolumes.GetParameters{
		Id:     state.ID.ValueString(),
		Expand: &sdkBlockStorageVolumes.GetParametersExpand{"volume_type"},
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading BS",
			"Could not read BS ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	r.setValuesFromServer(getCreatedResource, state)

	state.CreatedAt = types.StringValue(time.Now().Format(time.RFC850))
	state.UpdatedAt = types.StringValue(time.Now().Format(time.RFC850))
	// Set state to fully populated data
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *bsVolumes) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	data := bsVolumesResourceModel{}
	currState := &bsVolumesResourceModel{}

	req.State.Get(ctx, currState)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if data.Type.Name.ValueString() != currState.Type.Name.ValueString() {
		// retype
		err := r.bsVolumes.Retype(sdkBlockStorageVolumes.RetypeParameters{
			Id: data.ID.ValueString(),
			NewType: sdkBlockStorageVolumes.RetypeParametersNewType{
				Name: data.Type.Name.ValueStringPointer(),
			},
		},
			tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.RetypeConfigs{}))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error to retype the block storage volume",
				"Could not retype, unexpected error: "+err.Error(),
			)
			return
		}
	}

	if data.Size.ValueInt64() > currState.Size.ValueInt64() {
		err := r.bsVolumes.Extend(sdkBlockStorageVolumes.ExtendParameters{
			Id:   data.ID.ValueString(),
			Size: int(data.Size.ValueInt64()),
		},
			tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.ExtendConfigs{}))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error to resize the block storage volume",
				"Could not resize, unexpected error: "+err.Error(),
			)
			return
		}
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bsVolumes) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bsVolumesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	err := r.bsVolumes.Delete(
		sdkBlockStorageVolumes.DeleteParameters{
			Id: data.ID.ValueString(),
		},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.DeleteConfigs{}),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting block storage volume",
			"Could not delete block storage volume "+data.ID.ValueString()+": "+err.Error(),
		)
		return
	}

}
