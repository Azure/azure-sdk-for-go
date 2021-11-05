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
	props, err := adminClient.AddQueue(context.Background(), queueName)
	require.NoError(t, err)
	require.EqualValues(t, queueName, props.Value.Name)

	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), queueName)
		require.NoError(t, err)
	}()
}

func TestAdminClient_QueueWithMaxValues(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	es := EntityStatusReceiveDisabled

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())

	queueExists, err := adminClient.QueueExists(context.Background(), queueName)
	require.False(t, queueExists)
	require.NoError(t, err)

	_, err = adminClient.AddQueueWithProperties(context.Background(), &QueueProperties{
		Name:         queueName,
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
	})
	require.NoError(t, err)

	defer deleteQueue(t, adminClient, queueName)

	queueExists, err = adminClient.QueueExists(context.Background(), queueName)
	require.True(t, queueExists)
	require.NoError(t, err)

	resp, err := adminClient.GetQueue(context.Background(), queueName)
	require.NoError(t, err)

	require.EqualValues(t, &QueueProperties{
		Name:         queueName,
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
	}, resp.Value)
}

func TestAdminClient_AddQueue(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())

	es := EntityStatusReceiveDisabled
	createResp, err := adminClient.AddQueueWithProperties(context.Background(), &QueueProperties{
		Name:         queueName,
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
	})
	require.NoError(t, err)

	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), queueName)
		require.NoError(t, err)
	}()

	require.EqualValues(t, &QueueProperties{
		Name:         queueName,
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
	}, createResp.Value)

	getResp, err := adminClient.GetQueue(context.Background(), queueName)
	require.NoError(t, err)

	require.EqualValues(t, getResp.Value, createResp.Value)
}

func TestAdminClient_Queue_Forwarding(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	forwardToQueueName := fmt.Sprintf("queue-fwd-%X", time.Now().UnixNano())

	_, err = adminClient.AddQueue(context.Background(), forwardToQueueName)
	require.NoError(t, err)

	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), forwardToQueueName)
		require.NoError(t, err)
	}()

	formatted := fmt.Sprintf("%s%s", adminClient.em.Host, forwardToQueueName)

	createResp, err := adminClient.AddQueueWithProperties(context.Background(), &QueueProperties{
		Name:                          queueName,
		ForwardTo:                     &formatted,
		ForwardDeadLetteredMessagesTo: &formatted,
	})

	require.NoError(t, err)
	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), queueName)
		require.NoError(t, err)
	}()

	require.EqualValues(t, formatted, *createResp.Value.ForwardTo)
	require.EqualValues(t, formatted, *createResp.Value.ForwardDeadLetteredMessagesTo)

	getResp, err := adminClient.GetQueue(context.Background(), queueName)

	require.NoError(t, err)
	require.EqualValues(t, createResp.Value, getResp.Value)

	client, err := NewClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	sender, err := client.NewSender(queueName)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("this message will be auto-forwarded"),
	})
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(forwardToQueueName, nil)
	require.NoError(t, err)

	forwardedMessage, err := receiver.receiveMessage(context.Background(), nil)
	require.NoError(t, err)

	require.EqualValues(t, "this message will be auto-forwarded", string(forwardedMessage.Body))
}

func TestAdminClient_UpdateQueue(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	createdProps, err := adminClient.AddQueue(context.Background(), queueName)
	require.NoError(t, err)

	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), queueName)
		require.NoError(t, err)
	}()

	createdProps.Value.MaxDeliveryCount = to.Int32Ptr(101)
	updatedProps, err := adminClient.UpdateQueue(context.Background(), createdProps.Value)
	require.NoError(t, err)

	require.EqualValues(t, 101, *updatedProps.Value.MaxDeliveryCount)

	// try changing a value that's not allowed
	updatedProps.Value.RequiresSession = to.BoolPtr(true)
	updatedProps, err = adminClient.UpdateQueue(context.Background(), updatedProps.Value)
	require.Contains(t, err.Error(), "The value for the RequiresSession property of an existing Queue cannot be changed")
	require.Nil(t, updatedProps)

	createdProps.Value.Name = "non-existent-queue"
	updatedProps, err = adminClient.UpdateQueue(context.Background(), createdProps.Value)
	// a little awkward, we'll make these programatically inspectable as we add in better error handling.
	require.Contains(t, err.Error(), "error code: 404")
	require.Nil(t, updatedProps)
}

