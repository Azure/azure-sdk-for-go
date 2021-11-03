// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestSessionReceiver_acceptSession(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &QueueProperties{
		RequiresSession: to.BoolPtr(true),
	})
	defer cleanup()

	ctx := context.Background()

	// send a message to a specific session
	sender, err := client.NewSender(queueName)
	require.NoError(t, err)

	err = sender.SendMessage(ctx, &Message{
		Body:      []byte("session-based message"),
		SessionID: to.StringPtr("session-1"),
	})
	require.NoError(t, err)

	receiver, err := client.AcceptSessionForQueue(ctx, queueName, "session-1", nil)
	require.NoError(t, err)

	msg, err := receiver.inner.receiveMessage(ctx, nil)
	require.NoError(t, err)

	require.EqualValues(t, "session-based message", msg.Body)
	require.EqualValues(t, "session-1", *msg.SessionID)
	require.NoError(t, receiver.CompleteMessage(ctx, msg))

	require.EqualValues(t, "session-1", receiver.SessionID())

	sessionState, err := receiver.GetSessionState(ctx)
	require.NoError(t, err)
	require.Nil(t, sessionState)

	require.NoError(t, receiver.SetSessionState(ctx, []byte("hello")))
	sessionState, err = receiver.GetSessionState(ctx)
	require.NoError(t, err)
	require.EqualValues(t, "hello", string(sessionState))
}

func TestSessionReceiver_blankSessionIDs(t *testing.T) {
	t.Skip("Can't run blank session ID test because of issue")
	// errors while closing links: amqp sender close error: *Error{Condition: amqp:not-allowed, Description: The SessionId was not set on a message, and it cannot be sent to the entity. Entities that have session support enabled can only receive messages that have the SessionId set to a valid value.

	client, cleanup, queueName := setupLiveTest(t, &QueueProperties{
		RequiresSession: to.BoolPtr(true),
	})
	defer cleanup()

	ctx := context.Background()

	// send a message to a specific session
	sender, err := client.NewSender(queueName)
	require.NoError(t, err)

	err = sender.SendMessage(ctx, &Message{
		Body:      []byte("session-based message"),
		SessionID: to.StringPtr(""),
	})
	require.NoError(t, err)

	receiver, err := client.AcceptSessionForQueue(ctx, queueName, "", nil)
	require.NoError(t, err)

	msg, err := receiver.inner.receiveMessage(ctx, nil)
	require.NoError(t, err)

	require.EqualValues(t, "session-based message", msg.Body)
	require.EqualValues(t, "", *msg.SessionID)
	require.NoError(t, receiver.CompleteMessage(ctx, msg))

	require.EqualValues(t, "session-1", receiver.SessionID())
}

func TestSessionReceiver_acceptSessionButAlreadyLocked(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &QueueProperties{
		RequiresSession: to.BoolPtr(true),
	})
	defer cleanup()

	ctx := context.Background()

	receiver, err := client.AcceptSessionForQueue(ctx, queueName, "session-1", nil)
	require.NoError(t, err)
	require.NotNil(t, receiver)

	// You can address a session by name which makes lock contention possible (unlike
	// messages where the lock token is not a predefined value)
	receiver, err = client.AcceptSessionForQueue(ctx, queueName, "session-1", nil)
	require.True(t, internal.IsSessionLockedError(err))
	require.Nil(t, receiver)
}

