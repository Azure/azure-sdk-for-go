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
		receiveMode       amqp.ReceiverSettleMode
		prefetch          int
	}

	// ReceiverOptions provides a structure for configuring receivers
	ReceiverOptions func(receiver *receiver) error

	// ListenerHandle provides the ability to close or listen to the close of a Receiver
	ListenerHandle struct {
		r   *receiver
		ctx context.Context
	}

	// ReceiveMode represents the behavior when consuming a message from a queue
	ReceiveMode int
)

const (
	// ReceiveAndDeleteMode causes a receiver to pop messages off of the queue without waiting for DispositionAction
	ReceiveAndDeleteMode ReceiveMode = 0
	// PeekLockMode causes a receiver to peek at a message, lock it so no others can consume and have the queue wait for
	// the DispositionAction
	PeekLockMode ReceiveMode = 1
)

// newReceiver creates a new Service Bus message listener given an AMQP client and an entity path
func (ns *Namespace) newReceiver(ctx context.Context, entityPath string, opts ...ReceiverOptions) (*receiver, error) {
	span, ctx := ns.startSpanFromContext(ctx, "sb.Hub.newReceiver")
	defer span.Finish()

	receiver := &receiver{
		namespace:  ns,
		entityPath: entityPath,
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
	event := messageFromAMQPMessage(msg)
	var span opentracing.Span
	wireContext, err := extractWireContext(event)
	if err == nil {
		span, ctx = r.startConsumerSpanFromWire(ctx, "sb.receiver.handleMessage", wireContext)
	} else {
		span, ctx = r.startConsumerSpanFromContext(ctx, "sb.receiver.handleMessage")
	}
	defer span.Finish()

	id := messageID(msg)
	span.SetTag("amqp.message-id", id)

	dispositionAction := handler(ctx, event)
	if dispositionAction != nil {
		dispositionAction(ctx)
	} else {
		log.For(ctx).Info(fmt.Sprintf("disposition action not provided auto accepted message id %q", id))
		event.Accept()
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

	opts := []amqp.LinkOption{
		amqp.LinkSourceAddress(r.entityPath),
		amqp.LinkReceiverSettle(r.receiveMode),
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

// ReceiverWithSession configures a receiver to use a session
func ReceiverWithSession(sessionID string) ReceiverOptions {
	return func(r *receiver) error {
		r.requiredSessionID = &sessionID
		return nil
	}
}

// ReceiverWithReceiveMode configures a receiver to automatically pop messages from the Queue using ReceiveAndDeleteMode
// vs. peeking at a message, locking it and waiting for the receiver to provide a DispositionAction before popping the
// message
func ReceiverWithReceiveMode(mode ReceiveMode) ReceiverOptions {
	return func(r *receiver) error {
		if mode == ReceiveAndDeleteMode {
			r.receiveMode = amqp.ModeFirst
		} else {
			r.receiveMode = amqp.ModeSecond
		}
		return nil
	}
}

// ReceiverWithPrefetch configures a receiver to fetch a maximum number of unacknowledged messages
func ReceiverWithPrefetch(count int) ReceiverOptions {
	return func(r *receiver) error {
		r.prefetch = count
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
