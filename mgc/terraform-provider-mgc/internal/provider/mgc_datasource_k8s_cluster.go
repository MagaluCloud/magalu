package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	"magalu.cloud/lib/products/kubernetes/cluster"
	"magalu.cloud/sdk"
)

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

	sdk, ok := req.ProviderData.(*sdk.Sdk)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected provider config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

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
				Required:            true,
			},
			"name": schema.StringAttribute{
				Description: "Kubernetes cluster name.",
				Computed:            true,
			},
			"enabled_bastion": schema.BoolAttribute{
				Description: "Indicates if a bastion host is enabled for secure access to the cluster.",
				Computed:            true,
			},
			"node_pools": schema.ListNestedAttribute{
				Description: "An array representing a set of nodes within a Kubernetes cluster.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"flavor": schema.StringAttribute{
							Description: "Definition of the CPU, RAM, and storage capacity of the nodes.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the node pool.",
							Computed:            true,
						},
						"replicas": schema.Int64Attribute{
							Description: "Number of replicas of the nodes in the node pool.",
							Computed:            true,
						},
						"auto_scale": schema.SingleNestedAttribute{
							Description: "Object specifying properties for updating workload resources in the Kubernetes cluster.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"max_replicas": schema.Int64Attribute{
									Description: "Maximum number of replicas for autoscaling.",
									Computed:            true,
								},
								"min_replicas": schema.Int64Attribute{
									Description: "Minimum number of replicas for autoscaling.",
									Computed:            true,
								},
							},
						},
						"tags": schema.ListAttribute{
							Description: "List of tags applied to the node pool.",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"created_at": schema.StringAttribute{
							Description: "Date of creation of the Kubernetes Node.",
							Computed:            true,
						},
						"updated_at": schema.StringAttribute{
							Description: "Date of the last change to the Kubernetes Node.",
							Computed:            true,
						},
						"id": schema.StringAttribute{
							Description: "Node pool's UUID.",
							Computed:            true,
						},
						"taints": schema.ListNestedAttribute{
							Description: "Property associating a set of nodes.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"effect": schema.StringAttribute{
										Description: "The effect of the taint on pods that do not tolerate the taint.",
										Computed:            true,
									},
									"key": schema.StringAttribute{
										Description: "Key of the taint to be applied to the node.",
										Computed:            true,
									},
									"value": schema.StringAttribute{
										Description: "Value corresponding to the taint key.",
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
			"allowed_cidrs": schema.ListAttribute{
				Description: "List of allowed CIDR blocks for API server access.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A brief description of the Kubernetes cluster.",
				Computed:            true,
			},
			"enabled_server_group": schema.BoolAttribute{
				Description: "Indicates if a server group with anti-affinity policy is used for the cluster and its node pools.",
				Computed:            true,
			},
			"version": schema.StringAttribute{
				Description: "The native Kubernetes version of the cluster.",
				Computed:            true,
			},
			"zone": schema.StringAttribute{
				Description: "Identifier of the zone where the Kubernetes cluster is located.",
				Computed:            true,
			},
			"addons": schema.SingleNestedAttribute{
				Description: "Object representing addons that extend the functionality of the Kubernetes cluster.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"loadbalance": schema.StringAttribute{
						Description: "Flag indicating whether the load balancer service is enabled/disabled in the cluster.",
						Computed:            true,
					},
					"secrets": schema.StringAttribute{
						Description: "Native Kubernetes secret to be included in the cluster during provisioning.",
						Computed:            true,
					},
					"volume": schema.StringAttribute{
						Description: "Flag indicating whether the storage class service is configured by default.",
						Computed:            true,
					},
				},
			},
			"controlplane": schema.SingleNestedAttribute{
				Description: "Object of the node pool response.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"auto_scale": schema.SingleNestedAttribute{
						Description: "Object specifying properties for updating workload resources in the Kubernetes cluster.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"max_replicas": schema.Int64Attribute{
								Description: "Maximum number of replicas for autoscaling.",
								Computed:            true,
							},
							"min_replicas": schema.Int64Attribute{
								Description: "Minimum number of replicas for autoscaling.",
								Computed:            true,
							},
						},
					},
					"created_at": schema.StringAttribute{
						Description: "Date of creation of the Kubernetes cluster.",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						Description: "Node pool's UUID.",
						Computed:            true,
					},
					"instance_template": schema.SingleNestedAttribute{
						Description: "Template for the instance object used to create machine instances and managed instance groups.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"disk_size": schema.Int64Attribute{
								Description: "Size of the disk attached to each node.",
								Computed:            true,
							},
							"disk_type": schema.StringAttribute{
								Description: "Type of disk attached to each node.",
								Computed:            true,
							},
							"flavor": schema.SingleNestedAttribute{
								Description: "Definition of CPU capacity, RAM, and storage for nodes.",
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "Unique identifier for the Flavor.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										Description: "Name of the Flavor.",
										Computed:            true,
									},
									"ram": schema.Int64Attribute{
										Description: "Amount of RAM, measured in MB.",
										Computed:            true,
									},
									"size": schema.Int64Attribute{
										Description: "Amount of disk space, measured in GB.",
										Computed:            true,
									},
									"vcpu": schema.Int64Attribute{
										Description: "Number of available vCPUs.",
										Computed:            true,
									},
								},
							},
							"node_image": schema.StringAttribute{
								Description: "Operating system image running on each node.",
								Computed:            true,
							},
						},
					},
					"labels": schema.MapAttribute{
						Description: "Key/value pairs attached to the object and used for specification.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"name": schema.StringAttribute{
						Description: "Node pool name",
						Computed:            true,
					},
					"replicas": schema.Int64Attribute{
						Description: "Number of replicas of the nodes in the node pool.",
						Computed:            true,
					},
					"security_groups": schema.ListAttribute{
						Description: "Name of the security group to define rules allowing network traffic in the worker node pool.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"status": schema.SingleNestedAttribute{
						Description: "Details about the status of the node pool or control plane.",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"messages": schema.ListAttribute{
								Description: "Detailed message about the status of the node pool or control plane.",
								Computed:            true,
								ElementType:         types.StringType,
							},
							"state": schema.StringAttribute{
								Description: "Current state of the node pool or control plane.",
								Computed:            true,
							},
						},
					},
					"tags": schema.ListAttribute{
						Description: "List of tags applied to the node pool.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"taints": schema.ListNestedAttribute{
						Description: "Property for associating a set of nodes.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"effect": schema.StringAttribute{
									Description: "The effect of the taint on pods that do not tolerate the taint.",
									Computed:            true,
								},
								"key": schema.StringAttribute{
									Description: "Key of the taint to be applied to the node.",
									Computed:            true,
								},
								"value": schema.StringAttribute{
									Description: "Value corresponding to the taint key.",
									Computed:            true,
								},
							},
						},
					},
					"updated_at": schema.StringAttribute{
						Description: "Date of the last change to the Kubernetes cluster.",
						Computed:            true,
					},
					"zone": schema.ListAttribute{
						Description: "Availability zone for creating the Kubernetes cluster.",
						Computed:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"created_at": schema.StringAttribute{
				Description: "Creation date of the Kubernetes cluster.",
				Computed:            true,
			},
			"kube_api_server": schema.SingleNestedAttribute{
				Description: "Information about the Kubernetes API Server of the cluster.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"disable_api_server_fip": schema.BoolAttribute{
						Description: "Enables or disables the use of Floating IP on the API Server.",
						Computed:            true,
					},
					"fixed_ip": schema.StringAttribute{
						Description: "Fixed IP configured for the Kubernetes API Server.",
						Computed:            true,
					},
					"floating_ip": schema.StringAttribute{
						Description: "Floating IP created for the Kubernetes API Server.",
						Computed:            true,
					},
					"port": schema.Int64Attribute{
						Description: "Port used by the Kubernetes API Server.",
						Computed:            true,
					},
				},
			},
			"network": schema.SingleNestedAttribute{
				Description: "Response object for the Kubernetes cluster network resource request.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"cidr": schema.StringAttribute{
						Description: "Available IP addresses used for provisioning nodes in the cluster.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						Description: "Name of the node pool.",
						Computed:            true,
					},
					"subnet_id": schema.StringAttribute{
						Description: "Identifier of the internal subnet where the cluster will be provisioned.",
						Computed:            true,
					},
					"uuid": schema.StringAttribute{
						Description: "Nodepool's UUID.",
						Computed:            true,
					},
				},
			},
			"project_id": schema.StringAttribute{
				Description: "(Deprecated) Unique identifier of the project where the cluster was provisioned.",
				Computed:            true,
				DeprecationMessage:  "This field is deprecated and will be removed in a future version.",
			},
			"region": schema.StringAttribute{
				Description: "Identifier of the region where the Kubernetes cluster is located.",
				Computed:            true,
			},
			"status": schema.SingleNestedAttribute{
				Description: "Details about the status of the Kubernetes cluster or node.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"message": schema.StringAttribute{
						Description: "Detailed message about the status of the cluster or node.",
						Computed:            true,
					},
					"state": schema.StringAttribute{
						Description: "Current state of the cluster or node.",
						Computed:            true,
					},
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "Date of the last modification of the Kubernetes cluster.",
				Computed:            true,
			},
		},
	}
}

func (d *DataSourceKubernetesCluster) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KubernetesClusterResourceModel
	diags := resp.State.Get(ctx, &data)
	if diags.HasError() {
		resp.Diagnostics = diags
		return
	}

	cluster, err := d.cluster.Get(cluster.GetParameters{
		ClusterId: data.ID.ValueString(),
	}, cluster.GetConfigs{})
	if err != nil {
		resp.Diagnostics.AddError("Failed to get cluster", err.Error())
		return
	}
	converted := ConvertSDKGetResultToTerraformModel(&cluster)
	resp.Diagnostics.Append(resp.State.Set(ctx, &converted)...)
}
