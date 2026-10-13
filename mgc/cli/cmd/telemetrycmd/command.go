package telemetrycmd

import (
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/spf13/cobra"
)

const CommandName = "telemetry"

func NewCommand(getService func() *telemetry.Service, catalogGroupID string) *cobra.Command {
	telemetryCmd := &cobra.Command{
		Use:   CommandName,
		Short: "Manage usage data collection",
		Long: fmt.Sprintf(`Manage the pseudonymous usage data collected by the CLI.
After each command the CLI sends one event with the command path, the names of
the flags that were set (never their values), the outcome and a failure category,
the duration, the CLI version, the operating system, the install method, the
execution context (interactive, ci or agent), the tenant ID, the last request ID,
the effective region used by the call, and an anonymous installation ID. Flag
values, credentials, error messages and command output are never sent, except
region: a closed infrastructure enum (datacenter region), never user-entered
data.
The event is sent in the background by a separate process, so the command never
waits for the network. With --debug, the sending process writes its diagnostics
to telemetry-debug.txt in the CLI config directory.
To opt out, run 'mgc telemetry disable' or set %s=1 or %s=1.
Privacy policy: %s`, telemetry.EnvOptOut, telemetry.EnvDoNotTrack, telemetry.PrivacyPolicyURL),
		GroupID: "settings",
	}

	// catalogGroupID é o mesmo grupo "Commands:" dos subcomandos de cada produto, para
	// o Cobra mostrar "Commands:" em vez de "Available Commands:" e o gerador de docs
	// reconhecer a seção
	telemetryCmd.AddGroup(&cobra.Group{ID: catalogGroupID, Title: "Commands:"})
	for _, subcommand := range []*cobra.Command{
		newEnableCmd(getService),
		newDisableCmd(getService),
		newStatusCmd(getService),
	} {
		subcommand.GroupID = catalogGroupID
		telemetryCmd.AddCommand(subcommand)
	}
	return telemetryCmd
}

func newEnableCmd(getService func() *telemetry.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "enable",
		Short: "Enable usage data collection",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc := getService()
			if err := setDisabled(svc, false); err != nil {
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

func newDisableCmd(getService func() *telemetry.Service) *cobra.Command {
	return &cobra.Command{
		Use:   "disable",
		Short: "Disable usage data collection",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := setDisabled(getService(), true); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Telemetry disabled.")
			return nil
		},
	}
}

func newStatusCmd(getService func() *telemetry.Service) *cobra.Command {
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

func setDisabled(svc *telemetry.Service, disabled bool) error {
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
