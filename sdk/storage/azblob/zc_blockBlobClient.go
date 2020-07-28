package azblob

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	// BlockBlobMaxUploadBlobBytes indicates the maximum number of bytes that can be sent in a call to Upload.
	BlockBlobMaxUploadBlobBytes = 256 * 1024 * 1024 // 256MB

	// BlockBlobMaxStageBlockBytes indicates the maximum number of bytes that can be sent in a call to StageBlock.
	BlockBlobMaxStageBlockBytes = 100 * 1024 * 1024 // 100MB

	// BlockBlobMaxBlocks indicates the maximum number of blocks allowed in a block blob.
	BlockBlobMaxBlocks = 50000
)

// BlockBlobClient defines a set of operations applicable to block blobs.
type BlockBlobClient struct {
	BlobClient
	client *client
}

// NewBlockBlobClient creates a BlockBlobClient object using the specified URL and request policy pipeline.
func NewBlockBlobClient(blobURL string, cred azcore.Credential, options *clientOptions) (BlockBlobClient, error) {
	client, err := newClient(blobURL, cred, options)

	if err != nil {
		return BlockBlobClient{}, err
	}

	return BlockBlobClient{client: client}, err
}

// WithPipeline creates a new BlockBlobClient object identical to the source but with the specific request policy pipeline.
func (bb BlockBlobClient) WithPipeline(pipeline azcore.Pipeline) (BlockBlobClient, error) {
	client, err := newClientWithPipeline(bb.client.u.String(), pipeline)

	if err != nil {
		return BlockBlobClient{}, err
	}

	return BlockBlobClient{client: client}, err
}

//// WithSnapshot creates a new BlockBlobClient object identical to the source but with the specified snapshot timestamp.
//// Pass "" to remove the snapshot returning a URL to the base blob.
//func (bb BlockBlobClient) WithSnapshot(snapshot string) BlockBlobClient {
//	p := NewBlobURLParts(bb.URL())
//	p.Snapshot = snapshot
//	return NewBlockBlobClient(p.URL(), bb.blobClient.Pipeline())
//}

// Upload creates a new block blob or overwrites an existing block blob.
// Updating an existing block blob overwrites any existing metadata on the blob. Partial updates are not
// supported with Upload; the content of the existing blob is overwritten with the new content. To
// perform a partial update of a block blob, use StageBlock and CommitBlockList.
// This method panics if the stream is not at position 0.
// Note that the http client closes the body stream after the request is sent to the service.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-blob.
func (bb BlockBlobClient) Upload(ctx context.Context, body azcore.ReadSeekCloser, options *UploadBlockBlobOptions) (*BlockBlobUploadResponse, error) {
	count, err := validateSeekableStreamAt0AndGetCount(body)
	if err != nil {
		return nil, err
	}

	basics, httpHeaders, leaseInfo, cpkV, cpkN, accessConditions := options.pointers()

	return bb.client.BlockBlobOperations().Upload(ctx, count, body, basics, httpHeaders, leaseInfo, cpkV, cpkN, accessConditions)
}

