package cmd

import (
	"os"
	"path/filepath"

	"github.com/MagaluCloud/magalu/mgc/spec_manipulator/cmd/pipeline"
	"github.com/MagaluCloud/magalu/mgc/spec_manipulator/cmd/spec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Use:               "cicd",
		Short:             "Utilitário para auxiliar nos pipelines de CI/CD",
		Long:              `Uma, ou mais uma CLI para ajudar no processo de atualização das specs.`,
	}
	viperUsedFile = ""
)

// Execute executes the root command.
func Execute() error {
	rootCmd.AddCommand(spec.SpecCmd())
	rootCmd.AddCommand(pipeline.PipelineCmd())
	rootCmd.AddCommand(versionCmd) // version
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {

	ex, err := os.Executable()
	home := filepath.Dir(ex)
	cobra.CheckErr(err)

	// Search config in home directory with name ".cobra" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(spec.VIPER_FILE)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		viperUsedFile = viper.ConfigFileUsed()
	}

}
