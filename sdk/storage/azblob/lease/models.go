// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
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

func (o *BlobAcquireOptions) format() *generated.BlobClientAcquireLeaseOptions {
	if o == nil {
		return nil
	}
	return &generated.BlobClientAcquireLeaseOptions{
		IfMatch:           o.ModifiedAccessConditions.IfMatch,
		IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
		IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
	}
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

	return &generated.BlobClientBreakLeaseOptions{
		BreakPeriod:       period,
		IfMatch:           o.ModifiedAccessConditions.IfMatch,
		IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
		IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
	}
}

// BlobChangeOptions contains the optional parameters for the LeaseClient.ChangeLease method.
type BlobChangeOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BlobChangeOptions) format() *generated.BlobClientChangeLeaseOptions {
	if o == nil {
		return nil
	}
	return &generated.BlobClientChangeLeaseOptions{
		IfMatch:           o.ModifiedAccessConditions.IfMatch,
		IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
		IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
	}
}

// BlobRenewOptions contains the optional parameters for the LeaseClient.RenewLease method.
type BlobRenewOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BlobRenewOptions) format() *generated.BlobClientRenewLeaseOptions {
	if o == nil {
		return nil
	}
	return &generated.BlobClientRenewLeaseOptions{
		IfMatch:           o.ModifiedAccessConditions.IfMatch,
		IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
		IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
	}
}

// BlobReleaseOptions contains the optional parameters for the LeaseClient.ReleaseLease method.
type BlobReleaseOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BlobReleaseOptions) format() *generated.BlobClientReleaseLeaseOptions {
	if o == nil {
		return nil
	}

	return &generated.BlobClientReleaseLeaseOptions{
		IfMatch:           o.ModifiedAccessConditions.IfMatch,
		IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
		IfNoneMatch:       o.ModifiedAccessConditions.IfNoneMatch,
		IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
	}
}

// ContainerAcquireOptions contains the optional parameters for the LeaseClient.AcquireLease method.
type ContainerAcquireOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ContainerAcquireOptions) format() *generated.ContainerClientAcquireLeaseOptions {
	if o == nil {
		return nil
	}
	// Note: missing mapping for o.ModifiedAccessConditions.IfMatch, o.ModifiedAccessConditions.IfNoneMatch
	return &generated.ContainerClientAcquireLeaseOptions{
		IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
		IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
	}
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
	return &generated.ContainerClientBreakLeaseOptions{
		BreakPeriod:       period,
		IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
		IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
	}
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
	return &generated.ContainerClientChangeLeaseOptions{
		IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
		IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
	}
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
	return &generated.ContainerClientRenewLeaseOptions{
		IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
		IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
	}
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
	return &generated.ContainerClientReleaseLeaseOptions{
		IfModifiedSince:   o.ModifiedAccessConditions.IfModifiedSince,
		IfUnmodifiedSince: o.ModifiedAccessConditions.IfUnmodifiedSince,
	}
}

func leasePeriodPointer(period int32) *int32 {
	if period != BreakNaturally {
		return &period
	} else {
		return nil
	}
}
