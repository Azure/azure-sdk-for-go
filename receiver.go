package servicebus

//	MIT License
//
//	Copyright (c) Microsoft Corporation. All rights reserved.
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE

import (
	"context"
	"fmt"
	"time"

	"pack.ag/amqp"

	"github.com/Azure/azure-amqp-common-go"
	"github.com/Azure/azure-amqp-common-go/log"
	"github.com/opentracing/opentracing-go"
)

// receiver provides session and link handling for a receiving entity path
type (
	receiver struct {
		namespace         *Namespace
		connection        *amqp.Client
		session           *session
		receiver          *amqp.Receiver
		entityPath        string
		done              func()
		Name              string
		requiredSessionID *string
		lastError         error
		mode              ReceiveMode
		prefetch          uint32
	}

	// receiverOption provides a structure for configuring receivers
	receiverOption func(receiver *receiver) error

	// ListenerHandle provides the ability to close or listen to the close of a Receiver
	ListenerHandle struct {
		r   *receiver
		ctx context.Context
	}
)

// newReceiver creates a new Service Bus message listener given an AMQP client and an entity path
func (ns *Namespace) newReceiver(ctx context.Context, entityPath string, opts ...receiverOption) (*receiver, error) {
	span, ctx := ns.startSpanFromContext(ctx, "sb.Hub.newReceiver")
	defer span.Finish()

	receiver := &receiver{
		namespace:  ns,
		entityPath: entityPath,
		mode:       PeekLockMode,
		prefetch:   1,
	}

	for _, opt := range opts {
		if err := opt(receiver); err != nil {
			return nil, err
		}
	}

	err := receiver.newSessionAndLink(ctx)
	return receiver, err
}

// Close will close the AMQP session and link of the receiver
func (r *receiver) Close(ctx context.Context) error {
	if r.done != nil {
		r.done()
	}

	return r.connection.Close()
}

// Recover will attempt to close the current session and link, then rebuild them
func (r *receiver) Recover(ctx context.Context) error {
	_ = r.Close(ctx) // we expect the receiver is in an error state
	return r.newSessionAndLink(ctx)
}

func (r *receiver) ReceiveOne(ctx context.Context) (*MessageWithContext, error) {
	span, ctx := r.startConsumerSpanFromContext(ctx, "sb.receiver.ReceiveOne")
	defer span.Finish()

	amqpMsg, err := r.listenForMessage(ctx)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	msg, err := messageFromAMQPMessage(amqpMsg)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	return r.messageToMessageWithContext(ctx, msg), nil
}

func (r *receiver) messageToMessageWithContext(ctx context.Context, msg *Message) (*MessageWithContext) {
	const optName = "sb.receiver.amqpEventToMessageWithContext"
	var span opentracing.Span
	wireContext, err := extractWireContext(msg)
	if err == nil {
		span, ctx = r.startConsumerSpanFromWire(ctx, optName, wireContext)
	} else {
		span, ctx = r.startConsumerSpanFromContext(ctx, optName)
	}
	defer span.Finish()

	span.SetTag("amqp.message-id", msg.ID)
	return &MessageWithContext{
		Message: msg,
		Ctx: ctx,
	}
}

// Listen start a listener for messages sent to the entity path
func (r *receiver) Listen(handler Handler) *ListenerHandle {
	ctx, done := context.WithCancel(context.Background())
	r.done = done

	span, ctx := r.startConsumerSpanFromContext(ctx, "sb.receiver.Listen")
	defer span.Finish()

	messages := make(chan *amqp.Message)
	go r.listenForMessages(ctx, messages)
	go r.handleMessages(ctx, messages, handler)

	return &ListenerHandle{
		r:   r,
		ctx: ctx,
	}
}

func (r *receiver) handleMessages(ctx context.Context, messages chan *amqp.Message, handler Handler) {
	span, ctx := r.startConsumerSpanFromContext(ctx, "sb.receiver.handleMessages")
	defer span.Finish()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-messages:
			r.handleMessage(ctx, msg, handler)
		}
	}
}

