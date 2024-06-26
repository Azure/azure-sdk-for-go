//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"errors"
)

var errMissingImport = errors.New("import github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache to enable persistent caching")

// CacheFilePath returns the path to the cache file for the given name.
// Defining it in this package makes it available to azidentity tests.
var CacheFilePath = func(name string) (string, error) {
	return "", errMissingImport
}
