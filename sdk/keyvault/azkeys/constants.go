//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"

// DeletionRecoveryLevel - Reflects the deletion recovery level currently in effect for certificates in the current vault. If it contains 'Purgeable', the
// certificate can be permanently deleted by a privileged user; otherwise,
// only the system can purge the certificate, at the end of the retention interval.
type DeletionRecoveryLevel string

const (
	// DeletionRecoveryLevelCustomizedRecoverable - Denotes a vault state in which deletion is recoverable without the possibility for immediate and permanent
	// deletion (i.e. purge when 7<= SoftDeleteRetentionInDays < 90).This level guarantees the recoverability of the deleted entity during the retention interval
	// and while the subscription is still available.
	DeletionRecoveryLevelCustomizedRecoverable DeletionRecoveryLevel = "CustomizedRecoverable"

	// DeletionRecoveryLevelCustomizedRecoverableProtectedSubscription - Denotes a vault and subscription state in which deletion is recoverable, immediate
	// and permanent deletion (i.e. purge) is not permitted, and in which the subscription itself cannot be permanently canceled when 7<= SoftDeleteRetentionInDays
	// < 90. This level guarantees the recoverability of the deleted entity during the retention interval, and also reflects the fact that the subscription
	// itself cannot be cancelled.
	DeletionRecoveryLevelCustomizedRecoverableProtectedSubscription DeletionRecoveryLevel = "CustomizedRecoverable+ProtectedSubscription"

	// DeletionRecoveryLevelCustomizedRecoverablePurgeable - Denotes a vault state in which deletion is recoverable, and which also permits immediate and permanent
	// deletion (i.e. purge when 7<= SoftDeleteRetentionInDays < 90). This level guarantees the recoverability of the deleted entity during the retention interval,
	// unless a Purge operation is requested, or the subscription is cancelled.
	DeletionRecoveryLevelCustomizedRecoverablePurgeable DeletionRecoveryLevel = "CustomizedRecoverable+Purgeable"

	// DeletionRecoveryLevelPurgeable - Denotes a vault state in which deletion is an irreversible operation, without the possibility for recovery. This level
	// corresponds to no protection being available against a Delete operation; the data is irretrievably lost upon accepting a Delete operation at the entity
	// level or higher (vault, resource group, subscription etc.)
	DeletionRecoveryLevelPurgeable DeletionRecoveryLevel = "Purgeable"

	// DeletionRecoveryLevelRecoverable - Denotes a vault state in which deletion is recoverable without the possibility for immediate and permanent deletion
	// (i.e. purge). This level guarantees the recoverability of the deleted entity during the retention interval(90 days) and while the subscription is still
	// available. System wil permanently delete it after 90 days, if not recovered
	DeletionRecoveryLevelRecoverable DeletionRecoveryLevel = "Recoverable"

	// DeletionRecoveryLevelRecoverableProtectedSubscription - Denotes a vault and subscription state in which deletion is recoverable within retention interval
	// (90 days), immediate and permanent deletion (i.e. purge) is not permitted, and in which the subscription itself cannot be permanently canceled. System
	// wil permanently delete it after 90 days, if not recovered
	DeletionRecoveryLevelRecoverableProtectedSubscription DeletionRecoveryLevel = "Recoverable+ProtectedSubscription"

	// DeletionRecoveryLevelRecoverablePurgeable - Denotes a vault state in which deletion is recoverable, and which also permits immediate and permanent deletion
	// (i.e. purge). This level guarantees the recoverability of the deleted entity during the retention interval (90 days), unless a Purge operation is requested,
	// or the subscription is cancelled. System wil permanently delete it after 90 days, if not recovered
	DeletionRecoveryLevelRecoverablePurgeable DeletionRecoveryLevel = "Recoverable+Purgeable"
)

// ToPtr returns a *DeletionRecoveryLevel pointing to the current value.
func (d DeletionRecoveryLevel) ToPtr() *DeletionRecoveryLevel {
	return &d
}

