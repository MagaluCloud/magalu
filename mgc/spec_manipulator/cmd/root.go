package cmd

import (
	"os"
	"path"
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
	viperUsedFile = ""
)

const (
	VIPER_FILE = "specs.yaml"
	SPEC_DIR   = "cli_specs"
)

var currentDir = func() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	return path.Join(exPath, SPEC_DIR)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(versionCmd)       // version
	rootCmd.AddCommand(downloadSpecsCmd) // download all
	rootCmd.AddCommand(addSpecsCmd)      // add spec
	rootCmd.AddCommand(deleteSpecsCmd)   // delete spec
	rootCmd.AddCommand(listSpecsCmd)     // list specs
	rootCmd.AddCommand(prepareToGoCmd)   // convert spec to golang
	rootCmd.AddCommand(downgradeSpecCmd) // downgrade spec

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
		viperUsedFile = viper.ConfigFileUsed()
	}

}
