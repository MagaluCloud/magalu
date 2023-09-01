package bucket

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"magalu.cloud/core"
	"magalu.cloud/sdk/static/s3"
)

type listObjectsParams struct {
	Bucket string `json:"bucket" jsonschema:"description=Name of the bucket to list objects from"`
}

type bucketContent struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	Size         int    `xml:"Size"`
}

type listObjectsResponse struct {
	Name     string           `xml:"Name"`
	Contents []*bucketContent `xml:"Contents"`
}

func newListObjectsRequest(ctx context.Context, region, bucket string) (*http.Request, error) {
	host := s3.BuildHost(region)
	url, err := url.JoinPath(host, bucket)
	if err != nil {
		return nil, err
	}
	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

func newListObjects() core.Executor {
	return core.NewStaticExecute(
		"list-objects",
		"",
		"List all objects from a bucket",
		listObjects,
	)
}

func listObjects(ctx context.Context, params listObjectsParams, cfg s3.Config) (core.Value, error) {
	bucket, _ := strings.CutPrefix(params.Bucket, s3.URIPrefix)
	req, err := newListObjectsRequest(ctx, cfg.Region, bucket)
	if err != nil {
		return nil, err
	}

	return s3.SendRequest(ctx, req, cfg.AccessKeyID, cfg.SecretKey, &listObjectsResponse{})
}
