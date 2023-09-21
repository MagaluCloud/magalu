package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"magalu.cloud/core"
	mgcLogger "magalu.cloud/core/logger"
	mgcSdk "magalu.cloud/sdk"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stoewer/go-strcase"
	"golang.org/x/exp/slices"

	flag "github.com/spf13/pflag"
)

const loggerConfigKey = "logging"

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
				logger().Warnf("Error marking %s as required: %s\n", name, err)
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
		if val == nil && !core.IsSchemaNullable((*core.Schema)(propSchema)) {
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

		cfgVal := config.Get(name)
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
		logger().Warnf("Warning: result validation failed: %v", err)
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

			if err := loadDataFromFlags(cmd.Flags(), exec.ParametersSchema(), parameters); err != nil {
				return err
			}

			if err := loadDataFromConfig(sdk.Config(), cmd.PersistentFlags(), exec.ConfigsSchema(), configs); err != nil {
				return err
			}

			ctx := sdk.NewContext()
			if t := getTimeoutFlag(cmd); t > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, t)
				defer cancel()
			}

			waitTermination := getWaitTerminationFlag(cmd)
			var cb retryUntilCb
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
				return err
			}
			result, err := retry.run(ctx, cb)

			return handleExecutorResult(ctx, cmd, result, err)
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

func Execute() (err error) {
	rootCmd := &cobra.Command{
		Use:     "mgc",
		Version: "TODO",
		Short:   "CLI tool for OpenAPI integration",
		Long: `This CLI is a dynamic processor of OpenAPI files that
can generate a command line on-demand for Rest manipulation`,
		RunE: runHelpE,
	}
	rootCmd.SetGlobalNormalizationFunc(normalizeFlagName)
	rootCmd.AddGroup(&cobra.Group{
		ID:    "other",
		Title: "Other commands:",
	})
	rootCmd.SetHelpCommandGroupID("other")
	rootCmd.SetCompletionCommandGroupID("other")
	addOutputFlag(rootCmd)
	addLogFilterFlag(rootCmd)
	addTimeoutFlag(rootCmd)
	addWaitTerminationFlag(rootCmd)
	addRetryUntilFlag(rootCmd)

	sdk := &mgcSdk.Sdk{}
	if err = initLogger(sdk, getLogFilterFlag(rootCmd)); err != nil {
		return err
	}

	rootCmd.AddCommand(newDumpTreeCmd(sdk))

	err = DynamicLoadCommand(rootCmd, os.Args[1:], createGroupLoader(sdk, sdk.Group()))
	if err != nil {
		rootCmd.PrintErrln("Warning: loading dynamic arguments:", err)
	}

	defer func() {
		_ = mgcLogger.Root().Sync()
	}()

	return rootCmd.Execute()
}
