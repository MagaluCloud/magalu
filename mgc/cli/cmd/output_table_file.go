package cmd

import (
	"fmt"
	"os"

	"github.com/invopop/yaml"
)

type tableFileOutputFormatter struct{}

func (*tableFileOutputFormatter) Format(value any, options string, isRaw bool) error {
	filename := options
	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("table-file output formatter: %w", err)
	}

	tableOpts := &tableOptions{}
	err = yaml.Unmarshal(file, tableOpts)

	if err != nil {
		return err
	}

	return formatTableWithOptions(value, tableOpts)
}

func (*tableFileOutputFormatter) Description() string {
	return `Format as table using https://github.com/jedib0t/go-pretty/#table using "table-file=path-to-YAML-file-with-table-configuration".` +
		` The YAML file should contain "columns" array with elements specifying "name" and "jsonpath",` +
		` optionally may include column configuration. Other top-level configuration include "format", "style" and "rowLength".`
}

func init() {
	outputFormatters["table-file"] = &tableFileOutputFormatter{}
}
