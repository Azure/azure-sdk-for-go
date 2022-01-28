// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestAdminClient_UsingIdentity(t *testing.T) {
	// test with azure identity support
	ns := os.Getenv("SERVICEBUS_ENDPOINT")
	cred, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil || ns == "" {
		t.Skip("Azure Identity compatible credentials not configured")
	}

	adminClient, err := NewClient(ns, cred, nil)
	require.NoError(t, err)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	props, err := adminClient.CreateQueue(context.Background(), queueName, nil, nil)
	require.NoError(t, err)
	require.EqualValues(t, 10, *props.MaxDeliveryCount)

	defer func() {
		deleteQueue(t, adminClient, queueName)
	}()
}

func TestAdminClient_GetNamespaceProperties(t *testing.T) {
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)
	resp, err := adminClient.GetNamespaceProperties(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.RawResponse)

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
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	es := EntityStatusReceiveDisabled

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())

	_, err = adminClient.CreateQueue(context.Background(), queueName, &QueueProperties{
		LockDuration: toDurationPtr(45 * time.Second),
		// when you enable partitioning Service Bus will automatically create 16 partitions, each with the size
		// of MaxSizeInMegabytes. This means when we retrieve this queue we'll get 16*4096 as the size (ie: 64GB)
		EnablePartitioning:                  to.BoolPtr(true),
		MaxSizeInMegabytes:                  to.Int32Ptr(4096),
		RequiresDuplicateDetection:          to.BoolPtr(true),
		RequiresSession:                     to.BoolPtr(true),
		DefaultMessageTimeToLive:            toDurationPtr(time.Duration(1<<63 - 1)),
		DeadLetteringOnMessageExpiration:    to.BoolPtr(true),
		DuplicateDetectionHistoryTimeWindow: toDurationPtr(4 * time.Hour),
		MaxDeliveryCount:                    to.Int32Ptr(100),
		EnableBatchedOperations:             to.BoolPtr(false),
		Status:                              &es,
		AutoDeleteOnIdle:                    toDurationPtr(time.Duration(1<<63 - 1)),
		UserMetadata:                        to.StringPtr("some metadata"),
	}, nil)
	require.NoError(t, err)

	defer deleteQueue(t, adminClient, queueName)

	resp, err := adminClient.GetQueue(context.Background(), queueName, nil)
	require.NoError(t, err)

	require.EqualValues(t, QueueProperties{
		LockDuration: toDurationPtr(45 * time.Second),
		// ie: this response was from a partitioned queue so the size is the original max size * # of partitions
		EnablePartitioning:                  to.BoolPtr(true),
		MaxSizeInMegabytes:                  to.Int32Ptr(16 * 4096),
		RequiresDuplicateDetection:          to.BoolPtr(true),
		RequiresSession:                     to.BoolPtr(true),
		DefaultMessageTimeToLive:            toDurationPtr(time.Duration(1<<63 - 1)),
		DeadLetteringOnMessageExpiration:    to.BoolPtr(true),
		DuplicateDetectionHistoryTimeWindow: toDurationPtr(4 * time.Hour),
		MaxDeliveryCount:                    to.Int32Ptr(100),
		EnableBatchedOperations:             to.BoolPtr(false),
		Status:                              &es,
		AutoDeleteOnIdle:                    toDurationPtr(time.Duration(1<<63 - 1)),
		UserMetadata:                        to.StringPtr("some metadata"),
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

func TestAdminClient_CreateQueue(t *testing.T) {
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())

	es := EntityStatusReceiveDisabled
	createResp, err := adminClient.CreateQueue(context.Background(), queueName, &QueueProperties{
		LockDuration: toDurationPtr(45 * time.Second),
		// when you enable partitioning Service Bus will automatically create 16 partitions, each with the size
		// of MaxSizeInMegabytes. This means when we retrieve this queue we'll get 16*4096 as the size (ie: 64GB)
		EnablePartitioning:                  to.BoolPtr(true),
		MaxSizeInMegabytes:                  to.Int32Ptr(4096),
		RequiresDuplicateDetection:          to.BoolPtr(true),
		RequiresSession:                     to.BoolPtr(true),
		DefaultMessageTimeToLive:            toDurationPtr(time.Hour * 6),
		DeadLetteringOnMessageExpiration:    to.BoolPtr(true),
		DuplicateDetectionHistoryTimeWindow: toDurationPtr(4 * time.Hour),
		MaxDeliveryCount:                    to.Int32Ptr(100),
		EnableBatchedOperations:             to.BoolPtr(false),
		Status:                              &es,
		AutoDeleteOnIdle:                    toDurationPtr(10 * time.Minute),
	}, nil)
	require.NoError(t, err)

	defer func() {
		deleteQueue(t, adminClient, queueName)
	}()

	require.EqualValues(t, QueueProperties{
		LockDuration: toDurationPtr(45 * time.Second),
		// ie: this response was from a partitioned queue so the size is the original max size * # of partitions
		EnablePartitioning:                  to.BoolPtr(true),
		MaxSizeInMegabytes:                  to.Int32Ptr(16 * 4096),
		RequiresDuplicateDetection:          to.BoolPtr(true),
		RequiresSession:                     to.BoolPtr(true),
		DefaultMessageTimeToLive:            toDurationPtr(time.Hour * 6),
		DeadLetteringOnMessageExpiration:    to.BoolPtr(true),
		DuplicateDetectionHistoryTimeWindow: toDurationPtr(4 * time.Hour),
		MaxDeliveryCount:                    to.Int32Ptr(100),
		EnableBatchedOperations:             to.BoolPtr(false),
		Status:                              &es,
		AutoDeleteOnIdle:                    toDurationPtr(10 * time.Minute),
	}, createResp.QueueProperties)

	getResp, err := adminClient.GetQueue(context.Background(), queueName, nil)
	require.NoError(t, err)

	require.EqualValues(t, getResp.QueueProperties, createResp.QueueProperties)
}

func TestAdminClient_UpdateQueue(t *testing.T) {
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	createdProps, err := adminClient.CreateQueue(context.Background(), queueName, nil, nil)
	require.NoError(t, err)

	defer func() {
		deleteQueue(t, adminClient, queueName)
	}()

	createdProps.MaxDeliveryCount = to.Int32Ptr(101)
	updatedProps, err := adminClient.UpdateQueue(context.Background(), queueName, createdProps.QueueProperties, nil)
	require.NoError(t, err)

	require.EqualValues(t, 101, *updatedProps.MaxDeliveryCount)

	// try changing a value that's not allowed
	updatedProps.RequiresSession = to.BoolPtr(true)
	updatedProps, err = adminClient.UpdateQueue(context.Background(), queueName, updatedProps.QueueProperties, nil)
	require.Contains(t, err.Error(), "The value for the RequiresSession property of an existing Queue cannot be changed")
	require.Nil(t, updatedProps)

	updatedProps, err = adminClient.UpdateQueue(context.Background(), "non-existent-queue", createdProps.QueueProperties, nil)
	// a little awkward, we'll make these programatically inspectable as we add in better error handling.
	require.Contains(t, err.Error(), "404 Not Found")

	var asResponseErr *azcore.ResponseError
	require.ErrorAs(t, err, &asResponseErr)
	require.EqualValues(t, 404, asResponseErr.StatusCode)

	require.Nil(t, updatedProps)
}

func TestAdminClient_ListQueues(t *testing.T) {
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	var expectedQueues []string
	now := time.Now().UnixNano()

	for i := 0; i < 3; i++ {
		queueName := strings.ToLower(fmt.Sprintf("list-queues-%d-%X", i, now))
		expectedQueues = append(expectedQueues, queueName)

		_, err = adminClient.CreateQueue(context.Background(), queueName, &QueueProperties{
			MaxDeliveryCount: to.Int32Ptr(int32(i + 10)),
		}, nil)
		require.NoError(t, err)

		defer deleteQueue(t, adminClient, queueName)
	}

	// we skipped the first queue so it shouldn't come back in the results.
	pager := adminClient.ListQueues(nil)
	all := map[string]*QueueItem{}

	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()

		for _, props := range page.Items {
			all[props.QueueName] = props
		}
	}

	require.NoError(t, pager.Err())

	// sanity check - the queues we created exist and their deserialization is
	// working.
	for i, expectedQueue := range expectedQueues {
		props, exists := all[expectedQueue]
		require.True(t, exists)
		require.EqualValues(t, i+10, *props.MaxDeliveryCount)
	}
}

