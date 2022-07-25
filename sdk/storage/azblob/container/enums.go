//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

const (
	PublicAccessTypeBlob      = generated.PublicAccessTypeBlob
	PublicAccessTypeContainer = generated.PublicAccessTypeContainer
)

const (
	ListBlobsIncludeItemCopy                = generated.ListBlobsIncludeItemCopy
	ListBlobsIncludeItemDeleted             = generated.ListBlobsIncludeItemDeleted
	ListBlobsIncludeItemMetadata            = generated.ListBlobsIncludeItemMetadata
	ListBlobsIncludeItemSnapshots           = generated.ListBlobsIncludeItemSnapshots
	ListBlobsIncludeItemUncommittedblobs    = generated.ListBlobsIncludeItemUncommittedblobs
	ListBlobsIncludeItemVersions            = generated.ListBlobsIncludeItemVersions
	ListBlobsIncludeItemTags                = generated.ListBlobsIncludeItemTags
	ListBlobsIncludeItemImmutabilitypolicy  = generated.ListBlobsIncludeItemImmutabilitypolicy
	ListBlobsIncludeItemLegalhold           = generated.ListBlobsIncludeItemLegalhold
	ListBlobsIncludeItemDeletedwithversions = generated.ListBlobsIncludeItemDeletedwithversions
)

// PossibleListBlobsIncludeItemValues returns the possible values for the ListBlobsIncludeItem const type.
func PossibleListBlobsIncludeItemValues() []ListBlobsIncludeItem {
	return generated.PossibleListBlobsIncludeItemValues()
}
