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
type FilesystemAccessConditions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
	LeaseAccessConditions    *LeaseAccessConditions
}

// FormatContainerAccessConditions formats FilesystemAccessConditions into container's LeaseAccessConditions and ModifiedAccessConditions.
func FormatContainerAccessConditions(b *FilesystemAccessConditions) *container.AccessConditions {
	if b == nil {
		return nil
	}
	return &container.AccessConditions{
		LeaseAccessConditions: &container.LeaseAccessConditions{
			LeaseID: b.LeaseAccessConditions.LeaseID,
		},
		ModifiedAccessConditions: &container.ModifiedAccessConditions{
			IfMatch:           b.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       b.ModifiedAccessConditions.IfNoneMatch,
			IfModifiedSince:   b.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: b.ModifiedAccessConditions.IfUnmodifiedSince,
		},
	}
}

// PathAccessConditions identifies blob-specific access conditions which you optionally set.
type PathAccessConditions struct {
	LeaseAccessConditions    *LeaseAccessConditions
	ModifiedAccessConditions *ModifiedAccessConditions
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

// FormatBlobAccessConditions formats PathAccessConditions into blob's LeaseAccessConditions and ModifiedAccessConditions.
func FormatBlobAccessConditions(p *PathAccessConditions) *blob.AccessConditions {
	if p == nil {
		return nil
	}
	return &blob.AccessConditions{LeaseAccessConditions: &blob.LeaseAccessConditions{
		LeaseID: p.LeaseAccessConditions.LeaseID,
	}, ModifiedAccessConditions: &blob.ModifiedAccessConditions{
		IfMatch:           p.ModifiedAccessConditions.IfMatch,
		IfNoneMatch:       p.ModifiedAccessConditions.IfNoneMatch,
		IfModifiedSince:   p.ModifiedAccessConditions.IfModifiedSince,
		IfUnmodifiedSince: p.ModifiedAccessConditions.IfUnmodifiedSince,
	}}
}

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = generated.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = generated.ModifiedAccessConditions
