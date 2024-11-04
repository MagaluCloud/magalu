package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

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

type verify struct {
	path   string
	method string
	hidden bool
}

type rejected struct {
	verify
	spec string
}

const (
	DELETE = "DEL"
	GET    = "GET"
	PATCH  = "PATCH"
	POST   = "POST"
	PUT    = "PUT"
)

func processHiddenExtension(method, extKey, extValue, menu, path string, toVerify *[]verify) {
	// fmt.Println(menu, path, method, " {", extKey, ":", extValue, "}")
	hiddenValue, err := strconv.ParseBool(extValue)
	if err != nil {
		fmt.Println("Error parsing bool:", err)
		return
	}

	*toVerify = append(*toVerify, verify{
		path:   path,
		method: method,
		hidden: hiddenValue,
	})
}
func removeVersionFromURL(url string) (string, int, error) {
	re := regexp.MustCompile(`^/v(\d+)/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return url, 0, fmt.Errorf("no version found in URL")
	}
	version, err := strconv.Atoi(matches[1])
	if err != nil {
		return url, 0, fmt.Errorf("invalid version number: %v", err)
	}
	cleanURL := re.ReplaceAllString(url, "/")
	return cleanURL, version, nil
}

// prepareToGoCmd is a hidden command that prepares all available specs to golang
func runPrepare(cmd *cobra.Command, args []string) {
	_ = verificarEAtualizarDiretorio(currentDir())

	currentConfig, err := loadList()

	if err != nil {
		fmt.Println(err)
		return
	}

	// finalFile := filepath.Join(currentDir(), "specs.go.tmp")
	// newFileSpecs, err := os.OpenFile(finalFile, os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer newFileSpecs.Close()

	// _, _ = newFileSpecs.Write([]byte("package openapi\n\n"))
	// _, _ = newFileSpecs.Write([]byte("import (\n"))
	// _, _ = newFileSpecs.Write([]byte("	\"os\"\n"))
	// _, _ = newFileSpecs.Write([]byte("	\"syscall\"\n"))
	// _, _ = newFileSpecs.Write([]byte("	\"magalu.cloud/core/dataloader\"\n"))
	// _, _ = newFileSpecs.Write([]byte(")\n\n"))
	// _, _ = newFileSpecs.Write([]byte("type embedLoader map[string][]byte\n"))
	// _, _ = newFileSpecs.Write([]byte("func GetEmbedLoader() dataloader.Loader {\n"))
	// _, _ = newFileSpecs.Write([]byte("return embedLoaderInstance\n"))
	// _, _ = newFileSpecs.Write([]byte("		}\n"))
	// _, _ = newFileSpecs.Write([]byte("func (f embedLoader) Load(name string) ([]byte, error) {\n"))
	// _, _ = newFileSpecs.Write([]byte("if data, ok := embedLoaderInstance[name]; ok {\n"))
	// _, _ = newFileSpecs.Write([]byte("return data, nil\n"))
	// _, _ = newFileSpecs.Write([]byte("}\n"))
	// _, _ = newFileSpecs.Write([]byte("return nil, &os.PathError{Op: \"open\", Path: name, Err: syscall.ENOENT}\n"))
	// _, _ = newFileSpecs.Write([]byte("}\n"))
	// _, _ = newFileSpecs.Write([]byte("func (f embedLoader) String() string {\n"))
	// _, _ = newFileSpecs.Write([]byte("		return \"embedLoader\"\n"))
	// _, _ = newFileSpecs.Write([]byte("}\n"))
	// _, _ = newFileSpecs.Write([]byte("var embedLoaderInstance = embedLoader{\n"))

	// indexModules := []modules{}

	rejectedPaths := []rejected{}

	for _, v := range currentConfig {
		// fileStringBase64 := ""
		// fmt.Println(filepath.Join(currentDir(), v.File))
		//read file and convert to string and save in new generate a new go file
		if v.Enabled {
			file := filepath.Join(currentDir(), v.File)
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

			// indexModules = append(indexModules, modules{
			// 	Description: docModel.Model.Info.Description,
			// 	Name:        v.Menu,
			// 	Path:        v.File,
			// 	Summary:     docModel.Model.Info.Description,
			// 	URL:         v.Url,
			// 	Version:     docModel.Model.Info.Version,
			// 	CLI:         v.CLI,
			// 	TF:          v.TF,
			// 	SDK:         v.SDK,
			// })

			toVerify := []verify{}

			for pair := docModel.Model.Paths.PathItems.Oldest(); pair != nil; pair = pair.Next() {
				if pair.Value.Delete != nil {
					for ext := pair.Value.Delete.Extensions.Oldest(); ext != nil; ext = ext.Next() {
						if ext.Key == "x-mgc-hidden" {
							processHiddenExtension(DELETE, ext.Key, ext.Value.Value, v.Menu, pair.Key, &toVerify)
						}
					}
				}

				if pair.Value.Get != nil {
					for ext := pair.Value.Get.Extensions.Oldest(); ext != nil; ext = ext.Next() {
						if ext.Key == "x-mgc-hidden" {
							processHiddenExtension(GET, ext.Key, ext.Value.Value, v.Menu, pair.Key, &toVerify)
						}
					}
				}

				if pair.Value.Patch != nil {
					for ext := pair.Value.Patch.Extensions.Oldest(); ext != nil; ext = ext.Next() {
						if ext.Key == "x-mgc-hidden" {
							processHiddenExtension(PATCH, ext.Key, ext.Value.Value, v.Menu, pair.Key, &toVerify)
						}
					}
				}

				if pair.Value.Post != nil {
					for ext := pair.Value.Post.Extensions.Oldest(); ext != nil; ext = ext.Next() {
						if ext.Key == "x-mgc-hidden" {
							processHiddenExtension(POST, ext.Key, ext.Value.Value, v.Menu, pair.Key, &toVerify)
						}
					}
				}

				if pair.Value.Put != nil {
					for ext := pair.Value.Post.Extensions.Oldest(); ext != nil; ext = ext.Next() {
						if ext.Key == "x-mgc-hidden" {
							processHiddenExtension(PUT, ext.Key, ext.Value.Value, v.Menu, pair.Key, &toVerify)
						}
					}
				}
			}

			ccVerify := make([]verify, len(toVerify))
			rejectPaths := make([]verify, 0)

			copy(ccVerify, toVerify)
			for _, vv := range toVerify {
				suffix, vVersion, err := removeVersionFromURL(vv.path)

				if err != nil {
					fmt.Println(err)
					return
				}

				for _, c := range ccVerify {
					_, cVersion, _ := removeVersionFromURL(c.path)

					if c.method != vv.method {
						continue
					}

					if !strings.HasSuffix(c.path, suffix) {
						continue
					}

					if cVersion == vVersion {
						continue
					}

					if c.hidden != vv.hidden || (!c.hidden && !vv.hidden) {
						continue
					}

					rejectPaths = append(rejectPaths, vv)
				}
			}

			for _, xv := range rejectPaths {
				// fmt.Printf("Rejecting %s %s - Hidden: %t\n", xv.method, xv.path, xv.hidden)
				rejectedPaths = append(rejectedPaths, rejected{
					verify: verify{
						path:   xv.path,
						hidden: xv.hidden,
						method: xv.method,
					},
					spec: v.File,
				})

			}
			rejectPaths = nil
			toVerify = nil
			ccVerify = nil

			_, document, _, errs := document.RenderAndReload()
			if len(errors) > 0 {
				panic(fmt.Sprintf("cannot re-render document: %d errors reported", len(errs)))
			}

			_, errors = document.BuildV3Model()
			if len(errors) > 0 {
				for i := range errors {
					fmt.Printf("error: %e\n", errors[i])
				}
				panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported", len(errors)))
			}

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

			if len(rejectedPaths) == 0 {
				fileBytes, _, _, errs = document.RenderAndReload()
				if len(errors) > 0 {
					panic(fmt.Sprintf("cannot re-render document: %d errors reported", len(errs)))
				}

				err = os.WriteFile(filepath.Join(currentDir(), v.File), fileBytes, 0644)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		// _, _ = newFileSpecs.Write([]byte(fmt.Sprintf("\"%v\":([]byte)(\"%v\"),\n", v.File, fileStringBase64)))

	}

	if len(rejectedPaths) > 0 {
		fmt.Println("Rejected paths:")
		for _, v := range rejectedPaths {
			fmt.Printf("Spec: %s - %s - %s - Hidden: %t\n", v.spec, v.method, v.path, v.hidden)
		}
		os.Exit(1)
	}
	//convert to json

	// indexJson, err := json.Marshal(indexModules)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fileStringBase64 := b64.StdEncoding.EncodeToString(indexJson)
	// _, _ = newFileSpecs.Write([]byte(fmt.Sprintf("\"%v\":([]byte)(\"%v\"),\n", "index.openapi.json", fileStringBase64)))
	// _, _ = newFileSpecs.Write([]byte("\n}\n"))

}

// replace another python scripts
var prepareToGoCmd = &cobra.Command{
	Use:    "prepare",
	Short:  "Prepare all available specs to MgcSDK",
	Hidden: true,
	Run:    runPrepare,
}
