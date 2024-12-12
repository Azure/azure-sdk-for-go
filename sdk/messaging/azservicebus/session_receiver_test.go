// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	internal_errors "github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestSessionReceiver_acceptSession(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &liveTestOptions{
		QueueProperties: &admin.QueueProperties{
			RequiresSession: to.Ptr(true),
		}})
	defer cleanup()

	ctx := context.Background()

	// send a message to a specific session
	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendMessage(ctx, &Message{
		Body:      []byte("session-based message"),
		SessionID: to.Ptr("session-1"),
	}, nil)
	require.NoError(t, err)

	receiver, err := client.AcceptSessionForQueue(ctx, queueName, "session-1", nil)
	require.NoError(t, err)

	messages, err := receiver.inner.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)

	require.EqualValues(t, "session-based message", messages[0].Body)
	require.EqualValues(t, "session-1", *messages[0].SessionID)
	require.NoError(t, receiver.CompleteMessage(ctx, messages[0], nil))

	require.EqualValues(t, "session-1", receiver.SessionID())

	sessionState, err := receiver.GetSessionState(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, sessionState)

	require.NoError(t, receiver.SetSessionState(ctx, []byte("hello"), nil))
	sessionState, err = receiver.GetSessionState(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, "hello", string(sessionState))

	// sanity check that we can clear out session state as well.
	require.NoError(t, receiver.SetSessionState(ctx, nil, nil))
	sessionState, err = receiver.GetSessionState(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, sessionState)
}

func TestSessionReceiver_blankSessionIDs(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &liveTestOptions{
		QueueProperties: &admin.QueueProperties{
			RequiresSession: to.Ptr(true),
		}})
	defer cleanup()

	ctx := context.Background()

	// send a message to a specific session
	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendMessage(ctx, &Message{
		Body:      []byte("session-based message"),
		SessionID: to.Ptr(""),
	}, nil)
	require.NoError(t, err)

	sequenceNumbers, err := sender.ScheduleMessages(ctx, []*Message{{
		Body:      []byte("session-based message"),
		SessionID: to.Ptr(""),
	}}, time.Now(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, sequenceNumbers)

	// start a receiver with the "" session ID
	receiver, err := client.AcceptSessionForQueue(ctx, queueName, "", nil)
	require.NoError(t, err)
	require.EqualValues(t, "", receiver.SessionID())

	var received []*ReceivedMessage

	receiveCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for {
		messages, err := receiver.ReceiveMessages(receiveCtx, 2, nil)
		require.NoError(t, err)

		for _, msg := range messages {
			require.NoError(t, receiver.CompleteMessage(ctx, msg, nil))
			received = append(received, msg)
		}

		if len(received) == 2 {
			break
		}
	}

	for _, msg := range received {
		require.EqualValues(t, "", *msg.SessionID)
		require.EqualValues(t, "session-based message", string(msg.Body))
	}
}

func TestSessionReceiver_acceptSessionButAlreadyLocked(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &liveTestOptions{
		QueueProperties: &admin.QueueProperties{
			RequiresSession: to.Ptr(true),
		}})
	defer cleanup()

	ctx := context.Background()

	receiver, err := client.AcceptSessionForQueue(ctx, queueName, "session-1", nil)
	require.NoError(t, err)
	require.NotNil(t, receiver)

	// You can address a session by name which makes lock contention possible (unlike
	// messages where the lock token is not a predefined value)
	receiver, err = client.AcceptSessionForQueue(ctx, queueName, "session-1", nil)

	require.EqualValues(t, internal_errors.RecoveryKindFatal, internal_errors.GetRecoveryKind(err))
	require.Nil(t, receiver)
}

func TestSessionReceiver_acceptNextSession_sessionExists(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &liveTestOptions{
		QueueProperties: &admin.QueueProperties{
			RequiresSession: to.Ptr(true),
		}})
	defer cleanup()

	ctx := context.Background()

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendMessage(ctx, &Message{
		Body:      []byte("session-based message"),
		SessionID: to.Ptr("acceptnextsession-test"),
	}, nil)
	require.NoError(t, err)

	// Using AcceptNextSessionForQueue will let the service determine the next 'available' session
	// This is useful for just round-robining through all the sessions that have data.
	receiver, err := client.AcceptNextSessionForQueue(ctx, queueName, nil)
	require.NoError(t, err)

	messages, err := receiver.inner.ReceiveMessages(ctx, 1, nil)
	require.NoError(t, err)

	require.EqualValues(t, "session-based message", messages[0].Body)
	require.EqualValues(t, "acceptnextsession-test", *messages[0].SessionID)
	require.NoError(t, receiver.CompleteMessage(ctx, messages[0], nil))

	require.EqualValues(t, "acceptnextsession-test", receiver.SessionID())

	sessionState, err := receiver.GetSessionState(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, sessionState)

	require.NoError(t, receiver.SetSessionState(ctx, []byte("hello"), nil))
	sessionState, err = receiver.GetSessionState(ctx, nil)
	require.NoError(t, err)
	require.EqualValues(t, "hello", string(sessionState))
}

