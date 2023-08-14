package bucket

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"magalu.cloud/core"
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

// TODO: change `convertValue()` to correctly infer a *string and avoid validation errors
type BucketOutput struct {
	_ s3.CreateBucketOutput

	Location string
}

func create(ctx context.Context, p createParams, c bucketConfig) (*BucketOutput, error) {
	svc, err := getS3Client(ctx, c)
	if err != nil {
		return nil, err
	}
	input := &s3.CreateBucketInput{Bucket: &p.Name}
	if p.Location != "" {
		input.CreateBucketConfiguration = &s3.CreateBucketConfiguration{
			LocationConstraint: &p.Location,
		}
	}
	if p.ACL != "" {
		input.ACL = &p.ACL
	}
	res, err := svc.CreateBucket(input)
	if err != nil {
		return nil, fmt.Errorf("Failed to create bucket %w", err)
	}

	return &BucketOutput{
		Location: *res.Location,
	}, nil
}
