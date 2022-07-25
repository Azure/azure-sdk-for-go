//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// Client represents a URL to an Azure Storage blob; the blob may be a block blob, append blob, or page blob.
type Client base.Client[generated.BlobClient]

// NewClient creates a Client object using the specified URL, Azure AD credential, and options.
func NewClient(blobURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{shared.TokenScope}, nil)
	conOptions := exported.GetConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, conOptions)

	return (*Client)(base.NewBlobClient(blobURL, pl)), nil
}

// NewClientWithNoCredential creates a Client object using the specified URL and options.
func NewClientWithNoCredential(blobURL string, options *ClientOptions) (*Client, error) {
	conOptions := exported.GetConnectionOptions(options)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, conOptions)

	return (*Client)(base.NewBlobClient(blobURL, pl)), nil
}

// NewClientWithSharedKey creates a Client object using the specified URL, shared key, and options.
func NewClientWithSharedKey(blobURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := exported.GetConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, conOptions)

	return (*Client)(base.NewBlobClient(blobURL, pl)), nil
}

// NewClientFromConnectionString creates Client from a connection String
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
		return NewClientWithSharedKey(parsed.ServiceURL, credential, options)
	}

	return NewClientWithNoCredential(parsed.ServiceURL, options)
}

// NewLeaseClient generates blob lease.Client from the blob.Client
func (b *Client) NewLeaseClient(leaseID *string) (*LeaseClient, error) {
	leaseID, err := shared.GenerateLeaseID(leaseID)
	if err != nil {
		return nil, err
	}
	return &LeaseClient{
		blobClient: (*Client)(base.NewBlobClient(b.URL(), b.generated().Pipeline())),
		leaseID:    leaseID,
	}, nil
}

func (b *Client) generated() *generated.BlobClient {
	return base.InnerClient((*base.Client[generated.BlobClient])(b))
}

func (b *Client) sharedKey() *SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.BlobClient])(b))
}

// URL returns the URL endpoint used by the Client object.
func (b *Client) URL() string {
	return b.generated().Endpoint()
}

// WithSnapshot creates a new Client object identical to the source but with the specified snapshot timestamp.
// Pass "" to remove the snapshot returning a URL to the base blob.
func (b *Client) WithSnapshot(snapshot string) (*Client, error) {
	p, err := exported.ParseBlobURL(b.URL())
	if err != nil {
		return nil, err
	}
	p.Snapshot = snapshot

	return (*Client)(base.NewBlobClient(p.URL(), b.generated().Pipeline())), nil
}

// WithVersionID creates a new AppendBlobURL object identical to the source but with the specified version id.
// Pass "" to remove the versionID returning a URL to the base blob.
func (b *Client) WithVersionID(versionID string) (*Client, error) {
	p, err := exported.ParseBlobURL(b.URL())
	if err != nil {
		return nil, err
	}
	p.VersionID = versionID

	return (*Client)(base.NewBlobClient(p.URL(), b.generated().Pipeline())), nil
}

