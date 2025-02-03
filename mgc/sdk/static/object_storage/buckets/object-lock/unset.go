package object_lock

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type unsetBucketObjectLockParams struct {
	Bucket common.BucketName `json:"dst" jsonschema:"description=Name of the bucket to unset object locking for its objects,example=my-bucket" mgc:"positional"`
}

var getUnset = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "unset",
			Description: "unset object locking for the specified bucket",
		},
		unsetObjectLocking,
	)

	exec = core.NewExecuteFormat(exec, func(exec core.Executor, result core.Result) string {
		return fmt.Sprintf("Successfully removed Object Locking for bucket %q", result.Source().Parameters["dst"])
	})

	return exec
})

func unsetObjectLocking(ctx context.Context, params unsetBucketObjectLockParams, cfg common.Config) (result core.Value, err error) {
	req, err := newUnsetBucketObjectLockingRequest(ctx, params, cfg)
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

func newUnsetBucketObjectLockingRequest(ctx context.Context, p unsetBucketObjectLockParams, cfg common.Config) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, p.Bucket)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("object-lock", "")
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url.String(), nil)
	if err != nil {
		return nil, err
	}

	getBody := func() (io.ReadCloser, error) {
		body := `<ObjectLockConfiguration xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><ObjectLockEnabled>Enabled</ObjectLockEnabled></ObjectLockConfiguration>`

		reader := bytes.NewReader([]byte(body))
		return io.NopCloser(reader), nil
	}

	req.Body, err = getBody()
	if err != nil {
		return nil, err
	}
	req.GetBody = getBody

	req.Header.Set("Content-Type", "application/xml")

	return req, nil
}
