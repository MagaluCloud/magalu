package downgrader

import (
	"fmt"
	"os"
	"strings"

	"github.com/pb33f/libopenapi/datamodel"
	"gopkg.in/yaml.v3"
)

// OpenAPIDowngrader gerencia o processo de downgrade de OpenAPI 3.1.0 para 3.0.3
type OpenAPIDowngrader struct {
	// Configuração opcional
	config *datamodel.DocumentConfiguration
}

// NewOpenAPIDowngrader cria uma nova instância do downgrader
func NewOpenAPIDowngrader() *OpenAPIDowngrader {
	return &OpenAPIDowngrader{
		config: &datamodel.DocumentConfiguration{
			BypassDocumentCheck: true,
		},
	}
}

// DowngradeFile realiza o downgrade de um arquivo OpenAPI 3.1.0 para 3.0.3
func (d *OpenAPIDowngrader) DowngradeFile(inputFile, outputFile string) error {
	// Ler o arquivo de entrada
	inputBytes, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de entrada: %v", err)
	}

	// Realizar o downgrade
	outputBytes, err := d.Downgrade(inputBytes)
	if err != nil {
		return err
	}

	// Salvar o arquivo de saída
	err = os.WriteFile(outputFile, outputBytes, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo de saída: %v", err)
	}

	return nil
}

// Downgrade converte bytes de um documento OpenAPI 3.1.0 para 3.0.3
func (d *OpenAPIDowngrader) Downgrade(inputBytes []byte) ([]byte, error) {
	// Extrair informações do spec
	specInfo, err := datamodel.ExtractSpecInfoWithConfig(inputBytes, d.config)
	if err != nil {
		return nil, fmt.Errorf("erro ao extrair informações do spec: %v", err)
	}

	// Verificar se é OpenAPI 3.x
	if specInfo.SpecType != "openapi" {
		return nil, fmt.Errorf("o arquivo não é uma especificação OpenAPI válida")
	}

	// // Verificar se é OpenAPI 3.1.x
	// if specInfo.VersionNumeric < 3.1 {
	// 	return nil, fmt.Errorf("o arquivo de entrada deve ser OpenAPI 3.1.x, mas é %s", specInfo.Version)
	// }

	// Converter para mapa JSON para manipulação
	document := make(map[string]interface{})
	if err := convertNodeToMap(specInfo.RootNode, &document); err != nil {
		return nil, fmt.Errorf("erro ao converter documento para JSON: %v", err)
	}

	// Aplicar transformações para downgrade
	if err := d.applyTransformations(&document); err != nil {
		return nil, fmt.Errorf("erro ao aplicar transformações: %v", err)
	}

	// Converter de volta para YAML
	outputBytes, err := yaml.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter para YAML: %v", err)
	}

	return outputBytes, nil
}

// convertNodeToMap converte um nó YAML em um mapa
func convertNodeToMap(node *yaml.Node, out interface{}) error {
	// Serializar para bytes e deserializar no formato correto
	bytes, err := yaml.Marshal(node)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, out)
}

