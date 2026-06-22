// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

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
//   - endpoint - The URL of the service account, share, directory or file that is the target of the desired operation.
//   - azClient - azcore.Client is a basic HTTP client.  It consists of a pipeline and tracing provider.
func NewServiceClient(endpoint string, azClient *azcore.Client) *ServiceClient {
	client := &ServiceClient{
		internal: azClient,
		url:      endpoint,
	}
	return client
}

// ListSharesSegmentCreateRequest creates the ListSharesSegment request.
func (client *ServiceClient) ListSharesSegmentCreateRequest(ctx context.Context, options *ServiceClientListSharesSegmentOptions) (*policy.Request, error) {
	return client.listSharesSegmentCreateRequest(ctx, options)
}

// ListSharesSegmentHandleResponse handles the ListSharesSegment response.
func (client *ServiceClient) ListSharesSegmentHandleResponse(resp *http.Response) (ServiceClientListSharesSegmentResponse, error) {
	return client.listSharesSegmentHandleResponse(resp)
}
