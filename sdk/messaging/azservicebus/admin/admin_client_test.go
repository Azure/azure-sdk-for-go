// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"bytes"
	"context"
	cryptoRand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/conn"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestAdminClient_UsingIdentity(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	props, err := adminClient.CreateQueue(context.Background(), queueName, nil)
	require.NoError(t, err)
	require.EqualValues(t, 10, *props.MaxDeliveryCount)

	defer func() {
		deleteQueue(t, adminClient, queueName)
	}()
}

func TestAdminClient_GetNamespaceProperties(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)
	resp, err := adminClient.GetNamespaceProperties(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, resp)

	require.True(t, resp.SKU == "Standard" || resp.SKU == "Premium" || resp.SKU == "Basic")

	if resp.SKU == "Standard" || resp.SKU == "Basic" {
		// messaging units don't exist in the lower tiers
		require.Nil(t, resp.MessagingUnits)
	} else {
		require.NotNil(t, resp.MessagingUnits)
	}

	require.NotEmpty(t, resp.Name)
	require.False(t, resp.CreatedTime.IsZero())
	require.False(t, resp.ModifiedTime.IsZero())
}

func TestAdminClient_QueueWithMaxValues(t *testing.T) {
	// this is the magic TimeSpan.Max value that the portal uses to determine we want
	// infinite "time"
	var MaxTimeSpanForTests = to.Ptr("P10675199DT2H48M5.4775807S")

	adminClient := newAdminClientForTest(t, nil)

	es := EntityStatusReceiveDisabled

	authRules := createAuthorizationRulesForTest(t)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())

	_, err := adminClient.CreateQueue(context.Background(), queueName, &CreateQueueOptions{
		Properties: &QueueProperties{
			LockDuration: to.Ptr("PT45S"),
			// when you enable partitioning Service Bus will automatically create 16 partitions, each with the size
			// of MaxSizeInMegabytes. This means when we retrieve this queue we'll get 16*4096 as the size (ie: 64GB)
			EnablePartitioning:                  to.Ptr(true),
			MaxSizeInMegabytes:                  to.Ptr(int32(4096)),
			RequiresDuplicateDetection:          to.Ptr(true),
			RequiresSession:                     to.Ptr(true),
			DefaultMessageTimeToLive:            MaxTimeSpanForTests,
			DeadLetteringOnMessageExpiration:    to.Ptr(true),
			DuplicateDetectionHistoryTimeWindow: to.Ptr("PT4H"),
			MaxDeliveryCount:                    to.Ptr(int32(100)),
			EnableBatchedOperations:             to.Ptr(false),
			Status:                              &es,
			AutoDeleteOnIdle:                    MaxTimeSpanForTests,
			UserMetadata:                        to.Ptr("some metadata"),
			AuthorizationRules:                  authRules,
		},
	})
	require.NoError(t, err)

	defer deleteQueue(t, adminClient, queueName)

	resp, err := adminClient.GetQueue(context.Background(), queueName, nil)
	require.NoError(t, err)

	require.EqualValues(t, QueueProperties{
		LockDuration: to.Ptr("PT45S"),
		// ie: this response was from a partitioned queue so the size is the original max size * # of partitions
		EnablePartitioning:                  to.Ptr(true),
		MaxSizeInMegabytes:                  to.Ptr(int32(16 * 4096)),
		RequiresDuplicateDetection:          to.Ptr(true),
		RequiresSession:                     to.Ptr(true),
		DefaultMessageTimeToLive:            MaxTimeSpanForTests,
		DeadLetteringOnMessageExpiration:    to.Ptr(true),
		DuplicateDetectionHistoryTimeWindow: to.Ptr("PT4H"),
		MaxDeliveryCount:                    to.Ptr(int32(100)),
		EnableBatchedOperations:             to.Ptr(false),
		Status:                              &es,
		AutoDeleteOnIdle:                    MaxTimeSpanForTests,
		UserMetadata:                        to.Ptr("some metadata"),
		AuthorizationRules:                  authRules,
		MaxMessageSizeInKilobytes:           to.Ptr(int64(256)), // the default size for standard.
	}, resp.QueueProperties)

	runtimeResp, err := adminClient.GetQueueRuntimeProperties(context.Background(), queueName, nil)
	require.NoError(t, err)

	require.False(t, runtimeResp.CreatedAt.IsZero())
	require.False(t, runtimeResp.UpdatedAt.IsZero())
	require.True(t, runtimeResp.AccessedAt.IsZero())
	require.Zero(t, runtimeResp.ActiveMessageCount)
	require.Zero(t, runtimeResp.DeadLetterMessageCount)
	require.Zero(t, runtimeResp.ScheduledMessageCount)
	require.Zero(t, runtimeResp.SizeInBytes)
	require.Zero(t, runtimeResp.TotalMessageCount)
}

func TestAdminClient_CreateQueue_Standard(t *testing.T) {
	testCreateQueue(t, false)
}

func TestAdminClient_CreateQueue_Premium(t *testing.T) {
	testCreateQueue(t, true)
}

func testCreateQueue(t *testing.T, isPremium bool) {
	adminClient := newAdminClientForTest(t, &test.NewClientOptions[ClientOptions]{
		UsePremium: isPremium,
	})

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())

	es := EntityStatusReceiveDisabled

	createQueueOptions := &CreateQueueOptions{
		Properties: &QueueProperties{
			LockDuration:                        to.Ptr("PT45S"),
			RequiresDuplicateDetection:          to.Ptr(true),
			RequiresSession:                     to.Ptr(true),
			DefaultMessageTimeToLive:            to.Ptr("PT6H"),
			DeadLetteringOnMessageExpiration:    to.Ptr(true),
			DuplicateDetectionHistoryTimeWindow: to.Ptr("PT4H"),
			MaxDeliveryCount:                    to.Ptr(int32(100)),
			EnableBatchedOperations:             to.Ptr(false),
			Status:                              &es,
			AutoDeleteOnIdle:                    to.Ptr("PT10M"),
		},
	}

	if isPremium {
		createQueueOptions.Properties.MaxMessageSizeInKilobytes = to.Ptr(int64(102400))

		// no partitioning in premium
		createQueueOptions.Properties.EnablePartitioning = to.Ptr(false)
	} else {
		// can't update message size in standard
		createQueueOptions.Properties.MaxMessageSizeInKilobytes = nil

		// when you enable partitioning Service Bus will automatically create 16 partitions, each with the size
		// of MaxSizeInMegabytes. This means when we retrieve this queue we'll get 16*4096 as the size (ie: 64GB)
		createQueueOptions.Properties.EnablePartitioning = to.Ptr(true)
		createQueueOptions.Properties.MaxSizeInMegabytes = to.Ptr(int32(4096))
	}

	createResp, err := adminClient.CreateQueue(context.Background(), queueName, createQueueOptions)
	require.NoError(t, err)

	defer func() {
		deleteQueue(t, adminClient, queueName)
	}()

	expectedQueueProperties := QueueProperties{
		LockDuration:                        to.Ptr("PT45S"),
		RequiresDuplicateDetection:          to.Ptr(true),
		RequiresSession:                     to.Ptr(true),
		DefaultMessageTimeToLive:            to.Ptr("PT6H"),
		DeadLetteringOnMessageExpiration:    to.Ptr(true),
		DuplicateDetectionHistoryTimeWindow: to.Ptr("PT4H"),
		MaxDeliveryCount:                    to.Ptr(int32(100)),
		EnableBatchedOperations:             to.Ptr(false),
		Status:                              &es,
		AutoDeleteOnIdle:                    to.Ptr("PT10M"),
	}

	if isPremium {
		expectedQueueProperties.MaxMessageSizeInKilobytes = to.Ptr(int64(102400))

		expectedQueueProperties.EnablePartitioning = to.Ptr(false)
		expectedQueueProperties.MaxSizeInMegabytes = to.Ptr(int32(1024))
	} else {
		// (the message size for SB Standard)
		expectedQueueProperties.MaxMessageSizeInKilobytes = to.Ptr(int64(256))

		expectedQueueProperties.EnablePartitioning = to.Ptr(true)
		expectedQueueProperties.MaxSizeInMegabytes = to.Ptr(int32(16 * 4096))
	}

	require.Equal(t, createResp.QueueName, queueName)
	require.Equal(t, expectedQueueProperties, createResp.QueueProperties)

	getResp, err := adminClient.GetQueue(context.Background(), queueName, nil)
	require.NoError(t, err)

	require.Equal(t, getResp.QueueName, queueName)
	require.Equal(t, getResp.QueueProperties, createResp.QueueProperties)

	// ensure we can round-trip
	updateResp, err := adminClient.UpdateQueue(context.Background(), queueName, getResp.QueueProperties, nil)
	require.NoError(t, err)

	require.Equal(t, updateResp.QueueName, queueName)
	require.EqualValues(t, expectedQueueProperties, updateResp.QueueProperties)

	runtimeResp, err := adminClient.GetQueueRuntimeProperties(context.Background(), queueName, nil)
	require.NoError(t, err)

	require.Equal(t, queueName, runtimeResp.QueueName)
	require.NotZero(t, runtimeResp.CreatedAt)
	require.NotZero(t, runtimeResp.UpdatedAt)
	require.Zero(t, runtimeResp.AccessedAt)
	require.Zero(t, runtimeResp.TotalMessageCount)
	require.Zero(t, runtimeResp.SizeInBytes)
}

