package servicebus

import (
	"context"
	"errors"
	mgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
	"time"
)

// TopicOption represents an option for configuring a topic.
type TopicOption func(*mgmt.SBTopic) error

// TopicWithPartitioning configures the topic to be partitioned across multiple message brokers.
func TopicWithPartitioning() TopicOption {
	return func(t *mgmt.SBTopic) error {
		t.EnablePartitioning = ptrBool(true)
		return nil
	}
}

// TopicWithOrdering configures the topic to support ordering of messages.
func TopicWithOrdering() TopicOption {
	return func(t *mgmt.SBTopic) error {
		t.SupportOrdering = ptrBool(true)
		return nil
	}
}

// TopicWithDuplicateDetection configures the topic to detect duplicates for a given time window. If window
// is not specified, then it uses the default of 10 minutes.
func TopicWithDuplicateDetection(window *time.Duration) TopicOption {
	return func(t *mgmt.SBTopic) error {
		t.RequiresDuplicateDetection = ptrBool(true)
		if window != nil {
			t.DuplicateDetectionHistoryTimeWindow = durationTo8601Seconds(window)
		}
		return nil
	}
}

// TopicWithExpress configures the topic to hold a message in memory temporarily before writing it to persistent storage.
func TopicWithExpress() TopicOption {
	return func(t *mgmt.SBTopic) error {
		t.EnableExpress = ptrBool(true)
		return nil
	}
}

// TopicWithBatchedOperations configures the topic to batch server-side operations.
func TopicWithBatchedOperations() TopicOption {
	return func(t *mgmt.SBTopic) error {
		t.EnableBatchedOperations = ptrBool(true)
		return nil
	}
}

// TopicWithAutoDeleteOnIdle configures the topic to automatically delete after the specified idle interval. The
// minimum duration is 5 minutes.
func TopicWithAutoDeleteOnIdle(window *time.Duration) TopicOption {
	return func(t *mgmt.SBTopic) error {
		if window != nil {
			if window.Minutes() < 5 {
				return errors.New("TopicWithAutoDeleteOnIdle: window must be greater than 5 minutes")
			}
			t.AutoDeleteOnIdle = durationTo8601Seconds(window)
		}
		return nil
	}
}

// TopicWithMessageTimeToLive configures the topic to set a time to live on messages. This is the duration after which
// the message expires, starting from when the message is sent to Service Bus. This is the default value used when
// TimeToLive is not set on a message itself. If nil, defaults to 14 days.
func TopicWithMessageTimeToLive(window *time.Duration) TopicOption {
	return func(t *mgmt.SBTopic) error {
		if window == nil {
			duration := time.Duration(14 * 24 * time.Hour)
			window = &duration
		}
		t.DefaultMessageTimeToLive = durationTo8601Seconds(window)
		return nil
	}
}

func (sb *serviceBus) EnsureTopic(ctx context.Context, name string, opts ...TopicOption) (*mgmt.SBTopic, error) {
	log.Debugf("ensuring exists topic %s", name)
	topicClient := sb.getTopicMgmtClient()
	topic, err := topicClient.Get(ctx, sb.resourceGroup, sb.namespace, name)

	// TODO: check if the queue properties are the same as the requested. If not, throw error or build new queue??
	if err != nil {
		newTopic := &mgmt.SBTopic{
			Name: &name,
			SBTopicProperties: &mgmt.SBTopicProperties{
				EnablePartitioning:      ptrBool(false),
				EnableBatchedOperations: ptrBool(false),
				EnableExpress:           ptrBool(false),
				SupportOrdering:         ptrBool(false),
			},
		}

		for _, opt := range opts {
			err = opt(newTopic)
			if err != nil {
				return nil, err
			}
		}

		topic, err = topicClient.CreateOrUpdate(ctx, sb.resourceGroup, sb.namespace, name, *newTopic)
		if err != nil {
			return nil, err
		}
	}
	return &topic, nil
}

// DeleteQueue deletes an existing queue
func (sb *serviceBus) DeleteTopic(ctx context.Context, queueName string) error {
	queueClient := sb.getQueueMgmtClient()
	_, err := queueClient.Delete(ctx, sb.resourceGroup, sb.namespace, queueName)
	return err
}

func (sb *serviceBus) getTopicMgmtClient() *mgmt.TopicsClient {
	client := mgmt.NewTopicsClientWithBaseURI(sb.environment.ResourceManagerEndpoint, sb.subscriptionID)
	client.Authorizer = autorest.NewBearerAuthorizer(sb.token)
	return &client
}
