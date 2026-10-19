package preview

type Package struct {
	Name string
	Spec []byte
}

type Store interface {
	Names() ([]string, error)
	Read(name string) (Package, error)
	String() string
}

type Validator interface {
	Validate(level Level, pkg Package) error
}
