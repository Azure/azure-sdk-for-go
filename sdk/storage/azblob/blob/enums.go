//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

const (
	DeleteSnapshotsOptionTypeInclude = generated.DeleteSnapshotsOptionTypeInclude
	DeleteSnapshotsOptionTypeOnly    = generated.DeleteSnapshotsOptionTypeOnly
)

// PossibleDeleteSnapshotsOptionTypeValues returns the possible values for the DeleteSnapshotsOptionType const type.
func PossibleDeleteSnapshotsOptionTypeValues() []DeleteSnapshotsOptionType {
	return generated.PossibleDeleteSnapshotsOptionTypeValues()
}

const (
	AccessTierArchive = generated.AccessTierArchive
	AccessTierCool    = generated.AccessTierCool
	AccessTierHot     = generated.AccessTierHot
	AccessTierP10     = generated.AccessTierP10
	AccessTierP15     = generated.AccessTierP15
	AccessTierP20     = generated.AccessTierP20
	AccessTierP30     = generated.AccessTierP30
	AccessTierP4      = generated.AccessTierP4
	AccessTierP40     = generated.AccessTierP40
	AccessTierP50     = generated.AccessTierP50
	AccessTierP6      = generated.AccessTierP6
	AccessTierP60     = generated.AccessTierP60
	AccessTierP70     = generated.AccessTierP70
	AccessTierP80     = generated.AccessTierP80
	AccessTierPremium = generated.AccessTierPremium
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
		AccessTierPremium,
	}
}

const (
	RehydratePriorityHigh     = generated.RehydratePriorityHigh
	RehydratePriorityStandard = generated.RehydratePriorityStandard
)

// PossibleRehydratePriorityValues returns the possible values for the RehydratePriority const type.
func PossibleRehydratePriorityValues() []RehydratePriority {
	return []RehydratePriority{
		RehydratePriorityHigh,
		RehydratePriorityStandard,
	}
}

const (
	ImmutabilityPolicyModeMutable  = generated.BlobImmutabilityPolicyModeMutable
	ImmutabilityPolicyModeUnlocked = generated.BlobImmutabilityPolicyModeUnlocked
	ImmutabilityPolicyModeLocked   = generated.BlobImmutabilityPolicyModeLocked
)

// PossibleBlobImmutabilityPolicyModeValues returns the possible values for the BlobImmutabilityPolicyMode const type.
func PossibleBlobImmutabilityPolicyModeValues() []ImmutabilityPolicyMode {
	return []ImmutabilityPolicyMode{
		ImmutabilityPolicyModeMutable,
		ImmutabilityPolicyModeUnlocked,
		ImmutabilityPolicyModeLocked,
	}
}

const (
	CopyStatusTypePending = generated.CopyStatusTypePending
	CopyStatusTypeSuccess = generated.CopyStatusTypeSuccess
	CopyStatusTypeAborted = generated.CopyStatusTypeAborted
	CopyStatusTypeFailed  = generated.CopyStatusTypeFailed
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

const (
	EncryptionAlgorithmTypeNone   = generated.EncryptionAlgorithmTypeNone
	EncryptionAlgorithmTypeAES256 = generated.EncryptionAlgorithmTypeAES256
)

// PossibleEncryptionAlgorithmTypeValues returns the possible values for the EncryptionAlgorithmType const type.
func PossibleEncryptionAlgorithmTypeValues() []EncryptionAlgorithmType {
	return []EncryptionAlgorithmType{
		EncryptionAlgorithmTypeNone,
		EncryptionAlgorithmTypeAES256,
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

// DeleteType defines values for BlobDeleteType
type DeleteType = generated.BlobDeleteType

const (
	DeleteTypeNone      = generated.BlobDeleteTypeNone
	DeleteTypePermanent = generated.BlobDeleteTypePermanent
)

// PossibleBlobDeleteTypeValues returns the possible values for the BlobDeleteType const type.
func PossibleBlobDeleteTypeValues() []DeleteType {
	return []DeleteType{
		DeleteTypeNone,
		DeleteTypePermanent,
	}
}

// ExpiryOptions defines values for BlobExpiryOptions
type ExpiryOptions = generated.BlobExpiryOptions

const (
	ExpiryOptionsAbsolute           = generated.BlobExpiryOptionsAbsolute
	ExpiryOptionsNeverExpire        = generated.BlobExpiryOptionsNeverExpire
	ExpiryOptionsRelativeToCreation = generated.BlobExpiryOptionsRelativeToCreation
	ExpiryOptionsRelativeToNow      = generated.BlobExpiryOptionsRelativeToNow
)

// PossibleBlobExpiryOptionsValues returns the possible values for the BlobExpiryOptions const type.
func PossibleBlobExpiryOptionsValues() []ExpiryOptions {
	return []ExpiryOptions{
		ExpiryOptionsAbsolute,
		ExpiryOptionsNeverExpire,
		ExpiryOptionsRelativeToCreation,
		ExpiryOptionsRelativeToNow,
	}
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
