package azblob

import (
	"context"
	"io"
	"net/url"

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
	client *blockBlobClient
	u      url.URL
}

// NewBlockBlobClient creates a BlockBlobClient object using the specified URL and request policy pipeline.
func NewBlockBlobClient(blobURL string, cred azcore.Credential, options *connectionOptions) (BlockBlobClient, error) {
	u, err := url.Parse(blobURL)
	if err != nil {
		return BlockBlobClient{}, err
	}
	con := newConnection(blobURL, cred, options)
	return BlockBlobClient{
		client:     &blockBlobClient{con: con},
		u:          *u,
		BlobClient: BlobClient{client: &blobClient{con: con}},
	}, nil
	//return bbc, nil
}

// WithPipeline creates a new BlockBlobClient object identical to the source but with the specific request policy pipeline.
func (bb BlockBlobClient) WithPipeline(pipeline azcore.Pipeline) BlockBlobClient {
	con := newConnectionWithPipeline(bb.u.String(), pipeline)
	return BlockBlobClient{
		client:     &blockBlobClient{con},
		u:          bb.u,
		BlobClient: BlobClient{client: &blobClient{con: con}},
	}
}

// URL returns the URL endpoint used by the BlobClient object.
func (bb BlockBlobClient) URL() url.URL {
	return bb.u
}

// WithSnapshot creates a new BlockBlobClient object identical to the source but with the specified snapshot timestamp.
// Pass "" to remove the snapshot returning a URL to the base blob.
func (bb BlockBlobClient) WithSnapshot(snapshot string) BlockBlobClient {
	p := NewBlobURLParts(bb.URL())
	p.Snapshot = snapshot
	snapshotURL := p.URL()
	con := newConnectionWithPipeline(snapshotURL.String(), bb.client.con.p)
	return BlockBlobClient{
		client: &blockBlobClient{
			con: con,
		},
		u:          snapshotURL,
		BlobClient: BlobClient{client: &blobClient{con: con}},
	}
}

// WithVersionID creates a new AppendBlobURL object identical to the source but with the specified version id.
// Pass "" to remove the versionID returning a URL to the base blob.
func (ab BlockBlobClient) WithVersionID(versionID string) BlockBlobClient {
	p := NewBlobURLParts(ab.URL())
	p.VersionID = versionID
	versionIDURL := p.URL()
	con := newConnectionWithPipeline(versionIDURL.String(), ab.client.con.p)
	return BlockBlobClient{
		client:     &blockBlobClient{con: con},
		u:          versionIDURL,
		BlobClient: BlobClient{client: &blobClient{con: con}},
	}
}

// Upload creates a new block blob or overwrites an existing block blob.
// Updating an existing block blob overwrites any existing metadata on the blob. Partial updates are not
// supported with Upload; the content of the existing blob is overwritten with the new content. To
// perform a partial update of a block blob, use StageBlock and CommitBlockList.
// This method panics if the stream is not at position 0.
// Note that the http client closes the body stream after the request is sent to the service.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-blob.
func (bb BlockBlobClient) Upload(ctx context.Context, body io.ReadSeeker, options *UploadBlockBlobOptions) (BlockBlobUploadResponse, error) {
	count, err := validateSeekableStreamAt0AndGetCount(body)
	if err != nil {
		return BlockBlobUploadResponse{}, err
	}

	basics, httpHeaders, leaseInfo, cpkV, cpkN, accessConditions := options.pointers()

	resp, err := bb.client.Upload(ctx, count, azcore.NopCloser(body), basics, httpHeaders, leaseInfo, cpkV, cpkN, accessConditions)

	return resp, handleError(err)
}

// StageBlock uploads the specified block to the block blob's "staging area" to be later committed by a call to CommitBlockList.
// Note that the http client closes the body stream after the request is sent to the service.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-block.
func (bb BlockBlobClient) StageBlock(ctx context.Context, base64BlockID string, body io.ReadSeeker, options *StageBlockOptions) (BlockBlobStageBlockResponse, error) {
	count, err := validateSeekableStreamAt0AndGetCount(body)
	if err != nil {
		return BlockBlobStageBlockResponse{}, err
	}

	ac, stageBlockOptions, cpkInfo, cpkScopeInfo := options.pointers()
	resp, err := bb.client.StageBlock(ctx, base64BlockID, count, azcore.NopCloser(body), stageBlockOptions, ac, cpkInfo, cpkScopeInfo)

	return resp, handleError(err)
}

// StageBlockFromURL copies the specified block from a source URL to the block blob's "staging area" to be later committed by a call to CommitBlockList.
// If count is CountToEnd (0), then data is read from specified offset to the end.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/put-block-from-url.
func (bb BlockBlobClient) StageBlockFromURL(ctx context.Context, base64BlockID string, sourceURL url.URL, contentLength int64, options *StageBlockFromURLOptions) (BlockBlobStageBlockFromURLResponse, error) {
	ac, smac, stageOptions, cpkInfo, cpkScope := options.pointers()

	resp, err := bb.client.StageBlockFromURL(ctx, base64BlockID, contentLength, sourceURL, stageOptions, cpkInfo, cpkScope, ac, smac)

	return resp, handleError(err)
}

// CommitBlockList writes a blob by specifying the list of block IDs that make up the blob.
// In order to be written as part of a blob, a block must have been successfully written
// to the server in a prior PutBlock operation. You can call PutBlockList to update a blob
// by uploading only those blocks that have changed, then committing the new and existing
// blocks together. Any blocks not specified in the block list and permanently deleted.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-block-list.
func (bb BlockBlobClient) CommitBlockList(ctx context.Context, base64BlockIDs []string, options *CommitBlockListOptions) (BlockBlobCommitBlockListResponse, error) {
	commitOptions, headers, cpkInfo, cpkScope, modifiedAccess, leaseAccess := options.pointers()

	resp, err := bb.client.CommitBlockList(ctx, BlockLookupList{
		Latest: &base64BlockIDs,
	}, commitOptions, headers, leaseAccess, cpkInfo, cpkScope, modifiedAccess)

	return resp, handleError(err)
}

// GetBlockList returns the list of blocks that have been uploaded as part of a block blob using the specified block list filter.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-block-list.
func (bb BlockBlobClient) GetBlockList(ctx context.Context, listType BlockListType, options *GetBlockListOptions) (BlockListResponse, error) {
	o, mac, lac := options.pointers()

	resp, err := bb.client.GetBlockList(ctx, listType, o, lac, mac)

	return resp, handleError(err)
}

// CopyFromURL synchronously copies the data at the source URL to a block blob, with sizes up to 256 MB.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/copy-blob-from-url.
func (bb BlockBlobClient) CopyFromURL(ctx context.Context, source url.URL, options *CopyBlockBlobFromURLOptions) (BlobCopyFromURLResponse, error) {
	copyOptions, smac, mac, lac := options.pointers()

	bClient := blobClient{
		con: bb.client.con,
	}

	resp, err := bClient.CopyFromURL(ctx, source, copyOptions, smac, mac, lac)

	return resp, handleError(err)
}
