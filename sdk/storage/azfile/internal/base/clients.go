// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"strings"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions

	// AllowTrailingDot specifies if a trailing dot present in request url should be trimmed or not.
	AllowTrailingDot *bool

	// FileRequestIntent is required when using TokenCredential for authentication.
	// Acceptable value is backup.
	FileRequestIntent *generated.ShareTokenIntent

	// AllowSourceTrailingDot specifies if a trailing dot present in source url should be trimmed or not.
	AllowSourceTrailingDot *bool

	// Audience to use when requesting tokens for Azure Active Directory authentication.
	// Only has an effect when credential is of type TokenCredential. The value could be
	// https://storage.azure.com/ (default) or https://<account>.file.core.windows.net.
	Audience string

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

func GetAudience(clOpts *ClientOptions) string {
	if clOpts == nil || len(strings.TrimSpace(clOpts.Audience)) == 0 {
		return shared.TokenScope
	} else {
		return strings.TrimRight(clOpts.Audience, "/") + "/.default"
	}
}

func NewServiceClient(serviceURL string, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *Client[generated.ServiceClient] {
	return &Client[generated.ServiceClient]{
		inner:     generated.NewServiceClient(serviceURL, azClient),
		sharedKey: sharedKey,
		options:   options,
	}
}

func NewShareClient(shareURL string, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *Client[generated.ShareClient] {
	return &Client[generated.ShareClient]{
		inner:     generated.NewShareClient(shareURL, options.FileRequestIntent, azClient),
		sharedKey: sharedKey,
		options:   options,
	}
}

func NewDirectoryClient(directoryURL string, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *Client[generated.DirectoryClient] {
	return &Client[generated.DirectoryClient]{
		inner:     generated.NewDirectoryClient(directoryURL, options.AllowTrailingDot, options.FileRequestIntent, options.AllowSourceTrailingDot, azClient),
		sharedKey: sharedKey,
		options:   options,
	}
}

func NewFileClient(fileURL string, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *Client[generated.FileClient] {
	return &Client[generated.FileClient]{
		inner:     generated.NewFileClient(fileURL, options.AllowTrailingDot, options.FileRequestIntent, options.AllowSourceTrailingDot, azClient),
		sharedKey: sharedKey,
		options:   options,
	}
}
