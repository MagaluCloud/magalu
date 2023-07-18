package cmd

import (
	"fmt"
	"net/http"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/profusion/magalu/libs/parser"
	"github.com/spf13/cobra"

	flag "github.com/spf13/pflag"
)

// -- BEGIN: OpenAPI finder

var openAPIFileNameRe = regexp.MustCompile("^(?P<name>[^.]+)(?:|[.]openapi)[.](?P<ext>json|yaml|yml)$")

// TODO: investigate channels to provide range iterOpenApis()...
// will it stop early when closed? how to do it?
func IterOpenApis(dir string, cb func(*parser.OpenAPIFileInfo) (run bool, err error)) (finished bool, err error) {
	// TODO: load from an index with description + version information

	items, err := os.ReadDir(dir)
	if err != nil {
		return false, fmt.Errorf("Unable to read OpenAPI files at %s: %w", dir, err)
	}

	for _, item := range items {
		info, err := item.Info()
		if err != nil {
			continue
		}

		if info.IsDir() {
			continue
		}

		matches := openAPIFileNameRe.FindStringSubmatch(item.Name())

		if len(matches) == 0 {
			continue
		}

		fileInfo := &parser.OpenAPIFileInfo{
			Name:      matches[1],
			Extension: matches[2],
			Path:      filepath.Join(dir, item.Name()),
			// TODO: load from an index with description + version information
			Description: "TODO: load description from index",
			Version:     "TODO: load version from index",
		}

		run, err := cb(fileInfo)
		if err != nil {
			return false, err
		}
		if !run {
			return false, nil
		}
	}

	return true, nil
}

// -- END: OpenAPI finder

// -- BEGIN: create Dynamic Argument Loaders --

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

// -- BEGIN: OpenApi Dynamic Loaders: Action, Resource, Module --

func createOpenApiResourceCmdLoader(
	openApi *parser.OpenAPIFileInfo,
) DynamicArgLoader {
	return createCommonDynamicArgLoader(
		func(cmd *cobra.Command) error {
			module, err := parser.LoadOpenAPI(openApi)
			if err != nil {
				return err
			}
			for _, tag := range module.Tags {
				_, _, err = AddOpenApiResourceCmd(cmd, module, tag)
				if err != nil {
					return err
				}
			}
			return nil
		},
		func(target string, cmd *cobra.Command) (*cobra.Command, DynamicArgLoader, error) {
			module, err := parser.LoadOpenAPI(openApi)
			if err != nil {
				return nil, nil, err
			}
			for _, tag := range module.Tags {
				if target == tag.Name {
					return AddOpenApiResourceCmd(cmd, module, tag)
				}
			}
			return nil, nil, fmt.Errorf("Resource not found: %s", target)
		},
	)
}

func createOpenApiActionCmdLoader(
	module *parser.OpenAPIModule,
	tag *openapi3.Tag,
) DynamicArgLoader {
	return createCommonDynamicArgLoader(
		func(cmd *cobra.Command) error {
			// TODO: fix ActionsByTag() to be an existing map or iterate and stop at the tag
			for _, action := range module.ActionsByTag()[tag] {
				_, _, err := AddOpenApiActionCmd(cmd, action)
				if err != nil {
					return err
				}
			}
			return nil
		},
		func(target string, cmd *cobra.Command) (*cobra.Command, DynamicArgLoader, error) {
			// TODO: fix ActionsByTag() to be an existing map or iterate and stop at the tag
			for _, action := range module.ActionsByTag()[tag] {
				if target == action.Name {
					return AddOpenApiActionCmd(cmd, action)
				}
			}
			return nil, nil, fmt.Errorf("Action not found: %s", target)
		},
	)
}

func createOpenApiCmdLoader(openApiDir string) DynamicArgLoader {
	return createCommonDynamicArgLoader(
		func(cmd *cobra.Command) error {
			_, err := IterOpenApis(openApiDir, func(fileInfo *parser.OpenAPIFileInfo) (run bool, err error) {
				_, _, e := AddOpenApiModuleCmd(cmd, fileInfo)
				return true, e
			})
			return err
		},
		func(target string, cmd *cobra.Command) (*cobra.Command, DynamicArgLoader, error) {
			var childCmd *cobra.Command = nil
			var loader DynamicArgLoader = nil
			var err error = nil
			finished, err := IterOpenApis(openApiDir, func(fileInfo *parser.OpenAPIFileInfo) (run bool, err error) {
				if fileInfo.Name == target {
					childCmd, loader, err = AddOpenApiModuleCmd(cmd, fileInfo)
					return false, err
				}
				return true, nil
			})
			if err != nil {
				return nil, nil, err
			}
			if !finished {
				// stopped early = found an item
				return childCmd, loader, nil
			}
			return nil, nil, fmt.Errorf("OpenAPI %s not found at %s", target, openApiDir)
		},
	)
}

// -- BEGIN: OpenApi Dynamic Loaders: Action, Resource, Module --

// -- BEGIN: AddOpenApi command structure --

func addOpenApiGroup(parentCmd *cobra.Command, kind string) {
	parentCmd.AddGroup(&cobra.Group{
		ID:    "openapi",
		Title: fmt.Sprintf("OpenAPI generated %s:", kind),
	})
}

