//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// ListFilesIncludeType defines values for ListFilesIncludeType
type ListFilesIncludeType = generated.ListFilesIncludeType

const (
	ListFilesIncludeTypeTimestamps    ListFilesIncludeType = generated.ListFilesIncludeTypeTimestamps
	ListFilesIncludeTypeETag          ListFilesIncludeType = generated.ListFilesIncludeTypeEtag
	ListFilesIncludeTypeAttributes    ListFilesIncludeType = generated.ListFilesIncludeTypeAttributes
	ListFilesIncludeTypePermissionKey ListFilesIncludeType = generated.ListFilesIncludeTypePermissionKey
)

// PossibleListFilesIncludeTypeValues returns the possible values for the ListFilesIncludeType const type.
func PossibleListFilesIncludeTypeValues() []ListFilesIncludeType {
	return generated.PossibleListFilesIncludeTypeValues()
}