func TestAdminClient_UpdateQueue(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	createdProps, err := adminClient.CreateQueue(context.Background(), queueName, nil)
	require.NoError(t, err)

	defer func() {
		deleteQueue(t, adminClient, queueName)
	}()

	createdProps.MaxDeliveryCount = to.Ptr(int32(101))
	createdProps.QueueProperties.AuthorizationRules = createAuthorizationRulesForTest(t)

	updatedProps, err := adminClient.UpdateQueue(context.Background(), queueName, createdProps.QueueProperties, nil)
	require.NoError(t, err)

	require.EqualValues(t, 101, *updatedProps.MaxDeliveryCount)
	require.EqualValues(t, createdProps.QueueProperties.AuthorizationRules, updatedProps.AuthorizationRules)

	// try changing a value that's not allowed
	updatedProps.RequiresSession = to.Ptr(true)
	updatedProps, err = adminClient.UpdateQueue(context.Background(), queueName, updatedProps.QueueProperties, nil)
	require.Contains(t, err.Error(), "The value for the RequiresSession property of an existing Queue cannot be changed")

	updatedProps, err = adminClient.UpdateQueue(context.Background(), "non-existent-queue", createdProps.QueueProperties, nil)
	// a little awkward, we'll make these programatically inspectable as we add in better error handling.
	require.Contains(t, err.Error(), "404 Not Found")

	var asResponseErr *azcore.ResponseError
	require.ErrorAs(t, err, &asResponseErr)
	require.EqualValues(t, 404, asResponseErr.StatusCode)
}

func TestAdminClient_ListQueues(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)

	var expectedQueues []string
	now := time.Now().UnixNano()

	for i := 0; i < 3; i++ {
		queueName := strings.ToLower(fmt.Sprintf("list-queues-%d-%X", i, now))
		expectedQueues = append(expectedQueues, queueName)

		_, err := adminClient.CreateQueue(context.Background(), queueName, &CreateQueueOptions{
			Properties: &QueueProperties{
				MaxDeliveryCount: to.Ptr(int32(int32(i + 10))),
			},
		})
		require.NoError(t, err)

		defer deleteQueue(t, adminClient, queueName)
	}

	// we skipped the first queue so it shouldn't come back in the results.
	pager := adminClient.NewListQueuesPager(nil)
	all := map[string]QueueItem{}

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)

		for _, props := range page.Queues {
			all[props.QueueName] = props
		}
	}

	// sanity check - the queues we created exist and their deserialization is
	// working.
	for i, expectedQueue := range expectedQueues {
		props, exists := all[expectedQueue]
		require.True(t, exists)
		require.EqualValues(t, i+10, *props.MaxDeliveryCount)
	}
}

func TestAdminClient_ListQueuesRuntimeProperties(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)

	var expectedQueues []string
	now := time.Now().UnixNano()

	for i := 0; i < 3; i++ {
		queueName := strings.ToLower(fmt.Sprintf("list-queuert-%d-%X", i, now))
		expectedQueues = append(expectedQueues, queueName)

		_, err := adminClient.CreateQueue(context.Background(), queueName, &CreateQueueOptions{
			Properties: &QueueProperties{
				MaxDeliveryCount: to.Ptr(int32(int32(i + 10))),
			}})
		require.NoError(t, err)

		defer deleteQueue(t, adminClient, queueName)
	}

	// we skipped the first queue so it shouldn't come back in the results.
	pager := adminClient.NewListQueuesRuntimePropertiesPager(&ListQueuesRuntimePropertiesOptions{
		MaxPageSize: 2,
	})
	all := map[string]QueueRuntimePropertiesItem{}

	times := 0

	for pager.More() {
		times++

		resp, err := pager.NextPage(context.Background())
		require.NoError(t, err)

		require.LessOrEqual(t, len(resp.QueueRuntimeProperties), 2)

		for _, queueRuntimeItem := range resp.QueueRuntimeProperties {
			// _, exists := all[queueRuntimeItem.QueueName]
			// require.False(t, exists, fmt.Sprintf("Each queue result should be unique but found more than one of '%s'", queueRuntimeItem.QueueName))
			all[queueRuntimeItem.QueueName] = queueRuntimeItem
		}
	}

	require.GreaterOrEqual(t, times, 2)

	// sanity check - the queues we created exist and their deserialization is
	// working.
	for _, expectedQueue := range expectedQueues {
		props, exists := all[expectedQueue]
		require.True(t, exists)
		require.NotEqualValues(t, time.Time{}, props.CreatedAt)
	}
}

func TestAdminClient_DurationToStringPtr(t *testing.T) {
	// The actual value max is TimeSpan.Max so we just assume that's what the user wants if they specify our time.Duration value
	require.EqualValues(t, "P10675199DT2H48M5.4775807S", utils.DurationTo8601Seconds(utils.MaxTimeDuration), "Max time.Duration gets converted to TimeSpan.Max")

	require.EqualValues(t, "PT0M1S", utils.DurationTo8601Seconds(time.Second))
	require.EqualValues(t, "PT1M0S", utils.DurationTo8601Seconds(time.Minute))
	require.EqualValues(t, "PT1M1S", utils.DurationTo8601Seconds(time.Minute+time.Second))
	require.EqualValues(t, "PT60M0S", utils.DurationTo8601Seconds(time.Hour))
	require.EqualValues(t, "PT61M1S", utils.DurationTo8601Seconds(time.Hour+time.Minute+time.Second))
}

func TestAdminClient_ISO8601StringToDuration(t *testing.T) {
	str := "PT10M1S"
	duration, err := utils.ISO8601StringToDuration(&str)
	require.NoError(t, err)
	require.EqualValues(t, (10*time.Minute)+time.Second, *duration)

	duration, err = utils.ISO8601StringToDuration(nil)
	require.NoError(t, err)
	require.Nil(t, duration)

	str = "PT1S"
	duration, err = utils.ISO8601StringToDuration(&str)
	require.NoError(t, err)
	require.EqualValues(t, time.Second, *duration)

	str = "PT1M"
	duration, err = utils.ISO8601StringToDuration(&str)
	require.NoError(t, err)
	require.EqualValues(t, time.Minute, *duration)

	// this is the .NET timespan max
	str = "P10675199DT2H48M5.4775807S"
	duration, err = utils.ISO8601StringToDuration(&str)
	require.NoError(t, err)
	require.EqualValues(t, utils.MaxTimeDuration, *duration)

	// this is the Java equivalent
	str = "PT256204778H48M5.4775807S"
	duration, err = utils.ISO8601StringToDuration(&str)
	require.NoError(t, err)
	require.EqualValues(t, utils.MaxTimeDuration, *duration)
}

func TestAdminClient_TopicAndSubscription_Standard(t *testing.T) {
	testTopicCreation(t, false)
}

func TestAdminClient_TopicAndSubscription_Premium(t *testing.T) {
	testTopicCreation(t, true)
}

func testTopicCreation(t *testing.T, isPremium bool) {
	adminClient := newAdminClientForTest(t, &test.NewClientOptions[ClientOptions]{
		UsePremium: isPremium,
	})

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	subscriptionName := fmt.Sprintf("sub-%X", time.Now().UnixNano())

	status := EntityStatusActive

	createTopicOptions := &CreateTopicOptions{
		Properties: &TopicProperties{
			RequiresDuplicateDetection:          to.Ptr(true),
			DefaultMessageTimeToLive:            to.Ptr("PT3M"),
			DuplicateDetectionHistoryTimeWindow: to.Ptr("PT4M"),
			EnableBatchedOperations:             to.Ptr(true),
			Status:                              &status,
			AutoDeleteOnIdle:                    to.Ptr("PT7M"),
			SupportOrdering:                     to.Ptr(true),
			UserMetadata:                        to.Ptr("user metadata"),
		},
	}

	if isPremium {
		// premium doesn't support partitioning.
		createTopicOptions.Properties.EnablePartitioning = to.Ptr(false)

		// premium allows you to update the max message size
		createTopicOptions.Properties.MaxMessageSizeInKilobytes = to.Ptr(int64(102400))
	} else {
		createTopicOptions.Properties.EnablePartitioning = to.Ptr(true)
		createTopicOptions.Properties.MaxSizeInMegabytes = to.Ptr(int32(2048))

		// standard can't change MaxMessageSizeInKilobytes
		createTopicOptions.Properties.MaxMessageSizeInKilobytes = nil
	}

	// check topic properties, existence
	createResp, err := adminClient.CreateTopic(context.Background(), topicName, createTopicOptions)
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	expectedTopicProps := TopicProperties{
		MaxSizeInMegabytes:                  to.Ptr(int32(1024)),
		RequiresDuplicateDetection:          to.Ptr(true),
		DefaultMessageTimeToLive:            to.Ptr("PT3M"),
		DuplicateDetectionHistoryTimeWindow: to.Ptr("PT4M"),
		EnableBatchedOperations:             to.Ptr(true),
		Status:                              &status,
		AutoDeleteOnIdle:                    to.Ptr("PT7M"),
		SupportOrdering:                     to.Ptr(true),
		UserMetadata:                        to.Ptr("user metadata"),
	}

	if isPremium {
		expectedTopicProps.MaxMessageSizeInKilobytes = to.Ptr(int64(102400))

		// no partitioning in premium.
		expectedTopicProps.EnablePartitioning = to.Ptr(false)
		expectedTopicProps.MaxSizeInMegabytes = to.Ptr(int32(1024))
	} else {
		expectedTopicProps.MaxMessageSizeInKilobytes = to.Ptr(int64(256))

		// enabling partitioning increases our max size because of the 16 partition),
		expectedTopicProps.EnablePartitioning = to.Ptr(true)
		expectedTopicProps.MaxSizeInMegabytes = to.Ptr(int32(16 * 2048))
	}

	require.Equal(t, CreateTopicResponse{
		TopicName:       topicName,
		TopicProperties: expectedTopicProps,
	}, createResp)

	getResp, err := adminClient.GetTopic(context.Background(), topicName, nil)
	require.NoError(t, err)

	require.Equal(t, &GetTopicResponse{
		TopicName:       topicName,
		TopicProperties: getResp.TopicProperties,
	}, getResp)

	runtimeResp, err := adminClient.GetTopicRuntimeProperties(context.Background(), topicName, nil)
	require.NoError(t, err)

	require.Equal(t, topicName, runtimeResp.TopicName)
	require.False(t, runtimeResp.CreatedAt.IsZero())
	require.False(t, runtimeResp.UpdatedAt.IsZero())
	require.True(t, runtimeResp.AccessedAt.IsZero())
	require.Zero(t, runtimeResp.SubscriptionCount)
	require.Zero(t, runtimeResp.ScheduledMessageCount)
	require.Zero(t, runtimeResp.SizeInBytes)

	createSubWithPropsResp, err := adminClient.CreateSubscription(context.Background(), topicName, subscriptionName, &CreateSubscriptionOptions{
		Properties: &SubscriptionProperties{
			LockDuration:                                    to.Ptr("PT3M"),
			RequiresSession:                                 to.Ptr(false),
			DefaultMessageTimeToLive:                        to.Ptr("PT7M"),
			DeadLetteringOnMessageExpiration:                to.Ptr(true),
			EnableDeadLetteringOnFilterEvaluationExceptions: to.Ptr(false),
			MaxDeliveryCount:                                to.Ptr(int32(11)),
			Status:                                          &status,
			EnableBatchedOperations:                         to.Ptr(false),
			AutoDeleteOnIdle:                                to.Ptr("PT11M"),
			UserMetadata:                                    to.Ptr("user metadata"),
		},
	})
	require.NoError(t, err)

	defer deleteSubscription(t, adminClient, topicName, subscriptionName)

	require.Equal(t, CreateSubscriptionResponse{
		SubscriptionName: subscriptionName,
		TopicName:        topicName,
		SubscriptionProperties: SubscriptionProperties{
			LockDuration:                                    to.Ptr("PT3M"),
			RequiresSession:                                 to.Ptr(false),
			DefaultMessageTimeToLive:                        to.Ptr("PT7M"),
			DeadLetteringOnMessageExpiration:                to.Ptr(true),
			EnableDeadLetteringOnFilterEvaluationExceptions: to.Ptr(false),
			MaxDeliveryCount:                                to.Ptr(int32(11)),
			Status:                                          &status,
			EnableBatchedOperations:                         to.Ptr(false),
			AutoDeleteOnIdle:                                to.Ptr("PT11M"),
			UserMetadata:                                    to.Ptr("user metadata"),
		},
	}, createSubWithPropsResp)

	runtimePropsResp, err := adminClient.GetSubscriptionRuntimeProperties(context.Background(), topicName, subscriptionName, nil)
	require.NoError(t, err)

	require.Equal(t, topicName, runtimePropsResp.TopicName)
	require.Equal(t, subscriptionName, runtimePropsResp.SubscriptionName)
}

