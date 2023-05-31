//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Datalake Storage service allowing you to manipulate paths.
type Client base.Client[generated.PathClient]

func (p *Client) generated() *generated.PathClient {
	return base.InnerClient((*base.Client[generated.PathClient])(p))
}

func (p *Client) sharedKey() *exported.SharedKeyCredential {
	return base.SharedKey((*base.Client[generated.PathClient])(p))
}

// URL returns the URL endpoint used by the Client object.
func (p *Client) URL() string {
	return "s.generated().Endpoint()"
}

// SetAccessControl sets the owner, owning group, and permissions for a file or directory (dfs1).
func (p *Client) SetAccessControl(ctx context.Context, options *SetAccessControlOptions) (SetAccessControlResponse, error) {
	return SetAccessControlResponse{}, nil
}

// SetAccessControlRecursive sets the owner, owning group, and permissions for a file or directory (dfs1).
func (p *Client) SetAccessControlRecursive(ctx context.Context, options *SetAccessControlRecursiveOptions) (SetAccessControlRecursiveResponse, error) {
	// TODO explicitly pass SetAccessControlRecursiveMode
	return SetAccessControlRecursiveResponse{}, nil
}

// UpdateAccessControlRecursive updates the owner, owning group, and permissions for a file or directory (dfs1).
func (p *Client) UpdateAccessControlRecursive(ctx context.Context, options *UpdateAccessControlRecursiveOptions) (UpdateAccessControlRecursiveResponse, error) {
	// TODO explicitly pass SetAccessControlRecursiveMode
	return SetAccessControlRecursiveResponse{}, nil
}

// GetAccessControl gets the owner, owning group, and permissions for a file or directory (dfs1).
func (p *Client) GetAccessControl(ctx context.Context, options *GetAccessControlOptions) (GetAccessControlResponse, error) {
	return GetAccessControlResponse{}, nil
}

// RemoveAccessControlRecursive removes the owner, owning group, and permissions for a file or directory (dfs1).
func (p *Client) RemoveAccessControlRecursive(ctx context.Context, options *RemoveAccessControlRecursiveOptions) (RemoveAccessControlRecursiveResponse, error) {
	// TODO explicitly pass SetAccessControlRecursiveMode
	return SetAccessControlRecursiveResponse{}, nil
}

// SetMetadata sets the metadata for a file or directory (blob3).
func (p *Client) SetMetadata(ctx context.Context, options *SetMetadataOptions) (SetMetadataResponse, error) {
	// TODO: call directly into blob
	return SetMetadataResponse{}, nil
}

// SetHTTPHeaders sets the HTTP headers for a file or directory (blob3).
func (p *Client) SetHTTPHeaders(ctx context.Context, httpHeaders HTTPHeaders, options *SetHTTPHeadersOptions) (SetHTTPHeadersResponse, error) {
	// TODO: call formatBlobHTTPHeaders() since we want to add the blob prefix to our options before calling into blob
	// TODO: call into blob
	return SetHTTPHeadersResponse{}, nil
}