func TestSessionReceiver_acceptNextSession_noSessionsExist(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &liveTestOptions{
		ClientOptions: &ClientOptions{
			RetryOptions: RetryOptions{
				MaxRetryDelay: time.Millisecond,
			},
		},
		QueueProperties: &admin.QueueProperties{
			RequiresSession: to.Ptr(true),
		}})
	defer cleanup()

	ctx := context.Background()

	// This adds an extra property to our link that sets a shorter scanning interval on the
	// service side. If you comment this out it will let the service determine the timeout,
	// which is 1 minute.
	client.acceptNextTimeout = time.Second

	receiver, err := client.AcceptNextSessionForQueue(ctx, queueName, nil)
	var sbErr *Error
	require.ErrorAs(t, err, &sbErr)
	require.Equal(t, CodeTimeout, sbErr.Code, "CodeTimeout since no sessions are available")
	require.Nil(t, receiver)
}

func TestSessionReceiver_nonSessionReceiver(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &liveTestOptions{
		QueueProperties: &admin.QueueProperties{
			RequiresSession: to.Ptr(true),
		}})
	defer cleanup()

	// opening a non-session full receiver fails and it's at least understandable
	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	// normal receivers are lazy initialized so we need to do _something_ to make sure
	// the link gets spun up (and thus fails)
	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.Nil(t, messages)

	var amqpError *amqp.Error
	require.True(t, errors.As(err, &amqpError))
	require.EqualValues(t, amqpError.Condition, "amqp:not-allowed")
	require.Contains(t, amqpError.Description, "It is not possible for an entity that requires sessions to create a non-sessionful message receiver.")

	messages, err = receiver.PeekMessages(context.Background(), 1, nil)
	require.Nil(t, messages)

	require.True(t, errors.As(err, &amqpError))
	require.EqualValues(t, amqpError.Condition, "amqp:not-allowed")
	require.Contains(t, amqpError.Description, "It is not possible for an entity that requires sessions to create a non-sessionful message receiver.")
}

func TestSessionReceiver_subscription(t *testing.T) {
	topic, cleanupTopic := createSubscription(t, nil, &admin.SubscriptionProperties{
		RequiresSession: to.Ptr(true),
	})

	defer cleanupTopic()

	client := newServiceBusClientForTest(t, &test.NewClientOptions[ClientOptions]{
		ClientOptions: &ClientOptions{
			RetryOptions: RetryOptions{
				MaxRetries: -1,
			},
		}})

	client.acceptNextTimeout = time.Second

	receiver, err := client.AcceptNextSessionForSubscription(context.Background(), topic, "sub", nil)
	require.Nil(t, receiver)

	var sbError *Error
	require.ErrorAs(t, err, &sbError)
	require.Equal(t, CodeTimeout, sbError.Code, "CodeTimeout because there are no sessions (yet)")

	sender, err := client.NewSender(topic, nil)
	require.NoError(t, err)

	defer sender.Close(context.Background())

	err = sender.SendMessage(context.Background(), &Message{
		SessionID: to.Ptr("session1"),
	}, nil)
	require.NoError(t, err)

	// there's a session this time...
	receiver, err = client.AcceptNextSessionForSubscription(context.Background(), topic, "sub", nil)
	require.NoError(t, err)

	defer receiver.Close(context.Background())

	require.Equal(t, "session1", receiver.SessionID())
	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	require.Equal(t, "session1", *messages[0].SessionID)
	require.Equal(t, 1, len(messages))

	err = receiver.CompleteMessage(context.Background(), messages[0], nil)
	require.NoError(t, err)
}

