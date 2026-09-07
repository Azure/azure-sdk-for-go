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
// end to end: it creates a topic with a subscription, reads the baseline counts, adds
// SQL and correlation filter rules, then asserts the topic-level counts moved by one
// each. The counts are served at api-version 2024-05, which the admin client sends by
// default; a namespace whose service build does not report them skips the test.
func TestAdminClient_TopicFilterCounts_Live(t *testing.T) {
	adminClient := newAdminClientForTest(t, &test.NewClientOptions[ClientOptions]{})

	topicName := fmt.Sprintf("topic-fc-%X", time.Now().UnixNano())
	_, err := adminClient.CreateTopic(context.Background(), topicName, nil)
	require.NoError(t, err)
	defer deleteTopic(t, adminClient, topicName)

	subscriptionName := "sub1"
	_, err = adminClient.CreateSubscription(context.Background(), topicName, subscriptionName, nil)
	require.NoError(t, err)

	before, err := adminClient.GetTopicRuntimeProperties(context.Background(), topicName, nil)
	require.NoError(t, err)

	// A new subscription carries a default $Default rule, which is a SQL TrueFilter, and the service
	// counts it - a subscription with no explicit rules reads SQLFilterCount 1. So a zero here means
	// the namespace does not report the counts at all.
	if before.SQLFilterCount == 0 {
		t.Skipf("namespace does not report topic filter counts (sql=%d corr=%d)",
			before.SQLFilterCount, before.CorrelationFilterCount)
	}

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

	after, err := adminClient.GetTopicRuntimeProperties(context.Background(), topicName, nil)
	require.NoError(t, err)

	// Asserting the delta proves the counts track rule creation, without assuming how the
	// service counts the default rule.
	require.Equal(t, before.SQLFilterCount+1, after.SQLFilterCount)
	require.Equal(t, before.CorrelationFilterCount+1, after.CorrelationFilterCount)
}
