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

func (client *ServiceClient) ListQueuesSegmentCreateRequest(ctx context.Context, options *ServiceClientListQueuesSegmentOptions) (*policy.Request, error) {
	return client.listQueuesSegmentCreateRequest(ctx, options)
}

func (client *ServiceClient) ListQueuesSegmentHandleResponse(resp *http.Response) (ServiceClientListQueuesSegmentResponse, error) {
	return client.listQueuesSegmentHandleResponse(resp)
}

func NewServiceClient(url string, azClient *azcore.Client) *ServiceClient {
	client := &ServiceClient{
		internal: azClient,
		url:      url,
	}
	return client
}
