//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto

import "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/crypto/models"

// KeyWrapAlgorithm provides access to the predefined values of KeyWrapAlgorithm.
const KeyWrapAlgorithms = models.KeyWrapAlgorithm("")

// EncryptionAlgorithm provides access to the predefined values of EncryptionAlgorithm.
const EncryptionAlgorithms = models.EncryptionAlgorithm("")

// SignatureAlgorithms provides access to the predefined values of SignatureAlgorithm.
const SignatureAlgorithms = models.SignatureAlgorithm("")
