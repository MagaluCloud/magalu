package pipeline

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
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

func boundWalkDir(fsys fs.FS, root string, fn fs.WalkDirFunc, maxReadDir int) error {
	info, err := fs.Stat(fsys, root)
	if err != nil {
		err = fn(root, nil, err)
	} else {
		err = boundWalkDirInternal(fsys, root, fs.FileInfoToDirEntry(info), fn, maxReadDir)
	}
	if err == fs.SkipDir || err == fs.SkipAll {
		return nil
	}
	return err
}

// WalkDirEntries recursively scans files/directories from a root directory
//
// checkPath() may be used to return fs.SkipDir or fs.SkipAll and control the walk process.
// If provided (non-nil), it's called before anything else. See fs.WalkDirFunc documentation.
// It may be used to omit hidden folders (ie: ".git") and the likes
//
// Each file/directory may contain an associated error, it may be ignored or keep going.
// By default, if no checkPath is provided, it keeps going.
func WalkDirEntriesBound(
	ctx context.Context,
	root string,
	checkPath fs.WalkDirFunc,
	maxReadDir int,
) (outputChan <-chan WalkDirEntry) {
	ch := make(chan WalkDirEntry)
	outputChan = ch

	logger := FromContext(ctx).Named("WalkDirEntries").With(
		"root", root,
		"outputChan", fmt.Sprintf("%#v", outputChan),
	)
	ctx = NewContext(ctx, logger)

	generator := func() {
		defer func() {
			logger.Info("closing output channel")
			close(ch)
		}()

		_ = boundWalkDir(os.DirFS(root), ".", func(p string, d fs.DirEntry, err error) error {
			if d == nil {
				return filepath.SkipDir
			}
			p = path.Join(root, p)
			if checkPath != nil {
				e := checkPath(p, d, err)
				if e != nil {
					logger.Debugw("checkPath != nil", "err", err, "path", p, "dirEntry", d)
					return e
				}
			}
			dir := NewSimpleWalkDirEntry(p, d, err)
			select {
			case <-ctx.Done():
				logger.Debugw("context.Done()", "err", ctx.Err())
				return filepath.SkipAll

			case ch <- dir:
				logger.Debugw("entry", "err", err, "path", p, "dirEntry", d)
				return nil
			}
		}, maxReadDir)
		logger.Debug("finished walking entries")
	}

	logger.Info("start", root)
	go generator()
	return
}
