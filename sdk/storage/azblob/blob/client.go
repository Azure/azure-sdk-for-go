// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
	"fmt"
	"io"
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
	plOpts := runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}
	if p := base.NewExpectContinuePolicy(conOptions.ExpectContinueBehavior); p != nil {
		plOpts.PerRetry = append(plOpts.PerRetry, p)
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
	plOpts := runtime.PipelineOptions{}
	if p := base.NewExpectContinuePolicy(conOptions.ExpectContinueBehavior); p != nil {
		plOpts.PerRetry = append(plOpts.PerRetry, p)
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
	plOpts := runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}
	if p := base.NewExpectContinuePolicy(conOptions.ExpectContinueBehavior); p != nil {
		plOpts.PerRetry = append(plOpts.PerRetry, p)
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

type downloadProgress struct {
	byteCount     int64
	byteCountLock sync.Mutex
}

// downloadBuffer downloads an Azure blob to a WriterAt in parallel.
// It uses an initial GetBlob (GET) request instead of GetProperties (HEAD) to determine the blob size,
// eliminating an extra round trip for small blobs where the entire content is returned in the initial response.
func (b *Client) downloadBuffer(ctx context.Context, writer io.WriterAt, o downloadOptions) (int64, error) {
	if o.BlockSize == 0 {
		o.BlockSize = DefaultDownloadBlockSize
	}

	count := o.Range.Count
	if count != CountToEnd {
		return b.parallelDownload(ctx, writer, o, count)
	}

	dr, err := b.DownloadStream(ctx, o.getDownloadBlobOptions(HTTPRange{Offset: o.Range.Offset, Count: o.BlockSize}, nil))
	if err != nil {
		if bloberror.HasCode(err, bloberror.InvalidRange) {
			return 0, nil
		}
		return 0, err
	}

	var totalSize int64
	if dr.ContentRange != nil {
		totalSize = parseContentRangeTotal(*dr.ContentRange)
		if totalSize <= 0 {
			_ = dr.Body.Close()
			return 0, fmt.Errorf("unable to parse total size from Content-Range header: %s", *dr.ContentRange)
		}
	} else if dr.ContentLength != nil {
		totalSize = *dr.ContentLength + o.Range.Offset
	} else {
		_ = dr.Body.Close()
		return 0, fmt.Errorf("response missing both Content-Range and Content-Length headers")
	}

	count = totalSize - o.Range.Offset
	if count <= 0 {
		_ = dr.Body.Close()
		return 0, nil
	}

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
	_, err = io.Copy(shared.NewSectionWriter(writer, 0, initialChunkSize), body)
	if err != nil {
		_ = body.Close()
		return 0, err
	}
	if err = body.Close(); err != nil {
		return 0, err
	}

	initialDataLen := initialChunkSize
	if dr.StructuredBodyType != nil && *dr.StructuredBodyType != "" && dr.ContentRange != nil {
		initialDataLen = parseContentRangeLength(*dr.ContentRange)
	}

	if initialChunkSize >= count {
		return initialDataLen, nil
	}

	remaining := count - initialChunkSize
	remainingDownloaded, err := b.parallelDownloadFrom(ctx, writer, o, initialChunkSize, remaining, prog)
	if err != nil {
		return 0, err
	}
	return initialDataLen + remainingDownloaded, nil
}

func (b *Client) parallelDownload(ctx context.Context, writer io.WriterAt, o downloadOptions, count int64) (int64, error) {
	dataDownloaded := int64(0)
	prog := &downloadProgress{}

	err := shared.DoBatchTransfer(ctx, &shared.BatchTransferOptions{
		OperationName: "downloadBlobToWriterAt",
		TransferSize:  count,
		ChunkSize:     o.BlockSize,
		NumChunks:     uint64(((count - 1) / o.BlockSize) + 1),
		Concurrency:   o.Concurrency,
		Operation: func(ctx context.Context, chunkStart int64, count int64) error {
			dr, err := b.DownloadStream(ctx, o.getDownloadBlobOptions(HTTPRange{
				Offset: chunkStart + o.Range.Offset, Count: count}, nil))
			if err != nil {
				return err
			}
			var body io.ReadCloser = dr.NewRetryReader(ctx, &o.RetryReaderOptionsPerBlock)
			if o.Progress != nil {
				rangeProgress := int64(0)
				body = streaming.NewResponseProgress(body, func(bytesTransferred int64) {
					diff := bytesTransferred - rangeProgress
					rangeProgress = bytesTransferred
					prog.byteCountLock.Lock()
					prog.byteCount += diff
					o.Progress(prog.byteCount)
					prog.byteCountLock.Unlock()
				})
			}
			if _, err = io.Copy(shared.NewSectionWriter(writer, chunkStart, count), body); err != nil {
				_ = body.Close()
				return err
			}
			if dr.StructuredBodyType != nil && *dr.StructuredBodyType != "" && dr.ContentRange != nil {
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

func (b *Client) parallelDownloadFrom(ctx context.Context, writer io.WriterAt, o downloadOptions, writerOffset int64, remaining int64, prog *downloadProgress) (int64, error) {
	dataDownloaded := int64(0)

	err := shared.DoBatchTransfer(ctx, &shared.BatchTransferOptions{
		OperationName: "downloadBlobToWriterAt",
		TransferSize:  remaining,
		ChunkSize:     o.BlockSize,
		NumChunks:     uint64(((remaining - 1) / o.BlockSize) + 1),
		Concurrency:   o.Concurrency,
		Operation: func(ctx context.Context, chunkStart int64, count int64) error {
			dr, err := b.DownloadStream(ctx, o.getDownloadBlobOptions(HTTPRange{
				Offset: chunkStart + writerOffset + o.Range.Offset, Count: count}, nil))
			if err != nil {
				return err
			}
			var body io.ReadCloser = dr.NewRetryReader(ctx, &o.RetryReaderOptionsPerBlock)
			if o.Progress != nil {
				rangeProgress := int64(0)
				body = streaming.NewResponseProgress(body, func(bytesTransferred int64) {
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

	dr, err := b.generated().Download(ctx, downloadOptions, leaseAccessConditions, cpkInfo, modifiedAccessConditions)
	if err != nil {
		return DownloadStreamResponse{}, err
	}

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

func parseContentRangeTotal(contentRange string) int64 {
	var start, end, total int64
	if _, err := fmt.Sscanf(contentRange, "bytes %d-%d/%d", &start, &end, &total); err != nil {
		return 0
	}
	return total
}
