//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azwebpubsub

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azwebpubsub/internal"
)

// ClientOptions contains optional settings for [Client]
type ClientOptions struct {
	azcore.ClientOptions
}

// NewLogsClient creates a client that accesses Azure Monitor logs data.
func NewClient(endpoint string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	authPolicy := runtime.NewBearerTokenPolicy(credential, []string{internal.TokenScope}, nil)
	azcoreClient, err := azcore.NewClient(internal.ModuleName+".Client", internal.ModuleVersion,
		runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &Client{
		internal: azcoreClient,
		endpoint: endpoint,
	}, nil
}

// NewClientFromConnectionString creates a Client from a connection string.
//
//	Endpoint=https://<your-namespace>.webpubsub.azure.com/;AccessKey=<key>;
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	props, err := internal.ParseConnectionString(connectionString)

	if err != nil {
		return nil, err
	}

	azcoreClient, err := azcore.NewClient(internal.ModuleName+".Client", internal.ModuleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{internal.NewWebPubSubKeyCredentialPolicy(props.AccessKey)},
	}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return &Client{
		internal: azcoreClient,
		endpoint: props.Endpoint,
	}, nil
}
