package spec

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/viper"
)

const (
	VIPER_FILE = "specs.yaml"
	SPEC_DIR   = "cli_specs"
)

func verificarEAtualizarDiretorio(caminho string) error {
	_, err := os.Stat(caminho)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err := os.MkdirAll(caminho, 0755)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}

func validarEndpoint(url string) bool {
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

	fmt.Println("Endpoint válido.")
	return true
}

func getAndSaveFile(url, caminhoDestino string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("erro ao fazer o download do arquivo JSON: %v", err)
	}
	defer resp.Body.Close()

	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler o corpo da resposta: %v", err)
	}
	err = os.WriteFile(caminhoDestino, fileBytes, 0644)
	if err != nil {
		return fmt.Errorf("erro ao gravar os dados no arquivo: %v", err)
	}

	return nil
}

func loadList() ([]specList, error) {
	var currentConfig []specList
	config := viper.Get("jaxyendy")

	if config != nil {
		for _, v := range config.([]interface{}) {
			vv, ok := interfaceToMap(v)
			if !ok {
				return currentConfig, fmt.Errorf("fail to load current config")
			}
			currentConfig = append(currentConfig, specList{
				Url:     vv["url"].(string),
				Menu:    vv["menu"].(string),
				File:    vv["file"].(string),
				Enabled: vv["enabled"].(bool),
			})
		}

	}
	return currentConfig, nil
}
