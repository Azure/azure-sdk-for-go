// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/go-amqp"
	"github.com/devigned/tab"
)

// ReceiveMode represents the lock style to use for a receiver - either
// `PeekLock` or `ReceiveAndDelete`
type ReceiveMode = internal.ReceiveMode

const (
	// PeekLock will lock messages as they are received and can be settled
	// using the Receiver or Processor's (Complete|Abandon|DeadLetter|Defer)Message
	// functions.
	PeekLock ReceiveMode = internal.PeekLock
	// ReceiveAndDelete will delete messages as they are received.
	ReceiveAndDelete ReceiveMode = internal.ReceiveAndDelete
)

// SubQueue allows you to target a subqueue of a queue or subscription.
// Ex: the dead letter queue (SubQueueDeadLetter).
type SubQueue int

const (
	// SubQueueDeadLetter targets the dead letter queue for a queue or subscription.
	SubQueueDeadLetter SubQueue = 1
	// SubQueueTransfer targets the transfer dead letter queue for a queue or subscription.
	SubQueueTransfer SubQueue = 2
)

// Receiver receives messages using pull based functions (ReceiveMessages).
// For push-based receiving via callbacks look at the `Processor` type.
type Receiver struct {
	settler *messageSettler

	config struct {
		ReceiveMode    ReceiveMode
		FullEntityPath string
		Entity         entity

		RetryOptions struct {
			Times int
			Delay time.Duration
		}
	}

	lastPeekedSequenceNumber int64
	amqpLinks                internal.AMQPLinks
}

const defaultLinkRxBuffer = 2048

// ReceiverOption represents an option for a receiver.
// Some examples:
// - `ReceiverWithReceiveMode` to configure the receive mode,
// - `ReceiverWithQueue` to target a queue.
type ReceiverOption func(receiver *Receiver) error

// ReceiverWithSubQueue allows you to open the sub queue (ie: dead letter queues, transfer dead letter queues)
// for a queue or subscription.
func ReceiverWithSubQueue(subQueue SubQueue) ReceiverOption {
	return func(receiver *Receiver) error {
		if subQueue == SubQueueDeadLetter || subQueue == SubQueueTransfer {
			receiver.config.Entity.Subqueue = subQueue
		} else {
			return fmt.Errorf("unknown SubQueue %d", subQueue)
		}

		return nil
	}
}

// ReceiverWithReceiveMode controls the receive mode for the receiver.
func ReceiverWithReceiveMode(receiveMode ReceiveMode) ReceiverOption {
	return func(receiver *Receiver) error {
		if receiveMode != PeekLock && receiveMode != ReceiveAndDelete {
			return fmt.Errorf("invalid receive mode specified %d", receiveMode)
		}

		receiver.config.ReceiveMode = receiveMode
		return nil
	}
}

// ReceiverWithQueue configures a receiver to connect to a queue.
func ReceiverWithQueue(queue string) ReceiverOption {
	return func(receiver *Receiver) error {
		receiver.config.Entity.Queue = queue
		return nil
	}
}

// ReceiverWithSubscription configures a receiver to connect to a subscription
// associated with a topic.
func ReceiverWithSubscription(topic string, subscription string) ReceiverOption {
	return func(receiver *Receiver) error {
		receiver.config.Entity.Topic = topic
		receiver.config.Entity.Subscription = subscription
		return nil
	}
}

func newReceiver(ns internal.NamespaceWithNewAMQPLinks, options ...ReceiverOption) (*Receiver, error) {
	receiver := &Receiver{
		config: struct {
			ReceiveMode    ReceiveMode
			FullEntityPath string
			Entity         entity
			RetryOptions   struct {
				Times int
				Delay time.Duration
			}
		}{
			ReceiveMode: PeekLock,
		},
		lastPeekedSequenceNumber: 0,
	}

	for _, opt := range options {
		if err := opt(receiver); err != nil {
			return nil, err
		}
	}

	entityPath, err := receiver.config.Entity.String()

	if err != nil {
		return nil, err
	}

	receiver.amqpLinks = ns.NewAMQPLinks(entityPath, func(ctx context.Context, session internal.AMQPSession) (internal.AMQPSenderCloser, internal.AMQPReceiverCloser, error) {
		linkOptions := createLinkOptions(receiver.config.ReceiveMode, entityPath)
		return createReceiverLink(ctx, session, linkOptions)
	})

	// 'nil' settler handles returning an error message for receiveAndDelete links.
	if receiver.config.ReceiveMode == PeekLock {
		receiver.settler = &messageSettler{links: receiver.amqpLinks}
	}

	return receiver, nil
}

// ReceiveOptions are options for the ReceiveMessages function.
type ReceiveOptions struct {
	maxWaitTime                  time.Duration
	maxWaitTimeAfterFirstMessage time.Duration
}

