//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// CreateResponse contains the response from method Client.Create.
type CreateResponse = generated.DirectoryClientCreateResponse

// DeleteResponse contains the response from method Client.Delete.
type DeleteResponse = generated.DirectoryClientDeleteResponse

// RenameResponse contains the response from method Client.Rename.
type RenameResponse struct {
	generated.DirectoryClientRenameResponse
}

// GetPropertiesResponse contains the response from method Client.GetProperties.
type GetPropertiesResponse = generated.DirectoryClientGetPropertiesResponse

// SetPropertiesResponse contains the response from method Client.SetProperties.
type SetPropertiesResponse = generated.DirectoryClientSetPropertiesResponse

// SetMetadataResponse contains the response from method Client.SetMetadata.
type SetMetadataResponse = generated.DirectoryClientSetMetadataResponse

// ListFilesAndDirectoriesResponse contains the response from method Client.NewListFilesAndDirectoriesPager.
type ListFilesAndDirectoriesResponse = generated.DirectoryClientListFilesAndDirectoriesSegmentResponse

// ListFilesAndDirectoriesSegmentResponse - An enumeration of directories and files.
type ListFilesAndDirectoriesSegmentResponse = generated.ListFilesAndDirectoriesSegmentResponse

// ListHandlesResponse contains the response from method Client.ListHandles.
type ListHandlesResponse = generated.DirectoryClientListHandlesResponse

// ListHandlesSegmentResponse - An enumeration of handles.
type ListHandlesSegmentResponse = generated.ListHandlesResponse

// ForceCloseHandlesResponse contains the response from method Client.ForceCloseHandles.
type ForceCloseHandlesResponse = generated.DirectoryClientForceCloseHandlesResponse
