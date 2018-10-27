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
	"encoding/xml"
	"sync"

	"github.com/Azure/azure-amqp-common-go/log"
	"github.com/Azure/go-autorest/autorest/date"
)

type (
	// Subscription represents a Service Bus Subscription entity which are used to receive topic messages. A topic
	// subscription resembles a virtual queue that receives copies of the messages that are sent to the topic.
	//Messages are received from a subscription identically to the way they are received from a queue.
	Subscription struct {
		*entity
		Topic             *Topic
		receiver          *receiver
		receiverMu        sync.Mutex
		receiveMode       ReceiveMode
		requiredSessionID *string
	}

	// SubscriptionDescription is the content type for Subscription management requests
	SubscriptionDescription struct {
		XMLName xml.Name `xml:"SubscriptionDescription"`
		BaseEntityDescription
		LockDuration                              *string       `xml:"LockDuration,omitempty"` // LockDuration - ISO 8601 timespan duration of a peek-lock; that is, the amount of time that the message is locked for other receivers. The maximum value for LockDuration is 5 minutes; the default value is 1 minute.
		RequiresSession                           *bool         `xml:"RequiresSession,omitempty"`
		DefaultMessageTimeToLive                  *string       `xml:"DefaultMessageTimeToLive,omitempty"`         // DefaultMessageTimeToLive - ISO 8601 default message timespan to live value. This is the duration after which the message expires, starting from when the message is sent to Service Bus. This is the default value used when TimeToLive is not set on a message itself.
		DeadLetteringOnMessageExpiration          *bool         `xml:"DeadLetteringOnMessageExpiration,omitempty"` // DeadLetteringOnMessageExpiration - A value that indicates whether this queue has dead letter support when a message expires.
		DeadLetteringOnFilterEvaluationExceptions *bool         `xml:"DeadLetteringOnFilterEvaluationExceptions,omitempty"`
		MessageCount                              *int64        `xml:"MessageCount,omitempty"`            // MessageCount - The number of messages in the queue.
		MaxDeliveryCount                          *int32        `xml:"MaxDeliveryCount,omitempty"`        // MaxDeliveryCount - The maximum delivery count. A message is automatically deadlettered after this number of deliveries. default value is 10.
		EnableBatchedOperations                   *bool         `xml:"EnableBatchedOperations,omitempty"` // EnableBatchedOperations - Value that indicates whether server-side batched operations are enabled.
		Status                                    *EntityStatus `xml:"Status,omitempty"`
		CreatedAt                                 *date.Time    `xml:"CreatedAt,omitempty"`
		UpdatedAt                                 *date.Time    `xml:"UpdatedAt,omitempty"`
		AccessedAt                                *date.Time    `xml:"AccessedAt,omitempty"`
		AutoDeleteOnIdle                          *string       `xml:"AutoDeleteOnIdle,omitempty"`
	}

	// SubscriptionOption configures the Subscription Azure Service Bus client
	SubscriptionOption func(*Subscription) error
)

// SubscriptionWithReceiveAndDelete configures a subscription to pop and delete messages off of the queue upon receiving the message.
// This differs from the default, PeekLock, where PeekLock receives a message, locks it for a period of time, then sends
// a disposition to the broker when the message has been processed.
func SubscriptionWithReceiveAndDelete() SubscriptionOption {
	return func(s *Subscription) error {
		s.receiveMode = ReceiveAndDeleteMode
		return nil
	}
}

// NewSubscription creates a new Topic Subscription client
func (t *Topic) NewSubscription(name string, opts ...SubscriptionOption) (*Subscription, error) {
	sub := &Subscription{
		entity: &entity{
			namespace: t.namespace,
			Name:      name,
		},
		Topic: t,
	}

	for i := range opts {
		if err := opts[i](sub); err != nil {
			return nil, err
		}
	}
	return sub, nil
}

// Peek fetches a list of Messages from the Service Bus broker, with-out acquiring a lock or committing to a
// disposition. The messages are delivered as close to sequence order as possible.
//
// The MessageIterator that is returned has the following properties:
// - Messages are fetches from the server in pages. Page size is configurable with PeekOptions.
// - The MessageIterator will always return "false" for Done().
// - When Next() is called, it will return either: a slice of messages and no error, nil with an error related to being
// unable to complete the operation, or an empty slice of messages and an instance of "ErrNoMessages" signifying that
// there are currently no messages in the subscription with a sequence ID larger than previously viewed ones.
func (s *Subscription) Peek(ctx context.Context, options ...PeekOption) (MessageIterator, error) {
	err := s.ensureReceiver(ctx)
	if err != nil {
		return nil, err
	}

	return newPeekIterator(s.entity, s.receiver.connection, options...)
}

// ReceiveOne will listen to receive a single message. ReceiveOne will only wait as long as the context allows.
func (s *Subscription) ReceiveOne(ctx context.Context, handler Handler) error {
	span, ctx := s.startSpanFromContext(ctx, "sb.Subscription.ReceiveOne")
	defer span.Finish()

	if err := s.ensureReceiver(ctx); err != nil {
		return err
	}

	return s.receiver.ReceiveOne(ctx, handler)
}

// Receive subscribes for messages sent to the Subscription
func (s *Subscription) Receive(ctx context.Context, handler Handler) error {
	span, ctx := s.startSpanFromContext(ctx, "sb.Subscription.Receive")
	defer span.Finish()

	if err := s.ensureReceiver(ctx); err != nil {
		return err
	}
	handle := s.receiver.Listen(ctx, handler)
	<-handle.Done()
	return handle.Err()
}

// ReceiveOneSession waits for the lock on a particular session to become available, takes it, then process the session.
func (s *Subscription) ReceiveOneSession(ctx context.Context, sessionID *string, handler SessionHandler) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	span, ctx := s.startSpanFromContext(ctx, "sb.Subscription.ReceiveOneSession")
	defer span.Finish()

	options := make([]receiverOption, 0, 1)
	if sessionID != nil {
		options = append(options, receiverWithSession(sessionID))
	}
	s.requiredSessionID = sessionID
	if err := s.ensureReceiver(ctx, options...); err != nil {
		return err
	}

	ms, err := newMessageSession(s.receiver, s.entity, sessionID)
	if err != nil {
		return err
	}

	err = handler.Start(ms)
	if err != nil {
		return err
	}

	defer handler.End()
	handle := s.receiver.Listen(ctx, handler)

	select {
	case <-handle.Done():
		return handle.Err()
	case <-ms.done:
		return nil
	}
}

func (s *Subscription) ensureReceiver(ctx context.Context, options ...receiverOption) error {
	span, ctx := s.startSpanFromContext(ctx, "sb.Queue.ensureReceiver")
	defer span.Finish()

	s.receiverMu.Lock()
	defer s.receiverMu.Unlock()

	options = append(options, receiverWithReceiveMode(s.receiveMode))

	receiver, err := s.namespace.newReceiver(ctx, s.Topic.Name+"/Subscriptions/"+s.Name, options...)
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	s.receiver = receiver
	return nil
}

// Close the underlying connection to Service Bus
func (s *Subscription) Close(ctx context.Context) error {
	if s.receiver != nil {
		return s.receiver.Close(ctx)
	}
	return nil
}
