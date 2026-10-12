package preview

import (
	"context"
	"fmt"

	"github.com/MagaluCloud/magalu/mgc/core"
	mgcLoggerPkg "github.com/MagaluCloud/magalu/mgc/core/logger"
	"github.com/MagaluCloud/magalu/mgc/sdk/openapi"
)

var logger = mgcLoggerPkg.NewLazy[Level]()

const GroupID = "preview"

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
			GroupID: GroupID,
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
					child = checkedModule(&moduleGroup{module, level, store})
				}
				children = append(children, child)
				return true, nil
			})
			return children, err
		},
	)
}

func checkedModule(module *moduleGroup) core.Descriptor {
	err := loadAll(module)
	if err == nil {
		return module
	}

	logger().Warnw("ignoring broken preview module", "level", module.level.Name, "module", module.Name(), "error", err)
	return newBrokenModule(module.Name(), err)
}

func Check(level Level, pkg Package, validator Validator, extensionPrefix *string) error {
	store := packageStore{pkg}
	source := openapi.NewSource(&storeLoader{level: level, store: store, validator: validator}, extensionPrefix)

	child, err := source.GetChildByName(pkg.Name)
	if err != nil {
		return err
	}
	module, ok := child.(core.Grouper)
	if !ok {
		return fmt.Errorf("%s module %q is not a group of commands", level.Name, pkg.Name)
	}
	return loadAll(&moduleGroup{module, level, store})
}

func loadAll(module core.Grouper) error {
	_, err := core.VisitAllExecutors(module, nil, true, func(core.Executor, []string) (bool, error) {
		return true, nil
	})
	return err
}

type packageStore struct {
	pkg Package
}

func (s packageStore) Names() ([]string, error) {
	return []string{s.pkg.Name}, nil
}

func (s packageStore) Read(name string) (Package, error) {
	if name != s.pkg.Name {
		return Package{}, fmt.Errorf("package %q not found", name)
	}
	return s.pkg, nil
}

func (s packageStore) String() string {
	return "download"
}

func newBrokenModule(name string, err error) core.Executor {
	return core.NewStaticExecuteSimple(
		core.DescriptorSpec{
			Name:        name,
			Summary:     "Failed to load: run it to see why",
			Description: err.Error(),
		},
		func(context.Context) (core.Value, error) {
			return nil, err
		},
	)
}

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
