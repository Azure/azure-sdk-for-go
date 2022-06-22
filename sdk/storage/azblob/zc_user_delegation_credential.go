//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// NewUserDelegationCredential creates a new UserDelegationCredential using a Storage account's name and a user delegation key from it
func NewUserDelegationCredential(accountName string, key UserDelegationKey) *UserDelegationCredential {
	return &UserDelegationCredential{
		accountName: accountName,
		accountKey:  key,
	}
}

type UserDelegationCredential struct {
	accountName string
	accountKey  UserDelegationKey
}

// AccountName returns the Storage account's name
func (f *UserDelegationCredential) AccountName() string {
	return f.accountName
}

// ComputeHMAC
func (f *UserDelegationCredential) ComputeHMACSHA256(message string) (base64String string, err error) {
	bytes, _ := base64.StdEncoding.DecodeString(*f.accountKey.Value)
	h := hmac.New(sha256.New, bytes)
	_, err = h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), err
}

// Private method to return important parameters for NewSASQueryParameters
func (f *UserDelegationCredential) getUDKParams() *UserDelegationKey {
	return &f.accountKey
}
