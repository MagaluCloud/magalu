package bucket

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"magalu.cloud/core"
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

func delete(ctx context.Context, params deleteParams, config bucketConfig) (*s3.DeleteBucketOutput, error) {
	svc, err := getS3Client(ctx, config)
	if err != nil {
		return nil, err
	}
	res, err := svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: &params.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to delete bucket %w", err)
	}

	return res, nil
}
