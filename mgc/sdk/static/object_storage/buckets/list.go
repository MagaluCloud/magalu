package buckets

import (
	"context"
	"net/http"

	"magalu.cloud/core"
	"magalu.cloud/sdk/static/object_storage/s3"
)

type BucketResponse struct {
	CreationDate string `xml:"CreationDate"`
	Name         string `xml:"Name"`
}

// Container for the owner's display name and ID.
type Owner struct {
	DisplayName *string `xml:"DisplayName"`
	ID          *string `type:"ID"`
}

type ListResponse struct {
	Buckets []*BucketResponse `xml:"Buckets>Bucket"`
	Owner   *Owner            `xml:"Owner"`
}

func newListRequest(ctx context.Context, region string) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, http.MethodGet, s3.BuildHost(region), nil)
}

func newList() core.Executor {
	return core.NewStaticExecute(
		"list",
		"",
		"List all buckets",
		list,
	)
}

func list(ctx context.Context, _ struct{}, cfg s3.Config) (core.Value, error) {
	req, err := newListRequest(ctx, cfg.Region)
	if err != nil {
		return nil, err
	}

	return s3.SendRequest(ctx, req, cfg.AccessKeyID, cfg.SecretKey, &ListResponse{})
}
