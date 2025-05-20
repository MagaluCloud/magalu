package object_lock

import (
	"context"
	"errors"
	"net/http"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/sdk/static/object_storage/common"
)

type GetBucketObjectLockResponse struct {
	ObjectLockEnabled string
	Rule              common.ObjectLockRule
}

var ErrBucketMissingObjectLockConfiguration = errors.New("bucket missing object lock configuration")

type GetBucketObjectLockParams struct {
	Bucket common.BucketName `json:"dst" jsonschema:"description=Specifies the bucket whose ACL is being requested" mgc:"positional"`
}

var GetGet = utils.NewLazyLoader[core.Executor](func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "get",
			Description: "Get object locking configuration for the specified bucket",
		},
		GetObjectLocking,
	)
	exec = core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "json"
	})
	return exec
})

func GetObjectLocking(ctx context.Context, params GetBucketObjectLockParams, cfg common.Config) (result GetBucketObjectLockResponse, err error) {
	req, err := newGetObjectLockingRequest(ctx, cfg, params.Bucket)
	if err != nil {
		return
	}

	res, err := common.SendRequest(ctx, req, cfg)
	if err != nil {
		return
	}

	// Se a resposta de GET /bucket?locking for 400, isso quer dizer que o bucket não tem locking habilitado.
	// Como precisamos tratar esse caso de maneira específica, usamos um erro com tipo específico.
	if res.StatusCode == 400 {
		err = ErrBucketMissingObjectLockConfiguration
	}

	return
}

func newGetObjectLockingRequest(ctx context.Context, cfg common.Config, bucketName common.BucketName) (*http.Request, error) {
	url, err := common.BuildBucketHostURL(cfg, bucketName)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	query := url.Query()
	query.Add("object-lock", "")
	url.RawQuery = query.Encode()

	return http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
}
