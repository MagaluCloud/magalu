package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"magalu.cloud/cli/ui"
	"magalu.cloud/core"
	mgcHttp "magalu.cloud/core/http"
	mgcLoggerPkg "magalu.cloud/core/logger"
	mgcSchema "magalu.cloud/core/schema"
	mgcSdk "magalu.cloud/sdk"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stoewer/go-strcase"
	"golang.org/x/exp/slices"

	flag "github.com/spf13/pflag"
)

const loggerConfigKey = "logging"

var argParser = &osArgParser{}

// -- BEGIN: create Dynamic Argument Loaders --

// TODO: likely this DynamicArgLoader is not needed anymore,
// I just converted it to use the sdk stuff without checking in detail

type DynamicArgLoader func(cmd *cobra.Command, target *string) (*cobra.Command, DynamicArgLoader, error)

// What we really need so far:
// - no target (default) or a help target: load all
// - else: run specific target (sub command)
func createCommonDynamicArgLoader(
	loadAll func(cmd *cobra.Command) error,
	loadTarget func(target string, cmd *cobra.Command) (*cobra.Command, DynamicArgLoader, error),
) DynamicArgLoader {
	return func(cmd *cobra.Command, target *string) (*cobra.Command, DynamicArgLoader, error) {
		if target == nil {
			return nil, nil, loadAll(cmd)
		}

		return loadTarget(*target, cmd)
	}
}

// -- END: create Dynamic Argument Loaders --

func handleLoaderChild(sdk *mgcSdk.Sdk, cmd *cobra.Command, child mgcSdk.Descriptor) (*cobra.Command, DynamicArgLoader, error) {
	if childGroup, ok := child.(mgcSdk.Grouper); ok {
		return AddGroup(sdk, cmd, childGroup)
	} else if childExec, ok := child.(mgcSdk.Executor); ok {
		return AddAction(sdk, cmd, childExec)
	} else {
		return nil, nil, fmt.Errorf("child %v not group/executor", child)
	}
}

func createGroupLoader(sdk *mgcSdk.Sdk, group mgcSdk.Grouper) DynamicArgLoader {
	return createCommonDynamicArgLoader(
		func(cmd *cobra.Command) error {
			_, err := group.VisitChildren(func(child mgcSdk.Descriptor) (bool, error) {
				_, _, err := handleLoaderChild(sdk, cmd, child)
				return true, err
			})

			return err
		},
		func(target string, cmd *cobra.Command) (*cobra.Command, DynamicArgLoader, error) {
			child, err := group.GetChildByName(target)

			if err != nil {
				return nil, nil, err
			}

			return handleLoaderChild(sdk, cmd, child)
		},
	)
}

func getPropType(prop *mgcSdk.Schema) string {
	result := prop.Type

	if prop.Type == "array" && prop.Items != nil {
		result += fmt.Sprintf("(%v)", getPropType((*mgcSdk.Schema)(prop.Items.Value)))
	} else if len(prop.Enum) != 0 {
		result = "enum"
	}

	return result
}

func addFlags(flags *flag.FlagSet, schema *mgcSdk.Schema) {
	for name, propRef := range schema.Properties {
		prop := propRef.Value

		propType := getPropType((*mgcSdk.Schema)(prop))
		if propType == "boolean" {
			def, _ := prop.Default.(bool)
			flags.Bool(name, def, prop.Description)
		} else {
			var value any
			if prop.Default != nil {
				value = prop.Default
			}

			constraints := fmt.Sprintf("(%s)", schemaValueConstraints((*mgcSdk.Schema)(prop)))
			description := prop.Description
			if constraints != "()" {
				if description == "" {
					description += constraints
				} else {
					description += fmt.Sprintf(" %s", constraints)
				}
			}

			f := &AnyFlagValue{value: value, typeName: propType}
			flags.AddFlag(&flag.Flag{
				Name:     name,
				DefValue: f.String(),
				Usage:    description,
				Value:    f,
			})
		}

		if slices.Contains(schema.Required, name) {
			if err := cobra.MarkFlagRequired(flags, name); err != nil {
				// Will probably never happen
				logger().Warnw(
					"unable to mark flag as required, but it should be required",
					"flag name", name,
					"error", err.Error(),
				)
			}
		}
	}
}

