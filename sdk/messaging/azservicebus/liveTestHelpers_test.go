// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/stretchr/testify/require"
)

func enableDebugClientOptions(t *testing.T, baseClientOptions *ClientOptions) (*ClientOptions, func()) {
	// Setting this variable will cause the SB client to dump out (in TESTS ONLY)
	// the pre-master-key for your AMQP connection. This allows you decrypt a packet
	// capture from wireshark.
	//
	// If you want to do this just set SSLKEYLOGFILE_TEST env var to a path on disk and
	// Go will write out the key.
	keyLogFile := os.Getenv("SSLKEYLOGFILE_TEST")

	var clientOptions ClientOptions

	if baseClientOptions != nil {
		clientOptions = *baseClientOptions
	}

	if keyLogFile != "" {
		writer, err := os.Create(keyLogFile)

		if err != nil {
			require.Fail(t, fmt.Sprintf("SSLKEYLOGFILE_TEST was set but we failed to create a keylog file at %s: %s", keyLogFile, err))
		}

		clientOptions.TLSConfig = &tls.Config{
			KeyLogWriter: writer,
		}

		return &clientOptions, func() { _ = writer.Close() }
	}

	return &clientOptions, func() {}
}

type liveTestOptions struct {
	QueueProperties *admin.QueueProperties
	ClientOptions   *ClientOptions
}

func setupLiveTest(t *testing.T, options *liveTestOptions) (*Client, func(), string) {
	if options == nil {
		options = &liveTestOptions{}
	}

	clientOptions, flushKeyFn := enableDebugClientOptions(t, options.ClientOptions)
	serviceBusClient := newServiceBusClientForTest(t, &test.NewClientOptions[ClientOptions]{
		ClientOptions: clientOptions,
	})

	queueName, cleanupQueue := createQueue(t, nil, options.QueueProperties)

	testCleanup := func() {
		require.NoError(t, serviceBusClient.Close(context.Background()))
		flushKeyFn()
		cleanupQueue()

		// just a simple sanity check that closing twice doesn't cause errors.
		// it's basically zero cost since all the links and connection are gone from the
		// first Close().
		require.NoError(t, serviceBusClient.Close(context.Background()))
	}

	return serviceBusClient, testCleanup, queueName
}

type liveTestOptionsWithSubscription struct {
	SubscriptionProperties *admin.SubscriptionProperties
	TopicProperties        *admin.TopicProperties
	ClientOptions          *ClientOptions
}

func setupLiveTestWithSubscription(t *testing.T, options *liveTestOptionsWithSubscription) (client *Client, cleanup func(), topic string, subscription string) {
	if options == nil {
		options = &liveTestOptionsWithSubscription{}
	}

	clientOptions, flushKeyFn := enableDebugClientOptions(t, options.ClientOptions)

	serviceBusClient := newServiceBusClientForTest(t, &test.NewClientOptions[ClientOptions]{
		ClientOptions: clientOptions,
	})

	topic, cleanupTopic := createSubscription(t, options.TopicProperties, options.SubscriptionProperties)

	testCleanup := func() {
		require.NoError(t, serviceBusClient.Close(context.Background()))
		flushKeyFn()
		cleanupTopic()

		// just a simple sanity check that closing twice doesn't cause errors.
		// it's basically zero cost since all the links and connection are gone from the
		// first Close().
		require.NoError(t, serviceBusClient.Close(context.Background()))
	}

	return serviceBusClient, testCleanup, topic, "sub"
}

// createQueue creates a queue, automatically setting it to delete on idle in 5 minutes.
func createQueue(t *testing.T, options *test.NewClientOptions[admin.ClientOptions], queueProperties *admin.QueueProperties) (string, func()) {
	adminClient := newAdminClientForTest(t, options)
	return createQueueUsingClient(t, adminClient, queueProperties)
}

func createQueueUsingClient(t *testing.T, adminClient *admin.Client, queueProperties *admin.QueueProperties) (string, func()) {
	nanoSeconds := time.Now().UnixNano()
	queueName := fmt.Sprintf("queue-%X", nanoSeconds)

	if queueProperties == nil {
		queueProperties = &admin.QueueProperties{}
	}

	if queueProperties.AutoDeleteOnIdle == nil {
		autoDeleteOnIdle := "PT5M"
		queueProperties.AutoDeleteOnIdle = &autoDeleteOnIdle
	}

	_, err := adminClient.CreateQueue(context.Background(), queueName, &admin.CreateQueueOptions{
		Properties: queueProperties,
	})
	require.NoError(t, err)

	return queueName, func() {
		deleteQueue(t, adminClient, queueName)
	}
}

// createSubscription creates a topic, automatically setting it to delete on idle in 5 minutes.
// It also creates a subscription named 'sub'.
func createSubscription(t *testing.T, topicProperties *admin.TopicProperties, subscriptionProperties *admin.SubscriptionProperties) (string, func()) {
	nanoSeconds := time.Now().UnixNano()
	topicName := fmt.Sprintf("topic-%X", nanoSeconds)

	adminClient := newAdminClientForTest(t, nil)

	if topicProperties == nil {
		topicProperties = &admin.TopicProperties{}
	}

	autoDeleteOnIdle := "PT5M"
	topicProperties.AutoDeleteOnIdle = &autoDeleteOnIdle

	_, err := adminClient.CreateTopic(context.Background(), topicName, &admin.CreateTopicOptions{
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
	// TODO

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
	}, RetryOptions{}, nil)

	require.NoError(t, err)

	return msg
}

func requireScheduledMessageDisappears(ctx context.Context, t *testing.T, receiver *Receiver, sequenceNumber int64) {
	// this function will keep checking a particular sequence number until it's gone (ie, it was the last
	// sequence number so it's obvious) _or_ we end up retrieving the next message instead since
	// it auto-skips gaps.

	for {
		msgs, err := receiver.PeekMessages(ctx, 1, &PeekMessagesOptions{
			FromSequenceNumber: &sequenceNumber,
		})
		require.NoError(t, err)

		if len(msgs) == 0 {
			// no message exists at the sequence number, and there was nowhere to jump to
			return
		}

		if *msgs[0].SequenceNumber != sequenceNumber {
			// the message is gone, we've been pushed to the next message after the "gap"
			return
		}

		require.Equal(t, MessageStateScheduled, msgs[0].State)
		time.Sleep(100 * time.Millisecond)
	}
}

func newServiceBusClientForTest(t *testing.T, options *test.NewClientOptions[ClientOptions]) *Client {
	return test.NewClient(t, test.NewClientArgs[ClientOptions, Client]{
		NewClientFromConnectionString: NewClientFromConnectionString, // allowed connection string
		NewClient:                     NewClient,
	}, options)
}

func newAdminClientForTest(t *testing.T, options *test.NewClientOptions[admin.ClientOptions]) *admin.Client {
	return test.NewClient(t, test.NewClientArgs[admin.ClientOptions, admin.Client]{
		NewClientFromConnectionString: admin.NewClientFromConnectionString, // allowed connection string
		NewClient:                     admin.NewClient,
	}, options)
}