func TestAdminClient_GetQueueRuntimeProperties(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	client, err := NewClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)
	defer client.Close(context.Background())

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	_, err = adminClient.AddQueue(context.Background(), queueName)
	require.NoError(t, err)

	defer deleteQueue(t, adminClient, queueName)

	sender, err := client.NewSender(queueName)
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte("hello"),
		})
		require.NoError(t, err)
	}

	sequenceNumbers, err := sender.ScheduleMessages(context.Background(), []SendableMessage{
		&Message{Body: []byte("hello")},
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

	require.EqualValues(t, queueName, props.Value.Name)

	require.EqualValues(t, 4, props.Value.TotalMessageCount)

	require.EqualValues(t, 2, props.Value.ActiveMessageCount)
	require.EqualValues(t, 1, props.Value.DeadLetterMessageCount)
	require.EqualValues(t, 1, props.Value.ScheduledMessageCount)
	require.EqualValues(t, 0, props.Value.TransferDeadLetterMessageCount)
	require.EqualValues(t, 0, props.Value.TransferMessageCount)

	require.Greater(t, props.Value.SizeInBytes, int64(0))

	require.NotEqual(t, time.Time{}, props.Value.CreatedAt)
	require.NotEqual(t, time.Time{}, props.Value.UpdatedAt)
	require.NotEqual(t, time.Time{}, props.Value.AccessedAt)
}

func TestAdminClient_ListQueues(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	var expectedQueues []string
	now := time.Now().UnixNano()

	for i := 0; i < 3; i++ {
		queueName := strings.ToLower(fmt.Sprintf("queue-%d-%X", i, now))
		expectedQueues = append(expectedQueues, queueName)

		_, err = adminClient.AddQueueWithProperties(context.Background(), &QueueProperties{
			Name:             queueName,
			MaxDeliveryCount: to.Int32Ptr(int32(i + 10)),
		})
		require.NoError(t, err)

		defer deleteQueue(t, adminClient, queueName)
	}

	// we skipped the first queue so it shouldn't come back in the results.
	pager := adminClient.ListQueues(nil)
	all := map[string]*QueueProperties{}
	var allNames []string

	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()

		for _, props := range page.Value {
			_, exists := all[props.Name]
			require.False(t, exists, "Each queue result should be unique")
			all[props.Name] = props
			allNames = append(allNames, props.Name)
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

	// grab the second to last item
	pager = adminClient.ListQueues(&ListQueuesOptions{
		Skip: len(allNames) - 2,
		Top:  1,
	})

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, 1, len(pager.PageResponse().Value))
	require.EqualValues(t, allNames[len(allNames)-2], pager.PageResponse().Value[0].Name)
	require.NoError(t, pager.Err())

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, allNames[len(allNames)-1], pager.PageResponse().Value[0].Name)
	require.False(t, pager.NextPage(context.Background()))
}

