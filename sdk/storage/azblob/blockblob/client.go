//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob

import (
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
)

// ClientOptions adds additional client options while constructing connection
type ClientOptions = exported.ClientOptions

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// Client defines a set of operations applicable to block blobs.
type Client base.CompositeClient[generated.BlobClient, generated.BlockBlobClient]

// NewClient creates a Client object using the specified URL, Azure AD credential, and options.
func NewClient(blobURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{shared.TokenScope}, nil)
	conOptions := exported.GetConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(shared.ModuleName, shared.ModuleVersion, runtime.PipelineOptions{}, conOptions)

	return (*Client)(base.NewBlockBlobClient(blobURL, pl)), nil
}

// NewClientWithNoCredential creates a Client object using the specified URL and options.
func NewClientWithNoCredential(blobURL string, options *ClientOptions) (*Client, error) {
	conOptions := exported.GetConnectionOptions(options)
	pl := runtime.NewPipeline(shared.ModuleName, shared.ModuleVersion, runtime.PipelineOptions{}, conOptions)

	return (*Client)(base.NewBlockBlobClient(blobURL, pl)), nil
}

// NewClientWithSharedKey creates a Client object using the specified URL, shared key, and options.
func NewClientWithSharedKey(blobURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := exported.GetConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(shared.ModuleName, shared.ModuleVersion, runtime.PipelineOptions{}, conOptions)

	return (*Client)(base.NewBlockBlobClient(blobURL, pl)), nil
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

func (bb *Client) generated() *generated.BlockBlobClient {
	_, blockBlob := base.InnerClients((*base.CompositeClient[generated.BlobClient, generated.BlockBlobClient])(bb))
	return blockBlob
}

// URL returns the URL endpoint used by the Client object.
func (bb *Client) URL() string {
	return bb.generated().Endpoint()
}

func (bb *Client) blobClient() *blob.Client {
	blobClient, _ := base.InnerClients((*base.CompositeClient[generated.BlobClient, generated.BlockBlobClient])(bb))
	return (*blob.Client)(blobClient)
}

// WithSnapshot creates a new Client object identical to the source but with the specified snapshot timestamp.
// Pass "" to remove the snapshot returning a URL to the base blob.
func (bb *Client) WithSnapshot(snapshot string) (*Client, error) {
	p, err := exported.ParseBlobURL(bb.URL())
	if err != nil {
		return nil, err
	}
	p.Snapshot = snapshot

	return (*Client)(base.NewBlockBlobClient(p.URL(), bb.generated().Pipeline())), nil
}

// WithVersionID creates a new AppendBlobURL object identical to the source but with the specified version id.
// Pass "" to remove the versionID returning a URL to the base blob.
func (bb *Client) WithVersionID(versionID string) (*Client, error) {
	p, err := exported.ParseBlobURL(bb.URL())
	if err != nil {
		return nil, err
	}
	p.VersionID = versionID

	return (*Client)(base.NewBlockBlobClient(p.URL(), bb.generated().Pipeline())), nil
}

// Upload creates a new block blob or overwrites an existing block blob.
// Updating an existing block blob overwrites any existing metadata on the blob. Partial updates are not
// supported with Upload; the content of the existing blob is overwritten with the new content. To
// perform a partial update of a block blob, use StageBlock and CommitBlockList.
// This method panics if the stream is not at position 0.
// Note that the http client closes the body stream after the request is sent to the service.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-blob.
func (bb *Client) Upload(ctx context.Context, body io.ReadSeekCloser, options *UploadOptions) (UploadResponse, error) {
	count, err := shared.ValidateSeekableStreamAt0AndGetCount(body)
	if err != nil {
		return UploadResponse{}, err
	}

	var basics *generated.BlockBlobClientUploadOptions
	var httpHeaders *generated.BlobHTTPHeaders
	var leaseInfo *lease.AccessConditions
	var cpkInfo *generated.CpkInfo
	var cpkScopeInfo *generated.CpkScopeInfo
	var accessConditions *generated.ModifiedAccessConditions

	if options != nil {
		basics = &generated.BlockBlobClientUploadOptions{
			BlobTagsString:          shared.SerializeBlobTagsToStrPtr(options.Tags),
			Metadata:                options.Metadata,
			Tier:                    options.Tier,
			TransactionalContentMD5: options.TransactionalContentMD5,
		}

		httpHeaders = options.HTTPHeaders
		leaseInfo, accessConditions = exported.FormatBlobAccessConditions(options.AccessConditions)
		cpkInfo = options.CpkInfo
		cpkScopeInfo = options.CpkScopeInfo
	}

	resp, err := bb.generated().Upload(ctx, count, body, basics, httpHeaders, leaseInfo, cpkInfo, cpkScopeInfo, accessConditions)
	return resp, err
}

// StageBlock uploads the specified block to the block blob's "staging area" to be later committed by a call to CommitBlockList.
// Note that the http client closes the body stream after the request is sent to the service.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-block.
func (bb *Client) StageBlock(ctx context.Context, base64BlockID string, body io.ReadSeekCloser, options *StageBlockOptions) (StageBlockResponse, error) {
	count, err := shared.ValidateSeekableStreamAt0AndGetCount(body)
	if err != nil {
		return StageBlockResponse{}, err
	}

	var stageBlockOptions *generated.BlockBlobClientStageBlockOptions
	var accessConditions *lease.AccessConditions
	var cpkInfo *generated.CpkInfo
	var cpkScopeInfo *generated.CpkScopeInfo

	if options != nil {
		stageBlockOptions = &generated.BlockBlobClientStageBlockOptions{
			TransactionalContentCRC64: options.TransactionalContentCRC64,
			TransactionalContentMD5:   options.TransactionalContentMD5,
		}
		accessConditions = options.LeaseAccessConditions
		cpkInfo = options.CpkInfo
		cpkScopeInfo = options.CpkScopeInfo
	}

	resp, err := bb.generated().StageBlock(ctx, base64BlockID, count, body, stageBlockOptions, accessConditions, cpkInfo, cpkScopeInfo)
	return resp, err
}

// StageBlockFromURL copies the specified block from a source URL to the block blob's "staging area" to be later committed by a call to CommitBlockList.
// If count is CountToEnd (0), then data is read from specified offset to the end.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/put-block-from-url.
func (bb *Client) StageBlockFromURL(ctx context.Context, base64BlockID string, sourceURL string,
	contentLength int64, options *StageBlockFromURLOptions) (StageBlockFromURLResponse, error) {

	stageBlockFromURLOptions, cpkInfo, cpkScopeInfo, leaseAccessConditions, sourceModifiedAccessConditions := options.format()

	resp, err := bb.generated().StageBlockFromURL(ctx, base64BlockID, contentLength, sourceURL, stageBlockFromURLOptions,
		cpkInfo, cpkScopeInfo, leaseAccessConditions, sourceModifiedAccessConditions)

	return resp, err
}

// CommitBlockList writes a blob by specifying the list of block IDs that make up the blob.
// In order to be written as part of a blob, a block must have been successfully written
// to the server in a prior PutBlock operation. You can call PutBlockList to update a blob
// by uploading only those blocks that have changed, then committing the new and existing
// blocks together. Any blocks not specified in the block list and permanently deleted.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-block-list.
func (bb *Client) CommitBlockList(ctx context.Context, base64BlockIDs []string, options *CommitBlockListOptions) (CommitBlockListResponse, error) {
	// this is a code smell in the generated code
	blockIds := make([]*string, len(base64BlockIDs))
	for k, v := range base64BlockIDs {
		blockIds[k] = to.Ptr(v)
	}

	blockLookupList := generated.BlockLookupList{Latest: blockIds}

	var commitOptions *generated.BlockBlobClientCommitBlockListOptions
	var headers *generated.BlobHTTPHeaders
	var leaseAccess *lease.AccessConditions
	var cpkInfo *generated.CpkInfo
	var cpkScope *generated.CpkScopeInfo
	var modifiedAccess *generated.ModifiedAccessConditions

	if options != nil {
		commitOptions = &generated.BlockBlobClientCommitBlockListOptions{
			BlobTagsString:            shared.SerializeBlobTagsToStrPtr(options.Tags),
			Metadata:                  options.Metadata,
			RequestID:                 options.RequestID,
			Tier:                      options.Tier,
			Timeout:                   options.Timeout,
			TransactionalContentCRC64: options.TransactionalContentCRC64,
			TransactionalContentMD5:   options.TransactionalContentMD5,
		}

		headers = options.BlobHTTPHeaders
		leaseAccess, modifiedAccess = exported.FormatBlobAccessConditions(options.BlobAccessConditions)
		cpkInfo = options.CpkInfo
		cpkScope = options.CpkScopeInfo
	}

	resp, err := bb.generated().CommitBlockList(ctx, blockLookupList, commitOptions, headers, leaseAccess, cpkInfo, cpkScope, modifiedAccess)
	return resp, err
}

// GetBlockList returns the list of blocks that have been uploaded as part of a block blob using the specified block list filter.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-block-list.
func (bb *Client) GetBlockList(ctx context.Context, listType BlockListType, options *GetBlockListOptions) (GetBlockListResponse, error) {
	o, lac, mac := options.format()

	resp, err := bb.generated().GetBlockList(ctx, listType, o, lac, mac)

	return resp, err
}

// Redeclared APIs ----- Copy over to Append blob and Page blob as well.

// Download reads a range of bytes from a blob. The response also includes the blob's properties and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob.
func (bb *Client) Download(ctx context.Context, o *blob.DownloadOptions) (blob.DownloadResponse, error) {
	return bb.blobClient().Download(ctx, o)
}

// Delete marks the specified blob or snapshot for deletion. The blob is later deleted during garbage collection.
// Note that deleting a blob also deletes all its snapshots.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/delete-blob.
func (bb *Client) Delete(ctx context.Context, o *blob.DeleteOptions) (blob.DeleteResponse, error) {
	return bb.blobClient().Delete(ctx, o)
}

// Undelete restores the contents and metadata of a soft-deleted blob and any associated soft-deleted snapshots.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/undelete-blob.
func (bb *Client) Undelete(ctx context.Context, o *blob.UndeleteOptions) (blob.UndeleteResponse, error) {
	return bb.blobClient().Undelete(ctx, o)
}

// SetTier operation sets the tier on a blob. The operation is allowed on a page
// blob in a premium storage account and on a block blob in a blob storage account (locally
// redundant storage only). A premium page blob's tier determines the allowed size, IOPS, and
// bandwidth of the blob. A block blob's tier determines Hot/Cool/Archive storage type. This operation
// does not update the blob's ETag.
// For detailed information about block blob level tiering see https://docs.microsoft.com/en-us/azure/storage/blobs/storage-blob-storage-tiers.
func (bb *Client) SetTier(ctx context.Context, tier blob.AccessTier, o *blob.SetTierOptions) (blob.SetTierResponse, error) {
	return bb.blobClient().SetTier(ctx, tier, o)
}

// GetProperties returns the blob's properties.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob-properties.
func (bb *Client) GetProperties(ctx context.Context, o *blob.GetPropertiesOptions) (blob.GetPropertiesResponse, error) {
	return bb.blobClient().GetProperties(ctx, o)
}

// SetHTTPHeaders changes a blob's HTTP headers.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-blob-properties.
func (bb *Client) SetHTTPHeaders(ctx context.Context, HTTPHeaders blob.HTTPHeaders, o *blob.SetHTTPHeadersOptions) (blob.SetHTTPHeadersResponse, error) {
	return bb.blobClient().SetHTTPHeaders(ctx, HTTPHeaders, o)
}

// SetMetadata changes a blob's metadata.
// https://docs.microsoft.com/rest/api/storageservices/set-blob-metadata.
func (bb *Client) SetMetadata(ctx context.Context, metadata map[string]string, o *blob.SetMetadataOptions) (blob.SetMetadataResponse, error) {
	return bb.blobClient().SetMetadata(ctx, metadata, o)
}

// CreateSnapshot creates a read-only snapshot of a blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/snapshot-blob.
func (bb *Client) CreateSnapshot(ctx context.Context, o *blob.CreateSnapshotOptions) (blob.CreateSnapshotResponse, error) {
	return bb.blobClient().CreateSnapshot(ctx, o)
}

// StartCopyFromURL copies the data at the source URL to a blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/copy-blob.
func (bb *Client) StartCopyFromURL(ctx context.Context, copySource string, o *blob.StartCopyFromURLOptions) (blob.StartCopyFromURLResponse, error) {
	return bb.blobClient().StartCopyFromURL(ctx, copySource, o)
}

// AbortCopyFromURL stops a pending copy that was previously started and leaves a destination blob with 0 length and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/abort-copy-blob.
func (bb *Client) AbortCopyFromURL(ctx context.Context, copyID string, o *blob.AbortCopyFromURLOptions) (blob.AbortCopyFromURLResponse, error) {
	return bb.blobClient().AbortCopyFromURL(ctx, copyID, o)
}

// SetTags operation enables users to set tags on a blob or specific blob version, but not snapshot.
// Each call to this operation replaces all existing tags attached to the blob.
// To remove all tags from the blob, call this operation with no tags set.
// https://docs.microsoft.com/en-us/rest/api/storageservices/set-blob-tags
func (bb *Client) SetTags(ctx context.Context, o *blob.SetTagsOptions) (blob.SetTagsResponse, error) {
	return bb.blobClient().SetTags(ctx, o)
}

// GetTags operation enables users to get tags on a blob or specific blob version, or snapshot.
// https://docs.microsoft.com/en-us/rest/api/storageservices/get-blob-tags
func (bb *Client) GetTags(ctx context.Context, o *blob.GetTagsOptions) (blob.GetTagsResponse, error) {
	return bb.blobClient().GetTags(ctx, o)
}

// CopyFromURL synchronously copies the data at the source URL to a block blob, with sizes up to 256 MB.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/copy-blob-from-url.
func (bb *Client) CopyFromURL(ctx context.Context, copySource string, o *blob.CopyFromURLOptions) (blob.CopyFromURLResponse, error) {
	return bb.blobClient().CopyFromURL(ctx, copySource, o)
}