func TestSessionReceiver_RenewSessionLock(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &liveTestOptions{
		QueueProperties: &admin.QueueProperties{
			RequiresSession: to.Ptr(true),
		}})
	defer cleanup()

	sessionReceiver, err := client.AcceptSessionForQueue(context.Background(), queueName, "session-1", nil)
	require.NoError(t, err)

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body:      []byte("hello world"),
		SessionID: to.Ptr("session-1"),
	}, nil)
	require.NoError(t, err)

	messages, err := sessionReceiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.NotNil(t, messages)

	orig := sessionReceiver.LockedUntil()
	require.NoError(t, sessionReceiver.RenewSessionLock(context.Background(), nil))
	require.Greater(t, sessionReceiver.LockedUntil().UnixNano(), orig.UnixNano())
}

func TestSessionReceiver_Detach(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, &liveTestOptions{
		QueueProperties: &admin.QueueProperties{
			RequiresSession: to.Ptr(true),
		}})
	defer cleanup()

	test.EnableStdoutLogging(t)

	adminClient := newAdminClientForTest(t, nil)

	receiver, err := serviceBusClient.AcceptSessionForQueue(context.Background(), queueName, "test-session", nil)
	require.NoError(t, err)

	sender, err := serviceBusClient.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body:      []byte("hello"),
		SessionID: to.Ptr("test-session"),
	}, nil)
	require.NoError(t, err)
	require.NoError(t, sender.Close(context.Background()))

	state, err := receiver.GetSessionState(context.Background(), nil)
	require.NoError(t, err)
	require.Nil(t, state)

	// force a detach to happen
	_, err = adminClient.UpdateQueue(context.Background(), queueName, admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	}, nil)
	require.NoError(t, err)

	state, err = receiver.GetSessionState(context.Background(), nil)
	require.NoError(t, err)
	require.Nil(t, state)

	// force a detach to happen
	_, err = adminClient.UpdateQueue(context.Background(), queueName, admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	}, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.NotEmpty(t, messages)

	require.NoError(t, receiver.CompleteMessage(context.Background(), messages[0], nil))
	require.NoError(t, receiver.Close(context.Background()))
}

func TestSessionReceiver_roundRobin(t *testing.T) {
	client, cleanup, sessionEnabledQueueName := setupLiveTest(t, &liveTestOptions{
		QueueProperties: &admin.QueueProperties{
			RequiresSession: to.Ptr(true),
		},
		ClientOptions: &ClientOptions{
			RetryOptions: RetryOptions{
				MaxRetries: -1,
			},
		}})
	defer cleanup()

	sender, err := client.NewSender(sessionEnabledQueueName, nil)
	require.NoError(t, err)

	defer func() {
		err := sender.Close(context.Background())
		require.NoError(t, err)
	}()

	const maxSessions = 5

	for i := 0; i < maxSessions; i++ {
		sessionID := fmt.Sprintf("session%d", i)

		err = sender.SendMessage(context.Background(), &Message{
			Body:      []byte(sessionID),
			SessionID: to.Ptr(sessionID),
		}, nil)
		require.NoError(t, err)
	}

	receivedSessions := make(chan string, maxSessions)

	processMessageFromSession := func(ctx context.Context, receiver *SessionReceiver, message *ReceivedMessage) error {
		fmt.Printf("Completing messages for '%s'\n", *message.SessionID)
		err := receiver.CompleteMessage(ctx, message, nil)
		fmt.Printf("Completed messages for '%s'\n", *message.SessionID)

		require.Equal(t, receiver.SessionID(), *message.SessionID)
		receivedSessions <- *message.SessionID
		return err
	}

	client.acceptNextTimeout = time.Second

	// NOTE: this code is intentionally similar to `ExampleClient_AcceptNextSessionForQueue_roundrobin` so we can
	// test it.
	// BEGIN
	for {
		// You can have multiple active session receivers, provided they're each receiving
		// from different sessions.
		//
		// AcceptNextSessionForQueue (or AcceptNextSessionForSubscription) makes it simple to implement
		// this pattern, consuming multiple session receivers in parallel.
		sessionReceiver, err := client.AcceptNextSessionForQueue(context.TODO(), sessionEnabledQueueName, nil)

		if err != nil {
			var sbErr *Error

			if errors.As(err, &sbErr) && sbErr.Code == CodeTimeout {
				fmt.Printf("No sessions available\n")

				// NOTE: you could also continue here, which will block and wait again for a
				// session to become available.
				break
			}

			panic(err)
		}

		fmt.Printf("Got receiving for session '%s'\n", sessionReceiver.SessionID())

		// consume the session
		go func() {
			defer func() {
				fmt.Printf("Closing receiver for session '%s'\n", sessionReceiver.SessionID())
				err := sessionReceiver.Close(context.TODO())

				if err != nil {
					panic(err)
				}
			}()

			// we're only reading a few messages here, but you can also receive in a loop
			messages, err := sessionReceiver.ReceiveMessages(context.TODO(), 10, nil)

			if err != nil {
				panic(err)
			}

			fmt.Printf("Received %d messages from session '%s'\n", len(messages), sessionReceiver.SessionID())

			for _, m := range messages {
				if err := processMessageFromSession(context.TODO(), sessionReceiver, m); err != nil {
					panic(err)
				}
			}
		}()
	}

	// END

	t.Logf("Waiting for all sessions to complete")

	all := map[string]bool{}

	for i := 0; i < maxSessions; i++ {
		sessionID := <-receivedSessions
		t.Logf("Got %s's message", sessionID)
		all[sessionID] = true
	}

	require.Equal(t, maxSessions, len(all))
}

