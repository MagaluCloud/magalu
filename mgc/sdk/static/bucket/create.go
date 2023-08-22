package bucket

import (
	"context"
	"net/http"
	"net/url"

	"magalu.cloud/core"
	"magalu.cloud/sdk/static/s3"
)

type createParams struct {
	Name     string `json:"name" jsonschema:"description=Name of the bucket to be created"`
	ACL      string `json:"acl,omitempty" jsonschema:"description=ACL Rules for the bucket"`
	Location string `json:"location,omitempty" jsonschema:"description=Location constraint for the bucket,default=br-ne-1"`
}

func newCreate() core.Executor {
	return core.NewStaticExecute(
		"create",
		"",
		"Create a bucket",
		create,
	)
}

func newCreateRequest(region, bucket string) (*http.Request, error) {
	host := s3.BuildHost(region)
	url, err := url.JoinPath(host, bucket)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(http.MethodPut, url, nil)
}

func create(ctx context.Context, params createParams, cfg s3.Config) (core.Value, error) {
	req, err := newCreateRequest(cfg.Region, params.Name)
	if err != nil {
		return nil, err
	}

	return s3.SendRequest(ctx, req, cfg.AccessKeyID, cfg.SecretKey, nil)
}
