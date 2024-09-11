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
	"context"

	mgcCore "magalu.cloud/core"
	mgcHelpers "magalu.cloud/lib/helpers"
)

type GetSchemaParameters struct {
	Key string `json:"key"`
}

type GetSchemaResult struct {
	AdditionalProperties *GetSchemaResultAdditionalProperties                         `json:"additionalProperties,omitempty"`
	AllOf                *GetSchemaResultAdditionalPropertiesSchemaValueAllOf         `json:"allOf,omitempty"`
	AllowEmptyValue      *bool                                                        `json:"allowEmptyValue,omitempty"`
	AnyOf                *GetSchemaResultAdditionalPropertiesSchemaValueAnyOf         `json:"anyOf,omitempty"`
	Default              *GetSchemaResultAdditionalPropertiesSchemaValueDefault       `json:"default,omitempty"`
	Deprecated           *bool                                                        `json:"deprecated,omitempty"`
	Description          *string                                                      `json:"description,omitempty"`
	Discriminator        *GetSchemaResultAdditionalPropertiesSchemaValueDiscriminator `json:"discriminator,omitempty"`
	Enum                 *GetSchemaResultAdditionalPropertiesSchemaValueEnum          `json:"enum,omitempty"`
	Example              *GetSchemaResultAdditionalPropertiesSchemaValueExample       `json:"example,omitempty"`
	ExclusiveMax         *bool                                                        `json:"exclusiveMax,omitempty"`
	ExclusiveMin         *bool                                                        `json:"exclusiveMin,omitempty"`
	Extensions           *GetSchemaResultAdditionalPropertiesSchemaValueExtensions    `json:"extensions,omitempty"`
	ExternalDocs         *GetSchemaResultAdditionalPropertiesSchemaValueExternalDocs  `json:"externalDocs,omitempty"`
	Format               *string                                                      `json:"format,omitempty"`
	Items                *GetSchemaResultAdditionalPropertiesSchemaValueItems         `json:"items,omitempty"`
	Max                  *float64                                                     `json:"max,omitempty"`
	MaxItems             *int                                                         `json:"maxItems,omitempty"`
	MaxLength            *int                                                         `json:"maxLength,omitempty"`
	MaxProps             *int                                                         `json:"maxProps,omitempty"`
	Min                  *float64                                                     `json:"min,omitempty"`
	MinItems             *int                                                         `json:"minItems,omitempty"`
	MinLength            *int                                                         `json:"minLength,omitempty"`
	MinProps             *int                                                         `json:"minProps,omitempty"`
	MultipleOf           *float64                                                     `json:"multipleOf,omitempty"`
	Not                  *GetSchemaResultAdditionalPropertiesSchemaValueNot           `json:"not,omitempty"`
	Nullable             *bool                                                        `json:"nullable,omitempty"`
	OneOf                *GetSchemaResultAdditionalPropertiesSchemaValueOneOf         `json:"oneOf,omitempty"`
	Pattern              *string                                                      `json:"pattern,omitempty"`
	Properties           *GetSchemaResultAdditionalPropertiesSchemaValueProperties    `json:"properties,omitempty"`
	ReadOnly             *bool                                                        `json:"readOnly,omitempty"`
	Required             *GetSchemaResultAdditionalPropertiesSchemaValueRequired      `json:"required,omitempty"`
	Title                *string                                                      `json:"title,omitempty"`
	Type                 *string                                                      `json:"type,omitempty"`
	UniqueItems          *bool                                                        `json:"uniqueItems,omitempty"`
	WriteOnly            *bool                                                        `json:"writeOnly,omitempty"`
	Xml                  *GetSchemaResultAdditionalPropertiesSchemaValueXml           `json:"xml,omitempty"`
}

type GetSchemaResultAdditionalProperties struct {
	Has    *bool                                      `json:"has,omitempty"`
	Schema *GetSchemaResultAdditionalPropertiesSchema `json:"schema,omitempty"`
}

type GetSchemaResultAdditionalPropertiesSchema struct {
	Ref   *string          `json:"ref,omitempty"`
	Value *GetSchemaResult `json:"value,omitempty"`
}

type GetSchemaResultAdditionalPropertiesSchemaValueAllOfItem struct {
	Ref   *string          `json:"ref,omitempty"`
	Value *GetSchemaResult `json:"value,omitempty"`
}

type GetSchemaResultAdditionalPropertiesSchemaValueAllOf []*GetSchemaResultAdditionalPropertiesSchemaValueAllOfItem

type GetSchemaResultAdditionalPropertiesSchemaValueAnyOfItem struct {
	Ref   *string          `json:"ref,omitempty"`
	Value *GetSchemaResult `json:"value,omitempty"`
}

