package pipeline

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func genCliDumpTree(cli string) ([]interface{}, error) {
	args := []string{"dump-tree", "-o", "json", "--raw"}
	cmd := exec.Command(cli, args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing command: \n command: %s \n error: %w", cmd.String(), err)
	}

	var tree []interface{}
	err = json.Unmarshal(output, &tree)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return tree, nil
}

type DumpeMenu struct {
	cli    string
	output string
}

func dumpTree(options DumpeMenu) {

	if options.cli == "" {
		fmt.Println("Error: cli argument is required")
		flag.Usage()
		os.Exit(1)
	}

	tree, err := genCliDumpTree(options.cli)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var out *os.File
	if options.output == "" {
		out = os.Stdout
	} else {
		var err error
		out, err = os.Create(options.output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer out.Close()
	}

	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(tree)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}

}

func CliDumpTreeCmd() *cobra.Command {
	options := &DumpeMenu{}

	cmd := &cobra.Command{
		Use:   "dumptree",
		Short: "run dump tree",
		Run: func(cmd *cobra.Command, args []string) {
			dumpTree(*options)
		},
	}

	cmd.Flags().StringVarP(&options.cli, "cli", "c", "", "Local ou comando da CLI")
	cmd.Flags().StringVarP(&options.output, "output", "o", "", "Local de saida do dump file")

	return cmd
}
