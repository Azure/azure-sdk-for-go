// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
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
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestAdminClient_UsingIdentity(t *testing.T) {
	// test with azure identity support
	ns := os.Getenv("SERVICEBUS_ENDPOINT")
	envCred, err := azidentity.NewEnvironmentCredential(nil)

	if err != nil || ns == "" {
		t.Skip("Azure Identity compatible credentials not configured")
	}

	adminClient, err := NewAdminClient(ns, envCred, nil)
	require.NoError(t, err)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	props, err := adminClient.CreateQueue(context.Background(), queueName, nil, nil)
	require.NoError(t, err)
	require.EqualValues(t, 10, *props.MaxDeliveryCount)

	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), queueName)
		require.NoError(t, err)
	}()
}

func TestAdminClient_QueueWithMaxValues(t *testing.T) {
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
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

	resp, err := adminClient.GetQueue(context.Background(), queueName)
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
}

func TestAdminClient_CreateQueue(t *testing.T) {
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
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
		_, err := adminClient.DeleteQueue(context.Background(), queueName)
		require.NoError(t, err)
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

	getResp, err := adminClient.GetQueue(context.Background(), queueName)
	require.NoError(t, err)

	require.EqualValues(t, getResp.QueueProperties, createResp.QueueProperties)
}

func TestAdminClient_Queue_Forwarding(t *testing.T) {
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	forwardToQueueName := fmt.Sprintf("queue-fwd-%X", time.Now().UnixNano())

	_, err = adminClient.CreateQueue(context.Background(), forwardToQueueName, nil, nil)
	require.NoError(t, err)

	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), forwardToQueueName)
		require.NoError(t, err)
	}()

	formatted := fmt.Sprintf("%s%s", adminClient.em.Host, forwardToQueueName)

	createResp, err := adminClient.CreateQueue(context.Background(), queueName, &QueueProperties{
		ForwardTo:                     &formatted,
		ForwardDeadLetteredMessagesTo: &formatted,
	}, nil)

	require.NoError(t, err)
	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), queueName)
		require.NoError(t, err)
	}()

	require.EqualValues(t, formatted, *createResp.ForwardTo)
	require.EqualValues(t, formatted, *createResp.ForwardDeadLetteredMessagesTo)

	getResp, err := adminClient.GetQueue(context.Background(), queueName)

	require.NoError(t, err)
	require.EqualValues(t, createResp.QueueProperties, getResp.QueueProperties)

	client, err := NewClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("this message will be auto-forwarded"),
	})
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(forwardToQueueName, nil)
	require.NoError(t, err)

	forwardedMessage, err := receiver.receiveMessage(context.Background(), nil)
	require.NoError(t, err)

	body, err := forwardedMessage.Body()
	require.NoError(t, err)
	require.EqualValues(t, "this message will be auto-forwarded", string(body))
}

func TestAdminClient_UpdateQueue(t *testing.T) {
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	createdProps, err := adminClient.CreateQueue(context.Background(), queueName, nil, nil)
	require.NoError(t, err)

	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), queueName)
		require.NoError(t, err)
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
	require.Contains(t, err.Error(), "error code: 404")
	require.Nil(t, updatedProps)
}

func TestAdminClient_GetQueueRuntimeProperties(t *testing.T) {
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	client, err := NewClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)
	defer client.Close(context.Background())

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	_, err = adminClient.CreateQueue(context.Background(), queueName, nil, nil)
	require.NoError(t, err)

	defer deleteQueue(t, adminClient, queueName)

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte("hello"),
		})
		require.NoError(t, err)
	}

	sequenceNumbers, err := sender.ScheduleMessages(context.Background(), []*Message{
		{Body: []byte("hello")},
	}, time.Now().Add(2*time.Hour))
	require.NoError(t, err)
	require.NotEmpty(t, sequenceNumbers)

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 2, nil)
	require.NoError(t, err)

	require.NoError(t, receiver.DeadLetterMessage(context.Background(), messages[0], nil))

	props, err := adminClient.GetQueueRuntimeProperties(context.Background(), queueName)
	require.NoError(t, err)

	require.EqualValues(t, 4, props.TotalMessageCount)

	require.EqualValues(t, 2, props.ActiveMessageCount)
	require.EqualValues(t, 1, props.DeadLetterMessageCount)
	require.EqualValues(t, 1, props.ScheduledMessageCount)
	require.EqualValues(t, 0, props.TransferDeadLetterMessageCount)
	require.EqualValues(t, 0, props.TransferMessageCount)

	require.Greater(t, props.SizeInBytes, int64(0))

	require.NotEqual(t, time.Time{}, props.CreatedAt)
	require.NotEqual(t, time.Time{}, props.UpdatedAt)
	require.NotEqual(t, time.Time{}, props.AccessedAt)
}

