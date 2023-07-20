//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/path"
)

type EncryptionAlgorithmType = path.EncryptionAlgorithmType

const (
	EncryptionAlgorithmTypeNone   EncryptionAlgorithmType = path.EncryptionAlgorithmTypeNone
	EncryptionAlgorithmTypeAES256 EncryptionAlgorithmType = path.EncryptionAlgorithmTypeAES256
)

// response models:

type ImmutabilityPolicyMode = blob.ImmutabilityPolicyMode

const (
	ImmutabilityPolicyModeMutable  ImmutabilityPolicyMode = blob.ImmutabilityPolicyModeMutable
	ImmutabilityPolicyModeUnlocked ImmutabilityPolicyMode = blob.ImmutabilityPolicyModeUnlocked
	ImmutabilityPolicyModeLocked   ImmutabilityPolicyMode = blob.ImmutabilityPolicyModeLocked
)
