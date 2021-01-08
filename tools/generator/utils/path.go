package utils

import "strings"

// NormalizePath normalizes the path by replacing \ with /
func NormalizePath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
