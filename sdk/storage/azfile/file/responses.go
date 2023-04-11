//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
)

// CreateResponse contains the response from method Client.Create.
type CreateResponse = generated.FileClientCreateResponse

// DeleteResponse contains the response from method Client.Delete.
type DeleteResponse = generated.FileClientDeleteResponse

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
	client  *Client
	getInfo httpGetterInfo
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