func TestAdminClient_TopicAndSubscription_WithFalseFilterDefaultSubscriptionRule(t *testing.T) {
	adminClient, topicName := createAdminClientWithTestTopic(t)

	defer deleteTopic(t, adminClient, topicName)

	subscriptionName := createTestSubscriptionWithDefaultRule(
		t,
		adminClient,
		topicName,
		&RuleProperties{Filter: &FalseFilter{}},
	)

	// Even though a default rule is created on the subscription, the response body does not include it.
	// which is why we can discard this response and need to fetch the rule explicitly.
	// (It is also not included when fetching the subscription)
	defaultRule, err := adminClient.GetRule(context.Background(), topicName, subscriptionName, "$Default", nil)
	require.NoError(t, err)

	require.Equal(t, "$Default", defaultRule.Name)
	require.IsType(t, &FalseFilter{}, defaultRule.Filter)
	require.Nil(t, defaultRule.Action)
}

func TestAdminClient_TopicAndSubscription_WithCustomFilterDefaultSubscriptionRule(t *testing.T) {
	adminClient, topicName := createAdminClientWithTestTopic(t)

	defer deleteTopic(t, adminClient, topicName)

	customSqlFilter := &SQLFilter{
		Expression: "SomeProperty LIKE 'O%'",
	}

	subscriptionName := createTestSubscriptionWithDefaultRule(
		t,
		adminClient,
		topicName,
		&RuleProperties{
			Name:   "TestRule",
			Filter: customSqlFilter,
		},
	)

	defaultRule, err := adminClient.GetRule(context.Background(), topicName, subscriptionName, "TestRule", nil)
	require.NoError(t, err)

	require.Equal(t, "TestRule", defaultRule.Name)
	require.Equal(t, customSqlFilter, defaultRule.Filter)
	require.Nil(t, defaultRule.Action)
}

func TestAdminClient_TopicAndSubscription_WithActionDefaultSubscriptionRule(t *testing.T) {
	adminClient, topicName := createAdminClientWithTestTopic(t)

	defer deleteTopic(t, adminClient, topicName)

	ruleAction := &SQLAction{
		Expression: "SET MessageID=@stringVar",
		Parameters: map[string]any{
			"@stringVar": "hello world",
		},
	}

	subscriptionName := createTestSubscriptionWithDefaultRule(
		t,
		adminClient,
		topicName,
		&RuleProperties{
			Action: ruleAction,
		},
	)

	defaultRule, err := adminClient.GetRule(context.Background(), topicName, subscriptionName, "$Default", nil)
	require.NoError(t, err)

	require.Equal(t, "$Default", defaultRule.Name)
	require.Equal(t, ruleAction, defaultRule.Action)
	require.Equal(t, defaultRule.Filter, &TrueFilter{})
}

func TestAdminClient_TopicAndSubscription_WithActionAndFilterDefaultSubscriptionRule(t *testing.T) {
	adminClient, topicName := createAdminClientWithTestTopic(t)

	defer deleteTopic(t, adminClient, topicName)

	ruleAction := &SQLAction{
		Expression: "SET MessageID=@stringVar",
		Parameters: map[string]any{
			"@stringVar": "hello world",
		},
	}

	ruleFilter := &SQLFilter{
		Expression: "SomeProperty LIKE 'O%'",
	}

	subscriptionName := createTestSubscriptionWithDefaultRule(
		t,
		adminClient,
		topicName,
		&RuleProperties{
			Action: ruleAction,
			Filter: ruleFilter,
		},
	)

	defaultRule, err := adminClient.GetRule(context.Background(), topicName, subscriptionName, "$Default", nil)
	require.NoError(t, err)

	require.Equal(t, "$Default", defaultRule.Name)
	require.EqualValues(t, ruleAction, defaultRule.Action)
	require.EqualValues(t, ruleFilter, defaultRule.Filter)
}

func createAdminClientWithTestTopic(t *testing.T) (*Client, string) {
	adminClient := newAdminClientForTest(t, nil)

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	_, err := adminClient.CreateTopic(context.Background(), topicName, nil)
	require.NoError(t, err)

	return adminClient, topicName
}

func createTestSubscriptionWithDefaultRule(t *testing.T, adminClient *Client, topicName string, defaultRule *RuleProperties) string {
	subscriptionName := fmt.Sprintf("subscription-%X", time.Now().UnixNano())

	_, err := adminClient.CreateSubscription(context.Background(), topicName, subscriptionName, &CreateSubscriptionOptions{
		Properties: &SubscriptionProperties{
			LockDuration:                                    to.Ptr("PT3M"),
			RequiresSession:                                 to.Ptr(false),
			DefaultMessageTimeToLive:                        to.Ptr("PT7M"),
			DeadLetteringOnMessageExpiration:                to.Ptr(true),
			EnableDeadLetteringOnFilterEvaluationExceptions: to.Ptr(false),
			MaxDeliveryCount:                                to.Ptr(int32(11)),
			Status:                                          to.Ptr(EntityStatusActive),
			EnableBatchedOperations:                         to.Ptr(false),
			AutoDeleteOnIdle:                                to.Ptr("PT11M"),
			UserMetadata:                                    to.Ptr("user metadata"),
			DefaultRule:                                     defaultRule,
		},
	})
	require.NoError(t, err)

	return subscriptionName
}

func TestAdminClient_Forwarding(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	forwardToQueueName := fmt.Sprintf("queue-fwd-%X", time.Now().UnixNano())

	_, err := adminClient.CreateQueue(context.Background(), forwardToQueueName, nil)
	require.NoError(t, err)

	defer deleteQueue(t, adminClient, forwardToQueueName)

	_, err = adminClient.CreateTopic(context.Background(), topicName, nil)
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	cs := test.MustGetEnvVar(t, test.EnvKeyConnectionString)
	parsed, err := conn.ParseConnectionString(cs)
	require.NoError(t, err)
	forwardToName := to.Ptr(fmt.Sprintf("sb://%s/%s", parsed.FullyQualifiedNamespace, forwardToQueueName))

	_, err = adminClient.CreateSubscription(context.Background(), topicName, "sub1", &CreateSubscriptionOptions{
		Properties: &SubscriptionProperties{
			ForwardTo:                     forwardToName,
			ForwardDeadLetteredMessagesTo: forwardToName,
		},
	})
	require.NoError(t, err)

	sub, err := adminClient.GetSubscription(context.Background(), topicName, "sub1", nil)
	require.NoError(t, err)

	require.Equal(t, forwardToName, sub.ForwardTo)
	require.Equal(t, forwardToName, sub.ForwardDeadLetteredMessagesTo)

	_, err = adminClient.CreateQueue(context.Background(), queueName, &CreateQueueOptions{
		Properties: &QueueProperties{
			ForwardTo:                     forwardToName,
			ForwardDeadLetteredMessagesTo: forwardToName,
		},
	})
	require.NoError(t, err)

	_, err = adminClient.GetQueue(context.Background(), queueName, nil)
	require.NoError(t, err)

	defer deleteQueue(t, adminClient, queueName)

	require.Equal(t, forwardToName, sub.ForwardTo)
	require.Equal(t, forwardToName, sub.ForwardDeadLetteredMessagesTo)
}

