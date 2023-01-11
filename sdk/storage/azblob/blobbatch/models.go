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

type DeleteOptions struct {
	VersionID         *string
	Snapshot          *string
	BlobDeleteOptions *blob.DeleteOptions
}

type SetTierOptions struct {
	VersionID          *string
	Snapshot           *string
	BlobSetTierOptions *blob.SetTierOptions
}

type BatchDeleteOptions struct {
	BlobPath *string
	*DeleteOptions
}

type BatchSetTierOptions struct {
	BlobPath   *string
	AccessTier blob.AccessTier
	*SetTierOptions
}

type BatchBuilder struct {
	batchDeleteList  []*BatchDeleteOptions
	batchSetTierList []*BatchSetTierOptions
}

type ContainerClientSubmitBatchResponse = generated.ContainerClientSubmitBatchResponse

type ServiceClientSubmitBatchResponse = generated.ServiceClientSubmitBatchResponse
