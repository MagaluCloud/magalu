package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

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

// func verificarERenomearArquivo(caminho string) error {
// 	// Verifica se o arquivo já existe
// 	_, err := os.Stat(caminho)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			// Arquivo não existe
// 			return nil
// 		}
// 		// Outro erro ao verificar o arquivo
// 		return err
// 	}

// 	// Obtém a data de criação do arquivo
// 	info, err := os.Stat(caminho)
// 	if err != nil {
// 		return err
// 	}
// 	dataCriacao := info.ModTime()
// 	dataCriacaoFormatada := dataCriacao.Format("2006-01-02_15-04-05")

// 	// Obtém o nome e a extensão do arquivo
// 	nomeArquivo := filepath.Base(caminho)
// 	extensao := filepath.Ext(caminho)
// 	nomeArquivoSemExtensao := nomeArquivo[0 : len(nomeArquivo)-len(extensao)]

// 	// Renomeia o arquivo para incluir a data de criação
// 	novoNome := fmt.Sprintf("%s_%s.old%s", nomeArquivoSemExtensao, dataCriacaoFormatada, extensao)
// 	novoCaminho := filepath.Join(filepath.Dir(caminho), novoNome)
// 	err = os.Rename(caminho, novoCaminho)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("Arquivo renomeado para: %s\n", novoCaminho)
// 	return nil
// }

// func removerArquivosOld(diretorio string) error {
// 	// Abre o diretório especificado
// 	dir, err := os.Open(diretorio)
// 	if err != nil {
// 		return fmt.Errorf("erro ao abrir o diretório: %v", err)
// 	}
// 	defer dir.Close()

// 	// Lê o conteúdo do diretório
// 	arquivos, err := dir.Readdir(-1)
// 	if err != nil {
// 		return fmt.Errorf("erro ao ler o conteúdo do diretório: %v", err)
// 	}

// 	// Itera sobre os arquivos do diretório
// 	for _, arquivo := range arquivos {
// 		// Verifica se é um arquivo com extensão ".old"
// 		if !arquivo.IsDir() && filepath.Ext(arquivo.Name()) == ".old" {
// 			// Monta o caminho completo do arquivo
// 			caminhoArquivo := filepath.Join(diretorio, arquivo.Name())

// 			// Remove o arquivo
// 			err := os.Remove(caminhoArquivo)
// 			if err != nil {
// 				return fmt.Errorf("erro ao remover o arquivo %s: %v", caminhoArquivo, err)
// 			}

// 			fmt.Printf("Arquivo %s removido com sucesso.\n", caminhoArquivo)
// 		}
// 	}

// 	return nil
// }

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
	// Faz o download do arquivo JSON
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("erro ao fazer o download do arquivo JSON: %v", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler o corpo da resposta: %v", err)
	}
	// Grava os dados no arquivo local
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
				CLI:     vv["cli"].(bool),
				TF:      vv["tf"].(bool),
				SDK:     vv["sdk"].(bool),
			})
		}

	}
	return currentConfig, nil
}
