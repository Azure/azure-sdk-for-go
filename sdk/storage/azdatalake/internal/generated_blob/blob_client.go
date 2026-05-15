// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated_blob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func (client *BlobClient) Endpoint() string {
	return client.url
}

func (client *BlobClient) InternalClient() *azcore.Client {
	return client.internal
}

// BlobClient contains the methods for the Blob group.
// Don't use this type directly, use a constructor function instead.
type BlobClient struct {
	internal *azcore.Client
	url      string
}

// NewBlobClient creates a new instance of BlobClient with the specified values.
//   - endpoint - The URL of the service account, container, or blob that is the target of the desired operation.
//   - azClient - azcore.Client is a basic HTTP client. It consists of a pipeline and tracing provider.
func NewBlobClient(endpoint string, azClient *azcore.Client) *BlobClient {
	client := &BlobClient{
		internal: azClient,
		url:      endpoint,
	}
	return client
}
