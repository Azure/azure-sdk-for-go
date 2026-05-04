// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

// BreakNaturally tells ContainerClient's or BlobClient's BreakLease method to break the lease using service semantics.
const BreakNaturally = -1

// AccessConditions contains a group of parameters for specifying lease access conditions.
type AccessConditions = generated.LeaseAccessConditions

// ModifiedAccessConditions contains a group of parameters for specifying access conditions.
type ModifiedAccessConditions = exported.ModifiedAccessConditions

// BlobAcquireOptions contains the optional parameters for the LeaseClient.AcquireLease method.
type BlobAcquireOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BlobAcquireOptions) format(leaseID *string) *generated.BlobClientAcquireLeaseOptions {
	opts := &generated.BlobClientAcquireLeaseOptions{
		ProposedLeaseID: leaseID,
	}
	if o != nil && o.ModifiedAccessConditions != nil {
		opts.IfMatch = o.ModifiedAccessConditions.IfMatch
		opts.IfModifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfModifiedSince)
		opts.IfNoneMatch = o.ModifiedAccessConditions.IfNoneMatch
		opts.IfUnmodifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfUnmodifiedSince)
	}

	return opts
}

// BlobBreakOptions contains the optional parameters for the LeaseClient.BreakLease method.
type BlobBreakOptions struct {
	// For a break operation, proposed Duration the lease should continue before it is broken, in seconds, between 0 and 60. This
	// break period is only used if it is shorter than the time remaining on the lease. If longer, the time remaining on the lease
	// is used. A new lease will not be available before the break period has expired, but the lease may be held for longer than
	// the break period. If this header does not appear with a break operation, a fixed-Duration lease breaks after the remaining
	// lease period elapses, and an infinite lease breaks immediately.
	BreakPeriod              *int32
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BlobBreakOptions) format() *generated.BlobClientBreakLeaseOptions {
	if o == nil {
		return nil
	}

	var period *int32
	if o.BreakPeriod != nil {
		period = leasePeriodPointer(*o.BreakPeriod)
	}

	opts := &generated.BlobClientBreakLeaseOptions{
		BreakPeriod: period,
	}
	if o.ModifiedAccessConditions != nil {
		opts.IfMatch = o.ModifiedAccessConditions.IfMatch
		opts.IfModifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfModifiedSince)
		opts.IfNoneMatch = o.ModifiedAccessConditions.IfNoneMatch
		opts.IfUnmodifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfUnmodifiedSince)
	}

	return opts
}

// BlobChangeOptions contains the optional parameters for the LeaseClient.ChangeLease method.
type BlobChangeOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BlobChangeOptions) format() *generated.BlobClientChangeLeaseOptions {
	if o == nil {
		return nil
	}

	opts := &generated.BlobClientChangeLeaseOptions{}
	if o.ModifiedAccessConditions != nil {
		opts.IfMatch = o.ModifiedAccessConditions.IfMatch
		opts.IfModifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfModifiedSince)
		opts.IfNoneMatch = o.ModifiedAccessConditions.IfNoneMatch
		opts.IfUnmodifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfUnmodifiedSince)
	}

	return opts
}

// BlobRenewOptions contains the optional parameters for the LeaseClient.RenewLease method.
type BlobRenewOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BlobRenewOptions) format() *generated.BlobClientRenewLeaseOptions {
	if o == nil {
		return nil
	}

	opts := &generated.BlobClientRenewLeaseOptions{}
	if o.ModifiedAccessConditions != nil {
		opts.IfMatch = o.ModifiedAccessConditions.IfMatch
		opts.IfModifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfModifiedSince)
		opts.IfNoneMatch = o.ModifiedAccessConditions.IfNoneMatch
		opts.IfUnmodifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfUnmodifiedSince)
	}

	return opts
}

// BlobReleaseOptions contains the optional parameters for the LeaseClient.ReleaseLease method.
type BlobReleaseOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BlobReleaseOptions) format() *generated.BlobClientReleaseLeaseOptions {
	if o == nil {
		return nil
	}

	opts := &generated.BlobClientReleaseLeaseOptions{}
	if o.ModifiedAccessConditions != nil {
		opts.IfMatch = o.ModifiedAccessConditions.IfMatch
		opts.IfModifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfModifiedSince)
		opts.IfNoneMatch = o.ModifiedAccessConditions.IfNoneMatch
		opts.IfUnmodifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfUnmodifiedSince)
	}

	return opts
}

