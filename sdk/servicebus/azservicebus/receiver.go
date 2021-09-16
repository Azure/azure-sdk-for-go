// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
)

// ReceiveMode represents the lock style to use for a reciever - either
// `PeekLock` or `ReceiveAndDelete`
type ReceiveMode int

const (
	// PeekLock will lock messages as they are received and can be settled
	// using the Receiver or Processor's (Complete|Abandon|DeadLetter|Defer)Message
	// functions.
	PeekLock ReceiveMode = 0
	// ReceiveAndDelete will delete messages as they are received.
	ReceiveAndDelete ReceiveMode = 1
)

// SubQueue allows you to target a subqueue of a queue or subscription.
// Ex: the dead letter queue (SubQueueDeadLetter).
type SubQueue string

const (
	// SubQueueDeadLetter targets the dead letter queue for a queue or subscription.
	SubQueueDeadLetter = "deadLetter"
	// SubQueueTransfer targets the transfer dead letter queue for a queue or subscription.
	SubQueueTransfer = "transferDeadLetter"
)

// Receiver receives messages using pull based functions (ReceiveMessages).
// For push-based receiving via callbacks look at the `Processor` type.
type Receiver struct {
	config struct {
		ReceiveMode    ReceiveMode
		FullEntityPath string
		Entity         entity

		RetryOptions struct {
			Times int
			Delay time.Duration
		}
	}

	mu       sync.Mutex
	ns       legacyNamespace
	receiver internal.LegacyReceiver

	linkState *linkState
}

// ReceiverOption represents an option for a receiver.
// Some examples:
// - `ReceiverWithReceiveMode` to configure the receive mode,
// - `ReceiverWithQueue` to target a queue.
type ReceiverOption func(receiver *Receiver) error

// ReceiverWithSubQueue allows you to open the sub queue (ie: dead letter queues, transfer dead letter queues)
// for a queue or subscription.
func ReceiverWithSubQueue(subQueue SubQueue) ReceiverOption {
	return func(receiver *Receiver) error {
		switch subQueue {
		case SubQueueDeadLetter:
		case SubQueueTransfer:
		case "":
			receiver.config.Entity.Subqueue = subQueue
		default:
			return fmt.Errorf("unknown SubQueue %s", subQueue)
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

func newReceiver(ns legacyNamespace, options ...ReceiverOption) (*Receiver, error) {
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
		ns:        ns,
		linkState: newLinkState(context.Background(), errClosed{link: "receiver"}),
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

	receiver.config.FullEntityPath = entityPath
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
func (r *Receiver) ReceiveMessages(ctx context.Context, numMessages int, options ...ReceiveOption) ([]*ReceivedMessage, error) {
	if r.linkState.Closed() {
		return nil, r.linkState.Err()
	}

	ropts := &ReceiveOptions{
		maxWaitTime:                  time.Minute,
		maxWaitTimeAfterFirstMessage: time.Second,
	}

	for _, opt := range options {
		if err := opt(ropts); err != nil {
			return nil, err
		}
	}

	receiver, err := r.updateReceiver(ctx)

	if err != nil {
		return nil, err
	}

	if err := receiver.IssueCredit(uint32(numMessages)); err != nil {
		return nil, err
	}

	var messages []*ReceivedMessage
	var allMessagesReceivedError = errors.New("all messages received")

	receiveCh, afterFirstMessageFn := startReceiverTimer(ropts.maxWaitTime, ropts.maxWaitTimeAfterFirstMessage)

	listenCtx, cancelListenCtx := context.WithCancel(ctx)
	defer cancelListenCtx()

	listenHandle := r.receiver.Listen(listenCtx, internal.HandlerFunc(func(c context.Context, legacyMessage *internal.Message) error {
		// NOTE: the invocation of this function from the AMQP library is single-threaded so
		// no concurrency sync is required.
		messages = append(messages, convertToReceivedMessage(legacyMessage))

		if len(messages) == numMessages {
			return allMessagesReceivedError
		}

		if len(messages) == 1 {
			afterFirstMessageFn()
		}

		return nil
	}))

	select {
	case <-listenHandle.Done():
		if listenHandle.Err() != allMessagesReceivedError {
			err = listenHandle.Err()
		}
	case <-receiveCh:
		break
	case <-r.linkState.Done():
		err = r.linkState.Err()
	case <-ctx.Done():
		err = ctx.Err()
	}

	// make sure we leave the link in a consistent state.
	if err := r.receiver.DrainCredit(ctx); err != nil {
		return nil, err
	}

	return messages, err
}

// CompleteMessage completes a message, deleting it from the queue or subscription.
func (r *Receiver) CompleteMessage(ctx context.Context, message *ReceivedMessage) error {
	return message.legacyMessage.Complete(ctx)
}

// DeadLetterMessage settles a message by moving it to the dead letter queue for a
// queue or subscription.
func (r *Receiver) DeadLetterMessage(ctx context.Context, message *ReceivedMessage) error {
	// TODO: expand to let them set the reason and description.
	return message.legacyMessage.DeadLetter(ctx, nil)
}

// AbandonMessage will cause a message to be returned to the queue or subscription.
// This will increment its delivery count, and potentially cause it to be dead lettered
// depending on your queue or subscription's configuration.
func (r *Receiver) AbandonMessage(ctx context.Context, message *ReceivedMessage) error {
	return message.legacyMessage.Abandon(ctx)
}

// DeferMessage will cause a message to be deferred.
// Messages that are deferred by can be retrieved using `Receiver.ReceiveDeferredMessages()`.
func (r *Receiver) DeferMessage(ctx context.Context, message *ReceivedMessage) error {
	return message.legacyMessage.Defer(ctx)
}

// Close permanently closes the receiver.
func (r *Receiver) Close(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	defer r.linkState.Close()

	var err error

	if r.receiver != nil {
		err = r.receiver.Close(ctx)
		r.receiver = nil
	}

	return err
}

func (r *Receiver) updateReceiver(ctx context.Context) (internal.LegacyReceiver, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.receiver != nil {
		return r.receiver, nil
	}

	receiver, err := r.ns.NewReceiver(ctx,
		r.config.FullEntityPath,
		internal.ReceiverWithReceiveMode(internal.ReceiveMode(r.config.ReceiveMode)))

	if err != nil {
		return nil, err
	}

	r.receiver = receiver
	return r.receiver, nil
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

func startReceiverTimer(initial time.Duration, timeAfterFirstMessage time.Duration) (<-chan struct{}, func()) {
	ch := make(chan struct{}, 1)

	go func() {
		<-time.After(initial)
		select {
		case ch <- struct{}{}:
		default:
		}
	}()

	afterFirstMessage := func() {
		go func() {
			<-time.After(timeAfterFirstMessage)
			select {
			case ch <- struct{}{}:
			default:
			}
		}()
	}

	return ch, afterFirstMessage
}
