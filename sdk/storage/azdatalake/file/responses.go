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

// SetAccessControlResponse contains the response fields for the SetAccessControl operation.
type SetAccessControlResponse = path.SetAccessControlResponse

// SetAccessControlRecursiveResponse contains the response fields for the SetAccessControlRecursive operation.
type SetAccessControlRecursiveResponse = path.SetAccessControlRecursiveResponse

// UpdateAccessControlRecursiveResponse contains the response fields for the UpdateAccessControlRecursive operation.
type UpdateAccessControlRecursiveResponse = path.SetAccessControlRecursiveResponse

// RemoveAccessControlRecursiveResponse contains the response fields for the RemoveAccessControlRecursive operation.
type RemoveAccessControlRecursiveResponse = path.SetAccessControlRecursiveResponse

// GetPropertiesResponse contains the response fields for the GetProperties operation.
type GetPropertiesResponse = path.GetPropertiesResponse

// SetMetadataResponse contains the response fields for the SetMetadata operation.
type SetMetadataResponse = path.SetMetadataResponse

// SetHTTPHeadersResponse contains the response fields for the SetHTTPHeaders operation.
type SetHTTPHeadersResponse = path.SetHTTPHeadersResponse

// RenameResponse contains the response fields for the Rename operation.
type RenameResponse = path.CreateResponse

// GetAccessControlResponse contains the response fields for the GetAccessControl operation.
type GetAccessControlResponse = path.GetAccessControlResponse
