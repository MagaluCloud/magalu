//go:build !embed

package openapi

func GetEmbedLoader() Loader {
	return nil
}
