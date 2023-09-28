package utils

import "golang.org/x/exp/slices"

// Slice that will be copied before its first write, the original value is used if no writes were done
//
// This variant will create a COW for each element, making it easy for users. Writes to the
// element's COW will update the parent (container) slice (s) given at construction.
type COWSliceOfCOW[V any, C COW[V]] struct {
	s          []V
	cow        []C
	cowFactory func(V) C
	changed    bool
}

func (c *COWSliceOfCOW[K, C]) fillCOWEntry(i int) {
	c.cow[i] = c.cowFactory(c.s[i])
}

func (c *COWSliceOfCOW[K, C]) fillCOWRange(start, end int) {
	for i := start; i < end; i++ {
		c.fillCOWEntry(i)
	}
}

func (c *COWSliceOfCOW[K, C]) initCOW() {
	c.cow = make([]C, len(c.s))
	c.fillCOWRange(0, len(c.s))
}

func (c *COWSliceOfCOW[K, C]) isCOWChanged() bool {
	for _, v := range c.cow {
		if v.IsChanged() {
			return true
		}
	}
	return false
}

// Sub COW are handled apart, but whenever we need to return the slice
// we must copy the slice if needed and then set all
// public pointers to the latest value of each COW
func (c *COWSliceOfCOW[K, C]) materializeCOW() {
	if !c.isCOWChanged() {
		return
	}
	c.copyIfNeeded()
	for i, cow := range c.cow {
		c.s[i] = cow.Peek()
	}
}

func NewCOWSliceOfCOW[V any, C COW[V]](s []V, cowFactory func(V) C) *COWSliceOfCOW[V, C] {
	c := &COWSliceOfCOW[V, C]{
		s:          s,
		cow:        nil,
		cowFactory: cowFactory,
		changed:    false,
	}
	c.initCOW()
	return c
}

// Iterates over all slice items.
//
// If cb returns false the loop will stop and this function will also return false.
// Otherwise this function returns true (meaning it looped through all items)
//
// May be called on nil COW handle, won't iterate and returns true.
func (c *COWSliceOfCOW[V, C]) ForEach(cb func(index int, item V) (run bool)) (finished bool) {
	for i, cow := range c.cow {
		if !cb(i, cow.Peek()) {
			return false
		}
	}
	return true
}

// Iterates over all slice COW items.
//
// If cb returns false the loop will stop and this function will also return false.
// Otherwise this function returns true (meaning it looped through all items)
//
// May be called on nil COW handle, won't iterate and returns true.
func (c *COWSliceOfCOW[V, C]) ForEachCOW(cb func(index int, cow C) (run bool)) (finished bool) {
	for i, cow := range c.cow {
		if !cb(i, cow) {
			return false
		}
	}
	return true
}

// Gets the Copy-on-Write value at the given index.
//
// If the COW is mutated, the COWSliceOfCOW (c) will reflect the new value automatically.
//
// For the plain value, use Get().
//
// May be called on nil COW handle, returns the empty value of C and ok == false.
func (c *COWSliceOfCOW[V, C]) GetCOW(i int) (cow C, ok bool) {
	if i >= c.Len() {
		return
	}
	cow = c.cow[i]
	return cow, ok
}

// Get the item at index.
//
// May be called on nil COW handle, returns the empty value of V and ok == false.
func (c *COWSliceOfCOW[V, C]) Get(i int) (v V, ok bool) {
	cow, ok := c.GetCOW(i)
	if !ok {
		return
	}
	return cow.Peek(), ok
}

// Get the current slice length.
//
// May be called on nil COW handle, returns 0.
func (c *COWSliceOfCOW[V, C]) Len() int {
	if c == nil {
		return 0
	}
	return len(c.cow)
}

// Resize the slice, if needed
//
// **CANNOT** be called on nil COW handle
func (c *COWSliceOfCOW[V, C]) Resize(newSize int) {
	oldLen := c.Len()
	switch {
	case oldLen == newSize:
		return
	case oldLen > newSize:
		c.cow = c.cow[:newSize]
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
		nCow := make([]C, newSize)
		copy(ns, c.s)
		copy(nCow, c.cow)
		c.s = ns
		c.cow = nCow
		c.fillCOWRange(c.Len(), newSize) // fill with cow using empty values
		c.changed = true
		return
	}
}

