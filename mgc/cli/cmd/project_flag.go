package cmd

import (
	"github.com/MagaluCloud/magalu/mgc/core/config"
	"github.com/spf13/cobra"
)

const (
	projectFlag    = "project-id"
	iamProjectFlag = "iam-project-id"
)

func addProjectFlags(cmd *cobra.Command) {
	flags := cmd.Root().PersistentFlags()
	flags.String(
		projectFlag,
		"",
		"Project to scope the requests. Overrides the configured project",
	)
	flags.String(
		iamProjectFlag,
		"",
		"Project scope for IAM commands only. Overrides the configured IAM project",
	)
}

func getProjectFlag(cmd *cobra.Command, name string) string {
	value, err := cmd.Root().PersistentFlags().GetString(name)
	if err != nil {
		return ""
	}
	return value
}

func applyProjectFlags(cmd *cobra.Command, cfg *config.Config) {
	for flag, key := range map[string]string{
		projectFlag:    config.ProjectKey,
		iamProjectFlag: config.IamProjectKey,
	} {
		if value := getProjectFlag(cmd, flag); value != "" {
			_ = cfg.SetTempConfig(key, value)
		}
	}
}
