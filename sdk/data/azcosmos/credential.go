// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "errors"

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
func NewKeyCredential(accountKey string) (KeyCredential, error) {
	if accountKey == "" {
		return KeyCredential{}, errors.New("azcosmos: account key must not be empty")
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
