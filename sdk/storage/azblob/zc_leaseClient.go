// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

type BlobLeaseClient struct {
	BlobClient
	LeaseID *string
}

func NewBlobLeaseClient(blobURL string, leaseID *string, cred azcore.Credential, options *connectionOptions) (BlobLeaseClient, error) {
	con := newConnection(blobURL, cred, options)
	c, _ := cred.(*SharedKeyCredential)

	blobClient := BlobClient{
		client: &blobClient{con, nil},
		cred:   c,
	}

	if leaseID == nil {
		leaseID = to.StringPtr(newUUID().String())
	}

	return BlobLeaseClient{
		BlobClient: blobClient,
		LeaseID:    leaseID,
	}, nil
}

// URL returns the URL endpoint used by the BlobLeaseClient object.
func (blc BlobLeaseClient) URL() string {
	return blc.client.con.u
}

// AcquireLease acquires a lease on the blob for write and delete operations. The lease Duration must be between
// 15 to 60 seconds, or infinite (-1).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (blc *BlobLeaseClient) AcquireLease(ctx context.Context, options *AcquireLeaseBlobOptions) (BlobAcquireLeaseResponse, error) {
	blobAcquireLeaseOptions, modifiedAccessConditions := options.pointers()
	blobAcquireLeaseOptions.ProposedLeaseID = blc.LeaseID

	resp, err := blc.client.AcquireLease(ctx, blobAcquireLeaseOptions, modifiedAccessConditions)
	return resp, handleError(err)
}

// BreakLease breaks the blob's previously-acquired lease (if it exists). Pass the LeaseBreakDefault (-1)
// constant to break a fixed-Duration lease when it expires or an infinite lease immediately.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (blc *BlobLeaseClient) BreakLease(ctx context.Context, options *BreakLeaseBlobOptions) (BlobBreakLeaseResponse, error) {
	blobBreakLeaseOptions, modifiedAccessConditions := options.pointers()
	resp, err := blc.client.BreakLease(ctx, blobBreakLeaseOptions, modifiedAccessConditions)
	return resp, handleError(err)
}

// ChangeLease changes the blob's lease ID.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (blc *BlobLeaseClient) ChangeLease(ctx context.Context, options *ChangeLeaseBlobOptions) (BlobChangeLeaseResponse, error) {
	if blc.LeaseID == nil {
		return BlobChangeLeaseResponse{}, errors.New("LeaseID cannot be nil")
	}
	proposedLeaseId, modifiedAccessConditions := options.pointers()
	resp, err := blc.client.ChangeLease(ctx, *blc.LeaseID, *proposedLeaseId, nil, modifiedAccessConditions)

	// If lease has been changed successfully, set the LeaseID in client
	if err == nil {
		blc.LeaseID = proposedLeaseId
	}

	return resp, handleError(err)
}

// RenewLease renews the blob's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (blc *BlobLeaseClient) RenewLease(ctx context.Context, options *RenewLeaseBlobOptions) (BlobRenewLeaseResponse, error) {
	if blc.LeaseID == nil {
		return BlobRenewLeaseResponse{}, errors.New("LeaseID cannot be nil")
	}
	renewLeaseBlobOptions, modifiedAccessConditions := options.pointers()
	resp, err := blc.client.RenewLease(ctx, *blc.LeaseID, renewLeaseBlobOptions, modifiedAccessConditions)
	return resp, handleError(err)
}

// ReleaseLease releases the blob's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (blc *BlobLeaseClient) ReleaseLease(ctx context.Context, options *ReleaseLeaseBlobOptions) (BlobReleaseLeaseResponse, error) {
	if blc.LeaseID == nil {
		return BlobReleaseLeaseResponse{}, errors.New("LeaseID cannot be nil")
	}
	renewLeaseBlobOptions, modifiedAccessConditions := options.pointers()
	resp, err := blc.client.ReleaseLease(ctx, *blc.LeaseID, renewLeaseBlobOptions, modifiedAccessConditions)
	return resp, handleError(err)
}

type ContainerLeaseClient struct {
	ContainerClient
	LeaseID *string
}

func NewContainerLeaseClient(containerURL string, leaseID *string, cred azcore.Credential, options *connectionOptions) (ContainerLeaseClient, error) {
	c, _ := cred.(*SharedKeyCredential)

	containerClient := ContainerClient{
		client: &containerClient{
			con: newConnection(containerURL, cred, options),
		}, cred: c,
	}

	if leaseID == nil {
		leaseID = to.StringPtr(newUUID().String())
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
	proposedLeaseId, modifiedAccessConditions := options.pointers()
	resp, err := clc.client.ChangeLease(ctx, *clc.LeaseID, proposedLeaseId, nil, modifiedAccessConditions)

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