func TestAdminClient_UpdateTopic(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	addResp, err := adminClient.CreateTopic(context.Background(), topicName, nil)
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	addResp.AutoDeleteOnIdle = to.Ptr("PT11M")
	addResp.AuthorizationRules = createAuthorizationRulesForTest(t)

	updateResp, err := adminClient.UpdateTopic(context.Background(), topicName, addResp.TopicProperties, nil)
	require.NoError(t, err)

	require.Equal(t, topicName, updateResp.TopicName)
	require.EqualValues(t, "PT11M", *updateResp.AutoDeleteOnIdle)
	require.EqualValues(t, addResp.AuthorizationRules, updateResp.AuthorizationRules)

	getResp, err := adminClient.GetTopic(context.Background(), topicName, nil)
	require.NoError(t, err)
	require.Equal(t, getResp.TopicProperties, updateResp.TopicProperties)

	// try changing a value that's not allowed
	updateResp.EnablePartitioning = to.Ptr(true)
	updateResp, err = adminClient.UpdateTopic(context.Background(), topicName, updateResp.TopicProperties, nil)
	require.Contains(t, err.Error(), "Partitioning cannot be changed for Topic. ")

	updateResp, err = adminClient.UpdateTopic(context.Background(), "non-existent-topic", addResp.TopicProperties, nil)
	// a little awkward, we'll make these programatically inspectable as we add in better error handling.
	require.Contains(t, err.Error(), "404 Not Found")

	var asResponseErr *azcore.ResponseError
	require.ErrorAs(t, err, &asResponseErr)
	require.EqualValues(t, 404, asResponseErr.StatusCode)
}

func TestAdminClient_ListTopics(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)

	var expectedTopics []string
	now := time.Now().UnixNano()

	wg := sync.WaitGroup{}

	for i := 0; i < 3; i++ {
		topicName := strings.ToLower(fmt.Sprintf("list-topic-%d-%X", i, now))
		expectedTopics = append(expectedTopics, topicName)
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			_, err := adminClient.CreateTopic(context.Background(), topicName, &CreateTopicOptions{
				Properties: &TopicProperties{
					DefaultMessageTimeToLive: to.Ptr(fmt.Sprintf("PT%dM", i+1)),
				},
			})
			require.NoError(t, err)
		}(i)

		defer deleteTopic(t, adminClient, topicName)
	}

	wg.Wait()

	// we skipped the first topic so it shouldn't come back in the results.
	pager := adminClient.NewListTopicsPager(&ListTopicsOptions{
		MaxPageSize: 2,
	})
	all := map[string]TopicItem{}

	times := 0

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)

		times++

		require.LessOrEqual(t, len(page.Topics), 2)

		for _, topicItem := range page.Topics {
			// _, exists := all[topicItem.TopicName]
			// require.False(t, exists, fmt.Sprintf("Each topic result should be unique but found more than one of '%s'", topicItem.TopicName))
			all[topicItem.TopicName] = topicItem
		}
	}

	require.GreaterOrEqual(t, times, 2)

	// sanity check - the topics we created exist and their deserialization is
	// working.
	for i, expectedTopic := range expectedTopics {
		props, exists := all[expectedTopic]
		require.True(t, exists)
		require.EqualValues(t, fmt.Sprintf("PT%dM", i+1), *props.DefaultMessageTimeToLive)
	}
}

func TestAdminClient_ListTopicsRuntimeProperties(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)

	var expectedTopics []string
	now := time.Now().UnixNano()

	for i := 0; i < 3; i++ {
		topicName := strings.ToLower(fmt.Sprintf("list-topicrt-%d-%X", i, now))
		expectedTopics = append(expectedTopics, topicName)

		_, err := adminClient.CreateTopic(context.Background(), topicName, nil)
		require.NoError(t, err)

		defer deleteTopic(t, adminClient, topicName)
	}

	times := 0

	// we skipped the first topic so it shouldn't come back in the results.
	pager := adminClient.NewListTopicsRuntimePropertiesPager(&ListTopicsRuntimePropertiesOptions{
		MaxPageSize: 2,
	})
	all := map[string]TopicRuntimePropertiesItem{}

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)

		times++

		require.LessOrEqual(t, len(page.TopicRuntimeProperties), 2)

		for _, item := range page.TopicRuntimeProperties {
			// _, exists := all[item.TopicName]
			// require.False(t, exists, fmt.Sprintf("Each topic result should be unique but found more than one of '%s'", item.TopicName))
			all[item.TopicName] = item
		}
	}

	require.GreaterOrEqual(t, times, 2)

	// sanity check - the topics we created exist and their deserialization is
	// working.
	for _, expectedTopic := range expectedTopics {
		props, exists := all[expectedTopic]
		require.True(t, exists)
		require.NotEqualValues(t, time.Time{}, props.CreatedAt)
	}
}

func TestAdminClient_ListSubscriptions(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)

	now := time.Now().UnixNano()
	topicName := strings.ToLower(fmt.Sprintf("listsub-%X", now))

	_, err := adminClient.CreateTopic(context.Background(), topicName, nil)
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	var expectedSubscriptions []string

	for i := 0; i < 3; i++ {
		subName := strings.ToLower(fmt.Sprintf("sub-%d-%X", i, now))
		expectedSubscriptions = append(expectedSubscriptions, subName)

		_, err = adminClient.CreateSubscription(context.Background(), topicName, subName, &CreateSubscriptionOptions{
			Properties: &SubscriptionProperties{
				DefaultMessageTimeToLive: to.Ptr(fmt.Sprintf("PT%dM", i+1)),
			},
		})
		require.NoError(t, err)

		defer deleteSubscription(t, adminClient, topicName, subName)
	}

	// we skipped the first topic so it shouldn't come back in the results.
	pager := adminClient.NewListSubscriptionsPager(topicName, &ListSubscriptionsOptions{
		MaxPageSize: 2,
	})
	all := map[string]SubscriptionPropertiesItem{}

	times := 0

	for pager.More() {
		times++
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)

		require.LessOrEqual(t, len(page.Subscriptions), 2)

		for _, item := range page.Subscriptions {
			// _, exists := all[item.SubscriptionName]
			// require.False(t, exists, fmt.Sprintf("Each subscription result should be unique but found more than one of '%s'", item.SubscriptionName))
			all[item.SubscriptionName] = item
		}
	}

	require.GreaterOrEqual(t, times, 2)

	// sanity check - the subscriptions we created exist and their deserialization is
	// working.
	for i, expectedTopic := range expectedSubscriptions {
		props, exists := all[expectedTopic]
		require.True(t, exists)
		require.Equal(t, topicName, props.TopicName)
		require.EqualValues(t, fmt.Sprintf("PT%dM", i+1), *props.DefaultMessageTimeToLive)
	}
}

func TestAdminClient_ListSubscriptionRuntimeProperties(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)

	now := time.Now().UnixNano()
	topicName := strings.ToLower(fmt.Sprintf("listsubrt-%X", now))

	_, err := adminClient.CreateTopic(context.Background(), topicName, nil)
	require.NoError(t, err)

	var expectedSubs []string

	for i := 0; i < 3; i++ {
		subscriptionName := strings.ToLower(fmt.Sprintf("sub-%d-%X", i, now))
		expectedSubs = append(expectedSubs, subscriptionName)

		_, err = adminClient.CreateSubscription(context.Background(), topicName, subscriptionName, nil)
		require.NoError(t, err)

		defer deleteSubscription(t, adminClient, topicName, subscriptionName)
	}

	// we skipped the first subscription so it shouldn't come back in the results.
	pager := adminClient.NewListSubscriptionsRuntimePropertiesPager(topicName, &ListSubscriptionsRuntimePropertiesOptions{
		MaxPageSize: 2,
	})
	all := map[string]SubscriptionRuntimePropertiesItem{}
	times := 0

	for pager.More() {
		times++
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)

		require.LessOrEqual(t, len(page.SubscriptionRuntimeProperties), 2)

		for _, subItem := range page.SubscriptionRuntimeProperties {
			// _, exists := all[subItem.SubscriptionName]
			// require.False(t, exists, fmt.Sprintf("Each subscription result should be unique but found more than one of '%s'", subItem.SubscriptionName))
			all[subItem.SubscriptionName] = subItem

			require.False(t, subItem.CreatedAt.IsZero())
			require.False(t, subItem.UpdatedAt.IsZero())
			require.False(t, subItem.AccessedAt.IsZero())
			require.Zero(t, subItem.ActiveMessageCount)
			require.Zero(t, subItem.DeadLetterMessageCount)
			require.Zero(t, subItem.TotalMessageCount)
		}
	}

	require.GreaterOrEqual(t, times, 2)

	// sanity check - the topics we created exist and their deserialization is
	// working.
	for _, expectedSub := range expectedSubs {
		props, exists := all[expectedSub]
		require.True(t, exists)
		require.NotEqualValues(t, time.Time{}, props.CreatedAt)
	}
}

