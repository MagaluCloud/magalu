package spec

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	VIPER_FILE    = "specs.yaml"
	SPEC_DIR      = "cli_specs"
	minRetryWait  = 1 * time.Second
	maxRetryWait  = 10 * time.Second
	maxRetryCount = 5
)

func verificarEAtualizarDiretorio(caminho string) error {
	_, err := os.Stat(caminho)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(caminho, 0755) // 0755 é o modo padrão de permissão para diretórios
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func validarEndpoint(url string) bool {
	if strings.HasPrefix(url, "http") {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Erro ao acessar o endpoint: %v\n", err)
			return false
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Erro: Status code não OK: %d\n", resp.StatusCode)
			return false
		}

		return true
	}
	if _, err := validateGitlabUrl(url); err != nil {
		return false
	}
	return true
}

func getAndSaveFile(url, caminhoDestino, menu string) error {
	var err error
	var resp *http.Response

	for i := 0; i < maxRetryCount; i++ {
		resp, err = http.Get(url)
		if err != nil {
			wait := time.Duration(math.Pow(2, float64(i))) * minRetryWait
			if wait > maxRetryWait {
				wait = maxRetryWait
			}
			fmt.Printf("Erro ao fazer download do arquivo %s, tentando novamente em %s\n", menu, wait)
			time.Sleep(wait)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			wait := time.Duration(math.Pow(2, float64(i))) * minRetryWait
			if wait > maxRetryWait {
				wait = maxRetryWait
			}
			fmt.Printf("Erro ao fazer download do arquivo %s (status %d), tentando novamente em %s\n", menu, resp.StatusCode, wait)
			time.Sleep(wait)
			continue
		}

		fileBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("erro ao ler o corpo da resposta %s: %v", menu, err)
		}
		err = os.WriteFile(caminhoDestino, fileBytes, 0644)
		if err != nil {
			return fmt.Errorf("erro ao gravar os dados no arquivo %s: %v", menu, err)
		}

		return nil
	}

	if err != nil {
		return fmt.Errorf("erro ao fazer o download do arquivo %s após %d tentativas: %v", menu, maxRetryCount, err)
	}
	return fmt.Errorf("erro ao fazer o download do arquivo %s após %d tentativas: status code %d", menu, maxRetryCount, resp.StatusCode)
}
