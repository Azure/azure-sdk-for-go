// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cache

import (
	"github.com/AzureAD/microsoft-authentication-extensions-for-go/cache/accessor"
	"golang.org/x/sys/windows"
)

var (
	cacheDir = func() (string, error) {
		return windows.KnownFolderPath(windows.FOLDERID_LocalAppData, 0)
	}
	storage = func(name string) (accessor.Accessor, error) {
		p, err := cacheFilePath(name)
		if err != nil {
			return nil, err
		}
		return accessor.New(p)
	}
)
