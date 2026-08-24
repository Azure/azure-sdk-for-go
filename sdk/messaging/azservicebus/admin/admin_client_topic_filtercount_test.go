// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"encoding/xml"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/stretchr/testify/require"
)

func i32Ptr(v int32) *int32 { return &v }

func TestNewTopicRuntimePropertiesItem_FilterCounts(t *testing.T) {
	env := &atom.TopicEnvelope{
		Entry: &atom.Entry{Title: "my-topic"},
		Content: &atom.TopicContent{
			TopicDescription: atom.TopicDescription{
				SubscriptionCount:      i32Ptr(2),
				SQLFilterCount:         i32Ptr(7),
				CorrelationFilterCount: i32Ptr(9),
				CountDetails:           &atom.CountDetails{ScheduledMessageCount: i32Ptr(1)},
				CreatedAt:              "2026-01-01T00:00:00Z",
				UpdatedAt:              "2026-01-01T00:00:00Z",
				AccessedAt:             "2026-01-01T00:00:00Z",
			},
		},
	}

	item, err := newTopicRuntimePropertiesItem(env)
	require.NoError(t, err)
	require.Equal(t, int32(2), item.SubscriptionCount)
	require.Equal(t, int32(7), item.SQLFilterCount)
	require.Equal(t, int32(9), item.CorrelationFilterCount)
}

func TestNewTopicRuntimePropertiesItem_FilterCountsDefaultZero(t *testing.T) {
	// A namespace that does not report the counts omits the SqlFilterCount/CorrelationFilterCount
	// elements; the counts must default to zero.
	env := &atom.TopicEnvelope{
		Entry: &atom.Entry{Title: "my-topic"},
		Content: &atom.TopicContent{
			TopicDescription: atom.TopicDescription{
				SubscriptionCount: i32Ptr(1),
				CountDetails:      &atom.CountDetails{ScheduledMessageCount: i32Ptr(0)},
				CreatedAt:         "2026-01-01T00:00:00Z",
				UpdatedAt:         "2026-01-01T00:00:00Z",
				AccessedAt:        "2026-01-01T00:00:00Z",
			},
		},
	}

	item, err := newTopicRuntimePropertiesItem(env)
	require.NoError(t, err)
	require.Equal(t, int32(1), item.SubscriptionCount)
	require.Zero(t, item.SQLFilterCount)
	require.Zero(t, item.CorrelationFilterCount)
}

func TestNewTopicRuntimePropertiesItem_FilterCountsFromXML(t *testing.T) {
	// Unmarshal a real ATOM topic response through the same xml tags the client uses, so a typo in
	// the SqlFilterCount/CorrelationFilterCount xml tag is caught in CI. The other tests populate the
	// struct directly and would not detect a wrong tag.
	//
	// Captured 2026-08-13 from a live GET .../my-topic?api-version=2024-05 against a Standard
	// namespace, on a topic with one subscription carrying $Default, one SQL rule and one
	// correlation rule. Host and topic name scrubbed; every element is as the service emitted it.
	const topicXML = `<entry xmlns="http://www.w3.org/2005/Atom"><id>https://CONTOSO.servicebus.windows.net/my-topic?api-version=2024-05</id><title type="text">my-topic</title><published>2026-08-14T00:46:28Z</published><updated>2026-08-14T00:46:28Z</updated><author><name>CONTOSO</name></author><link rel="self" href="https://CONTOSO.servicebus.windows.net/my-topic?api-version=2024-05"/><content type="application/xml"><TopicDescription xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect" xmlns:i="http://www.w3.org/2001/XMLSchema-instance"><DefaultMessageTimeToLive>P10675199DT2H48M5.4775807S</DefaultMessageTimeToLive><MaxSizeInMegabytes>1024</MaxSizeInMegabytes><RequiresDuplicateDetection>false</RequiresDuplicateDetection><DuplicateDetectionHistoryTimeWindow>PT10M</DuplicateDetectionHistoryTimeWindow><EnableBatchedOperations>true</EnableBatchedOperations><SizeInBytes>0</SizeInBytes><FilteringMessagesBeforePublishing>false</FilteringMessagesBeforePublishing><IsAnonymousAccessible>false</IsAnonymousAccessible><AuthorizationRules></AuthorizationRules><Status>Active</Status><CreatedAt>2026-08-14T00:46:28.4314346Z</CreatedAt><UpdatedAt>2026-08-14T00:46:28.4314346Z</UpdatedAt><AccessedAt>2026-08-14T00:46:32.9628806Z</AccessedAt><SupportOrdering>true</SupportOrdering><CountDetails xmlns:d2p1="http://schemas.microsoft.com/netservices/2011/06/servicebus"><d2p1:ActiveMessageCount>0</d2p1:ActiveMessageCount><d2p1:DeadLetterMessageCount>0</d2p1:DeadLetterMessageCount><d2p1:ScheduledMessageCount>0</d2p1:ScheduledMessageCount><d2p1:TransferMessageCount>0</d2p1:TransferMessageCount><d2p1:TransferDeadLetterMessageCount>0</d2p1:TransferDeadLetterMessageCount></CountDetails><SubscriptionCount>1</SubscriptionCount><AutoDeleteOnIdle>P10675199DT2H48M5.4775807S</AutoDeleteOnIdle><EnablePartitioning>false</EnablePartitioning><EntityAvailabilityStatus>Available</EntityAvailabilityStatus><EnableSubscriptionPartitioning>false</EnableSubscriptionPartitioning><EnableExpress>false</EnableExpress><MaxMessageSizeInKilobytes>256</MaxMessageSizeInKilobytes><SqlFilterCount>2</SqlFilterCount><CorrelationFilterCount>1</CorrelationFilterCount></TopicDescription></content></entry>`

	var env *atom.TopicEnvelope
	require.NoError(t, xml.Unmarshal([]byte(topicXML), &env))

	item, err := newTopicRuntimePropertiesItem(env)
	require.NoError(t, err)
	require.Equal(t, int32(1), item.SubscriptionCount)
	// The subscription's $Default rule is a SQL TrueFilter and the service counts it, so one
	// explicit SQL rule reads as 2.
	require.Equal(t, int32(2), item.SQLFilterCount)
	require.Equal(t, int32(1), item.CorrelationFilterCount)
}

func TestNewTopicRuntimePropertiesItem_FilterCountsAbsentFromXML(t *testing.T) {
	// An older api-version omits the filter-count elements entirely; the counts must default to 0.
	const topicXML = `<entry xmlns="http://www.w3.org/2005/Atom"><title>my-topic</title><content type="application/xml"><TopicDescription xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect"><SubscriptionCount>1</SubscriptionCount><CountDetails><ScheduledMessageCount>0</ScheduledMessageCount></CountDetails><CreatedAt>2026-01-01T00:00:00Z</CreatedAt><UpdatedAt>2026-01-01T00:00:00Z</UpdatedAt><AccessedAt>2026-01-01T00:00:00Z</AccessedAt></TopicDescription></content></entry>`

	var env *atom.TopicEnvelope
	require.NoError(t, xml.Unmarshal([]byte(topicXML), &env))

	item, err := newTopicRuntimePropertiesItem(env)
	require.NoError(t, err)
	require.Equal(t, int32(1), item.SubscriptionCount)
	require.Zero(t, item.SQLFilterCount)
	require.Zero(t, item.CorrelationFilterCount)
}
