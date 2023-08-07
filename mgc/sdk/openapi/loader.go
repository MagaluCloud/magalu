package openapi

import (
	"fmt"
	"os"
	"path"
	"syscall"
)

type Loader interface {
	Load(name string) ([]byte, error)
}

type FileLoader struct {
	Dir string
}

func (f FileLoader) Load(name string) ([]byte, error) {
	return os.ReadFile(path.Join(f.Dir, name))
}

func (f FileLoader) String() string {
	return fmt.Sprintf("FileLoader(dir: %s)", f.Dir)
}

var _ Loader = (*FileLoader)(nil)

type MergeLoader struct {
	Loaders []Loader
}

func NewMergeLoader(loaders ...Loader) Loader {
	if len(loaders) == 1 {
		return loaders[0]
	}
	return &MergeLoader{loaders}
}

func (m *MergeLoader) Load(name string) (data []byte, err error) {
	for _, loader := range m.Loaders {
		data, err = loader.Load(name)
		if err == nil {
			return
		}
	}
	if err == nil {
		err = &os.PathError{Op: "open", Path: name, Err: syscall.ENOENT}
	}
	return nil, err
}

func (m *MergeLoader) String() string {
	return fmt.Sprintf("MergeLoader(loaders: %s)", m.Loaders)
}
