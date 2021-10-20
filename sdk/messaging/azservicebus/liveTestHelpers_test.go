// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/stretchr/testify/require"
)

func setupLiveTest(t *testing.T, qd *internal.QueueDescription) (*Client, func(), string) {
	cs := getConnectionString(t)

	serviceBusClient, err := NewClientFromConnectionString(cs, nil)
	require.NoError(t, err)

	queueName, cleanupQueue := createQueue(t, cs, qd)

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

// createQueue creates a queue using a subset of entries in 'queueDescription':
// - EnablePartitioning
// - RequiresSession
func createQueue(t *testing.T, connectionString string, queueDescription *internal.QueueDescription) (string, func()) {
	nanoSeconds := time.Now().UnixNano()
	queueName := fmt.Sprintf("queue-%X", nanoSeconds)

	qm, err := internal.NewQueueManagerWithConnectionString(connectionString)
	require.NoError(t, err)

	var opts []internal.QueueManagementOption

	if queueDescription != nil {
		if queueDescription.EnablePartitioning != nil && *queueDescription.EnablePartitioning {
			opts = append(opts, internal.QueueEntityWithPartitioning())
		}

		if queueDescription.RequiresSession != nil && *queueDescription.RequiresSession {
			opts = append(opts, internal.QueueEntityWithRequiredSessions())
		}

		if queueDescription.RequiresSession != nil && *queueDescription.RequiresSession {
			opts = append(opts, internal.QueueEntityWithRequiredSessions())
		}
	}

	_, err = qm.Put(context.TODO(), queueName, opts...)
	require.NoError(t, err)

	return queueName, func() {
		if err := qm.Delete(context.TODO(), queueName); err != nil {
			require.NoError(t, err)
		}
	}
}