func getFlagValue(flags *flag.FlagSet, name string) (mgcSdk.Value, *pflag.Flag, error) {
	flag := flags.Lookup(name)
	if flag == nil {
		return nil, nil, os.ErrNotExist
	}

	if f, ok := flag.Value.(*AnyFlagValue); ok {
		return f.Value(), flag, nil
	} else if val, err := flags.GetBool(name); err == nil {
		return val, flag, nil
	} else {
		return nil, flag, fmt.Errorf("Could not get flag value %q: %w", name, err)
	}
}

func loadDataFromFlags(flags *flag.FlagSet, schema *mgcSdk.Schema, dst map[string]core.Value) error {
	if flags == nil || schema == nil || dst == nil {
		return fmt.Errorf("invalid command or parameter schema")
	}

	for name, propRef := range schema.Properties {
		propSchema := propRef.Value
		val, flag, err := getFlagValue(flags, name)
		if flag == nil {
			continue
		}
		if err != nil {
			return err
		}
		if val == nil && !mgcSchema.IsSchemaNullable((*core.Schema)(propSchema)) {
			continue
		}
		dst[name] = val
	}

	return nil
}

func loadDataFromConfig(config *mgcSdk.Config, flags *flag.FlagSet, schema *mgcSdk.Schema, dst map[string]mgcSdk.Value) error {
	for name := range schema.Properties {
		val, flag, err := getFlagValue(flags, name)
		if flag == nil {
			continue
		}

		var cfgVal any
		errCfg := config.Get(name, &cfgVal)
		if errCfg != nil {
			return errCfg
		}

		if flag.Changed || cfgVal == nil {
			if err != nil {
				return err
			}
			dst[name] = val
		} else {
			dst[name] = cfgVal
		}
	}

	return nil
}

func handleExecutorResult(ctx context.Context, cmd *cobra.Command, result core.Result, err error) error {
	if err != nil {
		var failedTerminationError core.FailedTerminationError
		if errors.As(err, &failedTerminationError) {
			_ = formatResult(cmd, failedTerminationError.Result)
		}
		return err
	}

	return formatResult(cmd, result)
}

func formatResult(cmd *cobra.Command, result core.Result) error {
	output := getOutputFor(cmd, result)

	if resultWithReader, ok := core.ResultAs[core.ResultWithReader](result); ok {
		return handleReaderResult(resultWithReader.Reader(), output)
	}

	if resultWithValue, ok := core.ResultAs[core.ResultWithValue](result); ok {
		return handleJsonResult(resultWithValue, output)
	}

	return fmt.Errorf("unsupported result: %T %+v", result, result)
}

