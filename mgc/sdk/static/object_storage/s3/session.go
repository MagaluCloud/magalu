package s3

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"magalu.cloud/core"
	"magalu.cloud/core/auth"
	corehttp "magalu.cloud/core/http"
)

var excludedHeaders = HeaderMap{
	"Authorization":         nil,
	"Accept-Encoding":       nil,
	"Amz-Sdk-Invocation-Id": nil,
	"Amz-Sdk-Request":       nil,
	"User-Agent":            nil,
	"X-Amzn-Trace-Id":       nil,
	"Expect":                nil,
	"Content-Length":        nil,
}

func BuildHost(cfg Config) string {
	if cfg.ServerUrl != "" {
		return cfg.ServerUrl
	}
	return strings.ReplaceAll(templateUrl, "{{region}}", cfg.Region)
}

func SendRequest[T core.Value](ctx context.Context, req *http.Request) (result T, res *http.Response, err error) {
	httpClient := corehttp.ClientFromContext(ctx)
	if httpClient == nil {
		err = fmt.Errorf("couldn't get http client from context")
		return
	}

	var unsignedPayload bool
	if req.Method == http.MethodPut {
		unsignedPayload = true
	}

	accesskeyId, accessSecretKey := auth.FromContext(ctx).AccessKeyPair()
	if accesskeyId == "" || accessSecretKey == "" {
		err = fmt.Errorf("access key not set, see how to set it with \"auth set -h\"")
		return
	}

	if err = sign(req, accesskeyId, accessSecretKey, unsignedPayload, excludedHeaders); err != nil {
		return
	}

	res, err = httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("error to send HTTP request: %w", err)
		return
	}

	result, err = corehttp.UnwrapResponse[T](res)
	return
}
