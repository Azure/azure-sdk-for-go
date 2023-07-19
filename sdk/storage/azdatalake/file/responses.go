//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/path"
)

// SetExpiryResponse contains the response fields for the SetExpiry operation.
type SetExpiryResponse = generated.PathClientSetExpiryResponse

// CreateResponse contains the response fields for the Create operation.
type CreateResponse = generated.PathClientCreateResponse

// DeleteResponse contains the response fields for the Delete operation.
type DeleteResponse = generated.PathClientDeleteResponse

// UpdateAccessControlResponse contains the response fields for the UpdateAccessControlRecursive operation.
type UpdateAccessControlResponse = generated.PathClientSetAccessControlRecursiveResponse

// RemoveAccessControlResponse contains the response fields for the RemoveAccessControlRecursive operation.
type RemoveAccessControlResponse = generated.PathClientSetAccessControlRecursiveResponse

// RenameResponse contains the response fields for the Create operation.
type RenameResponse = generated.PathClientCreateResponse

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
