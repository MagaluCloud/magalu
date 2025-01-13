package schema

import (
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/getkin/kin-openapi/openapi3"
)

// NOTE: COW/COWContainer/COWContainerOfCOW allows some nil receivers, so not doing explicit nil checks in here

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

// Copy-on-Write for SchemaRef
//
// All Setters are smart enough to understand whenever a copy is required or not
// There is no need to do it manually.
type COWSchemaRef struct {
	s        *SchemaRef
	changed  bool
	cowValue *COWSchema
}

func (c *COWSchemaRef) initCOWValue() {
	var value *Schema = nil
	if c.s != nil {
		value = (*Schema)(c.s.Value)
	}
	c.cowValue = NewCOWSchema(value)
}

func (c *COWSchemaRef) initCOWValueIfNeeded() {
	if c.cowValue == nil {
		c.initCOWValue()
	}
}

func (c *COWSchemaRef) initCOW() {
	if c.s == nil {
		c.cowValue = nil
	} else {
		c.initCOWValue()
	}
}

func (c *COWSchemaRef) isCOWChanged() bool {
	// nil receiver is ok
	return c.cowValue.IsChanged()
}

// Sub COW are handled apart, but whenever we need to return the schema
// we must copy the schema if needed and then set all
// public pointers to the latest value of each COW
func (c *COWSchemaRef) materializeCOW() {
	if !c.isCOWChanged() {
		return
	}
	c.copyIfNeeded()
	c.s.Value = (*openapi3.Schema)(c.cowValue.Peek())
}

func NewCOWSchemaRef(s *SchemaRef) *COWSchemaRef {
	c := &COWSchemaRef{
		s:       s,
		changed: false,
	}
	return c
}

func (c *COWSchemaRef) copyIfNeeded() {
	if !c.changed {
		if c.s == nil {
			c.s = new(SchemaRef)
		} else {
			s := *c.s
			c.s = &s
		}

		c.initCOWValueIfNeeded()

		c.changed = true
	}
}

func (c *COWSchemaRef) Ref() string {
	if c == nil || c.s == nil {
		return ""
	}
	return c.s.Ref
}

func (c *COWSchemaRef) SetRef(v string) bool {
	if c.Ref() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.Ref = v
	return true
}

func (c *COWSchemaRef) UnsetRef() bool {
	return c.SetRef("")
}

func (c *COWSchemaRef) SetValue(v *Schema) bool {
	c.initCOWValueIfNeeded()
	return c.cowValue.Replace(v)
}

func (c *COWSchemaRef) Value() *Schema {
	if c == nil {
		return nil
	}
	if c.cowValue != nil {
		return c.cowValue.Peek()
	}
	if c.s != nil {
		return (*Schema)(c.s.Value)
	}
	return nil
}

func (c *COWSchemaRef) ValueCOW() *COWSchema {
	c.initCOWValueIfNeeded()
	return c.cowValue
}

func (c *COWSchemaRef) Equals(other *SchemaRef) bool {
	if c == nil {
		return other == nil
	}
	return equalSchemaRef(c.Peek(), other)
}

// Only does it if the schema references are not equal.
//
// The COWSchemaRef will be set as changed and other will be COPIED
func (c *COWSchemaRef) Replace(other *SchemaRef) bool {
	if c.Equals(other) {
		return false
	}
	c.changed = true
	s := *other
	c.s = &s
	c.initCOW()
	return true
}

func (c *COWSchemaRef) Release() (s *SchemaRef, changed bool) {
	if c == nil {
		return
	}
	s = c.Peek()
	changed = c.IsChanged()
	c.s = nil
	c.changed = false
	c.initCOW()
	return s, changed
}

// Get the pointer to the internal reference.
//
// DO NOT MODIFY THE RETURNED SCHEMA
func (c *COWSchemaRef) Peek() (s *SchemaRef) {
	if c == nil {
		return
	}
	c.materializeCOW()
	return c.s
}

func (c *COWSchemaRef) IsChanged() (changed bool) {
	if c == nil {
		return
	}
	return c.changed || c.isCOWChanged()
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
	cowOneOf      *utils.COWSliceOfCOW[*SchemaRef, *COWSchemaRef]
	cowAllOf      *utils.COWSliceOfCOW[*SchemaRef, *COWSchemaRef]
	cowAnyOf      *utils.COWSliceOfCOW[*SchemaRef, *COWSchemaRef]
	cowNot        *COWSchemaRef
	cowItems      *COWSchemaRef
	cowRequired   *utils.COWSlice[string]
	cowProperties *utils.COWMapOfCOW[string, *SchemaRef, *COWSchemaRef]
	cowExtensions *utils.COWMap[string, any]
}

