package servicebus

import (
	"context"
	"errors"
	"time"

	mgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
)

const (
	// Megabytes is a helper for specifying MaxSizeInMegabytes and equals 1024, thus 5 GB is 5 * Megabytes
	Megabytes = 1024
)

type (
	// TopicOption represents an option for configuring a topic.
	TopicOption func(*mgmt.SBTopic) error
)

// TopicWithMaxSizeInMegabytes configures the maximum size of the topic in megabytes (1 * 1024 - 5 * 1024), which is the size of
// the memory allocated for the topic. Default is 1 MB (1 * 1024).
func TopicWithMaxSizeInMegabytes(size int) TopicOption {
	return func(t *mgmt.SBTopic) error {
		if size < 1*Megabytes || size > 5*Megabytes {
			return errors.New("TopicWithMaxSizeInMegabytes: must be between 1 * Megabytes and 5 * Megabytes")
		}
		size32 := int32(size)
		t.MaxSizeInMegabytes = &size32
		return nil
	}
}

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
			duration := 14 * 24 * time.Hour
			window = &duration
		}
		t.DefaultMessageTimeToLive = durationTo8601Seconds(window)
		return nil
	}
}

// EnsureTopic creates a topic if an existing topic does not exist
func (sb *serviceBus) EnsureTopic(ctx context.Context, name string, opts ...TopicOption) (*mgmt.SBTopic, error) {
	log.Debugf("ensuring topic %s exists", name)
	topicClient := sb.getTopicMgmtClient()
	topic, err := topicClient.Get(ctx, sb.resourceGroup, sb.namespace, name)

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

// DeleteTopic deletes an existing topic
func (sb *serviceBus) DeleteTopic(ctx context.Context, topicName string) error {
	topicClient := sb.getTopicMgmtClient()
	_, err := topicClient.Delete(ctx, sb.resourceGroup, sb.namespace, topicName)
	return err
}

func (sb *serviceBus) getTopicMgmtClient() *mgmt.TopicsClient {
	client := mgmt.NewTopicsClientWithBaseURI(sb.environment.ResourceManagerEndpoint, sb.subscriptionID)
	client.Authorizer = autorest.NewBearerAuthorizer(sb.armToken)
	return &client
}
