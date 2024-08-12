package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	sdkCluster "magalu.cloud/lib/products/kubernetes/cluster"
)

type KubernetesClusterResourceModel struct {
	Name          types.String   `tfsdk:"name"`
	NodePools     []NodePool     `tfsdk:"node_pools"`
	AllowedCidrs  []types.String `tfsdk:"allowed_cidrs"`
	Description   types.String   `tfsdk:"description"`
	Version       types.String   `tfsdk:"version"`
	Zone          types.String   `tfsdk:"zone"`
	Addons        *Addons        `tfsdk:"addons"`
	ControlPlane  *ControlPlane  `tfsdk:"controlplane"`
	CreatedAt     types.String   `tfsdk:"created_at"`
	ID            types.String   `tfsdk:"id"`
	KubeAPIServer *KubeAPIServer `tfsdk:"kube_api_server"`
	Network       *Network       `tfsdk:"network"`
	ProjectID     types.String   `tfsdk:"project_id"` // Deprecated
	Region        types.String   `tfsdk:"region"`
	Status        *Status        `tfsdk:"status"`
	UpdatedAt     types.String   `tfsdk:"updated_at"`
}

type NodePool struct {
	Name      types.String   `tfsdk:"name"`
	Replicas  types.Int64    `tfsdk:"replicas"`
	Flavor    types.String   `tfsdk:"flavor"`
	AutoScale *AutoScale     `tfsdk:"auto_scale"`
	Tags      []types.String `tfsdk:"tags"`
	Taints    []Taint        `tfsdk:"taints"`
	ID        types.String   `tfsdk:"id"`
	CreatedAt types.String   `tfsdk:"created_at"`
	UpdatedAt types.String   `tfsdk:"updated_at"`
}

type Taint struct {
	Key    types.String `tfsdk:"key"`
	Value  types.String `tfsdk:"value"`
	Effect types.String `tfsdk:"effect"`
}

type AutoScale struct {
	MaxReplicas types.Int64 `tfsdk:"max_replicas"`
	MinReplicas types.Int64 `tfsdk:"min_replicas"`
}

type Addons struct {
	LoadBalance types.String `tfsdk:"loadbalance"`
	Secrets     types.String `tfsdk:"secrets"`
	Volume      types.String `tfsdk:"volume"`
}

type ControlPlane struct {
	AutoScale        *AutoScale        `tfsdk:"auto_scale"`
	CreatedAt        types.String      `tfsdk:"created_at"`
	ID               types.String      `tfsdk:"id"`
	InstanceTemplate *InstanceTemplate `tfsdk:"instance_template"`
	Labels           map[string]string `tfsdk:"labels"`
	Name             types.String      `tfsdk:"name"`
	Replicas         types.Int64       `tfsdk:"replicas"`
	SecurityGroups   []types.String    `tfsdk:"security_groups"`
	Status           *CPStatus         `tfsdk:"status"`
	Tags             []types.String    `tfsdk:"tags"`
	Taints           []Taint           `tfsdk:"taints"`
	UpdatedAt        types.String      `tfsdk:"updated_at"`
	Zone             []types.String    `tfsdk:"zone"`
}

type InstanceTemplate struct {
	DiskSize  types.Int64  `tfsdk:"disk_size"`
	DiskType  types.String `tfsdk:"disk_type"`
	Flavor    *Flavor      `tfsdk:"flavor"`
	NodeImage types.String `tfsdk:"node_image"`
}

type Flavor struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	RAM  types.Int64  `tfsdk:"ram"`
	Size types.Int64  `tfsdk:"size"`
	VCPU types.Int64  `tfsdk:"vcpu"`
}

type KubeAPIServer struct {
	DisableAPIServerFIP types.Bool   `tfsdk:"disable_api_server_fip"`
	FixedIP             types.String `tfsdk:"fixed_ip"`
	FloatingIP          types.String `tfsdk:"floating_ip"`
	Port                types.Int64  `tfsdk:"port"`
}

type Network struct {
	CIDR     types.String `tfsdk:"cidr"`
	Name     types.String `tfsdk:"name"`
	SubnetID types.String `tfsdk:"subnet_id"`
	UUID     types.String `tfsdk:"uuid"`
}

type Status struct {
	Message types.String `tfsdk:"message"`
	State   types.String `tfsdk:"state"`
}

type CPStatus struct {
	Messages []types.String `tfsdk:"messages"`
	State    types.String   `tfsdk:"state"`
}

