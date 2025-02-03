package common

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	mgcSchemaPkg "github.com/MagaluCloud/magalu/mgc/core/schema"
)

func GetAbsSystemURI(uri mgcSchemaPkg.URI) (mgcSchemaPkg.URI, error) {
	path := uri.String()

	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return uri, err
		}
		path = homeDir + path[1:]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return uri, errors.New("invalid local path")
	}

	return mgcSchemaPkg.URI(absPath), nil
}
