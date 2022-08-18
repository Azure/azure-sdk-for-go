//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package blob

const (
	moduleName    = "azblob"
	moduleVersion = "v0.4.1"
)

// AccessTier enum
type AccessTier string

const (
	AccessTierArchive AccessTier = "Archive"
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
)

// PossibleAccessTierValues returns the possible values for the AccessTier const type.
func PossibleAccessTierValues() []AccessTier {
	return []AccessTier{
		AccessTierArchive,
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
	}
}

// ToPtr returns a *AccessTier pointing to the current value.
func (c AccessTier) ToPtr() *AccessTier {
	return &c
}

// AccountKind enum
type AccountKind string

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

// ToPtr returns a *AccountKind pointing to the current value.
func (c AccountKind) ToPtr() *AccountKind {
	return &c
}

// ArchiveStatus enum
type ArchiveStatus string

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

// ToPtr returns a *ArchiveStatus pointing to the current value.
func (c ArchiveStatus) ToPtr() *ArchiveStatus {
	return &c
}

// BlobExpiryOptions enum
type BlobExpiryOptions string

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

// ToPtr returns a *BlobExpiryOptions pointing to the current value.
func (c BlobExpiryOptions) ToPtr() *BlobExpiryOptions {
	return &c
}

// BlobGeoReplicationStatus - The status of the secondary location
type BlobGeoReplicationStatus string

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

// ToPtr returns a *BlobGeoReplicationStatus pointing to the current value.
func (c BlobGeoReplicationStatus) ToPtr() *BlobGeoReplicationStatus {
	return &c
}

// BlobImmutabilityPolicyMode enum
type BlobImmutabilityPolicyMode string

const (
	BlobImmutabilityPolicyModeMutable  BlobImmutabilityPolicyMode = "Mutable"
	BlobImmutabilityPolicyModeUnlocked BlobImmutabilityPolicyMode = "Unlocked"
	BlobImmutabilityPolicyModeLocked   BlobImmutabilityPolicyMode = "Locked"
)

// PossibleBlobImmutabilityPolicyModeValues returns the possible values for the BlobImmutabilityPolicyMode const type.
func PossibleBlobImmutabilityPolicyModeValues() []BlobImmutabilityPolicyMode {
	return []BlobImmutabilityPolicyMode{
		BlobImmutabilityPolicyModeMutable,
		BlobImmutabilityPolicyModeUnlocked,
		BlobImmutabilityPolicyModeLocked,
	}
}

// ToPtr returns a *BlobImmutabilityPolicyMode pointing to the current value.
func (c BlobImmutabilityPolicyMode) ToPtr() *BlobImmutabilityPolicyMode {
	return &c
}

// BlobType enum
type BlobType string

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

// ToPtr returns a *BlobType pointing to the current value.
func (c BlobType) ToPtr() *BlobType {
	return &c
}

// BlockListType enum
type BlockListType string

const (
	BlockListTypeCommitted   BlockListType = "committed"
	BlockListTypeUncommitted BlockListType = "uncommitted"
	BlockListTypeAll         BlockListType = "all"
)

// PossibleBlockListTypeValues returns the possible values for the BlockListType const type.
func PossibleBlockListTypeValues() []BlockListType {
	return []BlockListType{
		BlockListTypeCommitted,
		BlockListTypeUncommitted,
		BlockListTypeAll,
	}
}

// ToPtr returns a *BlockListType pointing to the current value.
func (c BlockListType) ToPtr() *BlockListType {
	return &c
}

// CopyStatusType enum
type CopyStatusType string

const (
	CopyStatusTypePending CopyStatusType = "pending"
	CopyStatusTypeSuccess CopyStatusType = "success"
	CopyStatusTypeAborted CopyStatusType = "aborted"
	CopyStatusTypeFailed  CopyStatusType = "failed"
)

// PossibleCopyStatusTypeValues returns the possible values for the CopyStatusType const type.
func PossibleCopyStatusTypeValues() []CopyStatusType {
	return []CopyStatusType{
		CopyStatusTypePending,
		CopyStatusTypeSuccess,
		CopyStatusTypeAborted,
		CopyStatusTypeFailed,
	}
}

// ToPtr returns a *CopyStatusType pointing to the current value.
func (c CopyStatusType) ToPtr() *CopyStatusType {
	return &c
}

// DeleteSnapshotsOptionType enum
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

// ToPtr returns a *DeleteSnapshotsOptionType pointing to the current value.
func (c DeleteSnapshotsOptionType) ToPtr() *DeleteSnapshotsOptionType {
	return &c
}

