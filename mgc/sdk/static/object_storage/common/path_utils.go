package common

import (
	"net/url"
	"path"
	"strings"
)

// NormalizeFilePath normalizes file paths for consistent handling across different systems
// It decodes URL-encoded characters, normalizes directory separators, and removes various prefixes
func NormalizeFilePath(srcPath string) string {
	decodedPath, err := url.PathUnescape(srcPath)
	if err != nil {
		decodedPath = srcPath
	}
	unifiedPath := strings.ReplaceAll(decodedPath, "\\", "/")
	unifiedPath = strings.TrimPrefix(unifiedPath, "path://")
	cleanPath := path.Clean(unifiedPath)

	return cleanPath
}

// ExtractFileName extracts just the filename from a path
func ExtractFileName(p string) string {
	normalized := NormalizeFilePath(p)
	return path.Base(normalized)
}

// GetRelativePath gets the relative path between a base path and a file path
func GetRelativePath(basePath, filePath string) string {
	normalizedBase := NormalizeFilePath(basePath)
	normalizedFile := NormalizeFilePath(filePath)

	if !strings.HasSuffix(normalizedBase, "/") {
		normalizedBase += "/"
	}

	if strings.HasPrefix(normalizedFile, normalizedBase) {
		return strings.TrimPrefix(normalizedFile, normalizedBase)
	}

	return normalizedFile
}