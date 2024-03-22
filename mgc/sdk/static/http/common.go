package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"magalu.cloud/core"
	mgcAuthPkg "magalu.cloud/core/auth"
	mgcConfigPkg "magalu.cloud/core/config"
	mgcHttpPkg "magalu.cloud/core/http"
	"magalu.cloud/core/utils"
)

type httpParams struct {
	Url      string                `json:"url" jsonschema:"description=Golang template with the URL"`
	Headers  map[string][]string   `json:"headers,omitempty"`
	Security []map[string][]string `json:"security,omitempty"`
	Method   string                `json:"method" jsonschema:"enum=GET,enum=PUT,enum=POST,enum=DELETE,enum=OPTIONS,enum=HEAD,enum=PATCH,enum=TRACE"`
	Body     string                `json:"body,omitempty"`
}

type httpConfig struct {
	mgcConfigPkg.NetworkConfig `json:",squash"` // nolint

	Region string `json:"region,omitempty" jsonschema:"description=Region to reach the service,default=br-ne-1,enum=br-ne-1,enum=br-se-1,enum=br-mgl-1"`
}

type httpResult struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Headers http.Header `json:"headers"`
	Body    string      `json:"body,omitempty"`
}

func getUrl(params httpParams, configs httpConfig) (url string, err error) {
	t, err := utils.NewTemplate(params.Url)
	if err != nil {
		return
	}

	return utils.ExecuteTemplateTrimmed(t, configs)
}

var httpMethodsBody = map[string]bool{
	"GET":     false,
	"PUT":     true,
	"POST":    true,
	"DELETE":  false,
	"OPTIONS": true,
	"HEAD":    false,
	"PATCH":   true,
	"TRACE":   false,
}

func getBody(params httpParams) io.Reader {
	if httpMethodsBody[params.Method] {
		return strings.NewReader(params.Body)
	}
	return nil
}

func addHeaders(params httpParams, req *http.Request, defaultHeaders map[string]string) {
	for key, values := range params.Headers {
		for _, v := range values {
			req.Header.Add(key, v)
		}
	}
	for key, value := range defaultHeaders {
		if req.Header.Get(key) == "" {
			req.Header.Set(key, value)
		}
	}
}

func addAuth(params httpParams, ctx context.Context, req *http.Request) error {
	for _, item := range params.Security {
		for scheme, scopes := range item {
			switch strings.ToLower(scheme) {
			case "oauth2", "bearerauth":
				auth := mgcAuthPkg.FromContext(ctx)
				accessToken, err := auth.AccessToken(ctx)
				if err != nil {
					return err
				}

				req.Header.Set("Authorization", "Bearer "+accessToken)
				return nil

			default:
				return fmt.Errorf("security scheme %q with scopes %#v is not supported", scheme, scopes)
			}
		}
	}

	return nil
}

var getHttpExecutor = utils.NewLazyLoader[core.Executor](func() core.Executor {
	return core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "do",
			Description: "Execute generic http request",
		},
		executeHttp,
	)
})

func executeHttp(ctx context.Context, params httpParams, configs httpConfig) (result httpResult, err error) {
	return executeHttpWithDefaultHeaders(ctx, params, configs, nil)
}

func executeHttpWithDefaultHeaders(ctx context.Context, params httpParams, configs httpConfig, defaultHeaders map[string]string) (result httpResult, err error) {
	client := mgcHttpPkg.ClientFromContext(ctx)

	params.Method = strings.ToUpper(params.Method)

	url, err := getUrl(params, configs)
	if err != nil {
		err = &core.ChainedError{Name: "url", Err: core.UsageError{Err: err}}
		return
	}

	body := getBody(params)
	req, err := http.NewRequestWithContext(ctx, params.Method, url, body)
	if err != nil {
		err = fmt.Errorf("cannot create HTTP request: %w", err)
		return
	}

	addHeaders(params, req, defaultHeaders)
	err = addAuth(params, ctx, req)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("HTTP request error: %w", err)
		return
	}

	var resultBody []byte
	if resp.Body != nil {
		defer resp.Body.Close()
		resultBody, err = io.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("error reading result body: %w", err)
			return
		}
	}

	return httpResult{
		resp.StatusCode,
		resp.Status,
		resp.Header,
		string(resultBody),
	}, nil
}
