package pipeline

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type TreeNode struct {
	Name     string     `json:"name"`
	Children []TreeNode `json:"children,omitempty"`
}

func loadJSON(filename string) ([]TreeNode, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var tree []TreeNode
	err = json.Unmarshal(data, &tree)
	return tree, err
}

func printPaths(node TreeNode, path string, result *[]string) {
	currentPath := node.Name
	if path != "" {
		currentPath = path + " " + node.Name
	}

	*result = append(*result, currentPath)

	for _, child := range node.Children {
		printPaths(child, currentPath, result)
	}
}

func genCliPaths(tree []TreeNode) []string {
	results := &[]string{}
	// Se o caminho está vazio, começamos apenas com o nome do nó atual
	for _, node := range tree {
		printPaths(node, "", results)
	}

	return *results
}

func genOutput(cmd []string) (string, error) {
	cmd = append(cmd, "--raw")
	output, err := exec.Command(cmd[0], cmd[1:]...).Output()
	return string(output), err
}

func genHelpOutput(path []string) (string, error) {
	cmd := append([]string{path[0], "help"}, path[1:]...)
	return genOutput(cmd)
}

// func genOutputHFlag(path []string) (string, error) {
// 	cmd := append(path, "-h")
// 	return genOutput(cmd)
// }

func prepareOutputString(inputText string, replaceString string, prepareLN bool) string {
	inputText = strings.ReplaceAll(inputText, replaceString, "")

	inputText = strings.ReplaceAll(inputText, "</br>", "")

	if !prepareLN {
		inputText += "\n\n"
		return inputText
	}

	lines := strings.Split(inputText, "\n")
	for i, x := range lines {
		if strings.HasPrefix(x, "  ") {
			lines[i] = strings.TrimPrefix(x, "  ")
		}
	}
	return strings.Join(lines, "\n")
}

func convertToMarkdown(inputText string, fileHeader string) string {
	sections := strings.Split(inputText, "\n\n")
	var markdown strings.Builder

	if fileHeader != "" {
		caser := cases.Title(language.Portuguese)
		fileHeader := caser.String(fileHeader)
		markdown.WriteString(fmt.Sprintf("# %s\n\n", fileHeader))
	}

	headerCtt := sections[0]
	// Header section
	if strings.Contains(sections[0], "██") {
		headerCtt = sections[1]
	}

	markdown.WriteString(prepareOutputString(headerCtt, "", false))

	// Usage section
	for _, section := range sections {
		if strings.Contains(section, "Usage:") {
			markdown.WriteString("## Usage:\n```\n")
			section = strings.ReplaceAll(section, "Usage:\n", "")
			markdown.WriteString(prepareOutputString(section, "Usage:", true))
			markdown.WriteString("\n```\n\n")
			break
		}
	}

	// Examples section
	for _, section := range sections {
		if strings.HasPrefix(section, "Examples:") {
			markdown.WriteString("## Examples:\n```\n")
			section = strings.ReplaceAll(section, "Examples:\n", "")
			markdown.WriteString(prepareOutputString(section, "Examples:", true))
			markdown.WriteString("\n```\n\n")
			break
		}
	}

	// Commands
	for _, section := range sections {
		if strings.Contains(section, "Commands:") {
			markdown.WriteString("## Commands:\n```\n")
			section = strings.ReplaceAll(section, "Commands:\n", "")
			markdown.WriteString(prepareOutputString(section, "Commands:", true))
			markdown.WriteString("\n```\n\n")
			break
		}
	}

	// Product catalog section
	for _, section := range sections {
		if strings.Contains(section, "Products:") {
			markdown.WriteString("## Product catalog:\n```\n")
			section = strings.ReplaceAll(section, "Products:\n", "")
			markdown.WriteString(prepareOutputString(section, "Products:", true))
			markdown.WriteString("\n```\n\n")
			break
		}
	}

	// Other commands section
	for _, section := range sections {
		if strings.Contains(section, "Other commands:") {
			markdown.WriteString("## Other commands:\n```\n")
			section = strings.ReplaceAll(section, "Other commands:\n", "")
			markdown.WriteString(prepareOutputString(section, "Other commands:", true))
			markdown.WriteString("\n```\n\n")
			break
		}
	}

	// Flags section
	for _, section := range sections {
		if strings.Contains(section, "Flags:") {
			markdown.WriteString("## Flags:\n```\n")
			section = strings.ReplaceAll(section, "Flags:\n", "")
			markdown.WriteString(prepareOutputString(section, "Flags:", true))
			markdown.WriteString("\n```\n\n")
			break
		}
	}

	// Global flags section
	for _, section := range sections {
		if strings.Contains(section, "Global Flags:") {
			markdown.WriteString("## Global Flags:\n```\n")
			section = strings.ReplaceAll(section, "Global Flags:\n", "")
			markdown.WriteString(prepareOutputString(section, "Global Flags:", true))
			markdown.WriteString("\n```\n\n")
			break
		}
	}

	//Settings
	for _, section := range sections {
		if strings.Contains(section, "Settings:") {
			markdown.WriteString("## Settings:\n```\n")
			section = strings.ReplaceAll(section, "Settings:\n", "")
			markdown.WriteString(prepareOutputString(section, "Settings:", true))
			markdown.WriteString("\n```\n\n")
			break
		}
	}

	return markdown.String()
}

