//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

// CreateContainerOptions contains the optional parameters for the ContainerClient.Create method.
type CreateContainerOptions = service.CreateContainerOptions

// DeleteContainerOptions contains the optional parameters for the container.Client.Delete method.
type DeleteContainerOptions = service.DeleteContainerOptions

// DeleteBlobOptions contains the optional parameters for the Client.Delete method.
type DeleteBlobOptions = blob.DeleteOptions

// DownloadOptions contains the optional parameters for the Client.Download method.
type DownloadOptions struct {
	BlobOptions *BlobDownloadOptions
}

// BlobDownloadOptions contains the optional parameters for the Client.Download method.
type BlobDownloadOptions = blob.DownloadOptions

// ListBlobsOptions contains the optional parameters for the container.Client.ListBlobFlatSegment method.
type ListBlobsOptions = container.ListBlobsFlatOptions

// ListContainersOptions contains the optional parameters for the container.Client.ListContainers operation
type ListContainersOptions = service.ListContainersOptions
