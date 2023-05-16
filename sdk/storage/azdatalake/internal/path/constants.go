//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"

type ResourceType = generated.PathResourceType

const (
	ResourceTypeFile      ResourceType = generated.PathResourceTypeFile
	ResourceTypeDirectory ResourceType = generated.PathResourceTypeDirectory
)

type ExpiryOptions = generated.PathExpiryOptions

const (
	ExpiryOptionsAbsolute           ExpiryOptions = generated.PathExpiryOptionsAbsolute
	ExpiryOptionsNeverExpire        ExpiryOptions = generated.PathExpiryOptionsNeverExpire
	ExpiryOptionsRelativeToCreation ExpiryOptions = generated.PathExpiryOptionsRelativeToCreation
	ExpiryOptionsRelativeToNow      ExpiryOptions = generated.PathExpiryOptionsRelativeToNow
)

type RenameMode = generated.PathRenameMode

const (
	RenameModeLegacy RenameMode = generated.PathRenameModeLegacy
	RenameModePosix  RenameMode = generated.PathRenameModePosix
)

type SetAccessControlRecursiveMode = generated.PathSetAccessControlRecursiveMode

const (
	SetAccessControlRecursiveModeSet    SetAccessControlRecursiveMode = generated.PathSetAccessControlRecursiveModeSet
	SetAccessControlRecursiveModeModify SetAccessControlRecursiveMode = generated.PathSetAccessControlRecursiveModeModify
	SetAccessControlRecursiveModeRemove SetAccessControlRecursiveMode = generated.PathSetAccessControlRecursiveModeRemove
)