func ConvertSDKCreateResultToTerraformModel(sdkResult *sdkCluster.GetResult) *KubernetesClusterCreateResourceModel {
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

	// if sdkResult.NodePools != nil {
	// 	tfModel.NodePools = convertSDKNodePoolsToTerraformNodePools(*sdkResult.NodePools)
	// }

	if sdkResult.AllowedCidrs != nil {
		tfModel.AllowedCidrs = convertStringValueSliceToListValue(convertStringSliceToTypesStringSlice(*sdkResult.AllowedCidrs))
	}

	return tfModel

}

func convertStringValueSliceToListValue(stringValues []basetypes.StringValue) basetypes.ListValue {
	// Crie um slice de types.String
	stringSlice := make([]attr.Value, len(stringValues))
	for i, sv := range stringValues {
		stringSlice[i] = sv
	}

	// Crie o ListValue
	listValue := basetypes.NewListValueMust(
		types.StringType,
		stringSlice,
	)

	return listValue
}

func ConvertSDKGetResultToTerraformModel(sdkResult *sdkCluster.GetResult) *KubernetesClusterResourceModel {
	if sdkResult == nil {
		return nil
	}

	tfModel := &KubernetesClusterResourceModel{
		Name:    types.StringValue(sdkResult.Name),
		ID:      types.StringValue(sdkResult.Id),
		Version: types.StringValue(sdkResult.Version),
		Region:  types.StringValue(sdkResult.Region),
	}

	if sdkResult.Description != nil {
		tfModel.Description = types.StringValue(*sdkResult.Description)
	}

	if sdkResult.ProjectId != nil {
		tfModel.ProjectID = types.StringValue(*sdkResult.ProjectId)
	}

	if sdkResult.CreatedAt != nil {
		tfModel.CreatedAt = types.StringValue(*sdkResult.CreatedAt)
	}

	if sdkResult.UpdatedAt != nil {
		tfModel.UpdatedAt = types.StringValue(*sdkResult.UpdatedAt)
	}

	if sdkResult.AllowedCidrs != nil {
		tfModel.AllowedCidrs = convertStringSliceToTypesStringSlice(*sdkResult.AllowedCidrs)
	}

	if sdkResult.Addons != nil {
		tfModel.Addons = convertSDKAddonsToTerraformAddons(*sdkResult.Addons)
	}

	if sdkResult.Controlplane != nil {
		tfModel.ControlPlane = convertSDKControlplaneToTerraformControlPlane(*sdkResult.Controlplane)
	}

	if sdkResult.KubeApiServer != nil {
		tfModel.KubeAPIServer = convertSDKKubeApiServerToTerraformKubeAPIServer(*sdkResult.KubeApiServer)
	}

	if sdkResult.Network != nil {
		tfModel.Network = convertSDKNetworkToTerraformNetwork(*sdkResult.Network)
	}

	if sdkResult.Status != nil {
		tfModel.Status = convertSDKStatusToTerraformStatus(*sdkResult.Status)
	}

	if sdkResult.NodePools != nil {
		tfModel.NodePools = convertSDKNodePoolsToTerraformNodePools(*sdkResult.NodePools)
	}

	// Note: EnabledBastion and EnabledServerGroup are not present in the SDK struct
	// You may need to set these based on some other logic or leave them as their zero values

	return tfModel
}

func convertStringSliceToTypesStringSlice(input []string) []types.String {
	result := make([]types.String, len(input))
	for i, v := range input {
		result[i] = types.StringValue(v)
	}
	return result
}

func convertSDKTagsToTerraformTags(sdkTags sdkCluster.GetResultTags) []types.String {
	tfTags := make([]types.String, len(sdkTags))
	for i, tag := range sdkTags {
		if tag != nil {
			v := types.StringValue(*tag)
			tfTags[i] = v
		}
	}
	return tfTags
}

func convertSDKControlPlaneTagsToTerraformTags(sdkTags sdkCluster.GetResultControlplaneTags) []types.String {
	tfTags := make([]types.String, len(sdkTags))
	for i, tag := range sdkTags {
		if tag != nil {
			tfTags[i] = types.StringValue(*tag)
		}
	}
	return tfTags
}

func convertSDKAddonsToTerraformAddons(sdkAddons sdkCluster.GetResultAddons) *Addons {
	return &Addons{
		LoadBalance: types.StringValue(sdkAddons.Loadbalance),
		Secrets:     types.StringValue(sdkAddons.Secrets),
		Volume:      types.StringValue(sdkAddons.Volume),
	}
}

