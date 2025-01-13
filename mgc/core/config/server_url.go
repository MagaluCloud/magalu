package config

import (
	"github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/invopop/jsonschema"
)

var networkConfigSchema *schema.Schema

type NetworkConfig struct {
	ServerUrl string `json:"serverUrl,omitempty" jsonschema:"description=Manually specify the server to use,format=uri"`
}

func NetworkConfigSchema() *schema.Schema {
	if networkConfigSchema == nil {
		networkConfigSchema, _ = schema.ToCoreSchema(jsonschema.Reflect(NetworkConfig{}))
	}
	return networkConfigSchema
}
