package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"slices"

	"github.com/MagaluCloud/magalu/mgc/cli/cmd/schema_flags"
	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// It will solve name conflicts between existing flags, parameters and configs,
// will handle "positional arguments" flags and more.
//
// each flag.Value is guaranteed to be schema_flags.SchemaFlagValue

type cmdFlags struct {
	schemaFlags    []*flag.Flag // only schema_flags.SchemaFlagValue elements
	positionalArgs []*flag.Flag // subset schemaFlags that can be positional, in order
	extraFlags     []*flag.Flag
	childFlags     []*flag.Flag

	knownFlags map[flag.NormalizedName]*flag.Flag // all known flags, both existing and schemaFlags
}

const childFlagSeparator = '.'

// Public Methods:

func (cf *cmdFlags) example(cmdPath string) string {
	type flagExample struct{ name, example string }
	var examples []flagExample

	for _, f := range cf.knownFlags {
		if e := getFlagFormattedExample(f); e != "" {
			examples = append(examples, flagExample{f.Name, e})
		}
	}

	if len(examples) == 0 {
		return ""
	}

	slices.SortFunc(examples, func(a, b flagExample) int {
		return strings.Compare(a.example, b.example)
	})

	// build the usage string without duplicates when --flag.prop= is used alongside --flag, otherwise we get duplicates
	// prefer the --flag.prop version
	var flagUsageExample string
	var lastFlag string
	for _, e := range examples {
		// sort above uses --flag=x which comes after --flag.prop=y
		if strings.HasPrefix(lastFlag, e.name) && len(e.name) < len(lastFlag) && lastFlag[len(e.name)] == childFlagSeparator {
			continue
		}
		lastFlag = e.name
		if flagUsageExample != "" {
			flagUsageExample += " "
		}
		flagUsageExample += e.example
	}

	return fmt.Sprintf("  %s %s", cmdPath, flagUsageExample)
}

func (cf *cmdFlags) positionalArgsArrayToExpand() int {
	count := 0
	first := -1
	for i, f := range cf.positionalArgs {
		if fv, ok := f.Value.(schema_flags.SchemaFlagValue); ok {
			if fv.Desc().Schema.Type != nil && fv.Desc().Schema.Type.Includes(openapi3.TypeArray) {
				first = i
				count++
				if count > 1 {
					return -1
				}
			}
		}
	}

	return first
}

// these are the public/external/user-visible names
func (cf *cmdFlags) positionalArgsNames() (names []string) {
	if len(cf.positionalArgs) == 0 {
		return
	}

	hasExpandedArrays := cf.positionalArgsArrayToExpand() >= 0

	names = make([]string, len(cf.positionalArgs))
	for i, f := range cf.positionalArgs {
		if hasExpandedArrays {
			if fv, ok := f.Value.(schema_flags.SchemaFlagValue); ok {
				if fv.Desc().Schema.Type != nil && fv.Desc().Schema.Type.Includes(openapi3.TypeArray) {
					names[i] = f.Name + "..."
					continue
				}
			}
		}
		names[i] = f.Name
	}

	return
}

func applyPositionalArgs(positionalArgs []*flag.Flag, args []string) (err error) {
	if len(positionalArgs) < len(args) {
		panic("programming error: len(positionalArgs) < len(args)")
	}

	for i, value := range args {
		f := positionalArgs[i]
		if f.Value.String() == "" && !f.Changed {
			if err = f.Value.Set(value); err != nil {
				err = fmt.Errorf("invalid argument for %s: %s", f.Name, err.Error())
				return
			}
		}
	}

	return
}

