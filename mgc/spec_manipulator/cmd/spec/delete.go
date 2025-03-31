package spec

import (
	"fmt"

	"github.com/spf13/cobra"
)

func delete(cmd *cobra.Command, args []string) {
	fmt.Println("not implemented - remover diretamente no arquivo de config")
}

var DeleteSpecsCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete spec",
	Args:  cobra.MinimumNArgs(1),
	Run:   delete,
}
