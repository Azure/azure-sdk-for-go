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
	"fmt"
	"sync"

	"github.com/Azure/azure-amqp-common-go/log"
	"github.com/Azure/go-autorest/autorest/date"
)

type (
	entity struct {
		Name      string
		namespace *Namespace
	}

	// Queue represents a Service Bus Queue entity, which offers First In, First Out (FIFO) message delivery to one or
	// more competing consumers. That is, messages are typically expected to be received and processed by the receivers
	// in the order in which they were added to the queue, and each message is received and processed by only one
	// message consumer.
	Queue struct {
		*entity
		sender            *sender
		receiver          *receiver
		receiverMu        sync.Mutex
		senderMu          sync.Mutex
		receiveMode       ReceiveMode
		requiredSessionID *string
	}

	// queueContent is a specialized Queue body for an Atom entry
	queueContent struct {
		XMLName          xml.Name         `xml:"content"`
		Type             string           `xml:"type,attr"`
		QueueDescription QueueDescription `xml:"QueueDescription"`
	}

	// QueueDescription is the content type for Queue management requests
	QueueDescription struct {
		XMLName xml.Name `xml:"QueueDescription"`
		BaseEntityDescription
		LockDuration                        *string       `xml:"LockDuration,omitempty"`               // LockDuration - ISO 8601 timespan duration of a peek-lock; that is, the amount of time that the message is locked for other receivers. The maximum value for LockDuration is 5 minutes; the default value is 1 minute.
		MaxSizeInMegabytes                  *int32        `xml:"MaxSizeInMegabytes,omitempty"`         // MaxSizeInMegabytes - The maximum size of the queue in megabytes, which is the size of memory allocated for the queue. Default is 1024.
		RequiresDuplicateDetection          *bool         `xml:"RequiresDuplicateDetection,omitempty"` // RequiresDuplicateDetection - A value indicating if this queue requires duplicate detection.
		RequiresSession                     *bool         `xml:"RequiresSession,omitempty"`
		DefaultMessageTimeToLive            *string       `xml:"DefaultMessageTimeToLive,omitempty"`            // DefaultMessageTimeToLive - ISO 8601 default message timespan to live value. This is the duration after which the message expires, starting from when the message is sent to Service Bus. This is the default value used when TimeToLive is not set on a message itself.
		DeadLetteringOnMessageExpiration    *bool         `xml:"DeadLetteringOnMessageExpiration,omitempty"`    // DeadLetteringOnMessageExpiration - A value that indicates whether this queue has dead letter support when a message expires.
		DuplicateDetectionHistoryTimeWindow *string       `xml:"DuplicateDetectionHistoryTimeWindow,omitempty"` // DuplicateDetectionHistoryTimeWindow - ISO 8601 timeSpan structure that defines the duration of the duplicate detection history. The default value is 10 minutes.
		MaxDeliveryCount                    *int32        `xml:"MaxDeliveryCount,omitempty"`                    // MaxDeliveryCount - The maximum delivery count. A message is automatically deadlettered after this number of deliveries. default value is 10.
		EnableBatchedOperations             *bool         `xml:"EnableBatchedOperations,omitempty"`             // EnableBatchedOperations - Value that indicates whether server-side batched operations are enabled.
		SizeInBytes                         *int64        `xml:"SizeInBytes,omitempty"`                         // SizeInBytes - The size of the queue, in bytes.
		MessageCount                        *int64        `xml:"MessageCount,omitempty"`                        // MessageCount - The number of messages in the queue.
		IsAnonymousAccessible               *bool         `xml:"IsAnonymousAccessible,omitempty"`
		Status                              *EntityStatus `xml:"Status,omitempty"`
		CreatedAt                           *date.Time    `xml:"CreatedAt,omitempty"`
		UpdatedAt                           *date.Time    `xml:"UpdatedAt,omitempty"`
		SupportOrdering                     *bool         `xml:"SupportOrdering,omitempty"`
		AutoDeleteOnIdle                    *string       `xml:"AutoDeleteOnIdle,omitempty"`
		EnablePartitioning                  *bool         `xml:"EnablePartitioning,omitempty"`
		EnableExpress                       *bool         `xml:"EnableExpress,omitempty"`
		CountDetails                        *CountDetails `xml:"CountDetails,omitempty"`
	}

	// QueueOption represents named options for assisting Queue message handling
	QueueOption func(*Queue) error

	// ReceiveMode represents the behavior when consuming a message from a queue
	ReceiveMode int
)