// ReceiveOption represents an option for a `ReceiveMessages`.
// For example, `ReceiveWithMaxWaitTime` will let you configure the
// maxmimum amount of time to wait for messages to arrive.
type ReceiveOption func(options *ReceiveOptions) error

// ReceiveWithMaxWaitTime configures how long to wait for the first
// message in a set of messages to arrive.
// Default: 60 seconds
func ReceiveWithMaxWaitTime(max time.Duration) ReceiveOption {
	return func(options *ReceiveOptions) error {
		options.maxWaitTime = max
		return nil
	}
}

// ReceiveWithMaxTimeAfterFirstMessage confiures how long, after the first
// message arrives, to wait before returning.
// Default: 1 second
func ReceiveWithMaxTimeAfterFirstMessage(max time.Duration) ReceiveOption {
	return func(options *ReceiveOptions) error {
		options.maxWaitTimeAfterFirstMessage = max
		return nil
	}
}

// ReceiveMessages receives a fixed number of messages, up to numMessages.
// There are two timeouts involved in receiving messages:
// 1. An explicit timeout set with `ReceiveWithMaxWaitTime` (default: 60 seconds)
// 2. An implicit timeout (default: 1 second) that starts after the first
//    message has been received. This time can be adjusted with `ReceiveWithMaxTimeAfterFirstMessage`
func (r *Receiver) ReceiveMessages(ctx context.Context, maxMessages int, options ...ReceiveOption) ([]*ReceivedMessage, error) {
	ropts := &ReceiveOptions{
		maxWaitTime:                  time.Minute,
		maxWaitTimeAfterFirstMessage: time.Second,
	}

	for _, opt := range options {
		if err := opt(ropts); err != nil {
			return nil, err
		}
	}

	_, receiver, _, _, err := r.amqpLinks.Get(ctx)

	if err != nil {
		return nil, err
	}

	if err := receiver.IssueCredit(uint32(maxMessages)); err != nil {
		return nil, err
	}

	var messages []*ReceivedMessage

	userCtx := ctx
	ctx, cancel := context.WithTimeout(ctx, ropts.maxWaitTime)
	defer cancel()

	// the last "interesting" error to occur in the sequence. Used only if we end up receiving zero
	// messages (otherwise, the user won't even know to look at the return value).
	var fatalError error

	// receive until we get our first interruption:
	// - context is cancelled because the user cancelled it
	// - context is cancelled because one of our timeouts hit
	// - we get all the messages we asked for
	for {
		amqpMessage, err := receiver.Receive(ctx)

		if err != nil {
			fatalError = err

			if !isCancelled(err) {
				// link is dead. Close our local state so it can be recreated again on the next
				// receiveMessage and return whatever we have (or this error if that's the best we can do)
				if err := r.amqpLinks.Close(ctx, false); err != nil {
					tab.For(ctx).Debug(fmt.Sprintf("failed when closing local links (not fatal): %s", err.Error()))
				}

				if len(messages) == 0 {
					return nil, fatalError
				} else {
					// if they got _some_ messages we don't consider that fatal.
					return messages, nil
				}
			}
			break
		}

		messages = append(messages, newReceivedMessage(ctx, amqpMessage))

		if len(messages) == maxMessages {
			break
		}

		if len(messages) == 1 {
			// once we receive a message we start a secondary
			// timer - this one allows us to guarantee we don't hold
			// onto messages too long.
			go func() {
				select {
				case <-time.After(ropts.maxWaitTimeAfterFirstMessage):
					cancel()
				case <-ctx.Done():
					break
				}
			}()
		}
	}

	if fatalError != nil {
		if isCancelled(fatalError) {
			select {
			case <-userCtx.Done():
				// if _we_ cancelled the context it'll just because the normal timeouts
				// fired, which isn't an error. So we only care if the user did it.
				break
			default:
				fatalError = nil
			}
		} else {
			// link is dead. Close our local state so it can be recreated again on the next
			// receiveMessage and return whatever we have (or this error if that's the best we can do)
			if err := r.amqpLinks.Close(ctx, false); err != nil {
				tab.For(ctx).Debug(fmt.Sprintf("failed when closing local links (not fatal): %s", err.Error()))
			}

			if len(messages) == 0 {
				return nil, fatalError
			} else {
				// if they got _some_ messages we don't consider that fatal.
				return messages, nil
			}
		}
	}

	// Check if there are any outstanding messsages/credits. If there are we want to drain
	// so we can flush the link.
	if len(messages) < maxMessages {
		ctx, cancel = context.WithCancel(context.Background())

		// this drain phase is pretty critical - we don't allow it to be interrupted otherwise the link is
		// in an inconsistent state.
		go func() {
			if err := receiver.DrainCredit(ctx); err != nil {
				tab.For(ctx).Debug(fmt.Sprintf("Draining of credit failed. link will be closed and will re-open on next receive: %s", err.Error()))

				// if the drain fails we just close the link so it'll re-open at the next receive.
				if err := r.amqpLinks.Close(ctx, false); err != nil {
					tab.For(ctx).Debug(fmt.Sprintf("Failed to close links on ReceiveMessages cleanup. Not fatal: %s", err.Error()))
				}
			}
			cancel()
		}()

		// now just read until drain completes
		// That's a gap here where we need to be able to drain _only_ the internally cached messages
		// in the receiver. Filed as https://github.com/Azure/go-amqp/issues/71
		for {
			am, err := receiver.Receive(ctx)

			if isCancelled(err) {
				fatalError = err
				break
			}

			messages = append(messages, newReceivedMessage(ctx, am))
		}
	}

	if len(messages) > 0 {
		return messages, nil
	} else {
		return nil, fatalError
	}
}