func convertSDKControlplaneToTerraformControlPlane(sdkControlplane sdkCluster.GetResultControlplane) *ControlPlane {
	cp := &ControlPlane{
		AutoScale: &AutoScale{
			MaxReplicas: types.Int64Value(int64(sdkControlplane.AutoScale.MaxReplicas)),
			MinReplicas: types.Int64Value(int64(sdkControlplane.AutoScale.MinReplicas)),
		},
		ID:       types.StringValue(sdkControlplane.Id),
		Name:     types.StringValue(sdkControlplane.Name),
		Replicas: types.Int64Value(int64(sdkControlplane.Replicas)),
	}

	if sdkControlplane.CreatedAt != nil {
		cp.CreatedAt = types.StringValue(*sdkControlplane.CreatedAt)
	}

	if sdkControlplane.UpdatedAt != nil {
		cp.UpdatedAt = types.StringValue(*sdkControlplane.UpdatedAt)
	}

	cp.InstanceTemplate = &InstanceTemplate{
		DiskSize:  types.Int64Value(int64(sdkControlplane.InstanceTemplate.DiskSize)),
		DiskType:  types.StringValue(sdkControlplane.InstanceTemplate.DiskType),
		NodeImage: types.StringValue(sdkControlplane.InstanceTemplate.NodeImage),
		Flavor: &Flavor{
			ID:   types.StringValue(sdkControlplane.InstanceTemplate.Flavor.Id),
			Name: types.StringValue(sdkControlplane.InstanceTemplate.Flavor.Name),
			RAM:  types.Int64Value(int64(sdkControlplane.InstanceTemplate.Flavor.Ram)),
			Size: types.Int64Value(int64(sdkControlplane.InstanceTemplate.Flavor.Size)),
			VCPU: types.Int64Value(int64(sdkControlplane.InstanceTemplate.Flavor.Vcpu)),
		},
	}

	if sdkControlplane.SecurityGroups != nil {
		cp.SecurityGroups = convertStringSliceToTypesStringSlice(*sdkControlplane.SecurityGroups)
	}

	// Note: Labels conversion is omitted as the SDK type is empty struct
	// if sdkControlplane.Tags != nil {
	// }

	if sdkControlplane.Taints != nil {
		cp.Taints = convertSDKTaintsToTerraformTaints(*sdkControlplane.Taints)
	}

	if sdkControlplane.Zone != nil {
		cp.Zone = convertStringSliceToTypesStringSlice(*sdkControlplane.Zone)
	}

	cp.Status = &CPStatus{
		Messages: convertStringSliceToTypesStringSlice(sdkControlplane.Status.Messages),
		State:    types.StringValue(sdkControlplane.Status.State),
	}

	// Note: Labels conversion is omitted as the SDK type is empty struct

	return cp
}

func convertSDKKubeApiServerToTerraformKubeAPIServer(sdkKubeApiServer sdkCluster.GetResultKubeApiServer) *KubeAPIServer {
	kas := &KubeAPIServer{}

	if sdkKubeApiServer.DisableApiServerFip != nil {
		kas.DisableAPIServerFIP = types.BoolValue(*sdkKubeApiServer.DisableApiServerFip)
	}

	if sdkKubeApiServer.FixedIp != nil {
		kas.FixedIP = types.StringValue(*sdkKubeApiServer.FixedIp)
	}

	if sdkKubeApiServer.FloatingIp != nil {
		kas.FloatingIP = types.StringValue(*sdkKubeApiServer.FloatingIp)
	}

	if sdkKubeApiServer.Port != nil {
		kas.Port = types.Int64Value(int64(*sdkKubeApiServer.Port))
	}

	return kas
}

func convertSDKNetworkToTerraformNetwork(sdkNetwork sdkCluster.GetResultNetwork) *Network {
	return &Network{
		CIDR:     types.StringValue(sdkNetwork.Cidr),
		Name:     types.StringValue(sdkNetwork.Name),
		SubnetID: types.StringValue(sdkNetwork.SubnetId),
		UUID:     types.StringValue(sdkNetwork.Uuid),
	}
}

func convertSDKStatusToTerraformStatus(sdkStatus sdkCluster.GetResultStatus) *Status {
	return &Status{
		Message: types.StringValue(sdkStatus.Message),
		State:   types.StringValue(sdkStatus.State),
	}
}

