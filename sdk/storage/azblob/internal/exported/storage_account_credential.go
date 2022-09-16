//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

// StorageAccountCredential is a wrapper interface for SharedKeyCredential and UserDelegationCredential
type StorageAccountCredential interface {
	AccountName() string
	ComputeHMACSHA256(message string) (string, error)
	getUDKParams() *generated.UserDelegationKey
}
