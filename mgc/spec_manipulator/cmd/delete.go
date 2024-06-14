package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func delete(cmd *cobra.Command, args []string) {
	fmt.Println("not implemented - remover diretamente no arquivo de config")
}

var deleteSpecsCmd = &cobra.Command{
	Use:    "delete",
	Hidden: true,
	Short:  "Delete spec",
	Args:   cobra.MinimumNArgs(1),
	Run:    delete,
}
