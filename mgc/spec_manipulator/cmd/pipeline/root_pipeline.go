package pipeline

import (
	"github.com/spf13/cobra"
)

func PipelineCmd() *cobra.Command {
	pipeMenu := &cobra.Command{
		Use:   "pipeline",
		Short: "Menu com opções uteis ao pipeline e ao pre-commit",
	}

	pipeMenu.AddCommand(CliDumpTreeCmd()) // download all
	pipeMenu.AddCommand(CliDocOutputCmd())
	pipeMenu.AddCommand(NewOAPIIndexCommand())
	pipeMenu.AddCommand(GetGenDocsMagaluCmd())

	return pipeMenu
}