// applyTransformations aplica todas as transformações necessárias para downgrade
func (d *OpenAPIDowngrader) applyTransformations(document *map[string]interface{}) error {
	// Alterar a versão para 3.0.3
	(*document)["openapi"] = "3.0.3"

	// Aplicar todas as transformações necessárias
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

// removeWebhooks remove a seção webhooks (não existente em 3.0.3)
func (d *OpenAPIDowngrader) removeWebhooks(document *map[string]interface{}) error {
	// Remover webhooks no nível raiz
	delete(*document, "webhooks")

	return nil
}

// transformJSONSchema transforma esquemas JSON para compatibilidade com 3.0.3
func (d *OpenAPIDowngrader) transformJSONSchema(document *map[string]interface{}) error {
	// Processar recursivamente todos os schemas
	processAllSchemas(*document, d.transformSchema)

	return nil
}

// transformSchema transforma um schema individual para compatibilidade com 3.0.3
func (d *OpenAPIDowngrader) transformSchema(schema map[string]interface{}) {
	// Remover propriedades não suportadas em 3.0.3
	deleteKeys := []string{
		"contentEncoding", "contentMediaType", "contentSchema", "examples",
		"unevaluatedProperties", "if", "then", "else",
		"dependentSchemas", "dependentRequired", "maxContains", "minContains",
		"const", "propertyNames",
	}
	for _, key := range deleteKeys {
		delete(schema, key)
	}

	// Converter type: null ou arrays com null para nullable: true
	if typeVal, ok := schema["type"]; ok {
		// Se type for uma string igual a "null"
		if typeStr, ok := typeVal.(string); ok && typeStr == "null" {
			delete(schema, "type")
			schema["nullable"] = true
		}

		// Se type for um array que contém "null"
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
					// OAS 3.0 não suporta múltiplos tipos, então usamos o primeiro
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

	// Converter prefixItems para items (usado em 3.1.0)
	if prefixItems, ok := schema["prefixItems"]; ok {
		if _, ok := schema["items"]; !ok {
			schema["items"] = prefixItems
		}
		delete(schema, "prefixItems")
	}

	// Remover propriedade exclusiveMinimum/exclusiveMaximum como booleanos
	// Em 3.0, eles são numéricos
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

// transformComponents transforma definições de componentes para compatibilidade com 3.0.3
func (d *OpenAPIDowngrader) transformComponents(document *map[string]interface{}) error {
	// Verificar se existem componentes
	componentsIface, exists := (*document)["components"]
	if !exists {
		return nil
	}

	components, ok := componentsIface.(map[string]interface{})
	if !ok {
		return nil
	}

	// Remover seções não suportadas em 3.0.3
	deleteKeys := []string{
		"pathItems", "webhooks",
	}
	for _, key := range deleteKeys {
		delete(components, key)
	}

	// Transformar exemplos para formato 3.0.3 se existirem
	if examples, ok := components["examples"].(map[string]interface{}); ok {
		for name, example := range examples {
			if exampleMap, ok := example.(map[string]interface{}); ok {
				// Verificar o formato e converter se necessário
				if _, hasValue := exampleMap["value"]; !hasValue {
					// Criar estrutura compatível com 3.0.3
					newExample := map[string]interface{}{
						"value": exampleMap,
					}

					// Preservar summary e description
					if summary, ok := exampleMap["summary"]; ok {
						newExample["summary"] = summary
						delete(exampleMap, "summary")
					}
					if description, ok := exampleMap["description"]; ok {
						newExample["description"] = description
						delete(exampleMap, "description")
					}

					// Atualizar o exemplo
					examples[name] = newExample
				}
			}
		}
	}

	return nil
}

// transformPaths transforma definições de caminhos para compatibilidade com 3.0.3
func (d *OpenAPIDowngrader) transformPaths(document *map[string]interface{}) error {
	// Verificar se existem paths
	pathsIface, exists := (*document)["paths"]
	if !exists {
		return nil
	}

	paths, ok := pathsIface.(map[string]interface{})
	if !ok {
		return nil
	}

	// Procurar operações em cada caminho
	for _, pathItem := range paths {
		if pathItemMap, ok := pathItem.(map[string]interface{}); ok {
			// Processar cada operação HTTP
			operations := []string{"get", "post", "put", "delete", "options", "head", "patch", "trace"}
			for _, op := range operations {
				if operation, ok := pathItemMap[op].(map[string]interface{}); ok {
					// Transformar requestBody se existir
					if requestBody, ok := operation["requestBody"].(map[string]interface{}); ok {
						// Transformar content se existir
						if content, ok := requestBody["content"].(map[string]interface{}); ok {
							// Processar cada tipo de mídia
							for _, mediaTypeObj := range content {
								if mediaType, ok := mediaTypeObj.(map[string]interface{}); ok {
									// Transformar exemplos para formato 3.0.3
									d.transformMediaTypeExamples(mediaType)
								}
							}
						}
					}

					// Transformar respostas se existirem
					if responses, ok := operation["responses"].(map[string]interface{}); ok {
						for _, responseObj := range responses {
							if response, ok := responseObj.(map[string]interface{}); ok {
								// Transformar content se existir
								if content, ok := response["content"].(map[string]interface{}); ok {
									// Processar cada tipo de mídia
									for _, mediaTypeObj := range content {
										if mediaType, ok := mediaTypeObj.(map[string]interface{}); ok {
											// Transformar exemplos para formato 3.0.3
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

// transformMediaTypeExamples transforma exemplos em mediaType para formato 3.0.3
func (d *OpenAPIDowngrader) transformMediaTypeExamples(mediaType map[string]interface{}) {
	// Verificar se existem exemplos
	if examples, ok := mediaType["examples"].(map[string]interface{}); ok {
		// Converter exemplos individuais para exemplo único se possível
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

// transformInfo transforma informações do documento para compatibilidade com 3.0.3
func (d *OpenAPIDowngrader) transformInfo(document *map[string]interface{}) error {
	// Verificar se existe seção info
	infoIface, exists := (*document)["info"]
	if !exists {
		return nil
	}

	info, ok := infoIface.(map[string]interface{})
	if !ok {
		return nil
	}

	// Remover propriedades não suportadas em 3.0.3
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

// transformSecuritySchemes transforma esquemas de segurança para compatibilidade com 3.0.3
func (d *OpenAPIDowngrader) transformSecuritySchemes(document *map[string]interface{}) error {
	// Verificar se existem componentes
	componentsIface, exists := (*document)["components"]
	if !exists {
		return nil
	}

	components, ok := componentsIface.(map[string]interface{})
	if !ok {
		return nil
	}

	// Verificar se existem esquemas de segurança
	securitySchemesIface, exists := components["securitySchemes"]
	if !exists {
		return nil
	}

	securitySchemes, ok := securitySchemesIface.(map[string]interface{})
	if !ok {
		return nil
	}

	// Processar cada esquema de segurança
	for _, schemeIface := range securitySchemes {
		scheme, ok := schemeIface.(map[string]interface{})
		if !ok {
			continue
		}

		// Converter schemeType "mutualTLS" (3.1.0) para formato compatível
		if schemeType, ok := scheme["type"].(string); ok && schemeType == "mutualTLS" {
			// Não há equivalente direto em 3.0.3, converter para http
			scheme["type"] = "http"
			scheme["scheme"] = "https"
			scheme["description"] = fmt.Sprintf("%s (Convertido de mutualTLS)",
				scheme["description"])
		}
	}

	return nil
}

// transformServers transforma servidores para compatibilidade com 3.0.3
func (d *OpenAPIDowngrader) transformServers(document *map[string]interface{}) error {
	// Verificar se existem servidores
	serversIface, exists := (*document)["servers"]
	if !exists {
		return nil
	}

	servers, ok := serversIface.([]interface{})
	if !ok {
		return nil
	}

	// Processar cada servidor
	for _, serverIface := range servers {
		server, ok := serverIface.(map[string]interface{})
		if !ok {
			continue
		}

		// Verificar e ajustar variáveis se existirem
		if variablesIface, ok := server["variables"].(map[string]interface{}); ok {
			for _, varIface := range variablesIface {
				variable, ok := varIface.(map[string]interface{})
				if !ok {
					continue
				}

				// Remover propriedades não suportadas em 3.0.3
				delete(variable, "examples")
			}
		}
	}

	return nil
}

// processAllSchemas processa todos os esquemas no documento
func processAllSchemas(document interface{}, processor func(map[string]interface{})) {
	switch v := document.(type) {
	case map[string]interface{}:
		// Verificar se este objeto é um schema
		if isSchema(v) {
			processor(v)
		}

		// Processar recursivamente todos os valores
		for _, value := range v {
			processAllSchemas(value, processor)
		}

	case []interface{}:
		// Processar cada item do array recursivamente
		for _, item := range v {
			processAllSchemas(item, processor)
		}
	}
}

// isSchema verifica se um objeto é um esquema JSON
func isSchema(obj map[string]interface{}) bool {
	// Verificar se tem propriedades que indicam que é um schema
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
