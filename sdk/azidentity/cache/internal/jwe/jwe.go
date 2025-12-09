// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package jwe

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache/internal/aescbc"
)

type Header struct {
	Alg string `json:"alg"`
	Enc string `json:"enc"`
	KID string `json:"kid"`
}

// JWE implements a subset of JSON Web Encryption (https://datatracker.ietf.org/doc/html/rfc7516).
// It supports only direct encryption (https://datatracker.ietf.org/doc/html/rfc7518#section-4.5)
// with A128CBC-HS256 and de/serializes only the compact format.
type JWE struct {
	Ciphertext, IV, Tag []byte
	Header              Header
}

func Encrypt(plaintext []byte, kid string, alg *aescbc.AESCBCHMACSHA2) (JWE, error) {
	iv := make([]byte, 16)
	_, err := rand.Read(iv)
	if err != nil {
		return JWE{}, err
	}
	result, err := alg.Encrypt(iv, plaintext, nil)
	if err != nil {
		return JWE{}, err
	}
	return JWE{
		Ciphertext: result.Ciphertext,
		Header: Header{
			Alg: "dir",
			Enc: alg.Alg,
			KID: kid,
		},
		IV:  iv,
		Tag: result.Tag,
	}, nil
}

// ParseCompactFormat deserializes the compact format as returned by [JWE.Serialize]
func ParseCompactFormat(b []byte) (JWE, error) {
	s := bytes.Split(b, []byte("."))
	if len(s) != 5 {
		return JWE{}, errors.New("incorrectly formatted JWE")
	}
	hdr, err := decode(s[0])
	if err != nil {
		return JWE{}, err
	}
	h := Header{}
	err = json.Unmarshal(hdr, &h)
	if err != nil {
		return JWE{}, err
	}
	iv, err := decode(s[2])
	if err != nil {
		return JWE{}, err
	}
	ciphertext, err := decode(s[3])
	if err != nil {
		return JWE{}, err
	}
	tag, err := decode(s[4])
	if err != nil {
		return JWE{}, err
	}
	return JWE{Header: h, IV: iv, Ciphertext: ciphertext, Tag: tag}, nil
}

func (j *JWE) Decrypt(key []byte) ([]byte, error) {
	if j.Header.Alg != "dir" {
		return nil, fmt.Errorf("unsupported content encryption algorithm %q", j.Header.Alg)
	}
	alg, err := aescbc.NewAES128CBCHMACSHA256(key)
	if err != nil {
		return nil, err
	}
	if j.Header.Enc != alg.Alg {
		return nil, fmt.Errorf("unsupported encryption algorithm %q", j.Header.Enc)
	}
	return alg.Decrypt(j.IV, j.Ciphertext, nil, j.Tag)
}

// Serialize the JWE to compact format
func (j *JWE) Serialize() (string, error) {
	hdr, err := json.Marshal(j.Header)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		// second segment (encrypted key) is empty because direct encryption doesn't wrap a key
		"%s..%s.%s.%s",
		base64.RawURLEncoding.EncodeToString(hdr),
		base64.RawURLEncoding.EncodeToString(j.IV),
		base64.RawURLEncoding.EncodeToString(j.Ciphertext),
		base64.RawURLEncoding.EncodeToString(j.Tag),
	), nil
}

func decode(b []byte) ([]byte, error) {
	dst := make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	n, err := base64.RawURLEncoding.Decode(dst, b)
	return dst[:n], err
}
