// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	// ISO8601 is used for formatting file creation, last write and change time.
	ISO8601 = "2006-01-02T15:04:05.0000000Z07:00"
)

func (client *DirectoryClient) Endpoint() string {
	return client.url
}

func (client *DirectoryClient) InternalClient() *azcore.Client {
	return client.internal
}

// NewDirectoryClient creates a new instance of DirectoryClient with the specified values.
//   - endpoint - The URL of the service account, share, or directory that is the target of the desired operation.
//   - azClient - azcore.Client is a basic HTTP client.  It consists of a pipeline and tracing provider.
func NewDirectoryClient(endpoint string, azClient *azcore.Client) *DirectoryClient {
	client := &DirectoryClient{
		internal: azClient,
		url:      endpoint,
	}
	return client
}

// ListFilesAndDirectoriesSegmentCreateRequest creates the ListFilesAndDirectoriesSegment request.
func (client *DirectoryClient) ListFilesAndDirectoriesSegmentCreateRequest(ctx context.Context, options *DirectoryClientListFilesAndDirectoriesSegmentOptions) (*policy.Request, error) {
	return client.listFilesAndDirectoriesSegmentCreateRequest(ctx, options)
}

// ListFilesAndDirectoriesSegmentHandleResponse handles the ListFilesAndDirectoriesSegment response.
func (client *DirectoryClient) ListFilesAndDirectoriesSegmentHandleResponse(resp *http.Response) (DirectoryClientListFilesAndDirectoriesSegmentResponse, error) {
	return client.listFilesAndDirectoriesSegmentHandleResponse(resp)
}
