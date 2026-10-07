// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
)

// LayoutAwareRouting determines whether the parallel range reads of a DownloadBuffer or
// DownloadFile are routed to the endpoint that serves each range, based on the file's layout. It
// is a performance optimization only: the data downloaded is the same whatever the mode.
type LayoutAwareRouting = blob.LayoutAwareRouting

const (
	// LayoutAwareRoutingAuto is the zero value, and therefore the default when no value is
	// specified. The client library decides whether layout aware routing is used, and that
	// decision may change in a future release. Currently, LayoutAwareRoutingAuto resolves to
	// LayoutAwareRoutingDisabled; use LayoutAwareRoutingEnabled to opt in.
	LayoutAwareRoutingAuto = blob.LayoutAwareRoutingAuto

	// LayoutAwareRoutingEnabled opts in to layout aware routing. When the initial read of a download
	// carries the layout download hint and data remains, the layout of the remaining range is fetched
	// and cached (with automatic background refresh), and each remaining range is read from the
	// endpoint that serves it.
	LayoutAwareRoutingEnabled = blob.LayoutAwareRoutingEnabled

	// LayoutAwareRoutingDisabled never routes by layout; every read goes to the client's configured
	// endpoint.
	LayoutAwareRoutingDisabled = blob.LayoutAwareRoutingDisabled
)

// PossibleLayoutAwareRoutingValues returns the possible values for the LayoutAwareRouting const type.
func PossibleLayoutAwareRoutingValues() []LayoutAwareRouting {
	return blob.PossibleLayoutAwareRoutingValues()
}

// Layout is a page of a file's layout: the byte ranges making up the file and the storage
// endpoints that serve them.
type Layout = blob.Layout

// LayoutRanges contains the ranges of a Layout.
type LayoutRanges = blob.LayoutRanges

// LayoutRange is a range of a file, inclusive of Start and End, and the index of the endpoint in
// LayoutEndpoints that serves it.
type LayoutRange = blob.LayoutRange

// LayoutEndpoints contains the endpoints of a Layout.
type LayoutEndpoints = blob.LayoutEndpoints

// LayoutEndpoint is an endpoint that serves ranges of a file, referenced from LayoutRange by Index.
// Pass its Value as DownloadStreamOptions.LayoutEndpoint to read a range from it.
type LayoutEndpoint = blob.LayoutEndpoint

// GetLayoutOptions contains the optional parameters for the Client.GetLayoutPager method.
type GetLayoutOptions struct {
	// Range restricts the layout to a range of the file. The default is the whole file.
	Range *HTTPRange
	// AccessConditions contains parameters for accessing the file. Unless they include an If-Match
	// condition, every page after the first is locked to the ETag the first page returned.
	AccessConditions *AccessConditions
	// CPKInfo contains a group of parameters for client provided encryption key.
	CPKInfo *CPKInfo
}

func (o *GetLayoutOptions) format() *blob.GetLayoutOptions {
	if o == nil {
		return nil
	}
	opts := &blob.GetLayoutOptions{
		AccessConditions: exported.FormatBlobAccessConditions(o.AccessConditions),
	}
	if o.Range != nil {
		opts.Range = *o.Range
	}
	if o.CPKInfo != nil {
		opts.CPKInfo = &blob.CPKInfo{
			EncryptionKey:       o.CPKInfo.EncryptionKey,
			EncryptionKeySHA256: o.CPKInfo.EncryptionKeySHA256,
			EncryptionAlgorithm: (*blob.EncryptionAlgorithmType)(o.CPKInfo.EncryptionAlgorithm),
		}
	}
	return opts
}

// GetLayoutResponse contains a page of a file's layout returned by Client.GetLayoutPager, along
// with the file's properties.
type GetLayoutResponse struct {
	// Layout contains the ranges and endpoints of this page, and the marker for the next one.
	Layout

	// AcceptRanges contains the information returned from the Accept-Ranges header response.
	AcceptRanges *string
	// AccessTier contains the information returned from the x-ms-access-tier header response.
	AccessTier *string
	// AccessTierChangeTime contains the information returned from the x-ms-access-tier-change-time header response.
	AccessTierChangeTime *time.Time
	// AccessTierInferred contains the information returned from the x-ms-access-tier-inferred header response.
	AccessTierInferred *bool
	// ArchiveStatus contains the information returned from the x-ms-archive-status header response.
	ArchiveStatus *string
	// CacheControl contains the information returned from the Cache-Control header response.
	CacheControl *string
	// ClientRequestID contains the information returned from the x-ms-client-request-id header response.
	ClientRequestID *string
	// ContentDisposition contains the information returned from the Content-Disposition header response.
	ContentDisposition *string
	// ContentEncoding contains the information returned from the Content-Encoding header response.
	ContentEncoding *string
	// ContentLanguage contains the information returned from the Content-Language header response.
	ContentLanguage *string
	// ContentLength contains the length of the response body, the layout page.
	ContentLength *int64
	// ContentMD5 contains the information returned from the Content-MD5 header response.
	ContentMD5 []byte
	// ContentType contains the information returned from the Content-Type header response.
	ContentType *string
	// CopyCompletionTime contains the information returned from the x-ms-copy-completion-time header response.
	CopyCompletionTime *time.Time
	// CopyID contains the information returned from the x-ms-copy-id header response.
	CopyID *string
	// CopyProgress contains the information returned from the x-ms-copy-progress header response.
	CopyProgress *string
	// CopySource contains the information returned from the x-ms-copy-source header response.
	CopySource *string
	// CopyStatus contains the information returned from the x-ms-copy-status header response.
	CopyStatus *CopyStatusType
	// CopyStatusDescription contains the information returned from the x-ms-copy-status-description header response.
	CopyStatusDescription *string
	// CreationTime contains the information returned from the x-ms-creation-time header response.
	CreationTime *time.Time
	// Date contains the information returned from the Date header response.
	Date *time.Time
	// ETag contains the information returned from the ETag header response.
	ETag *azcore.ETag
	// EncryptionKeySHA256 contains the information returned from the x-ms-encryption-key-sha256 header response.
	EncryptionKeySHA256 *string
	// EncryptionScope contains the information returned from the x-ms-encryption-scope header response.
	EncryptionScope *string
	// ExpiresOn contains the information returned from the x-ms-expiry-time header response.
	ExpiresOn *time.Time
	// FileContentEncoding contains the content encoding of the file.
	FileContentEncoding *string
	// FileContentLength contains the length of the file.
	FileContentLength *int64
	// FileContentMD5 contains the MD5 of the file's content.
	FileContentMD5 []byte
	// FileContentType contains the content type of the file.
	FileContentType *string
	// FileCreationTime contains the creation time of the file.
	FileCreationTime *time.Time
	// IsIncrementalCopy contains the information returned from the x-ms-incremental-copy header response.
	IsIncrementalCopy *bool
	// IsServerEncrypted contains the information returned from the x-ms-server-encrypted header response.
	IsServerEncrypted *bool
	// LastModified contains the information returned from the Last-Modified header response.
	LastModified *time.Time
	// LeaseDuration contains the information returned from the x-ms-lease-duration header response.
	LeaseDuration *DurationType
	// LeaseState contains the information returned from the x-ms-lease-state header response.
	LeaseState *StateType
	// LeaseStatus contains the information returned from the x-ms-lease-status header response.
	LeaseStatus *StatusType
	// Metadata contains the information returned from the x-ms-meta header response.
	Metadata map[string]*string
	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string
	// SmartAccessTier contains the information returned from the x-ms-smart-access-tier header response.
	SmartAccessTier *string
	// Version contains the information returned from the x-ms-version header response.
	Version *string
}