func TestAdminClient_UpdateSubscription(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	_, err := adminClient.CreateTopic(context.Background(), topicName, nil)
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	subscriptionName := fmt.Sprintf("sub-%X", time.Now().UnixNano())
	createResp, err := adminClient.CreateSubscription(context.Background(), topicName, subscriptionName, nil)
	require.NoError(t, err)
	require.Equal(t, topicName, createResp.TopicName)
	require.Equal(t, subscriptionName, createResp.SubscriptionName)

	defer deleteSubscription(t, adminClient, topicName, subscriptionName)

	createResp.LockDuration = to.Ptr("PT4M")
	updateResp, err := adminClient.UpdateSubscription(context.Background(), topicName, subscriptionName, createResp.SubscriptionProperties, nil)
	require.NoError(t, err)

	require.Equal(t, subscriptionName, updateResp.SubscriptionName)
	require.Equal(t, topicName, updateResp.TopicName)
	require.EqualValues(t, "PT4M", *updateResp.LockDuration)

	// try changing a value that's not allowed
	updateResp.RequiresSession = to.Ptr(true)
	updateResp, err = adminClient.UpdateSubscription(context.Background(), topicName, subscriptionName, updateResp.SubscriptionProperties, nil)
	require.Contains(t, err.Error(), "The value for the RequiresSession property of an existing Subscription cannot be changed")

	updateResp, err = adminClient.UpdateSubscription(context.Background(), topicName, "non-existent-subscription", createResp.SubscriptionProperties, nil)
	require.Contains(t, err.Error(), "404 Not Found")

	var asResponseErr *azcore.ResponseError
	require.ErrorAs(t, err, &asResponseErr)
	require.EqualValues(t, 404, asResponseErr.StatusCode)
}

func TestAdminClient_LackPermissions_Queue(t *testing.T) {
	testData := setupLowPrivTest(t)
	defer testData.Cleanup()

	ctx := context.Background()

	queue, err := testData.Client.GetQueue(ctx, "not-found-queue", nil)
	require.NoError(t, err)
	require.Nil(t, queue)

	runtimeProps, err := testData.Client.GetQueueRuntimeProperties(ctx, "not-found-queue", nil)
	require.NoError(t, err)
	require.Nil(t, runtimeProps)

	var re *azcore.ResponseError

	_, err = testData.Client.GetQueue(ctx, testData.QueueName, nil)
	require.Contains(t, err.Error(), "Manage,EntityRead claims required for this operation")
	require.ErrorAs(t, err, &re)
	require.EqualValues(t, 401, re.StatusCode)

	pager := testData.Client.NewListQueuesPager(nil)
	page, err := pager.NextPage(context.Background())
	require.Empty(t, page.Queues)
	require.Contains(t, err.Error(), "Manage,EntityRead claims required for this operation")
	require.ErrorAs(t, err, &re)
	require.EqualValues(t, 401, re.StatusCode)

	_, err = testData.Client.CreateQueue(ctx, "canneverbecreated", nil)
	require.Contains(t, err.Error(), "Authorization failed for specified action: Manage,EntityWrite")
	require.ErrorAs(t, err, &re)
	require.EqualValues(t, 401, re.StatusCode)

	_, err = testData.Client.UpdateQueue(ctx, "canneverbecreated", QueueProperties{}, nil)
	require.Contains(t, err.Error(), "Authorization failed for specified action: Manage,EntityWrite")
	require.ErrorAs(t, err, &re)
	require.EqualValues(t, 401, re.StatusCode)

	_, err = testData.Client.DeleteQueue(ctx, testData.QueueName, nil)
	require.Contains(t, err.Error(), "Authorization failed for specified action: Manage,EntityDelete.")
	require.ErrorAs(t, err, &re)
	require.EqualValues(t, 401, re.StatusCode)
}

func TestAdminClient_LackPermissions_Topic(t *testing.T) {
	testData := setupLowPrivTest(t)
	defer testData.Cleanup()

	ctx := context.Background()

	resp, err := testData.Client.GetTopic(ctx, "not-found-topic", nil)
	require.NoError(t, err)
	require.Nil(t, resp)

	resprt, err := testData.Client.GetTopicRuntimeProperties(ctx, "not-found-topic", nil)
	require.NoError(t, err)
	require.Nil(t, resprt)

	var asResponseErr *azcore.ResponseError

	_, err = testData.Client.GetTopic(ctx, testData.TopicName, nil)
	require.Contains(t, err.Error(), ">Manage,EntityRead claims required for this operation")
	require.ErrorAs(t, err, &asResponseErr)
	require.EqualValues(t, 401, asResponseErr.StatusCode)

	pager := testData.Client.NewListTopicsPager(nil)
	_, err = pager.NextPage(context.Background())
	require.Contains(t, err.Error(), ">Manage,EntityRead claims required for this operation")
	require.ErrorAs(t, err, &asResponseErr)
	require.EqualValues(t, 401, asResponseErr.StatusCode)

	_, err = testData.Client.CreateTopic(ctx, "canneverbecreated", nil)
	require.Contains(t, err.Error(), "Authorization failed for specified action")
	require.ErrorAs(t, err, &asResponseErr)
	require.EqualValues(t, 401, asResponseErr.StatusCode)

	_, err = testData.Client.UpdateTopic(ctx, "canneverbecreated", TopicProperties{}, nil)
	require.Contains(t, err.Error(), "Authorization failed for specified action")
	require.ErrorAs(t, err, &asResponseErr)
	require.EqualValues(t, 401, asResponseErr.StatusCode)

	_, err = testData.Client.DeleteTopic(ctx, testData.TopicName, nil)
	require.Contains(t, err.Error(), "Authorization failed for specified action: Manage,EntityDelete.")
	require.ErrorAs(t, err, &asResponseErr)
	require.EqualValues(t, 401, asResponseErr.StatusCode)
}

// TestAdminClient_NotFound makes sure that the "GET as LIST" behavior maps to nil when we are trying
// to get an entity by name and it is not found.
func TestAdminClient_NotFound(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)

	queue, err := adminClient.GetQueue(context.Background(), "non-existent-queue", nil)
	require.NoError(t, err)
	require.Nil(t, queue)

	queuert, err := adminClient.GetQueueRuntimeProperties(context.Background(), "non-existent-queue", nil)
	require.NoError(t, err)
	require.Nil(t, queuert)

	topic, err := adminClient.GetTopic(context.Background(), "non-existent-topic", nil)
	require.NoError(t, err)
	require.Nil(t, topic)

	topicrt, err := adminClient.GetTopicRuntimeProperties(context.Background(), "non-existent-topic", nil)
	require.NoError(t, err)
	require.Nil(t, topicrt)

	sub, err := adminClient.GetSubscription(context.Background(), "non-existent-topic", "sub1", nil)
	require.NoError(t, err)
	require.Nil(t, sub)

	subrt, err := adminClient.GetSubscriptionRuntimeProperties(context.Background(), "non-existent-topic", "sub1", nil)
	require.NoError(t, err)
	require.Nil(t, subrt)

	nanoSeconds := time.Now().UnixNano()
	topicName := fmt.Sprintf("topic-%d", nanoSeconds)

	_, err = adminClient.CreateTopic(context.Background(), topicName, nil)
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	sub, err = adminClient.GetSubscription(context.Background(), topicName, "sub1", nil)
	require.NoError(t, err)
	require.Nil(t, sub)

	subrt, err = adminClient.GetSubscriptionRuntimeProperties(context.Background(), topicName, "sub1", nil)
	require.NoError(t, err)
	require.Nil(t, subrt)
}

func TestAdminClient_LackPermissions_Subscription(t *testing.T) {
	testData := setupLowPrivTest(t)
	defer testData.Cleanup()

	ctx := context.Background()

	_, err := testData.Client.GetSubscription(ctx, testData.TopicName, "not-found-sub", nil)
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'GetSubscription'")

	_, err = testData.Client.GetSubscription(ctx, testData.TopicName, testData.SubName, nil)
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'GetSubscription'")

	pager := testData.Client.NewListSubscriptionsPager(testData.TopicName, nil)
	_, err = pager.NextPage(context.Background())
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'EnumerateSubscriptions' operation")

	_, err = testData.Client.CreateSubscription(ctx, testData.TopicName, "canneverbecreated", nil)
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'CreateOrUpdateSubscription'")

	_, err = testData.Client.UpdateSubscription(ctx, testData.TopicName, "canneverbecreated", SubscriptionProperties{}, nil)
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'CreateOrUpdateSubscription'")

	_, err = testData.Client.DeleteSubscription(ctx, testData.TopicName, testData.SubName, nil)
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'DeleteSubscription'")
}

type fakeEM struct {
	atom.EntityManager

	getResponses []string
}

func (em *fakeEM) Get(ctx context.Context, entityPath string, respObj any) (*http.Response, error) {
	jsonPath := em.getResponses[0]
	em.getResponses = em.getResponses[1:]

	reader, err := os.Open(jsonPath)

	if err != nil {
		return nil, err
	}

	defer reader.Close()

	bytes, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, respObj); err != nil {
		return nil, err
	}

	return &http.Response{
		Body: http.NoBody,
	}, nil
}

func TestAdminClient_ruleWithDifferentXMLNamespaces(t *testing.T) {
	adminClient := &Client{
		em: &fakeEM{
			getResponses: []string{
				"testdata/rulefeed.json",
				"testdata/emptyrulefeed.json",
			},
		},
	}

	// we're stubbed out at this point, so no need to use "real" entities
	pager := adminClient.NewListRulesPager("test", "test", nil)

	require.True(t, pager.More())
	listRulesResp, err := pager.NextPage(context.Background())
	require.NoError(t, err)
	require.NotNil(t, listRulesResp)

	expectedTime, err := time.Parse(time.RFC3339, "2020-01-01T01:02:03Z")
	require.NoError(t, err)

	correlationFilter := listRulesResp.Rules[0].Filter.(*CorrelationFilter)
	require.Equal(t, map[string]any{
		"hello":         "world",
		"hellodatetime": expectedTime,
		"hellodouble":   1.1,
		"helloint":      int64(101),
	}, correlationFilter.ApplicationProperties)
	require.Equal(t, "11", *correlationFilter.CorrelationID)

	sqlFilter := listRulesResp.Rules[1].Filter.(*SQLFilter)
	require.Equal(t, "1=1", sqlFilter.Expression)

	require.True(t, pager.More())
	listRulesResp, err = pager.NextPage(context.Background())
	require.NoError(t, err)
	require.Empty(t, listRulesResp.Rules)

	require.False(t, pager.More())
}

