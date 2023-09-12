package openapi

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
	"magalu.cloud/core"
)

type IndexModule struct {
	Name        string
	Path        string
	Version     string
	Description string
}

type IndexFile struct {
	Version string
	Modules []IndexModule
}

const indexFile = "index.yaml"
const indexVersion = "1.0.0"

// Source -> Module -> Resource -> Operation

// -- ROOT: Source

type Source struct {
	Loader          Loader
	ExtensionPrefix *string
	modules         []*Module
	byName          map[string]*Module
}

// BEGIN: Descriptor interface:

func (o *Source) Name() string {
	return "OpenApis"
}

func (o *Source) Version() string {
	return ""
}

func (o *Source) Description() string {
	return fmt.Sprintf("OpenApis loaded using %v", o.Loader)
}

// END: Descriptor interface

// BEGIN: Grouper interface:

func (o *Source) getModules() (modules []*Module, byName map[string]*Module, err error) {
	if len(o.modules) > 0 {
		return o.modules, o.byName, nil
	}

	data, err := o.Loader.Load(indexFile)
	if err != nil {
		return nil, nil, err
	}

	var index IndexFile
	err = yaml.Unmarshal(data, &index)
	if err != nil {
		return nil, nil, err
	}
	if index.Version != indexVersion {
		return nil, nil, fmt.Errorf("Unsupported %q version %q, expected %q", indexFile, index.Version, indexVersion)
	}

	o.modules = make([]*Module, len(index.Modules))
	o.byName = make(map[string]*Module, len(index.Modules))

	for i, item := range index.Modules {
		module := &Module{
			name:            item.Name,
			path:            item.Path,
			version:         item.Version,
			description:     item.Description,
			extensionPrefix: o.ExtensionPrefix,
			loader:          o.Loader,
			logger:          logger().Named(item.Name),
		}
		o.modules[i] = module
		o.byName[module.Name()] = module
	}

	slices.SortFunc(o.modules, func(a, b *Module) int {
		return strings.Compare(a.Name(), b.Name())
	})

	return o.modules, o.byName, nil
}

func (o *Source) VisitChildren(visitor core.DescriptorVisitor) (finished bool, err error) {
	modules, _, err := o.getModules()
	if err != nil {
		return false, err
	}

	for _, module := range modules {
		run, err := visitor(module)
		if err != nil {
			return false, err
		}
		if !run {
			return false, nil
		}
	}

	return true, nil
}

func (o *Source) GetChildByName(name string) (child core.Descriptor, err error) {
	_, byName, err := o.getModules()
	if err != nil {
		return nil, err
	}

	if module, ok := byName[name]; ok {
		return module, nil
	} else {
		return nil, fmt.Errorf("Module not found: %s", name)
	}
}

var _ core.Grouper = (*Source)(nil)

// END: Grouper interface
