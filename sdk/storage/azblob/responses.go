//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

type CreateContainerResponse = service.CreateContainerResponse

type DeleteContainerResponse = service.DeleteContainerResponse

type DeleteBlobResponse = blob.DeleteResponse

type UploadResponse = blockblob.CommitBlockListResponse

type DownloadResponse = blob.DownloadResponse

type ListBlobsResponse = container.ListBlobsFlatResponse

type ListContainersResponse = service.ListContainersResponse
