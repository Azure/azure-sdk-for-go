// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"

	"github.com/Azure/go-amqp"
	"github.com/devigned/tab"
)

type amqpHandler interface {
	Handle(ctx context.Context, msg *amqp.Message) error
}

// amqpAdapterHandler is a middleware handler that translates amqp messages into servicebus messages
type amqpAdapterHandler struct {
	next     Handler
	receiver *Receiver
}

func newAmqpAdapterHandler(receiver *Receiver, next Handler) *amqpAdapterHandler {
	return &amqpAdapterHandler{
		next:     next,
		receiver: receiver,
	}
}

func (h *amqpAdapterHandler) Handle(ctx context.Context, msg *amqp.Message) error {
	const optName = "sb.amqpHandler.Handle"

	event, err := messageFromAMQPMessage(msg)
	if err != nil {
		_, span := h.receiver.startConsumerSpanFromContext(ctx, optName)
		span.Logger().Error(err)
		h.receiver.lastError = err
		if h.receiver.doneListening != nil {
			h.receiver.doneListening()
		}
		return err
	}

	ctx, span := tab.StartSpanWithRemoteParent(ctx, optName, event)
	defer span.End()

	id := messageID(msg)
	if idStr, ok := id.(string); ok {
		span.AddAttributes(tab.StringAttribute("amqp.message.id", idStr))
	}

	if err := h.next.Handle(ctx, event); err != nil {
		// stop handling messages since the message consumer ran into an unexpected error
		h.receiver.lastError = err
		if h.receiver.doneListening != nil {
			h.receiver.doneListening()
		}
		return err
	}

	// nothing more to be done. The message was settled when it was accepted by the Receiver
	if h.receiver.mode == ReceiveAndDeleteMode {
		return nil
	}

	// nothing more to be done. The Receiver has no default disposition, so the handler is solely responsible for
	// disposition
	if h.receiver.DefaultDisposition == nil {
		return nil
	}

	// default disposition is set, so try to send the disposition. If the message disposition has already been set, the
	// underlying AMQP library will ignore the second disposition respecting the disposition of the handler func.
	if err := h.receiver.DefaultDisposition(ctx); err != nil {
		// if an error is returned by the default disposition, then we must alert the message consumer as we can't
		// be sure the final message disposition.
		tab.For(ctx).Error(err)
		h.receiver.lastError = err
		if h.receiver.doneListening != nil {
			h.receiver.doneListening()
		}
		return nil
	}
	return nil
}
