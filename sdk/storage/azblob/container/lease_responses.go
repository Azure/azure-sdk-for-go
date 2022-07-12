//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

type AcquireResponse = generated.ContainerClientAcquireLeaseResponse

type BreakResponse = generated.ContainerClientBreakLeaseResponse

type ChangeResponse = generated.ContainerClientChangeLeaseResponse

type ReleaseResponse = generated.ContainerClientReleaseLeaseResponse

type RenewResponse = generated.ContainerClientRenewLeaseResponse
