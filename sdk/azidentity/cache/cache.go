//go:build darwin || linux || windows

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cache

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"path/filepath"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal"
	extcache "github.com/AzureAD/microsoft-authentication-extensions-for-go/cache"
	msal "github.com/AzureAD/microsoft-authentication-library-for-go/apps/cache"
)

var (
	// once ensures New tests the storage implementation only once
	once = &sync.Once{}
	// storageError is the error from the storage test
	storageError error
	// tryStorage tests the storage implementation by round-tripping data
	tryStorage = func() {
		const errFmt = "persistent storage isn't available due to error %q"
		// random content prevents conflict with concurrent processes executing this function
		n := fmt.Sprint(rand.Int())
		s, err := storage("azidentity-test" + n)
		if err != nil {
			storageError = fmt.Errorf(errFmt, err)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		in := []byte("test content")
		err = s.Write(ctx, in)
		if err != nil {
			storageError = fmt.Errorf(errFmt, err)
			return
		}
		out, err := s.Read(ctx)
		if err != nil {
			storageError = fmt.Errorf(errFmt, err)
			return
		}
		if !bytes.Equal(in, out) {
			storageError = fmt.Errorf(errFmt, "reading or writing cache data failed")
		}
		err = s.Delete(ctx)
		if err != nil {
			storageError = fmt.Errorf(errFmt, err)
		}
	}
)

// Options for persistent token caches.
type Options struct {
	// Name distinguishes caches. Set this to isolate data from other applications.
	Name string
}

// New constructs persistent token caches. See the [token caching guide] for details
// about the storage implementation.
//
// [token caching guide]: https://aka.ms/azsdk/go/identity/caching#Persistent-token-caching
func New(opts *Options) (azidentity.Cache, error) {
	once.Do(tryStorage)
	if storageError != nil {
		return azidentity.Cache{}, storageError
	}
	o := Options{}
	if opts != nil {
		o = *opts
	}
	if o.Name == "" {
		o.Name = "msal.cache"
	}
	factory := func(cae bool) (msal.ExportReplace, error) {
		name := o.Name
		if cae {
			name += ".cae"
		}
		p, err := cacheFilePath(name)
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

// cacheFilePath maps a cache name to a file path. This path is the base for a lockfile.
// Storage implementations may also use it directly to store cache data.
func cacheFilePath(name string) (string, error) {
	dir, err := cacheDir()
	if err != nil {
		return "", fmt.Errorf("couldn't create a cache file due to error %q", err)
	}
	return filepath.Join(dir, ".IdentityService", name), nil
}
