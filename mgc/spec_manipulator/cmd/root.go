package cmd

import (
	"github.com/MagaluCloud/magalu/mgc/spec_manipulator/cmd/pipeline"
	"github.com/MagaluCloud/magalu/mgc/spec_manipulator/cmd/spec"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Use:               "cicd",
		Short:             "Utilitário para auxiliar nos pipelines de CI/CD",
		Long:              `Uma, ou mais uma CLI para ajudar no processo de atualização das specs.`,
	}
)

func Execute() error {
	rootCmd.AddCommand(spec.SpecCmd())
	rootCmd.AddCommand(pipeline.PipelineCmd())
	rootCmd.AddCommand(versionCmd) // version
	return rootCmd.Execute()
}
