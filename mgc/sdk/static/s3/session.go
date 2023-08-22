package s3

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"magalu.cloud/core"
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

func BuildHost(region string) string {
	return strings.ReplaceAll(templateUrl, "{{region}}", region)
}

func SendRequest(ctx context.Context, req *http.Request, accessKey, secretKey string, out core.Value) (core.Value, error) {
	httpClient := core.HttpClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("couldn't get http client from context")
	}

	if err := sign(req, accessKey, secretKey, excludedHeaders); err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error to send HTTP request: %w", err)
	}

	return core.UnwrapResponse(res, out)
}
