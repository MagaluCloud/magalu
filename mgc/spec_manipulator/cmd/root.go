package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"magalu.cloud/actions/cmd/pipeline"
	"magalu.cloud/actions/cmd/spec"
)

var (
	rootCmd = &cobra.Command{
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Use:               "mgc_action",
		Short:             "Utilitário para auxiliar na atualização de specs",
		Long:              `Uma, ou mais uma CLI para ajudar no processo de atualização das specs.`,
	}

	// specsCmd = &cobra.Command{
	// 	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	// 	Use:               "specs",
	// 	Short:             "Utilitário para auxiliar na atualização de specs",
	// 	Long:              `Uma, ou mais uma CLI para ajudar no processo de atualização das specs.`,
	// 	GroupID:           "specs",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Println("Este é o menu principal")
	// 	},
	// }

	// pipelineCmd = &cobra.Command{
	// 	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	// 	Use:               "pipeline",
	// 	Short:             "Utilitário para auxiliar na atualização de specs",
	// 	Long:              `Uma, ou mais uma CLI para ajudar no processo de atualização das specs.`,
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		fmt.Println("Este é o menu principal")
	// 	},
	// }
)

const (
	VIPER_FILE = "specs.yaml"
	SPEC_DIR   = "cli_specs"
)

// Execute executes the root command.
func Execute() error {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(spec.SpecCmd())
	rootCmd.AddCommand(pipeline.PipelineCmd())
	rootCmd.AddCommand(versionCmd) // version

	return rootCmd.Execute()
}

func initConfig() {

	ex, err := os.Executable()
	home := filepath.Dir(ex)
	cobra.CheckErr(err)

	// Search config in home directory with name ".cobra" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(VIPER_FILE)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

}
