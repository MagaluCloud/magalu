package resources

import (
	"context"
	"regexp"
	"time"

	bws "github.com/geffersonFerraz/brazilian-words-sorter"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	"magalu.cloud/terraform-provider-mgc/mgc/tfutil"

	sdkBlockStorageSnapshots "magalu.cloud/lib/products/block_storage/snapshots"
)

var (
	_ resource.Resource              = &bsSnapshots{}
	_ resource.ResourceWithConfigure = &bsSnapshots{}
)

func NewBlockStorageSnapshotsResource() resource.Resource {
	return &bsSnapshots{}
}

type bsSnapshots struct {
	sdkClient   *mgcSdk.Client
	bsSnapshots sdkBlockStorageSnapshots.Service
}

func (r *bsSnapshots) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_storage_snapshots"
}

func (r *bsSnapshots) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.bsSnapshots = sdkBlockStorageSnapshots.NewService(ctx, r.sdkClient)
}

type bsSnapshotsResourceModel struct {
	ID                types.String              `tfsdk:"id"`
	Name              types.String              `tfsdk:"name"`
	NameIsPrefix      types.Bool                `tfsdk:"name_is_prefix"`
	Description       types.String              `tfsdk:"description"`
	FinalName         types.String              `tfsdk:"final_name"`
	UpdatedAt         types.String              `tfsdk:"updated_at"`
	CreatedAt         types.String              `tfsdk:"created_at"`
	Volume            *bsSnapshotsVolumeIDModel `tfsdk:"volume"`
	State             types.String              `tfsdk:"state"`
	Status            types.String              `tfsdk:"status"`
	Size              types.Int64               `tfsdk:"size"`
	SnapshotSourceID  types.String              `tfsdk:"snapshot_source_id"`
	Type              types.String              `tfsdk:"type"`
	AvailabilityZones types.List                `tfsdk:"availability_zones"`
}

type bsSnapshotsVolumeIDModel struct {
	ID types.String `tfsdk:"id"`
}

func (r *bsSnapshots) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The block storage snapshots resource allows you to manage block storage snapshots in the Magalu Cloud.",
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
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z0-9]+(?:[-_][a-z0-9]+)*$`),
						"The name must contain only lowercase letters, numbers, underlines and hyphens. Hyphens and underlines cannot be located at the edges either.",
					),
				},
				Required: true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the volume snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
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
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description: "ID of block storage volume",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
						Required: true,
					},
				},
			},
			"snapshot_source_id": schema.StringAttribute{
				Description: "The ID of the snapshot source.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Optional: true,
			},
			"type": schema.StringAttribute{
				Description: "The type of the snapshot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("instant", "object"),
				},
				Optional: true,
			},
			"availability_zones": schema.ListAttribute{
				Description: "The availability zones of the snapshot.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *bsSnapshots) setValuesFromServer(result sdkBlockStorageSnapshots.GetResult, state *bsSnapshotsResourceModel) {
	state.ID = types.StringValue(result.Id)
	state.FinalName = types.StringValue(result.Name)
	state.State = types.StringValue(result.State)
	state.Status = types.StringValue(result.Status)
	state.Type = types.StringValue(result.Type)

	if result.AvailabilityZones != nil {
		state.AvailabilityZones = types.List{}
		for _, az := range result.AvailabilityZones {
			lv, _ := types.ListValue(types.StringType, []attr.Value{types.StringValue(az)})
			state.AvailabilityZones = lv
		}
	}
}

func (r *bsSnapshots) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	data := &bsSnapshotsResourceModel{}
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	result, err := r.bsSnapshots.GetContext(ctx, sdkBlockStorageSnapshots.GetParameters{
		Id: data.ID.ValueString()},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageSnapshots.GetConfigs{}),
	)

	if err != nil {
		resp.Diagnostics.AddError("Error reading block storage snapshot", err.Error())
		return
	}

	r.setValuesFromServer(result, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

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

	createRequest := sdkBlockStorageSnapshots.CreateParameters{
		Description: plan.Description.ValueStringPointer(),
		Name:        plan.FinalName.ValueString(),
		Type:        plan.Type.ValueStringPointer(),
	}

	if plan.Volume != nil {
		createRequest.Volume = &sdkBlockStorageSnapshots.CreateParametersVolume{
			Id: plan.Volume.ID.ValueString(),
		}
	}

	if !plan.SnapshotSourceID.IsNull() {
		createRequest.SourceSnapshot = &sdkBlockStorageSnapshots.CreateParametersSourceSnapshot{
			Id: plan.SnapshotSourceID.ValueString(),
		}
	}

	createResult, err := r.bsSnapshots.CreateContext(ctx, createRequest,
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageSnapshots.CreateConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Error creating snapshot", err.Error())
		return
	}

	state.ID = types.StringValue(createResult.Id)
	state.CreatedAt = types.StringValue(time.Now().Format(time.RFC850))
	state.UpdatedAt = types.StringValue(time.Now().Format(time.RFC850))

	getCreatedResource, err := r.bsSnapshots.GetContext(ctx, sdkBlockStorageSnapshots.GetParameters{
		Id: state.ID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageSnapshots.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Error Reading BS", err.Error())
		return
	}
	r.checkStatusIsCreating(ctx, state.ID.ValueString())
	r.setValuesFromServer(getCreatedResource, state)

	state.Size = types.Int64Value(int64(getCreatedResource.Size))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *bsSnapshots) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

func (r *bsSnapshots) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data bsSnapshotsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	err := r.bsSnapshots.DeleteContext(ctx, sdkBlockStorageSnapshots.DeleteParameters{
		Id: data.ID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkBlockStorageSnapshots.DeleteConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Error deleting VM Snapshot", err.Error())
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
