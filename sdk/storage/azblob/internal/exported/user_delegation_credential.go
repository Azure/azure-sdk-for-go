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

// NewUserDelegationCredential creates a new UserDelegationCredential using a Storage account's Name and a user delegation Key from it
func NewUserDelegationCredential(accountName string, key generated.UserDelegationKey) *UserDelegationCredential {
	return &UserDelegationCredential{
		Name: accountName,
		Key:  key,
	}
}

type UserDelegationCredential struct {
	Name string
	Key  generated.UserDelegationKey
}

// AccountName returns the Storage account's Name
func (f *UserDelegationCredential) AccountName() string {
	return f.Name
}

// ComputeHMAC
func (f *UserDelegationCredential) ComputeHMACSHA256(message string) (string, error) {
	bytes, _ := base64.StdEncoding.DecodeString(*f.Key.Value)
	h := hmac.New(sha256.New, bytes)
	_, err := h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), err
}

// Private method to return important parameters for NewSASQueryParameters
func (f *UserDelegationCredential) getUDKParams() *generated.UserDelegationKey {
	return &f.Key
}

type UserDelegationKey = generated.UserDelegationKey
