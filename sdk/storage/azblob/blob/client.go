// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/temporal"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to an Azure Storage blob; the blob may be a block blob, append blob, or page blob.
type Client base.Client[generated.BlobClient]

// NewClient creates an instance of Client with the specified values.
//   - blobURL - the URL of the blob e.g. https://<account>.blob.core.windows.net/container/blob.txt
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options; pass nil to accept the default values
func NewClient(blobURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	conOptions := shared.GetClientOptions(options)

	azClient, err := base.GetAzClient(blobURL, cred, nil, (*base.ClientOptions)(conOptions))
	if err != nil {
		return nil, err
	}
	return (*Client)(base.NewBlobClient(blobURL, azClient, &cred, (*base.ClientOptions)(conOptions))), nil
}

// NewClientWithNoCredential creates an instance of Client with the specified values.
// This is used to anonymously access a blob or with a shared access signature (SAS) token.
//   - blobURL - the URL of the blob e.g. https://<account>.blob.core.windows.net/container/blob.txt?<sas token>
//   - options - client options; pass nil to accept the default values
func NewClientWithNoCredential(blobURL string, options *ClientOptions) (*Client, error) {
	conOptions := shared.GetClientOptions(options)

	azClient, err := base.GetAzClient(blobURL, nil, nil, (*base.ClientOptions)(conOptions))
	if err != nil {
		return nil, err
	}
	return (*Client)(base.NewBlobClient(blobURL, azClient, nil, (*base.ClientOptions)(conOptions))), nil
}

// NewClientWithSharedKeyCredential creates an instance of Client with the specified values.
//   - blobURL - the URL of the blob e.g. https://<account>.blob.core.windows.net/container/blob.txt
//   - cred - a SharedKeyCredential created with the matching blob's storage account and access key
//   - options - client options; pass nil to accept the default values
func NewClientWithSharedKeyCredential(blobURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	conOptions := shared.GetClientOptions(options)

	azClient, err := base.GetAzClient(blobURL, nil, cred, (*base.ClientOptions)(conOptions))
	if err != nil {
		return nil, err
	}
	return (*Client)(base.NewBlobClient(blobURL, azClient, cred, (*base.ClientOptions)(conOptions))), nil
}

// NewClientFromConnectionString creates an instance of Client with the specified values.
//   - connectionString - a connection string for the desired storage account
//   - containerName - the name of the container within the storage account
//   - blobName - the name of the blob within the container
//   - options - client options; pass nil to accept the default values
func NewClientFromConnectionString(connectionString, containerName, blobName string, options *ClientOptions) (*Client, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	parsed.ServiceURL = runtime.JoinPaths(parsed.ServiceURL, containerName, blobName)

	if parsed.AccountKey != "" && parsed.AccountName != "" {
		credential, err := exported.NewSharedKeyCredential(parsed.AccountName, parsed.AccountKey)
		if err != nil {
			return nil, err
		}
		return NewClientWithSharedKeyCredential(parsed.ServiceURL, credential, options)
	}

	return NewClientWithNoCredential(parsed.ServiceURL, options)
}

func (b *Client) generated() *generated.BlobClient {
	return base.InnerClient((*base.Client[generated.BlobClient])(b))
}

func (b *Client) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.BlobClient])(b))
}

func (b *Client) credential() any {
	return base.Credential((*base.Client[generated.BlobClient])(b))
}

func (b *Client) getClientOptions() *base.ClientOptions {
	return base.GetClientOptions((*base.Client[generated.BlobClient])(b))
}

// URL returns the URL endpoint used by the Client object.
func (b *Client) URL() string {
	return b.generated().Endpoint()
}

// WithSnapshot creates a new Client object identical to the source but with the specified snapshot timestamp.
// Pass "" to remove the snapshot returning a URL to the base blob.
func (b *Client) WithSnapshot(snapshot string) (*Client, error) {
	p, err := ParseURL(b.URL())
	if err != nil {
		return nil, err
	}
	p.Snapshot = snapshot

	return (*Client)(base.NewBlobClient(p.String(), b.generated().InternalClient(), b.credential(), b.getClientOptions())), nil
}

