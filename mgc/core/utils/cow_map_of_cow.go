package utils

import "golang.org/x/exp/maps"

// Map that will be copied before its first write, the original value is used if no writes were done
//
// This variant will create a COW for each element, making it easy for users. Writes to the
// element's COW will update the parent (container) map (m) given at construction.
type COWMapOfCOW[K comparable, V any, C COW[V]] struct {
	m          map[K]V
	cow        map[K]C
	cowFactory func(V) C
	changed    bool
}

func (c *COWMapOfCOW[K, V, C]) fillCOWEntry(key K) {
	c.cow[key] = c.cowFactory(c.m[key])
}

func (c *COWMapOfCOW[K, V, C]) initCOW() {
	c.cow = make(map[K]C, len(c.m))
	for k := range c.m {
		c.fillCOWEntry(k)
	}
}

func (c *COWMapOfCOW[K, V, C]) isCOWChanged() bool {
	for _, v := range c.cow {
		if v.IsChanged() {
			return true
		}
	}
	return false
}

// Sub COW are handled apart, but whenever we need to return the map
// we must copy the map if needed and then set all
// public pointers to the latest value of each COW
func (c *COWMapOfCOW[K, V, C]) materializeCOW() {
	if !c.isCOWChanged() {
		return
	}
	c.copyIfNeeded()
	for k, cow := range c.cow {
		c.m[k] = cow.Peek()
	}
}

func NewCOWMapOfCOW[K comparable, V any, C COW[V]](m map[K]V, cowFactory func(V) C) *COWMapOfCOW[K, V, C] {
	c := &COWMapOfCOW[K, V, C]{
		m:          m,
		cow:        nil,
		cowFactory: cowFactory,
		changed:    false,
	}
	c.initCOW()
	return c
}

// Iterates over all map items.
//
// If cb returns false the loop will stop and this function will also return false.
// Otherwise this function returns true (meaning it looped through all items)
//
// May be called on nil COW handle, won't iterate and return true.
func (c *COWMapOfCOW[K, V, C]) ForEach(cb func(key K, item V) (run bool)) (finished bool) {
	if c == nil {
		return true
	}
	for key, cow := range c.cow {
		if !cb(key, cow.Peek()) {
			return false
		}
	}
	return true
}

// Iterates over all map COW items.
//
// If cb returns false the loop will stop and this function will also return false.
// Otherwise this function returns true (meaning it looped through all items)
//
// May be called on nil COW handle, won't iterate and return true.
func (c *COWMapOfCOW[K, V, C]) ForEachCOW(cb func(key K, cow C) (run bool)) (finished bool) {
	if c == nil {
		return true
	}
	for key, cow := range c.cow {
		if !cb(key, cow) {
			return false
		}
	}
	return true
}

// Gets the Copy-on-Write value at the given key.
//
// If the COW is mutated, the COWMapOfCOW (c) will reflect the new value automatically.
//
// For the plain value, use Get().
func (c *COWMapOfCOW[K, V, C]) GetCOW(key K) (cow C, ok bool) {
	if c == nil {
		return
	}
	cow, ok = c.cow[key]
	return cow, ok
}

// Gets the plain value at the given key. Do not mutate it, if you want to mutate use GetCOW()
func (c *COWMapOfCOW[K, V, C]) Get(key K) (v V, ok bool) {
	cow, ok := c.GetCOW(key)
	if !ok {
		return
	}
	return cow.Peek(), ok
}

// Get the current map length.
//
// May be called on nil COW handle, returns 0.
func (c *COWMapOfCOW[K, V, C]) Len() int {
	if c == nil {
		return 0
	}
	return len(c.cow)
}

func (c *COWMapOfCOW[K, V, C]) copyIfNeeded() {
	if !c.changed {
		if c.m == nil {
			c.m = make(map[K]V)
		} else {
			c.m = maps.Clone(c.m)
		}
		c.changed = true
	}
}

func (c *COWMapOfCOW[K, V, C]) Equals(other map[K]V) bool {
	if c.Len() != len(other) {
		return false
	}
	for key, value := range other {
		if !c.ExistsAt(key, value) {
			return false
		}
	}
	return true
}

// Checks if the given value exists at the target key
//
// May be called on nil COW handle, returns false.
func (c *COWMapOfCOW[K, V, C]) ExistsAt(key K, value V) bool {
	if cow, ok := c.GetCOW(key); ok {
		return cow.Equals(value)
	}
	return false
}

// Set is smart to not modify the map if value is the same
//
// Returns true if the value was modified, false otherwise.
//
// **CANNOT** be called on nil COW handle.
func (c *COWMapOfCOW[K, V, C]) Set(key K, value V) (mutated bool) {
	if cow, ok := c.GetCOW(key); ok {
		return cow.Replace(value)
	} else {
		c.copyIfNeeded()
		c.m[key] = value
		c.fillCOWEntry(key)
		return true
	}
}

// Deletes the key (if it exists).
//
// May be called on nil COW handle, nothing is done
func (c *COWMapOfCOW[K, V, C]) Delete(key K) (mutated bool) {
	if c == nil {
		return
	}
	if _, ok := c.GetCOW(key); !ok {
		return
	}

	c.copyIfNeeded()
	delete(c.m, key)
	delete(c.cow, key)
	return true
}

// Only does it if the maps are not equal.
//
// The COWMap will be set as changed and other will be COPIED.
//
// **CANNOT** be called on nil COW handle.
func (c *COWMapOfCOW[K, V, C]) Replace(other map[K]V) bool {
	if c.Equals(other) {
		return false
	}
	c.changed = true
	c.m = maps.Clone(other)
	c.initCOW()
	return true
}

func (c *COWMapOfCOW[K, V, C]) Release() (m map[K]V, changed bool) {
	if c == nil {
		return
	}
	m = c.Peek()
	changed = c.IsChanged()
	c.m = nil
	c.changed = false
	c.initCOW()
	return m, changed
}

// Get the pointer to the internal reference.
//
// DO NOT MODIFY THE RETURNED MAP
func (c *COWMapOfCOW[K, V, C]) Peek() (m map[K]V) {
	if c == nil {
		return
	}
	c.materializeCOW()
	return c.m
}

func (c *COWMapOfCOW[K, V, C]) IsChanged() (changed bool) {
	if c == nil {
		return
	}
	return c.changed || c.isCOWChanged()
}

var _ COW[map[string]any] = (*COWMapOfCOW[string, any, COW[any]])(nil)
var _ COWContainer[string, any] = (*COWMapOfCOW[string, any, COW[any]])(nil)
var _ COWContainerOfCOW[string, any, COW[any]] = (*COWMapOfCOW[string, any, COW[any]])(nil)
