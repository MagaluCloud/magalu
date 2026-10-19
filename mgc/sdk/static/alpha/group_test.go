package alpha

import (
	"context"
	"testing"

	"github.com/MagaluCloud/magalu/mgc/core"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
)

const fooSpec = `openapi: 3.0.3
info:
  title: Foo
  version: '1.0'
paths:
  /v1/things:
    get:
      tags: [things]
      responses:
        '200':
          description: OK
`

type memStore map[string]string

func (s memStore) Names() ([]string, error) {
	names := []string{}
	for name := range s {
		names = append(names, name)
	}
	return names, nil
}

func (s memStore) Read(name string) (preview.Package, error) {
	return preview.Package{Name: name, Spec: []byte(s[name])}, nil
}

func (s memStore) String() string { return "mem-store" }

type acceptAll struct{}

func (acceptAll) Validate(preview.Level, preview.Package) error { return nil }

func newCommand() core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{Name: "status", Description: "A command that ships with the CLI"},
		func(context.Context) (core.Value, error) { return nil, nil },
	)
}

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

func TestAlphaIsHiddenUntilSpecsAreInstalled(t *testing.T) {
	prefix := "x-mgc"
	newRoot := func(installed memStore) core.Grouper {
		return core.NewMergeGroup(
			core.DescriptorSpec{Name: "products"},
			func() []core.Grouper {
				return []core.Grouper{
					preview.NewGroup(preview.Alpha, installed, acceptAll{}, &prefix),
					NewRoot(newCommand()),
				}
			},
		)
	}

	alpha := alphaOf(t, newRoot(memStore{}))
	if !alpha.IsInternal() {
		t.Error("without installed specs alpha should be hidden")
	}
	if _, err := alpha.GetChildByName("status"); err != nil {
		t.Errorf("hidden alpha should still run its own commands: %v", err)
	}

	alpha = alphaOf(t, newRoot(memStore{"foo": fooSpec}))
	if alpha.IsInternal() {
		t.Error("with installed specs alpha should be visible")
	}
	for _, name := range []string{"foo", "status"} {
		if _, err := alpha.GetChildByName(name); err != nil {
			t.Errorf("alpha should have %q: %v", name, err)
		}
	}
}