func TestAdminClient_ListQueuesRuntimeProperties(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	var expectedQueues []string
	now := time.Now().UnixNano()

	for i := 0; i < 3; i++ {
		queueName := strings.ToLower(fmt.Sprintf("queue-%d-%X", i, now))
		expectedQueues = append(expectedQueues, queueName)

		_, err = adminClient.AddQueueWithProperties(context.Background(), &QueueProperties{
			Name:             queueName,
			MaxDeliveryCount: to.Int32Ptr(int32(i + 10)),
		})
		require.NoError(t, err)

		defer deleteQueue(t, adminClient, queueName)
	}

	// we skipped the first queue so it shouldn't come back in the results.
	pager := adminClient.ListQueuesRuntimeProperties(nil)
	all := map[string]*QueueRuntimeProperties{}
	var allNames []string

	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()

		for _, props := range page.Value {
			_, exists := all[props.Name]
			require.False(t, exists, "Each queue result should be unique")
			all[props.Name] = props
			allNames = append(allNames, props.Name)
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

	// grab the second to last item
	pager = adminClient.ListQueuesRuntimeProperties(&ListQueuesRuntimePropertiesOptions{
		Skip: len(allNames) - 2,
		Top:  1,
	})

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, 1, len(pager.PageResponse().Value))
	require.EqualValues(t, allNames[len(allNames)-2], pager.PageResponse().Value[0].Name)
	require.NoError(t, pager.Err())

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, allNames[len(allNames)-1], pager.PageResponse().Value[0].Name)
	require.False(t, pager.NextPage(context.Background()))
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
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	subscriptionName := fmt.Sprintf("sub-%X", time.Now().UnixNano())
	forwardToQueueName := fmt.Sprintf("queue-fwd-%X", time.Now().UnixNano())

	_, err = adminClient.AddQueue(context.Background(), forwardToQueueName)
	require.NoError(t, err)

	defer deleteQueue(t, adminClient, forwardToQueueName)

	status := EntityStatusActive

	// check topic properties, existence
	topicExists, err := adminClient.TopicExists(context.Background(), topicName)
	require.False(t, topicExists)
	require.NoError(t, err)

	addResp, err := adminClient.AddTopicWithProperties(context.Background(), &TopicProperties{
		Name:                                topicName,
		EnablePartitioning:                  to.BoolPtr(true),
		MaxSizeInMegabytes:                  to.Int32Ptr(2048),
		RequiresDuplicateDetection:          to.BoolPtr(true),
		DefaultMessageTimeToLive:            toDurationPtr(time.Minute * 3),
		DuplicateDetectionHistoryTimeWindow: toDurationPtr(time.Minute * 4),
		EnableBatchedOperations:             to.BoolPtr(true),
		Status:                              &status,
		AutoDeleteOnIdle:                    toDurationPtr(time.Minute * 7),
		SupportsOrdering:                    to.BoolPtr(true),
	})
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	topicExists, err = adminClient.TopicExists(context.Background(), topicName)
	require.True(t, topicExists)
	require.NoError(t, err)

	require.EqualValues(t, &TopicProperties{
		Name:                                topicName,
		EnablePartitioning:                  to.BoolPtr(true),
		MaxSizeInMegabytes:                  to.Int32Ptr(16 * 2048), // enabling partitioning increases our max size because of the 16 partitions
		RequiresDuplicateDetection:          to.BoolPtr(true),
		DefaultMessageTimeToLive:            toDurationPtr(time.Minute * 3),
		DuplicateDetectionHistoryTimeWindow: toDurationPtr(time.Minute * 4),
		EnableBatchedOperations:             to.BoolPtr(true),
		Status:                              &status,
		AutoDeleteOnIdle:                    toDurationPtr(time.Minute * 7),
		SupportsOrdering:                    to.BoolPtr(true),
	}, addResp.Value)

	getResp, err := adminClient.GetTopic(context.Background(), topicName)
	require.NoError(t, err)

	require.EqualValues(t, &TopicProperties{
		Name:                                topicName,
		EnablePartitioning:                  to.BoolPtr(true),
		MaxSizeInMegabytes:                  to.Int32Ptr(16 * 2048), // enabling partitioning increases our max size because of the 16 partitions
		RequiresDuplicateDetection:          to.BoolPtr(true),
		DefaultMessageTimeToLive:            toDurationPtr(time.Minute * 3),
		DuplicateDetectionHistoryTimeWindow: toDurationPtr(time.Minute * 4),
		EnableBatchedOperations:             to.BoolPtr(true),
		Status:                              &status,
		AutoDeleteOnIdle:                    toDurationPtr(time.Minute * 7),
		SupportsOrdering:                    to.BoolPtr(true),
	}, getResp.Value)

	// check some subscriptions properties, existence
	subExists, err := adminClient.SubscriptionExists(context.Background(), topicName, subscriptionName)
	require.NoError(t, err)
	require.False(t, subExists)

	addSubWithPropsResp, err := adminClient.AddSubscriptionWithProperties(context.Background(), topicName, &SubscriptionProperties{
		Name:                             subscriptionName,
		LockDuration:                     toDurationPtr(3 * time.Minute),
		RequiresSession:                  to.BoolPtr(false),
		DefaultMessageTimeToLive:         toDurationPtr(7 * time.Minute),
		DeadLetteringOnMessageExpiration: to.BoolPtr(true),
		EnableDeadLetteringOnFilterEvaluationExceptions: to.BoolPtr(false),
		MaxDeliveryCount: to.Int32Ptr(11),
		Status:           &status,
		// ForwardTo:                     &forwardToQueueName,
		// ForwardDeadLetteredMessagesTo: &forwardToQueueName,
		EnableBatchedOperations: to.BoolPtr(false),
		UserMetadata:            to.StringPtr("some user metadata"),
	})

	require.NoError(t, err)
	require.EqualValues(t, &SubscriptionProperties{
		Name:                             subscriptionName,
		LockDuration:                     toDurationPtr(3 * time.Minute),
		RequiresSession:                  to.BoolPtr(false),
		DefaultMessageTimeToLive:         toDurationPtr(7 * time.Minute),
		DeadLetteringOnMessageExpiration: to.BoolPtr(true),
		EnableDeadLetteringOnFilterEvaluationExceptions: to.BoolPtr(false),
		MaxDeliveryCount: to.Int32Ptr(11),
		Status:           &status,
		// ForwardTo:                     &forwardToQueueName,
		// ForwardDeadLetteredMessagesTo: &forwardToQueueName,
		EnableBatchedOperations: to.BoolPtr(false),
		UserMetadata:            to.StringPtr("some user metadata"),
	}, addSubWithPropsResp.Value)

	subExists, err = adminClient.SubscriptionExists(context.Background(), topicName, subscriptionName)
	require.NoError(t, err)
	require.True(t, subExists)

	defer deleteSubscription(t, adminClient, topicName, subscriptionName)
}

