package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	mgcSdk "magalu.cloud/lib"
	sdkCluster "magalu.cloud/lib/products/kubernetes/cluster"
	"magalu.cloud/sdk"
)

const (
	ClusterPoolingTimeout = 100 * time.Minute
)

type KubernetesClusterCreateResourceModel struct {
	Name               types.String `tfsdk:"name"`
	AsyncCreation      types.Bool   `tfsdk:"async_creation"`
	NodePools          types.List   `tfsdk:"node_pools"` // Deprecated
	AllowedCidrs       types.List   `tfsdk:"allowed_cidrs"`
	Description        types.String `tfsdk:"description"`
	EnabledServerGroup types.Bool   `tfsdk:"enabled_server_group"`
	Version            types.String `tfsdk:"version"`
	CreatedAt          types.String `tfsdk:"created_at"`
	ID                 types.String `tfsdk:"id"`
	EnabledBastion     types.Bool   `tfsdk:"enabled_bastion"` // Deprecated
	Zone               types.String `tfsdk:"zone"`
}

type NodePoolCreatedResource struct {
	Name      types.String `tfsdk:"name"`
	Replicas  types.Int64  `tfsdk:"replicas"`
	Flavor    types.String `tfsdk:"flavor"`
	AutoScale types.Object `tfsdk:"auto_scale"`
	Tags      types.List   `tfsdk:"tags"`
	Taints    types.List   `tfsdk:"taints"`
	ID        types.String `tfsdk:"id"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

type k8sClusterResource struct {
	sdkClient  *mgcSdk.Client
	k8sCluster sdkCluster.Service
}

func NewK8sClusterResource() resource.Resource {
	return &k8sClusterResource{}
}

func (r *k8sClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_cluster"
}

func (r *k8sClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.k8sCluster = sdkCluster.NewService(ctx, r.sdkClient)
}

func (r *k8sClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	nameRule := regexp.MustCompile(`^[a-z]([-a-z0-9]*[a-z0-9])?$`)
	resp.Schema = schema.Schema{
		Description: "Kubernetes cluster resource in MGC",
		Attributes: map[string]schema.Attribute{
			"enabled_bastion": schema.BoolAttribute{
				Description:        "Enables the use of a bastion host for secure access to the cluster.",
				Optional:           true,
				DeprecationMessage: "This field is deprecated and will be removed in a future version.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"async_creation": schema.BoolAttribute{
				Description: "Enables asynchronous creation of the Kubernetes cluster.",
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Kubernetes cluster name. Must be unique within a namespace and follow naming rules.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(63),
					stringvalidator.RegexMatches(nameRule, "must contain only lowercase alphanumeric characters or '-'"),
				},
			},
			"node_pools": schema.ListNestedAttribute{
				Description:        "An array representing a set of nodes within a Kubernetes cluster.",
				Optional:           true,
				DeprecationMessage: "This field is deprecated and will be removed in a future version. Create node pools in a separately resource.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"flavor": schema.StringAttribute{
							Description: "Flavor name. Definition of the CPU, RAM, and storage capacity of the nodes.",
							Required:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the node pool. Must be unique and follow naming rules.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtMost(63),
								stringvalidator.RegexMatches(nameRule, "must contain only lowercase alphanumeric characters or '-'"),
							},
						},
						"replicas": schema.Int64Attribute{
							Description: "Number of replicas of the nodes in the node pool.",
							Required:    true,
						},
						"auto_scale": schema.SingleNestedAttribute{
							Description: "Object specifying properties for updating workload resources in the Kubernetes cluster.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Object{
								objectplanmodifier.UseStateForUnknown(),
							},
							Attributes: map[string]schema.Attribute{
								"max_replicas": schema.Int64Attribute{
									Description: "Maximum number of replicas for autoscaling.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
								"min_replicas": schema.Int64Attribute{
									Description: "Minimum number of replicas for autoscaling.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
							},
						},
						"tags": schema.ListAttribute{
							Description: "List of tags applied to the node pool.",
							Optional:    true,
							ElementType: types.StringType,
						},
						"taints": schema.ListNestedAttribute{
							Description: "Property associating a set of nodes.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.List{
								listplanmodifier.UseStateForUnknown(),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"effect": schema.StringAttribute{
										Description: "The effect of the taint on pods that do not tolerate the taint.",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									"key": schema.StringAttribute{
										Description: "Key of the taint to be applied to the node.",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									"value": schema.StringAttribute{
										Description: "Value corresponding to the taint key.",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
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
					},
				},
			},
			"allowed_cidrs": schema.ListAttribute{
				Description: "List of allowed CIDR blocks for API server access.",
				Optional:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A brief description of the Kubernetes cluster.",
				Optional:    true,
			},
			"enabled_server_group": schema.BoolAttribute{
				Description: "Enables the use of a server group with anti-affinity policy during the creation of the cluster and its node pools.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.StringAttribute{
				Description: "The native Kubernetes version of the cluster. Use the standard \"vX.Y.Z\" format.",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Creation date of the Kubernetes cluster.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Description: "Cluster's UUID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone": schema.StringAttribute{
				Description: "Identifier of the zone where the Kubernetes cluster is located.",
				Optional:    true,
			},
		},
	}
}

func (r *k8sClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KubernetesClusterCreateResourceModel
	diags := req.State.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	param := sdkCluster.GetParameters{
		ClusterId: data.ID.ValueString(),
	}
	cluster, err := r.k8sCluster.Get(param,
		GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkCluster.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get Kubernetes cluster", err.Error())
		return
	}

	out := ConvertSDKCreateResultToTerraformModel(&cluster)
	out.EnabledBastion = data.EnabledBastion
	out.AsyncCreation = data.AsyncCreation
	out.EnabledServerGroup = data.EnabledServerGroup
	out.Zone = data.Zone
	resp.State.Set(ctx, &out)
}

func (r *k8sClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KubernetesClusterCreateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var np []NodePoolCreatedResource
	diags := data.NodePools.ElementsAs(ctx, &np, false)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	createNP := []sdkCluster.CreateParametersNodePoolsItem{}
	if len(np) > 0 {
		for _, n := range np {
			var as AutoScale
			diags = n.AutoScale.As(ctx, &as, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    true,
				UnhandledUnknownAsEmpty: true,
			})
			if diags.HasError() {
				resp.Diagnostics = diags
				return
			}

			var t []Taint
			diags = n.Taints.ElementsAs(ctx, &t, false)
			if diags.HasError() {
				resp.Diagnostics = diags
				return
			}
			createNP = append(createNP, ConvertTerraformNodePoolToSDKNodePool(n, as, t))
		}
	}

	param := convertTerraformModelToSDKCreateParameters(&data, createNP)
	cluster, err := r.k8sCluster.Create(*param,
		GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkCluster.CreateConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Kubernetes cluster", err.Error())
		return
	}

	if cluster.Id == "" {
		resp.Diagnostics.AddError("Failed to create Kubernetes cluster", "ID is empty")
		return
	}

	createdCluster, err := r.GetClusterPooling(ctx, cluster.Id, data.AsyncCreation.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get Kubernetes cluster", err.Error())
		return
	}

	newState := ConvertSDKCreateResultToTerraformModel(&createdCluster)
	newState.EnabledBastion = data.EnabledBastion
	newState.AsyncCreation = data.AsyncCreation
	newState.EnabledServerGroup = data.EnabledServerGroup
	newState.Zone = data.Zone
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *k8sClusterResource) GetClusterPooling(ctx context.Context, clusterId string, isAssync bool) (sdkCluster.GetResult, error) {
	param := sdkCluster.GetParameters{
		ClusterId: clusterId,
	}

	var result sdkCluster.GetResult
	var err error
	for startTime := time.Now(); time.Since(startTime) < ClusterPoolingTimeout; {
		time.Sleep(1 * time.Minute)
		result, err = r.k8sCluster.Get(param,
			GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkCluster.GetConfigs{}))
		if err != nil {
			return sdkCluster.GetResult{}, err
		}
		state := strings.ToLower(result.Status.State)

		if state == "running" || state == "provisioned" || isAssync {
			return result, nil
		}
		if state == "failed" {
			return result, errors.New("cluster failed to provision")
		}
	}

	return result, errors.New("timeout waiting for cluster to provision")
}

func (r *k8sClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Update not supported, please delete and recreate the resource")
}

func (r *k8sClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KubernetesClusterResourceModel
	diags := req.State.Get(ctx, &data)

	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	param := sdkCluster.DeleteParameters{
		ClusterId: data.ID.ValueString(),
	}

	err := r.k8sCluster.Delete(param,
		GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkCluster.DeleteConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete Kubernetes cluster", err.Error())
		return
	}
}

func convertTerraformModelToSDKCreateParameters(data *KubernetesClusterCreateResourceModel, nd []sdkCluster.CreateParametersNodePoolsItem) *sdkCluster.CreateParameters {
	ac := createAllowedCidrs(data)

	return &sdkCluster.CreateParameters{
		NodePools:          nd,
		AllowedCidrs:       &ac,
		Description:        data.Description.ValueStringPointer(),
		Name:               data.Name.ValueString(),
		Version:            data.Version.ValueStringPointer(),
		Zone:               data.Zone.ValueStringPointer(),
		EnabledBastion:     data.EnabledBastion.ValueBool(),
		EnabledServerGroup: data.EnabledServerGroup.ValueBoolPointer(),
	}
}

func createAllowedCidrs(ctx context.Context, data *KubernetesClusterCreateResourceModel) sdkCluster.CreateParametersAllowedCidrs {
	allowedCidrs := []string{}
	data.AllowedCidrs.ElementsAs(ctx, &allowedCidrs, false)
	ac := sdkCluster.CreateParametersAllowedCidrs(allowedCidrs)
	return ac
}

func convertGetResultToTerraformModel(ctx context.Context, data *sdkCluster.GetResult) (rsult *KubernetesClusterCreateResourceModel, diags diag.Diagnostics) {
	attrTypes := map[string]attr.Type{
		"name":     types.StringType,
		"replicas": types.Int64Type,
		"flavor":   types.StringType,
		"auto_scale": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"enabled":  types.BoolType,
				"min_size": types.Int64Type,
				"max_size": types.Int64Type,
			},
		},
		"tags": types.ListType{
			ElemType: types.StringType,
		},
		"taints": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"key":    types.StringType,
					"value":  types.StringType,
					"effect": types.StringType,
				},
			},
		},
		"id":         types.StringType,
		"created_at": types.StringType,
		"updated_at": types.StringType,
	}

	np := []types.Object{}
	for _, n := range *data.NodePools {
		as, diag := convertAutoScaleToObjectValue((*sdkCluster.CreateParametersNodePoolsItemAutoScale)(&n.AutoScale))
		diags = append(diags, diag...)

		tags, diag := types.ListValueFrom(ctx, types.StringType, n.Tags)
		diags = append(diags, diag...)

		taints, diag := convertTaintsToList(ctx, n.Taints)
		types.ListValueMust(types.ObjectType{}, taints)

		types.

		attrValues := map[string]attr.Value{
			"name":       types.StringValue(n.Name),
			"replicas":   types.Int64Value(int64(n.Replicas)),
			"flavor":     types.StringValue(n.InstanceTemplate.Flavor.Name),
			"auto_scale": as,
			"tags":       tags,
			"taints": types.ListType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"key":    types.StringType,
						"value":  types.StringType,
						"effect": types.StringType,
					},
				},
			},
			"id":         types.StringType,
			"created_at": types.StringType,
			"updated_at": types.StringType,
		}

		np = append(np)
	}
}

func convertAutoScaleToObjectValue(data *sdkCluster.CreateParametersNodePoolsItemAutoScale) (types.Object, diag.Diagnostics) {
	attributeTypes := map[string]attr.Type{
		"max_replicas": types.Int64Type,
		"min_replicas": types.Int64Type,
	}

	attributeValues := map[string]attr.Value{
		"max_replicas": types.Int64Value(int64(data.MaxReplicas)),
		"min_replicas": types.Int64Value(int64(data.MinReplicas)),
	}

	obj, diag := types.ObjectValue(attributeTypes, attributeValues)
	if diag.HasError() {
		return types.ObjectNull(map[string]attr.Type{}), diag
	}

	return obj, diag
}

func convertTaintsToList(ctx context.Context, data *sdkCluster.GetResultNodePoolsItemTaints) (objs []types.Object, diag diag.Diagnostics) {
	attributeTypes := map[string]attr.Type{
		"effect": types.StringType,
		"key":    types.StringType,
		"value":  types.StringType,
	}

	for _, t := range *data {
		obj, diags := types.ObjectValue(attributeTypes, map[string]attr.Value{
			"effect": types.StringValue(t.Effect),
			"key":    types.StringValue(t.Key),
			"value":  types.StringValue(t.Value),
		})
		diag.Append(diags...)
		objs = append(objs, obj)
	}
	return
}

func convertStringArrayToListType(data *[]string) (tags []types.String) {
	for _, t := range *data {
		tags = append(tags, types.StringValue(t))
	}
	return
}
