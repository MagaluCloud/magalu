package schema

import (
	"github.com/getkin/kin-openapi/openapi3"
	"magalu.cloud/core/utils"
)

// Copy-on-Write for SchemaRef
//
// All Setters are smart enough to understand whenever a copy is required or not
// There is no need to do it manually.
type COWSchemaRef struct {
	s       *SchemaRef
	changed bool
}

func equalSchema(a, b *Schema) bool {
	return utils.IsPointerEqualFunc(a, b, func(v1, v2 *Schema) bool {
		return a.Equals(b)
	})
}

func equalSchemaRef(a, b *SchemaRef) bool {
	return utils.IsPointerEqualFunc(a, b, func(v1, v2 *SchemaRef) bool {
		return equalSchema((*Schema)(a.Value), (*Schema)(b.Value))
	})
}

func NewCOWSchemaRef(s *SchemaRef) *COWSchemaRef {
	return &COWSchemaRef{s, false}
}

func (c *COWSchemaRef) copyIfNeeded() {
	if !c.changed {
		if c.s == nil {
			c.s = new(SchemaRef)
		} else {
			s := *c.s
			c.s = &s
		}
		c.changed = true
	}
}

func (c *COWSchemaRef) SetValue(v *Schema) {
	if equalSchema(c.Value(), v) {
		return
	}
	c.copyIfNeeded()
	c.s.Ref = ""
	c.s.Value = (*openapi3.Schema)(v)
}

func (c *COWSchemaRef) Value() *Schema {
	if c.s == nil {
		return nil
	}
	return (*Schema)(c.s.Value)
}

// Only does it if the schema references are not equal.
//
// The COWSchemaRef will be set as changed and other will be COPIED
func (c *COWSchemaRef) Replace(other *SchemaRef) {
	if equalSchemaRef(c.s, other) {
		return
	}
	c.changed = true
	s := *other
	c.s = &s
}

func (c *COWSchemaRef) Release() (s *SchemaRef, changed bool) {
	s = c.s
	changed = c.changed
	c.s = nil
	c.changed = false
	return s, changed
}

// Get the pointer to the internal reference.
//
// DO NOT MODIFY THE RETURNED SCHEMA
func (c *COWSchemaRef) Peek() *SchemaRef {
	return c.s
}

func (c *COWSchemaRef) IsChanged() bool {
	return c.changed
}

var _ utils.COW[*SchemaRef] = (*COWSchemaRef)(nil)

// Copy-on-Write for Schema
//
// All Setters are smart enough to understand whenever a copy is required or not
// There is no need to do it manually.
type COWSchema struct {
	s             *Schema
	changed       bool
	cowEnum       *utils.COWSlice[any]
	cowOneOf      *utils.COWSlice[*SchemaRef]
	cowAllOf      *utils.COWSlice[*SchemaRef]
	cowAnyOf      *utils.COWSlice[*SchemaRef]
	cowRequired   *utils.COWSlice[string]
	cowProperties *utils.COWMap[string, *SchemaRef]
	cowExtensions *utils.COWMap[string, any]
}

func (c *COWSchema) initCOW() {
	s := c.s
	if s == nil {
		s = &Schema{}
	}
	c.cowEnum = utils.NewCOWSliceFunc(s.Enum, utils.IsSameValueOrPointer)
	c.cowOneOf = utils.NewCOWSliceFunc(s.OneOf, equalSchemaRef)
	c.cowAllOf = utils.NewCOWSliceFunc(s.AllOf, equalSchemaRef)
	c.cowAnyOf = utils.NewCOWSliceFunc(s.AnyOf, equalSchemaRef)
	c.cowRequired = utils.NewCOWSliceComparable(s.Required)
	c.cowProperties = utils.NewCOWMapFunc(s.Properties, equalSchemaRef)
	c.cowExtensions = utils.NewCOWMapFunc(s.Extensions, utils.IsSameValueOrPointer)
}

func (c *COWSchema) isCOWChanged() bool {
	return c.cowEnum.IsChanged() ||
		c.cowOneOf.IsChanged() ||
		c.cowAllOf.IsChanged() ||
		c.cowAnyOf.IsChanged() ||
		c.cowRequired.IsChanged() ||
		c.cowProperties.IsChanged() ||
		c.cowExtensions.IsChanged()
}

