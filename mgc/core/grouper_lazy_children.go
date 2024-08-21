package core

import (
	"fmt"
	"strings"
	"sync"

	"slices"
)

var childrenMutex sync.Mutex

// Implements VisitChildren() and GetChildByName() by calling createChildren()
// only once, when/if needed, and storing the results.
//
// Children will be automatically sorted by their names, so visitor are always in a predictable and stable order.
type GrouperLazyChildren[T Descriptor] struct {
	createChildren func() ([]T, error)
	children       []T
	childByName    map[string]T
}

func NewGrouperLazyChildren[T Descriptor](createChildren func() ([]T, error)) *GrouperLazyChildren[T] {
	if createChildren == nil {
		panic("createChildren == nil")
	}
	return &GrouperLazyChildren[T]{createChildren, nil, nil}
}

func (g *GrouperLazyChildren[T]) getChildren() (children []T, byName map[string]T, err error) {
	if g.createChildren == nil {
		return g.children, g.childByName, nil
	}

	children, err = g.createChildren()
	if err != nil {
		return nil, nil, err
	}

	g.createChildren = nil // avoid reloading and also free/GC any resources used by the loader function
	g.children = children
	g.childByName = make(map[string]T, len(children))

	for _, child := range children {
		childrenMutex.Lock()
		g.childByName[child.Name()] = child
		childrenMutex.Unlock()
	}

	slices.SortFunc(g.children, func(a, b T) int {
		return strings.Compare(a.Name(), b.Name())
	})

	return g.children, g.childByName, nil
}

func (g *GrouperLazyChildren[T]) VisitChildren(visitor DescriptorVisitor) (finished bool, err error) {
	children, _, err := g.getChildren()
	if err != nil {
		return false, err
	}

	for _, child := range children {
		run, err := visitor(child)
		if err != nil {
			return false, err
		}
		if !run {
			return false, nil
		}
	}
	return true, nil
}

func (g *GrouperLazyChildren[T]) GetChildByName(name string) (child Descriptor, err error) {
	_, childByName, err := g.getChildren()
	if err != nil {
		return nil, err
	}
	child, ok := childByName[name]
	if !ok {
		return nil, fmt.Errorf("child not found: %s", name)
	}

	return child, nil
}
