// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/sas"
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
		t.Logf("Env cred works, being added to our chained token credential")
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

	err = sender.SendMessage(context.TODO(), &Message{Body: []byte("hello - authenticating with a TokenCredential")}, nil)

	if err != nil && strings.Contains(err.Error(), "'Send' claim(s) are required to perform this operation") {
		const sleepDuration = time.Minute
		// it's possible we're just dealing with a propagation delay for our
		// configured identity and the newly created resource. We'll sleep
		// a bit to give it some time and try again.
		t.Logf("Enacting CI workaround to deal with RBAC propagation delays. Sleeping for %s...", sleepDuration)
		time.Sleep(sleepDuration)
		t.Logf("Done sleeping for %s", sleepDuration)

		err = sender.SendMessage(context.TODO(), &Message{Body: []byte("hello - authenticating with a TokenCredential")}, nil)
	}

	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queue, nil)
	require.NoError(t, err)
	actualSettler, _ := receiver.settler.(*messageSettler)
	actualSettler.onlyDoBackupSettlement = true // this'll also exercise the management link

	messages, err := receiver.ReceiveMessages(context.TODO(), 1, nil)
	require.NoError(t, err)

	require.EqualValues(t, []string{"hello - authenticating with a TokenCredential"}, getSortedBodies(messages))

	for _, m := range messages {
		err = receiver.CompleteMessage(context.TODO(), m, nil)
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
	}, nil)
	require.NoError(t, err)

	// we have to test this down here since the connection is lazy initialized.
	require.True(t, webSocketCreateCalled)

	receiver, err := client.NewReceiverForQueue(queue, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	require.EqualValues(t, "hello world", string(messages[0].Body))
}

func TestNewClientUsingSharedAccessSignature(t *testing.T) {
	sasCS, err := sas.CreateConnectionStringWithSAS(test.GetConnectionString(t), time.Hour)
	require.NoError(t, err)

	// sanity check - we did actually generate a connection string with an embedded SharedAccessSignature
	require.Contains(t, sasCS, "SharedAccessSignature=SharedAccessSignature")

	queue, cleanup := createQueue(t, sasCS, nil)
	defer cleanup()

	client, err := NewClientFromConnectionString(sasCS, nil)
	require.NoError(t, err)

	sender, err := client.NewSender(queue, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello world"),
	}, nil)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(queue, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	})
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	require.EqualValues(t, "hello world", string(messages[0].Body))
}

const fastNotFoundDuration = 10 * time.Second

func TestNewClientNewSenderNotFound(t *testing.T) {
	connectionString := test.GetConnectionString(t)
	client, err := NewClientFromConnectionString(connectionString, nil)
	require.NoError(t, err)

	defer client.Close(context.Background())

	sender, err := client.NewSender("non-existent-queue", nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), fastNotFoundDuration)
	defer cancel()

	err = sender.SendMessage(ctx, &Message{Body: []byte("hello")}, nil)
	assertRPCNotFound(t, err)
}

func TestNewClientNewReceiverNotFound(t *testing.T) {
	connectionString := test.GetConnectionString(t)
	client, err := NewClientFromConnectionString(connectionString, nil)
	require.NoError(t, err)

	defer client.Close(context.Background())

	receiver, err := client.NewReceiverForQueue("non-existent-queue", nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), fastNotFoundDuration)
	defer cancel()

	messages, err := receiver.ReceiveMessages(ctx, 1, nil)
	require.Nil(t, messages)
	assertRPCNotFound(t, err)

	receiver, err = client.NewReceiverForSubscription("non-existent-topic", "non-existent-subscription", nil)
	require.NoError(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), fastNotFoundDuration)
	defer cancel()

	messages, err = receiver.PeekMessages(ctx, 1, nil)
	require.Nil(t, messages)
	assertRPCNotFound(t, err)
}

func TestClientNewSessionReceiverNotFound(t *testing.T) {
	connectionString := test.GetConnectionString(t)
	client, err := NewClientFromConnectionString(connectionString, nil)
	require.NoError(t, err)

	defer client.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), fastNotFoundDuration)
	defer cancel()

	receiver, err := client.AcceptSessionForQueue(ctx, "non-existent-queue", "session-id", nil)
	require.Nil(t, receiver)
	assertRPCNotFound(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), fastNotFoundDuration)
	defer cancel()

	receiver, err = client.AcceptNextSessionForQueue(ctx, "non-existent-queue", nil)
	require.Nil(t, receiver)
	assertRPCNotFound(t, err)
}

func TestClientCloseVsClosePermanently(t *testing.T) {
	connectionString := test.GetConnectionString(t)
	client, err := NewClientFromConnectionString(connectionString, nil)
	require.NoError(t, err)

	require.NoError(t, client.Close(context.Background()))

	receiver, err := client.NewReceiverForQueue("queue", nil)
	require.EqualError(t, err, "client has been closed by user")
	require.Nil(t, receiver)

	receiver, err = client.NewReceiverForSubscription("topic", "subscription", nil)
	require.EqualError(t, err, "client has been closed by user")
	require.Nil(t, receiver)

	sender, err := client.NewSender("queue", nil)
	require.EqualError(t, err, "client has been closed by user")
	require.Nil(t, sender)

	sessionReceiver, err := client.AcceptSessionForQueue(context.Background(), "queue", "session-id-that-is-not-used", nil)
	require.EqualError(t, err, "client has been closed by user")
	require.Nil(t, sessionReceiver)

	sessionReceiver, err = client.AcceptSessionForSubscription(context.Background(), "topic", "subscription", "session-id-that-is-not-used", nil)
	require.EqualError(t, err, "client has been closed by user")
	require.Nil(t, sessionReceiver)

	sessionReceiver, err = client.AcceptNextSessionForSubscription(context.Background(), "topic", "subscription", nil)
	require.EqualError(t, err, "client has been closed by user")
	require.Nil(t, sessionReceiver)
}

