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
	"net/http"
	"sync"
	"time"

	"github.com/Azure/azure-amqp-common-go/log"
	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/azure-service-bus-go/atom"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/Azure/go-autorest/autorest/to"
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
		sender            *sender
		receiver          *receiver
		receiverMu        sync.Mutex
		senderMu          sync.Mutex
		receiveMode       ReceiveMode
		requiredSessionID *string
	}

	// QueueManager provides CRUD functionality for Service Bus Queues
	QueueManager struct {
		*entityManager
	}

	// QueueEntity is the Azure Service Bus description of a Queue for management activities
	QueueEntity struct {
		*QueueDescription
		Name string
	}

	// queueFeed is a specialized feed containing QueueEntries
	queueFeed struct {
		*atom.Feed
		Entries []queueEntry `xml:"entry"`
	}

	// queueEntry is a specialized Queue feed entry
	queueEntry struct {
		*atom.Entry
		Content *queueContent `xml:"content"`
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
		LockDuration                        *string                  `xml:"LockDuration,omitempty"`               // LockDuration - ISO 8601 timespan duration of a peek-lock; that is, the amount of time that the message is locked for other receivers. The maximum value for LockDuration is 5 minutes; the default value is 1 minute.
		MaxSizeInMegabytes                  *int32                   `xml:"MaxSizeInMegabytes,omitempty"`         // MaxSizeInMegabytes - The maximum size of the queue in megabytes, which is the size of memory allocated for the queue. Default is 1024.
		RequiresDuplicateDetection          *bool                    `xml:"RequiresDuplicateDetection,omitempty"` // RequiresDuplicateDetection - A value indicating if this queue requires duplicate detection.
		RequiresSession                     *bool                    `xml:"RequiresSession,omitempty"`
		DefaultMessageTimeToLive            *string                  `xml:"DefaultMessageTimeToLive,omitempty"`            // DefaultMessageTimeToLive - ISO 8601 default message timespan to live value. This is the duration after which the message expires, starting from when the message is sent to Service Bus. This is the default value used when TimeToLive is not set on a message itself.
		DeadLetteringOnMessageExpiration    *bool                    `xml:"DeadLetteringOnMessageExpiration,omitempty"`    // DeadLetteringOnMessageExpiration - A value that indicates whether this queue has dead letter support when a message expires.
		DuplicateDetectionHistoryTimeWindow *string                  `xml:"DuplicateDetectionHistoryTimeWindow,omitempty"` // DuplicateDetectionHistoryTimeWindow - ISO 8601 timeSpan structure that defines the duration of the duplicate detection history. The default value is 10 minutes.
		MaxDeliveryCount                    *int32                   `xml:"MaxDeliveryCount,omitempty"`                    // MaxDeliveryCount - The maximum delivery count. A message is automatically deadlettered after this number of deliveries. default value is 10.
		EnableBatchedOperations             *bool                    `xml:"EnableBatchedOperations,omitempty"`             // EnableBatchedOperations - Value that indicates whether server-side batched operations are enabled.
		SizeInBytes                         *int64                   `xml:"SizeInBytes,omitempty"`                         // SizeInBytes - The size of the queue, in bytes.
		MessageCount                        *int64                   `xml:"MessageCount,omitempty"`                        // MessageCount - The number of messages in the queue.
		IsAnonymousAccessible               *bool                    `xml:"IsAnonymousAccessible,omitempty"`
		Status                              *servicebus.EntityStatus `xml:"Status,omitempty"`
		CreatedAt                           *date.Time               `xml:"CreatedAt,omitempty"`
		UpdatedAt                           *date.Time               `xml:"UpdatedAt,omitempty"`
		SupportOrdering                     *bool                    `xml:"SupportOrdering,omitempty"`
		AutoDeleteOnIdle                    *string                  `xml:"AutoDeleteOnIdle,omitempty"`
		EnablePartitioning                  *bool                    `xml:"EnablePartitioning,omitempty"`
		EnableExpress                       *bool                    `xml:"EnableExpress,omitempty"`
		CountDetails                        *CountDetails            `xml:"CountDetails,omitempty"`
	}

	// QueueManagementOption represents named configuration options for queue mutation
	QueueManagementOption func(*QueueDescription) error

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

