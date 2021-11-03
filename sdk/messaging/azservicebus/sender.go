// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/Azure/go-amqp"
	"github.com/devigned/tab"
)

type (
	// Sender is used to send messages as well as schedule them to be delivered at a later date.
	Sender struct {
		queueOrTopic   string
		cleanupOnClose func()
		links          internal.AMQPLinks
	}

	// SendableMessage can be sent using `Sender.SendMessage` or `Sender.SendMessages`.
	// Message implements this interface.
	SendableMessage interface {
		toAMQPMessage() *amqp.Message
		messageType() string
	}
)

// tracing
const (
	spanNameSendMessageFmt string = "sb.sender.SendMessage.%s"
)

// MessageBatchOptions contains options for the `Sender.NewMessageBatch` function.
type MessageBatchOptions struct {
	// MaxSizeInBytes overrides the max size (in bytes) for a batch.
	// By default NewMessageBatch will use the max message size provided by the service.
	MaxSizeInBytes int
}

// NewMessageBatch can be used to create a batch that contain multiple
// messages. Sending a batch of messages is more efficient than sending the
// messages one at a time.
func (s *Sender) NewMessageBatch(ctx context.Context, options *MessageBatchOptions) (*MessageBatch, error) {
	sender, _, _, _, err := s.links.Get(ctx)

	if err != nil {
		return nil, err
	}

	maxBytes := int(sender.MaxMessageSize())

	if options != nil && options.MaxSizeInBytes != 0 {
		maxBytes = options.MaxSizeInBytes
	}

	return &MessageBatch{maxBytes: maxBytes}, nil
}

// SendMessage sends a SendableMessage (Message) to a queue or topic.
func (s *Sender) SendMessage(ctx context.Context, message SendableMessage) error {
	ctx, span := s.startProducerSpanFromContext(ctx, fmt.Sprintf(spanNameSendMessageFmt, message.messageType()))
	defer span.End()

	sender, _, _, _, err := s.links.Get(ctx)

	if err != nil {
		return err
	}

	return sender.Send(ctx, message.toAMQPMessage())
}

// SendMessageBatch sends a MessageBatch to a queue or topic.
// Message batches can be created using `Sender.NewMessageBatch`.
func (s *Sender) SendMessageBatch(ctx context.Context, batch *MessageBatch) error {
	ctx, span := s.startProducerSpanFromContext(ctx, fmt.Sprintf(spanNameSendMessageFmt, "batch"))
	defer span.End()

	sender, _, _, _, err := s.links.Get(ctx)

	if err != nil {
		return err
	}

	return sender.Send(ctx, batch.toAMQPMessage())
}

// SendMessages sends messages to a queue or topic, using a single MessageBatch.
// If the messages cannot fit into a single MessageBatch this function will fail.
func (s *Sender) SendMessages(ctx context.Context, messages []*Message) error {
	batch, err := s.NewMessageBatch(ctx, nil)

	if err != nil {
		return err
	}

	for _, m := range messages {
		added, err := batch.Add(m)

		if err != nil {
			return err
		}

		if !added {
			// to avoid partial failure scenarios we just bail if the messages are too large to fit
			// into a single batch.
			return errors.New("Messages were too big to fit in a single batch. Remove some messages and try again or create your own batch using Sender.NewMessageBatch(), which gives more fine-grained control.")
		}
	}

	return s.SendMessageBatch(ctx, batch)
}

// ScheduleMessage schedules a message to appear on Service Bus Queue/Subscription at a later time.
// Returns the sequence number of the message that was scheduled. If the message hasn't been
// delivered you can cancel using `Receiver.CancelScheduleMessage(s)`
func (s *Sender) ScheduleMessage(ctx context.Context, message SendableMessage, scheduledEnqueueTime time.Time) (int64, error) {
	sequenceNumbers, err := s.ScheduleMessages(ctx, []SendableMessage{message}, scheduledEnqueueTime)

	if err != nil {
		return 0, err
	}

	return sequenceNumbers[0], nil
}

// ScheduleMessages schedules a slice of messages to appear on Service Bus Queue/Subscription at a later time.
// Returns the sequence numbers of the messages that were scheduled.  Messages that haven't been
// delivered can be cancelled using `Receiver.CancelScheduleMessage(s)`
func (s *Sender) ScheduleMessages(ctx context.Context, messages []SendableMessage, scheduledEnqueueTime time.Time) ([]int64, error) {
	_, _, mgmt, _, err := s.links.Get(ctx)

	if err != nil {
		return nil, err
	}

	var amqpMessages []*amqp.Message

	for _, m := range messages {
		amqpMessages = append(amqpMessages, m.toAMQPMessage())
	}

	return mgmt.ScheduleMessages(ctx, scheduledEnqueueTime, amqpMessages...)
}

// CancelScheduledMessage cancels a message that was scheduled.
func (s *Sender) CancelScheduledMessage(ctx context.Context, sequenceNumber int64) error {
	_, _, mgmt, _, err := s.links.Get(ctx)

	if err != nil {
		return err
	}

	return mgmt.CancelScheduled(ctx, sequenceNumber)
}

// CancelScheduledMessages cancels multiple messages that were scheduled.
func (s *Sender) CancelScheduledMessages(ctx context.Context, sequenceNumber []int64) error {
	_, _, mgmt, _, err := s.links.Get(ctx)

	if err != nil {
		return err
	}

	return mgmt.CancelScheduled(ctx, sequenceNumber...)
}

// Close permanently closes the Sender.
func (s *Sender) Close(ctx context.Context) error {
	s.cleanupOnClose()
	return s.links.Close(ctx, true)
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

func newSender(ns internal.NamespaceWithNewAMQPLinks, queueOrTopic string, cleanupOnClose func()) (*Sender, error) {
	sender := &Sender{
		queueOrTopic:   queueOrTopic,
		cleanupOnClose: cleanupOnClose,
	}

	sender.links = ns.NewAMQPLinks(queueOrTopic, sender.createSenderLink)
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
