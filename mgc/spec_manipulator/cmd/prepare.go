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
	"github.com/pb33f/libopenapi/orderedmap"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

func runPrepare(cmd *cobra.Command, args []string) {
	_ = verificarEAtualizarDiretorio(currentDir())

	currentConfig, err := loadList()

	if err != nil {
		fmt.Println(err)
		return
	}

	rejectedPaths := []rejected{}

	for _, v := range currentConfig {
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

			toVerify := []verify{}
			toRemove := []string{}
			for pair := docModel.Model.Paths.PathItems.Oldest(); pair != nil; pair = pair.Next() {
				forceHidden := false
				if strings.Contains(strings.ToLower(pair.Key), "xaas") || strings.Contains(strings.ToLower(pair.Key), "/internal") {
					forceHidden = true
					/*
						BEGIN
						Esse código é apenas para manter a compatibilidade com o que tinhamos antes.
						Após garantir funcionamento, remove-lo e garantir que quando for xaas, o x-mgc-hidden fique true.
					*/
					toRemove = append(toRemove, pair.Key)
					continue
					// END
				}

				if pair.Value.Delete != nil {
					hasHidden := false
					for ext := pair.Value.Delete.Extensions.Oldest(); ext != nil; ext = ext.Next() {
						if ext.Key == "x-mgc-hidden" {
							processHiddenExtension(DELETE, ext.Value.Value, pair.Key, &toVerify)
							if forceHidden && ext.Value.Value != "true" {
								ext.Value.Value = "true"
							}
							hasHidden = true
						}
					}
					if !hasHidden && forceHidden {
						if pair.Value.Delete.Extensions == nil {
							pair.Value.Delete.Extensions = &orderedmap.Map[string, *yaml.Node]{}
						}
						pair.Value.Delete.Extensions.Set("x-mgc-hidden", &yaml.Node{
							Kind:  yaml.ScalarNode,
							Value: "true",
						})
					}
				}

				if pair.Value.Get != nil {
					hasHidden := false
					for ext := pair.Value.Get.Extensions.Oldest(); ext != nil; ext = ext.Next() {
						if ext.Key == "x-mgc-hidden" {
							processHiddenExtension(GET, ext.Value.Value, pair.Key, &toVerify)
							if forceHidden && ext.Value.Value != "true" {
								ext.Value.Value = "true"
							}
							hasHidden = true
						}
					}
					if !hasHidden && forceHidden {
						if pair.Value.Get.Extensions == nil {
							pair.Value.Get.Extensions = &orderedmap.Map[string, *yaml.Node]{}
						}
						pair.Value.Get.Extensions.Set("x-mgc-hidden", &yaml.Node{
							Kind:  yaml.ScalarNode,
							Value: "true",
						})
					}

				}

				if pair.Value.Patch != nil {
					hasHidden := false
					for ext := pair.Value.Patch.Extensions.Oldest(); ext != nil; ext = ext.Next() {
						if ext.Key == "x-mgc-hidden" {
							processHiddenExtension(PATCH, ext.Value.Value, pair.Key, &toVerify)
							if forceHidden && ext.Value.Value != "true" {
								ext.Value.Value = "true"
							}
							hasHidden = true
						}
					}
					if !hasHidden && forceHidden {
						if pair.Value.Patch.Extensions == nil {
							pair.Value.Patch.Extensions = &orderedmap.Map[string, *yaml.Node]{}
						}
						pair.Value.Patch.Extensions.Set("x-mgc-hidden", &yaml.Node{
							Kind:  yaml.ScalarNode,
							Value: "true",
						})
					}

				}

				if pair.Value.Post != nil {
					hasHidden := false
					for ext := pair.Value.Post.Extensions.Oldest(); ext != nil; ext = ext.Next() {
						if ext.Key == "x-mgc-hidden" {
							processHiddenExtension(POST, ext.Value.Value, pair.Key, &toVerify)
							if forceHidden && ext.Value.Value != "true" {
								ext.Value.Value = "true"
							}
							hasHidden = true
						}
					}
					if !hasHidden && forceHidden {
						if pair.Value.Post.Extensions == nil {
							pair.Value.Post.Extensions = &orderedmap.Map[string, *yaml.Node]{}
						}
						pair.Value.Post.Extensions.Set("x-mgc-hidden", &yaml.Node{
							Kind:  yaml.ScalarNode,
							Value: "true",
						})
					}

				}

				if pair.Value.Put != nil {
					hasHidden := false
					for ext := pair.Value.Put.Extensions.Oldest(); ext != nil; ext = ext.Next() {
						if ext.Key == "x-mgc-hidden" {
							processHiddenExtension(PUT, ext.Value.Value, pair.Key, &toVerify)
							if forceHidden && ext.Value.Value != "true" {
								ext.Value.Value = "true"
							}
							hasHidden = true
						}
					}
					if !hasHidden && forceHidden {
						if pair.Value.Put.Extensions == nil {
							pair.Value.Put.Extensions = &orderedmap.Map[string, *yaml.Node]{}
						}
						pair.Value.Put.Extensions.Set("x-mgc-hidden", &yaml.Node{
							Kind:  yaml.ScalarNode,
							Value: "true",
						})
					}

				}
			}

			/*
				BEGIN
				Aqui continua o código a ser removido
			*/
			for _, key := range toRemove {
				docModel.Model.Paths.PathItems.Delete(key)
			}
			//END

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
					cffix, cVersion, _ := removeVersionFromURL(c.path)

					if c.method != vv.method {
						continue
					}

					if c.path == vv.path {
						continue
					}

					if !strings.HasSuffix(c.path, suffix) {
						continue
					}

					if suffix != cffix {
						continue
					}

					if cVersion == vVersion {
						continue
					}

					if vVersion < cVersion {
						continue
					}

					if (!vv.hidden && c.hidden) || (c.hidden && vv.hidden) {
						continue
					}

					rejectPaths = append(rejectPaths, vv)
				}
			}

			for _, xv := range rejectPaths {
				rejectedPaths = append(rejectedPaths, rejected{
					verify: verify{
						path:   xv.path,
						hidden: xv.hidden,
						method: xv.method,
					},
					spec: v.File,
				})

			}
			toVerify = nil

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

	}

	if len(rejectedPaths) > 0 {
		fmt.Println("Rejected paths:")
		for _, v := range rejectedPaths {
			fmt.Printf("Spec: %s - %s - %s - Hidden: %t\n", v.spec, v.method, v.path, v.hidden)
		}
		os.Exit(1)
	}
}

// replace another python scripts
var prepareToGoCmd = &cobra.Command{
	Use:    "prepare",
	Short:  "Prepare all available specs to MgcSDK",
	Hidden: true,
	Run:    runPrepare,
}
