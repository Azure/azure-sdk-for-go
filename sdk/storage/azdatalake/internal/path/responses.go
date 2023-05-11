//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// CreateResponse contains the response fields for the Create operation.
type CreateResponse = generated.PathClientCreateResponse

// DeleteResponse contains the response fields for the Delete operation.
type DeleteResponse = generated.PathClientDeleteResponse

// SetAccessControlResponse contains the response fields for the SetAccessControl operation.
type SetAccessControlResponse = generated.PathClientSetAccessControlResponse

// SetAccessControlRecursiveResponse contains the response fields for the SetAccessControlRecursive operation.
type SetAccessControlRecursiveResponse = generated.PathClientSetAccessControlRecursiveResponse

// GetPropertiesResponse contains the response fields for the GetProperties operation.
type GetPropertiesResponse = blob.GetPropertiesResponse

// SetMetadataResponse contains the response fields for the SetMetadata operation.
type SetMetadataResponse = blob.SetMetadataResponse

// SetHTTPHeadersResponse contains the response fields for the SetHTTPHeaders operation.
type SetHTTPHeadersResponse = blob.SetHTTPHeadersResponse
