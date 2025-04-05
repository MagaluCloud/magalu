package customizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
	"gopkg.in/yaml.v3"
)

// MagaluCustomizer é responsável por personalizar especificações OpenAPI para o padrão Magalu
type MagaluCustomizer struct {
}

// CustomizeOptions define as opções para personalização de especificações
type CustomizeOptions struct {
	IncludeRegion       bool
	IncludeGlobalRegion bool
	ProductPathURL      string
	DowngradeToVersion  string
	ParamsToRemove      []string
	ConfigureSecurity   bool // Nova opção para configurar segurança nas rotas
}

// NewMagaluCustomizer cria uma nova instância do customizador Magalu
func NewMagaluCustomizer() *MagaluCustomizer {
	return &MagaluCustomizer{}
}

// CustomizeSpec personaliza uma especificação OpenAPI para o padrão Magalu
func (m *MagaluCustomizer) CustomizeSpec(specPath, outputPath string, options *CustomizeOptions) error {
	// Carregar a especificação
	file := filepath.Join(specPath)
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de especificação: %v", err)
	}

	doc, err := libopenapi.NewDocument(fileBytes)
	if err != nil {
		return fmt.Errorf("erro ao carregar a especificação: %v", err)
	}

	// Personalizar o spec
	customizedSpec, err := m.applyCustomizations(doc, options)
	if err != nil {
		return fmt.Errorf("erro ao personalizar especificação: %v", err)
	}

	// Salvar a especificação personalizada
	err = m.saveSpec(customizedSpec, outputPath)
	if err != nil {
		return fmt.Errorf("erro ao salvar a especificação personalizada: %v", err)
	}

	return nil
}

// applyCustomizations aplica as personalizações específicas da Magalu
func (m *MagaluCustomizer) applyCustomizations(spec libopenapi.Document, options *CustomizeOptions) (libopenapi.Document, error) {
	// Criar extensão para as transformações
	createExtension := func(value interface{}) *yaml.Node {
		node := &yaml.Node{}
		err := node.Encode(value)
		if err != nil {
			fmt.Printf("Aviso: erro ao criar extensão: %v\n", err)
			return nil
		}
		return node
	}

	// Construir o modelo V3
	model, errs := spec.BuildV3Model()
	if errs != nil {
		return spec, fmt.Errorf("erro ao construir modelo V3: %v", errs)
	}

	// Definir versão OpenAPI se solicitado
	if options.DowngradeToVersion != "" {
		model.Model.Version = options.DowngradeToVersion
	}

	// Configurar URL do servidor
	url := "https://{env}/" + options.ProductPathURL
	if options.IncludeRegion {
		url = "https://{env}/{region}/" + options.ProductPathURL
	}

	server := &v3.Server{
		URL:       url,
		Variables: orderedmap.New[string, *v3.ServerVariable](),
	}

	// Configurar variável de região
	if options.IncludeRegion {
		regionVar := &v3.ServerVariable{
			Enum:       []string{"br-ne1", "br-se1", "br-mgl1"},
			Extensions: orderedmap.New[string, *yaml.Node](),
		}

		if options.IncludeGlobalRegion {
			regionVar.Enum = append(regionVar.Enum, "global")
		}

		regionVar.Description = "Region to reach the service"
		regionVar.Default = "br-se1"

		regionTransforms := []map[string]interface{}{
			{
				"type":         "translate",
				"allowMissing": true,
				"translations": []map[string]string{
					{"from": "br-mgl1", "to": "br-se-1"},
				},
			},
		}
		regionVar.Extensions.Set("x-mgc-transforms", createExtension(regionTransforms))
		server.Variables.Set("region", regionVar)
	}

	// Configurar variável de ambiente
	envVar := &v3.ServerVariable{
		Enum:       []string{"api.magalu.cloud", "api.pre-prod.jaxyendy.com"},
		Extensions: orderedmap.New[string, *yaml.Node](),
	}
	envVar.Description = "Environment to use"
	envVar.Default = "api.magalu.cloud"
	envTransforms := []map[string]interface{}{
		{
			"type": "translate",
			"translations": []map[string]string{
				{"from": "prod", "to": "api.magalu.cloud"},
				{"from": "pre-prod", "to": "api.pre-prod.jaxyendy.com"},
			},
		},
	}
	envVar.Extensions.Set("x-mgc-transforms", createExtension(envTransforms))
	server.Variables.Set("env", envVar)

	// Definir os servidores
	model.Model.Servers = []*v3.Server{server}

	// Configurar segurança nas rotas se solicitado
	if options.ConfigureSecurity {
		err := m.configureSecurity(model, options.ProductPathURL)
		if err != nil {
			return spec, fmt.Errorf("erro ao configurar segurança nas rotas: %v", err)
		}
	}

	// Remover parâmetros especificados na lista ParamsToRemove
	if options.ParamsToRemove != nil && len(options.ParamsToRemove) > 0 {
		err := m.removeParameters(model, options.ParamsToRemove)
		if err != nil {
			return spec, fmt.Errorf("erro ao remover parâmetros: %v", err)
		}
	}

	return spec, nil
}

