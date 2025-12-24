// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// FileAcquireResponse contains the response from method FileClient.Acquire.
type FileAcquireResponse = generated.FileClientAcquireLeaseResponse

// FileBreakResponse contains the response from method FileClient.Break.
type FileBreakResponse = generated.FileClientBreakLeaseResponse

// FileChangeResponse contains the response from method FileClient.Change.
type FileChangeResponse = generated.FileClientChangeLeaseResponse

// FileReleaseResponse contains the response from method FileClient.Release.
type FileReleaseResponse = generated.FileClientReleaseLeaseResponse

// ShareAcquireResponse contains the response from method ShareClient.Acquire.
type ShareAcquireResponse = generated.ShareClientAcquireLeaseResponse

// ShareBreakResponse contains the response from method ShareClient.Break.
type ShareBreakResponse = generated.ShareClientBreakLeaseResponse

// ShareChangeResponse contains the response from method ShareClient.Change.
type ShareChangeResponse = generated.ShareClientChangeLeaseResponse

// ShareReleaseResponse contains the response from method ShareClient.Release.
type ShareReleaseResponse = generated.ShareClientReleaseLeaseResponse

// ShareRenewResponse contains the response from method ShareClient.Renew.
type ShareRenewResponse = generated.ShareClientRenewLeaseResponse