// Sub COW are handled apart, but whenever we need to return the schema
// we must copy the schema if needed and then set all
// public pointers to the latest value of each COW
func (c *COWSchema) materializeCOW() {
	if !c.isCOWChanged() {
		return
	}
	c.copyIfNeeded()
	c.s.Enum = c.cowEnum.Peek()
	c.s.OneOf = c.cowOneOf.Peek()
	c.s.AllOf = c.cowAllOf.Peek()
	c.s.AnyOf = c.cowAnyOf.Peek()
	c.s.Required = c.cowRequired.Peek()
	c.s.Properties = c.cowProperties.Peek()
	c.s.Extensions = c.cowExtensions.Peek()
}

func NewCOWSchema(s *Schema) *COWSchema {
	c := &COWSchema{
		s:       s,
		changed: false,
	}
	c.initCOW()
	return c
}

func (c *COWSchema) copyIfNeeded() {
	if !c.changed {
		if c.s == nil {
			c.s = new(Schema)
		} else {
			s := *c.s
			c.s = &s
		}
		c.changed = true
	}
}

func (c *COWSchema) Extensions() map[string]any {
	return c.cowExtensions.Peek()
}

func (c *COWSchema) SetExtensions(v map[string]any) {
	c.cowExtensions.Replace(v)
}

func (c *COWSchema) SetExtension(name string, value any) {
	c.cowExtensions.Set(name, value)
}

func (c *COWSchema) GetExtension(name string) (any, bool) {
	return c.cowExtensions.Get(name)
}

func (c *COWSchema) Type() string {
	if c.s == nil {
		return ""
	}
	return c.s.Type
}

func (c *COWSchema) SetType(v string) {
	if c.Type() == v {
		return
	}
	c.copyIfNeeded()
	c.s.Type = v
}

func (c *COWSchema) Format() string {
	if c.s == nil {
		return ""
	}
	return c.s.Format
}

func (c *COWSchema) SetFormat(v string) {
	if c.Format() == v {
		return
	}
	c.copyIfNeeded()
	c.s.Format = v
}

func (c *COWSchema) Description() string {
	if c.s == nil {
		return ""
	}
	return c.s.Description
}

func (c *COWSchema) SetDescription(v string) {
	if c.Description() == v {
		return
	}
	c.copyIfNeeded()
	c.s.Description = v
}

func (c *COWSchema) Default() any {
	if c.s == nil {
		return nil
	}
	return c.s.Default
}

func (c *COWSchema) SetDefault(v any) {
	if utils.IsSameValueOrPointer(c.Default(), v) {
		return
	}
	c.copyIfNeeded()
	c.s.Default = v
}

func (c *COWSchema) Example() any {
	if c.s == nil {
		return nil
	}
	return c.s.Example
}

func (c *COWSchema) SetExample(v any) {
	if utils.IsSameValueOrPointer(c.Example(), v) {
		return
	}
	c.copyIfNeeded()
	c.s.Example = v
}

func (c *COWSchema) Enum() []any {
	return c.cowEnum.Peek()
}

func (c *COWSchema) SetEnum(v []any) {
	c.cowEnum.Replace(v)
}

// Checks if exists, otherwise append to Enum
func (c *COWSchema) AddEnum(v any) {
	c.cowEnum.Add(v)
}

func (c *COWSchema) OneOf() SchemaRefs {
	return c.cowOneOf.Peek()
}

func (c *COWSchema) SetOneOf(v SchemaRefs) {
	c.cowOneOf.Replace(v)
}

// Checks if exists, otherwise append to OneOf
func (c *COWSchema) AddOneOf(v *SchemaRef) {
	c.cowOneOf.Add(v)
}

func (c *COWSchema) AnyOf() SchemaRefs {
	return c.cowAnyOf.Peek()
}

func (c *COWSchema) SetAnyOf(v SchemaRefs) {
	c.cowAnyOf.Replace(v)
}

// Checks if exists, otherwise append to AnyOf
func (c *COWSchema) AddAnyOf(v *SchemaRef) {
	c.cowAnyOf.Add(v)
}

func (c *COWSchema) AllOf() SchemaRefs {
	return c.cowAllOf.Peek()
}

func (c *COWSchema) SetAllOf(v SchemaRefs) {
	c.cowAllOf.Replace(v)
}

// Checks if exists, otherwise append to AllOf
func (c *COWSchema) AddAllOf(v *SchemaRef) {
	c.cowAllOf.Add(v)
}

func (c *COWSchema) Not() *SchemaRef {
	if c.s == nil {
		return nil
	}
	return c.s.Not
}

