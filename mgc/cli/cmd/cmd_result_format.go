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

func keepProperties(data any, properties []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, prop := range properties {
		parts := strings.Split(strings.TrimPrefix(prop, "$."), ".")
		keepPropertyRecursive(data, parts, result)
	}
	return result
}

func keepPropertyRecursive(data interface{}, parts []string, result map[string]interface{}) {
	if len(parts) == 0 || data == nil {
		return
	}

	currentPart := parts[0]
	remainingParts := parts[1:]

	if arrayPart, isArray := strings.CutSuffix(currentPart, "[*]"); isArray {
		if obj, ok := data.(map[string]interface{}); ok {
			if arr, exists := obj[arrayPart]; exists {
				if arrayData, ok := arr.([]interface{}); ok {
					newArray := make([]interface{}, len(arrayData))
					for i, item := range arrayData {
						newItem := make(map[string]interface{})
						if mapItem, ok := item.(map[string]interface{}); ok {
							keepPropertyRecursive(mapItem, remainingParts, newItem)
						}
						newArray[i] = newItem
					}
					setOrMergeValue(result, arrayPart, newArray)
				}
			}
		}
	} else {
		if obj, ok := data.(map[string]interface{}); ok {
			if value, exists := obj[currentPart]; exists {
				if len(remainingParts) == 0 {
					setOrMergeValue(result, currentPart, value)
				} else {
					subResult := make(map[string]interface{})
					keepPropertyRecursive(value, remainingParts, subResult)
					setOrMergeValue(result, currentPart, subResult)
				}
			}
		}
	}
}

func setOrMergeValue(result map[string]interface{}, key string, value interface{}) {
	if existing, exists := result[key]; exists {
		if existingMap, ok := existing.(map[string]interface{}); ok {
			if newMap, ok := value.(map[string]interface{}); ok {
				for k, v := range newMap {
					existingMap[k] = v
				}
				return
			}
		}
		if existingArray, ok := existing.([]interface{}); ok {
			if newArray, ok := value.([]interface{}); ok {
				if len(existingArray) == len(newArray) {
					for i, newItem := range newArray {
						if newMap, ok := newItem.(map[string]interface{}); ok {
							if existingMap, ok := existingArray[i].(map[string]interface{}); ok {
								for k, v := range newMap {
									existingMap[k] = v
								}
							} else {
								existingArray[i] = newMap
							}
						}
					}
				} else {
					result[key] = newArray
				}
				return
			}
		}
	}
	result[key] = value
}

func handleResultWithValue(result core.ResultWithValue, output string, cmd *cobra.Command) (err error) {
	err = result.ValidateSchema()
	if err != nil {
		logValidationErr(err)
	}

	outputs := strings.Split(output, ";")
	output = ""
	var remove string
	for _, ot := range outputs {
		if strings.HasPrefix(ot, "remove=") {
			remove = strings.Split(ot, "=")[1]
		} else {
			output = ot
		}
	}
	var fieldsToRemove []string
	if remove != "" {
		fieldsToRemove = strings.Split(remove, ",")
	}

	var allowed string
	for _, ot := range outputs {
		if strings.HasPrefix(ot, "allowfields=") {
			allowed = strings.Split(ot, "=")[1]
		} else {
			output = ot
		}
	}
	var allowedFields []string
	if allowed != "" {
		allowedFields = strings.Split(allowed, ",")
		for i, x := range allowedFields {
			allowedFields[i] = strings.Split(x, ":")[1]
		}
	}

	for _, ot := range outputs {
		if strings.HasPrefix(ot, "default=") {
			output = strings.Split(ot, "=")[1]
		}
	}

	value := result.Value()
	if value == nil {
		return nil
	}

	for _, path := range fieldsToRemove {
		value = removeProperty(value, path)
	}
	if len(allowedFields) > 0 {
		value = keepProperties(value, allowedFields)
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
