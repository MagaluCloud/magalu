package pipeline

/*
Comando gen-docs-magalu

Este comando converte a funcionalidade do script Python original para Golang usando Cobra.
Ele insere sidebar_position em arquivos markdown seguindo uma ordem específica:

Uso:
  pipeline gen-docs-magalu <diretório>

Exemplo:
  pipeline gen-docs-magalu docs/devops-tools/cli-mgc/commands-reference

Exemplo de saída:
  Atualizado docs/devops-tools/cli-mgc/commands-reference/help.md com sidebar_position: 0
  Atualizado docs/devops-tools/cli-mgc/commands-reference/list.md com sidebar_position: 1
  Atualizado docs/devops-tools/cli-mgc/commands-reference/create.md com sidebar_position: 2
  Atualizado docs/devops-tools/cli-mgc/commands-reference/delete.md com sidebar_position: 3

Funcionalidades:
- Percorre recursivamente o diretório especificado
- Identifica arquivos markdown (.md)
- Aplica ordem específica de sidebar_position:
  * help.md: posição 0
  * list.md: posição 1
  * create.md: posição 2
  * outros arquivos: posições subsequentes (3, 4, 5, ...)
- Insere ou atualiza frontmatter com sidebar_position
- Preserva conteúdo existente dos arquivos

Conversão do script Python original:
- os.walk() -> filepath.WalkDir()
- Manipulação de strings -> strings.Split() e strings.Join()
- Leitura/escrita de arquivos -> os.ReadFile() e os.WriteFile()
- Tratamento de erros -> error handling idiomático do Go
*/

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var genDocsMagaluCmd = &cobra.Command{
	Use:   "gen-docs-magalu",
	Short: "Gera documentação com sidebar_position para arquivos markdown",
	Long: `Insere sidebar_position em arquivos markdown seguindo uma ordem específica:
- help.md: posição 0
- list.md: posição 1  
- create.md: posição 2
- outros arquivos: posições subsequentes (3, 4, 5, ...)`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootDirectory := args[0]
		return insertSidebarPosition(rootDirectory)
	},
}

func insertSidebarPosition(rootDirectory string) error {
	return filepath.WalkDir(rootDirectory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Verifica se é um arquivo markdown
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		// Lista todos os arquivos markdown no diretório atual
		dir := filepath.Dir(path)
		entries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("erro ao ler diretório %s: %w", dir, err)
		}

		var mdFiles []string
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				mdFiles = append(mdFiles, entry.Name())
			}
		}

		if len(mdFiles) == 0 {
			return nil
		}

		// Classifica os arquivos de acordo com as regras
		orderedFiles := make([]struct {
			filename string
			position int
		}, 0)

		var unorderedFiles []string

		for _, mdFile := range mdFiles {
			switch mdFile {
			case "help.md":
				orderedFiles = append(orderedFiles, struct {
					filename string
					position int
				}{mdFile, 0})
			case "list.md":
				orderedFiles = append(orderedFiles, struct {
					filename string
					position int
				}{mdFile, 1})
			case "create.md":
				orderedFiles = append(orderedFiles, struct {
					filename string
					position int
				}{mdFile, 2})
			default:
				unorderedFiles = append(unorderedFiles, mdFile)
			}
		}

		// Ordena os arquivos restantes com índices subsequentes
		index := 3
		for _, mdFile := range unorderedFiles {
			orderedFiles = append(orderedFiles, struct {
				filename string
				position int
			}{mdFile, index})
			index++
		}

		// Atualiza cada arquivo com o índice apropriado
		for _, fileInfo := range orderedFiles {
			filePath := filepath.Join(dir, fileInfo.filename)

			// Lê o conteúdo do arquivo
			content, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("erro ao ler arquivo %s: %w", filePath, err)
			}

			// Verifica se já tem frontmatter
			contentStr := string(content)
			var newContent string

			if strings.HasPrefix(contentStr, "---") {
				// Se já tem frontmatter, verifica se sidebar_position já existe
				lines := strings.Split(contentStr, "\n")
				var newLines []string
				hasSidebarPosition := false
				inFrontmatter := false

				for _, line := range lines {
					trimmedLine := strings.TrimSpace(line)

					// Detecta início e fim do frontmatter
					if line == "---" {
						if !inFrontmatter {
							inFrontmatter = true
						} else {
							inFrontmatter = false
						}
						newLines = append(newLines, line)
						continue
					}

					// Se está no frontmatter e a linha já é sidebar_position, substitui
					if inFrontmatter && strings.HasPrefix(trimmedLine, "sidebar_position:") {
						newLines = append(newLines, fmt.Sprintf("sidebar_position: %d", fileInfo.position))
						hasSidebarPosition = true
					} else {
						newLines = append(newLines, line)
					}
				}

				// Se não encontrou sidebar_position no frontmatter, insere após a primeira "---"
				if !hasSidebarPosition {
					var finalLines []string
					inserted := false

					for _, line := range newLines {
						finalLines = append(finalLines, line)

						// Insere após a primeira linha "---"
						if line == "---" && !inserted {
							finalLines = append(finalLines, fmt.Sprintf("sidebar_position: %d", fileInfo.position))
							inserted = true
						}
					}

					newContent = strings.Join(finalLines, "\n")
				} else {
					newContent = strings.Join(newLines, "\n")
				}
			} else {
				// Se não tem frontmatter, cria um novo
				newContent = fmt.Sprintf("---\nsidebar_position: %d\n---\n%s", fileInfo.position, contentStr)
			}

			// Escreve o conteúdo atualizado
			err = os.WriteFile(filePath, []byte(newContent), 0644)
			if err != nil {
				return fmt.Errorf("erro ao escrever arquivo %s: %w", filePath, err)
			}

			// fmt.Printf("Atualizado %s com sidebar_position: %d\n", filePath, fileInfo.position)
		}

		return nil
	})
}

// GetGenDocsMagaluCmd retorna o comando para ser usado em outros lugares
func GetGenDocsMagaluCmd() *cobra.Command {
	return genDocsMagaluCmd
}
