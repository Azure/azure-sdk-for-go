// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package generated

type AccessTier string

const (
	AccessTierArchive AccessTier = "Archive"
	AccessTierCold    AccessTier = "Cold"
	AccessTierCool    AccessTier = "Cool"
	AccessTierHot     AccessTier = "Hot"
	AccessTierP10     AccessTier = "P10"
	AccessTierP15     AccessTier = "P15"
	AccessTierP20     AccessTier = "P20"
	AccessTierP30     AccessTier = "P30"
	AccessTierP4      AccessTier = "P4"
	AccessTierP40     AccessTier = "P40"
	AccessTierP50     AccessTier = "P50"
	AccessTierP6      AccessTier = "P6"
	AccessTierP60     AccessTier = "P60"
	AccessTierP70     AccessTier = "P70"
	AccessTierP80     AccessTier = "P80"
	AccessTierPremium AccessTier = "Premium"
)

// PossibleAccessTierValues returns the possible values for the AccessTier const type.
func PossibleAccessTierValues() []AccessTier {
	return []AccessTier{
		AccessTierArchive,
		AccessTierCold,
		AccessTierCool,
		AccessTierHot,
		AccessTierP10,
		AccessTierP15,
		AccessTierP20,
		AccessTierP30,
		AccessTierP4,
		AccessTierP40,
		AccessTierP50,
		AccessTierP6,
		AccessTierP60,
		AccessTierP70,
		AccessTierP80,
		AccessTierPremium,
	}
}

type AccountKind string

const (
	AccountKindBlobStorage      AccountKind = "BlobStorage"
	AccountKindBlockBlobStorage AccountKind = "BlockBlobStorage"
	AccountKindFileStorage      AccountKind = "FileStorage"
	AccountKindStorage          AccountKind = "Storage"
	AccountKindStorageV2        AccountKind = "StorageV2"
)

// PossibleAccountKindValues returns the possible values for the AccountKind const type.
func PossibleAccountKindValues() []AccountKind {
	return []AccountKind{
		AccountKindBlobStorage,
		AccountKindBlockBlobStorage,
		AccountKindFileStorage,
		AccountKindStorage,
		AccountKindStorageV2,
	}
}

type ArchiveStatus string

const (
	ArchiveStatusRehydratePendingToCold ArchiveStatus = "rehydrate-pending-to-cold"
	ArchiveStatusRehydratePendingToCool ArchiveStatus = "rehydrate-pending-to-cool"
	ArchiveStatusRehydratePendingToHot  ArchiveStatus = "rehydrate-pending-to-hot"
)

// PossibleArchiveStatusValues returns the possible values for the ArchiveStatus const type.
func PossibleArchiveStatusValues() []ArchiveStatus {
	return []ArchiveStatus{
		ArchiveStatusRehydratePendingToCold,
		ArchiveStatusRehydratePendingToCool,
		ArchiveStatusRehydratePendingToHot,
	}
}

type BlobCopySourceTags string

const (
	BlobCopySourceTagsCOPY    BlobCopySourceTags = "COPY"
	BlobCopySourceTagsREPLACE BlobCopySourceTags = "REPLACE"
)

// PossibleBlobCopySourceTagsValues returns the possible values for the BlobCopySourceTags const type.
func PossibleBlobCopySourceTagsValues() []BlobCopySourceTags {
	return []BlobCopySourceTags{
		BlobCopySourceTagsCOPY,
		BlobCopySourceTagsREPLACE,
	}
}

// BlobGeoReplicationStatus - The status of the secondary location
type BlobGeoReplicationStatus string

const (
	BlobGeoReplicationStatusBootstrap   BlobGeoReplicationStatus = "bootstrap"
	BlobGeoReplicationStatusLive        BlobGeoReplicationStatus = "live"
	BlobGeoReplicationStatusUnavailable BlobGeoReplicationStatus = "unavailable"
)

// PossibleBlobGeoReplicationStatusValues returns the possible values for the BlobGeoReplicationStatus const type.
func PossibleBlobGeoReplicationStatusValues() []BlobGeoReplicationStatus {
	return []BlobGeoReplicationStatus{
		BlobGeoReplicationStatusBootstrap,
		BlobGeoReplicationStatusLive,
		BlobGeoReplicationStatusUnavailable,
	}
}

type BlobType string

const (
	BlobTypeAppendBlob BlobType = "AppendBlob"
	BlobTypeBlockBlob  BlobType = "BlockBlob"
	BlobTypePageBlob   BlobType = "PageBlob"
)

// PossibleBlobTypeValues returns the possible values for the BlobType const type.
func PossibleBlobTypeValues() []BlobType {
	return []BlobType{
		BlobTypeAppendBlob,
		BlobTypeBlockBlob,
		BlobTypePageBlob,
	}
}

type BlockListType string

const (
	BlockListTypeAll         BlockListType = "all"
	BlockListTypeCommitted   BlockListType = "committed"
	BlockListTypeUncommitted BlockListType = "uncommitted"
)

