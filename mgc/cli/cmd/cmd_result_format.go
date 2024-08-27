package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

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

func removeProperty(data any, path string) any {
	parts := strings.Split(strings.TrimPrefix(path, "$."), ".")
	return removePropertyRecursive(data, parts)
}

func removePropertyRecursive(data any, parts []string) interface{} {
	if len(parts) == 0 {
		return data
	}

	currentPart := parts[0]
	remainingParts := parts[1:]

	// check if part is array
	if currentPart, ok := strings.CutSuffix(currentPart, "[*]"); ok {
		if obj, ok := data.(map[string]interface{}); ok {
			obj[currentPart] = removePropertyRecursive(obj[currentPart], parts[1:])
		}
	} else {
		if len(remainingParts) == 0 {
			if obj, ok := data.([]interface{}); ok {
				for _, item := range obj {
					if it, ok := item.(map[string]interface{}); ok {
						delete(it, currentPart)
					}

				}
			}

		}

	}

	return data
}

func handleResultWithValue(result core.ResultWithValue, output string, cmd *cobra.Command) (err error) {
	outputas := strings.Split(output, ",")
	remove := ""
	for _, ot := range outputas {
		if strings.HasPrefix(ot, "remove=") {
			remove = strings.Split(ot, "=")[1]
		} else {
			output = ot
		}
	}

	fieldsToRemove := strings.Split(remove, "|")

	value := result.Value()
	if value == nil {
		return nil
	}

	for _, path := range fieldsToRemove {
		value = removeProperty(value, path)
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
