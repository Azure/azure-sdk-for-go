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
	"io/ioutil"
	"sync"
	"time"

	"github.com/Azure/azure-amqp-common-go/log"
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
		sender     *sender
		receiver   *receiver
		receiverMu sync.Mutex
		senderMu   sync.Mutex
	}

	// QueueManager provides CRUD functionality for Service Bus Queues
	QueueManager struct {
		*EntityManager
	}

	// QueueFeed is a specialized Feed containing QueueEntries
	QueueFeed struct {
		*Feed
		Entries []QueueEntry `xml:"entry"`
	}

	// QueueEntry is a specialized Queue Feed Entry
	QueueEntry struct {
		*Entry
		Content *QueueContent `xml:"content"`
	}

	// QueueContent is a specialized Queue body for an Atom Entry
	QueueContent struct {
		XMLName          xml.Name         `xml:"content"`
		Type             string           `xml:"type,attr"`
		QueueDescription QueueDescription `xml:"QueueDescription"`
	}

	// QueueDescription is the content type for Queue management requests
	QueueDescription struct {
		XMLName xml.Name `xml:"QueueDescription"`
		ReceiveBaseDescription
		SendBaseDescription
		BaseEntityDescription
	}

	// QueueOption represents named options for assisting queue creation
	QueueOption func(queue *QueueDescription) error
)

/*
QueueWithPartitioning ensure the created queue will be a partitioned queue. Partitioned queues offer increased
storage and availability compared to non-partitioned queues with the trade-off of requiring the following to ensure
FIFO message retrieval:

SessionId. If a message has the SessionId property set, then Service Bus uses the SessionId property as the
partition key. This way, all messages that belong to the same session are assigned to the same fragment and handled
by the same message broker. This allows Service Bus to guarantee message ordering as well as the consistency of
session states.

PartitionKey. If a message has the PartitionKey property set but not the SessionId property, then Service Bus uses
the PartitionKey property as the partition key. Use the PartitionKey property to send non-sessionful transactional
messages. The partition key ensures that all messages that are sent within a transaction are handled by the same
messaging broker.

MessageId. If the queue has the RequiresDuplicationDetection property set to true, then the MessageId
property serves as the partition key if the SessionId or a PartitionKey properties are not set. This ensures that
all copies of the same message are handled by the same message broker and, thus, allows Service Bus to detect and
eliminate duplicate messages
*/
func QueueWithPartitioning() QueueOption {
	return func(queue *QueueDescription) error {
		queue.EnablePartitioning = ptrBool(true)
		return nil
	}
}

// QueueWithMaxSizeInMegabytes configures the maximum size of the queue in megabytes (1 * 1024 - 5 * 1024), which is the size of
// the memory allocated for the queue. Default is 1 MB (1 * 1024).
func QueueWithMaxSizeInMegabytes(size int) QueueOption {
	return func(q *QueueDescription) error {
		if size < 1*Megabytes || size > 5*Megabytes {
			return errors.New("QueueWithMaxSizeInMegabytes: must be between 1 * Megabytes and 5 * Megabytes")
		}
		int32Size := int32(size)
		q.MaxSizeInMegabytes = &int32Size
		return nil
	}
}

// QueueWithDuplicateDetection configures the queue to detect duplicates for a given time window. If window
// is not specified, then it uses the default of 10 minutes.
func QueueWithDuplicateDetection(window *time.Duration) QueueOption {
	return func(q *QueueDescription) error {
		q.RequiresDuplicateDetection = ptrBool(true)
		if window != nil {
			q.DuplicateDetectionHistoryTimeWindow = durationTo8601Seconds(window)
		}
		return nil
	}
}

// QueueWithRequiredSessions will ensure the queue requires senders and receivers to have sessionIDs
func QueueWithRequiredSessions() QueueOption {
	return func(q *QueueDescription) error {
		q.RequiresSession = ptrBool(true)
		return nil
	}
}

// QueueWithDeadLetteringOnMessageExpiration will ensure the queue sends expired messages to the dead letter queue
func QueueWithDeadLetteringOnMessageExpiration() QueueOption {
	return func(q *QueueDescription) error {
		q.DeadLetteringOnMessageExpiration = ptrBool(true)
		return nil
	}
}

// QueueWithAutoDeleteOnIdle configures the queue to automatically delete after the specified idle interval. The
// minimum duration is 5 minutes.
func QueueWithAutoDeleteOnIdle(window *time.Duration) QueueOption {
	return func(q *QueueDescription) error {
		if window != nil {
			if window.Minutes() < 5 {
				return errors.New("QueueWithAutoDeleteOnIdle: window must be greater than 5 minutes")
			}
			q.AutoDeleteOnIdle = durationTo8601Seconds(window)
		}
		return nil
	}
}

