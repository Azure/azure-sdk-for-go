// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package jwe

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache/internal/aescbc"
	"github.com/stretchr/testify/require"
)

func TestEncryptParseDecrypt(t *testing.T) {
	plaintext := []byte("plaintext")
	kid := "42"
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)
	alg, err := aescbc.NewAES128CBCHMACSHA256(key)
	require.NoError(t, err)

	j, err := Encrypt(plaintext, kid, alg)
	require.NoError(t, err)

	s, err := j.Serialize()
	require.NoError(t, err)
	segments := strings.Split(s, ".")
	require.Len(t, segments, 5, "compact format has 5 segments")

	p, err := ParseCompactFormat([]byte(s))
	require.NoError(t, err)
	require.Equal(t, j, p)

	h, err := base64.RawURLEncoding.DecodeString(segments[0])
	require.NoError(t, err, segments[0])
	hdr := Header{}
	require.NoError(t, json.Unmarshal(h, &hdr))
	require.Equal(t, alg.Alg, hdr.Enc)
	require.Equal(t, "dir", hdr.Alg)
	require.Equal(t, kid, hdr.KID)

	require.Empty(t, segments[1])

	iv, err := base64.RawURLEncoding.DecodeString(segments[2])
	require.NoError(t, err)
	require.Len(t, iv, 16)

	ciphertext, err := base64.RawURLEncoding.DecodeString(segments[3])
	require.NoError(t, err)
	require.Len(t, ciphertext, 16)

	tag, err := base64.RawURLEncoding.DecodeString(segments[4])
	require.NoError(t, err)
	require.Len(t, tag, 16)

	actual, err := j.Decrypt(key)
	require.NoError(t, err)
	require.Equal(t, actual, plaintext)
}
