package resources

import (
	"context"
	"errors"
	"fmt"
	"time"

	bws "github.com/geffersonFerraz/brazilian-words-sorter"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"

	sdkBlockStorageVolumes "magalu.cloud/lib/products/block_storage/volumes"
)

const (
	completedBsSttus = "completed"
)

var (
	_              resource.Resource              = &bsVolumes{}
	_              resource.ResourceWithConfigure = &bsVolumes{}
	expandBsVolume                                = &sdkBlockStorageVolumes.GetParametersExpand{"volume_type"}
)

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
	r.sdkClient, err = client.NewSDKClient(req)
	if err != nil {
		resp.Diagnostics.AddError(
			err.Error(),
			fmt.Sprintf("Expected provider config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.bsVolumes = sdkBlockStorageVolumes.NewService(ctx, r.sdkClient)
}

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
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the block storage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
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

func convertToState(result sdkBlockStorageVolumes.GetResult, originalName string, originalIsPrefix bool) *bsVolumesResourceModel {
	state := bsVolumesResourceModel{}
	state.ID = types.StringValue(result.Id)
	state.FinalName = types.StringValue(result.Name)
	state.NameIsPrefix = types.BoolValue(originalIsPrefix)
	state.Name = types.StringValue(originalName)
	state.Size = types.Int64Value(int64(result.Size))
	state.State = types.StringValue(result.State)
	state.Status = types.StringValue(result.Status)

	state.Type = bsVolumeType{
		DiskType: types.StringPointerValue(result.Type.DiskType),
		Id:       types.StringValue(result.Type.Id),
		Name:     types.StringPointerValue(result.Type.Name),
		Status:   types.StringPointerValue(result.Type.Status),
	}

	state.CreatedAt = types.StringValue(result.CreatedAt)
	return &state
}

func (r *bsVolumes) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	plan := &bsVolumesResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	getResult, err := r.bsVolumes.GetContext(ctx, sdkBlockStorageVolumes.GetParameters{
		Id:     plan.ID.ValueString(),
		Expand: expandBsVolume,
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading block storage",
			err.Error(),
		)
		return
	}
	plan = convertToState(getResult, plan.Name.ValueString(), plan.NameIsPrefix.ValueBool())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *bsVolumes) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	state := &bsVolumesResourceModel{}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.FinalName = types.StringValue(state.Name.ValueString())
	if state.NameIsPrefix.ValueBool() {
		bwords := bws.BrazilianWords(3, "-")
		state.FinalName = types.StringValue(state.Name.ValueString() + "-" + bwords.Sort())
	}

	createResult, err := r.bsVolumes.CreateContext(ctx, sdkBlockStorageVolumes.CreateParameters{
		Name: state.FinalName.ValueString(),
		Size: int(state.Size.ValueInt64()),
		Type: sdkBlockStorageVolumes.CreateParametersType{
			Name: state.Type.Name.ValueStringPointer(),
		},
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.CreateConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Error creating volume", err.Error())
		return
	}

	state.ID = types.StringValue(createResult.Id)
	getCreatedResource, err := r.waitCompletedVolume(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading block storage", err.Error())
		return
	}

	state = convertToState(*getCreatedResource, state.Name.ValueString(), state.NameIsPrefix.ValueBool())
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *bsVolumes) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	newState := &bsVolumesResourceModel{}
	currState := &bsVolumesResourceModel{}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, currState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if newState.Type.Name.ValueString() != currState.Type.Name.ValueString() {
		err := r.bsVolumes.RetypeContext(ctx, sdkBlockStorageVolumes.RetypeParameters{
			Id: newState.ID.ValueString(),
			NewType: sdkBlockStorageVolumes.RetypeParametersNewType{
				Name: newState.Type.Name.ValueStringPointer(),
			},
		},
			tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.RetypeConfigs{}))
		if err != nil {
			resp.Diagnostics.AddError("Error to retype the block storage volume", err.Error())
			return
		}
		tflog.Debug(ctx, "waiting retry completion")
		_, _ = r.waitCompletedVolume(ctx, currState.ID.ValueString())
		tflog.Info(ctx, "retype performed with success")
	}

	if newState.Size.ValueInt64() > currState.Size.ValueInt64() {
		err := r.bsVolumes.ExtendContext(ctx, sdkBlockStorageVolumes.ExtendParameters{
			Id:   newState.ID.ValueString(),
			Size: int(newState.Size.ValueInt64()),
		},
			tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.ExtendConfigs{}))
		if err != nil {
			resp.Diagnostics.AddError("Error to resize the block storage volume", err.Error())
			return
		}
		tflog.Info(ctx, "resize performed with success")
	}

	tflog.Debug(ctx, "waiting volume completion")
	getResult, err := r.waitCompletedVolume(ctx, currState.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading block storage", err.Error())
		return
	}
	newState = convertToState(*getResult, currState.Name.ValueString(), currState.NameIsPrefix.ValueBool())
	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *bsVolumes) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bsVolumesResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	err := r.bsVolumes.DeleteContext(ctx,
		sdkBlockStorageVolumes.DeleteParameters{
			Id: data.ID.ValueString(),
		},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.DeleteConfigs{}),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting block storage volume", err.Error())
		return
	}
}

func (r *bsVolumes) waitCompletedVolume(ctx context.Context, id string) (*sdkBlockStorageVolumes.GetResult, error) {
	for startTime := time.Now(); time.Since(startTime) < ClusterPoolingTimeout; {
		time.Sleep(3 * time.Second)
		getResult, err := r.bsVolumes.GetContext(ctx, sdkBlockStorageVolumes.GetParameters{
			Id:     id,
			Expand: expandBsVolume,
		}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageVolumes.GetConfigs{}))
		if err != nil {
			return nil, err
		}
		if getResult.Status == completedBsSttus {
			return &getResult, nil
		}
		tflog.Debug(ctx, fmt.Sprintf("volume current status: %s", getResult.Status))
	}
	return nil, errors.New("timeout fetching block storage resource")
}
