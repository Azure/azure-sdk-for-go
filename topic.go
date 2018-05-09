package servicebus

import (
	"context"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"sync"
	"time"
)

type (
	// Topic in contrast to queues, in which each message is processed by a single consumer, topics and subscriptions
	// provide a one-to-many form of communication, in a publish/subscribe pattern. Useful for scaling to very large
	// numbers of recipients, each published message is made available to each subscription registered with the topic.
	// Messages are sent to a topic and delivered to one or more associated subscriptions, depending on filter rules
	// that can be set on a per-subscription basis. The subscriptions can use additional filters to restrict the
	// messages that they want to receive. Messages are sent to a topic in the same way they are sent to a queue,
	// but messages are not received from the topic directly. Instead, they are received from subscriptions. A topic
	// subscription resembles a virtual queue that receives copies of the messages that are sent to the topic.
	// Messages are received from a subscription identically to the way they are received from a queue.
	Topic struct {
		Name      string
		namespace *Namespace
		sender    *sender
		senderMu  sync.Mutex
	}

	// TopicManager provides CRUD functionality for Service Bus Topics
	TopicManager struct {
		*EntityManager
	}

	// TopicFeed is a specialized Feed containing Topic Entries
	TopicFeed struct {
		*Feed
		Entries []TopicEntry `xml:"entry"`
	}
	// TopicEntry is a specialized Topic Feed Entry
	TopicEntry struct {
		*Entry
		Content *TopicContent `xml:"content"`
	}

	// TopicContent is a specialized Topic body for an Atom Entry
	TopicContent struct {
		XMLName          xml.Name         `xml:"content"`
		Type             string           `xml:"type,attr"`
		TopicDescription TopicDescription `xml:"TopicDescription"`
	}

	// TopicDescription is the content type for Topic management requests
	TopicDescription struct {
		XMLName xml.Name `xml:"TopicDescription"`
		SendBaseDescription
		BaseEntityDescription
		FilteringMessagesBeforePublishing *bool `xml:"FilteringMessagesBeforePublishing,omitempty"`
		EnableSubscriptionPartitioning    *bool `xml:"EnableSubscriptionPartitioning,omitempty"`
	}

	// TopicOption represents named options for assisting Topic creation
	TopicOption func(topic *TopicDescription) error
)

// NewTopicManager creates a new TopicManager for a Service Bus Namespace
func (ns *Namespace) NewTopicManager() *TopicManager {
	return &TopicManager{
		EntityManager: NewEntityManager(ns.getHTTPSHostURI(), ns.TokenProvider),
	}
}

// Delete deletes a Service Bus Topic entity by name
func (tm *TopicManager) Delete(ctx context.Context, name string) error {
	_, err := tm.EntityManager.Delete(ctx, "/"+name)
	return err
}

// Put creates or updates a Service Bus Topic
func (tm *TopicManager) Put(ctx context.Context, name string, opts ...TopicOption) (*TopicEntry, error) {
	topicDescription := new(TopicDescription)

	for _, opt := range opts {
		if err := opt(topicDescription); err != nil {
			return nil, err
		}
	}

	topicDescription.InstanceMetadataSchema = instanceMetadataSchema
	topicDescription.ServiceBusSchema = serviceBusSchema

	qe := &TopicEntry{
		Entry: &Entry{
			DataServiceSchema:         dataServiceSchema,
			DataServiceMetadataSchema: dataServiceMetadataSchema,
			AtomSchema:                atomSchema,
		},
		Content: &TopicContent{
			Type:             applicationXML,
			TopicDescription: *topicDescription,
		},
	}

	reqBytes, err := xml.Marshal(qe)
	if err != nil {
		return nil, err
	}

	reqBytes = xmlDoc(reqBytes)
	res, err := tm.EntityManager.Put(ctx, "/"+name, reqBytes)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var entry TopicEntry
	err = xml.Unmarshal(b, &entry)
	return &entry, err
}

// List fetches all of the Topics for a Service Bus Namespace
func (tm *TopicManager) List(ctx context.Context) (*TopicFeed, error) {
	res, err := tm.EntityManager.Get(ctx, `/$Resources/Topics`)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feed TopicFeed
	err = xml.Unmarshal(b, &feed)
	return &feed, err
}

// Get fetches a Service Bus Topic entity by name
func (tm *TopicManager) Get(ctx context.Context, name string) (*TopicEntry, error) {
	res, err := tm.EntityManager.Get(ctx, name)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var entry TopicEntry
	err = xml.Unmarshal(b, &entry)
	return &entry, err
}

// NewTopic creates a new Topic Sender
func (ns *Namespace) NewTopic(name string) *Topic {
	return &Topic{
		namespace: ns,
		Name:      name,
	}
}

// Send sends messages to the Topic
func (t *Topic) Send(ctx context.Context, event *Event, opts ...SendOption) error {
	err := t.ensureSender(ctx)
	if err != nil {
		return err
	}
	return t.sender.Send(ctx, event, opts...)
}

// Close the underlying connection to Service Bus
func (t *Topic) Close(ctx context.Context) error {
	if t.sender != nil {
		return t.sender.Close(ctx)
	}

	return nil
}

func (t *Topic) ensureSender(ctx context.Context) error {
	t.senderMu.Lock()
	defer t.senderMu.Unlock()

	if t.sender == nil {
		s, err := t.namespace.newSender(ctx, t.Name)
		if err != nil {
			return err
		}
		t.sender = s
	}
	return nil
}

// TopicWithMaxSizeInMegabytes configures the maximum size of the topic in megabytes (1 * 1024 - 5 * 1024), which is the size of
// the memory allocated for the topic. Default is 1 MB (1 * 1024).
func TopicWithMaxSizeInMegabytes(size int) TopicOption {
	return func(t *TopicDescription) error {
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
	return func(t *TopicDescription) error {
		t.EnablePartitioning = ptrBool(true)
		return nil
	}
}

// TopicWithOrdering configures the topic to support ordering of messages.
func TopicWithOrdering() TopicOption {
	return func(t *TopicDescription) error {
		t.SupportOrdering = ptrBool(true)
		return nil
	}
}

// TopicWithDuplicateDetection configures the topic to detect duplicates for a given time window. If window
// is not specified, then it uses the default of 10 minutes.
func TopicWithDuplicateDetection(window *time.Duration) TopicOption {
	return func(t *TopicDescription) error {
		t.RequiresDuplicateDetection = ptrBool(true)
		if window != nil {
			t.DuplicateDetectionHistoryTimeWindow = durationTo8601Seconds(window)
		}
		return nil
	}
}

// TopicWithExpress configures the topic to hold a message in memory temporarily before writing it to persistent storage.
func TopicWithExpress() TopicOption {
	return func(t *TopicDescription) error {
		t.EnableExpress = ptrBool(true)
		return nil
	}
}

// TopicWithBatchedOperations configures the topic to batch server-side operations.
func TopicWithBatchedOperations() TopicOption {
	return func(t *TopicDescription) error {
		t.EnableBatchedOperations = ptrBool(true)
		return nil
	}
}

// TopicWithAutoDeleteOnIdle configures the topic to automatically delete after the specified idle interval. The
// minimum duration is 5 minutes.
func TopicWithAutoDeleteOnIdle(window *time.Duration) TopicOption {
	return func(t *TopicDescription) error {
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
	return func(t *TopicDescription) error {
		if window == nil {
			duration := time.Duration(14 * 24 * time.Hour)
			window = &duration
		}
		t.DefaultMessageTimeToLive = durationTo8601Seconds(window)
		return nil
	}
}
