package bucket

import (
	"context"
	"net/http"
	"net/url"

	"magalu.cloud/core"
	"magalu.cloud/sdk/static/s3"
)

type deleteParams struct {
	Name string `json:"name" jsonschema:"description=Name of the bucket to be deleted"`
}

func newDelete() core.Executor {
	return core.NewStaticExecute(
		"delete",
		"",
		"Delete a bucket",
		delete,
	)
}

func newDeleteRequest(region, bucket string) (*http.Request, error) {
	host := s3.BuildHost(region)
	url, err := url.JoinPath(host, bucket)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(http.MethodDelete, url, nil)
}

func delete(ctx context.Context, params deleteParams, cfg s3.Config) (core.Value, error) {
	req, err := newDeleteRequest(cfg.Region, params.Name)
	if err != nil {
		return nil, err
	}

	return s3.SendRequest(ctx, req, cfg.AccessKeyID, cfg.SecretKey, nil)
}