func addParamFlagToCommand(cmd *cobra.Command, p *parser.OpenAPIParameter) {
	name := p.Name

	switch p.Type {
	case "boolean":
		cmd.Flags().Bool(name, false, p.Description)
	case "integer":
		cmd.Flags().Int(name, 0, p.Description)
	case "number":
		cmd.Flags().Float64(name, 0.0, p.Description)
	case "string":
		cmd.Flags().String(name, "", p.Description)
	case "array[boolean]":
		cmd.Flags().BoolSlice(name, []bool{}, p.Description)
	case "array[integer]":
		cmd.Flags().IntSlice(name, []int{}, p.Description)
	case "array[number]":
		log.Printf("number slice not implemented for param %s", p.Name)
	case "array[string]":
		cmd.Flags().StringSlice(name, []string{}, p.Description)
	}

	if p.Required {
		cmd.MarkFlagRequired(name)
	}
}

func loadParametersIntoCommand(action *parser.OpenAPIAction, cmd *cobra.Command) {
	for _, p := range action.PathParams {
		addParamFlagToCommand(cmd, p)
	}
	for _, p := range action.HeaderParams {
		addParamFlagToCommand(cmd, p)
	}
	for _, p := range action.RequestBodyParams {
		addParamFlagToCommand(cmd, p)
	}
}

func AddOpenApiActionCmd(
	parentCmd *cobra.Command,
	action *parser.OpenAPIAction,
) (*cobra.Command, DynamicArgLoader, error) {
	actionCmd := &cobra.Command{
		Use:     action.Name,
		Short:   action.Summary,
		Long:    action.Description,
		GroupID: "openapi",

		RunE: func(cmd *cobra.Command, args []string) error {
			println("\033[42mTODO: EXECUTE\033[0m", action.Name, "--", strings.Join(args, " "))

			// TODO: Actually execute action command here
			for _, param := range action.PathParams {
				if cmd.Flags().Changed(param.Name) {
					value := cmd.Flags().Lookup(param.Name).Value.String()
					action.PathName = strings.Replace(action.PathName, "{"+param.Name+"}", fmt.Sprintf("%v", value), 1)
				}
			}

			header := http.Header{}
			for _, param := range action.HeaderParams {
				if cmd.Flags().Changed(param.Name) {
					value := cmd.Flags().Lookup(param.Name)
					header.Add(param.Name, value.Value.String())
				}
			}

			return nil
		},
	}

	loadParametersIntoCommand(action, actionCmd)

	println("\033[1;36mACTION: ADDED CMD:\033[0m", actionCmd.Use)
	parentCmd.AddCommand(actionCmd)
	return actionCmd, nil, nil
}

func runHelpE(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func AddOpenApiResourceCmd(
	parentCmd *cobra.Command,
	module *parser.OpenAPIModule,
	tag *openapi3.Tag,
) (*cobra.Command, DynamicArgLoader, error) {
	resourceCmd := &cobra.Command{
		Use:     tag.Name,
		Short:   tag.Description,
		GroupID: "openapi",
		RunE:    runHelpE,
	}
	addOpenApiGroup(resourceCmd, "actions")

	loader := createOpenApiActionCmdLoader(module, tag)

	println("\033[1;35mRESOURCE: ADDED CMD:\033[0m", resourceCmd.Use)
	parentCmd.AddCommand(resourceCmd)
	return resourceCmd, loader, nil
}

func AddOpenApiModuleCmd(
	parentCmd *cobra.Command,
	openApi *parser.OpenAPIFileInfo,
) (*cobra.Command, DynamicArgLoader, error) {
	moduleCmd := &cobra.Command{
		Use:     openApi.Name,
		Short:   openApi.Description,
		Version: openApi.Version,
		GroupID: "openapi",
		RunE:    runHelpE,
	}
	addOpenApiGroup(moduleCmd, "resources")

	loader := createOpenApiResourceCmdLoader(openApi)

	println("\033[1;34mMODULE: ADDED CMD:\033[0m", moduleCmd.Use)
	parentCmd.AddCommand(moduleCmd)
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
	for true {
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

func Execute() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// TODO: this should be in viper + viper.SetDefault()
	openApiDir := filepath.Join(cwd, "openapis")

	rootCmd := &cobra.Command{
		Use:     "cloud",
		Version: "TODO",
		Short:   "CLI tool for OpenAPI integration",
		Long: `This CLI is a dynamic processor of OpenAPI files that
can generate a command line on-demand for Rest manipulation`,
		RunE: runHelpE,
	}
	addOpenApiGroup(rootCmd, "modules")
	rootCmd.AddGroup(&cobra.Group{
		ID:    "other",
		Title: "Other commands:",
	})
	rootCmd.SetHelpCommandGroupID("other")
	rootCmd.SetCompletionCommandGroupID("other")

	err = DynamicLoadCommand(rootCmd, os.Args[1:], createOpenApiCmdLoader(openApiDir))
	if err != nil {
		rootCmd.PrintErrln("Warning: loading dynamic arguments:", err)
	}

	return rootCmd.Execute()
}
