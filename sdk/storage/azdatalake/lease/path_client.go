//go:build go1.18

//

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated_blob"
)

// PathClient provides lease functionality for the underlying path client.
type PathClient struct {
	blobClient *lease.BlobClient
	leaseID    *string
}

// PathClientOptions contains the optional values when creating a PathClient.
type PathClientOptions = lease.BlobClientOptions

// NewPathClient creates a path lease client for the provided path client.
//   - client - an instance of a path client
//   - options - client options; pass nil to accept the default values
func NewPathClient[T directory.Client | file.Client](client *T, options *PathClientOptions) (*PathClient, error) {
	var blobClient *blockblob.Client
	switch t := any(client).(type) {
	case *directory.Client:
		_, _, blobClient = base.InnerClients((*base.CompositeClient[generated.PathClient, generated_blob.BlobClient, blockblob.Client])(t))
	case *file.Client:
		_, _, blobClient = base.InnerClients((*base.CompositeClient[generated.PathClient, generated_blob.BlobClient, blockblob.Client])(t))
	default:
		return nil, fmt.Errorf("unhandled client type %T", client)
	}
	blobLeaseClient, err := lease.NewBlobClient(blobClient, options)
	if err != nil {
		return nil, exported.ConvertToDFSError(err)
	}
	return &PathClient{
		blobClient: blobLeaseClient,
		leaseID:    blobLeaseClient.LeaseID(),
	}, nil
}

// LeaseID returns leaseID of the client.
func (c *PathClient) LeaseID() *string {
	return c.leaseID
}

// AcquireLease acquires a lease on the path for write and delete operations.
// The lease Duration must be between 15 and 60 seconds, or infinite (-1).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *PathClient) AcquireLease(ctx context.Context, duration int32, o *PathAcquireOptions) (PathAcquireResponse, error) {
	opts := o.format()
	resp, err := c.blobClient.AcquireLease(ctx, duration, opts)
	return resp, exported.ConvertToDFSError(err)
}

// BreakLease breaks the path's previously-acquired lease.
func (c *PathClient) BreakLease(ctx context.Context, o *PathBreakOptions) (PathBreakResponse, error) {
	opts := o.format()
	resp, err := c.blobClient.BreakLease(ctx, opts)
	return resp, exported.ConvertToDFSError(err)
}

// ChangeLease changes the path's lease ID.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *PathClient) ChangeLease(ctx context.Context, proposedID string, o *PathChangeOptions) (PathChangeResponse, error) {
	opts := o.format()
	resp, err := c.blobClient.ChangeLease(ctx, proposedID, opts)
	if err != nil {
		return resp, exported.ConvertToDFSError(err)
	}
	c.leaseID = &proposedID
	return resp, nil
}

// RenewLease renews the path's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *PathClient) RenewLease(ctx context.Context, o *PathRenewOptions) (PathRenewResponse, error) {
	opts := o.format()
	resp, err := c.blobClient.RenewLease(ctx, opts)
	return resp, exported.ConvertToDFSError(err)
}

// ReleaseLease releases the path's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *PathClient) ReleaseLease(ctx context.Context, o *PathReleaseOptions) (PathReleaseResponse, error) {
	opts := o.format()
	resp, err := c.blobClient.ReleaseLease(ctx, opts)
	return resp, exported.ConvertToDFSError(err)
}
