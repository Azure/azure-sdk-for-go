// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

type ContainerLeaseClient struct {
	ContainerClient
	LeaseID *string
}

func NewContainerLeaseClient(containerURL string, leaseID *string, cred azcore.Credential, options *connectionOptions) (ContainerLeaseClient, error) {

	containerClient := ContainerClient{
		client: &containerClient{
			con: newConnection(containerURL, cred, options),
		}, cred: cred,
	}

	if leaseID == nil {
		generatedUuid, err := uuid.New()
		if err != nil {
			return ContainerLeaseClient{}, err
		}
		leaseID = to.StringPtr(generatedUuid.String())
	}

	return ContainerLeaseClient{
		ContainerClient: containerClient,
		LeaseID:         leaseID,
	}, nil
}

// URL returns the URL endpoint used by the ContainerClient object.
func (clc ContainerLeaseClient) URL() string {
	return clc.client.con.u
}

// AcquireLease acquires a lease on the container for delete operations. The lease Duration must be between 15 to 60 seconds, or infinite (-1).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (clc *ContainerLeaseClient) AcquireLease(ctx context.Context, options *AcquireLeaseContainerOptions) (ContainerAcquireLeaseResponse, error) {
	containerAcquireLeaseOptions, modifiedAccessConditions := options.pointers()
	containerAcquireLeaseOptions.ProposedLeaseID = clc.LeaseID

	resp, err := clc.client.AcquireLease(ctx, containerAcquireLeaseOptions, modifiedAccessConditions)
	if err == nil && resp.LeaseID != nil {
		clc.LeaseID = resp.LeaseID
	}
	return resp, handleError(err)
}

// BreakLease breaks the container's previously-acquired lease (if it exists).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (clc *ContainerLeaseClient) BreakLease(ctx context.Context, options *BreakLeaseContainerOptions) (ContainerBreakLeaseResponse, error) {
	containerBreakLeaseOptions, modifiedAccessConditions := options.pointers()
	resp, err := clc.client.BreakLease(ctx, containerBreakLeaseOptions, modifiedAccessConditions)
	return resp, handleError(err)
}

// ChangeLease changes the container's lease ID.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (clc *ContainerLeaseClient) ChangeLease(ctx context.Context, options *ChangeLeaseContainerOptions) (ContainerChangeLeaseResponse, error) {
	if clc.LeaseID == nil {
		return ContainerChangeLeaseResponse{}, errors.New("LeaseID cannot be nil")
	}
	proposedLeaseID, modifiedAccessConditions, err := options.pointers()
	if err != nil {
		return ContainerChangeLeaseResponse{}, err
	}

	resp, err := clc.client.ChangeLease(ctx, *clc.LeaseID, *proposedLeaseID, nil, modifiedAccessConditions)
	if err == nil && resp.LeaseID != nil {
		clc.LeaseID = resp.LeaseID
	}
	return resp, handleError(err)
}

// ReleaseLease releases the container's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (clc *ContainerLeaseClient) ReleaseLease(ctx context.Context, options *ReleaseLeaseContainerOptions) (ContainerReleaseLeaseResponse, error) {
	containerReleaseLeaseOptions, modifiedAccessConditions := options.pointers()
	resp, err := clc.client.ReleaseLease(ctx, *clc.LeaseID, containerReleaseLeaseOptions, modifiedAccessConditions)
	return resp, handleError(err)
}

// RenewLease renews the container's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (clc *ContainerLeaseClient) RenewLease(ctx context.Context, options *RenewLeaseContainerOptions) (ContainerRenewLeaseResponse, error) {
	if clc.LeaseID == nil {
		return ContainerRenewLeaseResponse{}, errors.New("LeaseID cannot be nil")
	}
	renewLeaseBlobOptions, modifiedAccessConditions := options.pointers()
	resp, err := clc.client.RenewLease(ctx, *clc.LeaseID, renewLeaseBlobOptions, modifiedAccessConditions)
	if err == nil && resp.LeaseID != nil {
		clc.LeaseID = resp.LeaseID
	}
	return resp, handleError(err)
}
