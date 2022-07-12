//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

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
