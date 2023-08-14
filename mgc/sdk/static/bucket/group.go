package bucket

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"magalu.cloud/core"
)

const (
	templateUrl = "https://{{region}}.api.magalu.cloud/magaluobjects"
)

type bucketConfig struct {
	AccessKeyID string `json:"accessKeyId" mapstructure:"accessKeyId" jsonschema:"description=Access Key ID for S3 Credentials"`
	SecretKey   string `json:"secretKey" mapstructure:"secretKey" jsonschema:"description=Secret Key for S3 Credentials"`
	Token       string `json:"token,omitempty" mapstructure:"s3-token" jsonschema:"description=Token for S3 Credentials"`
	Region      string `json:"region,omitempty" jsonschema:"description=Region to reach the service,default=br-ne-1,enum=br-ne-1,enum=br-ne-2,enum=br-se-1"`
}

func getS3Client(ctx context.Context, c bucketConfig) (*s3.S3, error) {
	suffixPath := true
	// TODO: fix once MGC accepts "br-ne-1"
	hostedS3Region := "us-east-1"
	endpoint := strings.ReplaceAll(templateUrl, "{{region}}", c.Region)
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
			newCreate(), // cmd: "create"
			newDelete(), // cmd: "delete"
			newList(),   // cmd: "list"
		},
	)
}
