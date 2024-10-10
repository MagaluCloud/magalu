package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

type commandsList struct {
	Module  string `yaml:"module"`
	Command string `yaml:"command"`
}

func interfaceToMap(i interface{}) (map[string]interface{}, bool) {
	mapa, ok := i.(map[string]interface{})
	if !ok {
		fmt.Println("A interface não é um mapa ou mapa de interfaces.")
		return nil, false
	}
	return mapa, true
}

func loadList() ([]commandsList, error) {
	var currentConfig []commandsList
	config := viper.Get("commands")

	if config != nil {
		for _, v := range config.([]interface{}) {
			vv, ok := interfaceToMap(v)
			if !ok {
				return currentConfig, fmt.Errorf("fail to load current config")
			}
			currentConfig = append(currentConfig, commandsList{
				Module:  vv["module"].(string),
				Command: vv["command"].(string),
			})
		}

	}
	return currentConfig, nil
}

func ensureDirectoryExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

func createFile(content []byte, dir, filePath string) error {
	return os.WriteFile(filepath.Join(dir, filePath), content, 0644)
}

func loadFile(dir, filePath string) ([]byte, error) {
	return os.ReadFile(filepath.Join(dir, filePath))
}

func writeSnapshot(output []byte, dir string, id string) error {
	_ = createFile(output, dir, fmt.Sprintf("%s.cli", id))
	return nil
}

func compareBytes(expected, actual []byte, ignoreDateUUID bool) error {
	if bytes.Equal(expected, actual) {
		return nil
	}

	allEqual := true

	expectedLines := strings.Split(string(expected), "\n")
	actualLines := strings.Split(string(actual), "\n")

	var diff strings.Builder
	diff.WriteString("\nDiferenças encontradas:\n")

	dateRegex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z`)
	uuidRegex := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

	i, j := 0, 0
	for i < len(expectedLines) && j < len(actualLines) {
		expectedLine := expectedLines[i]
		actualLine := actualLines[j]

		if ignoreDateUUID {
			expectedLine = dateRegex.ReplaceAllString(expectedLine, "DATE")
			expectedLine = uuidRegex.ReplaceAllString(expectedLine, "UUID")
			actualLine = dateRegex.ReplaceAllString(actualLine, "DATE")
			actualLine = uuidRegex.ReplaceAllString(actualLine, "UUID")
		}

		if expectedLine == actualLine {
			diff.WriteString("  " + expectedLines[i] + "\n")
			i++
			j++
		} else {
			allEqual = false
			diff.WriteString("- " + expectedLines[i] + "\n")
			diff.WriteString("+ " + actualLines[j] + "\n")
			i++
			j++
		}
	}

	for ; i < len(expectedLines); i++ {
		diff.WriteString("- " + expectedLines[i] + "\n")
	}
	for ; j < len(actualLines); j++ {
		diff.WriteString("+ " + actualLines[j] + "\n")
	}

	if allEqual {
		return nil
	}

	return fmt.Errorf("%s", diff.String())
}

func compareSnapshot(output []byte, dir string, id string) error {
	snapContent, err := loadFile(dir, fmt.Sprintf("%s.cli", id))
	if err == nil {
		return compareBytes(snapContent, output, true)
	}

	if errors.Is(err, os.ErrNotExist) {
		_ = writeSnapshot(output, dir, id)
		return nil
	}

	return fmt.Errorf("Diosmio")
}

func normalizeCommandToFile(input string) string {
	words := strings.Fields(input)
	var filteredWords []string
	for _, word := range words {
		if !strings.HasPrefix(word, "--") {
			filteredWords = append(filteredWords, word)
		}
	}
	result := strings.Join(filteredWords, "-")
	return result
}
