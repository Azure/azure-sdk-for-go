// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package utils

import "strings"

// NormalizePath normalizes the path by replacing \ with /
func NormalizePath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
