package azblob

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/url"
	"time"
)

type BlobLeaseClient struct {
	BlobClient
	leaseId      *string
	lastModified *time.Time
	eTag         *string
}

func NewBlobLeaseClient(blobURL string, cred azcore.Credential, pathRenameMode *PathRenameMode, options *connectionOptions, leaseId *string) (BlobLeaseClient, error) {
	u, err := url.Parse(blobURL)
	if err != nil {
		return BlobLeaseClient{}, err
	}
	con := newConnection(blobURL, cred, options)
	blobClient := BlobClient{client: &blobClient{con, pathRenameMode}, u: *u}
	return BlobLeaseClient{BlobClient: blobClient, leaseId: leaseId}, nil
}

// URL returns the URL endpoint used by the BlobLeaseClient object.
func (blc BlobLeaseClient) URL() url.URL {
	return blc.u
}

// String returns the URL as a string.
func (blc BlobLeaseClient) String() string {
	u := blc.URL()
	return u.String()
}

// WithPipeline creates a new BlobLeaseClient object identical to the source but with the specified request policy pipeline.
func (blc BlobLeaseClient) WithPipeline(pipeline azcore.Pipeline) BlobLeaseClient {
	con := newConnectionWithPipeline(blc.u.String(), pipeline)
	blobClient := BlobClient{client: &blobClient{con, blc.client.pathRenameMode}, u: blc.u}
	return BlobLeaseClient{BlobClient: blobClient}
}

// AcquireLease acquires a lease on the blob for write and delete operations. The lease Duration must be between
// 15 to 60 seconds, or infinite (-1).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (blc *BlobLeaseClient) AcquireLease(ctx context.Context, options *AcquireLeaseBlobOptions) (BlobAcquireLeaseResponse, error) {
	blobAcquireLeaseOptions, modifiedAccessConditions := options.pointers()
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
	leaseId, proposedLeaseId, modifiedAccessConditions := options.pointers()
	resp, err := blc.client.ChangeLease(ctx, leaseId, proposedLeaseId, nil, modifiedAccessConditions)
	return resp, handleError(err)
}

// RenewLease renews the blob's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (blc *BlobLeaseClient) RenewLease(ctx context.Context, options *RenewLeaseBlobOptions) (BlobRenewLeaseResponse, error) {
	leaseId, renewLeaseBlobOptions, modifiedAccessConditions := options.pointers()
	resp, err := blc.client.RenewLease(ctx, leaseId, renewLeaseBlobOptions, modifiedAccessConditions)
	return resp, handleError(err)
}

// ReleaseLease releases the blob's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (blc *BlobLeaseClient) ReleaseLease(ctx context.Context, options *ReleaseLeaseBlobOptions) (BlobReleaseLeaseResponse, error) {
	leaseId, renewLeaseBlobOptions, modifiedAccessConditions := options.pointers()
	resp, err := blc.client.ReleaseLease(ctx, leaseId, renewLeaseBlobOptions, modifiedAccessConditions)
	return resp, handleError(err)
}

type ContainerLeaseClient struct {
	ContainerClient
	leaseId      *string
	lastModified *time.Time
	eTag         *string
}

func NewContainerLeaseClient(containerURL string, cred azcore.Credential, options *connectionOptions, leaseId *string) (ContainerLeaseClient, error) {
	u, err := url.Parse(containerURL)
	if err != nil {
		return ContainerLeaseClient{}, err
	}
	containerClient := ContainerClient{
		client: &containerClient{
			con: newConnection(containerURL, cred, options),
		}, u: *u,
	}
	return ContainerLeaseClient{
		ContainerClient: containerClient,
		leaseId:         leaseId,
	}, nil
}

// URL returns the URL endpoint used by the ContainerClient object.
func (clc ContainerLeaseClient) URL() url.URL {
	return clc.u
}

// String returns the URL as a string.
func (clc ContainerLeaseClient) String() string {
	u := clc.URL()
	return u.String()
}

// WithPipeline creates a new ContainerLeaseClient object identical to the source but with the specified request policy pipeline.
func (clc ContainerLeaseClient) WithPipeline(pipeline azcore.Pipeline) ContainerLeaseClient {
	con := newConnectionWithPipeline(clc.u.String(), pipeline)
	containerClient := ContainerClient{client: &containerClient{con: con}, u: clc.u}
	return ContainerLeaseClient{
		ContainerClient: containerClient,
	}
}

// AcquireLease acquires a lease on the container for delete operations. The lease Duration must be between 15 to 60 seconds, or infinite (-1).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (clc *ContainerLeaseClient) AcquireLease(ctx context.Context, options *AcquireLeaseContainerOptions) (ContainerAcquireLeaseResponse, error) {
	containerAcquireLeaseOptions, modifiedAccessConditions := options.pointers()
	resp, err := clc.client.AcquireLease(ctx, containerAcquireLeaseOptions, modifiedAccessConditions)
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
	leaseId, proposedLeaseId, modifiedAccessConditions := options.pointers()
	resp, err := clc.client.ChangeLease(ctx, leaseId, proposedLeaseId, nil, modifiedAccessConditions)
	return resp, handleError(err)
}

// ReleaseLease releases the container's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (clc *ContainerLeaseClient) ReleaseLease(ctx context.Context, options *ReleaseLeaseContainerOptions) (ContainerReleaseLeaseResponse, error) {
	leaseId, containerReleaseLeaseOptions, modifiedAccessConditions := options.pointers()
	resp, err := clc.client.ReleaseLease(ctx, leaseId, containerReleaseLeaseOptions, modifiedAccessConditions)
	return resp, handleError(err)
}

// RenewLease renews the container's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-container.
func (clc *ContainerLeaseClient) RenewLease(ctx context.Context, options *RenewLeaseContainerOptions) (ContainerRenewLeaseResponse, error) {
	leaseId, renewLeaseBlobOptions, modifiedAccessConditions := options.pointers()
	resp, err := clc.client.RenewLease(ctx, leaseId, renewLeaseBlobOptions, modifiedAccessConditions)
	return resp, handleError(err)
}
