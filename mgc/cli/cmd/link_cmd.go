package cmd

import (
	"fmt"

	flag "github.com/spf13/pflag"

	"github.com/MagaluCloud/magalu/mgc/cli/cmd/schema_flags"
	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

func newListLinkFlag() (f *flag.Flag) {
	const listLinkFlag = "cli.list-links"
	listLinkFlagSchema := mgcSchemaPkg.NewStringSchema()
	listLinkFlagSchema.Description = "List all available links for this command"
	listLinkFlagSchema.Enum = []any{
		"table",
		"json",
		"yaml",
	}

	f = schema_flags.NewSchemaFlag(
		mgcSchemaPkg.NewObjectSchema(map[string]*mgcSchemaPkg.Schema{
			listLinkFlag: listLinkFlagSchema,
		}, nil),
		listLinkFlag,
		listLinkFlag,
		false,
		false,
		false,
	)
	f.NoOptDefVal = "table"

	return
}

func listLinks(f *flag.Flag, links core.Links) (err error) {
	if f == nil {
		return
	}

	output := f.Value.String()
	if output == "" {
		return
	}

	type LinkerListEntry struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	type LinkerList []LinkerListEntry

	result := make(LinkerList, 0, len(links))

	for linkName, link := range links {
		if link.IsInternal() {
			continue
		}
		result = append(result, LinkerListEntry{Name: linkName, Description: link.Description()})
	}

	simplified, err := utils.SimplifyAny(result)
	if err != nil {
		return
	}

	err = handleSimpleResultValue(simplified, output)
	if err != nil {
		return
	}

	return schema_flags.ErrWantHelp
}

type cmdLinks struct {
	sdk     *mgcSdk.Sdk
	links   core.Links
	cmdPath string

	listLinksFlag *flag.Flag

	root *cobra.Command // set by initCommands()

	// these are set when resolve() calls root.Execute():
	resolvedCmd       *cobra.Command
	resolvedLink      core.Linker
	resolvedLinkFlags *cmdFlags
	args              []string

	next *cmdLinks
}

func newCmdLinks(sdk *mgcSdk.Sdk, links core.Links, cmdPath string) (c *cmdLinks) {
	if len(links) == 0 {
		return nil
	}

	c = &cmdLinks{
		sdk:           sdk,
		links:         links,
		cmdPath:       cmdPath,
		listLinksFlag: newListLinkFlag(),
	}
	c.initCommands()
	logger().Debugw("newCmdLinks", "links", links, "cmdPath", cmdPath)

	return
}

func (c *cmdLinks) resolve(chainedArgs [][]string) (err error) {
	if c == nil {
		logger().Debugw("no more links")
		return
	}

	if err = listLinks(c.listLinksFlag, c.links); err != nil {
		logger().Debugw("link list requested", "chainedArgs", chainedArgs)
		return
	}

	if len(chainedArgs) == 0 {
		logger().Debug("no more links requested")
		return
	}

	c.args = chainedArgs[0]
	if len(c.args) == 0 {
		logger().Debug("no link command")
		return
	}

	c.root.SetArgs(c.args)
	err = c.root.Execute() // will set resolved, cmdFlags and nextLinks
	if err == schema_flags.ErrWantHelp {
		return
	} else if c.resolvedCmd == nil {
		return schema_flags.ErrWantHelp // <link> -h/--help
	} else if err != nil {
		err = fmt.Errorf("unknown link %q. Use \"%s ! help\" for more information", c.args[0], c.cmdPath)
		return
	}
	logger().Debugw("resolved executor link", "args", c.args, "error", err, "links", c.links, "hasNext", c.next != nil, "chainedArgs", chainedArgs)

	chainedArgs = chainedArgs[1:]
	if getWatchFlag(c.resolvedCmd) {
		chainedArgs = append([][]string{{"get", "-w"}}, chainedArgs...)
	}

	return c.next.resolve(chainedArgs)
}

func (c *cmdLinks) handle(originalResult core.Result, parentOutputFlag string) (err error) {
	if c == nil {
		return
	}

	if len(c.args) == 0 {
		return
	}

	link := c.resolvedLink
	logger().Debugw("handling link", "link", link.Name(), "originalResult", originalResult.Source())

	ctx := originalResult.Source().Context
	exec, err := link.CreateExecutor(originalResult)
	if err != nil {
		logger().Debugw("could not create link executor", "originalResult", originalResult, "error", err, "link", link.Name())
		return
	}

	cmd := c.resolvedCmd
	setOutputFlag(cmd, parentOutputFlag)
	err = cmd.ParseFlags(c.args[1:])
	if err != nil {
		logger().Debugw("could not parse link flags", "args", c.args, "error", err, "link", link.Name())
		return
	}

	sdk := c.sdk

	config := sdk.Config()
	additionalParameters, additionalConfigs, err := c.resolvedLinkFlags.getValues(config, c.args)
	if err != nil {
		return
	}

	printLinkExecutionTable(link.Name(), link.Description(), additionalParameters, additionalConfigs)
	result, err := handleExecutor(ctx, sdk, c.resolvedCmd, exec, additionalParameters, additionalConfigs)
	if err != nil {
		logger().Debugw("handled handleExecutor link", "exec", exec, "args", c.args, "error", err)
		return
	}

	err = c.next.handle(result, getOutputFlag(c.resolvedCmd))
	logger().Debugw("handled next link", "exec", exec, "args", c.args, "error", err)
	return err
}

func (c *cmdLinks) addLinkCommand(link core.Linker) {
	linkFlags, _ := newCmdFlags(
		c.root,
		link.AdditionalParametersSchema(),
		link.AdditionalConfigsSchema(),
		nil,
		nil,
	)

	linkName, aliases := getCommandNameAndAliases(link.Name())

	linkCmd := &cobra.Command{
		Use:     linkName,
		Aliases: aliases,
		Short:   link.Description(),
		GroupID: "links",
		Run: func(cmd *cobra.Command, args []string) {
			c.resolvedCmd = cmd
			c.resolvedLink = link
			c.resolvedLinkFlags = linkFlags
			c.next = newCmdLinks(c.sdk, link.Links(), fmt.Sprintf("%s ! %s", c.cmdPath, linkName))

			logger().Debugw("resolved link", "link", link.Name(), "cmd", cmd.CommandPath(), "hasNext", c.next != nil)
		},
	}

	if getLink, ok := link.Links()["get"]; ok && getLink.IsTargetTerminatorExecutor() {
		linkFlags.addExtraFlag(newWatchFlag())
	}

	linkFlags.addFlags(linkCmd)
	c.root.AddCommand(linkCmd)
}

func (c *cmdLinks) initCommands() {
	c.root = &cobra.Command{
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		SilenceErrors:     true,
		SilenceUsage:      true,
	}
	c.root.AddGroup(&cobra.Group{
		ID:    "links",
		Title: "Available links:",
	})

	addOutputFlag(c.root)
	addWaitTerminationFlag(c.root)
	addRetryUntilFlag(c.root)
	addBypassConfirmationFlag(c.root)

	// standard UsageTemplate replacing `{{.CommandPath}} [command]` with `%[1]s ! <link>`
	c.root.SetUsageTemplate(fmt.Sprintf(`Usage:{{if .Runnable}}
  %[1]s !{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  %[1]s ! <link>{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "%[1]s ! <link> --help" for more information about a command.{{end}}
`, c.cmdPath))

	configureOutputColor(c.root)

	for _, link := range c.links {
		if link.IsInternal() {
			continue
		}
		c.addLinkCommand(link)
	}

	c.initHelp()
}

func (c *cmdLinks) initHelp() {
	c.root.InitDefaultHelpCmd()
	helpCmd, _, _ := c.root.Find([]string{"help"})
	helpCmd.Run = nil
	helpCmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
		cmd, _, e := cmd.Root().Find(args)
		if cmd == nil || e != nil {
			err = fmt.Errorf("unknown help topic %#q", args)
		} else {
			if err = cmd.Help(); err == nil {
				err = schema_flags.ErrWantHelp
			}
		}
		return
	}
}

func formatMapTable(m map[string]any) string {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Name", "Value"})
	for k, v := range m {
		t.AppendRow(table.Row{k, v})
	}
	t.SortBy([]table.SortBy{{Number: 0}})
	t.SetStyle(table.StyleDefault)
	return t.Render()
}

func printLinkExecutionTable(name, description string, parameters, configs map[string]any) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Executing link"})
	t.AppendRows([]table.Row{{"Name", name}, {"Description", description}})

	if len(parameters) > 0 {
		t.AppendRow(table.Row{"Parameters", formatMapTable(parameters)})
	}
	if len(configs) > 0 {
		t.AppendRow(table.Row{"Configs", formatMapTable(configs)})
	}

	t.SetStyle(table.StyleDefault)
	fmt.Println(t.Render())
}
