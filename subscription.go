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
	"github.com/Azure/go-autorest/autorest/date"
)

type (
	// Subscription represents a Service Bus Subscription entity which are used to receive topic messages. A topic
	// subscription resembles a virtual queue that receives copies of the messages that are sent to the topic.
	//Messages are received from a subscription identically to the way they are received from a queue.
	Subscription struct {
		*entity
		Topic      *Topic
		receiver   *receiver
		receiverMu sync.Mutex
	}

	// SubscriptionManager provides CRUD functionality for Service Bus Subscription
	SubscriptionManager struct {
		*EntityManager
		Topic *Topic
	}

	// SubscriptionEntity is the Azure Service Bus description of a topic Subscription for management activities
	SubscriptionEntity struct {
		*SubscriptionDescription
		Name string
	}

	// subscriptionFeed is a specialized Feed containing Topic Subscriptions
	subscriptionFeed struct {
		*Feed
		Entries []subscriptionEntry `xml:"entry"`
	}

	// subscriptionEntryContent is a specialized Topic Feed Subscription
	subscriptionEntry struct {
		*Entry
		Content *subscriptionContent `xml:"content"`
	}

	// subscriptionContent is a specialized Subscription body for an Atom Entry
	subscriptionContent struct {
		XMLName                 xml.Name                `xml:"content"`
		Type                    string                  `xml:"type,attr"`
		SubscriptionDescription SubscriptionDescription `xml:"SubscriptionDescription"`
	}

	// SubscriptionDescription is the content type for Subscription management requests
	SubscriptionDescription struct {
		XMLName xml.Name `xml:"SubscriptionDescription"`
		ReceiveBaseDescription
		BaseEntityDescription
		DeadLetteringOnFilterEvaluationExceptions *bool     `xml:"DeadLetteringOnFilterEvaluationExceptions,omitempty"`
		AccessedAt                                date.Time `xml:"AccessedAt,omitempty"`
	}

	// SubscriptionOption represents named options for assisting Subscription creation
	SubscriptionOption func(topic *SubscriptionDescription) error
)

// NewSubscriptionManager creates a new SubscriptionManager for a Service Bus Topic
func (t *Topic) NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		EntityManager: NewEntityManager(t.namespace.getHTTPSHostURI(), t.namespace.TokenProvider),
		Topic:         t,
	}
}

// NewSubscriptionManager creates a new SubscriptionManger for a Service Bus Namespace
func (ns *Namespace) NewSubscriptionManager(ctx context.Context, topicName string) (*SubscriptionManager, error) {
	t, err := ns.NewTopic(ctx, topicName)
	if err != nil {
		return nil, err
	}
	return &SubscriptionManager{
		EntityManager: NewEntityManager(t.namespace.getHTTPSHostURI(), t.namespace.TokenProvider),
		Topic:         t,
	}, nil
}

// Delete deletes a Service Bus Topic entity by name
func (sm *SubscriptionManager) Delete(ctx context.Context, name string) error {
	span, ctx := sm.startSpanFromContext(ctx, "sb.SubscriptionManager.Delete")
	defer span.Finish()

	_, err := sm.EntityManager.Delete(ctx, sm.getResourceURI(name))
	return err
}