func (r *receiver) handleMessage(ctx context.Context, msg *amqp.Message, handler Handler) {
	const optName = "sb.receiver.handleMessage"
	event, err := messageFromAMQPMessage(msg)
	if err != nil {
		_, ctx := r.startConsumerSpanFromContext(ctx, optName)
		log.For(ctx).Error(err)
	}
	var span opentracing.Span
	wireContext, err := extractWireContext(event)
	if err == nil {
		span, ctx = r.startConsumerSpanFromWire(ctx, optName, wireContext)
	} else {
		span, ctx = r.startConsumerSpanFromContext(ctx, optName)
	}
	defer span.Finish()

	id := messageID(msg)
	span.SetTag("amqp.message-id", id)

	dispositionAction := handler(ctx, event)

	if r.mode == ReceiveAndDeleteMode {
		return
	}

	if dispositionAction != nil {
		dispositionAction(ctx)
	} else {
		log.For(ctx).Info(fmt.Sprintf("disposition action not provided auto accepted message id %q", id))
		event.Complete()
	}
}

func extractWireContext(reader opentracing.TextMapReader) (opentracing.SpanContext, error) {
	return opentracing.GlobalTracer().Extract(opentracing.TextMap, reader)
}

func (r *receiver) listenForMessages(ctx context.Context, msgChan chan *amqp.Message) {
	span, ctx := r.startConsumerSpanFromContext(ctx, "sb.receiver.listenForMessages")
	defer span.Finish()

	for {
		msg, err := r.listenForMessage(ctx)
		if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
			return
		}

		if err != nil {
			_, retryErr := common.Retry(5, 10*time.Second, func() (interface{}, error) {
				sp, ctx := r.startConsumerSpanFromContext(ctx, "sb.receiver.listenForMessages.tryRecover")
				defer sp.Finish()

				err := r.Recover(ctx)
				if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
					return nil, ctx.Err()
				}

				if err != nil {
					log.For(ctx).Error(err)
					return nil, common.Retryable(err.Error())
				}
				return nil, nil
			})

			if retryErr != nil {
				r.lastError = retryErr
				r.Close(ctx)
				return
			}
			continue
		}
		select {
		case msgChan <- msg:
		case <-ctx.Done():
			return
		}
	}
}

func (r *receiver) listenForMessage(ctx context.Context) (*amqp.Message, error) {
	span, ctx := r.startConsumerSpanFromContext(ctx, "sb.receiver.listenForMessage")
	defer span.Finish()

	msg, err := r.receiver.Receive(ctx)
	if err != nil {
		log.For(ctx).Debug(err.Error())
		return nil, err
	}

	id := messageID(msg)
	span.SetTag("amqp.message-id", id)
	return msg, nil
}

// newSessionAndLink will replace the session and link on the receiver
func (r *receiver) newSessionAndLink(ctx context.Context) error {
	connection, err := r.namespace.newConnection()
	if err != nil {
		return err
	}
	r.connection = connection

	err = r.namespace.negotiateClaim(ctx, connection, r.entityPath)
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	amqpSession, err := connection.NewSession()
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	r.session, err = newSession(amqpSession)
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	receiveMode := amqp.ModeSecond
	sendMode := amqp.ModeUnsettled
	if r.mode == ReceiveAndDeleteMode {
		receiveMode = amqp.ModeFirst
		sendMode = amqp.ModeSettled
	}

	opts := []amqp.LinkOption{
		amqp.LinkSourceAddress(r.entityPath),
		amqp.LinkSenderSettle(sendMode),
		amqp.LinkReceiverSettle(receiveMode),
		amqp.LinkCredit(r.prefetch),
	}

	if r.requiredSessionID != nil {
		opts = append(opts, amqp.LinkSessionFilter(*r.requiredSessionID))
		r.session.SessionID = *r.requiredSessionID
	}

	amqpReceiver, err := amqpSession.NewReceiver(opts...)
	if err != nil {
		return err
	}

	r.receiver = amqpReceiver
	return nil
}

// receiverWithSession configures a receiver to use a session
func receiverWithSession(sessionID string) receiverOption {
	return func(r *receiver) error {
		r.requiredSessionID = &sessionID
		return nil
	}
}

func receiverWithReceiveMode(mode ReceiveMode) receiverOption {
	return func(r *receiver) error {
		r.mode = mode
		return nil
	}
}

func messageID(msg *amqp.Message) interface{} {
	var id interface{} = "null"
	if msg.Properties != nil {
		id = msg.Properties.MessageID
	}
	return id
}

// Close will close the listener
func (lc *ListenerHandle) Close(ctx context.Context) error {
	return lc.r.Close(ctx)
}

// Done will close the channel when the listener has stopped
func (lc *ListenerHandle) Done() <-chan struct{} {
	return lc.ctx.Done()
}

// Err will return the last error encountered
func (lc *ListenerHandle) Err() error {
	if lc.r.lastError != nil {
		return lc.r.lastError
	}
	return lc.ctx.Err()
}
