//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package lease

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"

// FileAcquireResponse contains the response from method FileClient.AcquireLease.
type FileAcquireResponse = generated.FileClientAcquireLeaseResponse

// FileBreakResponse contains the response from method FileClient.BreakLease.
type FileBreakResponse = generated.FileClientBreakLeaseResponse

// FileChangeResponse contains the response from method FileClient.ChangeLease.
type FileChangeResponse = generated.FileClientChangeLeaseResponse

// FileReleaseResponse contains the response from method FileClient.ReleaseLease.
type FileReleaseResponse = generated.FileClientReleaseLeaseResponse

// ShareAcquireResponse contains the response from method ShareClient.AcquireLease.
type ShareAcquireResponse = generated.ShareClientAcquireLeaseResponse

// ShareBreakResponse contains the response from method ShareClient.BreakLease.
type ShareBreakResponse = generated.ShareClientBreakLeaseResponse

// ShareChangeResponse contains the response from method ShareClient.ChangeLease.
type ShareChangeResponse = generated.ShareClientChangeLeaseResponse

// ShareReleaseResponse contains the response from method ShareClient.ReleaseLease.
type ShareReleaseResponse = generated.ShareClientReleaseLeaseResponse

// ShareRenewResponse contains the response from method ShareClient.RenewLease.
type ShareRenewResponse = generated.ShareClientRenewLeaseResponse
