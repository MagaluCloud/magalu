// Package serverpin keeps a downloaded spec from sending the user's token to a host it picks.
package serverpin

import (
	"fmt"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"gopkg.in/yaml.v3"
)

const (
	prodHost    = "api.magalu.cloud"
	preProdHost = "api.pre-prod.jaxyendy.com:8443"
)

var envVariable = map[string]any{
	"description": "Environment to use",
	"default":     prodHost,
	"enum":        []string{prodHost, preProdHost},
	"x-mgc-transforms": []map[string]any{{
		"type": "translate",
		"translations": []map[string]string{
			{"from": "prod", "to": prodHost},
			{"from": "pre-prod", "to": preProdHost},
		},
	}},
}

var operationMethods = []string{"get", "put", "post", "delete", "options", "head", "patch", "trace"}

type Store struct {
	preview.Store
}

func Wrap(store preview.Store) *Store {
	return &Store{store}
}

func (s *Store) Read(name string) (preview.Package, error) {
	pkg, err := s.Store.Read(name)
	if err != nil {
		return pkg, err
	}

	pkg.Spec, err = Pin(pkg.Spec)
	return pkg, err
}

func Pin(spec []byte) ([]byte, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(spec, &doc); err != nil {
		return nil, fmt.Errorf("spec does not parse: %w", err)
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return spec, nil
	}

	var env yaml.Node
	if err := env.Encode(envVariable); err != nil {
		return nil, err
	}

	root := doc.Content[0]
	pinServers(root, &env)
	if paths := value(root, "paths"); paths != nil && paths.Kind == yaml.MappingNode {
		for i := 1; i < len(paths.Content); i += 2 {
			item := paths.Content[i]
			if item.Kind != yaml.MappingNode {
				continue
			}
			pinServers(item, &env)
			for _, method := range operationMethods {
				if op := value(item, method); op != nil && op.Kind == yaml.MappingNode {
					pinServers(op, &env)
				}
			}
		}
	}

	return yaml.Marshal(&doc)
}

func pinServers(object *yaml.Node, env *yaml.Node) {
	servers := value(object, "servers")
	if servers == nil || servers.Kind != yaml.SequenceNode {
		return
	}

	for _, server := range servers.Content {
		if server.Kind != yaml.MappingNode {
			continue
		}
		if url := value(server, "url"); url != nil && url.Kind == yaml.ScalarNode {
			url.Value = PinURL(url.Value)
		}

		variables := value(server, "variables")
		if variables == nil || variables.Kind != yaml.MappingNode {
			variables = &yaml.Node{Kind: yaml.MappingNode}
			set(server, "variables", variables)
		}
		set(variables, "env", env)
	}
}

func PinURL(raw string) string {
	path := raw
	if _, rest, ok := strings.Cut(raw, "://"); ok {
		path = ""
		if i := strings.Index(rest, "/"); i >= 0 {
			path = rest[i:]
		}
	}
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return "https://{env}" + path
}

func value(object *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(object.Content); i += 2 {
		if object.Content[i].Value == key {
			return object.Content[i+1]
		}
	}
	return nil
}

func set(object *yaml.Node, key string, node *yaml.Node) {
	for i := 0; i+1 < len(object.Content); i += 2 {
		if object.Content[i].Value == key {
			object.Content[i+1] = node
			return
		}
	}
	object.Content = append(object.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key}, node)
}

var _ preview.Store = (*Store)(nil)
