//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

// ---------------------------------------------------------------------------------------------------------------------

// LeaseBreakNaturally tells ContainerClient's or BlobClient's Break method to break the lease using service semantics.
const LeaseBreakNaturally = -1

func leasePeriodPointer(period int32) *int32 {
	if period != LeaseBreakNaturally {
		return &period
	} else {
		return nil
	}
}

type AcquireOptions struct {
	LeaseID *string

	// Specifies the Duration of the lease, in seconds, or negative one (-1) for a lease that never expires. A non-infinite lease
	// can be between 15 and 60 seconds. A lease Duration cannot be changed using renew or change.
	Duration *int32

	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *AcquireOptions) format() (generated.ContainerClientAcquireLeaseOptions, *generated.ModifiedAccessConditions) {
	if o == nil {
		return generated.ContainerClientAcquireLeaseOptions{}, nil
	}
	containerAcquireLeaseOptions := generated.ContainerClientAcquireLeaseOptions{
		Duration: o.Duration,
	}

	return containerAcquireLeaseOptions, o.ModifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// BreakOptions provides set of configurations for BreakLeaseContainer operation
type BreakOptions struct {
	BreakPeriod              *int32
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BreakOptions) format() (*generated.ContainerClientBreakLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	containerBreakLeaseOptions := &generated.ContainerClientBreakLeaseOptions{
		BreakPeriod: o.BreakPeriod,
	}

	return containerBreakLeaseOptions, o.ModifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// ChangeOptions provides set of configurations for ChangeLeaseContainer operation
type ChangeOptions struct {
	ProposedLeaseID          *string
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ChangeOptions) format() (*string, *generated.ContainerClientChangeLeaseOptions, *generated.ModifiedAccessConditions, error) {
	generatedUuid, err := uuid.New()
	if err != nil {
		return nil, nil, nil, err
	}
	leaseID := to.Ptr(generatedUuid.String())
	if o == nil {
		return leaseID, nil, nil, err
	}

	if o.ProposedLeaseID == nil {
		o.ProposedLeaseID = leaseID
	}

	return o.ProposedLeaseID, nil, o.ModifiedAccessConditions, err
}

// ---------------------------------------------------------------------------------------------------------------------

// ReleaseOptions provides set of configurations for ReleaseLeaseContainer operation
type ReleaseOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ReleaseOptions) format() (*generated.ContainerClientReleaseLeaseOptions, *generated.ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.ModifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------

// RenewOptions provides set of configurations for RenewLeaseContainer operation
type RenewOptions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *RenewOptions) format() (*generated.ContainerClientRenewLeaseOptions, *generated.ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.ModifiedAccessConditions
}

// ---------------------------------------------------------------------------------------------------------------------
