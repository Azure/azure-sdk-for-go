// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package share

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// CreateResponse contains the response from method Client.Create.
type CreateResponse = generated.ShareClientCreateResponse

// DeleteResponse contains the response from method Client.Delete.
type DeleteResponse = generated.ShareClientDeleteResponse

// RestoreResponse contains the response from method Client.Restore.
type RestoreResponse = generated.ShareClientRestoreResponse

// GetPropertiesResponse contains the response from method Client.GetProperties.
type GetPropertiesResponse = generated.ShareClientGetPropertiesResponse

// SetPropertiesResponse contains the response from method Client.SetProperties.
type SetPropertiesResponse = generated.ShareClientSetPropertiesResponse

// CreateSnapshotResponse contains the response from method Client.CreateSnapshot.
type CreateSnapshotResponse = generated.ShareClientCreateSnapshotResponse

// GetAccessPolicyResponse contains the response from method Client.GetAccessPolicy.
type GetAccessPolicyResponse = generated.ShareClientGetAccessPolicyResponse

// SetAccessPolicyResponse contains the response from method Client.SetAccessPolicy.
type SetAccessPolicyResponse = generated.ShareClientSetAccessPolicyResponse

// CreatePermissionResponse contains the response from method Client.CreatePermission.
type CreatePermissionResponse = generated.ShareClientCreatePermissionResponse

// GetPermissionResponse contains the response from method Client.GetPermission.
type GetPermissionResponse = generated.ShareClientGetPermissionResponse

// SetMetadataResponse contains the response from method Client.SetMetadata.
type SetMetadataResponse = generated.ShareClientSetMetadataResponse

// GetStatisticsResponse contains the response from method Client.GetStatistics.
type GetStatisticsResponse = generated.ShareClientGetStatisticsResponse
