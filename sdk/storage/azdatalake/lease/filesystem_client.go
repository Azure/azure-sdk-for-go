//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/filesystem"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// FilesystemClient provides lease functionality for the underlying container client.
type FilesystemClient struct {
	filesystemClient *filesystem.Client
	leaseID          *string
	containerClient  *lease.ContainerClient
}

// FilesystemClientOptions contains the optional values when creating a ContainerClient.
type FilesystemClientOptions struct {
	// LeaseID contains a caller-provided lease ID.
	LeaseID *string
}

func (c *FilesystemClient) generated() *generated.FileSystemClient {
	return base.InnerClient((*base.Client[generated.FileSystemClient])(c.filesystemClient))
}

// NewFilesystemClient creates a container lease client for the provided container client.
//   - client - an instance of a container client
//   - options - client options; pass nil to accept the default values
func NewFilesystemClient(client *filesystem.Client, options *FilesystemClientOptions) (*FilesystemClient, error) {
	// TODO: set up container lease client
	return nil, nil
}

// LeaseID returns leaseID of the client.
func (c *FilesystemClient) LeaseID() *string {
	return c.leaseID
}

// AcquireLease acquires a lease on the blob for write and delete operations.
// The lease Duration must be between 15 and 60 seconds, or infinite (-1).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *FilesystemClient) AcquireLease(ctx context.Context, duration int32, o *FilesystemAcquireOptions) (FilesystemAcquireResponse, error) {
	opts := o.format()
	return c.containerClient.AcquireLease(ctx, duration, opts)
}

// BreakLease breaks the blob's previously-acquired lease (if it exists). Pass the LeaseBreakDefault (-1)
// constant to break a fixed-Duration lease when it expires or an infinite lease immediately.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *FilesystemClient) BreakLease(ctx context.Context, o *FilesystemBreakOptions) (FilesystemBreakResponse, error) {
	opts := o.format()
	return c.containerClient.BreakLease(ctx, opts)
}

// ChangeLease changes the blob's lease ID.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *FilesystemClient) ChangeLease(ctx context.Context, proposedLeaseID string, o *FilesystemChangeOptions) (FilesystemChangeResponse, error) {
	opts := o.format()
	return c.containerClient.ChangeLease(ctx, proposedLeaseID, opts)
}

// RenewLease renews the blob's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *FilesystemClient) RenewLease(ctx context.Context, o *FilesystemRenewOptions) (FilesystemRenewResponse, error) {
	opts := o.format()
	return c.containerClient.RenewLease(ctx, opts)
}

func (c *FilesystemClient) ReleaseLease(ctx context.Context, o *FilesystemReleaseOptions) (FilesystemReleaseResponse, error) {
	opts := o.format()
	return c.containerClient.ReleaseLease(ctx, opts)
}