// WithVersionID creates a new AppendBlobURL object identical to the source but with the specified version id.
// Pass "" to remove the versionID returning a URL to the base blob.
func (b *Client) WithVersionID(versionID string) (*Client, error) {
	p, err := ParseURL(b.URL())
	if err != nil {
		return nil, err
	}
	p.VersionID = versionID

	return (*Client)(base.NewBlobClient(p.String(), b.generated().InternalClient(), b.credential(), b.getClientOptions())), nil
}

// Delete marks the specified blob or snapshot for deletion. The blob is later deleted during garbage collection.
// Note that deleting a blob also deletes all its snapshots.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/delete-blob.
func (b *Client) Delete(ctx context.Context, o *DeleteOptions) (DeleteResponse, error) {
	deleteOptions, leaseInfo, accessConditions := o.format()
	resp, err := b.generated().Delete(ctx, deleteOptions, leaseInfo, accessConditions)
	return resp, err
}

// Undelete restores the contents and metadata of a soft-deleted blob and any associated soft-deleted snapshots.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/undelete-blob.
func (b *Client) Undelete(ctx context.Context, o *UndeleteOptions) (UndeleteResponse, error) {
	undeleteOptions := o.format()
	resp, err := b.generated().Undelete(ctx, undeleteOptions)
	return resp, err
}

// SetTier operation sets the tier on a blob. The operation is allowed on a page
// blob in a premium storage account and on a block blob in a blob storage account (locally
// redundant storage only). A premium page blob's tier determines the allowed size, IOPs, and
// bandwidth of the blob. A block blob's tier determines Hot/Cool/Archive storage type. This operation
// does not update the blob's ETag.
// For detailed information about block blob level tiers see https://docs.microsoft.com/en-us/azure/storage/blobs/storage-blob-storage-tiers.
func (b *Client) SetTier(ctx context.Context, tier AccessTier, o *SetTierOptions) (SetTierResponse, error) {
	opts, leaseAccessConditions, modifiedAccessConditions := o.format()
	resp, err := b.generated().SetTier(ctx, tier, opts, leaseAccessConditions, modifiedAccessConditions)
	return resp, err
}

// GetProperties returns the blob's properties.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob-properties.
func (b *Client) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	opts, leaseAccessConditions, cpkInfo, modifiedAccessConditions := options.format()
	resp, err := b.generated().GetProperties(ctx, opts, leaseAccessConditions, cpkInfo, modifiedAccessConditions)
	return resp, err
}

// SetHTTPHeaders changes a blob's HTTP headers.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-blob-properties.
func (b *Client) SetHTTPHeaders(ctx context.Context, httpHeaders HTTPHeaders, o *SetHTTPHeadersOptions) (SetHTTPHeadersResponse, error) {
	opts, leaseAccessConditions, modifiedAccessConditions := o.format()
	resp, err := b.generated().SetHTTPHeaders(ctx, opts, &httpHeaders, leaseAccessConditions, modifiedAccessConditions)
	return resp, err
}

// SetMetadata changes a blob's metadata.
// https://docs.microsoft.com/rest/api/storageservices/set-blob-metadata.
func (b *Client) SetMetadata(ctx context.Context, metadata map[string]*string, o *SetMetadataOptions) (SetMetadataResponse, error) {
	basics := generated.BlobClientSetMetadataOptions{Metadata: metadata}
	leaseAccessConditions, cpkInfo, cpkScope, modifiedAccessConditions := o.format()
	resp, err := b.generated().SetMetadata(ctx, &basics, leaseAccessConditions, cpkInfo, cpkScope, modifiedAccessConditions)
	return resp, err
}

