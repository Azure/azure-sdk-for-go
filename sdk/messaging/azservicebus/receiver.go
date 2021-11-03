// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/go-amqp"
	"github.com/devigned/tab"
)

// ReceiveMode represents the lock style to use for a receiver - either
// `PeekLock` or `ReceiveAndDelete`
type ReceiveMode = internal.ReceiveMode

const (
	// ReceiveModePeekLock will lock messages as they are received and can be settled
	// using the Receiver's (Complete|Abandon|DeadLetter|Defer)Message
	// functions.
	ReceiveModePeekLock ReceiveMode = internal.PeekLock
	// ReceiveModeReceiveAndDelete will delete messages as they are received.
	ReceiveModeReceiveAndDelete ReceiveMode = internal.ReceiveAndDelete
)

// SubQueue allows you to target a subqueue of a queue or subscription.
// Ex: the dead letter queue (SubQueueDeadLetter).
type SubQueue int

const (
	// SubQueueNone means no sub queue.
	SubQueueNone SubQueue = 0
	// SubQueueDeadLetter targets the dead letter queue for a queue or subscription.
	SubQueueDeadLetter SubQueue = 1
	// SubQueueTransfer targets the transfer dead letter queue for a queue or subscription.
	SubQueueTransfer SubQueue = 2
)

// Receiver receives messages using pull based functions (ReceiveMessages).
type Receiver struct {
	receiveMode ReceiveMode

	settler        settler
	baseRetrier    internal.Retrier
	cleanupOnClose func()

	lastPeekedSequenceNumber int64
	amqpLinks                internal.AMQPLinks

	mu        sync.Mutex
	receiving bool
}

// ReceiverOptions contains options for the `Client.NewReceiverForQueue` or `Client.NewReceiverForSubscription`
// functions.
type ReceiverOptions struct {
	// ReceiveMode controls when a message is deleted from Service Bus.
	//
	// `azservicebus.PeekLock` is the default. The message is locked, preventing multiple
	// receivers from processing the message at once. You control the lock state of the message
	// using one of the message settlement functions like Receiver.CompleteMessage(), which removes
	// it from Service Bus, or Receiver.AbandonMessage(), which makes it available again.
	//
	// `azservicebus.ReceiveAndDelete` causes Service Bus to remove the message as soon
	// as it's received.
	//
	// More information about receive modes:
	// https://docs.microsoft.com/azure/service-bus-messaging/message-transfers-locks-settlement#settling-receive-operations
	ReceiveMode ReceiveMode

	// SubQueue should be set to connect to the sub queue (ex: dead letter queue)
	// of the queue or subscription.
	SubQueue SubQueue
}

const defaultLinkRxBuffer = 2048

func applyReceiverOptions(receiver *Receiver, entity *entity, options *ReceiverOptions) error {
	if options == nil {
		receiver.receiveMode = ReceiveModePeekLock
		return nil
	}

	if err := checkReceiverMode(options.ReceiveMode); err != nil {
		return err
	}

	receiver.receiveMode = options.ReceiveMode

	if err := entity.SetSubQueue(options.SubQueue); err != nil {
		return err
	}

	return nil
}

func newReceiver(ns internal.NamespaceWithNewAMQPLinks, entity *entity, cleanupOnClose func(), options *ReceiverOptions, newLinksFn func(ctx context.Context, session internal.AMQPSession) (internal.AMQPSenderCloser, internal.AMQPReceiverCloser, error)) (*Receiver, error) {
	receiver := &Receiver{
		lastPeekedSequenceNumber: 0,
		// TODO: make this configurable
		baseRetrier: internal.NewBackoffRetrier(internal.BackoffRetrierParams{
			Factor:     1.5,
			Jitter:     true,
			Min:        time.Second,
			Max:        time.Minute,
			MaxRetries: 10,
		}),
		cleanupOnClose: cleanupOnClose,
	}

	if err := applyReceiverOptions(receiver, entity, options); err != nil {
		return nil, err
	}

	entityPath, err := entity.String()

	if err != nil {
		return nil, err
	}

	if newLinksFn == nil {
		newLinksFn = func(ctx context.Context, session internal.AMQPSession) (internal.AMQPSenderCloser, internal.AMQPReceiverCloser, error) {
			linkOptions := createLinkOptions(receiver.receiveMode, entityPath)
			return createReceiverLink(ctx, session, linkOptions)
		}
	}

	receiver.amqpLinks = ns.NewAMQPLinks(entityPath, newLinksFn)

	// 'nil' settler handles returning an error message for receiveAndDelete links.
	if receiver.receiveMode == ReceiveModePeekLock {
		receiver.settler = newMessageSettler(receiver.amqpLinks, receiver.baseRetrier)
	}

	return receiver, nil
}