// PossibleBlockListTypeValues returns the possible values for the BlockListType const type.
func PossibleBlockListTypeValues() []BlockListType {
	return []BlockListType{
		BlockListTypeAll,
		BlockListTypeCommitted,
		BlockListTypeUncommitted,
	}
}

type CopyStatusType string

const (
	CopyStatusTypeAborted CopyStatusType = "aborted"
	CopyStatusTypeFailed  CopyStatusType = "failed"
	CopyStatusTypePending CopyStatusType = "pending"
	CopyStatusTypeSuccess CopyStatusType = "success"
)

// PossibleCopyStatusTypeValues returns the possible values for the CopyStatusType const type.
func PossibleCopyStatusTypeValues() []CopyStatusType {
	return []CopyStatusType{
		CopyStatusTypeAborted,
		CopyStatusTypeFailed,
		CopyStatusTypePending,
		CopyStatusTypeSuccess,
	}
}

type DeleteSnapshotsOptionType string

const (
	DeleteSnapshotsOptionTypeInclude DeleteSnapshotsOptionType = "include"
	DeleteSnapshotsOptionTypeOnly    DeleteSnapshotsOptionType = "only"
)

// PossibleDeleteSnapshotsOptionTypeValues returns the possible values for the DeleteSnapshotsOptionType const type.
func PossibleDeleteSnapshotsOptionTypeValues() []DeleteSnapshotsOptionType {
	return []DeleteSnapshotsOptionType{
		DeleteSnapshotsOptionTypeInclude,
		DeleteSnapshotsOptionTypeOnly,
	}
}

type DeleteType string

const (
	DeleteTypeNone      DeleteType = "None"
	DeleteTypePermanent DeleteType = "Permanent"
)

// PossibleDeleteTypeValues returns the possible values for the DeleteType const type.
func PossibleDeleteTypeValues() []DeleteType {
	return []DeleteType{
		DeleteTypeNone,
		DeleteTypePermanent,
	}
}

type EncryptionAlgorithmType string

const (
	EncryptionAlgorithmTypeAES256 EncryptionAlgorithmType = "AES256"
	EncryptionAlgorithmTypeNone   EncryptionAlgorithmType = "None"
)

// PossibleEncryptionAlgorithmTypeValues returns the possible values for the EncryptionAlgorithmType const type.
func PossibleEncryptionAlgorithmTypeValues() []EncryptionAlgorithmType {
	return []EncryptionAlgorithmType{
		EncryptionAlgorithmTypeAES256,
		EncryptionAlgorithmTypeNone,
	}
}

type ExpiryOptions string

const (
	ExpiryOptionsAbsolute           ExpiryOptions = "Absolute"
	ExpiryOptionsNeverExpire        ExpiryOptions = "NeverExpire"
	ExpiryOptionsRelativeToCreation ExpiryOptions = "RelativeToCreation"
	ExpiryOptionsRelativeToNow      ExpiryOptions = "RelativeToNow"
)

// PossibleExpiryOptionsValues returns the possible values for the ExpiryOptions const type.
func PossibleExpiryOptionsValues() []ExpiryOptions {
	return []ExpiryOptions{
		ExpiryOptionsAbsolute,
		ExpiryOptionsNeverExpire,
		ExpiryOptionsRelativeToCreation,
		ExpiryOptionsRelativeToNow,
	}
}

type FileShareTokenIntent string

const (
	FileShareTokenIntentBackup FileShareTokenIntent = "backup"
)

// PossibleFileShareTokenIntentValues returns the possible values for the FileShareTokenIntent const type.
func PossibleFileShareTokenIntentValues() []FileShareTokenIntent {
	return []FileShareTokenIntent{
		FileShareTokenIntentBackup,
	}
}

type FilterBlobsIncludeItem string

const (
	FilterBlobsIncludeItemNone     FilterBlobsIncludeItem = "none"
	FilterBlobsIncludeItemVersions FilterBlobsIncludeItem = "versions"
)

// PossibleFilterBlobsIncludeItemValues returns the possible values for the FilterBlobsIncludeItem const type.
func PossibleFilterBlobsIncludeItemValues() []FilterBlobsIncludeItem {
	return []FilterBlobsIncludeItem{
		FilterBlobsIncludeItemNone,
		FilterBlobsIncludeItemVersions,
	}
}

type ImmutabilityPolicyMode string

const (
	ImmutabilityPolicyModeLocked   ImmutabilityPolicyMode = "Locked"
	ImmutabilityPolicyModeMutable  ImmutabilityPolicyMode = "Mutable"
	ImmutabilityPolicyModeUnlocked ImmutabilityPolicyMode = "Unlocked"
)

// PossibleImmutabilityPolicyModeValues returns the possible values for the ImmutabilityPolicyMode const type.
func PossibleImmutabilityPolicyModeValues() []ImmutabilityPolicyMode {
	return []ImmutabilityPolicyMode{
		ImmutabilityPolicyModeLocked,
		ImmutabilityPolicyModeMutable,
		ImmutabilityPolicyModeUnlocked,
	}
}

type ImmutabilityPolicySetting string

