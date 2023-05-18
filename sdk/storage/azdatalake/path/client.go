//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions base.ClientOptions

// Client represents a URL to the Azure Blob Storage service allowing you to manipulate blob containers.
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

// dfs
func (p *Client) SetAccessControl(options *SetAccessControlOptions) (SetAccessControlResponse, error) {
	return SetAccessControlResponse{}, nil
}

// dfs
func (p *Client) SetAccessControlRecursive(options *SetAccessControlRecursiveOptions) (SetAccessControlRecursiveResponse, error) {
	return SetAccessControlRecursiveResponse{}, nil
}

// dfs
func (p *Client) UpdateAccessControlRecursive(options *UpdateAccessControlRecursiveOptions) (UpdateAccessControlRecursiveResponse, error) {
	return SetAccessControlRecursiveResponse{}, nil
}

// dfs
func (p *Client) RemoveAccessControlRecursive(options *RemoveAccessControlRecursiveOptions) (RemoveAccessControlRecursiveResponse, error) {
	return SetAccessControlRecursiveResponse{}, nil
}

// blob
func (p *Client) SetMetadata(options *SetMetadataOptions) (SetMetadataResponse, error) {
	return SetMetadataResponse{}, nil
}

// blob
func (p *Client) SetHTTPHeaders(options *SetHTTPHeadersOptions) (SetHTTPHeadersResponse, error) {
	return SetHTTPHeadersResponse{}, nil
}
