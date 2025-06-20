package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

func deleteSpecCmd() *cobra.Command {
	var menu string

	cmd := &cobra.Command{
		Use:   "delete [menu]",
		Short: "Delete spec",
		Run: func(cmd *cobra.Command, args []string) {
			currentConfig, err := loadList("")
			if err != nil {
				fmt.Println(err)
				return
			}
			newConfig := []specList{}
			for _, v := range currentConfig {
				if v.Menu != menu {
					newConfig = append(newConfig, v)
				}
			}
			saveConfig(newConfig)
			fmt.Println("Spec removed successfully")
		},
	}

	cmd.Flags().StringVarP(&menu, "menu", "m", "", "Menu to delete")
	return cmd
}
