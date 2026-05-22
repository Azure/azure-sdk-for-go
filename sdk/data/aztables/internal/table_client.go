// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
)

func NewTableClient(endpoint string, client *azcore.Client) *TableClient {
	return &TableClient{
		endpoint: endpoint,
		internal: client,
	}
}

func (t *TableClient) Endpoint() string {
	return t.endpoint
}

func (t *TableClient) Pipeline() runtime.Pipeline {
	return t.internal.Pipeline()
}

func (t *TableClient) Tracer() tracing.Tracer {
	return t.internal.Tracer()
}
