package sdk

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"magalu.cloud/core"
	"magalu.cloud/core/auth"
	"magalu.cloud/core/config"
	"magalu.cloud/core/dataloader"
	mgcHttpPkg "magalu.cloud/core/http"
	"magalu.cloud/core/profile_manager"
	"magalu.cloud/sdk/blueprint"
	"magalu.cloud/sdk/openapi"
	"magalu.cloud/sdk/static"
)

// Re-exports from Core
type Descriptor = core.Descriptor
type DescriptorVisitor = core.DescriptorVisitor
type Example = core.Example
type Executor = core.Executor
type Grouper = core.Grouper
type Linker = core.Linker
type Schema = core.Schema
type Value = core.Value
type Result = core.Result
type Config = config.Config

type Sdk struct {
	group          *core.MergeGroup
	profileManager *profile_manager.ProfileManager
	auth           *auth.Auth
	httpClient     *mgcHttpPkg.Client
	config         *config.Config
	refResolver    core.RefPathResolver
}

type contextKey string

var ctxWrappedKey contextKey = "magalu.cloud/sdk/SdkWrapped"

func NewSdk() *Sdk {
	return &Sdk{}
}

// The Context is created with the following values:
// - use GrouperFromContext() to retrieve Sdk.Group() (root group)
// - use AuthFromContext() to retrieve Sdk.Auth()
// - use HttpClientFromContext() to retrieve Sdk.HttpClient()
// - use ConfigFromContext() to retrieve Sdk.Config()
func (o *Sdk) NewContext() context.Context {
	var ctx = context.Background()
	return o.WrapContext(ctx)
}

// The following values are added to the context:
// - use RefPathResolverFromContext() to retrieve root Sdk.RefResolver()
// - use GrouperFromContext() to retrieve Sdk.Group() (root group)
// - use AuthFromContext() to retrieve Sdk.Auth()
// - use HttpClientFromContext() to retrieve Sdk.HttpClient()
// - use ConfigFromContext() to retrieve Sdk.Config()
func (o *Sdk) WrapContext(ctx context.Context) context.Context {
	if wrap := ctx.Value(ctxWrappedKey); wrap != nil {
		return ctx
	}

	ctx = core.NewRefPathResolverContext(ctx, o.RefResolver())
	ctx = core.NewGrouperContext(ctx, o.Group)
	ctx = profile_manager.NewContext(ctx, o.ProfileManager())
	ctx = config.NewContext(ctx, o.Config())
	ctx = auth.NewContext(ctx, o.Auth())
	// Needs to be called after Auth, because we need the refresh token callback for the interceptor
	ctx = mgcHttpPkg.NewClientContext(ctx, o.HttpClient())
	ctx = context.WithValue(ctx, ctxWrappedKey, true)
	return ctx
}

func (o *Sdk) newOpenApiSource() core.Grouper {
	embedLoader := openapi.GetEmbedLoader()

	// TODO: are these going to be fixed? configurable?
	extensionPrefix := "x-mgc"
	openApiDir := os.Getenv("MGC_SDK_OPENAPI_DIR")
	if openApiDir == "" {
		cwd, err := os.Getwd()
		if err == nil {
			openApiDir = filepath.Join(cwd, "openapis")
		}
	}
	fileLoader := &dataloader.FileLoader{
		Dir: openApiDir,
	}

	var loader dataloader.Loader
	if embedLoader != nil {
		loader = dataloader.NewMergeLoader(fileLoader, embedLoader)
	} else {
		loader = fileLoader
	}

	return openapi.NewSource(loader, &extensionPrefix)
}

func (o *Sdk) newBlueprintSource(rootRefResolver core.RefPathResolver) core.Grouper {
	embedLoader := blueprint.GetEmbedLoader()

	blueprintsDir := os.Getenv("MGC_SDK_BLUEPRINTS_DIR")
	if blueprintsDir == "" {
		cwd, err := os.Getwd()
		if err == nil {
			blueprintsDir = filepath.Join(cwd, "blueprints")
		}
	}
	fileLoader := &dataloader.FileLoader{
		Dir: blueprintsDir,
	}

	var loader dataloader.Loader
	if embedLoader != nil {
		loader = dataloader.NewMergeLoader(fileLoader, embedLoader)
	} else {
		loader = fileLoader
	}

	return blueprint.NewSource(loader, rootRefResolver)
}

func (o *Sdk) RefResolver() core.RefPathResolver {
	if o.refResolver == nil {
		o.refResolver = core.NewDocumentRefPathResolver(func() (any, error) { return o.group, nil })
	}
	return o.refResolver
}

