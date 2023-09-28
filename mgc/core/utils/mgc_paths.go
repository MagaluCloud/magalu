package utils

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
)

// These two const are needed because only 07xx (executable) directories can be accessed
const FILE_PERMISSION = 0644
const DIR_PERMISSION = 0744

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
			home := os.Getenv("HOME")
			if home != "" {
				dir = path.Join(home, ".config")
			}
		}
		if dir == "" {
			return "", errors.New("neither $XDG_CONFIG_HOME nor $HOME are defined")
		}
	}
	mgcDir := path.Join(dir, "mgc")
	if err := os.MkdirAll(mgcDir, DIR_PERMISSION); err != nil {
		return "", fmt.Errorf("Error creating mgc dir at %s: %w", mgcDir, err)
	}
	return mgcDir, nil
}

func BuildMGCFilePath(filename string) (string, error) {
	mgcDir, err := BuildMGCPath()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(mgcDir, DIR_PERMISSION); err != nil {
		return "", fmt.Errorf("Error creating mgc dir at %s: %w", mgcDir, err)
	}
	return path.Join(mgcDir, filename), nil
}
