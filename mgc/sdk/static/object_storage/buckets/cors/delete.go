package cors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type deleteBucketCorsParams struct {
	Bucket common.BucketName `json:"dst" jsonschema:"description=Name of the bucket to delete CORS file from,example=my-bucket" mgc:"positional"`
}

var getDelete = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "delete",
			Description: "Delete CORS rules for the specified bucket",
		},
		deleteCors,
	)

	exec = core.NewExecuteFormat(exec, func(exec core.Executor, result core.Result) string {
		return fmt.Sprintf("Successfully deleted CORS for bucket %q", result.Source().Parameters["dst"])
	})

	return exec
})

func deleteCors(ctx context.Context, params deleteBucketCorsParams, cfg common.Config) (result core.Value, err error) {
	req, err := newDeleteBucketCorsRequest(ctx, params, cfg)
	if err != nil {
		return
	}

	resp, err := common.SendRequest(ctx, req)
	if err != nil {
		return
	}

	err = common.ExtractErr(resp, req)
	if err != nil {
		return
	}

	return
}

func newDeleteBucketCorsRequest(ctx context.Context, p deleteBucketCorsParams, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, p.Bucket)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("cors", "")
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
