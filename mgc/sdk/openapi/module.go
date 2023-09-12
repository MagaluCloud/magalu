package openapi

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"magalu.cloud/core"

	"github.com/getkin/kin-openapi/openapi3"
)

// Source -> Module -> Resource -> Operation

// Module

type Module struct {
	name            string
	path            string
	version         string
	description     string
	extensionPrefix *string
	doc             *openapi3.T
	loader          Loader
	logger          *zap.SugaredLogger
	resources       []*Resource
	byName          map[string]*Resource
	execResolver    executorResolver
}

// BEGIN: Descriptor interface:

func (m *Module) Name() string {
	return m.name
}

func (m *Module) Version() string {
	return m.version
}

func (m *Module) Description() string {
	return m.description
}

// END: Descriptor interface

// BEGIN: Grouper interface:

func (m *Module) getDoc() (*openapi3.T, error) {
	if m.doc == nil {
		ctx := context.Background()
		mData, err := m.loader.Load(m.path)
		if err != nil {
			return nil, err
		}

		oapiLoader := openapi3.Loader{Context: ctx, IsExternalRefsAllowed: false}
		doc, err := oapiLoader.LoadFromData(mData)
		if err != nil {
			return nil, err
		}
		m.doc = doc
	}

	return m.doc, nil
}

func (m *Module) getResources() (resources []*Resource, byName map[string]*Resource, err error) {
	if len(m.resources) > 0 {
		return m.resources, m.byName, nil
	}

	doc, err := m.getDoc()
	if err != nil {
		return nil, nil, err
	}

	m.resources = make([]*Resource, 0, len(doc.Tags))
	m.byName = make(map[string]*Resource, len(doc.Tags))

	for _, tag := range doc.Tags {
		if getHiddenExtension(m.extensionPrefix, tag.Extensions) {
			continue
		}

		resource := &Resource{
			tag:             tag,
			doc:             doc,
			extensionPrefix: m.extensionPrefix,
			servers:         doc.Servers,
			logger:          m.logger.Named(tag.Name),
			execResolver:    &m.execResolver,
		}
		m.resources = append(m.resources, resource)
		m.byName[resource.Name()] = resource
	}

	slices.SortFunc(m.resources, func(a, b *Resource) int {
		return strings.Compare(a.Name(), b.Name())
	})

	return m.resources, m.byName, nil
}

func (m *Module) VisitChildren(visitor core.DescriptorVisitor) (finished bool, err error) {
	resources, _, err := m.getResources()
	if err != nil {
		return false, err
	}

	for _, module := range resources {
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

func (m *Module) GetChildByName(name string) (child core.Descriptor, err error) {
	_, byName, err := m.getResources()
	if err != nil {
		return nil, err
	}

	if module, ok := byName[name]; ok {
		return module, nil
	} else {
		return nil, fmt.Errorf("Resource not found: %s", name)
	}
}

var _ core.Grouper = (*Module)(nil)

// END: Grouper interface