func TestClientNewSessionReceiverCancel(t *testing.T) {
	// Both the session APIs create the receiver immediately however AcceptNextSession() has a quirk
	// where it takes an excessively long time.
	connectionString := test.GetConnectionString(t)

	queue, cleanup := createQueue(t, connectionString, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})

	defer cleanup()

	client, err := NewClientFromConnectionString(connectionString, nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// non-cancelled version
	receiver, err := client.AcceptNextSessionForQueue(ctx, queue, nil)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Nil(t, receiver)
}

func TestClientPropagatesRetryOptionsForSessions(t *testing.T) {
	connectionString := test.GetConnectionString(t)

	queue, cleanupQueue := createQueue(t, connectionString, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})

	defer cleanupQueue()

	topic, cleanupTopic := createSubscription(t, connectionString, nil, &admin.SubscriptionProperties{
		RequiresSession: to.Ptr(true),
	})

	defer cleanupTopic()

	expectedRetryOptions := RetryOptions{
		MaxRetries:    1,
		RetryDelay:    time.Second,
		MaxRetryDelay: time.Millisecond,
	}

	client, err := NewClientFromConnectionString(connectionString, &ClientOptions{
		RetryOptions: expectedRetryOptions,
	})
	require.NoError(t, err)

	actualNS := client.namespace.(*internal.Namespace)
	require.Equal(t, expectedRetryOptions, actualNS.RetryOptions)

	queueSender, err := client.NewSender(queue, nil)
	require.NoError(t, err)

	topicSender, err := client.NewSender(topic, nil)
	require.NoError(t, err)

	err = queueSender.SendMessage(context.Background(), &Message{
		SessionID: to.Ptr("hello"),
	}, nil)
	require.NoError(t, err)

	err = topicSender.SendMessage(context.Background(), &Message{
		SessionID: to.Ptr("hello"),
	}, nil)
	require.NoError(t, err)

	sessionReceiver, err := client.AcceptSessionForQueue(context.Background(), queue, "hello", nil)
	require.NoError(t, err)
	require.NoError(t, sessionReceiver.Close(context.Background()))

	require.Equal(t, expectedRetryOptions, sessionReceiver.inner.retryOptions)

	sessionReceiver, err = client.AcceptSessionForSubscription(context.Background(), topic, "sub", "hello", nil)
	require.NoError(t, err)
	require.NoError(t, sessionReceiver.Close(context.Background()))

	require.Equal(t, expectedRetryOptions, sessionReceiver.inner.retryOptions)

	sessionReceiver, err = client.AcceptNextSessionForQueue(context.Background(), queue, nil)
	require.NoError(t, err)
	require.NoError(t, sessionReceiver.Close(context.Background()))

	require.Equal(t, expectedRetryOptions, sessionReceiver.inner.retryOptions)

	sessionReceiver, err = client.AcceptNextSessionForSubscription(context.Background(), topic, "sub", nil)
	require.NoError(t, err)
	require.NoError(t, sessionReceiver.Close(context.Background()))

	require.Equal(t, expectedRetryOptions, sessionReceiver.inner.retryOptions)
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

	t.Run("RetryOptionsArePropagated", func(t *testing.T) {
		// retry options are passed and copied along several routes, just make sure it's properly propagated.
		// NOTE: session receivers are checked in a separate test because they require actual SB access.
		client, err := NewClient("fake.something", struct{ azcore.TokenCredential }{}, &ClientOptions{
			RetryOptions: RetryOptions{
				MaxRetries:    101,
				RetryDelay:    6 * time.Hour,
				MaxRetryDelay: 12 * time.Hour,
			},
		})

		client.namespace = &internal.FakeNS{
			AMQPLinks: &internal.FakeAMQPLinks{
				Receiver: &internal.FakeAMQPReceiver{},
			},
		}

		require.NoError(t, err)

		require.Equal(t, RetryOptions{
			MaxRetries:    101,
			RetryDelay:    6 * time.Hour,
			MaxRetryDelay: 12 * time.Hour,
		}, client.retryOptions)

		sender, err := client.NewSender("hello", nil)
		require.NoError(t, err)

		require.Equal(t, RetryOptions{
			MaxRetries:    101,
			RetryDelay:    6 * time.Hour,
			MaxRetryDelay: 12 * time.Hour,
		}, sender.retryOptions)

		receiver, err := client.NewReceiverForQueue("hello", nil)
		require.NoError(t, err)

		require.Equal(t, RetryOptions{
			MaxRetries:    101,
			RetryDelay:    6 * time.Hour,
			MaxRetryDelay: 12 * time.Hour,
		}, receiver.retryOptions)

		actualSettler := receiver.settler.(*messageSettler)

		require.Equal(t, RetryOptions{
			MaxRetries:    101,
			RetryDelay:    6 * time.Hour,
			MaxRetryDelay: 12 * time.Hour,
		}, actualSettler.retryOptions)

		subscriptionReceiver, err := client.NewReceiverForSubscription("hello", "world", nil)
		require.NoError(t, err)

		require.Equal(t, RetryOptions{
			MaxRetries:    101,
			RetryDelay:    6 * time.Hour,
			MaxRetryDelay: 12 * time.Hour,
		}, subscriptionReceiver.retryOptions)
	})
}

func assertRPCNotFound(t *testing.T, err error) {
	require.NotNil(t, err)

	var rpcError interface {
		RPCCode() int
		error
	}

	require.ErrorAs(t, err, &rpcError)
	require.Equal(t, http.StatusNotFound, rpcError.RPCCode())
}
