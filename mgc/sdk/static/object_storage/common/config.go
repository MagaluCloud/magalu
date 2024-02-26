package common

import "magalu.cloud/core/config"

type Config struct {
	Workers   int    `json:"workers,omitempty" jsonschema:"description=Number of routines that spawn to do parallel operations within object_storage,default=5,minimum=1"`
	ChunkSize uint64 `json:"chunkSize,omitempty" jsonschema:"description=Chunk size to consider when doing multipart requests. Specified in Mb,default=5,minimum=5,maximum=5120"`
	Region    string `json:"region,omitempty" jsonschema:"description=Region to reach the service,default=br-ne1,enum=br-ne-1,enum=br-se-1,enum=br-mgl-1"`
	// See more about the 'squash' directive here: https://pkg.go.dev/github.com/mitchellh/mapstructure#hdr-Embedded_Structs_and_Squashing
	config.NetworkConfig `json:",squash"` // nolint
}

var regionMap = map[string]string{
	"br-ne-1":  "br-ne-1",
	"br-se-1":  "br-se-1",
	"br-mgl-1": "br-se-1",
}

func (c *Config) translateRegion() string {
	return regionMap[c.Region]
}

func (c *Config) chunkSizeInBytes() uint64 {
	if c.ChunkSize <= MIN_CHUNK_SIZE {
		return MIN_CHUNK_SIZE
	}
	if c.ChunkSize >= MAX_CHUNK_SIZE {
		return MAX_CHUNK_SIZE
	}

	return c.ChunkSize * (1024 * 1024)
}
