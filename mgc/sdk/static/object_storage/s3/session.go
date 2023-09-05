package s3

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"magalu.cloud/core"
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
}

func BuildHost(region string) string {
	return strings.ReplaceAll(templateUrl, "{{region}}", region)
}

func SendRequest(ctx context.Context, req *http.Request, accessKey, secretKey string, out core.Value) (core.Value, error) {
	httpClient := corehttp.ClientFromContext(ctx)
	if httpClient == nil {
		return nil, fmt.Errorf("couldn't get http client from context")
	}

	var unsignedPayload bool
	switch req.Body.(type) {
	case io.ReadSeeker:
		unsignedPayload = true
	}

	if err := sign(req, accessKey, secretKey, unsignedPayload, excludedHeaders); err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error to send HTTP request: %w", err)
	}

	return corehttp.UnwrapResponse(res, out)
}
