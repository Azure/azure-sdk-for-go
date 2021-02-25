// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/url"
	"time"
)

// A BlobClient represents a URL to an Azure Storage blob; the blob may be a block blob, append blob, or page blob.
type BlobClient struct {
	client *blobClient
	u      url.URL
	cred   StorageAccountCredential
}

// NewBlobClient creates a BlobClient object using the specified URL and request policy pipeline.
func NewBlobClient(blobURL string, cred azcore.Credential, options *ClientOptions) (BlobClient, error) {
	u, err := url.Parse(blobURL)
	if err != nil {
		return BlobClient{}, err
	}
	con := newConnection(blobURL, cred, options.getConnectionOptions())

	c, _ := cred.(*SharedKeyCredential)

	return BlobClient{client: &blobClient{con, nil}, u: *u, cred: c}, nil
}

// URL returns the URL endpoint used by the BlobClient object.
func (b BlobClient) URL() url.URL {
	return b.u
}

// string returns the URL as a string.
func (b BlobClient) string() string {
	u := b.URL()
	return u.String()
}

// WithSnapshot creates a new BlobClient object identical to the source but with the specified snapshot timestamp.
// Pass "" to remove the snapshot returning a URL to the base blob.
func (b BlobClient) WithSnapshot(snapshot string) BlobClient {
	p := NewBlobURLParts(b.URL())
	p.Snapshot = snapshot
	snapshotURL := p.URL()
	return BlobClient{
		client: &blobClient{
			newConnectionWithPipeline(snapshotURL.String(), b.client.con.p),
			b.client.pathRenameMode,
		},
		u: b.u,
	}
}

// ToAppendBlobClient creates an AppendBlobURL using the source's URL and pipeline.
func (b BlobClient) ToAppendBlobClient() AppendBlobClient {
	con := newConnectionWithPipeline(b.string(), b.client.con.p)
	return AppendBlobClient{
		client:     &appendBlobClient{con},
		u:          b.u,
		BlobClient: BlobClient{client: &blobClient{con: con}},
	}
}

// ToBlockBlobURL creates a BlockBlobClient using the source's URL and pipeline.
func (b BlobClient) ToBlockBlobClient() BlockBlobClient {
	con := newConnectionWithPipeline(b.string(), b.client.con.p)
	return BlockBlobClient{
		client:     &blockBlobClient{con},
		u:          b.u,
		BlobClient: BlobClient{client: &blobClient{con: con}},
	}
}

// ToPageBlobClient creates a PageBlobURL using the source's URL and pipeline.
func (b BlobClient) ToPageBlobClient() PageBlobClient {
	con := newConnectionWithPipeline(b.string(), b.client.con.p)
	return PageBlobClient{
		client:     &pageBlobClient{con},
		u:          b.u,
		BlobClient: BlobClient{client: &blobClient{con: con}},
	}
}

func (b BlobClient) GetAccountInfo(ctx context.Context) (BlobGetAccountInfoResponse, error) {
	return b.client.GetAccountInfo(ctx, nil)
}

// DownloadBlob reads a range of bytes from a blob. The response also includes the blob's properties and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob.
func (b BlobClient) Download(ctx context.Context, options *DownloadBlobOptions) (*DownloadResponse, error) {
	o, lease, cpk, accessConditions := options.pointers()
	dr, err := b.client.Download(ctx, o, lease, cpk, accessConditions)
	if err != nil {
		return nil, handleError(err)
	}

	offset := int64(0)
	count := int64(CountToEnd)

	if options != nil && options.Offset != nil {
		offset = *options.Offset
	}

	if options != nil && options.Count != nil {
		count = *options.Count
	}
	return &DownloadResponse{
		b:       b,
		BlobDownloadResponse:       dr,
		ctx:     ctx,
		getInfo: HTTPGetterInfo{Offset: offset, Count: count, ETag: *dr.ETag},
	}, err
}

// DeleteBlob marks the specified blob or snapshot for deletion. The blob is later deleted during garbage collection.
// Note that deleting a blob also deletes all its snapshots.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/delete-blob.
func (b BlobClient) Delete(ctx context.Context, options *DeleteBlobOptions) (BlobDeleteResponse, error) {
	basics, leaseInfo, accessConditions := options.pointers()
	resp, err := b.client.Delete(ctx, basics, leaseInfo, accessConditions)

	return resp, handleError(err)
}