func TestAdminClient_UpdateTopic(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	addResp, err := adminClient.AddTopic(context.Background(), topicName)
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	addResp.Value.AutoDeleteOnIdle = toDurationPtr(11 * time.Minute)
	updateResp, err := adminClient.UpdateTopic(context.Background(), addResp.Value)
	require.NoError(t, err)

	require.EqualValues(t, 11*time.Minute, *updateResp.Value.AutoDeleteOnIdle)

	// try changing a value that's not allowed
	updateResp.Value.EnablePartitioning = to.BoolPtr(true)
	updateResp, err = adminClient.UpdateTopic(context.Background(), updateResp.Value)
	require.Contains(t, err.Error(), "Partitioning cannot be changed for Topic. ")
	require.Nil(t, updateResp)

	addResp.Value.Name = "non-existent-topic"
	updateResp, err = adminClient.UpdateTopic(context.Background(), addResp.Value)
	// a little awkward, we'll make these programatically inspectable as we add in better error handling.
	require.Contains(t, err.Error(), "error code: 404")
	require.Nil(t, updateResp)
}

func TestAdminClient_TopicAndSubscriptionRuntimeProperties(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	client, err := NewClientFromConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	subscriptionName := fmt.Sprintf("sub-%X", time.Now().UnixNano())

	_, err = adminClient.AddTopic(context.Background(), topicName)
	require.NoError(t, err)

	addSubResp, err := adminClient.AddSubscription(context.Background(), topicName, subscriptionName)
	require.NoError(t, err)
	require.NotNil(t, addSubResp)
	require.EqualValues(t, subscriptionName, addSubResp.Value.Name)

	defer deleteSubscription(t, adminClient, topicName, subscriptionName)

	sender, err := client.NewSender(topicName)
	require.NoError(t, err)

	// trigger some stats

	//  Scheduled messages are accounted for in the topic stats.
	_, err = sender.ScheduleMessages(context.Background(), []SendableMessage{
		&Message{Body: []byte("hello")},
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

	require.EqualValues(t, subscriptionName, getSubResp.Value.Name)
	require.EqualValues(t, 0, getSubResp.Value.ActiveMessageCount)
	require.NotEqual(t, time.Time{}, getSubResp.Value.CreatedAt)
	require.NotEqual(t, time.Time{}, getSubResp.Value.UpdatedAt)
	require.NotEqual(t, time.Time{}, getSubResp.Value.AccessedAt)
}

func TestAdminClient_ListTopics(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
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
			_, err = adminClient.AddTopicWithProperties(context.Background(), &TopicProperties{
				Name:                     topicName,
				DefaultMessageTimeToLive: toDurationPtr(time.Duration(i+1) * time.Minute),
			})
			require.NoError(t, err)
		}(i)

		defer deleteTopic(t, adminClient, topicName)
	}

	wg.Wait()

	// we skipped the first topic so it shouldn't come back in the results.
	pager := adminClient.ListTopics(nil)
	all := map[string]*TopicProperties{}
	var allNames []string

	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()

		for _, props := range page.Value {
			_, exists := all[props.Name]
			require.False(t, exists, "Each topic result should be unique")
			all[props.Name] = props
			allNames = append(allNames, props.Name)
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

	// grab the second to last item
	pager = adminClient.ListTopics(&ListTopicsOptions{
		Skip: len(allNames) - 2,
		Top:  1,
	})

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, 1, len(pager.PageResponse().Value))
	require.EqualValues(t, allNames[len(allNames)-2], pager.PageResponse().Value[0].Name)
	require.NoError(t, pager.Err())

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, allNames[len(allNames)-1], pager.PageResponse().Value[0].Name)
	require.False(t, pager.NextPage(context.Background()))
}

