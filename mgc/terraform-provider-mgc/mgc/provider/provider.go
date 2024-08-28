package provider

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	datasources "magalu.cloud/terraform-provider-mgc/mgc/datasources"
	resources "magalu.cloud/terraform-provider-mgc/mgc/resources"

	mgcSdk "magalu.cloud/sdk"
)

var _ provider.Provider = (*MgcProvider)(nil)

const providerTypeName = "mgc"

var ignoredTFModules = []string{"profile"}

type MgcProvider struct {
	version string
	commit  string
	date    string
	sdk     *mgcSdk.Sdk
}

type KeyPair struct {
	KeyID     types.String `tfsdk:"key_id"`
	KeySecret types.String `tfsdk:"key_secret"`
}

type ObjectStorageConfig struct {
	ObjectKeyPair *KeyPair `tfsdk:"key_pair"`
}

type ProviderConfig struct {
	Region        types.String         `tfsdk:"region"`
	Env           types.String         `tfsdk:"env"`
	ApiKey        types.String         `tfsdk:"api_key"`
	ObjectStorage *ObjectStorageConfig `tfsdk:"object_storage"`
}

type MgcApiKey struct {
	ApiKey string
}

func (m *MgcApiKey) GetAPIKey() string {
	return m.ApiKey
}

func (p *MgcProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	tflog.Debug(ctx, "setting provider metadata")
	resp.TypeName = providerTypeName
	resp.Version = p.version
}

func (p *MgcProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	tflog.Debug(ctx, "setting provider schema")

	schemaApiKey := schema.SingleNestedAttribute{
		MarkdownDescription: "Specific Bucket Key Pair configuration",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"key_id": schema.StringAttribute{
				MarkdownDescription: "API Key ID\nOptionally you can set the environment variable MGC_OBJ_KEY_ID to override this value.",
				Required:            true,
			},
			"key_secret": schema.StringAttribute{
				MarkdownDescription: "API Key Secret\nOptionally you can set the environment variable MGC_OBJ_KEY_SECRET to override this value.",
				Required:            true,
			},
		},
	}

	schemaObjectStorage := schema.SingleNestedAttribute{
		MarkdownDescription: "Specific Object Storage configuration",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"key_pair": schemaApiKey,
		},
	}

	resp.Schema = schema.Schema{
		Description: "Terraform Provider for Magalu Cloud",
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				MarkdownDescription: "Region. Options: br-ne1 / br-se1\nDefault is br-se1.\nOptionally you can set the environment variable MGC_REGION to override this value.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("br-ne1", "br-se1", "br-mgl1"),
				},
			},
			"env": schema.StringAttribute{
				MarkdownDescription: "Environment. Options: prod / pre-prod\nDefault is prod.\nOptionally you can set the environment variable MGC_ENV to override this value.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("prod", "pre-prod"),
				},
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Magalu API Key for authentication\nOptionally you can set the environment variable MGC_API_KEY to override this value.",
				Optional:            true,
			},
			"object_storage": schemaObjectStorage,
		},
	}

}

func (p *MgcProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "configuring MGC provider")

	var data ProviderConfig

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "fail to get configs from provider")
	}

	if !data.Region.IsNull() || os.Getenv("MGC_REGION") != "" {
		region := os.Getenv("MGC_REGION")
		if region == "" && !data.Region.IsNull() {
			region = data.Region.ValueString()
		}
		if err := p.sdk.Config().SetTempConfig("region", region); err != nil {
			tflog.Error(ctx, "fail to set region")
		}
	} else {
		if err := p.sdk.Config().SetTempConfig("region", "br-se1"); err != nil {
			tflog.Error(ctx, "fail to set region")
		}
	}
	if !data.Env.IsNull() {
		env := os.Getenv("MGC_ENV")
		if env == "" && !data.Env.IsNull() {
			env = data.Env.String()
		}
		if err := p.sdk.Config().SetTempConfig("env", env); err != nil {
			tflog.Error(ctx, "fail to set env")
		}
	}
	if !data.ApiKey.IsNull() || os.Getenv("MGC_API_KEY") != "" {
		apiKey := os.Getenv("MGC_API_KEY")

		if apiKey == "" && !data.ApiKey.IsNull() {
			apiKey = data.ApiKey.ValueString()
		}

		err := p.sdk.Auth().SetAPIKey(apiKey)
		if err != nil {
			tflog.Error(ctx, "fail to set api key")
		}
	}

	keyId := ""
	keySecret := ""
	if os.Getenv("MGC_OBJ_KEY_ID") != "" && os.Getenv("MGC_OBJ_KEY_SECRET") != "" {
		keyId = os.Getenv("MGC_OBJ_KEY_ID")
		keySecret = os.Getenv("MGC_OBJ_KEY_SECRET")
	}

	if keyId == "" &&
		keySecret == "" &&
		data.ObjectStorage != nil &&
		data.ObjectStorage.ObjectKeyPair != nil &&
		!data.ObjectStorage.ObjectKeyPair.KeyID.IsNull() &&
		!data.ObjectStorage.ObjectKeyPair.KeySecret.IsNull() {
		keyId = data.ObjectStorage.ObjectKeyPair.KeyID.ValueString()
		keySecret = data.ObjectStorage.ObjectKeyPair.KeySecret.ValueString()
	}
	if keyId != "" && keySecret != "" {
		p.sdk.Config().AddTempKeyPair("apikey",
			keyId,
			keySecret,
		)
		tflog.Debug(ctx, "setting object storage key pair")
	}
	resp.DataSourceData = p.sdk
	resp.ResourceData = p.sdk
}

