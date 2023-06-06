//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package filesystem

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// CreateResponse contains the response from method FilesystemClient.Create.
type CreateResponse = container.CreateResponse

// DeleteResponse contains the response from method FilesystemClient.Delete.
type DeleteResponse = container.DeleteResponse

// SetMetadataResponse contains the response from method FilesystemClient.SetMetadata.
type SetMetadataResponse = container.SetMetadataResponse

// SetAccessPolicyResponse contains the response from method FilesystemClient.SetAccessPolicy.
type SetAccessPolicyResponse = container.SetAccessPolicyResponse

// GetAccessPolicyResponse contains the response from method FilesystemClient.GetAccessPolicy.
type GetAccessPolicyResponse = container.GetAccessPolicyResponse

// GetPropertiesResponse contains the response from method FilesystemClient.GetProperties.
type GetPropertiesResponse = generated.FileSystemClientGetPropertiesResponse

// ListPathsSegmentResponse contains the response from method FilesystemClient.ListPathsSegment.
type ListPathsSegmentResponse = generated.FileSystemClientListPathsResponse

// ListDeletedPathsSegmentResponse contains the response from method FilesystemClient.ListPathsSegment.
type ListDeletedPathsSegmentResponse = generated.FileSystemClientListBlobHierarchySegmentResponse

// UndeletePathResponse contains the response from method FilesystemClient.UndeletePath.
type UndeletePathResponse = generated.PathClientUndeleteResponse