const (
	// PeekLockMode causes a receiver to peek at a message, lock it so no others can consume and have the queue wait for
	// the DispositionAction
	PeekLockMode ReceiveMode = 0
	// ReceiveAndDeleteMode causes a receiver to pop messages off of the queue without waiting for DispositionAction
	ReceiveAndDeleteMode ReceiveMode = 1
)

// QueueWithReceiveAndDelete configures a queue to pop and delete messages off of the queue upon receiving the message.
// This differs from the default, PeekLock, where PeekLock receives a message, locks it for a period of time, then sends
// a disposition to the broker when the message has been processed.
func QueueWithReceiveAndDelete() QueueOption {
	return func(q *Queue) error {
		q.receiveMode = ReceiveAndDeleteMode
		return nil
	}
}

//// QueueWithRequiredSession configures a queue to use a session
//func QueueWithRequiredSession(sessionID string) QueueOption {
//	return func(q *Queue) error {
//		q.requiredSessionID = &sessionID
//		return nil
//	}
//}

// NewQueue creates a new Queue Sender / Receiver
func (ns *Namespace) NewQueue(name string, opts ...QueueOption) (*Queue, error) {
	queue := &Queue{
		entity: &entity{
			namespace: ns,
			Name:      name,
		},
		receiveMode: PeekLockMode,
	}

	for _, opt := range opts {
		if err := opt(queue); err != nil {
			return nil, err
		}
	}
	return queue, nil
}

// Send sends messages to the Queue
func (q *Queue) Send(ctx context.Context, event *Message) error {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.Send")
	defer span.Finish()

	err := q.ensureSender(ctx)
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}
	return q.sender.Send(ctx, event)
}

// ReceiveOne will listen to receive a single message. ReceiveOne will only wait as long as the context allows.
func (q *Queue) ReceiveOne(ctx context.Context, handler Handler) error {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.ReceiveOne")
	defer span.Finish()

	if err := q.ensureReceiver(ctx); err != nil {
		return err
	}

	return q.receiver.ReceiveOne(ctx, handler)
}

// Receive subscribes for messages sent to the Queue
func (q *Queue) Receive(ctx context.Context, handler Handler) error {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.Receive")
	defer span.Finish()

	err := q.ensureReceiver(ctx)
	if err != nil {
		return err
	}

	handle := q.receiver.Listen(ctx, handler)
	<-handle.Done()
	return handle.Err()
}

func (q *Queue) ensureReceiver(ctx context.Context) error {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.ensureReceiver")
	defer span.Finish()

	q.receiverMu.Lock()
	defer q.receiverMu.Unlock()

	opts := []receiverOption{receiverWithReceiveMode(q.receiveMode)}
	if q.requiredSessionID != nil {
		opts = append(opts, receiverWithSession(*q.requiredSessionID))
	}

	receiver, err := q.namespace.newReceiver(ctx, q.Name, opts...)
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	q.receiver = receiver
	return nil
}

// Close the underlying connection to Service Bus
func (q *Queue) Close(ctx context.Context) error {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.Close")
	defer span.Finish()

	if q.receiver != nil {
		if err := q.receiver.Close(ctx); err != nil {
			_ = q.sender.Close(ctx)
			log.For(ctx).Error(err)
			return err
		}
	}

	if q.sender != nil {
		return q.sender.Close(ctx)
	}

	return nil
}

func (q *Queue) ensureSender(ctx context.Context) error {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.ensureSender")
	defer span.Finish()

	q.senderMu.Lock()
	defer q.senderMu.Unlock()

	var opts []senderOption
	if q.requiredSessionID != nil {
		opts = append(opts, sendWithSession(*q.requiredSessionID))
	}

	if q.sender == nil {
		s, err := q.namespace.newSender(ctx, q.Name, opts...)
		if err != nil {
			log.For(ctx).Error(err)
			return err
		}
		q.sender = s
	}
	return nil
}

func (e *entity) ManagementPath() string {
	return fmt.Sprintf("%s/$management", e.Name)
}