const (
	ImmutabilityPolicySettingLocked   ImmutabilityPolicySetting = "Locked"
	ImmutabilityPolicySettingUnlocked ImmutabilityPolicySetting = "Unlocked"
)

// PossibleImmutabilityPolicySettingValues returns the possible values for the ImmutabilityPolicySetting const type.
func PossibleImmutabilityPolicySettingValues() []ImmutabilityPolicySetting {
	return []ImmutabilityPolicySetting{
		ImmutabilityPolicySettingLocked,
		ImmutabilityPolicySettingUnlocked,
	}
}

type LeaseDurationType string

const (
	LeaseDurationTypeFixed    LeaseDurationType = "fixed"
	LeaseDurationTypeInfinite LeaseDurationType = "infinite"
)

// PossibleLeaseDurationTypeValues returns the possible values for the LeaseDurationType const type.
func PossibleLeaseDurationTypeValues() []LeaseDurationType {
	return []LeaseDurationType{
		LeaseDurationTypeFixed,
		LeaseDurationTypeInfinite,
	}
}

type LeaseStateType string

const (
	LeaseStateTypeAvailable LeaseStateType = "available"
	LeaseStateTypeBreaking  LeaseStateType = "breaking"
	LeaseStateTypeBroken    LeaseStateType = "broken"
	LeaseStateTypeExpired   LeaseStateType = "expired"
	LeaseStateTypeLeased    LeaseStateType = "leased"
)

// PossibleLeaseStateTypeValues returns the possible values for the LeaseStateType const type.
func PossibleLeaseStateTypeValues() []LeaseStateType {
	return []LeaseStateType{
		LeaseStateTypeAvailable,
		LeaseStateTypeBreaking,
		LeaseStateTypeBroken,
		LeaseStateTypeExpired,
		LeaseStateTypeLeased,
	}
}

type LeaseStatusType string

const (
	LeaseStatusTypeLocked   LeaseStatusType = "locked"
	LeaseStatusTypeUnlocked LeaseStatusType = "unlocked"
)

// PossibleLeaseStatusTypeValues returns the possible values for the LeaseStatusType const type.
func PossibleLeaseStatusTypeValues() []LeaseStatusType {
	return []LeaseStatusType{
		LeaseStatusTypeLocked,
		LeaseStatusTypeUnlocked,
	}
}

type ListBlobsIncludeItem string

const (
	ListBlobsIncludeItemCopy                ListBlobsIncludeItem = "copy"
	ListBlobsIncludeItemDeleted             ListBlobsIncludeItem = "deleted"
	ListBlobsIncludeItemDeletedwithversions ListBlobsIncludeItem = "deletedwithversions"
	ListBlobsIncludeItemImmutabilitypolicy  ListBlobsIncludeItem = "immutabilitypolicy"
	ListBlobsIncludeItemLegalhold           ListBlobsIncludeItem = "legalhold"
	ListBlobsIncludeItemMetadata            ListBlobsIncludeItem = "metadata"
	ListBlobsIncludeItemPermissions         ListBlobsIncludeItem = "permissions"
	ListBlobsIncludeItemSnapshots           ListBlobsIncludeItem = "snapshots"
	ListBlobsIncludeItemTags                ListBlobsIncludeItem = "tags"
	ListBlobsIncludeItemUncommittedblobs    ListBlobsIncludeItem = "uncommittedblobs"
	ListBlobsIncludeItemVersions            ListBlobsIncludeItem = "versions"
)

// PossibleListBlobsIncludeItemValues returns the possible values for the ListBlobsIncludeItem const type.
func PossibleListBlobsIncludeItemValues() []ListBlobsIncludeItem {
	return []ListBlobsIncludeItem{
		ListBlobsIncludeItemCopy,
		ListBlobsIncludeItemDeleted,
		ListBlobsIncludeItemDeletedwithversions,
		ListBlobsIncludeItemImmutabilitypolicy,
		ListBlobsIncludeItemLegalhold,
		ListBlobsIncludeItemMetadata,
		ListBlobsIncludeItemPermissions,
		ListBlobsIncludeItemSnapshots,
		ListBlobsIncludeItemTags,
		ListBlobsIncludeItemUncommittedblobs,
		ListBlobsIncludeItemVersions,
	}
}

type ListContainersIncludeType string

const (
	ListContainersIncludeTypeDeleted  ListContainersIncludeType = "deleted"
	ListContainersIncludeTypeMetadata ListContainersIncludeType = "metadata"
	ListContainersIncludeTypeSystem   ListContainersIncludeType = "system"
)

// PossibleListContainersIncludeTypeValues returns the possible values for the ListContainersIncludeType const type.
func PossibleListContainersIncludeTypeValues() []ListContainersIncludeType {
	return []ListContainersIncludeType{
		ListContainersIncludeTypeDeleted,
		ListContainersIncludeTypeMetadata,
		ListContainersIncludeTypeSystem,
	}
}

type PremiumPageBlobAccessTier string

