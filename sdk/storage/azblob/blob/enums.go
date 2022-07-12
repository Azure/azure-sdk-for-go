//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

const (
	DeleteSnapshotsOptionTypeInclude DeleteSnapshotsOptionType = generated.DeleteSnapshotsOptionTypeInclude
	DeleteSnapshotsOptionTypeOnly    DeleteSnapshotsOptionType = generated.DeleteSnapshotsOptionTypeOnly
)

// PossibleDeleteSnapshotsOptionTypeValues returns the possible values for the DeleteSnapshotsOptionType const type.
func PossibleDeleteSnapshotsOptionTypeValues() []DeleteSnapshotsOptionType {
	return generated.PossibleDeleteSnapshotsOptionTypeValues()
}

const (
	AccessTierArchive AccessTier = generated.AccessTierArchive
	AccessTierCool    AccessTier = generated.AccessTierCool
	AccessTierHot     AccessTier = generated.AccessTierHot
	AccessTierP10     AccessTier = generated.AccessTierP10
	AccessTierP15     AccessTier = generated.AccessTierP15
	AccessTierP20     AccessTier = generated.AccessTierP20
	AccessTierP30     AccessTier = generated.AccessTierP30
	AccessTierP4      AccessTier = generated.AccessTierP4
	AccessTierP40     AccessTier = generated.AccessTierP40
	AccessTierP50     AccessTier = generated.AccessTierP50
	AccessTierP6      AccessTier = generated.AccessTierP6
	AccessTierP60     AccessTier = generated.AccessTierP60
	AccessTierP70     AccessTier = generated.AccessTierP70
	AccessTierP80     AccessTier = generated.AccessTierP80
	AccessTierPremium AccessTier = generated.AccessTierPremium
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
	RehydratePriorityHigh     RehydratePriority = generated.RehydratePriorityHigh
	RehydratePriorityStandard RehydratePriority = generated.RehydratePriorityStandard
)

// PossibleRehydratePriorityValues returns the possible values for the RehydratePriority const type.
func PossibleRehydratePriorityValues() []RehydratePriority {
	return []RehydratePriority{
		RehydratePriorityHigh,
		RehydratePriorityStandard,
	}
}

const (
	ImmutabilityPolicyModeMutable  ImmutabilityPolicyMode = generated.BlobImmutabilityPolicyModeMutable
	ImmutabilityPolicyModeUnlocked ImmutabilityPolicyMode = generated.BlobImmutabilityPolicyModeUnlocked
	ImmutabilityPolicyModeLocked   ImmutabilityPolicyMode = generated.BlobImmutabilityPolicyModeLocked
)

// PossibleBlobImmutabilityPolicyModeValues returns the possible values for the BlobImmutabilityPolicyMode const type.
func PossibleBlobImmutabilityPolicyModeValues() []ImmutabilityPolicyMode {
	return []ImmutabilityPolicyMode{
		ImmutabilityPolicyModeMutable,
		ImmutabilityPolicyModeUnlocked,
		ImmutabilityPolicyModeLocked,
	}
}

type CopyStatusType = generated.CopyStatusType

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

type EncryptionAlgorithmType = generated.EncryptionAlgorithmType

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
