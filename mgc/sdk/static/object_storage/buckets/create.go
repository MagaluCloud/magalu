package buckets

import (
	"context"
	"net/http"
	"net/url"

	"magalu.cloud/core"
	"magalu.cloud/sdk/static/object_storage/s3"
)

type createParams struct {
	Name     string `json:"name" jsonschema:"description=Name of the bucket to be created"`
	ACL      string `json:"acl,omitempty" jsonschema:"description=ACL Rules for the bucket"`
	Location string `json:"location,omitempty" jsonschema:"description=Location constraint for the bucket,default=br-ne-1"`
}

func newCreate() core.Executor {
	executor := core.NewStaticExecute(
		"create",
		"",
		"Create a bucket",
		create,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Created bucket {{.name}}\n"
	})
}

func newCreateRequest(ctx context.Context, cfg s3.Config, bucket string) (*http.Request, error) {
	host := s3.BuildHost(cfg)
	url, err := url.JoinPath(host, bucket)
	if err != nil {
		return nil, err
	}
	return http.NewRequestWithContext(ctx, http.MethodPut, url, nil)
}

func create(ctx context.Context, params createParams, cfg s3.Config) (core.Value, error) {
	req, err := newCreateRequest(ctx, cfg, params.Name)
	if err != nil {
		return nil, err
	}

	_, _, err = s3.SendRequest[core.Value](ctx, req, cfg.AccessKeyID, cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	return params, nil
}