// Undelete restores the contents and metadata of a soft-deleted blob and any associated soft-deleted snapshots.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/undelete-blob.
func (b BlobClient) Undelete(ctx context.Context) (BlobUndeleteResponse, error) {
	resp, err := b.client.Undelete(ctx, nil)

	return resp, handleError(err)
}

// SetTier operation sets the tier on a blob. The operation is allowed on a page
// blob in a premium storage account and on a block blob in a blob storage account (locally
// redundant storage only). A premium page blob's tier determines the allowed size, IOPS, and
// bandwidth of the blob. A block blob's tier determines Hot/Cool/Archive storage type. This operation
// does not update the blob's ETag.
// For detailed information about block blob level tiering see https://docs.microsoft.com/en-us/azure/storage/blobs/storage-blob-storage-tiers.
func (b BlobClient) SetTier(ctx context.Context, tier AccessTier, options *SetTierOptions) (BlobSetTierResponse, error) {
	basics, lease, accessConditions := options.pointers()
	resp, err := b.client.SetTier(ctx, tier, basics, lease, accessConditions)

	return resp, handleError(err)
}

// GetBlobProperties returns the blob's properties.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob-properties.
func (b BlobClient) GetProperties(ctx context.Context, options *GetBlobPropertiesOptions) (BlobGetPropertiesResponse, error) {
	basics, lease, cpk, access := options.pointers()
	resp, err := b.client.GetProperties(ctx, basics, lease, cpk, access)

	return resp, handleError(err)
}

// SetBlobHTTPHeaders changes a blob's HTTP headers.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-blob-properties.
func (b BlobClient) SetHTTPHeaders(ctx context.Context, blobHttpHeaders BlobHttpHeaders, options *SetBlobHTTPHeadersOptions) (BlobSetHTTPHeadersResponse, error) {
	basics, lease, access := options.pointers()
	resp, err := b.client.SetHTTPHeaders(ctx, basics, &blobHttpHeaders, lease, access)

	return resp, handleError(err)
}

// SetBlobMetadata changes a blob's metadata.
// https://docs.microsoft.com/rest/api/storageservices/set-blob-metadata.
func (b BlobClient) SetMetadata(ctx context.Context, metadata map[string]string, options *SetBlobMetadataOptions) (BlobSetMetadataResponse, error) {
	lease, cpk, cpkScope, access := options.pointers()
	basics := BlobSetMetadataOptions{
		Metadata: &metadata,
	}
	resp, err := b.client.SetMetadata(ctx, &basics, lease, cpk, cpkScope, access)

	return resp, handleError(err)
}

// CreateSnapshot creates a read-only snapshot of a blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/snapshot-blob.
func (b BlobClient) CreateSnapshot(ctx context.Context, options *CreateBlobSnapshotOptions) (BlobCreateSnapshotResponse, error) {
	// CreateSnapshot does NOT panic if the user tries to create a snapshot using a URL that already has a snapshot query parameter
	// because checking this would be a performance hit for a VERY unusual path and we don't think the common case should suffer this
	// performance hit.
	basics, cpk, cpkScope, access, lease := options.pointers()
	resp, err := b.client.CreateSnapshot(ctx, basics, cpk, cpkScope, access, lease)

	return resp, handleError(err)
}

// AcquireLease acquires a lease on the blob for write and delete operations. The lease duration must be between
// 15 to 60 seconds, or infinite (-1).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (b BlobClient) AcquireLease(ctx context.Context, options *AcquireBlobLeaseOptions) (BlobAcquireLeaseResponse, error) {
	basics, access := options.pointers()
	resp, err := b.client.AcquireLease(ctx, basics, access)

	return resp, handleError(err)
}

// RenewLease renews the blob's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (b BlobClient) RenewLease(ctx context.Context, leaseID string, options *RenewBlobLeaseOptions) (BlobRenewLeaseResponse, error) {
	basics, access := options.pointers()
	return b.client.RenewLease(ctx, leaseID, basics, access)
}

