package cmd

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"slices"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/cli/ui/progress_bar"
	mgcLoggerPkg "github.com/MagaluCloud/magalu/mgc/core/logger"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stoewer/go-strcase"
)

const (
	loggerConfigKey = "logging"
	defaultRegion   = "br-se1"
	apiKeyEnvVar    = "MGC_API_KEY"
)

var argParser = &osArgParser{}

var pb *progress_bar.ProgressBar

func normalizeFlagName(f *pflag.FlagSet, name string) pflag.NormalizedName {
	name = strcase.KebabCase(name)
	return pflag.NormalizedName(name)
}

func Execute() (err error) {
	sdk := &mgcSdk.Sdk{}

	vv := fmt.Sprintf("%s (%s/%s)",
		mgcSdk.Version,
		runtime.GOOS,
		runtime.GOARCH)

	use := argParser.FullProgramPath()
	use = strings.Replace(use, "./", "", 1)

	rootCmd := &cobra.Command{
		Use:     use,
		Version: vv,
		Short:   "Magalu Cloud CLI",
		Long: `
	███╗   ███╗ ██████╗  ██████╗     ██████╗██╗     ██╗
	████╗ ████║██╔════╝ ██╔════╝    ██╔════╝██║     ██║
	██╔████╔██║██║  ███╗██║         ██║     ██║     ██║
	██║╚██╔╝██║██║   ██║██║         ██║     ██║     ██║
	██║ ╚═╝ ██║╚██████╔╝╚██████╗    ╚██████╗███████╗██║
	╚═╝     ╚═╝ ╚═════╝  ╚═════╝     ╚═════╝╚══════╝╚═╝
       
Magalu Cloud CLI is a command-line interface for the Magalu Cloud. 
It allows you to interact with the Magalu Cloud to manage your resources.
`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	rootCmd.SetGlobalNormalizationFunc(normalizeFlagName)

	rootCmd.AddGroup(&cobra.Group{
		ID:    "catalog",
		Title: "Products:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "settings",
		Title: "Settings:",
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
	addShowInternalFlag(rootCmd)
	addShowHiddenFlag(rootCmd)
	addRawOutputFlag(rootCmd)
	addApiKeyFlag(rootCmd)

	rootCmd.InitDefaultHelpFlag()
	rootCmd.InitDefaultVersionFlag()

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

		if strings.HasPrefix(err.Error(), "flag needs an argument:") {
			break
		}

		flag, found := strings.CutPrefix(err.Error(), "unknown flag: ")
		if found && len(flag) > 0 {
			skipTo := slices.IndexFunc(args, func(arg string) bool {
				return strings.Split(arg, "=")[0] == flag
			})
			args = args[skipTo+1:]
			continue
		}
		flag, found = strings.CutPrefix(err.Error(), "unknown shorthand flag: ")
		if found && len(flag) > 0 {
			flag = getLastFlag(flag)
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

	rootCmd.AddCommand(newDumpTreeCmd(sdk))

	mainArgs := argParser.MainArgs()

	loadErr := loadSdkCommandTree(sdk, rootCmd, mainArgs)
	if loadErr != nil {
		logger().Debugw("failed to load command tree", "error", loadErr)
	}

	defer func() {
		_ = mgcLoggerPkg.Root().Sync()
	}()

	rootCmd.SetArgs(mainArgs)

	if !getRawOutputFlag(rootCmd) {
		pb = progress_bar.New()
		go pb.Render()
		defer pb.Finalize()
	}

	setDefaultRegion(sdk)
	setApiKey(rootCmd, sdk)
	setKeyPair(sdk)

	err = rootCmd.Execute()
	if err == nil && loadErr != nil {
		err = loadErr
	}

	err = showHelpForError(rootCmd, mainArgs, err) // since we SilenceUsage and SilenceErrors
	return err
}

func setKeyPair(sdk *mgcSdk.Sdk) {
	objId := os.Getenv("MGC_OBJ_KEY_ID")
	objKey := os.Getenv("MGC_OBJ_KEY_SECRET")

	if objId != "" && objKey != "" {
		sdk.Config().AddTempKeyPair("apikey",
			objId,
			objKey,
		)
	}
}

func setApiKey(rootCmd *cobra.Command, sdk *mgcSdk.Sdk) {
	if key := getApiKeyFlag(rootCmd); key != "" {
		_ = sdk.Auth().SetAPIKey(key)
		return
	}

	if key := os.Getenv(apiKeyEnvVar); key != "" {
		_ = sdk.Auth().SetAPIKey(key)
		return
	}
}

func getLastFlag(s string) string {
	re := regexp.MustCompile(`-(\w)`)
	matches := re.FindAllStringSubmatch(s, -1)
	if len(matches) > 0 {
		lastMatch := matches[len(matches)-1]
		if len(lastMatch) > 1 {
			return lastMatch[0]
		}
	}
	return ""
}

func setDefaultRegion(sdk *mgcSdk.Sdk) {
	var region string
	err := sdk.Config().Get("region", &region)
	if err != nil {
		logger().Debugw("failed to get region from config", "error", err)
		return
	}
	if region == "" {
		region = defaultRegion
		err = sdk.Config().Set("region", region)
		if err != nil {
			logger().Debugw("failed to set region in config", "error", err)
			return
		}
	}
}
