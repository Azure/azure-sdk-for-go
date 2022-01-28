//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/models"
)

// DeletionRecoveryLevels provides access to the predefined values of DeletionRecoveryLevel.
const DeletionRecoveryLevels = models.DeletionRecoveryLevel("")

// JSONWebKeyCurveNames provides access to the predefined values of JSONWebKeyCurveName.
const JSONWebKeyCurveNames = models.JSONWebKeyCurveName("")

// JSONWebKeyOperations provides access to the predefined values of JSONWebKeyOperation.
const JSONWebKeyOperations = models.JSONWebKeyOperation("")

// ActionTypes provides access to the predefined values of ActionType.
const ActionTypes = models.ActionType("")

// KeyEncryptionAlgorithms provides access to the predefined values of KeyEncryptionAlgorithm.
const KeyEncryptionAlgorithms = models.KeyEncryptionAlgorithm("")

// KeyTypes provides access to the predefined values of KeyType.
const KeyTypes = models.KeyType("")
