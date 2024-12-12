// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"io"
	"net"
	"net/http"
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
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/stretchr/testify/require"
	"nhooyr.io/websocket"
)

func TestNewClientWithAzureIdentity(t *testing.T) {
	queue, cleanup := createQueue(t, nil, nil)
	defer cleanup()

	client := newServiceBusClientForTest(t, nil)

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

	messages, err := receiver.ReceiveMessages(context.TODO(), 1, nil)
	require.NoError(t, err)

	require.EqualValues(t, []string{"hello - authenticating with a TokenCredential"}, getSortedBodies(messages))
	forceManagementSettlement(messages)

	for _, m := range messages {
		err = receiver.CompleteMessage(context.TODO(), m, nil)
		require.NoError(t, err)
	}

	client.Close(context.TODO())
}

func TestNewClientWithWebsockets(t *testing.T) {
	queue, cleanup := createQueue(t, nil, nil)
	defer cleanup()

	client := newServiceBusClientForTest(t, &test.NewClientOptions[ClientOptions]{
		ClientOptions: &ClientOptions{
			NewWebSocketConn: func(ctx context.Context, args NewWebSocketConnArgs) (net.Conn, error) {
				transport := http.DefaultTransport.(*http.Transport).Clone()
				transport.TLSClientConfig.KeyLogWriter = io.Discard

				opts := &websocket.DialOptions{Subprotocols: []string{"amqp"}, HTTPClient: &http.Client{
					Transport: transport,
				}}

				wssConn, _, err := websocket.Dial(ctx, args.Host, opts)

				if err != nil {
					return nil, err
				}

				conn := websocket.NetConn(context.Background(), wssConn, websocket.MessageBinary)
				return conn, nil
			},
		},
	})

	defer test.RequireClose(t, client)

	sender, err := client.NewSender(queue, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("hello world"),
	}, nil)
	require.NoError(t, err)
}

func TestNewClientUsingSharedAccessSignature(t *testing.T) {
	getLogsFn := test.CaptureLogsForTest(false)

	cs := test.MustGetEnvVar(t, test.EnvKeyConnectionString)
	sasCS, err := sas.CreateConnectionStringWithSASUsingExpiry(cs, time.Now().UTC().Add(time.Hour))
	require.NoError(t, err)

	// sanity check - we did actually generate a connection string with an embedded SharedAccessSignature
	require.Contains(t, sasCS, "SharedAccessSignature=SharedAccessSignature")

	adminClient, err := admin.NewClientFromConnectionString(sasCS, nil) // allowed connection string
	require.NoError(t, err)

	queue, cleanup := createQueueUsingClient(t, adminClient, nil)
	defer cleanup()

	client, err := NewClientFromConnectionString(sasCS, nil) // allowed connection string
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

	logs := getLogsFn()
	require.Contains(t, logs, backgroundRenewalDisabledMsg)
}

const backgroundRenewalDisabledMsg = "[azsb.Auth] Token does not have an expiration date, no background renewal needed."

const fastNotFoundDuration = 10 * time.Second

func TestNewClientNewSenderNotFound(t *testing.T) {
	client := newServiceBusClientForTest(t, nil)

	defer client.Close(context.Background())

	sender, err := client.NewSender("non-existent-queue", nil)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), fastNotFoundDuration)
	defer cancel()

	err = sender.SendMessage(ctx, &Message{Body: []byte("hello")}, nil)
	assertRPCNotFound(t, err)
}

func TestNewClientNewReceiverNotFound(t *testing.T) {
	client := newServiceBusClientForTest(t, nil)

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
	client := newServiceBusClientForTest(t, nil)

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
	client := newServiceBusClientForTest(t, nil)

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
	queue, cleanup := createQueue(t, nil, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})

	defer cleanup()

	client := newServiceBusClientForTest(t, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// non-cancelled version
	receiver, err := client.AcceptNextSessionForQueue(ctx, queue, nil)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Nil(t, receiver)
}

func TestClientPropagatesRetryOptionsForSessions(t *testing.T) {
	queue, cleanupQueue := createQueue(t, nil, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})

	defer cleanupQueue()

	topic, cleanupTopic := createSubscription(t, nil, &admin.SubscriptionProperties{
		RequiresSession: to.Ptr(true),
	})

	defer cleanupTopic()

	expectedRetryOptions := RetryOptions{
		MaxRetries:    1,
		RetryDelay:    time.Second,
		MaxRetryDelay: time.Millisecond,
	}

	client := newServiceBusClientForTest(t, &test.NewClientOptions[ClientOptions]{
		ClientOptions: &ClientOptions{
			RetryOptions: expectedRetryOptions,
		},
	})

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

