//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

// PublicAccessType defines values for AccessType - private (default) or blob or container.
type PublicAccessType = azblob.PublicAccessType

const (
	AccessTypeFile PublicAccessType = azblob.PublicAccessTypeBlob
	TypeFilesystem PublicAccessType = azblob.PublicAccessTypeContainer
)
