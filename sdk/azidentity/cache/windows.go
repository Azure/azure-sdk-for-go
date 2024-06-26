//go:build go1.18 && windows
// +build go1.18,windows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cache

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal"
	"github.com/AzureAD/microsoft-authentication-extensions-for-go/cache/accessor"
	"golang.org/x/sys/windows"
)

var cacheDir = func() (string, error) {
	return windows.KnownFolderPath(windows.FOLDERID_LocalAppData, 0)
}

func storage(name string) (accessor.Accessor, error) {
	p, err := internal.CacheFilePath(name)
	if err != nil {
		return nil, err
	}
	return accessor.New(p)
}