func handleReaderResult(reader io.Reader, outFile string) (err error) {
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	var writer io.WriteCloser
	if outFile == "" || outFile == "-" {
		writer = os.Stdout
	} else {
		writer, err = os.OpenFile(outFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
	}

	n, err := io.Copy(writer, reader)
	defer writer.Close()
	if err != nil {
		return fmt.Errorf("Wrote %d bytes. Error: %w\n", n, err)
	}
	return nil
}

func handleJsonResult(result core.ResultWithValue, output string) (err error) {
	value := result.Value()
	if value == nil {
		return nil
	}

	err = result.ValidateSchema()
	if err != nil {
		logger().Warnw("result validation failed", "error", err.Error())
	}

	name, options := parseOutputFormatter(output)
	if name == "" {
		if formatter, ok := core.ResultAs[core.ResultWithDefaultFormatter](result); ok {
			fmt.Println(formatter.DefaultFormatter())
			return nil
		}
	}

	formatter, err := getOutputFormatter(name, options)
	if err != nil {
		return err
	}
	return formatter.Format(value, options)
}

func getOutputFor(cmd *cobra.Command, result core.Result) string {
	output := getOutputFlag(cmd)
	if output == "" {
		if outputOptions, ok := core.ResultAs[core.ResultWithDefaultOutputOptions](result); ok {
			return outputOptions.DefaultOutputOptions()
		}
	}

	return output
}

func handleLinkArgs(
	ctx context.Context,
	parentCmd *cobra.Command,
	linkChainedArgs [][]string,
	links map[string]core.Linker,
	config *mgcSdk.Config,
	originalResult core.Result,
) error {
	if len(linkChainedArgs) == 0 {
		return nil
	}

	currentLinkArgs := linkChainedArgs[0]
	linkName := currentLinkArgs[0]

	if link, ok := links[linkName]; ok {
		linkCmd := AddLink(ctx, parentCmd, config, originalResult, link, linkChainedArgs[1:])
		err := linkCmd.ParseFlags(currentLinkArgs[1:])
		if err != nil {
			return err
		}
		return linkCmd.RunE(linkCmd, []string{})
	} else if linkName == "help" {
		linkHelpCmd := AddLinkHelp(parentCmd)
		linkHelpCmd.Run(linkHelpCmd, nil)
		return nil
	} else {
		return fmt.Errorf("Invalid link execution. Command '%s' has no link '%s'", parentCmd.Use, linkName)
	}
}

func handleExecutor(
	ctx context.Context,
	cmd *cobra.Command,
	exec core.Executor,
	parameters core.Parameters,
	configs core.Configs,
) (core.Result, error) {

	if cExec, ok := core.ExecutorAs[core.ConfirmableExecutor](exec); ok && !getBypassConfirmationFlag(cmd) {
		msg := cExec.ConfirmPrompt(parameters, configs)
		run, err := ui.Confirm(msg)
		if err != nil {
			return nil, err
		}

		if !run {
			return nil, core.UserDeniedConfirmationError{Prompt: msg}
		}
	}

	if t := getTimeoutFlag(cmd); t > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, t)
		defer cancel()
	}

	waitTermination := getWaitTerminationFlag(cmd)
	var cb core.RetryUntilCb
	if tExec, ok := core.ExecutorAs[core.TerminatorExecutor](exec); ok && waitTermination {
		cb = func() (result core.Result, err error) {
			return tExec.ExecuteUntilTermination(ctx, parameters, configs)
		}
	} else {
		cb = func() (result core.Result, err error) {
			return exec.Execute(ctx, parameters, configs)
		}
	}

	retry, err := getRetryUntilFlag(cmd)
	if err != nil {
		return nil, err
	}

	result, err := retry.Run(ctx, cb)

	err = handleExecutorResult(ctx, cmd, result, err)
	if err != nil {
		return nil, err
	}

	return result, err
}

func AddLinkHelp(
	parentCmd *cobra.Command,
) *cobra.Command {
	linkHelpCmd := &cobra.Command{
		Use:   "help",
		Short: "Get help on the usage of link chains",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("\nExecuting link\nName: %s\nDescription: %s\n\n", cmd.Use, cmd.Short)

			msg := `All Executors might have possible links to other Executors.
These Links use the output of said Executor to aid in deciding which
parameters and configs to pass to the next Executor without needing
to declare them explicitly.

Calling these links is as simple as adding a "!" to the end of a
command call, specifying the link name and passing in any additional
parameters via flags. For instance:

./cli initial command --some-flag flag-value ! link --other-flag other

In this case, 'link' may have access to '--some-flag' (and other internal
metadata), which can aid in its execution following the 'initial command'
execution.`
			fmt.Println(msg)
		},
	}
	parentCmd.AddCommand(linkHelpCmd)
	logger().Debugw("Link added to command tree", "name", "help")
	return linkHelpCmd
}

