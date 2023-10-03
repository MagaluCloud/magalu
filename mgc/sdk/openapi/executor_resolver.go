package openapi

import (
	"fmt"

	"github.com/go-openapi/jsonpointer"
	"magalu.cloud/core"
)

type executorTree struct {
	exec core.Executor
	tree map[string]*executorTree
}

// JSONLookup implements github.com/go-openapi/jsonpointer#JSONPointable
func (t executorTree) JSONLookup(token string) (interface{}, error) {
	ref, ok := t.tree[token]
	if !ok {
		return nil, fmt.Errorf("object has no field %q", token)
	}

	return ref, nil
}

var _ jsonpointer.JSONPointable = (*executorTree)(nil)

func (t *executorTree) add(key []string, exec core.Executor) error {
	if len(key) == 0 {
		if t.exec != nil {
			return fmt.Errorf("already exists as %+v, want to add %+v", t.exec, exec)
		}
		t.exec = exec
		return nil
	}

	if t.tree == nil {
		t.tree = map[string]*executorTree{}
	}

	current := key[0]
	childT, ok := t.tree[current]
	if !ok {
		childT = &executorTree{}
		t.tree[current] = childT
	}

	if err := childT.add(key[1:], exec); err != nil {
		return fmt.Errorf("%q %w", current, err)
	}

	return nil
}

type executorResolver struct {
	byId   map[string]core.Executor
	byPath executorTree
}

func (o *executorResolver) add(id string, path []string, exec core.Executor) error {
	if o.byId == nil {
		o.byId = map[string]core.Executor{}
	}

	if id != "" {
		if old, exists := o.byId[id]; exists {
			return fmt.Errorf("id %q already exists as %+v, want to add %+v", id, old, exec)
		}
		o.byId[id] = exec
	}

	return o.byPath.add(path, exec)
}

func (o *executorResolver) get(id string) core.Executor {
	if exec, ok := o.byId[id]; ok {
		return exec
	}
	return nil
}

func (o *executorResolver) resolve(ref string) (core.Executor, error) {
	jp, err := jsonpointer.New(ref)
	if err != nil {
		return nil, err
	}

	result, _, err := jp.Get(o.byPath)
	if err != nil {
		return nil, err
	}
	if tree, ok := result.(*executorTree); ok && tree.exec != nil {
		return tree.exec, nil
	}
	return nil, fmt.Errorf("reference %q doesn't resolve to executorTree with valid executor but %#v", ref, result)
}