// removeParameters remove parâmetros especificados em todas as operações
// Retorna o número total de parâmetros removidos
func (m *MagaluCustomizer) removeParameters(model *libopenapi.DocumentModel[v3.Document], paramsToRemove []string) error {
	if model == nil || model.Model.Paths == nil {
		return nil
	}

	// Criar mapa de parâmetros a serem removidos para busca eficiente (em lowercase)
	paramsMap := make(map[string]bool)
	for _, param := range paramsToRemove {
		paramsMap[strings.ToLower(param)] = true
	}

	// Estatísticas de remoção
	removedCount := make(map[string]int) // Contagem por parâmetro
	totalPathsAffected := 0              // Total de caminhos afetados
	totalOperationsAffected := 0         // Total de operações afetadas
	emptyPathsCount := 0                 // Caminhos que ficaram sem parâmetros
	emptyOperationsCount := 0            // Operações que ficaram sem parâmetros
	componentParamsRemoved := 0          // Parâmetros removidos da seção de componentes

	// Inicializar contador para cada parâmetro original (preservando o case original)
	for _, param := range paramsToRemove {
		removedCount[param] = 0
	}

	// Remover parâmetros da seção de componentes
	if model.Model.Components != nil && model.Model.Components.Parameters != nil {
		// Verificar cada parâmetro componente
		for name := model.Model.Components.Parameters.Oldest(); name != nil; {
			next := name.Next() // Guardar o próximo antes de potencialmente remover

			// Verificar se o nome do parâmetro está na lista para remover (case insensitive)
			paramNameLower := strings.ToLower(name.Key)
			if paramsMap[paramNameLower] {
				model.Model.Components.Parameters.Delete(name.Key)

				// Incrementar o contador do parâmetro original correspondente
				for _, originalParam := range paramsToRemove {
					if strings.ToLower(originalParam) == paramNameLower {
						removedCount[originalParam]++
						break
					}
				}

				componentParamsRemoved++
				fmt.Printf("Parâmetro '%s' removido da seção de componentes\n", name.Key)
			}

			name = next
		}

		// Se todos os parâmetros componentes foram removidos, definir como nil
		if model.Model.Components.Parameters.Len() == 0 {
			model.Model.Components.Parameters = nil
			fmt.Println("Seção de parâmetros em componentes foi completamente removida")
		}
	}

	// Percorrer todos os caminhos
	for path := model.Model.Paths.PathItems.Oldest(); path != nil; path = path.Next() {
		pathItem := path.Value
		pathAffected := false

		// Remover parâmetros no nível do caminho
		if pathItem.Parameters != nil {
			originalLen := len(pathItem.Parameters)
			pathItem.Parameters = m.filterParameters(pathItem.Parameters, paramsMap, removedCount, paramsToRemove)
			removedPathParams := originalLen - len(pathItem.Parameters)

			if removedPathParams > 0 {
				pathAffected = true
				totalPathsAffected++
			}

			// Se todos os parâmetros foram removidos, defina como nil
			if len(pathItem.Parameters) == 0 {
				pathItem.Parameters = nil
				emptyPathsCount++
				fmt.Printf("Todos os parâmetros do caminho '%s' foram removidos\n", path.Key)
			} else if removedPathParams > 0 {
				fmt.Printf("Removidos %d parâmetros do caminho '%s'\n", removedPathParams, path.Key)
			}
		}

		// Obter todas as operações deste caminho
		operations := pathItem.GetOperations()
		for op := operations.Oldest(); op != nil; op = op.Next() {
			operation := op.Value
			operationAffected := false

			// Remover parâmetros no nível da operação
			if operation.Parameters != nil {
				originalLen := len(operation.Parameters)
				operation.Parameters = m.filterParameters(operation.Parameters, paramsMap, removedCount, paramsToRemove)
				removedOpParams := originalLen - len(operation.Parameters)

				if removedOpParams > 0 {
					operationAffected = true
					totalOperationsAffected++
				}

				// Se todos os parâmetros foram removidos, defina como nil
				if len(operation.Parameters) == 0 {
					operation.Parameters = nil
					emptyOperationsCount++
					fmt.Printf("Todos os parâmetros da operação '%s' no caminho '%s' foram removidos\n",
						op.Key, path.Key)
				} else if removedOpParams > 0 {
					fmt.Printf("Removidos %d parâmetros da operação '%s' no caminho '%s'\n",
						removedOpParams, op.Key, path.Key)
				}
			}

			// Verificar se Request Body tem Content vazio e definir como nil se necessário
			m.cleanupRequestBody(operation, path.Key, op.Key)

			// Verificar se Response tem Content vazio e definir como nil se necessário
			m.cleanupResponses(operation, path.Key, op.Key)

			if operationAffected {
				// Adicionar log somente se foi afetado
			}
		}

		if pathAffected {
			// Adicionar log somente se foi afetado
		}
	}

	// Verificar se o componente está vazio e limpá-lo se necessário
	m.cleanupEmptyComponents(model)

	// Exibir estatísticas de remoção
	fmt.Printf("\n=== Resumo da remoção de parâmetros ===\n")
	fmt.Printf("Total de caminhos afetados: %d (ficaram vazios: %d)\n", totalPathsAffected, emptyPathsCount)
	fmt.Printf("Total de operações afetadas: %d (ficaram vazias: %d)\n", totalOperationsAffected, emptyOperationsCount)
	fmt.Printf("Parâmetros removidos de componentes: %d\n", componentParamsRemoved)

	fmt.Printf("\nParâmetros removidos:\n")
	for param, count := range removedCount {
		if count > 0 {
			fmt.Printf("- '%s': %d ocorrência(s)\n", param, count)
		} else {
			fmt.Printf("- '%s': nenhuma ocorrência encontrada\n", param)
		}
	}

	return nil
}

