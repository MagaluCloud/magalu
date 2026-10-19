package cmd

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func alphaTree() (root, list *cobra.Command) {
	root = &cobra.Command{Use: "mgc"}
	alpha := &cobra.Command{Use: "alpha"}
	iam := &cobra.Command{Use: "iam"}
	list = &cobra.Command{Use: "list", Run: func(*cobra.Command, []string) {}}
	dbaas := &cobra.Command{Use: "dbaas"}
	iam.AddCommand(list)
	alpha.AddCommand(iam)
	root.AddCommand(alpha, dbaas)
	return root, list
}

func TestIsAlphaCommand(t *testing.T) {
	root, list := alphaTree()
	assert.True(t, isAlphaCommand(list))
	assert.False(t, isAlphaCommand(root))

	dbaas, _, _ := root.Find([]string{"dbaas"})
	assert.False(t, isAlphaCommand(dbaas))
}

func TestUnknownAlphaModulePointsToPreviewList(t *testing.T) {
	root, _ := alphaTree()
	unknown := unknownCommandError{"foo"}

	err := showHelpForError(root, []string{"alpha", "foo"}, unknown)
	assert.EqualError(t, err, `no command with name "foo": if your tenant was invited to this preview, see 'mgc preview list'`)
	assert.True(t, errors.As(err, new(unknownCommandError)))

	err = showHelpForError(root, []string{"dbaas", "foo"}, unknown)
	assert.Equal(t, unknown, err, "only under 'mgc alpha'")
	err = showHelpForError(root, []string{"alpha", "iam", "foo"}, unknown)
	assert.Equal(t, unknown, err, "a module installed has its own commands")
}
