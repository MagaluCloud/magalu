package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Use:               "specs",
		Short:             "Utilitário para auxiliar na atualização de specs",
		Long:              `Uma, ou mais uma CLI para ajudar no processo de atualização das specs.`,
	}
)

const (
	VIPER_FILE = "specs.yaml"
	SPEC_DIR   = "cli_specs"
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

var fromViveiro bool

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(versionCmd)
	downloadSpecsCmd.Flags().BoolVarP(&fromViveiro, "viveiro", "v", false, "Função utilizando viveiro")
	rootCmd.AddCommand(downloadSpecsCmd)
	writeSpecsCmd.Flags().BoolVarP(&fromViveiro, "viveiro", "v", false, "Função utilizando viveiro")
	rootCmd.AddCommand(writeSpecsCmd)
	rootCmd.AddCommand(deleteSpecsCmd)
}

func initConfig() {

	home, err := filepath.Abs(".")

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
