package core

import (
	"errors"
	"fmt"
)

type DescriptorSpec struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Description  string `json:"description"`
	Summary      string `json:"summary"`
	IsInternal   *bool  `json:"isInternal,omitempty"`
	Scopes       Scopes `json:"scopes"`
	Observations string `json:"observation,omitempty"`
	GroupID      string `json:"groupId,omitempty"`
}

func (d *DescriptorSpec) Validate() error {
	if d.Name == "" {
		return &ChainedError{Name: fmt.Sprintf("<missing name %p>", d), Err: errors.New("missing name")}
	}
	if d.Description == "" {
		return &ChainedError{Name: d.Name, Err: errors.New("missing description")}
	}
	// Version and Summary are optional
	return nil
}

// General interface that describes both Executor and Grouper
type Descriptor interface {
	Name() string
	Version() string
	Description() string
	Summary() string
	IsInternal() bool
	Scopes() Scopes
	DescriptorSpec() DescriptorSpec
	GroupID() string
}

type SimpleDescriptor struct {
	Spec DescriptorSpec
}

func (d *SimpleDescriptor) GroupID() string {
	if d.Spec.GroupID == "" {
		return "catalog"
	}
	return d.Spec.GroupID
}

func (d *SimpleDescriptor) Name() string {
	return d.Spec.Name
}

func (d *SimpleDescriptor) Version() string {
	return d.Spec.Version
}

func (d *SimpleDescriptor) Description() string {
	return d.Spec.Description
}

func (d *SimpleDescriptor) IsInternal() bool {
	if d.Spec.IsInternal == nil {
		return false
	}
	return *d.Spec.IsInternal
}

func (d *SimpleDescriptor) Scopes() Scopes {
	return d.Spec.Scopes
}

func (d *SimpleDescriptor) DescriptorSpec() DescriptorSpec {
	return d.Spec
}

func (d *SimpleDescriptor) Summary() string {
	if d.Spec.Summary == "" {
		return d.Spec.Description
	}
	return d.Spec.Summary
}

var _ Descriptor = (*SimpleDescriptor)(nil)

type DescriptorVisitor func(child Descriptor) (run bool, err error)
