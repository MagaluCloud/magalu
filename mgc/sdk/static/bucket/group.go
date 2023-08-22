package bucket

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"magalu.cloud/core"
	mgcS3 "magalu.cloud/sdk/static/s3"
)

func getS3Client(ctx context.Context, c mgcS3.Config) (*s3.S3, error) {
	suffixPath := true
	// TODO: fix once MGC accepts "br-ne-1"
	hostedS3Region := "us-east-1"
	endpoint := mgcS3.BuildHost(c.Region)
	cfg := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(c.AccessKeyID, c.SecretKey, c.Token),
		Endpoint:         &endpoint,
		Region:           &hostedS3Region,
		S3ForcePathStyle: &suffixPath,
	}

	session, err := session.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to create bucket session %w", err)
	}

	svc := s3.New(session)
	if svc == nil {
		return nil, fmt.Errorf("Failed to create bucket service")
	}

	return svc, nil
}

func NewGroup() core.Grouper {
	return core.NewStaticGroup(
		"bucket",
		"",
		"Bucket operations for Object Storage API",
		[]core.Descriptor{
			newCreate(),
			newDelete(),
			newList(),
		},
	)
}
