//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

// AccountKind defines values for AccountKind
type AccountKind = generated.AccountKind

const (
	AccountKindStorage          = generated.AccountKindStorage
	AccountKindBlobStorage      = generated.AccountKindBlobStorage
	AccountKindStorageV2        = generated.AccountKindStorageV2
	AccountKindFileStorage      = generated.AccountKindFileStorage
	AccountKindBlockBlobStorage = generated.AccountKindBlockBlobStorage
)

// PossibleAccountKindValues returns the possible values for the AccountKind const type.
func PossibleAccountKindValues() []AccountKind {
	return []AccountKind{
		AccountKindStorage,
		AccountKindBlobStorage,
		AccountKindStorageV2,
		AccountKindFileStorage,
		AccountKindBlockBlobStorage,
	}
}

// ArchiveStatus defines values for ArchiveStatus
type ArchiveStatus = generated.ArchiveStatus

const (
	ArchiveStatusRehydratePendingToCool = generated.ArchiveStatusRehydratePendingToCool
	ArchiveStatusRehydratePendingToHot  = generated.ArchiveStatusRehydratePendingToHot
)

// PossibleArchiveStatusValues returns the possible values for the ArchiveStatus const type.
func PossibleArchiveStatusValues() []ArchiveStatus {
	return []ArchiveStatus{
		ArchiveStatusRehydratePendingToCool,
		ArchiveStatusRehydratePendingToHot,
	}
}

// BlobDeleteType defines values for BlobDeleteType
type BlobDeleteType = generated.BlobDeleteType

const (
	BlobDeleteTypeNone      = generated.BlobDeleteTypeNone
	BlobDeleteTypePermanent = generated.BlobDeleteTypePermanent
)

// PossibleBlobDeleteTypeValues returns the possible values for the BlobDeleteType const type.
func PossibleBlobDeleteTypeValues() []BlobDeleteType {
	return []BlobDeleteType{
		BlobDeleteTypeNone,
		BlobDeleteTypePermanent,
	}
}

// BlobExpiryOptions defines values for BlobExpiryOptions
type BlobExpiryOptions = generated.BlobExpiryOptions

const (
	BlobExpiryOptionsAbsolute           = generated.BlobExpiryOptionsAbsolute
	BlobExpiryOptionsNeverExpire        = generated.BlobExpiryOptionsNeverExpire
	BlobExpiryOptionsRelativeToCreation = generated.BlobExpiryOptionsRelativeToCreation
	BlobExpiryOptionsRelativeToNow      = generated.BlobExpiryOptionsRelativeToNow
)

// PossibleBlobExpiryOptionsValues returns the possible values for the BlobExpiryOptions const type.
func PossibleBlobExpiryOptionsValues() []BlobExpiryOptions {
	return []BlobExpiryOptions{
		BlobExpiryOptionsAbsolute,
		BlobExpiryOptionsNeverExpire,
		BlobExpiryOptionsRelativeToCreation,
		BlobExpiryOptionsRelativeToNow,
	}
}

// BlobGeoReplicationStatus - The status of the secondary location
type BlobGeoReplicationStatus = generated.BlobGeoReplicationStatus

const (
	BlobGeoReplicationStatusLive        = generated.BlobGeoReplicationStatusLive
	BlobGeoReplicationStatusBootstrap   = generated.BlobGeoReplicationStatusBootstrap
	BlobGeoReplicationStatusUnavailable = generated.BlobGeoReplicationStatusUnavailable
)

// PossibleBlobGeoReplicationStatusValues returns the possible values for the BlobGeoReplicationStatus const type.
func PossibleBlobGeoReplicationStatusValues() []BlobGeoReplicationStatus {
	return []BlobGeoReplicationStatus{
		BlobGeoReplicationStatusLive,
		BlobGeoReplicationStatusBootstrap,
		BlobGeoReplicationStatusUnavailable,
	}
}

// BlobType defines values for BlobType
type BlobType = generated.BlobType

const (
	BlobTypeBlockBlob  = generated.BlobTypeBlockBlob
	BlobTypePageBlob   = generated.BlobTypePageBlob
	BlobTypeAppendBlob = generated.BlobTypeAppendBlob
)

// PossibleBlobTypeValues returns the possible values for the BlobType const type.
func PossibleBlobTypeValues() []BlobType {
	return []BlobType{
		BlobTypeBlockBlob,
		BlobTypePageBlob,
		BlobTypeAppendBlob,
	}
}

// LeaseDurationType defines values for LeaseDurationType
type LeaseDurationType = generated.LeaseDurationType

const (
	LeaseDurationTypeInfinite = generated.LeaseDurationTypeInfinite
	LeaseDurationTypeFixed    = generated.LeaseDurationTypeFixed
)

