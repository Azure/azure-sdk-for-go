//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

// NewUserDelegationCredential creates a new UserDelegationCredential using a Storage account's name and a user delegation key from it
func NewUserDelegationCredential(accountName string, key generated.UserDelegationKey) *UserDelegationCredential {
	return &UserDelegationCredential{
		accountName: accountName,
		accountKey:  key,
	}
}

type UserDelegationCredential struct {
	accountName string
	accountKey  generated.UserDelegationKey
}

// AccountName returns the Storage account's name
func (f *UserDelegationCredential) AccountName() string {
	return f.accountName
}

// ComputeHMAC
func (f *UserDelegationCredential) ComputeHMACSHA256(message string) (string, error) {
	bytes, _ := base64.StdEncoding.DecodeString(*f.accountKey.Value)
	h := hmac.New(sha256.New, bytes)
	_, err := h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), err
}

// Private method to return important parameters for NewSASQueryParameters
func (f *UserDelegationCredential) getUDKParams() *generated.UserDelegationKey {
	return &f.accountKey
}

type UserDelegationKey = generated.UserDelegationKey
