package utils

import "golang.org/x/exp/maps"

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

// Iterates over all map items.
//
// If cb returns false the loop will stop and this function will also return false.
// Otherwise this function returns true (meaning it looped through all items)
//
// May be called on nil COW handle, won't iterate and return true.
func (c *COWMap[K, V]) ForEach(cb func(key K, value V) (run bool)) (finished bool) {
	if c == nil {
		return true
	}
	for key, value := range c.m {
		if !cb(key, value) {
			return false
		}
	}
	return true
}

// Get the item at key.
//
// May be called on nil COW handle, returns the empty value of V and ok == false.
func (c *COWMap[K, V]) Get(key K) (value V, ok bool) {
	if c == nil {
		return
	}
	value, ok = c.m[key]
	return value, ok
}

// Get the current map length.
//
// May be called on nil COW handle, returns 0.
func (c *COWMap[K, V]) Len() int {
	if c == nil {
		return 0
	}
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
	if c == nil {
		return len(other) == 0
	}
	if c.equals == nil {
		return false
	}
	return maps.EqualFunc(c.m, other, c.equals)
}

// Checks if the given value exists at the target key
//
// May be called on nil COW handle, returns false.
func (c *COWMap[K, V]) ExistsAt(key K, value V) bool {
	if c == nil || c.equals == nil {
		return false
	}
	if existing, ok := c.m[key]; ok {
		return c.equals(existing, value)
	}
	return false
}

// Set is smart to not modify the map if value is the same
//
// Returns true if the value was modified, false otherwise.
//
// **CANNOT** be called on nil COW handle.
func (c *COWMap[K, V]) Set(key K, value V) (mutated bool) {
	if c.ExistsAt(key, value) {
		return
	}

	c.copyIfNeeded()
	c.m[key] = value
	return true
}

// Deletes the key (if it exists).
//
// May be called on nil COW handle, nothing is done
func (c *COWMap[K, V]) Delete(key K) (mutated bool) {
	if c == nil {
		return
	}
	if _, ok := c.m[key]; !ok {
		return
	}

	c.copyIfNeeded()
	delete(c.m, key)
	return true
}

// Only does it if the maps are not equal.
//
// The COWMap will be set as changed and other will be COPIED.
//
// **CANNOT** be called on nil COW handle.
func (c *COWMap[K, V]) Replace(other map[K]V) bool {
	if c.Equals(other) {
		return false
	}
	c.changed = true
	c.m = maps.Clone(other)
	return true
}

func (c *COWMap[K, V]) Release() (m map[K]V, changed bool) {
	if c == nil {
		return
	}
	m = c.Peek()
	changed = c.IsChanged()
	c.m = nil
	c.changed = false
	return m, changed
}

// Get the pointer to the internal reference.
//
// DO NOT MODIFY THE RETURNED MAP
func (c *COWMap[K, V]) Peek() (m map[K]V) {
	if c == nil {
		return
	}
	return c.m
}

func (c *COWMap[K, V]) IsChanged() (changed bool) {
	if c == nil {
		return
	}
	return c.changed
}

var _ COW[map[string]any] = (*COWMap[string, any])(nil)
var _ COWContainer[string, any] = (*COWMap[string, any])(nil)
