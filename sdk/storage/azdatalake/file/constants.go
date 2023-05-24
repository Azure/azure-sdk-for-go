//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"

type ExpiryOptions = generated.PathExpiryOptions

const (
	Absolute           ExpiryOptions = generated.PathExpiryOptionsAbsolute
	NeverExpire        ExpiryOptions = generated.PathExpiryOptionsNeverExpire
	RelativeToCreation ExpiryOptions = generated.PathExpiryOptionsRelativeToCreation
	RelativeToNow      ExpiryOptions = generated.PathExpiryOptionsRelativeToNow
)
