package provider

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
				MarkdownDescription: "API Key ID",
				Required:            true,
			},
			"key_secret": schema.StringAttribute{
				MarkdownDescription: "API Key Secret",
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
				MarkdownDescription: "Region",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Magalu API Key for authentication",
				Optional:            true,
			},
			"object_storage": schemaObjectStorage,
		},
	}

}

var acceptedRegions = []string{"br-ne1", "br-se1", "br-mgl1"}

func (p *MgcProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "configuring MGC provider")

	var data ProviderConfig

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "fail to get configs from provider")
	}

	if !data.Region.IsNull() {
		if !slices.Contains(acceptedRegions, data.Region.ValueString()) {
			tflog.Error(ctx, "invalid region. Valid options: "+strings.Join(acceptedRegions, ", "))
		}
		if err := p.sdk.Config().SetTempConfig("region", data.Region.String()); err != nil {
			tflog.Error(ctx, "fail to set region")
		}
	}

	if !data.ApiKey.IsNull() {
		MgcApiKey := &MgcApiKey{data.ApiKey.ValueString()}
		err := p.sdk.Auth().SetAPIKey(MgcApiKey)
		if err != nil {
			tflog.Error(ctx, "fail to set api key")
		}
	}

	if data.ObjectStorage != nil && data.ObjectStorage.ObjectKeyPair != nil &&
		!data.ObjectStorage.ObjectKeyPair.KeyID.IsNull() &&
		!data.ObjectStorage.ObjectKeyPair.KeySecret.IsNull() {
		p.sdk.Config().AddTempKeyPair("apikey",
			data.ObjectStorage.ObjectKeyPair.KeyID.ValueString(),
			data.ObjectStorage.ObjectKeyPair.KeySecret.ValueString(),
		)
	}
	resp.DataSourceData = p.sdk
	resp.ResourceData = p.sdk
}

func (p *MgcProvider) Resources(ctx context.Context) []func() resource.Resource {
	tflog.Info(ctx, "configuring MGC provider resources")

	root := p.sdk.Group()
	resources, err := collectGroupResources(ctx, p.sdk, root, []string{providerTypeName})
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("An error occurred while generating the provider resource list: %v", err))
	}

	return resources
}

func collectGroupResources(
	ctx context.Context,
	sdk *mgcSdk.Sdk,
	group mgcSdk.Grouper,
	path []string,
) ([]func() resource.Resource, error) {
	// TODO: We should check if the version is correct in the Configuration call or Resource
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
			// TODO: see how this stands in practice
			// some resources have more than one action and we're de-duplicating them,
			// resulting in get-X + get-Y...
			// maybe something to check with scripts/spec_stats.py
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

	resourceName := tfName(strings.Join(path, "_"))
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
	return nil
}

func New(version string, commit string, date string) func() provider.Provider {
	sdk := mgcSdk.NewSdk()

	return func() provider.Provider {
		return &MgcProvider{
			sdk:     sdk,
			version: version,
			commit:  commit,
			date:    date,
		}
	}
}
