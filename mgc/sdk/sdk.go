package sdk

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"magalu.cloud/core"
	"magalu.cloud/sdk/openapi"
	"magalu.cloud/sdk/static"
)

// Re-exports from Core
type Descriptor = core.Descriptor
type DescriptorVisitor = core.DescriptorVisitor
type Example = core.Example
type Executor = core.Executor
type Grouper = core.Grouper
type Schema = core.Schema
type Value = core.Value

type Sdk struct {
	group      *core.MergeGroup
	auth       *core.Auth
	httpClient *core.HttpClient
	config     *core.Config
}

type contextKey string

var ctxWrappedKey contextKey = "magalu.cloud/sdk/SdkWrapped"

// TODO: Change config with build tags or from environment
var config core.AuthConfig = core.AuthConfig{
	ClientId:         "cw9qpaUl2nBiC8PVjNFN5jZeb2vTd_1S5cYs1FhEXh0",
	RedirectUri:      "http://localhost:8095/callback",
	LoginUrl:         "https://id.magalu.com/login",
	TokenUrl:         "https://id.magalu.com/oauth/token",
	ValidationUrl:    "https://id.magalu.com/oauth/introspect",
	RefreshUrl:       "https://id.magalu.com/oauth/token",
	TenantsListUrl:   "https://id.magalu.com/account/api/v2/whoami/tenants",
	TenantsSelectUrl: "https://id.magalu.com/oauth/token/exchange",
	Scopes: []string{
		"openid",
		"mke.read",
		"mke.write",
		"network.read",
		"network.write",
		"object-storage.read",
		"object-storage.write",
		"block-storage.read",
		"block-storage.write",
		"virtual-machine.read",
		"virtual-machine.write",
		"dbaas.read",
		"dbaas.write",
		"cpo:read",
		"cpo:write",
		"api-consulta.read",
	},
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
// - use GrouperFromContext() to retrieve Sdk.Group() (root group)
// - use AuthFromContext() to retrieve Sdk.Auth()
// - use HttpClientFromContext() to retrieve Sdk.HttpClient()
// - use ConfigFromContext() to retrieve Sdk.Config()
func (o *Sdk) WrapContext(ctx context.Context) context.Context {
	if wrap := ctx.Value(ctxWrappedKey); wrap != nil {
		return ctx
	}

	ctx = core.NewGrouperContext(ctx, o.Group())
	ctx = core.NewAuthContext(ctx, o.Auth())
	// Needs to be called after Auth, because we need the refresh token callback for the interceptor
	ctx = core.NewHttpClientContext(ctx, o.HttpClient())
	ctx = core.NewConfigContext(ctx, o.Config())

	ctx = context.WithValue(ctx, ctxWrappedKey, true)
	return ctx
}

func (o *Sdk) newOpenApiSource() *openapi.Source {
	embedLoader := openapi.GetEmbedLoader()

	// TODO: are these going to be fixed? configurable?
	extensionPrefix := "x-cli"
	openApiDir := os.Getenv("MGC_SDK_OPENAPI_DIR")
	if openApiDir == "" {
		cwd, err := os.Getwd()
		if err == nil {
			openApiDir = filepath.Join(cwd, "openapis")
		}
	}
	fileLoader := &openapi.FileLoader{
		Dir: openApiDir,
	}

	var loader openapi.Loader
	if embedLoader != nil {
		loader = openapi.NewMergeLoader(fileLoader, embedLoader)
	} else {
		loader = fileLoader
	}

	return &openapi.Source{
		Loader:          loader,
		ExtensionPrefix: &extensionPrefix,
	}
}

func (o *Sdk) Group() core.Grouper {
	if o.group == nil {
		o.group = core.NewMergeGroup(
			"MagaLu Cloud",
			"1.0",
			"All MagaLu Groups & Executors",
			[]core.Grouper{
				static.NewGroup(),
				o.newOpenApiSource(),
			},
		)
	}
	return o.group
}

func newHttpTransport() http.RoundTripper {
	var transport http.RoundTripper = &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	if os.Getenv("MGC_SDK_HTTP_LOG") == "1" {
		transport = core.NewDefaultHttpClientLogger(transport)
	}
	return transport
}

func (o *Sdk) addHttpRefreshHandler(t http.RoundTripper) http.RoundTripper {
	return core.NewDefaultHttpRefreshLogger(t, o.Auth().RefreshAccessToken)
}

func (o *Sdk) Auth() *core.Auth {
	if o.auth == nil {
		client := &http.Client{Transport: newHttpTransport()}
		o.auth = core.NewAuth(config, client)
	}
	return o.auth
}

func (o *Sdk) HttpClient() *core.HttpClient {
	if o.httpClient == nil {
		transport := o.addHttpRefreshHandler(newHttpTransport())
		o.httpClient = core.NewHttpClient(transport)
	}
	return o.httpClient
}

func (o *Sdk) Config() *core.Config {
	if o.config == nil {
		o.config = core.NewConfig()
	}
	return o.config
}
