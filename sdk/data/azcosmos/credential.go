// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import "errors"

// KeyCredential is an account key used to authenticate with Cosmos DB.
//
// Prefer Microsoft Entra ID authentication with [NewClient] where possible; account keys grant
// full access to the account and cannot be scoped down.
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
