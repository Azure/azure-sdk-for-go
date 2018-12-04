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
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-amqp-common-go/log"
	"github.com/Azure/azure-amqp-common-go/rpc"
	"github.com/Azure/azure-amqp-common-go/uuid"
	"github.com/Azure/go-autorest/autorest/date"
	"pack.ag/amqp"
)

type (
	entity struct {
		Name                  string
		namespace             *Namespace
		renewMessageLockMutex sync.Mutex
	}

	// Queue represents a Service Bus Queue entity, which offers First In, First Out (FIFO) message delivery to one or
	// more competing consumers. That is, messages are typically expected to be received and processed by the receivers
	// in the order in which they were added to the queue, and each message is received and processed by only one
	// message consumer.
	Queue struct {
		*entity
		sender            *Sender
		receiver          *Receiver
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
	// PeekLockMode causes a Receiver to peek at a message, lock it so no others can consume and have the queue wait for
	// the DispositionAction
	PeekLockMode ReceiveMode = 0
	// ReceiveAndDeleteMode causes a Receiver to pop messages off of the queue without waiting for DispositionAction
	ReceiveAndDeleteMode ReceiveMode = 1

	// DeadLetterQueueName is the name of the dead letter queue to be appended to the entity path
	DeadLetterQueueName = "$DeadLetterQueue"

	TransferDeadLetterQueueName = "$Transfer/" + DeadLetterQueueName
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

// ScheduleAt will send a batch of messages to a Queue, schedule them to be enqueued, and return the sequence numbers
// that can be used to cancel each message.
func (q *Queue) ScheduleAt(ctx context.Context, enqueueTime time.Time, messages ...*Message) ([]int64, error) {
	if len(messages) <= 0 {
		return nil, errors.New("expected one or more messages")
	}

	transformed := make([]interface{}, 0, len(messages))
	for i := range messages {
		messages[i].ScheduleAt(enqueueTime)

		if messages[i].ID == "" {
			id, err := uuid.NewV4()
			if err != nil {
				return nil, err
			}
			messages[i].ID = id.String()
		}

		rawAmqp, err := messages[i].toMsg()
		if err != nil {
			return nil, err
		}
		encoded, err := rawAmqp.MarshalBinary()
		if err != nil {
			return nil, err
		}

		individualMessage := map[string]interface{}{
			"message-id": messages[i].ID,
			"message":    encoded,
		}
		if messages[i].GroupID != nil {
			individualMessage["session-id"] = *messages[i].GroupID
		}
		if partitionKey := messages[i].SystemProperties.PartitionKey; partitionKey != nil {
			individualMessage["partition-key"] = *partitionKey
		}
		if viaPartitionKey := messages[i].SystemProperties.ViaPartitionKey; viaPartitionKey != nil {
			individualMessage["via-partition-key"] = *viaPartitionKey
		}

		transformed = append(transformed, individualMessage)
	}

	msg := &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			operationFieldName: scheduleMessageOperationID,
		},
		Value: map[string]interface{}{
			"messages": transformed,
		},
	}

	if deadline, ok := ctx.Deadline(); ok {
		msg.ApplicationProperties[serverTimeoutFieldName] = uint(time.Until(deadline) / time.Millisecond)
	}

	err := q.ensureSender(ctx)
	if err != nil {
		return nil, err
	}

	link, err := rpc.NewLink(q.sender.connection, q.ManagementPath())
	if err != nil {
		return nil, err
	}

	resp, err := link.RetryableRPC(ctx, 5, 5*time.Second, msg)
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, ErrAMQP(*resp)
	}

	retval := make([]int64, 0, len(messages))

	if rawVal, ok := resp.Message.Value.(map[string]interface{}); ok {
		const sequenceFieldName = "sequence-numbers"
		if rawArr, ok := rawVal[sequenceFieldName]; ok {
			if arr, ok := rawArr.([]int64); ok {
				for i := range arr {
					retval = append(retval, arr[i])
				}
				return retval, nil
			}
			return nil, newErrIncorrectType(sequenceFieldName, []int64{}, rawArr)
		}
		return nil, ErrMissingField(sequenceFieldName)
	}
	return nil, newErrIncorrectType("value", map[string]interface{}{}, resp.Message.Value)
}

