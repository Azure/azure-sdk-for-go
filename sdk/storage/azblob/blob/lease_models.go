//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

// LeaseBreakNaturally tells ContainerClient's or BlobClient's Break method to break the lease using service semantics.
const LeaseBreakNaturally = -1

func leasePeriodPointer(period int32) *int32 {
	if period != LeaseBreakNaturally {
		return &period
	} else {
		return nil
	}
}

// ---------------------------------------------------------------------------------------------------------------------

// AcquireOptions provides set of configurations for AcquireLeaseBlob operation
type AcquireOptions struct {
	// Specifies the Duration of the lease, in seconds, or negative one (-1) for a lease that never expires. A non-infinite lease
	// can be between 15 and 60 seconds. A lease Duration cannot be changed using renew or change.
	Duration *int32

	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *AcquireOptions) format() (generated.BlobClientAcquireLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return generated.BlobClientAcquireLeaseOptions{}, nil
	}
	return generated.BlobClientAcquireLeaseOptions{
		Duration: o.Duration,
	}, o.ModifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// BreakOptions provides set of configurations for BreakLeaseBlob operation
type BreakOptions struct {
	// For a break operation, proposed Duration the lease should continue before it is broken, in seconds, between 0 and 60. This
	// break period is only used if it is shorter than the time remaining on the lease. If longer, the time remaining on the lease
	// is used. A new lease will not be available before the break period has expired, but the lease may be held for longer than
	// the break period. If this header does not appear with a break operation, a fixed-Duration lease breaks after the remaining
	// lease period elapses, and an infinite lease breaks immediately.
	BreakPeriod              *int32
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BreakOptions) format() (*generated.BlobClientBreakLeaseOptions, *ModifiedAccessConditions) {
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

// ---------------------------------------------------------------------------------------------------------------------

// ChangeOptions provides set of configurations for ChangeLeaseBlob operation
type ChangeOptions struct {
	ProposedLeaseID          *string
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ChangeOptions) format() (*string, *generated.BlobClientChangeLeaseOptions, *ModifiedAccessConditions, error) {
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

// ---------------------------------------------------------------------------------------------------------------------

// RenewOptions provides set of configurations for RenewLeaseBlob operation
type RenewOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *RenewOptions) format() (*generated.BlobClientRenewLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.ModifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// ReleaseOptions provides set of configurations for ReleaseLeaseBlob operation
type ReleaseOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ReleaseOptions) format() (*generated.BlobClientReleaseLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.ModifiedAccessConditions
}
