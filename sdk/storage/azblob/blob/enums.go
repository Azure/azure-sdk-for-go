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
