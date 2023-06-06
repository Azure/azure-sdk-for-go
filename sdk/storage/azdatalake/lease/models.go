//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
)

// FilesystemAcquireOptions contains the optional parameters for the LeaseClient.AcquireLease method.
type FilesystemAcquireOptions struct {
	ModifiedAccessConditions *azdatalake.ModifiedAccessConditions
}

func (o *FilesystemAcquireOptions) format() *lease.ContainerAcquireOptions {
	return &lease.ContainerAcquireOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

// FilesystemBreakOptions contains the optional parameters for the LeaseClient.BreakLease method.
type FilesystemBreakOptions struct {
	BreakPeriod              *int32
	ModifiedAccessConditions *azdatalake.ModifiedAccessConditions
}

func (o *FilesystemBreakOptions) format() *lease.ContainerBreakOptions {
	return &lease.ContainerBreakOptions{
		BreakPeriod: o.BreakPeriod,
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

// FilesystemChangeOptions contains the optional parameters for the LeaseClient.ChangeLease method.
type FilesystemChangeOptions struct {
	ModifiedAccessConditions *azdatalake.ModifiedAccessConditions
}

func (o *FilesystemChangeOptions) format() *lease.ContainerChangeOptions {
	return &lease.ContainerChangeOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

type FilesystemReleaseOptions struct {
	ModifiedAccessConditions *azdatalake.ModifiedAccessConditions
}

func (o *FilesystemReleaseOptions) format() *lease.ContainerReleaseOptions {
	return &lease.ContainerReleaseOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

type FilesystemRenewOptions struct {
	ModifiedAccessConditions *azdatalake.ModifiedAccessConditions
}

func (o *FilesystemRenewOptions) format() *lease.ContainerRenewOptions {
	return &lease.ContainerRenewOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

// PathAcquireOptions contains the optional parameters for the LeaseClient.AcquireLease method.
type PathAcquireOptions struct {
	ModifiedAccessConditions *azdatalake.ModifiedAccessConditions
}

func (o *PathAcquireOptions) format() *lease.BlobAcquireOptions {
	return &lease.BlobAcquireOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

// PathBreakOptions contains the optional parameters for the LeaseClient.BreakLease method.
type PathBreakOptions struct {
	BreakPeriod              *int32
	ModifiedAccessConditions *azdatalake.ModifiedAccessConditions
}

func (o *PathBreakOptions) format() *lease.BlobBreakOptions {
	return &lease.BlobBreakOptions{
		BreakPeriod: o.BreakPeriod,
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

// PathChangeOptions contains the optional parameters for the LeaseClient.ChangeLease method.
type PathChangeOptions struct {
	ModifiedAccessConditions *azdatalake.ModifiedAccessConditions
}

func (o *PathChangeOptions) format() *lease.BlobChangeOptions {
	return &lease.BlobChangeOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

type PathReleaseOptions struct {
	ModifiedAccessConditions *azdatalake.ModifiedAccessConditions
}

func (o *PathReleaseOptions) format() *lease.BlobReleaseOptions {
	return &lease.BlobReleaseOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

type PathRenewOptions struct {
	ModifiedAccessConditions *azdatalake.ModifiedAccessConditions
}

func (o *PathRenewOptions) format() *lease.BlobRenewOptions {
	return &lease.BlobRenewOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}
