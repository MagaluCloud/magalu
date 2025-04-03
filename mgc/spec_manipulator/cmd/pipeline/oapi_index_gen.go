package pipeline

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type IndexModule struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	Path        string `yaml:"path"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Summary     string `yaml:"summary,omitempty"`
}

type IndexFile struct {
	Version string        `yaml:"version"`
	Modules []IndexModule `yaml:"modules"`
}

const (
	indexFilename = "index.openapi.yaml"
	indexVersion  = "1.0.0"
)

var modnameRegex = regexp.MustCompile(`^(?P<name>[a-z0-9-]+)\.openapi\.yaml$`)

func NewOAPIIndexCommand() *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "oapi-index [diretório]",
		Short: "Gera arquivo de índice para todos os arquivos OAPI YAML no diretório",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputDir := args[0]
			if outputDir == "" {
				outputDir = inputDir
			}

			mods, err := loadModules(inputDir, outputDir)
			if err != nil {
				return err
			}

			return saveIndex(mods, outputDir)
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", "", "Diretório para salvar o novo índice YAML. Padrão é o diretório de entrada")

	return cmd
}

func loadModules(oapiDir, outDir string) ([]IndexModule, error) {
	files, err := os.ReadDir(oapiDir)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler diretório: %w", err)
	}

	var modules []IndexModule
	for _, file := range files {
		if file.IsDir() || file.Name() == indexFilename {
			continue
		}

		matches := modnameRegex.FindStringSubmatch(file.Name())
		if matches == nil {
			fmt.Printf("arquivo ignorado: %s\n", file.Name())
			continue
		}

		fullPath := filepath.Join(oapiDir, file.Name())
		relpath, err := filepath.Rel(outDir, fullPath)
		if err != nil {
			return nil, fmt.Errorf("erro ao obter caminho relativo: %w", err)
		}

		data, err := loadYAML(fullPath)
		if err != nil {
			return nil, fmt.Errorf("erro ao carregar YAML %s: %w", file.Name(), err)
		}

		info, ok := data["info"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("campo 'info' inválido em %s", file.Name())
		}

		description := ""
		if desc, ok := info["x-mgc-description"].(string); ok {
			description = desc
		} else if desc, ok := info["description"].(string); ok {
			description = desc
		}

		summary := description
		if sum, ok := info["summary"].(string); ok {
			summary = sum
		}

		version := ""
		if ver, ok := info["version"].(string); ok {
			version = ver
		}

		url := ""
		if id, ok := data["$id"].(string); ok {
			url = id
		}

		modules = append(modules, IndexModule{
			Name:        matches[1],
			URL:         url,
			Path:        relpath,
			Description: description,
			Version:     version,
			Summary:     summary,
		})
	}

	sort.Slice(modules, func(i, j int) bool {
		return modules[i].Name < modules[j].Name
	})

	return modules, nil
}

func loadYAML(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func saveIndex(modules []IndexModule, outDir string) error {
	indexFile := IndexFile{
		Version: indexVersion,
		Modules: modules,
	}

	data, err := yaml.Marshal(indexFile)
	if err != nil {
		return fmt.Errorf("erro ao serializar índice: %w", err)
	}

	outPath := filepath.Join(outDir, indexFilename)
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		return fmt.Errorf("erro ao salvar arquivo de índice: %w", err)
	}

	return nil
}
