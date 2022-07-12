//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

const (
	PublicAccessTypeBlob      PublicAccessType = generated.PublicAccessTypeBlob
	PublicAccessTypeContainer PublicAccessType = generated.PublicAccessTypeContainer
)

// ListBlobsIncludeItem enum
type ListBlobsIncludeItem = generated.ListBlobsIncludeItem

const (
	ListBlobsIncludeItemCopy                ListBlobsIncludeItem = generated.ListBlobsIncludeItemCopy
	ListBlobsIncludeItemDeleted             ListBlobsIncludeItem = generated.ListBlobsIncludeItemDeleted
	ListBlobsIncludeItemMetadata            ListBlobsIncludeItem = generated.ListBlobsIncludeItemMetadata
	ListBlobsIncludeItemSnapshots           ListBlobsIncludeItem = generated.ListBlobsIncludeItemSnapshots
	ListBlobsIncludeItemUncommittedblobs    ListBlobsIncludeItem = generated.ListBlobsIncludeItemUncommittedblobs
	ListBlobsIncludeItemVersions            ListBlobsIncludeItem = generated.ListBlobsIncludeItemVersions
	ListBlobsIncludeItemTags                ListBlobsIncludeItem = generated.ListBlobsIncludeItemTags
	ListBlobsIncludeItemImmutabilitypolicy  ListBlobsIncludeItem = generated.ListBlobsIncludeItemImmutabilitypolicy
	ListBlobsIncludeItemLegalhold           ListBlobsIncludeItem = generated.ListBlobsIncludeItemLegalhold
	ListBlobsIncludeItemDeletedwithversions ListBlobsIncludeItem = generated.ListBlobsIncludeItemDeletedwithversions
)

// PossibleListBlobsIncludeItemValues returns the possible values for the ListBlobsIncludeItem const type.
func PossibleListBlobsIncludeItemValues() []ListBlobsIncludeItem {
	return generated.PossibleListBlobsIncludeItemValues()
}
