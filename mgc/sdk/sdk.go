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

// TODO: Change config with build tags or from environment
var config core.AuthConfig = core.AuthConfig{
	ClientId:      "cw9qpaUl2nBiC8PVjNFN5jZeb2vTd_1S5cYs1FhEXh0",
	RedirectUri:   "http://localhost:8095/callback",
	LoginUrl:      "https://id.magalu.com/login",
	TokenUrl:      "https://id.magalu.com/oauth/token",
	ValidationUrl: "https://id.magalu.com/oauth/introspect",
	RefreshUrl:    "https://id.magalu.com/oauth/token",
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
func (o *Sdk) NewContext() context.Context {
	var ctx = context.Background()
	ctx = core.NewGrouperContext(ctx, o.Group())
	ctx = core.NewAuthContext(ctx, o.Auth())
	ctx = core.NewHttpClientContext(ctx, o.HttpClient())
	ctx = core.NewConfigContext(ctx, o.Config())
	return ctx
}

func (o *Sdk) newOpenApiSource() *openapi.Source {
	// TODO: are these going to be fixed? configurable?
	extensionPrefix := "x-cli"
	openApiDir := os.Getenv("MGC_SDK_OPENAPI_DIR")
	if openApiDir == "" {
		cwd, err := os.Getwd()
		if err == nil {
			openApiDir = filepath.Join(cwd, "openapis")
		}
	}

	return &openapi.Source{
		Loader: openapi.FileLoader{
			Dir: openApiDir,
		},
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

func (o *Sdk) Auth() *core.Auth {
	if o.auth == nil {
		client := &http.Client{Transport: newHttpTransport()}
		o.auth = core.NewAuth(config, client)
	}
	return o.auth
}

func (o *Sdk) HttpClient() *core.HttpClient {
	if o.httpClient == nil {
		transport := newHttpTransport()
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
