// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/go-amqp"
)

// SessionReceiver is a Receiver that handles sessions.
type SessionReceiver struct {
	inner       *Receiver
	sessionID   *string
	lockedUntil time.Time
}

// SessionReceiverOptions contains options for the `Client.AcceptSessionForQueue/Subscription` or `Client.AcceptNextSessionForQueue/Subscription`
// functions.
type SessionReceiverOptions struct {
	// ReceiveMode controls when a message is deleted from Service Bus.
	//
	// `azservicebus.PeekLock` is the default. The message is locked, preventing multiple
	// receivers from processing the message at once. You control the lock state of the message
	// using one of the message settlement functions like SessionReceiver.CompleteMessage(), which removes
	// it from Service Bus, or SessionReceiver..AbandonMessage(), which makes it available again.
	//
	// `azservicebus.ReceiveAndDelete` causes Service Bus to remove the message as soon
	// as it's received.
	//
	// More information about receive modes:
	// https://docs.microsoft.com/azure/service-bus-messaging/message-transfers-locks-settlement#settling-receive-operations
	ReceiveMode ReceiveMode
}

func toReceiverOptions(sropts *SessionReceiverOptions) *ReceiverOptions {
	if sropts == nil {
		return nil
	}

	return &ReceiverOptions{
		ReceiveMode: sropts.ReceiveMode,
	}
}

func newSessionReceiver(ctx context.Context, sessionID *string, ns internal.NamespaceWithNewAMQPLinks, entity *entity, cleanupOnClose func(), options *ReceiverOptions) (*SessionReceiver, error) {
	const sessionFilterName = "com.microsoft:session-filter"
	const code = uint64(0x00000137000000C)

	sessionReceiver := &SessionReceiver{
		sessionID:   sessionID,
		lockedUntil: time.Time{},
	}

	var err error

	sessionReceiver.inner, err = newReceiver(ns, entity, cleanupOnClose, options, func(ctx context.Context, session internal.AMQPSession) (internal.AMQPSenderCloser, internal.AMQPReceiverCloser, error) {
		linkOptions := createLinkOptions(sessionReceiver.inner.receiveMode, sessionReceiver.inner.amqpLinks.EntityPath())

		if sessionID == nil {
			linkOptions = append(linkOptions, amqp.LinkSourceFilter(sessionFilterName, code, nil))
		} else {
			linkOptions = append(linkOptions, amqp.LinkSourceFilter(sessionFilterName, code, sessionID))
		}

		_, link, err := createReceiverLink(ctx, session, linkOptions)

		if err != nil {
			return nil, nil, err
		}

		// check the session ID that came back - if we asked for a named session ID and didn't get it then
		// we failed to get the lock.
		// if we specified nil then we can _set_ our internally held session ID now that we know the value.
		receivedSessionID := link.LinkSourceFilterValue(sessionFilterName)
		asStr, ok := receivedSessionID.(string)

		if !ok || (sessionID != nil && asStr != *sessionID) {
			return nil, nil, fmt.Errorf("invalid type/value for returned sessionID(type:%T, value:%v)", receivedSessionID, receivedSessionID)
		}

		sessionReceiver.sessionID = &asStr
		return nil, link, nil
	})

	if err != nil {
		return nil, err
	}

	// temp workaround until we expose the session expiration time from the receiver in go-amqp
	if err := sessionReceiver.RenewSessionLock(ctx); err != nil {
		_ = sessionReceiver.Close(context.Background())
		return nil, err
	}

	return sessionReceiver, nil
}

// ReceiveMessages receives a fixed number of messages, up to numMessages.
// There are two ways to stop receiving messages:
// 1. Cancelling the `ctx` parameter.
// 2. An implicit timeout (default: 1 second) that starts after the first
//    message has been received.
func (r *SessionReceiver) ReceiveMessages(ctx context.Context, maxMessages int, options *ReceiveOptions) ([]*ReceivedMessage, error) {
	return r.inner.ReceiveMessages(ctx, maxMessages, options)
}

// ReceiveDeferredMessages receives messages that were deferred using `Receiver.DeferMessage`.
func (r *SessionReceiver) ReceiveDeferredMessages(ctx context.Context, sequenceNumbers []int64) ([]*ReceivedMessage, error) {
	return r.inner.ReceiveDeferredMessages(ctx, sequenceNumbers)
}

