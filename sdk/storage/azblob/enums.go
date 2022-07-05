package azblob

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

type AccountKind = generated.AccountKind

const (
	AccountKindStorage          AccountKind = "Storage"
	AccountKindBlobStorage      AccountKind = "BlobStorage"
	AccountKindStorageV2        AccountKind = "StorageV2"
	AccountKindFileStorage      AccountKind = "FileStorage"
	AccountKindBlockBlobStorage AccountKind = "BlockBlobStorage"
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

type ArchiveStatus = generated.ArchiveStatus

const (
	ArchiveStatusRehydratePendingToCool ArchiveStatus = "rehydrate-pending-to-cool"
	ArchiveStatusRehydratePendingToHot  ArchiveStatus = "rehydrate-pending-to-hot"
)

// PossibleArchiveStatusValues returns the possible values for the ArchiveStatus const type.
func PossibleArchiveStatusValues() []ArchiveStatus {
	return []ArchiveStatus{
		ArchiveStatusRehydratePendingToCool,
		ArchiveStatusRehydratePendingToHot,
	}
}

type BlobDeleteType = generated.BlobDeleteType

const (
	BlobDeleteTypeNone      BlobDeleteType = "None"
	BlobDeleteTypePermanent BlobDeleteType = "Permanent"
)

// PossibleBlobDeleteTypeValues returns the possible values for the BlobDeleteType const type.
func PossibleBlobDeleteTypeValues() []BlobDeleteType {
	return []BlobDeleteType{
		BlobDeleteTypeNone,
		BlobDeleteTypePermanent,
	}
}

type BlobExpiryOptions = generated.BlobExpiryOptions

const (
	BlobExpiryOptionsAbsolute           BlobExpiryOptions = "Absolute"
	BlobExpiryOptionsNeverExpire        BlobExpiryOptions = "NeverExpire"
	BlobExpiryOptionsRelativeToCreation BlobExpiryOptions = "RelativeToCreation"
	BlobExpiryOptionsRelativeToNow      BlobExpiryOptions = "RelativeToNow"
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
	BlobGeoReplicationStatusLive        BlobGeoReplicationStatus = "live"
	BlobGeoReplicationStatusBootstrap   BlobGeoReplicationStatus = "bootstrap"
	BlobGeoReplicationStatusUnavailable BlobGeoReplicationStatus = "unavailable"
)

// PossibleBlobGeoReplicationStatusValues returns the possible values for the BlobGeoReplicationStatus const type.
func PossibleBlobGeoReplicationStatusValues() []BlobGeoReplicationStatus {
	return []BlobGeoReplicationStatus{
		BlobGeoReplicationStatusLive,
		BlobGeoReplicationStatusBootstrap,
		BlobGeoReplicationStatusUnavailable,
	}
}

type BlobType = generated.BlobType

const (
	BlobTypeBlockBlob  BlobType = "BlockBlob"
	BlobTypePageBlob   BlobType = "PageBlob"
	BlobTypeAppendBlob BlobType = "AppendBlob"
)

// PossibleBlobTypeValues returns the possible values for the BlobType const type.
func PossibleBlobTypeValues() []BlobType {
	return []BlobType{
		BlobTypeBlockBlob,
		BlobTypePageBlob,
		BlobTypeAppendBlob,
	}
}

type LeaseDurationType = generated.LeaseDurationType

const (
	LeaseDurationTypeInfinite LeaseDurationType = "infinite"
	LeaseDurationTypeFixed    LeaseDurationType = "fixed"
)

// PossibleLeaseDurationTypeValues returns the possible values for the LeaseDurationType const type.
func PossibleLeaseDurationTypeValues() []LeaseDurationType {
	return []LeaseDurationType{
		LeaseDurationTypeInfinite,
		LeaseDurationTypeFixed,
	}
}

type LeaseStateType = generated.LeaseStateType

const (
	LeaseStateTypeAvailable LeaseStateType = "available"
	LeaseStateTypeLeased    LeaseStateType = "leased"
	LeaseStateTypeExpired   LeaseStateType = "expired"
	LeaseStateTypeBreaking  LeaseStateType = "breaking"
	LeaseStateTypeBroken    LeaseStateType = "broken"
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

type LeaseStatusType = generated.LeaseStatusType

const (
	LeaseStatusTypeLocked   LeaseStatusType = generated.LeaseStatusTypeLocked
	LeaseStatusTypeUnlocked LeaseStatusType = generated.LeaseStatusTypeUnlocked
)

// PossibleLeaseStatusTypeValues returns the possible values for the LeaseStatusType const type.
func PossibleLeaseStatusTypeValues() []LeaseStatusType {
	return []LeaseStatusType{
		LeaseStatusTypeLocked,
		LeaseStatusTypeUnlocked,
	}
}

type ListBlobsIncludeItem = generated.ListBlobsIncludeItem

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

type ListContainersIncludeType = generated.ListContainersIncludeType

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

// PossiblePremiumPageBlobAccessTierValues returns the possible values for the PremiumPageBlobAccessTier const type.

type PublicAccessType = generated.PublicAccessType

const (
	PublicAccessTypeBlob      PublicAccessType = generated.PublicAccessTypeBlob
	PublicAccessTypeContainer PublicAccessType = generated.PublicAccessTypeContainer
)

// PossiblePublicAccessTypeValues returns the possible values for the PublicAccessType const type.
func PossiblePublicAccessTypeValues() []PublicAccessType {
	return generated.PossiblePublicAccessTypeValues()
}

// QueryFormatType - The quick query format type.
type QueryFormatType = generated.QueryFormatType

const (
	QueryFormatTypeDelimited QueryFormatType = "delimited"
	QueryFormatTypeJSON      QueryFormatType = "json"
	QueryFormatTypeArrow     QueryFormatType = "arrow"
	QueryFormatTypeParquet   QueryFormatType = "parquet"
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

type SKUName = generated.SKUName

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