func (c *COWSchema) initCOWEnum() {
	var value []any = nil
	if c.s != nil {
		value = c.s.Enum
	}
	c.cowEnum = utils.NewCOWSliceFunc(value, utils.IsSameValueOrPointer)
}

func (c *COWSchema) initCOWEnumIfNeeded() {
	if c.cowEnum == nil {
		c.initCOWEnum()
	}
}

func (c *COWSchema) initCOWOneOf() {
	var value []*SchemaRef = nil
	if c.s != nil {
		value = c.s.OneOf
	}
	c.cowOneOf = utils.NewCOWSliceOfCOW(value, NewCOWSchemaRef)
}

func (c *COWSchema) initCOWOneOfIfNeeded() {
	if c.cowOneOf == nil {
		c.initCOWOneOf()
	}
}

func (c *COWSchema) initCOWAllOf() {
	var value []*SchemaRef = nil
	if c.s != nil {
		value = c.s.AllOf
	}
	c.cowAllOf = utils.NewCOWSliceOfCOW(value, NewCOWSchemaRef)
}

func (c *COWSchema) initCOWAllOfIfNeeded() {
	if c.cowAllOf == nil {
		c.initCOWAllOf()
	}
}

func (c *COWSchema) initCOWAnyOf() {
	var value []*SchemaRef = nil
	if c.s != nil {
		value = c.s.AnyOf
	}
	c.cowAnyOf = utils.NewCOWSliceOfCOW(value, NewCOWSchemaRef)
}

func (c *COWSchema) initCOWAnyOfIfNeeded() {
	if c.cowAnyOf == nil {
		c.initCOWAnyOf()
	}
}

func (c *COWSchema) initCOWNot() {
	var value *SchemaRef = nil
	if c.s != nil {
		value = c.s.Not
	}
	c.cowNot = NewCOWSchemaRef(value)
}

func (c *COWSchema) initCOWNotIfNeeded() {
	if c.cowNot == nil {
		c.initCOWNot()
	}
}

func (c *COWSchema) initCOWItems() {
	var value *SchemaRef = nil
	if c.s != nil {
		value = c.s.Items
	}
	c.cowItems = NewCOWSchemaRef(value)
}

func (c *COWSchema) initCOWItemsIfNeeded() {
	if c.cowItems == nil {
		c.initCOWItems()
	}
}

func (c *COWSchema) initCOWRequired() {
	var value []string = nil
	if c.s != nil {
		value = c.s.Required
	}
	c.cowRequired = utils.NewCOWSliceComparable(value)
}

func (c *COWSchema) initCOWRequiredIfNeeded() {
	if c.cowRequired == nil {
		c.initCOWRequired()
	}
}

func (c *COWSchema) initCOWProperties() {
	var value map[string]*SchemaRef = nil
	if c.s != nil {
		value = c.s.Properties
	}
	c.cowProperties = utils.NewCOWMapOfCOW(value, NewCOWSchemaRef)
}

func (c *COWSchema) initCOWPropertiesIfNeeded() {
	if c.cowProperties == nil {
		c.initCOWProperties()
	}
}

func (c *COWSchema) initCOWExtensions() {
	var value map[string]any = nil
	if c.s != nil {
		value = c.s.Extensions
	}
	c.cowExtensions = utils.NewCOWMapFunc(value, utils.IsSameValueOrPointer)
}

func (c *COWSchema) initCOWExtensionsIfNeeded() {
	if c.cowExtensions == nil {
		c.initCOWExtensions()
	}
}

func (c *COWSchema) initCOW() {
	if c.s == nil {
		c.cowEnum = nil
		c.cowOneOf = nil
		c.cowAllOf = nil
		c.cowAnyOf = nil
		c.cowNot = nil
		c.cowItems = nil
		c.cowRequired = nil
		c.cowProperties = nil
		c.cowExtensions = nil
	} else {
		c.initCOWEnum()
		c.initCOWOneOf()
		c.initCOWAllOf()
		c.initCOWAnyOf()
		c.initCOWNot()
		c.initCOWItems()
		c.initCOWRequired()
		c.initCOWProperties()
		c.initCOWExtensions()
	}
}