// Download reads a range of bytes from a blob. The response also includes the blob's properties and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob.
func (b *Client) Download(ctx context.Context, o *DownloadOptions) (DownloadResponse, error) {
	downloadOptions, leaseAccessConditions, cpkInfo, modifiedAccessConditions := o.format()

	dr, err := b.generated().Download(ctx, downloadOptions, leaseAccessConditions, cpkInfo, modifiedAccessConditions)
	if err != nil {
		return DownloadResponse{}, err
	}

	offset := int64(0)
	count := int64(CountToEnd)

	if o != nil && o.Offset != nil {
		offset = *o.Offset
	}

	if o != nil && o.Count != nil {
		count = *o.Count
	}

	eTag := ""
	if dr.ETag != nil {
		eTag = *dr.ETag
	}
	return DownloadResponse{
		b:                          b,
		BlobClientDownloadResponse: dr,
		getInfo:                    HTTPGetterInfo{Offset: offset, Count: count, ETag: eTag},
		ObjectReplicationRules:     deserializeORSPolicies(dr.ObjectReplicationRules),
	}, err
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
// redundant storage only). A premium page blob's tier determines the allowed size, IOPS, and
// bandwidth of the blob. A block blob's tier determines Hot/Cool/Archive storage type. This operation
// does not update the blob's ETag.
// For detailed information about block blob level tiering see https://docs.microsoft.com/en-us/azure/storage/blobs/storage-blob-storage-tiers.
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
func (b *Client) SetHTTPHeaders(ctx context.Context, HTTPHeaders HTTPHeaders, o *SetHTTPHeadersOptions) (SetHTTPHeadersResponse, error) {
	opts, leaseAccessConditions, modifiedAccessConditions := o.format()
	resp, err := b.generated().SetHTTPHeaders(ctx, opts, &HTTPHeaders, leaseAccessConditions, modifiedAccessConditions)
	return resp, err
}

// SetMetadata changes a blob's metadata.
// https://docs.microsoft.com/rest/api/storageservices/set-blob-metadata.
func (b *Client) SetMetadata(ctx context.Context, metadata map[string]string, o *SetMetadataOptions) (SetMetadataResponse, error) {
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
func (b *Client) SetTags(ctx context.Context, options *SetTagsOptions) (SetTagsResponse, error) {
	blobSetTagsOptions, modifiedAccessConditions, leaseAccessConditions := options.format()
	resp, err := b.generated().SetTags(ctx, blobSetTagsOptions, modifiedAccessConditions, leaseAccessConditions)

	return resp, err
}

// GetTags operation enables users to get tags on a blob or specific blob version, or snapshot.
// https://docs.microsoft.com/en-us/rest/api/storageservices/get-blob-tags
func (b *Client) GetTags(ctx context.Context, options *GetTagsOptions) (GetTagsResponse, error) {
	blobGetTagsOptions, modifiedAccessConditions, leaseAccessConditions := options.format()
	resp, err := b.generated().GetTags(ctx, blobGetTagsOptions, modifiedAccessConditions, leaseAccessConditions)
	return resp, err

}

// CopyFromURL synchronously copies the data at the source URL to a block blob, with sizes up to 256 MB.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/copy-blob-from-url.
func (b *Client) CopyFromURL(ctx context.Context, copySource string, options *CopyFromURLOptions) (CopyFromURLResponse, error) {
	copyOptions, smac, mac, lac := options.format()
	resp, err := b.generated().CopyFromURL(ctx, copySource, copyOptions, smac, mac, lac)
	return resp, err
}

// GetSASToken is a convenience method for generating a SAS token for the currently pointed at blob.
// It can only be used if the credential supplied during creation was a SharedKeyCredential.
func (b *Client) GetSASToken(permissions SASPermissions, start time.Time, expiry time.Time) (string, error) {
	urlParts, _ := exported.ParseBlobURL(b.URL())

	t, err := time.Parse(exported.SnapshotTimeFormat, urlParts.Snapshot)

	if err != nil {
		t = time.Time{}
	}

	if b.sharedKey() == nil {
		return "", errors.New("credential is not a SharedKeyCredential. SAS can only be signed with a SharedKeyCredential")
	}

	qps, err := exported.BlobSASSignatureValues{
		ContainerName: urlParts.ContainerName,
		BlobName:      urlParts.BlobName,
		SnapshotTime:  t,
		Version:       exported.SASVersion,

		Permissions: permissions.String(),

		StartTime:  start.UTC(),
		ExpiryTime: expiry.UTC(),
	}.NewSASQueryParameters(b.sharedKey())

	if err != nil {
		return "", err
	}

	endpoint := b.URL()
	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}
	endpoint += "?" + qps.Encode()

	return endpoint, nil
}

// Concurrent Download Functions -----------------------------------------------------------------------------------------

// DownloadToWriterAt downloads an Azure blob to a WriterAt in parallel.
// Offset and count are optional, pass 0 for both to download the entire blob.
func (b *Client) DownloadToWriterAt(ctx context.Context, offset, count int64, writer io.WriterAt, o *DownloadToWriterAtOptions) error {
	if o.BlockSize == 0 {
		o.BlockSize = DefaultDownloadBlockSize
	}

	if count == CountToEnd { // If size not specified, calculate it
		// If we don't have the length at all, get it
		downloadBlobOptions := o.getDownloadBlobOptions(0, CountToEnd, nil)
		dr, err := b.Download(ctx, downloadBlobOptions)
		if err != nil {
			return err
		}
		count = *dr.ContentLength - offset
	}

	if count <= 0 {
		// The file is empty, there is nothing to download.
		return nil
	}

	// Prepare and do parallel download.
	progress := int64(0)
	progressLock := &sync.Mutex{}

	err := shared.DoBatchTransfer(ctx, &shared.BatchTransferOptions{
		OperationName: "downloadBlobToWriterAt",
		TransferSize:  count,
		ChunkSize:     o.BlockSize,
		Parallelism:   o.Parallelism,
		Operation: func(chunkStart int64, count int64, ctx context.Context) error {

			downloadBlobOptions := o.getDownloadBlobOptions(chunkStart+offset, count, nil)
			dr, err := b.Download(ctx, downloadBlobOptions)
			if err != nil {
				return err
			}
			body := dr.BodyReader(&o.RetryReaderOptionsPerBlock)
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
			err = body.Close()
			return err
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// DownloadToBuffer downloads an Azure blob to a buffer with parallel.
// Offset and count are optional, pass 0 for both to download the entire blob.
func (b *Client) DownloadToBuffer(ctx context.Context, offset, count int64, _bytes []byte, o *DownloadToBufferOptions) error {
	return b.DownloadToWriterAt(ctx, offset, count, shared.NewBytesWriter(_bytes), o)
}

// DownloadToFile downloads an Azure blob to a local file.
// The file would be truncated if the size doesn't match.
// Offset and count are optional, pass 0 for both to download the entire blob.
func (b *Client) DownloadToFile(ctx context.Context, offset, count int64, file *os.File, o *DownloadToFileOptions) error {
	// 1. Calculate the size of the destination file
	var size int64

	if count == CountToEnd {
		// Try to get Azure blob's size
		getBlobPropertiesOptions := o.getBlobPropertiesOptions()
		props, err := b.GetProperties(ctx, getBlobPropertiesOptions)
		if err != nil {
			return err
		}
		size = *props.ContentLength - offset
	} else {
		size = count
	}

	// 2. Compare and try to resize local file's size if it doesn't match Azure blob's size.
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if stat.Size() != size {
		if err = file.Truncate(size); err != nil {
			return err
		}
	}

	if size > 0 {
		return b.DownloadToWriterAt(ctx, offset, size, file, o)
	} else { // if the blob's size is 0, there is no need in downloading it
		return nil
	}
}
