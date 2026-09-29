// Package fsstore keeps preview packages as files in a dir, such as
// ~/.config/mgc/beta: the package "<name>" is the file "<name>.openapi.yaml".
package fsstore

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
)

const specSuffix = ".openapi.yaml"

type Store struct {
	dir string
}

func New(dir string) *Store {
	return &Store{dir: dir}
}

// Names lists the spec files kept in the dir, sorted. A missing dir has none.
func (s *Store) Names() ([]string, error) {
	entries, err := os.ReadDir(s.dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		name, ok := strings.CutSuffix(entry.Name(), specSuffix)
		if entry.IsDir() || !ok || name == "" {
			continue
		}
		names = append(names, name)
	}
	return names, nil
}

func (s *Store) Read(name string) (preview.Package, error) {
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return preview.Package{}, fmt.Errorf("preview package %q must be a file name inside %s", name, s.dir)
	}

	spec, err := os.ReadFile(filepath.Join(s.dir, name+specSuffix))
	if err != nil {
		return preview.Package{}, err
	}
	return preview.Package{Name: name, Spec: spec}, nil
}

func (s *Store) String() string {
	return s.dir
}

var _ preview.Store = (*Store)(nil)
