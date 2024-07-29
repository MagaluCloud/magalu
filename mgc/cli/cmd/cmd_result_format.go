package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"magalu.cloud/core"
	mgcSdk "magalu.cloud/sdk"
)

func formatResult(sdk *mgcSdk.Sdk, cmd *cobra.Command, result core.Result) error {
	output := getOutputFor(sdk, cmd, result)

	if resultWithReader, ok := core.ResultAs[core.ResultWithReader](result); ok {
		return handleResultWithReader(resultWithReader.Reader(), output, cmd)
	}

	if resultWithValue, ok := core.ResultAs[core.ResultWithValue](result); ok {
		return handleResultWithValue(resultWithValue, output, cmd)
	}

	return fmt.Errorf("unsupported result: %T %+v", result, result)
}

func handleResultWithReader(reader io.Reader, outFile string, cmd *cobra.Command) (err error) {
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	if outFile == "yaml" {
		yamlOutputFormatter := yamlOutputFormatter{}
		outputObject := make(map[string]interface{})
		outputBytes, err := io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("error reading from reader: %w", err)
		}
		err = yaml.Unmarshal(outputBytes, outputObject)
		if err != nil {
			return fmt.Errorf("error unmarshaling YAML: %w", err)
		}
		err = yamlOutputFormatter.Format(outputObject, "", getRawOutputFlag(cmd))
		if err != nil {
			return fmt.Errorf("error formatting YAML output: %w", err)
		}
		return nil
	}

	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		return fmt.Errorf("error copying reader to stdout: %w", err)
	}
	return nil
}

func handleResultWithValue(result core.ResultWithValue, output string, cmd *cobra.Command) (err error) {
	value := result.Value()
	if value == nil {
		return nil
	}

	err = result.ValidateSchema()
	if err != nil {
		logValidationErr(err)
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
	return formatter.Format(value, options, getRawOutputFlag(cmd))
}

func handleSimpleResultValue(value core.Value, output string) error {
	name, options := parseOutputFormatter(output)
	formatter, err := getOutputFormatter(name, options)

	if err != nil {
		return err
	}

	return formatter.Format(value, options, true)
}
