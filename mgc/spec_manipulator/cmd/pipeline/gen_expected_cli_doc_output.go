package pipeline

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
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

func iterTree(children []TreeNode, parentPath []string) [][]string {
	var result [][]string
	for _, child := range children {
		path := append(parentPath, child.Name)
		result = append(result, path)

		if len(child.Children) > 0 {
			result = append(result, iterTree(child.Children, path)...)
		}
	}
	return result
}

func genCliPaths(cli string, tree []TreeNode) [][]string {
	result := [][]string{{cli}}
	result = append(result, iterTree(tree, []string{cli})...)
	return result
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

func convertToMarkdown(inputText string) string {
	sections := strings.Split(inputText, "\n\n")
	var markdown strings.Builder

	// Header section
	markdown.WriteString(fmt.Sprintf("# %s\n\n", strings.TrimSpace(sections[0])))

	// Usage section
	markdown.WriteString("## Usage:\n```bash\n")
	markdown.WriteString(strings.TrimSpace(sections[1]))
	markdown.WriteString("\n```\n\n")

	// Product catalog section
	markdown.WriteString("## Product catalog:\n")
	for _, item := range strings.Split(strings.TrimSpace(sections[2]), "\n") {
		markdown.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
	}
	markdown.WriteString("\n")

	// Other commands section
	markdown.WriteString("## Other commands:\n")
	for _, item := range strings.Split(strings.TrimSpace(sections[3]), "\n") {
		markdown.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
	}
	markdown.WriteString("\n")

	// Flags section
	markdown.WriteString("## Flags:\n```bash\n")
	markdown.WriteString(strings.TrimSpace(sections[4]))
	markdown.WriteString("\n```\n\n")

	return markdown.String()
}

type CliDocParams struct {
	cli         string
	dumpCliJson string
	outputDir   string
	verbose     int
}

func runDocParams(params CliDocParams) {

	log.SetFlags(0)
	log.SetPrefix("INF ")

	if params.verbose > 0 {
		log.SetPrefix("DBG ")
	}

	tree, err := loadJSON(params.dumpCliJson)
	if err != nil {
		log.Fatalf("Failed to load JSON: %v", err)
	}

	rootDir, _ := filepath.Abs(params.outputDir)
	log.Printf("removing output-dir: %s", rootDir)
	os.RemoveAll(rootDir)

	for _, path := range genCliPaths(params.cli, tree) {
		log.Printf("processing: %s", strings.Join(path, " "))
		helpOutput, err := genHelpOutput(path)
		if err != nil {
			log.Printf("Error generating help output: %v", err)
			continue
		}
		markdownOutput := convertToMarkdown(helpOutput)

		outDir := filepath.Join(rootDir, filepath.Join(path[1:]...))
		_ = os.MkdirAll(outDir, os.ModePerm)
		filePath := filepath.Join(outDir, "help.md")
		err = os.WriteFile(filePath, []byte(markdownOutput), 0644)
		if err != nil {
			log.Printf("Error writing file: %v", err)
		} else {
			log.Printf("wrote %s", filePath)
		}
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

	return cmd
}
