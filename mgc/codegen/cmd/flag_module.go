package cmd

import (
	"github.com/spf13/cobra"
)

var (
	moduleName string
)

const (
	flagModuleName    = "module"
	defaultModuleName = "magalu.cloud/lib"
)

func addModuleNameFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&moduleName, flagModuleName, "m", defaultModuleName, "SDK module name to generate")
}

func getModuleNameFlag(cmd *cobra.Command) string {
	if moduleName == "" {
		return defaultModuleName
	}
	return moduleName
}
