package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	"magalu.cloud/lib/products/kubernetes/cluster"
	"magalu.cloud/sdk"
	tfutil "magalu.cloud/terraform-provider-mgc/mgc/tfutil"
)

type KubernetesCluster struct {
	ID                         types.String      `tfsdk:"id"`
	Name                       types.String      `tfsdk:"name"`
	EnabledBastion             types.Bool        `tfsdk:"enabled_bastion"`
	NodePools                  []tfutil.NodePool `tfsdk:"node_pools"`
	AllowedCIDRs               []types.String    `tfsdk:"allowed_cidrs"`
	Description                types.String      `tfsdk:"description"`
	EnabledServerGroup         types.Bool        `tfsdk:"enabled_server_group"`
	Version                    types.String      `tfsdk:"version"`
	Zone                       types.String      `tfsdk:"zone"`
	AddonsLoadbalance          types.String      `tfsdk:"addons_loadbalance"`
	AddonsSecrets              types.String      `tfsdk:"addons_secrets"`
	AddonsVolume               types.String      `tfsdk:"addons_volume"`
	Controlplane               *Controlplane     `tfsdk:"controlplane"`
	CreatedAt                  types.String      `tfsdk:"created_at"`
	KubeAPIDisableAPIServerFIP types.Bool        `tfsdk:"kube_api_disable_api_server_fip"`
	KubeAPIFixedIP             types.String      `tfsdk:"kube_api_fixed_ip"`
	KubeAPIFloatingIP          types.String      `tfsdk:"kube_api_floating_ip"`
	KubeAPIPort                types.Int64       `tfsdk:"kube_api_port"`
	CIDR                       types.String      `tfsdk:"cidr"`
	ClusterName                types.String      `tfsdk:"cluster_name"`
	SubnetID                   types.String      `tfsdk:"subnet_id"`
	ProjectID                  types.String      `tfsdk:"project_id"`
	Region                     types.String      `tfsdk:"region"`
	Message                    types.String      `tfsdk:"message"`
	State                      types.String      `tfsdk:"state"`
	UpdatedAt                  types.String      `tfsdk:"updated_at"`
}

type Controlplane struct {
	MaxReplicas    types.Int64    `tfsdk:"max_replicas"`
	MinReplicas    types.Int64    `tfsdk:"min_replicas"`
	CreatedAt      types.String   `tfsdk:"created_at"`
	ID             types.String   `tfsdk:"id"`
	DiskSize       types.Int64    `tfsdk:"disk_size"`
	DiskType       types.String   `tfsdk:"disk_type"`
	FlavorID       types.String   `tfsdk:"flavor_id"`
	FlavorName     types.String   `tfsdk:"flavor_name"`
	NodeImage      types.String   `tfsdk:"node_image"`
	Labels         types.Map      `tfsdk:"labels"`
	Name           types.String   `tfsdk:"name"`
	Replicas       types.Int64    `tfsdk:"replicas"`
	SecurityGroups []types.String `tfsdk:"security_groups"`
	StatusMessages []types.String `tfsdk:"status_messages"`
	State          types.String   `tfsdk:"state"`
	Tags           []types.String `tfsdk:"tags"`
	Taints         []tfutil.Taint `tfsdk:"taints"`
	UpdatedAt      types.String   `tfsdk:"updated_at"`
	Zone           []types.String `tfsdk:"zone"`
}

func NewDataSourceKubernetesCluster() datasource.DataSource {
	return &DataSourceKubernetesCluster{}
}

type DataSourceKubernetesCluster struct {
	sdkClient *mgcSdk.Client
	cluster   cluster.Service
}

func (r *DataSourceKubernetesCluster) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	r.cluster = cluster.NewService(ctx, r.sdkClient)
}

func (d *DataSourceKubernetesCluster) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_cluster"
}

