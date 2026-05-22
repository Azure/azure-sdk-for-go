// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cache

import (
	"os"

	"github.com/AzureAD/microsoft-authentication-extensions-for-go/cache/accessor"
)

var (
	cacheDir = os.UserHomeDir
	storage  = func(name string) (accessor.Accessor, error) {
		return accessor.New(name, accessor.WithAccount("MSALCache"))
	}
)
