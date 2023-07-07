//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/service"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions
	pipelineOptions *runtime.PipelineOptions
}

func GetPipelineOptions(clOpts *ClientOptions) *runtime.PipelineOptions {
	return clOpts.pipelineOptions
}

func SetPipelineOptions(clOpts *ClientOptions, plOpts *runtime.PipelineOptions) {
	clOpts.pipelineOptions = plOpts
}

type CompositeClient[T, K, U any] struct {
	// generated client with dfs
	innerT *T
	// generated client with blob
	innerK *K
	// blob client
	innerU    *U
	sharedKey *exported.SharedKeyCredential
	options   *ClientOptions
}

func InnerClients[T, K, U any](client *CompositeClient[T, K, U]) (*T, *K, *U) {
	return client.innerT, client.innerK, client.innerU
}

func SharedKeyComposite[T, K, U any](client *CompositeClient[T, K, U]) *exported.SharedKeyCredential {
	return client.sharedKey
}

func NewFilesystemClient(fsURL string, fsURLWithBlobEndpoint string, client *container.Client, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *CompositeClient[generated.FileSystemClient, generated.FileSystemClient, container.Client] {
	return &CompositeClient[generated.FileSystemClient, generated.FileSystemClient, container.Client]{
		innerT:    generated.NewFilesystemClient(fsURL, azClient),
		innerK:    generated.NewFilesystemClient(fsURLWithBlobEndpoint, azClient),
		sharedKey: sharedKey,
		innerU:    client,
		options:   options,
	}
}

func NewServiceClient(serviceURL string, serviceURLWithBlobEndpoint string, client *service.Client, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *CompositeClient[generated.ServiceClient, generated.ServiceClient, service.Client] {
	return &CompositeClient[generated.ServiceClient, generated.ServiceClient, service.Client]{
		innerT:    generated.NewServiceClient(serviceURL, azClient),
		innerK:    generated.NewServiceClient(serviceURLWithBlobEndpoint, azClient),
		sharedKey: sharedKey,
		innerU:    client,
		options:   options,
	}
}

func NewPathClient(dirURL string, dirURLWithBlobEndpoint string, client *blob.Client, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *CompositeClient[generated.PathClient, generated.PathClient, blob.Client] {
	return &CompositeClient[generated.PathClient, generated.PathClient, blob.Client]{
		innerT:    generated.NewPathClient(dirURL, azClient),
		innerK:    generated.NewPathClient(dirURLWithBlobEndpoint, azClient),
		sharedKey: sharedKey,
		innerU:    client,
		options:   options,
	}
}

func GetCompositeClientOptions[T, K, U any](client *CompositeClient[T, K, U]) *ClientOptions {
	return client.options
}
