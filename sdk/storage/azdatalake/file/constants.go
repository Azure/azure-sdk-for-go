//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

type ResourceType = generated.PathResourceType

// TODO: consider the possibility of not exposing this and just pass it under the hood
const (
	ResourceTypeFile      ResourceType = generated.PathResourceTypeFile
	ResourceTypeDirectory ResourceType = generated.PathResourceTypeDirectory
)

type RenameMode = generated.PathRenameMode

// TODO: consider the possibility of not exposing this and just pass it under the hood
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

type EncryptionAlgorithmType = blob.EncryptionAlgorithmType

const (
	EncryptionAlgorithmTypeNone   EncryptionAlgorithmType = blob.EncryptionAlgorithmTypeNone
	EncryptionAlgorithmTypeAES256 EncryptionAlgorithmType = blob.EncryptionAlgorithmTypeAES256
)
