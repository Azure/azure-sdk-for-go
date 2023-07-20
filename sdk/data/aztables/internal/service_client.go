//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func NewServiceClient(endpoint string, plOpts runtime.PipelineOptions, options *azcore.ClientOptions) (*ServiceClient, error) {
	client, err := azcore.NewClient(moduleName+".ServiceClient", version, plOpts, options)
	if err != nil {
		return nil, err
	}
	return &ServiceClient{
		endpoint: endpoint,
		internal: client,
		version:  Enum0TwoThousandNineteen0202,
	}, nil
}

func (s *ServiceClient) Endpoint() string {
	return s.endpoint
}

func (s *ServiceClient) Pipeline() runtime.Pipeline {
	return s.internal.Pipeline()
}
