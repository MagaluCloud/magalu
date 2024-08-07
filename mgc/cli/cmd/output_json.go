package cmd

import (
	"encoding/json"
	"os"

	"github.com/mattn/go-colorable"
	jsonColor "github.com/neilotoole/jsoncolor"
)

type jsonOutputFormatter struct{}

func (*jsonOutputFormatter) Format(value any, options string, isRaw bool) error {
	if isRaw {
		enc := json.NewEncoder(os.Stdout)
		if options == "compact" {
			enc.SetIndent("", "")
		} else {
			enc.SetIndent("", " ")
		}
		enc.SetEscapeHTML(false)
		return enc.Encode(value)
	}
	out := colorable.NewColorable(os.Stdout)
	enc := jsonColor.NewEncoder(out)

	if options == "compact" {
		enc.SetIndent("", "")
	} else {
		enc.SetIndent("", " ")
	}

	clrs := jsonColor.DefaultColors()
	clrs.Bool = jsonColor.Color("\x1b[95m")
	clrs.Number = jsonColor.Color("\x1b[95m")
	clrs.Key = jsonColor.Color("\x1b[96m")
	clrs.String = jsonColor.Color("\x1b[92m")

	enc.SetColors(clrs)
	enc.SetEscapeHTML(false)

	return enc.Encode(value)
}

func (*jsonOutputFormatter) Description() string {
	return `Format as JSON.` +
		` Use "json=compact" to use the compact encoding without spaces and indentation.`
}

func init() {
	outputFormatters["json"] = &jsonOutputFormatter{}
}