func (cf *cmdFlags) positionalArgsArrays(toExpand int, args []string) (err error) {
	nArgs := len(args)
	nPositionalArgs := len(cf.positionalArgs)

	// pre:
	//   [flag1] [flag2] [array...]     => toExpand == 2, nPositionalArgs == 3
	//   value1                         => nArgs = 1, endPre = 1 (array is not populated)
	//   value1  value2  value3 value4  => nArgs = 4, endPre = 2 (array: value3, value4)
	endPre := toExpand
	if endPre > nArgs {
		endPre = nArgs
	}
	if endPre > 0 {
		if err = applyPositionalArgs(cf.positionalArgs[:endPre], args[:endPre]); err != nil {
			return
		}
	}

	// post:
	//   [flag1] [flag2] [array...] [flag3] [flag4]   => toExpand == 2, nPositionalArgs == 5, nPost = 2
	//   value1                                       => nArgs = 1, startPost = -1 (array is not populated)
	//   value1  value2  value3 value4                => nArgs = 4, startPost = 2  (array is not populated)
	//   value1  value2  value3 value4 value5 value6  => nArgs = 6, startPost = 4  (array: value3, value4)
	nPost := nPositionalArgs - toExpand - 1
	startPost := nArgs - nPost
	if toExpand <= startPost {
		if err = applyPositionalArgs(cf.positionalArgs[toExpand+1:], args[startPost:]); err != nil {
			return
		}

		// actual array to handle
		f := cf.positionalArgs[toExpand]
		for _, value := range args[toExpand:startPost] {
			if err = f.Value.Set(value); err != nil {
				err = fmt.Errorf("invalid argument for %s: %s", f.Name, err.Error())
				return
			}
		}
	}

	return
}

func (cf *cmdFlags) positionalArgsFunction(cmd *cobra.Command, args []string) (err error) {
	if toExpand := cf.positionalArgsArrayToExpand(); toExpand >= 0 {
		return cf.positionalArgsArrays(toExpand, args)
	}

	numArgs := len(args)
	numPositionalArgs := len(cf.positionalArgs)

	if numArgs > numPositionalArgs {
		if numPositionalArgs == 0 {
			return fmt.Errorf("this command does not accept positional arguments, %d given: %s", numArgs, strings.Join(args, ", "))
		}
		return fmt.Errorf("this command receives at most %d positional arguments, %d given", numPositionalArgs, numArgs)
	}

	positionalToPopulate := 0
	for _, f := range cf.positionalArgs {
		if f.Value.String() == "" && !f.Changed {
			positionalToPopulate++
		}
	}
	if numArgs > positionalToPopulate {
		return fmt.Errorf("this command has one or more invalid positional arguments, given: %s", strings.Join(args, ", "))
	}

	return applyPositionalArgs(cf.positionalArgs[:len(args)], args)
}

func completeEnum(f *flag.Flag, toComplete string, completions []string) []string {
	fv, ok := f.Value.(schema_flags.SchemaFlagValue)
	if !ok {
		return completions
	}

	var prefixMatches, containsMatches, nonMatches []string

	for _, v := range fv.Desc().Schema.Enum {
		s, ok := v.(string)
		if !ok {
			if data, err := json.Marshal(v); err != nil {
				s = string(data)
			}
		}
		if strings.HasPrefix(s, toComplete) {
			prefixMatches = append(prefixMatches, s)
		} else if strings.Contains(s, toComplete) {
			containsMatches = append(containsMatches, s)
		} else {
			nonMatches = append(nonMatches, s)
		}
	}

	if len(prefixMatches) > 0 {
		return append(completions, prefixMatches...)
	}
	if len(containsMatches) > 0 {
		return append(completions, containsMatches...)
	}

	return append(completions, nonMatches...)
}

func (cf *cmdFlags) validateArgsArrays(toExpand int, cmd *cobra.Command, args []string, toComplete string) (completions []string, directive cobra.ShellCompDirective) {
	if len(args) < toExpand {
		f := cf.positionalArgs[len(args)]

		return cf.completeFlag(f, cmd, args, toComplete, completions)
	}

	f := cf.positionalArgs[toExpand]

	completions = cobra.AppendActiveHelp(completions, "The following arguments are accepted: multiple "+getFlagActiveHelp(f))
	completions, directive = cf.completeFlagValues(f, cmd, args, toComplete, completions)

	for _, f := range cf.positionalArgs[toExpand+1:] {
		var curDirective cobra.ShellCompDirective
		completions[0] += " or " + getFlagActiveHelp(f)
		completions, curDirective = cf.completeFlagValues(f, cmd, args, toComplete, completions)
		directive |= curDirective
	}

	return
}

func (cf *cmdFlags) validateArgs(cmd *cobra.Command, args []string, toComplete string) (completions []string, directive cobra.ShellCompDirective) {
	logger().Debug("validateArgs", "cmd", cmd.Use, "args", args, "toComplete", toComplete)
	if toExpand := cf.positionalArgsArrayToExpand(); toExpand >= 0 {
		return cf.validateArgsArrays(toExpand, cmd, args, toComplete)
	}

	directive = cobra.ShellCompDirectiveNoFileComp
	if len(cf.positionalArgs) == 0 {
		completions = cobra.AppendActiveHelp(completions, "This command does not take any arguments")
		return
	}
	if len(args) >= len(cf.positionalArgs) {
		completions = cobra.AppendActiveHelp(completions, "This command does not take any more arguments")
		return
	}

	f := cf.positionalArgs[len(args)]

	return cf.completeFlag(f, cmd, args, toComplete, completions)
}

