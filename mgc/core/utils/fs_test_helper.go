package utils

import (
	"bytes"
	"fmt"
	"io/fs"

	"github.com/spf13/afero"
)

type TestFsEntry struct {
	Path string
	Mode fs.FileMode
	Data []byte
}

func FindFsEntry(path string, entries []TestFsEntry) (TestFsEntry, error) {
	for _, e := range entries {
		if e.Path == path {
			return e, nil
		}
	}
	return TestFsEntry{}, fmt.Errorf("%q: %w", path, fs.ErrNotExist)
}
func PrepareFs(afs afero.Fs, provided []TestFsEntry) (err error) {
	for _, p := range provided {
		if p.Mode&fs.ModeDir != 0 {
			err = afs.MkdirAll(p.Path, p.Mode)
		} else {
			err = afero.WriteFile(afs, p.Path, p.Data, p.Mode)
		}
		if err != nil {
			return
		}
	}
	return
}

func CheckFs(afs afero.Fs, expected []TestFsEntry) (err error) {
	existingFiles := 0
	err = afero.Walk(afs, "/", func(path string, info fs.FileInfo, e error) (err error) {
		if e != nil {
			return e
		}
		if path == "/" {
			return nil
		}
		fsEntry, err := FindFsEntry(path, expected)
		if err != nil {
			return
		}
		if fsEntry.Mode != info.Mode() {
			return fmt.Errorf("%s: expected mode %x, got %x", path, fsEntry.Mode, info.Mode())
		}
		if fsEntry.Mode&fs.ModeDir == 0 {
			var data []byte
			if data, err = afero.ReadFile(afs, path); err != nil {
				return
			}
			if !bytes.Equal(fsEntry.Data, data) {
				return fmt.Errorf("%s: expected data %q, got %q", path, fsEntry.Data, data)
			}
		}
		existingFiles++
		return nil
	})
	if err != nil {
		return
	}

	if len(expected) != existingFiles {
		return fmt.Errorf("expected %d FS entries, got %d", len(expected), existingFiles)
	}

	return
}

func MergeFsEntries(toBeMerged ...[]TestFsEntry) (merged []TestFsEntry) {
	knownPaths := map[string]bool{}
	for _, entries := range toBeMerged {
		for _, e := range entries {
			if !knownPaths[e.Path] {
				knownPaths[e.Path] = true
				merged = append(merged, e)
			}
		}
	}
	return
}

func getDirs(p string) (dirs []string) {
	for i, c := range p {
		if c == '/' && i > 0 {
			dirs = append(dirs, p[:i])
		}
	}
	return
}

func AutoMkdirAll(entries []TestFsEntry) (expanded []TestFsEntry) {
	knownPaths := map[string]bool{}
	for _, e := range entries {
		knownPaths[e.Path] = true
	}

	for _, e := range entries {
		for _, d := range getDirs(e.Path) {
			if !knownPaths[d] {
				knownPaths[d] = true
				expanded = append(expanded, TestFsEntry{
					Path: d,
					Mode: fs.ModeDir | DIR_PERMISSION,
					Data: nil,
				})
			}
		}
		expanded = append(expanded, e)
	}

	return expanded
}
