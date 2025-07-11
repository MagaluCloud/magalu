package policy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type setBucketPolicyParams struct {
	Bucket common.BucketName `json:"dst" jsonschema:"description=Name of the bucket to set permissions for,example=my-bucket" mgc:"positional"`
	Policy map[string]any    `json:"policy" jsonschema:"description=Policy file path to be uploaded,example=@./policy.json or ./policy.json" mgc:"positional"`
}

var getSet = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "set",
			Description: "Set policy document for the specified bucket. The policy can be provided as a direct JSON string or a file path using @./policy.json.",
		},
		setPolicy,
	)

	exec = core.NewExecuteFormat(exec, func(exec core.Executor, result core.Result) string {
		return fmt.Sprintf("Successfully set policy for bucket %q", result.Source().Parameters["dst"])
	})

	return exec
})

func setPolicy(ctx context.Context, params setBucketPolicyParams, cfg common.Config) (result core.Value, err error) {
	req, err := newSetBucketPolicyRequest(ctx, params, cfg)
	if err != nil {
		return
	}

	resp, err := common.SendRequest(ctx, req, cfg)
	if err != nil {
		return
	}

	err = common.ExtractErr(resp, req)
	if err != nil {
		return
	}

	return
}

func newSetBucketPolicyRequest(ctx context.Context, p setBucketPolicyParams, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, p.Bucket)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("policy", "")
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url.String(), nil)
	if err != nil {
		return nil, err
	}

	getBody := func() (io.ReadCloser, error) {
		body, err := json.Marshal(p.Policy)
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(body)
		return io.NopCloser(reader), nil
	}

	req.Body, err = getBody()
	if err != nil {
		return nil, err
	}
	req.GetBody = getBody

	return req, nil
}