func getFlagActiveHelp(f *flag.Flag) string {
	if description := getFlagDescription(f); description != "" {
		return fmt.Sprintf("%s (%s)", f.Name, description)
	}
	return f.Name
}

func (cf *cmdFlags) completeFlag(f *flag.Flag, cmd *cobra.Command, args []string, toComplete string, completions []string) ([]string, cobra.ShellCompDirective) {
	completions = cobra.AppendActiveHelp(completions, getFlagActiveHelp(f))
	return cf.completeFlagValues(f, cmd, args, toComplete, completions)
}

func (cf *cmdFlags) completeFlagValues(f *flag.Flag, cmd *cobra.Command, args []string, toComplete string, completions []string) ([]string, cobra.ShellCompDirective) {
	var directive cobra.ShellCompDirective

	switch f.Value.Type() {
	case "enum":
		completions = completeEnum(f, toComplete, completions)
		directive = cobra.ShellCompDirectiveNoFileComp

	case schema_flags.FlagTypeFile:
		directive = cobra.ShellCompDirectiveDefault

	case schema_flags.FlagTypeDirectory:
		directive = cobra.ShellCompDirectiveFilterDirs

	default:
		if f.DefValue != "" {
			completions = append(completions, f.DefValue)
		}
		if strings.HasPrefix(toComplete, schema_flags.ValueLoadJSONFromFilePrefix) || strings.HasPrefix(toComplete, schema_flags.ValueLoadVerbatimFromFilePrefix) {
			directive = cobra.ShellCompDirectiveDefault
		} else {
			directive = cobra.ShellCompDirectiveNoFileComp
		}
	}

	return completions, directive
}

func (cf *cmdFlags) newCompleteFlagFunc(f *flag.Flag) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cf.completeFlag(f, cmd, args, toComplete, nil)
	}
}

func (cf *cmdFlags) addFlags(cmd *cobra.Command) {
	configFlags := cmd.Root().PersistentFlags()
	parametersFlags := cmd.Flags()
	showHiddenFlags := getShowHiddenFlag(cmd)
	for _, f := range cf.schemaFlags {
		var flags *flag.FlagSet
		desc := f.Value.(schema_flags.SchemaFlagValue).Desc()
		if desc.IsConfig {
			flags = configFlags
		} else {
			flags = parametersFlags
		}
		logger().Debugw("adding schema flag", "flag", f.Name, "desc", desc)
		if showHiddenFlags {
			f.Hidden = !showHiddenFlags
		}
		flags.AddFlag(f)
		_ = cmd.RegisterFlagCompletionFunc(f.Name, cf.newCompleteFlagFunc(f))
	}

	for _, f := range cf.childFlags {
		var flags *flag.FlagSet
		desc := f.Value.(schema_flags.SchemaFlagValue).Desc()
		if desc.IsConfig {
			flags = configFlags
		} else {
			flags = parametersFlags
		}
		logger().Debugw("adding child flag", "flag", f.Name, "desc", desc)
		flags.AddFlag(f)
		_ = cmd.RegisterFlagCompletionFunc(f.Name, cf.newCompleteFlagFunc(f))
	}

	for _, f := range cf.extraFlags {
		logger().Debugw("adding extra flag", "flag", f.Name, "value", f.Value)
		parametersFlags.AddFlag(f)
		_ = cmd.RegisterFlagCompletionFunc(f.Name, cf.newCompleteFlagFunc(f))
	}
}

func (cf *cmdFlags) addExtraFlag(f *flag.Flag) {
	cf.knownFlags[flag.NormalizedName(f.Name)] = f
	cf.extraFlags = append(cf.extraFlags, f)
}

