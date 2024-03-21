package versioning

import (
	"context"
	"net/http"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type SuspendBucketVersioningParams struct {
	Bucket common.BucketName `json:"bucket" jsonschema:"description=Bucket name to suspend versioning" mgc:"positional"`
}

var getSuspend = utils.NewLazyLoader(func() core.Executor {
	var exec core.Executor = core.NewStaticExecute(
		core.DescriptorSpec{
			Name:        "suspend",
			Description: "Suspend versioning for a Bucket",
		},
		SuspendBucketVersioning,
	)

	return core.NewExecuteResultOutputOptions(exec, func(exec core.Executor, result core.Result) string {
		return "template=Suspended versioning for {{.bucket}}\n"
	})
})

func SuspendBucketVersioning(ctx context.Context, params SuspendBucketVersioningParams, cfg common.Config) (core.Value, error) {
	req, err := newSuspendBucketVersioningRequest(ctx, params.Bucket, cfg)
	if err != nil {
		return nil, err
	}

	res, err := common.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	_, err = common.UnwrapResponse[core.Value](res, req)
	if err != nil {
		return nil, err
	}
	return params, nil
}

func newSuspendBucketVersioningRequest(ctx context.Context, bucketName common.BucketName, cfg common.Config) (*http.Request, error) {
	return newSetBucketVersioningRequest(
		ctx,
		bucketName,
		cfg,
		versioningConfiguration{
			Status:    "Suspended",
			Namespace: "http://s3.amazonaws.com/doc/2006-03-01/",
		},
	)
}