func TestAdminClient_ListTopicsRuntimeProperties(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	var expectedTopics []string
	now := time.Now().UnixNano()

	for i := 0; i < 3; i++ {
		topicName := strings.ToLower(fmt.Sprintf("topic-%d-%X", i, now))
		expectedTopics = append(expectedTopics, topicName)

		_, err = adminClient.AddTopic(context.Background(), topicName)
		require.NoError(t, err)

		defer deleteTopic(t, adminClient, topicName)
	}

	// we skipped the first topic so it shouldn't come back in the results.
	pager := adminClient.ListTopicsRuntimeProperties(nil)
	all := map[string]*TopicRuntimeProperties{}
	var allNames []string

	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()

		for _, props := range page.Value {
			_, exists := all[props.Name]
			require.False(t, exists, "Each topic result should be unique")
			all[props.Name] = props
			allNames = append(allNames, props.Name)
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

	// grab the second to last item
	pager = adminClient.ListTopicsRuntimeProperties(&ListTopicsRuntimePropertiesOptions{
		Skip: len(allNames) - 2,
		Top:  1,
	})

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, 1, len(pager.PageResponse().Value))
	require.EqualValues(t, allNames[len(allNames)-2], pager.PageResponse().Value[0].Name)
	require.NoError(t, pager.Err())

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, allNames[len(allNames)-1], pager.PageResponse().Value[0].Name)
	require.False(t, pager.NextPage(context.Background()))
}

