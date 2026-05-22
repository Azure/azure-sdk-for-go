// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aescbc

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"encoding/binary"
	"errors"
	"hash"
)

// AESCBCHMACSHA2 implements AES_CBC_HMAC_SHA2 as defined in https://tools.ietf.org/html/rfc7518#section-5.2.2
type AESCBCHMACSHA2 struct {
	Alg            string
	encKey, macKey []byte
	hasher         func() hash.Hash
	tLen           int
}

type EncryptResult struct {
	Ciphertext, Tag []byte
}

// NewAES128CBCHMACSHA256 returns an implementation of AES_128_CBC_HMAC_SHA_256
// (https://tools.ietf.org/html/rfc7518#section-5.2.3)
func NewAES128CBCHMACSHA256(key []byte) (*AESCBCHMACSHA2, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes")
	}
	cp := make([]byte, 32)
	copy(cp, key)
	return newAESCBCHMACSHA2("A128CBC-HS256", cp, crypto.SHA256.New)
}

func newAESCBCHMACSHA2(alg string, k []byte, hasher func() hash.Hash) (*AESCBCHMACSHA2, error) {
	return &AESCBCHMACSHA2{
		Alg:    alg,
		encKey: k[len(k)/2:],
		hasher: hasher,
		macKey: k[:len(k)/2],
		tLen:   len(k) / 2,
	}, nil
}

func (a *AESCBCHMACSHA2) Decrypt(iv, ciphertext, additionalData, tag []byte) ([]byte, error) {
	expected := a.tag(iv, ciphertext, additionalData)
	if !hmac.Equal(tag, expected) {
		return nil, errors.New("decryption failed")
	}
	block, err := aes.NewCipher(a.encKey)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(out, ciphertext)
	return unpad(out)
}

func (a *AESCBCHMACSHA2) Encrypt(iv, plaintext, additionalData []byte) (EncryptResult, error) {
	result := EncryptResult{}
	block, err := aes.NewCipher(a.encKey)
	if err != nil {
		return result, err
	}
	plaintext = pad(plaintext)
	result.Ciphertext = make([]byte, len(plaintext))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(result.Ciphertext, plaintext)
	result.Tag = a.tag(iv, result.Ciphertext, additionalData)
	return result, nil
}

func (a *AESCBCHMACSHA2) tag(iv, ciphertext, aad []byte) []byte {
	h := hmac.New(a.hasher, a.macKey)
	h.Write(aad)
	h.Write(iv)
	h.Write(ciphertext)
	// aadBits is AL from step 4 of https://datatracker.ietf.org/doc/html/rfc7518#section-5.2.2.1
	aadBits := make([]byte, 8)
	binary.BigEndian.PutUint64(aadBits, uint64(len(aad)*8))
	h.Write(aadBits)
	return h.Sum(nil)[:a.tLen]
}

// pad adds PKCS#7 padding (https://datatracker.ietf.org/doc/html/rfc5652#section-6.3)
func pad(b []byte) []byte {
	n := aes.BlockSize - (len(b) % aes.BlockSize)
	padding := bytes.Repeat([]byte{byte(n)}, n)
	return append(b, padding...)
}

// unpad checks and removes PKCS#7 padding
func unpad(b []byte) ([]byte, error) {
	l := len(b)
	if l == 0 {
		return nil, nil
	}
	n := int(b[l-1])
	if n < 1 || n > aes.BlockSize || l%aes.BlockSize != 0 {
		return nil, errors.New("decryption failed")
	}
	for i := l - n; i < len(b); i++ {
		if b[i] != byte(n) {
			return nil, errors.New("decryption failed")
		}
	}
	return b[:l-n], nil
}
