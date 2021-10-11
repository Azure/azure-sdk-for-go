//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal"
)

// Attributes - The object attributes managed by the KeyVault service.
type Attributes struct {
	// Determines whether the object is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Expiry date in UTC.
	Expires *time.Time `json:"exp,omitempty"`

	// Not before date in UTC.
	NotBefore *time.Time `json:"nbf,omitempty"`

	// READ-ONLY; Creation time in UTC.
	Created *time.Time `json:"created,omitempty" azure:"ro"`

	// READ-ONLY; Last updated time in UTC.
	Updated *time.Time `json:"updated,omitempty" azure:"ro"`
}

// KeyAttributes - The attributes of a key managed by the key vault service.
type KeyAttributes struct {
	Attributes
	// READ-ONLY; softDelete data retention days. Value should be >=7 and <=90 when softDelete enabled, otherwise 0.
	RecoverableDays *int32 `json:"recoverableDays,omitempty" azure:"ro"`

	// READ-ONLY; Reflects the deletion recovery level currently in effect for keys in the current vault. If it contains 'Purgeable' the key can be permanently
	// deleted by a privileged user; otherwise, only the system
	// can purge the key, at the end of the retention interval.
	RecoveryLevel *DeletionRecoveryLevel `json:"recoveryLevel,omitempty" azure:"ro"`
}

func keyAttributesFromGenerated(i *internal.KeyAttributes) *KeyAttributes {
	if i == nil {
		return &KeyAttributes{}
	}

	return &KeyAttributes{
		RecoverableDays: i.RecoverableDays,
		RecoveryLevel:   DeletionRecoveryLevel(*i.RecoveryLevel).ToPtr(),
		Attributes: Attributes{
			Enabled:   i.Enabled,
			Expires:   i.Expires,
			NotBefore: i.NotBefore,
			Created:   i.Created,
			Updated:   i.Updated,
		},
	}
}

// KeyBundle - A KeyBundle consisting of a WebKey plus its attributes.
type KeyBundle struct {
	// The key management attributes.
	Attributes *KeyAttributes `json:"attributes,omitempty"`

	// The Json web key.
	Key *JSONWebKey `json:"key,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; True if the key's lifetime is managed by key vault. If this is a key backing a certificate, then managed will be true.
	Managed *bool `json:"managed,omitempty" azure:"ro"`
}

// JSONWebKey - As of http://tools.ietf.org/html/draft-ietf-jose-json-web-key-18
type JSONWebKey struct {
	// Elliptic curve name. For valid values, see JsonWebKeyCurveName.
	Crv *JSONWebKeyCurveName `json:"crv,omitempty"`

	// RSA private exponent, or the D component of an EC private key.
	D []byte `json:"d,omitempty"`

	// RSA private key parameter.
	DP []byte `json:"dp,omitempty"`

	// RSA private key parameter.
	DQ []byte `json:"dq,omitempty"`

	// RSA public exponent.
	E []byte `json:"e,omitempty"`

	// Symmetric key.
	K      []byte    `json:"k,omitempty"`
	KeyOps []*string `json:"key_ops,omitempty"`

	// Key identifier.
	Kid *string `json:"kid,omitempty"`

	// JsonWebKey Key Type (kty), as defined in https://tools.ietf.org/html/draft-ietf-jose-json-web-algorithms-40.
	Kty *JSONWebKeyType `json:"kty,omitempty"`

	// RSA modulus.
	N []byte `json:"n,omitempty"`

	// RSA secret prime.
	P []byte `json:"p,omitempty"`

	// RSA secret prime, with p < q.
	Q []byte `json:"q,omitempty"`

	// RSA private key parameter.
	QI []byte `json:"qi,omitempty"`

	// Protected Key, used with 'Bring Your Own Key'.
	T []byte `json:"key_hsm,omitempty"`

	// X component of an EC public key.
	X []byte `json:"x,omitempty"`

	// Y component of an EC public key.
	Y []byte `json:"y,omitempty"`
}

func jsonWebKeyToGenerated(i *internal.JSONWebKey) *JSONWebKey {
	if i == nil {
		return &JSONWebKey{}
	}

	return &JSONWebKey{
		Crv:    (*JSONWebKeyCurveName)(i.Crv),
		D:      i.D,
		DP:     i.DP,
		DQ:     i.DQ,
		E:      i.E,
		K:      i.K,
		KeyOps: i.KeyOps,
		Kid:    i.Kid,
		Kty:    (*JSONWebKeyType)(i.Kty),
		N:      i.N,
		P:      i.P,
		Q:      i.Q,
		QI:     i.QI,
		T:      i.T,
		X:      i.X,
		Y:      i.Y,
	}
}

// JSONWebKeyType - JsonWebKey Key Type (kty), as defined in https://tools.ietf.org/html/draft-ietf-jose-json-web-algorithms-40.
type JSONWebKeyType string

const (
	// JSONWebKeyTypeEC - Elliptic Curve.
	JSONWebKeyTypeEC JSONWebKeyType = "EC"

	// JSONWebKeyTypeECHSM - Elliptic Curve with a private key which is not exportable from the HSM.
	JSONWebKeyTypeECHSM JSONWebKeyType = "EC-HSM"

	// JSONWebKeyTypeOct - Octet sequence (used to represent symmetric keys)
	JSONWebKeyTypeOct JSONWebKeyType = "oct"

	// JSONWebKeyTypeOctHSM - Octet sequence (used to represent symmetric keys) which is not exportable from the HSM.
	JSONWebKeyTypeOctHSM JSONWebKeyType = "oct-HSM"

	// JSONWebKeyTypeRSA - RSA (https://tools.ietf.org/html/rfc3447)
	JSONWebKeyTypeRSA JSONWebKeyType = "RSA"

	// JSONWebKeyTypeRSAHSM - RSA with a private key which is not exportable from the HSM.
	JSONWebKeyTypeRSAHSM JSONWebKeyType = "RSA-HSM"
)

func (j JSONWebKeyType) toGenerated() *internal.JSONWebKeyType {
	return internal.JSONWebKeyType(j).ToPtr()
}

// KeyItem - The key item containing key metadata.
type KeyItem struct {
	// The key management attributes.
	Attributes *KeyAttributes `json:"attributes,omitempty"`

	// Key identifier.
	Kid *string `json:"kid,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`

	// READ-ONLY; True if the key's lifetime is managed by key vault. If this is a key backing a certificate, then managed will be true.
	Managed *bool `json:"managed,omitempty" azure:"ro"`
}

func keyItemFromGenerated(i *internal.KeyItem) *KeyItem {
	if i == nil {
		return nil
	}

	return &KeyItem{
		Attributes: keyAttributesFromGenerated(i.Attributes),
		Kid:        i.Kid,
		Tags:       i.Tags,
		Managed:    i.Managed,
	}
}
