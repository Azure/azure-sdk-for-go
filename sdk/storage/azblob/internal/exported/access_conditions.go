//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

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

type LeaseAccessConditions = generated.LeaseAccessConditions

type ModifiedAccessConditions = generated.ModifiedAccessConditions