func (p *MgcProvider) Resources(ctx context.Context) []func() resource.Resource {
	tflog.Info(ctx, "configuring MGC provider resources")

	root := p.sdk.Group()
	rsrc, err := collectGroupResources(ctx, p.sdk, root, []string{providerTypeName})

	rsrc = append(rsrc,
		resources.NewNewNodePoolResource,
		resources.NewK8sClusterResource,
		resources.NewObjectStorageBucketsResource,
		resources.NewVirtualMachineInstancesResource,
		resources.NewVirtualMachineSnapshotsResource,
		resources.NewVolumeAttachResource,
		resources.NewBlockStorageSnapshotsResource,
		resources.NewBlockStorageVolumesResource,
		resources.NewSshKeysResource,
	)

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("An error occurred while generating the provider resource list: %v", err))
	}

	return rsrc
}

func collectGroupResources(
	ctx context.Context,
	sdk *mgcSdk.Sdk,
	group mgcSdk.Grouper,
	path []string,
) ([]func() resource.Resource, error) {
	debugMap := map[string]any{"path": path}
	tflog.Debug(ctx, "Collecting resources", debugMap)
	var resources []func() resource.Resource
	var create, read, update, delete, list mgcSdk.Executor
	var connectionExecs []mgcSdk.Executor
	_, err := group.VisitChildren(func(child mgcSdk.Descriptor) (run bool, err error) {
		if childGroup, ok := child.(mgcSdk.Grouper); ok {
			if slices.Contains(ignoredTFModules, childGroup.Name()) {
				return true, nil
			}

			childResources, err := collectGroupResources(ctx, sdk, childGroup, append(path, childGroup.Name()))
			if err != nil {
				return false, err
			}

			resources = append(resources, childResources...)
			return true, err
		} else if exec, ok := child.(mgcSdk.Executor); ok {
			switch exec.Name() {
			case "create":
				tflog.Debug(ctx, "found create operation", debugMap)
				create = exec
			case "get":
				tflog.Debug(ctx, "found get/read operation", debugMap)
				read = exec
			case "update":
				tflog.Debug(ctx, "found update operation", debugMap)
				update = exec
			case "delete":
				tflog.Debug(ctx, "found delete operation", debugMap)
				delete = exec
			case "list":
				tflog.Debug(ctx, "found list operation", debugMap)
				list = exec
			default:
				connectionExecs = append(connectionExecs, exec)
			}
			return true, nil
		} else {
			return false, fmt.Errorf("Unsupported descriptor child type %T", child)
		}
	})
	if err != nil || create == nil {
		return resources, err
	}

	strResourceName := strings.Join(path, "_")
	strResourceName = strings.Replace(strResourceName, "-", "_", -1)

	ignoredTFModules := []string{
		"mgc_kubernetes_nodepool",
		"mgc_object_storage_buckets",
		"mgc_virtual_machine_instances",
		"mgc_virtual_machine_snapshots",
		"mgc_block_storage_volume_attachment",
		"mgc_block_storage_snapshots",
		"mgc_block_storage_volumes",
		"mgc_kubernetes_cluster",
		"mgc_ssh_public_keys",
		"mgc_workspace",
		"mgc_virtual_machine_backups",
	}

	if slices.Contains(ignoredTFModules, strResourceName) {
		tflog.Debug(ctx, fmt.Sprintf("resource %q is ignored", strResourceName), debugMap)
		return resources, nil
	}

	resourceName := tfName(strResourceName)

	tflog.Debug(ctx, fmt.Sprintf("found resource %q", resourceName), debugMap)

	res, err := newMgcResource(ctx, sdk, resourceName, mgcName(group.Name()), group.Description(), create, read, update, delete, list)
	if err != nil {
		tflog.Warn(ctx, err.Error(), debugMap)
		return resources, nil
	}

	tflog.Debug(ctx, fmt.Sprintf("export resource %q", resourceName), debugMap)
	resources = append(resources, func() resource.Resource { return res })

	for _, connectionCreate := range connectionExecs {
		connectionPath := append(path, connectionCreate.Name())
		name := tfName(strings.Join(connectionPath, "_"))
		if strings.Contains(string(name), "get") {
			tflog.Debug(ctx, fmt.Sprintf("connection creation %s is a non-modifying action, it can't be turned into a resource", name))
			continue
		}

		connectionResource, err := newMgcConnectionResource(ctx, sdk, name, connectionCreate.Description(), connectionCreate, delete)
		if err != nil {
			tflog.Warn(ctx, err.Error(), debugMap)
			continue
		}

		resources = append(resources, func() resource.Resource { return connectionResource })
		tflog.Debug(ctx, fmt.Sprintf("export connection resource %q", resourceName), debugMap)
	}

	return resources, err
}

func (p *MgcProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	tflog.Info(ctx, "configuring MGC provider data sources")

	var dataSources []func() datasource.DataSource
	dataSources = append(dataSources,
		datasources.NewDataSourceKubernetesClusterKubeConfig,
		datasources.NewDataSourceKubernetesCluster,
		datasources.NewDataSourceKubernetesFlavor,
		datasources.NewDataSourceKubernetesVersion,
		datasources.NewDataSourceKubernetesNodepool,
		datasources.NewDataSourceKubernetesNode,
		// datasources.NewDataSourceSSH,
		datasources.NewDataSourceVmMachineType,
		datasources.NewDataSourceVMIMages,
		datasources.NewDataSourceVmInstance,
		datasources.NewDataSourceVmInstances,
		// datasources.NewDatasourceBucket,
		// datasources.NewDatasourceBuckets,
	)

	return dataSources
}

func New(version string, commit string, date string) func() provider.Provider {
	sdk := mgcSdk.NewSdk()
	mgcSdk.SetUserAgent("MgcTF")

	return func() provider.Provider {
		return &MgcProvider{
			sdk:     sdk,
			version: version,
			commit:  commit,
			date:    date,
		}
	}
}
