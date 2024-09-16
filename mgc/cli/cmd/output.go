package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"magalu.cloud/core"
	mgcSdk "magalu.cloud/sdk"
)

const outputFlag = "output"
const helpFormatter = "help"
const defaultFormatter = "yaml"

type OutputFormatter interface {
	Format(value any, options string, isRaw bool) error
	Description() string
}

var outputFormatters = map[string]OutputFormatter{}

func addOutputFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().StringP(
		outputFlag,
		"o",
		"",
		`Change the output format. Use '--output=help' to know more details.`)
}

func getOutputFlag(cmd *cobra.Command) string {
	return cmd.Root().PersistentFlags().Lookup(outputFlag).Value.String()
}

func setOutputFlag(cmd *cobra.Command, value string) {
	_ = cmd.Root().PersistentFlags().Lookup(outputFlag).Value.Set(value)
}

// TODO: Bind config to PFlag. Investigate how to make it work correctly
func getOutputConfig(sdk *mgcSdk.Sdk) string {
	var defaultOutput string
	err := sdk.Config().Get("defaultOutput", &defaultOutput)
	if err != nil {
		return ""
	}
	return defaultOutput
}

func hasOutputFormatHelp(cmd *cobra.Command) bool {
	value := getOutputFlag(cmd)
	if value == helpFormatter {
		showFormatHelp()
		return true
	}
	return false
}

func parseOutputFormatter(output string) (name, options string) {
	parts := strings.SplitN(output, "=", 2)
	name = parts[0]
	if len(parts) == 2 {
		options = parts[1]
	}
	return name, options
}

// NOTE: use parseOutputFormatter() to get both name and options
func getOutputFormatter(name, options string) (formatter OutputFormatter, err error) {
	if name == "" {
		name = defaultFormatter
	}

	if formatter, ok := outputFormatters[name]; ok {
		return formatter, nil
	}
	return nil, fmt.Errorf("unknown formatter %q", name)
}

func getOutputFor(sdk *mgcSdk.Sdk, cmd *cobra.Command, result core.Result) string {
	var output string
	var defaultConfigOutput string
	var configFlag string

	if defaultConfigOutput = getOutputConfig(sdk); defaultConfigOutput != "" {
		output = defaultConfigOutput
	}

	if configFlag = getOutputFlag(cmd); configFlag != "" {
		output = configFlag
	}

	if outputOptions, ok := core.ResultAs[core.ResultWithDefaultOutputOptions](result); ok {
		outputFromSpec := outputOptions.DefaultOutputOptions()
		if strings.Contains(outputFromSpec, "default=") {
			outs := strings.Split(outputFromSpec, ";")
			for i, ot := range outs {
				if strings.HasPrefix(ot, "default=") {
					if defaultConfigOutput != "" {
						outs[i] = "default=" + defaultConfigOutput
					} else if configFlag != "" {
						outs[i] = "default=" + configFlag
					} else {
						outs[i] = ot
					}
					break
				}
			}
			output = strings.Join(outs, ";")
		}
	}

	if output == "" {
		return defaultFormatter
	}

	return output
}
