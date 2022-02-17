//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/models"
)

// DeletionRecoveryLevel - Reflects the deletion recovery level currently in effect for certificates in the current vault. If it contains 'Purgeable', the
// certificate can be permanently deleted by a privileged user; otherwise,
// only the system can purge the certificate, at the end of the retention interval.
//
// * Access the predefined values through the DeletionRecoveryLevelValues instance.
type DeletionRecoveryLevel = models.DeletionRecoveryLevel

// DeletionRecoveryLevelValues provides access to the predefined values of DeletionRecoveryLevel.
const DeletionRecoveryLevelValues = DeletionRecoveryLevel("")

// JSONWebKeyCurveName - Elliptic curve name. For valid values, see JsonWebKeyCurveName.
//
// * Access the predefined values through the JSONWebKeyCurveNameValues instance.
type JSONWebKeyCurveName = models.JSONWebKeyCurveName

// JSONWebKeyCurveNameValues provides access to the predefined values of JSONWebKeyCurveName.
const JSONWebKeyCurveNameValues = JSONWebKeyCurveName("")

// JSONWebKeyOperation - JSON web key operations. For more information, see JsonWebKeyOperation.
//
// * Access the predefined values through the JSONWebKeyOperationValues instance.
type JSONWebKeyOperation = models.JSONWebKeyOperation

// JSONWebKeyOperationValues provides access to the predefined values of JSONWebKeyOperation.
const JSONWebKeyOperationValues = JSONWebKeyOperation("")

// ActionType - The type of the action.
//
// * Access the predefined values through the ActionTypeValues instance.
type ActionType = models.ActionType

// ActionTypeValues provides access to the predefined values of ActionType.
const ActionTypeValues = ActionType("")

// KeyEncryptionAlgorithm - The encryption algorithm to use to protected the exported key material
//
// * Access the predefined values through the KeyEncryptionAlgorithmValues instance.
type KeyEncryptionAlgorithm = models.KeyEncryptionAlgorithm

// KeyEncryptionAlgorithmValues provides access to the predefined values of KeyEncryptionAlgorithm.
const KeyEncryptionAlgorithmValues = KeyEncryptionAlgorithm("")

// KeyType - JsonWebKey Key Type (kty), as defined in https://tools.ietf.org/html/draft-ietf-jose-json-web-algorithms-40.
//
// * Access the predefined values through the KeyTypeValues instance.
type KeyType = models.KeyType

// KeyTypeValues provides access to the predefined values of KeyType.
const KeyTypeValues = KeyType("")
