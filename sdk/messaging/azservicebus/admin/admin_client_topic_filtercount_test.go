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
	// A service region that has not yet deployed the topic filter-count feature omits
	// the SqlFilterCount/CorrelationFilterCount elements; the counts must default to zero.
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
	// Unmarshal a representative ATOM topic response through the same xml tags the client uses, so a
	// typo in the SqlFilterCount/CorrelationFilterCount xml tag is caught in CI. The other tests
	// populate the struct directly and would not detect a wrong tag.
	const topicXML = `<entry xmlns="http://www.w3.org/2005/Atom"><title>my-topic</title><content type="application/xml"><TopicDescription xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect"><SubscriptionCount>2</SubscriptionCount><SqlFilterCount>3</SqlFilterCount><CorrelationFilterCount>1</CorrelationFilterCount><CountDetails><ScheduledMessageCount>0</ScheduledMessageCount></CountDetails><CreatedAt>2026-01-01T00:00:00Z</CreatedAt><UpdatedAt>2026-01-01T00:00:00Z</UpdatedAt><AccessedAt>2026-01-01T00:00:00Z</AccessedAt></TopicDescription></content></entry>`

	var env *atom.TopicEnvelope
	require.NoError(t, xml.Unmarshal([]byte(topicXML), &env))

	item, err := newTopicRuntimePropertiesItem(env)
	require.NoError(t, err)
	require.Equal(t, int32(2), item.SubscriptionCount)
	require.Equal(t, int32(3), item.SQLFilterCount)
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
