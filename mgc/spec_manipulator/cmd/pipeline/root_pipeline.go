package pipeline

import (
	"github.com/spf13/cobra"
)

func PipelineCmd() *cobra.Command {
	pipeMenu := &cobra.Command{
		Use:   "pipeline",
		Short: "Exibe o submenu",
	}

	pipeMenu.AddCommand(CliDumpTreeCmd()) // download all
	pipeMenu.AddCommand(CliDocOutputCmd())

	return pipeMenu
}
