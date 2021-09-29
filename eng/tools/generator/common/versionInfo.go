// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import "fmt"

type VersionInfo struct {
	LatestVersion string
	NewVersion    string
}

func (i VersionInfo) String() string {
	return fmt.Sprintf("Latest version: %s, new version: %s", i.LatestVersion, i.NewVersion)
}
