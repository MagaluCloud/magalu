package buckets

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/buckets/versioning"
	"magalu.cloud/sdk/static/object_storage/common"
)

var createLogger = utils.NewLazyLoader(func() *zap.SugaredLogger {
	return logger().Named("create")
})

type createParams struct {
	BucketName            common.BucketName `json:"bucket" jsonschema:"description=Name of the bucket to be created" mgc:"positional"`
	IsPrefix              bool              `json:"bucket_is_prefix" jsonschema:"description=Use bucket name as prefix value to generate a unique bucket name,default=false" mgc:"hidden"`
	EnableVersioning      bool              `json:"enable_versioning,omitempty" jsonschema:"description=Enable versioning for this bucket,default=true"`
	common.ACLPermissions `json:",squash"`  // nolint
}

var getCreate = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewReflectedSimpleExecutor[createParams, common.Config, *createParams](
		core.ExecutorSpec{
			DescriptorSpec: core.DescriptorSpec{
				Name:        "create",
				Summary:     "Create a new Bucket",
				Description: `Buckets are "containers" that are able to store various Objects inside`,
				// Scopes:      core.Scopes{"object-storage.write"},
			},
			Links: utils.NewLazyLoaderWithArg(func(e core.Executor) core.Links {
				return core.Links{
					"delete": core.NewSimpleLink(
						core.SimpleLinkSpec{
							Owner:     e,
							Target:    getDelete(),
							FromOwner: map[string]string{"name": "bucket"},
						},
					),
					"list": core.NewSimpleLink(
						core.SimpleLinkSpec{
							Owner:  e,
							Target: getList(),
						},
					),
					"read": core.NewSimpleLink(
						core.SimpleLinkSpec{
							Owner:     e,
							Target:    getBucket(),
							FromOwner: map[string]string{"name": "bucket"},
						},
					),
				}
			}),
		},
		create,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Created bucket {{.bucket}}\n"
	})
})

func newCreateRequest(ctx context.Context, cfg common.Config, bucket common.BucketName, aclPermissions common.ACLPermissions) (*http.Request, error) {
	url, err := common.BuildBucketHost(cfg, bucket)
	if err != nil {
		return nil, core.UsageError{Err: err}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, string(url), nil)
	if err != nil {
		return nil, err
	}

	err = aclPermissions.SetHeaders(req, cfg)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func create(ctx context.Context, params createParams, cfg common.Config) (*createParams, error) {

	logger := createLogger().With(
		"bucket", params.BucketName,
		"location", cfg.Region,
	)

	err := params.ACLPermissions.Validate()
	if err != nil {
		return nil, err
	}

	if params.IsPrefix {
		sufix := common.GetRandomWords(3, "-")
		params.BucketName = common.BucketName(fmt.Sprintf("%s-%s", params.BucketName.String(), sufix))
	}

	req, err := newCreateRequest(ctx, cfg, params.BucketName, params.ACLPermissions)
	if err != nil {
		return nil, err
	}

	resp, err := common.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}

	err = common.ExtractErr(resp, req)
	if err != nil {
		return nil, err
	}

	logger.Info("bucket created successfully")

	if !params.EnableVersioning {
		logger.Info("suspending bucket versioning, as 'enable_versioning' was passed as false")
		_, err := versioning.SuspendBucketVersioning(ctx, versioning.SuspendBucketVersioningParams{Bucket: params.BucketName}, cfg)
		if err != nil {
			return nil, err
		}
	}

	return &params, nil
}