func TestAdminClient_ListQueues(t *testing.T) {
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	var expectedQueues []string
	now := time.Now().UnixNano()

	for i := 0; i < 3; i++ {
		queueName := strings.ToLower(fmt.Sprintf("queue-%d-%X", i, now))
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
			_, exists := all[props.QueueName]
			require.False(t, exists, "Each queue result should be unique")
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
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	var expectedQueues []string
	now := time.Now().UnixNano()

	for i := 0; i < 3; i++ {
		queueName := strings.ToLower(fmt.Sprintf("queue-%d-%X", i, now))
		expectedQueues = append(expectedQueues, queueName)

		_, err = adminClient.CreateQueue(context.Background(), queueName, &QueueProperties{
			MaxDeliveryCount: to.Int32Ptr(int32(i + 10)),
		}, nil)
		require.NoError(t, err)

		defer deleteQueue(t, adminClient, queueName)
	}

	// we skipped the first queue so it shouldn't come back in the results.
	pager := adminClient.ListQueuesRuntimeProperties(nil)
	all := map[string]*QueueRuntimePropertiesItem{}

	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()

		for _, queueRuntimeItem := range page.Items {
			_, exists := all[queueRuntimeItem.QueueName]
			require.False(t, exists, "Each queue result should be unique")
			all[queueRuntimeItem.QueueName] = queueRuntimeItem
		}
	}

	require.NoError(t, pager.Err())

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
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
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
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
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
	require.Contains(t, err.Error(), "error code: 404")
	require.Nil(t, updateResp)
}

func TestAdminClient_TopicAndSubscriptionRuntimeProperties(t *testing.T) {
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	client, err := NewClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	subscriptionName := fmt.Sprintf("sub-%X", time.Now().UnixNano())

	_, err = adminClient.CreateTopic(context.Background(), topicName, nil, nil)
	require.NoError(t, err)

	addSubResp, err := adminClient.CreateSubscription(context.Background(), topicName, subscriptionName, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, addSubResp)
	require.EqualValues(t, 10, *addSubResp.MaxDeliveryCount)

	defer deleteSubscription(t, adminClient, topicName, subscriptionName)

	sender, err := client.NewSender(topicName, nil)
	require.NoError(t, err)

	// trigger some stats

	//  Scheduled messages are accounted for in the topic stats.
	_, err = sender.ScheduleMessages(context.Background(), []*Message{
		{Body: []byte("hello")},
	}, time.Now().Add(2*time.Hour))
	require.NoError(t, err)

	// validate the topic runtime properties
	getRuntimeResp, err := adminClient.GetTopicRuntimeProperties(context.Background(), topicName)
	require.NoError(t, err)

	require.EqualValues(t, topicName, getRuntimeResp.Value.Name)

	require.EqualValues(t, 1, getRuntimeResp.Value.SubscriptionCount)
	require.NotEqual(t, time.Time{}, getRuntimeResp.Value.CreatedAt)
	require.NotEqual(t, time.Time{}, getRuntimeResp.Value.UpdatedAt)
	require.NotEqual(t, time.Time{}, getRuntimeResp.Value.AccessedAt)

	require.Greater(t, getRuntimeResp.Value.SizeInBytes, int64(0))
	require.EqualValues(t, int32(1), getRuntimeResp.Value.ScheduledMessageCount)

	// validate subscription runtime properties
	getSubResp, err := adminClient.GetSubscriptionRuntimeProperties(context.Background(), topicName, subscriptionName)
	require.NoError(t, err)

	require.EqualValues(t, 0, getSubResp.ActiveMessageCount)
	require.NotEqual(t, time.Time{}, getSubResp.CreatedAt)
	require.NotEqual(t, time.Time{}, getSubResp.UpdatedAt)
	require.NotEqual(t, time.Time{}, getSubResp.AccessedAt)
}