func (d *DataSourceKubernetesCluster) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for Kubernetes cluster in MGC",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Cluster's UUID.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Kubernetes cluster name.",
				Computed:    true,
			},
			"enabled_bastion": schema.BoolAttribute{
				Description: "Indicates if a bastion host is enabled for secure access to the cluster.",
				Computed:    true,
			},
			"node_pools": schema.ListNestedAttribute{
				Description: "An array representing a set of nodes within a Kubernetes cluster.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"flavor_name": schema.StringAttribute{
							Description: "Definition of the CPU, RAM, and storage capacity of the nodes.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the node pool.",
							Computed:    true,
						},
						"replicas": schema.Int64Attribute{
							Description: "Number of replicas of the nodes in the node pool.",
							Computed:    true,
						},
						"max_replicas": schema.Int64Attribute{
							Description: "Maximum number of replicas for autoscaling.",
							Computed:    true,
						},
						"min_replicas": schema.Int64Attribute{
							Description: "Minimum number of replicas for autoscaling.",
							Computed:    true,
						},
						"tags": schema.ListAttribute{
							Description: "List of tags applied to the node pool.",
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
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"effect": schema.StringAttribute{
										Description: "The effect of the taint on pods that do not tolerate the taint.",
										Computed:    true,
									},
									"key": schema.StringAttribute{
										Description: "Key of the taint to be applied to the node.",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Value corresponding to the taint key.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"allowed_cidrs": schema.ListAttribute{
				Description: "List of allowed CIDR blocks for API server access.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A brief description of the Kubernetes cluster.",
				Computed:    true,
			},
			"enabled_server_group": schema.BoolAttribute{
				Description: "Indicates if a server group with anti-affinity policy is used for the cluster and its node pools.",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "The native Kubernetes version of the cluster.",
				Computed:    true,
			},
			"zone": schema.StringAttribute{
				Description: "Identifier of the zone where the Kubernetes cluster is located.",
				Computed:    true,
			},
			"addons_loadbalance": schema.StringAttribute{
				Description: "Flag indicating whether the load balancer service is enabled/disabled in the cluster.",
				Computed:    true,
			},
			"addons_secrets": schema.StringAttribute{
				Description: "Native Kubernetes secret to be included in the cluster during provisioning.",
				Computed:    true,
			},
			"addons_volume": schema.StringAttribute{
				Description: "Flag indicating whether the storage class service is configured by default.",
				Computed:    true,
			},
			"controlplane": schema.SingleNestedAttribute{
				Description: "Object of the node pool response.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"max_replicas": schema.Int64Attribute{
						Description: "Maximum number of replicas for autoscaling.",
						Computed:    true,
					},
					"min_replicas": schema.Int64Attribute{
						Description: "Minimum number of replicas for autoscaling.",
						Computed:    true,
					},
					"created_at": schema.StringAttribute{
						Description: "Date of creation of the Kubernetes cluster.",
						Computed:    true,
					},
					"id": schema.StringAttribute{
						Description: "Node pool's UUID.",
						Computed:    true,
					},
					"disk_size": schema.Int64Attribute{
						Description: "Size of the disk attached to each node.",
						Computed:    true,
					},
					"disk_type": schema.StringAttribute{
						Description: "Type of disk attached to each node.",
						Computed:    true,
					},
					"flavor_id": schema.StringAttribute{
						Description: "Unique identifier for the Flavor.",
						Computed:    true,
					},
					"flavor_name": schema.StringAttribute{
						Description: "Name of the Flavor.",
						Computed:    true,
					},
					"node_image": schema.StringAttribute{
						Description: "Operating system image running on each node.",
						Computed:    true,
					},
					"labels": schema.MapAttribute{
						Description: "Key/value pairs attached to the object and used for specification.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"name": schema.StringAttribute{
						Description: "Node pool name",
						Computed:    true,
					},
					"replicas": schema.Int64Attribute{
						Description: "Number of replicas of the nodes in the node pool.",
						Computed:    true,
					},
					"security_groups": schema.ListAttribute{
						Description: "Name of the security group to define rules allowing network traffic in the worker node pool.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"status_messages": schema.ListAttribute{
						Description: "Detailed message about the status of the node pool or control plane.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"state": schema.StringAttribute{
						Description: "Current state of the node pool or control plane.",
						Computed:    true,
					},
					"tags": schema.ListAttribute{
						Description: "List of tags applied to the node pool.",
						Computed:    true,
						ElementType: types.StringType,
					},
					"taints": schema.ListNestedAttribute{
						Description: "Property for associating a set of nodes.",
						Computed:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"effect": schema.StringAttribute{
									Description: "The effect of the taint on pods that do not tolerate the taint.",
									Computed:    true,
								},
								"key": schema.StringAttribute{
									Description: "Key of the taint to be applied to the node.",
									Computed:    true,
								},
								"value": schema.StringAttribute{
									Description: "Value corresponding to the taint key.",
									Computed:    true,
								},
							},
						},
					},
					"updated_at": schema.StringAttribute{
						Description: "Date of the last change to the Kubernetes cluster.",
						Computed:    true,
					},
					"zone": schema.ListAttribute{
						Description: "Availability zone for creating the Kubernetes cluster.",
						Computed:    true,
						ElementType: types.StringType,
					},
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Creation date of the Kubernetes cluster.",
				Computed:    true,
			},
			"kube_api_disable_api_server_fip": schema.BoolAttribute{
				Description: "Enables or disables the use of Floating IP on the API Server.",
				Computed:    true,
			},
			"kube_api_fixed_ip": schema.StringAttribute{
				Description: "Fixed IP configured for the Kubernetes API Server.",
				Computed:    true,
			},
			"kube_api_floating_ip": schema.StringAttribute{
				Description: "Floating IP created for the Kubernetes API Server.",
				Computed:    true,
			},
			"kube_api_port": schema.Int64Attribute{
				Description: "Port used by the Kubernetes API Server.",
				Computed:    true,
			},
			"cidr": schema.StringAttribute{
				Description: "Available IP addresses used for provisioning nodes in the cluster.",
				Computed:    true,
			},
			"cluster_name": schema.StringAttribute{
				Description: "Name of the node pool.",
				Computed:    true,
			},
			"subnet_id": schema.StringAttribute{
				Description: "Identifier of the internal subnet where the cluster will be provisioned.",
				Computed:    true,
			},
			"project_id": schema.StringAttribute{
				Description:        "(Deprecated) Unique identifier of the project where the cluster was provisioned.",
				Computed:           true,
				DeprecationMessage: "This field is deprecated and will be removed in a future version.",
			},
			"region": schema.StringAttribute{
				Description: "Identifier of the region where the Kubernetes cluster is located.",
				Computed:    true,
			},
			"message": schema.StringAttribute{
				Description: "Detailed message about the status of the cluster or node.",
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: "Current state of the cluster or node.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Date of the last modification of the Kubernetes cluster.",
				Computed:    true,
			},
		},
	}
}