func AddLink(
	ctx context.Context,
	parentCmd *cobra.Command,
	config *mgcSdk.Config,
	originalResult core.Result,
	link core.Linker,
	followingLinkArgs [][]string,
) *cobra.Command {
	linkCmd := &cobra.Command{
		Use:   link.Name(),
		Short: link.Description(),
		RunE: func(cmd *cobra.Command, args []string) error {
			t := table.NewWriter()
			t.AppendHeader(table.Row{"Executing link"})
			t.AppendRows([]table.Row{{"Name", link.Name()}, {"Description", link.Description()}})
			t.SetStyle(table.StyleRounded)
			fmt.Println()
			fmt.Println(t.Render())
			fmt.Println()

			additionalParameters := core.Parameters{}
			additionalConfigs := core.Configs{}

			if err := loadDataFromFlags(cmd.Flags(), link.AdditionalParametersSchema(), additionalParameters); err != nil {
				return err
			}

			if err := loadDataFromConfig(config, cmd.PersistentFlags(), link.AdditionalConfigsSchema(), additionalConfigs); err != nil {
				return err
			}

			exec, err := link.CreateExecutor(originalResult)
			if err != nil {
				return fmt.Errorf("unable to resolve link %s: %w", link.Name(), err)
			}

			result, err := handleExecutor(ctx, cmd, exec, additionalParameters, additionalConfigs)
			if err != nil {
				return err
			}

			return handleLinkArgs(ctx, cmd, followingLinkArgs, exec.Links(), config, result)
		},
	}

	addFlags(linkCmd.Flags(), link.AdditionalParametersSchema())
	addFlags(linkCmd.PersistentFlags(), link.AdditionalConfigsSchema())

	parentCmd.AddCommand(linkCmd)
	logger().Debugw("Link added to command tree", "name", link.Name())

	// Reset values of persistent flags to avoid inheriting the values set from previous actions/links
	linkCmd.PersistentFlags().Visit(func(f *flag.Flag) {
		_ = f.Value.Set(f.DefValue)
	})

	return linkCmd
}

func AddAction(
	sdk *mgcSdk.Sdk,
	parentCmd *cobra.Command,
	exec mgcSdk.Executor,
) (*cobra.Command, DynamicArgLoader, error) {
	desc := exec.(mgcSdk.Descriptor)

	actionCmd := &cobra.Command{
		Use:   desc.Name(),
		Short: desc.Description(),
		// TODO: Long:    desc.Description,

		RunE: func(cmd *cobra.Command, args []string) error {
			parameters := core.Parameters{}
			configs := core.Configs{}

			config := sdk.Config()

			if err := loadDataFromFlags(cmd.Flags(), exec.ParametersSchema(), parameters); err != nil {
				return err
			}

			if err := loadDataFromConfig(config, cmd.PersistentFlags(), exec.ConfigsSchema(), configs); err != nil {
				return err
			}

			ctx := sdk.NewContext()
			result, err := handleExecutor(ctx, cmd, exec, parameters, configs)
			if err != nil {
				return err
			}

			// First chained args structure is MainArgs
			linkChainedArgs := argParser.ChainedArgs()[1:]
			return handleLinkArgs(ctx, cmd, linkChainedArgs, exec.Links(), config, result)
		},
	}

	addFlags(actionCmd.Flags(), exec.ParametersSchema())
	addFlags(actionCmd.PersistentFlags(), exec.ConfigsSchema())

	parentCmd.AddCommand(actionCmd)
	logger().Debugw("Executor added to command tree", "name", exec.Name())
	return actionCmd, nil, nil
}

