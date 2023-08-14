package bucket

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"magalu.cloud/core"
)

func newList() core.Executor {
	return core.NewStaticExecute(
		"list",
		"",
		"List all buckets",
		list,
	)
}

func list(ctx context.Context, _ struct{}, config bucketConfig) (*s3.ListBucketsOutput, error) {
	svc, err := getS3Client(ctx, config)
	if err != nil {
		return nil, err
	}
	result, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("Failed to list buckets %w", err)
	}

	return result, nil
}
