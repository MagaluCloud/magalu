package objects

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"magalu.cloud/core"
	"magalu.cloud/sdk/static/object_storage/s3"
)

type DeleteObjectParams struct {
	Destination string `json:"dst" jsonschema:"description=Path of the object to be deleted" example:"s3://bucket1/file1"`
}

func newDelete() core.Executor {
	return core.NewStaticExecute(
		"delete",
		"",
		"Delete an object from a bucket",
		Delete,
	)
}

func newDeleteRequest(ctx context.Context, region string, pathURIs ...string) (*http.Request, error) {
	host := s3.BuildHost(region)
	url, err := url.JoinPath(host, pathURIs...)
	if err != nil {
		return nil, err
	}
	return http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
}

func Delete(ctx context.Context, params DeleteObjectParams, cfg s3.Config) (core.Value, error) {
	bucketURI, _ := strings.CutPrefix(params.Destination, s3.URIPrefix)
	req, err := newDeleteRequest(ctx, cfg.Region, bucketURI)
	if err != nil {
		return nil, err
	}

	return s3.SendRequest(ctx, req, cfg.AccessKeyID, cfg.SecretKey, (*core.Value)(nil))
}
