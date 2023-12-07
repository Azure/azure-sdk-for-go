//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

func NewAzureAppConfigurationClient(endpoint string, client *azcore.Client) *AzureAppConfigurationClient {
	return &AzureAppConfigurationClient{
		internal: client,
		endpoint: endpoint,
	}
}

func (a *AzureAppConfigurationClient) Tracer() tracing.Tracer {
	return a.internal.Tracer()
}

func NewCreateSnapshotPoller[T any](ctx context.Context, client *AzureAppConfigurationClient, name string, entity Snapshot, options *AzureAppConfigurationClientBeginCreateSnapshotOptions) (*runtime.Poller[T], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createSnapshot(ctx, name, entity, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[T](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[T](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}