func TestAdminClient_ListTopics(t *testing.T) {
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	var expectedTopics []string
	now := time.Now().UnixNano()

	wg := sync.WaitGroup{}

	for i := 0; i < 3; i++ {
		topicName := strings.ToLower(fmt.Sprintf("topic-%d-%X", i, now))
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
	pager := adminClient.ListTopics(nil)
	all := map[string]*TopicItem{}

	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()

		for _, topicItem := range page.Items {
			_, exists := all[topicItem.TopicName]
			require.False(t, exists, "Each topic result should be unique")
			all[topicItem.TopicName] = topicItem
		}
	}

	require.NoError(t, pager.Err())

	// sanity check - the topics we created exist and their deserialization is
	// working.
	for i, expectedTopic := range expectedTopics {
		props, exists := all[expectedTopic]
		require.True(t, exists)
		require.EqualValues(t, time.Duration(i+1)*time.Minute, *props.DefaultMessageTimeToLive)
	}
}

func TestAdminClient_ListTopicsRuntimeProperties(t *testing.T) {
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	var expectedTopics []string
	now := time.Now().UnixNano()

	for i := 0; i < 3; i++ {
		topicName := strings.ToLower(fmt.Sprintf("topic-%d-%X", i, now))
		expectedTopics = append(expectedTopics, topicName)

		_, err = adminClient.CreateTopic(context.Background(), topicName, nil, nil)
		require.NoError(t, err)

		defer deleteTopic(t, adminClient, topicName)
	}

	// we skipped the first topic so it shouldn't come back in the results.
	pager := adminClient.ListTopicsRuntimeProperties(nil)
	all := map[string]*TopicRuntimeProperties{}

	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()

		for _, props := range page.Value {
			_, exists := all[props.Name]
			require.False(t, exists, "Each topic result should be unique")
			all[props.Name] = props
		}
	}

	require.NoError(t, pager.Err())

	// sanity check - the topics we created exist and their deserialization is
	// working.
	for _, expectedTopic := range expectedTopics {
		props, exists := all[expectedTopic]
		require.True(t, exists)
		require.NotEqualValues(t, time.Time{}, props.CreatedAt)
	}
}

func TestAdminClient_ListSubscriptions(t *testing.T) {
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	now := time.Now().UnixNano()
	topicName := strings.ToLower(fmt.Sprintf("topic-%X", now))

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
	pager := adminClient.ListSubscriptions(topicName, nil)
	all := map[string]*SubscriptionPropertiesItem{}

	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()

		for _, item := range page.Items {
			_, exists := all[item.SubscriptionName]
			require.False(t, exists, "Each subscription result should be unique")
			all[item.SubscriptionName] = item
		}
	}

	require.NoError(t, pager.Err())

	// sanity check - the subscriptions we created exist and their deserialization is
	// working.
	for i, expectedTopic := range expectedSubscriptions {
		props, exists := all[expectedTopic]
		require.True(t, exists)
		require.EqualValues(t, time.Duration(i+1)*time.Minute, *props.DefaultMessageTimeToLive)
	}
}

func TestAdminClient_ListSubscriptionRuntimeProperties(t *testing.T) {
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	now := time.Now().UnixNano()
	topicName := strings.ToLower(fmt.Sprintf("topic-%X", now))

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
	pager := adminClient.ListSubscriptionsRuntimeProperties(topicName, nil)
	all := map[string]*SubscriptionRuntimePropertiesItem{}

	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()

		for _, subItem := range page.Items {
			_, exists := all[subItem.SubscriptionName]
			require.False(t, exists, "Each subscription result should be unique")
			all[subItem.SubscriptionName] = subItem
		}
	}

	require.NoError(t, pager.Err())

	// sanity check - the topics we created exist and their deserialization is
	// working.
	for _, expectedSub := range expectedSubs {
		props, exists := all[expectedSub]
		require.True(t, exists)
		require.NotEqualValues(t, time.Time{}, props.CreatedAt)
	}
}

func TestAdminClient_UpdateSubscription(t *testing.T) {
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
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
	// a little awkward, we'll make these programatically inspectable as we add in better error handling.
	require.Contains(t, err.Error(), "error code: 404")
	require.Nil(t, updateResp)
}