// ReceiveOptions are options for the ReceiveMessages function.
type ReceiveOptions struct {
	maxWaitTimeAfterFirstMessage time.Duration
}

// ReceiveMessages receives a fixed number of messages, up to numMessages.
// There are two ways to stop receiving messages:
// 1. Cancelling the `ctx` parameter.
// 2. An implicit timeout (default: 1 second) that starts after the first
//    message has been received.
func (r *Receiver) ReceiveMessages(ctx context.Context, maxMessages int, options *ReceiveOptions) ([]*ReceivedMessage, error) {
	r.mu.Lock()
	isReceiving := r.receiving

	if !isReceiving {
		r.receiving = true

		defer func() {
			r.mu.Lock()
			r.receiving = false
			r.mu.Unlock()
		}()
	}
	r.mu.Unlock()

	if isReceiving {
		return nil, errors.New("receiver is already receiving messages. ReceiveMessages() cannot be called concurrently")
	}

	return r.receiveMessagesImpl(ctx, maxMessages, options)
}

// ReceiveDeferredMessages receives messages that were deferred using `Receiver.DeferMessage`.
func (r *Receiver) ReceiveDeferredMessages(ctx context.Context, sequenceNumbers []int64) ([]*ReceivedMessage, error) {
	_, _, mgmt, _, err := r.amqpLinks.Get(ctx)

	if err != nil {
		return nil, err
	}

	amqpMessages, err := mgmt.ReceiveDeferred(ctx, r.receiveMode, sequenceNumbers)

	if err != nil {
		return nil, err
	}

	var receivedMessages []*ReceivedMessage

	for _, amqpMsg := range amqpMessages {
		receivedMsg := newReceivedMessage(ctx, amqpMsg)
		receivedMsg.deferred = true

		receivedMessages = append(receivedMessages, receivedMsg)
	}

	return receivedMessages, nil
}

// PeekMessagesOptions contains options for the `Receiver.PeekMessages`
// function.
type PeekMessagesOptions struct {
	// FromSequenceNumber is the sequence number to start with when peeking messages.
	FromSequenceNumber *int64
}

// PeekMessages will peek messages without locking or deleting messages.
// Messages that are peeked do not have lock tokens, so settlement methods
// like CompleteMessage, AbandonMessage, DeferMessage or DeadLetterMessage
// will not work with them.
func (r *Receiver) PeekMessages(ctx context.Context, maxMessageCount int, options *PeekMessagesOptions) ([]*ReceivedMessage, error) {
	_, _, mgmt, _, err := r.amqpLinks.Get(ctx)

	if err != nil {
		return nil, err
	}

	var sequenceNumber = r.lastPeekedSequenceNumber + 1
	updateInternalSequenceNumber := true

	if options != nil && options.FromSequenceNumber != nil {
		sequenceNumber = *options.FromSequenceNumber
		updateInternalSequenceNumber = false
	}

	messages, err := mgmt.PeekMessages(ctx, sequenceNumber, int32(maxMessageCount))

	if err != nil {
		return nil, err
	}

	receivedMessages := make([]*ReceivedMessage, len(messages))

	for i := 0; i < len(messages); i++ {
		receivedMessages[i] = newReceivedMessage(ctx, messages[i])
	}

	if len(receivedMessages) > 0 && updateInternalSequenceNumber {
		// only update this if they're doing the implicit iteration as part of the receiver.
		r.lastPeekedSequenceNumber = *receivedMessages[len(receivedMessages)-1].SequenceNumber
	}

	return receivedMessages, nil
}

// RenewLock renews the lock on a message, updating the `LockedUntil` field on `msg`.
func (r *Receiver) RenewMessageLock(ctx context.Context, msg *ReceivedMessage) error {
	_, _, mgmt, _, err := r.amqpLinks.Get(ctx)

	if err != nil {
		return err
	}

	newExpirationTime, err := mgmt.RenewLocks(ctx, msg.rawAMQPMessage.LinkName(), []amqp.UUID{
		(amqp.UUID)(msg.LockToken),
	})

	if err != nil {
		return err
	}

	msg.LockedUntil = &newExpirationTime[0]
	return nil
}

