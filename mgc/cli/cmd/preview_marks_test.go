package cmd

import (
	"testing"

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
