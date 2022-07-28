// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func TestUnitNewConsumerClient(t *testing.T) {
	t.Run("ConnectionStringNoEntityPath", func(t *testing.T) {
		connectionStringNoEntityPath := "Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>"

		client, err := NewConsumerClientFromConnectionString(connectionStringNoEntityPath, "eventHubName", "0", "$Default", nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, "eventHubName", client.eventHub)

		client, err = NewConsumerClientFromConnectionString(connectionStringNoEntityPath, "", "0", "$Default", nil)
		require.EqualError(t, err, "connection string does not contain an EntityPath. eventHub cannot be an empty string")
		require.Nil(t, client)
	})

	t.Run("ConnectionStringWithEntityPath", func(t *testing.T) {
		connectionStringWithEntityPath := "Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>;EntityPath=eventHubName"

		client, err := NewConsumerClientFromConnectionString(connectionStringWithEntityPath, "", "0", "$Default", nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, "eventHubName", client.eventHub)

		client, err = NewConsumerClientFromConnectionString(connectionStringWithEntityPath, "eventHubName", "0", "$Default", nil)
		require.EqualError(t, err, "connection string contains an EntityPath. eventHub must be an empty string")
		require.Nil(t, client)
	})

	t.Run("TokenCredential", func(t *testing.T) {
		tokenCredential := fakeTokenCredential{}
		client, err := NewConsumerClient("ripark.servicebus.windows.net", "eventHubName", "0", "$Default", tokenCredential, nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, "eventHubName", client.eventHub)
	})
}

type fakeTokenCredential struct {
	azcore.TokenCredential
}