// EncryptionAlgorithmType enum
type EncryptionAlgorithmType string

const (
	EncryptionAlgorithmTypeNone   EncryptionAlgorithmType = "None"
	EncryptionAlgorithmTypeAES256 EncryptionAlgorithmType = "AES256"
)

// PossibleEncryptionAlgorithmTypeValues returns the possible values for the EncryptionAlgorithmType const type.
func PossibleEncryptionAlgorithmTypeValues() []EncryptionAlgorithmType {
	return []EncryptionAlgorithmType{
		EncryptionAlgorithmTypeNone,
		EncryptionAlgorithmTypeAES256,
	}
}

// ToPtr returns a *EncryptionAlgorithmType pointing to the current value.
func (c EncryptionAlgorithmType) ToPtr() *EncryptionAlgorithmType {
	return &c
}

// LeaseDurationType enum
type LeaseDurationType string

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

// ToPtr returns a *LeaseDurationType pointing to the current value.
func (c LeaseDurationType) ToPtr() *LeaseDurationType {
	return &c
}

// LeaseStateType enum
type LeaseStateType string

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

// ToPtr returns a *LeaseStateType pointing to the current value.
func (c LeaseStateType) ToPtr() *LeaseStateType {
	return &c
}

// LeaseStatusType enum
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

// ToPtr returns a *LeaseStatusType pointing to the current value.
func (c LeaseStatusType) ToPtr() *LeaseStatusType {
	return &c
}

// ListBlobsIncludeItem enum
type ListBlobsIncludeItem string

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

// ToPtr returns a *ListBlobsIncludeItem pointing to the current value.
func (c ListBlobsIncludeItem) ToPtr() *ListBlobsIncludeItem {
	return &c
}

// ListContainersIncludeType enum
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

// ToPtr returns a *ListContainersIncludeType pointing to the current value.
func (c ListContainersIncludeType) ToPtr() *ListContainersIncludeType {
	return &c
}

// PremiumPageBlobAccessTier enum
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

// ToPtr returns a *PremiumPageBlobAccessTier pointing to the current value.
func (c PremiumPageBlobAccessTier) ToPtr() *PremiumPageBlobAccessTier {
	return &c
}

// PublicAccessType enum
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

// ToPtr returns a *PublicAccessType pointing to the current value.
func (c PublicAccessType) ToPtr() *PublicAccessType {
	return &c
}

// QueryFormatType - The quick query format type.
type QueryFormatType string

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

// ToPtr returns a *QueryFormatType pointing to the current value.
func (c QueryFormatType) ToPtr() *QueryFormatType {
	return &c
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

// ToPtr returns a *RehydratePriority pointing to the current value.
func (c RehydratePriority) ToPtr() *RehydratePriority {
	return &c
}

// SKUName enum
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

// ToPtr returns a *SKUName pointing to the current value.
func (c SKUName) ToPtr() *SKUName {
	return &c
}

// SequenceNumberActionType enum
type SequenceNumberActionType string

const (
	SequenceNumberActionTypeMax       SequenceNumberActionType = "max"
	SequenceNumberActionTypeUpdate    SequenceNumberActionType = "update"
	SequenceNumberActionTypeIncrement SequenceNumberActionType = "increment"
)

// PossibleSequenceNumberActionTypeValues returns the possible values for the SequenceNumberActionType const type.
func PossibleSequenceNumberActionTypeValues() []SequenceNumberActionType {
	return []SequenceNumberActionType{
		SequenceNumberActionTypeMax,
		SequenceNumberActionTypeUpdate,
		SequenceNumberActionTypeIncrement,
	}
}

// ToPtr returns a *SequenceNumberActionType pointing to the current value.
func (c SequenceNumberActionType) ToPtr() *SequenceNumberActionType {
	return &c
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
	StorageErrorCodeIncrementalCopyOfEralierVersionSnapshotNotAllowed StorageErrorCode = "IncrementalCopyOfEralierVersionSnapshotNotAllowed"
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
		StorageErrorCodeIncrementalCopyOfEralierVersionSnapshotNotAllowed,
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

// ToPtr returns a *StorageErrorCode pointing to the current value.
func (c StorageErrorCode) ToPtr() *StorageErrorCode {
	return &c
}

// BlobDeleteType enum
type BlobDeleteType string

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

// ToPtr returns a *BlobDeleteType pointing to the current value.
func (c BlobDeleteType) ToPtr() *BlobDeleteType {
	return &c
}