func TestAdminClient_CreateRules(t *testing.T) {
	adminClient, topicName := createTestSub(t)
	defer func() {
		_, err := adminClient.DeleteTopic(context.Background(), topicName, nil)
		require.NoError(t, err)
	}()

	t.Run("ruleThatDoesNotExist", func(t *testing.T) {
		getResp, err := adminClient.GetRule(context.Background(), topicName, "sub", "non-existent-rule", nil)
		require.NoError(t, err)
		require.Nil(t, getResp)
	})

	// (simple all in one - create the rule, update the rule, get the rule and make sure it
	// all does the right stuff.
	assertRuleCRUD := func(t *testing.T, rp RuleProperties) {
		defer func() {
			resp, err := adminClient.DeleteRule(context.Background(), topicName, "sub", rp.Name, nil)
			require.NoError(t, err)
			require.NotNil(t, resp)
		}()

		createdRule, err := adminClient.CreateRule(context.Background(), topicName, "sub", &CreateRuleOptions{
			Name:   &rp.Name,
			Filter: rp.Filter,
			Action: rp.Action,
		})
		require.NoError(t, err, fmt.Sprintf("Created rule %s", rp.Name))

		if rp.Filter == nil {
			// Service Bus will automatically add in a 'TrueFilter' to our
			// rule. We'll add it to our local copy just for assert purposes.
			rp.Filter = &TrueFilter{}
		}

		require.Equal(t, createdRule, CreateRuleResponse{
			RuleProperties: rp,
		}, fmt.Sprintf("Created rule %s matches our rule", rp.Name))

		updateResp, err := adminClient.UpdateRule(context.Background(), topicName, "sub", createdRule.RuleProperties)
		require.NoError(t, err, fmt.Sprintf("Updated rule %s succeeds", rp.Name))

		require.Equal(t, updateResp, UpdateRuleResponse{
			RuleProperties: rp,
		}, fmt.Sprintf("Updated rule %s matches our rule", rp.Name))

		getResp, err := adminClient.GetRule(context.Background(), topicName, "sub", rp.Name, nil)
		require.NoError(t, err, fmt.Sprintf("Get rule %s succeeds", rp.Name))

		require.Equal(t, getResp, &GetRuleResponse{
			RuleProperties: updateResp.RuleProperties,
		}, fmt.Sprintf("Get rule %s matches our rule", rp.Name))
	}

	t.Run("ruleWithNoActionOrFilter", func(t *testing.T) {
		assertRuleCRUD(t, RuleProperties{
			Name: "ruleWithNoActionOrFilter",
		})
	})

	t.Run("ruleWithFalseFilter", func(t *testing.T) {
		assertRuleCRUD(t, RuleProperties{
			Name:   "ruleWithFalseFilter",
			Filter: &FalseFilter{},
		})
	})

	t.Run("ruleWithTrueFilter", func(t *testing.T) {
		assertRuleCRUD(t, RuleProperties{
			Name:   "ruleWithTrueFilter",
			Filter: &FalseFilter{},
		})
	})

	t.Run("ruleWithSQLFilterNoParams", func(t *testing.T) {
		assertRuleCRUD(t, RuleProperties{
			Name: "ruleWithSQLFilterNoParams",
			Filter: &SQLFilter{
				Expression: "MessageID='hello'",
			},
		})
	})

	t.Run("ruleWithSQLFilterWithParams", func(t *testing.T) {
		dt, err := time.Parse(time.RFC3339, "2001-01-01T01:02:03Z")
		require.NoError(t, err)

		assertRuleCRUD(t, RuleProperties{
			Name: "ruleWithSQLFilterWithParams",
			Filter: &SQLFilter{
				Expression: "MessageID=@stringVar OR MessageID=@intVar OR MessageID=@floatVar OR MessageID=@dateTimeVar OR MessageID=@boolVar",
				Parameters: map[string]any{
					"@stringVar":   "hello world",
					"@intVar":      int64(100),
					"@floatVar":    float64(100.1),
					"@dateTimeVar": dt,
					"@boolVar":     true,
				},
			},
		})
	})

	t.Run("ruleWithCorrelationFilter", func(t *testing.T) {
		assertRuleCRUD(t, RuleProperties{
			Name: "ruleWithCorrelationFilter",
			Filter: &CorrelationFilter{
				ContentType:      to.Ptr("application/xml"),
				CorrelationID:    to.Ptr("correlationID"),
				MessageID:        to.Ptr("messageID"),
				ReplyTo:          to.Ptr("replyTo"),
				ReplyToSessionID: to.Ptr("replyToSessionID"),
				SessionID:        to.Ptr("sessionID"),
				Subject:          to.Ptr("subject"),
				To:               to.Ptr("to"),
				ApplicationProperties: map[string]any{
					"CustomProp1": "hello",
				},
			},
		})
	})

	t.Run("ruleWithAction", func(t *testing.T) {
		dt, err := time.Parse(time.RFC3339, "2001-01-01T01:02:03Z")
		require.NoError(t, err)

		assertRuleCRUD(t, RuleProperties{
			Name: "ruleWithAction",
			Action: &SQLAction{
				Expression: "SET MessageID=@stringVar SET MessageID=@intVar SET MessageID=@floatVar SET MessageID=@dateTimeVar SET MessageID=@boolVar",
				Parameters: map[string]any{
					"@stringVar":   "hello world",
					"@intVar":      int64(100),
					"@floatVar":    float64(100.1),
					"@dateTimeVar": dt,
					"@boolVar":     true,
				},
			},
		})
	})

	t.Run("ruleWithFilterAndAction", func(t *testing.T) {
		dt, err := time.Parse(time.RFC3339, "2001-01-01T01:02:03Z")
		require.NoError(t, err)

		assertRuleCRUD(t, RuleProperties{
			Name: "ruleWithFilterAndAction",
			Filter: &SQLFilter{
				Expression: "MessageID=@stringVar OR MessageID=@intVar OR MessageID=@floatVar OR MessageID=@dateTimeVar OR MessageID=@boolVar",
				Parameters: map[string]any{
					"@stringVar":   "hello world",
					"@intVar":      int64(100),
					"@floatVar":    float64(100.1),
					"@dateTimeVar": dt,
					"@boolVar":     true,
				},
			},
			Action: &SQLAction{
				Expression: "SET MessageID=@stringVar SET MessageID=@intVar SET MessageID=@floatVar SET MessageID=@dateTimeVar SET MessageID=@boolVar",
				Parameters: map[string]any{
					"@stringVar":   "hello world",
					"@intVar":      int64(100),
					"@floatVar":    float64(100.1),
					"@dateTimeVar": dt,
					"@boolVar":     true,
				},
			},
		})
	})
}

func TestAdminClient_ListRulesWithOnlyDefault(t *testing.T) {
	adminClient, topicName := createTestSub(t)
	defer func() {
		_, err := adminClient.DeleteTopic(context.Background(), topicName, nil)
		require.NoError(t, err)
	}()

	rulesPager := adminClient.NewListRulesPager(topicName, "sub", nil)
	require.True(t, rulesPager.More())
	resp, err := rulesPager.NextPage(context.Background())
	require.NoError(t, err)

	require.Equal(t, []RuleProperties{
		{Name: "$Default", Filter: &TrueFilter{}},
	}, resp.Rules)

	// documenting this behavior - we let the service dictate the
	// default page size so we don't know (yet) if there are any more results
	// remaining.
	require.True(t, rulesPager.More())

	resp, err = rulesPager.NextPage(context.Background())
	require.NoError(t, err)
	require.Empty(t, resp)

	require.False(t, rulesPager.More())
}

func TestAdminClient_ListRules_MaxPageSize(t *testing.T) {
	adminClient, topicName := createTestSub(t)
	defer func() {
		_, err := adminClient.DeleteTopic(context.Background(), topicName, nil)
		require.NoError(t, err)
	}()

	for _, rule := range []string{"rule1", "rule2", "rule3"} {
		_, err := adminClient.CreateRule(context.Background(), topicName, "sub", &CreateRuleOptions{
			Name: to.Ptr(rule),
			Filter: &SQLFilter{
				Expression: fmt.Sprintf("MessageID=%s", rule),
			},
		})
		require.NoError(t, err)

		defer func(rule string) {
			_, err := adminClient.DeleteRule(context.Background(), topicName, "sub", rule, nil)
			require.NoError(t, err)
		}(rule)
	}

	rulesPager := adminClient.NewListRulesPager(topicName, "sub", &ListRulesOptions{
		// there are actually 4 rules on the subscription right now - the 3 I just added
		// _and_ the $Default rule, which was auto-generated when the subscription
		// was created.
		MaxPageSize: 3,
	})

	var all []RuleProperties

	// first page
	require.True(t, rulesPager.More())
	resp, err := rulesPager.NextPage(context.Background())
	require.NoError(t, err)
	require.Equal(t, 3, len(resp.Rules))

	all = append(all, resp.Rules...)

	// second page
	require.True(t, rulesPager.More())
	resp, err = rulesPager.NextPage(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, len(resp.Rules))

	// since we explicitly configured a page size we know that this one came back
	// light, so we know this is also the last page.
	require.False(t, rulesPager.More())

	all = append(all, resp.Rules...)

	sort.Slice(all, func(i, j int) bool {
		return strings.Compare(all[i].Name, all[j].Name) < 0
	})

	require.Equal(t, []RuleProperties{
		{
			Name:   "$Default",
			Filter: &TrueFilter{},
		},
		{
			Name: "rule1",
			Filter: &SQLFilter{
				Expression: "MessageID=rule1",
			},
		},
		{
			Name: "rule2",
			Filter: &SQLFilter{
				Expression: "MessageID=rule2",
			},
		},
		{
			Name: "rule3",
			Filter: &SQLFilter{
				Expression: "MessageID=rule3",
			},
		},
	}, all)
}

