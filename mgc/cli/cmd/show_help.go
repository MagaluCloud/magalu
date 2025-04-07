package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
	flag "github.com/spf13/pflag"

	"github.com/MagaluCloud/magalu/mgc/cli/cmd/schema_flags"
	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

func showFormatHelp() {
	writer := table.NewWriter()
	writer.AppendHeader(table.Row{"Formatter", "Description"})
	maxLen := 0
	for k, f := range outputFormatters {
		if maxLen < len(k) {
			maxLen = len(k)
		}
		writer.AppendRow(table.Row{k, f.Description()})
	}

	termColumns := getTermColumns()
	tablePadding := 4 // 2 columns x 2 spaces per column
	writer.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft, VAlign: text.VAlignTop, WidthMax: maxLen},
		{Number: 2, Align: text.AlignLeft, VAlign: text.VAlignTop, WidthMax: termColumns - maxLen - tablePadding},
	})
	style := table.StyleDefault
	style.Options = table.OptionsNoBordersAndSeparators
	style.Options.SeparateHeader = true

	writer.SetStyle(style)
	writer.SortBy([]table.SortBy{{Name: "Formatter", Mode: table.Asc}})

	fmt.Println("For plain data types, the following values are accepted:")

	fmt.Println(writer.Render())

	fmt.Println("\nFor streams, use the file name to save to or '-' to write to stdout (default).")
}

func showHelpForError(cmd *cobra.Command, args []string, err error) error {
	switch {
	case err == schema_flags.ErrWantHelp:
		return nil

	case errors.As(err, new(core.UsageError)):
		// we can't call UsageString() on the root, we need to find the actual leaf command that failed:
		subCmd, _, _ := cmd.Find(args)
		cmd.PrintErrln(subCmd.UsageString())
	default:
		break
	}

	return err
}

func textIndent(s, firstPrefix, siblingPrefix string) string {
	return firstPrefix + strings.Join(strings.Split(s, "\n"), "\n"+siblingPrefix)
}

func textReflow(s, firstPrefix, siblingPrefix string, columns int) string {
	if columns > 0 {
		s = text.WrapText(s, columns)
	}

	if siblingPrefix != "" {
		s = textIndent(s, firstPrefix, siblingPrefix)
	}

	return s
}

func cleanJSONSchemaCOW(cowSchema *mgcSchemaPkg.COWSchema) {
	_ = cowSchema.SetExample(nil)
	_ = cowSchema.SetExtensions(nil)

	if cowSchema.Items() != nil {
		cleanJSONSchemaRefCOW(cowSchema.ItemsCOW())
	}

	if len(cowSchema.Properties()) > 0 {
		cowSchema.PropertiesCOW().ForEachCOW(func(_ string, cow *mgcSchemaPkg.COWSchemaRef) (run bool) {
			cleanJSONSchemaRefCOW(cow)
			return true
		})
	}

	if len(cowSchema.AllOf()) > 0 {
		cowSchema.AllOfCOW().ForEachCOW(func(_ int, cow *mgcSchemaPkg.COWSchemaRef) (run bool) {
			cleanJSONSchemaRefCOW(cow)
			return true
		})
	}

	if len(cowSchema.AnyOf()) > 0 {
		cowSchema.AnyOfCOW().ForEachCOW(func(_ int, cow *mgcSchemaPkg.COWSchemaRef) (run bool) {
			cleanJSONSchemaRefCOW(cow)
			return true
		})
	}

	if len(cowSchema.OneOf()) > 0 {
		cowSchema.OneOfCOW().ForEachCOW(func(_ int, cow *mgcSchemaPkg.COWSchemaRef) (run bool) {
			cleanJSONSchemaRefCOW(cow)
			return true
		})
	}
}

func cleanJSONSchemaRefCOW(cowSchemaRef *mgcSchemaPkg.COWSchemaRef) {
	if cowSchemaRef.Value() == nil {
		return
	}
	cleanJSONSchemaCOW(cowSchemaRef.ValueCOW())
}

// Removes examples, extensions and whatever else may not be relevant to the user
func getCleanJSONSchema(schema *core.Schema) *core.Schema {
	cowSchema := mgcSchemaPkg.NewCOWSchema(schema)
	cowSchema.SetDescription("") // the root description was already displayed
	cleanJSONSchemaCOW(cowSchema)
	return cowSchema.Peek()
}

func getExample(schema, container *core.Schema, propName string) (example any) {
	example = schema.Example
	if example != nil {
		return
	}

	if schema.Type != nil && schema.Type.Includes(openapi3.TypeArray) &&
		schema.Items != nil && schema.Items.Value != nil && schema.Items.Value.Example != nil {
		return []any{schema.Items.Value.Example}
	}

	if container.Example == nil {
		return
	}

	if containerExample, ok := container.Example.(map[string]any); ok {
		example = containerExample[propName]
		if example != nil {
			return
		}
	}

	return
}

