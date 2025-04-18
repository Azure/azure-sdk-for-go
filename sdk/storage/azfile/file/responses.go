//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"io"
)

// CreateResponse contains the response from method Client.Create.
type CreateResponse = generated.FileClientCreateResponse

// DeleteResponse contains the response from method Client.Delete.
type DeleteResponse = generated.FileClientDeleteResponse

// RenameResponse contains the response from method Client.Rename.
type RenameResponse struct {
	generated.FileClientRenameResponse
}

// GetPropertiesResponse contains the response from method Client.GetProperties.
type GetPropertiesResponse = generated.FileClientGetPropertiesResponse

// SetMetadataResponse contains the response from method Client.SetMetadata.
type SetMetadataResponse = generated.FileClientSetMetadataResponse

// SetHTTPHeadersResponse contains the response from method Client.SetHTTPHeaders.
type SetHTTPHeadersResponse = generated.FileClientSetHTTPHeadersResponse

// StartCopyFromURLResponse contains the response from method Client.StartCopyFromURL.
type StartCopyFromURLResponse = generated.FileClientStartCopyResponse

// AbortCopyResponse contains the response from method Client.AbortCopy.
type AbortCopyResponse = generated.FileClientAbortCopyResponse

// DownloadResponse contains the response from method FileClient.Download.
type DownloadResponse = generated.FileClientDownloadResponse

// DownloadStreamResponse contains the response from method Client.DownloadStream.
// To read from the stream, read from the Body field, or call the NewRetryReader method.
type DownloadStreamResponse struct {
	DownloadResponse

	client                *Client
	getInfo               httpGetterInfo
	leaseAccessConditions *LeaseAccessConditions
}

// NewRetryReader constructs new RetryReader stream for reading data. If a connection fails while
// reading, it will make additional requests to reestablish a connection and continue reading.
// Pass nil for options to accept the default options.
// Callers of this method should not access the DownloadStreamResponse.Body field.
func (r *DownloadStreamResponse) NewRetryReader(ctx context.Context, options *RetryReaderOptions) *RetryReader {
	if options == nil {
		options = &RetryReaderOptions{}
	}

	return newRetryReader(ctx, r.Body, r.getInfo, func(ctx context.Context, getInfo httpGetterInfo) (io.ReadCloser, error) {
		options := DownloadStreamOptions{
			Range:                 getInfo.Range,
			LeaseAccessConditions: r.leaseAccessConditions,
		}
		resp, err := r.client.DownloadStream(ctx, &options)
		if err != nil {
			return nil, err
		}
		return resp.Body, err
	}, *options)
}

// ResizeResponse contains the response from method Client.Resize.
type ResizeResponse = generated.FileClientSetHTTPHeadersResponse

// UploadRangeResponse contains the response from method Client.UploadRange.
type UploadRangeResponse = generated.FileClientUploadRangeResponse

// ClearRangeResponse contains the response from method Client.ClearRange.
type ClearRangeResponse = generated.FileClientUploadRangeResponse

// UploadRangeFromURLResponse contains the response from method Client.UploadRangeFromURL.
type UploadRangeFromURLResponse = generated.FileClientUploadRangeFromURLResponse

// GetRangeListResponse contains the response from method Client.GetRangeList.
type GetRangeListResponse = generated.FileClientGetRangeListResponse

// ForceCloseHandlesResponse contains the response from method Client.ForceCloseHandles.
type ForceCloseHandlesResponse = generated.FileClientForceCloseHandlesResponse

// ListHandlesResponse contains the response from method Client.ListHandles.
type ListHandlesResponse = generated.FileClientListHandlesResponse

// ListHandlesSegmentResponse - An enumeration of handles.
type ListHandlesSegmentResponse = generated.ListHandlesResponse

// CreateHardLinkResponse contains response from method Client.CreateHardLink
type CreateHardLinkResponse = generated.FileClientCreateHardLinkResponse