type CliDocParams struct {
	cli         string
	dumpCliJson string
	outputDir   string
	verbose     int
	goroutine   bool
}

func runDocParams(params CliDocParams) {
	if params.verbose == 0 {
		log.SetOutput(io.Discard)
	} else if params.verbose == 1 {
		log.SetFlags(0)
		log.SetPrefix("INF ")
	}
	tree, err := loadJSON(params.dumpCliJson)
	if err != nil {
		log.Fatalf("Failed to load JSON: %v", err)
	}

	rootDir, _ := filepath.Abs(params.outputDir)
	log.Printf("removing output-dir: %s", rootDir)
	os.RemoveAll(rootDir)

	insideRunDocParams(rootDir, []string{params.cli})

	if !params.goroutine {
		for _, path := range genCliPaths(tree) {
			path = fmt.Sprintf("%s %s", params.cli, path)
			insideRunDocParams(rootDir, strings.Split(path, ""))
		}
		return
	}

	if params.goroutine {
		wg := &sync.WaitGroup{}
		count := 0
		for _, path := range genCliPaths(tree) {
			count++
			wg.Add(1)
			go func(count int, rootDir string, path string) {
				defer wg.Done()
				count--
				path = fmt.Sprintf("%s %s", params.cli, path)
				insideRunDocParams(rootDir, strings.Split(path, ""))
			}(count, rootDir, path)
			if count >= 50 {
				wg.Wait()
			}
		}
		wg.Wait()
	}
}

func insideRunDocParams(rootDir string, path []string) {
	pathBingo := strings.Join(path, "")
	log.Printf("processing: %s", pathBingo)
	path = strings.Split(pathBingo, " ")
	helpOutput, err := genHelpOutput(path)
	if err != nil {
		log.Printf("Error generating help output: %v\npath: %v", err, path)
	}
	outDir := filepath.Join(rootDir, filepath.Join(path[1:]...))

	prepareHeader := strings.Split(pathBingo, "/")
	fileHeader := prepareHeader[len(prepareHeader)-1]
	prepareHeader = strings.Split(fileHeader, " ")
	fileHeader = prepareHeader[len(prepareHeader)-1]
	if prepareHeader[len(prepareHeader)-1] == "mgc" {
		fileHeader = ""
	}

	markdownOutput := convertToMarkdown(helpOutput, fileHeader)

	_ = os.MkdirAll(outDir, os.ModePerm)
	filePath := filepath.Join(outDir, "help.md")
	err = os.WriteFile(filePath, []byte(markdownOutput), 0644)
	if err != nil {
		log.Printf("Error writing file: %v", err)
	} else {
		log.Printf("wrote %s", filePath)
	}
}

func CliDocOutputCmd() *cobra.Command {
	options := &CliDocParams{}

	cmd := &cobra.Command{
		Use:   "cligendoc",
		Short: "run gen doc cli",
		Run: func(cmd *cobra.Command, args []string) {
			runDocParams(*options)
		},
	}

	cmd.Flags().StringVarP(&options.cli, "cli", "c", "", "Local ou comando da CLI")
	cmd.Flags().StringVarP(&options.outputDir, "outputdir", "o", "", "Local de saida do dump file")
	cmd.Flags().StringVarP(&options.dumpCliJson, "dump", "d", "", "CLI Dump file json")
	cmd.Flags().IntVarP(&options.verbose, "verbose", "v", 0, "Verbose")
	cmd.Flags().BoolVarP(&options.goroutine, "goroutine", "g", false, "Goroutine")

	return cmd
}
