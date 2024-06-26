//go:build go1.18 && darwin
// +build go1.18,darwin

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cache

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/AzureAD/microsoft-authentication-extensions-for-go/cache/accessor"
)

var cacheDir = os.UserHomeDir

func storage(name string) (accessor.Accessor, error) {
	if err := tryAccessor(); err != nil {
		return nil, errors.New("cache encryption is impossible because the keychain isn't usable: " + err.Error())
	}
	return accessor.New(name, accessor.WithAccount("MSALCache"))
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
