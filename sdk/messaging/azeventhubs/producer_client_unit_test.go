// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnitNewProducerClient(t *testing.T) {
	t.Run("ConnectionStringNoEntityPath", func(t *testing.T) {
		connectionStringNoEntityPath := "Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>"

		client, err := NewProducerClientFromConnectionString(connectionStringNoEntityPath, "eventHubName", nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, "eventHubName", client.eventHub)

		client, err = NewProducerClientFromConnectionString(connectionStringNoEntityPath, "", nil)
		require.EqualError(t, err, "connection string does not contain an EntityPath. eventHub cannot be an empty string")
		require.Nil(t, client)
	})

	t.Run("ConnectionStringWithEntityPath", func(t *testing.T) {
		connectionStringWithEntityPath := "Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>;EntityPath=eventHubName"

		client, err := NewProducerClientFromConnectionString(connectionStringWithEntityPath, "", nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, "eventHubName", client.eventHub)

		client, err = NewProducerClientFromConnectionString(connectionStringWithEntityPath, "eventHubName", nil)
		require.EqualError(t, err, "connection string contains an EntityPath. eventHub must be an empty string")
		require.Nil(t, client)
	})

	t.Run("TokenCredential", func(t *testing.T) {
		tokenCredential := fakeTokenCredential{}
		client, err := NewProducerClient("ripark.servicebus.windows.net", "eventHubName", tokenCredential, nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, "eventHubName", client.eventHub)
	})
}
