// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

const SnapshotTimeFormat = "2006-01-02T15:04:05.0000000Z07:00"

// ContainerAccessConditions identifies container-specific access conditions which you optionally set.
type ContainerAccessConditions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
	LeaseAccessConditions    *LeaseAccessConditions
}

func FormatContainerAccessConditions(b *ContainerAccessConditions) (*LeaseAccessConditions, *ModifiedAccessConditions) {
	if b == nil {
		return nil, nil
	}
	return b.LeaseAccessConditions, b.ModifiedAccessConditions
}

// BlobAccessConditions identifies blob-specific access conditions which you optionally set.
type BlobAccessConditions struct {
	LeaseAccessConditions    *LeaseAccessConditions
	ModifiedAccessConditions *ModifiedAccessConditions
}

func FormatBlobAccessConditions(b *BlobAccessConditions) (*LeaseAccessConditions, *ModifiedAccessConditions) {
	if b == nil {
		return nil, nil
	}
	return b.LeaseAccessConditions, b.ModifiedAccessConditions
}

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = generated.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = generated.ModifiedAccessConditions

type BlobModifiedAccessConditions = generated.BlobModifiedAccessConditions

func FormatBlobModifiedAccessConditions(src *BlobModifiedAccessConditions, mac *ModifiedAccessConditions) *BlobModifiedAccessConditions {
	if src == nil && mac == nil {
		return nil
	}

	result := &BlobModifiedAccessConditions{}

	// Start with values from ModifiedAccessConditions as fallback
	if mac != nil {
		result.IfMatch = mac.IfMatch
		result.IfNoneMatch = mac.IfNoneMatch
		result.IfModifiedSince = mac.IfModifiedSince
		result.IfUnmodifiedSince = mac.IfUnmodifiedSince
	}

	// Override with values from BlobModifiedAccessConditions (takes precedence)
	if src != nil {
		if src.IfMatch != nil {
			result.IfMatch = src.IfMatch
		}
		if src.IfNoneMatch != nil {
			result.IfNoneMatch = src.IfNoneMatch
		}
		if src.IfModifiedSince != nil {
			result.IfModifiedSince = src.IfModifiedSince
		}
		if src.IfUnmodifiedSince != nil {
			result.IfUnmodifiedSince = src.IfUnmodifiedSince
		}
	}

	return result
}
