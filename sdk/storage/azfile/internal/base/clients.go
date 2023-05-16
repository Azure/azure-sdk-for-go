//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions
	AllowTrailingDot       *bool
	FileRequestIntent      *generated.ShareTokenIntent
	AllowSourceTrailingDot *bool
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

func NewServiceClient(serviceURL string, pipeline runtime.Pipeline, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *Client[generated.ServiceClient] {
	return &Client[generated.ServiceClient]{
		inner:     generated.NewServiceClient(serviceURL, pipeline),
		sharedKey: sharedKey,
		options:   options,
	}
}

func NewShareClient(shareURL string, pipeline runtime.Pipeline, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *Client[generated.ShareClient] {
	return &Client[generated.ShareClient]{
		inner:     generated.NewShareClient(shareURL, options.FileRequestIntent, pipeline),
		sharedKey: sharedKey,
		options:   options,
	}
}

func NewDirectoryClient(directoryURL string, pipeline runtime.Pipeline, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *Client[generated.DirectoryClient] {
	return &Client[generated.DirectoryClient]{
		inner:     generated.NewDirectoryClient(directoryURL, options.AllowTrailingDot, options.FileRequestIntent, options.AllowSourceTrailingDot, pipeline),
		sharedKey: sharedKey,
		options:   options,
	}
}

func NewFileClient(fileURL string, pipeline runtime.Pipeline, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *Client[generated.FileClient] {
	return &Client[generated.FileClient]{
		inner:     generated.NewFileClient(fileURL, options.AllowTrailingDot, options.FileRequestIntent, options.AllowSourceTrailingDot, pipeline),
		sharedKey: sharedKey,
		options:   options,
	}
}
