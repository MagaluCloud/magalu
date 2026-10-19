package preview

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
	"github.com/stretchr/testify/assert"
)

func command(name string, flags ...string) core.Executor {
	props := map[string]*core.Schema{}
	for _, flag := range flags {
		props[flag] = mgcSchemaPkg.NewStringSchema()
	}
	return core.NewSimpleExecutor(core.ExecutorSpec{
		DescriptorSpec:   core.DescriptorSpec{Name: name, Description: name},
		ParametersSchema: mgcSchemaPkg.NewObjectSchema(props, nil),
		ConfigsSchema:    mgcSchemaPkg.NewObjectSchema(map[string]*core.Schema{}, nil),
		ResultSchema:     mgcSchemaPkg.NewNullSchema(),
		Execute: func(core.Executor, context.Context, core.Parameters, core.Configs) (core.Result, error) {
			return nil, nil
		},
	})
}

func staticGroup(name, groupID string, children ...core.Descriptor) core.Grouper {
	return core.NewStaticGroup(
		core.DescriptorSpec{Name: name, Description: name, GroupID: groupID},
		func() []core.Descriptor { return children },
	)
}

func propertyDiff(preview, official core.Executor) (added []string) {
	for name := range preview.ParametersSchema().Properties {
		if official.ParametersSchema().Properties[name] == nil {
			added = append(added, name)
		}
	}
	slices.Sort(added)
	return added
}

// Like the alpha IAM: --project on every command, --organization on roles,
// one command and one group the official product lacks.
func iamTrees() (root core.Grouper, alpha core.Grouper) {
	official := staticGroup("iam", productsGroupID,
		staticGroup("roles", "", command("list", "role-name"), command("create", "name")),
		staticGroup("members", "", command("list"), command("delete", "id")),
	)
	module := staticGroup("iam", "",
		staticGroup("roles", "", command("list", "role-name", "organization", "project"), command("create", "name", "project")),
		staticGroup("members", "", command("list", "project"), command("delete", "id", "project"), command("create", "email", "project")),
		staticGroup("projects", "", command("list", "project")),
	)
	alpha = staticGroup("alpha", "", module, staticGroup("connect", "", staticGroup("networks", "", command("list"))))
	root = staticGroup("root", "", official, alpha)
	return root, alpha
}

func chainOf(root core.Grouper, names ...string) []core.Descriptor {
	var chain []core.Descriptor
	var node core.Descriptor = root
	for _, name := range names {
		child, err := node.(core.Grouper).GetChildByName(name)
		if err != nil {
			panic(err)
		}
		chain = append(chain, child)
		node = child
	}
	return chain
}

func TestMarksTellWhatTheOfficialProductLacks(t *testing.T) {
	root, _ := iamTrees()
	marks := NewMarks(func() core.Grouper { return root }, propertyDiff)
	short := func(names ...string) string {
		return shown(marks, chainOf(root, names...))
	}

	assert.Equal(t, "iam (preview commands)", short("alpha", "iam"))
	assert.Equal(t, "roles (preview flags)", short("alpha", "iam", "roles"))
	assert.Equal(t, "list (preview flags)", short("alpha", "iam", "roles", "list"))
	assert.Equal(t, "create", short("alpha", "iam", "roles", "create"), "--project is everywhere, so it leads nowhere")
	assert.Equal(t, "members (preview commands)", short("alpha", "iam", "members"))
	assert.Equal(t, "[PREVIEW] create", short("alpha", "iam", "members", "create"))
	assert.Equal(t, "[PREVIEW] projects", short("alpha", "iam", "projects"))
	assert.Equal(t, "list", short("alpha", "iam", "projects", "list"), "inside something new, nothing is marked")

	assert.Equal(t, "[PREVIEW] connect", short("alpha", "connect"), "a product the official tree lacks")
	assert.Equal(t, "networks", short("alpha", "connect", "networks"))

	assert.Equal(t, []string{"organization", "project"}, marks.NewFlags(chainOf(root, "alpha", "iam", "roles", "list")))
	assert.Empty(t, marks.NewFlags(chainOf(root, "alpha", "iam", "members", "create")), "a new command is new as a whole")
	assert.Equal(t, []string{"project"}, marks.Everywhere(chainOf(root, "alpha", "iam")))
	assert.Empty(t, marks.Everywhere(chainOf(root, "alpha", "iam", "roles")), "told once, in the module")
}

func TestMarksNothingWhenTheComparisonFails(t *testing.T) {
	failing := core.NewSimpleGrouper(
		core.DescriptorSpec{Name: "iam", Description: "iam"},
		func() ([]core.Descriptor, error) { return nil, errors.New("spec does not load") },
	)
	root := staticGroup("root", "", staticGroup("iam", productsGroupID, command("list")), staticGroup("alpha", "", failing))
	marks := NewMarks(func() core.Grouper { return root }, propertyDiff)

	assert.Equal(t, "iam", shown(marks, chainOf(root, "alpha", "iam")), "the official product exists, so the module is not new")
}

func TestMarksLeaveTheOfficialTreeAlone(t *testing.T) {
	root, _ := iamTrees()
	marks := NewMarks(func() core.Grouper { return root }, propertyDiff)

	chain := chainOf(root, "iam", "roles", "list")
	assert.Equal(t, "list", shown(marks, chain))
	assert.Empty(t, marks.NewFlags(chain))
	assert.Equal(t, "alpha", shown(marks, chainOf(root, "alpha")), "the level itself")
}

// The name as the help shows it, without styles.
func shown(marks *Marks, chain []core.Descriptor) string {
	name := chain[len(chain)-1].Name()
	isNew, trail := marks.Mark(chain)
	switch {
	case isNew:
		return NewMark + " " + name
	case trail != "":
		return name + " " + trail
	}
	return name
}
