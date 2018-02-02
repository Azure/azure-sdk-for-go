package servicebus

import (
	"context"
	"errors"
	mgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
	"time"
)

type (
	// SubscriptionOption represents an option for configuring a topic.
	SubscriptionOption func(subscription *mgmt.SBSubscription) error
)

// SubscriptionWithBatchedOperations configures the subscription to batch server-side operations.
func SubscriptionWithBatchedOperations() SubscriptionOption {
	return func(t *mgmt.SBSubscription) error {
		t.EnableBatchedOperations = ptrBool(true)
		return nil
	}
}

// SubscriptionWithLockDuration configures the subscription to have a duration of a peek-lock; that is, the amount of
// time that the message is locked for other receivers. The maximum value for LockDuration is 5 minutes; the default
// value is 1 minute.
func SubscriptionWithLockDuration(window *time.Duration) SubscriptionOption {
	return func(q *mgmt.SBSubscription) error {
		if window == nil {
			duration := time.Duration(1 * time.Minute)
			window = &duration
		}
		q.LockDuration = durationTo8601Seconds(window)
		return nil
	}
}

// SubscriptionWithRequiredSessions will ensure the subscription requires senders and receivers to have sessionIDs
func SubscriptionWithRequiredSessions() SubscriptionOption {
	return func(q *mgmt.SBSubscription) error {
		q.RequiresSession = ptrBool(true)
		return nil
	}
}

// SubscriptionWithDeadLetteringOnMessageExpiration will ensure the Subscription sends expired messages to the dead
// letter queue
func SubscriptionWithDeadLetteringOnMessageExpiration() SubscriptionOption {
	return func(q *mgmt.SBSubscription) error {
		q.DeadLetteringOnMessageExpiration = ptrBool(true)
		return nil
	}
}

// SubscriptionWithAutoDeleteOnIdle configures the subscription to automatically delete after the specified idle
// interval. The minimum duration is 5 minutes.
func SubscriptionWithAutoDeleteOnIdle(window *time.Duration) SubscriptionOption {
	return func(q *mgmt.SBSubscription) error {
		if window != nil {
			if window.Minutes() < 5 {
				return errors.New("SubscriptionWithAutoDeleteOnIdle: window must be greater than 5 minutes")
			}
			q.AutoDeleteOnIdle = durationTo8601Seconds(window)
		}
		return nil
	}
}

// SubscriptionWithMessageTimeToLive configures the subscription to set a time to live on messages. This is the duration
// after which the message expires, starting from when the message is sent to Service Bus. This is the default value
// used when TimeToLive is not set on a message itself. If nil, defaults to 14 days.
func SubscriptionWithMessageTimeToLive(window *time.Duration) SubscriptionOption {
	return func(q *mgmt.SBSubscription) error {
		if window == nil {
			duration := time.Duration(14 * 24 * time.Hour)
			window = &duration
		}
		q.DefaultMessageTimeToLive = durationTo8601Seconds(window)
		return nil
	}
}

// EnsureSubscription creates a subscription if the subscription does not exist
func (sb *serviceBus) EnsureSubscription(ctx context.Context, topicName, name string, opts ...SubscriptionOption) (*mgmt.SBSubscription, error) {
	log.Debugf("ensuring subscription %s exists", name)
	subClient := sb.getSubscriptionMgmtClient()
	subscription, err := subClient.Get(ctx, sb.resourceGroup, sb.namespace, topicName, name)

	if err != nil {
		newSub := &mgmt.SBSubscription{
			Name: &name,
			SBSubscriptionProperties: &mgmt.SBSubscriptionProperties{
				EnableBatchedOperations: ptrBool(false),
			},
		}

		for _, opt := range opts {
			err = opt(newSub)
			if err != nil {
				return nil, err
			}
		}

		subscription, err = subClient.CreateOrUpdate(ctx, sb.resourceGroup, sb.namespace, topicName, name, *newSub)
		if err != nil {
			return nil, err
		}
	}
	return &subscription, nil
}

// DeleteSubscription deletes an existing subscription
func (sb *serviceBus) DeleteSubscription(ctx context.Context, topicName, name string) error {
	subscriptionClient := sb.getSubscriptionMgmtClient()
	_, err := subscriptionClient.Delete(ctx, sb.resourceGroup, sb.namespace, topicName, name)
	return err
}

func (sb *serviceBus) getSubscriptionMgmtClient() *mgmt.SubscriptionsClient {
	client := mgmt.NewSubscriptionsClientWithBaseURI(sb.environment.ResourceManagerEndpoint, sb.subscriptionID)
	client.Authorizer = autorest.NewBearerAuthorizer(sb.armToken)
	return &client
}