// CreateSnapshot creates a read-only snapshot of a blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/snapshot-blob.
func (b *Client) CreateSnapshot(ctx context.Context, options *CreateSnapshotOptions) (CreateSnapshotResponse, error) {
	// CreateSnapshot does NOT panic if the user tries to create a snapshot using a URL that already has a snapshot query parameter
	// because checking this would be a performance hit for a VERY unusual path, and we don't think the common case should suffer this
	// performance hit.
	opts, cpkInfo, cpkScope, modifiedAccessConditions, leaseAccessConditions := options.format()
	resp, err := b.generated().CreateSnapshot(ctx, opts, cpkInfo, cpkScope, modifiedAccessConditions, leaseAccessConditions)

	return resp, err
}

// StartCopyFromURL copies the data at the source URL to a blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/copy-blob.
func (b *Client) StartCopyFromURL(ctx context.Context, copySource string, options *StartCopyFromURLOptions) (StartCopyFromURLResponse, error) {
	opts, sourceModifiedAccessConditions, modifiedAccessConditions, leaseAccessConditions := options.format()
	resp, err := b.generated().StartCopyFromURL(ctx, copySource, opts, sourceModifiedAccessConditions, modifiedAccessConditions, leaseAccessConditions)
	return resp, err
}

// AbortCopyFromURL stops a pending copy that was previously started and leaves a destination blob with 0 length and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/abort-copy-blob.
func (b *Client) AbortCopyFromURL(ctx context.Context, copyID string, options *AbortCopyFromURLOptions) (AbortCopyFromURLResponse, error) {
	opts, leaseAccessConditions := options.format()
	resp, err := b.generated().AbortCopyFromURL(ctx, copyID, opts, leaseAccessConditions)
	return resp, err
}

// SetTags operation enables users to set tags on a blob or specific blob version, but not snapshot.
// Each call to this operation replaces all existing tags attached to the blob.
// To remove all tags from the blob, call this operation with no tags set.
// https://docs.microsoft.com/en-us/rest/api/storageservices/set-blob-tags
func (b *Client) SetTags(ctx context.Context, tags map[string]string, options *SetTagsOptions) (SetTagsResponse, error) {
	serializedTags := shared.SerializeBlobTags(tags)
	blobSetTagsOptions, modifiedAccessConditions, leaseAccessConditions, blobModifiedAccessConditions := options.format()
	resp, err := b.generated().SetTags(ctx, *serializedTags, blobSetTagsOptions, modifiedAccessConditions, leaseAccessConditions, blobModifiedAccessConditions)
	return resp, err
}

// GetTags operation enables users to get tags on a blob or specific blob version, or snapshot.
// https://docs.microsoft.com/en-us/rest/api/storageservices/get-blob-tags
func (b *Client) GetTags(ctx context.Context, options *GetTagsOptions) (GetTagsResponse, error) {
	blobGetTagsOptions, modifiedAccessConditions, leaseAccessConditions, blobModifiedAccessConditions := options.format()
	resp, err := b.generated().GetTags(ctx, blobGetTagsOptions, modifiedAccessConditions, leaseAccessConditions, blobModifiedAccessConditions)
	return resp, err

}

// SetImmutabilityPolicy operation enables users to set the immutability policy on a blob. Mode defaults to "Unlocked".
// https://learn.microsoft.com/en-us/azure/storage/blobs/immutable-storage-overview
func (b *Client) SetImmutabilityPolicy(ctx context.Context, expiryTime time.Time, options *SetImmutabilityPolicyOptions) (SetImmutabilityPolicyResponse, error) {
	blobSetImmutabilityPolicyOptions, modifiedAccessConditions := options.format()
	blobSetImmutabilityPolicyOptions.ImmutabilityPolicyExpiry = &expiryTime
	resp, err := b.generated().SetImmutabilityPolicy(ctx, blobSetImmutabilityPolicyOptions, modifiedAccessConditions)
	return resp, err
}

