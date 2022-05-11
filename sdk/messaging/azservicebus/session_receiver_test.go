// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestSessionReceiver_acceptSession(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})
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
}

func TestSessionReceiver_blankSessionIDs(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})
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
	client, cleanup, queueName := setupLiveTest(t, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})
	defer cleanup()

	ctx := context.Background()

	receiver, err := client.AcceptSessionForQueue(ctx, queueName, "session-1", nil)
	require.NoError(t, err)
	require.NotNil(t, receiver)

	// You can address a session by name which makes lock contention possible (unlike
	// messages where the lock token is not a predefined value)
	receiver, err = client.AcceptSessionForQueue(ctx, queueName, "session-1", nil)

	require.EqualValues(t, internal.RecoveryKindFatal, internal.GetRecoveryKind(err))
	require.Nil(t, receiver)
}

func TestSessionReceiver_acceptNextSession(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})
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

func TestSessionReceiver_noSessionsAvailable(t *testing.T) {
	t.Skip("Really slow test (since it has to wait for a timeout from the service)")

	client, cleanup, queueName := setupLiveTest(t, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})
	defer cleanup()

	ctx := context.Background()

	// now try again (there are no sessions available with messages so we expect a failure
	receiver, err := client.AcceptNextSessionForQueue(ctx, queueName, nil)

	var amqpError *amqp.Error
	require.True(t, errors.As(err, &amqpError))

	// just documenting this for now - if this is consistent enough we'll want to provide some guidance as this is a normal condition but
	// it looks like it's an error instead.
	// (actual error)
	// *Error{Condition: com.microsoft:timeout, Description: The operation did not complete within the allotted timeout of 00:01:04.9800000. The time allotted to this operation may have been a portion of a longer timeout.
	require.EqualValues(t, amqpError.Condition, "com.microsoft:timeout")

	require.Nil(t, receiver)
}

func TestSessionReceiver_nonSessionReceiver(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})
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

func TestSessionReceiver_RenewSessionLock(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})
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
	serviceBusClient, cleanup, queueName := setupLiveTest(t, &admin.QueueProperties{
		RequiresSession: to.Ptr(true),
	})
	defer cleanup()

	azlog.SetListener(func(e azlog.Event, s string) {
		fmt.Printf("%s %s\n", e, s)
	})

	defer azlog.SetListener(nil)

	adminClient, err := admin.NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

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

func Test_toReceiverOptions(t *testing.T) {
	require.Nil(t, toReceiverOptions(nil))

	require.EqualValues(t, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	}, toReceiverOptions(&SessionReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	}))
}
