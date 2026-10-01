// Package fsstore keeps preview packages as files in a file system: the package
// "<name>" is the file "<name>.openapi.yaml". The file system is either embedded
// in the binary (the beta specs, which ship with the CLI) or a dir on disk (the
// alpha specs a customer installs, see NewDir).
package fsstore

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
)

const specSuffix = ".openapi.yaml"

type Store struct {
	fsys  fs.FS
	where string
}

// New reads the packages at the root of fsys. where says where they live, for
// messages.
func New(fsys fs.FS, where string) *Store {
	return &Store{fsys: fsys, where: where}
}

// NewDir reads the packages kept in dir. A missing dir has none.
func NewDir(dir string) *Store {
	return New(os.DirFS(dir), dir)
}

// Names lists the spec files, sorted. A missing dir has none.
func (s *Store) Names() ([]string, error) {
	entries, err := fs.ReadDir(s.fsys, ".")
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
	file := name + specSuffix
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) || !fs.ValidPath(file) {
		return preview.Package{}, fmt.Errorf("preview package %q must be a file name inside %s", name, s.where)
	}

	spec, err := fs.ReadFile(s.fsys, file)
	if err != nil {
		return preview.Package{}, err
	}
	return preview.Package{Name: name, Spec: spec}, nil
}

func (s *Store) String() string {
	return s.where
}

var _ preview.Store = (*Store)(nil)
