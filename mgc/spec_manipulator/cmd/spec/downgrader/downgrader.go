package downgrader

import (
	"fmt"
	"os"
	"strings"

	"github.com/pb33f/libopenapi/datamodel"
	"gopkg.in/yaml.v3"
)

type OpenAPIDowngrader struct {
	config *datamodel.DocumentConfiguration
}

func NewOpenAPIDowngrader() *OpenAPIDowngrader {
	return &OpenAPIDowngrader{
		config: &datamodel.DocumentConfiguration{
			BypassDocumentCheck: true,
		},
	}
}

func (d *OpenAPIDowngrader) DowngradeFile(inputFile, outputFile string) error {
	inputBytes, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de entrada: %v", err)
	}

	outputBytes, err := d.Downgrade(inputBytes)
	if err != nil {
		return err
	}

	err = os.WriteFile(outputFile, outputBytes, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo de saída: %v", err)
	}

	return nil
}

func (d *OpenAPIDowngrader) Downgrade(inputBytes []byte) ([]byte, error) {
	specInfo, err := datamodel.ExtractSpecInfoWithConfig(inputBytes, d.config)
	if err != nil {
		return nil, fmt.Errorf("erro ao extrair informações do spec: %v", err)
	}

	if specInfo.SpecType != "openapi" {
		return nil, fmt.Errorf("o arquivo não é uma especificação OpenAPI válida")
	}

	document := make(map[string]interface{})
	if err := convertNodeToMap(specInfo.RootNode, &document); err != nil {
		return nil, fmt.Errorf("erro ao converter documento para JSON: %v", err)
	}

	if err := d.applyTransformations(&document); err != nil {
		return nil, fmt.Errorf("erro ao aplicar transformações: %v", err)
	}

	outputBytes, err := yaml.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter para YAML: %v", err)
	}

	return outputBytes, nil
}

func convertNodeToMap(node *yaml.Node, out interface{}) error {
	bytes, err := yaml.Marshal(node)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, out)
}

func (d *OpenAPIDowngrader) applyTransformations(document *map[string]interface{}) error {
	(*document)["openapi"] = "3.0.3"

	transformers := []func(*map[string]interface{}) error{
		d.removeWebhooks,
		d.transformJSONSchema,
		d.transformComponents,
		d.transformPaths,
		d.transformInfo,
		d.transformSecuritySchemes,
		d.transformServers,
	}

	for _, transformer := range transformers {
		if err := transformer(document); err != nil {
			return err
		}
	}

	return nil
}

func (d *OpenAPIDowngrader) removeWebhooks(document *map[string]interface{}) error {
	delete(*document, "webhooks")

	return nil
}

func (d *OpenAPIDowngrader) transformJSONSchema(document *map[string]interface{}) error {
	processAllSchemas(*document, d.transformSchema)

	return nil
}

func (d *OpenAPIDowngrader) transformSchema(schema map[string]interface{}) {
	deleteKeys := []string{
		"contentEncoding", "contentMediaType", "contentSchema", "examples",
		"unevaluatedProperties", "if", "then", "else",
		"dependentSchemas", "dependentRequired", "maxContains", "minContains",
		"const", "propertyNames",
	}
	for _, key := range deleteKeys {
		delete(schema, key)
	}

	if typeVal, ok := schema["type"]; ok {
		if typeStr, ok := typeVal.(string); ok && typeStr == "null" {
			delete(schema, "type")
			schema["nullable"] = true
		}

		if typeArr, ok := typeVal.([]interface{}); ok {
			newTypes := make([]interface{}, 0)
			hasNull := false

			for _, t := range typeArr {
				if tStr, ok := t.(string); ok {
					if tStr == "null" {
						hasNull = true
					} else {
						newTypes = append(newTypes, tStr)
					}
				}
			}

			if hasNull {
				schema["nullable"] = true
				if len(newTypes) == 1 {
					schema["type"] = newTypes[0]
				} else if len(newTypes) > 1 {
					if len(newTypes) > 0 {
						schema["type"] = newTypes[0]
					} else {
						delete(schema, "type")
					}
				} else {
					delete(schema, "type")
				}
			}
		}
	}

	if prefixItems, ok := schema["prefixItems"]; ok {
		if _, ok := schema["items"]; !ok {
			schema["items"] = prefixItems
		}
		delete(schema, "prefixItems")
	}

	if excMin, ok := schema["exclusiveMinimum"]; ok {
		if _, isBool := excMin.(bool); isBool {
			delete(schema, "exclusiveMinimum")
		}
	}
	if excMax, ok := schema["exclusiveMaximum"]; ok {
		if _, isBool := excMax.(bool); isBool {
			delete(schema, "exclusiveMaximum")
		}
	}
}

func (d *OpenAPIDowngrader) transformComponents(document *map[string]interface{}) error {
	componentsIface, exists := (*document)["components"]
	if !exists {
		return nil
	}

	components, ok := componentsIface.(map[string]interface{})
	if !ok {
		return nil
	}

	deleteKeys := []string{
		"pathItems", "webhooks",
	}
	for _, key := range deleteKeys {
		delete(components, key)
	}

	if examples, ok := components["examples"].(map[string]interface{}); ok {
		for name, example := range examples {
			if exampleMap, ok := example.(map[string]interface{}); ok {
				if _, hasValue := exampleMap["value"]; !hasValue {
					newExample := map[string]interface{}{
						"value": exampleMap,
					}

					if summary, ok := exampleMap["summary"]; ok {
						newExample["summary"] = summary
						delete(exampleMap, "summary")
					}
					if description, ok := exampleMap["description"]; ok {
						newExample["description"] = description
						delete(exampleMap, "description")
					}

					examples[name] = newExample
				}
			}
		}
	}

	return nil
}

