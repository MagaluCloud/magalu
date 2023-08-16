package core

import (
	"fmt"
	"os"
	"path"
)

func buildMGCFilePath(filename string) (string, error) {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", fmt.Errorf("Neither $XDG_CONFIG_HOME nor $HOME are defined")
		}
		dir = path.Join(dir, ".config")
	}
	mgcDir := path.Join(dir, "mgc")
	if err := os.MkdirAll(mgcDir, FILE_PERMISSION); err != nil {
		return "", fmt.Errorf("Error creating mgc dir at %s: %w", mgcDir, err)
	}

	return path.Join(mgcDir, filename), nil
}
