package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	b64 "encoding/base64"
	"encoding/json"

	"github.com/pb33f/libopenapi"
	"github.com/spf13/cobra"
)

type modules struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Summary     string `json:"summary"`
	URL         string `json:"url"`
	Version     string `json:"version"`
	CLI         bool   `json:"cli"`
	TF          bool   `json:"tf"`
	SDK         bool   `json:"sdk"`
}

// WIP WIP WIP
// replace another python scripts
var prepareToGoCmd = &cobra.Command{
	Use:    "prepare",
	Short:  "Prepare all available specs to golang",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = verificarEAtualizarDiretorio(SPEC_DIR)

		currentConfig, err := loadList()

		if err != nil {
			fmt.Println(err)
			return
		}

		finalFile := filepath.Join(SPEC_DIR, "specs.go")
		newFileSpecs, err := os.OpenFile(finalFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer newFileSpecs.Close()

		newFileSpecs.Write([]byte("package openapi\n\n"))
		newFileSpecs.Write([]byte("import (\n"))
		newFileSpecs.Write([]byte("	\"os\"\n"))
		newFileSpecs.Write([]byte("	\"syscall\"\n"))
		newFileSpecs.Write([]byte("	\"magalu.cloud/core/dataloader\"\n"))
		newFileSpecs.Write([]byte(")\n\n"))
		newFileSpecs.Write([]byte("type embedLoader map[string][]byte\n"))
		newFileSpecs.Write([]byte("func GetEmbedLoader() dataloader.Loader {\n"))
		newFileSpecs.Write([]byte("return embedLoaderInstance\n"))
		newFileSpecs.Write([]byte("		}\n"))
		newFileSpecs.Write([]byte("func (f embedLoader) Load(name string) ([]byte, error) {\n"))
		newFileSpecs.Write([]byte("if data, ok := embedLoaderInstance[name]; ok {\n"))
		newFileSpecs.Write([]byte("return data, nil\n"))
		newFileSpecs.Write([]byte("}\n"))
		newFileSpecs.Write([]byte("return nil, &os.PathError{Op: \"open\", Path: name, Err: syscall.ENOENT}\n"))
		newFileSpecs.Write([]byte("}\n"))
		newFileSpecs.Write([]byte("func (f embedLoader) String() string {\n"))
		newFileSpecs.Write([]byte("		return \"embedLoader\"\n"))
		newFileSpecs.Write([]byte("}\n"))
		newFileSpecs.Write([]byte("var embedLoaderInstance = embedLoader{\n"))

		indexModules := []modules{}

		for _, v := range currentConfig {
			fileStringBase64 := ""
			fmt.Println(filepath.Join(SPEC_DIR, v.File))
			//read file and convert to string and save in new generate a new go file
			if !v.Enabled {
				fileStringBase64 = ""
			} else {
				file := filepath.Join(SPEC_DIR, v.File)
				fileBytes, err := os.ReadFile(file)
				if err != nil {
					fmt.Println(err)
					return
				}

				document, err := libopenapi.NewDocument(fileBytes)
				if err != nil {
					panic(fmt.Sprintf("cannot read document: %e", err))
				}
				docModel, errors := document.BuildV3Model()
				if len(errors) > 0 {
					for i := range errors {
						fmt.Printf("error: %e\n", errors[i])
					}
					panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported", len(errors)))
				}

				indexModules = append(indexModules, modules{
					Description: docModel.Model.Info.Description,
					Name:        v.Menu,
					Path:        v.File,
					Summary:     docModel.Model.Info.Description,
					URL:         v.Url,
					Version:     docModel.Model.Info.Version,
					CLI:         v.CLI,
					TF:          v.TF,
					SDK:         v.SDK,
				})

				//remove all paths that contains xaas
				toRemove := []string{}
				for pair := docModel.Model.Paths.PathItems.Oldest(); pair != nil; pair = pair.Next() {
					if strings.Contains(strings.ToLower(pair.Key), "xaas") {
						toRemove = append(toRemove, pair.Key)
					}
				}

				for _, key := range toRemove {
					docModel.Model.Paths.PathItems.Delete(key)
				}

				fmt.Printf("Total PATH removed: %v\n", len(toRemove))

				toRemove = []string{}

				fileBytes, document, _, errs := document.RenderAndReload()
				if len(errors) > 0 {
					panic(fmt.Sprintf("cannot re-render document: %d errors reported", len(errs)))
				}

				docModel, errors = document.BuildV3Model()
				if len(errors) > 0 {
					for i := range errors {
						fmt.Printf("error: %e\n", errors[i])
					}
					panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported", len(errors)))
				}

				for pair := docModel.Model.Components.Schemas.Oldest(); pair != nil; pair = pair.Next() {
					if strings.Contains(strings.ToLower(pair.Key), "xaas") {
						toRemove = append(toRemove, pair.Key)
					}
				}

				for _, key := range toRemove {
					docModel.Model.Components.Schemas.Delete(key)
				}

				fmt.Printf("Total COMPONENT removed: %v\n", len(toRemove))

				fileBytes, _, _, errs = document.RenderAndReload()
				if len(errors) > 0 {
					panic(fmt.Sprintf("cannot re-render document: %d errors reported", len(errs)))
				}

				fileStringBase64 = b64.StdEncoding.EncodeToString(fileBytes)

				os.WriteFile(filepath.Join(SPEC_DIR, v.File), fileBytes, 0644)

			}

			newFileSpecs.Write([]byte(fmt.Sprintf("\"%v\":([]byte)(\"%v\"),\n", v.File, fileStringBase64)))

		}

		//convert to json

		indexJson, err := json.Marshal(indexModules)
		if err != nil {
			fmt.Println(err)
			return
		}

		fileStringBase64 := b64.StdEncoding.EncodeToString(indexJson)
		newFileSpecs.Write([]byte(fmt.Sprintf("\"%v\":([]byte)(\"%v\"),\n", "index.openapi.json", fileStringBase64)))
		newFileSpecs.Write([]byte("\n}\n"))

	},
}
