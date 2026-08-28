package core

import (
	"errors"
	"fmt"
)

// ProjectScope diz de qual chave de configuração um comando tira o escopo de
// projeto. Vazio (o default) significa que o produto não é escopável — não
// recebe 'x-project-id' nem ganha as flags de escopo.
//
// São dois valores porque o IAM codifica o escopo de outro jeito: para ele,
// omitir o header significa "o tenant inteiro", enquanto para os demais
// significa "o projeto default". Ver ProjectScopeIAM.
type ProjectScope string

const (
	// ProjectScopeProduct lê a chave `project`. Omitir = projeto default.
	ProjectScopeProduct ProjectScope = "project"
	// ProjectScopeIAM lê a chave `iamProject`. Omitir = tenant inteiro, e é por
	// isso que escrita de IAM exige escopo explícito.
	ProjectScopeIAM ProjectScope = "iam"
)

type DescriptorSpec struct {
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Description  string       `json:"description"`
	Summary      string       `json:"summary"`
	IsInternal   *bool        `json:"isInternal,omitempty"`
	Scopes       Scopes       `json:"scopes"`
	Observations string       `json:"observation,omitempty"`
	GroupID      string       `json:"groupId,omitempty"`
	ProjectScope ProjectScope `json:"projectScope,omitempty"`
	// ScopeRequired diz que este comando não age sem escopo explícito. Vale para
	// escrita de IAM, onde omitir o escopo significa "o tenant inteiro" — errar
	// ali não falha uma request, muda configuração da organização toda.
	//
	// Declarado por serviço (quais métodos) e desligável por operação, nunca
	// ligável: a marcação é opt-out para que endpoint novo nasça protegido.
	ScopeRequired bool `json:"scopeRequired,omitempty"`
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
