package pipeline

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"

	"magalu.cloud/core/utils"
)

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

		_ = utils.BoundWalkDir(root, func(p string, d fs.DirEntry, err error) error {
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
