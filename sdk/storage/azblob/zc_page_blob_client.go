//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"io"
	"net/http"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// PageBlobClient represents a client to an Azure Storage page blob;
type PageBlobClient struct {
	BlobClient
	client *pageBlobClient
}

// NewPageBlobClient creates a ServiceClient object using the specified URL, Azure AD credential, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewPageBlobClient(blobURL string, cred azcore.TokenCredential, options *ClientOptions) (*PageBlobClient, error) {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{internal.TokenScope}, nil)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := internal.NewConnection(blobURL, conOptions)

	return &PageBlobClient{
		client: newPageBlobClient(conn.Endpoint(), conn.Pipeline()),
		BlobClient: BlobClient{
			client: newBlobClient(conn.Endpoint(), conn.Pipeline()),
		},
	}, nil
}

// NewPageBlobClientWithNoCredential creates a ServiceClient object using the specified URL and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net?<SAS token>
func NewPageBlobClientWithNoCredential(blobURL string, options *ClientOptions) (*PageBlobClient, error) {
	conOptions := getConnectionOptions(options)
	conn := internal.NewConnection(blobURL, conOptions)

	return &PageBlobClient{
		client: newPageBlobClient(conn.Endpoint(), conn.Pipeline()),
		BlobClient: BlobClient{
			client: newBlobClient(conn.Endpoint(), conn.Pipeline()),
		},
	}, nil
}

