package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkNodepool "magalu.cloud/lib/products/kubernetes/nodepool"
	"magalu.cloud/terraform-provider-mgc/mgc/client"
	tfutil "magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type NodesDataSourceModel struct {
	ClusterID  types.String      `tfsdk:"cluster_id"`
	NodepoolID types.String      `tfsdk:"nodepool_id"`
	Nodes      []NodesResultFlat `tfsdk:"nodes"`
}

type NodesResultFlat struct {
	ID                             types.String   `tfsdk:"id"`
	ClusterName                    types.String   `tfsdk:"cluster_name"`
	CreatedAt                      types.String   `tfsdk:"created_at"`
	Flavor                         types.String   `tfsdk:"flavor"`
	Name                           types.String   `tfsdk:"name"`
	Namespace                      types.String   `tfsdk:"namespace"`
	NodeImage                      types.String   `tfsdk:"node_image"`
	NodepoolName                   types.String   `tfsdk:"nodepool_name"`
	Zone                           types.String   `tfsdk:"zone"`
	Addresses                      []AddressItem  `tfsdk:"addresses"`
	StatusMessage                  types.String   `tfsdk:"status_message"`
	StatusState                    types.String   `tfsdk:"status_state"`
	InfrastructureArchitecture     types.String   `tfsdk:"infrastructure_architecture"`
	InfrastructureContainerRuntime types.String   `tfsdk:"infrastructure_container_runtime"`
	InfrastructureKernelVersion    types.String   `tfsdk:"infrastructure_kernel_version"`
	InfrastructureKubeProxyVersion types.String   `tfsdk:"infrastructure_kube_proxy_version"`
	InfrastructureKubeletVersion   types.String   `tfsdk:"infrastructure_kubelet_version"`
	InfrastructureOperatingSystem  types.String   `tfsdk:"infrastructure_operating_system"`
	InfrastructureOsImage          types.String   `tfsdk:"infrastructure_os_image"`
	AllocatableCPU                 types.String   `tfsdk:"allocatable_cpu"`
	AllocatableEphemeralStorage    types.String   `tfsdk:"allocatable_ephemeral_storage"`
	AllocatableHugepages1Gi        types.String   `tfsdk:"allocatable_hugepages_1gi"`
	AllocatableHugepages2Mi        types.String   `tfsdk:"allocatable_hugepages_2mi"`
	AllocatableMemory              types.String   `tfsdk:"allocatable_memory"`
	AllocatablePods                types.String   `tfsdk:"allocatable_pods"`
	CapacityCPU                    types.String   `tfsdk:"capacity_cpu"`
	CapacityEphemeralStorage       types.String   `tfsdk:"capacity_ephemeral_storage"`
	CapacityHugepages1Gi           types.String   `tfsdk:"capacity_hugepages_1gi"`
	CapacityHugepages2Mi           types.String   `tfsdk:"capacity_hugepages_2mi"`
	CapacityMemory                 types.String   `tfsdk:"capacity_memory"`
	CapacityPods                   types.String   `tfsdk:"capacity_pods"`
	Taints                         []tfutil.Taint `tfsdk:"taints"`
}

type AddressItem struct {
	Address types.String `tfsdk:"address"`
	Type    types.String `tfsdk:"type"`
}

type DataSourceKubernetesNode struct {
	sdkClient *mgcSdk.Client
	nodepool  sdkNodepool.Service
}

func NewDataSourceKubernetesNode() datasource.DataSource {
	return &DataSourceKubernetesNode{}
}

func (d *DataSourceKubernetesNode) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_node"
}

