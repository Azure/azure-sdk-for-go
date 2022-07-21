//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container

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

type LeaseClient struct {
	containerClient *Client
	leaseID         *string
}

// NewLeaseClient creates a Client object using the specified URL, Azure AD credential, and options.
func NewLeaseClient(containerURL string, cred azcore.TokenCredential, leaseID *string, options *ClientOptions) (*LeaseClient, error) {
	leaseID, err := shared.GenerateLeaseID(leaseID)
	if err != nil {
		return nil, err
	}
	containerClient, err := NewClient(containerURL, cred, options)
	return &LeaseClient{
		containerClient: containerClient,
		leaseID:         leaseID,
	}, err
}

// NewLeaseClientWithNoCredential creates a Client object using the specified URL and options.
func NewLeaseClientWithNoCredential(containerURL string, leaseID *string, options *ClientOptions) (*LeaseClient, error) {
	leaseID, err := shared.GenerateLeaseID(leaseID)
	if err != nil {
		return nil, err
	}
	containerClient, err := NewClientWithNoCredential(containerURL, options)
	return &LeaseClient{
		containerClient: containerClient,
		leaseID:         leaseID,
	}, err
}

// NewLeaseClientWithSharedKey creates a Client object using the specified URL, shared key, and options.
func NewLeaseClientWithSharedKey(containerURL string, cred *SharedKeyCredential, leaseID *string, options *ClientOptions) (*LeaseClient, error) {
	leaseID, err := shared.GenerateLeaseID(leaseID)
	if err != nil {
		return nil, err
	}
	containerClient, err := NewClientWithSharedKey(containerURL, cred, options)
	return &LeaseClient{
		containerClient: containerClient,
		leaseID:         leaseID,
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

func (c *LeaseClient) generated() *generated.ContainerClient {
	return base.InnerClient((*base.Client[generated.ContainerClient])(c.containerClient))
}

func (c *LeaseClient) AcquireLease(ctx context.Context, o *AcquireLeaseOptions) (AcquireResponse, error) {
	opts, modifiedAccessConditions := o.format()
	opts.ProposedLeaseID = c.leaseID

	resp, err := c.generated().AcquireLease(ctx, &generated.ContainerClientAcquireLeaseOptions{
		Duration:        opts.Duration,
		ProposedLeaseID: opts.ProposedLeaseID,
	}, modifiedAccessConditions)
	if err != nil {
		return AcquireResponse{}, err
	}
	return resp, err
}

// BreakLease breaks the container's previously-acquired lease (if it exists).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (c *LeaseClient) BreakLease(ctx context.Context, options *BreakLeaseOptions) (BreakResponse, error) {
	containerBreakLeaseOptions, modifiedAccessConditions := options.format()
	resp, err := c.generated().BreakLease(ctx, containerBreakLeaseOptions, modifiedAccessConditions)
	return resp, err
}

// ChangeLease changes the container's lease ID.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (c *LeaseClient) ChangeLease(ctx context.Context, options *ChangeLeaseOptions) (ChangeResponse, error) {
	if c.leaseID == nil {
		return ChangeResponse{}, errors.New("leaseID cannot be nil")
	}

	proposedLeaseID, changeLeaseOptions, modifiedAccessConditions, err := options.format()
	if err != nil {
		return ChangeResponse{}, err
	}

	resp, err := c.generated().ChangeLease(ctx, *c.leaseID, *proposedLeaseID, changeLeaseOptions, modifiedAccessConditions)
	if err == nil && resp.LeaseID != nil {
		c.leaseID = resp.LeaseID
	}
	return resp, err
}

// ReleaseLease releases the container's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (c *LeaseClient) ReleaseLease(ctx context.Context, options *ReleaseLeaseOptions) (ReleaseResponse, error) {
	if c.leaseID == nil {
		return ReleaseResponse{}, errors.New("leaseID cannot be nil")
	}
	containerReleaseLeaseOptions, modifiedAccessConditions := options.format()
	resp, err := c.generated().ReleaseLease(ctx, *c.leaseID, containerReleaseLeaseOptions, modifiedAccessConditions)

	return resp, err
}

// RenewLease renews the container's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (c *LeaseClient) RenewLease(ctx context.Context, options *RenewLeaseOptions) (RenewResponse, error) {
	if c.leaseID == nil {
		return RenewResponse{}, errors.New("leaseID cannot be nil")
	}
	renewLeaseBlobOptions, modifiedAccessConditions := options.format()
	resp, err := c.generated().RenewLease(ctx, *c.leaseID, renewLeaseBlobOptions, modifiedAccessConditions)
	if err == nil && resp.LeaseID != nil {
		c.leaseID = resp.LeaseID
	}
	return resp, err
}
