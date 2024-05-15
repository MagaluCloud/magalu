package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func verificarEAtualizarDiretorio(caminho string) error {
	// Verifica se o diretório já existe
	_, err := os.Stat(caminho)
	if err == nil {
		// O diretório já existe
		return nil
	}
	if os.IsNotExist(err) {
		// O diretório não existe, então tentamos criar
		err := os.MkdirAll(caminho, 0755) // 0755 é o modo padrão de permissão para diretórios
		if err != nil {
			return err
		}
		return nil
	}
	// Se ocorrer algum outro erro ao verificar o diretório, retorna o erro
	return err
}

func verificarERenomearArquivo(caminho string) error {
	// Verifica se o arquivo já existe
	_, err := os.Stat(caminho)
	if err != nil {
		if os.IsNotExist(err) {
			// Arquivo não existe
			return nil
		}
		// Outro erro ao verificar o arquivo
		return err
	}

	// Obtém a data de criação do arquivo
	info, err := os.Stat(caminho)
	if err != nil {
		return err
	}
	dataCriacao := info.ModTime()
	dataCriacaoFormatada := dataCriacao.Format("2006-01-02_15-04-05")

	// Obtém o nome e a extensão do arquivo
	nomeArquivo := filepath.Base(caminho)
	extensao := filepath.Ext(caminho)
	nomeArquivoSemExtensao := nomeArquivo[0 : len(nomeArquivo)-len(extensao)]

	// Renomeia o arquivo para incluir a data de criação
	novoNome := fmt.Sprintf("%s_%s.old%s", nomeArquivoSemExtensao, dataCriacaoFormatada, extensao)
	novoCaminho := filepath.Join(filepath.Dir(caminho), novoNome)
	err = os.Rename(caminho, novoCaminho)
	if err != nil {
		return err
	}

	fmt.Printf("Arquivo renomeado para: %s\n", novoCaminho)
	return nil
}

func removerArquivosOld(diretorio string) error {
	// Abre o diretório especificado
	dir, err := os.Open(diretorio)
	if err != nil {
		return fmt.Errorf("erro ao abrir o diretório: %v", err)
	}
	defer dir.Close()

	// Lê o conteúdo do diretório
	arquivos, err := dir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("erro ao ler o conteúdo do diretório: %v", err)
	}

	// Itera sobre os arquivos do diretório
	for _, arquivo := range arquivos {
		// Verifica se é um arquivo com extensão ".old"
		if !arquivo.IsDir() && filepath.Ext(arquivo.Name()) == ".old" {
			// Monta o caminho completo do arquivo
			caminhoArquivo := filepath.Join(diretorio, arquivo.Name())

			// Remove o arquivo
			err := os.Remove(caminhoArquivo)
			if err != nil {
				return fmt.Errorf("erro ao remover o arquivo %s: %v", caminhoArquivo, err)
			}

			fmt.Printf("Arquivo %s removido com sucesso.\n", caminhoArquivo)
		}
	}

	return nil
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

func getAndSaveFromGitlab(url, caminhoDestino string) error {

	destine := strings.Replace(url, "https://gitlab.luizalabs.com/open-platform/pcx/u0/-/raw/main/api_products", "", 1)
	destine = strings.Replace(destine, "?ref_type=heads", "", -1)
	destine = strings.Replace(destine, "/", "%2F", -1)

	url = "https://gitlab.luizalabs.com/api/v4/projects/7739/repository/files/api_products"
	url += destine + "/raw?ref=main"

	// Faz o download do arquivo JSON

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar a solicitação HTTP: %v", err)
	}

	// Define o cabeçalho de autenticação com o token fornecido

	req.Header.Set("PRIVATE-TOKEN", os.Getenv("TOKEN_GITLAB"))

	// Envia a solicitação HTTP
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer a solicitação HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	dados, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler o corpo da resposta: %v", err)
	}

	// Grava os dados no arquivo local
	err = os.WriteFile(caminhoDestino, dados, 0644)
	if err != nil {
		return fmt.Errorf("erro ao gravar os dados no arquivo: %v", err)
	}

	fmt.Println("Arquivo JSON baixado e salvo com sucesso.")
	return nil
}

func getAndSaveFile(url, caminhoDestino string) error {
	// Faz o download do arquivo JSON
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("erro ao fazer o download do arquivo JSON: %v", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	dados, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler o corpo da resposta: %v", err)
	}

	// Grava os dados no arquivo local
	err = os.WriteFile(caminhoDestino, dados, 0644)
	if err != nil {
		return fmt.Errorf("erro ao gravar os dados no arquivo: %v", err)
	}

	fmt.Println("Arquivo JSON baixado e salvo com sucesso.")
	return nil
}

func loadListFromViper(origin string) ([]specList, error) {
	var currentConfig []specList
	fromViper := viper.Get(origin)

	if fromViper != nil {
		for _, v := range fromViper.([]interface{}) {
			vv, ok := interfaceToMap(v)
			if !ok {
				return currentConfig, fmt.Errorf("fail to load current config")
			}
			currentConfig = append(currentConfig, specList{Url: vv["url"].(string), Menu: vv["menu"].(string), File: vv["file"].(string)})
		}

	}
	return currentConfig, nil
}

func toWriteViveiro(isViveiro bool) string {
	if isViveiro {
		return "viveiro"
	}
	return "jaxyendy"
}
