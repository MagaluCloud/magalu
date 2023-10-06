package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"magalu.cloud/core"
	"magalu.cloud/core/auth"
	mgcSdk "magalu.cloud/sdk"
)

var _ provider.Provider = (*MgcProvider)(nil)

const providerTypeName = "magalu"

type apiSpec struct {
	name    string
	version string
}

func (a apiSpec) String() string {
	return fmt.Sprintf("%s@%s", a.name, a.version)
}

func parseApiSpec(spec string) *apiSpec {
	parts := strings.Split(spec, "@")
	name := parts[0]
	version := parts[1]
	return &apiSpec{name, version}
}

type MgcProvider struct {
	version string
	commit  string
	date    string
	sdk     *mgcSdk.Sdk
	apis    []*apiSpec
}

type MgcProviderModel struct {
	RefreshToken types.String   `tfsdk:"refresh_token"`
	AccessToken  types.String   `tfsdk:"access_token"`
	Apis         []types.String `tfsdk:"apis"`
}

func (p *MgcProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	tflog.Debug(ctx, "setting provider metadata")
	resp.TypeName = providerTypeName
	resp.Version = p.version
}

func (p *MgcProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	tflog.Debug(ctx, "setting provider schema")
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"refresh_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Refresh Token used to authenticate with the MagaluID platform.",
			},
			"access_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Access Token used to authenticate with the MagaluID platform.",
			},
			"apis": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "Which MagaluCloud products need to be supported by the provider.",
				// Should we use validators to know if the apis exist? Or return an error later
				// Validators: ,
			},
		},
	}
}

func (p *MgcProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "configuring MGC provider")
	var config MgcProviderModel

	// Load all configurations
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	refreshToken := config.RefreshToken.ValueString()
	accessToken := config.AccessToken.ValueString()
	if accessToken != "" || refreshToken != "" {
		_ = p.sdk.Auth().SetTokens(&auth.LoginResult{AccessToken: accessToken, RefreshToken: refreshToken})
	}

	if len(config.Apis) == 0 {
		p.attrErr(resp, "apis")
		return
	}
	p.apis = make([]*apiSpec, len(config.Apis))
	for i, spec := range config.Apis {
		p.apis[i] = parseApiSpec(spec.ValueString())
	}
	tflog.Debug(ctx, "using `apis` property from tf file")

	resp.DataSourceData = p.sdk
	resp.ResourceData = p.sdk
}

