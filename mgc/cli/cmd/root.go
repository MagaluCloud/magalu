package cmd

import (
	"magalu.cloud/cli/ui/progress_bar"
	mgcLoggerPkg "magalu.cloud/core/logger"
	mgcSdk "magalu.cloud/sdk"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stoewer/go-strcase"
)

const loggerConfigKey = "logging"

var argParser = &osArgParser{}

var pb *progress_bar.ProgressBar

func normalizeFlagName(f *pflag.FlagSet, name string) pflag.NormalizedName {
	name = strcase.KebabCase(name)
	return pflag.NormalizedName(name)
}

func Execute() (err error) {
	sdk := &mgcSdk.Sdk{}

	rootCmd := &cobra.Command{
		Use:     argParser.FullProgramPath(),
		Version: mgcSdk.Version,
		Short:   "CLI tool for OpenAPI integration",
		Long: `This CLI is a dynamic processor of OpenAPI files that
can generate a command line on-demand for Rest manipulation`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	rootCmd.SetGlobalNormalizationFunc(normalizeFlagName)
	rootCmd.AddGroup(&cobra.Group{
		ID:    "catalog",
		Title: "Product catalog:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "other",
		Title: "Other commands:",
	})
	rootCmd.SetHelpCommandGroupID("other")
	rootCmd.SetCompletionCommandGroupID("other")
	addOutputFlag(rootCmd)
	addLogFilterFlag(rootCmd, getLogFilterConfig(sdk))
	addTimeoutFlag(rootCmd)
	addWaitTerminationFlag(rootCmd)
	addRetryUntilFlag(rootCmd)
	addBypassConfirmationFlag(rootCmd)
	addHideProgressFlag(rootCmd)
	addShowInternalFlag(rootCmd)

	rootCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) { f.Hidden = true })

	rootCmd.InitDefaultHelpFlag()
	rootCmd.InitDefaultVersionFlag()
	addShowCliGlobalFlags(rootCmd)

	// Immediately parse flags for root command because we'll access the global flags prior
	// to calling Execute (which is when Cobra parses the flags)
	_ = rootCmd.ParseFlags(argParser.MainArgs())

	if hasOutputFormatHelp(rootCmd) {
		return nil
	}

	if err = initLogger(sdk, getLogFilterFlag(rootCmd)); err != nil {
		return err
	}

	if getShowCliGlobalFlags(rootCmd) {
		rootCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) { f.Hidden = false })
	}

	rootCmd.AddCommand(newDumpTreeCmd(sdk))

	mainArgs := argParser.MainArgs()

	err = loadSdkCommandTree(sdk, rootCmd, mainArgs)
	if err != nil {
		rootCmd.PrintErrln("Warning: loading dynamic arguments:", err)
	}

	defer func() {
		_ = mgcLoggerPkg.Root().Sync()
	}()

	rootCmd.SetArgs(mainArgs)

	if !getHideProgressFlag(rootCmd) {
		pb = progress_bar.New()
		defer pb.Stop()
	}

	err = rootCmd.Execute()
	showHelpForError(rootCmd, mainArgs, err) // since we SilenceUsage and SilenceErrors
	return err
}