//
//// StageBlock uploads the specified block to the block blob's "staging area" to be later committed by a call to CommitBlockList.
//// Note that the http client closes the body stream after the request is sent to the service.
//// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-block.
//func (bb BlockBlobClient) StageBlock(ctx context.Context, base64BlockID string, body io.ReadSeeker, ac LeaseAccessConditions, transactionalMD5 []byte) (*BlockBlobStageBlockResponse, error) {
//	count, err := validateSeekableStreamAt0AndGetCount(body)
//	if err != nil {
//		return nil, err
//	}
//	return bb.client.StageBlock(ctx, base64BlockID, count, body, transactionalMD5, nil, nil, ac.pointers(),
//		nil, nil, EncryptionAlgorithmNone, // CPK
//		nil)
//}
//
//// StageBlockFromURL copies the specified block from a source URL to the block blob's "staging area" to be later committed by a call to CommitBlockList.
//// If count is CountToEnd (0), then data is read from specified offset to the end.
//// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/put-block-from-url.
//func (bb BlockBlobClient) StageBlockFromURL(ctx context.Context, base64BlockID string, sourceURL url.URL, offset int64, count int64, destinationAccessConditions LeaseAccessConditions, sourceAccessConditions ModifiedAccessConditions) (*BlockBlobStageBlockFromURLResponse, error) {
//	sourceIfModifiedSince, sourceIfUnmodifiedSince, sourceIfMatchETag, sourceIfNoneMatchETag := sourceAccessConditions.pointers()
//	return bb.client.StageBlockFromURL(ctx, base64BlockID, 0, sourceURL.String(), httpRange{offset: offset, count: count}.pointers(), nil, nil, nil,
//		nil, nil, EncryptionAlgorithmNone, // CPK
//		destinationAccessConditions.pointers(), sourceIfModifiedSince, sourceIfUnmodifiedSince, sourceIfMatchETag, sourceIfNoneMatchETag, nil)
//}
//
//// CommitBlockList writes a blob by specifying the list of block IDs that make up the blob.
//// In order to be written as part of a blob, a block must have been successfully written
//// to the server in a prior PutBlock operation. You can call PutBlockList to update a blob
//// by uploading only those blocks that have changed, then committing the new and existing
//// blocks together. Any blocks not specified in the block list and permanently deleted.
//// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-block-list.
//func (bb BlockBlobClient) CommitBlockList(ctx context.Context, base64BlockIDs []string, h BlobHTTPHeaders,
//	metadata Metadata, ac BlobAccessConditions) (*BlockBlobCommitBlockListResponse, error) {
//	ifModifiedSince, ifUnmodifiedSince, ifMatchETag, ifNoneMatchETag := ac.ModifiedAccessConditions.pointers()
//	return bb.client.CommitBlockList(ctx, BlockLookupList{Latest: base64BlockIDs}, nil,
//		&h.CacheControl, &h.ContentType, &h.ContentEncoding, &h.ContentLanguage, h.ContentMD5, nil, nil,
//		metadata, ac.LeaseAccessConditions.pointers(), &h.ContentDisposition,
//		nil, nil, EncryptionAlgorithmNone, // CPK
//		AccessTierNone,
//		ifModifiedSince, ifUnmodifiedSince, ifMatchETag, ifNoneMatchETag, nil)
//}
//
//// GetBlockList returns the list of blocks that have been uploaded as part of a block blob using the specified block list filter.
//// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-block-list.
//func (bb BlockBlobClient) GetBlockList(ctx context.Context, listType BlockListType, ac LeaseAccessConditions) (*BlockList, error) {
//	return bb.client.GetBlockList(ctx, listType, nil, nil, ac.pointers(), nil)
//}
//
//// CopyFromURL synchronously copies the data at the source URL to a block blob, with sizes up to 256 MB.
//// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/copy-blob-from-url.
//func (bb BlockBlobClient) CopyFromURL(ctx context.Context, source url.URL, metadata Metadata,
//	srcac ModifiedAccessConditions, dstac BlobAccessConditions, srcContentMD5 []byte) (*BlobCopyFromURLResponse, error) {
//
//	srcIfModifiedSince, srcIfUnmodifiedSince, srcIfMatchETag, srcIfNoneMatchETag := srcac.pointers()
//	dstIfModifiedSince, dstIfUnmodifiedSince, dstIfMatchETag, dstIfNoneMatchETag := dstac.ModifiedAccessConditions.pointers()
//	dstLeaseID := dstac.LeaseAccessConditions.pointers()
//
//	return bb.blobClient.CopyFromURL(ctx, source.String(), nil, metadata, AccessTierNone,
//		srcIfModifiedSince, srcIfUnmodifiedSince,
//		srcIfMatchETag, srcIfNoneMatchETag,
//		dstIfModifiedSince, dstIfUnmodifiedSince,
//		dstIfMatchETag, dstIfNoneMatchETag,
//		dstLeaseID, nil, srcContentMD5)
//}
