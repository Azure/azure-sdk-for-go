// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/stretchr/testify/require"
)

// TestAdminClient_TopicFilterCounts_Live exercises the topic filter-count feature
// end to end: it creates a topic with a subscription carrying SQL and correlation
// filter rules, then reads the topic runtime properties and asserts the aggregated
// counts. The counts are served by the 2024-05 service API version, which the admin
// client sends by default, so this requires a service build with the feature.
func TestAdminClient_TopicFilterCounts_Live(t *testing.T) {
	adminClient := newAdminClientForTest(t, &test.NewClientOptions[ClientOptions]{})

	topicName := fmt.Sprintf("topic-fc-%X", time.Now().UnixNano())
	_, err := adminClient.CreateTopic(context.Background(), topicName, nil)
	require.NoError(t, err)
	defer deleteTopic(t, adminClient, topicName)

	subscriptionName := "sub1"
	_, err = adminClient.CreateSubscription(context.Background(), topicName, subscriptionName, nil)
	require.NoError(t, err)

	// A new subscription has a default $Default rule (a SQL TrueFilter). Add an
	// explicit SQL rule and a correlation rule so the topic-level counts are non-zero.
	_, err = adminClient.CreateRule(context.Background(), topicName, subscriptionName, &CreateRuleOptions{
		Name:   to.Ptr("sqlrule"),
		Filter: &SQLFilter{Expression: "1=1"},
	})
	require.NoError(t, err)
	_, err = adminClient.CreateRule(context.Background(), topicName, subscriptionName, &CreateRuleOptions{
		Name:   to.Ptr("corrrule"),
		Filter: &CorrelationFilter{CorrelationID: to.Ptr("abc")},
	})
	require.NoError(t, err)

	resp, err := adminClient.GetTopicRuntimeProperties(context.Background(), topicName, nil)
	require.NoError(t, err)

	// $Default (TrueFilter) + sqlrule = 2 SQL filters; corrrule = 1 correlation filter.
	require.Equal(t, int32(2), resp.SQLFilterCount)
	require.Equal(t, int32(1), resp.CorrelationFilterCount)
}