func setupLowPrivTest(t *testing.T) *struct {
	Client    *AdminClient
	TopicName string
	SubName   string
	QueueName string
	Cleanup   func()
} {
	adminClient, err := NewAdminClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	lowPrivAdminClient, err := NewAdminClientFromConnectionString(getConnectionStringWithoutManagePerms(t), nil)
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
		Client    *AdminClient
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

func TestAdminClient_LackPermissions_Queue(t *testing.T) {
	testData := setupLowPrivTest(t)
	defer testData.Cleanup()

	ctx := context.Background()

	_, err := testData.Client.GetQueue(ctx, "not-found-queue")
	notFound, resp := atom.NotFound(err)
	require.True(t, notFound)
	require.NotNil(t, resp)

	_, err = testData.Client.GetQueue(ctx, testData.QueueName)
	require.Contains(t, err.Error(), "error code: 401, Details: Manage,EntityRead claims")

	pager := testData.Client.ListQueues(nil)
	require.False(t, pager.NextPage(context.Background()))
	require.Contains(t, pager.Err().Error(), "error code: 401, Details: Manage,EntityRead claims required for this operation")

	_, err = testData.Client.CreateQueue(ctx, "canneverbecreated", nil, nil)
	require.Contains(t, err.Error(), "error code: 401, Details: Authorization failed for specified action: Manage,EntityWrite")

	_, err = testData.Client.UpdateQueue(ctx, "canneverbecreated", QueueProperties{}, nil)
	require.Contains(t, err.Error(), "error code: 401, Details: Authorization failed for specified action: Manage,EntityWrite")

	_, err = testData.Client.DeleteQueue(ctx, testData.QueueName)
	require.Contains(t, err.Error(), "error code: 401, Details: Authorization failed for specified action: Manage,EntityDelete.")
}

func TestAdminClient_LackPermissions_Topic(t *testing.T) {
	testData := setupLowPrivTest(t)
	defer testData.Cleanup()

	ctx := context.Background()

	_, err := testData.Client.GetTopic(ctx, "not-found-topic", nil)
	notFound, resp := atom.NotFound(err)
	require.True(t, notFound)
	require.NotNil(t, resp)

	_, err = testData.Client.GetTopic(ctx, testData.TopicName, nil)
	require.Contains(t, err.Error(), "error code: 401, Details: Manage,EntityRead claims")

	pager := testData.Client.ListTopics(nil)
	require.False(t, pager.NextPage(context.Background()))
	require.Contains(t, pager.Err().Error(), "error code: 401, Details: Manage,EntityRead claims required for this operation")

	_, err = testData.Client.CreateTopic(ctx, "canneverbecreated", nil, nil)
	require.Contains(t, err.Error(), "error code: 401, Details: Authorization failed for specified action")

	_, err = testData.Client.UpdateTopic(ctx, "canneverbecreated", TopicProperties{}, nil)
	require.Contains(t, err.Error(), "error code: 401, Details: Authorization failed for specified action")

	_, err = testData.Client.DeleteTopic(ctx, testData.TopicName)
	require.Contains(t, err.Error(), "error code: 401, Details: Authorization failed for specified action: Manage,EntityDelete.")

	// sanity check that the http response is getting bundled into these errors, should it be needed.
	var httpResponse azcore.HTTPResponse
	require.True(t, errors.As(err, &httpResponse))
	require.EqualValues(t, http.StatusUnauthorized, httpResponse.RawResponse().StatusCode)
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

func toDurationPtr(d time.Duration) *time.Duration {
	return &d
}

func deleteQueue(t *testing.T, ac *AdminClient, queueName string) {
	_, err := ac.DeleteQueue(context.Background(), queueName)
	require.NoError(t, err)
}

func deleteTopic(t *testing.T, ac *AdminClient, topicName string) {
	_, err := ac.DeleteTopic(context.Background(), topicName)
	require.NoError(t, err)
}

func deleteSubscription(t *testing.T, ac *AdminClient, topicName string, subscriptionName string) {
	_, err := ac.DeleteSubscription(context.Background(), topicName, subscriptionName, nil)
	require.NoError(t, err)
}
