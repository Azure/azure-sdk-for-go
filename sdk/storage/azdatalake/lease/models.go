// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
)

// FileSystemAcquireOptions contains the optional parameters for the LeaseClient.AcquireLease method.
type FileSystemAcquireOptions struct {
	// ModifiedAccessConditions contains optional parameters to access filesystem.
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *FileSystemAcquireOptions) format() *lease.ContainerAcquireOptions {
	if o == nil || o.ModifiedAccessConditions == nil {
		return nil
	}
	return &lease.ContainerAcquireOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

// FileSystemBreakOptions contains the optional parameters for the LeaseClient.BreakLease method.
type FileSystemBreakOptions struct {
	// BreakPeriod is the proposed duration of seconds that the lease should continue before it is broken.
	BreakPeriod *int32
	// ModifiedAccessConditions contains optional parameters to access filesystem.
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *FileSystemBreakOptions) format() *lease.ContainerBreakOptions {
	opts := &lease.ContainerBreakOptions{}
	if o == nil {
		return opts
	}
	if o.ModifiedAccessConditions == nil {
		opts.BreakPeriod = o.BreakPeriod
		return opts
	}
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

// FileSystemChangeOptions contains the optional parameters for the LeaseClient.ChangeLease method.
type FileSystemChangeOptions struct {
	// ModifiedAccessConditions contains optional parameters to access filesystem.
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *FileSystemChangeOptions) format() *lease.ContainerChangeOptions {
	if o == nil || o.ModifiedAccessConditions == nil {
		return nil
	}
	return &lease.ContainerChangeOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

// FileSystemReleaseOptions contains the optional parameters for the LeaseClient.ReleaseLease method.
type FileSystemReleaseOptions struct {
	// ModifiedAccessConditions contains optional parameters to access filesystem.
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *FileSystemReleaseOptions) format() *lease.ContainerReleaseOptions {
	if o == nil || o.ModifiedAccessConditions == nil {
		return nil
	}
	return &lease.ContainerReleaseOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

// FileSystemRenewOptions contains the optional parameters for the LeaseClient.RenewLease method.
type FileSystemRenewOptions struct {
	// ModifiedAccessConditions contains optional parameters to access filesystem.
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *FileSystemRenewOptions) format() *lease.ContainerRenewOptions {
	if o == nil || o.ModifiedAccessConditions == nil {
		return nil
	}
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
	// ModifiedAccessConditions contains optional parameters to access path.
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *PathAcquireOptions) format() *lease.BlobAcquireOptions {
	if o == nil || o.ModifiedAccessConditions == nil {
		return nil
	}
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
	// BreakPeriod is the proposed duration of seconds that the lease should continue before it is broken.
	BreakPeriod *int32
	// ModifiedAccessConditions contains optional parameters to access path.
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *PathBreakOptions) format() *lease.BlobBreakOptions {
	opts := &lease.BlobBreakOptions{}
	if o == nil {
		return opts
	}
	if o.ModifiedAccessConditions == nil {
		opts.BreakPeriod = o.BreakPeriod
		return opts
	}
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
	// ModifiedAccessConditions contains optional parameters to access path.
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *PathChangeOptions) format() *lease.BlobChangeOptions {
	if o == nil || o.ModifiedAccessConditions == nil {
		return nil
	}
	return &lease.BlobChangeOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

// PathReleaseOptions contains the optional parameters for the LeaseClient.ReleaseLease method.
type PathReleaseOptions struct {
	// ModifiedAccessConditions contains optional parameters to access path.
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *PathReleaseOptions) format() *lease.BlobReleaseOptions {
	if o == nil || o.ModifiedAccessConditions == nil {
		return nil
	}
	return &lease.BlobReleaseOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

// PathRenewOptions contains the optional parameters for the LeaseClient.RenewLease method.
type PathRenewOptions struct {
	// ModifiedAccessConditions contains optional parameters to access path.
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *PathRenewOptions) format() *lease.BlobRenewOptions {
	if o == nil || o.ModifiedAccessConditions == nil {
		return nil
	}
	return &lease.BlobRenewOptions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
			IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
			IfMatch:           o.ModifiedAccessConditions.IfMatch,
			IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		},
	}
}

// AccessConditions identifies blob-specific access conditions which you optionally set.
type AccessConditions = exported.AccessConditions

// LeaseAccessConditions contains optional parameters to access leased entity.
type LeaseAccessConditions = exported.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = exported.ModifiedAccessConditions