// CancelScheduled allows for removal of messages that have been handed to the Service Bus broker for later delivery,
// but have not yet ben enqueued.
func (q *Queue) CancelScheduled(ctx context.Context, seq ...int64) error {
	msg := &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			operationFieldName: cancelScheduledOperationID,
		},
		Value: map[string]interface{}{
			"sequence-numbers": seq,
		},
	}

	if deadline, ok := ctx.Deadline(); ok {
		msg.ApplicationProperties[serverTimeoutFieldName] = uint(time.Until(deadline) / time.Millisecond)
	}

	err := q.ensureSender(ctx)
	if err != nil {
		return err
	}

	link, err := rpc.NewLink(q.sender.connection, q.ManagementPath())
	if err != nil {
		return err
	}

	resp, err := link.RetryableRPC(ctx, 5, 5*time.Second, msg)
	if err != nil {
		return err
	}

	if resp.Code != 200 {
		return ErrAMQP(*resp)
	}

	return nil
}

// Peek fetches a list of Messages from the Service Bus broker without acquiring a lock or committing to a disposition.
// The messages are delivered as close to sequence order as possible.
//
// The MessageIterator that is returned has the following properties:
// - Messages are fetches from the server in pages. Page size is configurable with PeekOptions.
// - The MessageIterator will always return "false" for Done().
// - When Next() is called, it will return either: a slice of messages and no error, nil with an error related to being
// unable to complete the operation, or an empty slice of messages and an instance of "ErrNoMessages" signifying that
// there are currently no messages in the queue with a sequence ID larger than previously viewed ones.
func (q *Queue) Peek(ctx context.Context, options ...PeekOption) (MessageIterator, error) {
	err := q.ensureReceiver(ctx)
	if err != nil {
		return nil, err
	}

	return newPeekIterator(q.entity, q.receiver.connection, options...)
}

// PeekOne fetches a single Message from the Service Bus broker without acquiring a lock or committing to a disposition.
func (q *Queue) PeekOne(ctx context.Context, options ...PeekOption) (*Message, error) {
	err := q.ensureReceiver(ctx)
	if err != nil {
		return nil, err
	}

	// Adding PeekWithPageSize(1) as the last option assures that either:
	// - creating the iterator will fail because two of the same option will be applied.
	// - PeekWithPageSize(1) will be applied after all others, so we will not wastefully pull down messages destined to
	//   be unread.
	options = append(options, PeekWithPageSize(1))

	it, err := newPeekIterator(q.entity, q.receiver.connection, options...)
	if err != nil {
		return nil, err
	}
	return it.Next(ctx)
}

// ReceiveOne will listen to receive a single message. ReceiveOne will only wait as long as the context allows.
//
// Handler must call a disposition action such as Complete, Abandon, Deadletter on the message. If the messages does not
// have a disposition set, the Queue's DefaultDisposition will be used.
func (q *Queue) ReceiveOne(ctx context.Context, handler Handler) error {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.ReceiveOne")
	defer span.Finish()

	if err := q.ensureReceiver(ctx); err != nil {
		return err
	}

	return q.receiver.ReceiveOne(ctx, handler)
}

// Receive subscribes for messages sent to the Queue. If the messages not within a session, messages will arrive
// unordered.
//
// Handler must call a disposition action such as Complete, Abandon, Deadletter on the message. If the messages does not
// have a disposition set, the Queue's DefaultDisposition will be used.
//
// If the handler returns an error, the receive loop will be terminated.
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

// NewSession will create a new session based receiver and sender for the queue
//
// Microsoft Azure Service Bus sessions enable joint and ordered handling of unbounded sequences of related messages.
// To realize a FIFO guarantee in Service Bus, use Sessions. Service Bus is not prescriptive about the nature of the
// relationship between the messages, and also does not define a particular model for determining where a message
// sequence starts or ends.
func (q *Queue) NewSession(sessionID *string) *QueueSession {
	return NewQueueSession(q, sessionID)
}

// NewReceiver will create a new Receiver for receiving messages off of a queue
func (q *Queue) NewReceiver(ctx context.Context, opts ...ReceiverOption) (*Receiver, error) {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.NewReceiver")
	defer span.Finish()

	opts = append(opts, ReceiverWithReceiveMode(q.receiveMode))
	return q.namespace.NewReceiver(ctx, q.Name, opts...)
}

// NewSender will create a new Sender for sending messages to the queue
func (q *Queue) NewSender(ctx context.Context, opts ...SenderOption) (*Sender, error) {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.NewSender")
	defer span.Finish()

	return q.namespace.NewSender(ctx, q.Name)
}