// QueueWithMessageTimeToLive configures the queue to set a time to live on messages. This is the duration after which
// the message expires, starting from when the message is sent to Service Bus. This is the default value used when
// TimeToLive is not set on a message itself. If nil, defaults to 14 days.
func QueueWithMessageTimeToLive(window *time.Duration) QueueOption {
	return func(q *QueueDescription) error {
		if window == nil {
			duration := time.Duration(14 * 24 * time.Hour)
			window = &duration
		}
		q.DefaultMessageTimeToLive = durationTo8601Seconds(window)
		return nil
	}
}

// QueueWithLockDuration configures the queue to have a duration of a peek-lock; that is, the amount of time that the
// message is locked for other receivers. The maximum value for LockDuration is 5 minutes; the default value is 1
// minute.
func QueueWithLockDuration(window *time.Duration) QueueOption {
	return func(q *QueueDescription) error {
		if window == nil {
			duration := time.Duration(1 * time.Minute)
			window = &duration
		}
		q.LockDuration = durationTo8601Seconds(window)
		return nil
	}
}

// NewQueueManager creates a new QueueManager for a Service Bus Namespace
func (ns *Namespace) NewQueueManager() *QueueManager {
	return &QueueManager{
		EntityManager: NewEntityManager(ns.getHTTPSHostURI(), ns.TokenProvider),
	}
}

// Delete deletes a Service Bus Queue entity by name
func (qm *QueueManager) Delete(ctx context.Context, name string) error {
	span, ctx := qm.startSpanFromContext(ctx, "sb.QueueManager.Delete")
	defer span.Finish()

	_, err := qm.EntityManager.Delete(ctx, "/"+name)
	return err
}

// Put creates or updates a Service Bus Queue
func (qm *QueueManager) Put(ctx context.Context, name string, opts ...QueueOption) (*QueueEntry, error) {
	span, ctx := qm.startSpanFromContext(ctx, "sb.QueueManager.Put")
	defer span.Finish()

	queueDescription := new(QueueDescription)

	for _, opt := range opts {
		if err := opt(queueDescription); err != nil {
			log.For(ctx).Error(err)
			return nil, err
		}
	}

	queueDescription.InstanceMetadataSchema = instanceMetadataSchema
	queueDescription.ServiceBusSchema = serviceBusSchema

	qe := &QueueEntry{
		Entry: &Entry{
			DataServiceSchema:         dataServiceSchema,
			DataServiceMetadataSchema: dataServiceMetadataSchema,
			AtomSchema:                atomSchema,
		},
		Content: &QueueContent{
			Type:             applicationXML,
			QueueDescription: *queueDescription,
		},
	}

	reqBytes, err := xml.Marshal(qe)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	reqBytes = xmlDoc(reqBytes)
	res, err := qm.EntityManager.Put(ctx, "/"+name, reqBytes)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	var entry QueueEntry
	err = xml.Unmarshal(b, &entry)
	return &entry, err
}

// List fetches all of the queues for a Service Bus Namespace
func (qm *QueueManager) List(ctx context.Context) (*QueueFeed, error) {
	span, ctx := qm.startSpanFromContext(ctx, "sb.QueueManager.List")
	defer span.Finish()

	res, err := qm.EntityManager.Get(ctx, `/$Resources/Queues`)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	var feed QueueFeed
	err = xml.Unmarshal(b, &feed)
	return &feed, err
}

// Get fetches a Service Bus Queue entity by name
func (qm *QueueManager) Get(ctx context.Context, name string) (*QueueEntry, error) {
	span, ctx := qm.startSpanFromContext(ctx, "sb.QueueManager.Get")
	defer span.Finish()

	res, err := qm.EntityManager.Get(ctx, name)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	var entry QueueEntry
	err = xml.Unmarshal(b, &entry)
	return &entry, err
}

// NewQueue creates a new Queue Sender / Receiver
func (ns *Namespace) NewQueue(name string) *Queue {
	return &Queue{
		entity: &entity{
			namespace: ns,
			Name:      name,
		},
	}
}

// Send sends messages to the Queue
func (q *Queue) Send(ctx context.Context, event *Event, opts ...SendOption) error {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.Send")
	defer span.Finish()

	err := q.ensureSender(ctx)
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}
	return q.sender.Send(ctx, event, opts...)
}

// Receive subscribes for messages sent to the Queue
func (q *Queue) Receive(ctx context.Context, handler Handler, opts ...ReceiverOptions) (*ListenerHandle, error) {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.Receive")
	defer span.Finish()

	q.receiverMu.Lock()
	defer q.receiverMu.Unlock()

	if q.receiver != nil {
		if err := q.receiver.Close(ctx); err != nil {
			log.For(ctx).Error(err)
			return nil, err
		}
	}

	receiver, err := q.namespace.newReceiver(ctx, q.Name, opts...)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	q.receiver = receiver
	return receiver.Listen(handler), err
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

	if q.sender == nil {
		s, err := q.namespace.newSender(ctx, q.Name)
		if err != nil {
			log.For(ctx).Error(err)
			return err
		}
		q.sender = s
	}
	return nil
}
