//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
)

type EncryptionAlgorithmType = blob.EncryptionAlgorithmType

const (
	EncryptionAlgorithmTypeNone   EncryptionAlgorithmType = blob.EncryptionAlgorithmTypeNone
	EncryptionAlgorithmTypeAES256 EncryptionAlgorithmType = blob.EncryptionAlgorithmTypeAES256
)
