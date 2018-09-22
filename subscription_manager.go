package servicebus

import (
	"context"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Azure/azure-service-bus-go/atom"
	"github.com/Azure/go-autorest/autorest/to"
)

type (
	// SubscriptionManager provides CRUD functionality for Service Bus Subscription
	SubscriptionManager struct {
		*entityManager
		Topic *Topic
	}

	// SubscriptionEntity is the Azure Service Bus description of a topic Subscription for management activities
	SubscriptionEntity struct {
		*SubscriptionDescription
		Name string
	}

	// subscriptionFeed is a specialized feed containing Topic Subscriptions
	subscriptionFeed struct {
		*atom.Feed
		Entries []subscriptionEntry `xml:"entry"`
	}

	// subscriptionEntryContent is a specialized Topic feed Subscription
	subscriptionEntry struct {
		*atom.Entry
		Content *subscriptionContent `xml:"content"`
	}

	// subscriptionContent is a specialized Subscription body for an Atom entry
	subscriptionContent struct {
		XMLName                 xml.Name                `xml:"content"`
		Type                    string                  `xml:"type,attr"`
		SubscriptionDescription SubscriptionDescription `xml:"SubscriptionDescription"`
	}

	// SubscriptionManagementOption represents named options for assisting Subscription creation
	SubscriptionManagementOption func(*SubscriptionDescription) error
)

// NewSubscriptionManager creates a new SubscriptionManager for a Service Bus Topic
func (t *Topic) NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		entityManager: newEntityManager(t.namespace.getHTTPSHostURI(), t.namespace.TokenProvider),
		Topic:         t,
	}
}

// NewSubscriptionManager creates a new SubscriptionManger for a Service Bus Namespace
func (ns *Namespace) NewSubscriptionManager(topicName string) (*SubscriptionManager, error) {
	t, err := ns.NewTopic(topicName)
	if err != nil {
		return nil, err
	}
	return &SubscriptionManager{
		entityManager: newEntityManager(t.namespace.getHTTPSHostURI(), t.namespace.TokenProvider),
		Topic:         t,
	}, nil
}

// Delete deletes a Service Bus Topic entity by name
func (sm *SubscriptionManager) Delete(ctx context.Context, name string) error {
	span, ctx := sm.startSpanFromContext(ctx, "sb.SubscriptionManager.Delete")
	defer span.Finish()

	res, err := sm.entityManager.Delete(ctx, sm.getResourceURI(name))
	if res != nil {
		defer res.Body.Close()
	}

	return err
}

// Put creates or updates a Service Bus Topic
func (sm *SubscriptionManager) Put(ctx context.Context, name string, opts ...SubscriptionManagementOption) (*SubscriptionEntity, error) {
	span, ctx := sm.startSpanFromContext(ctx, "sb.SubscriptionManager.Put")
	defer span.Finish()

	sd := new(SubscriptionDescription)
	for _, opt := range opts {
		if err := opt(sd); err != nil {
			return nil, err
		}
	}

	sd.ServiceBusSchema = to.StringPtr(serviceBusSchema)

	qe := &subscriptionEntry{
		Entry: &atom.Entry{
			AtomSchema: atomSchema,
		},
		Content: &subscriptionContent{
			Type:                    applicationXML,
			SubscriptionDescription: *sd,
		},
	}

	reqBytes, err := xml.Marshal(qe)
	if err != nil {
		return nil, err
	}

	reqBytes = xmlDoc(reqBytes)
	res, err := sm.entityManager.Put(ctx, sm.getResourceURI(name), reqBytes)
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var entry subscriptionEntry
	err = xml.Unmarshal(b, &entry)
	if err != nil {
		return nil, formatManagementError(b)
	}
	return subscriptionEntryToEntity(&entry), nil
}