func (c *COWSchema) isCOWChanged() bool {
	// nil receivers are ok
	return (c.cowEnum.IsChanged() ||
		c.cowOneOf.IsChanged() ||
		c.cowAllOf.IsChanged() ||
		c.cowAnyOf.IsChanged() ||
		c.cowNot.IsChanged() ||
		c.cowItems.IsChanged() ||
		c.cowRequired.IsChanged() ||
		c.cowProperties.IsChanged() ||
		c.cowExtensions.IsChanged())
}

// Sub COW are handled apart, but whenever we need to return the schema
// we must copy the schema if needed and then set all
// public pointers to the latest value of each COW
func (c *COWSchema) materializeCOW() {
	if !c.isCOWChanged() {
		return
	}
	c.copyIfNeeded()
	// nil receivers are ok
	c.s.Enum = c.cowEnum.Peek()
	c.s.OneOf = c.cowOneOf.Peek()
	c.s.AllOf = c.cowAllOf.Peek()
	c.s.AnyOf = c.cowAnyOf.Peek()
	c.s.Not = c.cowNot.Peek()
	c.s.Items = c.cowItems.Peek()
	c.s.Required = c.cowRequired.Peek()
	c.s.Properties = c.cowProperties.Peek()
	c.s.Extensions = c.cowExtensions.Peek()
}

func NewCOWSchema(s *Schema) *COWSchema {
	c := &COWSchema{
		s:       s,
		changed: false,
	}
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

		c.initCOWEnumIfNeeded()
		c.initCOWOneOfIfNeeded()
		c.initCOWAllOfIfNeeded()
		c.initCOWAnyOfIfNeeded()
		c.initCOWNotIfNeeded()
		c.initCOWItemsIfNeeded()
		c.initCOWRequiredIfNeeded()
		c.initCOWPropertiesIfNeeded()
		c.initCOWExtensionsIfNeeded()

		c.changed = true
	}
}

func (c *COWSchema) ExtensionsCOW() *utils.COWMap[string, any] {
	c.initCOWExtensionsIfNeeded()
	return c.cowExtensions
}

// Do not mutate the returned handle, see ExtensionsCOW() if you want to mutate
func (c *COWSchema) Extensions() map[string]any {
	if c == nil {
		return nil
	}
	if c.cowExtensions != nil {
		return c.cowExtensions.Peek()
	}
	if c.s != nil {
		return c.s.Extensions
	}
	return nil
}

// Replace the Extensions map, if it's different from the existing value
//
// In order to do more fine grained operations such as Set/Delete, use ExtensionsCOW()
func (c *COWSchema) SetExtensions(v map[string]any) bool {
	c.initCOWExtensionsIfNeeded()
	return c.cowExtensions.Replace(v)
}

func (c *COWSchema) Type() string {
	if c == nil || c.s == nil {
		return ""
	}
	return c.s.Type
}

func (c *COWSchema) SetType(v string) bool {
	if c.Type() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.Type = v
	return true
}

func (c *COWSchema) Format() string {
	if c == nil || c.s == nil {
		return ""
	}
	return c.s.Format
}

func (c *COWSchema) SetFormat(v string) bool {
	if c.Format() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.Format = v
	return true
}

func (c *COWSchema) Description() string {
	if c == nil || c.s == nil {
		return ""
	}
	return c.s.Description
}

func (c *COWSchema) SetDescription(v string) bool {
	if c.Description() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.Description = v
	return true
}

func (c *COWSchema) Default() any {
	if c == nil || c.s == nil {
		return nil
	}
	return c.s.Default
}

func (c *COWSchema) SetDefault(v any) bool {
	if utils.IsSameValueOrPointer(c.Default(), v) {
		return false
	}
	c.copyIfNeeded()
	c.s.Default = v
	return true
}

func (c *COWSchema) Example() any {
	if c == nil || c.s == nil {
		return nil
	}
	return c.s.Example
}

func (c *COWSchema) SetExample(v any) bool {
	if utils.IsSameValueOrPointer(c.Example(), v) {
		return false
	}
	c.copyIfNeeded()
	c.s.Example = v
	return true
}

func (c *COWSchema) EnumCOW() *utils.COWSlice[any] {
	c.initCOWEnumIfNeeded()
	return c.cowEnum
}

// Do not mutate the returned handle, see EnumCOW() if you want to mutate
func (c *COWSchema) Enum() []any {
	if c == nil {
		return nil
	}
	if c.cowEnum != nil {
		return c.cowEnum.Peek()
	}
	if c.s != nil {
		return c.s.Enum
	}
	return nil
}

