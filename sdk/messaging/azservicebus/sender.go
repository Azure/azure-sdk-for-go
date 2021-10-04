// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
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

	// SendableMessage are sendable using Sender.SendMessage.
	// Message, MessageBatch implement this interface.
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

// SendMessage sends a message to a queue or topic.
// Message can be a MessageBatch (created using `Sender.CreateMessageBatch`) or
// a Message.
func (s *Sender) SendMessage(ctx context.Context, message SendableMessage) error {
	ctx, span := s.startProducerSpanFromContext(ctx, fmt.Sprintf(spanNameSendMessageFmt, message.messageType()))
	defer span.End()

	sender, _, _, _, err := s.links.Get(ctx)

	if err != nil {
		return err
	}

	return sender.Send(ctx, message.toAMQPMessage())
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
	internal.ApplyComponentInfo(span)
	span.AddAttributes(
		tab.StringAttribute("span.kind", "producer"),
		tab.StringAttribute("message_bus.destination", s.links.Audience()),
	)
	return ctx, span
}
