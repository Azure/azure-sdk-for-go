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

// CreateContainerOptions contains the optional parameters for the ContainerClient.Create method.
type CreateContainerOptions = service.CreateContainerOptions

// DeleteContainerOptions contains the optional parameters for the container.Client.Delete method.
type DeleteContainerOptions = service.DeleteContainerOptions

// DeleteBlobOptions contains the optional parameters for the Client.Delete method.
type DeleteBlobOptions = blob.DeleteOptions

// DownloadToStreamOptions contains the optional parameters for the Client.DownloadToSTream method.
type DownloadToStreamOptions = blob.DownloadToStreamOptions

// ListBlobsOptions contains the optional parameters for the container.Client.ListBlobFlatSegment method.
type ListBlobsOptions = container.ListBlobsFlatOptions

// ListContainersOptions contains the optional parameters for the container.Client.ListContainers operation
type ListContainersOptions = service.ListContainersOptions

// UploadBufferOptions provides set of configurations for UploadBuffer operation
type UploadBufferOptions = blockblob.UploadBufferOptions

// UploadFileOptions provides set of configurations for UploadFile operation
type UploadFileOptions = blockblob.UploadReaderAtToBlockBlobOptions

// UploadStreamOptions provides set of configurations for UploadStream operation
type UploadStreamOptions = blockblob.UploadStreamOptions

// DownloadToWriterAtOptions identifies options used by the DownloadToBuffer and DownloadToFile functions.
type DownloadToWriterAtOptions = blob.DownloadToWriterAtOptions

// DownloadToBufferOptions identifies options used by the DownloadToBuffer and DownloadToFile functions.
type DownloadToBufferOptions = blob.DownloadToBufferOptions

// DownloadToFileOptions identifies options used by the DownloadToBuffer and DownloadToFile functions.
type DownloadToFileOptions = blob.DownloadToFileOptions
