package cmd

import (
	"fmt"
	"os"
	"regexp"
	"runtime"

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

func Execute(version string) (err error) {
	sdk := &mgcSdk.Sdk{}
	sdk.SetVersion(version)

	vv := fmt.Sprintf("%s (%s/%s)",
		version,
		runtime.GOOS,
		runtime.GOARCH)

	rootCmd := &cobra.Command{
		Use:     "mgc",
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
		SilenceErrors: true, // ####    Hack: true to avoid panic on error / false to debug error
		SilenceUsage:  true, // ####
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
	configureOutputColor(rootCmd)
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
