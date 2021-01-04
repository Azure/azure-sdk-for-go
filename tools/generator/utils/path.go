package utils

import "strings"

func NormalizePath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
