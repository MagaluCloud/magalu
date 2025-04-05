package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pb33f/libopenapi"
)

// SpecDownloader é responsável por baixar especificações OpenAPI
type SpecDownloader struct {
	// Opcionalmente, podemos adicionar configurações aqui, como:
	// - Timeout para o download
	// - Headers customizados
	// - Configurações de proxy
	// etc.
	HTTPClient *http.Client
	// Verifica se o arquivo é uma especificação OpenAPI válida
	ValidateSpec bool
}

// NewSpecDownloader cria uma nova instância do downloader de specs
func NewSpecDownloader() *SpecDownloader {
	return &SpecDownloader{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		ValidateSpec: true,
	}
}

// DownloadSpec baixa uma especificação OpenAPI de uma URL ou carrega de um arquivo local
// e salva no destino especificado
func (d *SpecDownloader) DownloadSpec(source, destination string) error {
	// Verificar se o destino é um diretório
	destInfo, err := os.Stat(destination)
	if err == nil && destInfo.IsDir() {
		// Se o destino é um diretório, precisamos extrair o nome do arquivo da fonte
		var fileName string

		if strings.HasPrefix(source, "@") {
			// Para caminhos de arquivo local, usar o nome do arquivo original
			fileName = filepath.Base(source[1:])
		} else {
			// Para URLs, extrair o nome do arquivo da URL
			fileName = getFileNameFromURL(source)
		}

		destination = filepath.Join(destination, fileName)
	}

	// Criar o diretório de destino se não existir
	destinationDir := filepath.Dir(destination)
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de destino: %v", err)
	}

	var content []byte

	// Verificar se a fonte é um arquivo local (prefixo @)
	if strings.HasPrefix(source, "@") {
		// Remover o prefixo @ e carregar o arquivo
		filePath := source[1:]
		content, err = os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("erro ao ler arquivo local '%s': %v", filePath, err)
		}

		fmt.Printf("Arquivo local carregado: %s (%d bytes)\n", filePath, len(content))
	} else {
		// Baixar a especificação da URL
		resp, err := d.HTTPClient.Get(source)
		if err != nil {
			return fmt.Errorf("erro ao baixar de '%s': %v", source, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("falha ao baixar de '%s': status %d", source, resp.StatusCode)
		}

		content, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("erro ao ler conteúdo de '%s': %v", source, err)
		}

		fmt.Printf("Arquivo baixado de: %s (%d bytes)\n", source, len(content))
	}

	// Validar se é uma especificação OpenAPI válida
	if d.ValidateSpec && len(content) > 0 {
		isValid, specVersion, err := d.validateOpenAPISpec(content)
		if err != nil {
			return fmt.Errorf("erro ao validar especificação: %v", err)
		}

		if !isValid {
			return fmt.Errorf("o arquivo não parece ser uma especificação OpenAPI válida")
		}

		fmt.Printf("Versão da especificação OpenAPI: %s\n", specVersion)
	}

	// Salvar o conteúdo no destino
	err = os.WriteFile(destination, content, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo em '%s': %v", destination, err)
	}

	return nil
}

// validateOpenAPISpec verifica se o conteúdo é uma especificação OpenAPI válida
// e retorna a versão da especificação
func (d *SpecDownloader) validateOpenAPISpec(content []byte) (bool, string, error) {
	// Tentar parsear a especificação usando libopenapi
	document, err := libopenapi.NewDocument(content)
	if err != nil {
		return false, "", fmt.Errorf("falha ao parsear a especificação: %v", err)
	}

	// Verificar se é uma especificação OpenAPI 3.x
	v3Model, v3Errs := document.BuildV3Model()
	if v3Errs == nil && v3Model != nil {
		return true, v3Model.Model.Version, nil
	}

	// Verificar se é uma especificação Swagger 2.0
	var swagger map[string]interface{}
	if err := json.Unmarshal(content, &swagger); err == nil {
		if version, ok := swagger["swagger"].(string); ok && version == "2.0" {
			return true, "2.0 (Swagger)", nil
		}
	}

	// Se chegou aqui, não é uma especificação OpenAPI/Swagger válida
	return false, "", nil
}

// Extrai o nome do arquivo de uma URL
func getFileNameFromURL(url string) string {
	// Remover parâmetros de query string
	if idx := strings.Index(url, "?"); idx > 0 {
		url = url[:idx]
	}

	// Obter o último segmento da URL
	parts := strings.Split(url, "/")
	fileName := parts[len(parts)-1]

	// Se o nome do arquivo estiver vazio, usar um nome padrão
	if fileName == "" {
		fileName = "openapi.json"
	}

	return fileName
}