const (
	PremiumPageBlobAccessTierP10 PremiumPageBlobAccessTier = "P10"
	PremiumPageBlobAccessTierP15 PremiumPageBlobAccessTier = "P15"
	PremiumPageBlobAccessTierP20 PremiumPageBlobAccessTier = "P20"
	PremiumPageBlobAccessTierP30 PremiumPageBlobAccessTier = "P30"
	PremiumPageBlobAccessTierP4  PremiumPageBlobAccessTier = "P4"
	PremiumPageBlobAccessTierP40 PremiumPageBlobAccessTier = "P40"
	PremiumPageBlobAccessTierP50 PremiumPageBlobAccessTier = "P50"
	PremiumPageBlobAccessTierP6  PremiumPageBlobAccessTier = "P6"
	PremiumPageBlobAccessTierP60 PremiumPageBlobAccessTier = "P60"
	PremiumPageBlobAccessTierP70 PremiumPageBlobAccessTier = "P70"
	PremiumPageBlobAccessTierP80 PremiumPageBlobAccessTier = "P80"
)

// PossiblePremiumPageBlobAccessTierValues returns the possible values for the PremiumPageBlobAccessTier const type.
func PossiblePremiumPageBlobAccessTierValues() []PremiumPageBlobAccessTier {
	return []PremiumPageBlobAccessTier{
		PremiumPageBlobAccessTierP10,
		PremiumPageBlobAccessTierP15,
		PremiumPageBlobAccessTierP20,
		PremiumPageBlobAccessTierP30,
		PremiumPageBlobAccessTierP4,
		PremiumPageBlobAccessTierP40,
		PremiumPageBlobAccessTierP50,
		PremiumPageBlobAccessTierP6,
		PremiumPageBlobAccessTierP60,
		PremiumPageBlobAccessTierP70,
		PremiumPageBlobAccessTierP80,
	}
}

type PublicAccessType string

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

// QueryFormatType - The quick query format type.
type QueryFormatType string

const (
	QueryFormatTypeArrow     QueryFormatType = "arrow"
	QueryFormatTypeDelimited QueryFormatType = "delimited"
	QueryFormatTypeJSON      QueryFormatType = "json"
	QueryFormatTypeParquet   QueryFormatType = "parquet"
)

// PossibleQueryFormatTypeValues returns the possible values for the QueryFormatType const type.
func PossibleQueryFormatTypeValues() []QueryFormatType {
	return []QueryFormatType{
		QueryFormatTypeArrow,
		QueryFormatTypeDelimited,
		QueryFormatTypeJSON,
		QueryFormatTypeParquet,
	}
}

// RehydratePriority - If an object is in rehydrate pending state then this header is returned with priority of rehydrate.
// Valid values are High and Standard.
type RehydratePriority string

const (
	RehydratePriorityHigh     RehydratePriority = "High"
	RehydratePriorityStandard RehydratePriority = "Standard"
)

// PossibleRehydratePriorityValues returns the possible values for the RehydratePriority const type.
func PossibleRehydratePriorityValues() []RehydratePriority {
	return []RehydratePriority{
		RehydratePriorityHigh,
		RehydratePriorityStandard,
	}
}

type SKUName string

const (
	SKUNamePremiumLRS    SKUName = "Premium_LRS"
	SKUNameStandardGRS   SKUName = "Standard_GRS"
	SKUNameStandardLRS   SKUName = "Standard_LRS"
	SKUNameStandardRAGRS SKUName = "Standard_RAGRS"
	SKUNameStandardZRS   SKUName = "Standard_ZRS"
)

// PossibleSKUNameValues returns the possible values for the SKUName const type.
func PossibleSKUNameValues() []SKUName {
	return []SKUName{
		SKUNamePremiumLRS,
		SKUNameStandardGRS,
		SKUNameStandardLRS,
		SKUNameStandardRAGRS,
		SKUNameStandardZRS,
	}
}

type SequenceNumberActionType string

const (
	SequenceNumberActionTypeIncrement SequenceNumberActionType = "increment"
	SequenceNumberActionTypeMax       SequenceNumberActionType = "max"
	SequenceNumberActionTypeUpdate    SequenceNumberActionType = "update"
)

// PossibleSequenceNumberActionTypeValues returns the possible values for the SequenceNumberActionType const type.
func PossibleSequenceNumberActionTypeValues() []SequenceNumberActionType {
	return []SequenceNumberActionType{
		SequenceNumberActionTypeIncrement,
		SequenceNumberActionTypeMax,
		SequenceNumberActionTypeUpdate,
	}
}

// StorageErrorCode - Error codes returned by the service
type StorageErrorCode string

