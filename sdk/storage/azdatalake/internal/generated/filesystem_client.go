// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func (client *FileSystemClient) Endpoint() string {
	return client.url
}

func (client *FileSystemClient) InternalClient() *azcore.Client {
	return client.internal
}

// NewFileSystemClient creates a new instance of ServiceClient with the specified values.
//   - endpoint - The URL of the service account, share, directory or file that is the target of the desired operation.
//   - azClient - azcore.Client is a basic HTTP client.  It consists of a pipeline and tracing provider.
func NewFileSystemClient(endpoint string, azClient *azcore.Client) *FileSystemClient {
	client := &FileSystemClient{
		internal: azClient,
		url:      endpoint,
	}
	return client
}

func (client *FileSystemClient) ListBlobHierarchySegmentCreateRequest(ctx context.Context, options *FileSystemClientListBlobHierarchySegmentOptions) (*policy.Request, error) {
	return client.listBlobHierarchySegmentCreateRequest(ctx, options)
}

// ListBlobHierarchySegmentHandleResponse handles the ListBlobHierarchySegment response.
func (client *FileSystemClient) ListBlobHierarchySegmentHandleResponse(resp *http.Response) (FileSystemClientListPathHierarchySegmentResponse, error) {
	return client.listBlobHierarchySegmentHandleResponse(resp)
}

func (client *FileSystemClient) ListPathsCreateRequest(ctx context.Context, recursive bool, options *FileSystemClientListPathsOptions) (*policy.Request, error) {
	return client.listPathsCreateRequest(ctx, recursive, options)
}

// ListPathsHandleResponse handles the ListPaths response.
func (client *FileSystemClient) ListPathsHandleResponse(resp *http.Response) (FileSystemClientListPathsResponse, error) {
	return client.listPathsHandleResponse(resp)
}
