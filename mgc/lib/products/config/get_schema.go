/*
Executor: get-schema

# Summary

# Get the JSON Schema for the specified Config

# Description

Get the JSON Schema for the specified Config. The Schema has
information about the accepted values for the Config, constraints, type, description, etc.

import "magalu.cloud/lib/products/config"
*/
package config

import (
	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetSchemaParameters struct {
	Key string `json:"key"`
}

type GetSchemaResult struct {
	AdditionalProperties GetSchemaResultAdditionalProperties `json:"additionalProperties,omitempty"`
	AllOf                GetSchemaResultAllOf                `json:"allOf,omitempty"`
	AllowEmptyValue      bool                                `json:"allowEmptyValue,omitempty"`
	AnyOf                GetSchemaResultAnyOf                `json:"anyOf,omitempty"`
	Default              *GetSchemaResultDefault             `json:"default,omitempty"`
	Deprecated           bool                                `json:"deprecated,omitempty"`
	Description          string                              `json:"description,omitempty"`
	Discriminator        GetSchemaResultDiscriminator        `json:"discriminator,omitempty"`
	Enum                 GetSchemaResultEnum                 `json:"enum,omitempty"`
	Example              *GetSchemaResultExample             `json:"example,omitempty"`
	ExclusiveMax         bool                                `json:"exclusiveMax,omitempty"`
	ExclusiveMin         bool                                `json:"exclusiveMin,omitempty"`
	Extensions           GetSchemaResultExtensions           `json:"extensions,omitempty"`
	ExternalDocs         GetSchemaResultExternalDocs         `json:"externalDocs,omitempty"`
	Format               string                              `json:"format,omitempty"`
	Items                *GetSchemaResultItems               `json:"items,omitempty"`
	Max                  float64                             `json:"max,omitempty"`
	MaxItems             int                                 `json:"maxItems,omitempty"`
	MaxLength            int                                 `json:"maxLength,omitempty"`
	MaxProps             int                                 `json:"maxProps,omitempty"`
	Min                  float64                             `json:"min,omitempty"`
	MinItems             int                                 `json:"minItems,omitempty"`
	MinLength            int                                 `json:"minLength,omitempty"`
	MinProps             int                                 `json:"minProps,omitempty"`
	MultipleOf           float64                             `json:"multipleOf,omitempty"`
	Not                  *GetSchemaResultNot                 `json:"not,omitempty"`
	Nullable             bool                                `json:"nullable,omitempty"`
	OneOf                GetSchemaResultOneOf                `json:"oneOf,omitempty"`
	Pattern              string                              `json:"pattern,omitempty"`
	Properties           GetSchemaResultProperties           `json:"properties,omitempty"`
	ReadOnly             bool                                `json:"readOnly,omitempty"`
	Required             GetSchemaResultRequired             `json:"required,omitempty"`
	Title                string                              `json:"title,omitempty"`
	Type                 string                              `json:"type,omitempty"`
	UniqueItems          bool                                `json:"uniqueItems,omitempty"`
	WriteOnly            bool                                `json:"writeOnly,omitempty"`
	Xml                  GetSchemaResultXml                  `json:"xml,omitempty"`
}

type GetSchemaResultAdditionalProperties struct {
	Has    *bool                                      `json:"has,omitempty"`
	Schema *GetSchemaResultAdditionalPropertiesSchema `json:"schema,omitempty"`
}

type GetSchemaResultAdditionalPropertiesSchema struct {
	Ref   string           `json:"ref,omitempty"`
	Value *GetSchemaResult `json:"value,omitempty"`
}

type GetSchemaResultAllOfItem struct {
	Ref   string           `json:"ref,omitempty"`
	Value *GetSchemaResult `json:"value,omitempty"`
}

type GetSchemaResultAllOf []*GetSchemaResultAllOfItem

type GetSchemaResultAnyOfItem struct {
	Ref   string           `json:"ref,omitempty"`
	Value *GetSchemaResult `json:"value,omitempty"`
}

type GetSchemaResultAnyOf []*GetSchemaResultAnyOfItem

// any of: bool, string, float64, int, GetSchemaResultDefault4, GetSchemaResultDefault5
type GetSchemaResultDefault any

type GetSchemaResultDefault4 []any

type GetSchemaResultDefault5 struct {
}

type GetSchemaResultDiscriminator struct {
	Extensions   GetSchemaResultDiscriminatorExtensions `json:"extensions,omitempty"`
	Mapping      GetSchemaResultDiscriminatorMapping    `json:"mapping,omitempty"`
	PropertyName string                                 `json:"propertyName,omitempty"`
}

type GetSchemaResultDiscriminatorExtensions struct {
}

type GetSchemaResultDiscriminatorMapping struct {
}

// any of: bool, string, float64, int, GetSchemaResultEnumItem4, GetSchemaResultEnumItem5
type GetSchemaResultEnumItem any

type GetSchemaResultEnumItem4 []any

type GetSchemaResultEnumItem5 struct {
}

type GetSchemaResultEnum []*GetSchemaResultEnumItem

// any of: bool, string, float64, int, GetSchemaResultExample4, GetSchemaResultExample5
type GetSchemaResultExample any

type GetSchemaResultExample4 []any

type GetSchemaResultExample5 struct {
}

type GetSchemaResultExtensions struct {
}

type GetSchemaResultExternalDocs struct {
	Description string                                `json:"description,omitempty"`
	Extensions  GetSchemaResultExternalDocsExtensions `json:"extensions,omitempty"`
	Url         string                                `json:"url,omitempty"`
}

type GetSchemaResultExternalDocsExtensions struct {
}

type GetSchemaResultItems struct {
	Ref   string           `json:"ref,omitempty"`
	Value *GetSchemaResult `json:"value,omitempty"`
}

type GetSchemaResultNot struct {
	Ref   string           `json:"ref,omitempty"`
	Value *GetSchemaResult `json:"value,omitempty"`
}

type GetSchemaResultOneOfItem struct {
	Ref   string           `json:"ref,omitempty"`
	Value *GetSchemaResult `json:"value,omitempty"`
}

type GetSchemaResultOneOf []*GetSchemaResultOneOfItem

type GetSchemaResultProperties struct {
}

type GetSchemaResultRequired []string

type GetSchemaResultXml struct {
	Attribute  bool                         `json:"attribute,omitempty"`
	Extensions GetSchemaResultXmlExtensions `json:"extensions,omitempty"`
	Name       string                       `json:"name,omitempty"`
	Namespace  string                       `json:"namespace,omitempty"`
	Prefix     string                       `json:"prefix,omitempty"`
	Wrapped    bool                         `json:"wrapped,omitempty"`
}

type GetSchemaResultXmlExtensions struct {
}

func (s *service) GetSchema(
	parameters GetSchemaParameters,
) (
	result GetSchemaResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("GetSchema", mgcCore.RefPath("/config/get-schema"), s.client, s.ctx)
	if err != nil {
		return
	}

	var p mgcCore.Parameters
	if p, err = mgcHelpers.ConvertParameters[GetSchemaParameters](parameters); err != nil {
		return
	}

	var c mgcCore.Configs

	r, err := exec.Execute(ctx, p, c)
	if err != nil {
		return
	}
	return mgcHelpers.ConvertResult[GetSchemaResult](r)
}

// TODO: links
// TODO: related
