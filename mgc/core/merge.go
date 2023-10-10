package core

import "fmt"

type MergeGroup struct {
	name        string
	version     string
	description string
	merge       []Grouper
	children    []Descriptor
	childByName map[string]Descriptor
}

func NewMergeGroup(name string, version string, description string, merge []Grouper) *MergeGroup {
	return &MergeGroup{name, version, description, merge, nil, nil}
}

// BEGIN: Descriptor interface:

func (o *MergeGroup) Name() string {
	return o.name
}

func (o *MergeGroup) Version() string {
	return o.version
}

func (o *MergeGroup) Description() string {
	return o.description
}

// END: Descriptor interface

// BEGIN: Grouper interface:

func (o *MergeGroup) mergeAfter(target Grouper, start int) []Grouper {
	result := make([]Grouper, 1, len(o.merge)-start+1)
	result[0] = target

	name := target.Name()

	for i := start; i < len(o.merge); i += 1 {
		child, err := o.merge[i].GetChildByName(name)
		if err == nil {
			if group, ok := child.(Grouper); ok {
				result = append(result, group)
			}
		}
	}
	return result
}

func (o *MergeGroup) getChildren() (children []Descriptor, byName map[string]Descriptor, err error) {
	if o.children != nil {
		return o.children, o.childByName, nil
	}

	o.children = make([]Descriptor, 0)
	o.childByName = make(map[string]Descriptor)

	for i, c := range o.merge {
		finished, err := c.VisitChildren(func(child Descriptor) (run bool, err error) {
			name := child.Name()

			if o.childByName[name] != nil {
				return true, nil
			}

			if group, ok := child.(Grouper); ok {
				merge := o.mergeAfter(group, i+1)
				if len(merge) > 1 {
					child = &MergeGroup{
						name:        name,
						version:     group.Version(),
						description: group.Description(),
						merge:       merge,
					}
				}
			}

			o.children = append(o.children, child)
			o.childByName[name] = child
			return true, nil
		})

		if err != nil {
			return nil, nil, err
		}
		if !finished {
			return nil, nil, err
		}
	}

	return o.children, o.childByName, nil
}

func (o *MergeGroup) VisitChildren(visitor DescriptorVisitor) (finished bool, err error) {
	children, _, err := o.getChildren()
	if err != nil {
		return false, err
	}

	for _, res := range children {
		run, err := visitor(res)
		if err != nil {
			return false, err
		}
		if !run {
			return false, nil
		}
	}
	return true, nil
}

func (o *MergeGroup) GetChildByName(name string) (child Descriptor, err error) {
	_, byName, err := o.getChildren()
	if err != nil {
		return nil, err
	}
	res, ok := byName[name]
	if !ok {
		return nil, fmt.Errorf("child not found: %s", name)
	}

	return res, nil
}

var _ Grouper = (*MergeGroup)(nil)

// END: Grouper interface
