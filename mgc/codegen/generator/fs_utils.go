package generator

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"

	"github.com/spf13/afero"
)

func createDir(ctx *GeneratorContext, dirname string) (err error) {
	if err = ctx.FS.Mkdir(dirname, dirMode); err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	} else if err != nil {
		ctx.Reporter.Generate(dirname, "exists")
	} else {
		ctx.Reporter.Generate(dirname, "created")
	}

	var ok bool
	if ok, err = afero.IsDir(ctx.FS, dirname); err != nil {
		return
	} else if !ok {
		return fmt.Errorf("%s: %w", dirname, errNotDir)
	}

	return
}

func moveFileToBackup(afs afero.Fs, name string) (err error) {
	return afs.Rename(name, fmt.Sprintf(backupFmt, name))
}

func createFileSafe(afs afero.Fs, name string, flags int, perm fs.FileMode) (f afero.File, err error) {
	for {
		f, err = afs.OpenFile(name, flags, perm)
		if err == nil {
			return
		}
		if !errors.Is(err, fs.ErrExist) {
			return
		}

		bErr := moveFileToBackup(afs, name)
		if bErr != nil {
			return
		}
	}
}

func replaceFileIfNeeded(ctx *GeneratorContext, name string, newData []byte) (err error) {
	oldData, rErr := afero.ReadFile(ctx.FS, name)
	if rErr == nil {
		if bytes.Equal(newData, oldData) {
			ctx.Reporter.Generate(name, "up to date")
			return nil
		}
	}

	f, err := createFileSafe(ctx.FS, name, createFileFlags, fileMode)
	if err != nil {
		ctx.Reporter.Error(name, "failed to create file", err)
		return
	}
	_, err = f.Write(newData)
	if cErr := f.Close(); err == nil && cErr != nil {
		err = cErr
		ctx.Reporter.Error(name, "failed to close file", err)
	} else if err != nil {
		ctx.Reporter.Error(name, "failed to write file", err)
	} else {
		ctx.Reporter.Generate(name, "created")
	}
	return err
}
