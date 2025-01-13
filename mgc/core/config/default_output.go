package config

import mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"

func defaultOutputSchema() *mgcSchemaPkg.Schema {
	s := mgcSchemaPkg.NewStringSchema()
	s.Description = "Default output string to be used when no other is specified"
	s.Pattern = "[A-Za-z]"
	return s
}
