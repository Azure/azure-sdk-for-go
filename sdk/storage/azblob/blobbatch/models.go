//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blobbatch

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// DeleteOptions contains the optional parameters for the BlobClient.Delete method.
type DeleteOptions struct {
	VersionID         *string
	Snapshot          *string
	BlobDeleteOptions *blob.DeleteOptions
}

// SetTierOptions contains the optional parameters for the BlobClient.SetTier method.
type SetTierOptions struct {
	VersionID          *string
	Snapshot           *string
	BlobSetTierOptions *blob.SetTierOptions
}

// BatchDeleteOptions contains the options for batch delete operation
//   - BlobPath: Must be in the following format.
//   - blobName when using ContainerBatchClient, e.g. blob.txt
//   - /containerName/blobName when using ServiceBatchClient, e.g. /container/blob.txt
//   - DeleteOptions - optional parameters for the BlobClient.Delete method.
type BatchDeleteOptions struct {
	BlobPath *string
	*DeleteOptions
}

// BatchSetTierOptions contains the options for batch set tier operation
//   - BlobPath: Must be in the following format.
//   - blobName when using ContainerBatchClient, e.g. blob.txt
//   - /containerName/blobName when using ServiceBatchClient, e.g. /container/blob.txt
//   - AccessTier - defines values for Blob Access Tier
//   - DeleteOptions - optional parameters for the BlobClient.Delete method.
type BatchSetTierOptions struct {
	BlobPath   *string
	AccessTier blob.AccessTier
	*SetTierOptions
}

// BatchBuilder is used for creating the batch operations list. It contains delete and set tier lists for blobs.
type BatchBuilder struct {
	batchDeleteList  []*BatchDeleteOptions
	batchSetTierList []*BatchSetTierOptions
}

// ContainerClientSubmitBatchResponse contains the response from method ContainerBatchClient.SubmitBatch.
type ContainerClientSubmitBatchResponse = generated.ContainerClientSubmitBatchResponse

// ServiceClientSubmitBatchResponse contains the response from method ServiceBatchClient.SubmitBatch.
type ServiceClientSubmitBatchResponse = generated.ServiceClientSubmitBatchResponse