func Test_toReceiverOptions(t *testing.T) {
	require.Nil(t, toReceiverOptions(nil))

	require.EqualValues(t, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	}, toReceiverOptions(&SessionReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	}))
}

func TestSessionReceiverSendFiveReceiveFive_Queue(t *testing.T) {
	serviceBusClient, cleanup, queueName := setupLiveTest(t, &liveTestOptions{
		QueueProperties: &admin.QueueProperties{
			RequiresSession: to.Ptr(true),
		},
	})
	defer cleanup()

	sender, err := serviceBusClient.NewSender(queueName, nil)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	for i := 0; i < 5; i++ {
		err = sender.SendMessage(context.Background(), &Message{
			Body:      []byte(fmt.Sprintf("[%d]: send five, receive five", i)),
			SessionID: to.Ptr("session-1"),
		}, nil)
		require.NoError(t, err)
	}

	receiver, err := serviceBusClient.AcceptNextSessionForQueue(context.Background(), queueName, nil)
	require.NoError(t, err)

	messages := mustReceiveMessages(t, receiver, 5, time.Minute)

	for i := 0; i < 5; i++ {
		require.EqualValues(t,
			fmt.Sprintf("[%d]: send five, receive five", i),
			string(messages[i].Body))

		require.Equal(t, "session-1", *messages[i].SessionID)

		require.NoError(t, receiver.CompleteMessage(context.Background(), messages[i], nil))
	}
}

func TestSessionReceiverSendFiveReceiveFive_Subscription(t *testing.T) {
	serviceBusClient, cleanup, topicName, subscriptionName := setupLiveTestWithSubscription(t, &liveTestOptionsWithSubscription{
		SubscriptionProperties: &admin.SubscriptionProperties{
			RequiresSession: to.Ptr(true),
		},
	})
	defer cleanup()

	sender, err := serviceBusClient.NewSender(topicName, nil)
	require.NoError(t, err)
	defer sender.Close(context.Background())

	for i := 0; i < 5; i++ {
		err = sender.SendMessage(context.Background(), &Message{
			Body:      []byte(fmt.Sprintf("[%d]: send five, receive five", i)),
			SessionID: to.Ptr("session-1"),
		}, nil)
		require.NoError(t, err)
	}

	receiver, err := serviceBusClient.AcceptNextSessionForSubscription(context.Background(), topicName, subscriptionName, nil)
	require.NoError(t, err)

	messages := mustReceiveMessages(t, receiver, 5, time.Minute)

	for i := 0; i < 5; i++ {
		require.EqualValues(t,
			fmt.Sprintf("[%d]: send five, receive five", i),
			string(messages[i].Body))

		require.Equal(t, "session-1", *messages[i].SessionID)

		require.NoError(t, receiver.CompleteMessage(context.Background(), messages[i], nil))
	}
}

func mustReceiveMessages(t *testing.T, receiver interface {
	ReceiveMessages(ctx context.Context, count int, options *ReceiveMessagesOptions) ([]*ReceivedMessage, error)
}, count int, waitTime time.Duration) []*ReceivedMessage {
	var messages []*ReceivedMessage

	ctx, cancel := context.WithTimeout(context.Background(), waitTime)
	defer cancel()

	for {
		all, err := receiver.ReceiveMessages(ctx, count, nil)
		require.NoError(t, err)

		messages = append(messages, all...)

		if len(messages) > count {
			require.FailNowf(t, "Too many messages received", "Received more messages than expected: %d", len(messages))
		}

		if len(messages) == count {
			sort.Sort(receivedMessageSlice(messages))
			return messages
		}
	}
}
