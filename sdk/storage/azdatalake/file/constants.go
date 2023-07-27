//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/path"
)

type EncryptionAlgorithmType = path.EncryptionAlgorithmType

const (
	EncryptionAlgorithmTypeNone   EncryptionAlgorithmType = path.EncryptionAlgorithmTypeNone
	EncryptionAlgorithmTypeAES256 EncryptionAlgorithmType = path.EncryptionAlgorithmTypeAES256
)

// response models:

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
