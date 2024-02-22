package common

import "magalu.cloud/core/config"

type Config struct {
	Workers int    `json:"workers,omitempty" jsonschema:"description=Number of routines that spawn to do parallel operations within object_storage,default=5,minimum=1"`
	Region  string `json:"region,omitempty" jsonschema:"description=Region to reach the service,default=br-ne1,enum=br-ne1,enum=br-se1,enum=br-mgl1"`
	// See more about the 'squash' directive here: https://pkg.go.dev/github.com/mitchellh/mapstructure#hdr-Embedded_Structs_and_Squashing
	config.NetworkConfig `json:",squash"` // nolint
}

var regionMap = map[string]string{
	"br-ne1":  "br-ne-1",
	"br-se1":  "br-se-1",
	"br-mgl1": "br-mgl-1",
}

func (c *Config) translateRegion() string {
	return regionMap[c.Region]
}
