package openapi

import (
	"embed"
	"io/fs"
)

//go:embed openapis/beta
var betaFolder embed.FS

func BetaSpecs() fs.FS {
	sub, err := fs.Sub(betaFolder, "openapis/beta")
	if err != nil {
		panic(err)
	}
	return sub
}
