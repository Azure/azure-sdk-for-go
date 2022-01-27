//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"

// DeletionRecoveryLevel - Reflects the deletion recovery level currently in effect for certificates in the current vault. If it contains 'Purgeable', the
// certificate can be permanently deleted by a privileged user; otherwise,
// only the system can purge the certificate, at the end of the retention interval.
type DeletionRecoveryLevel = generated.DeletionRecoveryLevel

// DeletionRecoveryLevels provides access to the predefined values of DeletionRecoveryLevel.
const DeletionRecoveryLevels = DeletionRecoveryLevel("")

// JSONWebKeyCurveName - Elliptic curve name. For valid values, see JsonWebKeyCurveName.
type JSONWebKeyCurveName = generated.JSONWebKeyCurveName

// JSONWebKeyCurveNames provides access to the predefined values of JSONWebKeyCurveName.
const JSONWebKeyCurveNames = JSONWebKeyCurveName("")

// JSONWebKeyOperation - JSON web key operations. For more information, see JsonWebKeyOperation.
type JSONWebKeyOperation = generated.JSONWebKeyOperation

// JSONWebKeyOperations provides access to the predefined values of JSONWebKeyOperation.
const JSONWebKeyOperations = JSONWebKeyOperation("")

// ActionType - The type of the action.
type ActionType = generated.ActionType

// ActionTypes provides access to the predefined values of ActionType.
const ActionTypes = generated.ActionType("")

// KeyEncryptionAlgorithm - The encryption algorithm to use to protected the exported key material
type KeyEncryptionAlgorithm = generated.KeyEncryptionAlgorithm

// KeyEncryptionAlgorithms provides access to the predefined values of KeyEncryptionAlgorithm.
const KeyEncryptionAlgorithms = generated.KeyEncryptionAlgorithm("")
