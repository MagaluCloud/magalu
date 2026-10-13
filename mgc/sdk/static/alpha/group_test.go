package alpha

import (
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
)

func alphaOf(t *testing.T, root core.Grouper) core.Grouper {
	t.Helper()

	child, err := root.GetChildByName(Name)
	if err != nil {
		t.Fatalf("root has no %q: %v", Name, err)
	}
	group, ok := child.(core.Grouper)
	if !ok {
		t.Fatalf("%q is not a group", Name)
	}
	return group
}

func TestAlphaIsHiddenButHasProject(t *testing.T) {
	alpha := alphaOf(t, GetRoot())

	if !alpha.IsInternal() {
		t.Error("alpha should be hidden")
	}
	if _, err := alpha.GetChildByName("project"); err != nil {
		t.Errorf("alpha should have project: %v", err)
	}
}

func TestVisibleAlphaMergedFirstWins(t *testing.T) {
	visible := core.NewStaticGroup(
		core.DescriptorSpec{Name: "installed root"},
		func() []core.Descriptor {
			return []core.Descriptor{
				core.NewStaticGroup(core.DescriptorSpec{Name: Name}, func() []core.Descriptor { return nil }),
			}
		},
	)

	merged := core.NewMergeGroup(
		core.DescriptorSpec{Name: "products"},
		func() []core.Grouper { return []core.Grouper{visible, GetRoot()} },
	)

	alpha := alphaOf(t, merged)
	if alpha.IsInternal() {
		t.Error("alpha merged after a visible one should be visible")
	}
	if _, err := alpha.GetChildByName("project"); err != nil {
		t.Errorf("merged alpha should still have project: %v", err)
	}
}