func (c *COWSchema) SetNot(v *SchemaRef) {
	if equalSchemaRef(c.Not(), v) {
		return
	}
	c.copyIfNeeded()
	c.s.Not = v
}

// Array-related, here for struct compactness

func (c *COWSchema) UniqueItems() bool {
	if c.s == nil {
		return false
	}
	return c.s.UniqueItems
}

func (c *COWSchema) SetUniqueItems(v bool) {
	if c.UniqueItems() == v {
		return
	}
	c.copyIfNeeded()
	c.s.UniqueItems = v
}

// Number-related, here for struct compactness

func (c *COWSchema) ExclusiveMin() bool {
	if c.s == nil {
		return false
	}
	return c.s.ExclusiveMin
}

func (c *COWSchema) SetExclusiveMin(v bool) {
	if c.ExclusiveMin() == v {
		return
	}
	c.copyIfNeeded()
	c.s.ExclusiveMin = v
}

func (c *COWSchema) ExclusiveMax() bool {
	if c.s == nil {
		return false
	}
	return c.s.ExclusiveMax
}

func (c *COWSchema) SetExclusiveMax(v bool) {
	if c.ExclusiveMax() == v {
		return
	}
	c.copyIfNeeded()
	c.s.ExclusiveMax = v
}

// Properties

func (c *COWSchema) Nullable() bool {
	if c.s == nil {
		return false
	}
	return c.s.Nullable
}

func (c *COWSchema) SetNullable(v bool) {
	if c.Nullable() == v {
		return
	}
	c.copyIfNeeded()
	c.s.Nullable = v
}

func (c *COWSchema) ReadOnly() bool {
	if c.s == nil {
		return false
	}
	return c.s.ReadOnly
}

func (c *COWSchema) SetReadOnly(v bool) {
	if c.ReadOnly() == v {
		return
	}
	c.copyIfNeeded()
	c.s.ReadOnly = v
}

func (c *COWSchema) WriteOnly() bool {
	if c.s == nil {
		return false
	}
	return c.s.WriteOnly
}

func (c *COWSchema) SetWriteOnly(v bool) {
	if c.WriteOnly() == v {
		return
	}
	c.copyIfNeeded()
	c.s.WriteOnly = v
}

func (c *COWSchema) AllowEmptyValue() bool {
	if c.s == nil {
		return false
	}
	return c.s.AllowEmptyValue
}

func (c *COWSchema) SetAllowEmptyValue(v bool) {
	if c.AllowEmptyValue() == v {
		return
	}
	c.copyIfNeeded()
	c.s.AllowEmptyValue = v
}

func (c *COWSchema) Deprecated() bool {
	if c.s == nil {
		return false
	}
	return c.s.Deprecated
}

func (c *COWSchema) SetDeprecated(v bool) {
	if c.Deprecated() == v {
		return
	}
	c.copyIfNeeded()
	c.s.Deprecated = v
}

// Number

func (c *COWSchema) Min() *float64 {
	if c.s == nil {
		return nil
	}
	return c.s.Min
}

func (c *COWSchema) SetMin(v *float64) {
	if utils.IsComparablePointerEqual(c.Min(), v) {
		return
	}
	c.copyIfNeeded()
	c.s.Min = v
}

func (c *COWSchema) Max() *float64 {
	if c.s == nil {
		return nil
	}
	return c.s.Max
}

func (c *COWSchema) SetMax(v *float64) {
	if utils.IsComparablePointerEqual(c.Max(), v) {
		return
	}
	c.copyIfNeeded()
	c.s.Max = v
}

func (c *COWSchema) MultipleOf() *float64 {
	if c.s == nil {
		return nil
	}
	return c.s.MultipleOf
}

func (c *COWSchema) SetMultipleOf(v *float64) {
	if utils.IsComparablePointerEqual(c.MultipleOf(), v) {
		return
	}
	c.copyIfNeeded()
	c.s.MultipleOf = v
}

// String

func (c *COWSchema) MinLength() uint64 {
	if c.s == nil {
		return 0
	}
	return c.s.MinLength
}

func (c *COWSchema) SetMinLength(v uint64) {
	if c.MinLength() == v {
		return
	}
	c.copyIfNeeded()
	c.s.MinLength = v
}

func (c *COWSchema) MaxLength() *uint64 {
	if c.s == nil {
		return nil
	}
	return c.s.MaxLength
}

