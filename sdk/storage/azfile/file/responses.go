//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// CreateResponse contains the response from method Client.Create.
type CreateResponse = generated.FileClientCreateResponse

// DeleteResponse contains the response from method Client.Delete.
type DeleteResponse = generated.FileClientDeleteResponse

// GetPropertiesResponse contains the response from method Client.GetProperties.
type GetPropertiesResponse = generated.FileClientGetPropertiesResponse

// SetMetadataResponse contains the response from method Client.SetMetadata.
type SetMetadataResponse = generated.FileClientSetMetadataResponse

// SetPropertiesResponse contains the response from method Client.SetProperties.
type SetPropertiesResponse = generated.FileClientSetHTTPHeadersResponse

// StartCopyFromURLResponse contains the response from method Client.StartCopyFromURL.
type StartCopyFromURLResponse = generated.FileClientStartCopyResponse

// DownloadResponse contains the response from method Client.Download.
type DownloadResponse = generated.FileClientDownloadResponse

// AbortCopyResponse contains the response from method Client.AbortCopy.
type AbortCopyResponse = generated.FileClientAbortCopyResponse

// AcquireLeaseResponse contains the response from method Client.AcquireLease.
type AcquireLeaseResponse = generated.FileClientAcquireLeaseResponse

// BreakLeaseResponse contains the response from method Client.BreakLease.
type BreakLeaseResponse = generated.FileClientBreakLeaseResponse

// ChangeLeaseResponse contains the response from method Client.ChangeLease.
type ChangeLeaseResponse = generated.FileClientChangeLeaseResponse

// ReleaseLeaseResponse contains the response from method Client.ReleaseLease.
type ReleaseLeaseResponse = generated.FileClientReleaseLeaseResponse