// ReceiveDeferredMessages receives messages that were deferred using `Receiver.DeferMessage`.
func (r *Receiver) ReceiveDeferredMessages(ctx context.Context, sequenceNumbers []int64) ([]*ReceivedMessage, error) {
	_, _, mgmt, _, err := r.amqpLinks.Get(ctx)

	amqpMessages, err := mgmt.ReceiveDeferred(ctx, r.config.ReceiveMode, sequenceNumbers)

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

type peekOptions struct {
	fromSequenceNumber *int64
}

type PeekOption func(p *peekOptions) error

// PeekFromSequenceNumber sets the sequence number to start peeking
// messages.
func PeekFromSequenceNumber(sequenceNumber int64) PeekOption {
	return func(p *peekOptions) error {
		p.fromSequenceNumber = &sequenceNumber
		return nil
	}
}

// PeekMessages will peek messages without locking or deleting messages.
// Messages that are peeked do not have lock tokens, so settlement methods
// like CompleteMessage, AbandonMessage, DeferMessage or DeadLetterMessage
// will not work with them.
func (r *Receiver) PeekMessages(ctx context.Context, maxMessageCount int, options ...PeekOption) ([]*ReceivedMessage, error) {

	peekOptions := &peekOptions{}

	for _, opt := range options {
		if err := opt(peekOptions); err != nil {
			return nil, err
		}
	}

	_, _, mgmt, _, err := r.amqpLinks.Get(ctx)

	if err != nil {
		return nil, err
	}

	var sequenceNumber = r.lastPeekedSequenceNumber + 1

	if peekOptions.fromSequenceNumber != nil {
		sequenceNumber = *peekOptions.fromSequenceNumber
	}

	messages, err := mgmt.PeekMessages(ctx, sequenceNumber, int32(maxMessageCount))

	if err != nil {
		return nil, err
	}

	receivedMessages := make([]*ReceivedMessage, len(messages))

	for i := 0; i < len(messages); i++ {
		receivedMessages[i] = newReceivedMessage(ctx, messages[i])
	}

	if len(receivedMessages) > 0 && peekOptions.fromSequenceNumber == nil {
		// only update this if they're doing the implicit iteration as part of the receiver.
		r.lastPeekedSequenceNumber = *receivedMessages[len(receivedMessages)-1].SequenceNumber
	}

	return receivedMessages, nil
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

// receiveMessage receives a single message.
func (r *Receiver) receiveMessage(ctx context.Context, options ...ReceiveOption) (*ReceivedMessage, error) {
	messages, err := r.ReceiveMessages(ctx, 1, options...)

	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages[0], nil
}

// Close permanently closes the receiver.
func (r *Receiver) Close(ctx context.Context) error {
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
// queue or subscription. To receive these messages create a receiver with `Client.NewReceiver()`
// using the `ReceiverWithSubQueue()` option.
func (r *Receiver) DeadLetterMessage(ctx context.Context, message *ReceivedMessage, options ...DeadLetterOption) error {
	return r.settler.DeadLetterMessage(ctx, message, options...)
}

type entity struct {
	Subqueue     SubQueue
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

	if e.Subqueue == SubQueueDeadLetter {
		entityPath += "/$DeadLetterQueue"
	} else if e.Subqueue == SubQueueTransfer {
		entityPath += "/$Transfer/$DeadLetterQueue"
	}

	return entityPath, nil
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

	if mode == ReceiveAndDelete {
		receiveMode = amqp.ModeFirst
	}

	opts := []amqp.LinkOption{
		amqp.LinkSourceAddress(entityPath),
		amqp.LinkReceiverSettle(receiveMode),
		amqp.LinkWithManualCredits(),
		amqp.LinkCredit(defaultLinkRxBuffer),
	}

	if mode == ReceiveAndDelete {
		opts = append(opts, amqp.LinkSenderSettle(amqp.ModeSettled))
	}

	return opts
}

func isCancelled(err error) bool {
	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) ||
		err.Error() == "context canceled" { // go-amqp is returning this when I cancel
		return true
	}

	return false
}
