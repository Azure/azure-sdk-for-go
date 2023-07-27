//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

type EncryptionAlgorithmType = generated.EncryptionAlgorithmType

const (
	EncryptionAlgorithmTypeNone   EncryptionAlgorithmType = generated.EncryptionAlgorithmTypeNone
	EncryptionAlgorithmTypeAES256 EncryptionAlgorithmType = generated.EncryptionAlgorithmTypeAES256
)

type ImmutabilityPolicyMode = blob.ImmutabilityPolicyMode

const (
	ImmutabilityPolicyModeMutable  ImmutabilityPolicyMode = blob.ImmutabilityPolicyModeMutable
	ImmutabilityPolicyModeUnlocked ImmutabilityPolicyMode = blob.ImmutabilityPolicyModeUnlocked
	ImmutabilityPolicyModeLocked   ImmutabilityPolicyMode = blob.ImmutabilityPolicyModeLocked
)

// CopyStatusType defines values for CopyStatusType
type CopyStatusType = blob.CopyStatusType

const (
	CopyStatusTypePending CopyStatusType = blob.CopyStatusTypePending
	CopyStatusTypeSuccess CopyStatusType = blob.CopyStatusTypeSuccess
	CopyStatusTypeAborted CopyStatusType = blob.CopyStatusTypeAborted
	CopyStatusTypeFailed  CopyStatusType = blob.CopyStatusTypeFailed
)
