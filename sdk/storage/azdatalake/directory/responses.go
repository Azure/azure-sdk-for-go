//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/path"
)

// RenameResponse contains the response fields for the Create operation.
type RenameResponse struct {
	Response           CreateResponse
	NewDirectoryClient *Client
}

// SetAccessControlRecursiveResponse contains the response fields for the SetAccessControlRecursive operation.
type SetAccessControlRecursiveResponse = generated.PathClientSetAccessControlRecursiveResponse

// ========================================== path imports ===========================================================

// SetAccessControlResponse contains the response fields for the SetAccessControl operation.
type SetAccessControlResponse = path.SetAccessControlResponse

// SetHTTPHeadersResponse contains the response from method Client.SetHTTPHeaders.
type SetHTTPHeadersResponse = path.SetHTTPHeadersResponse

// GetAccessControlResponse contains the response fields for the GetAccessControl operation.
type GetAccessControlResponse = path.GetAccessControlResponse

// GetPropertiesResponse contains the response fields for the GetProperties operation.
type GetPropertiesResponse = path.GetPropertiesResponse

// SetMetadataResponse contains the response fields for the SetMetadata operation.
type SetMetadataResponse = path.SetMetadataResponse

// CreateResponse contains the response fields for the Create operation.
type CreateResponse = path.CreateResponse

// DeleteResponse contains the response fields for the Delete operation.
type DeleteResponse = path.DeleteResponse

// UpdateAccessControlRecursiveResponse contains the response fields for the UpdateAccessControlRecursive operation.
type UpdateAccessControlRecursiveResponse = path.UpdateAccessControlResponse

// RemoveAccessControlRecursiveResponse contains the response fields for the RemoveAccessControlRecursive operation.
type RemoveAccessControlRecursiveResponse = path.RemoveAccessControlResponse
