// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package utils

import (
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/sdk"
)

func CanIncludeInMinor(r autorest.ChangelogResult) bool {
	if sdk.IsPreviewPackage(r.PackageName) {
		return true
	}
	return !r.Changelog.HasBreakingChanges()
}
