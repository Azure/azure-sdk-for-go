//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/path"
)

type RenameMode = path.RenameMode

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
