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
		*entity
		sender   *sender
		senderMu sync.Mutex
	}

	// TopicManager provides CRUD functionality for Service Bus Topics
	TopicManager struct {
		*EntityManager
	}

	// TopicEntity is the Azure Service Bus description of a Topic for management activities
	TopicEntity struct {
		*TopicDescription
		Name string
	}

	// topicFeed is a specialized Feed containing Topic Entries
	topicFeed struct {
		*Feed
		Entries []topicEntry `xml:"entry"`
	}

	// topicEntry is a specialized Topic Feed Entry
	topicEntry struct {
		*Entry
		Content *topicContent `xml:"content"`
	}

	// topicContent is a specialized Topic body for an Atom Entry
	topicContent struct {
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
	span, ctx := tm.startSpanFromContext(ctx, "sb.TopicManager.Delete")
	defer span.Finish()

	_, err := tm.EntityManager.Delete(ctx, "/"+name)
	return err
}

// Put creates or updates a Service Bus Topic
func (tm *TopicManager) Put(ctx context.Context, name string, opts ...TopicOption) (*TopicEntity, error) {
	span, ctx := tm.startSpanFromContext(ctx, "sb.TopicManager.Put")
	defer span.Finish()

	td := new(TopicDescription)
	for _, opt := range opts {
		if err := opt(td); err != nil {
			log.For(ctx).Error(err)
			return nil, err
		}
	}

	td.InstanceMetadataSchema = instanceMetadataSchema
	td.ServiceBusSchema = serviceBusSchema

	qe := &topicEntry{
		Entry: &Entry{
			DataServiceSchema:         dataServiceSchema,
			DataServiceMetadataSchema: dataServiceMetadataSchema,
			AtomSchema:                atomSchema,
		},
		Content: &topicContent{
			Type:             applicationXML,
			TopicDescription: *td,
		},
	}

	reqBytes, err := xml.Marshal(qe)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	reqBytes = xmlDoc(reqBytes)
	res, err := tm.EntityManager.Put(ctx, "/"+name, reqBytes)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	var entry topicEntry
	err = xml.Unmarshal(b, &entry)
	if err != nil {
		return nil, err
	}
	return topicEntryToEntity(&entry), nil
}

// List fetches all of the Topics for a Service Bus Namespace
func (tm *TopicManager) List(ctx context.Context) ([]*TopicEntity, error) {
	span, ctx := tm.startSpanFromContext(ctx, "sb.TopicManager.List")
	defer span.Finish()

	res, err := tm.EntityManager.Get(ctx, `/$Resources/Topics`)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	var feed topicFeed
	err = xml.Unmarshal(b, &feed)
	if err != nil {
		return nil, err
	}

	topics := make([]*TopicEntity, len(feed.Entries))
	for idx, entry := range feed.Entries {
		topics[idx] = topicEntryToEntity(&entry)
	}
	return topics, nil
}

// Get fetches a Service Bus Topic entity by name
func (tm *TopicManager) Get(ctx context.Context, name string) (*TopicEntity, error) {
	span, ctx := tm.startSpanFromContext(ctx, "sb.TopicManager.Get")
	defer span.Finish()

	res, err := tm.EntityManager.Get(ctx, name)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.For(ctx).Error(err)
		return nil, err
	}

	var entry topicEntry
	err = xml.Unmarshal(b, &entry)
	if err != nil {
		return nil, err
	}
	return topicEntryToEntity(&entry), nil
}

func topicEntryToEntity(entry *topicEntry) *TopicEntity {
	return &TopicEntity{
		TopicDescription: &entry.Content.TopicDescription,
		Name:             entry.Title,
	}
}

// NewTopic creates a new Topic Sender
func (ns *Namespace) NewTopic(name string) *Topic {
	return &Topic{
		entity: &entity{
			namespace: ns,
			Name:      name,
		},
	}
}

// Send sends messages to the Topic
func (t *Topic) Send(ctx context.Context, event *Message, opts ...SendOption) error {
	span, ctx := t.startSpanFromContext(ctx, "sb.Topic.Send")
	defer span.Finish()

	err := t.ensureSender(ctx)
	if err != nil {
		log.For(ctx).Error(err)
		return err
	}
	return t.sender.Send(ctx, event, opts...)
}

// Close the underlying connection to Service Bus
func (t *Topic) Close(ctx context.Context) error {
	span, ctx := t.startSpanFromContext(ctx, "sb.Topic.Close")
	defer span.Finish()

	if t.sender != nil {
		return t.sender.Close(ctx)
	}

	return nil
}

func (t *Topic) ensureSender(ctx context.Context) error {
	span, ctx := t.startSpanFromContext(ctx, "sb.Topic.ensureSender")
	defer span.Finish()

	t.senderMu.Lock()
	defer t.senderMu.Unlock()

	if t.sender == nil {
		s, err := t.namespace.newSender(ctx, t.Name)
		if err != nil {
			log.For(ctx).Error(err)
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