func formatGetLayoutResponse(r blob.GetLayoutResponse) GetLayoutResponse {
	return GetLayoutResponse{
		Layout:                r.BlobLayout,
		AcceptRanges:          r.AcceptRanges,
		AccessTier:            r.AccessTier,
		AccessTierChangeTime:  r.AccessTierChangeTime,
		AccessTierInferred:    r.AccessTierInferred,
		ArchiveStatus:         r.ArchiveStatus,
		CacheControl:          r.CacheControl,
		ClientRequestID:       r.ClientRequestID,
		ContentDisposition:    r.ContentDisposition,
		ContentEncoding:       r.ContentEncoding,
		ContentLanguage:       r.ContentLanguage,
		ContentLength:         r.ContentLength,
		ContentMD5:            r.ContentMD5,
		ContentType:           r.ContentType,
		CopyCompletionTime:    r.CopyCompletionTime,
		CopyID:                r.CopyID,
		CopyProgress:          r.CopyProgress,
		CopySource:            r.CopySource,
		CopyStatus:            r.CopyStatus,
		CopyStatusDescription: r.CopyStatusDescription,
		CreationTime:          r.CreationTime,
		Date:                  r.Date,
		ETag:                  r.ETag,
		EncryptionKeySHA256:   r.EncryptionKeySHA256,
		EncryptionScope:       r.EncryptionScope,
		ExpiresOn:             r.ExpiresOn,
		FileContentEncoding:   r.BlobContentEncoding,
		FileContentLength:     r.BlobContentLength,
		FileContentMD5:        r.BlobContentMD5,
		FileContentType:       r.BlobContentType,
		FileCreationTime:      r.BlobCreationTime,
		IsIncrementalCopy:     r.IsIncrementalCopy,
		IsServerEncrypted:     r.IsServerEncrypted,
		LastModified:          r.LastModified,
		LeaseDuration:         r.LeaseDuration,
		LeaseState:            r.LeaseState,
		LeaseStatus:           r.LeaseStatus,
		Metadata:              r.Metadata,
		RequestID:             r.RequestID,
		SmartAccessTier:       r.SmartAccessTier,
		Version:               r.Version,
	}
}

// GetLayoutPager returns the file's layout: the set of byte ranges making up the file and the
// storage endpoint that serves each one. Pass the endpoint covering a given offset as
// DownloadStreamOptions.LayoutEndpoint to route that read for better locality.
//
// The layout comes from the blob endpoint's Get Blob Layout operation; there is no DFS equivalent.
// A single enumeration describes the whole requested range, so callers implementing a custom
// chunked download should enumerate once and reuse the result rather than paging per chunk. A
// file's layout can change over time; refresh a cached layout roughly every 5 minutes, which is the
// interval Client.DownloadBuffer and Client.DownloadFile use internally.
//
// Unless options.AccessConditions includes an If-Match condition, every page after the first is
// locked to the ETag the first page returned, so one enumeration describes one version of the file.
func (f *Client) GetLayoutPager(options *GetLayoutOptions) *runtime.Pager[GetLayoutResponse] {
	blobPager := f.blobClient().GetLayoutPager(options.format())
	return runtime.NewPager(runtime.PagingHandler[GetLayoutResponse]{
		More: func(GetLayoutResponse) bool {
			return blobPager.More()
		},
		Fetcher: func(ctx context.Context, _ *GetLayoutResponse) (GetLayoutResponse, error) {
			resp, err := blobPager.NextPage(ctx)
			if err != nil {
				return GetLayoutResponse{}, exported.ConvertToDFSError(err)
			}
			return formatGetLayoutResponse(resp), nil
		},
	})
}
