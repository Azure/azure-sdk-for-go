//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"
)

func toGeneratedDeletionRecoveryLevel(s *string) *generated.DeletionRecoveryLevel {
	if s == nil {
		return nil
	}
	return to.Ptr(generated.DeletionRecoveryLevel(*s))
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

// RotationAction - The type of the action.
type RotationAction string

const (
	// RotationActionRotate - Rotate the key based on the key policy.
	RotationActionRotate RotationAction = "rotate"
	// RotationActionNotify - Trigger event grid events. For preview, the notification time is not configurable and it is default to 30 days before expiry.
	RotationActionNotify RotationAction = "notify"
)

// PossibleActionTypeValues provides a slice of all possible ActionTypes
func PossibleActionTypeValues() []RotationAction {
	return []RotationAction{
		RotationActionRotate,
		RotationActionNotify,
	}
}

// ExportEncryptionAlg - The encryption algorithm to use to protected the exported key material
type ExportEncryptionAlg string

const (
	ExportEncryptionAlgCKMRSAAESKEYWRAP ExportEncryptionAlg = "CKM_RSA_AES_KEY_WRAP"
	ExportEncryptionAlgRSAAESKEYWRAP256 ExportEncryptionAlg = "RSA_AES_KEY_WRAP_256"
	ExportEncryptionAlgRSAAESKEYWRAP384 ExportEncryptionAlg = "RSA_AES_KEY_WRAP_384"
)

// PossibleExportEncryptionAlgValues provides a slice of all possible ExportEncryptionAlgs
func PossibleExportEncryptionAlgValues() []ExportEncryptionAlg {
	return []ExportEncryptionAlg{
		ExportEncryptionAlgCKMRSAAESKEYWRAP,
		ExportEncryptionAlgRSAAESKEYWRAP256,
		ExportEncryptionAlgRSAAESKEYWRAP384,
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
	return to.Ptr(generated.JSONWebKeyType(j))
}
