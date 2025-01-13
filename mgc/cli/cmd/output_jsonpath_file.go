package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type jsonpathFileOutputFormatter struct{}

func (*jsonpathFileOutputFormatter) Format(value any, options string, isRaw bool) error {
	filename := options
	path, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("jsonpath-file output formatter: %w", err)
	}
	target, err := utils.GetJsonPath(string(path), value)
	if err != nil {
		return fmt.Errorf("jsonpath-file output formatter: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", " ")
	return enc.Encode(target)
}

func (*jsonpathFileOutputFormatter) Description() string {
	return `Same as jsonpath, but reads the expression from the given file: "jsonpath-file=path-to-file-with-jsonpath-expression".`
}

func init() {
	outputFormatters["jsonpath-file"] = &jsonpathFileOutputFormatter{}
}
