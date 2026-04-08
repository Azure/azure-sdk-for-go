// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// NewMessagesClient creates a new instance of MessagesClient with the specified values.
//   - endpoint - The URL of the service account, queue, or message that is the target of the desired operation.
//   - azClient - azcore.Client is a basic HTTP client. It consists of a pipeline and tracing provider.
func NewMessagesClient(endpoint string, azClient *azcore.Client) *MessagesClient {
	client := &MessagesClient{
		internal: azClient,
		endpoint: endpoint,
		version:  ServiceVersion,
	}
	return client
}