// convert a pointer to exported DeletionRecoveryLevel to the generated version
func recoveryLevelToGenerated(d *DeletionRecoveryLevel) *generated.DeletionRecoveryLevel {
	if d == nil {
		return nil
	}
	if *d == DeletionRecoveryLevelCustomizedRecoverable {
		return generated.DeletionRecoveryLevelCustomizedRecoverable.ToPtr()
	} else if *d == DeletionRecoveryLevelCustomizedRecoverableProtectedSubscription {
		return generated.DeletionRecoveryLevelCustomizedRecoverableProtectedSubscription.ToPtr()
	} else if *d == DeletionRecoveryLevelCustomizedRecoverablePurgeable {
		return generated.DeletionRecoveryLevelCustomizedRecoverablePurgeable.ToPtr()
	} else if *d == DeletionRecoveryLevelPurgeable {
		return generated.DeletionRecoveryLevelPurgeable.ToPtr()
	} else if *d == DeletionRecoveryLevelRecoverableProtectedSubscription {
		return generated.DeletionRecoveryLevelRecoverableProtectedSubscription.ToPtr()
	} else if *d == DeletionRecoveryLevelRecoverable {
		return generated.DeletionRecoveryLevelRecoverable.ToPtr()
	} else {
		return generated.DeletionRecoveryLevelRecoverablePurgeable.ToPtr()
	}
}

// JSONWebKeyCurveName - Elliptic curve name. For valid values, see JsonWebKeyCurveName.
type JSONWebKeyCurveName string

const (
	// JSONWebKeyCurveNameP256 - The NIST P-256 elliptic curve, AKA SECG curve SECP256R1.
	JSONWebKeyCurveNameP256 JSONWebKeyCurveName = "P-256"

	// JSONWebKeyCurveNameP256K - The SECG SECP256K1 elliptic curve.
	JSONWebKeyCurveNameP256K JSONWebKeyCurveName = "P-256K"

	// JSONWebKeyCurveNameP384 - The NIST P-384 elliptic curve, AKA SECG curve SECP384R1.
	JSONWebKeyCurveNameP384 JSONWebKeyCurveName = "P-384"

	// JSONWebKeyCurveNameP521 - The NIST P-521 elliptic curve, AKA SECG curve SECP521R1.
	JSONWebKeyCurveNameP521 JSONWebKeyCurveName = "P-521"
)

// ToPtr returns a *JSONWebKeyCurveName pointing to the current value.
func (c JSONWebKeyCurveName) ToPtr() *JSONWebKeyCurveName {
	return &c
}

// JSONWebKeyOperation - JSON web key operations. For more information, see JsonWebKeyOperation.
type JSONWebKeyOperation string

const (
	JSONWebKeyOperationDecrypt   JSONWebKeyOperation = "decrypt"
	JSONWebKeyOperationEncrypt   JSONWebKeyOperation = "encrypt"
	JSONWebKeyOperationImport    JSONWebKeyOperation = "import"
	JSONWebKeyOperationSign      JSONWebKeyOperation = "sign"
	JSONWebKeyOperationUnwrapKey JSONWebKeyOperation = "unwrapKey"
	JSONWebKeyOperationVerify    JSONWebKeyOperation = "verify"
	JSONWebKeyOperationWrapKey   JSONWebKeyOperation = "wrapKey"
)

// ToPtr returns a *JSONWebKeyOperation pointing to the current value.
func (c JSONWebKeyOperation) ToPtr() *JSONWebKeyOperation {
	return &c
}

// ActionType - The type of the action.
type ActionType string

const (
	// ActionTypeRotate - Rotate the key based on the key policy.
	ActionTypeRotate ActionType = "rotate"
	// ActionTypeNotify - Trigger event grid events. For preview, the notification time is not configurable and it is default to 30 days before expiry.
	ActionTypeNotify ActionType = "notify"
)

// ToPtr returns a *ActionType pointing to the current value.
func (c ActionType) ToPtr() *ActionType {
	return &c
}

// KeyEncryptionAlgorithm - The encryption algorithm to use to protected the exported key material
type KeyEncryptionAlgorithm string

const (
	KeyEncryptionAlgorithmCKMRSAAESKEYWRAP KeyEncryptionAlgorithm = "CKM_RSA_AES_KEY_WRAP"
	KeyEncryptionAlgorithmRSAAESKEYWRAP256 KeyEncryptionAlgorithm = "RSA_AES_KEY_WRAP_256"
	KeyEncryptionAlgorithmRSAAESKEYWRAP384 KeyEncryptionAlgorithm = "RSA_AES_KEY_WRAP_384"
)

// ToPtr returns a *KeyEncryptionAlgorithm pointing to the current value.
func (c KeyEncryptionAlgorithm) ToPtr() *KeyEncryptionAlgorithm {
	return &c
}
