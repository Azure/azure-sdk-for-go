//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem

import (
	"net/url"
	"strings"
)

// escapeSplitPaths is utility function to escape the individual strings by eliminating "/" in the path
func escapeSplitPaths(filePath string) string {
	names := strings.Split(filePath, "/")
	path := make([]string, len(names))
	for i, name := range names {
		path[i] = url.PathEscape(name)
	}
	escapedPathUrl := strings.Join(path, "/")
	return escapedPathUrl
}
