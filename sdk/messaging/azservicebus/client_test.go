// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/stretchr/testify/require"
	"nhooyr.io/websocket"
)

func TestNewClientWithAzureIdentity(t *testing.T) {
	queue, cleanup := createQueue(t, test.GetConnectionString(t), nil)
	defer cleanup()

	// test with azure identity support
	ns := os.Getenv("SERVICEBUS_ENDPOINT")

	var credsToAdd []azcore.TokenCredential

	cliCred, err := azidentity.NewAzureCLICredential(nil)
	require.NoError(t, err)

	envCred, err := azidentity.NewEnvironmentCredential(nil)

	if err == nil {
		fmt.Printf("Env cred works, being added to our chained token credential")
		credsToAdd = append(credsToAdd, envCred)
	}

	credsToAdd = append(credsToAdd, cliCred)

	cred, err := azidentity.NewChainedTokenCredential(credsToAdd, nil)
	require.NoError(t, err)

	if err != nil || ns == "" {
		t.Skip("Azure Identity compatible credentials not configured")
	}

	client, err := NewClient(ns, cred, nil)
	require.NoError(t, err)

	sender, err := client.NewSender(queue, nil)
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

func TestNewClientWithWebsockets(t *testing.T) {
	connectionString := test.GetConnectionString(t)

	queue, cleanup := createQueue(t, connectionString, nil)
	defer cleanup()

	webSocketCreateCalled := false

	client, err := NewClientFromConnectionString(connectionString, &ClientOptions{
		NewWebSocketConn: func(ctx context.Context, args NewWebSocketConnArgs) (net.Conn, error) {
			webSocketCreateCalled = true
			opts := &websocket.DialOptions{Subprotocols: []string{"amqp"}}
			wssConn, _, err := websocket.Dial(ctx, args.Host, opts)

			if err != nil {
				return nil, err
			}

			return websocket.NetConn(context.Background(), wssConn, websocket.MessageBinary), nil
		},
	})
	require.NoError(t, err)

	sender, err := client.NewSender(queue, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello world"),
	})
	require.NoError(t, err)

	// we have to test this down here since the connection is lazy initialized.
	require.True(t, webSocketCreateCalled)

	receiver, err := client.NewReceiverForQueue(queue, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	bytes, err := messages[0].Body()
	require.NoError(t, err)
	require.EqualValues(t, "hello world", string(bytes))
}

func TestNewClientUnitTests(t *testing.T) {
	t.Run("WithTokenCredential", func(t *testing.T) {
		fakeTokenCredential := struct{ azcore.TokenCredential }{}

		client, err := NewClient("fake.something", fakeTokenCredential, nil)
		require.NoError(t, err)

		require.NoError(t, err)
		require.EqualValues(t, fakeTokenCredential, client.creds.credential)
		require.EqualValues(t, "fake.something", client.creds.fullyQualifiedNamespace)

		client, err = NewClient("mysb.windows.servicebus.net", fakeTokenCredential, nil)
		require.NoError(t, err)
		require.EqualValues(t, fakeTokenCredential, client.creds.credential)
		require.EqualValues(t, "mysb.windows.servicebus.net", client.creds.fullyQualifiedNamespace)

		_, err = NewClientFromConnectionString("", nil)
		require.EqualError(t, err, "connectionString must not be empty")

		_, err = NewClient("", fakeTokenCredential, nil)
		require.EqualError(t, err, "fullyQualifiedNamespace must not be empty")

		_, err = NewClient("fake.something", nil, nil)
		require.EqualError(t, err, "credential was nil")

		// (really all part of the same functionality)
		ns := &internal.Namespace{}
		require.NoError(t, internal.NamespaceWithTokenCredential("mysb.windows.servicebus.net",
			fakeTokenCredential)(ns))

		require.EqualValues(t, ns.FQDN, "mysb.windows.servicebus.net")
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
		_, err := client.NewSender("hello", nil)

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
	})
}