// ReleaseLease releases the blob's previously-acquired lease.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (b BlobClient) ReleaseLease(ctx context.Context, leaseID string, options *ReleaseBlobLeaseOptions) (BlobReleaseLeaseResponse, error) {
	basics, access := options.pointers()
	resp, err := b.client.ReleaseLease(ctx, leaseID, basics, access)

	return resp, handleError(err)
}

// BreakLease breaks the blob's previously-acquired lease (if it exists). Pass the LeaseBreakDefault (-1)
// constant to break a fixed-duration lease when it expires or an infinite lease immediately.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (b BlobClient) BreakLease(ctx context.Context, options *BreakBlobLeaseOptions) (BlobBreakLeaseResponse, error) {
	basics, access := options.pointers()
	resp, err := b.client.BreakLease(ctx, basics, access)

	return resp, handleError(err)
}

// ChangeLease changes the blob's lease ID.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/lease-blob.
func (b BlobClient) ChangeLease(ctx context.Context, leaseID string, proposedID string, options *ChangeBlobLeaseOptions) (BlobChangeLeaseResponse, error) {
	basics, access := options.pointers()
	resp, err := b.client.ChangeLease(ctx, leaseID, proposedID, basics, access)

	return resp, handleError(err)
}

// StartCopyFromURL copies the data at the source URL to a blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/copy-blob.
func (b BlobClient) StartCopyFromURL(ctx context.Context, copySource url.URL, options *StartCopyBlobOptions) (BlobStartCopyFromURLResponse, error) {
	basics, srcAccess, destAccess, lease := options.pointers()
	resp, err := b.client.StartCopyFromURL(ctx, copySource, basics, srcAccess, destAccess, lease)

	return resp, handleError(err)
}

// AbortCopyFromURL stops a pending copy that was previously started and leaves a destination blob with 0 length and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/abort-copy-blob.
func (b BlobClient) AbortCopyFromURL(ctx context.Context, copyID string, options *AbortCopyBlobOptions) (BlobAbortCopyFromURLResponse, error) {
	basics, lease := options.pointers()
	resp, err := b.client.AbortCopyFromURL(ctx, copyID, basics, lease)

	return resp, handleError(err)
}

// SetTags operation enables users to set tags on a blob or specific blob version, but not snapshot.
// Each call to this operation replaces all existing tags attached to the blob.
// To remove all tags from the blob, call this operation with no tags set.
// https://docs.microsoft.com/en-us/rest/api/storageservices/set-blob-tags
func (b BlobClient) SetTags(ctx context.Context, options *SetTagsBlobOptions) (BlobSetTagsResponse, error) {
	blobSetTagsOptions, modifiedAccessConditions := options.pointers()
	resp, err := b.client.SetTags(ctx, blobSetTagsOptions, modifiedAccessConditions)

	return resp, handleError(err)
}

// GetTags operation enables users to get tags on a blob or specific blob version, or snapshot.
// https://docs.microsoft.com/en-us/rest/api/storageservices/get-blob-tags
func (b BlobClient) GetTags(ctx context.Context, options *GetTagsBlobOptions) (BlobTagsResponse, error) {
	blobGetTagsOptions, modifiedAccessConditions := options.pointers()
	resp, err := b.client.GetTags(ctx, blobGetTagsOptions, modifiedAccessConditions)

	return resp, handleError(err)

}

func (b BlobClient) CanGetBlobSASToken() bool {
	return b.cred != nil
}

// GetBlobSASToken is a convenience method for generating a SAS token for the currently pointed at blob.
// It can only be used if the supplied azcore.Credential during creation was a SharedKeyCredential.
// This validity can be checked with CanGetBlobSASToken().
func (b BlobClient) GetBlobSASToken(permissions BlobSASPermissions, validityTime time.Duration) (SASQueryParameters, error) {
	urlParts := NewBlobURLParts(b.URL())

	t, err := time.Parse(SnapshotTimeFormat, urlParts.Snapshot)

	if err != nil {
		t = time.Time{}
	}

	return BlobSASSignatureValues{
		ContainerName: urlParts.ContainerName,
		BlobName: urlParts.BlobName,
		SnapshotTime: t,
		Version: SASVersion,

		Permissions: permissions.String(),

		StartTime: time.Now(),
		ExpiryTime: time.Now().Add(validityTime),
	}.NewSASQueryParameters(b.cred)
}
