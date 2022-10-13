//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package hashing

type StorageTransferValidationOption uint8

const (
	StorageTransferValidationOptionNone  StorageTransferValidationOption = 0
	StorageTransferValidationOptionCRC64 StorageTransferValidationOption = 1 << iota
)
