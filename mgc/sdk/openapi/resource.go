package openapi

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
	"magalu.cloud/core"

	"github.com/getkin/kin-openapi/openapi3"
)

// Source -> Module -> Resource -> Operation

// Resource

type Resource struct {
	tag             *openapi3.Tag
	doc             *openapi3.T
	extensionPrefix *string
	servers         openapi3.Servers
	operations      *map[string]*Operation
}

// BEGIN: Descriptor interface:

func (o *Resource) Name() string {
	return getNameExtension(o.extensionPrefix, o.tag.Extensions, o.tag.Name)
}

func (o *Resource) Version() string {
	return ""
}

func (o *Resource) Description() string {
	return getDescriptionExtension(o.extensionPrefix, o.tag.Extensions, o.tag.Description)
}

// END: Descriptor interface

// BEGIN: Grouper interface:

func getServers(p *openapi3.PathItem, op *openapi3.Operation) openapi3.Servers {
	var servers openapi3.Servers
	if op.Servers != nil && len(*op.Servers) > 0 {
		servers = *op.Servers
	}
	if servers == nil && len(p.Servers) > 0 {
		servers = p.Servers
	}

	return servers
}

// NOTE: some OpenAPIs may have have similar operations with the same tag (== Resource),
// in order to disambiguate, we need to pass the whole set of names and then find the
// the action name, simplifying if no collisions
type operationDesc struct {
	path   *openapi3.PathItem
	op     *openapi3.Operation
	method string
	key    string
}

// a tree based on other maps or operationDesc
type operationTree struct {
	tree map[string]*operationTree
	desc *operationDesc
}

func (t *operationTree) Add(key []string, desc *operationDesc) error {
	if len(key) == 0 {
		t.desc = desc
		return nil
	}

	if t.tree == nil {
		t.tree = map[string]*operationTree{}
	}

	current := key[0]
	childT, ok := t.tree[current]
	if !ok {
		childT = &operationTree{}
		t.tree[current] = childT
	}

	return childT.Add(key[1:], desc)
}

type operationTreePath struct {
	key    string
	parent *operationTree
}

func (t *operationTree) VisitDesc(path []operationTreePath, visitor func(path []operationTreePath, desc *operationDesc) bool) bool {
	if t.desc != nil {
		if !visitor(path, t.desc) {
			return false
		}
	}

	for k, childT := range t.tree {
		oldLen := len(path)
		path = append(path, operationTreePath{k, t})

		if !childT.VisitDesc(path, visitor) {
			return false
		}

		path = path[:oldLen]
	}

	return true
}

var openAPIPathArgRegex = regexp.MustCompile("[{](?P<name>[^}]+)[}]")

func getPathEntry(pathEntry string) (string, bool) {
	match := openAPIPathArgRegex.FindStringSubmatch(pathEntry)
	if len(match) > 0 {
		for i, substr := range match {
			if openAPIPathArgRegex.SubexpNames()[i] == "name" {
				return substr, true
			}
		}
	}

	return pathEntry, false
}

func getCoalescedPath(path []operationTreePath) []string {
	parts := []string{}
	wasVariable := false
	for _, p := range path {
		pathEntry, isVariable := getPathEntry(p.key)
		if len(p.parent.tree) > 1 || wasVariable {
			parts = append(parts, pathEntry)
		}
		wasVariable = isVariable
	}
	return parts
}

func getFullPath(path []operationTreePath) []string {
	parts := []string{}
	for _, p := range path {
		pathEntry, _ := getPathEntry(p.key)
		parts = append(parts, pathEntry)
	}
	return parts
}

func renamePath(httpMethod string, pathName string) string {
	switch httpMethod {
	case "post":
		return "create"
	case "put":
		return "replace"
	case "patch":
		return "update"
	case "get":
		// only consider "get" if ends with, mid-path are still list, ex:
		// GET:  /resource/{id}
		// LIST: /{containerId}/resource
		// GET:  /{containerId}/resource/{id}
		if strings.HasSuffix(pathName, "}") {
			return "get"
		}
		return "list"
	}

	return httpMethod
}

func getFullOperationName(httpMethod string, pathName string) []string {
	actionName := renamePath(httpMethod, pathName)
	name := []string{actionName}

	for _, pathEntry := range strings.Split(pathName, "/") {
		if pathEntry == "" {
			continue
		}
		name = append(name, pathEntry)
	}

	return name
}

func (o *Resource) collectOperations() *operationTree {
	tree := &operationTree{}
	for key, path := range o.doc.Paths {
		if getHiddenExtension(o.extensionPrefix, path.Extensions) {
			continue
		}

		pathOps := map[string]*openapi3.Operation{
			"get":    path.Get,
			"post":   path.Post,
			"put":    path.Put,
			"patch":  path.Patch,
			"delete": path.Delete,
		}

		for method, op := range pathOps {
			if op == nil || getHiddenExtension(o.extensionPrefix, op.Extensions) {
				continue
			}

			if !slices.Contains(op.Tags, o.tag.Name) {
				continue
			}

			name := getFullOperationName(method, key)
			if err := tree.Add(name, &operationDesc{path, op, method, key}); err != nil {
				log.Printf("Failed to add %s %s: %s", method, key, err)
			}
		}
	}

	return tree
}

func (o *Resource) getOperations() map[string]*Operation {
	if o.operations == nil {
		opMap := map[string]*Operation{}
		opTree := o.collectOperations()
		opTree.VisitDesc([]operationTreePath{}, func(path []operationTreePath, desc *operationDesc) bool {
			opName := getNameExtension(o.extensionPrefix, desc.op.Extensions, "")
			if opName == "" {
				opName = strings.Join(getCoalescedPath(path), "-")
				if _, ok := opMap[opName]; ok {
					opName = strings.Join(getFullPath(path), "-")
				}
			}

			servers := getServers(desc.path, desc.op)
			if servers == nil {
				servers = o.servers
			}

			opMap[opName] = &Operation{
				name:            opName,
				key:             desc.key,
				method:          desc.method,
				path:            desc.path,
				operation:       desc.op,
				doc:             o.doc,
				extensionPrefix: o.extensionPrefix,
				servers:         servers,
			}

			return true
		})
		o.operations = &opMap
	}
	return *o.operations
}

func (o *Resource) VisitChildren(visitor core.DescriptorVisitor) (finished bool, err error) {
	for _, op := range o.getOperations() {
		run, err := visitor(op)
		if err != nil {
			return false, err
		}
		if !run {
			return false, nil
		}
	}
	return true, nil
}

func (o *Resource) GetChildByName(name string) (child core.Descriptor, err error) {
	op, ok := o.getOperations()[name]
	if !ok {
		return nil, fmt.Errorf("Action not found: %s", name)
	}

	return op, nil
}

var _ core.Grouper = (*Resource)(nil)

// END: Grouper interface
