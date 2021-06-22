// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package sdk

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/internal/utils"
)

func IsPreviewPackage(pkgName string) bool {
	return strings.HasPrefix(utils.NormalizePath(pkgName), "services/preview/")
}
