package profile_manager

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path"
	"strings"

	"slices"

	"github.com/spf13/afero"

	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

type ProfileManager struct {
	dir string
	fs  afero.Fs
}

type contextKey string

var profileKey contextKey = "github.com/MagaluCloud/magalu/mgc/core/Profile"

func NewContext(parentCtx context.Context, profile *ProfileManager) context.Context {
	return context.WithValue(parentCtx, profileKey, profile)
}

func FromContext(ctx context.Context) *ProfileManager {
	a := ctx.Value(profileKey).(*ProfileManager)
	return a
}

func New() *ProfileManager {
	dir, err := buildMGCPath()
	if err != nil {
		dir = "."
	}

	return &ProfileManager{dir, afero.NewOsFs()}
}

func NewInMemoryProfileManager() (*ProfileManager, afero.Fs) {
	fs := afero.NewMemMapFs()
	pf := &ProfileManager{"/", fs}
	return pf, fs
}

func (m *ProfileManager) buildPath(name string) string {
	s := sanitizePath(name)
	return path.Join(m.dir, s)
}

func (m *ProfileManager) read(name string) ([]byte, error) {
	return afero.ReadFile(m.fs, m.buildPath(name))
}

func (m *ProfileManager) prepareWrite(name string) (fullPath, dirName string, err error) {
	fullPath = m.buildPath(name)
	dirName = path.Dir(fullPath)
	err = m.fs.MkdirAll(dirName, utils.DIR_PERMISSION)
	return
}

func (m *ProfileManager) write(name string, data []byte) (err error) {
	fullPath, _, err := m.prepareWrite(name)
	if err != nil {
		return
	}
	return afero.WriteFile(m.fs, fullPath, data, utils.FILE_PERMISSION)
}

func (m *ProfileManager) remove(name string) error {
	return m.fs.RemoveAll(m.buildPath(name))
}

// NOTE: cb(path) won't have full prefix, since we'll remove it in order for a simple m.read(path) work
// only files are reported.
// err can be fs.SkipAll or fs.SkipDir to stop iterating, no error will be returned
func (m *ProfileManager) walk(name string, cb func(path string) error) error {
	root := m.buildPath(name)
	prefixLen := len(root) + 1 // "dir/name" + "/"
	return afero.Walk(m.fs, root, func(path string, info fs.FileInfo, err error) error {
		if path == root {
			return nil
		}
		if err != nil {
			return err
		}
		path = path[prefixLen:]
		return cb(path)
	})
}

func (m *ProfileManager) Get(name string) (p *Profile, err error) {
	if err = checkProfileName(name); err != nil {
		return
	}

	return newProfile(name, m), nil
}

func (m *ProfileManager) Current() *Profile {
	var name string

	data, err := m.read(currentProfileNameFile)
	if err != nil || len(data) == 0 {
		name = defaultProfileName
	} else {
		name = string(data)
	}

	p, err := m.Get(name)
	if err != nil {
		p, err = m.Get(defaultProfileName)
		if err != nil {
			// Should never happen
			panic("default profile should always work")
		}
	}

	return p
}

func (m *ProfileManager) SetCurrent(p *Profile) error {
	return m.write(currentProfileNameFile, []byte(p.Name))
}

func (m *ProfileManager) Create(name string) (p *Profile, err error) {
	if p, err = m.Get(name); err != nil {
		return
	}

	fullPath, _, err := m.prepareWrite(name)
	if err != nil {
		return
	}

	err = m.fs.Mkdir(fullPath, utils.DIR_PERMISSION)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			err = errorProfileAlreadyExists
		}
	}

	return
}

func (m *ProfileManager) Copy(src, dst *Profile) error {
	if src.Name == dst.Name {
		return errorCopyToSelf
	}
	return m.walk(src.Name, func(name string) (err error) {
		data, err := m.read(path.Join(src.Name, name))
		if err != nil {
			return
		}

		return m.write(path.Join(dst.Name, name), data)
	})
}

func (m *ProfileManager) Delete(p *Profile) error {
	if m.Current().Name == p.Name {
		return errorDeleteCurrentNotAllowed
	}
	return m.remove(p.Name)
}

func (m *ProfileManager) List() (profiles []*Profile) {
	entries, err := afero.ReadDir(m.fs, m.dir)
	current := m.Current()
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				if p, err := m.Get(e.Name()); err == nil {
					if current != nil && p.Name == current.Name {
						current = nil
					}
					profiles = append(profiles, p)
				}
			}
		}
	}

	if current != nil {
		profiles = append(profiles, current)
	}

	slices.SortFunc(profiles, func(a, b *Profile) int {
		return strings.Compare(a.Name, b.Name)
	})

	return
}