func runHelpE(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func AddGroup(
	sdk *mgcSdk.Sdk,
	parentCmd *cobra.Command,
	group mgcSdk.Grouper,
) (*cobra.Command, DynamicArgLoader, error) {
	desc := group.(mgcSdk.Descriptor)
	moduleCmd := &cobra.Command{
		Use:     desc.Name(),
		Short:   desc.Description(),
		Version: desc.Version(),
		RunE:    runHelpE,
	}

	loader := createGroupLoader(sdk, group)

	parentCmd.AddCommand(moduleCmd)
	logger().Debugw("Groupper added to command tree", "name", group.Name())
	return moduleCmd, loader, nil
}

// BEGIN: copied from cobra: command.go

// These are required to get the next unknown command and its arguments without the flags

func hasNoOptDefVal(name string, fs *flag.FlagSet) bool {
	flag := fs.Lookup(name)
	if flag == nil {
		return false
	}
	return flag.NoOptDefVal != ""
}

func shortHasNoOptDefVal(name string, fs *flag.FlagSet) bool {
	if len(name) == 0 {
		return false
	}

	flag := fs.ShorthandLookup(name[:1])
	if flag == nil {
		return false
	}
	return flag.NoOptDefVal != ""
}

func isFlagArg(arg string) bool {
	return ((len(arg) >= 3 && arg[0:2] == "--") ||
		(len(arg) >= 2 && arg[0] == '-' && arg[1] != '-'))
}

// END: copied from cobra: command.go

// this is a modified Traverse() from cobra: command.go
// we need it to get the next command skipping flags
// it's far from ideal, for instance this is a bug:
//
//	$ main-cmd sub-cmd -h sub-sub-cmd
//
// it will take `sub-sub-cmd` as -h argument (but it's not a boolean)
func getNextUnknownCommand(c *cobra.Command, args []string) (*string, []string) {
	inFlag := false

	for i, arg := range args {
		switch {
		// A long flag with a space separated value
		case strings.HasPrefix(arg, "--") && !strings.Contains(arg, "="):
			// TODO: this isn't quite right, we should really check ahead for 'true' or 'false'
			inFlag = !hasNoOptDefVal(arg[2:], c.Flags())
			continue
		// A short flag with a space separated value
		case strings.HasPrefix(arg, "-") && !strings.Contains(arg, "=") && len(arg) == 2 && !shortHasNoOptDefVal(arg[1:], c.Flags()):
			inFlag = true
			continue
		// The value for a flag
		case inFlag:
			inFlag = false
			continue
		// A flag without a value, or with an `=` separated value
		case isFlagArg(arg):
			continue
		}

		return &arg, args[i+1:]
	}

	return nil, args
}

func DynamicLoadCommand(cmd *cobra.Command, args []string, loader DynamicArgLoader) error {
	childCmd, childArgs, err := cmd.Traverse(args)
	if err != nil {
		return err
	}

	if cmd != childCmd {
		return nil
	}

	var childCmdName *string
	for {
		// NOTE: this replicates cmd.Traverse(), but returns the command and the args without flags
		childCmdName, childArgs = getNextUnknownCommand(cmd, childArgs)
		if childCmdName == nil || *childCmdName != "help" {
			break
		}
	}

	childCmd, childLoader, err := loader(cmd, childCmdName)
	if err != nil || childLoader == nil || childCmd == nil {
		return err
	}

	return DynamicLoadCommand(childCmd, childArgs, childLoader)
}

func normalizeFlagName(f *pflag.FlagSet, name string) pflag.NormalizedName {
	name = strcase.KebabCase(name)
	return pflag.NormalizedName(name)
}

func showHelpForError(cmd *cobra.Command, args []string, err error) {
	switch {
	case err == nil:
		break

	case errors.As(err, new(*mgcHttp.HttpError)),
		errors.As(err, new(*url.Error)),
		errors.As(err, new(core.FailedTerminationError)),
		errors.As(err, new(core.UserDeniedConfirmationError)),
		errors.Is(err, context.Canceled),
		errors.Is(err, context.DeadlineExceeded):
		break

	default:
		// we can't call UsageString() on the root, we need to find the actual leaf command that failed:
		subCmd, _, _ := cmd.Find(args)
		cmd.PrintErrln(subCmd.UsageString())
	}
}

// TODO: Bind config to PFlag. Investigate how to make it work correctly
func getLogFilterConfig(sdk *mgcSdk.Sdk) string {
	var logfilter string
	err := sdk.Config().Get("logfilter", &logfilter)
	if err != nil {
		return ""
	}
	return logfilter
}

func Execute() (err error) {
	sdk := &mgcSdk.Sdk{}

	rootCmd := &cobra.Command{
		Use:     "mgc",
		Version: mgcSdk.Version,
		Short:   "CLI tool for OpenAPI integration",
		Long: `This CLI is a dynamic processor of OpenAPI files that
can generate a command line on-demand for Rest manipulation`,
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE:          runHelpE,
	}
	rootCmd.SetGlobalNormalizationFunc(normalizeFlagName)
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

	if hasOutputFormatHelp(rootCmd) {
		return nil
	}

	if err = initLogger(sdk, getLogFilterFlag(rootCmd)); err != nil {
		return err
	}

	rootCmd.AddCommand(newDumpTreeCmd(sdk))

	mainArgs := argParser.MainArgs()

	err = DynamicLoadCommand(rootCmd, mainArgs, createGroupLoader(sdk, sdk.Group()))
	if err != nil {
		rootCmd.PrintErrln("Warning: loading dynamic arguments:", err)
	}

	defer func() {
		_ = mgcLoggerPkg.Root().Sync()
	}()

	rootCmd.SetArgs(mainArgs)
	err = rootCmd.Execute()
	showHelpForError(rootCmd, mainArgs, err) // since we SilenceUsage and SilenceErrors
	return err
}
