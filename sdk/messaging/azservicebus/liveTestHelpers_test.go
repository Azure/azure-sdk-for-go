// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
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

		// just a simple sanity check that closing twice doesn't cause errors.
		// it's basically zero cost since all the links and connection are gone from the
		// first Close().
		require.NoError(t, serviceBusClient.Close(context.Background()))
	}

	return serviceBusClient, testCleanup, queueName
}

// createQueue creates a queue, automatically setting it to delete on idle in 5 minutes.
func createQueue(t *testing.T, connectionString string, queueProperties *admin.QueueProperties) (string, func()) {
	nanoSeconds := time.Now().UnixNano()
	queueName := fmt.Sprintf("queue-%X", nanoSeconds)

	adminClient, err := admin.NewClientFromConnectionString(connectionString, nil)
	require.NoError(t, err)

	if queueProperties == nil {
		queueProperties = &admin.QueueProperties{}
	}

	autoDeleteOnIdle := "PT5M"
	queueProperties.AutoDeleteOnIdle = &autoDeleteOnIdle

	_, err = adminClient.CreateQueue(context.Background(), queueName, &admin.CreateQueueOptions{
		Properties: queueProperties,
	})
	require.NoError(t, err)

	return queueName, func() {
		deleteQueue(t, adminClient, queueName)
	}
}

// createSubscription creates a queue, automatically setting it to delete on idle in 5 minutes.
func createSubscription(t *testing.T, connectionString string, topicProperties *admin.TopicProperties, subscriptionProperties *admin.SubscriptionProperties) (string, func()) {
	nanoSeconds := time.Now().UnixNano()
	topicName := fmt.Sprintf("topic-%X", nanoSeconds)

	adminClient, err := admin.NewClientFromConnectionString(connectionString, nil)
	require.NoError(t, err)

	if topicProperties == nil {
		topicProperties = &admin.TopicProperties{}
	}

	autoDeleteOnIdle := "PT5M"
	topicProperties.AutoDeleteOnIdle = &autoDeleteOnIdle

	_, err = adminClient.CreateTopic(context.Background(), topicName, &admin.CreateTopicOptions{
		Properties: topicProperties,
	})
	require.NoError(t, err)

	_, err = adminClient.CreateSubscription(context.Background(), topicName, "sub", &admin.CreateSubscriptionOptions{Properties: subscriptionProperties})
	require.NoError(t, err)

	return topicName, func() {
		_, err := adminClient.DeleteTopic(context.Background(), topicName, nil)
		require.NoError(t, err)
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

// peekSingleMessageForTest wraps a standard Receiver.Peek() call so it returns at least one message
// and fails tests otherwise.
func peekSingleMessageForTest(t *testing.T, receiver *Receiver) *ReceivedMessage {
	var msg *ReceivedMessage

	// Peek, unlike Receive, doesn't block until at least one message has arrived, so we have to poll
	// to get a similar effect.
	err := utils.Retry(context.Background(), EventReceiver, "peekSingleForTest", func(ctx context.Context, args *utils.RetryFnArgs) error {
		peekedMessages, err := receiver.PeekMessages(context.Background(), 1, nil)
		require.NoError(t, err)

		if len(peekedMessages) == 1 {
			msg = peekedMessages[0]
			return nil
		} else {
			return errors.New("No peekable messages available")
		}
	}, func(err error) bool {
		return false
	}, RetryOptions{})

	require.NoError(t, err)

	return msg
}