func (d *OpenAPIDowngrader) transformPaths(document *map[string]interface{}) error {
	pathsIface, exists := (*document)["paths"]
	if !exists {
		return nil
	}

	paths, ok := pathsIface.(map[string]interface{})
	if !ok {
		return nil
	}

	for _, pathItem := range paths {
		if pathItemMap, ok := pathItem.(map[string]interface{}); ok {
			operations := []string{"get", "post", "put", "delete", "options", "head", "patch", "trace"}
			for _, op := range operations {
				if operation, ok := pathItemMap[op].(map[string]interface{}); ok {
					if requestBody, ok := operation["requestBody"].(map[string]interface{}); ok {
						if content, ok := requestBody["content"].(map[string]interface{}); ok {
							for _, mediaTypeObj := range content {
								if mediaType, ok := mediaTypeObj.(map[string]interface{}); ok {
									d.transformMediaTypeExamples(mediaType)
								}
							}
						}
					}

					if responses, ok := operation["responses"].(map[string]interface{}); ok {
						for _, responseObj := range responses {
							if response, ok := responseObj.(map[string]interface{}); ok {
								if content, ok := response["content"].(map[string]interface{}); ok {
									for _, mediaTypeObj := range content {
										if mediaType, ok := mediaTypeObj.(map[string]interface{}); ok {
											d.transformMediaTypeExamples(mediaType)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (d *OpenAPIDowngrader) transformMediaTypeExamples(mediaType map[string]interface{}) {
	if examples, ok := mediaType["examples"].(map[string]interface{}); ok {
		if len(examples) == 1 {
			for _, example := range examples {
				if exampleMap, ok := example.(map[string]interface{}); ok {
					if value, ok := exampleMap["value"]; ok {
						mediaType["example"] = value
					}
				}
			}
		}
	}
}

func (d *OpenAPIDowngrader) transformInfo(document *map[string]interface{}) error {
	infoIface, exists := (*document)["info"]
	if !exists {
		return nil
	}

	info, ok := infoIface.(map[string]interface{})
	if !ok {
		return nil
	}

	deleteKeys := []string{
		"summary", "license.identifier",
	}
	for _, key := range deleteKeys {
		parts := strings.Split(key, ".")

		if len(parts) == 1 {
			delete(info, key)
		} else if len(parts) == 2 {
			if subObj, ok := info[parts[0]].(map[string]interface{}); ok {
				delete(subObj, parts[1])
			}
		}
	}

	return nil
}

func (d *OpenAPIDowngrader) transformSecuritySchemes(document *map[string]interface{}) error {
	componentsIface, exists := (*document)["components"]
	if !exists {
		return nil
	}

	components, ok := componentsIface.(map[string]interface{})
	if !ok {
		return nil
	}

	securitySchemesIface, exists := components["securitySchemes"]
	if !exists {
		return nil
	}

	securitySchemes, ok := securitySchemesIface.(map[string]interface{})
	if !ok {
		return nil
	}

	for _, schemeIface := range securitySchemes {
		scheme, ok := schemeIface.(map[string]interface{})
		if !ok {
			continue
		}

		if schemeType, ok := scheme["type"].(string); ok && schemeType == "mutualTLS" {
			scheme["type"] = "http"
			scheme["scheme"] = "https"
			scheme["description"] = fmt.Sprintf("%s (Convertido de mutualTLS)",
				scheme["description"])
		}
	}

	return nil
}

func (d *OpenAPIDowngrader) transformServers(document *map[string]interface{}) error {
	serversIface, exists := (*document)["servers"]
	if !exists {
		return nil
	}

	servers, ok := serversIface.([]interface{})
	if !ok {
		return nil
	}

	for _, serverIface := range servers {
		server, ok := serverIface.(map[string]interface{})
		if !ok {
			continue
		}

		if variablesIface, ok := server["variables"].(map[string]interface{}); ok {
			for _, varIface := range variablesIface {
				variable, ok := varIface.(map[string]interface{})
				if !ok {
					continue
				}

				delete(variable, "examples")
			}
		}
	}

	return nil
}

func processAllSchemas(document interface{}, processor func(map[string]interface{})) {
	switch v := document.(type) {
	case map[string]interface{}:
		if isSchema(v) {
			processor(v)
		}

		for _, value := range v {
			processAllSchemas(value, processor)
		}

	case []interface{}:
		for _, item := range v {
			processAllSchemas(item, processor)
		}
	}
}

func isSchema(obj map[string]interface{}) bool {
	schemaProperties := []string{
		"type", "properties", "items", "additionalProperties",
		"allOf", "oneOf", "anyOf", "not", "$ref",
	}

	for _, prop := range schemaProperties {
		if _, exists := obj[prop]; exists {
			return true
		}
	}

	return false
}
