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

// HealthAPIClientOptions contains optional settings for [HealthAPIClient]
type HealthAPIClientOptions struct {
	azcore.ClientOptions
}

// NewHealthAPIClient creates a client that checks the healthy status of Web PubSub service
func NewHealthAPIClient(endpoint string, options *HealthAPIClientOptions) (*HealthAPIClient, error) {
	if options == nil {
		options = &HealthAPIClientOptions{}
	}

	azcoreClient, err := azcore.NewClient(internal.ModuleName, internal.ModuleVersion, runtime.PipelineOptions{}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	return &HealthAPIClient{
		internal: azcoreClient,
		endpoint: endpoint,
	}, nil
}
