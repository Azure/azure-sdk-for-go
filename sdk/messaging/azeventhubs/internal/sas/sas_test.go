// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package sas

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	sas = "SharedAccessSignature"
)

type (
	sig struct {
		sr  string
		se  string
		skn string
		sig string
	}
)

func TestNewSigner(t *testing.T) {
	keyName, key := "foo", "superSecret"
	signer := NewSigner(keyName, key)
	before := time.Now().UTC().Add(-2 * time.Second)

	// the URL is lowercased and escaped when used as the audience in our signature.
	sigStr, expiry, err := signer.SignWithDuration("http://MiCrosoft.com", 1*time.Hour)

	require.NoError(t, err)
	nixExpiry, err := strconv.ParseInt(expiry, 10, 64)
	require.NoError(t, err)
	assert.WithinDuration(t, before.Add(1*time.Hour), time.Unix(nixExpiry, 0), 10*time.Second, "signing expiry")

	sig, err := parseSig(sigStr)
	assert.Nil(t, err)
	assert.Equal(t, "http%3a%2f%2fmicrosoft.com", sig.sr)
	assert.Equal(t, keyName, sig.skn)
	assert.Equal(t, expiry, sig.se)
	assert.NotNil(t, sig.sig)
}

func TestTokenProviderWithSAS(t *testing.T) {
	tp, err := NewTokenProvider(TokenProviderWithSAS("hello"))
	require.NoError(t, err)

	token, err := tp.GetToken("audience")
	require.NoError(t, err)

	require.Equal(t, &auth.Token{
		TokenType: auth.CBSTokenTypeSAS,
		Expiry:    "0",
		Token:     "hello",
	}, token)
}

func TestTokenProviderWithKey(t *testing.T) {
	tp, err := NewTokenProvider(TokenProviderWithKey("keyName", "key", 3*24*time.Hour))
	require.NoError(t, err)

	now, err := time.Parse(time.RFC3339, "2020-01-01T01:02:03Z")
	require.NoError(t, err)

	// hardcodes a particular date so our test is consistent.
	tp.signer.getNow = func() time.Time {
		return now
	}

	token, err := tp.GetToken("audience")
	require.NoError(t, err)

	require.Equal(t, &auth.Token{
		TokenType: auth.CBSTokenTypeSAS,
		Expiry:    fmt.Sprintf("%d", now.UTC().Add(3*24*time.Hour).Unix()),
		// NOTE: this is just literally the signature, using the key "key". Nothing secret or interesting here.
		Token: "SharedAccessSignature sr=audience&sig=8UM0iIfFCfeBSqxSdBMW8pUbhAm7mnjSUaIZTZx8V0g%3D&se=1578099723&skn=keyName",
	}, token)
}

func parseSig(sigStr string) (*sig, error) {
	if !strings.HasPrefix(sigStr, sas+" ") {
		return nil, errors.New("should start with " + sas)
	}
	sigStr = strings.TrimPrefix(sigStr, sas+" ")
	parts := strings.Split(sigStr, "&")
	parsed := new(sig)
	for _, part := range parts {
		keyValue := strings.Split(part, "=")
		if len(keyValue) != 2 {
			return nil, errors.New("key value is malformed")
		}
		switch keyValue[0] {
		case "sr":
			parsed.sr = keyValue[1]
		case "se":
			parsed.se = keyValue[1]
		case "sig":
			parsed.sig = keyValue[1]
		case "skn":
			parsed.skn = keyValue[1]
		default:
			return nil, fmt.Errorf("unknown key / value: %q", keyValue)
		}
	}
	return parsed, nil
}
