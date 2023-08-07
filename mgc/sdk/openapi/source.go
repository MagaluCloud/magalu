package openapi

import (
	"fmt"

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

func (o *Source) VisitChildren(visitor core.DescriptorVisitor) (finished bool, err error) {
	data, err := o.Loader.Load(indexFile)
	if err != nil {
		return false, err
	}

	var index IndexFile
	err = yaml.Unmarshal(data, &index)
	if err != nil {
		return false, err
	}
	if index.Version != indexVersion {
		return false, fmt.Errorf("Unsupported %q version %q, expected %q", indexFile, index.Version, indexVersion)
	}

	for _, item := range index.Modules {
		module := &Module{
			name:            item.Name,
			path:            item.Path,
			version:         item.Version,
			description:     item.Description,
			extensionPrefix: o.ExtensionPrefix,
			loader:          o.Loader,
		}

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
	// TODO: write O(1) version that doesn't list
	var found core.Descriptor
	finished, err := o.VisitChildren(func(child core.Descriptor) (run bool, err error) {
		if child.Name() == name {
			found = child
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		return nil, err
	}

	if finished {
		return nil, fmt.Errorf("Module not found: %s", name)
	}

	return found, err
}

var _ core.Grouper = (*Source)(nil)

// END: Grouper interface
