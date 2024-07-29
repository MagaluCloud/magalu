package cmd

import (
	_ "embed"
	"fmt"
	"log"
	"slices"

	"github.com/spf13/cobra"
	"github.com/stoewer/go-strcase"
	"gopkg.in/yaml.v2"
	"magalu.cloud/core"
	mgcSdk "magalu.cloud/sdk"
)

//go:embed alias.yml
var aliasFile []byte

func addChildDesc(sdk *mgcSdk.Sdk, parentCmd *cobra.Command, child core.Descriptor) (cmd *cobra.Command, flags *cmdFlags, err error) {
	if childGroup, ok := child.(mgcSdk.Grouper); ok {
		cmd, err = addGroup(sdk, parentCmd, childGroup)
		return
	} else if childExec, ok := child.(mgcSdk.Executor); ok {
		return addAction(sdk, parentCmd, childExec)
	} else {
		err = fmt.Errorf("child %v not group/executor", child)
		return
	}
}

func findChildByNameOrAliases(cmdGrouper core.Grouper, childName string) (child core.Descriptor, err error) {
	notFound, _ := cmdGrouper.VisitChildren(func(desc core.Descriptor) (run bool, err error) {
		name, aliases := getCommandNameAndAliases(desc.Name())
		if name == childName {
			child = desc
			return false, nil
		}
		for _, name := range aliases {
			if name == childName {
				child = desc
				return false, nil
			}
		}

		return true, nil
	})

	if notFound {
		err = fmt.Errorf("no command with name %q", childName)
		return
	}

	return
}

func loadGrouperChild(sdk *mgcSdk.Sdk, cmd *cobra.Command, cmdGrouper core.Grouper, childName string) (*cobra.Command, core.Descriptor, error) {
	child, err := findChildByNameOrAliases(cmdGrouper, childName)
	if err != nil {
		return nil, nil, err
	}

	childCmd, flags, err := addChildDesc(sdk, cmd, child)
	if err != nil {
		return nil, nil, err
	}

	if flags != nil {
		flags.addFlags(childCmd)
	}

	return childCmd, child, nil
}

func loadChild(sdk *mgcSdk.Sdk, cmd *cobra.Command, cmdDesc core.Descriptor, childName string) (*cobra.Command, core.Descriptor, error) {
	if cmdGrouper, ok := cmdDesc.(core.Grouper); ok {
		return loadGrouperChild(sdk, cmd, cmdGrouper, childName)
	}

	return nil, nil, fmt.Errorf("command %q has no child named %q", cmd.Name(), childName)
}

func loadAllGrouperChildren(sdk *mgcSdk.Sdk, cmd *cobra.Command, cmdGrouper core.Grouper) error {
	_, err := cmdGrouper.VisitChildren(func(child core.Descriptor) (run bool, err error) {
		if child.IsInternal() && !getShowInternalFlag(cmd.Root()) {
			return true, nil
		}

		_, _, err = addChildDesc(sdk, cmd, child)
		return true, err
	})
	if err == nil {
		err = loadSelectHelperCommand(sdk, cmd, cmdGrouper)
	}

	return err
}

