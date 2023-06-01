//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

// PublicAccessType defines values for AccessType - private (default) or file or filesystem.
type PublicAccessType = azblob.PublicAccessType

const (
	File       PublicAccessType = azblob.PublicAccessTypeBlob
	Filesystem PublicAccessType = azblob.PublicAccessTypeContainer
)
