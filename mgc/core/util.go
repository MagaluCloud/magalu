package core

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
)

const FILE_PERMISSION = 0777 // TODO: investigate how to lower permission

// Copied code from https://pkg.go.dev/os#UserConfigDir but modified to treat
// Darwin the same way as Unix by setting to "~/.config"
func BuildMGCPath() (string, error) {
	dir := ""
	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("AppData")
		if dir == "" {
			return "", errors.New("%AppData% is not defined")
		}

	default: // Unix
		dir = os.Getenv("XDG_CONFIG_HOME")
		if dir == "" {
			dir = os.Getenv("HOME")
			if dir == "" {
				return "", errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
			}
			dir = path.Join(dir, ".config")
		}
	}
	mgcDir := path.Join(dir, "mgc")
	if err := os.MkdirAll(mgcDir, FILE_PERMISSION); err != nil {
		return "", fmt.Errorf("Error creating mgc dir at %s: %w", mgcDir, err)
	}
	return mgcDir, nil
}

func BuildMGCFilePath(filename string) (string, error) {
	mgcDir, err := BuildMGCPath()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(mgcDir, FILE_PERMISSION); err != nil {
		return "", fmt.Errorf("Error creating mgc dir at %s: %w", mgcDir, err)
	}
	return path.Join(mgcDir, filename), nil
}
