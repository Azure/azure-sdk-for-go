// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func (client *QueueClient) Endpoint() string {
	return client.url
}

func (client *QueueClient) InternalClient() *azcore.Client {
	return client.internal
}

type MessagesClient struct {
	// client just here for back compatibility. Use QueueClient instead.
}

func NewQueueClient(url string, azClient *azcore.Client) *QueueClient {
	client := &QueueClient{
		internal: azClient,
		url:      url,
	}
	return client
}