// DeleteImmutabilityPolicy operation enables users to delete the immutability policy on a blob.
// https://learn.microsoft.com/en-us/azure/storage/blobs/immutable-storage-overview
func (b *Client) DeleteImmutabilityPolicy(ctx context.Context, options *DeleteImmutabilityPolicyOptions) (DeleteImmutabilityPolicyResponse, error) {
	deleteImmutabilityOptions := options.format()
	resp, err := b.generated().DeleteImmutabilityPolicy(ctx, deleteImmutabilityOptions)
	return resp, err
}

// SetLegalHold operation enables users to set legal hold on a blob.
// https://learn.microsoft.com/en-us/azure/storage/blobs/immutable-storage-overview
func (b *Client) SetLegalHold(ctx context.Context, legalHold bool, options *SetLegalHoldOptions) (SetLegalHoldResponse, error) {
	setLegalHoldOptions := options.format()
	resp, err := b.generated().SetLegalHold(ctx, legalHold, setLegalHoldOptions)
	return resp, err
}

// CopyFromURL synchronously copies the data at the source URL to a block blob, with sizes up to 256 MB.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/copy-blob-from-url.
func (b *Client) CopyFromURL(ctx context.Context, copySource string, options *CopyFromURLOptions) (CopyFromURLResponse, error) {
	copyOptions, smac, mac, lac, cpkScopeInfo := options.format()
	resp, err := b.generated().CopyFromURL(ctx, copySource, copyOptions, smac, mac, lac, cpkScopeInfo)
	return resp, err
}

// GetAccountInfo provides account level information
// For more information, see https://learn.microsoft.com/en-us/rest/api/storageservices/get-account-information?tabs=shared-access-signatures.
func (b *Client) GetAccountInfo(ctx context.Context, o *GetAccountInfoOptions) (GetAccountInfoResponse, error) {
	getAccountInfoOptions := o.format()
	resp, err := b.generated().GetAccountInfo(ctx, getAccountInfoOptions)
	return resp, err
}

// GetSASURL is a convenience method for generating a SAS token for the currently pointed at blob.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
func (b *Client) GetSASURL(permissions sas.BlobPermissions, expiry time.Time, o *GetSASURLOptions) (string, error) {
	if b.sharedKey() == nil {
		return "", bloberror.MissingSharedKeyCredential
	}

	urlParts, err := ParseURL(b.URL())
	if err != nil {
		return "", err
	}

	t, err := time.Parse(SnapshotTimeFormat, urlParts.Snapshot)

	if err != nil {
		t = time.Time{}
	}
	st := o.format()

	qps, err := sas.BlobSignatureValues{
		ContainerName: urlParts.ContainerName,
		BlobName:      urlParts.BlobName,
		SnapshotTime:  t,
		Version:       sas.Version,
		Permissions:   permissions.String(),
		StartTime:     st,
		ExpiryTime:    expiry.UTC(),
	}.SignWithSharedKey(b.sharedKey())

	if err != nil {
		return "", err
	}

	endpoint := b.URL() + "?" + qps.Encode()

	return endpoint, nil
}

// Concurrent Download Functions -----------------------------------------------------------------------------------------

type downloadProgress struct {
	byteCount     int64
	byteCountLock sync.Mutex
}

