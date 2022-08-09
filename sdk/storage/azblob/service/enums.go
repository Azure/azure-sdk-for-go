//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

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

// ListContainersIncludeType defines values for ListContainersIncludeType
type ListContainersIncludeType string

const (
	ListContainersIncludeTypeMetadata ListContainersIncludeType = "metadata"
	ListContainersIncludeTypeDeleted  ListContainersIncludeType = "deleted"
	ListContainersIncludeTypeSystem   ListContainersIncludeType = "system"
)

// PossibleListContainersIncludeTypeValues returns the possible values for the ListContainersIncludeType const type.
func PossibleListContainersIncludeTypeValues() []ListContainersIncludeType {
	return []ListContainersIncludeType{
		ListContainersIncludeTypeMetadata,
		ListContainersIncludeTypeDeleted,
		ListContainersIncludeTypeSystem,
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
