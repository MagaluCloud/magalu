package policy

import (
	"context"
	"fmt"
	"net/http"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type deleteBucketPolicyParams struct {
	Bucket common.BucketName `json:"dst" jsonschema:"description=Name of the bucket to delete policy file from,example=my-bucket" mgc:"positional"`
}

var getDelete = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "delete",
			Description: "Delete policy document for the specified bucket",
		},
		deletePolicy,
	)

	exec = core.NewExecuteFormat(exec, func(exec core.Executor, result core.Result) string {
		return fmt.Sprintf("Successfully deleted policy for bucket %q", result.Source().Parameters["dst"])
	})

	return exec
})

func deletePolicy(ctx context.Context, params deleteBucketPolicyParams, cfg common.Config) (result core.Value, err error) {
	req, err := newDeleteBucketPolicyRequest(ctx, params, cfg)
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

func newDeleteBucketPolicyRequest(ctx context.Context, p deleteBucketPolicyParams, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, p.Bucket)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("policy", "")
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
