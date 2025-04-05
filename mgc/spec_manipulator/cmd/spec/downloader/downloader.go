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

type SpecDownloader struct {
	HTTPClient   *http.Client
	ValidateSpec bool
}

func NewSpecDownloader() *SpecDownloader {
	return &SpecDownloader{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		ValidateSpec: true,
	}
}

func (d *SpecDownloader) DownloadSpec(source, destination string) error {
	destInfo, err := os.Stat(destination)
	if err == nil && destInfo.IsDir() {
		var fileName string

		if strings.HasPrefix(source, "@") {
			fileName = filepath.Base(source[1:])
		} else {
			fileName = getFileNameFromURL(source)
		}

		destination = filepath.Join(destination, fileName)
	}

	destinationDir := filepath.Dir(destination)
	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de destino: %v", err)
	}

	var content []byte

	if strings.HasPrefix(source, "@") {
		filePath := source[1:]
		content, err = os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("erro ao ler arquivo local '%s': %v", filePath, err)
		}

		fmt.Printf("Arquivo local carregado: %s (%d bytes)\n", filePath, len(content))
	} else {
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

	err = os.WriteFile(destination, content, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo em '%s': %v", destination, err)
	}

	return nil
}

func (d *SpecDownloader) validateOpenAPISpec(content []byte) (bool, string, error) {
	document, err := libopenapi.NewDocument(content)
	if err != nil {
		return false, "", fmt.Errorf("falha ao parsear a especificação: %v", err)
	}

	v3Model, v3Errs := document.BuildV3Model()
	if v3Errs == nil && v3Model != nil {
		return true, v3Model.Model.Version, nil
	}

	var swagger map[string]interface{}
	if err := json.Unmarshal(content, &swagger); err == nil {
		if version, ok := swagger["swagger"].(string); ok && version == "2.0" {
			return true, "2.0 (Swagger)", nil
		}
	}

	return false, "", nil
}

func getFileNameFromURL(url string) string {
	if idx := strings.Index(url, "?"); idx > 0 {
		url = url[:idx]
	}

	parts := strings.Split(url, "/")
	fileName := parts[len(parts)-1]

	if fileName == "" {
		fileName = "openapi.json"
	}

	return fileName
}