func (c *COWSliceOfCOW[V, C]) copyIfNeeded() {
	if !c.changed {
		if c.s == nil {
			c.s = make([]V, 0)
		} else {
			c.s = slices.Clone(c.s)
		}
		c.changed = true
	}
}

func (c *COWSliceOfCOW[V, C]) Equals(other []V) bool {
	if c.Len() != len(other) {
		return false
	}
	for i, value := range other {
		if !c.ExistsAt(i, value) {
			return false
		}
	}
	return true
}

// Checks if the given value exists at the target position
//
// May be called on nil COW handle, returns false.
func (c *COWSliceOfCOW[V, C]) ExistsAt(i int, value V) bool {
	if cow, ok := c.GetCOW(i); ok {
		return cow.Equals(value)
	}
	return false
}

// Set is smart to not modify the slice if value is the same
//
// Will grow the slice if needed.
//
// Returns true if the value was modified, false otherwise.
//
// **CANNOT** be called on nil COW handle.
func (c *COWSliceOfCOW[V, C]) Set(i int, value V) (mutated bool) {
	if c.ExistsAt(i, value) {
		return
	}

	if i >= c.Len() {
		c.Resize(i + 1)
	}
	c.copyIfNeeded()
	c.s[i] = value
	c.fillCOWEntry(i)
	return true
}

// Deletes the index (if it exists).
//
// May be called on nil COW handle, nothing is done
func (c *COWSliceOfCOW[V, C]) Delete(i int) (mutated bool) {
	if i >= c.Len() {
		return
	}

	c.copyIfNeeded()
	c.s = slices.Delete(c.s, i, i+1)
	c.cow = slices.Delete(c.cow, i, i+1)
	return true
}

// Checks if the given value exists in the slice
//
// May be called on nil COW handle, returns false.
func (c *COWSliceOfCOW[V, C]) Contains(value V) bool {
	if c == nil {
		return false
	}
	for _, cow := range c.cow {
		if cow.Equals(value) {
			return true
		}
	}
	return false
}

// If not Contains(value), then Append(value)
//
// **CANNOT** be called on nil COW handle.
func (c *COWSliceOfCOW[V, C]) Add(value V) {
	if !c.Contains(value) {
		c.Append(value)
	}
}

// Unconditionally appends the value.
//
// **CANNOT** be called on nil COW handle.
func (c *COWSliceOfCOW[V, C]) Append(value V) {
	c.copyIfNeeded()
	i := len(c.s)
	c.s = append(c.s, value)
	c.fillCOWEntry(i)
}

// Only does it if the slices are not equal.
//
// # The COWSlice will be set as changed and other will be COPIED
//
// **CANNOT** be called on nil COW handle.
func (c *COWSliceOfCOW[V, C]) Replace(other []V) bool {
	if c.Equals(other) {
		return false
	}
	c.changed = true
	c.s = slices.Clone(other)
	c.initCOW()
	return true
}

func (c *COWSliceOfCOW[V, C]) Release() (s []V, changed bool) {
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
// DO NOT MODIFY THE RETURNED SLICE
func (c *COWSliceOfCOW[V, C]) Peek() (s []V) {
	if c == nil {
		return
	}
	c.materializeCOW()
	return c.s
}

func (c *COWSliceOfCOW[V, C]) IsChanged() (changed bool) {
	if c == nil {
		return
	}
	return c.changed || c.isCOWChanged()
}

var _ COW[[]any] = (*COWSliceOfCOW[any, COW[any]])(nil)
var _ COWContainer[int, any] = (*COWSliceOfCOW[any, COW[any]])(nil)
var _ COWContainerOfCOW[int, any, COW[any]] = (*COWSliceOfCOW[any, COW[any]])(nil)
