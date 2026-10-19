package cmd

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/cli/cmd/schema_flags"
	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSdk "github.com/MagaluCloud/magalu/mgc/sdk"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// The flag part of a usage line, up to where its description starts.
var flagUsageRe = regexp.MustCompile(`^\s+(?:-\S, )?--([^\s=\[]+)\S*(?: \S+)?\s{2,}`)

type previewMarks struct {
	sdk   *mgcSdk.Sdk
	marks *preview.Marks
	// Descriptors from the root's child down to each command, resolved once.
	chains map[*cobra.Command][]core.Descriptor
}

// installPreviewMarks marks in the help what the preview modules bring that
// the official products have not. Outside the preview levels the help stays
// the same, byte for byte. It goes after the colors, which rewrite the same
// template.
func installPreviewMarks(rootCmd *cobra.Command, sdk *mgcSdk.Sdk) {
	m := &previewMarks{
		sdk:    sdk,
		marks:  preview.NewMarks(sdk.Group, newPreviewFlagDiff(rootCmd)),
		chains: map[*cobra.Command][]core.Descriptor{},
	}
	cobra.AddTemplateFunc("previewShort", m.short)
	cobra.AddTemplateFunc("previewFlags", m.flags)
	cobra.AddTemplateFunc("previewNote", m.note)

	rootCmd.SetUsageTemplate(strings.NewReplacer(
		".Short", "(previewShort .)",
		".LocalFlags.FlagUsages", "(previewFlags . .LocalFlags.FlagUsages)",
		".InheritedFlags.FlagUsages", "(previewFlags . .InheritedFlags.FlagUsages)",
	).Replace(rootCmd.UsageTemplate()) + "{{previewNote .}}")
}

func (m *previewMarks) short(cmd *cobra.Command) string {
	if !underPreviewLevel(cmd) {
		return cmd.Short
	}
	isNew, trail := m.marks.Mark(m.chain(cmd))
	switch {
	case isNew:
		return previewMark() + cmd.Short
	case trail != "":
		return cmd.Short + " " + color.New(color.Italic).Sprint(trail)
	}
	return cmd.Short
}

func (m *previewMarks) flags(cmd *cobra.Command, usages string) string {
	if !underPreviewLevel(cmd) {
		return usages
	}
	return markFlagUsages(usages, m.marks.NewFlags(m.chain(cmd)))
}

func (m *previewMarks) note(cmd *cobra.Command) string {
	if !underPreviewLevel(cmd) {
		return ""
	}
	flags := m.marks.Everywhere(m.chain(cmd))
	if len(flags) == 0 {
		return ""
	}
	return fmt.Sprintf("\nEvery command here also takes --%s %s\n", strings.Join(flags, ", --"), color.New(color.Bold).Sprint(preview.NewMark))
}

// The mark in bold, and the trail in italic, stand out only in a terminal
// with colors; elsewhere, and with NO_COLOR, they are plain text.
func previewMark() string {
	return color.New(color.Bold).Sprint(preview.NewMark) + " "
}

func markFlagUsages(usages string, flags []string) string {
	if len(flags) == 0 {
		return usages
	}
	lines := strings.Split(usages, "\n")
	for i, line := range lines {
		if match := flagUsageRe.FindStringSubmatchIndex(line); match != nil && slices.Contains(flags, line[match[2]:match[3]]) {
			lines[i] = line[:match[1]] + previewMark() + line[match[1]:]
		}
	}
	return strings.Join(lines, "\n")
}

func underPreviewLevel(cmd *cobra.Command) bool {
	top := cmd
	for top.HasParent() && top.Parent().HasParent() {
		top = top.Parent()
	}
	return top.HasParent() && preview.IsLevel(top.Name())
}

// Resolved the way the command tree is loaded, since descriptor and command
// names differ. Commands that are not descriptors, like the helpers, get none.
func (m *previewMarks) chain(cmd *cobra.Command) []core.Descriptor {
	if chain, ok := m.chains[cmd]; ok {
		return chain
	}

	var chain []core.Descriptor
	if parent := cmd.Parent(); parent != nil {
		var group core.Grouper
		var prefix []core.Descriptor
		if parent.HasParent() {
			prefix = m.chain(parent)
			if len(prefix) > 0 {
				group, _ = prefix[len(prefix)-1].(core.Grouper)
			}
		} else {
			group = m.sdk.Group()
		}

		if group != nil {
			if desc, err := findChildByNameOrAliases(group, cmd.Name()); err == nil {
				chain = append(slices.Clip(prefix), desc)
			}
		}
	}

	m.chains[cmd] = chain
	return chain
}

// Compares the flags as the commands build them, so the names are the ones the
// help shows, conflicts and child flags included. The root they hang from is
// built at the first comparison, once every flag of the CLI is there.
func newPreviewFlagDiff(rootCmd *cobra.Command) preview.FlagDiff {
	var parent *cobra.Command
	return func(previewExec, official core.Executor) []string {
		if parent == nil {
			parent = withoutConfigFlags(rootCmd)
		}
		previewFlags, err := visibleSchemaFlags(parent, previewExec)
		if err != nil {
			return nil
		}
		officialFlags, err := visibleSchemaFlags(parent, official)
		if err != nil {
			return nil
		}
		return slices.DeleteFunc(previewFlags, func(name string) bool { return slices.Contains(officialFlags, name) })
	}
}

// A root with only the flags the CLI defines. The config flags of each command
// join the root once it loads, and would rename the flags compared after it.
func withoutConfigFlags(rootCmd *cobra.Command) *cobra.Command {
	parent := &cobra.Command{Use: rootCmd.Use}
	parent.SetGlobalNormalizationFunc(rootCmd.GlobalNormalizationFunc())
	add := func(f *flag.Flag) {
		if _, isSchema := f.Value.(schema_flags.SchemaFlagValue); !isSchema && parent.Flags().Lookup(f.Name) == nil {
			parent.Flags().AddFlag(f)
		}
	}
	rootCmd.PersistentFlags().VisitAll(add)
	rootCmd.Flags().VisitAll(add)
	return parent
}

func visibleSchemaFlags(rootCmd *cobra.Command, exec core.Executor) ([]string, error) {
	flags, err := newExecutorCmdFlags(rootCmd, exec)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, f := range slices.Concat(flags.schemaFlags, flags.childFlags) {
		if !f.Hidden {
			names = append(names, f.Name)
		}
	}
	slices.Sort(names)
	return names, nil
}
