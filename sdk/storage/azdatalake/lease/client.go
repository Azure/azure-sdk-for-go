//go:build go1.18

//

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// FileSystemClient provides lease functionality for the underlying filesystem client.
type FileSystemClient struct {
	containerClient *lease.ContainerClient
	leaseID         *string
}

// FileSystemClientOptions contains the optional values when creating a FileSystemClient.
type FileSystemClientOptions = lease.ContainerClientOptions

// NewFileSystemClient creates a filesystem lease client for the provided filesystem client.
//   - client - an instance of a filesystem client
//   - options - client options; pass nil to accept the default values
func NewFileSystemClient(client *filesystem.Client, options *FileSystemClientOptions) (*FileSystemClient, error) {
	_, _, containerClient := base.InnerClients((*base.CompositeClient[generated.FileSystemClient, generated.FileSystemClient, container.Client])(client))
	containerLeaseClient, err := lease.NewContainerClient(containerClient, options)
	if err != nil {
		return nil, exported.ConvertToDFSError(err)
	}
	return &FileSystemClient{
		containerClient: containerLeaseClient,
		leaseID:         containerLeaseClient.LeaseID(),
	}, nil
}

// LeaseID returns leaseID of the client.
func (c *FileSystemClient) LeaseID() *string {
	return c.leaseID
}

// AcquireLease acquires a lease on the filesystem for write and delete operations.
// The lease Duration must be between 15 and 60 seconds, or infinite (-1).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *FileSystemClient) AcquireLease(ctx context.Context, duration int32, o *FileSystemAcquireOptions) (FileSystemAcquireResponse, error) {
	opts := o.format()
	resp, err := c.containerClient.AcquireLease(ctx, duration, opts)
	return resp, exported.ConvertToDFSError(err)
}

// BreakLease breaks the filesystem's previously-acquired lease (if it exists). Pass the LeaseBreakDefault (-1)
// constant to break a fixed-Duration lease when it expires or an infinite lease immediately.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *FileSystemClient) BreakLease(ctx context.Context, o *FileSystemBreakOptions) (FileSystemBreakResponse, error) {
	opts := o.format()
	resp, err := c.containerClient.BreakLease(ctx, opts)
	return resp, exported.ConvertToDFSError(err)
}

// ChangeLease changes the filesystem's lease ID.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *FileSystemClient) ChangeLease(ctx context.Context, proposedLeaseID string, o *FileSystemChangeOptions) (FileSystemChangeResponse, error) {
	opts := o.format()
	resp, err := c.containerClient.ChangeLease(ctx, proposedLeaseID, opts)
	if err != nil {
		return resp, exported.ConvertToDFSError(err)
	}
	c.leaseID = &proposedLeaseID
	return resp, nil
}

// RenewLease renews the filesystem's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *FileSystemClient) RenewLease(ctx context.Context, o *FileSystemRenewOptions) (FileSystemRenewResponse, error) {
	opts := o.format()
	resp, err := c.containerClient.RenewLease(ctx, opts)
	return resp, exported.ConvertToDFSError(err)
}

// ReleaseLease releases the filesystem's previously-acquired lease.
func (c *FileSystemClient) ReleaseLease(ctx context.Context, o *FileSystemReleaseOptions) (FileSystemReleaseResponse, error) {
	opts := o.format()
	resp, err := c.containerClient.ReleaseLease(ctx, opts)
	return resp, exported.ConvertToDFSError(err)
}
