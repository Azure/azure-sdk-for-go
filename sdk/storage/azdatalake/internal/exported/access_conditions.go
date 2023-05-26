//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

const SnapshotTimeFormat = "2006-01-02T15:04:05.0000000Z07:00"

// FilesystemAccessConditions identifies container-specific access conditions which you optionally set.
type FilesystemAccessConditions = container.AccessConditions

// PathAccessConditions identifies blob-specific access conditions which you optionally set.
type PathAccessConditions = blob.AccessConditions

// FormatContainerAccessConditions formats FilesystemAccessConditions into container's LeaseAccessConditions and ModifiedAccessConditions.
func FormatContainerAccessConditions(b *FilesystemAccessConditions) (*LeaseAccessConditions, *ModifiedAccessConditions) {
	if b == nil {
		return nil, nil
	}
	return b.LeaseAccessConditions, b.ModifiedAccessConditions
}

// FormatPathAccessConditions formats PathAccessConditions into path's LeaseAccessConditions and ModifiedAccessConditions.
func FormatPathAccessConditions(p *PathAccessConditions) (*generated.LeaseAccessConditions, *generated.ModifiedAccessConditions) {
	if p == nil {
		return nil, nil
	}
	return &generated.LeaseAccessConditions{
			LeaseID: p.LeaseAccessConditions.LeaseID,
		}, &generated.ModifiedAccessConditions{
			IfMatch:           p.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       p.ModifiedAccessConditions.IfNoneMatch,
			IfModifiedSince:   p.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: p.ModifiedAccessConditions.IfUnmodifiedSince,
		}
}

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = blob.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = blob.ModifiedAccessConditions
