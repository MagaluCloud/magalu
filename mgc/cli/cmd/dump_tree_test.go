package cmd

import (
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDumpTreeLeavesAlphaOut(t *testing.T) {
	newGroup := func(name string) core.Grouper {
		return core.NewSimpleGrouper(
			core.DescriptorSpec{Name: name, Description: name},
			func() ([]core.Descriptor, error) { return nil, nil },
		)
	}
	root := core.NewSimpleGrouper(
		core.DescriptorSpec{Name: "products", Description: "root"},
		func() ([]core.Descriptor, error) {
			return []core.Descriptor{newGroup(preview.Beta.Name), newGroup("iam"), newGroup(preview.Alpha.Name)}, nil
		},
	)

	tree, err := collectAllChildren(withoutAlpha(root))
	require.NoError(t, err)

	names := []string{}
	for _, child := range tree["children"].([]map[string]any) {
		names = append(names, child["name"].(string))
	}
	assert.Equal(t, []string{preview.Beta.Name, "iam"}, names)
}
