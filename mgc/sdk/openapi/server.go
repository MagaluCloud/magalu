package openapi

import (
	"fmt"
	"strings"

	"slices"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/config"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/getkin/kin-openapi/openapi3"
)

type server struct {
	servers         openapi3.Servers
	extensionPrefix *string
}

func newServer(servers openapi3.Servers, extensionPrefix *string) *server {
	return &server{servers, extensionPrefix}
}

func (o *server) addToSchema(schema *mgcSchemaPkg.Schema) error {
	s := config.NetworkConfigSchema()
	for k, v := range s.Properties {
		if existing := schema.Properties[k]; existing != nil {
			return &utils.ChainedError{Name: k, Err: fmt.Errorf("already exists as schema %v", existing)}
		}
		schema.Properties[k] = v

		if slices.Contains(s.Required, k) {
			schema.Required = append(schema.Required, k)
		}
	}

	_, err := o.forEachVariable(func(externalName, internalName string, spec *openapi3.ServerVariable, server *openapi3.Server) (run bool, err error) {
		varSchema := openapi3.NewStringSchema()
		varSchema.Default = spec.Default

		varSchema.Description = getDescriptionExtension(o.extensionPrefix, spec.Extensions, spec.Description)
		for _, e := range spec.Enum {
			varSchema.Enum = append(varSchema.Enum, e)
		}
		varSchema.Extensions = spec.Extensions

		if existing := schema.Properties[externalName]; existing != nil {
			return false, &utils.ChainedError{Name: externalName, Err: fmt.Errorf("already exists as schema %v", existing)}
		}

		schema.Properties[externalName] = &openapi3.SchemaRef{Value: varSchema}
		return true, nil
	})

	slices.Sort(schema.Required)

	return err
}

type cbForEachVariable func(externalName string, internalName string, spec *openapi3.ServerVariable, server *openapi3.Server) (run bool, err error)

func (o *server) forEachVariable(cb cbForEachVariable) (finished bool, err error) {
	var s *openapi3.Server
	if len(o.servers) > 0 {
		s = o.servers[0]
	}

	if s == nil {
		return false, fmt.Errorf("no available servers in spec")
	}

	for internalName, spec := range s.Variables {
		externalName := getNameExtension(o.extensionPrefix, spec.Extensions, internalName)
		run, err := cb(externalName, internalName, spec, s)
		if err != nil {
			return false, err
		}
		if !run {
			return false, nil
		}
	}

	return true, nil
}

func (o *server) url(configs core.Configs) (string, error) {
	nc, _ := utils.DecodeNewValue[config.NetworkConfig](configs)

	if nc.ServerUrl != "" {
		return nc.ServerUrl, nil
	}
	if len(o.servers) == 0 {
		return "", fmt.Errorf("no available servers in spec")
	}
	if len(o.servers[0].Variables) == 0 {
		return o.servers[0].URL, nil
	}

	url := ""
	_, err := o.forEachVariable(func(externalName, internalName string, spec *openapi3.ServerVariable, server *openapi3.Server) (run bool, err error) {
		val, ok := configs[externalName]
		if !ok {
			val = spec.Default
		}
		tmpl := "{" + internalName + "}"

		if url == "" {
			url = server.URL
		}
		url = strings.ReplaceAll(url, tmpl, fmt.Sprintf("%v", val))

		return true, nil
	})

	if err != nil {
		return "", err
	}

	return url, nil
}