type GetSchemaResultAdditionalPropertiesSchemaValueAnyOf []*GetSchemaResultAdditionalPropertiesSchemaValueAnyOfItem

// any of: *bool, *string, *float64, *int, *GetSchemaResultAdditionalPropertiesSchemaValueDefault4, *GetSchemaResultAdditionalPropertiesSchemaValueDefault5
type GetSchemaResultAdditionalPropertiesSchemaValueDefault any

type GetSchemaResultAdditionalPropertiesSchemaValueDefault4 []any

type GetSchemaResultAdditionalPropertiesSchemaValueDefault5 struct {
}

type GetSchemaResultAdditionalPropertiesSchemaValueDiscriminator struct {
	Extensions   *GetSchemaResultAdditionalPropertiesSchemaValueDiscriminatorExtensions `json:"extensions,omitempty"`
	Mapping      *GetSchemaResultAdditionalPropertiesSchemaValueDiscriminatorMapping    `json:"mapping,omitempty"`
	PropertyName *string                                                                `json:"propertyName,omitempty"`
}

type GetSchemaResultAdditionalPropertiesSchemaValueDiscriminatorExtensions struct {
}

type GetSchemaResultAdditionalPropertiesSchemaValueDiscriminatorMapping struct {
}

// any of: *bool, *string, *float64, *int, *GetSchemaResultAdditionalPropertiesSchemaValueEnumItem4, *GetSchemaResultAdditionalPropertiesSchemaValueEnumItem5
type GetSchemaResultAdditionalPropertiesSchemaValueEnumItem any

type GetSchemaResultAdditionalPropertiesSchemaValueEnumItem4 []any

type GetSchemaResultAdditionalPropertiesSchemaValueEnumItem5 struct {
}

type GetSchemaResultAdditionalPropertiesSchemaValueEnum []*GetSchemaResultAdditionalPropertiesSchemaValueEnumItem

// any of: *bool, *string, *float64, *int, *GetSchemaResultAdditionalPropertiesSchemaValueExample4, *GetSchemaResultAdditionalPropertiesSchemaValueExample5
type GetSchemaResultAdditionalPropertiesSchemaValueExample any

type GetSchemaResultAdditionalPropertiesSchemaValueExample4 []any

type GetSchemaResultAdditionalPropertiesSchemaValueExample5 struct {
}

type GetSchemaResultAdditionalPropertiesSchemaValueExtensions struct {
}

type GetSchemaResultAdditionalPropertiesSchemaValueExternalDocs struct {
	Description *string                                                               `json:"description,omitempty"`
	Extensions  *GetSchemaResultAdditionalPropertiesSchemaValueExternalDocsExtensions `json:"extensions,omitempty"`
	Url         *string                                                               `json:"url,omitempty"`
}

type GetSchemaResultAdditionalPropertiesSchemaValueExternalDocsExtensions struct {
}

type GetSchemaResultAdditionalPropertiesSchemaValueItems struct {
	Ref   *string          `json:"ref,omitempty"`
	Value *GetSchemaResult `json:"value,omitempty"`
}

type GetSchemaResultAdditionalPropertiesSchemaValueNot struct {
	Ref   *string          `json:"ref,omitempty"`
	Value *GetSchemaResult `json:"value,omitempty"`
}

type GetSchemaResultAdditionalPropertiesSchemaValueOneOfItem struct {
	Ref   *string          `json:"ref,omitempty"`
	Value *GetSchemaResult `json:"value,omitempty"`
}

type GetSchemaResultAdditionalPropertiesSchemaValueOneOf []*GetSchemaResultAdditionalPropertiesSchemaValueOneOfItem

type GetSchemaResultAdditionalPropertiesSchemaValueProperties struct {
}

type GetSchemaResultAdditionalPropertiesSchemaValueRequired []string

type GetSchemaResultAdditionalPropertiesSchemaValueXml struct {
	Attribute  *bool                                                        `json:"attribute,omitempty"`
	Extensions *GetSchemaResultAdditionalPropertiesSchemaValueXmlExtensions `json:"extensions,omitempty"`
	Name       *string                                                      `json:"name,omitempty"`
	Namespace  *string                                                      `json:"namespace,omitempty"`
	Prefix     *string                                                      `json:"prefix,omitempty"`
	Wrapped    *bool                                                        `json:"wrapped,omitempty"`
}

type GetSchemaResultAdditionalPropertiesSchemaValueXmlExtensions struct {
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

// Context from caller is used to allow cancellation of long-running requests
func (s *service) GetSchemaContext(
	ctx context.Context,
	parameters GetSchemaParameters,
) (
	result GetSchemaResult,
	err error,
) {
	exec, ctx, err := mgcHelpers.PrepareExecutor("GetSchema", mgcCore.RefPath("/config/get-schema"), s.client, ctx)
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
