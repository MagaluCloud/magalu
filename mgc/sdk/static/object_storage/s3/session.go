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

func SendRequest[T core.Value](ctx context.Context, req *http.Request, accessKey, secretKey string, dataPtr *T) (result T, err error) {
	httpClient := corehttp.ClientFromContext(ctx)
	if httpClient == nil {
		return result, fmt.Errorf("couldn't get http client from context")
	}

	var unsignedPayload bool
	switch req.Body.(type) {
	case io.ReadSeeker:
		unsignedPayload = true
	}

	if err := sign(req, accessKey, secretKey, unsignedPayload, excludedHeaders); err != nil {
		return result, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("error to send HTTP request: %w", err)
	}

	data, err := corehttp.UnwrapResponse(res, dataPtr)
	if err != nil || data == nil {
		return result, err
	}
	convertedVal, ok := data.(T)
	if !ok {
		return result, fmt.Errorf("failed to convert response value from %T to %T", data, result)
	}
	return convertedVal, err
}
