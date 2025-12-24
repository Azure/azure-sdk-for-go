// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package service

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// CreateShareResponse contains the response from method share.Client.Create.
type CreateShareResponse = generated.ShareClientCreateResponse

// DeleteShareResponse contains the response from method share.Client.Delete.
type DeleteShareResponse = generated.ShareClientDeleteResponse

// RestoreShareResponse contains the response from method share.Client.Restore.
type RestoreShareResponse = generated.ShareClientRestoreResponse

// GetPropertiesResponse contains the response from method Client.GetProperties.
type GetPropertiesResponse = generated.ServiceClientGetPropertiesResponse

// SetPropertiesResponse contains the response from method Client.SetProperties.
type SetPropertiesResponse = generated.ServiceClientSetPropertiesResponse

// ListSharesSegmentResponse contains the response from method Client.NewListSharesPager.
type ListSharesSegmentResponse = generated.ServiceClientListSharesSegmentResponse

// ListSharesResponse - An enumeration of shares.
type ListSharesResponse = generated.ListSharesResponse
