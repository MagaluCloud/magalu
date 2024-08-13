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

var _ datasource.DataSource = &DataSourceKubernetesClusterKubeConfig{}

func NewDataSourceKubernetesClusterKubeConfig() datasource.DataSource {
	return &DataSourceKubernetesClusterKubeConfig{}
}

type DataSourceKubernetesClusterKubeConfig struct {
	sdkClient *mgcSdk.Client
	cluster   cluster.Service
}

type DataSourceKubernetesClusterKubeConfigModel struct {
	ClusterID types.String `tfsdk:"cluster_id"`
	RawConfig types.String `tfsdk:"kubeconfig"`
}

func (d *DataSourceKubernetesClusterKubeConfig) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_cluster_kubeconfig"
}

func (d *DataSourceKubernetesClusterKubeConfig) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the Kubernetes cluster.",
			},
			"kubeconfig": schema.StringAttribute{
				Computed:    true,
				Description: "The full contents of the Kubernetes cluster's kubeconfig yaml file.",
			},
		},
	}
	resp.Schema.Description = "Get the kubeconfig of a Kubernetes cluster by cluster_id."
}

func (d *DataSourceKubernetesClusterKubeConfig) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DataSourceKubernetesClusterKubeConfigModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	sdkOuput, err := d.cluster.Kubeconfig(cluster.KubeconfigParameters{
		ClusterId: data.ClusterID.ValueString(),
	}, GetConfigsFromTags(d.sdkClient.Sdk().Config().Get, cluster.KubeconfigConfigs{}))
	if err != nil {
		resp.Diagnostics.AddError("Failed to get kubeconfig", err.Error())
		return
	}

	data.RawConfig = types.StringValue(sdkOuput)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DataSourceKubernetesClusterKubeConfig) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
