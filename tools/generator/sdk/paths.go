// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sdk

import (
	"path/filepath"
)

func ServicesPath(root string) string {
	return filepath.Join(root, services)
}

func ChangelogPath(pkg string) string {
	return filepath.Join(pkg, changelogFilename)
}

func MetadataPath(pkg string) string {
	return filepath.Join(pkg, metadataFilename)
}

func VersionGoPath(root string) string {
	return filepath.Join(root, relativeVersionGo)
}
