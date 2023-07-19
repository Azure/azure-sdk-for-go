//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"time"
)

// SetAccessControlResponse contains the response fields for the SetAccessControl operation.
type SetAccessControlResponse = generated.PathClientSetAccessControlResponse

// GetAccessControlResponse contains the response fields for the GetAccessControl operation.
type GetAccessControlResponse = generated.PathClientGetPropertiesResponse

// GetPropertiesResponse contains the response fields for the GetProperties operation.
type GetPropertiesResponse = blob.GetPropertiesResponse

// SetMetadataResponse contains the response fields for the SetMetadata operation.
type SetMetadataResponse = blob.SetMetadataResponse

//// SetHTTPHeadersResponse contains the response fields for the SetHTTPHeaders operation.
//type SetHTTPHeadersResponse = blob.SetHTTPHeadersResponse
// we need to remove the blob sequence number from the response

// SetHTTPHeadersResponse contains the response from method Client.SetHTTPHeaders.
type SetHTTPHeadersResponse struct {
	// ClientRequestID contains the information returned from the x-ms-client-request-id header response.
	ClientRequestID *string

	// Date contains the information returned from the Date header response.
	Date *time.Time

	// ETag contains the information returned from the ETag header response.
	ETag *azcore.ETag

	// LastModified contains the information returned from the Last-Modified header response.
	LastModified *time.Time

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// Version contains the information returned from the x-ms-version header response.
	Version *string
}

func FormatSetHTTPHeadersResponse(r *SetHTTPHeadersResponse, blobResp *blob.SetHTTPHeadersResponse) {
	r.ClientRequestID = blobResp.ClientRequestID
	r.Date = blobResp.Date
	r.ETag = blobResp.ETag
	r.LastModified = blobResp.LastModified
	r.RequestID = blobResp.RequestID
	r.Version = blobResp.Version
}