// cleanupRequestBody verifica se o Request Body está vazio e define como nil se necessário
func (m *MagaluCustomizer) cleanupRequestBody(operation *v3.Operation, pathKey, opKey string) {
	if operation.RequestBody != nil {
		// Verificar se content está vazio
		if operation.RequestBody.Content != nil && operation.RequestBody.Content.Len() == 0 {
			operation.RequestBody.Content = nil
		}

		// Se não tiver content e outras propriedades importantes, definir como nil
		if operation.RequestBody.Content == nil &&
			operation.RequestBody.Description == "" &&
			operation.RequestBody.Required == nil &&
			operation.RequestBody.Extensions.Len() == 0 {
			operation.RequestBody = nil
			fmt.Printf("RequestBody vazio da operação '%s' no caminho '%s' foi removido\n", opKey, pathKey)
		}
	}
}

// cleanupResponses verifica se as Responses estão vazias e define como nil se necessário
func (m *MagaluCustomizer) cleanupResponses(operation *v3.Operation, pathKey, opKey string) {
	if operation.Responses != nil {
		// Verificar se há códigos de resposta
		if operation.Responses.Codes.Len() == 0 &&
			operation.Responses.Default == nil &&
			operation.Responses.Extensions.Len() == 0 {
			operation.Responses = nil
			fmt.Printf("Responses vazias da operação '%s' no caminho '%s' foram removidas\n", opKey, pathKey)
		}
	}
}

// filterParameters filtra uma lista de parâmetros removendo os parâmetros especificados
// Atualiza o contador de parâmetros removidos
func (m *MagaluCustomizer) filterParameters(parameters []*v3.Parameter, paramsToRemove map[string]bool, removedCount map[string]int, originalParams []string) []*v3.Parameter {
	result := make([]*v3.Parameter, 0)

	for _, param := range parameters {
		paramNameLower := strings.ToLower(param.Name)
		if !paramsToRemove[paramNameLower] {
			result = append(result, param)
		} else {
			// Incrementar contador do parâmetro original correspondente
			for _, originalParam := range originalParams {
				if strings.ToLower(originalParam) == paramNameLower {
					removedCount[originalParam]++
					break
				}
			}
		}
	}

	return result
}

// saveSpec salva a especificação personalizada em um arquivo
func (m *MagaluCustomizer) saveSpec(spec libopenapi.Document, filename string) error {
	model, errs := spec.BuildV3Model()
	if errs != nil {
		return fmt.Errorf("erro ao gerar modelo do spec personalizado: %v", errs)
	}

	byteFile, err := model.Model.RenderJSON("  ")
	if err != nil {
		return fmt.Errorf("erro ao renderizar spec personalizado: %v", err)
	}

	err = os.WriteFile(filename, byteFile, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo: %v", err)
	}

	return nil
}

