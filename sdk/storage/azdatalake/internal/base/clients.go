//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions
	pipelineOptions *runtime.PipelineOptions
}

type Client[T any] struct {
	inner     *T
	sharedKey *exported.SharedKeyCredential
	options   *ClientOptions
}

func InnerClient[T any](client *Client[T]) *T {
	return client.inner
}

func SharedKey[T any](client *Client[T]) *exported.SharedKeyCredential {
	return client.sharedKey
}

func GetClientOptions[T any](client *Client[T]) *ClientOptions {
	return client.options
}

func GetPipelineOptions(clOpts *ClientOptions) *runtime.PipelineOptions {
	return clOpts.pipelineOptions
}

func SetPipelineOptions(clOpts *ClientOptions, plOpts *runtime.PipelineOptions) {
	clOpts.pipelineOptions = plOpts
}

func NewServiceClient(serviceURL string, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *Client[generated.ServiceClient] {
	return &Client[generated.ServiceClient]{
		inner:     generated.NewServiceClient(serviceURL, azClient),
		sharedKey: sharedKey,
		options:   options,
	}
}

func NewFilesystemClient(containerURL string, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *Client[generated.FileSystemClient] {
	return &Client[generated.FileSystemClient]{
		inner:     generated.NewFilesystemClient(containerURL, azClient),
		sharedKey: sharedKey,
		options:   options,
	}
}

func NewPathClient(containerURL string, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *Client[generated.PathClient] {
	return &Client[generated.PathClient]{
		inner:     generated.NewPathClient(containerURL, azClient),
		sharedKey: sharedKey,
		options:   options,
	}
}
