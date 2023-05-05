//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"encoding/base64"
	"fmt"
	"sync/atomic"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential struct {
	// Only the NewSharedKeyCredential method should set these; all other methods should treat them as read-only
	accountName string
	accountKey  atomic.Value // []byte
}

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName string, accountKey string) (*SharedKeyCredential, error) {
	c := SharedKeyCredential{accountName: accountName}
	if err := c.SetAccountKey(accountKey); err != nil {
		return nil, err
	}
	return &c, nil
}

// AccountName returns the Storage account's name.
func (c *SharedKeyCredential) AccountName() string {
	return c.accountName
}

// SetAccountKey replaces the existing account key with the specified account key.
func (c *SharedKeyCredential) SetAccountKey(accountKey string) error {
	_bytes, err := base64.StdEncoding.DecodeString(accountKey)
	if err != nil {
		return fmt.Errorf("decode account key: %w", err)
	}
	c.accountKey.Store(_bytes)
	return nil
}
