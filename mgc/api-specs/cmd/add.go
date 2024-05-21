package cmd

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type specList struct {
	Url  string `json:"url"`
	File string `json:"file"`
	Menu string `json:"menu"`
}

func interfaceToMap(i interface{}) (map[string]interface{}, bool) {
	mapa, ok := i.(map[string]interface{})
	if !ok {
		fmt.Println("A interface não é um mapa ou mapa de interfaces.")
		return nil, false
	}
	return mapa, true
}

func write(cmd *cobra.Command, args []string) {
	var fromViveiro bool
	cmd.Flags().BoolVarP(&fromViveiro, "viveiro", "v", false, "Função utilizando viveiro")

	var toSave []specList
	var file string
	if fromViveiro {
		file = fmt.Sprintf("%s.viveiro.openapi.json", args[1])
	} else {
		file = fmt.Sprintf("%s.jaxyendy.openapi.json", args[1])

	}
	toSave = append(toSave, specList{Url: args[0], File: file, Menu: args[1]})

	currentConfig, err := loadListFromViper(toWriteViveiro(fromViveiro))
	if err != nil {
		fmt.Println(err)
		return
	}
	if slices.Contains(currentConfig, toSave[0]) {
		fmt.Println("url already exists")
		return
	}
	if !validarEndpoint(args[0]) {
		fmt.Println("url is invalid")
		return
	}

	toSave = append(toSave, currentConfig...)
	viper.Set(toWriteViveiro(fromViveiro), toSave)
	err = viper.WriteConfigAs(VIPER_FILE)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("done")
}

var writeSpecsCmd = &cobra.Command{
	Use:     "add [url] [menu]",
	Short:   "Add new spec",
	Example: "specs add https://block-storage.br-ne-1.jaxyendy.com/v1/openapi.json block-storage",
	Args:    cobra.MinimumNArgs(2),
	Run:     write,
}