func TestSessionReceiver_acceptNextSession(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &QueueProperties{
		RequiresSession: to.BoolPtr(true),
	})
	defer cleanup()

	ctx := context.Background()

	sender, err := client.NewSender(queueName)
	require.NoError(t, err)

	err = sender.SendMessage(ctx, &Message{
		Body:      []byte("session-based message"),
		SessionID: to.StringPtr("acceptnextsession-test"),
	})
	require.NoError(t, err)

	// Using AcceptNextSessionForQueue will let the service determine the next 'available' session
	// This is useful for just round-robining through all the sessions that have data.
	receiver, err := client.AcceptNextSessionForQueue(ctx, queueName, nil)
	require.NoError(t, err)

	msg, err := receiver.inner.receiveMessage(ctx, nil)
	require.NoError(t, err)

	require.EqualValues(t, "session-based message", msg.Body)
	require.EqualValues(t, "acceptnextsession-test", *msg.SessionID)
	require.NoError(t, receiver.CompleteMessage(ctx, msg))

	require.EqualValues(t, "acceptnextsession-test", receiver.SessionID())

	sessionState, err := receiver.GetSessionState(ctx)
	require.NoError(t, err)
	require.Nil(t, sessionState)

	require.NoError(t, receiver.SetSessionState(ctx, []byte("hello")))
	sessionState, err = receiver.GetSessionState(ctx)
	require.NoError(t, err)
	require.EqualValues(t, "hello", string(sessionState))
}

func TestSessionReceiver_noSessionsAvailable(t *testing.T) {
	t.Skip("Really slow test (since it has to wait for a timeout from the service)")

	client, cleanup, queueName := setupLiveTest(t, &QueueProperties{
		RequiresSession: to.BoolPtr(true),
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
	client, cleanup, queueName := setupLiveTest(t, &QueueProperties{
		RequiresSession: to.BoolPtr(true),
	})
	defer cleanup()

	// opening a non-session full receiver fails and it's at least understandable
	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	// normal receivers are lazy initialized so we need to do _something_ to make sure
	// the link gets spun up (and thus fails)
	message, err := receiver.receiveMessage(context.Background(), nil)
	require.Nil(t, message)

	var amqpError *amqp.Error
	require.True(t, errors.As(err, &amqpError))
	require.EqualValues(t, amqpError.Condition, "amqp:not-allowed")
	require.Contains(t, amqpError.Description, "It is not possible for an entity that requires sessions to create a non-sessionful message receiver.")

	messages, err := receiver.PeekMessages(context.Background(), 1, nil)
	require.Nil(t, messages)

	require.True(t, errors.As(err, &amqpError))
	require.EqualValues(t, amqpError.Condition, "amqp:not-allowed")
	require.Contains(t, amqpError.Description, "It is not possible for an entity that requires sessions to create a non-sessionful message receiver.")
}

func TestSessionReceiver_RenewSessionLock(t *testing.T) {
	client, cleanup, queueName := setupLiveTest(t, &QueueProperties{
		RequiresSession: to.BoolPtr(true),
	})
	defer cleanup()

	sessionReceiver, err := client.AcceptSessionForQueue(context.Background(), queueName, "session-1", nil)
	require.NoError(t, err)

	sender, err := client.NewSender(queueName)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body:      []byte("hello world"),
		SessionID: to.StringPtr("session-1"),
	})
	require.NoError(t, err)

	messages, err := sessionReceiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)
	require.NotNil(t, messages)

	// surprisingly this works. Not sure what it accomplishes though. C# has a manual check for it.
	// err = sessionReceiver.RenewMessageLock(context.Background(), messages[0])
	// require.NoError(t, err)

	orig := sessionReceiver.LockedUntil()
	require.NoError(t, sessionReceiver.RenewSessionLock(context.Background()))
	require.Greater(t, sessionReceiver.LockedUntil().UnixNano(), orig.UnixNano())

	// bogus renewal
	sessionReceiver.sessionID = to.StringPtr("bogus")

	err = sessionReceiver.RenewSessionLock(context.Background())
	require.Contains(t, err.Error(), "status code 410 and description: The session lock has expired on the MessageSession")
}

func Test_toReceiverOptions(t *testing.T) {
	require.Nil(t, toReceiverOptions(nil))

	require.EqualValues(t, &ReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	}, toReceiverOptions(&SessionReceiverOptions{
		ReceiveMode: ReceiveModeReceiveAndDelete,
	}))
}