func loadAllChildren(sdk *mgcSdk.Sdk, cmd *cobra.Command, cmdDesc core.Descriptor) error {
	if cmdGrouper, ok := cmdDesc.(core.Grouper); ok {
		return loadAllGrouperChildren(sdk, cmd, cmdGrouper)
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

	if !getRawOutputFlag(cmd) {
		core.NewVersionChecker(
			sdk.HttpClient().Get,
			sdk.Config().Get,
			sdk.Config().Set,
		).
			CheckVersion(mgcSdk.Version, args...)
	}

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
	keepLoadingChildren := true
	for {
		childName, childArgs = getNextUnknownCommand(cmd, childArgs)
		if childName == nil {
			break
		} else if !slices.Contains(keepLoadingCommands, *childName) {
			keepLoadingChildren = false
			break
		}
	}

	if childName == nil {
		logger().Debugw(
			"no childName, load all children",
			"descriptor", cmdDesc,
			"childArgs", childArgs,
		)
		return loadAllChildren(sdk, cmd, cmdDesc)
	} else if isExistingCommand(cmd, *childName) {
		logger().Debugw(
			"childName is an existing command",
			"childName", *childName,
			"descriptor", cmdDesc,
			"childArgs", childArgs,
		)
		return nil
	}

	childCmd, childCmdDesc, err := loadChild(sdk, cmd, cmdDesc, *childName)
	if err != nil {
		// If loading specified child fails, force load all children to print in help command
		// as all available child commands
		if loadAllErr := loadAllChildren(sdk, cmd, cmdDesc); loadAllErr != nil {
			logger().Debugw(
				"childName wasn't found and load all children failed.",
				"childName", *childName,
				"descriptor", cmdDesc,
				"childArgs", childArgs,
				"loadChildError", err,
				"loadAllChildrenError", loadAllErr,
			)
			return loadAllErr
		}

		if keepLoadingChildren {
			logger().Debugw(
				"childName wasn't found, loaded all children.",
				"childName", *childName,
				"descriptor", cmdDesc,
				"childArgs", childArgs,
				"loadChildError", err,
			)
			return nil
		}

		if isExistingCommand(cmd, *childName) {
			// for instance: loadSelectHelperCommand()
			logger().Debugw(
				"childName was created after loading all children.",
				"childName", *childName,
				"descriptor", cmdDesc,
				"childArgs", childArgs,
				"loadChildError", err,
			)
			return nil
		}

		if _, ok := cmdDesc.(core.Executor); !ok {
			logger().Debugw(
				"childName wasn't found, not an executor.",
				"childName", *childName,
				"descriptor", cmdDesc,
				"childArgs", childArgs,
				"loadChildError", err,
			)
			return err
		}
		logger().Debugw(
			"childName wasn't found, process the executor.",
			"childName", *childName,
			"executor", cmdDesc,
			"childArgs", childArgs,
			"loadChildError", err,
		)
	}

	return loadCommandTree(sdk, childCmd, childCmdDesc, childArgs)
}

func buildUse(name string, args []string) string {
	use := name
	for _, name := range args {
		use += " " + fmt.Sprintf("[%s]", name)
	}
	return use
}

func getCommandNameAndAliases(origName string) (name string, aliases []string) {
	name = strcase.KebabCase(origName)
	alias := readAlias()
	if alias[name] != nil {
		aliases = alias[name]
	}
	if name != origName {
		aliases = append(aliases, origName)
	}
	return
}

func readAlias() map[string][]string {
	m := make(map[string][]string)
	err := yaml.Unmarshal(aliasFile, &m)
	if err != nil {
		log.Fatalf("error parsing alias: %v", err)
	}
	return m
}

const mgcFooterKey string = "mgc_footer_help"

func setCustomHelpTemplate(c *cobra.Command) {
	customTemplate := c.HelpTemplate() + fmt.Sprintf("{{with $footer := index .Annotations \"%s\"}}{{if $footer}}	** {{$footer}}{{end}}\n{{end}}", mgcFooterKey)
	c.SetHelpTemplate(customTemplate)
}

func addAction(
	sdk *mgcSdk.Sdk,
	parentCmd *cobra.Command,
	exec mgcSdk.Executor,
) (actionCmd *cobra.Command, flags *cmdFlags, err error) {
	flags, err = newExecutorCmdFlags(parentCmd, exec)
	if err != nil {
		return
	}

	name, aliases := getCommandNameAndAliases(exec.Name())
	cmdPath := fmt.Sprintf("%s %s", parentCmd.CommandPath(), name)

	links := newCmdLinks(sdk, exec.Links(), cmdPath)
	if links != nil {
		flags.addExtraFlag(links.listLinksFlag)
	}

	actionCmd = &cobra.Command{
		Use:               buildUse(name, flags.positionalArgsNames()),
		Aliases:           aliases,
		Args:              flags.positionalArgsFunction,
		ValidArgsFunction: flags.validateArgs,
		Example:           flags.example(cmdPath),
		Short:             exec.Summary(),
		Long:              exec.Description(),
		Version:           exec.Version(),
		GroupID:           "catalog",
		Annotations:       make(map[string]string),

		RunE: func(cmd *cobra.Command, args []string) error {
			// First chained args structure is MainArgs
			linkChainedArgs := argParser.ChainedArgs()[1:]
			if getWatchFlag(cmd) {
				linkChainedArgs = append([][]string{{"get", "-w"}}, linkChainedArgs...)
			}

			if err := links.resolve(linkChainedArgs); err != nil {
				return err
			}

			config := sdk.Config()
			parameters, configs, err := flags.getValues(config, args)
			if err != nil {
				return err
			}

			ctx := sdk.NewContext()
			result, err := handleExecutor(ctx, sdk, cmd, exec, parameters, configs)
			if err != nil {
				return err
			}

			return links.handle(result, getOutputFlag(cmd))
		},
	}

	observations := exec.DescriptorSpec().Observations
	if len(observations) != 0 {
		actionCmd.Annotations[mgcFooterKey] = exec.DescriptorSpec().Observations
	}

	setCustomHelpTemplate(actionCmd)

	if getLink, ok := exec.Links()["get"]; ok && getLink.IsTargetTerminatorExecutor() {
		flags.addExtraFlag(newWatchFlag())
	}

	parentCmd.AddCommand(actionCmd)

	logger().Debugw("Executor added to command tree", "name", exec.Name())

	return
}

func addGroup(
	sdk *mgcSdk.Sdk,
	parentCmd *cobra.Command,
	group mgcSdk.Grouper,
) (*cobra.Command, error) {
	name, aliases := getCommandNameAndAliases(group.Name())
	moduleCmd := &cobra.Command{
		Use:     name,
		Aliases: aliases,
		Short:   group.Summary(),
		Long:    group.Description(),
		Version: group.Version(),
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
