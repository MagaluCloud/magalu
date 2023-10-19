package s3

import "magalu.cloud/core/config"

type Config struct {
	Token   string `json:"token,omitempty" jsonschema:"description=Token for S3 Credentials"`
	Workers int    `json:"workers,omitempty" jsonschema:"description=Number of routines that spawn to do parallel operations within object_storage,default=5,exlusiveMinimum=0"`
	Region  string `json:"region,omitempty" jsonschema:"description=Region to reach the service,default=br-ne-1,enum=br-ne-1,enum=br-ne-2,enum=br-se-1"`
	// See more about the 'squash' directive here: https://pkg.go.dev/github.com/mitchellh/mapstructure#hdr-Embedded_Structs_and_Squashing
	config.NetworkConfig `json:",squash"` // nolint
}
