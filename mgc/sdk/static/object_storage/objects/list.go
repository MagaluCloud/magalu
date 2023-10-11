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

type prefix struct {
	Path string `xml:"Prefix"`
}

type ListObjectsResponse struct {
	Name           string           `xml:"Name"`
	Contents       []*bucketContent `xml:"Contents"`
	CommonPrefixes []*prefix        `xml:"CommonPrefixes" json:"SubDirectories"`
}

func newListRequest(ctx context.Context, cfg s3.Config, bucket string) (*http.Request, error) {
	parsedUrl, err := parseURL(cfg, bucket)
	if err != nil {
		return nil, err
	}
	return http.NewRequestWithContext(ctx, http.MethodGet, parsedUrl.String(), nil)
}

func newList() core.Executor {
	return core.NewStaticExecute(
		"list",
		"",
		"List all objects from a bucket",
		List,
	)
}

func parseURL(cfg s3.Config, bucketURI string) (*url.URL, error) {
	dirs := strings.Split(bucketURI, "/")
	path, err := url.JoinPath(s3.BuildHost(cfg), dirs[0])
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	if len(dirs) <= 1 {
		return u, nil
	}
	q := u.Query()
	delimiter := "/"
	prefixQ := strings.Join(dirs[1:], delimiter)
	lastChar := string(prefixQ[len(prefixQ)-1])
	if lastChar != delimiter {
		prefixQ += delimiter
	}
	q.Set("prefix", prefixQ)
	q.Set("delimiter", delimiter)
	q.Set("encoding-type", "url")
	u.RawQuery = q.Encode()
	return u, nil
}

func List(ctx context.Context, params ListObjectsParams, cfg s3.Config) (result ListObjectsResponse, err error) {
	bucket, _ := strings.CutPrefix(params.Destination, s3.URIPrefix)
	req, err := newListRequest(ctx, cfg, bucket)
	if err != nil {
		return
	}

	result, _, err = s3.SendRequest[ListObjectsResponse](ctx, req, cfg.AccessKeyID, cfg.SecretKey)
	return
}
