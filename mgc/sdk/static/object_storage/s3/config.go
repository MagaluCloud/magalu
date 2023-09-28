package s3

import "magalu.cloud/core/config"

type Config struct {
	AccessKeyID string `json:"accessKeyId" jsonschema:"description=Access Key ID for S3 Credentials"`
	SecretKey   string `json:"secretKey" jsonschema:"description=Secret Key for S3 Credentials"`
	Token       string `json:"token,omitempty" jsonschema:"description=Token for S3 Credentials"`
	Region      string `json:"region,omitempty" jsonschema:"description=Region to reach the service,default=br-ne-1,enum=br-ne-1,enum=br-ne-2,enum=br-se-1"`
	// See more about the 'squash' directive here: https://pkg.go.dev/github.com/mitchellh/mapstructure#hdr-Embedded_Structs_and_Squashing
	config.NetworkConfig `json:",squash"` // nolint
}
