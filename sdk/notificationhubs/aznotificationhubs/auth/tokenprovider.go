//go:build go1.20
// +build go1.20

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	// AZNHTokenTypeJWT is the type of token to be used for JWTs. For example Azure Active Directory tokens.
	AZNHTokenTypeJWT TokenType = "jwt"
	// AZNHTokenTypeSAS is the type of token to be used for SAS tokens.
	AZNHTokenTypeSAS TokenType = "servicebus.windows.net:sastoken"

	AZNHEndpointKey            = "Endpoint"
	AZNHSharedAccessKeyNameKey = "SharedAccessKeyName"
	AZNHSharedAccessKeyKey     = "SharedAccessKey"
)

type (
	// TokenType represents types of tokens known for claims-based auth
	TokenType string

	// Token contains all of the information to negotiate authentication
	Token struct {
		// TokenType is the type of AZNH token
		TokenType TokenType
		Token     string
		Expiry    string
	}

	// TokenProvider abstracts the fetching of authentication tokens
	TokenProvider interface {
		GetToken(uri string) (*Token, error)
	}

	// The NotificationHubsTokenProvider is a TokenProvider that uses a shared access signature to authenticate with Azure Notification Hubs.
	NotificationHubsTokenProvider struct {
		keyName  string
		keyValue string
	}

	// ParsedConnection is a struct that contains the parsed connection string
	ParsedConnection struct {
		Endpoint string
		KeyName  string
		KeyValue string
	}
)

// NewToken constructs a new auth token
func NewToken(tokenType TokenType, token, expiry string) *Token {
	return &Token{
		TokenType: tokenType,
		Token:     token,
		Expiry:    expiry,
	}
}

// Createst a new NotificationHubsTokenProvider with the SAS key name and key value.
func NewNotificationHubsTokenProvider(keyName string, keyValue string) *NotificationHubsTokenProvider {
	return &NotificationHubsTokenProvider{
		keyName:  keyName,
		keyValue: keyValue,
	}
}

// GetToken returns a token for the given audience URI
func (t *NotificationHubsTokenProvider) GetToken(uri string) (*Token, error) {
	audience := strings.ToLower(uri)
	expiry := time.Now().UTC().Unix() + int64(3600)
	sts := createStringToSign(audience, expiry)
	sig := t.signString(sts)
	tokenParams := url.Values{
		"sr":  {audience},
		"sig": {sig},
		"se":  {fmt.Sprintf("%d", expiry)},
		"skn": {t.keyName},
	}

	return &Token{
		TokenType: AZNHTokenTypeSAS,
		Token:     fmt.Sprintf("SharedAccessSignature %s", tokenParams.Encode()),
		Expiry:    fmt.Sprintf("%d", expiry),
	}, nil
}

func createStringToSign(uri string, expiry int64) string {
	return fmt.Sprintf("%s\n%d", url.QueryEscape(uri), expiry)
}

func (t *NotificationHubsTokenProvider) signString(str string) string {
	h := hmac.New(sha256.New, []byte(t.keyValue))
	h.Write([]byte(str))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// FromConnectionString parses a connection string and returns a ParsedConnection
func FromConnectionString(connectionString string) (*ParsedConnection, error) {
	var endpoint, keyName, keyValue string
	splits := strings.Split(connectionString, ";")
	for _, split := range splits {
		keyValuePair := strings.SplitN(split, "=", 2)
		if len(keyValuePair) < 2 {
			return nil, errors.New("failed parsing connection string due to unmatched key value separated by '='")
		}

		key := keyValuePair[0]
		value := keyValuePair[1]
		switch {
		case strings.EqualFold(AZNHEndpointKey, key):
			endpoint = value
		case strings.EqualFold(AZNHSharedAccessKeyNameKey, key):
			keyName = value
		case strings.EqualFold(AZNHSharedAccessKeyKey, key):
			keyValue = value
		}
	}

	if endpoint == "" {
		return nil, fmt.Errorf("key %q must not be empty", AZNHEndpointKey)
	}

	if keyName == "" {
		return nil, fmt.Errorf("key %q must not be empty", AZNHSharedAccessKeyNameKey)
	}

	if keyValue == "" {
		return nil, fmt.Errorf("key %q must not be empty", AZNHSharedAccessKeyKey)
	}

	return &ParsedConnection{
		Endpoint: endpoint,
		KeyName:  keyName,
		KeyValue: keyValue,
	}, nil
}