// PossibleLeaseDurationTypeValues returns the possible values for the LeaseDurationType const type.
func PossibleLeaseDurationTypeValues() []LeaseDurationType {
	return []LeaseDurationType{
		LeaseDurationTypeInfinite,
		LeaseDurationTypeFixed,
	}
}

// LeaseStateType defines values for LeaseStateType
type LeaseStateType = generated.LeaseStateType

const (
	LeaseStateTypeAvailable = generated.LeaseStateTypeAvailable
	LeaseStateTypeLeased    = generated.LeaseStateTypeLeased
	LeaseStateTypeExpired   = generated.LeaseStateTypeExpired
	LeaseStateTypeBreaking  = generated.LeaseStateTypeBreaking
	LeaseStateTypeBroken    = generated.LeaseStateTypeBroken
)

// PossibleLeaseStateTypeValues returns the possible values for the LeaseStateType const type.
func PossibleLeaseStateTypeValues() []LeaseStateType {
	return []LeaseStateType{
		LeaseStateTypeAvailable,
		LeaseStateTypeLeased,
		LeaseStateTypeExpired,
		LeaseStateTypeBreaking,
		LeaseStateTypeBroken,
	}
}

// LeaseStatusType defines values for LeaseStatusType
type LeaseStatusType = generated.LeaseStatusType

const (
	LeaseStatusTypeLocked   = generated.LeaseStatusTypeLocked
	LeaseStatusTypeUnlocked = generated.LeaseStatusTypeUnlocked
)

// PossibleLeaseStatusTypeValues returns the possible values for the LeaseStatusType const type.
func PossibleLeaseStatusTypeValues() []LeaseStatusType {
	return []LeaseStatusType{
		LeaseStatusTypeLocked,
		LeaseStatusTypeUnlocked,
	}
}

// ListBlobsIncludeItem defines values for ListBlobsIncludeItem
type ListBlobsIncludeItem = generated.ListBlobsIncludeItem

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

// ListContainersIncludeType defines values for ListContainersIncludeType
type ListContainersIncludeType = generated.ListContainersIncludeType

const (
	ListContainersIncludeTypeMetadata = generated.ListContainersIncludeTypeMetadata
	ListContainersIncludeTypeDeleted  = generated.ListContainersIncludeTypeDeleted
	ListContainersIncludeTypeSystem   = generated.ListContainersIncludeTypeSystem
)

// PossibleListContainersIncludeTypeValues returns the possible values for the ListContainersIncludeType const type.
func PossibleListContainersIncludeTypeValues() []ListContainersIncludeType {
	return []ListContainersIncludeType{
		ListContainersIncludeTypeMetadata,
		ListContainersIncludeTypeDeleted,
		ListContainersIncludeTypeSystem,
	}
}

// PublicAccessType defines values for AccessType - private (default) or blob or container
type PublicAccessType = generated.PublicAccessType

const (
	PublicAccessTypeBlob      = generated.PublicAccessTypeBlob
	PublicAccessTypeContainer = generated.PublicAccessTypeContainer
)

// PossiblePublicAccessTypeValues returns the possible values for the PublicAccessType const type.
func PossiblePublicAccessTypeValues() []PublicAccessType {
	return generated.PossiblePublicAccessTypeValues()
}

// QueryFormatType - The quick query format type.
type QueryFormatType = generated.QueryFormatType

const (
	QueryFormatTypeDelimited = generated.QueryFormatTypeDelimited
	QueryFormatTypeJSON      = generated.QueryFormatTypeJSON
	QueryFormatTypeArrow     = generated.QueryFormatTypeArrow
	QueryFormatTypeParquet   = generated.QueryFormatTypeParquet
)

// PossibleQueryFormatTypeValues returns the possible values for the QueryFormatType const type.
func PossibleQueryFormatTypeValues() []QueryFormatType {
	return []QueryFormatType{
		QueryFormatTypeDelimited,
		QueryFormatTypeJSON,
		QueryFormatTypeArrow,
		QueryFormatTypeParquet,
	}
}

// SKUName defines values for SkuName - LRS, GRS, RAGRS, ZRS, Premium LRS
type SKUName = generated.SKUName

const (
	SKUNameStandardLRS   = generated.SKUNameStandardLRS
	SKUNameStandardGRS   = generated.SKUNameStandardGRS
	SKUNameStandardRAGRS = generated.SKUNameStandardRAGRS
	SKUNameStandardZRS   = generated.SKUNameStandardZRS
	SKUNamePremiumLRS    = generated.SKUNamePremiumLRS
)

// PossibleSKUNameValues returns the possible values for the SKUName const type.
func PossibleSKUNameValues() []SKUName {
	return []SKUName{
		SKUNameStandardLRS,
		SKUNameStandardGRS,
		SKUNameStandardRAGRS,
		SKUNameStandardZRS,
		SKUNamePremiumLRS,
	}
}
