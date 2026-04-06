// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
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
	audience := base.GetAudience((*base.ClientOptions)(options))
	conOptions := shared.GetClientOptions(options)
	authPolicy := shared.NewStorageChallengePolicy(cred, audience, conOptions.InsecureAllowCredentialWithHTTP)
	plOpts := runtime.PipelineOptions{
		PerCall:  []policy.Policy{shared.NewLayoutPolicy()},
		PerRetry: []policy.Policy{authPolicy},
	}

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
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
	plOpts := runtime.PipelineOptions{
		PerCall: []policy.Policy{shared.NewLayoutPolicy()},
	}

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
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
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	plOpts := runtime.PipelineOptions{
		PerCall:  []policy.Policy{shared.NewLayoutPolicy()},
		PerRetry: []policy.Policy{authPolicy},
	}

	azClient, err := azcore.NewClient(exported.ModuleName, exported.ModuleVersion, plOpts, &conOptions.ClientOptions)
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

// downloadBuffer downloads an Azure blob to a WriterAt in parallel.
func (b *Client) downloadBuffer(ctx context.Context, writer io.WriterAt, o downloadOptions, resizeFile func(int64) error) (int64, error) {
	if o.BlockSize == 0 {
		o.BlockSize = DefaultDownloadBlockSize
	}
	dataDownloaded := int64(0)
	computeReadLength := true
	count := o.Range.Count

	// TODO : SDK should ideally start with an initial download instead of get properties to optimize for small blobs.
	state := layoutState{
		ctx: ctx,
	}
	useLayout := o.EnableLayoutAwareRouting
	var l layout
	// If we don't have the length at all, get it
	var length int64
	var initialIfMatch *azcore.ETag
	// Try layout-aware routing first if enabled, otherwise use GetProperties
	if o.EnableLayoutAwareRouting {
		var err error
		l, _, err = getLayout(state, b.GetLayoutPager(o.getBlobLayoutOptions()))
		sc := bloberror.GetStatusCode(err)
		if err != nil {
			if sc == 400 || sc >= 500 { // fall back to old behavior if service doesn't support layout or layout wasn't fetched
				useLayout = false
			} else { // fail the operation
				return 0, err
			}
		} else {
			length = l.contentLength
			initialIfMatch = l.eTag
		}
	}

	if !useLayout {
		gr, err := b.GetProperties(ctx, o.getBlobPropertiesOptions())
		if err != nil {
			return 0, err
		}
		length = *gr.ContentLength
		initialIfMatch = gr.ETag
	}
	if len(l.layoutRanges) == 0 { // fall back to old behavior if layout doesn't exist
		useLayout = false
	}

	if count == CountToEnd { // If size not specified, calculate it
		count = length - o.Range.Offset
		dataDownloaded = count
		computeReadLength = false
	}

	if resizeFile != nil {
		if err := resizeFile(count); err != nil {
			return 0, err
		}
	}

	if count <= 0 {
		// The file is empty, there is nothing to download.
		return 0, nil
	}

	// If unspecified by the user, eTag lock on the initial call to ensure consistency of the blob through the download.
	if o.AccessConditions == nil {
		o.AccessConditions = &AccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: initialIfMatch,
			},
		}
	} else if o.AccessConditions.ModifiedAccessConditions == nil {
		o.AccessConditions.ModifiedAccessConditions = &ModifiedAccessConditions{
			IfMatch: initialIfMatch,
		}
	} else if o.AccessConditions.ModifiedAccessConditions.IfMatch == nil {
		o.AccessConditions.ModifiedAccessConditions.IfMatch = initialIfMatch
	}

	// Prepare and do parallel download.
	progress := int64(0)
	progressLock := &sync.Mutex{}

	err := shared.DoBatchTransfer(ctx, &shared.BatchTransferOptions{
		OperationName: "downloadBlobToWriterAt",
		TransferSize:  count,
		ChunkSize:     o.BlockSize,
		NumChunks:     uint64(((count - 1) / o.BlockSize) + 1),
		Concurrency:   o.Concurrency,
		Operation: func(ctx context.Context, chunkStart int64, count int64) error {
			// Fetch ideal endpoint for this chunk from layout
			if useLayout {
				endpoint := l.getIdealEndpoint(chunkStart + o.Range.Offset)
				if endpoint != "" {
					ctx = shared.WithLayoutEndpoint(ctx, endpoint)
				}
			}
			downloadBlobOptions := o.getDownloadBlobOptions(HTTPRange{
				Offset: chunkStart + o.Range.Offset,
				Count:  count,
			}, nil)
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
						progressLock.Lock()
						progress += diff
						o.Progress(progress)
						progressLock.Unlock()
					})
			}
			_, err = io.Copy(shared.NewSectionWriter(writer, chunkStart, count), body)
			if err != nil {
				return err
			}
			if computeReadLength {
				atomic.AddInt64(&dataDownloaded, *dr.ContentLength)
			}
			err = body.Close()
			return err
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

	dr, err := b.generated().Download(ctx, downloadOptions, leaseAccessConditions, cpkInfo, modifiedAccessConditions)
	if err != nil {
		return DownloadStreamResponse{}, err
	}

	return DownloadStreamResponse{
		client:                 b,
		DownloadResponse:       dr,
		getInfo:                httpGetterInfo{Range: o.Range, ETag: dr.ETag},
		ObjectReplicationRules: deserializeORSPolicies(dr.ObjectReplicationRules),
		cpkInfo:                o.CPKInfo,
		cpkScope:               o.CPKScopeInfo,
	}, err
}

// DownloadBuffer downloads an Azure blob to a buffer with parallel.
func (b *Client) DownloadBuffer(ctx context.Context, buffer []byte, o *DownloadBufferOptions) (int64, error) {
	if o == nil {
		o = &DownloadBufferOptions{}
	}
	return b.downloadBuffer(ctx, shared.NewBytesWriter(buffer), (downloadOptions)(*o), nil)
}

// DownloadFile downloads an Azure blob to a local file.
// The file would be truncated if the size doesn't match.
func (b *Client) DownloadFile(ctx context.Context, file *os.File, o *DownloadFileOptions) (int64, error) {
	if o == nil {
		o = &DownloadFileOptions{}
	}
	do := (*downloadOptions)(o)

	// Compare and try to resize local file's size if it doesn't match Azure blob's size.
	resizeFile := func(size int64) error {
		stat, err := file.Stat()
		if err != nil {
			return err
		}
		if stat.Size() != size {
			if err = file.Truncate(size); err != nil {
				return err
			}
		}
		return nil
	}
	return b.downloadBuffer(ctx, file, *do, resizeFile)
}

// GetLayoutPager returns the blob's layout.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob-layout.
func (b *Client) GetLayoutPager(options *GetLayoutOptions) *runtime.Pager[GetLayoutResponse] {
	opts, leaseAccessConditions, cpkInfo, modifiedAccessConditions := options.format()
	// Use user's IfMatch if provided, otherwise we'll capture the ETag from the initial response
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
				mac := modifiedAccessConditions
				if mac == nil {
					mac = &generated.ModifiedAccessConditions{}
				}
				mac.IfMatch = initialIfMatch
				req, err = b.generated().GetLayoutCreateRequest(ctx, opts, leaseAccessConditions, mac, cpkInfo)
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
