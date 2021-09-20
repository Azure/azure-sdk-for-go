// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/go-amqp"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/devigned/tab"
)

type (
	// SenderOption specifies an option that can configure a Sender.
	SenderOption func(sender *Sender) error

	// Sender is used to send messages as well as schedule them to be delivered at a later date.
	Sender struct {
		queueOrTopic string
		links        internal.AMQPLinks
	}

	// SendableMessage are sendable using Sender.SendMessage.
	// Message, MessageBatch implement this interface.
	SendableMessage interface {
		ToAMQPMessage() *amqp.Message
		MessageType() string
	}
)

// tracing
const (
	spanNameSendMessageFmt string = "sb.sender.SendMessage.%s"
)

type messageBatchOptions struct {
	maxSizeInBytes *int
}

// MessageBatchOption is an option for configuring batch creation in
// `NewMessageBatch`.
type MessageBatchOption func(options *messageBatchOptions) error

// MessageBatchWithMaxSize overrides the max size (in bytes) for a batch.
// By default NewMessageBatch will use the max message size provided by the service.
func MessageBatchWithMaxSize(maxSizeInBytes int) func(options *messageBatchOptions) error {
	return func(options *messageBatchOptions) error {
		options.maxSizeInBytes = &maxSizeInBytes
		return nil
	}
}

// NewMessageBatch can be used to create a batch that contain multiple
// messages. Sending a batch of messages is more efficient than sending the
// messages one at a time.
func (s *Sender) NewMessageBatch(ctx context.Context, options ...MessageBatchOption) (*MessageBatch, error) {
	sender, _, _, _, err := s.links.Get(ctx)

	if err != nil {
		return nil, err
	}

	opts := &messageBatchOptions{
		maxSizeInBytes: to.IntPtr(int(sender.MaxMessageSize())),
	}

	for _, opt := range options {
		if err := opt(opts); err != nil {
			return nil, err
		}
	}

	return &MessageBatch{maxBytes: *opts.maxSizeInBytes}, nil
}

// SendMessage sends a message to a queue or topic.
// Message can be a MessageBatch (created using `Sender.CreateMessageBatch`) or
// a Message.
func (s *Sender) SendMessage(ctx context.Context, message SendableMessage) error {
	ctx, span := s.startProducerSpanFromContext(ctx, fmt.Sprintf(spanNameSendMessageFmt, message.MessageType()))
	defer span.End()

	sender, _, _, _, err := s.links.Get(ctx)

	if err != nil {
		return err
	}

	return sender.Send(ctx, message.ToAMQPMessage())
}

// Close permanently closes the Sender.
func (s *Sender) Close(ctx context.Context) error {
	return s.links.Close(ctx, true)
}

func (sender *Sender) createSenderLink(ctx context.Context, session internal.AMQPSession) (internal.AMQPSenderCloser, internal.AMQPReceiverCloser, error) {
	amqpSender, err := session.NewSender(
		amqp.LinkSenderSettle(amqp.ModeMixed),
		amqp.LinkReceiverSettle(amqp.ModeFirst),
		amqp.LinkTargetAddress(sender.queueOrTopic))

	return amqpSender, nil, err
}

func newSender(ns *internal.Namespace, queueOrTopic string) (*Sender, error) {
	sender := &Sender{
		queueOrTopic: queueOrTopic,
	}

	sender.links = ns.NewAMQPLinks(queueOrTopic, sender.createSenderLink)
	return sender, nil
}

// handleAMQPError is called internally when an event has failed to send so we
// can parse the error to determine whether we should attempt to retry sending the event again.
// func (s *Sender) handleAMQPError(ctx context.Context, err error) error {
// 	var amqpError *amqp.Error
// 	if errors.As(err, &amqpError) {
// 		switch amqpError.Condition {
// 		case errorServerBusy:
// 			return s.retryRetryableAmqpError(ctx, amqpRetryDefaultTimes, amqpRetryBusyServerDelay)
// 		case errorTimeout:
// 			return s.retryRetryableAmqpError(ctx, amqpRetryDefaultTimes, amqpRetryDefaultDelay)
// 		case errorOperationCancelled:
// 			return s.retryRetryableAmqpError(ctx, amqpRetryDefaultTimes, amqpRetryDefaultDelay)
// 		case errorContainerClose:
// 			return s.retryRetryableAmqpError(ctx, amqpRetryDefaultTimes, amqpRetryDefaultDelay)
// 		default:
// 			return err
// 		}
// 	}
// 	return s.retryRetryableAmqpError(ctx, amqpRetryDefaultTimes, amqpRetryDefaultDelay)
// }

func (s *Sender) startProducerSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	internal.ApplyComponentInfo(span)
	span.AddAttributes(
		tab.StringAttribute("span.kind", "producer"),
		tab.StringAttribute("message_bus.destination", s.links.Audience()),
	)
	return ctx, span
}
