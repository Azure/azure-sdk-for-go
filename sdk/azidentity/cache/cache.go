//go:build go1.18 && (darwin || linux || windows)
// +build go1.18
// +build darwin linux windows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cache

import (
	"fmt"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal"
	extcache "github.com/AzureAD/microsoft-authentication-extensions-for-go/cache"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
)

const defaultName = "msal.cache"

func init() {
	internal.NewCache = func(o *internal.TokenCachePersistenceOptions, enableCAE bool) (cache.ExportReplace, error) {
		if o == nil {
			return nil, nil
		}
		cp := *o
		if cp.Name == "" {
			cp.Name = defaultName
		}
		suffix := ".nocae"
		if enableCAE {
			suffix = ".cae"
		}
		cp.Name += suffix
		a, err := storage(cp)
		if err != nil {
			return nil, err
		}
		p, err := internal.CacheFilePath(cp.Name)
		if err != nil {
			return nil, err
		}
		return extcache.New(a, p)
	}
	internal.CacheFilePath = func(name string) (string, error) {
		dir, err := cacheDir()
		if err != nil {
			return "", fmt.Errorf("couldn't create a cache file due to error %q", err)
		}
		return filepath.Join(dir, ".IdentityService", name), nil
	}
}
