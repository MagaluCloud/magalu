package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"magalu.cloud/core"
	mgcSdk "magalu.cloud/sdk"
)

const outputFlag = "output"
const defaultFormatter = "yaml"
const helpFormatter = "help"

type OutputFormatter interface {
	Format(value any, options string, isRaw bool) error
	Description() string
}

var outputFormatters = map[string]OutputFormatter{}

func getOutputFormats() []string {
	keys := make([]string, 0, len(outputFormatters))
	for k := range outputFormatters {
		keys = append(keys, k)
	}
	return keys
}

func addOutputFlag(cmd *cobra.Command) {
	cmd.Root().PersistentFlags().StringP(
		outputFlag,
		"o",
		"",
		fmt.Sprintf(
			`If the result is plain data types, it's the output format.
One of %s, use 'help' to know more details.
Otherwise it's the file name to save to, use '-' to write to stdout (default)`,
			strings.Join(getOutputFormats(), "|")))
}

func getOutputFlag(cmd *cobra.Command) string {
	return cmd.Root().PersistentFlags().Lookup(outputFlag).Value.String()
}

func setOutputFlag(cmd *cobra.Command, value string) {
	_ = cmd.Root().PersistentFlags().Lookup(outputFlag).Value.Set(value)
}

// TODO: Bind config to PFlag. Investigate how to make it work correctly
func getOutpuConfig(sdk *mgcSdk.Sdk) string {
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
	output := getOutputFlag(cmd)
	if output == "" {
		output = getOutpuConfig(sdk)
	}

	if output == "" {
		if outputOptions, ok := core.ResultAs[core.ResultWithDefaultOutputOptions](result); ok {
			return outputOptions.DefaultOutputOptions()
		}
	}

	return output
}
