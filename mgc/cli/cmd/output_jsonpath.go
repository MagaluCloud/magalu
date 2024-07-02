package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"magalu.cloud/core/utils"
)

type jsonpathOutputFormatter struct{}

func (*jsonpathOutputFormatter) Format(value any, options string, isRaw bool) error {
	path := options
	target, err := utils.GetJsonPath(path, value)
	if err != nil {
		return fmt.Errorf("jsonpath output formatter: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	return enc.Encode(target)
}

func (*jsonpathOutputFormatter) Description() string {
	return `Use JSON Path expression to select elements: "jsonpath=jsonpath-expression".` +
		` For more complex specifications, see "jsonpath-file".`
}

func init() {
	outputFormatters["jsonpath"] = &jsonpathOutputFormatter{}
}
