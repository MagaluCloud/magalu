package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func add(cmd *cobra.Command, args []string) {
	var toSave []commandsList

	module := strings.Split(args[0], " ")[1]
	if module == "" {
		fmt.Println(`fail! Command syntax eg.: "mgc auth login"`)
		return
	}
	toSave = append(toSave, commandsList{Command: args[0], Module: module})

	currentConfig, err := loadList()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, x := range currentConfig {
		if x.Command == args[0] {
			fmt.Println(`fail! command already exists`)
			return
		}
	}

	toSave = append(toSave, currentConfig...)
	viper.Set("commands", toSave)
	err = viper.WriteConfigAs(VIPER_FILE)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("done")
}

var addCommandCmd = &cobra.Command{
	Use:     "add [command]",
	Short:   "Add new command",
	Example: "specs add 'mgc auth login'",
	Args:    cobra.MinimumNArgs(1),
	Hidden:  false,
	Run:     add,
}
