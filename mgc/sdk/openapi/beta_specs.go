package openapi

import (
	"embed"
	"io/fs"
)

// The beta specs ship with the CLI in their own folder. The official loader
// embeds only openapis/*.yaml, so they never become official commands: they are
// read by the preview package, under `mgc beta`.
//
//go:embed openapis/beta
var betaFolder embed.FS

// BetaSpecs is the folder with the compiled beta specs, one
// "<module>.openapi.yaml" per module.
func BetaSpecs() fs.FS {
	sub, err := fs.Sub(betaFolder, "openapis/beta")
	if err != nil {
		panic(err) // the path is the constant embedded above
	}
	return sub
}
