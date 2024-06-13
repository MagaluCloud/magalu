package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	b64 "encoding/base64"
	"encoding/json"

	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"
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

// prepareToGoCmd is a hidden command that prepares all available specs to golang
func runPrepare(cmd *cobra.Command, args []string) {
	_ = verificarEAtualizarDiretorio(SPEC_DIR)

	currentConfig, err := loadList()

	if err != nil {
		fmt.Println(err)
		return
	}

	finalFile := filepath.Join(SPEC_DIR, "specs.go.tmp")
	newFileSpecs, err := os.OpenFile(finalFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer newFileSpecs.Close()

	_, _ = newFileSpecs.Write([]byte("package openapi\n\n"))
	_, _ = newFileSpecs.Write([]byte("import (\n"))
	_, _ = newFileSpecs.Write([]byte("	\"os\"\n"))
	_, _ = newFileSpecs.Write([]byte("	\"syscall\"\n"))
	_, _ = newFileSpecs.Write([]byte("	\"magalu.cloud/core/dataloader\"\n"))
	_, _ = newFileSpecs.Write([]byte(")\n\n"))
	_, _ = newFileSpecs.Write([]byte("type embedLoader map[string][]byte\n"))
	_, _ = newFileSpecs.Write([]byte("func GetEmbedLoader() dataloader.Loader {\n"))
	_, _ = newFileSpecs.Write([]byte("return embedLoaderInstance\n"))
	_, _ = newFileSpecs.Write([]byte("		}\n"))
	_, _ = newFileSpecs.Write([]byte("func (f embedLoader) Load(name string) ([]byte, error) {\n"))
	_, _ = newFileSpecs.Write([]byte("if data, ok := embedLoaderInstance[name]; ok {\n"))
	_, _ = newFileSpecs.Write([]byte("return data, nil\n"))
	_, _ = newFileSpecs.Write([]byte("}\n"))
	_, _ = newFileSpecs.Write([]byte("return nil, &os.PathError{Op: \"open\", Path: name, Err: syscall.ENOENT}\n"))
	_, _ = newFileSpecs.Write([]byte("}\n"))
	_, _ = newFileSpecs.Write([]byte("func (f embedLoader) String() string {\n"))
	_, _ = newFileSpecs.Write([]byte("		return \"embedLoader\"\n"))
	_, _ = newFileSpecs.Write([]byte("}\n"))
	_, _ = newFileSpecs.Write([]byte("var embedLoaderInstance = embedLoader{\n"))

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

			_, document, _, errs := document.RenderAndReload()
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

			//todo - remove from py
			// svar := orderedmap.New[string, *v3.ServerVariable]()
			// svar.Set("region", &v3.ServerVariable{
			// 	Default:     "br-se1",
			// 	Description: "Region to reach the service",
			// 	Enum: []string{
			// 		"br-ne-1",
			// 		"br-se1",
			// 		"br-mgl1",
			// 	},
			// })

			// svar.Set("env", &v3.ServerVariable{
			// 	Description: "Environment to use",
			// 	Default:     "api.magalu.cloud",
			// 	Enum: []string{
			// 		"api.magalu.cloud",
			// 		"api.pre-prod.jaxyendy.com",
			// 	},
			// })

			// servers := []*v3.Server{}
			// servers = append(servers, &v3.Server{
			// 	URL:         "https://{env}/{region}/$API_ENDPOINT_NAME",
			// 	Description: "",
			// 	Variables:   svar,
			// })

			// docModel.Model.Servers = servers

			_, document, _, errs = document.RenderAndReload()
			if len(errors) > 0 {
				panic(fmt.Sprintf("cannot re-render document: %d errors reported", len(errs)))
			}
			docValidator, validatorErrs := validator.NewValidator(document)
			if len(validatorErrs) > 0 {
				panic(fmt.Sprintf("cannot create validator: %d errors reported", len(validatorErrs)))
			}

			valid, validationErrs := docValidator.ValidateDocument()

			if !valid {
				for _, e := range validationErrs {
					// 5. Handle the error
					fmt.Printf("Type: %s, Failure: %s\n", e.ValidationType, e.Message)
					fmt.Printf("Fix: %s\n\n", e.HowToFix)
				}
			}

			fileBytes, _, _, errs = document.RenderAndReload()
			if len(errors) > 0 {
				panic(fmt.Sprintf("cannot re-render document: %d errors reported", len(errs)))
			}

			fileStringBase64 = b64.StdEncoding.EncodeToString(fileBytes)

			err = os.WriteFile(filepath.Join(SPEC_DIR, v.File), fileBytes, 0644)
			if err != nil {
				fmt.Println(err)
				return
			}

		}

		_, _ = newFileSpecs.Write([]byte(fmt.Sprintf("\"%v\":([]byte)(\"%v\"),\n", v.File, fileStringBase64)))

	}

	//convert to json

	indexJson, err := json.Marshal(indexModules)
	if err != nil {
		fmt.Println(err)
		return
	}

	fileStringBase64 := b64.StdEncoding.EncodeToString(indexJson)
	_, _ = newFileSpecs.Write([]byte(fmt.Sprintf("\"%v\":([]byte)(\"%v\"),\n", "index.openapi.json", fileStringBase64)))
	_, _ = newFileSpecs.Write([]byte("\n}\n"))

}

// replace another python scripts
var prepareToGoCmd = &cobra.Command{
	Use:    "prepare",
	Short:  "Prepare all available specs to golang",
	Hidden: true,
	Run:    runPrepare,
}
