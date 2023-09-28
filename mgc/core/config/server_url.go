package config

import (
	"github.com/invopop/jsonschema"
	"magalu.cloud/core/schema"
)

var networkConfigSchema *schema.Schema

type NetworkConfig struct {
	ServerUrl string `json:"serverUrl,omitempty" jsonschema:"description=Manually specify the server to use,format=uri,default="`
}

func NetworkConfigSchema() *schema.Schema {
	if networkConfigSchema == nil {
		networkConfigSchema, _ = schema.ToCoreSchema(jsonschema.Reflect(NetworkConfig{}))
	}
	return networkConfigSchema
}
