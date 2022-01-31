// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
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
	// SubQueueDeadLetter targets the dead letter queue for a queue or subscription.
	SubQueueDeadLetter SubQueue = 1
	// SubQueueTransfer targets the transfer dead letter queue for a queue or subscription.
	SubQueueTransfer SubQueue = 2
)

// Receiver receives messages using pull based functions (ReceiveMessages).
type Receiver struct {
	receiveMode ReceiveMode
	entityPath  string

	settler        settler
	retryOptions   utils.RetryOptions
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

	retryOptions utils.RetryOptions
}

const defaultLinkRxBuffer = 2048

func applyReceiverOptions(receiver *Receiver, entity *entity, options *ReceiverOptions) error {

	if options == nil {
		receiver.receiveMode = ReceiveModePeekLock
	} else {
		if err := checkReceiverMode(options.ReceiveMode); err != nil {
			return err
		}

		receiver.receiveMode = options.ReceiveMode

		if err := entity.SetSubQueue(options.SubQueue); err != nil {
			return err
		}

		receiver.retryOptions = options.retryOptions
	}

	entityPath, err := entity.String()

	if err != nil {
		return err
	}

	receiver.entityPath = entityPath
	return nil
}

type newReceiverArgs struct {
	ns             internal.NamespaceWithNewAMQPLinks
	entity         entity
	cleanupOnClose func()
	newLinkFn      func(ctx context.Context, session internal.AMQPSession) (internal.AMQPSenderCloser, internal.AMQPReceiverCloser, error)
}

func newReceiver(args newReceiverArgs, options *ReceiverOptions) (*Receiver, error) {
	receiver := &Receiver{
		lastPeekedSequenceNumber: 0,
		cleanupOnClose:           args.cleanupOnClose,
	}

	if err := applyReceiverOptions(receiver, &args.entity, options); err != nil {
		return nil, err
	}

	newLinkFn := receiver.newReceiverLink

	if args.newLinkFn != nil {
		newLinkFn = args.newLinkFn
	}

	receiver.amqpLinks = args.ns.NewAMQPLinks(receiver.entityPath, newLinkFn)

	// 'nil' settler handles returning an error message for receiveAndDelete links.
	if receiver.receiveMode == ReceiveModePeekLock {
		receiver.settler = newMessageSettler(receiver.amqpLinks, receiver.retryOptions)
	} else {
		receiver.settler = (*messageSettler)(nil)
	}

	return receiver, nil
}

func (r *Receiver) newReceiverLink(ctx context.Context, session internal.AMQPSession) (internal.AMQPSenderCloser, internal.AMQPReceiverCloser, error) {
	linkOptions := createLinkOptions(r.receiveMode, r.entityPath)
	link, err := createReceiverLink(ctx, session, linkOptions)
	return nil, link, err
}

// ReceiveMessagesOptions are options for the ReceiveMessages function.
type ReceiveMessagesOptions struct {
	// For future expansion
}

