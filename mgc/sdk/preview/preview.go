// Package preview exposes preview commands under "mgc beta" and "mgc alpha".
// They are built at runtime from the compiled OpenAPI specs kept on this
// machine, never from the specs embedded in the binary, so the official
// commands stay untouched.
//
// Where the specs come from and how they are checked are ports (Store,
// Validator); the adapters live in subpackages and sdk.go picks them.
package preview

import (
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcLoggerPkg "github.com/MagaluCloud/magalu/mgc/core/logger"
	"github.com/MagaluCloud/magalu/mgc/sdk/openapi"
)

var logger = mgcLoggerPkg.NewLazy[Level]()

// NewGroup returns a root to be merged with the official ones. It has a single
// child named after the level when the store has packages, and no children
// otherwise, so the command stays hidden. It never fails: a broken store must
// not break the CLI.
func NewGroup(level Level, store Store, validator Validator, extensionPrefix *string) core.Grouper {
	return core.NewSimpleGrouper(
		core.DescriptorSpec{Name: level.Name + " root", Description: level.Summary},
		func() ([]core.Descriptor, error) {
			names, err := store.Names()
			if err != nil {
				logger().Debugw("ignoring preview store", "level", level.Name, "store", store.String(), "error", err)
				return nil, nil
			}
			if len(names) == 0 {
				return nil, nil
			}
			return []core.Descriptor{newLevelGroup(level, store, validator, extensionPrefix)}, nil
		},
	)
}

func newLevelGroup(level Level, store Store, validator Validator, extensionPrefix *string) core.Grouper {
	loader := &storeLoader{level: level, store: store, validator: validator}
	source := openapi.NewSource(loader, extensionPrefix)
	return core.NewSimpleGrouper(
		core.DescriptorSpec{
			Name:    level.Name,
			Summary: level.Summary,
			Description: fmt.Sprintf(
				"%s, loaded from %s. Commands may change or stop working. ",
				level.Summary,
				store,
			),
		},
		func() (children []core.Descriptor, err error) {
			_, err = source.VisitChildren(func(child core.Descriptor) (bool, error) {
				if module, ok := child.(core.Grouper); ok {
					child = &moduleGroup{module, level, store}
				}
				children = append(children, child)
				return true, nil
			})
			return children, err
		},
	)
}

// moduleGroup says which preview module failed and where it lives, since
// the openapi errors do not.
type moduleGroup struct {
	core.Grouper
	level Level
	store Store
}

func (m *moduleGroup) wrap(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s module %q (%s): %w", m.level.Name, m.Name(), m.store, err)
}

func (m *moduleGroup) VisitChildren(visitor core.DescriptorVisitor) (bool, error) {
	finished, err := m.Grouper.VisitChildren(visitor)
	return finished, m.wrap(err)
}

func (m *moduleGroup) GetChildByName(name string) (core.Descriptor, error) {
	child, err := m.Grouper.GetChildByName(name)
	return child, m.wrap(err)
}