// downloadBuffer downloads an Azure blob to a WriterAt in parallel.
//
// It opens with a ranged Get Blob rather than a Get Properties: that one response carries the
// blob's size, its ETag and the first chunk of data, so a blob that fits in a single block costs
// one round trip instead of two. The same response carries the download hint, which is what
// decides whether the remaining chunks are worth routing with the blob's layout.
func (b *Client) downloadBuffer(ctx context.Context, writer io.WriterAt, o downloadOptions) (int64, error) {
	if o.BlockSize == 0 {
		o.BlockSize = DefaultDownloadBlockSize
	}

	count := o.Range.Count
	layoutAware := o.layoutAwareRoutingEnabled()

	// When the caller gave a count and layout aware routing is off, there is nothing for an
	// initial request to discover, so the download goes straight to the parallel chunks. This
	// leaves the request pattern of an explicit-range download exactly as it was.
	if count != CountToEnd && !layoutAware {
		if count <= 0 {
			return 0, nil
		}
		return b.parallelDownloadFrom(ctx, writer, o, 0, count, &downloadProgress{}, nil)
	}

	// The initial request never reads past what the caller asked for.
	initialCount := o.BlockSize
	if count != CountToEnd {
		if count <= 0 {
			return 0, nil
		}
		if count < initialCount {
			initialCount = count
		}
	}

	dr, err := b.DownloadStream(ctx, o.getDownloadBlobOptions(HTTPRange{Offset: o.Range.Offset, Count: initialCount}, nil))
	if err != nil {
		if bloberror.HasCode(err, bloberror.InvalidRange) {
			// an empty blob has no range to read, so there is nothing to download
			return 0, nil
		}
		return 0, err
	}

	if dr.ContentRange == nil && dr.ContentLength == nil {
		// A 304 Not Modified response (from If-None-Match / If-Modified-Since conditions)
		// has no body or size headers. Close the body and surface as an error so callers
		// see the same behavior as the prior GetProperties path.
		_ = dr.Body.Close()
		return 0, fmt.Errorf("response contained no content headers; this may indicate a 304 Not Modified due to access conditions")
	}

	if count == CountToEnd {
		var totalSize int64
		if dr.ContentRange != nil {
			totalSize = parseContentRangeTotal(*dr.ContentRange)
			if totalSize <= 0 {
				_ = dr.Body.Close()
				return 0, fmt.Errorf("unable to parse total size from Content-Range header: %s", *dr.ContentRange)
			}
		} else {
			totalSize = *dr.ContentLength + o.Range.Offset
		}
		count = totalSize - o.Range.Offset
		if count <= 0 {
			_ = dr.Body.Close()
			return 0, nil
		}
	}

	// The initial response is the consistency anchor for everything that follows: the remaining
	// chunks, and the layout enumeration when layout aware routing is used. The caller's own
	// conditions were already applied to this request, so pinning the ETag it returned can only
	// narrow them, and it closes the gap a caller-supplied ETagAny would otherwise leave open,
	// where later chunks could come from a different version of the blob.
	if dr.ETag != nil {
		ac := &AccessConditions{}
		if o.AccessConditions != nil {
			clone := *o.AccessConditions
			ac = &clone
		}
		mac := &ModifiedAccessConditions{}
		if ac.ModifiedAccessConditions != nil {
			macClone := *ac.ModifiedAccessConditions
			mac = &macClone
		}
		mac.IfMatch = dr.ETag
		ac.ModifiedAccessConditions = mac
		o.AccessConditions = ac
	}

	var initialChunkSize int64
	if dr.ContentRange != nil {
		initialChunkSize = parseContentRangeLength(*dr.ContentRange)
	} else if dr.ContentLength != nil {
		initialChunkSize = *dr.ContentLength
	}
	if initialChunkSize <= 0 {
		_ = dr.Body.Close()
		return 0, nil
	}

	prog := &downloadProgress{}
	var body io.ReadCloser = dr.NewRetryReader(ctx, &o.RetryReaderOptionsPerBlock)
	if o.Progress != nil {
		body = streaming.NewResponseProgress(body, func(bytesTransferred int64) {
			prog.byteCountLock.Lock()
			prog.byteCount = bytesTransferred
			o.Progress(prog.byteCount)
			prog.byteCountLock.Unlock()
		})
	}
	if _, err = io.Copy(shared.NewSectionWriter(writer, 0, initialChunkSize), body); err != nil {
		_ = body.Close()
		return 0, err
	}
	if err = body.Close(); err != nil {
		return 0, err
	}

	initialDataLen := initialChunkSize
	if dr.StructuredBodyType != nil && *dr.StructuredBodyType != "" && dr.ContentRange != nil {
		// For structured message responses, ContentLength reflects the encoded size.
		// Use ContentRange to get the original data length.
		initialDataLen = parseContentRangeLength(*dr.ContentRange)
	}

	if initialChunkSize >= count {
		// The initial request returned everything that was asked for. Fetching a layout now
		// would cost a round trip for a download that is already finished, which is exactly what
		// starting with a Get Blob is meant to avoid.
		return initialDataLen, nil
	}

	// More data remains, so the layout decides where the remaining chunks are read from.
	var layoutResource *temporal.Resource[layout, context.Context]
	if layoutAware {
		layoutResource, err = b.resolveLayout(ctx, o)
		if err != nil {
			return 0, err
		}
	}

	remaining := count - initialChunkSize
	remainingDownloaded, err := b.parallelDownloadFrom(ctx, writer, o, initialChunkSize, remaining, prog, layoutResource)
	if err != nil {
		return 0, err
	}
	return initialDataLen + remainingDownloaded, nil
}

