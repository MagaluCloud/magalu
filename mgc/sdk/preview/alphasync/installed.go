package alphasync

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
)

const (
	indexFile  = "index.json"
	specSuffix = ".openapi.yaml"
)

// Installed is an alpha module on this machine. It keeps the manifest that
// listed it, so that it upgrades whatever tenant the CLI uses later.
type Installed struct {
	ManifestSpec
	Manifest string `json:"manifest"`
	Pinned   bool   `json:"pinned,omitempty"`
}

// specFile names the spec of an installed module after its content, so that a
// new one never overwrites the one the index vouches for: writing the index
// switches from one to the other at once.
func specFile(entry Installed) string {
	return entry.Module + "@" + strings.ToLower(entry.SHA256) + specSuffix
}

// InstalledStore serves the alpha modules the index lists, checking each one
// against its signed item every time it is read. The index is read once.
func InstalledStore(dir string, keys Keys) preview.Store {
	return &verifiedStore{dir: dir, keys: keys}
}

type verifiedStore struct {
	dir   string
	keys  Keys
	once  sync.Once
	index map[string]Installed
	err   error
}

func (s *verifiedStore) readIndex() (map[string]Installed, error) {
	s.once.Do(func() { s.index, s.err = readIndex(s.dir) })
	return s.index, s.err
}

func (s *verifiedStore) Names() ([]string, error) {
	index, err := s.readIndex()
	if err != nil {
		return nil, err
	}
	return slices.Sorted(maps.Keys(index)), nil
}

func (s *verifiedStore) Read(name string) (preview.Package, error) {
	index, err := s.readIndex()
	if err != nil {
		return preview.Package{}, err
	}
	entry, ok := index[name]
	if !ok || entry.Module != name {
		return preview.Package{}, fmt.Errorf("alpha module %q is not installed", name)
	}
	if err := s.keys.Verify(entry.ManifestSpec); err != nil {
		return preview.Package{}, err
	}

	spec, err := os.ReadFile(filepath.Join(s.dir, specFile(entry)))
	if errors.Is(err, fs.ErrNotExist) {
		return preview.Package{}, fmt.Errorf("its spec is missing: run 'mgc preview install %s'", name)
	}
	if err != nil {
		return preview.Package{}, err
	}
	sum := sha256.Sum256(spec)
	if hex.EncodeToString(sum[:]) != strings.ToLower(entry.SHA256) {
		return preview.Package{}, fmt.Errorf("it changed since it was installed: run 'mgc preview install %s'", name)
	}
	return preview.Package{Name: name, Spec: spec}, nil
}

func (s *verifiedStore) String() string {
	return s.dir
}

// intact tells whether the installed spec is still the one its entry vouches
// for, so that installing it again repairs it.
func intact(dir string, entry Installed) bool {
	data, err := os.ReadFile(filepath.Join(dir, specFile(entry)))
	if err != nil {
		return false
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]) == strings.ToLower(entry.SHA256)
}

func readIndex(dir string) (map[string]Installed, error) {
	file := filepath.Join(dir, indexFile)
	data, err := os.ReadFile(file)
	if errors.Is(err, fs.ErrNotExist) {
		return map[string]Installed{}, nil
	}
	if err != nil {
		return nil, err
	}

	index := map[string]Installed{}
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("%s is damaged: remove it and install the modules again: %w", file, err)
	}
	return index, nil
}

// apply writes the specs of the steps under new names, and then the index that
// switches to them in one rename: a failure before it leaves everything as it
// was. The specs the index no longer names go last; one left behind is unused.
func apply(dir string, index map[string]Installed, steps []step) error {
	if len(steps) == 0 {
		return nil
	}
	next := maps.Clone(index)
	var written []string
	for _, s := range steps {
		if s.entry == nil {
			delete(next, s.module)
			continue
		}
		next[s.module] = *s.entry
		if s.spec == nil {
			continue
		}
		// Repairing a spec rewrites the file the index already names.
		name := specFile(*s.entry)
		if err := writeFile(dir, name, s.spec); err != nil {
			removeFiles(dir, written)
			return err
		}
		if current, ok := index[s.module]; !ok || specFile(current) != name {
			written = append(written, name)
		}
	}
	if err := writeIndex(dir, next); err != nil {
		removeFiles(dir, written)
		return err
	}

	var unused []string
	for module, entry := range index {
		if current, ok := next[module]; !ok || specFile(current) != specFile(entry) {
			unused = append(unused, specFile(entry))
		}
	}
	removeFiles(dir, unused)
	return nil
}

func removeFiles(dir string, names []string) {
	for _, name := range names {
		if err := os.Remove(filepath.Join(dir, name)); err != nil && !errors.Is(err, fs.ErrNotExist) {
			logger().Debugw("leaving an unused alpha spec behind", "file", name, "error", err)
		}
	}
}

func writeIndex(dir string, index map[string]Installed) error {
	data, err := json.Marshal(index)
	if err != nil {
		return err
	}
	return writeFile(dir, indexFile, data)
}

// writeFile replaces a file in one rename, readable only by its owner.
func writeFile(dir, name string, data []byte) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "."+name+"-")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), filepath.Join(dir, name))
}