// ContainerAcquireOptions contains the optional parameters for the LeaseClient.AcquireLease method.
type ContainerAcquireOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ContainerAcquireOptions) format(leaseID *string) *generated.ContainerClientAcquireLeaseOptions {
	// Note: missing mapping for o.ModifiedAccessConditions.IfMatch, o.ModifiedAccessConditions.IfNoneMatch
	opts := &generated.ContainerClientAcquireLeaseOptions{
		ProposedLeaseID: leaseID,
	}
	if o != nil && o.ModifiedAccessConditions != nil {
		opts.IfModifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfModifiedSince)
		opts.IfUnmodifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfUnmodifiedSince)
	}

	return opts
}

// ContainerBreakOptions contains the optional parameters for the LeaseClient.BreakLease method.
type ContainerBreakOptions struct {
	// For a break operation, proposed Duration the lease should continue before it is broken, in seconds, between 0 and 60. This
	// break period is only used if it is shorter than the time remaining on the lease. If longer, the time remaining on the lease
	// is used. A new lease will not be available before the break period has expired, but the lease may be held for longer than
	// the break period. If this header does not appear with a break operation, a fixed-Duration lease breaks after the remaining
	// lease period elapses, and an infinite lease breaks immediately.
	BreakPeriod              *int32
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ContainerBreakOptions) format() *generated.ContainerClientBreakLeaseOptions {
	if o == nil {
		return nil
	}

	var period *int32
	if o.BreakPeriod != nil {
		period = leasePeriodPointer(*o.BreakPeriod)
	}

	// Note: missing mapping for o.ModifiedAccessConditions.IfMatch, o.ModifiedAccessConditions.IfNoneMatch
	opts := &generated.ContainerClientBreakLeaseOptions{
		BreakPeriod: period,
	}
	if o.ModifiedAccessConditions != nil {
		opts.IfModifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfModifiedSince)
		opts.IfUnmodifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfUnmodifiedSince)
	}

	return opts
}

// ContainerChangeOptions contains the optional parameters for the LeaseClient.ChangeLease method.
type ContainerChangeOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ContainerChangeOptions) format() *generated.ContainerClientChangeLeaseOptions {
	if o == nil {
		return nil
	}
	// Note: missing mapping for o.ModifiedAccessConditions.IfMatch, o.ModifiedAccessConditions.IfNoneMatch

	opts := &generated.ContainerClientChangeLeaseOptions{}
	if o.ModifiedAccessConditions != nil {
		opts.IfModifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfModifiedSince)
		opts.IfUnmodifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfUnmodifiedSince)
	}

	return opts
}

// ContainerRenewOptions contains the optional parameters for the LeaseClient.RenewLease method.
type ContainerRenewOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ContainerRenewOptions) format() *generated.ContainerClientRenewLeaseOptions {
	if o == nil {
		return nil
	}
	// Note: missing mapping for o.ModifiedAccessConditions.IfMatch, o.ModifiedAccessConditions.IfNoneMatch

	opts := &generated.ContainerClientRenewLeaseOptions{}
	if o.ModifiedAccessConditions != nil {
		opts.IfModifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfModifiedSince)
		opts.IfUnmodifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfUnmodifiedSince)
	}

	return opts
}

// ContainerReleaseOptions contains the optional parameters for the LeaseClient.ReleaseLease method.
type ContainerReleaseOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ContainerReleaseOptions) format() *generated.ContainerClientReleaseLeaseOptions {
	if o == nil {
		return nil
	}
	// Note: missing mapping for o.ModifiedAccessConditions.IfMatch, o.ModifiedAccessConditions.IfNoneMatch

	opts := &generated.ContainerClientReleaseLeaseOptions{}
	if o.ModifiedAccessConditions != nil {
		opts.IfModifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfModifiedSince)
		opts.IfUnmodifiedSince = shared.ConvertToGMT(o.ModifiedAccessConditions.IfUnmodifiedSince)
	}

	return opts
}

func leasePeriodPointer(period int32) *int32 {
	if period != BreakNaturally {
		return &period
	} else {
		return nil
	}
}
