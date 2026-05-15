// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated_blob

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

func (client *ServiceClient) Endpoint() string {
	return client.url
}

func (client *ServiceClient) InternalClient() *azcore.Client {
	return client.internal
}

// NewServiceClient creates a new instance of ServiceClient with the specified values.
//   - endpoint - The URL of the service account, container, or blob that is the target of the desired operation.
//   - azClient - azcore.Client is a basic HTTP client. It consists of a pipeline and tracing provider.
func NewServiceClient(endpoint string, azClient *azcore.Client) *ServiceClient {
	client := &ServiceClient{
		internal: azClient,
		url:      endpoint,
	}
	return client
}

func (client *ServiceClient) ListContainersSegmentCreateRequest(ctx context.Context, options *ServiceClientListContainersSegmentOptions) (*policy.Request, error) {
	return client.listContainersSegmentCreateRequest(ctx, options)
}

// listContainersSegmentHandleResponse handles the ListContainersSegment response.
func (client *ServiceClient) ListContainersSegmentHandleResponse(resp *http.Response) (ServiceClientListFileSystemsSegmentResponse, error) {
	return client.listContainersSegmentHandleResponse(resp)
}
