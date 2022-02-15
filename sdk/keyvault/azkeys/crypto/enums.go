//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto

import "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/crypto/models"

// Key wrapping algorithms
// * Access the predefined values through the KeyWrapAlgorithms instance.
type KeyWrapAlgorithm = models.KeyWrapAlgorithm

// KeyWrapAlgorithm provides access to the predefined values of KeyWrapAlgorithm.
const KeyWrapAlgorithms = KeyWrapAlgorithm("")

// EncryptionAlgorithm - algorithm identifier
// * Access the predefined values through the EncryptionAlgorithms instance.
type EncryptionAlgorithm = models.EncryptionAlgorithm

// EncryptionAlgorithm provides access to the predefined values of EncryptionAlgorithm.
const EncryptionAlgorithms = EncryptionAlgorithm("")

// SignatureAlgorithm - The signing/verification algorithm identifier.
// * Access the predefined values through the SignatureAlgorithms instance.
type SignatureAlgorithm = models.SignatureAlgorithm

// SignatureAlgorithms provides access to the predefined values of SignatureAlgorithm.
const SignatureAlgorithms = SignatureAlgorithm("")
