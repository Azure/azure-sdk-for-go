//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/models"
)

// Attributes - The object attributes managed by the KeyVault service.
type Attributes = models.Attributes

// KeyAttributes - The attributes of a key managed by the key vault service.
type KeyAttributes = models.KeyAttributes

// KeyBundle - A KeyBundle consisting of a WebKey plus its attributes.
type KeyBundle = models.KeyBundle

// JSONWebKey - As of http://tools.ietf.org/html/draft-ietf-jose-json-web-key-18
type JSONWebKey = models.JSONWebKey

// KeyItem - The key item containing key metadata.
type KeyItem = models.KeyItem

// DeletedKeyBundle - A DeletedKeyBundle consisting of a WebKey plus its Attributes and deletion info
type DeletedKeyBundle = models.DeletedKeyBundle

// DeletedKeyItem - The deleted key item containing the deleted key metadata and information about deletion.
type DeletedKeyItem = models.DeletedKeyItem

type KeyReleasePolicy = models.KeyReleasePolicy

// KeyRotationPolicy - Management policy for a key.
type KeyRotationPolicy = models.KeyRotationPolicy

// KeyRotationPolicyAttributes - The key rotation policy attributes.
type KeyRotationPolicyAttributes = models.KeyRotationPolicyAttributes

// LifetimeActions - Action and its trigger that will be performed by Key Vault over the lifetime of a key.
type LifetimeActions = models.LifetimeActions

// LifetimeActionsType - The action that will be executed.
type LifetimeActionsType = models.LifetimeActionsType

// LifetimeActionsTrigger - A condition to be satisfied for an action to be executed.
type LifetimeActionsTrigger = models.LifetimeActionsTrigger
