package sdk

import (
	"context"
	"net/http"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/auth"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/MagaluCloud/magalu/mgc/core/dataloader"
	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/profile_manager"
	"github.com/MagaluCloud/magalu/mgc/sdk/openapi"
	"github.com/MagaluCloud/magalu/mgc/sdk/static"
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

var ctxWrappedKey contextKey = "github.com/MagaluCloud/magalu/mgc/sdk/SdkWrapped"

var currentUserAgent string = "MgcSDK"

func SetUserAgent(userAgent string) {
	currentUserAgent = userAgent
}

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
	var loader dataloader.Loader
	if embedLoader != nil {
		loader = dataloader.NewMergeLoader(embedLoader)
	}

	return openapi.NewSource(loader, &extensionPrefix)
}

func (o *Sdk) RefResolver() core.RefPathResolver {
	if o.refResolver == nil {
		o.refResolver = core.NewDocumentRefPathResolver(func() (any, error) { return o.Group(), nil })
	}
	return o.refResolver
}

func (o *Sdk) Group() core.Grouper {
	if o.group == nil {
		o.group = core.NewMergeGroup(
			core.DescriptorSpec{
				Name:        "products",
				Version:     Version,
				Description: "All MagaLu Groups & Executors",
			},
			func() []core.Grouper {
				return []core.Grouper{
					static.GetGroup(),
					o.newOpenApiSource(),
				}
			},
		)
	}
	return o.group
}

func newHttpTransport() http.RoundTripper {
	userAgent := currentUserAgent + "/" + Version

	if strings.HasPrefix(currentUserAgent, "MgcTF") {
		userAgent = currentUserAgent
	}
	// To avoid creating a transport with zero values, we leverage
	// DefaultTransport (exemple: `Proxy: ProxyFromEnvironment`)
	transport := mgcHttpPkg.DefaultTransport()
	transport = mgcHttpPkg.NewDefaultClientLogger(transport)
	transport = newDefaultSdkTransport(transport, userAgent)
	transport = mgcHttpPkg.NewDefaultClientRetryer(transport)
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
			ClientId:              "cw9qpaUl2nBiC8PVjNFN5jZeb2vTd_1S5cYs1FhEXh0",
			ObjectStoreScopeIDs:   []string{"b6afac7e-0afd-42de-b4aa-1bc82a27e307", "5ea6d1f7-20eb-4e80-9a9c-c7923636a4bd"},
			PublicClientsScopeIDs: map[string]string{"openid": "2836b3ba-093c-416a-92f0-7fc4ee5ac961", "profile": "50447cbf-8a42-4426-8e53-fe84bf0726ad"},
			RedirectUri:           "http://localhost:8095/callback",
			LoginUrl:              "https://id.magalu.com/login",
			TokenUrl:              "https://id.magalu.com/oauth/token",
			ValidationUrl:         "https://id.magalu.com/oauth/introspect",
			RefreshUrl:            "https://id.magalu.com/oauth/token",
			TenantsListUrl:        "https://id.magalu.com/account/api/v2/whoami/tenants",
			TokenExchangeUrl:      "https://id.magalu.com/oauth/token/exchange",
			ApiKeysUrlV1:          "https://id.magalu.com/account/api/v1/api-keys",
			ApiKeysUrlV2:          "https://id.magalu.com/account/api/v2/api-keys",
			PublicClientsUrl:      "https://id.magalu.com/account/api/v1/external/clients",
			ClientsV2Url:          "https://id.magalu.com/account/api/v2/clients",
		},
		"pre-prod": { // TODO update this links to the correct ones
			ClientId:              "dByqQVtHcs07b_O9jpUDgfV5UCskh9TbC64WUXEdVHE",
			ObjectStoreScopeIDs:   []string{"b6afac7e-0afd-42de-b4aa-1bc82a27e307", "5ea6d1f7-20eb-4e80-9a9c-c7923636a4bd"},
			PublicClientsScopeIDs: map[string]string{"openid": "4bdb7c8e-6006-478a-ba90-f8313f88bbb8", "profile": "8614f807-9aea-462c-bade-6c08fa52a272"},
			RedirectUri:           "http://localhost:8095/callback",
			LoginUrl:              "https://idmagalu-preprod.luizalabs.com/login",
			TokenUrl:              "https://idpa-api-preprod.luizalabs.com/oauth/token",
			ValidationUrl:         "https://idpa-api-preprod.luizalabs.com/oauth/introspect",
			RefreshUrl:            "https://idpa-api-preprod.luizalabs.com/oauth/token",
			TenantsListUrl:        "https://platform-account-api-preprod.luizalabs.com/api/v2/whoami/tenants",
			TokenExchangeUrl:      "https://idpa-api-preprod.luizalabs.com/oauth/token/exchange",
			ApiKeysUrlV1:          "https://platform-account-api-preprod.luizalabs.com/api/v1/api-keys",
			ApiKeysUrlV2:          "https://platform-account-api-preprod.luizalabs.com/api/v2/api-keys",
			PublicClientsUrl:      "https://platform-account-api-preprod.luizalabs.com/api/v1/external/clients",
			ClientsV2Url:          "https://platform-account-api-preprod.luizalabs.com/api/v2/clients",
		},
	}
	authConfigMap["default"] = authConfigMap["prod"]
}