func TestAdminClient_ListQueuesRuntimeProperties(t *testing.T) {
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	var expectedQueues []string
	now := time.Now().UnixNano()

	for i := 0; i < 3; i++ {
		queueName := strings.ToLower(fmt.Sprintf("list-queuert-%d-%X", i, now))
		expectedQueues = append(expectedQueues, queueName)

		_, err = adminClient.CreateQueue(context.Background(), queueName, &QueueProperties{
			MaxDeliveryCount: to.Int32Ptr(int32(i + 10)),
		}, nil)
		require.NoError(t, err)

		defer deleteQueue(t, adminClient, queueName)
	}

	// we skipped the first queue so it shouldn't come back in the results.
	pager := adminClient.ListQueuesRuntimeProperties(&ListQueuesRuntimePropertiesOptions{
		MaxPageSize: 2,
	})
	all := map[string]*QueueRuntimePropertiesItem{}

	times := 0

	for pager.NextPage(context.Background()) {
		times++
		page := pager.PageResponse()

		require.LessOrEqual(t, len(page.Items), 2)

		for _, queueRuntimeItem := range page.Items {
			// _, exists := all[queueRuntimeItem.QueueName]
			// require.False(t, exists, fmt.Sprintf("Each queue result should be unique but found more than one of '%s'", queueRuntimeItem.QueueName))
			all[queueRuntimeItem.QueueName] = queueRuntimeItem
		}
	}

	require.NoError(t, pager.Err())
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

func TestAdminClient_TopicAndSubscription(t *testing.T) {
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	subscriptionName := fmt.Sprintf("sub-%X", time.Now().UnixNano())
	forwardToQueueName := fmt.Sprintf("queue-fwd-%X", time.Now().UnixNano())

	_, err = adminClient.CreateQueue(context.Background(), forwardToQueueName, nil, nil)
	require.NoError(t, err)

	defer deleteQueue(t, adminClient, forwardToQueueName)

	status := EntityStatusActive

	// check topic properties, existence
	addResp, err := adminClient.CreateTopic(context.Background(), topicName, &TopicProperties{
		EnablePartitioning:                  to.BoolPtr(true),
		MaxSizeInMegabytes:                  to.Int32Ptr(2048),
		RequiresDuplicateDetection:          to.BoolPtr(true),
		DefaultMessageTimeToLive:            toDurationPtr(time.Minute * 3),
		DuplicateDetectionHistoryTimeWindow: toDurationPtr(time.Minute * 4),
		EnableBatchedOperations:             to.BoolPtr(true),
		Status:                              &status,
		AutoDeleteOnIdle:                    toDurationPtr(time.Minute * 7),
		SupportOrdering:                     to.BoolPtr(true),
		UserMetadata:                        to.StringPtr("user metadata"),
	}, nil)
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	require.EqualValues(t, TopicProperties{
		EnablePartitioning:                  to.BoolPtr(true),
		MaxSizeInMegabytes:                  to.Int32Ptr(16 * 2048), // enabling partitioning increases our max size because of the 16 partitions
		RequiresDuplicateDetection:          to.BoolPtr(true),
		DefaultMessageTimeToLive:            toDurationPtr(time.Minute * 3),
		DuplicateDetectionHistoryTimeWindow: toDurationPtr(time.Minute * 4),
		EnableBatchedOperations:             to.BoolPtr(true),
		Status:                              &status,
		AutoDeleteOnIdle:                    toDurationPtr(time.Minute * 7),
		SupportOrdering:                     to.BoolPtr(true),
		UserMetadata:                        to.StringPtr("user metadata"),
	}, addResp.TopicProperties)

	getResp, err := adminClient.GetTopic(context.Background(), topicName, nil)
	require.NoError(t, err)

	require.EqualValues(t, TopicProperties{
		EnablePartitioning:                  to.BoolPtr(true),
		MaxSizeInMegabytes:                  to.Int32Ptr(16 * 2048), // enabling partitioning increases our max size because of the 16 partitions
		RequiresDuplicateDetection:          to.BoolPtr(true),
		DefaultMessageTimeToLive:            toDurationPtr(time.Minute * 3),
		DuplicateDetectionHistoryTimeWindow: toDurationPtr(time.Minute * 4),
		EnableBatchedOperations:             to.BoolPtr(true),
		Status:                              &status,
		AutoDeleteOnIdle:                    toDurationPtr(time.Minute * 7),
		SupportOrdering:                     to.BoolPtr(true),
		UserMetadata:                        to.StringPtr("user metadata"),
	}, getResp.TopicProperties)

	runtimeResp, err := adminClient.GetTopicRuntimeProperties(context.Background(), topicName, nil)
	require.NoError(t, err)

	require.False(t, runtimeResp.CreatedAt.IsZero())
	require.False(t, runtimeResp.UpdatedAt.IsZero())
	require.True(t, runtimeResp.AccessedAt.IsZero())
	require.Zero(t, runtimeResp.SubscriptionCount)
	require.Zero(t, runtimeResp.ScheduledMessageCount)
	require.Zero(t, runtimeResp.SizeInBytes)

	addSubWithPropsResp, err := adminClient.CreateSubscription(context.Background(), topicName, subscriptionName, &SubscriptionProperties{
		LockDuration:                                    toDurationPtr(3 * time.Minute),
		RequiresSession:                                 to.BoolPtr(false),
		DefaultMessageTimeToLive:                        toDurationPtr(7 * time.Minute),
		DeadLetteringOnMessageExpiration:                to.BoolPtr(true),
		EnableDeadLetteringOnFilterEvaluationExceptions: to.BoolPtr(false),
		MaxDeliveryCount:                                to.Int32Ptr(11),
		Status:                                          &status,
		// ForwardTo:                     &forwardToQueueName,
		// ForwardDeadLetteredMessagesTo: &forwardToQueueName,
		EnableBatchedOperations: to.BoolPtr(false),
		AutoDeleteOnIdle:        toDurationPtr(11 * time.Minute),
		UserMetadata:            to.StringPtr("user metadata"),
	}, nil)
	require.NoError(t, err)

	defer deleteSubscription(t, adminClient, topicName, subscriptionName)

	require.EqualValues(t, SubscriptionProperties{
		LockDuration:                                    toDurationPtr(3 * time.Minute),
		RequiresSession:                                 to.BoolPtr(false),
		DefaultMessageTimeToLive:                        toDurationPtr(7 * time.Minute),
		DeadLetteringOnMessageExpiration:                to.BoolPtr(true),
		EnableDeadLetteringOnFilterEvaluationExceptions: to.BoolPtr(false),
		MaxDeliveryCount:                                to.Int32Ptr(11),
		Status:                                          &status,
		// ForwardTo:                     &forwardToQueueName,
		// ForwardDeadLetteredMessagesTo: &forwardToQueueName,
		EnableBatchedOperations: to.BoolPtr(false),
		AutoDeleteOnIdle:        toDurationPtr(11 * time.Minute),
		UserMetadata:            to.StringPtr("user metadata"),
	}, addSubWithPropsResp.CreateSubscriptionResult.SubscriptionProperties)
}

func TestAdminClient_UpdateTopic(t *testing.T) {
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	addResp, err := adminClient.CreateTopic(context.Background(), topicName, nil, nil)
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	addResp.AutoDeleteOnIdle = toDurationPtr(11 * time.Minute)
	updateResp, err := adminClient.UpdateTopic(context.Background(), topicName, addResp.TopicProperties, nil)
	require.NoError(t, err)

	require.EqualValues(t, 11*time.Minute, *updateResp.AutoDeleteOnIdle)

	// try changing a value that's not allowed
	updateResp.EnablePartitioning = to.BoolPtr(true)
	updateResp, err = adminClient.UpdateTopic(context.Background(), topicName, updateResp.TopicProperties, nil)
	require.Contains(t, err.Error(), "Partitioning cannot be changed for Topic. ")
	require.Nil(t, updateResp)

	updateResp, err = adminClient.UpdateTopic(context.Background(), "non-existent-topic", addResp.TopicProperties, nil)
	// a little awkward, we'll make these programatically inspectable as we add in better error handling.
	require.Contains(t, err.Error(), "404 Not Found")

	var asResponseErr *azcore.ResponseError
	require.ErrorAs(t, err, &asResponseErr)
	require.EqualValues(t, 404, asResponseErr.StatusCode)

	require.Nil(t, updateResp)
}

func TestAdminClient_ListTopics(t *testing.T) {
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	var expectedTopics []string
	now := time.Now().UnixNano()

	wg := sync.WaitGroup{}

	for i := 0; i < 3; i++ {
		topicName := strings.ToLower(fmt.Sprintf("list-topic-%d-%X", i, now))
		expectedTopics = append(expectedTopics, topicName)
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			_, err = adminClient.CreateTopic(context.Background(), topicName, &TopicProperties{
				DefaultMessageTimeToLive: toDurationPtr(time.Duration(i+1) * time.Minute),
			}, nil)
			require.NoError(t, err)
		}(i)

		defer deleteTopic(t, adminClient, topicName)
	}

	wg.Wait()

	// we skipped the first topic so it shouldn't come back in the results.
	pager := adminClient.ListTopics(&ListTopicsOptions{
		MaxPageSize: 2,
	})
	all := map[string]*TopicItem{}

	times := 0

	for pager.NextPage(context.Background()) {
		times++
		page := pager.PageResponse()

		require.LessOrEqual(t, len(page.Items), 2)

		for _, topicItem := range page.Items {
			// _, exists := all[topicItem.TopicName]
			// require.False(t, exists, fmt.Sprintf("Each topic result should be unique but found more than one of '%s'", topicItem.TopicName))
			all[topicItem.TopicName] = topicItem
		}
	}

	require.NoError(t, pager.Err())
	require.GreaterOrEqual(t, times, 2)

	// sanity check - the topics we created exist and their deserialization is
	// working.
	for i, expectedTopic := range expectedTopics {
		props, exists := all[expectedTopic]
		require.True(t, exists)
		require.EqualValues(t, time.Duration(i+1)*time.Minute, *props.DefaultMessageTimeToLive)
	}
}

