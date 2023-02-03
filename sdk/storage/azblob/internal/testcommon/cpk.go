//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Contains common helpers for TESTS ONLY
package testcommon

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/stretchr/testify/require"
)

const (
	EncryptionScopeEnvVar = "AZURE_STORAGE_ENCRYPTION_SCOPE"
)

var testEncryptedKey = "MDEyMzQ1NjcwMTIzNDU2NzAxMjM0NTY3MDEyMzQ1Njc="
var testEncryptedHash = "3QFFFpRA5+XANHqwwbT4yXDmrT/2JaLt/FKHjzhOdoE="
var testEncryptionAlgorithm = blob.EncryptionAlgorithmTypeAES256
var TestCPKByValue = blob.CPKInfo{
	EncryptionKey:       &testEncryptedKey,
	EncryptionKeySHA256: &testEncryptedHash,
	EncryptionAlgorithm: &testEncryptionAlgorithm,
}

var testInvalidEncryptedKey = "mumbojumbo"
var testInvalidEncryptedHash = "mumbojumbohash"
var TestInvalidCPKByValue = blob.CPKInfo{
	EncryptionKey:       &testInvalidEncryptedKey,
	EncryptionKeySHA256: &testInvalidEncryptedHash,
	EncryptionAlgorithm: &testEncryptionAlgorithm,
}

func GetCPKScopeInfo(t *testing.T) blob.CPKScopeInfo {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return blob.CPKScopeInfo{EncryptionScope: to.Ptr("blobgokeytestscope")}
	}

	encryptionScope, err := GetRequiredEnv(EncryptionScopeEnvVar)
	require.NoError(t, err)
	return blob.CPKScopeInfo{EncryptionScope: &encryptionScope}
}

var testInvalidEncryptedScope = "mumbojumboscope"
var TestInvalidCPKByScope = blob.CPKScopeInfo{
	EncryptionScope: &testInvalidEncryptedScope,
}
