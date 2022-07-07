//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
)

// CreateContainerOptions contains the optional parameters for the ContainerClient.Create method.
type CreateContainerOptions = service.CreateContainerOptions

type DeleteContainerOptions = service.DeleteContainerOptions

type DeleteBlobOptions = blob.DeleteOptions

//type UploadOptions = blockblob.ConcurrentUploadOptions

type DownloadOptions struct {
	BlobOptions *BlobDownloadOptions
}

type BlobDownloadOptions = blob.DownloadOptions

type ListBlobsOptions = container.ListBlobsFlatOptions

type ListContainersOptions = service.ListContainersOptions
