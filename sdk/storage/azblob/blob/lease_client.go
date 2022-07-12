//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

type LeaseClient struct {
	blobClient *Client
	leaseID    *string
}

// NewLeaseClient creates a Client object using the specified URL, Azure AD credential, and options.
func NewLeaseClient(containerURL string, cred azcore.TokenCredential, leaseID *string, options *ClientOptions) (*LeaseClient, error) {
	if leaseID == nil {
		generatedUuid, err := uuid.New()
		if err != nil {
			return nil, err
		}
		leaseID = to.Ptr(generatedUuid.String())
	}
	blobClient, err := NewClient(containerURL, cred, options)
	return &LeaseClient{
		blobClient: blobClient,
		leaseID:    leaseID,
	}, err
}

// NewLeaseClientWithNoCredential creates a Client object using the specified URL and options.
func NewLeaseClientWithNoCredential(containerURL string, leaseID *string, options *ClientOptions) (*LeaseClient, error) {
	if leaseID == nil {
		generatedUuid, err := uuid.New()
		if err != nil {
			return nil, err
		}
		leaseID = to.Ptr(generatedUuid.String())
	}
	blobClient, err := NewClientWithNoCredential(containerURL, options)
	return &LeaseClient{
		blobClient: blobClient,
		leaseID:    leaseID,
	}, err
}

// NewLeaseClientWithSharedKey creates a Client object using the specified URL, shared key, and options.
func NewLeaseClientWithSharedKey(containerURL string, cred *SharedKeyCredential, leaseID *string, options *ClientOptions) (*LeaseClient, error) {
	if leaseID == nil {
		generatedUuid, err := uuid.New()
		if err != nil {
			return nil, err
		}
		leaseID = to.Ptr(generatedUuid.String())
	}
	blobClient, err := NewClientWithSharedKey(containerURL, cred, options)
	return &LeaseClient{
		blobClient: blobClient,
		leaseID:    leaseID,
	}, err
}

// NewLeaseClientFromConnectionString creates a Client object using connection string of an account
func NewLeaseClientFromConnectionString(connectionString string, containerName string, leaseID *string, options *ClientOptions) (*LeaseClient, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	parsed.ServiceURL = runtime.JoinPaths(parsed.ServiceURL, containerName)

	if parsed.AccountKey != "" && parsed.AccountName != "" {
		credential, err := exported.NewSharedKeyCredential(parsed.AccountName, parsed.AccountKey)
		if err != nil {
			return nil, err
		}
		return NewLeaseClientWithSharedKey(parsed.ServiceURL, credential, leaseID, options)
	}

	return NewLeaseClientWithNoCredential(parsed.ServiceURL, leaseID, options)
}

func (c *LeaseClient) generated() *generated.BlobClient {
	return base.InnerClient((*base.Client[generated.BlobClient])(c.blobClient))
}

// Acquire acquires a lease on the blob for write and delete operations.
//The lease Duration must be between 15 and 60 seconds, or infinite (-1).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *LeaseClient) Acquire(ctx context.Context, o *AcquireOptions) (AcquireResponse, error) {
	blobAcquireLeaseOptions, modifiedAccessConditions := o.format()
	blobAcquireLeaseOptions.ProposedLeaseID = c.leaseID

	resp, err := c.generated().AcquireLease(ctx, &blobAcquireLeaseOptions, modifiedAccessConditions)
	return resp, err
}

// Break breaks the blob's previously-acquired lease (if it exists). Pass the LeaseBreakDefault (-1)
// constant to break a fixed-Duration lease when it expires or an infinite lease immediately.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *LeaseClient) Break(ctx context.Context, o *BreakOptions) (BreakResponse, error) {
	blobBreakLeaseOptions, modifiedAccessConditions := o.format()
	resp, err := c.generated().BreakLease(ctx, blobBreakLeaseOptions, modifiedAccessConditions)
	return resp, err
}

// Change changes the blob's lease ID.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *LeaseClient) Change(ctx context.Context, o *ChangeOptions) (ChangeResponse, error) {
	if c.leaseID == nil {
		return ChangeResponse{}, errors.New("leaseID cannot be nil")
	}
	proposedLeaseID, changeLeaseOptions, modifiedAccessConditions, err := o.format()
	if err != nil {
		return ChangeResponse{}, err
	}
	resp, err := c.generated().ChangeLease(ctx, *c.leaseID, *proposedLeaseID, changeLeaseOptions, modifiedAccessConditions)

	// If lease has been changed successfully, set the leaseID in client
	if err == nil {
		c.leaseID = proposedLeaseID
	}

	return resp, err
}

// Renew renews the blob's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *LeaseClient) Renew(ctx context.Context, o *RenewOptions) (RenewResponse, error) {
	if c.leaseID == nil {
		return RenewResponse{}, errors.New("leaseID cannot be nil")
	}
	renewLeaseBlobOptions, modifiedAccessConditions := o.format()
	resp, err := c.generated().RenewLease(ctx, *c.leaseID, renewLeaseBlobOptions, modifiedAccessConditions)
	return resp, err
}

// Release releases the blob's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *LeaseClient) Release(ctx context.Context, o *ReleaseOptions) (ReleaseResponse, error) {
	if c.leaseID == nil {
		return ReleaseResponse{}, errors.New("leaseID cannot be nil")
	}
	renewLeaseBlobOptions, modifiedAccessConditions := o.format()
	resp, err := c.generated().ReleaseLease(ctx, *c.leaseID, renewLeaseBlobOptions, modifiedAccessConditions)
	return resp, err
}
