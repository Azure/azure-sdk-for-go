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

// DeleteResponse contains the response from method Client.Delete.
type DeleteResponse = generated.ContainerClientDeleteResponse

// GetPropertiesResponse contains the response from method Client.GetProperties.
type GetPropertiesResponse = generated.ContainerClientGetPropertiesResponse

// ListBlobsFlatResponse contains the response from method Client.ListBlobFlatSegment.
type ListBlobsFlatResponse = generated.ContainerClientListBlobFlatSegmentResponse

// ListBlobsHierarchyResponse contains the response from method Client.ListBlobHierarchySegment.
type ListBlobsHierarchyResponse = generated.ContainerClientListBlobHierarchySegmentResponse

// SetMetadataResponse contains the response from method Client.SetMetadata.
type SetMetadataResponse = generated.ContainerClientSetMetadataResponse

// GetAccessPolicyResponse contains the response from method Client.GetAccessPolicy.
type GetAccessPolicyResponse = generated.ContainerClientGetAccessPolicyResponse

// SetAccessPolicyResponse contains the response from method Client.SetAccessPolicy.
type SetAccessPolicyResponse = generated.ContainerClientSetAccessPolicyResponse

// ---------------------------------------------------------------------------------------------------------------------

// AcquireResponse contains the response from method LeaseClient.AcquireLease.
type AcquireResponse = generated.ContainerClientAcquireLeaseResponse

// BreakResponse contains the response from method LeaseClient.BreakLease.
type BreakResponse = generated.ContainerClientBreakLeaseResponse

// ChangeResponse contains the response from method LeaseClient.ChangeLease.
type ChangeResponse = generated.ContainerClientChangeLeaseResponse

// ReleaseResponse contains the response from method LeaseClient.ReleaseLease.
type ReleaseResponse = generated.ContainerClientReleaseLeaseResponse

// RenewResponse contains the response from method LeaseClient.RenewLease.
type RenewResponse = generated.ContainerClientRenewLeaseResponse
