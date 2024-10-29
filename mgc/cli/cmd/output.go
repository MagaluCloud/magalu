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
	name = ""
	newOptions := []string{}
	for _, x := range strings.Split(output, ";") {
		if value, _ := strings.CutPrefix(x, "default="); value != "" {
			if _, ok := outputFormatters[value]; ok {
				name = value
				continue
			}
		}

		if strings.HasPrefix(x, "allowfields=") || strings.HasPrefix(x, "remove=") {
			newOptions = append(newOptions, x)
		}
	}
	return name, strings.Join(newOptions, ";")
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
	outputFromSpec := ""

	formatFromSpec := []string{}
	if outputOptions, ok := core.ResultAs[core.ResultWithDefaultOutputOptions](result); ok {
		for _, x := range strings.Split(outputOptions.DefaultOutputOptions(), ";") {
			if strings.HasPrefix(x, "allowfields=") || strings.HasPrefix(x, "remove=") {
				formatFromSpec = append(formatFromSpec, x)
				continue
			}
			if fromSpec, ok := strings.CutPrefix(x, "default="); ok {
				outputFromSpec = fromSpec
				continue
			}
		}
	}

	fmtFromSpec := strings.Join(formatFromSpec, ";")
	if fmtFromSpec != "" {
		fmtFromSpec = ";" + fmtFromSpec
	}

	outputFromConfig := getOutputConfig(sdk)
	outputFromFlag := getOutputFlag(cmd)

	if outputFromFlag != "" {
		return outputFromFlag + fmtFromSpec
	}

	if outputFromConfig != "" {
		return outputFromConfig + fmtFromSpec
	}

	return outputFromSpec + fmtFromSpec
}
