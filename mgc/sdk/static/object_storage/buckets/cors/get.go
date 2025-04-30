package cors

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type GetBucketCorsParams struct {
	Bucket common.BucketName `json:"dst" jsonschema:"description=Specifies the bucket whose CORS rules is being requested" mgc:"positional"`
}

var getGet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get the CORS rules for the specified bucket",
		},
		getCors,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "json"
	})
	return exec
})

func getCors(ctx context.Context, params GetBucketCorsParams, cfg common.Config) (result map[string]any, err error) {
	req, err := newGetCorsRequest(ctx, cfg, params.Bucket)
	if err != nil {
		return nil, err
	}

	res, err := common.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		bodyStr := string(bodyBytes)
		return nil, fmt.Errorf("Status: %s\n\n%s", res.Status, bodyStr)
	}

	var corsConfig CORSConfiguration
	if err = xml.Unmarshal(bodyBytes, &corsConfig); err != nil {
		return nil, err
	}

	result = map[string]any{
		"CORSRules": corsConfig.CORSRules,
	}
	return
}

func newGetCorsRequest(ctx context.Context, cfg common.Config, bucketName common.BucketName) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, bucketName)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("cors", "")
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	return req, nil
}