// parse, then validate flags and return the final values.
//
// flags that can be positional will be loaded from `argValues`,
// flags that were not set but could be set via configuration, will
// be loaded from `config`.
func (cf *cmdFlags) getValues(config *mgcSdk.Config, argValues []string) (core.Parameters, core.Configs, error) {
	parameters := core.Parameters{}
	configs := core.Configs{}

	var loadErrors utils.MultiError
	var missingRequiredFlags requiredFlagsError

	for _, f := range cf.childFlags {
		if f.Value.String() == schema_flags.ValueHelpIsRequired {
			showFlagHelp(f)
			return nil, nil, schema_flags.ErrWantHelp
		}
	}

	for _, f := range cf.schemaFlags {
		var value any
		value, err := schema_flags.GetFlagValue(f, config)
		logger().Debugw("parsed flag", "flag", f.Name, "desc", f.Value.(schema_flags.SchemaFlagValue).Desc(), "value", value, "error", err)
		if err == schema_flags.ErrNoFlagValue {
			continue
		} else if err == schema_flags.ErrRequiredFlag {
			missingRequiredFlags = append(missingRequiredFlags, f)
		} else if err == schema_flags.ErrWantHelp {
			showFlagHelp(f)
			return nil, nil, err
		} else if err != nil {
			loadErrors = append(loadErrors, &flagError{Flag: f, Err: err})
		} else {
			desc := f.Value.(schema_flags.SchemaFlagValue).Desc()
			if desc.IsConfig {
				configs[desc.PropName] = value
			} else {
				parameters[desc.PropName] = value
			}
		}
	}

	if len(missingRequiredFlags) > 0 {
		loadErrors = append(loadErrors, missingRequiredFlags)
	}

	if len(loadErrors) > 0 {
		return nil, nil, core.UsageError{Err: loadErrors}
	}

	return parameters, configs, nil
}

func newCmdFlags(
	parentCmd *cobra.Command, // used to discover existing flags
	parametersSchema, configsSchema *mgcSdk.Schema,
	positionalArgs, hiddenFlags []string, // names must match parameterSchema.Properties keys

) (cf *cmdFlags, err error) {
	schemaFlagsLen := len(parametersSchema.Properties) + len(configsSchema.Properties)

	cf = &cmdFlags{
		knownFlags:  make(map[flag.NormalizedName]*flag.Flag, schemaFlagsLen),
		schemaFlags: make([]*flag.Flag, 0, schemaFlagsLen),
	}

	parentFlags := parentCmd.Flags()
	parentFlags.VisitAll(cf.addExistingFlag)
	parentCmd.Root().Flags().VisitAll(cf.addExistingFlag)

	normalizeFunc := parentFlags.GetNormalizeFunc()
	normalizeName := func(name string) flag.NormalizedName {
		return normalizeFunc(parentFlags, name)
	}

	err = cf.addParametersFlags(parametersSchema, positionalArgs, hiddenFlags, normalizeName)
	cf.addConfigsFlags(configsSchema, normalizeName)

	return
}

func newExecutorCmdFlags(parentCmd *cobra.Command, exec core.Executor) (*cmdFlags, error) {
	return newCmdFlags(
		parentCmd,
		exec.ParametersSchema(),
		exec.ConfigsSchema(),
		exec.PositionalArgs(),
		exec.HiddenFlags(),
	)
}

// Internal Methods:

func (cf *cmdFlags) addExistingFlag(existingFlag *flag.Flag) {
	cf.knownFlags[flag.NormalizedName(existingFlag.Name)] = existingFlag
}

const (
	originalControlPrefix = "_"
	targetControlPrefix   = "control."
)

func (cf *cmdFlags) addSchemaFlag(
	container *mgcSdk.Schema,
	propName string,
	conflictPrefix flag.NormalizedName, // used if propName already exists
	normalizeName func(name string) flag.NormalizedName,
	isRequired bool,
	isConfig bool,
	isHidden bool,
) (f *flag.Flag) {
	baseFlagName, isControl := strings.CutPrefix(propName, originalControlPrefix)
	if isControl {
		baseFlagName = targetControlPrefix + baseFlagName
	}

	flagName := normalizeName(baseFlagName)
	for cf.knownFlags[flagName] != nil {
		flagName = conflictPrefix + flagName
	}

	f = schema_flags.NewSchemaFlag(
		container,
		propName,
		flagName,
		isRequired,
		isConfig,
		isHidden,
	)
	cf.knownFlags[flagName] = f
	cf.schemaFlags = append(cf.schemaFlags, f)

	cf.addFlagChildren(f, normalizeName)
	return
}

