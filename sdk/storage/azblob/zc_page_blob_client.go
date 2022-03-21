// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"io"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const (
	// PageBlobPageBytes indicates the number of bytes in a page (512).
	PageBlobPageBytes = 512
)

// PageBlobClient represents a client to an Azure Storage page blob;
type PageBlobClient struct {
	BlobClient
	client *pageBlobClient
}

// NewPageBlobClient creates a ServiceClient object using the specified URL, Azure AD credential, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewPageBlobClient(blobURL string, cred azcore.TokenCredential, options *ClientOptions) (*PageBlobClient, error) {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{tokenScope}, nil)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := newConnection(blobURL, authPolicy, conOptions)

	return &PageBlobClient{
		client: newPageBlobClient(conn.Endpoint(), conn.Pipeline()),
		BlobClient: BlobClient{
			client: newBlobClient(conn.Endpoint(), conn.Pipeline()),
			conn:   conn,
		},
	}, nil
}

// NewPageBlobClientWithNoCredential creates a ServiceClient object using the specified URL and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net?<SAS token>
func NewPageBlobClientWithNoCredential(blobURL string, options *ClientOptions) (*PageBlobClient, error) {
	conOptions := getConnectionOptions(options)
	conn := newConnection(blobURL, nil, conOptions)

	return &PageBlobClient{
		client: newPageBlobClient(conn.Endpoint(), conn.Pipeline()),
		BlobClient: BlobClient{
			client: newBlobClient(conn.Endpoint(), conn.Pipeline()),
			conn:   conn,
		},
	}, nil
}

// NewPageBlobClientWithSharedKey creates a ServiceClient object using the specified URL, shared key, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewPageBlobClientWithSharedKey(blobURL string, cred *SharedKeyCredential, options *ClientOptions) (*PageBlobClient, error) {
	authPolicy := newSharedKeyCredPolicy(cred)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := newConnection(blobURL, authPolicy, conOptions)

	return &PageBlobClient{
		client: newPageBlobClient(conn.Endpoint(), conn.Pipeline()),
		BlobClient: BlobClient{
			client:    newBlobClient(conn.Endpoint(), conn.Pipeline()),
			conn:      conn,
			sharedKey: cred,
		},
	}, nil
}

// WithSnapshot creates a new PageBlobURL object identical to the source but with the specified snapshot timestamp.
// Pass "" to remove the snapshot returning a URL to the base blob.
func (pb *PageBlobClient) WithSnapshot(snapshot string) (*PageBlobClient, error) {
	p, err := NewBlobURLParts(pb.URL())
	if err != nil {
		return nil, err
	}
	p.Snapshot = snapshot
	conn := pb.conn
	conn.u = p.URL()

	return &PageBlobClient{
		client: newPageBlobClient(conn.Endpoint(), conn.Pipeline()),
		BlobClient: BlobClient{
			client:    newBlobClient(conn.Endpoint(), conn.Pipeline()),
			conn:      conn,
			sharedKey: pb.sharedKey,
		},
	}, nil
}

// WithVersionID creates a new PageBlobURL object identical to the source but with the specified snapshot timestamp.
// Pass "" to remove the version returning a URL to the base blob.
func (pb *PageBlobClient) WithVersionID(versionID string) (*PageBlobClient, error) {
	p, err := NewBlobURLParts(pb.URL())
	if err != nil {
		return nil, err
	}
	p.VersionID = versionID
	conn := pb.conn
	conn.u = p.URL()

	return &PageBlobClient{
		client: newPageBlobClient(conn.Endpoint(), conn.Pipeline()),
		BlobClient: BlobClient{
			client:    newBlobClient(conn.Endpoint(), conn.Pipeline()),
			conn:      conn,
			sharedKey: pb.sharedKey,
		},
	}, nil
}

// Create creates a page blob of the specified length. Call PutPage to upload data to a page blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-blob.
func (pb *PageBlobClient) Create(ctx context.Context, size int64, o *PageBlobCreateOptions) (PageBlobCreateResponse, error) {
	createOptions, HTTPHeaders, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions := o.pointers()

	resp, err := pb.client.Create(ctx, 0, size, createOptions, HTTPHeaders, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions)

	return toPageBlobCreateResponse(resp), handleError(err)
}

// UploadPages writes 1 or more pages to the page blob. The start offset and the stream size must be a multiple of 512 bytes.
// This method panics if the stream is not at position 0.
// Note that the http client closes the body stream after the request is sent to the service.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-page.
func (pb *PageBlobClient) UploadPages(ctx context.Context, body io.ReadSeekCloser, options *PageBlobUploadPagesOptions) (PageBlobUploadPagesResponse, error) {
	count, err := validateSeekableStreamAt0AndGetCount(body)

	if err != nil {
		return PageBlobUploadPagesResponse{}, err
	}

	uploadPagesOptions, leaseAccessConditions, cpkInfo, cpkScopeInfo, sequenceNumberAccessConditions, modifiedAccessConditions := options.pointers()

	resp, err := pb.client.UploadPages(ctx, count, body, uploadPagesOptions, leaseAccessConditions,
		cpkInfo, cpkScopeInfo, sequenceNumberAccessConditions, modifiedAccessConditions)

	return toPageBlobUploadPagesResponse(resp), handleError(err)
}

