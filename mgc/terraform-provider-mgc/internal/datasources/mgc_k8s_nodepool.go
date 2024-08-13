package datasources

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	mgcSdk "magalu.cloud/lib"
	sdkNodepool "magalu.cloud/lib/products/kubernetes/nodepool"
	"magalu.cloud/sdk"
	tfutil "magalu.cloud/terraform-provider-mgc/internal/tfutil"
)

type FlattenedGetResult struct {
	ID                         types.String   `tfsdk:"id"`
	ClusterID                  types.String   `tfsdk:"cluster_id"`
	Name                       types.String   `tfsdk:"name"`
	CreatedAt                  types.String   `tfsdk:"created_at"`
	UpdatedAt                  types.String   `tfsdk:"updated_at"`
	Replicas                   types.Int64    `tfsdk:"replicas"`
	AutoScaleMaxReplicas       types.Int64    `tfsdk:"auto_scale_max_replicas"`
	AutoScaleMinReplicas       types.Int64    `tfsdk:"auto_scale_min_replicas"`
	InstanceTemplateDiskSize   types.Int64    `tfsdk:"instance_template_disk_size"`
	InstanceTemplateDiskType   types.String   `tfsdk:"instance_template_disk_type"`
	InstanceTemplateNodeImage  types.String   `tfsdk:"instance_template_node_image"`
	InstanceTemplateFlavorID   types.String   `tfsdk:"instance_template_flavor_id"`
	InstanceTemplateFlavorName types.String   `tfsdk:"instance_template_flavor_name"`
	InstanceTemplateFlavorRam  types.Int64    `tfsdk:"instance_template_flavor_ram"`
	InstanceTemplateFlavorSize types.Int64    `tfsdk:"instance_template_flavor_size"`
	InstanceTemplateFlavorVcpu types.Int64    `tfsdk:"instance_template_flavor_vcpu"`
	Labels                     types.Map      `tfsdk:"labels"`
	SecurityGroups             []types.String `tfsdk:"security_groups"`
	StatusState                types.String   `tfsdk:"status_state"`
	StatusMessages             []types.String `tfsdk:"status_messages"`
	Tags                       types.List     `tfsdk:"tags"`
	Taints                     []tfutil.Taint `tfsdk:"taints"`
	Zone                       types.List     `tfsdk:"zone"`
}

type DataSourceKubernetesNodepool struct {
	sdkClient *mgcSdk.Client
	nodepool  sdkNodepool.Service
}

func NewDataSourceKubernetesNodepool() datasource.DataSource {
	return &DataSourceKubernetesNodepool{}
}

func (r *DataSourceKubernetesNodepool) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_nodepool"
}

func (r *DataSourceKubernetesNodepool) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	sdk, ok := req.ProviderData.(*sdk.Sdk)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected provider config, got: %T. Please report this issue to the provider developers.",
		)
		return
	}

	r.sdkClient = mgcSdk.NewClient(sdk)
	r.nodepool = sdkNodepool.NewService(ctx, r.sdkClient)
}

