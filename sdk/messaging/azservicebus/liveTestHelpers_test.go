// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/stretchr/testify/require"
)

func setupLiveTest(t *testing.T, props *admin.QueueProperties) (*Client, func(), string) {
	cs := test.GetConnectionString(t)

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

// createQueue creates a queue using a subset of entries in 'queueDescription':
// - EnablePartitioning
// - RequiresSession
func createQueue(t *testing.T, connectionString string, queueProperties *admin.QueueProperties) (string, func()) {
	nanoSeconds := time.Now().UnixNano()
	queueName := fmt.Sprintf("queue-%X", nanoSeconds)

	adminClient, err := admin.NewClientFromConnectionString(connectionString, nil)
	require.NoError(t, err)

	if queueProperties == nil {
		queueProperties = &admin.QueueProperties{}
	}

	autoDeleteOnIdle := 5 * time.Minute
	queueProperties.AutoDeleteOnIdle = &autoDeleteOnIdle

	_, err = adminClient.CreateQueue(context.Background(), queueName, queueProperties, nil)
	require.NoError(t, err)

	return queueName, func() {
		deleteQueue(t, adminClient, queueName)
	}
}

func deleteQueue(t *testing.T, ac *admin.Client, queueName string) {
	_, err := ac.DeleteQueue(context.Background(), queueName, nil)
	require.NoError(t, err)
}

func deleteSubscription(t *testing.T, ac *admin.Client, topicName string, subscriptionName string) {
	_, err := ac.DeleteSubscription(context.Background(), topicName, subscriptionName, nil)
	require.NoError(t, err)
}