// NewPageBlobClientWithSharedKey creates a ServiceClient object using the specified URL, shared key, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewPageBlobClientWithSharedKey(blobURL string, cred *SharedKeyCredential, options *ClientOptions) (*PageBlobClient, error) {
	authPolicy := newSharedKeyCredPolicy(cred)
	conOptions := getConnectionOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	conn := internal.NewConnection(blobURL, conOptions)

	return &PageBlobClient{
		client: newPageBlobClient(conn.Endpoint(), conn.Pipeline()),
		BlobClient: BlobClient{
			client:    newBlobClient(conn.Endpoint(), conn.Pipeline()),
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

	endpoint := p.URL()
	pipeline := pb.client.pl
	return &PageBlobClient{
		client: newPageBlobClient(endpoint, pipeline),
		BlobClient: BlobClient{
			client:    newBlobClient(endpoint, pipeline),
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
	endpoint := p.URL()

	pipeline := pb.client.pl
	return &PageBlobClient{
		client: newPageBlobClient(endpoint, pipeline),
		BlobClient: BlobClient{
			client:    newBlobClient(endpoint, pipeline),
			sharedKey: pb.sharedKey,
		},
	}, nil
}

// Create creates a page blob of the specified length. Call PutPage to upload data to a page blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-blob.
func (pb *PageBlobClient) Create(ctx context.Context, size int64, o *PageBlobCreateOptions) (PageBlobCreateResponse, error) {
	createOptions, HTTPHeaders, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions := o.format()

	resp, err := pb.client.Create(ctx, 0, size, createOptions, HTTPHeaders, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions)

	return toPageBlobCreateResponse(resp), err
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

	uploadPagesOptions, leaseAccessConditions, cpkInfo, cpkScopeInfo, sequenceNumberAccessConditions, modifiedAccessConditions := options.format()

	resp, err := pb.client.UploadPages(ctx, count, body, uploadPagesOptions, leaseAccessConditions,
		cpkInfo, cpkScopeInfo, sequenceNumberAccessConditions, modifiedAccessConditions)

	return toPageBlobUploadPagesResponse(resp), err
}

// UploadPagesFromURL copies 1 or more pages from a source URL to the page blob.
// The sourceOffset specifies the start offset of source data to copy from.
// The destOffset specifies the start offset of data in page blob will be written to.
// The count must be a multiple of 512 bytes.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-page-from-url.
func (pb *PageBlobClient) UploadPagesFromURL(ctx context.Context, source string, sourceOffset, destOffset, count int64,
	options *PageBlobUploadPagesFromURLOptions) (PageBlobUploadPagesFromURLResponse, error) {

	uploadPagesFromURLOptions, cpkInfo, cpkScopeInfo, leaseAccessConditions, sequenceNumberAccessConditions, modifiedAccessConditions, sourceModifiedAccessConditions := options.format()

	resp, err := pb.client.UploadPagesFromURL(ctx, source, rangeToString(sourceOffset, count), 0,
		rangeToString(destOffset, count), uploadPagesFromURLOptions, cpkInfo, cpkScopeInfo, leaseAccessConditions,
		sequenceNumberAccessConditions, modifiedAccessConditions, sourceModifiedAccessConditions)

	return toPageBlobUploadPagesFromURLResponse(resp), err
}

// ClearPages frees the specified pages from the page blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-page.
func (pb *PageBlobClient) ClearPages(ctx context.Context, pageRange HttpRange, options *PageBlobClearPagesOptions) (PageBlobClearPagesResponse, error) {
	clearOptions := &pageBlobClientClearPagesOptions{
		Range: pageRange.format(),
	}

	leaseAccessConditions, cpkInfo, cpkScopeInfo, sequenceNumberAccessConditions, modifiedAccessConditions := options.format()

	resp, err := pb.client.ClearPages(ctx, 0, clearOptions, leaseAccessConditions, cpkInfo,
		cpkScopeInfo, sequenceNumberAccessConditions, modifiedAccessConditions)

	return toPageBlobClearPagesResponse(resp), err
}

// NewGetPageRangesPager returns the list of valid page ranges for a page blob or snapshot of a page blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-page-ranges.
func (pb *PageBlobClient) NewGetPageRangesPager(options *PageBlobGetPageRangesOptions) *runtime.Pager[PageBlobGetPageRangesResponse] {
	getPageRangesOptions, leaseAccessConditions, modifiedAccessConditions := options.format()

	return runtime.NewPager(runtime.PagingHandler[PageBlobGetPageRangesResponse]{
		More: func(page PageBlobGetPageRangesResponse) bool {
			if page.NextMarker == nil || len(*page.NextMarker) == 0 {
				return false
			}
			return true
		},
		Fetcher: func(ctx context.Context, page *PageBlobGetPageRangesResponse) (PageBlobGetPageRangesResponse, error) {
			var marker *string
			if page != nil {
				if page.NextMarker != nil {
					marker = page.NextMarker
				}
			} else {
				// If provided by the user, then use the one from options bag
				marker = getPageRangesOptions.Marker
			}

			req, err := pb.client.getPageRangesCreateRequest(ctx, &getPageRangesOptions, leaseAccessConditions, modifiedAccessConditions)
			if err != nil {
				return PageBlobGetPageRangesResponse{}, err
			}
			if marker != nil {
				queryValues, err := url.ParseQuery(req.Raw().URL.RawQuery)
				if err != nil {
					return PageBlobGetPageRangesResponse{}, err
				}
				queryValues.Set("marker", *marker)
				req.Raw().URL.RawQuery = queryValues.Encode()
				if err != nil {
					return PageBlobGetPageRangesResponse{}, err
				}
			}

			resp, err := pb.client.pl.Do(req)
			if err != nil {
				return PageBlobGetPageRangesResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PageBlobGetPageRangesResponse{}, runtime.NewResponseError(resp)
			}
			generatedResp, err := pb.client.getPageRangesHandleResponse(resp)
			return toPageBlobGetPageRangesResponse(generatedResp), err
		},
	})
}

// NewGetPageRangesDiffPager gets the collection of page ranges that differ between a specified snapshot and this page blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-page-ranges.
func (pb *PageBlobClient) NewGetPageRangesDiffPager(options *PageBlobGetPageRangesDiffOptions) *runtime.Pager[PageBlobGetPageRangesDiffResponse] {
	getPageRangesDiffOptions, leaseAccessConditions, modifiedAccessConditions := options.format()

	return runtime.NewPager(runtime.PagingHandler[PageBlobGetPageRangesDiffResponse]{
		More: func(page PageBlobGetPageRangesDiffResponse) bool {
			if page.NextMarker == nil || len(*page.NextMarker) == 0 {
				return false
			}
			return true
		},
		Fetcher: func(ctx context.Context, page *PageBlobGetPageRangesDiffResponse) (PageBlobGetPageRangesDiffResponse, error) {
			var marker *string
			if page != nil {
				if page.NextMarker != nil {
					marker = page.NextMarker
				}
			} else {
				// If provided by the user, then use the one from options bag
				marker = getPageRangesDiffOptions.Marker
			}

			req, err := pb.client.getPageRangesDiffCreateRequest(ctx, &getPageRangesDiffOptions, leaseAccessConditions, modifiedAccessConditions)
			if err != nil {
				return PageBlobGetPageRangesDiffResponse{}, err
			}
			if marker != nil {
				queryValues, err := url.ParseQuery(req.Raw().URL.RawQuery)
				if err != nil {
					return PageBlobGetPageRangesDiffResponse{}, err
				}
				queryValues.Set("marker", *marker)
				req.Raw().URL.RawQuery = queryValues.Encode()
				if err != nil {
					return PageBlobGetPageRangesDiffResponse{}, err
				}
			}

			resp, err := pb.client.pl.Do(req)
			if err != nil {
				return PageBlobGetPageRangesDiffResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PageBlobGetPageRangesDiffResponse{}, runtime.NewResponseError(resp)
			}
			generatedResp, err := pb.client.getPageRangesDiffHandleResponse(resp)
			return toPageBlobGetPageRangesDiffResponse(generatedResp), err
		},
	})
}

// Resize resizes the page blob to the specified size (which must be a multiple of 512).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-blob-properties.
func (pb *PageBlobClient) Resize(ctx context.Context, size int64, options *PageBlobResizeOptions) (PageBlobResizeResponse, error) {
	resizeOptions, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions := options.format()

	resp, err := pb.client.Resize(ctx, size, resizeOptions, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions)

	return toPageBlobResizeResponse(resp), err
}

// UpdateSequenceNumber sets the page blob's sequence number.
func (pb *PageBlobClient) UpdateSequenceNumber(ctx context.Context, options *PageBlobUpdateSequenceNumberOptions) (PageBlobUpdateSequenceNumberResponse, error) {
	actionType, updateOptions, lac, mac := options.format()
	resp, err := pb.client.UpdateSequenceNumber(ctx, *actionType, updateOptions, lac, mac)

	return toPageBlobUpdateSequenceNumberResponse(resp), err
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

	pageBlobCopyIncrementalOptions, modifiedAccessConditions := options.format()
	resp, err := pb.client.CopyIncremental(ctx, copySourceURL.String(), pageBlobCopyIncrementalOptions, modifiedAccessConditions)

	return toPageBlobCopyIncrementalResponse(resp), err
}
