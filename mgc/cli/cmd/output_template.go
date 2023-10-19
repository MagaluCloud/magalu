package cmd

import (
	"fmt"
	"os"

	"magalu.cloud/core/utils"
)

type templateOutputFormatter struct{}

func (*templateOutputFormatter) Format(value any, options string) error {
	text := options
	tmpl, err := utils.NewTemplate(text)
	if err != nil {
		return fmt.Errorf("template output formatter: %w", err)
	}

	return tmpl.Execute(os.Stdout, value)
}

func (*templateOutputFormatter) Description() string {
	return `Format using https://pkg.go.dev/text/template. Use "template=your-template-here."` +
		` For more complex specifications, see "template-file".`
}

func init() {
	outputFormatters["template"] = &templateOutputFormatter{}
}
