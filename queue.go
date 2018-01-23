package servicebus

import (
	"context"
	mgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
)

// QueueOption represents named options for assisting queue creation
type QueueOption func(queue *mgmt.SBQueue) error

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

MessageId. If the queue or topic has the RequiresDuplicationDetection property set to true, then the MessageId
property serves as the partition key if the SessionId or a PartitionKey properties are not set. This ensures that
all copies of the same message are handled by the same message broker and, thus, allows Service Bus to detect and
eliminate duplicate messages
*/
func QueueWithPartitioning() QueueOption {
	return func(queue *mgmt.SBQueue) error {
		queue.SBQueueProperties.EnablePartitioning = ptrBool(true)
		return nil
	}
}

// QueueWithDuplicateDetection will ensure the queue has the ability to detected duplicate messages based on
// the message's MessageID
func QueueWithDuplicateDetection() QueueOption {
	return func(queue *mgmt.SBQueue) error {
		queue.SBQueueProperties.RequiresDuplicateDetection = ptrBool(true)
		return nil
	}
}

// QueueWithRequiredSessions will ensure the queue requires senders and receivers to have sessionIDs
func QueueWithRequiredSessions() QueueOption {
	return func(queue *mgmt.SBQueue) error {
		queue.SBQueueProperties.RequiresSession = ptrBool(true)
		return nil
	}
}

// QueueWithMessageExpiration will ensure the queue sends expired messages to the dead letter queue
func QueueWithMessageExpiration() QueueOption {
	return func(queue *mgmt.SBQueue) error {
		queue.DeadLetteringOnMessageExpiration = ptrBool(true)
		return nil
	}
}

// EnsureQueue makes sure a queue exists in the given namespace. If the queue doesn't exist, it will create it with
// the specified name and properties. If properties are not specified, it will build a default partitioned queue.
func (sb *serviceBus) EnsureQueue(ctx context.Context, queueName string, opts ...QueueOption) (*mgmt.SBQueue, error) {
	log.Debugf("ensuring exists queue %s", queueName)
	queueClient := sb.getQueueMgmtClient()
	queue, err := queueClient.Get(ctx, sb.resourceGroup, sb.namespace, queueName)

	// TODO: check if the queue properties are the same as the requested. If not, throw error or build new queue??
	if err != nil {
		newQueue := &mgmt.SBQueue{
			Name:              &queueName,
			SBQueueProperties: &mgmt.SBQueueProperties{},
		}

		for _, opt := range opts {
			opt(newQueue)
		}

		queue, err = queueClient.CreateOrUpdate(ctx, sb.resourceGroup, sb.namespace, queueName, *newQueue)
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
	client.Authorizer = autorest.NewBearerAuthorizer(sb.token)
	return client
}
