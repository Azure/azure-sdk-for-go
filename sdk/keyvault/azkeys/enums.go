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
// * Access the predefined values through the DeletionRecoveryLevels instance.
type DeletionRecoveryLevel = models.DeletionRecoveryLevel

// DeletionRecoveryLevels provides access to the predefined values of DeletionRecoveryLevel.
const DeletionRecoveryLevels = DeletionRecoveryLevel("")

// JSONWebKeyCurveName - Elliptic curve name. For valid values, see JsonWebKeyCurveName.
// * Access the predefined values through the JSONWebKeyCurveNames instance.
type JSONWebKeyCurveName = models.JSONWebKeyCurveName

// JSONWebKeyCurveNames provides access to the predefined values of JSONWebKeyCurveName.
const JSONWebKeyCurveNames = JSONWebKeyCurveName("")

// JSONWebKeyOperation - JSON web key operations. For more information, see JsonWebKeyOperation.
// * Access the predefined values through the JSONWebKeyOperations instance.
type JSONWebKeyOperation = models.JSONWebKeyOperation

// JSONWebKeyOperations provides access to the predefined values of JSONWebKeyOperation.
const JSONWebKeyOperations = JSONWebKeyOperation("")

// ActionType - The type of the action.
// * Access the predefined values through the ActionTypes instance.
type ActionType = models.ActionType

// ActionTypes provides access to the predefined values of ActionType.
const ActionTypes = ActionType("")

// KeyEncryptionAlgorithm - The encryption algorithm to use to protected the exported key material
// * Access the predefined values through the KeyEncryptionAlgorithms instance.
type KeyEncryptionAlgorithm = models.KeyEncryptionAlgorithm

// KeyEncryptionAlgorithms provides access to the predefined values of KeyEncryptionAlgorithm.
const KeyEncryptionAlgorithms = KeyEncryptionAlgorithm("")

// KeyType - JsonWebKey Key Type (kty), as defined in https://tools.ietf.org/html/draft-ietf-jose-json-web-algorithms-40.
// * Access the predefined values through the KeyTypes instance.
type KeyType = models.KeyType

// KeyTypes provides access to the predefined values of KeyType.
const KeyTypes = KeyType("")
