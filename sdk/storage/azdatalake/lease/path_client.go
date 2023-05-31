package lease

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/path"
)

// PathClient provides lease functionality for the underlying blob client.
type PathClient struct {
	blobClient *lease.BlobClient
	pathClient *path.Client
	leaseID    *string
}

// PathClientOptions contains the optional values when creating a BlobClient.
type PathClientOptions struct {
	// LeaseID contains a caller-provided lease ID.
	LeaseID *string
}

// NewPathClient creates a blob lease client for the provided blob client.
//   - client - an instance of a blob client
//   - options - client options; pass nil to accept the default values
func NewPathClient[T directory.Client | file.Client](client *T, options *PathClientOptions) (*PathClient, error) {
	// TODO: set up blob lease client
	return nil, nil
}

func (c *PathClient) generated() *generated.PathClient {
	return base.InnerClient((*base.Client[generated.PathClient])(c.pathClient))
}

// LeaseID returns leaseID of the client.
func (c *PathClient) LeaseID() *string {
	return c.leaseID
}

// AcquireLease acquires a lease on the blob for write and delete operations.
// The lease Duration must be between 15 and 60 seconds, or infinite (-1).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *PathClient) AcquireLease(ctx context.Context, duration int32, o *PathAcquireOptions) (PathAcquireResponse, error) {
	opts := o.format()
	return c.blobClient.AcquireLease(ctx, duration, opts)
}

func (c *PathClient) BreakLease(ctx context.Context, o *PathBreakOptions) (PathBreakResponse, error) {
	opts := o.format()
	return c.blobClient.BreakLease(ctx, opts)
}

// ChangeLease changes the blob's lease ID.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *PathClient) ChangeLease(ctx context.Context, proposedID string, o *PathChangeOptions) (PathChangeResponse, error) {
	opts := o.format()
	return c.blobClient.ChangeLease(ctx, proposedID, opts)
}

// RenewLease renews the blob's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *PathClient) RenewLease(ctx context.Context, o *PathRenewOptions) (PathRenewResponse, error) {
	opts := o.format()
	return c.blobClient.RenewLease(ctx, opts)
}

// ReleaseLease releases the blob's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (c *PathClient) ReleaseLease(ctx context.Context, o *PathReleaseOptions) (PathReleaseResponse, error) {
	opts := o.format()
	return c.blobClient.ReleaseLease(ctx, opts)
}