// Close permanently closes the receiver.
func (r *Receiver) Close(ctx context.Context) error {
	r.cleanupOnClose()
	return r.amqpLinks.Close(ctx, true)
}

// CompleteMessage completes a message, deleting it from the queue or subscription.
func (r *Receiver) CompleteMessage(ctx context.Context, message *ReceivedMessage) error {
	return r.settler.CompleteMessage(ctx, message)
}

// AbandonMessage will cause a message to be returned to the queue or subscription.
// This will increment its delivery count, and potentially cause it to be dead lettered
// depending on your queue or subscription's configuration.
func (r *Receiver) AbandonMessage(ctx context.Context, message *ReceivedMessage) error {
	return r.settler.AbandonMessage(ctx, message)
}

// DeferMessage will cause a message to be deferred. Deferred messages
// can be received using `Receiver.ReceiveDeferredMessages`.
func (r *Receiver) DeferMessage(ctx context.Context, message *ReceivedMessage) error {
	return r.settler.DeferMessage(ctx, message)
}

// DeadLetterMessage settles a message by moving it to the dead letter queue for a
// queue or subscription. To receive these messages create a receiver with `Client.NewReceiverForQueue()`
// or `Client.NewReceiverForSubscription()` using the `ReceiverOptions.SubQueue` option.
func (r *Receiver) DeadLetterMessage(ctx context.Context, message *ReceivedMessage, options *DeadLetterOptions) error {
	return r.settler.DeadLetterMessage(ctx, message, options)
}

// receiveDeferredMessage receives a single message that was deferred using `Receiver.DeferMessage`.
func (r *Receiver) receiveDeferredMessage(ctx context.Context, sequenceNumber int64) (*ReceivedMessage, error) {
	messages, err := r.ReceiveDeferredMessages(ctx, []int64{sequenceNumber})

	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages[0], nil
}

// receiveMessage receives a single message, waiting up to `ReceiveOptions.MaxWaitTime` (default: 60 seconds)
func (r *Receiver) receiveMessage(ctx context.Context, options *ReceiveOptions) (*ReceivedMessage, error) {
	messages, err := r.ReceiveMessages(ctx, 1, options)

	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages[0], nil
}

func (r *Receiver) receiveMessagesImpl(ctx context.Context, maxMessages int, options *ReceiveOptions) ([]*ReceivedMessage, error) {
	// There are three phases for this function:
	// Phase 1. <receive, respecting user cancellation>
	// Phase 2. <check error and exit if fatal>
	//    NOTE: We don't exit here so we don't end up buffering messages internally that the
	//    user isn't actually waiting for anymore. So we make sure that #3 runs if the
	//    link is still valid.
	// Phase 3. <drain the link and leave it in a good state>
	localOpts := &ReceiveOptions{
		maxWaitTimeAfterFirstMessage: time.Second,
	}

	if options != nil {
		if options.maxWaitTimeAfterFirstMessage != 0 {
			localOpts.maxWaitTimeAfterFirstMessage = options.maxWaitTimeAfterFirstMessage
		}
	}

	_, receiver, _, linksRevision, err := r.amqpLinks.Get(ctx)

	if err != nil {
		if err := r.amqpLinks.RecoverIfNeeded(ctx, linksRevision, err); err != nil {
			return nil, err
		}

		return nil, err
	}

	if err := receiver.IssueCredit(uint32(maxMessages)); err != nil {
		_ = r.amqpLinks.RecoverIfNeeded(ctx, linksRevision, err)
		return nil, err
	}

	messages, err := r.getMessages(ctx, receiver, maxMessages, localOpts)

	if err != nil {
		return nil, err
	}

	if len(messages) == maxMessages {
		// no drain needed, all messages arrived.
		return messages, nil
	}

	return r.drainLink(receiver, messages)
}