// ReceiveMessages receives a fixed number of messages, up to numMessages.
// There are two ways to stop receiving messages:
// 1. Cancelling the `ctx` parameter.
// 2. An implicit timeout (default: 1 second) that starts after the first
//    message has been received.
func (r *Receiver) ReceiveMessages(ctx context.Context, maxMessages int, options *ReceiveMessagesOptions) ([]*ReceivedMessage, error) {
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
	var receivedMessages []*ReceivedMessage

	err := r.amqpLinks.Retry(ctx, "receiveDeferredMessage", func(ctx context.Context, lwid *internal.LinksWithID, args *utils.RetryFnArgs) error {
		amqpMessages, err := internal.ReceiveDeferred(ctx, lwid.RPC, r.receiveMode, sequenceNumbers)

		if err != nil {
			return err
		}

		for _, amqpMsg := range amqpMessages {
			receivedMsg := newReceivedMessage(amqpMsg)
			receivedMsg.deferred = true

			receivedMessages = append(receivedMessages, receivedMsg)
		}

		return nil
	}, utils.RetryOptions(r.retryOptions))

	if err != nil {
		return nil, err
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
	var receivedMessages []*ReceivedMessage

	err := r.amqpLinks.Retry(ctx, "peekMessages", func(ctx context.Context, links *internal.LinksWithID, args *utils.RetryFnArgs) error {
		var sequenceNumber = r.lastPeekedSequenceNumber + 1
		updateInternalSequenceNumber := true

		if options != nil && options.FromSequenceNumber != nil {
			sequenceNumber = *options.FromSequenceNumber
			updateInternalSequenceNumber = false
		}

		messages, err := internal.PeekMessages(ctx, links.RPC, sequenceNumber, int32(maxMessageCount))

		if err != nil {
			return err
		}

		receivedMessages = make([]*ReceivedMessage, len(messages))

		for i := 0; i < len(messages); i++ {
			receivedMessages[i] = newReceivedMessage(messages[i])
		}

		if len(receivedMessages) > 0 && updateInternalSequenceNumber {
			// only update this if they're doing the implicit iteration as part of the receiver.
			r.lastPeekedSequenceNumber = *receivedMessages[len(receivedMessages)-1].SequenceNumber
		}

		return nil
	}, r.retryOptions)

	if err != nil {
		return nil, err
	}

	return receivedMessages, nil
}

// RenewLock renews the lock on a message, updating the `LockedUntil` field on `msg`.
func (r *Receiver) RenewMessageLock(ctx context.Context, msg *ReceivedMessage) error {
	return r.amqpLinks.Retry(ctx, "renewMessageLock", func(ctx context.Context, linksWithVersion *internal.LinksWithID, args *utils.RetryFnArgs) error {
		newExpirationTime, err := internal.RenewLocks(ctx, linksWithVersion.RPC, msg.rawAMQPMessage.LinkName(), []amqp.UUID{
			(amqp.UUID)(msg.LockToken),
		})

		if err != nil {
			return err
		}

		msg.LockedUntil = &newExpirationTime[0]
		return nil
	}, r.retryOptions)

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
func (r *Receiver) AbandonMessage(ctx context.Context, message *ReceivedMessage, options *AbandonMessageOptions) error {
	return r.settler.AbandonMessage(ctx, message, options)
}

// DeferMessage will cause a message to be deferred. Deferred messages
// can be received using `Receiver.ReceiveDeferredMessages`.
func (r *Receiver) DeferMessage(ctx context.Context, message *ReceivedMessage, options *DeferMessageOptions) error {
	return r.settler.DeferMessage(ctx, message, options)
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
func (r *Receiver) receiveMessage(ctx context.Context, options *ReceiveMessagesOptions) (*ReceivedMessage, error) {
	messages, err := r.ReceiveMessages(ctx, 1, options)

	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages[0], nil
}

func (r *Receiver) receiveMessagesImpl(ctx context.Context, maxMessages int, options *ReceiveMessagesOptions) ([]*ReceivedMessage, error) {
	// There are three phases for this function:
	// Phase 1. <receive, respecting user cancellation>
	// Phase 2. <check error and exit if fatal>
	//    NOTE: We don't exit here so we don't end up buffering messages internally that the
	//    user isn't actually waiting for anymore. So we make sure that #3 runs if the
	//    link is still valid.
	// Phase 3. <drain the link and leave it in a good state>

	var all []*ReceivedMessage
	var cancel context.CancelFunc
	needsDrain := true

	defaultTimeAfterFirstMessage := 20 * time.Millisecond

	if r.receiveMode == ReceiveModeReceiveAndDelete {
		defaultTimeAfterFirstMessage = time.Second
	}

	var linksWithID *internal.LinksWithID

	err := r.amqpLinks.Retry(ctx, "receiveMessages", func(ctx context.Context, lwid *internal.LinksWithID, args *utils.RetryFnArgs) error {
		if args.LastErr != nil {
			// we recovered succesfully (amqpLinks does it for us), we can reset our retry attempts.
			// This fixes a potential problem where something like this happens:
			// a. amqplink.Receive() returns an error
			// b. we attempt recovery a few times and recover just in the nick of time at attempt #3.
			// d. amqplink.Receive() blocks (no messages in queue, for instance)
			//   <years pass>
			// e. amqpLink.Receive() returns an error
			//
			// We don't want to the # of attempts to already be '3' when we get to
			// step 5, so we reset here so the next set of retry attempts starts fresh.
			args.ResetAttempts()
		}

		linksWithID = lwid
		credits := maxMessages - len(all)

		if err := lwid.Receiver.IssueCredit(uint32(credits)); err != nil {
			return err
		}

		got := 0

		for {
			amqpMessage, err := lwid.Receiver.Receive(ctx)

			if err != nil {
				return err
			}

			all = append(all, newReceivedMessage(amqpMessage))
			got++

			if got == credits {
				// no excess credits on link
				needsDrain = false
				break
			}

			if cancel == nil {
				// replace the context that we're using for everything with a new one that will cancel
				// after a period of time.
				ctx, cancel = context.WithTimeout(ctx, defaultTimeAfterFirstMessage)
				defer cancel()
			}
		}

		return nil
	}, utils.RetryOptions(r.retryOptions))

	ret := func(err error) ([]*ReceivedMessage, error) {
		if len(all) > 0 {
			// we don't return the error here because we did retrieve _some_ messages and you can still
			// use them.
			return all, nil
		} else {
			return nil, err
		}
	}

	if err != nil && !internal.IsCancelError(err) {
		return ret(err)
	}

	if needsDrain && linksWithID != nil {
		// start the drain asynchronously. Note that we ignore the user's context at this point
		// since draining makes sure we don't get messages when nobody is receiving.
		if err := linksWithID.Receiver.DrainCredit(context.Background()); err != nil {
			if err := r.amqpLinks.RecoverIfNeeded(context.Background(), linksWithID.ID, err); err != nil {
				log.Writef(internal.EventReceiver, "Failed to recover links after a failed drain: %s", err.Error())
				return ret(err)
			}

			return ret(err)
		}

		// Draining data from the receiver's prefetched queue. This won't wait for new messages to
		// arrive, so it'll only receive messages that arrived prior to the drain.
		for {
			am, err := linksWithID.Receiver.Prefetched(context.Background())

			if am == nil || internal.IsCancelError(err) {
				return all, nil
			}

			if err != nil {
				return ret(err)
			}

			all = append(all, newReceivedMessage(am))
		}
	}

	return all, nil
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
	if subQueue == 0 {
		return nil
	} else if subQueue == SubQueueDeadLetter || subQueue == SubQueueTransfer {
		e.subqueue = subQueue
		return nil
	}

	return fmt.Errorf("unknown SubQueue %d", subQueue)
}

func createReceiverLink(ctx context.Context, session internal.AMQPSession, linkOptions []amqp.LinkOption) (internal.AMQPReceiverCloser, error) {
	amqpReceiver, err := session.NewReceiver(linkOptions...)

	if err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	return amqpReceiver, nil
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

func checkReceiverMode(receiveMode ReceiveMode) error {
	if receiveMode == ReceiveModePeekLock || receiveMode == ReceiveModeReceiveAndDelete {
		return nil
	}

	return fmt.Errorf("invalid receive mode %d, must be either azservicebus.PeekLock or azservicebus.ReceiveAndDelete", receiveMode)
}
