package cmd

import (
	"context"
	"fmt"

	"slices"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/stoewer/go-strcase"
	"magalu.cloud/core"
	"magalu.cloud/core/utils"
	mgcSdk "magalu.cloud/sdk"
)

const (
	listLinksCmd = "cli.list-links"
)

var allExecutorChildren = []string{listLinksCmd}

func isConfigRequired(string) bool { return false }

func addChildDesc(sdk *mgcSdk.Sdk, parentCmd *cobra.Command, child core.Descriptor) (*cobra.Command, error) {
	if childGroup, ok := child.(mgcSdk.Grouper); ok {
		cmd, err := addGroup(sdk, parentCmd, childGroup)
		return cmd, err
	} else if childExec, ok := child.(mgcSdk.Executor); ok {
		cmd, err := addAction(sdk, parentCmd, childExec)
		return cmd, err
	} else {
		return nil, fmt.Errorf("child %v not group/executor", child)
	}
}

func loadGrouperChild(sdk *mgcSdk.Sdk, cmd *cobra.Command, cmdGrouper core.Grouper, childName string) (*cobra.Command, core.Descriptor, error) {
	child, err := cmdGrouper.GetChildByName(childName)
	if err != nil {
		return nil, nil, err
	}

	childCmd, err := addChildDesc(sdk, cmd, child)
	if err != nil {
		return nil, nil, err
	}

	if childExec, ok := child.(mgcSdk.Executor); ok {
		isParameterRequired := func(name string) bool {
			if slices.Contains(childExec.ParametersSchema().Required, name) {
				return !slices.Contains(childExec.PositionalArgs(), name)
			}

			return false
		}

		addFlags(childCmd.Flags(), childExec.ParametersSchema(), isParameterRequired)
		addFlags(childCmd.Root().PersistentFlags(), childExec.ConfigsSchema(), isConfigRequired)
	}

	return childCmd, child, nil
}

func loadChild(sdk *mgcSdk.Sdk, cmd *cobra.Command, cmdDesc core.Descriptor, childName string) (*cobra.Command, core.Descriptor, error) {
	if cmdGrouper, ok := cmdDesc.(core.Grouper); ok {
		return loadGrouperChild(sdk, cmd, cmdGrouper, childName)
	} else if cmdExec, ok := cmdDesc.(core.Executor); ok {
		// Manually fail if specified child is invalid
		if !slices.Contains(allExecutorChildren, childName) {
			return nil, nil, fmt.Errorf("command %q has no child named %q", cmd.Name(), childName)
		}

		// Always load all executor children
		return nil, nil, loadAllExecChildren(sdk, cmd, cmdExec)
	}

	return nil, nil, fmt.Errorf("command %q has no child named %q", cmd.Name(), childName)
}

func loadAllGrouperChildren(sdk *mgcSdk.Sdk, cmd *cobra.Command, cmdGrouper core.Grouper) error {
	_, err := cmdGrouper.VisitChildren(func(child core.Descriptor) (run bool, err error) {
		if child.IsInternal() && !getShowInternalFlag(cmd.Root()) {
			return true, nil
		}
		_, err = addChildDesc(sdk, cmd, child)
		return true, err
	})
	return err
}

func loadAllExecChildren(sdk *mgcSdk.Sdk, cmd *cobra.Command, cmdExec core.Executor) error {
	addListLinks(cmd, cmdExec)
	return nil
}

func loadAllChildren(sdk *mgcSdk.Sdk, cmd *cobra.Command, cmdDesc core.Descriptor) error {
	if cmdGrouper, ok := cmdDesc.(core.Grouper); ok {
		return loadAllGrouperChildren(sdk, cmd, cmdGrouper)
	} else if cmdExec, ok := cmdDesc.(core.Executor); ok {
		return loadAllExecChildren(sdk, cmd, cmdExec)
	}

	return nil
}

func isExistingCommand(cmd *cobra.Command, name string) bool {
	for _, c := range cmd.Commands() {
		if c.Name() == name || c.HasAlias(name) {
			return true
		}
	}
	return false
}

// these are not added before the SDK loads, but will be added afterwards and we must ignore:
var builtInCommands = []string{
	"completion", // added by cobra.Command.InitDefaultCompletionCmd() only if there are sub-commands
}

