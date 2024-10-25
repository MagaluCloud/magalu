package tfutil

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"magalu.cloud/lib/products/kubernetes/cluster"
	sdkNodepool "magalu.cloud/lib/products/kubernetes/nodepool"
)

type NodePoolCreate struct {
	ClusterID types.String `tfsdk:"cluster_id"`
	NodePool
}

type NodePool struct {
	Flavor      types.String    `tfsdk:"flavor_name"`
	Name        types.String    `tfsdk:"name"`
	Replicas    types.Int64     `tfsdk:"replicas"`
	MaxReplicas types.Int64     `tfsdk:"max_replicas"`
	MinReplicas types.Int64     `tfsdk:"min_replicas"`
	Tags        *[]types.String `tfsdk:"tags"`
	CreatedAt   types.String    `tfsdk:"created_at"`
	UpdatedAt   types.String    `tfsdk:"updated_at"`
	ID          types.String    `tfsdk:"id"`
	Taints      *[]Taint        `tfsdk:"taints"`
}

type Taint struct {
	Effect types.String `tfsdk:"effect"`
	Key    types.String `tfsdk:"key"`
	Value  types.String `tfsdk:"value"`
}

func ConvertToNodePool(np *cluster.GetResultNodePoolsItem) NodePool {
	if np == nil {
		return NodePool{}
	}

	nodePool := NodePool{
		Name:        types.StringValue(np.Name),
		Replicas:    types.Int64Value(int64(np.Replicas)),
		MaxReplicas: types.Int64Value(int64(np.AutoScale.MaxReplicas)),
		MinReplicas: types.Int64Value(int64(np.AutoScale.MinReplicas)),
		CreatedAt:   types.StringPointerValue(np.CreatedAt),
		UpdatedAt:   types.StringPointerValue(np.UpdatedAt),
		ID:          types.StringValue(np.Id),
	}

	if np.InstanceTemplate.Flavor.Name != "" {
		nodePool.Flavor = types.StringValue(np.InstanceTemplate.Flavor.Name)
	}

	if np.Tags != nil {
		tags := make([]types.String, len(*np.Tags))
		for i, tag := range *np.Tags {
			tags[i] = types.StringPointerValue(tag)
		}
		if len(tags) > 0 {
			nodePool.Tags = &tags
		}
	}

	if np.Taints != nil {
		var taints []Taint
		for _, taint := range *np.Taints {
			taints = append(taints, Taint{
				Effect: types.StringValue(taint.Effect),
				Key:    types.StringValue(taint.Key),
				Value:  types.StringValue(taint.Value),
			})
		}
		nodePool.Taints = &taints
	}

	return nodePool
}

func ConvertToNodePoolGet(np *sdkNodepool.GetResult) NodePool {
	if np == nil {
		return NodePool{}
	}

	nodePool := NodePool{
		Name:        types.StringValue(np.Name),
		Replicas:    types.Int64Value(int64(np.Replicas)),
		MaxReplicas: types.Int64Value(int64(np.AutoScale.MaxReplicas)),
		MinReplicas: types.Int64Value(int64(np.AutoScale.MinReplicas)),
		CreatedAt:   types.StringPointerValue(np.CreatedAt),
		UpdatedAt:   types.StringPointerValue(np.UpdatedAt),
		ID:          types.StringValue(np.Id),
	}

	if np.InstanceTemplate.Flavor.Name != "" {
		nodePool.Flavor = types.StringValue(np.InstanceTemplate.Flavor.Name)
	}

	if np.Tags != nil {
		tags := make([]types.String, len(*np.Tags))
		for i, tag := range *np.Tags {
			tags[i] = types.StringPointerValue(tag)
		}
		if len(tags) > 0 {
			nodePool.Tags = &tags
		}
	}

	if np.Taints != nil {
		var taints []Taint
		for _, taint := range *np.Taints {
			taints = append(taints, Taint{
				Effect: types.StringValue(taint.Effect),
				Key:    types.StringValue(taint.Key),
				Value:  types.StringValue(taint.Value),
			})
		}
		nodePool.Taints = &taints
	}

	return nodePool
}

func ConvertStringSliceToTypesStringSlice(input []string) []types.String {
	result := make([]types.String, len(input))
	for i, v := range input {
		result[i] = types.StringValue(v)
	}
	return result
}

func ConvertToNodePoolCreate(np *sdkNodepool.CreateResult) NodePool {
	if np == nil {
		return NodePool{}
	}

	nodePool := NodePool{
		Name:        types.StringValue(np.Name),
		Replicas:    types.Int64Value(int64(np.Replicas)),
		MaxReplicas: types.Int64Value(int64(np.AutoScale.MaxReplicas)),
		MinReplicas: types.Int64Value(int64(np.AutoScale.MinReplicas)),
		CreatedAt:   types.StringPointerValue(np.CreatedAt),
		UpdatedAt:   types.StringPointerValue(np.UpdatedAt),
		ID:          types.StringValue(np.Id),
	}

	if np.InstanceTemplate.Flavor.Name != "" {
		nodePool.Flavor = types.StringValue(np.InstanceTemplate.Flavor.Name)
	}

	if np.Tags != nil {
		tags := make([]types.String, len(*np.Tags))
		for i, tag := range *np.Tags {
			tags[i] = types.StringPointerValue(tag)
		}
		if len(tags) > 0 {
			nodePool.Tags = &tags
		}
	}

	if np.Taints != nil {
		var taints []Taint
		for _, taint := range *np.Taints {
			taints = append(taints, Taint{
				Effect: types.StringValue(taint.Effect),
				Key:    types.StringValue(taint.Key),
				Value:  types.StringValue(taint.Value),
			})
		}
		nodePool.Taints = &taints
	}

	return nodePool
}

func ConvertToNodePoolUpdate(np *sdkNodepool.UpdateResult) NodePool {
	if np == nil {
		return NodePool{}
	}

	nodePool := NodePool{
		Name:        types.StringValue(np.Name),
		Replicas:    types.Int64Value(int64(np.Replicas)),
		MaxReplicas: types.Int64Value(int64(np.AutoScale.MaxReplicas)),
		MinReplicas: types.Int64Value(int64(np.AutoScale.MinReplicas)),
		CreatedAt:   types.StringPointerValue(np.CreatedAt),
		UpdatedAt:   types.StringPointerValue(np.UpdatedAt),
		ID:          types.StringValue(np.Id),
	}

	if np.InstanceTemplate.Flavor.Name != "" {
		nodePool.Flavor = types.StringValue(np.InstanceTemplate.Flavor.Name)
	}

	if np.Tags != nil {
		tags := make([]types.String, len(*np.Tags))
		for i, tag := range *np.Tags {
			tags[i] = types.StringPointerValue(tag)
		}
		if len(tags) > 0 {
			nodePool.Tags = &tags
		}
	}

	if np.Taints != nil {
		var taints []Taint
		for _, taint := range *np.Taints {
			taints = append(taints, Taint{
				Effect: types.StringValue(taint.Effect),
				Key:    types.StringValue(taint.Key),
				Value:  types.StringValue(taint.Value),
			})
		}
		nodePool.Taints = &taints
	}

	return nodePool
}
