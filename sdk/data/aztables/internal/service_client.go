// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func NewServiceClient(endpoint string, client *azcore.Client) *ServiceClient {
	return &ServiceClient{
		endpoint: endpoint,
		internal: client,
	}
}

func (s *ServiceClient) Endpoint() string {
	return s.endpoint
}

func (s *ServiceClient) Pipeline() runtime.Pipeline {
	return s.internal.Pipeline()
}
