//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azwebpubsub

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub/internal"
)

// NewHealthClientWithNoCredential creates a client that checks the healthy status of Web PubSub service
func NewHealthClientWithNoCredential(endpoint string, options *ClientOptions) (*HealthClient, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	azcoreClient, err := azcore.NewClient(internal.ModuleName, internal.ModuleVersion, runtime.PipelineOptions{}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	return &HealthClient{
		internal: azcoreClient,
		endpoint: endpoint,
	}, nil
}
