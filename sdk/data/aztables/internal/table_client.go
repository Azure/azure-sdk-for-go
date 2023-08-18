//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func NewTableClient(endpoint string, plOpts runtime.PipelineOptions, options *azcore.ClientOptions) (*TableClient, error) {
	client, err := azcore.NewClient(moduleName+".TableClient", version, plOpts, options)
	if err != nil {
		return nil, err
	}
	return &TableClient{
		endpoint: endpoint,
		internal: client,
	}, nil
}

func (t *TableClient) Endpoint() string {
	return t.endpoint
}

func (t *TableClient) Pipeline() runtime.Pipeline {
	return t.internal.Pipeline()
}