// NewDeadLetter creates an entity that represents the dead letter sub queue of the queue
//
// Azure Service Bus queues and topic subscriptions provide a secondary sub-queue, called a dead-letter queue
// (DLQ). The dead-letter queue does not need to be explicitly created and cannot be deleted or otherwise managed
// independent of the main entity.
//
// The purpose of the dead-letter queue is to hold messages that cannot be delivered to any receiver, or messages
// that could not be processed. Messages can then be removed from the DLQ and inspected. An application might, with
// help of an operator, correct issues and resubmit the message, log the fact that there was an error, and take
// corrective action.
//
// From an API and protocol perspective, the DLQ is mostly similar to any other queue, except that messages can only
// be submitted via the dead-letter operation of the parent entity. In addition, time-to-live is not observed, and
// you can't dead-letter a message from a DLQ. The dead-letter queue fully supports peek-lock delivery and
// transactional operations.
//
// Note that there is no automatic cleanup of the DLQ. Messages remain in the DLQ until you explicitly retrieve
// them from the DLQ and call Complete() on the dead-letter message.
func (q *Queue) NewDeadLetter() *DeadLetter {
	return NewDeadLetter(q)
}

// NewDeadLetterReceiver builds a receiver for the Queue's dead letter queue
func (q *Queue) NewDeadLetterReceiver(ctx context.Context, opts ...ReceiverOption) (ReceiveOner, error) {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.NewDeadLetterReceiver")
	defer span.Finish()

	deadLetterEntityPath := strings.Join([]string{q.Name, DeadLetterQueueName}, "/")
	return q.namespace.NewReceiver(ctx, deadLetterEntityPath, opts...)
}

// NewTransferDeadLetter creates an entity that represents the transfer dead letter sub queue of the queue
//
// Messages will be sent to the transfer dead-letter queue under the following conditions:
//   - A message passes through more than 3 queues or topics that are chained together.
//   - The destination queue or topic is disabled or deleted.
//   - The destination queue or topic exceeds the maximum entity size.
func (q Queue) NewTransferDeadLetter() *TransferDeadLetter {
	return NewTransferDeadLetter(q)
}

// NewTransferDeadLetterReceiver builds a receiver for the Queue's transfer dead letter queue
//
// Messages will be sent to the transfer dead-letter queue under the following conditions:
//   - A message passes through more than 3 queues or topics that are chained together.
//   - The destination queue or topic is disabled or deleted.
//   - The destination queue or topic exceeds the maximum entity size.
func (q Queue) NewTransferDeadLetterReceiver(ctx context.Context, opts ...ReceiverOption) (ReceiveOner, error) {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.NewTransferDeadLetterReceiver")
	defer span.Finish()

	transferDeadLetterEntityPath := strings.Join([]string{q.Name, TransferDeadLetterQueueName}, "/")
	return q.namespace.NewReceiver(ctx, transferDeadLetterEntityPath, opts...)
}

// Close the underlying connection to Service Bus
func (q *Queue) Close(ctx context.Context) error {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.Close")
	defer span.Finish()

	if q.receiver != nil {
		if err := q.receiver.Close(ctx); err != nil {
			if q.sender != nil {
				if err := q.sender.Close(ctx); err != nil {
					log.For(ctx).Error(err)
				}
			}
			log.For(ctx).Error(err)
			return err
		}
	}

	if q.sender != nil {
		return q.sender.Close(ctx)
	}

	return nil
}

func (q *Queue) newReceiver(ctx context.Context, opts ...ReceiverOption) (*Receiver, error) {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.NewReceiver")
	defer span.Finish()

	opts = append(opts, ReceiverWithReceiveMode(q.receiveMode))
	return q.namespace.NewReceiver(ctx, q.Name, opts...)
}

func (q *Queue) ensureReceiver(ctx context.Context, opts ...ReceiverOption) error {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.ensureReceiver")
	defer span.Finish()

	q.receiverMu.Lock()
	defer q.receiverMu.Unlock()

	if q.receiver != nil {
		return nil
	}

	receiver, err := q.newReceiver(ctx, opts...)
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}

	q.receiver = receiver
	return nil
}

func (q *Queue) ensureSender(ctx context.Context) error {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.ensureSender")
	defer span.Finish()

	q.senderMu.Lock()
	defer q.senderMu.Unlock()

	if q.sender != nil {
		return nil
	}

	s, err := q.NewSender(ctx)
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}
	q.sender = s
	return nil
}

func (e *entity) ManagementPath() string {
	return fmt.Sprintf("%s/$management", e.Name)
}