func TestAdminClient_ListTopicsRuntimeProperties(t *testing.T) {
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	var expectedTopics []string
	now := time.Now().UnixNano()

	for i := 0; i < 3; i++ {
		topicName := strings.ToLower(fmt.Sprintf("list-topicrt-%d-%X", i, now))
		expectedTopics = append(expectedTopics, topicName)

		_, err = adminClient.CreateTopic(context.Background(), topicName, nil, nil)
		require.NoError(t, err)

		defer deleteTopic(t, adminClient, topicName)
	}

	times := 0

	// we skipped the first topic so it shouldn't come back in the results.
	pager := adminClient.ListTopicsRuntimeProperties(&ListTopicsRuntimePropertiesOptions{
		MaxPageSize: 2,
	})
	all := map[string]*TopicRuntimePropertiesItem{}

	for pager.NextPage(context.Background()) {
		times++
		page := pager.PageResponse()

		require.LessOrEqual(t, len(page.Items), 2)

		for _, item := range page.Items {
			// _, exists := all[item.TopicName]
			// require.False(t, exists, fmt.Sprintf("Each topic result should be unique but found more than one of '%s'", item.TopicName))
			all[item.TopicName] = item
		}
	}

	require.NoError(t, pager.Err())
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
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	now := time.Now().UnixNano()
	topicName := strings.ToLower(fmt.Sprintf("listsub-%X", now))

	_, err = adminClient.CreateTopic(context.Background(), topicName, nil, nil)
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	var expectedSubscriptions []string

	for i := 0; i < 3; i++ {
		subName := strings.ToLower(fmt.Sprintf("sub-%d-%X", i, now))
		expectedSubscriptions = append(expectedSubscriptions, subName)

		_, err = adminClient.CreateSubscription(context.Background(), topicName, subName, &SubscriptionProperties{
			DefaultMessageTimeToLive: toDurationPtr(time.Duration(i+1) * time.Minute),
		}, nil)
		require.NoError(t, err)

		defer deleteSubscription(t, adminClient, topicName, subName)
	}

	// we skipped the first topic so it shouldn't come back in the results.
	pager := adminClient.ListSubscriptions(topicName, &ListSubscriptionsOptions{
		MaxPageSize: 2,
	})
	all := map[string]*SubscriptionPropertiesItem{}

	times := 0

	for pager.NextPage(context.Background()) {
		times++
		page := pager.PageResponse()

		require.LessOrEqual(t, len(page.Items), 2)

		for _, item := range page.Items {
			// _, exists := all[item.SubscriptionName]
			// require.False(t, exists, fmt.Sprintf("Each subscription result should be unique but found more than one of '%s'", item.SubscriptionName))
			all[item.SubscriptionName] = item
		}
	}

	require.NoError(t, pager.Err())
	require.GreaterOrEqual(t, times, 2)

	// sanity check - the subscriptions we created exist and their deserialization is
	// working.
	for i, expectedTopic := range expectedSubscriptions {
		props, exists := all[expectedTopic]
		require.True(t, exists)
		require.EqualValues(t, time.Duration(i+1)*time.Minute, *props.DefaultMessageTimeToLive)
	}
}