// Put creates or updates a Service Bus Topic
func (sm *SubscriptionManager) Put(ctx context.Context, name string, opts ...SubscriptionOption) (*SubscriptionEntity, error) {
	span, ctx := sm.startSpanFromContext(ctx, "sb.SubscriptionManager.Put")
	defer span.Finish()

	sd := new(SubscriptionDescription)
	for _, opt := range opts {
		if err := opt(sd); err != nil {
			return nil, err
		}
	}

	sd.InstanceMetadataSchema = instanceMetadataSchema
	sd.ServiceBusSchema = serviceBusSchema

	qe := &subscriptionEntry{
		Entry: &Entry{
			DataServiceSchema:         dataServiceSchema,
			DataServiceMetadataSchema: dataServiceMetadataSchema,
			AtomSchema:                atomSchema,
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
	res, err := sm.EntityManager.Put(ctx, sm.getResourceURI(name), reqBytes)
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
		return nil, err
	}
	return subscriptionEntryToEntity(&entry), nil
}

// List fetches all of the Topics for a Service Bus Namespace
func (sm *SubscriptionManager) List(ctx context.Context) ([]*SubscriptionEntity, error) {
	span, ctx := sm.startSpanFromContext(ctx, "sb.SubscriptionManager.List")
	defer span.Finish()

	res, err := sm.EntityManager.Get(ctx, "/"+sm.Topic.Name+"/subscriptions")
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
		return nil, err
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

	res, err := sm.EntityManager.Get(ctx, sm.getResourceURI(name))
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
		return nil, err
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

// NewSubscription creates a new Topic Subscription client
func (t *Topic) NewSubscription(ctx context.Context, name string, opts ...SubscriptionOption) (*Subscription, error) {
	sm := t.NewSubscriptionManager()
	qe, err := sm.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	if qe == nil {
		_, err := sm.Put(ctx, name, opts...)
		if err != nil {
			return nil, err
		}
	}

	return &Subscription{
		entity: &entity{
			namespace: t.namespace,
			Name:      name,
		},
		Topic: t,
	}, nil
}

// Receive subscribes for messages sent to the Subscription
func (s *Subscription) Receive(ctx context.Context, handler Handler, opts ...ReceiverOptions) (*ListenerHandle, error) {
	span, ctx := s.startSpanFromContext(ctx, "sb.Subscription.Receive")
	defer span.Finish()

	s.receiverMu.Lock()
	defer s.receiverMu.Unlock()

	if s.receiver != nil {
		if err := s.receiver.Close(ctx); err != nil {
			log.For(ctx).Error(err)
			return nil, err
		}
	}

	receiver, err := s.namespace.newReceiver(ctx, s.Topic.Name+"/Subscriptions/"+s.Name)
	for _, opt := range opts {
		if err := opt(receiver); err != nil {
			log.For(ctx).Error(err)
			return nil, err
		}
	}

	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	s.receiver = receiver
	return receiver.Listen(handler), err
}

// Close the underlying connection to Service Bus
func (s *Subscription) Close(ctx context.Context) error {
	if s.receiver != nil {
		return s.receiver.Close(ctx)
	}
	return nil
}

// SubscriptionWithBatchedOperations configures the subscription to batch server-side operations.
func SubscriptionWithBatchedOperations() SubscriptionOption {
	return func(s *SubscriptionDescription) error {
		s.EnableBatchedOperations = ptrBool(true)
		return nil
	}
}

// SubscriptionWithLockDuration configures the subscription to have a duration of a peek-lock; that is, the amount of
// time that the message is locked for other receivers. The maximum value for LockDuration is 5 minutes; the default
// value is 1 minute.
func SubscriptionWithLockDuration(window *time.Duration) SubscriptionOption {
	return func(s *SubscriptionDescription) error {
		if window == nil {
			duration := time.Duration(1 * time.Minute)
			window = &duration
		}
		s.LockDuration = durationTo8601Seconds(window)
		return nil
	}
}

// SubscriptionWithRequiredSessions will ensure the subscription requires senders and receivers to have sessionIDs
func SubscriptionWithRequiredSessions() SubscriptionOption {
	return func(s *SubscriptionDescription) error {
		s.RequiresSession = ptrBool(true)
		return nil
	}
}

// SubscriptionWithDeadLetteringOnMessageExpiration will ensure the Subscription sends expired messages to the dead
// letter queue
func SubscriptionWithDeadLetteringOnMessageExpiration() SubscriptionOption {
	return func(s *SubscriptionDescription) error {
		s.DeadLetteringOnMessageExpiration = ptrBool(true)
		return nil
	}
}

// SubscriptionWithAutoDeleteOnIdle configures the subscription to automatically delete after the specified idle
// interval. The minimum duration is 5 minutes.
func SubscriptionWithAutoDeleteOnIdle(window *time.Duration) SubscriptionOption {
	return func(s *SubscriptionDescription) error {
		if window != nil {
			if window.Minutes() < 5 {
				return errors.New("window must be greater than 5 minutes")
			}
			s.AutoDeleteOnIdle = durationTo8601Seconds(window)
		}
		return nil
	}
}

// SubscriptionWithMessageTimeToLive configures the subscription to set a time to live on messages. This is the duration
// after which the message expires, starting from when the message is sent to Service Bus. This is the default value
// used when TimeToLive is not set on a message itself. If nil, defaults to 14 days.
func SubscriptionWithMessageTimeToLive(window *time.Duration) SubscriptionOption {
	return func(s *SubscriptionDescription) error {
		if window == nil {
			duration := time.Duration(14 * 24 * time.Hour)
			window = &duration
		}
		s.DefaultMessageTimeToLive = durationTo8601Seconds(window)
		return nil
	}
}
