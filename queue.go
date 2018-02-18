package servicebus

import (
	"context"
	"errors"
	"time"

	mgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
)

type (
	// QueueOption represents named options for assisting queue creation
	QueueOption func(queue *mgmt.SBQueue) error
)

/*
QueueWithPartitioning ensure the created queue will be a partitioned queue. Partitioned queues offer increased
storage and availability compared to non-partitioned queues with the trade-off of requiring the following to ensure
FIFO message retreival:

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
	return func(queue *mgmt.SBQueue) error {
		queue.EnablePartitioning = ptrBool(true)
		return nil
	}
}

// QueueWithMaxSizeInMegabytes configures the maximum size of the queue in megabytes (1 * 1024 - 5 * 1024), which is the size of
// the memory allocated for the queue. Default is 1 MB (1 * 1024).
func QueueWithMaxSizeInMegabytes(size int) QueueOption {
	return func(q *mgmt.SBQueue) error {
		if size < 1*Megabytes || size > 5*Megabytes {
			return errors.New("QueueWithMaxSizeInMegabytes: must be between 1 * Megabytes and 5 * Megabytes")
		}
		size32 := int32(size)
		q.MaxSizeInMegabytes = &size32
		return nil
	}
}

// QueueWithDuplicateDetection configures the queue to detect duplicates for a given time window. If window
// is not specified, then it uses the default of 10 minutes.
func QueueWithDuplicateDetection(window *time.Duration) QueueOption {
	return func(q *mgmt.SBQueue) error {
		q.RequiresDuplicateDetection = ptrBool(true)
		if window != nil {
			q.DuplicateDetectionHistoryTimeWindow = durationTo8601Seconds(window)
		}
		return nil
	}
}

// QueueWithRequiredSessions will ensure the queue requires senders and receivers to have sessionIDs
func QueueWithRequiredSessions() QueueOption {
	return func(q *mgmt.SBQueue) error {
		q.RequiresSession = ptrBool(true)
		return nil
	}
}

// QueueWithDeadLetteringOnMessageExpiration will ensure the queue sends expired messages to the dead letter queue
func QueueWithDeadLetteringOnMessageExpiration() QueueOption {
	return func(q *mgmt.SBQueue) error {
		q.DeadLetteringOnMessageExpiration = ptrBool(true)
		return nil
	}
}

// QueueWithAutoDeleteOnIdle configures the queue to automatically delete after the specified idle interval. The
// minimum duration is 5 minutes.
func QueueWithAutoDeleteOnIdle(window *time.Duration) QueueOption {
	return func(q *mgmt.SBQueue) error {
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
	return func(q *mgmt.SBQueue) error {
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
	return func(q *mgmt.SBQueue) error {
		if window == nil {
			duration := time.Duration(1 * time.Minute)
			window = &duration
		}
		q.LockDuration = durationTo8601Seconds(window)
		return nil
	}
}

// EnsureQueue makes sure a queue exists in the given namespace. If the queue doesn't exist, it will create it with
// the specified name and properties. If properties are not specified, it will build a default partitioned queue.
func (sb *serviceBus) EnsureQueue(ctx context.Context, name string, opts ...QueueOption) (*mgmt.SBQueue, error) {
	log.Debugf("ensuring exists queue %s", name)
	queueClient := sb.getQueueMgmtClient()
	queue, err := queueClient.Get(ctx, sb.resourceGroup, sb.namespace, name)

	// TODO: check if the queue properties are the same as the requested. If not, throw error or build new queue??
	if err != nil {
		newQueue := &mgmt.SBQueue{
			Name:              &name,
			SBQueueProperties: &mgmt.SBQueueProperties{},
		}

		for _, opt := range opts {
			err = opt(newQueue)
			if err != nil {
				return nil, err
			}
		}

		queue, err = queueClient.CreateOrUpdate(ctx, sb.resourceGroup, sb.namespace, name, *newQueue)
		if err != nil {
			return nil, err
		}
	}
	return &queue, nil
}

// DeleteQueue deletes an existing queue
func (sb *serviceBus) DeleteQueue(ctx context.Context, queueName string) error {
	queueClient := sb.getQueueMgmtClient()
	_, err := queueClient.Delete(ctx, sb.resourceGroup, sb.namespace, queueName)
	return err
}

func (sb *serviceBus) getQueueMgmtClient() mgmt.QueuesClient {
	client := mgmt.NewQueuesClientWithBaseURI(sb.environment.ResourceManagerEndpoint, sb.subscriptionID)
	client.Authorizer = autorest.NewBearerAuthorizer(sb.armToken)
	return client
}
