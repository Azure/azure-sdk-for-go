//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/path"
)

// EncryptionAlgorithmType defines values for EncryptionAlgorithmType.
type EncryptionAlgorithmType = path.EncryptionAlgorithmType

const (
	EncryptionAlgorithmTypeNone   EncryptionAlgorithmType = path.EncryptionAlgorithmTypeNone
	EncryptionAlgorithmTypeAES256 EncryptionAlgorithmType = path.EncryptionAlgorithmTypeAES256
)

// response models:

// ImmutabilityPolicyMode Specifies the immutability policy mode to set on the file.
type ImmutabilityPolicyMode = path.ImmutabilityPolicyMode

const (
	ImmutabilityPolicyModeMutable  ImmutabilityPolicyMode = path.ImmutabilityPolicyModeMutable
	ImmutabilityPolicyModeUnlocked ImmutabilityPolicyMode = path.ImmutabilityPolicyModeUnlocked
	ImmutabilityPolicyModeLocked   ImmutabilityPolicyMode = path.ImmutabilityPolicyModeLocked
)

// CopyStatusType defines values for CopyStatusType
type CopyStatusType = path.CopyStatusType

const (
	CopyStatusTypePending CopyStatusType = path.CopyStatusTypePending
	CopyStatusTypeSuccess CopyStatusType = path.CopyStatusTypeSuccess
	CopyStatusTypeAborted CopyStatusType = path.CopyStatusTypeAborted
	CopyStatusTypeFailed  CopyStatusType = path.CopyStatusTypeFailed
)

// TransferValidationType abstracts the various mechanisms used to verify a transfer.
type TransferValidationType = exported.TransferValidationType

// TransferValidationTypeCRC64 is a TransferValidationType used to provide a precomputed crc64.
type TransferValidationTypeCRC64 = exported.TransferValidationTypeCRC64

// TransferValidationTypeComputeCRC64 is a TransferValidationType that indicates a CRC64 should be computed during transfer.
func TransferValidationTypeComputeCRC64() TransferValidationType {
	return exported.TransferValidationTypeComputeCRC64()
}

// StatusType defines values for StatusType
type StatusType = azdatalake.StatusType

const (
	StatusTypeLocked   StatusType = azdatalake.StatusTypeLocked
	StatusTypeUnlocked StatusType = azdatalake.StatusTypeUnlocked
)

// PossibleStatusTypeValues returns the possible values for the StatusType const type.
func PossibleStatusTypeValues() []StatusType {
	return azdatalake.PossibleStatusTypeValues()
}

// DurationType defines values for DurationType
type DurationType = azdatalake.DurationType

const (
	DurationTypeInfinite DurationType = azdatalake.DurationTypeInfinite
	DurationTypeFixed    DurationType = azdatalake.DurationTypeFixed
)

// PossibleDurationTypeValues returns the possible values for the DurationType const type.
func PossibleDurationTypeValues() []DurationType {
	return azdatalake.PossibleDurationTypeValues()
}

// StateType defines values for StateType
type StateType = azdatalake.StateType

const (
	StateTypeAvailable StateType = azdatalake.StateTypeAvailable
	StateTypeLeased    StateType = azdatalake.StateTypeLeased
	StateTypeExpired   StateType = azdatalake.StateTypeExpired
	StateTypeBreaking  StateType = azdatalake.StateTypeBreaking
	StateTypeBroken    StateType = azdatalake.StateTypeBroken
)
