// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

// StorageAccountCredential is a wrapper interface for SharedKeyCredential and UserDelegationCredential
type StorageAccountCredential interface {
	AccountName() string
	ComputeHMACSHA256(message string) (base64String string)
}
