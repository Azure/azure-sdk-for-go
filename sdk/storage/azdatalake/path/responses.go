//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

type CreateResponse = generated.PathClientCreateResponse

type DeleteResponse = generated.PathClientDeleteResponse

type SetAccessControlResponse = generated.PathClientSetAccessControlResponse

type SetAccessControlRecursiveResponse = generated.PathClientSetAccessControlRecursiveResponse

type GetPropertiesResponse = blob.GetPropertiesResponse

type SetMetadataResponse = blob.SetMetadataResponse

type SetHTTPHeadersResponse = blob.SetHTTPHeadersResponse
