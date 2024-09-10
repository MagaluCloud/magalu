package utils

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"slices"
	"strings"
)

func boundReadDir(dir fs.ReadDirFile, name string, maxReadDir int) ([]fs.DirEntry, error) {
	list, err := dir.ReadDir(maxReadDir)
	slices.SortFunc(list, func(a, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})
	return list, err
}

func boundWalkDirInternal(fsys fs.FS, name string, d fs.DirEntry, walkDirFn fs.WalkDirFunc, maxReadDir int) error {
	if err := walkDirFn(name, d, nil); err != nil || !d.IsDir() {
		if err == fs.SkipDir && d.IsDir() {
			// Successfully skipped directory.
			err = nil
		}
		return err
	}

	file, err := fsys.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	dir, ok := file.(fs.ReadDirFile)
	if !ok {
		return &fs.PathError{Op: "readdir", Path: name, Err: errors.New("not implemented")}
	}

	for {
		dirs, err := boundReadDir(dir, name, maxReadDir)
		if err != nil {
			if err == io.EOF {
				break
			}
			// Second call, to report ReadDir error.
			err = walkDirFn(name, d, err)
			if err != nil {
				if err == fs.SkipDir && d.IsDir() {
					err = nil
				}
				return err
			}
		}

		for _, d1 := range dirs {
			name1 := path.Join(name, d1.Name())
			if err := boundWalkDirInternal(fsys, name1, d1, walkDirFn, maxReadDir); err != nil {
				if err == fs.SkipDir {
					break
				}
				return err
			}
		}
	}

	return nil
}

func BoundWalkDir(dir string, fn fs.WalkDirFunc, maxReadDir int) error {
	fsys := os.DirFS(dir)
	info, err := fs.Stat(fsys, ".")
	if err != nil {
		err = fn(dir, nil, err)
	} else {
		err = boundWalkDirInternal(fsys, ".", fs.FileInfoToDirEntry(info), fn, maxReadDir)
	}
	if err == fs.SkipDir || err == fs.SkipAll {
		return nil
	}
	return err
}