func TestAdminClient_ListSubscriptions(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	now := time.Now().UnixNano()
	topicName := strings.ToLower(fmt.Sprintf("topic-%X", now))

	_, err = adminClient.AddTopic(context.Background(), topicName)
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	var expectedSubscriptions []string

	for i := 0; i < 3; i++ {
		subName := strings.ToLower(fmt.Sprintf("sub-%d-%X", i, now))
		expectedSubscriptions = append(expectedSubscriptions, subName)

		_, err = adminClient.AddSubscriptionWithProperties(context.Background(), topicName, &SubscriptionProperties{
			Name:                     subName,
			DefaultMessageTimeToLive: toDurationPtr(time.Duration(i+1) * time.Minute),
		})
		require.NoError(t, err)

		defer deleteSubscription(t, adminClient, topicName, subName)
	}

	// we skipped the first topic so it shouldn't come back in the results.
	pager := adminClient.ListSubscriptions(topicName, nil)
	all := map[string]*SubscriptionProperties{}
	var allNames []string

	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()

		for _, props := range page.Value {
			_, exists := all[props.Name]
			require.False(t, exists, "Each subscription result should be unique")
			all[props.Name] = props
			allNames = append(allNames, props.Name)
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

	// grab the second to last item
	pager = adminClient.ListSubscriptions(topicName, &ListSubscriptionsOptions{
		Skip: len(allNames) - 2,
		Top:  1,
	})

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, 1, len(pager.PageResponse().Value))
	require.EqualValues(t, allNames[len(allNames)-2], pager.PageResponse().Value[0].Name)
	require.NoError(t, pager.Err())

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, allNames[len(allNames)-1], pager.PageResponse().Value[0].Name)
	require.False(t, pager.NextPage(context.Background()))
}

func TestAdminClient_ListSubscriptionRuntimeProperties(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	now := time.Now().UnixNano()
	topicName := strings.ToLower(fmt.Sprintf("topic-%X", now))

	_, err = adminClient.AddTopic(context.Background(), topicName)
	require.NoError(t, err)

	var expectedSubs []string

	for i := 0; i < 3; i++ {
		subscriptionName := strings.ToLower(fmt.Sprintf("sub-%d-%X", i, now))
		expectedSubs = append(expectedSubs, subscriptionName)

		_, err = adminClient.AddSubscription(context.Background(), topicName, subscriptionName)
		require.NoError(t, err)

		defer deleteSubscription(t, adminClient, topicName, subscriptionName)
	}

	// we skipped the first subscription so it shouldn't come back in the results.
	pager := adminClient.ListSubscriptionsRuntimeProperties(topicName, nil)
	all := map[string]*SubscriptionRuntimeProperties{}
	var allNames []string

	for pager.NextPage(context.Background()) {
		page := pager.PageResponse()

		for _, props := range page.Value {
			_, exists := all[props.Name]
			require.False(t, exists, "Each subscription result should be unique")
			all[props.Name] = props
			allNames = append(allNames, props.Name)
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

	// grab the second to last item
	pager = adminClient.ListSubscriptionsRuntimeProperties(topicName, &ListSubscriptionsRuntimePropertiesOptions{
		Skip: len(allNames) - 2,
		Top:  1,
	})

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, 1, len(pager.PageResponse().Value))
	require.EqualValues(t, allNames[len(allNames)-2], pager.PageResponse().Value[0].Name)
	require.NoError(t, pager.Err())

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, allNames[len(allNames)-1], pager.PageResponse().Value[0].Name)
	require.False(t, pager.NextPage(context.Background()))
}

