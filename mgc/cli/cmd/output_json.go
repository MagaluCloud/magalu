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
		return enc.Encode(value)
	}
	out := colorable.NewColorable(os.Stdout) // needed for Windows
	enc := jsonColor.NewEncoder(out)

	if options == "compact" {
		enc.SetIndent("", "")
	} else {
		enc.SetIndent("", " ")
	}

	clrs := jsonColor.DefaultColors()
    // Update colors based on provided settings
    clrs.Bool = jsonColor.Color("\x1b[95m") // FgHiMagenta for boolean
    clrs.Number = jsonColor.Color("\x1b[95m") // FgHiMagenta for number
    clrs.Key = jsonColor.Color("\x1b[96m") // FgHiCyan for map keys
    // Assuming Anchor and Alias are not directly supported by jsoncolor, so we'll skip these
    clrs.String = jsonColor.Color("\x1b[92m") // FgHiGreen for strings
    // Reset color is not needed here as jsoncolor handles it


	enc.SetColors(clrs)

	return enc.Encode(value)
}

func (*jsonOutputFormatter) Description() string {
	return `Format as JSON.` +
		` Use "json=compact" to use the compact encoding without spaces and indentation.`
}

func init() {
	outputFormatters["json"] = &jsonOutputFormatter{}
}
