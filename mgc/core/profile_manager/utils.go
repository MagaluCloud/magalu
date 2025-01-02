package profile_manager

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/core/utils"
)

var isProfileNameValid = regexp.MustCompile(`^[\w-]*.$`).MatchString

func sanitizePath(p string) string {
	pathEntries := strings.Split(p, string(os.PathSeparator))

	if len(pathEntries) == 1 {
		return pathEntries[0]
	}

	result := []string{}
	for _, entry := range pathEntries {
		if entry != "" && entry != ".." && entry != "." {
			result = append(result, entry)
		}
	}

	return strings.Join(result, "/")
}

func checkProfileName(name string) error {
	if name == currentProfileNameFile {
		return errorNameNotAllowed
	}
	if !isProfileNameValid(name) {
		return errorInvalidName
	}

	return nil
}

func buildMGCPath() (string, error) {
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
	if err := os.MkdirAll(mgcDir, utils.DIR_PERMISSION); err != nil {
		return "", fmt.Errorf("Error creating mgc dir at %s: %w", mgcDir, err)
	}
	return mgcDir, nil
}
