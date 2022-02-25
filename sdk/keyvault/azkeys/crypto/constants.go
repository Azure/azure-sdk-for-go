//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto

// Key wrapping algorithms
type KeyWrapAlgorithm string

const (
	AES128     KeyWrapAlgorithm = "A128KW"
	AES192     KeyWrapAlgorithm = "A192KW"
	AES256     KeyWrapAlgorithm = "A256KW"
	RSAOAEP    KeyWrapAlgorithm = "RSA-OAEP"
	RSAOAEP256 KeyWrapAlgorithm = "RSA-OAEP-256"
	RSA15      KeyWrapAlgorithm = "RSA1_5"
)

// Returns a pointer to a KeyWrapAlgorithm constant
func (k KeyWrapAlgorithm) ToPtr() *KeyWrapAlgorithm {
	return &k
}

// EncryptionAlgorithm - algorithm identifier
type EncryptionAlgorithm string

const (
	AlgorithmA128CBC    EncryptionAlgorithm = "A128CBC"
	AlgorithmA128CBCPAD EncryptionAlgorithm = "A128CBCPAD"
	AlgorithmA128GCM    EncryptionAlgorithm = "A128GCM"
	AlgorithmA128KW     EncryptionAlgorithm = "A128KW"
	AlgorithmA192CBC    EncryptionAlgorithm = "A192CBC"
	AlgorithmA192CBCPAD EncryptionAlgorithm = "A192CBCPAD"
	AlgorithmA192GCM    EncryptionAlgorithm = "A192GCM"
	AlgorithmA192KW     EncryptionAlgorithm = "A192KW"
	AlgorithmA256CBC    EncryptionAlgorithm = "A256CBC"
	AlgorithmA256CBCPAD EncryptionAlgorithm = "A256CBCPAD"
	AlgorithmA256GCM    EncryptionAlgorithm = "A256GCM"
	AlgorithmA256KW     EncryptionAlgorithm = "A256KW"
	AlgorithmRSA15      EncryptionAlgorithm = "RSA1_5"
	AlgorithmRSAOAEP    EncryptionAlgorithm = "RSA-OAEP"
	AlgorithmRSAOAEP256 EncryptionAlgorithm = "RSA-OAEP-256"
)

// ToPtr returns a *EncryptionAlgorithm pointing to the current value.
func (c EncryptionAlgorithm) ToPtr() *EncryptionAlgorithm {
	return &c
}

// SignatureAlgorithm - The signing/verification algorithm identifier.
type SignatureAlgorithm string

const (
	// ES256 - ECDSA using P-256 and SHA-256, as described in https://tools.ietf.org/html/rfc7518.
	ES256 SignatureAlgorithm = "ES256"
	// ES256K - ECDSA using P-256K and SHA-256, as described in https://tools.ietf.org/html/rfc7518
	ES256K SignatureAlgorithm = "ES256K"
	// ES384 - ECDSA using P-384 and SHA-384, as described in https://tools.ietf.org/html/rfc7518
	ES384 SignatureAlgorithm = "ES384"
	// ES512 - ECDSA using P-521 and SHA-512, as described in https://tools.ietf.org/html/rfc7518
	ES512 SignatureAlgorithm = "ES512"
	// PS256 - RSASSA-PSS using SHA-256 and MGF1 with SHA-256, as described in https://tools.ietf.org/html/rfc7518
	PS256 SignatureAlgorithm = "PS256"
	// PS384 - RSASSA-PSS using SHA-384 and MGF1 with SHA-384, as described in https://tools.ietf.org/html/rfc7518
	PS384 SignatureAlgorithm = "PS384"
	// PS512 - RSASSA-PSS using SHA-512 and MGF1 with SHA-512, as described in https://tools.ietf.org/html/rfc7518
	PS512 SignatureAlgorithm = "PS512"
	// RS256 - RSASSA-PKCS1-v1_5 using SHA-256, as described in https://tools.ietf.org/html/rfc7518
	RS256 SignatureAlgorithm = "RS256"
	// RS384 - RSASSA-PKCS1-v1_5 using SHA-384, as described in https://tools.ietf.org/html/rfc7518
	RS384 SignatureAlgorithm = "RS384"
	// RS512 - RSASSA-PKCS1-v1_5 using SHA-512, as described in https://tools.ietf.org/html/rfc7518
	RS512 SignatureAlgorithm = "RS512"
	// RSNULL - Reserved
	RSNULL SignatureAlgorithm = "RSNULL"
)

// ToPtr returns a *SignatureAlgorithm pointing to the current value.
func (c SignatureAlgorithm) ToPtr() *SignatureAlgorithm {
	return &c
}
