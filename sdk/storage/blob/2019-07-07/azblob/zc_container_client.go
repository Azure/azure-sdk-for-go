// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ContainerOperations contains the methods for the Container group.
type ContainerOperations interface {
	// ListBlobs - [Update] The List Blobs operation returns a list of the blobs under the specified container
	ListBlobs(containerListBlobFlatSegmentOptions *ContainerListBlobFlatSegmentOptions) (ListBlobsFlatSegmentResponsePager, error)
}

// NewContainerClient creates an instance of the client type with the specified endpoint.
func NewContainerClient(endpoint string, cred azcore.Credential) ContainerOperations {
	options := defaultClientOptions()
	client := newClient(endpoint, cred, &options)
	return &myContainerClient{
		containerClient: &containerClient{
			client: client,
		},
	}
}

// NewContainerClientWithPipeline creates an instance of the client type with the specified endpoint and pipeline.
func NewContainerClientWithPipeline(endpoint string, p azcore.Pipeline) ContainerOperations {
	client := newClientWithPipeline(endpoint, p)
	return &myContainerClient{
		containerClient: &containerClient{
			client: client,
		},
	}
}

type myContainerClient struct {
	*containerClient
}

// ListBlobFlatSegment - [Update] The List Blobs operation returns a list of the blobs under the specified container
func (client *myContainerClient) ListBlobs(containerListBlobFlatSegmentOptions *ContainerListBlobFlatSegmentOptions) (ListBlobsFlatSegmentResponsePager, error) {
	return &myListBlobsFlatSegmentResponsePager{
		client:  client,
		options: containerListBlobFlatSegmentOptions,
	}, nil
}

func (client *myContainerClient) ListBlobFlatSegmentHandleError(resp *azcore.Response) error {
	panic("nyi")
}

type myListBlobsFlatSegmentResponsePager struct {
	client *myContainerClient
	// request parameters
	options *ContainerListBlobFlatSegmentOptions
	// contains the current response
	current *ListBlobsFlatSegmentResponseResponse
	// any error encountered
	err error
}

func (p *myListBlobsFlatSegmentResponsePager) NextPage(ctx context.Context) bool {
	var nextMarker *string
	if p.current != nil {
		if p.current.EnumerationResults.NextMarker == nil || len(*p.current.EnumerationResults.NextMarker) == 0 {
			return false
		}
		nextMarker = p.current.EnumerationResults.NextMarker
	}
	options := *p.options
	options.Marker = nextMarker
	req, err := p.client.ListBlobFlatSegmentCreateRequest(ctx, &options)
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.client.Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !resp.HasStatusCode(http.StatusOK) {
		p.err = p.client.ListBlobFlatSegmentHandleError(resp)
		return false
	}
	result, err := p.client.ListBlobFlatSegmentHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

func (p *myListBlobsFlatSegmentResponsePager) Err() error {
	return p.err
}

func (p *myListBlobsFlatSegmentResponsePager) PageResponse() *ListBlobsFlatSegmentResponseResponse {
	return p.current
}