// UploadPagesFromURL copies 1 or more pages from a source URL to the page blob.
// The sourceOffset specifies the start offset of source data to copy from.
// The destOffset specifies the start offset of data in page blob will be written to.
// The count must be a multiple of 512 bytes.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-page-from-url.
func (pb *PageBlobClient) UploadPagesFromURL(ctx context.Context, source string, sourceOffset, destOffset, count int64,
	options *PageBlobUploadPagesFromURLOptions) (PageBlobUploadPagesFromURLResponse, error) {

	uploadPagesFromURLOptions, cpkInfo, cpkScopeInfo, leaseAccessConditions, sequenceNumberAccessConditions, modifiedAccessConditions, sourceModifiedAccessConditions := options.pointers()

	resp, err := pb.client.UploadPagesFromURL(ctx, source, rangeToString(sourceOffset, count), 0,
		rangeToString(destOffset, count), uploadPagesFromURLOptions, cpkInfo, cpkScopeInfo, leaseAccessConditions,
		sequenceNumberAccessConditions, modifiedAccessConditions, sourceModifiedAccessConditions)

	return toPageBlobUploadPagesFromURLResponse(resp), handleError(err)
}

// ClearPages frees the specified pages from the page blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-page.
func (pb *PageBlobClient) ClearPages(ctx context.Context, pageRange HttpRange, options *PageBlobClearPagesOptions) (PageBlobClearPagesResponse, error) {
	clearOptions := &pageBlobClientClearPagesOptions{
		Range: pageRange.pointers(),
	}

	leaseAccessConditions, cpkInfo, cpkScopeInfo, sequenceNumberAccessConditions, modifiedAccessConditions := options.pointers()

	resp, err := pb.client.ClearPages(ctx, 0, clearOptions, leaseAccessConditions, cpkInfo,
		cpkScopeInfo, sequenceNumberAccessConditions, modifiedAccessConditions)

	return toPageBlobClearPagesResponse(resp), handleError(err)
}

// GetPageRanges returns the list of valid page ranges for a page blob or snapshot of a page blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-page-ranges.
func (pb *PageBlobClient) GetPageRanges(options *PageBlobGetPageRangesOptions) *PageBlobGetPageRangesPager {
	getPageRangesOptions, leaseAccessConditions, modifiedAccessConditions := options.pointers()

	pageBlobGetPageRangesPager := pb.client.GetPageRanges(getPageRangesOptions, leaseAccessConditions, modifiedAccessConditions)

	return toPageBlobGetPageRangesPager(pageBlobGetPageRangesPager)
}

// GetPageRangesDiff gets the collection of page ranges that differ between a specified snapshot and this page blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-page-ranges.
func (pb *PageBlobClient) GetPageRangesDiff(options *PageBlobGetPageRangesDiffOptions) *PageBlobGetPageRangesDiffPager {
	getPageRangesDiffOptions, leaseAccessConditions, modifiedAccessConditions := options.pointers()

	getPageRangesDiffPager := pb.client.GetPageRangesDiff(&getPageRangesDiffOptions, leaseAccessConditions, modifiedAccessConditions)

	return toPageBlobGetPageRangesDiffPager(getPageRangesDiffPager)
}

// Resize resizes the page blob to the specified size (which must be a multiple of 512).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-blob-properties.
func (pb *PageBlobClient) Resize(ctx context.Context, size int64, options *PageBlobResizeOptions) (PageBlobResizeResponse, error) {
	resizeOptions, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions := options.pointers()

	resp, err := pb.client.Resize(ctx, size, resizeOptions, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions)

	return toPageBlobResizeResponse(resp), handleError(err)
}

// UpdateSequenceNumber sets the page blob's sequence number.
func (pb *PageBlobClient) UpdateSequenceNumber(ctx context.Context, options *PageBlobUpdateSequenceNumberOptions) (PageBlobUpdateSequenceNumberResponse, error) {
	actionType, updateOptions, lac, mac := options.pointers()
	resp, err := pb.client.UpdateSequenceNumber(ctx, *actionType, updateOptions, lac, mac)

	return toPageBlobUpdateSequenceNumberResponse(resp), handleError(err)
}

// StartCopyIncremental begins an operation to start an incremental copy from one page blob's snapshot to this page blob.
// The snapshot is copied such that only the differential changes between the previously copied snapshot are transferred to the destination.
// The copied snapshots are complete copies of the original snapshot and can be read or copied from as usual.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/incremental-copy-blob and
// https://docs.microsoft.com/en-us/azure/virtual-machines/windows/incremental-snapshots.
func (pb *PageBlobClient) StartCopyIncremental(ctx context.Context, copySource string, prevSnapshot string, options *PageBlobCopyIncrementalOptions) (PageBlobCopyIncrementalResponse, error) {
	copySourceURL, err := url.Parse(copySource)
	if err != nil {
		return PageBlobCopyIncrementalResponse{}, err
	}

	queryParams := copySourceURL.Query()
	queryParams.Set("snapshot", prevSnapshot)
	copySourceURL.RawQuery = queryParams.Encode()

	pageBlobCopyIncrementalOptions, modifiedAccessConditions := options.pointers()
	resp, err := pb.client.CopyIncremental(ctx, copySourceURL.String(), pageBlobCopyIncrementalOptions, modifiedAccessConditions)

	return toPageBlobCopyIncrementalResponse(resp), handleError(err)
}
