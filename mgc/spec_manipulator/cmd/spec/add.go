package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

func add(options AddMenu) {

	file := fmt.Sprintf("%s.jaxyendy.openapi.json", options.menu)

	currentConfig, err := loadList("")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range currentConfig {
		if v.Url == options.url {
			fmt.Println("url already exists with menu: " + v.Menu)
			return
		}
	}
	currentConfig = append(currentConfig, specList{Url: options.url, File: file, Menu: options.menu, Enabled: true})
	saveConfig(currentConfig)
	fmt.Println("Added successfully")
}

type AddMenu struct {
	url  string
	menu string
}

func specAddNewCmd() *cobra.Command {
	options := &AddMenu{}

	cmd := &cobra.Command{
		Use:     "add [url] [menu]",
		Short:   "Add new spec",
		Example: "specs add https://block-storage.br-ne-1.jaxyendy.com/v1/openapi.json block-storage",
		Run: func(cmd *cobra.Command, args []string) {
			if options.menu == "" {
				fmt.Println(cmd.UsageString())
				fmt.Println(">> menu is required")
				return
			}
			if !validarEndpoint(options.url) {
				fmt.Println(cmd.UsageString())

				fmt.Print(">> url is invalid\n\n")
				fmt.Print("Gitlab example: \n    main/api_products/mcr-api/br-ne1-prod-yel-1/openapi.yaml\n\n")
				fmt.Print("Gitlab example: \n    https://gitlab.luizalabs.com/open-platform/pcx/u0/-/blob/main/api_products/mcr-api/br-ne1-prod-yel-1/openapi.yaml?plain=1\n\n")
				fmt.Print("Deployed spec example: \n    https://block-storage.br-ne-1.jaxyendy.com/v1/openapi.json\n\n")
				return
			}

			add(*options)
		},
	}

	cmd.Flags().StringVarP(&options.url, "url", "u", "", "URL")
	cmd.Flags().StringVarP(&options.menu, "menu", "m", "", "Menu")

	return cmd
}