func (c *COWSchema) SetMaxLength(v *uint64) {
	if utils.IsComparablePointerEqual(c.MaxLength(), v) {
		return
	}
	c.copyIfNeeded()
	c.s.MaxLength = v
}

func (c *COWSchema) Pattern() string {
	if c.s == nil {
		return ""
	}
	return c.s.Pattern
}

func (c *COWSchema) SetPattern(v string) {
	if c.Pattern() == v {
		return
	}
	c.copyIfNeeded()
	(*openapi3.Schema)(c.s).WithPattern(v) // resets compiledPattern
}

// Array

func (c *COWSchema) MinItems() uint64 {
	if c.s == nil {
		return 0
	}
	return c.s.MinItems
}

func (c *COWSchema) SetMinItems(v uint64) {
	if c.MinItems() == v {
		return
	}
	c.copyIfNeeded()
	c.s.MinItems = v
}

func (c *COWSchema) MaxItems() *uint64 {
	if c.s == nil {
		return nil
	}
	return c.s.MaxItems
}

func (c *COWSchema) SetMaxItems(v *uint64) {
	if utils.IsComparablePointerEqual(c.MaxItems(), v) {
		return
	}
	c.copyIfNeeded()
	c.s.MaxItems = v
}

func (c *COWSchema) Items() *SchemaRef {
	if c.s == nil {
		return nil
	}
	return c.s.Items
}

func (c *COWSchema) SetItems(v *SchemaRef) {
	if equalSchemaRef(c.Items(), v) {
		return
	}
	c.copyIfNeeded()
	c.s.Items = v
}

// Object
func (c *COWSchema) Properties() map[string]*SchemaRef {
	return c.cowProperties.Peek()
}

func (c *COWSchema) SetProperties(v map[string]*SchemaRef) {
	c.cowProperties.Replace(v)
}

func (c *COWSchema) SetProperty(name string, schemaRef *SchemaRef) {
	c.cowProperties.Set(name, schemaRef)
}

func (c *COWSchema) GetProperty(name string) (*SchemaRef, bool) {
	return c.cowProperties.Get(name)
}

func (c *COWSchema) Required() []string {
	return c.cowRequired.Peek()
}

func (c *COWSchema) SetRequired(v []string) {
	c.cowRequired.Replace(v)
}

// Checks if exists, otherwise append to Required
func (c *COWSchema) AddRequired(v string) {
	c.cowRequired.Add(v)
}

func (c *COWSchema) MinProps() uint64 {
	if c.s == nil {
		return 0
	}
	return c.s.MinProps
}

func (c *COWSchema) SetMinProps(v uint64) {
	if c.MinProps() == v {
		return
	}
	c.copyIfNeeded()
	c.s.MinProps = v
}

func (c *COWSchema) MaxProps() *uint64 {
	if c.s == nil {
		return nil
	}
	return c.s.MaxProps
}

func (c *COWSchema) SetMaxProps(v *uint64) {
	if utils.IsComparablePointerEqual(c.MaxProps(), v) {
		return
	}
	c.copyIfNeeded()
	c.s.MaxProps = v
}

func (c *COWSchema) AdditionalProperties() openapi3.AdditionalProperties {
	if c.s == nil {
		return openapi3.AdditionalProperties{}
	}
	return c.s.AdditionalProperties
}

func (c *COWSchema) SetAdditionalProperties(v openapi3.AdditionalProperties) {
	existing := c.AdditionalProperties()
	if utils.IsComparablePointerEqual(existing.Has, v.Has) && equalSchemaRef(existing.Schema, v.Schema) {
		return
	}
	c.copyIfNeeded()
	c.s.AdditionalProperties = v
}

// Only does it if the schemas are not equal.
//
// The COWSchema will be set as changed and other will be COPIED
func (c *COWSchema) Replace(other *Schema) {
	if equalSchema(c.s, other) {
		return
	}
	c.changed = true
	s := *other
	c.s = &s
	c.initCOW()
}

func (c *COWSchema) Release() (s *Schema, changed bool) {
	c.materializeCOW()
	s = c.s
	changed = c.changed
	c.s = nil
	c.changed = false
	c.initCOW()
	return s, changed
}

// Get the pointer to the internal reference.
//
// DO NOT MODIFY THE RETURNED SCHEMA
func (c *COWSchema) Peek() *Schema {
	c.materializeCOW()
	return c.s
}

func (c *COWSchema) IsChanged() bool {
	return c.changed || c.isCOWChanged()
}

var _ utils.COW[*Schema] = (*COWSchema)(nil)
