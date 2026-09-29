package preview

// Package is one preview module: its name and its compiled spec.
type Package struct {
	Name string
	Spec []byte
}

// Store keeps the packages of one level on this machine.
type Store interface {
	// Names lists the installed packages, sorted. A store with nothing
	// installed has no names and no error.
	Names() ([]string, error)
	Read(name string) (Package, error)
	// String says where the packages live, for messages.
	String() string
}

// Validator decides whether a package may be used. It runs every time a
// module is loaded, since a file may have been put there or changed by hand.
type Validator interface {
	Validate(level Level, pkg Package) error
}
