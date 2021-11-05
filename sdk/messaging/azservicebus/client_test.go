// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/stretchr/testify/require"
)

func TestNewClientWithAzureIdentity(t *testing.T) {
	queue, cleanup := createQueue(t, getConnectionString(t), nil)
	defer cleanup()

	// test with azure identity support
	ns := os.Getenv("SERVICEBUS_ENDPOINT")
	envCred, err := azidentity.NewEnvironmentCredential(nil)

	if err != nil || ns == "" {
		t.Skip("Azure Identity compatible credentials not configured")
	}

	client, err := NewClient(ns, envCred, nil)
	require.NoError(t, err)

	sender, err := client.NewSender(queue)
	require.NoError(t, err)

	err = sender.SendMessage(context.TODO(), &Message{Body: []byte("hello - authenticating with a TokenCredential")})
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queue, nil)
	require.NoError(t, err)
	actualSettler, _ := receiver.settler.(*messageSettler)
	actualSettler.onlyDoBackupSettlement = true // this'll also exercise the management link

	messages, err := receiver.ReceiveMessages(context.TODO(), 1, nil)
	require.NoError(t, err)

	require.EqualValues(t, []string{"hello - authenticating with a TokenCredential"}, getSortedBodies(messages))

	for _, m := range messages {
		err = receiver.CompleteMessage(context.TODO(), m)
		require.NoError(t, err)
	}

	client.Close(context.TODO())
}

func TestNewClientUnitTests(t *testing.T) {
	t.Run("WithTokenCredential", func(t *testing.T) {
		fakeTokenCredential := struct{ azcore.TokenCredential }{}

		client, err := NewClient("fake.something", fakeTokenCredential, nil)
		require.NoError(t, err)

		require.NoError(t, err)
		require.EqualValues(t, fakeTokenCredential, client.config.credential)
		require.EqualValues(t, "fake.something", client.config.fullyQualifiedNamespace)

		client, err = NewClient("mysb.windows.servicebus.net", fakeTokenCredential, nil)
		require.NoError(t, err)
		require.EqualValues(t, fakeTokenCredential, client.config.credential)
		require.EqualValues(t, "mysb.windows.servicebus.net", client.config.fullyQualifiedNamespace)

		_, err = NewClientFromConnectionString("", nil)
		require.EqualError(t, err, "connectionString must not be empty")

		_, err = NewClient("", fakeTokenCredential, nil)
		require.EqualError(t, err, "fullyQualifiedNamespace must not be empty")

		_, err = NewClient("mysb", fakeTokenCredential, nil)
		require.EqualError(t, err, "fullyQualifiedNamespace is not properly formed. Should be similar to 'myservicebus.servicebus.windows.net'")

		_, err = NewClient("fake.something", nil, nil)
		require.EqualError(t, err, "credential was nil")

		// (really all part of the same functionality)
		ns := &internal.Namespace{}
		require.NoError(t, internal.NamespacesWithTokenCredential("mysb.windows.servicebus.net",
			fakeTokenCredential)(ns))

		require.EqualValues(t, ns.Name, "mysb")
		require.EqualValues(t, ns.Suffix, "windows.servicebus.net")

		err = internal.NamespacesWithTokenCredential("mysb", fakeTokenCredential)(&internal.Namespace{})
		require.EqualError(t, err, "fullyQualifiedNamespace is not properly formed. Should be similar to 'myservicebus.servicebus.windows.net'")
	})

	t.Run("CloseAndLinkTracking", func(t *testing.T) {
		setupClient := func() (*Client, *internal.FakeNS) {
			client, err := NewClient("fake.something", struct{ azcore.TokenCredential }{}, nil)
			require.NoError(t, err)

			ns := &internal.FakeNS{
				AMQPLinks: &internal.FakeAMQPLinks{
					Receiver: &internal.FakeAMQPReceiver{},
				},
			}

			client.namespace = ns
			return client, ns
		}

		client, ns := setupClient()
		_, err := client.NewSender("hello")

		require.NoError(t, err)
		require.EqualValues(t, 1, len(client.links))
		require.NotNil(t, client.links[1])
		require.NoError(t, client.Close(context.Background()))
		require.Empty(t, client.links)
		require.True(t, ns.AMQPLinks.ClosedPermanently())

		client, ns = setupClient()
		_, err = client.NewReceiverForQueue("hello", nil)

		require.NoError(t, err)
		require.EqualValues(t, 1, len(client.links))
		require.NotNil(t, client.links[1])
		require.NoError(t, client.Close(context.Background()))
		require.Empty(t, client.links)
		require.True(t, ns.AMQPLinks.ClosedPermanently())

		client, ns = setupClient()
		_, err = client.NewReceiverForSubscription("hello", "world", nil)

		require.NoError(t, err)
		require.EqualValues(t, 1, len(client.links))
		require.NotNil(t, client.links[1])
		require.NoError(t, client.Close(context.Background()))
		require.Empty(t, client.links)
		require.EqualValues(t, 1, ns.AMQPLinks.Closed)

		client, ns = setupClient()
		_, err = newProcessorForQueue(client, "hello", nil)

		require.NoError(t, err)
		require.EqualValues(t, 1, len(client.links))
		require.NotNil(t, client.links[1])
		require.NoError(t, client.Close(context.Background()))
		require.Empty(t, client.links)
		require.EqualValues(t, 1, ns.AMQPLinks.Closed)

		client, ns = setupClient()
		_, err = newProcessorForSubscription(client, "hello", "world", nil)

		require.NoError(t, err)
		require.EqualValues(t, 1, len(client.links))
		require.NotNil(t, client.links[1])
		require.NoError(t, client.Close(context.Background()))
		require.Empty(t, client.links)
		require.EqualValues(t, 1, ns.AMQPLinks.Closed)
	})
}
