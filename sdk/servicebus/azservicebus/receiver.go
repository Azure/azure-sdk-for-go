package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
)

type ReceiveMode int

const (
	PeekLock         ReceiveMode = 0
	ReceiveAndDelete ReceiveMode = 1
)

type SubQueue string

const (
	SubQueueDeadLetter = "deadLetter"
	SubQueueTransfer   = "transferDeadLetter"
)

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
}

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

func ReceiverWithQueue(queue string) ReceiverOption {
	return func(receiver *Receiver) error {
		receiver.config.Entity.Queue = queue
		return nil
	}
}

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
		ns: ns,
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

type ReceiveOptions struct {
	maxWaitTime time.Duration
}

type ReceiveOption func(options *ReceiveOptions) error

func ReceiveWithMaxWaitTime(max time.Duration) ReceiveOption {
	return func(options *ReceiveOptions) error {
		options.maxWaitTime = max
		return nil
	}
}

func (r *Receiver) ReceiveMessages(ctx context.Context, numMessages int, options ...ReceiveOption) ([]*ReceivedMessage, error) {
	ropts := &ReceiveOptions{
		maxWaitTime: time.Minute,
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

	if err := receiver.AddCredit(uint32(numMessages)); err != nil {
		return nil, err
	}

	mu := sync.Mutex{}
	var messages []*ReceivedMessage

	var allMessagesReceivedError = errors.New("all messages received")

	listenHandle := r.receiver.Listen(ctx, internal.HandlerFunc(func(c context.Context, legacyMessage *internal.Message) error {
		mu.Lock()
		defer mu.Unlock()
		messages = append(messages, convertToReceivedMessage(legacyMessage))

		if len(messages) == numMessages {
			return allMessagesReceivedError
		}

		return nil
	}))

	select {
	case <-listenHandle.Done():
		if listenHandle.Err() != allMessagesReceivedError {
			err = listenHandle.Err()
		}
	case <-time.After(ropts.maxWaitTime):
	case <-ctx.Done():
		err = ctx.Err()
	}

	// make sure we leave the link in a consistent state.
	if err := r.receiver.Drain(ctx); err != nil {
		return nil, err
	}

	return messages, err
}

func (r *Receiver) CompleteMessage(ctx context.Context, message *ReceivedMessage) error {
	return message.legacyMessage.Complete(ctx)
}

func (r *Receiver) DeadLetterMessage(ctx context.Context, message *ReceivedMessage) error {
	// TODO: expand to let them set the reason and description.
	return message.legacyMessage.DeadLetter(ctx, nil)
}

func (r *Receiver) AbandonMessage(ctx context.Context, message *ReceivedMessage) error {
	return message.legacyMessage.Abandon(ctx)
}

func (r *Receiver) DeferMessage(ctx context.Context, message *ReceivedMessage) error {
	return message.legacyMessage.Defer(ctx)
}

func (r *Receiver) Close(ctx context.Context) error {
	return r.Close(ctx)
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
