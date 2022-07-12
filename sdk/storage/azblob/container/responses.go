//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

// CreateResponse contains the response from method Client.Create.
type CreateResponse = generated.ContainerClientCreateResponse

type DeleteResponse = generated.ContainerClientDeleteResponse

type GetPropertiesResponse = generated.ContainerClientGetPropertiesResponse

type ListBlobsFlatResponse = generated.ContainerClientListBlobFlatSegmentResponse

type ListBlobsHierarchyResponse = generated.ContainerClientListBlobHierarchySegmentResponse

type SetMetadataResponse = generated.ContainerClientSetMetadataResponse

type GetAccessPolicyResponse = generated.ContainerClientGetAccessPolicyResponse

type SetAccessPolicyResponse = generated.ContainerClientSetAccessPolicyResponse

// --------------------

type AcquireResponse = generated.ContainerClientAcquireLeaseResponse

type BreakResponse = generated.ContainerClientBreakLeaseResponse

type ChangeResponse = generated.ContainerClientChangeLeaseResponse

type ReleaseResponse = generated.ContainerClientReleaseLeaseResponse

type RenewResponse = generated.ContainerClientRenewLeaseResponse
