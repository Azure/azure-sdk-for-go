// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package base

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
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
		inner:     generated.NewShareClient(shareURL, azClient),
		sharedKey: sharedKey,
		options:   options,
	}
}

func NewDirectoryClient(directoryURL string, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *Client[generated.DirectoryClient] {
	return &Client[generated.DirectoryClient]{
		inner:     generated.NewDirectoryClient(directoryURL, azClient),
		sharedKey: sharedKey,
		options:   options,
	}
}

func NewFileClient(fileURL string, azClient *azcore.Client, sharedKey *exported.SharedKeyCredential, options *ClientOptions) *Client[generated.FileClient] {
	return &Client[generated.FileClient]{
		inner:     generated.NewFileClient(fileURL, azClient),
		sharedKey: sharedKey,
		options:   options,
	}
}

// Addressing describes how an operation requires a file or directory client to address its resource.
type Addressing int

const (
	// PathOnly operations require the client to address the resource by its path.
	PathOnly Addressing = iota

	// PathOrFileID operations support a client that addresses the resource either by its path or by its file ID.
	PathOrFileID

	// FileIDOnly operations require the client to address the resource by its file ID.
	FileIDOnly
)

// FileID returns the file ID in the given file or directory URL, or an empty string if the URL addresses the
// resource by its path.
func FileID(resourceURL string) string {
	urlParts, err := sas.ParseURL(resourceURL)
	if err != nil {
		return ""
	}
	return urlParts.FileID
}

// CheckAddressing returns an error if the operation cannot be called on a client with the given URL, so that the
// operation fails before a request is sent to the service.
func CheckAddressing(resourceURL string, operation string, addressing Addressing) error {
	urlParts, err := sas.ParseURL(resourceURL)
	if err != nil {
		// An unparsable URL is reported by the request itself, as it was before file ID support was added.
		return nil
	}

	if urlParts.FileID == "" {
		if addressing == FileIDOnly {
			return fmt.Errorf("%s: %w", operation, fileerror.ErrFileIDRequired)
		}
		return nil
	}

	if addressing == PathOnly {
		return fmt.Errorf("%s: %w", operation, fileerror.ErrFileIDNotSupported)
	}
	if urlParts.DirectoryOrFilePath != "" {
		// A file ID identifies the resource within the share on its own, e.g. a client navigated to with
		// directory.Client.NewFileClient from a client that addresses a directory by its file ID.
		return fmt.Errorf("%s: the URL must not contain both a path and a file ID: %w", operation, fileerror.ErrFileIDNotSupported)
	}
	return nil
}