const (
	StorageErrorCodeAccountAlreadyExists                              StorageErrorCode = "AccountAlreadyExists"
	StorageErrorCodeAccountBeingCreated                               StorageErrorCode = "AccountBeingCreated"
	StorageErrorCodeAccountIsDisabled                                 StorageErrorCode = "AccountIsDisabled"
	StorageErrorCodeAppendPositionConditionNotMet                     StorageErrorCode = "AppendPositionConditionNotMet"
	StorageErrorCodeAuthenticationFailed                              StorageErrorCode = "AuthenticationFailed"
	StorageErrorCodeAuthorizationFailure                              StorageErrorCode = "AuthorizationFailure"
	StorageErrorCodeAuthorizationPermissionMismatch                   StorageErrorCode = "AuthorizationPermissionMismatch"
	StorageErrorCodeAuthorizationProtocolMismatch                     StorageErrorCode = "AuthorizationProtocolMismatch"
	StorageErrorCodeAuthorizationResourceTypeMismatch                 StorageErrorCode = "AuthorizationResourceTypeMismatch"
	StorageErrorCodeAuthorizationServiceMismatch                      StorageErrorCode = "AuthorizationServiceMismatch"
	StorageErrorCodeAuthorizationSourceIPMismatch                     StorageErrorCode = "AuthorizationSourceIPMismatch"
	StorageErrorCodeBlobAccessTierNotSupportedForAccountType          StorageErrorCode = "BlobAccessTierNotSupportedForAccountType"
	StorageErrorCodeBlobAlreadyExists                                 StorageErrorCode = "BlobAlreadyExists"
	StorageErrorCodeBlobArchived                                      StorageErrorCode = "BlobArchived"
	StorageErrorCodeBlobBeingRehydrated                               StorageErrorCode = "BlobBeingRehydrated"
	StorageErrorCodeBlobImmutableDueToPolicy                          StorageErrorCode = "BlobImmutableDueToPolicy"
	StorageErrorCodeBlobNotArchived                                   StorageErrorCode = "BlobNotArchived"
	StorageErrorCodeBlobNotFound                                      StorageErrorCode = "BlobNotFound"
	StorageErrorCodeBlobOverwritten                                   StorageErrorCode = "BlobOverwritten"
	StorageErrorCodeBlobTierInadequateForContentLength                StorageErrorCode = "BlobTierInadequateForContentLength"
	StorageErrorCodeBlobUsesCustomerSpecifiedEncryption               StorageErrorCode = "BlobUsesCustomerSpecifiedEncryption"
	StorageErrorCodeBlockCountExceedsLimit                            StorageErrorCode = "BlockCountExceedsLimit"
	StorageErrorCodeBlockListTooLong                                  StorageErrorCode = "BlockListTooLong"
	StorageErrorCodeCannotChangeToLowerTier                           StorageErrorCode = "CannotChangeToLowerTier"
	StorageErrorCodeCannotVerifyCopySource                            StorageErrorCode = "CannotVerifyCopySource"
	StorageErrorCodeConditionHeadersNotSupported                      StorageErrorCode = "ConditionHeadersNotSupported"
	StorageErrorCodeConditionNotMet                                   StorageErrorCode = "ConditionNotMet"
	StorageErrorCodeContainerAlreadyExists                            StorageErrorCode = "ContainerAlreadyExists"
	StorageErrorCodeContainerBeingDeleted                             StorageErrorCode = "ContainerBeingDeleted"
	StorageErrorCodeContainerDisabled                                 StorageErrorCode = "ContainerDisabled"
	StorageErrorCodeContainerNotFound                                 StorageErrorCode = "ContainerNotFound"
	StorageErrorCodeContentLengthLargerThanTierLimit                  StorageErrorCode = "ContentLengthLargerThanTierLimit"
	StorageErrorCodeCopyAcrossAccountsNotSupported                    StorageErrorCode = "CopyAcrossAccountsNotSupported"
	StorageErrorCodeCopyIDMismatch                                    StorageErrorCode = "CopyIdMismatch"
	StorageErrorCodeEmptyMetadataKey                                  StorageErrorCode = "EmptyMetadataKey"
	StorageErrorCodeFeatureVersionMismatch                            StorageErrorCode = "FeatureVersionMismatch"
	StorageErrorCodeIncrementalCopyBlobMismatch                       StorageErrorCode = "IncrementalCopyBlobMismatch"
	StorageErrorCodeIncrementalCopyOfEarlierVersionSnapshotNotAllowed StorageErrorCode = "IncrementalCopyOfEarlierVersionSnapshotNotAllowed"
	StorageErrorCodeIncrementalCopySourceMustBeSnapshot               StorageErrorCode = "IncrementalCopySourceMustBeSnapshot"
	StorageErrorCodeInfiniteLeaseDurationRequired                     StorageErrorCode = "InfiniteLeaseDurationRequired"
	StorageErrorCodeInsufficientAccountPermissions                    StorageErrorCode = "InsufficientAccountPermissions"
	StorageErrorCodeInternalError                                     StorageErrorCode = "InternalError"
	StorageErrorCodeInvalidAuthenticationInfo                         StorageErrorCode = "InvalidAuthenticationInfo"
	StorageErrorCodeInvalidBlobOrBlock                                StorageErrorCode = "InvalidBlobOrBlock"
	StorageErrorCodeInvalidBlobTier                                   StorageErrorCode = "InvalidBlobTier"
	StorageErrorCodeInvalidBlobType                                   StorageErrorCode = "InvalidBlobType"
	StorageErrorCodeInvalidBlockID                                    StorageErrorCode = "InvalidBlockId"
	StorageErrorCodeInvalidBlockList                                  StorageErrorCode = "InvalidBlockList"
	StorageErrorCodeInvalidHTTPVerb                                   StorageErrorCode = "InvalidHttpVerb"
	StorageErrorCodeInvalidHeaderValue                                StorageErrorCode = "InvalidHeaderValue"
	StorageErrorCodeInvalidInput                                      StorageErrorCode = "InvalidInput"
	StorageErrorCodeInvalidMD5                                        StorageErrorCode = "InvalidMd5"
	StorageErrorCodeInvalidMetadata                                   StorageErrorCode = "InvalidMetadata"
	StorageErrorCodeInvalidOperation                                  StorageErrorCode = "InvalidOperation"
	StorageErrorCodeInvalidPageRange                                  StorageErrorCode = "InvalidPageRange"
	StorageErrorCodeInvalidQueryParameterValue                        StorageErrorCode = "InvalidQueryParameterValue"
	StorageErrorCodeInvalidRange                                      StorageErrorCode = "InvalidRange"
	StorageErrorCodeInvalidResourceName                               StorageErrorCode = "InvalidResourceName"
	StorageErrorCodeInvalidSourceBlobType                             StorageErrorCode = "InvalidSourceBlobType"
	StorageErrorCodeInvalidSourceBlobURL                              StorageErrorCode = "InvalidSourceBlobUrl"
	StorageErrorCodeInvalidURI                                        StorageErrorCode = "InvalidUri"
	StorageErrorCodeInvalidVersionForPageBlobOperation                StorageErrorCode = "InvalidVersionForPageBlobOperation"
	StorageErrorCodeInvalidXMLDocument                                StorageErrorCode = "InvalidXmlDocument"
	StorageErrorCodeInvalidXMLNodeValue                               StorageErrorCode = "InvalidXmlNodeValue"
	StorageErrorCodeLeaseAlreadyBroken                                StorageErrorCode = "LeaseAlreadyBroken"
	StorageErrorCodeLeaseAlreadyPresent                               StorageErrorCode = "LeaseAlreadyPresent"
	StorageErrorCodeLeaseIDMismatchWithBlobOperation                  StorageErrorCode = "LeaseIdMismatchWithBlobOperation"
	StorageErrorCodeLeaseIDMismatchWithContainerOperation             StorageErrorCode = "LeaseIdMismatchWithContainerOperation"
	StorageErrorCodeLeaseIDMismatchWithLeaseOperation                 StorageErrorCode = "LeaseIdMismatchWithLeaseOperation"
	StorageErrorCodeLeaseIDMissing                                    StorageErrorCode = "LeaseIdMissing"
	StorageErrorCodeLeaseIsBreakingAndCannotBeAcquired                StorageErrorCode = "LeaseIsBreakingAndCannotBeAcquired"
	StorageErrorCodeLeaseIsBreakingAndCannotBeChanged                 StorageErrorCode = "LeaseIsBreakingAndCannotBeChanged"
	StorageErrorCodeLeaseIsBrokenAndCannotBeRenewed                   StorageErrorCode = "LeaseIsBrokenAndCannotBeRenewed"
	StorageErrorCodeLeaseLost                                         StorageErrorCode = "LeaseLost"
	StorageErrorCodeLeaseNotPresentWithBlobOperation                  StorageErrorCode = "LeaseNotPresentWithBlobOperation"
	StorageErrorCodeLeaseNotPresentWithContainerOperation             StorageErrorCode = "LeaseNotPresentWithContainerOperation"
	StorageErrorCodeLeaseNotPresentWithLeaseOperation                 StorageErrorCode = "LeaseNotPresentWithLeaseOperation"
	StorageErrorCodeMD5Mismatch                                       StorageErrorCode = "Md5Mismatch"
	StorageErrorCodeMaxBlobSizeConditionNotMet                        StorageErrorCode = "MaxBlobSizeConditionNotMet"
	StorageErrorCodeMetadataTooLarge                                  StorageErrorCode = "MetadataTooLarge"
	StorageErrorCodeMissingContentLengthHeader                        StorageErrorCode = "MissingContentLengthHeader"
	StorageErrorCodeMissingRequiredHeader                             StorageErrorCode = "MissingRequiredHeader"
	StorageErrorCodeMissingRequiredQueryParameter                     StorageErrorCode = "MissingRequiredQueryParameter"
	StorageErrorCodeMissingRequiredXMLNode                            StorageErrorCode = "MissingRequiredXmlNode"
	StorageErrorCodeMultipleConditionHeadersNotSupported              StorageErrorCode = "MultipleConditionHeadersNotSupported"
	StorageErrorCodeNoAuthenticationInformation                       StorageErrorCode = "NoAuthenticationInformation"
	StorageErrorCodeNoPendingCopyOperation                            StorageErrorCode = "NoPendingCopyOperation"
	StorageErrorCodeOperationNotAllowedOnIncrementalCopyBlob          StorageErrorCode = "OperationNotAllowedOnIncrementalCopyBlob"
	StorageErrorCodeOperationTimedOut                                 StorageErrorCode = "OperationTimedOut"
	StorageErrorCodeOutOfRangeInput                                   StorageErrorCode = "OutOfRangeInput"
	StorageErrorCodeOutOfRangeQueryParameterValue                     StorageErrorCode = "OutOfRangeQueryParameterValue"
	StorageErrorCodePendingCopyOperation                              StorageErrorCode = "PendingCopyOperation"
	StorageErrorCodePreviousSnapshotCannotBeNewer                     StorageErrorCode = "PreviousSnapshotCannotBeNewer"
	StorageErrorCodePreviousSnapshotNotFound                          StorageErrorCode = "PreviousSnapshotNotFound"
	StorageErrorCodePreviousSnapshotOperationNotSupported             StorageErrorCode = "PreviousSnapshotOperationNotSupported"
	StorageErrorCodeRequestBodyTooLarge                               StorageErrorCode = "RequestBodyTooLarge"
	StorageErrorCodeRequestURLFailedToParse                           StorageErrorCode = "RequestUrlFailedToParse"
	StorageErrorCodeResourceAlreadyExists                             StorageErrorCode = "ResourceAlreadyExists"
	StorageErrorCodeResourceNotFound                                  StorageErrorCode = "ResourceNotFound"
	StorageErrorCodeResourceTypeMismatch                              StorageErrorCode = "ResourceTypeMismatch"
	StorageErrorCodeSequenceNumberConditionNotMet                     StorageErrorCode = "SequenceNumberConditionNotMet"
	StorageErrorCodeSequenceNumberIncrementTooLarge                   StorageErrorCode = "SequenceNumberIncrementTooLarge"
	StorageErrorCodeServerBusy                                        StorageErrorCode = "ServerBusy"
	StorageErrorCodeSnapshotCountExceeded                             StorageErrorCode = "SnapshotCountExceeded"
	StorageErrorCodeSnapshotOperationRateExceeded                     StorageErrorCode = "SnapshotOperationRateExceeded"
	StorageErrorCodeSnapshotsPresent                                  StorageErrorCode = "SnapshotsPresent"
	StorageErrorCodeSourceConditionNotMet                             StorageErrorCode = "SourceConditionNotMet"
	StorageErrorCodeSystemInUse                                       StorageErrorCode = "SystemInUse"
	StorageErrorCodeTargetConditionNotMet                             StorageErrorCode = "TargetConditionNotMet"
	StorageErrorCodeUnauthorizedBlobOverwrite                         StorageErrorCode = "UnauthorizedBlobOverwrite"
	StorageErrorCodeUnsupportedHTTPVerb                               StorageErrorCode = "UnsupportedHttpVerb"
	StorageErrorCodeUnsupportedHeader                                 StorageErrorCode = "UnsupportedHeader"
	StorageErrorCodeUnsupportedQueryParameter                         StorageErrorCode = "UnsupportedQueryParameter"
	StorageErrorCodeUnsupportedXMLNode                                StorageErrorCode = "UnsupportedXmlNode"
)