// List fetches all of the Topics for a Service Bus Namespace
func (sm *SubscriptionManager) List(ctx context.Context) ([]*SubscriptionEntity, error) {
	span, ctx := sm.startSpanFromContext(ctx, "sb.SubscriptionManager.List")
	defer span.Finish()

	res, err := sm.entityManager.Get(ctx, "/"+sm.Topic.Name+"/subscriptions")
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feed subscriptionFeed
	err = xml.Unmarshal(b, &feed)
	if err != nil {
		return nil, formatManagementError(b)
	}

	subs := make([]*SubscriptionEntity, len(feed.Entries))
	for idx, entry := range feed.Entries {
		subs[idx] = subscriptionEntryToEntity(&entry)
	}
	return subs, nil
}

// Get fetches a Service Bus Topic entity by name
func (sm *SubscriptionManager) Get(ctx context.Context, name string) (*SubscriptionEntity, error) {
	span, ctx := sm.startSpanFromContext(ctx, "sb.SubscriptionManager.Get")
	defer span.Finish()

	res, err := sm.entityManager.Get(ctx, sm.getResourceURI(name))
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var entry subscriptionEntry
	err = xml.Unmarshal(b, &entry)
	if err != nil {
		if isEmptyFeed(b) {
			return nil, nil
		}
		return nil, formatManagementError(b)
	}
	return subscriptionEntryToEntity(&entry), nil
}

func subscriptionEntryToEntity(entry *subscriptionEntry) *SubscriptionEntity {
	return &SubscriptionEntity{
		SubscriptionDescription: &entry.Content.SubscriptionDescription,
		Name:                    entry.Title,
	}
}

func (sm *SubscriptionManager) getResourceURI(name string) string {
	return "/" + sm.Topic.Name + "/subscriptions/" + name
}

// SubscriptionWithBatchedOperations configures the subscription to batch server-side operations.
func SubscriptionWithBatchedOperations() SubscriptionManagementOption {
	return func(s *SubscriptionDescription) error {
		s.EnableBatchedOperations = ptrBool(true)
		return nil
	}
}

// SubscriptionWithLockDuration configures the subscription to have a duration of a peek-lock; that is, the amount of
// time that the message is locked for other receivers. The maximum value for LockDuration is 5 minutes; the default
// value is 1 minute.
func SubscriptionWithLockDuration(window *time.Duration) SubscriptionManagementOption {
	return func(s *SubscriptionDescription) error {
		if window == nil {
			duration := time.Duration(1 * time.Minute)
			window = &duration
		}
		s.LockDuration = ptrString(durationTo8601Seconds(*window))
		return nil
	}
}

// SubscriptionWithRequiredSessions will ensure the subscription requires senders and receivers to have sessionIDs
func SubscriptionWithRequiredSessions() SubscriptionManagementOption {
	return func(s *SubscriptionDescription) error {
		s.RequiresSession = ptrBool(true)
		return nil
	}
}

// SubscriptionWithDeadLetteringOnMessageExpiration will ensure the Subscription sends expired messages to the dead
// letter queue
func SubscriptionWithDeadLetteringOnMessageExpiration() SubscriptionManagementOption {
	return func(s *SubscriptionDescription) error {
		s.DeadLetteringOnMessageExpiration = ptrBool(true)
		return nil
	}
}

// SubscriptionWithAutoDeleteOnIdle configures the subscription to automatically delete after the specified idle
// interval. The minimum duration is 5 minutes.
func SubscriptionWithAutoDeleteOnIdle(window *time.Duration) SubscriptionManagementOption {
	return func(s *SubscriptionDescription) error {
		if window != nil {
			if window.Minutes() < 5 {
				return errors.New("window must be greater than 5 minutes")
			}
			s.AutoDeleteOnIdle = ptrString(durationTo8601Seconds(*window))
		}
		return nil
	}
}

// SubscriptionWithMessageTimeToLive configures the subscription to set a time to live on messages. This is the duration
// after which the message expires, starting from when the message is sent to Service Bus. This is the default value
// used when TimeToLive is not set on a message itself. If nil, defaults to 14 days.
func SubscriptionWithMessageTimeToLive(window *time.Duration) SubscriptionManagementOption {
	return func(s *SubscriptionDescription) error {
		if window == nil {
			duration := time.Duration(14 * 24 * time.Hour)
			window = &duration
		}
		s.DefaultMessageTimeToLive = ptrString(durationTo8601Seconds(*window))
		return nil
	}
}
