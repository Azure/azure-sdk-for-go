// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setupLiveTest(t *testing.T, props *QueueProperties) (*Client, func(), string) {
	cs := getConnectionString(t)

	serviceBusClient, err := NewClientFromConnectionString(cs, nil)
	require.NoError(t, err)

	queueName, cleanupQueue := createQueue(t, cs, props)

	testCleanup := func() {
		require.NoError(t, serviceBusClient.Close(context.Background()))
		cleanupQueue()
		require.NoError(t, serviceBusClient.Close(context.Background()))
	}

	return serviceBusClient, testCleanup, queueName
}

func getConnectionString(t *testing.T) string {
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	if cs == "" {
		t.Skip()
	}

	return cs
}

func getConnectionStringWithoutManagePerms(t *testing.T) string {
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING_NO_MANAGE")

	if cs == "" {
		t.Skip()
	}

	return cs
}

// createQueue creates a queue using a subset of entries in 'queueDescription':
// - EnablePartitioning
// - RequiresSession
func createQueue(t *testing.T, connectionString string, queueProperties *QueueProperties) (string, func()) {
	nanoSeconds := time.Now().UnixNano()
	queueName := fmt.Sprintf("queue-%X", nanoSeconds)

	adminClient, err := NewAdminClientWithConnectionString(connectionString, nil)
	require.NoError(t, err)

	if queueProperties == nil {
		queueProperties = &QueueProperties{}
	}

	queueProperties.Name = queueName
	_, err = adminClient.AddQueueWithProperties(context.Background(), queueProperties)
	require.NoError(t, err)

	return queueName, func() {
		if _, err := adminClient.DeleteQueue(context.TODO(), queueProperties.Name); err != nil {
			require.NoError(t, err)
		}
	}
}
