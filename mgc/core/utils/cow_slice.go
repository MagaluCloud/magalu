package utils

import "golang.org/x/exp/slices"

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

// Iterates over all slice items.
//
// If cb returns false the loop will stop and this function will also return false.
// Otherwise this function returns true (meaning it looped through all items)
//
// May be called on nil COW handle, won't iterate and returns true.
func (c *COWSlice[V]) ForEach(cb func(index int, value V) (run bool)) (finished bool) {
	if c == nil {
		return true
	}
	for i, value := range c.s {
		if !cb(i, value) {
			return false
		}
	}
	return true
}

// Get the item at index.
//
// May be called on nil COW handle, returns the empty value of V and ok == false.
func (c *COWSlice[V]) Get(i int) (value V, ok bool) {
	if i >= c.Len() {
		return
	}
	return c.s[i], true
}

// Get the current slice length.
//
// May be called on nil COW handle, returns 0.
func (c *COWSlice[V]) Len() int {
	if c == nil {
		return 0
	}
	return len(c.s)
}

// Resize the slice, if needed
//
// **CANNOT** be called on nil COW handle
func (c *COWSlice[V]) Resize(newSize int) {
	oldLen := c.Len()
	switch {
	case oldLen == newSize:
		return
	case oldLen > newSize:
		ns := c.s[:newSize]
		if c.changed {
			c.s = ns
		} else {
			c.s = slices.Clone(ns)
			c.changed = true
		}
		return
	case oldLen < newSize:
		ns := make([]V, newSize)
		copy(ns, c.s)
		c.s = ns
		c.changed = true
		return
	}
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
	if c == nil {
		return len(other) == 0
	}
	if c.equals == nil {
		return false
	}
	return slices.EqualFunc(c.s, other, c.equals)
}

// Checks if the given value exists at the target position
//
// May be called on nil COW handle, returns false.
func (c *COWSlice[V]) ExistsAt(i int, value V) bool {
	if c == nil || c.equals == nil {
		return false
	}
	if i >= c.Len() {
		return false
	}
	existing := c.s[i]
	return c.equals(existing, value)
}

// Set is smart to not modify the slice if value is the same
//
// Will grow the slice if needed.
//
// Returns true if the value was modified, false otherwise.
//
// **CANNOT** be called on nil COW handle.
func (c *COWSlice[V]) Set(i int, value V) (mutated bool) {
	if c.ExistsAt(i, value) {
		return
	}

	if i >= c.Len() {
		c.Resize(i + 1)
	}
	c.copyIfNeeded()
	c.s[i] = value
	return true
}

// Deletes the index (if it exists).
//
// May be called on nil COW handle, nothing is done
func (c *COWSlice[V]) Delete(i int) (mutated bool) {
	if i >= c.Len() {
		return
	}

	c.copyIfNeeded()
	c.s = slices.Delete(c.s, i, i+1)
	return true
}

// Checks if the given value exists in the slice
//
// May be called on nil COW handle, returns false.
func (c *COWSlice[V]) Contains(value V) bool {
	if c == nil || c.equals == nil {
		return false
	}
	return slices.ContainsFunc(c.s, func(existing V) bool {
		return c.equals(existing, value)
	})
}

// If not Contains(value), then Append(value)
//
// **CANNOT** be called on nil COW handle.
func (c *COWSlice[V]) Add(value V) {
	if !c.Contains(value) {
		c.Append(value)
	}
}

// Unconditionally appends the value.
//
// **CANNOT** be called on nil COW handle.
func (c *COWSlice[V]) Append(value V) {
	c.copyIfNeeded()
	c.s = append(c.s, value)
}

// Only does it if the slices are not equal.
//
// The COWSlice will be set as changed and other will be COPIED.
//
// **CANNOT** be called on nil COW handle.
func (c *COWSlice[V]) Replace(other []V) bool {
	if c.Equals(other) {
		return false
	}
	c.changed = true
	c.s = slices.Clone(other)
	return true
}

func (c *COWSlice[V]) Release() (s []V, changed bool) {
	if c == nil {
		return
	}
	s = c.Peek()
	changed = c.IsChanged()
	c.s = nil
	c.changed = false
	return s, changed
}

// Get the pointer to the internal reference.
//
// DO NOT MODIFY THE RETURNED SLICE
func (c *COWSlice[V]) Peek() (s []V) {
	if c == nil {
		return
	}
	return c.s
}

func (c *COWSlice[V]) IsChanged() (changed bool) {
	if c == nil {
		return
	}
	return c.changed
}

var _ COW[[]any] = (*COWSlice[any])(nil)
var _ COWContainer[int, any] = (*COWSlice[any])(nil)
