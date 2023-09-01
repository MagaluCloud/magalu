package bucket

import (
	"context"
	"strings"

	"magalu.cloud/core"
	"magalu.cloud/sdk/static/s3"
)

type deleteObjectParams struct {
	Destination string `json:"dst" jsonschema:"description=Path of the object to be deleted" example:"s3://bucket1/file1"`
}

func newDeleteObject() core.Executor {
	return core.NewStaticExecute(
		"delete-object",
		"",
		"Delete an object from a bucket",
		deleteObject,
	)
}

func deleteObject(ctx context.Context, params deleteObjectParams, cfg s3.Config) (core.Value, error) {
	bucketURI, _ := strings.CutPrefix(params.Destination, s3.URIPrefix)
	req, err := newDeleteRequest(ctx, cfg.Region, bucketURI)
	if err != nil {
		return nil, err
	}

	return s3.SendRequest(ctx, req, cfg.AccessKeyID, cfg.SecretKey, nil)
}