// PeekMessages will peek messages without locking or deleting messages.
// Messages that are peeked do not have lock tokens, so settlement methods
// like CompleteMessage, AbandonMessage, DeferMessage or DeadLetterMessage
// will not work with them.
func (r *SessionReceiver) PeekMessages(ctx context.Context, maxMessageCount int, options *PeekMessagesOptions) ([]*ReceivedMessage, error) {
	return r.inner.PeekMessages(ctx, maxMessageCount, options)
}

// RenewLock renews the lock on a message, updating the `LockedUntil` field on `msg`.
func (r *SessionReceiver) RenewMessageLock(ctx context.Context, msg *ReceivedMessage) error {
	return r.inner.RenewMessageLock(ctx, msg)
}

// Close permanently closes the receiver.
func (r *SessionReceiver) Close(ctx context.Context) error {
	return r.inner.Close(ctx)
}

// CompleteMessage completes a message, deleting it from the queue or subscription.
func (r *SessionReceiver) CompleteMessage(ctx context.Context, message *ReceivedMessage) error {
	return r.inner.CompleteMessage(ctx, message)
}

// AbandonMessage will cause a message to be returned to the queue or subscription.
// This will increment its delivery count, and potentially cause it to be dead lettered
// depending on your queue or subscription's configuration.
func (r *SessionReceiver) AbandonMessage(ctx context.Context, message *ReceivedMessage) error {
	return r.inner.AbandonMessage(ctx, message)
}

// DeferMessage will cause a message to be deferred. Deferred messages
// can be received using `Receiver.ReceiveDeferredMessages`.
func (r *SessionReceiver) DeferMessage(ctx context.Context, message *ReceivedMessage) error {
	return r.inner.DeferMessage(ctx, message)
}

// DeadLetterMessage settles a message by moving it to the dead letter queue for a
// queue or subscription. To receive these messages create a receiver with `Client.NewReceiverForQueue()`
// or `Client.NewReceiverForSubscription()` using the `ReceiverOptions.SubQueue` option.
func (r *SessionReceiver) DeadLetterMessage(ctx context.Context, message *ReceivedMessage, options *DeadLetterOptions) error {
	return r.inner.DeadLetterMessage(ctx, message, options)
}

// SessionID is the session ID for this SessionReceiver.
func (sr *SessionReceiver) SessionID() string {
	// return the ultimately assigned session ID for this link (anonymous will get it from the
	// link filter options, non-anonymous is set in newSessionReceiver)
	return *sr.sessionID
}

// LockedUntil is the time the lock on this session expires.
// The lock can be renewed using `SessionReceiver.RenewSessionLock`.
func (sr *SessionReceiver) LockedUntil() time.Time {
	return sr.lockedUntil
}

// GetSessionState retrieves state associated with the session.
func (sr *SessionReceiver) GetSessionState(ctx context.Context) ([]byte, error) {
	_, _, mgmt, _, err := sr.inner.amqpLinks.Get(ctx)

	if err != nil {
		return nil, err
	}

	return mgmt.GetSessionState(ctx, sr.SessionID())
}

// SetSessionState sets the state associated with the session.
func (sr *SessionReceiver) SetSessionState(ctx context.Context, state []byte) error {
	_, _, mgmt, _, err := sr.inner.amqpLinks.Get(ctx)

	if err != nil {
		return err
	}

	return mgmt.SetSessionState(ctx, sr.SessionID(), state)
}

// RenewSessionLock renews this session's lock. The new expiration time is available
// using `LockedUntil`.
func (sr *SessionReceiver) RenewSessionLock(ctx context.Context) error {
	_, _, mgmt, _, err := sr.inner.amqpLinks.Get(ctx)

	if err != nil {
		return err
	}

	newLockedUntil, err := mgmt.RenewSessionLock(ctx, *sr.sessionID)

	if err != nil {
		return err
	}

	sr.lockedUntil = newLockedUntil
	return nil
}

// init ensures the link was created, guaranteeing that we get our expected session lock.
func (sr *SessionReceiver) init(ctx context.Context) error {
	// initialize the links
	_, _, _, _, err := sr.inner.amqpLinks.Get(ctx)
	return err
}