// Replace the Enum slice, if it's different from the existing value
//
// In order to do more fine grained operations such as Add/Append, use EnumCOW()
func (c *COWSchema) SetEnum(v []any) bool {
	c.initCOWEnumIfNeeded()
	return c.cowEnum.Replace(v)
}

func (c *COWSchema) OneOfCOW() *utils.COWSliceOfCOW[*SchemaRef, *COWSchemaRef] {
	c.initCOWOneOfIfNeeded()
	return c.cowOneOf
}

// Do not mutate the returned handle, see OneOfCOW() if you want to mutate
func (c *COWSchema) OneOf() SchemaRefs {
	if c == nil {
		return nil
	}
	if c.cowOneOf != nil {
		return c.cowOneOf.Peek()
	}
	if c.s != nil {
		return c.s.OneOf
	}
	return nil
}

// Replace the OneOf slice, if it's different from the existing value
//
// In order to do more fine grained operations such as Add/Append, use OneOfCOW()
func (c *COWSchema) SetOneOf(v SchemaRefs) bool {
	c.initCOWOneOfIfNeeded()
	return c.cowOneOf.Replace(v)
}

func (c *COWSchema) AnyOfCOW() *utils.COWSliceOfCOW[*SchemaRef, *COWSchemaRef] {
	c.initCOWAnyOfIfNeeded()
	return c.cowAnyOf
}

// Do not mutate the returned handle, see AnyOfCOW() if you want to mutate
func (c *COWSchema) AnyOf() SchemaRefs {
	if c == nil {
		return nil
	}
	if c.cowAnyOf != nil {
		return c.cowAnyOf.Peek()
	}
	if c.s != nil {
		return c.s.AnyOf
	}
	return nil
}

// Replace the AnyOf slice, if it's different from the existing value
//
// In order to do more fine grained operations such as Add/Append, use AnyOfCOW()
func (c *COWSchema) SetAnyOf(v SchemaRefs) bool {
	c.initCOWAnyOfIfNeeded()
	return c.cowAnyOf.Replace(v)
}

func (c *COWSchema) AllOfCOW() *utils.COWSliceOfCOW[*SchemaRef, *COWSchemaRef] {
	c.initCOWAllOfIfNeeded()
	return c.cowAllOf
}

// Do not mutate the returned handle, see AllOfCOW() if you want to mutate
func (c *COWSchema) AllOf() SchemaRefs {
	if c == nil {
		return nil
	}
	if c.cowAllOf != nil {
		return c.cowAllOf.Peek()
	}
	if c.s != nil {
		return c.s.AllOf
	}
	return nil
}

// Replace the AllOf slice, if it's different from the existing value
//
// In order to do more fine grained operations such as Add/Append, use AllCOW()
func (c *COWSchema) SetAllOf(v SchemaRefs) bool {
	c.initCOWAllOfIfNeeded()
	return c.cowAllOf.Replace(v)
}

func (c *COWSchema) NotCOW() *COWSchemaRef {
	c.initCOWNotIfNeeded()
	return c.cowNot
}

// Do not mutate the returned handle, see NotCOW() if you want to mutate
func (c *COWSchema) Not() *SchemaRef {
	if c == nil {
		return nil
	}
	if c.cowNot != nil {
		return c.cowNot.Peek()
	}
	if c.s != nil {
		return c.s.Not
	}
	return nil
}

func (c *COWSchema) SetNot(v *SchemaRef) bool {
	c.initCOWNotIfNeeded()
	return c.cowNot.Replace(v)
}

// Array-related, here for struct compactness

func (c *COWSchema) UniqueItems() bool {
	if c == nil || c.s == nil {
		return false
	}
	return c.s.UniqueItems
}

func (c *COWSchema) SetUniqueItems(v bool) bool {
	if c.UniqueItems() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.UniqueItems = v
	return true
}

// Number-related, here for struct compactness

func (c *COWSchema) ExclusiveMin() bool {
	if c == nil || c.s == nil {
		return false
	}
	return c.s.ExclusiveMin
}

func (c *COWSchema) SetExclusiveMin(v bool) bool {
	if c.ExclusiveMin() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.ExclusiveMin = v
	return true
}

func (c *COWSchema) ExclusiveMax() bool {
	if c == nil || c.s == nil {
		return false
	}
	return c.s.ExclusiveMax
}

func (c *COWSchema) SetExclusiveMax(v bool) bool {
	if c.ExclusiveMax() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.ExclusiveMax = v
	return true
}

