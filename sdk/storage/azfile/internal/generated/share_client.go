//go:build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func (client *ShareClient) Endpoint() string {
	return client.endpoint
}

func (client *ShareClient) InternalClient() *azcore.Client {
	return client.internal
}

// NewShareClient creates a new instance of ShareClient with the specified values.
//   - endpoint - The URL of the service account, share, directory or file that is the target of the desired operation.
//   - fileRequestIntent - Valid value is backup
//   - azClient - azcore.Client is a basic HTTP client.  It consists of a pipeline and tracing provider.
func NewShareClient(endpoint string, fileRequestIntent *ShareTokenIntent, azClient *azcore.Client) *ShareClient {
	client := &ShareClient{
		internal:          azClient,
		endpoint:          endpoint,
		fileRequestIntent: fileRequestIntent,
	}
	return client
}
