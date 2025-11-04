// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus_test

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/stretchr/testify/require"
)

func TestReceiveWithDelay(t *testing.T) {
	setupQueueWithMessages := func(t *testing.T, afterNumMessages int, delay time.Duration) (*azservicebus.Client, string) {
		iv := test.GetIdentityVars(t)
		faultInjectorEndpoint := test.StartSlowTransferFaultInjector(t, iv.Endpoint, afterNumMessages, delay)

		client := test.NewClient(t, test.NewClientArgs[azservicebus.ClientOptions, azservicebus.Client]{
			NewClientFromConnectionString: azservicebus.NewClientFromConnectionString,
			NewClient:                     azservicebus.NewClient,
		}, &test.NewClientOptions[azservicebus.ClientOptions]{
			ClientOptions: &azservicebus.ClientOptions{
				TLSConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				CustomEndpoint: faultInjectorEndpoint,
			},
		})

		t.Cleanup(func() { test.RequireClose(t, client) })

		queueName, cleanupQueue := test.CreateExpiringQueue(t, &atom.QueueDescription{
			LockDuration: to.Ptr("PT1M"),
		})

		t.Cleanup(cleanupQueue)

		sender, err := client.NewSender(queueName, nil)
		require.NoError(t, err)

		for range 10 {
			err = sender.SendMessage(context.Background(), &azservicebus.Message{Body: []byte("hello world")}, nil)
			require.NoError(t, err)
		}

		return client, queueName
	}

	t.Run("OnlySingleMessageReceived", func(t *testing.T) {
		client, queueName := setupQueueWithMessages(t, 1, 5*time.Second)

		receiver, err := client.NewReceiverForQueue(queueName, nil)
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(context.Background(), 10, nil)
		require.NoError(t, err)

		require.Equal(t, 1, len(messages))
	})

	t.Run("TwoMessagesExactly", func(t *testing.T) {
		client, queueName := setupQueueWithMessages(t, 2, 5*time.Second)

		receiver, err := client.NewReceiverForQueue(queueName, nil)
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(context.Background(), 10, nil)
		require.NoError(t, err)

		require.Equal(t, 2, len(messages))
	})

	t.Run("AllTen", func(t *testing.T) {
		client, queueName := setupQueueWithMessages(t, 10, 5*time.Second)

		receiver, err := client.NewReceiverForQueue(queueName, nil)
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(context.Background(), 10, nil)
		require.NoError(t, err)

		require.Equal(t, 10, len(messages))
	})
}
