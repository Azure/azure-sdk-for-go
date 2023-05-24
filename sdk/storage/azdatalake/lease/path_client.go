//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

//
//import (
//	"context"
//	"errors"
//	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/directory"
//	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"
//	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
//	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/path"
//)
//
//// PathClient provides lease functionality for the underlying blob client.
//type PathClient struct {
//	pathClient *path.Client
//	leaseID    *string
//}
//
//// PathClientOptions contains the optional values when creating a BlobClient.
//type PathClientOptions struct {
//	// LeaseID contains a caller-provided lease ID.
//	LeaseID *string
//}
//
//// NewPathClient creates a blob lease client for the provided blob client.
////   - client - an instance of a blob client
////   - options - client options; pass nil to accept the default values
//func NewPathClient[T file.Client | directory.Client](client *T, options *PathClientOptions) (*PathClient, error) {
//	return &PathClient{}, nil
//}
//
//func (c *PathClient) generated() *generated.PathClient {
//	return &generated.PathClient{}
//}
//
//// LeaseID returns leaseID of the client.
//func (c *PathClient) LeaseID() *string {
//	return c.leaseID
//}
//
//// AcquireLease acquires a lease on the blob for write and delete operations.
//// The lease Duration must be between 15 and 60 seconds, or infinite (-1).
//// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
//func (c *PathClient) AcquireLease(ctx context.Context, duration int32, o *BlobAcquireOptions) (BlobAcquireResponse, error) {
//	blobAcquireLeaseOptions, modifiedAccessConditions := o.format()
//	blobAcquireLeaseOptions.ProposedLeaseID = c.LeaseID()
//
//	resp, err := c.generated().Lease(ctx, duration, &blobAcquireLeaseOptions, modifiedAccessConditions)
//	return resp, err
//}
//
//// BreakLease breaks the blob's previously-acquired lease (if it exists). Pass the LeaseBreakDefault (-1)
//// constant to break a fixed-Duration lease when it expires or an infinite lease immediately.
//// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
//func (c *PathClient) BreakLease(ctx context.Context, o *BlobBreakOptions) (BlobBreakResponse, error) {
//	blobBreakLeaseOptions, modifiedAccessConditions := o.format()
//	resp, err := c.generated().BreakLease(ctx, blobBreakLeaseOptions, modifiedAccessConditions)
//	return resp, err
//}
//
//// ChangeLease changes the blob's lease ID.
//// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
//func (c *BlobClient) ChangeLease(ctx context.Context, proposedLeaseID string, o *BlobChangeOptions) (BlobChangeResponse, error) {
//	if c.LeaseID() == nil {
//		return BlobChangeResponse{}, errors.New("leaseID cannot be nil")
//	}
//	changeLeaseOptions, modifiedAccessConditions, err := o.format()
//	if err != nil {
//		return BlobChangeResponse{}, err
//	}
//	resp, err := c.generated().ChangeLease(ctx, *c.LeaseID(), proposedLeaseID, changeLeaseOptions, modifiedAccessConditions)
//
//	// If lease has been changed successfully, set the leaseID in client
//	if err == nil {
//		c.leaseID = &proposedLeaseID
//	}
//
//	return resp, err
//}
//
//// RenewLease renews the blob's previously-acquired lease.
//// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
//func (c *BlobClient) RenewLease(ctx context.Context, o *BlobRenewOptions) (BlobRenewResponse, error) {
//	if c.LeaseID() == nil {
//		return BlobRenewResponse{}, errors.New("leaseID cannot be nil")
//	}
//	renewLeaseBlobOptions, modifiedAccessConditions := o.format()
//	resp, err := c.generated().RenewLease(ctx, *c.LeaseID(), renewLeaseBlobOptions, modifiedAccessConditions)
//	return resp, err
//}
//
//// ReleaseLease releases the blob's previously-acquired lease.
//// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
//func (c *BlobClient) ReleaseLease(ctx context.Context, o *BlobReleaseOptions) (BlobReleaseResponse, error) {
//	if c.LeaseID() == nil {
//		return BlobReleaseResponse{}, errors.New("leaseID cannot be nil")
//	}
//	renewLeaseBlobOptions, modifiedAccessConditions := o.format()
//	resp, err := c.generated().ReleaseLease(ctx, *c.LeaseID(), renewLeaseBlobOptions, modifiedAccessConditions)
//	return resp, err
//}