func (d *DataSourceKubernetesCluster) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KubernetesCluster
	diags := resp.State.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	cluster, err := d.cluster.GetContext(ctx, cluster.GetParameters{
		ClusterId: data.ID.ValueString(),
	}, tfutil.GetConfigsFromTags(d.sdkClient.Sdk().Config().Get, cluster.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get cluster", err.Error())
		return
	}
	converted := convertToKubernetesCluster(&cluster)
	resp.Diagnostics.Append(resp.State.Set(ctx, &converted)...)
}

func convertToKubernetesCluster(getResult *cluster.GetResult) *KubernetesCluster {
	if getResult == nil {
		return nil
	}

	kubernetesCluster := &KubernetesCluster{
		ID:   types.StringValue(getResult.Id),
		Name: types.StringValue(getResult.Name),
		// EnabledBastion is not present in GetResult, set to nil
		EnabledBastion: types.BoolNull(),
		// EnabledServerGroup is not present in GetResult, set to nil
		EnabledServerGroup: types.BoolNull(),
		Version:            types.StringValue(getResult.Version),
		// Zone is not directly available in GetResult, set to nil
		Zone:        types.StringNull(),
		Description: types.StringPointerValue(getResult.Description),
		CreatedAt:   types.StringPointerValue(getResult.CreatedAt),
		ProjectID:   types.StringPointerValue(getResult.ProjectId),
		Region:      types.StringValue(getResult.Region),
		UpdatedAt:   types.StringPointerValue(getResult.UpdatedAt),
	}

	// Convert AllowedCIDRs
	if getResult.AllowedCidrs != nil {
		kubernetesCluster.AllowedCIDRs = make([]types.String, len(*getResult.AllowedCidrs))
		for i, cidr := range *getResult.AllowedCidrs {
			kubernetesCluster.AllowedCIDRs[i] = types.StringPointerValue(&cidr)
		}
	}

	// Convert Addons
	if getResult.Addons != nil {
		kubernetesCluster.AddonsLoadbalance = types.StringValue(getResult.Addons.Loadbalance)
		kubernetesCluster.AddonsSecrets = types.StringValue(getResult.Addons.Secrets)
		kubernetesCluster.AddonsVolume = types.StringValue(getResult.Addons.Volume)
	}

	// Convert KubeApiServer
	if getResult.KubeApiServer != nil {
		kubernetesCluster.KubeAPIDisableAPIServerFIP = types.BoolPointerValue(getResult.KubeApiServer.DisableApiServerFip)
		kubernetesCluster.KubeAPIFixedIP = types.StringPointerValue(getResult.KubeApiServer.FixedIp)
		kubernetesCluster.KubeAPIFloatingIP = types.StringPointerValue(getResult.KubeApiServer.FloatingIp)
		if getResult.KubeApiServer.Port != nil {
			kubernetesCluster.KubeAPIPort = types.Int64Value(int64(*getResult.KubeApiServer.Port))
		}
	}

	// Convert Network
	if getResult.Network != nil {
		kubernetesCluster.CIDR = types.StringValue(getResult.Network.Cidr)
		kubernetesCluster.ClusterName = types.StringValue(getResult.Network.Name)
		kubernetesCluster.SubnetID = types.StringValue(getResult.Network.SubnetId)
	}

	// Convert Status
	if getResult.Status != nil {
		kubernetesCluster.Message = types.StringValue(getResult.Status.Message)
		kubernetesCluster.State = types.StringValue(getResult.Status.State)
	}

	// Convert NodePools
	if getResult.NodePools != nil {
		kubernetesCluster.NodePools = make([]tfutil.NodePool, len(*getResult.NodePools))
		for i, np := range *getResult.NodePools {
			kubernetesCluster.NodePools[i] = tfutil.ConvertToNodePool(&np)
		}
	}

	// Convert Controlplane
	if getResult.Controlplane != nil {
		kubernetesCluster.Controlplane = convertToControlplane(getResult.Controlplane)
	}

	return kubernetesCluster
}