func (d *DataSourceKubernetesNodepool) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for Kubernetes Nodepool",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Nodepool UUID.",
				Required:    true,
			},
			"cluster_id": schema.StringAttribute{
				Description: "Cluster UUID.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the nodepool.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Creation timestamp.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Last update timestamp.",
				Computed:    true,
			},
			"replicas": schema.Int64Attribute{
				Description: "Number of replicas.",
				Computed:    true,
			},
			"auto_scale_max_replicas": schema.Int64Attribute{
				Description: "Maximum number of replicas for auto-scaling.",
				Computed:    true,
			},
			"auto_scale_min_replicas": schema.Int64Attribute{
				Description: "Minimum number of replicas for auto-scaling.",
				Computed:    true,
			},
			"instance_template_disk_size": schema.Int64Attribute{
				Description: "Disk size in GB for the instance template.",
				Computed:    true,
			},
			"instance_template_disk_type": schema.StringAttribute{
				Description: "Disk type for the instance template.",
				Computed:    true,
			},
			"instance_template_node_image": schema.StringAttribute{
				Description: "Node image for the instance template.",
				Computed:    true,
			},
			"instance_template_flavor_id": schema.StringAttribute{
				Description: "Flavor ID for the instance template.",
				Computed:    true,
			},
			"instance_template_flavor_name": schema.StringAttribute{
				Description: "Flavor name for the instance template.",
				Computed:    true,
			},
			"instance_template_flavor_ram": schema.Int64Attribute{
				Description: "RAM in MB for the instance template flavor.",
				Computed:    true,
			},
			"instance_template_flavor_size": schema.Int64Attribute{
				Description: "Size in GB for the instance template flavor.",
				Computed:    true,
			},
			"instance_template_flavor_vcpu": schema.Int64Attribute{
				Description: "Number of vCPUs for the instance template flavor.",
				Computed:    true,
			},
			"labels": schema.MapAttribute{
				Description: "Labels attached to the nodepool.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"security_groups": schema.ListAttribute{
				Description: "List of security groups.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"status_state": schema.StringAttribute{
				Description: "Current state of the nodepool.",
				Computed:    true,
			},
			"status_messages": schema.ListAttribute{
				Description: "Status messages.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"tags": schema.ListAttribute{
				Description: "List of tags.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"taints": schema.ListNestedAttribute{
				Description: "List of taints.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"effect": schema.StringAttribute{
							Description: "Taint effect.",
							Computed:    true,
						},
						"key": schema.StringAttribute{
							Description: "Taint key.",
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "Taint value.",
							Computed:    true,
						},
					},
				},
			},
			"zone": schema.ListAttribute{
				Description: "List of zones.",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (r *DataSourceKubernetesNodepool) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FlattenedGetResult

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOutput, err := r.nodepool.Get(sdkNodepool.GetParameters{
		ClusterId:  data.ClusterID.ValueString(),
		NodePoolId: data.ID.ValueString(),
	}, tfutil.GetConfigsFromTags(r.sdkClient.Sdk().Config().Get, sdkNodepool.GetConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get nodepool", err.Error())
		return
	}

	flattened, err := ConvertGetResultToFlattened(ctx, &sdkOutput, data.ClusterID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to convert nodepool", err.Error())
		return
	}

	flattened.ClusterID = data.ClusterID
	flattened.ID = data.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func ConvertGetResultToFlattened(ctx context.Context, original *sdkNodepool.GetResult, clusterID string) (*FlattenedGetResult, error) {
	if original == nil {
		return nil, nil
	}

	flattened := &FlattenedGetResult{
		ID:                         types.StringValue(original.Id),
		ClusterID:                  types.StringValue(clusterID),
		Name:                       types.StringValue(original.Name),
		CreatedAt:                  types.StringPointerValue(original.CreatedAt),
		UpdatedAt:                  types.StringPointerValue(original.UpdatedAt),
		Replicas:                   types.Int64Value(int64(original.Replicas)),
		InstanceTemplateDiskSize:   types.Int64Value(int64(original.InstanceTemplate.DiskSize)),
		InstanceTemplateDiskType:   types.StringValue(original.InstanceTemplate.DiskType),
		InstanceTemplateNodeImage:  types.StringValue(original.InstanceTemplate.NodeImage),
		InstanceTemplateFlavorID:   types.StringValue(original.InstanceTemplate.Flavor.Id),
		InstanceTemplateFlavorName: types.StringValue(original.InstanceTemplate.Flavor.Name),
		InstanceTemplateFlavorRam:  types.Int64Value(int64(original.InstanceTemplate.Flavor.Ram)),
		InstanceTemplateFlavorSize: types.Int64Value(int64(original.InstanceTemplate.Flavor.Size)),
		InstanceTemplateFlavorVcpu: types.Int64Value(int64(original.InstanceTemplate.Flavor.Vcpu)),
		StatusState:                types.StringValue(original.Status.State),
	}

	flattened.AutoScaleMaxReplicas = types.Int64Value(int64(original.AutoScale.MaxReplicas))
	flattened.AutoScaleMinReplicas = types.Int64Value(int64(original.AutoScale.MinReplicas))

	labelsMap, _ := types.MapValueFrom(ctx, types.StringType, original.Labels)
	flattened.Labels = labelsMap

	if original.SecurityGroups != nil {
		flattened.SecurityGroups = make([]types.String, len(*original.SecurityGroups))
		for i, sg := range *original.SecurityGroups {
			strVal := types.StringValue(sg)
			flattened.SecurityGroups[i] = strVal
		}
	}

	if original.Status.Messages != nil {
		flattened.StatusMessages = make([]types.String, len(original.Status.Messages))
		for i, msg := range original.Status.Messages {
			strVal := types.StringValue(msg)
			flattened.StatusMessages[i] = strVal
		}
	}

	if original.Tags != nil {
		tags, _ := types.ListValueFrom(ctx, types.StringType, *original.Tags)
		flattened.Tags = tags
	}

	if original.Taints != nil {
		flattened.Taints = make([]tfutil.Taint, len(*original.Taints))
		for i, taint := range *original.Taints {
			flattened.Taints[i] = tfutil.Taint{
				Effect: types.StringValue(taint.Effect),
				Key:    types.StringValue(taint.Key),
				Value:  types.StringValue(taint.Value),
			}
		}
	}

	if original.Zone != nil {
		zone, _ := types.ListValueFrom(ctx, types.StringType, *original.Zone)
		flattened.Zone = zone
	}

	return flattened, nil
}