func (cf *cmdFlags) addFlagChildren(
	parent *flag.Flag,
	normalizeName func(name string) flag.NormalizedName,
) {
	fv, ok := parent.Value.(schema_flags.SchemaFlagValue)
	if !ok {
		return
	}
	for _, childDesc := range fv.Desc().ChildrenFlags() {
		cf.addChildSchemaFlag(parent, childDesc, normalizeName)
	}
}

func (cf *cmdFlags) addChildSchemaFlag(
	parent *flag.Flag,
	desc schema_flags.SchemaFlagValueDesc,
	normalizeName func(name string) flag.NormalizedName,
) {
	fv := parent.Value.(schema_flags.SchemaFlagValue)
	parentDesc := fv.Desc()

	desc.FlagName = parentDesc.FlagName + flag.NormalizedName(childFlagSeparator) + desc.FlagName + normalizeName(desc.PropName)

	parentShortDescription := parentDesc.Schema.Title
	if parentShortDescription == "" {
		parentShortDescription = parentDesc.Schema.Description
	}
	if i := strings.IndexFunc(parentShortDescription, func(r rune) bool { return r == ',' || r == '.' || r == ':' || r == '(' }); i > 0 {
		parentShortDescription = strings.Trim(parentShortDescription[:i], "\t\n\r ")
	}
	if parentShortDescription == "" {
		parentShortDescription = fmt.Sprintf("%s's %s property", parent.Name, desc.PropName)
	}

	proxy := schema_flags.ProxyFlagSpec{
		Usage: func() string {
			return fmt.Sprintf(
				"%s: %s\nThis is the same as '--%s=%s:%s'.",
				parentShortDescription,
				desc.Usage(),
				parent.Name,
				desc.PropName,
				desc.FlagType(),
			)
		},
		Set: func(rawValue string) error {
			return parent.Value.Set(fmt.Sprintf(`%q=%s`, desc.PropName, rawValue))
		},

		Parse: func(rawValue string) (value any, err error) {
			container, err := fv.Parse()
			if err != nil || container == nil {
				return
			}
			m, ok := container.(map[string]any)
			if !ok {
				err = fmt.Errorf("expected object flag to return map[string]any, got %T", container)
				return
			}
			value = m[desc.PropName]
			return
		},
	}

	f := schema_flags.NewProxyFlag(desc, proxy)
	cf.knownFlags[flag.NormalizedName(f.Name)] = f
	cf.childFlags = append(cf.childFlags, f)

	cf.addFlagChildren(f, normalizeName)
}

func (cf *cmdFlags) addParametersFlags(
	parametersSchema *mgcSdk.Schema,
	positionalArgs []string,
	hiddenFlags []string,
	normalizeName func(name string) flag.NormalizedName,
) error {
	if len(positionalArgs) > 0 {
		cf.positionalArgs = make([]*flag.Flag, len(positionalArgs))
	}

	for propName := range parametersSchema.Properties {
		f := cf.addSchemaFlag(
			parametersSchema,
			propName,
			normalizeName("param."),
			normalizeName,
			slices.Contains(parametersSchema.Required, propName),
			false,
			slices.Contains(hiddenFlags, propName),
		)
		position := slices.Index(positionalArgs, propName)
		if position >= 0 {
			cf.positionalArgs[position] = f
		}
	}

	for i, f := range cf.positionalArgs {
		if f == nil {
			// these must not happen in practice, unless we did a mistake in our sdk (static, blueprint, openapi)
			return fmt.Errorf("programming error: positionalArgs[%d] %q is not an existing schema property", i, positionalArgs[i])
		}
	}

	return nil
}

func (cf *cmdFlags) addConfigsFlags(
	configsSchema *mgcSdk.Schema,
	normalizeName func(name string) flag.NormalizedName,
) {
	for propName := range configsSchema.Properties {
		_ = cf.addSchemaFlag(
			configsSchema,
			propName,
			normalizeName("config."),
			normalizeName,
			slices.Contains(configsSchema.Required, propName),
			true,
			false,
		)
	}
}

func getFlagFormattedExample(f *flag.Flag) (example string) {
	fv, ok := f.Value.(schema_flags.SchemaFlagValue)
	if !ok {
		return
	}

	desc := fv.Desc()
	example = getExampleFormattedValue(desc.Schema, desc.Container, desc.PropName)
	if example == "" {
		return
	}

	return fmt.Sprintf("--%s=%s", f.Name, example)
}
