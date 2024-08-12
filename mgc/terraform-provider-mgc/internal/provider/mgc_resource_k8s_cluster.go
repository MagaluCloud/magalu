package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkCluster "magalu.cloud/lib/products/kubernetes/cluster"
	"magalu.cloud/sdk"
)

const (
	ClusterPoolingTimeout = 100 * time.Minute
)

type KubernetesClusterCreateResourceModel struct {
	Name               types.String   `tfsdk:"name"`
	AsyncCreation      types.Bool     `tfsdk:"async_creation"`
	AllowedCidrs       []types.String `tfsdk:"allowed_cidrs"`
	Description        types.String   `tfsdk:"description"`
	EnabledServerGroup types.Bool     `tfsdk:"enabled_server_group"`
	Version            types.String   `tfsdk:"version"`
	CreatedAt          types.String   `tfsdk:"created_at"`
	ID                 types.String   `tfsdk:"id"`
	EnabledBastion     types.Bool     `tfsdk:"enabled_bastion"` // Deprecated
	Zone               types.String   `tfsdk:"zone"`
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
	nameRule := regexp.MustCompile(`^[a-z]([-a-z0-9]{0,61}[a-z0-9])?$`)
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

	out := ConvertSDKCreateResultToTerraformCreateClsuterModel(&cluster)
	out.EnabledBastion = data.EnabledBastion
	out.AsyncCreation = data.AsyncCreation
	out.EnabledServerGroup = data.EnabledServerGroup
	out.Zone = data.Zone
	diags = resp.State.Set(ctx, &out)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}
}

func (r *k8sClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KubernetesClusterCreateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	param := convertTerraformModelToSDKCreateParameters(&data)
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

	newState := ConvertSDKCreateResultToTerraformCreateClsuterModel(&createdCluster)
	newState.EnabledBastion = data.EnabledBastion
	newState.AsyncCreation = data.AsyncCreation
	newState.EnabledServerGroup = data.EnabledServerGroup
	newState.Zone = data.Zone
	diags := resp.State.Set(ctx, &newState)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}
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
	var data KubernetesClusterCreateResourceModel
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

func convertTerraformModelToSDKCreateParameters(data *KubernetesClusterCreateResourceModel) *sdkCluster.CreateParameters {
	ac := createAllowedCidrs(data)
	return &sdkCluster.CreateParameters{
		AllowedCidrs:       ac,
		Description:        data.Description.ValueStringPointer(),
		Name:               data.Name.ValueString(),
		Version:            data.Version.ValueStringPointer(),
		Zone:               data.Zone.ValueStringPointer(),
		EnabledBastion:     data.EnabledBastion.ValueBool(),
		EnabledServerGroup: data.EnabledServerGroup.ValueBoolPointer(),
	}
}

func createAllowedCidrs(data *KubernetesClusterCreateResourceModel) *sdkCluster.CreateParametersAllowedCidrs {
	allowedCidrs := []string{}
	for _, c := range data.AllowedCidrs {
		allowedCidrs = append(allowedCidrs, c.ValueString())
	}
	ac := sdkCluster.CreateParametersAllowedCidrs(allowedCidrs)

	if len(ac) == 0 {
		return nil
	}

	return &ac
}

func ConvertSDKCreateResultToTerraformCreateClsuterModel(sdkResult *sdkCluster.GetResult) *KubernetesClusterCreateResourceModel {
	if sdkResult == nil {
		return nil
	}

	tfModel := &KubernetesClusterCreateResourceModel{
		Name:      types.StringValue(sdkResult.Name),
		ID:        types.StringValue(sdkResult.Id),
		CreatedAt: types.StringValue(*sdkResult.CreatedAt),
		Version:   types.StringValue(sdkResult.Version),
	}

	if sdkResult.Description != nil {
		tfModel.Description = types.StringValue(*sdkResult.Description)
	}

	if sdkResult.AllowedCidrs != nil {
		tfModel.AllowedCidrs = convertStringSliceToTypesStringSlice(*sdkResult.AllowedCidrs)
	}

	return tfModel
}
