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
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

// LeaseClient represents a URL to the Azure Storage container allowing you to manipulate its blobs.
type LeaseClient struct {
	blobClient *Client
	leaseID    *string
}

// NewLeaseClient creates a Client object using the specified URL, Azure AD credential, and options.
func NewLeaseClient(containerURL string, cred azcore.TokenCredential, leaseID *string, options *ClientOptions) (*LeaseClient, error) {
	leaseID, err := shared.GenerateLeaseID(leaseID)
	if err != nil {
		return nil, err
	}
	blobClient, err := NewClient(containerURL, cred, options)
	return &LeaseClient{
		blobClient: blobClient,
		leaseID:    leaseID,
	}, err
}

// NewLeaseClientWithNoCredential creates a Client object using the specified URL and options.
func NewLeaseClientWithNoCredential(containerURL string, leaseID *string, options *ClientOptions) (*LeaseClient, error) {
	leaseID, err := shared.GenerateLeaseID(leaseID)
	if err != nil {
		return nil, err
	}
	blobClient, err := NewClientWithNoCredential(containerURL, options)
	return &LeaseClient{
		blobClient: blobClient,
		leaseID:    leaseID,
	}, err
}

// NewLeaseClientWithSharedKey creates a Client object using the specified URL, shared key, and options.
func NewLeaseClientWithSharedKey(containerURL string, cred *SharedKeyCredential, leaseID *string, options *ClientOptions) (*LeaseClient, error) {
	leaseID, err := shared.GenerateLeaseID(leaseID)
	if err != nil {
		return nil, err
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

// BlobClient returns the embedded blobClient from leaseClient
func (c *LeaseClient) BlobClient() *Client {
	return c.blobClient
}

// LeaseID returns leaseID of the client.
func (c *LeaseClient) LeaseID() *string {
	return c.leaseID
}

func (c *LeaseClient) generated() *generated.BlobClient {
	return base.InnerClient((*base.Client[generated.BlobClient])(c.blobClient))
}

// AcquireLease acquires a lease on the blob for write and delete operations.
//The lease Duration must be between 15 and 60 seconds, or infinite (-1).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *LeaseClient) AcquireLease(ctx context.Context, o *AcquireLeaseOptions) (AcquireLeaseResponse, error) {
	blobAcquireLeaseOptions, modifiedAccessConditions := o.format()
	blobAcquireLeaseOptions.ProposedLeaseID = c.leaseID

	resp, err := c.generated().AcquireLease(ctx, &blobAcquireLeaseOptions, modifiedAccessConditions)
	return resp, err
}

// BreakLease breaks the blob's previously-acquired lease (if it exists). Pass the LeaseBreakDefault (-1)
// constant to break a fixed-Duration lease when it expires or an infinite lease immediately.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *LeaseClient) BreakLease(ctx context.Context, o *BreakLeaseOptions) (BreakLeaseResponse, error) {
	blobBreakLeaseOptions, modifiedAccessConditions := o.format()
	resp, err := c.generated().BreakLease(ctx, blobBreakLeaseOptions, modifiedAccessConditions)
	return resp, err
}

// ChangeLease changes the blob's lease ID.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *LeaseClient) ChangeLease(ctx context.Context, o *ChangeLeaseOptions) (ChangeLeaseResponse, error) {
	if c.leaseID == nil {
		return ChangeLeaseResponse{}, errors.New("leaseID cannot be nil")
	}
	proposedLeaseID, changeLeaseOptions, modifiedAccessConditions, err := o.format()
	if err != nil {
		return ChangeLeaseResponse{}, err
	}
	resp, err := c.generated().ChangeLease(ctx, *c.leaseID, *proposedLeaseID, changeLeaseOptions, modifiedAccessConditions)

	// If lease has been changed successfully, set the leaseID in client
	if err == nil {
		c.leaseID = proposedLeaseID
	}

	return resp, err
}

// RenewLease renews the blob's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *LeaseClient) RenewLease(ctx context.Context, o *RenewLeaseOptions) (RenewLeaseResponse, error) {
	if c.leaseID == nil {
		return RenewLeaseResponse{}, errors.New("leaseID cannot be nil")
	}
	renewLeaseBlobOptions, modifiedAccessConditions := o.format()
	resp, err := c.generated().RenewLease(ctx, *c.leaseID, renewLeaseBlobOptions, modifiedAccessConditions)
	return resp, err
}

// ReleaseLease releases the blob's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *LeaseClient) ReleaseLease(ctx context.Context, o *ReleaseLeaseOptions) (ReleaseLeaseResponse, error) {
	if c.leaseID == nil {
		return ReleaseLeaseResponse{}, errors.New("leaseID cannot be nil")
	}
	renewLeaseBlobOptions, modifiedAccessConditions := o.format()
	resp, err := c.generated().ReleaseLease(ctx, *c.leaseID, renewLeaseBlobOptions, modifiedAccessConditions)
	return resp, err
}
