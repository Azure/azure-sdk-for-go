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
}

type Client[T any] struct {
	inner      *T
	credential any
}

func InnerClient[T any](client *Client[T]) *T {
	return client.inner
}

func SharedKey[T any](client *Client[T]) *exported.SharedKeyCredential {
	switch cred := client.credential.(type) {
	case *exported.SharedKeyCredential:
		return cred
	default:
		return nil
	}
}

func Credential[T any](client *Client[T]) any {
	return client.credential
}

func NewServiceClient(serviceURL string, pipeline runtime.Pipeline, credential any) *Client[generated.ServiceClient] {
	return &Client[generated.ServiceClient]{
		inner:      generated.NewServiceClient(serviceURL, pipeline),
		credential: credential,
	}
}

func NewShareClient(shareURL string, pipeline runtime.Pipeline, credential any) *Client[generated.ShareClient] {
	return &Client[generated.ShareClient]{
		inner:      generated.NewShareClient(shareURL, pipeline),
		credential: credential,
	}
}

func NewDirectoryClient(directoryURL string, pipeline runtime.Pipeline, credential any) *Client[generated.DirectoryClient] {
	return &Client[generated.DirectoryClient]{
		inner:      generated.NewDirectoryClient(directoryURL, pipeline),
		credential: credential,
	}
}

func NewFileClient(fileURL string, pipeline runtime.Pipeline, credential any) *Client[generated.FileClient] {
	return &Client[generated.FileClient]{
		inner:      generated.NewFileClient(fileURL, pipeline),
		credential: credential,
	}
}
