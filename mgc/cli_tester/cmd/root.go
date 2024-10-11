package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	mock "magalu.cloud/cli_tester/cmd/mock"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		Use:               "cli_tester",
		Short:             "Utilitário para auxiliar nos testes da CLI",
	}
)

const (
	VIPER_FILE = "commands.yaml"
	SNAP_DIR   = "snapshot"
)

var currentDir = func() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(lisCommandsCmd)
	rootCmd.AddCommand(addCommandCmd)
	rootCmd.AddCommand(runTestsCmd)
	rootCmd.AddCommand(mock.MockCmd())
}

func initConfig() {

	ex, err := os.Executable()
	home := filepath.Dir(ex)
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(VIPER_FILE)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Command list at file:", viper.ConfigFileUsed())
	}

}