func TestClientUnauthorizedCreds(t *testing.T) {
	allPowerfulCS := test.MustGetEnvVar(t, test.EnvKeyConnectionString)
	queueName := "testqueue"

	t.Run("ListenOnly with Sender", func(t *testing.T) {
		cs := test.MustGetEnvVar(t, test.EnvKeyConnectionStringListenOnly)

		client, err := NewClientFromConnectionString(cs, nil) // allowed connection string
		require.NoError(t, err)

		defer test.RequireClose(t, client)

		sender, err := client.NewSender(queueName, nil)
		require.NoError(t, err)

		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte("hello world"),
		}, nil)

		var sbErr *Error
		require.ErrorAs(t, err, &sbErr)
		require.Equal(t, CodeUnauthorizedAccess, sbErr.Code)
		require.Contains(t, err.Error(), "Description: Unauthorized access. 'Send' claim(s) are required to perform this operation")
	})

	t.Run("SenderOnly with Receiver", func(t *testing.T) {
		cs := test.MustGetEnvVar(t, test.EnvKeyConnectionStringSendOnly)

		client, err := NewClientFromConnectionString(cs, nil) // allowed connection string
		require.NoError(t, err)

		defer test.RequireClose(t, client)

		receiver, err := client.NewReceiverForQueue(queueName, nil)
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
		require.Empty(t, messages)

		var sbErr *Error
		require.ErrorAs(t, err, &sbErr)
		require.Equal(t, CodeUnauthorizedAccess, sbErr.Code)
		require.Contains(t, err.Error(), "Description: Unauthorized access. 'Listen' claim(s) are required to perform this operation")
	})

	t.Run("Expired SAS", func(t *testing.T) {
		expiredCS, err := sas.CreateConnectionStringWithSASUsingExpiry(allPowerfulCS, time.Now().Add(-10*time.Minute))
		require.NoError(t, err)

		client, err := NewClientFromConnectionString(expiredCS, nil) // allowed connection string
		require.NoError(t, err)

		defer test.RequireClose(t, client)

		sender, err := client.NewSender(queueName, nil)
		require.NoError(t, err)

		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte("hello world"),
		}, nil)

		var sbErr *Error
		require.ErrorAs(t, err, &sbErr)
		require.Equal(t, CodeUnauthorizedAccess, sbErr.Code)
		require.Contains(t, err.Error(), "rpc: failed, status code 401 and description: ExpiredToken: The token is expired. Expiration time:")

		receiver, err := client.NewReceiverForQueue(queueName, nil)
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
		require.Empty(t, messages)

		sbErr = nil
		require.ErrorAs(t, err, &sbErr)
		require.Equal(t, CodeUnauthorizedAccess, sbErr.Code)
		require.Contains(t, err.Error(), "rpc: failed, status code 401 and description: ExpiredToken: The token is expired. Expiration time:")
	})

	t.Run("invalid identity creds", func(t *testing.T) {
		identityVars := test.GetIdentityVars(t)

		if identityVars == nil {
			return
		}

		cliCred, err := azidentity.NewClientSecretCredential("00000000-0000-0000-0000-000000000000", "00000000-0000-0000-0000-000000000000", "bogus-client-secret", nil)
		require.NoError(t, err)

		client, err := NewClient(identityVars.Endpoint, cliCred, nil)
		require.NoError(t, err)

		defer test.RequireClose(t, client)

		sender, err := client.NewSender(queueName, nil)
		require.NoError(t, err)

		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte("hello world"),
		}, nil)
		var authFailedErr *azidentity.AuthenticationFailedError
		require.ErrorAs(t, err, &authFailedErr)

		receiver, err := client.NewReceiverForQueue(queueName, nil)
		require.NoError(t, err)

		messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
		require.ErrorAs(t, err, &authFailedErr)
		require.Empty(t, messages)
	})
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

		_, err = NewClientFromConnectionString("", nil) // allowed connection string
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

	t.Run("TracerIsSetUp", func(t *testing.T) {
		hostName := "fake.servicebus.windows.net"
		// when tracing provider is not set, use a no-op tracer.
		client, err := NewClient(hostName, struct{ azcore.TokenCredential }{}, nil)
		require.NoError(t, err)
		require.Zero(t, client.tracer)
		require.False(t, client.tracer.Enabled())

		// when tracing provider is set, the tracer is set up with the provider.
		provider := tracing.NewSpanValidator(t, tracing.SpanMatcher{
			Name:   "TestSpan",
			Status: tracing.SpanStatusUnset,
			Attributes: []tracing.Attribute{
				{Key: tracing.MessagingSystem, Value: "servicebus"},
				{Key: tracing.ServerAddress, Value: hostName},
			},
		})
		client, err = NewClient(hostName, struct{ azcore.TokenCredential }{}, &ClientOptions{
			TracingProvider: provider,
		})
		require.NoError(t, err)
		require.NotZero(t, client.tracer)
		require.True(t, client.tracer.Enabled())

		// ensure attributes are set up correctly.
		_, endSpan := tracing.StartSpan(context.Background(), client.tracer, tracing.NewSpanConfig("TestSpan"))
		endSpan(nil)

		// attributes should be set up when using a connection string.
		fakeConnectionString := "Endpoint=sb://fake.servicebus.windows.net/;SharedAccessKeyName=TestName;SharedAccessKey=TestKey"
		client, err = NewClientFromConnectionString(fakeConnectionString, &ClientOptions{
			TracingProvider: provider,
		})
		require.NoError(t, err)
		require.NotZero(t, client.tracer)
		require.True(t, client.tracer.Enabled())

		// ensure attributes are set up correctly.
		_, endSpan = tracing.StartSpan(context.Background(), client.tracer, tracing.NewSpanConfig("TestSpan"))
		endSpan(nil)
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

		actualSettler := receiver.settler

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

func forceManagementSettlement(messages []*ReceivedMessage) {
	for _, m := range messages {
		m.settleOnMgmtLink = true
	}
}
