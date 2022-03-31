//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto

// WrapAlg represents the key wrapping algorithms
type WrapAlg string

const (
	WrapAlgAES128     WrapAlg = "A128KW"
	WrapAlgAES192     WrapAlg = "A192KW"
	WrapAlgAES256     WrapAlg = "A256KW"
	WrapAlgRSAOAEP    WrapAlg = "RSA-OAEP"
	WrapAlgRSAOAEP256 WrapAlg = "RSA-OAEP-256"
	WrapAlgRSA15      WrapAlg = "RSA1_5"
)

// PossibleWrapAlgValues returns a slice of all possible WrapAlg values
func PossibleWrapAlgValues() []WrapAlg {
	return []WrapAlg{
		WrapAlgAES128,
		WrapAlgAES192,
		WrapAlgAES256,
		WrapAlgRSAOAEP,
		WrapAlgRSAOAEP256,
		WrapAlgRSA15,
	}
}

// EncryptionAlg - algorithm identifier
type EncryptionAlg string

const (
	EncryptionAlgA128CBC    EncryptionAlg = "A128CBC"
	EncryptionAlgA128CBCPAD EncryptionAlg = "A128CBCPAD"
	EncryptionAlgA128GCM    EncryptionAlg = "A128GCM"
	EncryptionAlgA128KW     EncryptionAlg = "A128KW"
	EncryptionAlgA192CBC    EncryptionAlg = "A192CBC"
	EncryptionAlgA192CBCPAD EncryptionAlg = "A192CBCPAD"
	EncryptionAlgA192GCM    EncryptionAlg = "A192GCM"
	EncryptionAlgA192KW     EncryptionAlg = "A192KW"
	EncryptionAlgA256CBC    EncryptionAlg = "A256CBC"
	EncryptionAlgA256CBCPAD EncryptionAlg = "A256CBCPAD"
	EncryptionAlgA256GCM    EncryptionAlg = "A256GCM"
	EncryptionAlgA256KW     EncryptionAlg = "A256KW"
	EncryptionAlgRSA15      EncryptionAlg = "RSA1_5"
	EncryptionAlgRSAOAEP    EncryptionAlg = "RSA-OAEP"
	EncryptionAlgRSAOAEP256 EncryptionAlg = "RSA-OAEP-256"
)

// PossibleEncryptionAlgValues returns a slice of all possible EncryptionAlg values
func PossibleEncryptionAlgValues() []EncryptionAlg {
	return []EncryptionAlg{
		EncryptionAlgA128CBC,
		EncryptionAlgA128CBCPAD,
		EncryptionAlgA128GCM,
		EncryptionAlgA128KW,
		EncryptionAlgA192CBC,
		EncryptionAlgA192CBCPAD,
		EncryptionAlgA192GCM,
		EncryptionAlgA192KW,
		EncryptionAlgA256CBC,
		EncryptionAlgA256CBCPAD,
		EncryptionAlgA256GCM,
		EncryptionAlgA256KW,
		EncryptionAlgRSA15,
		EncryptionAlgRSAOAEP,
		EncryptionAlgRSAOAEP256,
	}
}

// SignatureAlg - The signing/verification algorithm identifier.
type SignatureAlg string

const (
	// ES256 - ECDSA using P-256 and SHA-256, as described in https://tools.ietf.org/html/rfc7518.
	SignatureAlgES256 SignatureAlg = "ES256"
	// ES256K - ECDSA using P-256K and SHA-256, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgES256K SignatureAlg = "ES256K"
	// ES384 - ECDSA using P-384 and SHA-384, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgES384 SignatureAlg = "ES384"
	// ES512 - ECDSA using P-521 and SHA-512, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgES512 SignatureAlg = "ES512"
	// PS256 - RSASSA-PSS using SHA-256 and MGF1 with SHA-256, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgPS256 SignatureAlg = "PS256"
	// PS384 - RSASSA-PSS using SHA-384 and MGF1 with SHA-384, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgPS384 SignatureAlg = "PS384"
	// PS512 - RSASSA-PSS using SHA-512 and MGF1 with SHA-512, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgPS512 SignatureAlg = "PS512"
	// RS256 - RSASSA-PKCS1-v1_5 using SHA-256, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgRS256 SignatureAlg = "RS256"
	// RS384 - RSASSA-PKCS1-v1_5 using SHA-384, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgRS384 SignatureAlg = "RS384"
	// RS512 - RSASSA-PKCS1-v1_5 using SHA-512, as described in https://tools.ietf.org/html/rfc7518
	SignatureAlgRS512 SignatureAlg = "RS512"
	// RSNULL - Reserved
	SignatureAlgRSNULL SignatureAlg = "RSNULL"
)

// PossibleSignatureAlgValues returns a slice of all possible SignatureAlg values
func PossibleSignatureAlgValues() []SignatureAlg {
	return []SignatureAlg{
		SignatureAlgES256,
		SignatureAlgES256K,
		SignatureAlgES384,
		SignatureAlgES512,
		SignatureAlgPS256,
		SignatureAlgPS384,
		SignatureAlgPS512,
		SignatureAlgRS256,
		SignatureAlgRS384,
		SignatureAlgRS512,
		SignatureAlgRSNULL,
	}
}