// resolveLayout enumerates the blob's layout and returns the cache holding it. It returns a nil
// resource when the layout cannot be used, in which case the remaining chunks go to the client's
// configured endpoint.
//
// The layout is resolved once here rather than on the first chunk so that a fallback or
// no-layout answer costs a single enumeration instead of one per chunk, and the resource is
// handed on to the chunk downloads so a transfer running past the layout's expiry refreshes it
// once rather than per chunk.
func (b *Client) resolveLayout(ctx context.Context, o downloadOptions) (*temporal.Resource[layout, context.Context], error) {
	resource := temporal.NewResourceWithOptions(
		func(ctx context.Context) (layout, time.Time, error) {
			return getLayout(ctx, b.GetLayoutPager(o.getBlobLayoutOptions()))
		}, temporal.ResourceOptions[layout, context.Context]{
			ShouldRefresh: shouldRefreshLayout,
		})

	l, err := resource.Get(ctx)
	if err != nil {
		// getLayout caches "layout unavailable" as a fallback layout, so any error here is fatal.
		return nil, err
	}
	if l.fallback {
		// The service can't provide a layout; fall back to the client's configured endpoint.
		return nil, nil
	}
	if len(l.layoutRanges) == 0 {
		// The blob has no layout: download everything from the primary endpoint.
		return nil, nil
	}
	return resource, nil
}

// parallelDownloadFrom downloads remaining bytes in parallel chunks, writing each one at
// writerOffset plus its offset within the requested range. When layoutResource is non-nil each
// chunk is routed to the endpoint its offset maps to in the cached layout.
func (b *Client) parallelDownloadFrom(ctx context.Context, writer io.WriterAt, o downloadOptions, writerOffset, remaining int64, prog *downloadProgress, layoutResource *temporal.Resource[layout, context.Context]) (int64, error) {
	dataDownloaded := int64(0)

	err := shared.DoBatchTransfer(ctx, &shared.BatchTransferOptions{
		OperationName: "downloadBlobToWriterAt",
		TransferSize:  remaining,
		ChunkSize:     o.BlockSize,
		NumChunks:     uint64(((remaining - 1) / o.BlockSize) + 1),
		Concurrency:   o.Concurrency,
		Operation: func(ctx context.Context, chunkStart int64, count int64) error {
			blobOffset := chunkStart + writerOffset + o.Range.Offset
			downloadBlobOptions := o.getDownloadBlobOptions(HTTPRange{
				Offset: blobOffset,
				Count:  count,
			}, nil)
			// Fetch ideal endpoint for this chunk from layout. A refresh that fails leaves the
			// chunk on the client's configured endpoint rather than failing the download.
			if layoutResource != nil {
				if chunkLayout, err := layoutResource.Get(ctx); err == nil && !chunkLayout.fallback {
					downloadBlobOptions.LayoutEndpoint = getIdealEndpoint(blobOffset, chunkLayout)
				}
			}
			dr, err := b.DownloadStream(ctx, downloadBlobOptions)
			if err != nil {
				return err
			}
			var body io.ReadCloser = dr.NewRetryReader(ctx, &o.RetryReaderOptionsPerBlock)
			if o.Progress != nil {
				rangeProgress := int64(0)
				body = streaming.NewResponseProgress(
					body,
					func(bytesTransferred int64) {
						diff := bytesTransferred - rangeProgress
						rangeProgress = bytesTransferred
						prog.byteCountLock.Lock()
						prog.byteCount += diff
						o.Progress(prog.byteCount)
						prog.byteCountLock.Unlock()
					})
			}
			if _, err = io.Copy(shared.NewSectionWriter(writer, chunkStart+writerOffset, count), body); err != nil {
				_ = body.Close()
				return err
			}
			if dr.StructuredBodyType != nil && *dr.StructuredBodyType != "" && dr.ContentRange != nil {
				// For structured message responses, ContentLength reflects the encoded size.
				// Use ContentRange to get the original data length.
				atomic.AddInt64(&dataDownloaded, parseContentRangeLength(*dr.ContentRange))
			} else {
				atomic.AddInt64(&dataDownloaded, *dr.ContentLength)
			}
			return body.Close()
		},
	})
	if err != nil {
		return 0, err
	}
	return dataDownloaded, nil
}

