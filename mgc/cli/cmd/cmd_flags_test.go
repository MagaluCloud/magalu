package cmd

import (
	"fmt"
	"testing"

	"slices"

	"magalu.cloud/cli/cmd/schema_flags"
	mgcSchemaPkg "magalu.cloud/core/schema"

	"github.com/spf13/cobra"
)

func checkError(t *testing.T, prefix string, expected, got error) {
	if expected == nil {
		if got != nil {
			t.Errorf("%s: unexpected error %q", prefix, got.Error())
		}
	} else {
		if got == nil {
			t.Errorf("%s: expected error %q, got nothing", prefix, expected.Error())
		} else if expected.Error() != got.Error() {
			t.Errorf("%s: expected error %q, got %q", prefix, expected.Error(), got.Error())
		}
	}
}

func checkExpectedString(t *testing.T, prefix, expected, got string) {
	if expected == got {
		return
	}
	t.Errorf("%s: expected %q, got %q", prefix, expected, got)
}

func checkExpectedArray(t *testing.T, prefix string, expected, got []string) {
	if slices.Equal(expected, got) {
		return
	}
	t.Errorf("%s: expected %v, got %v", prefix, expected, got)
}

func Test_cmdFlags_positionalArray(t *testing.T) {
	type parameters struct {
		Files   []mgcSchemaPkg.FilePath `json:"files" jsonschema:"description=files to use"`
		Other   string                  `json:"other"`
		Another string                  `json:"another"`
		Array   []string                `json:"array" jsonschema:"description=another array"`
	}

	schema, err := mgcSchemaPkg.SchemaFromType[parameters]()
	checkError(t, "SchemaFromType", nil, err)

	type testCase struct {
		name           string
		positionalArgs []string
		hiddenFlags    []string
		input          []string
		files          []string
		array          []string
		use            string
		other, another string
		completions    []string
		err            error
	}
	tests := []testCase{
		{
			name:           "alone/empty",
			positionalArgs: []string{"files"},
			input:          nil,
			files:          nil,
			use:            "testing [files...]",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use)"),
		},
		{
			name:           "alone/single",
			positionalArgs: []string{"files"},
			input:          []string{"a"},
			files:          []string{"a"},
			use:            "testing [files...]",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use)"),
		},
		{
			name:           "alone/multiple",
			positionalArgs: []string{"files"},
			input:          []string{"a", "b"},
			files:          []string{"a", "b"},
			use:            "testing [files...]",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use)"),
		},
		{
			name:           "alone/mixed",
			positionalArgs: []string{"files"},
			input:          []string{"a", "b,c", `["d", "e"]`},
			files:          []string{"a", "b", "c", "d", "e"},
			use:            "testing [files...]",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use)"),
		},

		// leading:
		{
			name:           "leading/empty",
			positionalArgs: []string{"files", "other"},
			input:          nil,
			files:          nil,
			use:            "testing [files...] [other]",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use) or other"),
		},
		{
			name:           "leading/otherValue",
			positionalArgs: []string{"files", "other"},
			input:          []string{"otherValue"},
			files:          nil,
			use:            "testing [files...] [other]",
			other:          "otherValue",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use) or other"),
		},
		{
			name:           "leading/single",
			positionalArgs: []string{"files", "other"},
			input:          []string{"a", "otherValue"},
			files:          []string{"a"},
			use:            "testing [files...] [other]",
			other:          "otherValue",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use) or other"),
		},
		{
			name:           "leading/multiple",
			positionalArgs: []string{"files", "other"},
			input:          []string{"a", "b", "otherValue"},
			files:          []string{"a", "b"},
			use:            "testing [files...] [other]",
			other:          "otherValue",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use) or other"),
		},
		{
			name:           "leading/mixed",
			positionalArgs: []string{"files", "other"},
			input:          []string{"a", "b,c", `["d", "e"]`, "otherValue"},
			files:          []string{"a", "b", "c", "d", "e"},
			use:            "testing [files...] [other]",
			other:          "otherValue",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use) or other"),
		},

		// trailing:
		{
			name:           "trailing/empty",
			positionalArgs: []string{"other", "files"},
			input:          nil,
			files:          nil,
			use:            "testing [other] [files...]",
			completions:    cobra.AppendActiveHelp(nil, "other"),
		},
		{
			name:           "trailing/otherValue",
			positionalArgs: []string{"other", "files"},
			input:          []string{"otherValue"},
			files:          nil,
			use:            "testing [other] [files...]",
			other:          "otherValue",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use)"),
		},
		{
			name:           "trailing/single",
			positionalArgs: []string{"other", "files"},
			input:          []string{"otherValue", "a"},
			files:          []string{"a"},
			use:            "testing [other] [files...]",
			other:          "otherValue",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use)"),
		},
		{
			name:           "trailing/multiple",
			positionalArgs: []string{"other", "files"},
			input:          []string{"otherValue", "a", "b"},
			files:          []string{"a", "b"},
			use:            "testing [other] [files...]",
			other:          "otherValue",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use)"),
		},
		{
			name:           "trailing/mixed",
			positionalArgs: []string{"other", "files"},
			input:          []string{"otherValue", "a", "b,c", `["d", "e"]`},
			files:          []string{"a", "b", "c", "d", "e"},
			use:            "testing [other] [files...]",
			other:          "otherValue",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use)"),
		},

		// between:
		{
			name:           "between/empty",
			positionalArgs: []string{"other", "files", "another"},
			input:          nil,
			files:          nil,
			use:            "testing [other] [files...] [another]",
			completions:    cobra.AppendActiveHelp(nil, "other"),
		},
		{
			name:           "between/otherValue",
			positionalArgs: []string{"other", "files", "another"},
			input:          []string{"otherValue"},
			files:          nil,
			use:            "testing [other] [files...] [another]",
			other:          "otherValue",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use) or another"),
		},
		{
			name:           "between/anotherValue",
			positionalArgs: []string{"other", "files", "another"},
			input:          []string{"otherValue", "anotherValue"},
			files:          nil,
			use:            "testing [other] [files...] [another]",
			other:          "otherValue",
			another:        "anotherValue",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use) or another"),
		},
		{
			name:           "between/single",
			positionalArgs: []string{"other", "files", "another"},
			input:          []string{"otherValue", "a", "anotherValue"},
			files:          []string{"a"},
			use:            "testing [other] [files...] [another]",
			other:          "otherValue",
			another:        "anotherValue",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use) or another"),
		},
		{
			name:           "between/multiple",
			positionalArgs: []string{"other", "files", "another"},
			input:          []string{"otherValue", "a", "b", "anotherValue"},
			files:          []string{"a", "b"},
			use:            "testing [other] [files...] [another]",
			other:          "otherValue",
			another:        "anotherValue",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use) or another"),
		},
		{
			name:           "between/mixed",
			positionalArgs: []string{"other", "files", "another"},
			input:          []string{"otherValue", "a", "b,c", `["d", "e"]`, "anotherValue"},
			files:          []string{"a", "b", "c", "d", "e"},
			use:            "testing [other] [files...] [another]",
			other:          "otherValue",
			another:        "anotherValue",
			completions:    cobra.AppendActiveHelp(nil, "The following arguments are accepted: multiple files (files to use) or another"),
		},

		// another array
		{
			name:           "2-arrays/empty",
			positionalArgs: []string{"files", "array"},
			input:          nil,
			files:          nil,
			array:          nil,
			use:            "testing [files] [array]",
			completions:    cobra.AppendActiveHelp(nil, "files (files to use)"),
		},
		{
			name:           "2-arrays/single",
			positionalArgs: []string{"files", "array"},
			input:          []string{"a"},
			files:          []string{"a"},
			array:          nil,
			use:            "testing [files] [array]",
			completions:    cobra.AppendActiveHelp(nil, "array (another array)"),
		},
		{
			name:           "2-arrays/both",
			positionalArgs: []string{"files", "array"},
			input:          []string{"a", "arrayValue"},
			files:          []string{"a"},
			array:          []string{"arrayValue"},
			use:            "testing [files] [array]",
			completions:    cobra.AppendActiveHelp(nil, "This command does not take any more arguments"),
		},
		{
			name:           "2-arrays/extra(error)",
			positionalArgs: []string{"files", "array"},
			input:          []string{"a", "arrayValue", "extra(error)"},
			files:          nil,
			array:          nil,
			use:            "testing [files] [array]",
			completions:    cobra.AppendActiveHelp(nil, "This command does not take any more arguments"),
			err:            fmt.Errorf("this command receives at most 2 positional arguments, 3 given"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			flags, err := newCmdFlags(&cobra.Command{}, schema, &mgcSchemaPkg.Schema{}, tc.positionalArgs, tc.hiddenFlags)
			checkError(t, "newCmdFlags", nil, err)
			cmd := &cobra.Command{
				Use:               buildUse("testing", flags.positionalArgsNames()),
				Args:              flags.positionalArgsFunction,
				ValidArgsFunction: flags.validateArgs,
			}

			checkExpectedString(t, "cmd.Use", tc.use, cmd.Use)

			err = cmd.Args(cmd, tc.input)
			checkError(t, "cmd.Args", tc.err, err)

			v, err := flags.knownFlags["files"].Value.(schema_flags.SchemaFlagValue).Parse()
			checkError(t, "parse files", nil, err)

			var files []string
			for _, item := range v.([]any) {
				files = append(files, item.(string))
			}

			v, err = flags.knownFlags["array"].Value.(schema_flags.SchemaFlagValue).Parse()
			checkError(t, "parse array", nil, err)

			var array []string
			for _, item := range v.([]any) {
				array = append(array, item.(string))
			}

			checkExpectedArray(t, "files", tc.files, files)
			checkExpectedArray(t, "array", tc.array, array)
			checkExpectedString(t, "other", tc.other, flags.knownFlags["other"].Value.String())
			checkExpectedString(t, "another", tc.another, flags.knownFlags["another"].Value.String())

			comp, _ := cmd.ValidArgsFunction(cmd, tc.input, "")
			checkExpectedArray(t, "completions", tc.completions, comp)
		})
	}
}
