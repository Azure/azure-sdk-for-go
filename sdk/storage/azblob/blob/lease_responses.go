//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blob

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

type AcquireResponse = generated.BlobClientAcquireLeaseResponse

type BreakResponse = generated.BlobClientBreakLeaseResponse

type ChangeResponse = generated.BlobClientChangeLeaseResponse

type ReleaseResponse = generated.BlobClientReleaseLeaseResponse

type RenewResponse = generated.BlobClientRenewLeaseResponse