// cleanupEmptyComponents verifica se o componente está vazio e define como nil se necessário
func (m *MagaluCustomizer) cleanupEmptyComponents(model *libopenapi.DocumentModel[v3.Document]) {
	if model.Model.Components == nil {
		return
	}

	// Verificar se todos os componentes estão vazios
	isEmpty := true

	if model.Model.Components.Schemas != nil && model.Model.Components.Schemas.Len() > 0 {
		isEmpty = false
	}
	if model.Model.Components.Responses != nil && model.Model.Components.Responses.Len() > 0 {
		isEmpty = false
	}
	if model.Model.Components.Parameters != nil && model.Model.Components.Parameters.Len() > 0 {
		isEmpty = false
	}
	if model.Model.Components.Examples != nil && model.Model.Components.Examples.Len() > 0 {
		isEmpty = false
	}
	if model.Model.Components.RequestBodies != nil && model.Model.Components.RequestBodies.Len() > 0 {
		isEmpty = false
	}
	if model.Model.Components.Headers != nil && model.Model.Components.Headers.Len() > 0 {
		isEmpty = false
	}
	if model.Model.Components.SecuritySchemes != nil && model.Model.Components.SecuritySchemes.Len() > 0 {
		isEmpty = false
	}
	if model.Model.Components.Links != nil && model.Model.Components.Links.Len() > 0 {
		isEmpty = false
	}
	if model.Model.Components.Callbacks != nil && model.Model.Components.Callbacks.Len() > 0 {
		isEmpty = false
	}
	if model.Model.Components.Extensions != nil && model.Model.Components.Extensions.Len() > 0 {
		isEmpty = false
	}

	if isEmpty {
		model.Model.Components = nil
		fmt.Println("Seção de componentes completamente vazia, removida da especificação")
	}
}

// configureSecurity configura segurança nas rotas que não possuem
// Para rotas GET, HEAD e OPTIONS é usado o escopo "productPathURL:read", para as demais "productPathURL:write"
func (m *MagaluCustomizer) configureSecurity(model *libopenapi.DocumentModel[v3.Document], productPathURL string) error {
	if model == nil || model.Model.Paths == nil {
		return nil
	}

	// Contar quantas rotas já possuem e quantas foram configuradas
	var totalSecurity, addedSecurity int

	// Métodos de leitura (que não modificam recursos)
	readMethods := map[string]bool{
		"get":     true,
		"head":    true,
		"options": true,
	}

	// Percorrer todos os caminhos
	for path := model.Model.Paths.PathItems.Oldest(); path != nil; path = path.Next() {
		pathItem := path.Value

		// Obter todas as operações deste caminho
		operations := pathItem.GetOperations()
		for op := operations.Oldest(); op != nil; op = op.Next() {
			operation := op.Value

			// Verificar se a operação já possui definição de segurança
			if operation.Security != nil && len(operation.Security) > 0 {
				totalSecurity++
				continue
			}

			// Configurar segurança com base no tipo de operação
			var scope string
			opMethod := strings.ToLower(op.Key)

			if readMethods[opMethod] {
				scope = fmt.Sprintf("%s:read", productPathURL)
			} else {
				scope = fmt.Sprintf("%s:write", productPathURL)
			}

			sec := orderedmap.New[string, []string]()
			sec.Set("OAuth2", []string{scope})

			op.Value.Security = []*base.SecurityRequirement{
				{
					Requirements: sec,
				},
			}
			addedSecurity++

			fmt.Printf("Segurança adicionada à operação '%s' no caminho '%s' com escopo '%s'\n",
				op.Key, path.Key, scope)
		}
	}

	// Adicionar esquema de segurança OAuth2 aos componentes se não existir
	if model.Model.Components == nil {
		model.Model.Components = &v3.Components{
			SecuritySchemes: orderedmap.New[string, *v3.SecurityScheme](),
		}
	} else if model.Model.Components.SecuritySchemes == nil {
		model.Model.Components.SecuritySchemes = orderedmap.New[string, *v3.SecurityScheme]()
	}

	// Verificar se já existe o esquema OAuth2
	if _, exists := model.Model.Components.SecuritySchemes.Get("OAuth2"); !exists {
		oauthScheme := &v3.SecurityScheme{
			Type:        "oauth2",
			Description: "Segurança OAuth2 para autenticação na API Magalu Cloud",
		}

		model.Model.Components.SecuritySchemes.Set("OAuth2", oauthScheme)
		fmt.Printf("Esquema de segurança OAuth2 adicionado à especificação\n")
	}

	// Configurar segurança global padrão se não existir
	if model.Model.Security == nil || len(model.Model.Security) == 0 {
		// Criar segurança global com ambos os escopos
		readScope := fmt.Sprintf("%s:read", productPathURL)
		writeScope := fmt.Sprintf("%s:write", productPathURL)

		sec := orderedmap.New[string, []string]()
		sec.Set("OAuth2", []string{readScope, writeScope})

		model.Model.Security = []*base.SecurityRequirement{
			{
				Requirements: sec,
			},
		}
		fmt.Printf("Segurança global padrão adicionada à especificação\n")
	}

	// Exibir estatísticas
	fmt.Printf("\n=== Resumo da configuração de segurança ===\n")
	fmt.Printf("Total de operações que já possuíam segurança: %d\n", totalSecurity)
	fmt.Printf("Total de operações com segurança adicionada: %d\n", addedSecurity)

	return nil
}