// Properties

func (c *COWSchema) Nullable() bool {
	if c == nil || c.s == nil {
		return false
	}
	return c.s.Nullable
}

func (c *COWSchema) SetNullable(v bool) bool {
	if c.Nullable() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.Nullable = v
	return true
}

func (c *COWSchema) ReadOnly() bool {
	if c == nil || c.s == nil {
		return false
	}
	return c.s.ReadOnly
}

func (c *COWSchema) SetReadOnly(v bool) bool {
	if c.ReadOnly() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.ReadOnly = v
	return true
}

func (c *COWSchema) WriteOnly() bool {
	if c == nil || c.s == nil {
		return false
	}
	return c.s.WriteOnly
}

func (c *COWSchema) SetWriteOnly(v bool) bool {
	if c.WriteOnly() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.WriteOnly = v
	return true
}

func (c *COWSchema) AllowEmptyValue() bool {
	if c == nil || c.s == nil {
		return false
	}
	return c.s.AllowEmptyValue
}

func (c *COWSchema) SetAllowEmptyValue(v bool) bool {
	if c.AllowEmptyValue() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.AllowEmptyValue = v
	return true
}

func (c *COWSchema) Deprecated() bool {
	if c == nil || c.s == nil {
		return false
	}
	return c.s.Deprecated
}

func (c *COWSchema) SetDeprecated(v bool) bool {
	if c.Deprecated() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.Deprecated = v
	return true
}

// Number

func (c *COWSchema) Min() *float64 {
	if c == nil || c.s == nil {
		return nil
	}
	return c.s.Min
}

func (c *COWSchema) SetMin(v *float64) bool {
	if utils.IsComparablePointerEqual(c.Min(), v) {
		return false
	}
	c.copyIfNeeded()
	c.s.Min = v
	return true
}

func (c *COWSchema) Max() *float64 {
	if c == nil || c.s == nil {
		return nil
	}
	return c.s.Max
}

func (c *COWSchema) SetMax(v *float64) bool {
	if utils.IsComparablePointerEqual(c.Max(), v) {
		return false
	}
	c.copyIfNeeded()
	c.s.Max = v
	return true
}

func (c *COWSchema) MultipleOf() *float64 {
	if c == nil || c.s == nil {
		return nil
	}
	return c.s.MultipleOf
}

func (c *COWSchema) SetMultipleOf(v *float64) bool {
	if utils.IsComparablePointerEqual(c.MultipleOf(), v) {
		return false
	}
	c.copyIfNeeded()
	c.s.MultipleOf = v
	return true
}

// String

func (c *COWSchema) MinLength() uint64 {
	if c == nil || c.s == nil {
		return 0
	}
	return c.s.MinLength
}

func (c *COWSchema) SetMinLength(v uint64) bool {
	if c.MinLength() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.MinLength = v
	return true
}

func (c *COWSchema) MaxLength() *uint64 {
	if c == nil || c.s == nil {
		return nil
	}
	return c.s.MaxLength
}

func (c *COWSchema) SetMaxLength(v *uint64) bool {
	if utils.IsComparablePointerEqual(c.MaxLength(), v) {
		return false
	}
	c.copyIfNeeded()
	c.s.MaxLength = v
	return true
}

func (c *COWSchema) Pattern() string {
	if c == nil || c.s == nil {
		return ""
	}
	return c.s.Pattern
}

func (c *COWSchema) SetPattern(v string) bool {
	if c.Pattern() == v {
		return false
	}
	c.copyIfNeeded()
	(*openapi3.Schema)(c.s).WithPattern(v) // resets compiledPattern
	return true
}

// Array

func (c *COWSchema) MinItems() uint64 {
	if c == nil || c.s == nil {
		return 0
	}
	return c.s.MinItems
}

func (c *COWSchema) SetMinItems(v uint64) bool {
	if c.MinItems() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.MinItems = v
	return true
}

func (c *COWSchema) MaxItems() *uint64 {
	if c == nil || c.s == nil {
		return nil
	}
	return c.s.MaxItems
}

func (c *COWSchema) SetMaxItems(v *uint64) bool {
	if utils.IsComparablePointerEqual(c.MaxItems(), v) {
		return false
	}
	c.copyIfNeeded()
	c.s.MaxItems = v
	return true
}

func (c *COWSchema) ItemsCOW() *COWSchemaRef {
	c.initCOWItemsIfNeeded()
	return c.cowItems
}

