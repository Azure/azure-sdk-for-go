//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

const (
	PublicAccessTypeBlob      PublicAccessType = "blob"
	PublicAccessTypeContainer PublicAccessType = "container"
)

// PossiblePublicAccessTypeValues returns the possible values for the PublicAccessType const type.
func PossiblePublicAccessTypeValues() []PublicAccessType {
	return []PublicAccessType{
		PublicAccessTypeBlob,
		PublicAccessTypeContainer,
	}
}

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

// SKUName defines values for SkuName - LRS, GRS, RAGRS, ZRS, Premium LRS
type SKUName string

const (
	SKUNameStandardLRS   SKUName = "Standard_LRS"
	SKUNameStandardGRS   SKUName = "Standard_GRS"
	SKUNameStandardRAGRS SKUName = "Standard_RAGRS"
	SKUNameStandardZRS   SKUName = "Standard_ZRS"
	SKUNamePremiumLRS    SKUName = "Premium_LRS"
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

// LeaseStatusType defines values for LeaseStatusType
type LeaseStatusType = generated.LeaseStatusType

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
