package common

import (
	"net/url"
	"strings"
)

// NormalizeFilePath normalizes file paths for consistent handling across different systems
// It decodes URL-encoded characters, normalizes directory separators, and removes various prefixes
func NormalizeFilePath(srcPath string) string {
	decodedPath, err := url.PathUnescape(srcPath)
	if err != nil {
		decodedPath = srcPath
	}	
	normalizedPath := strings.ReplaceAll(decodedPath, "\\", "/")
	normalizedPath = strings.TrimPrefix(normalizedPath, "path://")
	normalizedPath = strings.TrimPrefix(normalizedPath, "./")
	normalizedPath = strings.TrimPrefix(normalizedPath, ".")
	
	return normalizedPath
}

// ExtractFileName extracts just the filename from a path
func ExtractFileName(path string) string {
	normalizedPath := NormalizeFilePath(path)

	var fileName string
	if lastSlash := strings.LastIndex(normalizedPath, "/"); lastSlash >= 0 {
		fileName = normalizedPath[lastSlash+1:]
	} else {
		fileName = normalizedPath
	}
	
	return fileName
}

// GetRelativePath gets the relative path between a base path and a file path
func GetRelativePath(basePath, filePath string) string {
	normalizedBasePath := NormalizeFilePath(basePath)
	normalizedFilePath := NormalizeFilePath(filePath)
	
	relPath := strings.TrimPrefix(normalizedFilePath, normalizedBasePath)
	relPath = strings.TrimPrefix(relPath, "/")
	
	return relPath
}