func TestAdminClient_GetDefaultRule(t *testing.T) {
	adminClient, topicName := createTestSub(t)
	defer func() {
		_, err := adminClient.DeleteTopic(context.Background(), topicName, nil)
		require.NoError(t, err)
	}()

	getResp, err := adminClient.GetRule(context.Background(), topicName, "sub", "$Default", nil)
	require.NoError(t, err)

	// by default a subscription has a filter that lets every
	// message through (ie, the TrueFilter)
	require.Equal(t, getResp, &GetRuleResponse{
		RuleProperties: RuleProperties{
			Name:   "$Default",
			Filter: &TrueFilter{},
		},
	})

	// switch to a filter that _rejects_ every message instead
	getResp.RuleProperties.Filter = &FalseFilter{}

	updateRuleResp, err := adminClient.UpdateRule(context.Background(), topicName, "sub", getResp.RuleProperties)
	require.NoError(t, err)

	require.Equal(t, updateRuleResp.RuleProperties, getResp.RuleProperties)
}

type emwrap struct {
	inner atom.EntityManager
}

func (em *emwrap) Put(ctx context.Context, entityPath string, body any, respObj any, options *atom.ExecuteOptions) (*http.Response, error) {
	resp, err := em.inner.Put(ctx, entityPath, body, respObj, options)

	if err != nil {
		return resp, err
	}

	em.makeFilterAndActionUnknown(respObj)
	return resp, err
}

func (em *emwrap) Delete(ctx context.Context, entityPath string) (*http.Response, error) {
	return em.inner.Delete(ctx, entityPath)
}
func (em *emwrap) TokenProvider() auth.TokenProvider { return em.inner.TokenProvider() }

func (em *emwrap) Get(ctx context.Context, entityPath string, respObj any) (*http.Response, error) {
	resp, err := em.inner.Get(ctx, entityPath, respObj)

	if err != nil {
		return resp, err
	}

	em.makeFilterAndActionUnknown(respObj)
	return resp, err
}

func (*emwrap) makeFilterAndActionUnknown(respObj any) {
	actual, ok := respObj.(**atom.RuleEnvelope)

	if !ok {
		return
	}

	f := (*actual).Content.RuleDescription.Filter
	f.Type = "PurposefullyChangedFilterType_" + f.Type

	a := (*actual).Content.RuleDescription.Action
	a.Type = "PurposefullyChangedActionType_" + a.Type
}

func TestAdminClient_UnknownFilterRoundtrippingWorks(t *testing.T) {
	// NOTE: This test is a little weird - we basically override all "known" type handling for filters and
	// actions and force them to go through our "unknown" filter handling.
	//
	// This allows the service to potentially upgrade in the future without breaking older clients. They get a
	// relatively primitive object but they won't accidentally delete or slice filters when doing updates.
	//
	// Also, if they're willing to deserialize the XML themselves they can interact with filters until they
	// update their azservicebus dependency.

	adminClient, topicName := createTestSub(t)
	defer func() {
		_, err := adminClient.DeleteTopic(context.Background(), topicName, nil)
		require.NoError(t, err)
	}()

	dt, err := time.Parse(time.RFC3339, "2001-01-01T01:02:03Z")
	require.NoError(t, err)

	rp := RuleProperties{
		Name: "ruleWithFilterAndAction",
		Filter: &SQLFilter{
			Expression: "MessageID=@stringVar OR MessageID=@intVar OR MessageID=@floatVar OR MessageID=@dateTimeVar OR MessageID=@boolVar",
			Parameters: map[string]any{
				"@stringVar":   "hello world",
				"@intVar":      int64(100),
				"@floatVar":    float64(100.1),
				"@dateTimeVar": dt,
				"@boolVar":     true,
			},
		},
		Action: &SQLAction{
			Expression: "SET MessageID=@stringVar SET MessageID=@intVar SET MessageID=@floatVar SET MessageID=@dateTimeVar SET MessageID=@boolVar",
			Parameters: map[string]any{
				"@stringVar":   "hello world",
				"@intVar":      int64(100),
				"@floatVar":    float64(100.1),
				"@dateTimeVar": dt,
				"@boolVar":     true,
			},
		},
	}

	origEM := adminClient.em

	adminClient.em = &emwrap{
		inner: origEM,
	}

	createdRule, err := adminClient.CreateRule(context.Background(), topicName, "sub", &CreateRuleOptions{
		Name:   &rp.Name,
		Filter: rp.Filter,
		Action: rp.Action,
	})
	require.NoError(t, err, fmt.Sprintf("Created rule %s", rp.Name))

	urf := createdRule.Filter.(*UnknownRuleFilter)
	require.Regexp(t, "^<Filter.*", string(urf.RawXML))
	require.Regexp(t, "PurposefullyChangedFilterType_SqlFilter", urf.Type)

	ura := createdRule.Action.(*UnknownRuleAction)
	require.Regexp(t, "^<Action.*", string(ura.RawXML))
	require.Regexp(t, "PurposefullyChangedActionType_SqlRuleAction", ura.Type)

	// now we'll change them to what they actually should be. We're still testing the
	// UnknownRule(Action|Filter) logic on PUT, we just did it in a round-a-bout way.
	urf.RawXML = bytes.Replace(urf.RawXML, []byte("PurposefullyChangedFilterType_SqlFilter"), []byte("SqlFilter"), 1)
	ura.RawXML = bytes.Replace(ura.RawXML, []byte("PurposefullyChangedActionType_SqlRuleAction"), []byte("SqlRuleAction"), 1)

	_, err = adminClient.UpdateRule(context.Background(), topicName, "sub", createdRule.RuleProperties)
	require.NoError(t, err, fmt.Sprintf("Updated rule %s succeeds", rp.Name))

	adminClient.em = origEM

	getResp, err := adminClient.GetRule(context.Background(), topicName, "sub", rp.Name, nil)
	require.NoError(t, err, fmt.Sprintf("Get rule %s succeeds", rp.Name))

	require.Equal(t, getResp, &GetRuleResponse{
		RuleProperties: rp,
	}, fmt.Sprintf("Get rule %s matches our rule", rp.Name))
}

func createTestSub(t *testing.T) (*Client, string) {
	adminClient := newAdminClientForTest(t, nil)

	topicName := fmt.Sprintf("rule-topic-%X", time.Now().UnixNano())
	_, err := adminClient.CreateTopic(context.Background(), topicName, &CreateTopicOptions{
		Properties: &TopicProperties{
			AutoDeleteOnIdle: to.Ptr("PT5M"),
		},
	})
	require.NoError(t, err)

	_, err = adminClient.CreateSubscription(context.Background(), topicName, "sub", nil)
	require.NoError(t, err)

	return adminClient, topicName
}

type entityManagerForPagerTests struct {
	atom.EntityManager
	getPaths []string
}

func (em *entityManagerForPagerTests) Get(ctx context.Context, entityPath string, respObj any) (*http.Response, error) {
	em.getPaths = append(em.getPaths, entityPath)

	switch feedPtrPtr := respObj.(type) {
	case **atom.TopicFeed:
		*feedPtrPtr = &atom.TopicFeed{
			Entries: []atom.TopicEnvelope{
				{},
				{},
				{},
			},
		}
	default:
		panic(fmt.Sprintf("Unknown feed type: %T", respObj))
	}

	return &http.Response{}, nil
}

func TestAdminClient_pagerWithLightPage(t *testing.T) {
	adminClient, err := NewClientFromConnectionString("Endpoint=sb://fakeendpoint.something/;SharedAccessKeyName=fakekeyname;SharedAccessKey=CHANGEME", nil) // allowed connection string
	require.NoError(t, err)

	em := &entityManagerForPagerTests{}
	adminClient.em = em

	pager := adminClient.newPagerFunc("/$Resources/Topics", 10, func(pv any) int {
		// note that we're returning fewer results than the max page size
		// in ATOM < max page size means this is the last page of results.
		return 3
	})

	var feed *atom.TopicFeed
	resp, err := pager(context.Background(), &feed)

	// first page should be good
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, feed)

	require.EqualValues(t, []string{
		"/$Resources/Topics?&$top=10",
	}, em.getPaths)

	// but from this point on it's considered EOF since the
	// previous page of results was "light"
	feed = nil

	resp, err = pager(context.Background(), &feed)
	require.Nil(t, resp)
	require.Nil(t, err)
	require.Nil(t, feed)

	// no new calls will be made to the service
	require.EqualValues(t, []string{
		"/$Resources/Topics?&$top=10",
	}, em.getPaths)
}

