//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"time"
)

// SetExpiryResponse contains the response fields for the SetExpiry operation.
type SetExpiryResponse = generated.PathClientSetExpiryResponse

// CreateResponse contains the response fields for the Create operation.
type CreateResponse = generated.PathClientCreateResponse

// DeleteResponse contains the response fields for the Delete operation.
type DeleteResponse = generated.PathClientDeleteResponse

// SetAccessControlResponse contains the response fields for the SetAccessControl operation.
type SetAccessControlResponse = generated.PathClientSetAccessControlResponse

// UpdateAccessControlResponse contains the response fields for the UpdateAccessControlRecursive operation.
type UpdateAccessControlResponse = generated.PathClientSetAccessControlRecursiveResponse

// RemoveAccessControlResponse contains the response fields for the RemoveAccessControlRecursive operation.
type RemoveAccessControlResponse = generated.PathClientSetAccessControlRecursiveResponse

// GetAccessControlResponse contains the response fields for the GetAccessControl operation.
type GetAccessControlResponse = generated.PathClientGetPropertiesResponse

// GetPropertiesResponse contains the response fields for the GetProperties operation.
type GetPropertiesResponse = blob.GetPropertiesResponse

// SetMetadataResponse contains the response fields for the SetMetadata operation.
type SetMetadataResponse = blob.SetMetadataResponse

// RenameResponse contains the response fields for the Create operation.
type RenameResponse = generated.PathClientCreateResponse

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

func formatSetHTTPHeadersResponse(r *SetHTTPHeadersResponse, blobResp *blob.SetHTTPHeadersResponse) {
	r.ClientRequestID = blobResp.ClientRequestID
	r.Date = blobResp.Date
	r.ETag = blobResp.ETag
	r.LastModified = blobResp.LastModified
	r.RequestID = blobResp.RequestID
	r.Version = blobResp.Version
}