func (o *Sdk) Group() core.Grouper {
	if o.group == nil {
		o.group = core.NewMergeGroup(
			core.DescriptorSpec{
				Name:        "MagaLu Cloud",
				Version:     "1.0",
				Description: "All MagaLu Groups & Executors",
			},
			func() []core.Grouper {
				return []core.Grouper{
					static.GetGroup(),
					o.newOpenApiSource(),
					o.newBlueprintSource(o.RefResolver()),
				}
			},
		)
	}
	return o.group
}

func newHttpTransport() http.RoundTripper {
	userAgent := "MgcSDK/" + Version

	// To avoid creating a transport with zero values, we leverage
	// DefaultTransport (exemple: `Proxy: ProxyFromEnvironment`)
	transport := mgcHttpPkg.DefaultTransport()
	transport = mgcHttpPkg.NewDefaultClientLogger(transport)
	transport = newDefaultSdkTransport(transport, userAgent)
	return transport
}

func (o *Sdk) addHttpRefreshHandler(t http.RoundTripper) http.RoundTripper {
	return mgcHttpPkg.NewDefaultRefreshLogger(t, o.Auth().RefreshAccessToken)
}

func (o *Sdk) ProfileManager() *profile_manager.ProfileManager {
	if o.profileManager == nil {
		o.profileManager = profile_manager.New()
	}
	return o.profileManager
}

func (o *Sdk) Auth() *auth.Auth {
	if o.auth == nil {
		client := &http.Client{Transport: newHttpTransport()}
		o.auth = auth.New(authConfigMap, client, o.ProfileManager(), o.Config())
	}

	return o.auth
}

func (o *Sdk) HttpClient() *mgcHttpPkg.Client {
	if o.httpClient == nil {
		transport := o.addHttpRefreshHandler(newHttpTransport())
		o.httpClient = mgcHttpPkg.NewClient(transport)
	}
	return o.httpClient
}

func (o *Sdk) Config() *config.Config {
	if o.config == nil {
		o.config = config.New(o.ProfileManager())
	}
	return o.config
}

var authConfigMap map[string]auth.Config

func init() {
	authConfigMap = map[string]auth.Config{
		"prod": {
			ClientId:            "cw9qpaUl2nBiC8PVjNFN5jZeb2vTd_1S5cYs1FhEXh0",
			ObjectStoreScopeIDs: []string{"b6afac7e-0afd-42de-b4aa-1bc82a27e307", "5ea6d1f7-20eb-4e80-9a9c-c7923636a4bd"},
			RedirectUri:         "http://localhost:8095/callback",
			LoginUrl:            "https://id.magalu.com/login",
			TokenUrl:            "https://id.magalu.com/oauth/token",
			ValidationUrl:       "https://id.magalu.com/oauth/introspect",
			RefreshUrl:          "https://id.magalu.com/oauth/token",
			TenantsListUrl:      "https://id.magalu.com/account/api/v2/whoami/tenants",
			TokenExchangeUrl:    "https://id.magalu.com/oauth/token/exchange",
			ApiKeysUrlV1:        "https://id.magalu.com/account/api/v1/api-keys",
			ApiKeysUrlV2:        "https://id.magalu.com/account/api/v2/api-keys",
		},
		"pre-prod": { // TODO update this links to the correct ones
			ClientId:            "dByqQVtHcs07b_O9jpUDgfV5UCskh9TbC64WUXEdVHE",
			ObjectStoreScopeIDs: []string{"b6afac7e-0afd-42de-b4aa-1bc82a27e307", "5ea6d1f7-20eb-4e80-9a9c-c7923636a4bd"},
			RedirectUri:         "http://localhost:8095/callback",
			LoginUrl:            "https://idmagalu-preprod.luizalabs.com/login",
			TokenUrl:            "https://idpa-api-preprod.luizalabs.com/oauth/token",
			ValidationUrl:       "https://idpa-api-preprod.luizalabs.com/oauth/introspect",
			RefreshUrl:          "https://idpa-api-preprod.luizalabs.com/oauth/token",
			TenantsListUrl:      "https://platform-account-api-preprod.luizalabs.com/api/v2/whoami/tenants",
			TokenExchangeUrl:    "https://idpa-api-preprod.luizalabs.com/oauth/token/exchange",
			ApiKeysUrlV1:        "https://idpa-api-preprod.luizalabs.com/account/api/v1/api-keys",
			ApiKeysUrlV2:        "https://idpa-api-preprod.luizalabs.com/account/api/v2/api-keys",
		},
	}
	authConfigMap["default"] = authConfigMap["prod"]
}
