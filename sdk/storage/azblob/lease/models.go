//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
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
	// Specifies the Duration of the lease, in seconds, or negative one (-1) for a lease that never expires. A non-infinite lease
	// can be between 15 and 60 seconds. A lease Duration cannot be changed using renew or change.
	Duration *int32

	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BlobAcquireOptions) format() (generated.BlobClientAcquireLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return generated.BlobClientAcquireLeaseOptions{}, nil
	}
	return generated.BlobClientAcquireLeaseOptions{
		Duration: o.Duration,
	}, o.ModifiedAccessConditions
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

func (o *BlobBreakOptions) format() (*generated.BlobClientBreakLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	if o.BreakPeriod != nil {
		period := leasePeriodPointer(*o.BreakPeriod)
		return &generated.BlobClientBreakLeaseOptions{
			BreakPeriod: period,
		}, o.ModifiedAccessConditions
	}

	return nil, o.ModifiedAccessConditions
}

// BlobChangeOptions contains the optional parameters for the LeaseClient.ChangeLease method.
type BlobChangeOptions struct {
	ProposedLeaseID          *string
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BlobChangeOptions) format() (*string, *generated.BlobClientChangeLeaseOptions, *ModifiedAccessConditions, error) {
	generatedUuid, err := uuid.New()
	if err != nil {
		return nil, nil, nil, err
	}
	leaseID := to.Ptr(generatedUuid.String())
	if o == nil {
		return leaseID, nil, nil, nil
	}

	if o.ProposedLeaseID == nil {
		o.ProposedLeaseID = leaseID
	}

	return o.ProposedLeaseID, nil, o.ModifiedAccessConditions, nil
}

// BlobRenewOptions contains the optional parameters for the LeaseClient.RenewLease method.
type BlobRenewOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BlobRenewOptions) format() (*generated.BlobClientRenewLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.ModifiedAccessConditions
}

// BlobReleaseOptions contains the optional parameters for the LeaseClient.ReleaseLease method.
type BlobReleaseOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BlobReleaseOptions) format() (*generated.BlobClientReleaseLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.ModifiedAccessConditions
}

///////////////////////////

// ContainerAcquireOptions contains the optional parameters for the LeaseClient.AcquireLease method.
type ContainerAcquireOptions struct {
	// Specifies the Duration of the lease, in seconds, or negative one (-1) for a lease that never expires. A non-infinite lease
	// can be between 15 and 60 seconds. A lease Duration cannot be changed using renew or change.
	Duration *int32

	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ContainerAcquireOptions) format() (generated.ContainerClientAcquireLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return generated.ContainerClientAcquireLeaseOptions{}, nil
	}
	return generated.ContainerClientAcquireLeaseOptions{
		Duration: o.Duration,
	}, o.ModifiedAccessConditions
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

func (o *ContainerBreakOptions) format() (*generated.ContainerClientBreakLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	if o.BreakPeriod != nil {
		period := leasePeriodPointer(*o.BreakPeriod)
		return &generated.ContainerClientBreakLeaseOptions{
			BreakPeriod: period,
		}, o.ModifiedAccessConditions
	}

	return nil, o.ModifiedAccessConditions
}

// ContainerChangeOptions contains the optional parameters for the LeaseClient.ChangeLease method.
type ContainerChangeOptions struct {
	ProposedLeaseID          *string
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ContainerChangeOptions) format() (*string, *generated.ContainerClientChangeLeaseOptions, *ModifiedAccessConditions, error) {
	generatedUuid, err := uuid.New()
	if err != nil {
		return nil, nil, nil, err
	}
	leaseID := to.Ptr(generatedUuid.String())
	if o == nil {
		return leaseID, nil, nil, nil
	}

	if o.ProposedLeaseID == nil {
		o.ProposedLeaseID = leaseID
	}

	return o.ProposedLeaseID, nil, o.ModifiedAccessConditions, nil
}

// ContainerRenewOptions contains the optional parameters for the LeaseClient.RenewLease method.
type ContainerRenewOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ContainerRenewOptions) format() (*generated.ContainerClientRenewLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.ModifiedAccessConditions
}

// ContainerReleaseOptions contains the optional parameters for the LeaseClient.ReleaseLease method.
type ContainerReleaseOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ContainerReleaseOptions) format() (*generated.ContainerClientReleaseLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.ModifiedAccessConditions
}

func leasePeriodPointer(period int32) *int32 {
	if period != BreakNaturally {
		return &period
	} else {
		return nil
	}
}