// PossibleStorageErrorCodeValues returns the possible values for the StorageErrorCode const type.
func PossibleStorageErrorCodeValues() []StorageErrorCode {
	return []StorageErrorCode{
		StorageErrorCodeAccountAlreadyExists,
		StorageErrorCodeAccountBeingCreated,
		StorageErrorCodeAccountIsDisabled,
		StorageErrorCodeAppendPositionConditionNotMet,
		StorageErrorCodeAuthenticationFailed,
		StorageErrorCodeAuthorizationFailure,
		StorageErrorCodeAuthorizationPermissionMismatch,
		StorageErrorCodeAuthorizationProtocolMismatch,
		StorageErrorCodeAuthorizationResourceTypeMismatch,
		StorageErrorCodeAuthorizationServiceMismatch,
		StorageErrorCodeAuthorizationSourceIPMismatch,
		StorageErrorCodeBlobAccessTierNotSupportedForAccountType,
		StorageErrorCodeBlobAlreadyExists,
		StorageErrorCodeBlobArchived,
		StorageErrorCodeBlobBeingRehydrated,
		StorageErrorCodeBlobImmutableDueToPolicy,
		StorageErrorCodeBlobNotArchived,
		StorageErrorCodeBlobNotFound,
		StorageErrorCodeBlobOverwritten,
		StorageErrorCodeBlobTierInadequateForContentLength,
		StorageErrorCodeBlobUsesCustomerSpecifiedEncryption,
		StorageErrorCodeBlockCountExceedsLimit,
		StorageErrorCodeBlockListTooLong,
		StorageErrorCodeCannotChangeToLowerTier,
		StorageErrorCodeCannotVerifyCopySource,
		StorageErrorCodeConditionHeadersNotSupported,
		StorageErrorCodeConditionNotMet,
		StorageErrorCodeContainerAlreadyExists,
		StorageErrorCodeContainerBeingDeleted,
		StorageErrorCodeContainerDisabled,
		StorageErrorCodeContainerNotFound,
		StorageErrorCodeContentLengthLargerThanTierLimit,
		StorageErrorCodeCopyAcrossAccountsNotSupported,
		StorageErrorCodeCopyIDMismatch,
		StorageErrorCodeEmptyMetadataKey,
		StorageErrorCodeFeatureVersionMismatch,
		StorageErrorCodeIncrementalCopyBlobMismatch,
		StorageErrorCodeIncrementalCopyOfEarlierVersionSnapshotNotAllowed,
		StorageErrorCodeIncrementalCopySourceMustBeSnapshot,
		StorageErrorCodeInfiniteLeaseDurationRequired,
		StorageErrorCodeInsufficientAccountPermissions,
		StorageErrorCodeInternalError,
		StorageErrorCodeInvalidAuthenticationInfo,
		StorageErrorCodeInvalidBlobOrBlock,
		StorageErrorCodeInvalidBlobTier,
		StorageErrorCodeInvalidBlobType,
		StorageErrorCodeInvalidBlockID,
		StorageErrorCodeInvalidBlockList,
		StorageErrorCodeInvalidHTTPVerb,
		StorageErrorCodeInvalidHeaderValue,
		StorageErrorCodeInvalidInput,
		StorageErrorCodeInvalidMD5,
		StorageErrorCodeInvalidMetadata,
		StorageErrorCodeInvalidOperation,
		StorageErrorCodeInvalidPageRange,
		StorageErrorCodeInvalidQueryParameterValue,
		StorageErrorCodeInvalidRange,
		StorageErrorCodeInvalidResourceName,
		StorageErrorCodeInvalidSourceBlobType,
		StorageErrorCodeInvalidSourceBlobURL,
		StorageErrorCodeInvalidURI,
		StorageErrorCodeInvalidVersionForPageBlobOperation,
		StorageErrorCodeInvalidXMLDocument,
		StorageErrorCodeInvalidXMLNodeValue,
		StorageErrorCodeLeaseAlreadyBroken,
		StorageErrorCodeLeaseAlreadyPresent,
		StorageErrorCodeLeaseIDMismatchWithBlobOperation,
		StorageErrorCodeLeaseIDMismatchWithContainerOperation,
		StorageErrorCodeLeaseIDMismatchWithLeaseOperation,
		StorageErrorCodeLeaseIDMissing,
		StorageErrorCodeLeaseIsBreakingAndCannotBeAcquired,
		StorageErrorCodeLeaseIsBreakingAndCannotBeChanged,
		StorageErrorCodeLeaseIsBrokenAndCannotBeRenewed,
		StorageErrorCodeLeaseLost,
		StorageErrorCodeLeaseNotPresentWithBlobOperation,
		StorageErrorCodeLeaseNotPresentWithContainerOperation,
		StorageErrorCodeLeaseNotPresentWithLeaseOperation,
		StorageErrorCodeMD5Mismatch,
		StorageErrorCodeMaxBlobSizeConditionNotMet,
		StorageErrorCodeMetadataTooLarge,
		StorageErrorCodeMissingContentLengthHeader,
		StorageErrorCodeMissingRequiredHeader,
		StorageErrorCodeMissingRequiredQueryParameter,
		StorageErrorCodeMissingRequiredXMLNode,
		StorageErrorCodeMultipleConditionHeadersNotSupported,
		StorageErrorCodeNoAuthenticationInformation,
		StorageErrorCodeNoPendingCopyOperation,
		StorageErrorCodeOperationNotAllowedOnIncrementalCopyBlob,
		StorageErrorCodeOperationTimedOut,
		StorageErrorCodeOutOfRangeInput,
		StorageErrorCodeOutOfRangeQueryParameterValue,
		StorageErrorCodePendingCopyOperation,
		StorageErrorCodePreviousSnapshotCannotBeNewer,
		StorageErrorCodePreviousSnapshotNotFound,
		StorageErrorCodePreviousSnapshotOperationNotSupported,
		StorageErrorCodeRequestBodyTooLarge,
		StorageErrorCodeRequestURLFailedToParse,
		StorageErrorCodeResourceAlreadyExists,
		StorageErrorCodeResourceNotFound,
		StorageErrorCodeResourceTypeMismatch,
		StorageErrorCodeSequenceNumberConditionNotMet,
		StorageErrorCodeSequenceNumberIncrementTooLarge,
		StorageErrorCodeServerBusy,
		StorageErrorCodeSnapshotCountExceeded,
		StorageErrorCodeSnapshotOperationRateExceeded,
		StorageErrorCodeSnapshotsPresent,
		StorageErrorCodeSourceConditionNotMet,
		StorageErrorCodeSystemInUse,
		StorageErrorCodeTargetConditionNotMet,
		StorageErrorCodeUnauthorizedBlobOverwrite,
		StorageErrorCodeUnsupportedHTTPVerb,
		StorageErrorCodeUnsupportedHeader,
		StorageErrorCodeUnsupportedQueryParameter,
		StorageErrorCodeUnsupportedXMLNode,
	}
}
