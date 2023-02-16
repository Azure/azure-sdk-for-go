//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Contains common helpers for TESTS ONLY
package testcommon

type TestAccountType string

const (
	TestAccountDefault    TestAccountType = ""
	TestAccountSecondary  TestAccountType = "SECONDARY_"
	TestAccountPremium    TestAccountType = "PREMIUM_"
	TestAccountSoftDelete TestAccountType = "SOFT_DELETE_"
)

const (
	DefaultEndpointSuffix       = "core.windows.net/"
	DefaultFileEndpointSuffix   = "file.core.windows.net/"
	AccountNameEnvVar           = "AZURE_STORAGE_ACCOUNT_NAME"
	AccountKeyEnvVar            = "AZURE_STORAGE_ACCOUNT_KEY"
	DefaultEndpointSuffixEnvVar = "AZURE_STORAGE_ENDPOINT_SUFFIX"
)

const (
	FakeStorageAccount = "fakestorage"
	FakeStorageURL     = "https://fakestorage.blob.core.windows.net"
)
