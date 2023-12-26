//go:build go1.18 && darwin
// +build go1.18,darwin

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cache

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/AzureAD/microsoft-authentication-extensions-for-go/cache/accessor"
	"github.com/AzureAD/microsoft-authentication-extensions-for-go/cache/accessor/file"
)

var cacheDir = os.UserHomeDir

func storage(o internal.TokenCachePersistenceOptions) (accessor.Accessor, error) {
	name := o.Name
	if name == "" {
		name = defaultName
	}
	var a accessor.Accessor
	err := tryAccessor()
	if err != nil {
		msg := fmt.Sprintf("cache encryption is impossible because the keychain isn't usable: %s", err)
		if o.AllowUnencryptedStorage {
			f := ""
			f, err = internal.CacheFilePath(name)
			if err == nil {
				log.Write(azidentity.EventAuthentication, msg+". Falling back to unencrypted storage")
				a, err = file.New(f)
			}
		} else {
			err = errors.New(msg + ". Set AllowUnencryptedStorage on TokenCachePersistenceOptions to store the cache in plaintext instead of returning this error")
		}
	} else {
		a, err = accessor.New(name, accessor.WithAccount("MSALCache"))
	}
	return a, err
}

func tryAccessor() error {
	a, err := accessor.New("azidentity-test-cache")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = a.Write(ctx, []byte("test"))
	if err != nil {
		return err
	}
	_, err = a.Read(ctx)
	if err != nil {
		return err
	}
	return a.Delete(ctx)
}
