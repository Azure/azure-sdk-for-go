//go:build go1.18 && (darwin || linux || windows)
// +build go1.18
// +build darwin linux windows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cache

import (
	"fmt"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal"
	extcache "github.com/AzureAD/microsoft-authentication-extensions-for-go/cache"
	msal "github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
)

const defaultName = "msal.cache"

func init() {
	internal.CacheFilePath = func(name string) (string, error) {
		dir, err := cacheDir()
		if err != nil {
			return "", fmt.Errorf("couldn't create a cache file due to error %q", err)
		}
		return filepath.Join(dir, ".IdentityService", name), nil
	}
}

// Options for persistent token caches.
type Options struct {
	// Name distinguishes the cache from other caches.
	// Set this to isolate data from other applications.
	Name string
}

// New is the constructor for persistent token caches.
func New(opts *Options) (azidentity.Cache, error) {
	// TODO: try the storage implementation, return any error

	o := Options{}
	if opts != nil {
		o = *opts
	}
	if o.Name == "" {
		o.Name = defaultName
	}
	factory := func(cae bool) (msal.ExportReplace, error) {
		name := o.Name
		if cae {
			name += ".cae"
		}
		p, err := internal.CacheFilePath(name)
		if err != nil {
			return nil, err
		}
		s, err := storage(name)
		if err != nil {
			return nil, err
		}
		return extcache.New(s, p)
	}
	return internal.NewCache(factory), nil
}
