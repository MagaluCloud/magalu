package preview

import (
	"errors"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func childNames(t *testing.T, g core.Grouper) []string {
	t.Helper()
	names := []string{}
	_, err := g.VisitChildren(func(child core.Descriptor) (bool, error) {
		names = append(names, child.Name())
		return true, nil
	})
	require.NoError(t, err)
	return names
}

func child(t *testing.T, g core.Grouper, name string) core.Descriptor {
	t.Helper()
	c, err := g.GetChildByName(name)
	require.NoError(t, err, name)
	return c
}

func group(t *testing.T, g core.Grouper, name string) core.Grouper {
	t.Helper()
	c, ok := child(t, g, name).(core.Grouper)
	require.True(t, ok, "%q is not a group", name)
	return c
}

func visitErr(g core.Grouper) error {
	_, err := g.VisitChildren(func(core.Descriptor) (bool, error) { return true, nil })
	return err
}

func TestNewGroupIsInvisibleWithoutPackages(t *testing.T) {
	cases := map[string]*memStore{
		"empty store":  newMemStore(nil),
		"broken store": {err: errors.New("permission denied")},
	}

	for name, store := range cases {
		t.Run(name, func(t *testing.T) {
			root := NewGroup(Beta, store, &fakeValidator{}, &extensionPrefix)
			assert.Empty(t, childNames(t, root))
		})
	}
}

func TestNewGroupExposesPackagesUnderTheLevel(t *testing.T) {
	for _, level := range Levels {
		t.Run(level.Name, func(t *testing.T) {
			validator := &fakeValidator{}
			root := NewGroup(level, newMemStore(map[string]string{"foo": fooSpec}), validator, &extensionPrefix)
			assert.Equal(t, []string{level.Name}, childNames(t, root))

			levelGroup := group(t, root, level.Name)
			assert.Equal(t, level.Summary, levelGroup.Summary())
			assert.Contains(t, levelGroup.Description(), "mem-store")
			assert.Equal(t, []string{"foo"}, childNames(t, levelGroup))

			foo := group(t, levelGroup, "foo")
			assert.Equal(t, "Foo preview API.", foo.Description())

			things := group(t, foo, "things")
			_, ok := child(t, things, "list").(core.Executor)
			assert.True(t, ok, "%s foo things list should be an executor", level.Name)

			assert.Contains(t, validator.seen, level.Name+"/foo")
		})
	}
}

func TestNewGroupIsolatesBrokenPackages(t *testing.T) {
	store := newMemStore(map[string]string{
		"foo":     fooSpec,
		"broken":  brokenSpec,
		"refused": fooSpec,
	})
	validator := &fakeValidator{refused: map[string]string{"refused": "spec has no operations"}}

	beta := group(t, NewGroup(Beta, store, validator, &extensionPrefix), Beta.Name)
	assert.Equal(t, []string{"broken", "foo", "refused"}, childNames(t, beta))

	// the valid module keeps working
	assert.Equal(t, []string{"things"}, childNames(t, group(t, beta, "foo")))

	// the others fail only when used, saying which module and where
	err := visitErr(group(t, beta, "broken"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), `beta module "broken" (mem-store)`)

	err = visitErr(group(t, beta, "refused"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), `beta module "refused" (mem-store)`)
	assert.Contains(t, err.Error(), "spec has no operations")
}

func TestNewGroupDoesNotTouchOfficialGroups(t *testing.T) {
	officialFoo := core.NewSimpleGrouper(
		core.DescriptorSpec{Name: "foo", Description: "Official foo"},
		func() ([]core.Descriptor, error) { return nil, nil },
	)
	official := core.NewSimpleGrouper(
		core.DescriptorSpec{Name: "official", Description: "Official root"},
		func() ([]core.Grouper, error) { return []core.Grouper{officialFoo}, nil },
	)

	newRoot := func(beta, alpha Store) core.Grouper {
		return core.NewMergeGroup(
			core.DescriptorSpec{Name: "products", Description: "root"},
			func() []core.Grouper {
				return []core.Grouper{
					official,
					NewGroup(Beta, beta, &fakeValidator{}, &extensionPrefix),
					NewGroup(Alpha, alpha, &fakeValidator{}, &extensionPrefix),
				}
			},
		)
	}

	t.Run("without packages the tree is the official one", func(t *testing.T) {
		root := newRoot(newMemStore(nil), newMemStore(nil))
		assert.Equal(t, []string{"foo"}, childNames(t, root))
		assert.Same(t, officialFoo, child(t, root, "foo"))
	})

	t.Run("a preview package with an official name does not merge into it", func(t *testing.T) {
		specs := map[string]string{"foo": fooSpec}
		root := newRoot(newMemStore(specs), newMemStore(specs))
		assert.Equal(t, []string{Alpha.Name, Beta.Name, "foo"}, childNames(t, root))
		assert.Same(t, officialFoo, child(t, root, "foo"))
		assert.Empty(t, childNames(t, officialFoo))
		assert.Equal(t, []string{"foo"}, childNames(t, group(t, root, Beta.Name)))
		assert.Equal(t, []string{"foo"}, childNames(t, group(t, root, Alpha.Name)))
	})

	t.Run("each level shows only its own packages", func(t *testing.T) {
		root := newRoot(newMemStore(nil), newMemStore(map[string]string{"foo": fooSpec}))
		assert.Equal(t, []string{Alpha.Name, "foo"}, childNames(t, root))
	})
}

func TestIsLevelName(t *testing.T) {
	assert.True(t, IsLevelName("beta"))
	assert.True(t, IsLevelName("alpha"))
	assert.False(t, IsLevelName("iam"))
}