func (d *DataSourceKubernetesNode) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for Kubernetes cluster in MGC",
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Description: "ID of the cluster.",
				Required:    true,
			},
			"nodepool_id": schema.StringAttribute{
				Description: "ID of the node pool.",
				Required:    true,
			},
			"nodes": schema.ListNestedAttribute{
				Description: "List of nodes in the cluster.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Node's UUID.",
							Computed:    true,
						},
						"cluster_name": schema.StringAttribute{
							Description: "Name of the cluster.",
							Computed:    true,
						},
						"created_at": schema.StringAttribute{
							Description: "Creation timestamp.",
							Computed:    true,
						},
						"flavor": schema.StringAttribute{
							Description: "Node flavor.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Node name.",
							Computed:    true,
						},
						"namespace": schema.StringAttribute{
							Description: "Namespace of the node.",
							Computed:    true,
						},
						"node_image": schema.StringAttribute{
							Description: "Image of the node.",
							Computed:    true,
						},
						"nodepool_name": schema.StringAttribute{
							Description: "Name of the nodepool.",
							Computed:    true,
						},
						"zone": schema.StringAttribute{
							Description: "Zone of the node.",
							Computed:    true,
						},
						"addresses": schema.ListNestedAttribute{
							Description: "List of node addresses.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"address": schema.StringAttribute{
										Description: "The address of the node.",
										Computed:    true,
									},
									"type": schema.StringAttribute{
										Description: "The type of the address.",
										Computed:    true,
									},
								},
							},
						},
						"status_message": schema.StringAttribute{
							Description: "Status message of the node.",
							Computed:    true,
						},
						"status_state": schema.StringAttribute{
							Description: "State of the node.",
							Computed:    true,
						},
						"infrastructure_architecture": schema.StringAttribute{
							Description: "Node architecture.",
							Computed:    true,
						},
						"infrastructure_container_runtime": schema.StringAttribute{
							Description: "Container runtime version.",
							Computed:    true,
						},
						"infrastructure_kernel_version": schema.StringAttribute{
							Description: "Kernel version of the node.",
							Computed:    true,
						},
						"infrastructure_kube_proxy_version": schema.StringAttribute{
							Description: "Kube-proxy version.",
							Computed:    true,
						},
						"infrastructure_kubelet_version": schema.StringAttribute{
							Description: "Kubelet version.",
							Computed:    true,
						},
						"infrastructure_operating_system": schema.StringAttribute{
							Description: "Operating system of the node.",
							Computed:    true,
						},
						"infrastructure_os_image": schema.StringAttribute{
							Description: "OS image of the node.",
							Computed:    true,
						},
						"allocatable_cpu": schema.StringAttribute{
							Description: "Allocatable CPU.",
							Computed:    true,
						},
						"allocatable_ephemeral_storage": schema.StringAttribute{
							Description: "Allocatable ephemeral storage.",
							Computed:    true,
						},
						"allocatable_hugepages_1gi": schema.StringAttribute{
							Description: "Allocatable hugepages (1Gi).",
							Computed:    true,
						},
						"allocatable_hugepages_2mi": schema.StringAttribute{
							Description: "Allocatable hugepages (2Mi).",
							Computed:    true,
						},
						"allocatable_memory": schema.StringAttribute{
							Description: "Allocatable memory.",
							Computed:    true,
						},
						"allocatable_pods": schema.StringAttribute{
							Description: "Allocatable pods.",
							Computed:    true,
						},
						"capacity_cpu": schema.StringAttribute{
							Description: "CPU capacity.",
							Computed:    true,
						},
						"capacity_ephemeral_storage": schema.StringAttribute{
							Description: "Ephemeral storage capacity.",
							Computed:    true,
						},
						"capacity_hugepages_1gi": schema.StringAttribute{
							Description: "Hugepages capacity (1Gi).",
							Computed:    true,
						},
						"capacity_hugepages_2mi": schema.StringAttribute{
							Description: "Hugepages capacity (2Mi).",
							Computed:    true,
						},
						"capacity_memory": schema.StringAttribute{
							Description: "Memory capacity.",
							Computed:    true,
						},
						"capacity_pods": schema.StringAttribute{
							Description: "Pods capacity.",
							Computed:    true,
						},
						"taints": schema.ListNestedAttribute{
							Description: "List of node taints.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"effect": schema.StringAttribute{
										Description: "The effect of the taint.",
										Computed:    true,
									},
									"key": schema.StringAttribute{
										Description: "The key of the taint.",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "The value of the taint.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *DataSourceKubernetesNode) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	var err error
	d.sdkClient, err = client.NewSDKClient(req)
	if err != nil {
		resp.Diagnostics.AddError(
			err.Error(),
			fmt.Sprintf("Expected provider config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.nodepool = sdkNodepool.NewService(ctx, d.sdkClient)
}

func (d *DataSourceKubernetesNode) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NodesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	node, err := d.nodepool.NodesContext(ctx, sdkNodepool.NodesParameters{
		ClusterId:  data.ClusterID.ValueString(),
		NodePoolId: data.NodepoolID.ValueString(),
	}, tfutil.GetConfigsFromTags(d.sdkClient.Sdk().Config().Get, sdkNodepool.NodesConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get node", err.Error())
		return
	}

	data.Nodes = make([]NodesResultFlat, len(node.Results))
	for i, n := range node.Results {
		data.Nodes[i] = *convertToTerraformKubernetesCluster(&n)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func convertToTerraformKubernetesCluster(original *sdkNodepool.NodesResultResultsItem) *NodesResultFlat {
	if original == nil {
		return nil
	}

	result := &NodesResultFlat{
		ID:           types.StringValue(original.Id),
		ClusterName:  types.StringValue(original.ClusterName),
		CreatedAt:    types.StringValue(original.CreatedAt),
		Flavor:       types.StringValue(original.Flavor),
		Name:         types.StringValue(original.Name),
		Namespace:    types.StringValue(original.Namespace),
		NodeImage:    types.StringValue(original.NodeImage),
		NodepoolName: types.StringValue(original.NodepoolName),
		Zone:         types.StringValue(stringPtrToString(original.Zone)),
		Addresses:    convertAddresses(original.Addresses),
		Taints:       convertTaints(original.Taints),
	}

	// Convert Status
	result.StatusMessage = types.StringValue(original.Status.Message)
	result.StatusState = types.StringValue(original.Status.State)

	// Convert Infrastructure
	result.InfrastructureArchitecture = types.StringValue(original.Infrastructure.Architecture)
	result.InfrastructureContainerRuntime = types.StringValue(original.Infrastructure.ContainerRuntimeVersion)
	result.InfrastructureKernelVersion = types.StringValue(original.Infrastructure.KernelVersion)
	result.InfrastructureKubeProxyVersion = types.StringValue(original.Infrastructure.KubeProxyVersion)
	result.InfrastructureKubeletVersion = types.StringValue(original.Infrastructure.KubeletVersion)
	result.InfrastructureOperatingSystem = types.StringValue(original.Infrastructure.OperatingSystem)
	result.InfrastructureOsImage = types.StringValue(original.Infrastructure.OsImage)

	// Convert Allocatable
	result.AllocatableCPU = types.StringValue(original.Infrastructure.Allocatable.Cpu)
	result.AllocatableEphemeralStorage = types.StringValue(original.Infrastructure.Allocatable.EphemeralStorage)
	result.AllocatableHugepages1Gi = types.StringValue(original.Infrastructure.Allocatable.Hugepages1gi)
	result.AllocatableHugepages2Mi = types.StringValue(original.Infrastructure.Allocatable.Hugepages2mi)
	result.AllocatableMemory = types.StringValue(original.Infrastructure.Allocatable.Memory)
	result.AllocatablePods = types.StringValue(original.Infrastructure.Allocatable.Pods)

	// Convert Capacity
	result.CapacityCPU = types.StringValue(original.Infrastructure.Capacity.Cpu)
	result.CapacityEphemeralStorage = types.StringValue(original.Infrastructure.Capacity.EphemeralStorage)
	result.CapacityHugepages1Gi = types.StringValue(original.Infrastructure.Capacity.Hugepages1gi)
	result.CapacityHugepages2Mi = types.StringValue(original.Infrastructure.Capacity.Hugepages2mi)
	result.CapacityMemory = types.StringValue(original.Infrastructure.Capacity.Memory)
	result.CapacityPods = types.StringValue(original.Infrastructure.Capacity.Pods)

	return result
}

// Helper functions

func convertAddresses(addresses sdkNodepool.NodesResultResultsItemAddresses) []AddressItem {
	result := make([]AddressItem, len(addresses))
	for i, addr := range addresses {
		result[i] = AddressItem{
			Address: types.StringValue(addr.Address),
			Type:    types.StringValue(addr.Type),
		}
	}
	return result
}

func convertTaints(taints sdkNodepool.NodesResultResultsItemTaints) []tfutil.Taint {
	result := make([]tfutil.Taint, len(taints))
	for i, taint := range taints {
		result[i] = tfutil.Taint{
			Effect: types.StringValue(taint.Effect),
			Key:    types.StringValue(taint.Key),
			Value:  types.StringValue(taint.Value),
		}
	}
	return result
}

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
