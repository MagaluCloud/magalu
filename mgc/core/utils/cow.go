package utils

import (
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// Copy-on-Write are helper to manage a mutable value, generating a copy whenever it will be written the first time
// If no writes are done, then the original object is never copied, saving resources

type COW[T any] interface {
	Replace(other T)
	Release() (value T, changed bool)
	Peek() T
	IsChanged() bool
}

// Map that will be copied before its first write, the original value is used if no writes were done
type COWMap[K comparable, V any] struct {
	m       map[K]V
	changed bool
	// How to compare values of the map
	equals func(V, V) bool
}

func NewCOWMapFunc[K comparable, V any](m map[K]V, equals func(V, V) bool) *COWMap[K, V] {
	return &COWMap[K, V]{m, false, equals}
}

func NewCOWMapComparable[K comparable, V comparable](m map[K]V) *COWMap[K, V] {
	return &COWMap[K, V]{m, false, IsComparableEqual[V]}
}

func (c *COWMap[K, V]) Get(key K) (V, bool) {
	v, ok := c.m[key]
	return v, ok
}

func (c *COWMap[K, V]) Len() int {
	return len(c.m)
}

func (c *COWMap[K, V]) copyIfNeeded() {
	if !c.changed {
		if c.m == nil {
			c.m = make(map[K]V)
		} else {
			c.m = maps.Clone(c.m)
		}
		c.changed = true
	}
}

func (c *COWMap[K, V]) Equals(other map[K]V) bool {
	if c.equals == nil {
		return false
	}
	return maps.EqualFunc(c.m, other, c.equals)
}

// Checks if the given value exists at the target position
func (c *COWMap[K, V]) ExistsAs(key K, value V) bool {
	if c.equals == nil {
		return false
	}
	if existing, ok := c.m[key]; ok {
		return c.equals(existing, value)
	}
	return false
}

func (c *COWMap[K, V]) Set(key K, value V) {
	if c.ExistsAs(key, value) {
		return
	}

	c.copyIfNeeded()
	c.m[key] = value
}

// Only does it if the maps are not equal.
//
// The COWMap will be set as changed and other will be COPIED
func (c *COWMap[K, V]) Replace(other map[K]V) {
	if c.Equals(other) {
		return
	}
	c.changed = true
	c.m = maps.Clone(other)
}

func (c *COWMap[K, V]) Release() (m map[K]V, changed bool) {
	m = c.m
	changed = c.changed
	c.m = nil
	c.changed = false
	return m, changed
}

// Get the pointer to the internal reference.
//
// DO NOT MODIFY THE RETURNED MAP
func (c *COWMap[K, V]) Peek() (m map[K]V) {
	return c.m
}

func (c *COWMap[K, V]) IsChanged() bool {
	return c.changed
}

var _ COW[map[string]any] = (*COWMap[string, any])(nil)

// Slice that will be copied before its first write, the original value is used if no writes were done
type COWSlice[V any] struct {
	s       []V
	changed bool
	// How to compare values of the slice
	equals func(V, V) bool
}

func NewCOWSliceFunc[V any](s []V, equals func(V, V) bool) *COWSlice[V] {
	return &COWSlice[V]{s, false, equals}
}

func NewCOWSliceComparable[V comparable](s []V) *COWSlice[V] {
	return &COWSlice[V]{s, false, IsComparableEqual[V]}
}

func (c *COWSlice[V]) Get(i int) (v V, ok bool) {
	if i >= len(c.s) {
		return
	}
	return c.s[i], true
}

func (c *COWSlice[V]) Len() int {
	return len(c.s)
}

func (c *COWSlice[V]) copyIfNeeded() {
	if !c.changed {
		if c.s == nil {
			c.s = make([]V, 0)
		} else {
			c.s = slices.Clone(c.s)
		}
		c.changed = true
	}
}

func (c *COWSlice[V]) Equals(other []V) bool {
	if c.equals == nil {
		return false
	}
	return slices.EqualFunc(c.s, other, c.equals)
}

// Checks if the given value exists at the target position
func (c *COWSlice[V]) ExistsAt(i int, value V) bool {
	if c.equals == nil {
		return false
	}
	if i >= len(c.s) {
		return false
	}
	existing := c.s[i]
	return c.equals(existing, value)
}

// Set is smart to not modify the slice if value is the same
//
// Will grow the slice if needed
func (c *COWSlice[V]) Set(i int, value V) {
	if c.ExistsAt(i, value) {
		return
	}

	c.copyIfNeeded()
	if i >= len(c.s) {
		ns := make([]V, i+1)
		copy(ns, c.s)
		c.s = ns
	}
	c.s[i] = value
}

func (c *COWSlice[V]) Contains(value V) bool {
	if c.equals == nil {
		return false
	}
	return slices.ContainsFunc(c.s, func(existing V) bool {
		return c.equals(existing, value)
	})
}

// If not Contains(value), then Append(value)
func (c *COWSlice[V]) Add(value V) {
	if !c.Contains(value) {
		c.Append(value)
	}
}

func (c *COWSlice[V]) Append(value V) {
	c.copyIfNeeded()
	c.s = append(c.s, value)
}

// Only does it if the slices are not equal.
//
// The COWSlice will be set as changed and other will be COPIED
func (c *COWSlice[V]) Replace(other []V) {
	if c.Equals(other) {
		return
	}
	c.changed = true
	c.s = slices.Clone(other)
}

func (c *COWSlice[V]) Release() (s []V, changed bool) {
	s = c.s
	changed = c.changed
	c.s = nil
	c.changed = false
	return s, changed
}

// Get the pointer to the internal reference.
//
// DO NOT MODIFY THE RETURNED SLICE
func (c *COWSlice[V]) Peek() (s []V) {
	return c.s
}

func (c *COWSlice[V]) IsChanged() bool {
	return c.changed
}

var _ COW[[]any] = (*COWSlice[any])(nil)
