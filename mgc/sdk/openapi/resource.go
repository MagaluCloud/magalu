package openapi

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"magalu.cloud/core"

	"github.com/PaesslerAG/gval"
	"github.com/PaesslerAG/jsonpath"
	"github.com/getkin/kin-openapi/openapi3"
)

type WaitTermination struct {
	MaxRetries        int           `json:"maxRetries,omitempty"`
	IntervalInSeconds time.Duration `json:"intervalInSeconds,omitempty"`
	JSONPathQuery     string        `json:"jsonPathQuery"`
}

var defaultWaitTerminaton = WaitTermination{MaxRetries: 30, IntervalInSeconds: time.Second}

// Source -> Module -> Resource -> Operation

// Resource

type Resource struct {
	tag             *openapi3.Tag
	doc             *openapi3.T
	extensionPrefix *string
	servers         openapi3.Servers
	operations      []core.Executor
	byName          map[string]core.Executor
	logger          *zap.SugaredLogger
	execResolver    *executorResolver
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

func (t *operationTree) VisitDesc(path []operationTreePath, visitor func(path []operationTreePath, desc *operationDesc) (bool, error)) (bool, error) {
	if t.desc != nil {
		if run, err := visitor(path, t.desc); !run || err != nil {
			return false, err
		}
	}

	for k, childT := range t.tree {
		oldLen := len(path)
		path = append(path, operationTreePath{k, t})

		if run, err := childT.VisitDesc(path, visitor); !run || err != nil {
			return false, err
		}

		path = path[:oldLen]
	}

	return true, nil
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
				o.logger.Warnw("failed to add operation", "method", method, "key", key, "error", err)
			}
		}
	}

	return tree
}

func (o *Resource) getOperations() (operations []core.Executor, byName map[string]core.Executor, err error) {
	if len(o.operations) == 0 {
		return o.operations, o.byName, nil
	}

	o.operations = make([]core.Executor, 0)
	o.byName = make(map[string]core.Executor)

	opTree := o.collectOperations()
	_, err = opTree.VisitDesc([]operationTreePath{}, func(path []operationTreePath, desc *operationDesc) (bool, error) {
		opName := getNameExtension(o.extensionPrefix, desc.op.Extensions, "")
		if opName == "" {
			opName = strings.Join(getCoalescedPath(path), "-")
			if _, ok := o.byName[opName]; ok {
				opName = strings.Join(getFullPath(path), "-")
			}
		}

		servers := getServers(desc.path, desc.op)
		if servers == nil {
			servers = o.servers
		}

		var operation core.Executor = &Operation{
			name:            opName,
			key:             desc.key,
			method:          strings.ToUpper(desc.method),
			path:            desc.path,
			operation:       desc.op,
			doc:             o.doc,
			extensionPrefix: o.extensionPrefix,
			servers:         servers,
			logger:          o.logger.Named(opName),
			execResolver:    o.execResolver,
		}

		if wtExt, ok := getExtensionObject(o.extensionPrefix, "wait-termination", desc.op.Extensions, nil); ok && wtExt != nil {
			if tExec, err := o.wrapInTerminatorExecutor(wtExt, operation); err == nil {
				operation = tExec
			} else {
				return false, err
			}
		}

		if output, ok := getExtensionString(o.extensionPrefix, "output-flag", desc.op.Extensions, ""); ok && output != "" {
			operation = core.NewExecuteResultOutputOptions(operation, func(exec core.Executor, result core.Value) string {
				return output
			})
		}

		err = o.execResolver.add(
			desc.op.OperationID,
			[]string{"paths", desc.key, desc.method},
			operation,
		)
		if err != nil {
			return false, err
		}

		o.operations = append(o.operations, operation)
		o.byName[opName] = operation

		return true, nil
	})

	if err != nil {
		o.operations = nil
		o.byName = nil
		return nil, nil, err
	}

	slices.SortFunc(o.operations, func(a, b core.Executor) int {
		return strings.Compare(a.Name(), b.Name())
	})

	return o.operations, o.byName, nil
}

func (o *Resource) wrapInTerminatorExecutor(wtExt map[string]any, exec core.Executor) (core.TerminatorExecutor, error) {
	wt := &defaultWaitTerminaton
	if err := core.DecodeValue(wtExt, wt); err != nil {
		o.logger.Warnw("error decoding extension wait-termination", "data", wtExt, "error", err)
		return nil, err
	}

	if wt.MaxRetries <= 0 {
		wt.MaxRetries = defaultWaitTerminaton.MaxRetries
	}
	if wt.IntervalInSeconds <= 0 {
		wt.IntervalInSeconds = defaultWaitTerminaton.IntervalInSeconds
	}

	builder := gval.Full(jsonpath.PlaceholderExtension())
	jp, err := builder.NewEvaluable(wt.JSONPathQuery)
	if err == nil {
		tExec := core.NewTerminatorExecutorWithCheck(exec, wt.MaxRetries, wt.IntervalInSeconds, func(ctx context.Context, exec core.Executor, result core.Value) bool {
			v, err := jp(ctx, result)
			if err != nil {
				o.logger.Warnw("error evaluating jsonpath query", "query", wt.JSONPathQuery, "target", result, "error", err)
				return false
			}

			if v == nil {
				return false
			} else if lst, ok := v.([]any); ok {
				return len(lst) > 0
			} else if m, ok := v.(map[string]any); ok {
				return len(m) > 0
			} else if b, ok := v.(bool); ok {
				return b
			} else {
				o.logger.Warnw("unknown jsonpath result. Expected list, map or boolean", "result", result)
				return false
			}
		})
		return tExec, nil
	} else {
		o.logger.Warnw("error parsing jsonpath. Executing without polling", "expression", wt.JSONPathQuery, "error", err)
		return nil, err
	}
}

func (o *Resource) VisitChildren(visitor core.DescriptorVisitor) (finished bool, err error) {
	operations, _, err := o.getOperations()
	if err != nil {
		return false, err
	}

	for _, op := range operations {
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
	_, byName, err := o.getOperations()
	if err != nil {
		return nil, err
	}

	op, ok := byName[name]
	if !ok {
		return nil, fmt.Errorf("Action not found: %s", name)
	}

	return op, nil
}

var _ core.Grouper = (*Resource)(nil)

// END: Grouper interface