var keepLoadingCommands = []string{
	"help",
	"__complete", // actual command that does the completions
}

func loadSdkCommandTree(sdk *mgcSdk.Sdk, cmd *cobra.Command, args []string) error {
	root := sdk.Group()
	if len(args) > 0 && slices.Contains(builtInCommands, args[0]) {
		return loadAllChildren(sdk, cmd, root)
	}

	return loadCommandTree(sdk, cmd, root, args)
}

func loadCommandTree(sdk *mgcSdk.Sdk, cmd *cobra.Command, cmdDesc core.Descriptor, args []string) error {
	if cmd == nil {
		return nil
	}

	var childName *string
	var childArgs = args
	keepLoadingChildren := false
	for {
		childName, childArgs = getNextUnknownCommand(cmd, childArgs)
		if childName == nil {
			break
		} else if !slices.Contains(keepLoadingCommands, *childName) {
			keepLoadingChildren = true
			break
		}
	}

	if childName == nil {
		return loadAllChildren(sdk, cmd, cmdDesc)
	} else if isExistingCommand(cmd, *childName) {
		return nil
	}

	childCmd, childCmdDesc, err := loadChild(sdk, cmd, cmdDesc, *childName)
	if err != nil {
		// If loading specified child fails, force load all children to print in help command
		// as all available child commands
		if loadAllErr := loadAllChildren(sdk, cmd, cmdDesc); loadAllErr != nil {
			return loadAllErr
		}

		if keepLoadingChildren {
			return nil
		}

		if _, ok := cmdDesc.(core.Executor); !ok {
			return err
		}
	}

	return loadCommandTree(sdk, childCmd, childCmdDesc, childArgs)
}

