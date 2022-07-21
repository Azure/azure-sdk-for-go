//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

type CreateOptions struct {
	// Specifies whether data in the container may be accessed publicly and the level of access
	Access *PublicAccessType

	// Optional. Specifies a user-defined name-value pair associated with the blob.
	Metadata map[string]string

	// Optional. Specifies the encryption scope settings to set on the container.
	CpkScope *ContainerCpkScopeInfo
}

type AccessConditions = exported.ContainerAccessConditions
type LeaseAccessConditions = exported.LeaseAccessConditions
type ModifiedAccessConditions = exported.ModifiedAccessConditions

type DeleteOptions struct {
	AccessConditions *AccessConditions
}

func (o *DeleteOptions) format() (*generated.ContainerClientDeleteOptions, *LeaseAccessConditions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	leaseAccessConditions, modifiedAccessConditions := exported.FormatContainerAccessConditions(o.AccessConditions)
	return nil, leaseAccessConditions, modifiedAccessConditions
}

type GetPropertiesOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetPropertiesOptions) format() (*generated.ContainerClientGetPropertiesOptions, *generated.LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.LeaseAccessConditions
}

type ListBlobsFlatOptions struct {
	// Include this parameter to specify one or more datasets to include in the response.
	Include []ListBlobsIncludeItem
	// A string value that identifies the portion of the list of containers to be returned with the next listing operation. The
	// operation returns the NextMarker value within the response body if the listing
	// operation did not return all containers remaining to be listed with the current page. The NextMarker value can be used
	// as the value for the marker parameter in a subsequent call to request the next
	// page of list items. The marker value is opaque to the client.
	Marker *string
	// Specifies the maximum number of containers to return. If the request does not specify maxresults, or specifies a value
	// greater than 5000, the server will return up to 5000 items. Note that if the
	// listing operation crosses a partition boundary, then the service will return a continuation token for retrieving the remainder
	// of the results. For this reason, it is possible that the service will
	// return fewer results than specified by maxresults, or than the default of 5000.
	MaxResults *int32
	// Filters the results to return only containers whose name begins with the specified prefix.
	Prefix *string
}

//ListBlobsHierarchyOptions provides set of configurations for ContainerClient.NewListBlobsHierarchyPager
type ListBlobsHierarchyOptions struct {
	// Include this parameter to specify one or more datasets to include in the response.
	Include []ListBlobsIncludeItem
	// A string value that identifies the portion of the list of containers to be returned with the next listing operation. The
	// operation returns the NextMarker value within the response body if the listing
	// operation did not return all containers remaining to be listed with the current page. The NextMarker value can be used
	// as the value for the marker parameter in a subsequent call to request the next
	// page of list items. The marker value is opaque to the client.
	Marker *string
	// Specifies the maximum number of containers to return. If the request does not specify maxresults, or specifies a value
	// greater than 5000, the server will return up to 5000 items. Note that if the
	// listing operation crosses a partition boundary, then the service will return a continuation token for retrieving the remainder
	// of the results. For this reason, it is possible that the service will
	// return fewer results than specified by maxresults, or than the default of 5000.
	MaxResults *int32
	// Filters the results to return only containers whose name begins with the specified prefix.
	Prefix *string
}

func (o *ListBlobsHierarchyOptions) format() generated.ContainerClientListBlobHierarchySegmentOptions {
	if o == nil {
		return generated.ContainerClientListBlobHierarchySegmentOptions{}
	}

	return generated.ContainerClientListBlobHierarchySegmentOptions{
		Include:    o.Include,
		Marker:     o.Marker,
		Maxresults: o.MaxResults,
		Prefix:     o.Prefix,
	}
}

type ContainerCpkScopeInfo = generated.ContainerCpkScopeInfo

type PublicAccessType = generated.PublicAccessType

type BlobItem = generated.BlobItemInternal

// SetMetadataOptions provides set of configurations for SetMetadataContainer operation
type SetMetadataOptions struct {
	Metadata                 map[string]string
	LeaseAccessConditions    *LeaseAccessConditions
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *SetMetadataOptions) format() (*generated.ContainerClientSetMetadataOptions, *generated.LeaseAccessConditions, *generated.ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	return &generated.ContainerClientSetMetadataOptions{Metadata: o.Metadata}, o.LeaseAccessConditions, o.ModifiedAccessConditions
}

// GetAccessPolicyOptions provides set of configurations for GetAccessPolicy operation
type GetAccessPolicyOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetAccessPolicyOptions) format() (*generated.ContainerClientGetAccessPolicyOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.LeaseAccessConditions
}

type SignedIdentifier = generated.SignedIdentifier

// SetAccessPolicyOptions provides set of configurations for ContainerClient.SetAccessPolicy operation
type SetAccessPolicyOptions struct {
	// Specifies whether data in the container may be accessed publicly and the level of access
	Access *PublicAccessType
	// the acls for the container
	ContainerACL []*SignedIdentifier
	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	AccessConditions *AccessConditions
}

func (o *SetAccessPolicyOptions) format() (*generated.ContainerClientSetAccessPolicyOptions, *LeaseAccessConditions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}
	lac, mac := exported.FormatContainerAccessConditions(o.AccessConditions)
	return &generated.ContainerClientSetAccessPolicyOptions{
		Access:       o.Access,
		ContainerACL: o.ContainerACL,
	}, lac, mac
}

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