func TestAdminClient_ListSubscriptionRuntimeProperties(t *testing.T) {
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	now := time.Now().UnixNano()
	topicName := strings.ToLower(fmt.Sprintf("listsubrt-%X", now))

	_, err = adminClient.CreateTopic(context.Background(), topicName, nil, nil)
	require.NoError(t, err)

	var expectedSubs []string

	for i := 0; i < 3; i++ {
		subscriptionName := strings.ToLower(fmt.Sprintf("sub-%d-%X", i, now))
		expectedSubs = append(expectedSubs, subscriptionName)

		_, err = adminClient.CreateSubscription(context.Background(), topicName, subscriptionName, nil, nil)
		require.NoError(t, err)

		defer deleteSubscription(t, adminClient, topicName, subscriptionName)
	}

	// we skipped the first subscription so it shouldn't come back in the results.
	pager := adminClient.ListSubscriptionsRuntimeProperties(topicName, &ListSubscriptionsRuntimePropertiesOptions{
		MaxPageSize: 2,
	})
	all := map[string]*SubscriptionRuntimePropertiesItem{}
	times := 0

	for pager.NextPage(context.Background()) {
		times++
		page := pager.PageResponse()

		require.LessOrEqual(t, len(page.Items), 2)

		for _, subItem := range page.Items {
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

	require.NoError(t, pager.Err())
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
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	_, err = adminClient.CreateTopic(context.Background(), topicName, nil, nil)
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	subscriptionName := fmt.Sprintf("sub-%X", time.Now().UnixNano())
	addResp, err := adminClient.CreateSubscription(context.Background(), topicName, subscriptionName, nil, nil)
	require.NoError(t, err)

	defer deleteSubscription(t, adminClient, topicName, subscriptionName)

	addResp.LockDuration = toDurationPtr(4 * time.Minute)
	updateResp, err := adminClient.UpdateSubscription(context.Background(), topicName, subscriptionName, addResp.SubscriptionProperties, nil)
	require.NoError(t, err)

	require.EqualValues(t, 4*time.Minute, *updateResp.LockDuration)

	// try changing a value that's not allowed
	updateResp.RequiresSession = to.BoolPtr(true)
	updateResp, err = adminClient.UpdateSubscription(context.Background(), topicName, subscriptionName, updateResp.SubscriptionProperties, nil)
	require.Contains(t, err.Error(), "The value for the RequiresSession property of an existing Subscription cannot be changed")
	require.Nil(t, updateResp)

	updateResp, err = adminClient.UpdateSubscription(context.Background(), topicName, "non-existent-subscription", addResp.CreateSubscriptionResult.SubscriptionProperties, nil)
	require.Contains(t, err.Error(), "404 Not Found")

	var asResponseErr *azcore.ResponseError
	require.ErrorAs(t, err, &asResponseErr)
	require.EqualValues(t, 404, asResponseErr.StatusCode)

	require.Nil(t, updateResp)
}

func TestAdminClient_LackPermissions_Queue(t *testing.T) {
	testData := setupLowPrivTest(t)
	defer testData.Cleanup()

	ctx := context.Background()

	_, err := testData.Client.GetQueue(ctx, "not-found-queue", nil)
	notFound, resp := atom.NotFound(err)
	require.True(t, notFound)
	require.NotNil(t, resp)

	var re *azcore.ResponseError

	_, err = testData.Client.GetQueue(ctx, testData.QueueName, nil)
	require.Contains(t, err.Error(), "Manage,EntityRead claims required for this operation")
	require.ErrorAs(t, err, &re)
	require.EqualValues(t, 401, re.StatusCode)

	pager := testData.Client.ListQueues(nil)
	require.False(t, pager.NextPage(context.Background()))
	require.Contains(t, pager.Err().Error(), "Manage,EntityRead claims required for this operation")
	require.ErrorAs(t, err, &re)
	require.EqualValues(t, 401, re.StatusCode)

	_, err = testData.Client.CreateQueue(ctx, "canneverbecreated", nil, nil)
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

	_, err := testData.Client.GetTopic(ctx, "not-found-topic", nil)
	notFound, resp := atom.NotFound(err)
	require.True(t, notFound)
	require.NotNil(t, resp)

	var asResponseErr *azcore.ResponseError

	_, err = testData.Client.GetTopic(ctx, testData.TopicName, nil)
	require.Contains(t, err.Error(), ">Manage,EntityRead claims required for this operation")
	require.ErrorAs(t, err, &asResponseErr)
	require.EqualValues(t, 401, asResponseErr.StatusCode)

	pager := testData.Client.ListTopics(nil)
	require.False(t, pager.NextPage(context.Background()))
	require.Contains(t, pager.Err().Error(), ">Manage,EntityRead claims required for this operation")
	require.ErrorAs(t, err, &asResponseErr)
	require.EqualValues(t, 401, asResponseErr.StatusCode)

	_, err = testData.Client.CreateTopic(ctx, "canneverbecreated", nil, nil)
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

func TestAdminClient_LackPermissions_Subscription(t *testing.T) {
	testData := setupLowPrivTest(t)
	defer testData.Cleanup()

	ctx := context.Background()

	_, err := testData.Client.GetSubscription(ctx, testData.TopicName, "not-found-sub", nil)
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'GetSubscription'")

	_, err = testData.Client.GetSubscription(ctx, testData.TopicName, testData.SubName, nil)
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'GetSubscription'")

	pager := testData.Client.ListSubscriptions(testData.TopicName, nil)
	require.False(t, pager.NextPage(context.Background()))
	require.Contains(t, pager.Err().Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'EnumerateSubscriptions' operation")

	_, err = testData.Client.CreateSubscription(ctx, testData.TopicName, "canneverbecreated", nil, nil)
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'CreateOrUpdateSubscription'")

	_, err = testData.Client.UpdateSubscription(ctx, testData.TopicName, "canneverbecreated", SubscriptionProperties{}, nil)
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'CreateOrUpdateSubscription'")

	_, err = testData.Client.DeleteSubscription(ctx, testData.TopicName, testData.SubName, nil)
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'DeleteSubscription'")
}

type entityManagerForPagerTests struct {
	atom.EntityManager
	getPaths []string
}

func (em *entityManagerForPagerTests) Get(ctx context.Context, entityPath string, respObj interface{}, mw ...atom.MiddlewareFunc) (*http.Response, error) {
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
	adminClient, err := NewClientFromConnectionString("Endpoint=sb://fakeendpoint.something/;SharedAccessKeyName=fakekeyname;SharedAccessKey=CHANGEME", nil)
	require.NoError(t, err)

	em := &entityManagerForPagerTests{}
	adminClient.em = em

	pager := adminClient.newPagerFunc("/$Resources/Topics", 10, func(pv interface{}) int {
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
	adminClient, err := NewClientFromConnectionString("Endpoint=sb://fakeendpoint.something/;SharedAccessKeyName=fakekeyname;SharedAccessKey=CHANGEME", nil)
	require.NoError(t, err)

	em := &entityManagerForPagerTests{}
	adminClient.em = em

	// first request - got 10 results back, not EOF
	simulatedPageSize := 10

	pager := adminClient.newPagerFunc("/$Resources/Topics", 10, func(pv interface{}) int {
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

func toDurationPtr(d time.Duration) *time.Duration {
	return &d
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
	adminClient, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	lowPrivAdminClient, err := NewClientFromConnectionString(test.GetConnectionStringWithoutManagePerms(t), nil)
	require.NoError(t, err)

	nanoSeconds := time.Now().UnixNano()

	topicName := fmt.Sprintf("topic-%d", nanoSeconds)
	queueName := fmt.Sprintf("queue-%d", nanoSeconds)
	subName := "subscription1"

	// TODO: add in rule management
	//ruleName := "rule"

	// create some entities that we need (there's a diff between something not being
	// found and something failing because of lack of authorization)
	cleanup := func() func() {
		_, err = adminClient.CreateQueue(context.Background(), queueName, nil, nil)
		require.NoError(t, err)

		_, err = adminClient.CreateTopic(context.Background(), topicName, nil, nil)
		require.NoError(t, err)

		_, err = adminClient.CreateSubscription(context.Background(), topicName, subName, nil, nil)
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
