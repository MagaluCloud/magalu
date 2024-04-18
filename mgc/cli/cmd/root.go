package cmd

import (
	"slices"
	"strings"

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
	configureOutputColor(rootCmd, nil)
	addOutputFlag(rootCmd)
	addLogFilterFlag(rootCmd, getLogFilterConfig(sdk))
	addLogDebugFlag(rootCmd)
	addTimeoutFlag(rootCmd)
	addWaitTerminationFlag(rootCmd)
	addRetryUntilFlag(rootCmd)
	addBypassConfirmationFlag(rootCmd)
	addHideProgressFlag(rootCmd)
	addShowInternalFlag(rootCmd)
	addShowHiddenFlag(rootCmd)

	rootCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) { f.Hidden = true })

	rootCmd.InitDefaultHelpFlag()
	rootCmd.InitDefaultVersionFlag()
	addShowCliGlobalFlags(rootCmd)

	// Immediately parse flags for root command because we'll access the global flags prior
	// to calling Execute (which is when Cobra parses the flags)
	args := argParser.MainArgs()
	// This loop will parse flags even if unknown flag error arises.
	// A flag error means that ParseFlags will early return and not parse the rest of the args.
	// This happens because some flags aren't available until further down the code.
	for {
		err = rootCmd.ParseFlags(args)
		// Either we parsed all the flags or there are no more args to parse
		if err == nil || len(args) == 0 {
			break
		}
		flag, found := strings.CutPrefix(err.Error(), "unknown flag: ")
		if found && len(flag) > 0 {
			skipTo := slices.Index(args, flag)
			args = args[skipTo+1:]
		}
	}

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

	loadErr := loadSdkCommandTree(sdk, rootCmd, mainArgs)
	if loadErr != nil {
		rootCmd.PrintErrln("Warning: loading dynamic arguments:", loadErr)
	}

	defer func() {
		_ = mgcLoggerPkg.Root().Sync()
	}()

	rootCmd.SetArgs(mainArgs)

	if !getHideProgressFlag(rootCmd) {
		pb = progress_bar.New()
		go pb.Render()
		defer pb.Finalize()
	}

	err = rootCmd.Execute()
	if err == nil && loadErr != nil {
		err = loadErr
	}

	err = showHelpForError(rootCmd, mainArgs, err) // since we SilenceUsage and SilenceErrors
	return err
}
