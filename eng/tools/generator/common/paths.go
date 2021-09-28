// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import "path/filepath"

func ServicesPath(root string) string {
	return filepath.Join(root, Services)
}

func ChangelogPath(pkg string) string {
	return filepath.Join(pkg, ChangelogFilename)
}

func MetadataPath(pkg string) string {
	return filepath.Join(pkg, MetadataFilename)
}

func VersionGoPath(root string) string {
	return filepath.Join(root, RelativeVersionGo)
}