// drainLink initiates a drainLink on the link. Service Bus will send whatever messages it might have still had and
// set our link credit to 0.
func (r *Receiver) drainLink(receiver internal.AMQPReceiver, messages []*ReceivedMessage) ([]*ReceivedMessage, error) {
	receiveCtx, cancelReceive := context.WithCancel(context.Background())

	// start the drain asynchronously. Note that we ignore the user's context at this point
	// since draining makes sure we don't get messages when nobody is receiving.
	go func() {
		if err := receiver.DrainCredit(context.Background()); err != nil {
			tab.For(receiveCtx).Debug(fmt.Sprintf("Draining of credit failed. link will be closed and will re-open on next receive: %s", err.Error()))

			// if the drain fails we just close the link so it'll re-open at the next receive.
			if err := r.amqpLinks.Close(context.Background(), false); err != nil {
				tab.For(receiveCtx).Debug(fmt.Sprintf("Failed to close links on ReceiveMessages cleanup. Not fatal: %s", err.Error()))
			}
		}
		cancelReceive()
	}()

	// Receive until the drain completes, at which point it'll cancel
	// our context.
	// NOTE: That's a gap here where we need to be able to drain _only_ the internally cached messages
	// in the receiver. Filed as https://github.com/Azure/go-amqp/issues/71
	for {
		am, err := receiver.Receive(receiveCtx)

		if internal.IsCancelError(err) {
			break
		} else if err != nil {
			// something fatal happened, we will just
			_ = r.amqpLinks.Close(context.TODO(), false)

			if len(messages) > 0 {
				return messages, nil
			} else {
				return nil, err
			}
		}

		messages = append(messages, newReceivedMessage(receiveCtx, am))
	}

	return messages, nil
}

// getMessages receives messages until a link failure, timeout or the user
// cancels their context.
func (r *Receiver) getMessages(ctx context.Context, receiver internal.AMQPReceiver, maxMessages int, ropts *ReceiveOptions) ([]*ReceivedMessage, error) {
	var messages []*ReceivedMessage

	for {
		var amqpMessage *amqp.Message
		amqpMessage, err := receiver.Receive(ctx)

		if err != nil {
			if internal.IsCancelError(err) {
				return messages, nil
			}

			// we'll close (instead of recovering) since we are holding onto messages
			// and want to get them back to the user ASAP. (recovery will just happen
			// on the next call to receive)
			if err := r.amqpLinks.Close(context.Background(), false); err != nil {
				tab.For(ctx).Debug(fmt.Sprintf("Failed to close links on ReceiveMessages cleanup. Not fatal: %s", err.Error()))
			}
			return nil, err
		}

		messages = append(messages, newReceivedMessage(ctx, amqpMessage))

		if len(messages) == maxMessages {
			return messages, nil
		}

		if len(messages) == 1 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, time.Second)
			defer cancel()
		}
	}
}

type entity struct {
	subqueue     SubQueue
	Queue        string
	Topic        string
	Subscription string
}

func (e *entity) String() (string, error) {
	entityPath := ""

	if e.Queue != "" {
		entityPath = e.Queue
	} else if e.Topic != "" && e.Subscription != "" {
		entityPath = fmt.Sprintf("%s/Subscriptions/%s", e.Topic, e.Subscription)
	} else {
		return "", errors.New("a queue or subscription was not specified")
	}

	if e.subqueue == SubQueueDeadLetter {
		entityPath += "/$DeadLetterQueue"
	} else if e.subqueue == SubQueueTransfer {
		entityPath += "/$Transfer/$DeadLetterQueue"
	}

	return entityPath, nil
}

func (e *entity) SetSubQueue(subQueue SubQueue) error {
	if subQueue == SubQueueNone {
		return nil
	} else if subQueue == SubQueueDeadLetter || subQueue == SubQueueTransfer {
		e.subqueue = subQueue
		return nil
	}

	return fmt.Errorf("unknown SubQueue %d", subQueue)
}

func createReceiverLink(ctx context.Context, session internal.AMQPSession, linkOptions []amqp.LinkOption) (internal.AMQPSenderCloser, internal.AMQPReceiverCloser, error) {
	amqpReceiver, err := session.NewReceiver(linkOptions...)

	if err != nil {
		tab.For(ctx).Error(err)
		return nil, nil, err
	}

	return nil, amqpReceiver, nil
}

func createLinkOptions(mode ReceiveMode, entityPath string) []amqp.LinkOption {
	receiveMode := amqp.ModeSecond

	if mode == ReceiveModeReceiveAndDelete {
		receiveMode = amqp.ModeFirst
	}

	opts := []amqp.LinkOption{
		amqp.LinkSourceAddress(entityPath),
		amqp.LinkReceiverSettle(receiveMode),
		amqp.LinkWithManualCredits(),
		amqp.LinkCredit(defaultLinkRxBuffer),
	}

	if mode == ReceiveModeReceiveAndDelete {
		opts = append(opts, amqp.LinkSenderSettle(amqp.ModeSettled))
	}

	return opts
}