func convertToControlplane(cp *cluster.GetResultControlplane) *Controlplane {
	if cp == nil {
		return &Controlplane{}
	}

	controlplane := &Controlplane{
		MaxReplicas: types.Int64Value(int64(cp.AutoScale.MaxReplicas)),
		MinReplicas: types.Int64Value(int64(cp.AutoScale.MinReplicas)),
		CreatedAt:   types.StringPointerValue(cp.CreatedAt),
		ID:          types.StringValue(cp.Id),
		DiskSize:    types.Int64Value(int64(cp.InstanceTemplate.DiskSize)),
		DiskType:    types.StringValue(cp.InstanceTemplate.DiskType),
		FlavorID:    types.StringValue(cp.InstanceTemplate.Flavor.Id),
		FlavorName:  types.StringValue(cp.InstanceTemplate.Flavor.Name),
		NodeImage:   types.StringValue(cp.InstanceTemplate.NodeImage),
		Name:        types.StringValue(cp.Name),
		Replicas:    types.Int64Value(int64(cp.Replicas)),
		State:       types.StringValue(cp.Status.State),
		UpdatedAt:   types.StringPointerValue(cp.UpdatedAt),
	}

	// Convert Labels
	controlplane.Labels = types.MapNull(types.StringType)

	// Convert SecurityGroups
	if cp.SecurityGroups != nil {
		controlplane.SecurityGroups = make([]types.String, len(*cp.SecurityGroups))
		for i, sg := range *cp.SecurityGroups {
			controlplane.SecurityGroups[i] = types.StringPointerValue(&sg)
		}
	}

	// Convert StatusMessages
	if cp.Status.Messages != nil {
		controlplane.StatusMessages = make([]types.String, len(cp.Status.Messages))
		for i, msg := range cp.Status.Messages {
			controlplane.StatusMessages[i] = types.StringPointerValue(&msg)
		}
	}

	// Convert Tags
	if cp.Tags != nil {
		controlplane.Tags = make([]types.String, len(*cp.Tags))
		for i, tag := range *cp.Tags {
			controlplane.Tags[i] = types.StringPointerValue(tag)
		}
	}

	// Convert Taints
	if cp.Taints != nil {
		controlplane.Taints = make([]tfutil.Taint, len(*cp.Taints))
		for i, taint := range *cp.Taints {
			controlplane.Taints[i] = tfutil.Taint{
				Effect: types.StringValue(taint.Effect),
				Key:    types.StringValue(taint.Key),
				Value:  types.StringValue(taint.Value),
			}
		}
	}

	// Convert Zone
	if cp.Zone != nil {
		controlplane.Zone = make([]types.String, len(*cp.Zone))
		for i, zone := range *cp.Zone {
			controlplane.Zone[i] = types.StringPointerValue(&zone)
		}
	}

	return controlplane
}
