package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/spec_manipulator/cmd/tui"
	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"

	"github.com/spf13/cobra"
)

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

func processHiddenExtension(method, extValue, path string, toVerify *[]verify) {
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

func runPrepare(dir string, menu string) {

	_ = verificarEAtualizarDiretorio(dir)

	var configToRun []specList
	var err error

	if menu != "" {
		configToRun, err = loadList(menu)
	} else {
		configToRun, err = getConfigToRun()
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	spinner := tui.NewSpinner()
	spinner.Start("Preparing ...")
	for _, v := range configToRun {
		spinner.UpdateText("Preparing " + v.File)
		justRunPrepare(dir, v)
	}

	spinner.Success("Specs prepared successfully")
}

func justRunPrepare(dir string, v specList) {
	if v.Enabled {
		file := filepath.Join(dir, v.File)
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

		for _, remove := range v.ToRemove {
			if remove.PathPrefix != nil {
				for pair := docModel.Model.Paths.PathItems.Oldest(); pair != nil; pair = pair.Next() {
					if strings.HasPrefix(pair.Key, *remove.PathPrefix) {
						docModel.Model.Paths.PathItems.Delete(pair.Key)
					}
				}
			}

			if remove.Path != nil && remove.Method == nil {
				for pair := docModel.Model.Paths.PathItems.Oldest(); pair != nil; pair = pair.Next() {
					if strings.EqualFold(pair.Key, *remove.Path) {
						docModel.Model.Paths.PathItems.Delete(pair.Key)
					}
				}
			}

			if remove.Path != nil && remove.Method != nil {
				for pair := docModel.Model.Paths.PathItems.Oldest(); pair != nil; pair = pair.Next() {
					if strings.EqualFold(pair.Key, *remove.Path) {
						for _, method := range remove.Method {
							switch strings.ToUpper(method) {
							case "GET":
								pair.Value.Get = nil
								// docModel.Model.Paths.PathItems.Set(pair.Key, pair.Value)
							case "POST":
								pair.Value.Post = nil
								// docModel.Model.Paths.PathItems.Set(pair.Key, pair.Value)
							case "PUT":
								pair.Value.Put = nil
								// docModel.Model.Paths.PathItems.Set(pair.Key, pair.Value)
							case "DELETE":
								pair.Value.Delete = nil
								// docModel.Model.Paths.PathItems.Set(pair.Key, pair.Value)
							case "PATCH":
								pair.Value.Patch = nil
								// docModel.Model.Paths.PathItems.Set(pair.Key, pair.Value)
							case "OPTIONS":
								pair.Value.Options = nil
								// docModel.Model.Paths.PathItems.Set(pair.Key, pair.Value)
							}
						}
					}
				}
			}

		}

		_, errors = document.BuildV3Model()
		if len(errors) > 0 {
			for i := range errors {
				fmt.Printf("error: %e\n", errors[i])
			}
			panic(fmt.Sprintf("cannot create v3 model from document: %d errors reported", len(errors)))
		}
		_, document, _, errs := document.RenderAndReload()
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
				fmt.Printf("Type: %s, Failure: %s\n", e.ValidationType, e.Message)
				fmt.Printf("Fix: %s\n\n", e.HowToFix)
			}
		}

		fileBytes, _, _, errs = document.RenderAndReload()
		if len(errors) > 0 {
			panic(fmt.Sprintf("cannot re-render document: %d errors reported", len(errs)))
		}

		err = os.WriteFile(filepath.Join(dir, v.File), fileBytes, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}

// replace another python scripts
func prepareToGoCmd() *cobra.Command {
	var dir string
	var menu string
	cmd := &cobra.Command{
		Use:   "prepare",
		Short: "Prepare all available specs to MgcSDK",
		Run: func(cmd *cobra.Command, args []string) {
			runPrepare(dir, menu)
		},
	}
	cmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to save the converted specs")
	cmd.Flags().StringVarP(&menu, "menu", "m", "", "Menu to prepare")
	return cmd
}
