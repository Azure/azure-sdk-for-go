//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

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

type DeleteResponse = generated.BlobClientDeleteResponse

type UndeleteResponse = generated.BlobClientUndeleteResponse

type SetTierResponse = generated.BlobClientSetTierResponse

type GetPropertiesResponse = generated.BlobClientGetPropertiesResponse

type SetHTTPHeadersResponse = generated.BlobClientSetHTTPHeadersResponse

type SetMetadataResponse = generated.BlobClientSetMetadataResponse

type CreateSnapshotResponse = generated.BlobClientCreateSnapshotResponse

type StartCopyFromURLResponse = generated.BlobClientStartCopyFromURLResponse

type AbortCopyFromURLResponse = generated.BlobClientAbortCopyFromURLResponse

type SetTagsResponse = generated.BlobClientSetTagsResponse

type GetTagsResponse = generated.BlobClientGetTagsResponse

type CopyFromURLResponse = generated.BlobClientCopyFromURLResponse
