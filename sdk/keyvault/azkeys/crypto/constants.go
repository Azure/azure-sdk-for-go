//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto

// WrapAlgorithm represents the key wrapping algorithms
type WrapAlgorithm string

const (
	WrapAlgorithmAES128     WrapAlgorithm = "A128KW"
	WrapAlgorithmAES192     WrapAlgorithm = "A192KW"
	WrapAlgorithmAES256     WrapAlgorithm = "A256KW"
	WrapAlgorithmRSAOAEP    WrapAlgorithm = "RSA-OAEP"
	WrapAlgorithmRSAOAEP256 WrapAlgorithm = "RSA-OAEP-256"
	WrapAlgorithmRSA15      WrapAlgorithm = "RSA1_5"
)

// ToPtr returns a pointer to WrapAlgorithm
func (k WrapAlgorithm) ToPtr() *WrapAlgorithm {
	return &k
}

// PossibleWrapAlgorithmValues returns a slice of all possible WrapAlgorithm values
func PossibleWrapAlgorithmValues() []WrapAlgorithm {
	return []WrapAlgorithm{
		WrapAlgorithmAES128,
		WrapAlgorithmAES192,
		WrapAlgorithmAES256,
		WrapAlgorithmRSAOAEP,
		WrapAlgorithmRSAOAEP256,
		WrapAlgorithmRSA15,
	}
}

// EncryptionAlgorithm - algorithm identifier
type EncryptionAlgorithm string

const (
	EncryptionAlgorithmA128CBC    EncryptionAlgorithm = "A128CBC"
	EncryptionAlgorithmA128CBCPAD EncryptionAlgorithm = "A128CBCPAD"
	EncryptionAlgorithmA128GCM    EncryptionAlgorithm = "A128GCM"
	EncryptionAlgorithmA128KW     EncryptionAlgorithm = "A128KW"
	EncryptionAlgorithmA192CBC    EncryptionAlgorithm = "A192CBC"
	EncryptionAlgorithmA192CBCPAD EncryptionAlgorithm = "A192CBCPAD"
	EncryptionAlgorithmA192GCM    EncryptionAlgorithm = "A192GCM"
	EncryptionAlgorithmA192KW     EncryptionAlgorithm = "A192KW"
	EncryptionAlgorithmA256CBC    EncryptionAlgorithm = "A256CBC"
	EncryptionAlgorithmA256CBCPAD EncryptionAlgorithm = "A256CBCPAD"
	EncryptionAlgorithmA256GCM    EncryptionAlgorithm = "A256GCM"
	EncryptionAlgorithmA256KW     EncryptionAlgorithm = "A256KW"
	EncryptionAlgorithmRSA15      EncryptionAlgorithm = "RSA1_5"
	EncryptionAlgorithmRSAOAEP    EncryptionAlgorithm = "RSA-OAEP"
	EncryptionAlgorithmRSAOAEP256 EncryptionAlgorithm = "RSA-OAEP-256"
)

// ToPtr returns a *EncryptionAlgorithm pointing to the current value.
func (c EncryptionAlgorithm) ToPtr() *EncryptionAlgorithm {
	return &c
}

// PossibleEncryptionAlgorithmValues returns a slice of all possible EncryptionAlgorithm values
func PossibleEncryptionAlgorithmValues() []EncryptionAlgorithm {
	return []EncryptionAlgorithm{
		EncryptionAlgorithmA128CBC,
		EncryptionAlgorithmA128CBCPAD,
		EncryptionAlgorithmA128GCM,
		EncryptionAlgorithmA128KW,
		EncryptionAlgorithmA192CBC,
		EncryptionAlgorithmA192CBCPAD,
		EncryptionAlgorithmA192GCM,
		EncryptionAlgorithmA192KW,
		EncryptionAlgorithmA256CBC,
		EncryptionAlgorithmA256CBCPAD,
		EncryptionAlgorithmA256GCM,
		EncryptionAlgorithmA256KW,
		EncryptionAlgorithmRSA15,
		EncryptionAlgorithmRSAOAEP,
		EncryptionAlgorithmRSAOAEP256,
	}
}

// SignatureAlgorithm - The signing/verification algorithm identifier.
type SignatureAlgorithm string

const (
	// ES256 - ECDSA using P-256 and SHA-256, as described in https://tools.ietf.org/html/rfc7518.
	SignatureAlgorithmES256 SignatureAlgorithm = "ES256"
	// ES256K - ECDSA using P-256K and SHA-256, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgorithmES256K SignatureAlgorithm = "ES256K"
	// ES384 - ECDSA using P-384 and SHA-384, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgorithmES384 SignatureAlgorithm = "ES384"
	// ES512 - ECDSA using P-521 and SHA-512, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgorithmES512 SignatureAlgorithm = "ES512"
	// PS256 - RSASSA-PSS using SHA-256 and MGF1 with SHA-256, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgorithmPS256 SignatureAlgorithm = "PS256"
	// PS384 - RSASSA-PSS using SHA-384 and MGF1 with SHA-384, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgorithmPS384 SignatureAlgorithm = "PS384"
	// PS512 - RSASSA-PSS using SHA-512 and MGF1 with SHA-512, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgorithmPS512 SignatureAlgorithm = "PS512"
	// RS256 - RSASSA-PKCS1-v1_5 using SHA-256, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgorithmRS256 SignatureAlgorithm = "RS256"
	// RS384 - RSASSA-PKCS1-v1_5 using SHA-384, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgorithmRS384 SignatureAlgorithm = "RS384"
	// RS512 - RSASSA-PKCS1-v1_5 using SHA-512, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgorithmRS512 SignatureAlgorithm = "RS512"
	// RSNULL - Reserved
	SignatureAlgorithmRSNULL SignatureAlgorithm = "RSNULL"
)

// ToPtr returns a *SignatureAlgorithm pointing to the current value.
func (c SignatureAlgorithm) ToPtr() *SignatureAlgorithm {
	return &c
}

// PossibleSignatureAlgorithmValues returns a slice of all possible SignatureAlgorithm values
func PossibleSignatureAlgorithmValues() []SignatureAlgorithm {
	return []SignatureAlgorithm{
		SignatureAlgorithmES256,
		SignatureAlgorithmES256K,
		SignatureAlgorithmES384,
		SignatureAlgorithmES512,
		SignatureAlgorithmPS256,
		SignatureAlgorithmPS384,
		SignatureAlgorithmPS512,
		SignatureAlgorithmRS256,
		SignatureAlgorithmRS384,
		SignatureAlgorithmRS512,
		SignatureAlgorithmRSNULL,
	}
}
