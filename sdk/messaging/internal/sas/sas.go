// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package sas provides SAS token functionality which implements TokenProvider from package auth for use with Azure
// Event Hubs and Service Bus.
package sas

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/internal/auth"
)

type (
	// Signer provides SAS token generation for use in Service Bus and Event Hub
	Signer struct {
		KeyName string
		Key     string
	}

	// TokenProvider is a SAS claims-based security token provider
	TokenProvider struct {
		signer *Signer
	}

	// TokenProviderOption provides configuration options for SAS Token Providers
	TokenProviderOption func(*TokenProvider) error
)

// TokenProviderWithKey configures a SAS TokenProvider to use the given key name and key (secret) for signing
func TokenProviderWithKey(keyName, key string) TokenProviderOption {
	return func(provider *TokenProvider) error {
		provider.signer = NewSigner(keyName, key)
		return nil
	}
}

// NewTokenProvider builds a SAS claims-based security token provider
func NewTokenProvider(opts ...TokenProviderOption) (*TokenProvider, error) {
	provider := new(TokenProvider)
	for _, opt := range opts {
		err := opt(provider)
		if err != nil {
			return nil, err
		}
	}
	return provider, nil
}

// GetToken gets a CBS SAS token
func (t *TokenProvider) GetToken(audience string) (*auth.Token, error) {
	signature, expiry, err := t.signer.SignWithDuration(audience, 2*time.Hour)

	if err != nil {
		return nil, err
	}

	return auth.NewToken(auth.CBSTokenTypeSAS, signature, expiry), nil
}

// NewSigner builds a new SAS signer for use in generation Service Bus and Event Hub SAS tokens
func NewSigner(keyName, key string) *Signer {
	return &Signer{
		KeyName: keyName,
		Key:     key,
	}
}

// SignWithDuration signs a given for a period of time from now
func (s *Signer) SignWithDuration(uri string, interval time.Duration) (signature, expiry string, err error) {
	expiry = signatureExpiry(time.Now().UTC(), interval)
	sig, err := s.SignWithExpiry(uri, expiry)

	if err != nil {
		return "", "", err
	}

	return sig, expiry, nil
}

// SignWithExpiry signs a given uri with a given expiry string
func (s *Signer) SignWithExpiry(uri, expiry string) (string, error) {
	audience := strings.ToLower(url.QueryEscape(uri))
	sts := stringToSign(audience, expiry)
	sig, err := s.signString(sts)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("SharedAccessSignature sr=%s&sig=%s&se=%s&skn=%s", audience, sig, expiry, s.KeyName), nil
}

func signatureExpiry(from time.Time, interval time.Duration) string {
	t := from.Add(interval).Round(time.Second).Unix()
	return strconv.FormatInt(t, 10)
}

func stringToSign(uri, expiry string) string {
	return uri + "\n" + expiry
}

func (s *Signer) signString(str string) (string, error) {
	h := hmac.New(sha256.New, []byte(s.Key))
	_, err := h.Write([]byte(str))

	if err != nil {
		return "", err
	}

	encodedSig := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return url.QueryEscape(encodedSig), nil
}