/*
QueueEntityWithPartitioning ensure the created queue will be a partitioned queue. Partitioned queues offer increased
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
func QueueEntityWithPartitioning() QueueManagementOption {
	return func(queue *QueueDescription) error {
		queue.EnablePartitioning = ptrBool(true)
		return nil
	}
}

// QueueEntityWithMaxSizeInMegabytes configures the maximum size of the queue in megabytes (1 * 1024 - 5 * 1024), which is the size of
// the memory allocated for the queue. Default is 1 MB (1 * 1024).
func QueueEntityWithMaxSizeInMegabytes(size int) QueueManagementOption {
	return func(q *QueueDescription) error {
		if size < 1*Megabytes || size > 5*Megabytes {
			return errors.New("QueueEntityWithMaxSizeInMegabytes: must be between 1 * Megabytes and 5 * Megabytes")
		}
		int32Size := int32(size)
		q.MaxSizeInMegabytes = &int32Size
		return nil
	}
}

// QueueEntityWithDuplicateDetection configures the queue to detect duplicates for a given time window. If window
// is not specified, then it uses the default of 10 minutes.
func QueueEntityWithDuplicateDetection(window *time.Duration) QueueManagementOption {
	return func(q *QueueDescription) error {
		q.RequiresDuplicateDetection = ptrBool(true)
		if window != nil {
			q.DuplicateDetectionHistoryTimeWindow = durationTo8601Seconds(window)
		}
		return nil
	}
}

// QueueEntityWithRequiredSessions will ensure the queue requires senders and receivers to have sessionIDs
func QueueEntityWithRequiredSessions() QueueManagementOption {
	return func(q *QueueDescription) error {
		q.RequiresSession = ptrBool(true)
		return nil
	}
}

// QueueEntityWithDeadLetteringOnMessageExpiration will ensure the queue sends expired messages to the dead letter queue
func QueueEntityWithDeadLetteringOnMessageExpiration() QueueManagementOption {
	return func(q *QueueDescription) error {
		q.DeadLetteringOnMessageExpiration = ptrBool(true)
		return nil
	}
}

// QueueEntityWithAutoDeleteOnIdle configures the queue to automatically delete after the specified idle interval. The
// minimum duration is 5 minutes.
func QueueEntityWithAutoDeleteOnIdle(window *time.Duration) QueueManagementOption {
	return func(q *QueueDescription) error {
		if window != nil {
			if window.Minutes() < 5 {
				return errors.New("QueueEntityWithAutoDeleteOnIdle: window must be greater than 5 minutes")
			}
			q.AutoDeleteOnIdle = durationTo8601Seconds(window)
		}
		return nil
	}
}

// QueueEntityWithMessageTimeToLive configures the queue to set a time to live on messages. This is the duration after which
// the message expires, starting from when the message is sent to Service Bus. This is the default value used when
// TimeToLive is not set on a message itself. If nil, defaults to 14 days.
func QueueEntityWithMessageTimeToLive(window *time.Duration) QueueManagementOption {
	return func(q *QueueDescription) error {
		if window == nil {
			duration := time.Duration(14 * 24 * time.Hour)
			window = &duration
		}
		q.DefaultMessageTimeToLive = durationTo8601Seconds(window)
		return nil
	}
}

// QueueEntityWithLockDuration configures the queue to have a duration of a peek-lock; that is, the amount of time that the
// message is locked for other receivers. The maximum value for LockDuration is 5 minutes; the default value is 1
// minute.
func QueueEntityWithLockDuration(window *time.Duration) QueueManagementOption {
	return func(q *QueueDescription) error {
		if window == nil {
			duration := time.Duration(1 * time.Minute)
			window = &duration
		}
		q.LockDuration = durationTo8601Seconds(window)
		return nil
	}
}

// QueueEntityWithMaxDeliveryCount configures the queue to have a maximum number of delivery attempts before
// dead-lettering the message
func QueueEntityWithMaxDeliveryCount(count int32) QueueManagementOption {
	return func(q *QueueDescription) error {
		q.MaxDeliveryCount = &count
		return nil
	}
}

// NewQueueManager creates a new QueueManager for a Service Bus Namespace
func (ns *Namespace) NewQueueManager() *QueueManager {
	return &QueueManager{
		entityManager: newEntityManager(ns.getHTTPSHostURI(), ns.TokenProvider),
	}
}

// Delete deletes a Service Bus Queue entity by name
func (qm *QueueManager) Delete(ctx context.Context, name string) error {
	span, ctx := qm.startSpanFromContext(ctx, "sb.QueueManager.Delete")
	defer span.Finish()

	res, err := qm.entityManager.Delete(ctx, "/"+name)
	if res != nil {
		defer res.Body.Close()
	}

	return err
}

// Put creates or updates a Service Bus Queue
func (qm *QueueManager) Put(ctx context.Context, name string, opts ...QueueManagementOption) (*QueueEntity, error) {
	span, ctx := qm.startSpanFromContext(ctx, "sb.QueueManager.Put")
	defer span.Finish()

	qd := new(QueueDescription)
	for _, opt := range opts {
		if err := opt(qd); err != nil {
			log.For(ctx).Error(err)
			return nil, err
		}
	}

	qd.ServiceBusSchema = to.StringPtr(serviceBusSchema)

	qe := &queueEntry{
		Entry: &atom.Entry{
			AtomSchema: atomSchema,
		},
		Content: &queueContent{
			Type:             applicationXML,
			QueueDescription: *qd,
		},
	}

	reqBytes, err := xml.Marshal(qe)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	reqBytes = xmlDoc(reqBytes)
	res, err := qm.entityManager.Put(ctx, "/"+name, reqBytes)
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	var entry queueEntry
	err = xml.Unmarshal(b, &entry)
	if err != nil {
		return nil, formatManagementError(b)
	}
	return queueEntryToEntity(&entry), nil
}

// List fetches all of the queues for a Service Bus Namespace
func (qm *QueueManager) List(ctx context.Context) ([]*QueueEntity, error) {
	span, ctx := qm.startSpanFromContext(ctx, "sb.QueueManager.List")
	defer span.Finish()

	res, err := qm.entityManager.Get(ctx, `/$Resources/Queues`)
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	var feed queueFeed
	err = xml.Unmarshal(b, &feed)
	if err != nil {
		return nil, formatManagementError(b)
	}

	qd := make([]*QueueEntity, len(feed.Entries))
	for idx, entry := range feed.Entries {
		qd[idx] = queueEntryToEntity(&entry)
	}
	return qd, nil
}

// Get fetches a Service Bus Queue entity by name
func (qm *QueueManager) Get(ctx context.Context, name string) (*QueueEntity, error) {
	span, ctx := qm.startSpanFromContext(ctx, "sb.QueueManager.Get")
	defer span.Finish()

	res, err := qm.entityManager.Get(ctx, name)
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	var entry queueEntry
	err = xml.Unmarshal(b, &entry)
	if err != nil {
		if isEmptyFeed(b) {
			return nil, nil
		}
		return nil, formatManagementError(b)
	}

	return queueEntryToEntity(&entry), nil
}

func queueEntryToEntity(entry *queueEntry) *QueueEntity {
	return &QueueEntity{
		QueueDescription: &entry.Content.QueueDescription,
		Name:             entry.Title,
	}
}

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
func (q *Queue) ReceiveOne(ctx context.Context) (*MessageWithContext, error) {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.ReceiveOne")
	defer span.Finish()

	err := q.ensureReceiver(ctx)
	if err != nil {
		return nil, err
	}

	return q.receiver.ReceiveOne(ctx)
}

// Receive subscribes for messages sent to the Queue
func (q *Queue) Receive(ctx context.Context, handler Handler) (*ListenerHandle, error) {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.Receive")
	defer span.Finish()

	err := q.ensureReceiver(ctx)
	if err != nil {
		return nil, err
	}
	return q.receiver.Listen(handler), nil
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
