package http

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type jsonParams struct {
	Url      string                `json:"url" jsonschema:"description=Golang template with the URL"`
	Headers  map[string][]string   `json:"headers,omitempty"`
	Security []map[string][]string `json:"security,omitempty"`
}

type jsonBodyParams struct {
	Url      string                `json:"url" jsonschema:"description=Golang template with the URL"`
	Headers  map[string][]string   `json:"headers,omitempty"`
	Security []map[string][]string `json:"security,omitempty"`
	Body     any                   `json:"body,omitempty"`
}

type jsonResult struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Headers http.Header `json:"headers"`
	Body    any         `json:"body,omitempty"`
}

func newJsonExecutor(method string, hasBody bool) core.Executor {
	spec := core.DescriptorSpec{
		Name:        strings.ToLower(method),
		Description: fmt.Sprintf("Call using HTTP %s", method),
	}

	if hasBody {
		return core.NewStaticExecute(spec, newExecuteJsonBody(method))
	}
	return core.NewStaticExecute(spec, newExecuteJson(method))
}

const targetContentType = "application/json"

var defaultJsonHeaders = map[string]string{
	"Content-Type": targetContentType,
	"Accept":       targetContentType,
	"Connection":   "keep-alive",
}

func jsonExecute(ctx context.Context, params httpParams, configs httpConfig) (result jsonResult, err error) {
	r, err := executeHttpWithDefaultHeaders(ctx, params, configs, defaultJsonHeaders)
	if err != nil {
		return
	}

	var resultBody any
	if len(r.Body) > 0 {
		contentType, _, _ := mime.ParseMediaType(r.Headers.Get("Content-Type"))
		err = json.Unmarshal([]byte(r.Body), &resultBody)
		if err != nil {
			err = fmt.Errorf("unable to decode JSON response (Content-Type: %q): %w", contentType, err)
			return
		}
	}

	return jsonResult{r.Code, r.Status, r.Headers, resultBody}, nil
}

func newExecuteJson(method string) func(ctx context.Context, params jsonParams, configs httpConfig) (result jsonResult, err error) {
	return func(ctx context.Context, params jsonParams, configs httpConfig) (result jsonResult, err error) {
		p := httpParams{
			Method:   method,
			Url:      params.Url,
			Headers:  params.Headers,
			Security: params.Security,
		}
		return jsonExecute(ctx, p, configs)
	}
}

func newExecuteJsonBody(method string) func(ctx context.Context, params jsonBodyParams, configs httpConfig) (result jsonResult, err error) {
	return func(ctx context.Context, params jsonBodyParams, configs httpConfig) (result jsonResult, err error) {
		reqBody, err := json.Marshal(params.Body)
		if err != nil {
			err = fmt.Errorf("unable to encode JSON response: %w", err)
			return
		}

		p := httpParams{
			Method:   method,
			Url:      params.Url,
			Headers:  params.Headers,
			Security: params.Security,
			Body:     string(reqBody),
		}

		return jsonExecute(ctx, p, configs)
	}
}

var getJsonGroup = utils.NewLazyLoader(func() core.Grouper {
	executors := []func() core.Descriptor{}
	for k, v := range httpMethodsBody {
		// Do this because of Go's weird loops lol
		// https://go.dev/blog/loopvar-preview
		method := k
		hasBody := v
		executors = append(executors, func() core.Descriptor { return newJsonExecutor(method, hasBody) })
	}

	return core.NewStaticGroup(
		core.DescriptorSpec{
			Name:        "json",
			Description: "JSON HTTP access",
		},
		func() (children []core.Descriptor) {
			for _, exec := range executors {
				children = append(children, exec())
			}
			return children
		},
	)
})
