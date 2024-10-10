package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Command struct {
	ID      int    `yaml:"id"`
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
}

type Config struct {
	Commands []Command `yaml:"commands"`
}

const (
	snapshotDir  = "snapshot"
	commandsYaml = "commands.yaml"
)

var (
	rewriteSnap bool
)

type resultError struct {
	Command
	Error string
}

type result struct {
	errors  []resultError
	success []Command
}

func main() {
	flag.BoolVar(&rewriteSnap, "bool", false, "Rewrite all failed snapshots")
	flag.Parse()

	yamlFile, err := os.ReadFile(commandsYaml)
	if err != nil {
		fmt.Printf("Erro ao ler o arquivo YAML: %v\n", err)
		return
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("Erro ao parsear o YAML: %v\n", err)
		return
	}

	_ = EnsureDirectoryExists(snapshotDir)

	var result result

	for _, cmd := range config.Commands {
		// fmt.Printf(`Executando comando: "%v - %s" -> $%s`, cmd.ID, cmd.Name, cmd.Command)

		output, err := exec.Command("sh", "-c", cmd.Command).CombinedOutput()

		if rewriteSnap {
			_ = WriteSnapshot(output, snapshotDir, cmd.ID)
		}

		if err == nil {
			err = CompareSnapshot(output, snapshotDir, cmd.ID)
		}

		if err != nil {
			result.errors = append(result.errors, resultError{
				Command: cmd,
				Error:   err.Error(),
			})
			continue
		}

		result.success = append(result.success, cmd)
	}

	// TODO- Tratar melhor os results
	for _, x := range result.errors {
		fmt.Println(x.Command)
		fmt.Println(x.Error)

	}

}
func EnsureDirectoryExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

func CreateFile(content []byte, dir, filePath string) error {
	return os.WriteFile(filepath.Join(dir, filePath), content, 0644)
}

func LoadFile(dir, filePath string) ([]byte, error) {
	return os.ReadFile(filepath.Join(dir, filePath))
}

func WriteSnapshot(output []byte, dir string, id int) error {
	_ = CreateFile(output, dir, fmt.Sprintf("%v.cli", id))
	return nil
}

func CompareBytes(expected, actual []byte) error {
	if bytes.Equal(expected, actual) {
		return nil
	}

	expectedLines := strings.Split(string(expected), "\n")
	actualLines := strings.Split(string(actual), "\n")

	var diff strings.Builder
	diff.WriteString("\nDiferenças encontradas:\n")

	i, j := 0, 0
	for i < len(expectedLines) && j < len(actualLines) {
		if expectedLines[i] == actualLines[j] {
			diff.WriteString("  " + expectedLines[i] + "\n")
			i++
			j++
		} else {
			diff.WriteString("- " + expectedLines[i] + "\n")
			diff.WriteString("+ " + actualLines[j] + "\n")
			i++
			j++
		}
	}

	// Adicionar linhas restantes, se houver
	for ; i < len(expectedLines); i++ {
		diff.WriteString("- " + expectedLines[i] + "\n")
	}
	for ; j < len(actualLines); j++ {
		diff.WriteString("+ " + actualLines[j] + "\n")
	}

	return fmt.Errorf("%s", diff.String())
}

func CompareSnapshot(output []byte, dir string, id int) error {
	snapContent, err := LoadFile(dir, fmt.Sprintf("%v.cli", id))
	if err == nil {
		return CompareBytes(snapContent, output)
	}

	if errors.Is(err, os.ErrNotExist) {
		_ = WriteSnapshot(output, dir, id)
		return nil
	}

	return fmt.Errorf("Diosmio")
}
