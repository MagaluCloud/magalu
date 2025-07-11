package spec

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pb33f/libopenapi"

	validator "github.com/pb33f/libopenapi-validator"

	"github.com/spf13/cobra"
)

func justRunValidate(dir string, v specList) {
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

	fileBytes, _, _, errs := document.RenderAndReload()
	if len(errs) > 0 {
		panic(fmt.Sprintf("cannot re-render document: %d errors reported", len(errs)))
	}

	_ = os.WriteFile(filepath.Join(dir, v.File), fileBytes, 0644)
}

func validateSpec() *cobra.Command {
	var dir string
	var menu string
	cmd := &cobra.Command{
		Use:    "validate",
		Short:  "Validate and prettify specs",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			_ = verificarEAtualizarDiretorio(dir)

			var currentConfig []specList
			var err error

			if menu != "" {
				currentConfig, err = loadList(menu)
			} else {
				currentConfig, err = getConfigToRun()
			}
			if err != nil {
				return
			}

			for _, v := range currentConfig {
				justRunValidate(dir, v)
			}
		},
	}

	cmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to save the converted specs")
	cmd.Flags().StringVarP(&menu, "menu", "m", "", "Menu to validate")
	return cmd
}
