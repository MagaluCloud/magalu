package resources

import (
	"context"
	"fmt"
	"strings"
	"time"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/lib"
	sdkNodepool "github.com/MagaluCloud/magalu/mgc/lib/products/kubernetes/nodepool"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/client"
	tfutil "github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/tfutil"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type NewNodePoolResource struct {
	sdkClient   *mgcSdk.Client
	sdkNodepool sdkNodepool.Service
}

func NewNewNodePoolResource() resource.Resource {
	return &NewNodePoolResource{}
}

func (r *NewNodePoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_nodepool"
}

func (r *NewNodePoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.sdkNodepool = sdkNodepool.NewService(ctx, r.sdkClient)
}

func (r *NewNodePoolResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "An array representing a set of nodes within a Kubernetes cluster.",
		Attributes: map[string]schema.Attribute{
			"flavor_name": schema.StringAttribute{
				Description: "Definition of the CPU, RAM, and storage capacity of the nodes.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_id": schema.StringAttribute{
				Description: "UUID of the Kubernetes cluster.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the node pool.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"replicas": schema.Int64Attribute{
				Description: "Number of replicas of the nodes in the node pool.",
				Required:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"max_replicas": schema.Int64Attribute{
				Description: "Maximum number of replicas for autoscaling.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"min_replicas": schema.Int64Attribute{
				Description: "Minimum number of replicas for autoscaling.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"tags": schema.ListAttribute{
				Description: "List of tags applied to the node pool.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"created_at": schema.StringAttribute{
				Description: "Date of creation of the Kubernetes Node.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Date of the last change to the Kubernetes Node.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "Node pool's UUID.",
				Computed:    true,
			},
			"taints": schema.ListNestedAttribute{
				Description: "Property associating a set of nodes.",
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"effect": schema.StringAttribute{
							Description: "The effect of the taint on pods that do not tolerate the taint.",
							Optional:    true,
						},
						"key": schema.StringAttribute{
							Description: "Key of the taint to be applied to the node.",
							Optional:    true,
						},
						"value": schema.StringAttribute{
							Description: "Value corresponding to the taint key.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *NewNodePoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data tfutil.NodePoolCreate
	diags := req.State.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}
	nodepool, err := r.sdkNodepool.GetContext(ctx, sdkNodepool.GetParameters{
		ClusterId:  data.ClusterID.ValueString(),
		NodePoolId: data.ID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkNodepool.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to get node pool", err.Error())
		return
	}
	data.NodePool = tfutil.ConvertToNodePoolGet(&nodepool)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *NewNodePoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data tfutil.NodePoolCreate
	diags := req.Config.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	var tags sdkNodepool.CreateParametersTags
	if data.Tags != nil {
		tags = sdkNodepool.CreateParametersTags(*convertStringArrayTFToSliceString(data.Tags))
	}
	createParams := sdkNodepool.CreateParameters{
		ClusterId: data.ClusterID.ValueString(),
		Flavor:    data.Flavor.ValueString(),
		Name:      data.Name.ValueString(),
		Replicas:  int(data.Replicas.ValueInt64()),
		Tags:      &tags,
		Taints:    convertTaintsNP(data.Taints),
	}

	if !data.MaxReplicas.IsNull() || !data.MinReplicas.IsNull() {
		createParams.AutoScale = &sdkNodepool.CreateParametersAutoScale{}
	}

	if !data.MaxReplicas.IsNull() {
		createParams.AutoScale.MaxReplicas = tfutil.ConvertInt64PointerToIntPointer(data.MaxReplicas.ValueInt64Pointer())
	}

	if !data.MinReplicas.IsNull() {
		createParams.AutoScale.MinReplicas = tfutil.ConvertInt64PointerToIntPointer(data.MinReplicas.ValueInt64Pointer())
	}

	nodepool, err := r.sdkNodepool.CreateContext(ctx, createParams,
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkNodepool.CreateConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to create node pool", err.Error())
		return
	}

	data.NodePool = tfutil.ConvertToNodePoolCreate(&nodepool)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	err = r.waitNodePoolCreation(ctx, nodepool.Id, data.ClusterID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to wait for node pool creation", err.Error())
		return
	}
}

func (r *NewNodePoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data tfutil.NodePoolCreate
	diags := req.Config.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	var state tfutil.NodePoolCreate
	diags = req.State.Get(ctx, &state)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	repli := int(data.Replicas.ValueInt64())
	updateParam := sdkNodepool.UpdateParameters{
		ClusterId:  data.ClusterID.ValueString(),
		NodePoolId: state.ID.ValueString(),
		Replicas:   &repli,
	}

	if !data.MaxReplicas.IsNull() || !data.MinReplicas.IsNull() {
		updateParam.AutoScale = &sdkNodepool.UpdateParametersAutoScale{}
	}

	if !data.MaxReplicas.IsNull() {
		updateParam.AutoScale.MaxReplicas = tfutil.ConvertInt64PointerToIntPointer(data.MaxReplicas.ValueInt64Pointer())
	}
	if !data.MinReplicas.IsNull() {
		updateParam.AutoScale.MinReplicas = tfutil.ConvertInt64PointerToIntPointer(data.MinReplicas.ValueInt64Pointer())
	}

	nodepool, err := r.sdkNodepool.UpdateContext(ctx, updateParam,
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkNodepool.UpdateConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to update node pool", err.Error())
		return
	}
	data.NodePool = tfutil.ConvertToNodePoolUpdate(&nodepool)

	err = r.waitNodePoolCreation(ctx, data.ID.ValueString(), data.ClusterID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to wait for node pool creation", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *NewNodePoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data tfutil.NodePoolCreate
	diags := req.State.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	err := r.sdkNodepool.DeleteContext(ctx, sdkNodepool.DeleteParameters{
		ClusterId:  data.ClusterID.ValueString(),
		NodePoolId: data.ID.ValueString(),
	},
		tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkNodepool.DeleteConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to delete node pool", err.Error())
		return
	}
}

func (r *NewNodePoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var data tfutil.NodePoolCreate
	ids := strings.Split(req.ID, ",")
	if len(ids) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", "Expected format: cluster_id,node_pool_id")
		return
	}

	nodepool, err := r.sdkNodepool.GetContext(ctx, sdkNodepool.GetParameters{
		ClusterId:  ids[0],
		NodePoolId: ids[1],
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkNodepool.GetConfigs{}))

	if err != nil {
		resp.Diagnostics.AddError("Failed to get node pool", err.Error())
		return
	}

	data.ClusterID = types.StringValue(ids[0])
	data.NodePool = tfutil.ConvertToNodePoolGet(&nodepool)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func convertTaintsNP(taints *[]tfutil.Taint) *sdkNodepool.CreateParametersTaints {
	if taints == nil {
		return nil
	}
	taintsNP := make([]sdkNodepool.CreateParametersTaintsItem, len(*taints))
	for i, taint := range *taints {
		taintsNP[i] = sdkNodepool.CreateParametersTaintsItem{
			Effect: taint.Effect.ValueString(),
			Key:    taint.Key.ValueString(),
			Value:  taint.Value.ValueString(),
		}
	}
	rt := sdkNodepool.CreateParametersTaints(taintsNP)
	return &rt
}

func convertStringArrayTFToSliceString(tags *[]types.String) *[]string {
	if tags == nil {
		return nil
	}
	tagsSlice := make([]string, len(*tags))
	for i, tag := range *tags {
		tagsSlice[i] = tag.ValueString()
	}
	return &tagsSlice
}

func (r *NewNodePoolResource) waitNodePoolCreation(ctx context.Context, nodepoolid, clusterId string) error {
	for startTime := time.Now(); time.Since(startTime) < ClusterPoolingTimeout; {
		time.Sleep(30 * time.Second)

		nodepool, err := r.sdkNodepool.GetContext(ctx, sdkNodepool.GetParameters{
			ClusterId:  clusterId,
			NodePoolId: nodepoolid,
		}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkNodepool.GetConfigs{}))

		if err != nil {
			return err
		}

		if nodepool.Status.State == "Running" {
			return nil
		}

		tflog.Debug(ctx, fmt.Sprintf("Node pool %s is in state %s", nodepoolid, nodepool.Status.State))
	}
	return fmt.Errorf("timeout waiting for node pool creation")
}
