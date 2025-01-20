package mgc

import (
	"context"
	"os"

	datasources "github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/datasources"
	resources "github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/resources"
	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc/tfutil"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
)

const providerTypeName = "mgc"

type mgcProvider struct {
	version string
	commit  string
	date    string
	sdk     *mgcSdk.Sdk
}

func (p *mgcProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = providerTypeName
	resp.Version = p.version
}

func (p *mgcProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform Provider for Magalu Cloud",
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Description: "The region to use for resources. Options: br-ne1 / br-se1. Default is br-se1. Can be set with environment variable MGC_REGION.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("br-ne1", "br-se1", "br-mgl1"),
				},
			},
			"env": schema.StringAttribute{
				Description: "The environment to use. Options: prod / pre-prod. Default is prod. Can be set with environment variable MGC_ENV.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("prod", "pre-prod"),
				},
			},
			"api_key": schema.StringAttribute{
				Description: "The Magalu API Key for authentication. Can be set with environment variable MGC_API_KEY.",
				Optional:    true,
			},
			"object_storage": schema.SingleNestedAttribute{
				Description: "Configuration settings for Object Storage",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"key_pair": schema.SingleNestedAttribute{
						Description: "Bucket Key Pair configuration",
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"key_id": schema.StringAttribute{
								Description: "The API Key ID. Can be set with environment variable MGC_OBJ_KEY_ID.",
								Required:    true,
							},
							"key_secret": schema.StringAttribute{
								Description: "The API Key Secret. Can be set with environment variable MGC_OBJ_KEY_SECRET.",
								Required:    true,
							},
						},
					},
				},
			},
		},
	}
}

func (p *mgcProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data tfutil.ProviderConfig
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "fail to get configs from provider")
	}

	if data.ApiKey.ValueString() == "" {
		if apiKeyFromOS := os.Getenv("MGC_API_KEY"); apiKeyFromOS != "" {
			data.ApiKey = types.StringValue(apiKeyFromOS)
		} else {
			data.ApiKey = types.StringValue("")
		}
	}

	if data.Env.ValueString() == "" {
		if envFromOS := os.Getenv("MGC_ENV"); envFromOS != "" {
			data.Env = types.StringValue(envFromOS)
		} else {
			data.Env = types.StringValue("prod")
		}
	}

	if data.Region.ValueString() == "" {
		if regionFromOS := os.Getenv("MGC_REGION"); regionFromOS != "" {
			data.Region = types.StringValue(regionFromOS)
		} else {
			data.Region = types.StringValue("br-se1")
		}
	}

	if data.ObjectStorage == nil || (os.Getenv("MGC_OBJ_KEY_ID") != "" && os.Getenv("MGC_OBJ_KEY_SECRET") != "") {
		data.ObjectStorage = &tfutil.ObjectStorageConfig{
			ObjectKeyPair: &tfutil.KeyPair{
				KeyID:     types.StringValue(os.Getenv("MGC_OBJ_KEY_ID")),
				KeySecret: types.StringValue(os.Getenv("MGC_OBJ_KEY_SECRET")),
			},
		}
	}

	resp.DataSourceData = data
	resp.ResourceData = data
}

func (p *mgcProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewNewNodePoolResource,
		resources.NewK8sClusterResource,
		resources.NewObjectStorageBucketsResource,
		resources.NewVirtualMachineInstancesResource,
		resources.NewVirtualMachineSnapshotsResource,
		resources.NewVolumeAttachResource,
		resources.NewBlockStorageSnapshotsResource,
		resources.NewBlockStorageVolumesResource,
		resources.NewSshKeysResource,
		resources.NewNetworkPublicIPResource,
		resources.NewNetworkVPCResource,
		resources.NewNetworkSecurityGroupsResource,
		resources.NewNetworkSecurityGroupsRulesResource,
		resources.NewNetworkVPCInterfaceResource,
		resources.NewNetworkSecurityGroupsAttachResource,
		resources.NewNetworkPublicIPAttachResource,
		resources.NewNetworkVpcsSubnetsResource,
		resources.NewNetworkSubnetpoolsResource,
		resources.NewDBaaSInstanceResource,
		resources.NewDBaaSInstanceSnapshotResource,
		resources.NewContainerRegistryRegistriesResource,
		resources.NewVirtualMachineInterfaceAttachResource,
	}
}

func (p *mgcProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewDataSourceKubernetesClusterKubeConfig,
		datasources.NewDataSourceKubernetesCluster,
		datasources.NewDataSourceKubernetesFlavor,
		datasources.NewDataSourceKubernetesVersion,
		datasources.NewDataSourceKubernetesNodepool,
		datasources.NewDataSourceKubernetesNode,
		datasources.NewDataSourceVmMachineType,
		datasources.NewDataSourceVMIMages,
		datasources.NewDataSourceVmInstance,
		datasources.NewDataSourceVmInstances,
		datasources.NewDataSourceNetworkVPC,
		datasources.NewDataSourceNetworkSecurityGroup,
		datasources.NewDataSourceNetworkVPCInterface,
		datasources.NewDataSourceNetworkVpcsSubnet,
		datasources.NewDataSourceNetworkSubnetpool,
		datasources.NewDataSourceNetworkPublicIP,
		datasources.NewDataSourceBSSnapshots,
		datasources.NewDataSourceBsVolumes,
		datasources.NewDataSourceBSSnapshot,
		datasources.NewDataSourceBsVolume,
		datasources.NewDataSourceBsVolumeTypes,
		datasources.NewDataSourceDbaasEngines,
		datasources.NewDataSourceDbaasInstanceTypes,
		datasources.NewDataSourceDbaasInstances,
		datasources.NewDataSourceDbaasInstance,
		datasources.NewDataSourceDbaasInstancesSnapshots,
		datasources.NewDataSourceDbaasInstancesSnapshot,
		datasources.NewDataSourceAvailabilityZones,
		datasources.NewDataSourceCRImages,
		datasources.NewDataSourceCRRepositories,
		datasources.NewDataSourceCRRegistries,
		datasources.NewDataSourceCRCredentials,
	}
}

func New(version string, commit string, date string) func() provider.Provider {
	sdk := mgcSdk.NewSdk()
	mgcSdk.SetUserAgent("MgcTF")

	return func() provider.Provider {
		return &mgcProvider{
			sdk:     sdk,
			version: version,
			commit:  commit,
			date:    date,
		}
	}
}
