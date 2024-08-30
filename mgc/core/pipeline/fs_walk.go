package pipeline

import (
	"context"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// WalkDirEntries -> FilterWalkDirEntries
// See fs.WalkDirFunc documentation
type WalkDirEntry interface {
	Path() string
	DirEntry() fs.DirEntry
	Err() error
}

type SimpleWalkDirEntry[T fs.DirEntry] struct {
	path   string
	Object T
	err    error
}

func (e *SimpleWalkDirEntry[T]) Path() string {
	return e.path
}

func (e *SimpleWalkDirEntry[T]) DirEntry() fs.DirEntry {
	return e.Object
}

func (e *SimpleWalkDirEntry[T]) Err() error {
	return e.err
}

func NewSimpleWalkDirEntry[T fs.DirEntry](path string, dirEntry T, err error) *SimpleWalkDirEntry[T] {
	return &SimpleWalkDirEntry[T]{path, dirEntry, err}
}

// Do not process any entry that name stars with "."
func WalkDirFilterHiddenDirs(path string, d fs.DirEntry, err error) error {
	if d == nil {
		return nil
	}
	if name := d.Name(); strings.HasPrefix(name, ".") && d.IsDir() {
		return fs.SkipDir
	}
	return nil
}

// WalkDirEntries recursively scans files/directories from a root directory
//
// checkPath() may be used to return fs.SkipDir or fs.SkipAll and control the walk process.
// If provided (non-nil), it's called before anything else. See fs.WalkDirFunc documentation.
// It may be used to omit hidden folders (ie: ".git") and the likes
//
// Each file/directory may contain an associated error, it may be ignored or keep going.
// By default, if no checkPath is provided, it keeps going.
func WalkDirEntries(
    ctx context.Context,
    root string,
    checkPath fs.WalkDirFunc,
) (<-chan WalkDirEntry) {
    ch := make(chan WalkDirEntry, 100) // Buffered channel
    logger := FromContext(ctx).Named("WalkDirEntries").With(
        "root", root,
    )
    ctx = NewContext(ctx, logger)

    go func() {
        defer close(ch)
        logger.Info("Starting directory walk")

        err := fs.WalkDir(os.DirFS(root), ".", func(p string, d fs.DirEntry, err error) error {
            if d == nil {
                return filepath.SkipDir
            }
            p = path.Join(root, p)
            if checkPath != nil {
                if e := checkPath(p, d, err); e != nil {
                    logger.Debugw("checkPath error", "err", e, "path", p)
                    return e
                }
            }
            dir := NewSimpleWalkDirEntry(p, d, err)
            select {
            case <-ctx.Done():
                logger.Debug("Context cancelled", "err", ctx.Err())
                return filepath.SkipAll
            case ch <- dir:
                return nil
            }
        })

        if err != nil {
            logger.Error("Error during directory walk", "err", err)
        }
    }()

    return ch
}
