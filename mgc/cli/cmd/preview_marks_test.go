package cmd

import (
	"context"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestMarkFlagUsages(t *testing.T) {
	usages := "" +
		"  -h, --help                          help for list\n" +
		"      --organization                  When True, it sets the scope\n" +
		"      --cli.list-links enum[=table]   List all available links\n" +
		"      --project string                Project to scope this command to.\n" +
		"                                      Prefer 'mgc alpha project set'\n"

	assert.Equal(t, ""+
		"  -h, --help                          help for list\n"+
		"      --organization                  [PREVIEW] When True, it sets the scope\n"+
		"      --cli.list-links enum[=table]   List all available links\n"+
		"      --project string                [PREVIEW] Project to scope this command to.\n"+
		"                                      Prefer 'mgc alpha project set'\n",
		markFlagUsages(usages, []string{"organization", "project"}))
	assert.Equal(t, usages, markFlagUsages(usages, nil))

	saved := color.NoColor
	color.NoColor = false
	defer func() { color.NoColor = saved }()
	assert.Contains(t, markFlagUsages(usages, []string{"organization"}), "\x1b[1m[PREVIEW]\x1b[22m When True", "bold on a terminal")
}

func TestUnderPreviewLevel(t *testing.T) {
	root := &cobra.Command{Use: "mgc"}
	alpha, iam, list := &cobra.Command{Use: "alpha"}, &cobra.Command{Use: "iam"}, &cobra.Command{Use: "list"}
	official := &cobra.Command{Use: "iam"}
	iam.AddCommand(list)
	alpha.AddCommand(iam)
	root.AddCommand(alpha, official)

	assert.True(t, underPreviewLevel(list))
	assert.True(t, underPreviewLevel(alpha))
	assert.False(t, underPreviewLevel(official))
	assert.False(t, underPreviewLevel(root))
}

func TestPreviewMarksStarTheCommandsWithSomethingNew(t *testing.T) {
	saved := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = saved }()

	command := func(name string) core.Descriptor {
		return core.NewStaticExecuteSimple(
			core.DescriptorSpec{Name: name, Description: name},
			func(context.Context) (core.Value, error) { return nil, nil },
		)
	}
	group := func(name, groupID string, children ...core.Descriptor) core.Grouper {
		return core.NewStaticGroup(
			core.DescriptorSpec{Name: name, Description: name, GroupID: groupID},
			func() []core.Descriptor { return children },
		)
	}
	tree := group("root", "",
		group("iam", "catalog", group("roles", "", command("list"))),
		group("dns", "catalog", command("list")),
		group("alpha", "",
			group("iam", "", group("roles", "", command("list"), command("create"))),
			group("dns", "", command("list")),
			group("connect", "", command("list")),
		),
	)

	run := func(*cobra.Command, []string) {}
	rootCmd := &cobra.Command{Use: "mgc"}
	alphaCmd := &cobra.Command{Use: "alpha", Short: "Early preview commands", Run: run}
	for _, name := range []string{"iam", "dns", "connect"} {
		alphaCmd.AddCommand(&cobra.Command{Use: name, Short: name, Run: run})
	}
	rootCmd.AddCommand(alphaCmd, &cobra.Command{Use: "iam", Short: "iam", Run: run})
	installPreviewMarks(rootCmd, func() core.Grouper { return tree })

	usage := alphaCmd.UsageString()
	assert.Contains(t, usage, "\n *connect     [PREVIEW] connect\n", "a module the official tree lacks")
	assert.Contains(t, usage, "\n  dns         dns\n", "the same as the official one")
	assert.Contains(t, usage, "\n *iam         iam (preview commands)\n", "something new inside")

	usage = rootCmd.UsageString()
	assert.Contains(t, usage, "\n  alpha       Early preview commands\n")
	assert.Contains(t, usage, "\n  iam         iam\n", "the official tree is never marked")
}
