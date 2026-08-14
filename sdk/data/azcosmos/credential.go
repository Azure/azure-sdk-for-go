// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// KeyCredential is an account key used to authenticate with Cosmos DB.
//
// Prefer Microsoft Entra ID authentication with [NewClient] where possible; account keys grant
// full access to the account and cannot be scoped down.
//
// This is a package type rather than [azcore.KeyCredential] because the key has to be handed to
// the Cosmos driver, and azcore's type has no accessor outside azcore itself. The storage packages
// define their own for the same reason.
type KeyCredential struct {
	accountKey string
}

// NewKeyCredential creates a credential from a Cosmos DB account key.
//
// The key is validated here rather than on first use. The driver takes it as an opaque string and
// decodes it once per signature, so a malformed key would otherwise surface as an authentication
// failure on every operation instead of as a configuration error at startup.
//
// Validation matches what the driver accepts, which is not quite what [base64.StdEncoding] does:
// padding is optional there, and whitespace is rejected rather than skipped. A key that survives
// this constructor is one the driver can sign with.
func NewKeyCredential(accountKey string) (KeyCredential, error) {
	if accountKey == "" {
		return KeyCredential{}, errors.New("azcosmos: account key must not be empty")
	}
	if strings.ContainsAny(accountKey, " \t\r\n") {
		// Go's decoder skips newlines, so a key wrapped by a config file would decode here and
		// then fail on every request. The driver rejects them, so reject them too.
		return KeyCredential{}, errors.New("azcosmos: account key must not contain whitespace")
	}
	// The driver's decoder treats padding as optional, so accept both forms rather than rejecting
	// a key it would have signed with.
	encoding := base64.StdEncoding
	if len(accountKey)%4 != 0 {
		encoding = base64.RawStdEncoding
	}
	if _, err := encoding.DecodeString(accountKey); err != nil {
		return KeyCredential{}, fmt.Errorf("azcosmos: decoding account key: %w", err)
	}
	return KeyCredential{accountKey: accountKey}, nil
}

// key returns the account key for the driver. It is unexported so the key does not become part of
// the public surface.
//
//nolint:unused // consumed once the driver binding lands.
func (k KeyCredential) key() string {
	return k.accountKey
}
