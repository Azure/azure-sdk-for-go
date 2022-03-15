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

// PossibleDeletionRecoveryLevelValues provides a slice of all possible DeletionRecoveryLevels
func PossibleDeletionRecoveryLevelValues() []DeletionRecoveryLevel {
	return []DeletionRecoveryLevel{
		DeletionRecoveryLevelCustomizedRecoverable,
		DeletionRecoveryLevelCustomizedRecoverableProtectedSubscription,
		DeletionRecoveryLevelCustomizedRecoverablePurgeable,
		DeletionRecoveryLevelPurgeable,
		DeletionRecoveryLevelRecoverable,
		DeletionRecoveryLevelRecoverableProtectedSubscription,
		DeletionRecoveryLevelRecoverablePurgeable,
	}
}

// CurveName - Elliptic curve name. For valid values, see PossibleCurveNameValues.
type CurveName string

const (
	// CurveNameP256 - The NIST P-256 elliptic curve, AKA SECG curve SECP256R1.
	CurveNameP256 CurveName = "P-256"

	// CurveNameP256K - The SECG SECP256K1 elliptic curve.
	CurveNameP256K CurveName = "P-256K"

	// CurveNameP384 - The NIST P-384 elliptic curve, AKA SECG curve SECP384R1.
	CurveNameP384 CurveName = "P-384"

	// CurveNameP521 - The NIST P-521 elliptic curve, AKA SECG curve SECP521R1.
	CurveNameP521 CurveName = "P-521"
)

// ToPtr returns a *CurveName pointing to the current value.
func (c CurveName) ToPtr() *CurveName {
	return &c
}

// PossibleCurveNameValues provides a slice of all possible CurveNames
func PossibleCurveNameValues() []CurveName {
	return []CurveName{
		CurveNameP256,
		CurveNameP256K,
		CurveNameP384,
		CurveNameP521,
	}
}

// Operation - JSON web key operations. For more information, see Operation.
type Operation string

const (
	OperationDecrypt   Operation = "decrypt"
	OperationEncrypt   Operation = "encrypt"
	OperationImport    Operation = "import"
	OperationSign      Operation = "sign"
	OperationUnwrapKey Operation = "unwrapKey"
	OperationVerify    Operation = "verify"
	OperationWrapKey   Operation = "wrapKey"
)

// ToPtr returns a *KeyOperation pointing to the current value.
func (c Operation) ToPtr() *Operation {
	return &c
}

// PossibleOperationValues provides a slice of all possible Operations
func PossibleOperationValues() []Operation {
	return []Operation{
		OperationDecrypt,
		OperationEncrypt,
		OperationImport,
		OperationSign,
		OperationUnwrapKey,
		OperationVerify,
		OperationWrapKey,
	}
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

// PossibleActionTypeValues provides a slice of all possible ActionTypes
func PossibleActionTypeValues() []ActionType {
	return []ActionType{
		ActionTypeRotate,
		ActionTypeNotify,
	}
}

// ExportEncryptionAlgorithm - The encryption algorithm to use to protected the exported key material
type ExportEncryptionAlgorithm string

const (
	ExportEncryptionAlgorithmCKMRSAAESKEYWRAP ExportEncryptionAlgorithm = "CKM_RSA_AES_KEY_WRAP"
	ExportEncryptionAlgorithmRSAAESKEYWRAP256 ExportEncryptionAlgorithm = "RSA_AES_KEY_WRAP_256"
	ExportEncryptionAlgorithmRSAAESKEYWRAP384 ExportEncryptionAlgorithm = "RSA_AES_KEY_WRAP_384"
)

// ToPtr returns a *ExportEncryptionAlgorithm pointing to the current value.
func (c ExportEncryptionAlgorithm) ToPtr() *ExportEncryptionAlgorithm {
	return &c
}

// PossibleExportEncryptionAlgorithmValues provides a slice of all possible ExportEncryptionAlgorithms
func PossibleExportEncryptionAlgorithmValues() []ExportEncryptionAlgorithm {
	return []ExportEncryptionAlgorithm{
		ExportEncryptionAlgorithmCKMRSAAESKEYWRAP,
		ExportEncryptionAlgorithmRSAAESKEYWRAP256,
		ExportEncryptionAlgorithmRSAAESKEYWRAP384,
	}
}

// KeyType - JsonWebKey Key Type (kty), as defined in https://tools.ietf.org/html/draft-ietf-jose-json-web-algorithms-40.
type KeyType string

const (
	// EC - Elliptic Curve.
	KeyTypeEC KeyType = "EC"

	// ECHSM - Elliptic Curve with a private key which is not exportable from the HSM.
	KeyTypeECHSM KeyType = "EC-HSM"

	// Oct - Octet sequence (used to represent symmetric keys)
	KeyTypeOct KeyType = "oct"

	// OctHSM - Octet sequence (used to represent symmetric keys) which is not exportable from the HSM.
	KeyTypeOctHSM KeyType = "oct-HSM"

	// RSA - RSA (https://tools.ietf.org/html/rfc3447)
	KeyTypeRSA KeyType = "RSA"

	// RSAHSM - RSA with a private key which is not exportable from the HSM.
	KeyTypeRSAHSM KeyType = "RSA-HSM"
)

// ToPtr returns a pointer to a KeyType
func (k KeyType) ToPtr() *KeyType {
	return &k
}

// PossibleKeyTypeValues provides a slice of all possible KeyTypes
func PossibleKeyTypeValues() []KeyType {
	return []KeyType{
		KeyTypeEC,
		KeyTypeECHSM,
		KeyTypeOct,
		KeyTypeOctHSM,
		KeyTypeRSA,
		KeyTypeRSAHSM,
	}
}

// convert KeyType to *generated.JSONWebKeyType
func (j KeyType) toGenerated() *generated.JSONWebKeyType {
	return generated.JSONWebKeyType(j).ToPtr()
}
