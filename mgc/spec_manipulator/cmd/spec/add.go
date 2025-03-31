package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type specList struct {
	Url     string `json:"url"`
	File    string `json:"file"`
	Menu    string `json:"menu"`
	Enabled bool   `json:"enabled"`
}

func interfaceToMap(i interface{}) (map[string]interface{}, bool) {
	mapa, ok := i.(map[string]interface{})
	if !ok {
		fmt.Println("A interface não é um mapa ou mapa de interfaces.")
		return nil, false
	}
	return mapa, true
}

func add(options AddMenu) {

	var toSave []specList
	file := fmt.Sprintf("%s.jaxyendy.openapi.json", options.menu)

	toSave = append(toSave, specList{Url: options.url, File: file, Menu: options.menu, Enabled: true})
	currentConfig, err := loadList()
	if err != nil {
		fmt.Println(err)
		return
	}
	if slices.Contains(currentConfig, toSave[0]) {
		fmt.Println("url already exists")
		return
	}
	if !validarEndpoint(options.url) {
		fmt.Println("url is invalid")
		return
	}

	toSave = append(toSave, currentConfig...)
	// move it to common function
	ex, err := os.Executable()
	home := filepath.Dir(ex)
	cobra.CheckErr(err)

	// Search config in home directory with name ".cobra" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(VIPER_FILE)

	viper.Set("jaxyendy", toSave)
	err = viper.WriteConfigAs(SPEC_DIR)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("done")
}

type AddMenu struct {
	url  string
	menu string
}

func SpecAddNewCmd() *cobra.Command {
	options := &AddMenu{}

	cmd := &cobra.Command{
		Use:     "add [url] [menu]",
		Short:   "Add new spec",
		Example: "specs add https://block-storage.br-ne-1.jaxyendy.com/v1/openapi.json block-storage",
		Run: func(cmd *cobra.Command, args []string) {
			add(*options)
		},
	}

	cmd.Flags().StringVarP(&options.url, "url", "u", "", "URL")
	cmd.Flags().StringVarP(&options.menu, "menu", "m", "", "Menu")

	return cmd
}
