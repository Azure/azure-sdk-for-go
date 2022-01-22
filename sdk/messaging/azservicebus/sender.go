// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/Azure/go-amqp"
	"github.com/devigned/tab"
)

type (
	// Sender is used to send messages as well as schedule them to be delivered at a later date.
	Sender struct {
		queueOrTopic   string
		cleanupOnClose func()
		links          internal.AMQPLinks
		retryOptions   utils.RetryOptions
	}
)

// tracing
const (
	spanNameSendMessage string = "sb.sender.SendMessage"
	spanNameSendBatch   string = "sb.sender.SendBatch"
)

// MessageBatchOptions contains options for the `Sender.NewMessageBatch` function.
type MessageBatchOptions struct {
	// MaxBytes overrides the max size (in bytes) for a batch.
	// By default NewMessageBatch will use the max message size provided by the service.
	MaxBytes uint64
}

// NewMessageBatch can be used to create a batch that contain multiple
// messages. Sending a batch of messages is more efficient than sending the
// messages one at a time.
func (s *Sender) NewMessageBatch(ctx context.Context, options *MessageBatchOptions) (*MessageBatch, error) {
	var batch *MessageBatch

	err := s.links.Retry(ctx, "send", func(ctx context.Context, lwid *internal.LinksWithID, args *utils.RetryFnArgs) error {
		maxBytes := lwid.Sender.MaxMessageSize()

		if options != nil && options.MaxBytes != 0 {
			maxBytes = options.MaxBytes
		}

		batch = newMessageBatch(maxBytes)
		return nil
	}, s.retryOptions)

	if err != nil {
		return nil, err
	}

	return batch, nil
}

// SendMessage sends a Message to a queue or topic.
func (s *Sender) SendMessage(ctx context.Context, message *Message) error {
	return s.links.Retry(ctx, "SendMessage", func(ctx context.Context, lwid *internal.LinksWithID, args *utils.RetryFnArgs) error {
		ctx, span := s.startProducerSpanFromContext(ctx, spanNameSendMessage)
		defer span.End()

		return lwid.Sender.Send(ctx, message.toAMQPMessage())
	}, utils.RetryOptions(s.retryOptions))
}

// SendMessageBatch sends a MessageBatch to a queue or topic.
// Message batches can be created using `Sender.NewMessageBatch`.
func (s *Sender) SendMessageBatch(ctx context.Context, batch *MessageBatch) error {
	return s.links.Retry(ctx, "SendMessageBatch", func(ctx context.Context, lwid *internal.LinksWithID, args *utils.RetryFnArgs) error {
		ctx, span := s.startProducerSpanFromContext(ctx, spanNameSendBatch)
		defer span.End()

		return lwid.Sender.Send(ctx, batch.toAMQPMessage())
	}, utils.RetryOptions(s.retryOptions))
}

// ScheduleMessages schedules a slice of Messages to appear on Service Bus Queue/Subscription at a later time.
// Returns the sequence numbers of the messages that were scheduled.  Messages that haven't been
// delivered can be cancelled using `Receiver.CancelScheduleMessage(s)`
func (s *Sender) ScheduleMessages(ctx context.Context, messages []*Message, scheduledEnqueueTime time.Time) ([]int64, error) {
	var amqpMessages []*amqp.Message

	for _, m := range messages {
		amqpMessages = append(amqpMessages, m.toAMQPMessage())
	}

	return s.scheduleAMQPMessages(ctx, amqpMessages, scheduledEnqueueTime)
}

// MessageBatch changes

// CancelScheduledMessages cancels multiple messages that were scheduled.
func (s *Sender) CancelScheduledMessages(ctx context.Context, sequenceNumbers []int64) error {
	return s.links.Retry(ctx, "cancelScheduledMessage", func(ctx context.Context, lwv *internal.LinksWithID, args *utils.RetryFnArgs) error {
		return internal.CancelScheduledMessages(ctx, lwv.RPC, sequenceNumbers)
	}, s.retryOptions)
}

// Close permanently closes the Sender.
func (s *Sender) Close(ctx context.Context) error {
	s.cleanupOnClose()
	return s.links.Close(ctx, true)
}

func (s *Sender) scheduleAMQPMessages(ctx context.Context, messages []*amqp.Message, scheduledEnqueueTime time.Time) ([]int64, error) {
	var sequenceNumbers []int64

	err := s.links.Retry(ctx, "scheduleMessages", func(ctx context.Context, lwv *internal.LinksWithID, args *utils.RetryFnArgs) error {
		sn, err := internal.ScheduleMessages(ctx, lwv.RPC, scheduledEnqueueTime, messages)

		if err != nil {
			return err
		}
		sequenceNumbers = sn
		return nil
	}, s.retryOptions)

	return sequenceNumbers, err
}

func (sender *Sender) createSenderLink(ctx context.Context, session internal.AMQPSession) (internal.AMQPSenderCloser, internal.AMQPReceiverCloser, error) {
	amqpSender, err := session.NewSender(
		amqp.LinkSenderSettle(amqp.ModeMixed),
		amqp.LinkReceiverSettle(amqp.ModeFirst),
		amqp.LinkTargetAddress(sender.queueOrTopic))

	if err != nil {
		tab.For(ctx).Error(err)
		return nil, nil, err
	}

	return amqpSender, nil, nil
}

type newSenderArgs struct {
	ns             internal.NamespaceWithNewAMQPLinks
	queueOrTopic   string
	cleanupOnClose func()
}

func newSender(args newSenderArgs, retryOptions RetryOptions) (*Sender, error) {
	sender := &Sender{
		queueOrTopic:   args.queueOrTopic,
		cleanupOnClose: args.cleanupOnClose,
		retryOptions:   utils.RetryOptions(retryOptions),
	}

	sender.links = args.ns.NewAMQPLinks(args.queueOrTopic, sender.createSenderLink)
	return sender, nil
}

func (s *Sender) startProducerSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	tracing.ApplyComponentInfo(span, internal.Version)
	span.AddAttributes(
		tab.StringAttribute("span.kind", "producer"),
		tab.StringAttribute("message_bus.destination", s.links.Audience()),
	)
	return ctx, span
}
