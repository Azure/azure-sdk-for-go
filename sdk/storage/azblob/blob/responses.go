//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import (
	"context"
	"io"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

// DownloadResponse wraps AutoRest generated BlobDownloadResponse and helps to provide info for retry.
type DownloadResponse struct {
	generated.BlobClientDownloadResponse
	ctx                    context.Context
	b                      *Client
	getInfo                HTTPGetterInfo
	ObjectReplicationRules []ObjectReplicationPolicy
}

// NewRetryReader constructs new RetryReader stream for reading data. If a connection fails
// while reading, it will make additional requests to reestablish a connection and
// continue reading. Specifying a RetryReaderOption's with MaxRetryRequests set to 0
// (the default), returns the original response body and no retries will be performed.
// Pass in nil for options to accept the default options.
func (r *DownloadResponse) NewRetryReader(ctx context.Context, options *RetryReaderOptions) io.ReadCloser {
	if options == nil {
		options = &RetryReaderOptions{}
	}

	if options.MaxRetryRequests == 0 { // No additional retries
		return r.BlobClientDownloadResponse.Body
	}
	return NewRetryReader(ctx, r.Body, r.getInfo, func(ctx context.Context, getInfo HTTPGetterInfo) (io.ReadCloser, error) {
		accessConditions := &AccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: &getInfo.ETag},
		}
		options := DownloadOptions{
			Offset:           &getInfo.Offset,
			Count:            &getInfo.Count,
			AccessConditions: accessConditions,
			CpkInfo:          options.CpkInfo,
			//CpkScopeInfo:         options.CpkScopeInfo,
		}
		resp, err := r.b.Download(ctx, &options)
		if err != nil {
			return nil, err
		}
		return resp.Body, err
	}, options)
}

// BodyReader constructs new RetryReader stream for reading data. If a connection fails
// while reading, it will make additional requests to reestablish a connection and
// continue reading. Specifying a RetryReaderOption's with MaxRetryRequests set to 0
// (the default), returns the original response body and no retries will be performed.
// Pass in nil for options to accept the default options.
func (r *DownloadResponse) BodyReader(options *RetryReaderOptions) io.ReadCloser {
	if options == nil {
		options = &RetryReaderOptions{}
	}

	if options.MaxRetryRequests == 0 { // No additional retries
		return r.BlobClientDownloadResponse.Body
	}
	return NewRetryReader(r.ctx, r.BlobClientDownloadResponse.Body, r.getInfo, func(ctx context.Context, getInfo HTTPGetterInfo) (io.ReadCloser, error) {
		accessConditions := &AccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: &getInfo.ETag},
		}
		options := DownloadOptions{
			Offset:           &getInfo.Offset,
			Count:            &getInfo.Count,
			AccessConditions: accessConditions,
			CpkInfo:          options.CpkInfo,
			CpkScopeInfo:     options.CpkScopeInfo,
		}
		resp, err := r.b.Download(ctx, &options)
		if err != nil {
			return nil, err
		}
		return resp.BlobClientDownloadResponse.Body, err
	}, options)
}

// DeleteResponse contains the response from method BlobClient.Delete.
type DeleteResponse = generated.BlobClientDeleteResponse

// UndeleteResponse contains the response from method BlobClient.Undelete.
type UndeleteResponse = generated.BlobClientUndeleteResponse

// SetTierResponse contains the response from method BlobClient.SetTier.
type SetTierResponse = generated.BlobClientSetTierResponse

// GetPropertiesResponse contains the response from method BlobClient.GetProperties.
type GetPropertiesResponse = generated.BlobClientGetPropertiesResponse

// SetHTTPHeadersResponse contains the response from method BlobClient.SetHTTPHeaders.
type SetHTTPHeadersResponse = generated.BlobClientSetHTTPHeadersResponse

// SetMetadataResponse contains the response from method BlobClient.SetMetadata.
type SetMetadataResponse = generated.BlobClientSetMetadataResponse

// CreateSnapshotResponse contains the response from method BlobClient.CreateSnapshot.
type CreateSnapshotResponse = generated.BlobClientCreateSnapshotResponse

// StartCopyFromURLResponse contains the response from method BlobClient.StartCopyFromURL.
type StartCopyFromURLResponse = generated.BlobClientStartCopyFromURLResponse

// AbortCopyFromURLResponse contains the response from method BlobClient.AbortCopyFromURL.
type AbortCopyFromURLResponse = generated.BlobClientAbortCopyFromURLResponse

// SetTagsResponse contains the response from method BlobClient.SetTags.
type SetTagsResponse = generated.BlobClientSetTagsResponse

// GetTagsResponse contains the response from method BlobClient.GetTags.
type GetTagsResponse = generated.BlobClientGetTagsResponse

// CopyFromURLResponse contains the response from method BlobClient.CopyFromURL.
type CopyFromURLResponse = generated.BlobClientCopyFromURLResponse

// AcquireLeaseResponse contains the response from method BlobClient.AcquireLease.
type AcquireLeaseResponse = generated.BlobClientAcquireLeaseResponse

// BreakLeaseResponse contains the response from method BlobClient.BreakLease.
type BreakLeaseResponse = generated.BlobClientBreakLeaseResponse

// ChangeLeaseResponse contains the response from method BlobClient.ChangeLease.
type ChangeLeaseResponse = generated.BlobClientChangeLeaseResponse

// ReleaseLeaseResponse contains the response from method BlobClient.ReleaseLease.
type ReleaseLeaseResponse = generated.BlobClientReleaseLeaseResponse

// RenewLeaseResponse contains the response from method BlobClient.RenewLease.
type RenewLeaseResponse = generated.BlobClientRenewLeaseResponse
