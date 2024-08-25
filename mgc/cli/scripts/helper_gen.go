package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Node struct {
	Name     string  `json:"name"`
	Children []*Node `json:"children,omitempty"`
}

func loadJSON(filename string) ([]*Node, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var nodes []*Node
	err = json.Unmarshal(data, &nodes)
	return nodes, err
}

func iterTree(children []*Node, parentPath []string) [][]string {
	var paths [][]string
	for _, child := range children {
		path := append(parentPath, child.Name)
		paths = append(paths, path)
		if child.Children != nil {
			paths = append(paths, iterTree(child.Children, path)...)
		}
	}
	return paths
}

func genCliPaths(cli string, tree []*Node) [][]string {
	paths := [][]string{{cli}}
	for _, p := range iterTree(tree, []string{cli}) {
		paths = append(paths, p)
	}
	return paths
}

func genOutput(cmd []string) (string, error) {
	cmd = append(cmd, "--raw", "-h")
	log.Printf("running %v", cmd)
	output, err := exec.Command(cmd[0], cmd[1:]...).Output()
	return string(output), err
}

func main() {
	cli := flag.String("cli", "", "the binary to use during executions")
	cliDumpTreeJSON := flag.String("cli-dump-tree-json", "", "path to the result of cli dump-tree")
	outputDirectory := flag.String("output-directory", "", "the root folder where to write help output")
	verbose := flag.Int("verbose", 0, "verbose level")
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("INFO: ")

	if *verbose > 0 {
		log.SetFlags(log.Lshortfile)
	}

	rootDir, err := filepath.Abs(*outputDirectory)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	log.Printf("removing output-dir: %s", rootDir)
	os.RemoveAll(rootDir)

	tree, err := loadJSON(*cliDumpTreeJSON)
	if err != nil {
		log.Fatalf("Error loading JSON: %v", err)
	}

	for _, path := range genCliPaths(*cli, tree) {
		log.Printf("processing: %v", path)
		helpOutput, err := genOutput(path)
		if err != nil {
			log.Printf("Error generating help output: %v", err)
			continue
		}

		outDir := filepath.Join(rootDir, filepath.Join(path[1:]...))
		err = os.MkdirAll(outDir, 0755)
		if err != nil {
			log.Printf("Error creating directory: %v", err)
			continue
		}

		filePath := filepath.Join(outDir, "help.md")
		err = ioutil.WriteFile(filePath, []byte(helpOutput), 0644)
		if err != nil {
			log.Printf("Error writing file: %v", err)
			continue
		}
		log.Printf("wrote %s", filePath)
	}
}
