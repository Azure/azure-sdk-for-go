//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/path"
)

type ExpiryOptions = generated.PathExpiryOptions

const (
	Absolute           ExpiryOptions = generated.PathExpiryOptionsAbsolute
	NeverExpire        ExpiryOptions = generated.PathExpiryOptionsNeverExpire
	RelativeToCreation ExpiryOptions = generated.PathExpiryOptionsRelativeToCreation
	RelativeToNow      ExpiryOptions = generated.PathExpiryOptionsRelativeToNow
)

type ResourceType = path.ResourceType

// TODO: consider the possibility of not exposing this and just pass it under the hood
const (
	ResourceTypeFile      ResourceType = path.ResourceTypeFile
	ResourceTypeDirectory ResourceType = path.ResourceTypeDirectory
)

type RenameMode = path.RenameMode

// TODO: consider the possibility of not exposing this and just pass it under the hood
const (
	RenameModeLegacy RenameMode = path.RenameModeLegacy
	RenameModePosix  RenameMode = path.RenameModePosix
)

type SetAccessControlRecursiveMode = path.SetAccessControlRecursiveMode

const (
	SetAccessControlRecursiveModeSet    SetAccessControlRecursiveMode = path.SetAccessControlRecursiveModeSet
	SetAccessControlRecursiveModeModify SetAccessControlRecursiveMode = path.SetAccessControlRecursiveModeModify
	SetAccessControlRecursiveModeRemove SetAccessControlRecursiveMode = path.SetAccessControlRecursiveModeRemove
)

type EncryptionAlgorithmType = path.EncryptionAlgorithmType

const (
	EncryptionAlgorithmTypeNone   EncryptionAlgorithmType = path.EncryptionAlgorithmTypeNone
	EncryptionAlgorithmTypeAES256 EncryptionAlgorithmType = path.EncryptionAlgorithmTypeAES256
)