func convertSDKNodePoolsToTerraformNodePools(sdkNodePools sdkCluster.GetResultNodePools) []NodePool {
	tfNodePools := make([]NodePool, len(sdkNodePools))
	for i, sdkNodePool := range sdkNodePools {
		a := AutoScale{
			MaxReplicas: types.Int64Value(int64(sdkNodePool.AutoScale.MaxReplicas)),
			MinReplicas: types.Int64Value(int64(sdkNodePool.AutoScale.MinReplicas)),
		}

		tfNodePool := NodePool{
			Name:      types.StringValue(sdkNodePool.Name),
			Replicas:  types.Int64Value(int64(sdkNodePool.Replicas)),
			Flavor:    types.StringValue(sdkNodePool.InstanceTemplate.Flavor.Name),
			ID:        types.StringValue(sdkNodePool.Id),
			AutoScale: &a,
		}

		if sdkNodePool.CreatedAt != nil {
			tfNodePool.CreatedAt = types.StringValue(*sdkNodePool.CreatedAt)
		}

		if sdkNodePool.UpdatedAt != nil {
			tfNodePool.UpdatedAt = types.StringValue(*sdkNodePool.UpdatedAt)
		}

		if sdkNodePool.Tags != nil {

			tfNodePool.Tags = convertSDKControlPlaneTagsToTerraformTags(*sdkNodePool.Tags)
		}

		if sdkNodePool.Taints != nil {
			tfNodePool.Taints = convertSDKNodePoolsTaintsToTerraformTaints(*sdkNodePool.Taints)
		}

		tfNodePools[i] = tfNodePool
	}
	return tfNodePools
}

func convertSDKTaintsToTerraformTaints(sdkTaints sdkCluster.GetResultControlplaneTaints) []Taint {
	tfTaints := make([]Taint, len(sdkTaints))
	for i, sdkTaint := range sdkTaints {
		key := types.StringValue(sdkTaint.Key)
		value := types.StringValue(sdkTaint.Value)
		effect := types.StringValue(sdkTaint.Effect)
		tfTaints[i] = Taint{
			Key:    key,
			Value:  value,
			Effect: effect,
		}
	}
	return tfTaints
}

func convertSDKNodePoolsTaintsToTerraformTaints(sdkTaints sdkCluster.GetResultNodePoolsItemTaints) []Taint {
	tfTaints := make([]Taint, len(sdkTaints))
	for i, sdkTaint := range sdkTaints {
		key := types.StringValue(sdkTaint.Key)
		value := types.StringValue(sdkTaint.Value)
		effect := types.StringValue(sdkTaint.Effect)
		tfTaints[i] = Taint{
			Key:    key,
			Value:  value,
			Effect: effect,
		}
	}
	return tfTaints
}

func ConvertTerraformNodePoolToSDKNodePool(tfNodePool NodePoolCreatedResource, as AutoScale, ts []Taint) sdkCluster.CreateParametersNodePoolsItem {
	sdkNodePool := &sdkCluster.CreateParametersNodePoolsItem{
		Name:     tfNodePool.Name.ValueString(),
		Flavor:   tfNodePool.Flavor.ValueString(),
		Replicas: int(tfNodePool.Replicas.ValueInt64()),
	}

	if as.MaxReplicas.IsUnknown() && as.MinReplicas.IsUnknown() {
		sdkNodePool.AutoScale = &sdkCluster.CreateParametersNodePoolsItemAutoScale{
			MaxReplicas: int(as.MaxReplicas.ValueInt64()),
			MinReplicas: int(as.MinReplicas.ValueInt64()),
		}
	}

	if !tfNodePool.Tags.IsNull() {
		tags := make(sdkCluster.CreateParametersNodePoolsItemTags, len(tfNodePool.Tags.Elements()))
		for i, tag := range tfNodePool.Tags.Elements() {
			tags[i] = tag.String()
		}
		sdkNodePool.Tags = &tags
	}

	if !tfNodePool.Taints.IsUnknown() && len(ts) > 0 {
		taints := make(sdkCluster.CreateParametersNodePoolsItemTaints, len(ts))
		for i, t := range ts {
			v := sdkCluster.CreateParametersNodePoolsItemTaintsItem{
				Effect: t.Effect.ValueString(),
				Key:    t.Key.ValueString(),
				Value:  t.Value.ValueString(),
			}
			taints[i] = v
		}
		sdkNodePool.Taints = &taints
	}

	return *sdkNodePool
}
