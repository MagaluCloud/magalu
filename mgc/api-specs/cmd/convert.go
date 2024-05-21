package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/pb33f/libopenapi"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var convertCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare viveiro with jaxyendy",
	Run: func(cmd *cobra.Command, args []string) {

		_ = verificarEAtualizarDiretorio(SPEC_DIR)

		jaxyendy, err := loadListFromViper(toWriteViveiro(false))
		if err != nil {
			fmt.Println(err)
			return
		}

		doConvert(jaxyendy)
	},
}

func doConvert(jaxyendy []specList) {
	runPath, _ := os.Getwd()
	pathSpecs := "cli_specs"

	for _, j := range jaxyendy {

		jeFile, _ := os.ReadFile(path.Join(runPath, pathSpecs, j.File))
		document, err := libopenapi.NewDocument(jeFile)
		if err != nil {
			panic(fmt.Sprintf("cannot create new %s document: %e", j.Menu, err))
		}

		v3Model, errors := document.BuildV3Model()

		if len(errors) > 0 {
			for i := range errors {
				fmt.Printf("error: %e\n", errors[i])
			}
			panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported",
				len(errors)))
		}

		v3Model.Model.Version = "3.0.3"

		_, _, newModel, errs := document.RenderAndReload()
		if errs != nil {
			panic(err)
		}

		bytess, err := newModel.Model.Render()
		if err == nil {
			savedFileName := fmt.Sprintf("%s.converted.openapi.yaml", j.Menu)
			saveFile := path.Join(runPath, pathSpecs, savedFileName)
			err := os.WriteFile(saveFile, bytess, fs.FileMode(0644))
			if err != nil {
				panic("fail to save converted file")
			}
			From310to303(bytess, saveFile)

		}
		//change to goroutine

	}

}

func From310to303(fileBytes []byte, filePath string) {

	var dockerCompose yaml.Node
	yaml.Unmarshal(fileBytes, &dockerCompose)
	// Swap out redis:alpine for redis:golang
	examples := findChildNode("examples", &dockerCompose)
	if examples != nil {
		examples.SetString("redis:golang")
	}
	// Create a modified yaml file
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Problem creating file: %v", err)
	}
	defer f.Close()
	yaml.NewEncoder(f).Encode(dockerCompose.Content[0])
}

// Recusive function to find the child node by value that we care about.
// Probably needs tweaking so use with caution.
func findChildNode(value string, node *yaml.Node) *yaml.Node {
	for _, v := range node.Content {
		// If we found the value we are looking for, return it.
		if v.Value == value {
			return v
		}
		// Otherwise recursively look more
		if child := findChildNode(value, v); child != nil {
			return child
		}
	}
	return nil
}
