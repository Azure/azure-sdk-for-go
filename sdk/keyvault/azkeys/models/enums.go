//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package models

import "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"

// DeletionRecoveryLevel - Reflects the deletion recovery level currently in effect for certificates in the current vault. If it contains 'Purgeable', the
// certificate can be permanently deleted by a privileged user; otherwise,
// only the system can purge the certificate, at the end of the retention interval.
// Access the predefined values through the azkeys.DeletionRecoveryLevels instance.
type DeletionRecoveryLevel = generated.DeletionRecoveryLevel

// JSONWebKeyCurveName - Elliptic curve name. For valid values, see JsonWebKeyCurveName.
// Access the predefined values through the azkeys.JSONWebKeyCurveNames instance.
type JSONWebKeyCurveName = generated.JSONWebKeyCurveName

// JSONWebKeyOperation - JSON web key operations. For more information, see JsonWebKeyOperation.
// Access the predefined values through the azkeys.JSONWebKeyOperations instance.
type JSONWebKeyOperation = generated.JSONWebKeyOperation

// ActionType - The type of the action.
// Access the predefined values through the azkeys.ActionTypes instance.
type ActionType = generated.ActionType

// KeyEncryptionAlgorithm - The encryption algorithm to use to protected the exported key material
// Access the predefined values through the azkeys.KeyEncryptionAlgorithms instance.
type KeyEncryptionAlgorithm = generated.KeyEncryptionAlgorithm

// KeyType - JsonWebKey Key Type (kty), as defined in https://tools.ietf.org/html/draft-ietf-jose-json-web-algorithms-40.
// Access the predefined values through the azkeys.KeyTypes instance.
type KeyType string

// EC - Elliptic Curve.
func (KeyType) EC() KeyType {
	return "EC"
}

// ECHSM - Elliptic Curve with a private key which is not exportable from the HSM.
func (KeyType) ECHSM() KeyType {
	return "EC-HSM"
}

// Oct - Octet sequence (used to represent symmetric keys)
func (KeyType) Oct() KeyType {
	return "oct"
}

// OctHSM - Octet sequence (used to represent symmetric keys) which is not exportable from the HSM.
func (KeyType) OctHSM() KeyType {
	return "oct-HSM"
}

// RSA - RSA (https://tools.ietf.org/html/rfc3447)
func (KeyType) RSA() KeyType {
	return "RSA"
}

// RSAHSM - RSA with a private key which is not exportable from the HSM.
func (KeyType) RSAHSM() KeyType {
	return "RSA-HSM"
}

// Values returns the predefined values for the KeyType type.
func (KeyType) Values() []KeyType {
	return []KeyType{
		"EC",
		"EC-HSM",
		"oct",
		"oct-HSM",
		"RSA",
		"RSA-HSM",
	}
}
