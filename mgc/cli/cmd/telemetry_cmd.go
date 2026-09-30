package cmd

import (
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/spf13/cobra"
)

const telemetryCmdName = "telemetry"

func newTelemetryCmd(getService func() *telemetry.Service) *cobra.Command {
	telemetryCmd := &cobra.Command{
		Use:     telemetryCmdName,
		Short:   "Manage usage data collection",
		Long:    fmt.Sprintf("Manage the pseudonymous usage data collected by the CLI. Privacy policy: %s", telemetry.PrivacyPolicyURL),
		GroupID: "settings",
	}

	telemetryCmd.AddCommand(
		newTelemetryEnableCmd(getService),
		newTelemetryDisableCmd(getService),
		newTelemetryStatusCmd(getService),
	)
	return telemetryCmd
}

func newTelemetryEnableCmd(getService func() *telemetry.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "enable",
		Short: "Enable usage data collection",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc := getService()
			if err := setTelemetryDisabled(svc, false); err != nil {
				return err
			}
			if by := svc.DisabledBy(); by != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Telemetry enabled in settings, but still disabled by %s.\n", describeDisabledBy(by))
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Telemetry enabled.")
			return nil
		},
	}
}

func newTelemetryDisableCmd(getService func() *telemetry.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Short: "Disable usage data collection",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := setTelemetryDisabled(getService(), true); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Telemetry disabled.")
			return nil
		},
	}
}

func newTelemetryStatusCmd(getService func() *telemetry.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show whether usage data collection is enabled",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if by := getService().DisabledBy(); by != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Telemetry is disabled by %s.\n", describeDisabledBy(by))
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Telemetry is enabled.")
			return nil
		},
	}
}

func setTelemetryDisabled(svc *telemetry.Service, disabled bool) error {
	if err := svc.SetDisabled(disabled); err != nil {
		return fmt.Errorf("unable to save telemetry settings: %w", err)
	}
	return nil
}

func describeDisabledBy(by string) string {
	switch by {
	case telemetry.DisabledByConfig:
		return "the CLI settings (mgc telemetry disable)"
	case telemetry.DisabledByDoNotTrack:
		return "the DO_NOT_TRACK environment variable"
	case telemetry.DisabledByOptOut:
		return "the MGC_CLI_TELEMETRY_OPTOUT environment variable"
	case telemetry.DisabledByUnreadableState:
		return "an unreadable " + telemetry.StateFileName + " file (run mgc telemetry enable or disable to rewrite it)"
	}
	return by
}