// DownloadStream reads a range of bytes from a blob. The response also includes the blob's properties and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob.
func (b *Client) DownloadStream(ctx context.Context, o *DownloadStreamOptions) (DownloadStreamResponse, error) {
	downloadOptions, leaseAccessConditions, cpkInfo, modifiedAccessConditions := o.format()
	if o == nil {
		o = &DownloadStreamOptions{}
	}
	if o.LayoutEndpoint != "" {
		ctx = shared.WithLayoutEndpoint(ctx, o.LayoutEndpoint)
	}

	dr, err := b.generated().Download(ctx, downloadOptions, leaseAccessConditions, cpkInfo, modifiedAccessConditions)
	if err != nil {
		return DownloadStreamResponse{}, err
	}

	// If the response contains a structured message body, wrap it with SMDecoder
	// to validate CRC64 checksums and extract the raw data.
	if dr.StructuredBodyType != nil && *dr.StructuredBodyType != "" {
		dr.Body = shared.NewSMDecoder(dr.Body)
	}

	return DownloadStreamResponse{
		client:                  b,
		DownloadResponse:        dr,
		getInfo:                 httpGetterInfo{Range: o.Range, ETag: dr.ETag},
		ObjectReplicationRules:  deserializeORSPolicies(dr.ObjectReplicationRules),
		cpkInfo:                 o.CPKInfo,
		cpkScope:                o.CPKScopeInfo,
		transactionalValidation: o.TransactionalValidation,
		// carried so a retry of this read stays on the endpoint the layout chose for it
		layoutEndpoint: o.LayoutEndpoint,
	}, err
}

// DownloadBuffer downloads an Azure blob to a buffer with parallel.
func (b *Client) DownloadBuffer(ctx context.Context, buffer []byte, o *DownloadBufferOptions) (int64, error) {
	if o == nil {
		o = &DownloadBufferOptions{}
	}
	return b.downloadBuffer(ctx, shared.NewBytesWriter(buffer), (downloadOptions)(*o))
}

// DownloadFile downloads an Azure blob to a local file.
// The file would be truncated if the size doesn't match.
func (b *Client) DownloadFile(ctx context.Context, file *os.File, o *DownloadFileOptions) (int64, error) {
	if o == nil {
		o = &DownloadFileOptions{}
	}
	do := (*downloadOptions)(o)

	downloaded, err := b.downloadBuffer(ctx, file, *do)
	if err != nil {
		return 0, err
	}

	// Compare and try to resize the local file's size if it doesn't match what was downloaded.
	stat, err := file.Stat()
	if err != nil {
		return downloaded, err
	}
	if stat.Size() != downloaded {
		if err = file.Truncate(downloaded); err != nil {
			return downloaded, err
		}
	}

	return downloaded, nil
}

