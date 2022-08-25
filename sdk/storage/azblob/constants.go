//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

// PublicAccessType defines values for AccessType - private (default) or blob or container
type PublicAccessType = generated.PublicAccessType

const (
	PublicAccessTypeBlob      PublicAccessType = generated.PublicAccessTypeBlob
	PublicAccessTypeContainer PublicAccessType = generated.PublicAccessTypeContainer
)

// PossiblePublicAccessTypeValues returns the possible values for the PublicAccessType const type.
func PossiblePublicAccessTypeValues() []PublicAccessType {
	return generated.PossiblePublicAccessTypeValues()
}

// DeleteSnapshotsOptionType defines values for DeleteSnapshotsOptionType
type DeleteSnapshotsOptionType = generated.DeleteSnapshotsOptionType

const (
	DeleteSnapshotsOptionTypeInclude DeleteSnapshotsOptionType = generated.DeleteSnapshotsOptionTypeInclude
	DeleteSnapshotsOptionTypeOnly    DeleteSnapshotsOptionType = generated.DeleteSnapshotsOptionTypeOnly
)

// PossibleDeleteSnapshotsOptionTypeValues returns the possible values for the DeleteSnapshotsOptionType const type.
func PossibleDeleteSnapshotsOptionTypeValues() []DeleteSnapshotsOptionType {
	return generated.PossibleDeleteSnapshotsOptionTypeValues()
}

const (
	ListBlobsIncludeItemCopy                ListBlobsIncludeItem = "copy"
	ListBlobsIncludeItemDeleted             ListBlobsIncludeItem = "deleted"
	ListBlobsIncludeItemMetadata            ListBlobsIncludeItem = "metadata"
	ListBlobsIncludeItemSnapshots           ListBlobsIncludeItem = "snapshots"
	ListBlobsIncludeItemUncommittedblobs    ListBlobsIncludeItem = "uncommittedblobs"
	ListBlobsIncludeItemVersions            ListBlobsIncludeItem = "versions"
	ListBlobsIncludeItemTags                ListBlobsIncludeItem = "tags"
	ListBlobsIncludeItemImmutabilitypolicy  ListBlobsIncludeItem = "immutabilitypolicy"
	ListBlobsIncludeItemLegalhold           ListBlobsIncludeItem = "legalhold"
	ListBlobsIncludeItemDeletedwithversions ListBlobsIncludeItem = "deletedwithversions"
)

// PossibleListBlobsIncludeItemValues returns the possible values for the ListBlobsIncludeItem const type.
func PossibleListBlobsIncludeItemValues() []ListBlobsIncludeItem {
	return []ListBlobsIncludeItem{
		ListBlobsIncludeItemCopy,
		ListBlobsIncludeItemDeleted,
		ListBlobsIncludeItemMetadata,
		ListBlobsIncludeItemSnapshots,
		ListBlobsIncludeItemUncommittedblobs,
		ListBlobsIncludeItemVersions,
		ListBlobsIncludeItemTags,
		ListBlobsIncludeItemImmutabilitypolicy,
		ListBlobsIncludeItemLegalhold,
		ListBlobsIncludeItemDeletedwithversions,
	}
}