func (c *COWSchema) Items() *SchemaRef {
	if c == nil {
		return nil
	}
	if c.cowItems != nil {
		return c.cowItems.Peek()
	}
	if c.s != nil {
		return c.s.Items
	}
	return nil
}

func (c *COWSchema) SetItems(v *SchemaRef) bool {
	c.initCOWItemsIfNeeded()
	return c.cowItems.Replace(v)
}

// Object
func (c *COWSchema) PropertiesCOW() *utils.COWMapOfCOW[string, *SchemaRef, *COWSchemaRef] {
	c.initCOWPropertiesIfNeeded()
	return c.cowProperties
}

// Do not mutate the returned handle, see PropertiesCOW() if you want to mutate
func (c *COWSchema) Properties() map[string]*SchemaRef {
	if c == nil {
		return nil
	}
	if c.cowProperties != nil {
		return c.cowProperties.Peek()
	}
	if c.s != nil {
		return c.s.Properties
	}
	return nil
}

// In order to do more fine grained operations such as Set/Delete, use PropertiesCOW()
func (c *COWSchema) SetProperties(v map[string]*SchemaRef) bool {
	c.initCOWPropertiesIfNeeded()
	return c.cowProperties.Replace(v)
}

func (c *COWSchema) RequiredCOW() *utils.COWSlice[string] {
	c.initCOWRequiredIfNeeded()
	return c.cowRequired
}

// Do not mutate the returned handle, see RequiredCOW() if you want to mutate
func (c *COWSchema) Required() []string {
	if c == nil {
		return nil
	}
	if c.cowRequired != nil {
		return c.cowRequired.Peek()
	}
	if c.s != nil {
		return c.s.Required
	}
	return nil
}

// Replace the Required slice, if it's different from the existing value
//
// In order to do more fine grained operations such as Add/Append, use RequiredCOW()
func (c *COWSchema) SetRequired(v []string) bool {
	c.initCOWRequiredIfNeeded()
	return c.cowRequired.Replace(v)
}

func (c *COWSchema) MinProps() uint64 {
	if c == nil || c.s == nil {
		return 0
	}
	return c.s.MinProps
}

func (c *COWSchema) SetMinProps(v uint64) bool {
	if c.MinProps() == v {
		return false
	}
	c.copyIfNeeded()
	c.s.MinProps = v
	return true
}

func (c *COWSchema) MaxProps() *uint64 {
	if c == nil || c.s == nil {
		return nil
	}
	return c.s.MaxProps
}

func (c *COWSchema) SetMaxProps(v *uint64) bool {
	if utils.IsComparablePointerEqual(c.MaxProps(), v) {
		return false
	}
	c.copyIfNeeded()
	c.s.MaxProps = v
	return true
}

func (c *COWSchema) AdditionalProperties() openapi3.AdditionalProperties {
	if c == nil || c.s == nil {
		return openapi3.AdditionalProperties{}
	}
	return c.s.AdditionalProperties
}

func (c *COWSchema) SetAdditionalProperties(v openapi3.AdditionalProperties) bool {
	existing := c.AdditionalProperties()
	if utils.IsComparablePointerEqual(existing.Has, v.Has) && equalSchemaRef(existing.Schema, v.Schema) {
		return false
	}
	c.copyIfNeeded()
	c.s.AdditionalProperties = v
	return true
}

func (c *COWSchema) Equals(other *Schema) bool {
	if c == nil {
		return other == nil
	}
	return equalSchema(c.Peek(), other)
}

// Only does it if the schemas are not equal.
//
// The COWSchema will be set as changed and other will be COPIED
func (c *COWSchema) Replace(other *Schema) bool {
	if c.Equals(other) {
		return false
	}
	c.changed = true
	s := *other
	c.s = &s
	c.initCOW()
	return true
}

func (c *COWSchema) Release() (s *Schema, changed bool) {
	if c == nil {
		return
	}
	s = c.Peek()
	changed = c.IsChanged()
	c.s = nil
	c.changed = false
	c.initCOW()
	return s, changed
}

// Get the pointer to the internal reference.
//
// DO NOT MODIFY THE RETURNED SCHEMA
func (c *COWSchema) Peek() (s *Schema) {
	if c == nil {
		return
	}
	c.materializeCOW()
	return c.s
}

func (c *COWSchema) IsChanged() (changed bool) {
	if c == nil {
		return
	}
	return c.changed || c.isCOWChanged()
}

var _ utils.COW[*Schema] = (*COWSchema)(nil)
