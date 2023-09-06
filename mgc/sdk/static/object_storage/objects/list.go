package objects

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"magalu.cloud/core"
	"magalu.cloud/sdk/static/object_storage/s3"
)

type ListObjectsParams struct {
	Destination string `json:"dst" jsonschema:"description=Path of the bucket to list objects from" example:"s3://bucket1/"`
}

type bucketContent struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	Size         int    `xml:"Size"`
}

type ListObjectsResponse struct {
	Name     string           `xml:"Name"`
	Contents []*bucketContent `xml:"Contents"`
}

func newListRequest(ctx context.Context, region, bucket string) (*http.Request, error) {
	host := s3.BuildHost(region)
	url, err := url.JoinPath(host, bucket)
	if err != nil {
		return nil, err
	}
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

func newList() core.Executor {
	return core.NewStaticExecute(
		"list",
		"",
		"List all objects from a bucket",
		List,
	)
}

func List(ctx context.Context, params ListObjectsParams, cfg s3.Config) (result ListObjectsResponse, err error) {
	bucket, _ := strings.CutPrefix(params.Destination, s3.URIPrefix)
	req, err := newListRequest(ctx, cfg.Region, bucket)
	if err != nil {
		return result, err
	}

	return s3.SendRequest(ctx, req, cfg.AccessKeyID, cfg.SecretKey, &result)
}