func (p *MgcProvider) Resources(ctx context.Context) []func() resource.Resource {
	tflog.Info(ctx, "configuring MGC provider resources")

	resources := make([]func() resource.Resource, 0)
	root := p.sdk.Group()
	_, err := root.VisitChildren(func(child core.Descriptor) (run bool, err error) {
		// TODO: We Should we check if the version is correct in the Configuration call or Resource
		childName := fmt.Sprintf("%q@%q", child.Name(), child.Version())
		if group, ok := child.(mgcSdk.Grouper); !ok {
			// Warning since we don't want to stop the discovery process
			tflog.Warn(ctx, fmt.Sprintf("Invalid API %q: invalid format", childName))
			return true, nil
		} else if groupResources, err := collectGroupResources(ctx, p.sdk, group, []string{providerTypeName, group.Name()}); err != nil {
			// Warning since we don't want to stop the discovery process
			tflog.Warn(ctx, fmt.Sprintf("Could not add API %q: %v", childName, err))
			return true, nil
		} else {
			tflog.Info(ctx, fmt.Sprintf("Resources %v", groupResources))
			resources = append(resources, groupResources...)
			return true, nil
		}
	})
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
) (resources []func() resource.Resource, err error) {
	resources = make([]func() resource.Resource, 0)
	var create, read, update, delete mgcSdk.Executor
	_, err = group.VisitChildren(func(child mgcSdk.Descriptor) (run bool, err error) {
		if childGroup, ok := child.(mgcSdk.Grouper); ok {
			oldLen := len(path)
			path = append(path, childGroup.Name())
			childResources, err := collectGroupResources(ctx, sdk, childGroup, path)
			resources = append(resources, childResources...)
			path = path[:oldLen]
			return true, err
		} else if exec, ok := child.(mgcSdk.Executor); ok {
			// TODO: see how this stands in practice
			// some resources have more than one action and we're de-duplicating them,
			// resulting in get-X + get-Y...
			// maybe something to check with scripts/spec_stats.py
			switch exec.Name() {
			case "create":
				create = exec
			case "get":
				read = exec
			case "update":
				update = exec
			case "delete":
				delete = exec
			default:
				actionResource := collectActionResource(ctx, sdk, group, exec, path)
				if len(actionResource) > 0 {
					resources = append(resources, actionResource...)
				} else {
					tflog.Warn(ctx, fmt.Sprintf("TODO: uncovered action %s", exec.Name()))
				}
			}
			return true, nil
		} else {
			return false, fmt.Errorf("unsupported grouper child type %T", child)
		}
	})
	if err != nil {
		return resources, err
	}

	if create != nil {
		name := strings.Join(path, "_")
		if read == nil {
			tflog.Warn(ctx, fmt.Sprintf("Resource %s misses read", name))
			return resources, nil
		}
		if update == nil {
			update = core.NewStaticExecuteSimple("empty update", "1.0", "empty update", func(context context.Context) (any, error) { return nil, nil })
			tflog.Warn(ctx, fmt.Sprintf("Resource %s misses update using empty update", name))
		}
		if delete == nil {
			tflog.Warn(ctx, fmt.Sprintf("Resource %s misses delete", name))
			return resources, nil
		}
		res := &MgcResource{
			sdk:    sdk,
			name:   name,
			group:  group,
			create: create,
			read:   read,
			update: update,
			delete: delete,
		}

		tflog.Debug(ctx, fmt.Sprintf("Export resource %q", name))

		resources = append(resources, func() resource.Resource { return res })
	}

	return resources, err
}

func collectActionResource(ctx context.Context, sdk *mgcSdk.Sdk, group mgcSdk.Grouper, exec mgcSdk.Executor, path []string) (resources [](func() resource.Resource)) {
	resources = make([]func() resource.Resource, 0)

	// Should handle only actions that perform changes
	// TODO: Improve naming, the terraform resource name should be more concise and easier to understand
	path = append(path, exec.Name())
	name := strings.Join(path, "_")
	if strings.Contains(name, "get") {
		tflog.Debug(ctx, fmt.Sprintf("action %s is a non-modifying action, it can't be turned into a resource", name))
		return
	}

	var create, read, update, delete mgcSdk.Executor
	create = exec

	for k, link := range exec.Links() {
		switch k {
		case "read":
			read = link.Target()
		case "update":
			update = link.Target()
		case "delete":
			delete = link.Target()
		}
	}

	if delete == nil {
		tflog.Warn(ctx, fmt.Sprintf("Action Resource %s miss delete operations", name))
		return
	}
	if read == nil {
		read = core.NewStaticExecuteSimple("read", "1.0", "read", func(context context.Context) (any, error) { return nil, nil })
	}
	if update == nil {
		update = core.NewStaticExecuteSimple("update", "1.0", "update", func(context context.Context) (any, error) { return nil, nil })
	}

	res := &MgcResource{
		sdk:    sdk,
		name:   name,
		group:  group,
		create: create,
		read:   read,
		update: update,
		delete: delete,
	}
	resources = append(resources, func() resource.Resource { return res })
	return
}

func (p *MgcProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *MgcProvider) attrErr(resp *provider.ConfigureResponse, attr string, def ...string) {
	exArg := ""
	if len(def) > 0 {
		exArg = fmt.Sprintf("You can set this configuration use one of the following: set %s attribute", attr)
		exArg += strings.Join(def, ", ")
		exArg += "."
	}
	resp.Diagnostics.AddAttributeError(
		path.Root(attr),
		fmt.Sprintf("missing %s", attr),
		fmt.Sprintf("unable to create %s provider due to missing attributes.\n%s", providerTypeName, exArg),
	)
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
