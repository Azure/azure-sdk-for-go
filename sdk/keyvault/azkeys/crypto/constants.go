//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto

// Key wrapping algorithms
type KeyWrapAlgorithm string

// KeyWrapAlgorithm provides access to the predefined values of KeyWrapAlgorithm.
const KeyWrapAlgorithms = KeyWrapAlgorithm("")

func (KeyWrapAlgorithm) AES128() KeyWrapAlgorithm {
	return "A128KW"
}

func (KeyWrapAlgorithm) AES192() KeyWrapAlgorithm {
	return "A192KW"
}

func (KeyWrapAlgorithm) AES256() KeyWrapAlgorithm {
	return "A256KW"
}

func (KeyWrapAlgorithm) RSAOAEP() KeyWrapAlgorithm {
	return "RSA-OAEP"
}

func (KeyWrapAlgorithm) RSAOAEP256() KeyWrapAlgorithm {
	return "RSA-OAEP-256"
}

func (KeyWrapAlgorithm) RSA15() KeyWrapAlgorithm {
	return "RSA1_5"
}

// Cast converts the specified string to a KeyWrapAlgorithm value.
func (KeyWrapAlgorithm) Cast(s string) KeyWrapAlgorithm {
	return KeyWrapAlgorithm(s)
}

// Returns a pointer to a KeyWrapAlgorithm constant
func (k KeyWrapAlgorithm) ToPtr() *KeyWrapAlgorithm {
	return &k
}

// EncryptionAlgorithm - algorithm identifier
type EncryptionAlgorithm string

// EncryptionAlgorithm provides access to the predefined values of EncryptionAlgorithm.
const EncryptionAlgorithms = EncryptionAlgorithm("")

func (EncryptionAlgorithm) A128CBC() EncryptionAlgorithm {
	return "A128CBC"
}

func (EncryptionAlgorithm) A128CBCPAD() EncryptionAlgorithm {
	return "A128CBCPAD"
}

func (EncryptionAlgorithm) A128GCM() EncryptionAlgorithm {
	return "A128GCM"
}

func (EncryptionAlgorithm) A128KW() EncryptionAlgorithm {
	return "A128KW"
}

func (EncryptionAlgorithm) A192CBC() EncryptionAlgorithm {
	return "A192CBC"
}

func (EncryptionAlgorithm) A192CBCPAD() EncryptionAlgorithm {
	return "A192CBCPAD"
}

func (EncryptionAlgorithm) A192GCM() EncryptionAlgorithm {
	return "A192GCM"
}

func (EncryptionAlgorithm) A192KW() EncryptionAlgorithm {
	return "A192KW"
}

func (EncryptionAlgorithm) A256CBC() EncryptionAlgorithm {
	return "A256CBC"
}

func (EncryptionAlgorithm) A256CBCPAD() EncryptionAlgorithm {
	return "A256CBCPAD"
}

func (EncryptionAlgorithm) A256GCM() EncryptionAlgorithm {
	return "A256GCM"
}

func (EncryptionAlgorithm) A256KW() EncryptionAlgorithm {
	return "A256KW"
}

func (EncryptionAlgorithm) RSA15() EncryptionAlgorithm {
	return "RSA1_5"
}

func (EncryptionAlgorithm) RSAOAEP() EncryptionAlgorithm {
	return "RSA-OAEP"
}

func (EncryptionAlgorithm) RSAOAEP256() EncryptionAlgorithm {
	return "RSA-OAEP-256"
}

// Cast converts the specified string to a EncryptionAlgorithm value.
func (EncryptionAlgorithm) Cast(s string) EncryptionAlgorithm {
	return EncryptionAlgorithm(s)
}

// ToPtr returns a *EncryptionAlgorithm pointing to the current value.
func (c EncryptionAlgorithm) ToPtr() *EncryptionAlgorithm {
	return &c
}

// SignatureAlgorithm - The signing/verification algorithm identifier.
type SignatureAlgorithm string

// SignatureAlgorithms provides access to the predefined values of SignatureAlgorithm.
const SignatureAlgorithms = SignatureAlgorithm("")

// ES256 - ECDSA using P-256 and SHA-256, as described in https://tools.ietf.org/html/rfc7518.
func (SignatureAlgorithm) ES256() SignatureAlgorithm {
	return "ES256"
}

// ES256K - ECDSA using P-256K and SHA-256, as described in https://tools.ietf.org/html/rfc7518
func (SignatureAlgorithm) ES256K() SignatureAlgorithm {
	return "ES256K"
}

// ES384 - ECDSA using P-384 and SHA-384, as described in https://tools.ietf.org/html/rfc7518
func (SignatureAlgorithm) ES384() SignatureAlgorithm {
	return "ES384"
}

// ES512 - ECDSA using P-521 and SHA-512, as described in https://tools.ietf.org/html/rfc7518
func (SignatureAlgorithm) ES512() SignatureAlgorithm {
	return "ES512"
}

// PS256 - RSASSA-PSS using SHA-256 and MGF1 with SHA-256, as described in https://tools.ietf.org/html/rfc7518
func (SignatureAlgorithm) PS256() SignatureAlgorithm {
	return "PS256"
}

// PS384 - RSASSA-PSS using SHA-384 and MGF1 with SHA-384, as described in https://tools.ietf.org/html/rfc7518
func (SignatureAlgorithm) PS384() SignatureAlgorithm {
	return "PS384"
}

// PS512 - RSASSA-PSS using SHA-512 and MGF1 with SHA-512, as described in https://tools.ietf.org/html/rfc7518
func (SignatureAlgorithm) PS512() SignatureAlgorithm {
	return "PS512"
}

// RS256 - RSASSA-PKCS1-v1_5 using SHA-256, as described in https://tools.ietf.org/html/rfc7518
func (SignatureAlgorithm) RS256() SignatureAlgorithm {
	return "RS256"
}

// RS384 - RSASSA-PKCS1-v1_5 using SHA-384, as described in https://tools.ietf.org/html/rfc7518
func (SignatureAlgorithm) RS384() SignatureAlgorithm {
	return "RS384"
}

// RS512 - RSASSA-PKCS1-v1_5 using SHA-512, as described in https://tools.ietf.org/html/rfc7518
func (SignatureAlgorithm) RS512() SignatureAlgorithm {
	return "RS512"
}

// RSNULL - Reserved
func (SignatureAlgorithm) RSNULL() SignatureAlgorithm {
	return "RSNULL"
}

// Cast converts the specified string to a SignatureAlgorithm value.
func (SignatureAlgorithm) Cast(s string) SignatureAlgorithm {
	return SignatureAlgorithm(s)
}

// ToPtr returns a *SignatureAlgorithm pointing to the current value.
func (c SignatureAlgorithm) ToPtr() *SignatureAlgorithm {
	return &c
}
