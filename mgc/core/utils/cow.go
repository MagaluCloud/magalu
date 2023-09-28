package utils

// Copy-on-Write are helper to manage a mutable value, generating a copy whenever it will be written the first time
// If no writes are done, then the original object is never copied, saving resources

type COW[T any] interface {
	// Compare current value to another.
	//
	// May be called on nil COW handle, compares to the empty value of T.
	Equals(other T) bool
	// Return the current value and if it was changed, resets the internal pointers
	//
	// May be called on nil COW handle, returns the empty value of T.
	Release() (value T, changed bool)
	// Return the current value, do not modify it!
	//
	// May be called on nil COW handle, returns the empty value of T
	Peek() (value T)
	// Checks if the COW was mutated from its original value.
	//
	// May be called on nil COW handle, returns false.
	IsChanged() (changed bool)

	// Replace the internal value with another, if it's different.
	//
	// **CANNOT** be called on nil COW handle.
	Replace(other T) bool
}

// COW that is a container of other values (ie: slice, map)
type COWContainer[K comparable, V any] interface {
	// Iterates over all items of the container.
	//
	// If cb returns false the loop will stop and this function will also return false.
	// Otherwise this function returns true (meaning it looped through all items)
	//
	// May be called on nil COW handle, won't iterate and return true.
	ForEach(cb func(key K, value V) (run bool)) (finished bool)

	// Get the item at key.
	//
	// May be called on nil COW handle, returns the empty value of V and ok == false.
	Get(key K) (value V, ok bool)

	// Get the current container length.
	//
	// May be called on nil COW handle, returns 0.
	Len() int

	// Checks if the given value exists at the target key
	//
	// May be called on nil COW handle, returns false.
	ExistsAt(key K, value V) bool

	// Set is smart to not modify the map if value is the same
	//
	// Returns true if the value was modified, false otherwise.
	//
	// **CANNOT** be called on nil COW handle.
	Set(key K, value V) bool

	// Deletes the key (if it exists).
	//
	// May be called on nil COW handle, nothing is done
	Delete(key K) bool
}

// COW that is a container of other values (ie: slice, map) that are exposed as children COW.
type COWContainerOfCOW[K comparable, V any, C COW[V]] interface {
	// Iterates over all COW items of the container.
	//
	// If cb returns false the loop will stop and this function will also return false.
	// Otherwise this function returns true (meaning it looped through all items)
	//
	// May be called on nil COW handle, won't iterate and return true.
	ForEachCOW(cb func(key K, cow C) (run bool)) (finished bool)

	// Get the COW item at key.
	//
	// May be called on nil COW handle, returns the empty COW of C and ok == false.
	GetCOW(key K) (cow C, ok bool)
}
