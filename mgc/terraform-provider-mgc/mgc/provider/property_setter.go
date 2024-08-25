package provider

import (
	"fmt"
	"strings"

	"magalu.cloud/core"
	"magalu.cloud/core/utils"
)

type propertySetter interface {
	propertyName() mgcName

	allParametersSchemas() []*core.Schema
	allResultSchemas() []*core.Schema

	getTarget(currentValue, targetValue core.Value) (core.Linker, error)
}

// update:field
type defaultPropertySetter struct {
	propName mgcName
	target   core.Linker
}

func newDefaultPropertySetter(propName mgcName, target core.Linker) *defaultPropertySetter {
	return &defaultPropertySetter{propName, target}
}

func (o *defaultPropertySetter) propertyName() mgcName {
	return o.propName
}

func (o *defaultPropertySetter) allParametersSchemas() []*core.Schema {
	return []*core.Schema{o.target.AdditionalParametersSchema()}
}

func (o *defaultPropertySetter) allResultSchemas() []*core.Schema {
	return []*core.Schema{o.target.ResultSchema()}
}

func (a *defaultPropertySetter) getTarget(currentValue, targetValue core.Value) (core.Linker, error) {
	if utils.IsSameValueOrPointer(currentValue, targetValue) {
		return nil, fmt.Errorf("trying to call property setter with current value")
	}
	return a.target, nil
}

var _ propertySetter = (*defaultPropertySetter)(nil)

// update:field:value
// update:field:other
type enumPropertySetter struct {
	propName      mgcName
	targetByValue map[any]core.Linker
}

func newEnumPropertySetter(propName mgcName, containerEntries []*propertySetterContainerEntry) *enumPropertySetter {
	targetsByValue := make(map[any]core.Linker, len(containerEntries))
	for _, entry := range containerEntries {
		targetsByValue[entry.args[0]] = entry.target
	}
	return &enumPropertySetter{propName, targetsByValue}
}

func (o *enumPropertySetter) propertyName() mgcName {
	return o.propName
}

func (o *enumPropertySetter) allParametersSchemas() []*core.Schema {
	result := make([]*core.Schema, 0, len(o.targetByValue))
	for _, target := range o.targetByValue {
		result = append(result, target.AdditionalParametersSchema())
	}
	return result
}

func (o *enumPropertySetter) allResultSchemas() []*core.Schema {
	result := make([]*core.Schema, 0, len(o.targetByValue))
	for _, target := range o.targetByValue {
		result = append(result, target.ResultSchema())
	}
	return result
}

func (a *enumPropertySetter) getTarget(currentValue, targetValue core.Value) (core.Linker, error) {
	if utils.IsSameValueOrPointer(currentValue, targetValue) {
		return nil, fmt.Errorf("trying to call property setter with current value")
	}

	target, ok := a.targetByValue[targetValue]
	if !ok {
		enum := make([]any, 0, len(a.targetByValue))
		for k := range a.targetByValue {
			enum = append(enum, k)
		}
		return nil, fmt.Errorf("%+v is not one of the allowed values: %v", targetValue, enum)
	}

	return target, nil
}

var _ propertySetter = (*enumPropertySetter)(nil)

type valueTransition struct {
	current string
	target  string
}

func (t *valueTransition) String() string {
	return fmt.Sprintf("%s to %s", t.current, t.target)
}

// update:field:from:to
// update:field:from_other:to_other
type strTransitionPropertySetter struct {
	propName            mgcName
	targetsByTransition map[valueTransition]core.Linker
}

func newStrTransitionPropertySetter(propName mgcName, containerEntries []*propertySetterContainerEntry) *strTransitionPropertySetter {
	targetsByTransition := make(map[valueTransition]core.Linker, len(containerEntries))
	for _, entry := range containerEntries {
		targetsByTransition[valueTransition{entry.args[0], entry.args[1]}] = entry.target
	}
	return &strTransitionPropertySetter{propName, targetsByTransition}
}

func (o *strTransitionPropertySetter) propertyName() mgcName {
	return o.propName
}

func (o *strTransitionPropertySetter) allParametersSchemas() []*core.Schema {
	result := make([]*core.Schema, 0, len(o.targetsByTransition))
	for _, target := range o.targetsByTransition {
		result = append(result, target.AdditionalParametersSchema())
	}
	return result
}

func (o *strTransitionPropertySetter) allResultSchemas() []*core.Schema {
	result := make([]*core.Schema, 0, len(o.targetsByTransition))
	for _, target := range o.targetsByTransition {
		result = append(result, target.ResultSchema())
	}
	return result
}

func (a *strTransitionPropertySetter) getTarget(currentValue, targetValue core.Value) (core.Linker, error) {
	if utils.IsSameValueOrPointer(currentValue, targetValue) {
		return nil, fmt.Errorf("trying to call property setter with current value")
	}

	currentStr, ok := currentValue.(string)
	if !ok {
		return nil, fmt.Errorf("strTransitionPropSetter only accepts string values, got %T for 'currentValue': %+v", currentValue, currentValue)
	}
	targetStr, ok := targetValue.(string)
	if !ok {
		return nil, fmt.Errorf("strTransitionPropSetter only accepts string values, got %T for 'targetValue': %+v", targetValue, targetValue)
	}

	transition := valueTransition{currentStr, targetStr}
	target, ok := a.targetsByTransition[transition]
	if !ok {
		allTransitions := make([]valueTransition, 0, len(a.targetsByTransition))
		for t := range a.targetsByTransition {
			allTransitions = append(allTransitions, t)
		}
		return nil, fmt.Errorf("%s to %s is not one of the allowed transitions: %+v", currentStr, targetStr, allTransitions)
	}

	return target, nil
}

var _ propertySetter = (*strTransitionPropertySetter)(nil)

type propertySetterContainerEntry struct {
	args   []string
	target core.Linker
}

type propertySetterContainer struct {
	key      mgcName
	entries  []*propertySetterContainerEntry
	argCount int
}

func newPropertySetterContainer(key mgcName, entry *propertySetterContainerEntry) *propertySetterContainer {
	return &propertySetterContainer{key, []*propertySetterContainerEntry{entry}, len(entry.args)}
}

func (c *propertySetterContainer) add(entry *propertySetterContainerEntry) error {
	if c.argCount == 0 {
		return fmt.Errorf("can't have more than one simple property setter (no args) with the same key")
	}
	if len(entry.args) != c.argCount {
		return fmt.Errorf("all property setters with the same key must have the same arg count")
	}
	c.entries = append(c.entries, entry)
	return nil
}

func collectPropertySetterContainers(readLinks core.Links) (containersByKey map[mgcName]*propertySetterContainer, err error) {
	containersByKey = map[mgcName]*propertySetterContainer{}
	for linkName, link := range readLinks {
		if !strings.HasPrefix(linkName, "update") {
			continue
		}

		setterParts := strings.Split(linkName, "/")
		if len(setterParts) < 2 {
			continue
		}

		key := mgcName(setterParts[1])
		args := setterParts[2:]

		if container, ok := containersByKey[key]; ok {
			err := container.add(&propertySetterContainerEntry{args, link})
			if err != nil {
				return nil, fmt.Errorf("unable to parse property setter link %q: %w", linkName, err)
			}
		} else {
			containersByKey[key] = newPropertySetterContainer(key, &propertySetterContainerEntry{args, link})
		}
	}
	return containersByKey, nil
}