func addFlags(flags *flag.FlagSet, schema *mgcSdk.Schema, isRequired func(string) bool) {
	for name, propRef := range schema.Properties {
		prop := propRef.Value

		propType := getPropType((*mgcSdk.Schema)(prop))

		// Prevents flags be added twice by Link command
		if flags.Lookup(name) != nil {
			continue
		}

		if propType == "boolean" {
			def, _ := prop.Default.(bool)
			flags.Bool(name, def, prop.Description)
		} else {
			var value any
			if prop.Default != nil {
				value = prop.Default
			}

			constraints := fmt.Sprintf("(%s)", schemaJSONRepAndConstraints((*mgcSdk.Schema)(prop), false))
			description := prop.Description
			if constraints != "()" {
				if description == "" {
					description += constraints
				} else {
					description += fmt.Sprintf(" %s", constraints)
				}
			}

			f := &anyFlagValue{value: value, typeName: propType}
			flags.AddFlag(&flag.Flag{
				Name:     name,
				DefValue: f.String(),
				Usage:    description,
				Value:    f,
			})
		}

		if isRequired(name) {
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

func buildUse(desc core.Descriptor, args []string) string {
	use := desc.Name()
	for _, name := range args {
		use += " " + fmt.Sprintf("[%s]", strcase.KebabCase(name))
	}
	return use
}

func addAction(
	sdk *mgcSdk.Sdk,
	parentCmd *cobra.Command,
	exec mgcSdk.Executor,
) (*cobra.Command, error) {
	desc := exec.(mgcSdk.Descriptor)
	links := exec.Links()

	actionCmd := &cobra.Command{
		Use:     buildUse(desc, exec.PositionalArgs()),
		Args:    cobra.MaximumNArgs(len(exec.PositionalArgs())),
		Short:   desc.Summary(),
		Long:    desc.Description(),
		Version: desc.Version(),
		GroupID: "catalog",

		RunE: func(cmd *cobra.Command, args []string) error {
			parameters := core.Parameters{}
			configs := core.Configs{}

			config := sdk.Config()

			if err := loadDataFromArgs(exec.PositionalArgs(), args, cmd.Flags()); err != nil {
				return err
			}

			if err := loadDataFromFlags(cmd.Flags(), exec.ParametersSchema(), parameters); err != nil {
				return err
			}

			// Load from 'Flags' instead of 'PersistentFlags' because Cobra merges them before executing the command
			// (and in other scenarios too). The canonical way to load the flags is always via cmd.Flags(), PersistentFlags
			// are only to be used for inserting new flags
			if err := loadDataFromConfig(config, cmd.Flags(), exec.ConfigsSchema(), configs); err != nil {
				return err
			}

			ctx := sdk.NewContext()
			result, err := handleExecutor(ctx, sdk, cmd, exec, parameters, configs)
			if err != nil {
				return err
			}

			// First chained args structure is MainArgs
			linkChainedArgs := argParser.ChainedArgs()[1:]
			return handleLinkArgs(ctx, sdk, cmd, linkChainedArgs, links, config, result)
		},
	}

	parentCmd.AddCommand(actionCmd)

	logger().Debugw("Executor added to command tree", "name", exec.Name())

	// TODO: Parse this command's flags right after its creation
	return actionCmd, nil
}

func addGroup(
	sdk *mgcSdk.Sdk,
	parentCmd *cobra.Command,
	group mgcSdk.Grouper,
) (*cobra.Command, error) {
	desc := group.(mgcSdk.Descriptor)
	moduleCmd := &cobra.Command{
		Use:     desc.Name(),
		Short:   desc.Summary(),
		Long:    desc.Description(),
		Version: desc.Version(),
		GroupID: "catalog",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	moduleCmd.AddGroup(&cobra.Group{
		ID:    "catalog",
		Title: "Commands:",
	})

	parentCmd.AddCommand(moduleCmd)
	logger().Debugw("Groupper added to command tree", "name", group.Name())
	// TODO: Parse this command's flags right after its creation
	return moduleCmd, nil
}

func addLink(
	ctx context.Context,
	sdk *mgcSdk.Sdk,
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
			printLinkExecutionTable(link.Name(), link.Description())

			additionalParameters := core.Parameters{}
			additionalConfigs := core.Configs{}

			if err := loadDataFromFlags(cmd.Flags(), link.AdditionalParametersSchema(), additionalParameters); err != nil {
				return err
			}

			// Load from 'Flags' instead of 'PersistentFlags' because Cobra merges them before executing the command
			// (and in other scenarios too). The canonical way to load the flags is always via cmd.Flags(), PersistentFlags
			// are only to be used for inserting new flags
			if err := loadDataFromConfig(config, cmd.Flags(), link.AdditionalConfigsSchema(), additionalConfigs); err != nil {
				return err
			}

			exec, err := link.CreateExecutor(originalResult)
			if err != nil {
				return fmt.Errorf("unable to resolve link %s: %w", link.Name(), err)
			}

			result, err := handleExecutor(ctx, sdk, cmd, exec, additionalParameters, additionalConfigs)
			if err != nil {
				return err
			}

			return handleLinkArgs(ctx, sdk, cmd, followingLinkArgs, exec.Links(), config, result)
		},
	}

	parentCmd.AddCommand(linkCmd)

	isParameterRequired := func(name string) bool {
		return slices.Contains(link.AdditionalParametersSchema().Required, name)
	}

	addFlags(linkCmd.Flags(), link.AdditionalParametersSchema(), isParameterRequired)
	addFlags(linkCmd.Root().PersistentFlags(), link.AdditionalConfigsSchema(), isConfigRequired)

	logger().Debugw("Link added to command tree", "name", link.Name())

	// Reset values of persistent flags to avoid inheriting the values set from previous actions/links
	linkCmd.PersistentFlags().Visit(func(f *flag.Flag) {
		_ = f.Value.Set(f.DefValue)
	})

	return linkCmd
}

func addListLinks(
	parentCmd *cobra.Command,
	sourceExec core.Executor,
) *cobra.Command {
	type LinkerListEntry struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	type LinkerList []LinkerListEntry

	listLinksCmd := &cobra.Command{
		Use:   listLinksCmd,
		Short: "List all available links for this command",
		RunE: func(cmd *cobra.Command, args []string) error {
			links := sourceExec.Links()
			result := make(LinkerList, 0, len(links))

			for linkName, link := range links {
				result = append(result, LinkerListEntry{Name: linkName, Description: link.Description()})
			}

			simplified, err := utils.SimplifyAny(result)
			if err != nil {
				return err
			}

			return handleSimpleResultValue(simplified, getOutputFlag(cmd))
		},
	}

	parentCmd.AddCommand(listLinksCmd)
	return listLinksCmd
}