func TestAdminClient_pagerWithFullPage(t *testing.T) {
	adminClient, err := NewClientFromConnectionString("Endpoint=sb://fakeendpoint.something/;SharedAccessKeyName=fakekeyname;SharedAccessKey=CHANGEME", nil) // allowed connection string
	require.NoError(t, err)

	em := &entityManagerForPagerTests{}
	adminClient.em = em

	// first request - got 10 results back, not EOF
	simulatedPageSize := 10

	pager := adminClient.newPagerFunc("/$Resources/Topics", 10, func(pv any) int {
		return simulatedPageSize
	})

	var feed *atom.TopicFeed
	simulatedPageSize = 10 // pretend the first page was full
	resp, err := pager(context.Background(), &feed)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, feed)

	require.EqualValues(t, []string{
		"/$Resources/Topics?&$top=10",
	}, em.getPaths)

	// grabbing the next page now, we'll also get 10 results back (still not EOF)
	simulatedPageSize = 10
	feed = nil

	simulatedPageSize = 10 // pretend the first page was full
	resp, err = pager(context.Background(), &feed)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, feed)

	require.EqualValues(t, []string{
		"/$Resources/Topics?&$top=10",
		"/$Resources/Topics?&$top=10&$skip=10",
	}, em.getPaths)

	// we'll return zero results for this next request, which should stop this pager.
	simulatedPageSize = 0
	feed = nil

	// nil, nil across the board indicates there wasn't an error, but we're
	// definitely done _now_, rather than having to check with another request.
	resp, err = pager(context.Background(), &feed)
	require.Nil(t, resp)
	require.Nil(t, err)

	// no new calls will be made to the service
	require.EqualValues(t, []string{
		"/$Resources/Topics?&$top=10",
		"/$Resources/Topics?&$top=10&$skip=10",
		"/$Resources/Topics?&$top=10&$skip=20",
	}, em.getPaths)

	// no more requests will be made, the pager is shut down after EOF.
	simulatedPageSize = 10
	feed = nil

	resp, err = pager(context.Background(), &feed)
	require.Nil(t, resp)
	require.Nil(t, err)

	// no new calls will be made to the service
	require.EqualValues(t, []string{
		"/$Resources/Topics?&$top=10",
		"/$Resources/Topics?&$top=10&$skip=10",
		"/$Resources/Topics?&$top=10&$skip=20",
	}, em.getPaths)
}

func TestAdminClient_unknownActionSerde(t *testing.T) {
	ura, err := newUnknownRuleActionFromActionDescription(&atom.ActionDescription{
		Type:   "SomeNewAction",
		RawXML: []byte("<someNewActionXML></someNewActionXML>"),
		RawAttrs: []xml.Attr{
			{
				Name: xml.Name{
					Local: "some-custom-attribute",
				},
				Value: "some-custom-attribute-value",
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t,
		`<Action some-custom-attribute="some-custom-attribute-value" i:type="SomeNewAction">`+
			`<someNewActionXML></someNewActionXML>`+
			`</Action>`, string(ura.RawXML))
	require.Equal(t, "SomeNewAction", ura.Type)

	// and now the inverse
	ad, err := convertUnknownRuleActionToActionDescription(ura)
	require.NoError(t, err)

	require.Equal(t, "<someNewActionXML></someNewActionXML>", string(ad.RawXML))

	require.Equal(t, []xml.Attr{
		{Name: xml.Name{Local: "some-custom-attribute"}, Value: "some-custom-attribute-value"}}, ad.RawAttrs)
	require.Equal(t, "SomeNewAction", string(ad.Type))

	_, err = convertUnknownRuleActionToActionDescription(&UnknownRuleAction{
		Type:   "something",
		RawXML: []byte("invalid &xml"),
	})
	require.Error(t, err)
}

func TestAdminClient_unknownFilterSerde(t *testing.T) {
	urf, err := newUnknownRuleFilterFromFilterDescription(&atom.FilterDescription{
		Type:   "SomeNewFilter",
		RawXML: []byte("<someNewFilterXML></someNewFilterXML>"),
		RawAttrs: []xml.Attr{
			{
				Name: xml.Name{
					Local: "some-custom-attribute",
				},
				Value: "some-custom-attribute-value",
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t,
		`<Filter some-custom-attribute="some-custom-attribute-value" i:type="SomeNewFilter">`+
			`<someNewFilterXML></someNewFilterXML>`+
			`</Filter>`, string(urf.RawXML))
	require.Equal(t, "SomeNewFilter", urf.Type)

	// and now the inverse
	ad, err := convertUnknownRuleFilterToFilterDescription(urf)
	require.NoError(t, err)

	require.Equal(t, "<someNewFilterXML></someNewFilterXML>", string(ad.RawXML))

	require.Equal(t, []xml.Attr{
		{Name: xml.Name{Local: "some-custom-attribute"}, Value: "some-custom-attribute-value"}}, ad.RawAttrs)
	require.Equal(t, "SomeNewFilter", string(ad.Type))

	_, err = convertUnknownRuleFilterToFilterDescription(&UnknownRuleFilter{
		Type:   "something",
		RawXML: []byte("invalid &xml"),
	})
	require.Error(t, err)
}

func deleteQueue(t *testing.T, ac *Client, queueName string) {
	_, err := ac.DeleteQueue(context.Background(), queueName, nil)
	require.NoError(t, err)
}

func deleteTopic(t *testing.T, ac *Client, topicName string) {
	_, err := ac.DeleteTopic(context.Background(), topicName, nil)
	require.NoError(t, err)
}

func deleteSubscription(t *testing.T, ac *Client, topicName string, subscriptionName string) {
	_, err := ac.DeleteSubscription(context.Background(), topicName, subscriptionName, nil)
	require.NoError(t, err)
}

func setupLowPrivTest(t *testing.T) *struct {
	Client    *Client
	TopicName string
	SubName   string
	QueueName string
	Cleanup   func()
} {
	adminClient := newAdminClientForTest(t, nil)

	lowPrivAdminClient, err := NewClientFromConnectionString(test.MustGetEnvVar(t, test.EnvKeyConnectionStringNoManage), nil) // allowed connection string
	require.NoError(t, err)

	nanoSeconds := time.Now().UnixNano()

	topicName := fmt.Sprintf("topic-%d", nanoSeconds)
	queueName := fmt.Sprintf("queue-%d", nanoSeconds)
	subName := "subscription1"

	// create some entities that we need (there's a diff between something not being
	// found and something failing because of lack of authorization)
	cleanup := func() func() {
		_, err = adminClient.CreateQueue(context.Background(), queueName, nil)
		require.NoError(t, err)

		_, err = adminClient.CreateTopic(context.Background(), topicName, nil)
		require.NoError(t, err)

		_, err = adminClient.CreateSubscription(context.Background(), topicName, subName, nil)
		require.NoError(t, err)

		// _, err = sm.PutRule(context.Background(), subName, ruleName, TrueFilter{})
		// require.NoError(t, err)

		return func() {
			deleteTopic(t, adminClient, topicName) // will also delete the subscription
			deleteQueue(t, adminClient, queueName)
		}
	}()

	return &struct {
		Client    *Client
		TopicName string
		SubName   string
		QueueName string
		Cleanup   func()
	}{
		Client:    lowPrivAdminClient,
		QueueName: queueName,
		TopicName: topicName,
		SubName:   subName,
		Cleanup:   cleanup,
	}
}

func createRandomKeyForSB(t *testing.T) string {
	tempPassword := make([]byte, 22)
	n, err := cryptoRand.Read(tempPassword)
	require.NoError(t, err)
	require.Equal(t, cap(tempPassword), n)

	return hex.EncodeToString(tempPassword)
}

func createAuthorizationRulesForTest(t *testing.T) []AuthorizationRule {
	primary := createRandomKeyForSB(t)
	secondary := createRandomKeyForSB(t)

	return []AuthorizationRule{
		{
			AccessRights: []AccessRight{AccessRightSend},
			KeyName:      to.Ptr("keyName1"),
			PrimaryKey:   &primary,
			SecondaryKey: &secondary,
		},
	}
}

func TestATOMNoCountDetails(t *testing.T) {
	subRP, err := newSubscriptionRuntimePropertiesItem(&atom.SubscriptionEnvelope{Content: &atom.SubscriptionContent{}}, "topic")
	require.Nil(t, subRP)
	require.Error(t, err, "invalid subscription runtime properties: no CountDetails element")

	qRP, err := newQueueRuntimePropertiesItem(&atom.QueueEnvelope{Content: &atom.QueueContent{}})
	require.Nil(t, qRP)
	require.Error(t, err, "invalid queue runtime properties: no CountDetails element")

	tRP, err := newTopicRuntimePropertiesItem(&atom.TopicEnvelope{Content: &atom.TopicContent{}})
	require.Nil(t, tRP)
	require.Error(t, err, "invalid topic runtime properties: no CountDetails element")
}

func TestATOMEntitiesHaveNames(t *testing.T) {
	adminClient := newAdminClientForTest(t, nil)

	nano := time.Now().UnixNano()
	topicName := fmt.Sprintf("topic-%X", nano)

	t.Run("topic", func(t *testing.T) {
		createTopicResp, err := adminClient.CreateTopic(context.Background(), topicName, &CreateTopicOptions{
			Properties: &TopicProperties{
				AutoDeleteOnIdle: to.Ptr("PT5M"),
			},
		})
		require.NoError(t, err)
		require.Equal(t, createTopicResp.TopicName, topicName)

		topicResp, err := adminClient.GetTopic(context.Background(), topicName, nil)
		require.NoError(t, err)
		require.Equal(t, topicResp.TopicName, topicName)
	})

	t.Run("sub", func(t *testing.T) {
		createSubResp, err := adminClient.CreateSubscription(context.Background(), topicName, "sub1", &CreateSubscriptionOptions{
			Properties: &SubscriptionProperties{
				AutoDeleteOnIdle: to.Ptr("PT5M"),
			},
		})
		require.NoError(t, err)
		require.Equal(t, createSubResp.TopicName, topicName)
		require.Equal(t, createSubResp.SubscriptionName, "sub1")

		subResp, err := adminClient.GetSubscription(context.Background(), topicName, "sub1", nil)
		require.NoError(t, err)

		require.Equal(t, subResp.TopicName, topicName)
		require.Equal(t, subResp.SubscriptionName, "sub1")
	})
}

func newAdminClientForTest(t *testing.T, options *test.NewClientOptions[ClientOptions]) *Client {
	return test.NewClient(t, test.NewClientArgs[ClientOptions, Client]{
		NewClientFromConnectionString: NewClientFromConnectionString, // allowed connection string
		NewClient:                     NewClient,
	}, options)
}