// GetLayoutPager returns the blob's layout: the set of byte ranges making up the blob and the
// storage endpoint that serves each one. Pass the endpoint covering a given offset as
// DownloadStreamOptions.LayoutEndpoint to route that read for better locality.
//
// A single enumeration describes the whole blob, so callers implementing a custom chunked
// download should enumerate once and reuse the result across chunks rather than paging per
// chunk. A blob's layout can change over time; refresh the cached layout roughly every
// 5 minutes, which is the interval Client.DownloadBuffer and Client.DownloadFile use internally.
//
// Unless the caller supplies an If-Match condition in options.AccessConditions, the ETag returned
// by the first page is sent as If-Match on every subsequent page, so that a single enumeration
// always describes one version of the blob.
//
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob-layout.
func (b *Client) GetLayoutPager(options *GetLayoutOptions) *runtime.Pager[GetLayoutResponse] {
	opts, leaseAccessConditions, cpkInfo, modifiedAccessConditions := options.format()
	if opts == nil {
		opts = &generated.BlobClientGetLayoutOptions{}
	}
	// Use the caller's IfMatch if provided, otherwise capture the ETag from the initial response so
	// that every subsequent page is locked to the same version of the blob. The caller's
	// ModifiedAccessConditions are never mutated; a copy is made when the lock is applied.
	var initialIfMatch *azcore.ETag
	if modifiedAccessConditions != nil {
		initialIfMatch = modifiedAccessConditions.IfMatch
	}
	return runtime.NewPager(runtime.PagingHandler[GetLayoutResponse]{
		More: func(page GetLayoutResponse) bool {
			return page.NextMarker != nil && len(*page.NextMarker) > 0
		},
		Fetcher: func(ctx context.Context, page *GetLayoutResponse) (GetLayoutResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = b.generated().GetLayoutCreateRequest(ctx, opts, leaseAccessConditions, modifiedAccessConditions, cpkInfo)
			} else {
				opts.Marker = page.NextMarker
				// Use the ETag to ensure consistency across all pages
				mac := generated.ModifiedAccessConditions{}
				if modifiedAccessConditions != nil {
					mac = *modifiedAccessConditions
				}
				mac.IfMatch = initialIfMatch
				req, err = b.generated().GetLayoutCreateRequest(ctx, opts, leaseAccessConditions, &mac, cpkInfo)
			}
			if err != nil {
				return GetLayoutResponse{}, err
			}
			resp, err := b.generated().InternalClient().Pipeline().Do(req)
			if err != nil {
				return GetLayoutResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
				return GetLayoutResponse{}, runtime.NewResponseError(resp)
			}
			result, err := b.generated().GetLayoutHandleResponse(resp)
			if err != nil {
				return GetLayoutResponse{}, err
			}
			// Capture the ETag from the initial response for all subsequent requests
			if page == nil && initialIfMatch == nil {
				initialIfMatch = result.ETag
			}
			return result, nil
		},
	})
}

// parseContentRangeLength parses the range length from a Content-Range header value.
// Format: "bytes start-end/total" → returns end - start + 1.
// Returns 0 if the header cannot be parsed.
func parseContentRangeLength(contentRange string) int64 {
	var start, end int64
	if _, err := fmt.Sscanf(contentRange, "bytes %d-%d/", &start, &end); err != nil {
		return 0
	}
	return end - start + 1
}

// parseContentRangeTotal parses the blob's total size from a Content-Range header value.
// Format: "bytes start-end/total" → returns total.
// Returns 0 if the header cannot be parsed.
func parseContentRangeTotal(contentRange string) int64 {
	var start, end, total int64
	if _, err := fmt.Sscanf(contentRange, "bytes %d-%d/%d", &start, &end, &total); err != nil {
		return 0
	}
	return total
}