func getExampleFormattedValue(schema, container *core.Schema, propName string) (value string) {
	example := getExample(schema, container, propName)
	if example == nil {
		return
	}

	data, err := json.Marshal(example)
	if err != nil {
		return
	}

	if schema.Type != nil {
		switch {
		case schema.Type.Includes("integer"), schema.Type.Includes("number"), schema.Type.Includes("boolean"):
			return string(data)

		case schema.Type.Includes("string"):
			value = string(data)
			if value == schema_flags.ValueHelpIsRequired {
				value = schema_flags.ValueVerbatimStringPrefix + value
			} else if strings.HasPrefix(value, schema_flags.ValueVerbatimStringPrefix) || strings.Contains(value, "$") {
				value = fmt.Sprintf("'%s'", data) // keep quotes and wrap in single, so shell doesn't replace variables
			}
			return
		default:
			return fmt.Sprintf("'%s'", data)
		}
	}
	return fmt.Sprintf("'%s'", data)
}

func getFlagHelpUsageLine(f *flag.Flag) (s string) {
	s = "--" + f.Name
	t := f.Value.Type()

	if f.NoOptDefVal != "" {
		s += fmt.Sprintf("[=%s] (if no explicit value is given, uses: %s)", t, f.NoOptDefVal)
	} else {
		s += "=" + t
	}

	return
}

func getFlagDescription(f *flag.Flag) (description string) {
	if fv, ok := f.Value.(schema_flags.SchemaFlagValue); ok {
		description = fv.Desc().Description()
	} else if f.Usage != "" {
		description = f.Usage
	}

	return description
}

func endsWithPunctuation(text string) bool {
	if text == "" {
		return false
	}
	return unicode.IsPunct(rune(text[len(text)-1]))
}

func addPunctuation(text, punctuation string) string {
	if endsWithPunctuation(text) {
		return text
	}
	return text + punctuation
}

func forcePunctuation(text, punctuation string) string {
	if endsWithPunctuation(text) {
		return text[:len(text)-1] + punctuation
	}
	return text + punctuation
}

func getConstraintsFormatted(constraints *schema_flags.HumanReadableConstraints, indentPrefix string, isListItem bool) (text string) {
	text = indentPrefix
	childIdentPrefix := indentPrefix + "  "
	if isListItem {
		text += "- "
		childIdentPrefix += "  " // must be same length as line above
	}

	if constraints.Description == "" {
		if constraints.Message != "" {
			text += addPunctuation(constraints.Message, ".")
		}
	} else {
		text += forcePunctuation(textIndent(constraints.Description, "", indentPrefix+"  | "), "")
		if constraints.Message != "" {
			text += fmt.Sprintf(" (%s)", constraints.Message)
		}
		text += "."
	}

	if len(constraints.Children) == 0 {
		text += "\n"
		return
	}

	if strings.Contains(constraints.Description, "\n") {
		text += "\n" + indentPrefix + "  "
	} else if constraints.Description != "" || constraints.Message != "" {
		text = addPunctuation(text, ".")
	}

	if constraints.ChildrenMessage != "" {
		if endsWithPunctuation(text) {
			text += " "
		}
		text += forcePunctuation(constraints.ChildrenMessage, ":")
	}

	if len(constraints.Children) == 1 {
		child := constraints.Children[0]
		if endsWithPunctuation(text) {
			text += " "
		}
		text += strings.Trim(getConstraintsFormatted(child, childIdentPrefix, false), "\t\n\r ")
		text += "\n"
		return
	}

	text = forcePunctuation(text, ":")

	text += "\n"
	for _, c := range constraints.Children {
		text += getConstraintsFormatted(c, childIdentPrefix, true)
	}

	return
}

func showFlagHelp(f *flag.Flag) {
	const (
		indentPrefix = "    "
	)
	columns := getTermColumns()
	reflowColumns := columns - len(indentPrefix)
	var output string

	reflow := func(s string) string {
		return textReflow(s, indentPrefix, indentPrefix, reflowColumns)
	}

	addSection := func(title string, rest ...string) {
		if output != "" {
			output += "\n"
		}

		output += fmt.Sprintf("%s:", title)
		if len(rest) > 0 {
			for _, s := range rest {
				output += " " + s
			}
		}
		output += "\n"
	}
	addSectionBodyText := func(body string) {
		output += reflow(body) + "\n"
	}
	addSectionBodyRaw := func(body string) {
		output += indentPrefix + body + "\n"
	}

	addSection("Flag usage")
	addSectionBodyRaw(getFlagHelpUsageLine(f))

	if fv, ok := f.Value.(schema_flags.SchemaFlagValue); ok {
		desc := fv.Desc()

		if constraints := desc.HumanReadableConstraints(); constraints != nil {
			addSection("Description")
			output += getConstraintsFormatted(constraints, indentPrefix, false)
		}

		cleanSchema := getCleanJSONSchema(desc.Schema)
		if data, err := json.MarshalIndent(cleanSchema, indentPrefix, "  "); err == nil {
			schema := strings.Trim(string(data), "\t\n\r ")
			if schema != "{}" {
				addSection("JSON Schema")
				addSectionBodyRaw(schema)
			}
		}

		if example := getExampleFormattedValue(desc.Schema, desc.Container, desc.PropName); example != "" {
			addSection("Example")
			addSectionBodyRaw(fmt.Sprintf("--%s=%s", f.Name, example))
		}
	} else {
		if description := getFlagDescription(f); description != "" {
			addSection("Description")
			addSectionBodyText(description)
		}

		if f.DefValue != "" {
			addSection("Default value", f.DefValue)
		}
	}

	fmt.Println(output)
}