func TestAdminClient_UpdateSubscription(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	_, err = adminClient.AddTopic(context.Background(), topicName)
	require.NoError(t, err)

	defer deleteTopic(t, adminClient, topicName)

	subscriptionName := fmt.Sprintf("sub-%X", time.Now().UnixNano())
	addResp, err := adminClient.AddSubscription(context.Background(), topicName, subscriptionName)
	require.NoError(t, err)

	defer deleteSubscription(t, adminClient, topicName, subscriptionName)

	addResp.Value.LockDuration = toDurationPtr(4 * time.Minute)
	updateResp, err := adminClient.UpdateSubscription(context.Background(), topicName, addResp.Value)
	require.NoError(t, err)

	require.EqualValues(t, 4*time.Minute, *updateResp.Value.LockDuration)

	// try changing a value that's not allowed
	updateResp.Value.RequiresSession = to.BoolPtr(true)
	updateResp, err = adminClient.UpdateSubscription(context.Background(), topicName, updateResp.Value)
	require.Contains(t, err.Error(), "The value for the RequiresSession property of an existing Subscription cannot be changed")
	require.Nil(t, updateResp)

	addResp.Value.Name = "non-existent-subscription"
	updateResp, err = adminClient.UpdateSubscription(context.Background(), topicName, addResp.Value)
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
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	lowPrivAdminClient, err := NewAdminClientWithConnectionString(getConnectionStringWithoutManagePerms(t), nil)
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
		_, err = adminClient.AddQueue(context.Background(), queueName)
		require.NoError(t, err)

		_, err = adminClient.AddTopic(context.Background(), topicName)
		require.NoError(t, err)

		_, err = adminClient.AddSubscription(context.Background(), topicName, subName)
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
	require.True(t, atom.NotFound(err))

	_, err = testData.Client.GetQueue(ctx, testData.QueueName)
	require.Contains(t, err.Error(), "error code: 401, Details: Manage,EntityRead claims")

	pager := testData.Client.ListQueues(nil)
	require.False(t, pager.NextPage(context.Background()))
	require.Contains(t, pager.Err().Error(), "error code: 401, Details: Manage,EntityRead claims required for this operation")

	_, err = testData.Client.AddQueue(ctx, "canneverbecreated")
	require.Contains(t, err.Error(), "error code: 401, Details: Authorization failed for specified action: Manage,EntityWrite")

	_, err = testData.Client.UpdateQueue(ctx, &QueueProperties{
		Name: "canneverbecreated",
	})
	require.Contains(t, err.Error(), "error code: 401, Details: Authorization failed for specified action: Manage,EntityWrite")

	_, err = testData.Client.DeleteQueue(ctx, testData.QueueName)
	require.Contains(t, err.Error(), "error code: 401, Details: Authorization failed for specified action: Manage,EntityDelete.")
}

func TestAdminClient_LackPermissions_Topic(t *testing.T) {
	testData := setupLowPrivTest(t)
	defer testData.Cleanup()

	ctx := context.Background()

	_, err := testData.Client.GetTopic(ctx, "not-found-topic")
	require.True(t, atom.NotFound(err))

	_, err = testData.Client.GetTopic(ctx, testData.TopicName)
	require.Contains(t, err.Error(), "error code: 401, Details: Manage,EntityRead claims")

	pager := testData.Client.ListTopics(nil)
	require.False(t, pager.NextPage(context.Background()))
	require.Contains(t, pager.Err().Error(), "error code: 401, Details: Manage,EntityRead claims required for this operation")

	_, err = testData.Client.AddTopic(ctx, "canneverbecreated")
	require.Contains(t, err.Error(), "error code: 401, Details: Authorization failed for specified action")

	_, err = testData.Client.UpdateTopic(ctx, &TopicProperties{
		Name: "canneverbecreated",
	})
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

	_, err := testData.Client.GetSubscription(ctx, testData.TopicName, "not-found-sub")
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'GetSubscription'")

	_, err = testData.Client.GetSubscription(ctx, testData.TopicName, testData.SubName)
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'GetSubscription'")

	pager := testData.Client.ListSubscriptions(testData.TopicName, nil)
	require.False(t, pager.NextPage(context.Background()))
	require.Contains(t, pager.Err().Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'EnumerateSubscriptions' operation")

	_, err = testData.Client.AddSubscription(ctx, testData.TopicName, "canneverbecreated")
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'CreateOrUpdateSubscription'")

	_, err = testData.Client.UpdateSubscription(ctx, testData.TopicName, &SubscriptionProperties{
		Name: "canneverbecreated",
	})
	require.Contains(t, err.Error(), "401 SubCode=40100: Unauthorized : Unauthorized access for 'CreateOrUpdateSubscription'")

	_, err = testData.Client.DeleteSubscription(ctx, testData.TopicName, testData.SubName)
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
	_, err := ac.DeleteSubscription(context.Background(), topicName, subscriptionName)
	require.NoError(t, err)
}
