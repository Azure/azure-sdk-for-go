// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
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
	require.EqualValues(t, queueName, props.Name)

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

	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), queueName)
		require.NoError(t, err)
	}()

	props, err := adminClient.GetQueue(context.Background(), queueName)
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
	}, props)
}

func TestAdminClient_Queue(t *testing.T) {
	adminClient, err := NewAdminClientWithConnectionString(getConnectionString(t), nil)
	require.NoError(t, err)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())

	es := EntityStatusReceiveDisabled
	propsFromCreate, err := adminClient.AddQueueWithProperties(context.Background(), &QueueProperties{
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
	}, propsFromCreate)

	propsFromGet, err := adminClient.GetQueue(context.Background(), queueName)
	require.NoError(t, err)

	require.EqualValues(t, propsFromGet, propsFromCreate)
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

	propsFromCreate, err := adminClient.AddQueueWithProperties(context.Background(), &QueueProperties{
		Name:                          queueName,
		ForwardTo:                     &formatted,
		ForwardDeadLetteredMessagesTo: &formatted,
	})

	require.NoError(t, err)
	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), queueName)
		require.NoError(t, err)
	}()

	require.EqualValues(t, formatted, *propsFromCreate.ForwardTo)
	require.EqualValues(t, formatted, *propsFromCreate.ForwardDeadLetteredMessagesTo)

	propsFromGet, err := adminClient.GetQueue(context.Background(), queueName)

	require.NoError(t, err)
	require.EqualValues(t, propsFromCreate, propsFromGet)

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

	createdProps.MaxDeliveryCount = to.Int32Ptr(101)
	updatedProps, err := adminClient.UpdateQueue(context.Background(), createdProps)
	require.NoError(t, err)

	require.EqualValues(t, 101, *updatedProps.MaxDeliveryCount)

	// try changing a value that's not allowed
	updatedProps.RequiresSession = to.BoolPtr(true)
	updatedProps, err = adminClient.UpdateQueue(context.Background(), updatedProps)
	require.Contains(t, err.Error(), "The value for the RequiresSession property of an existing Queue cannot be changed")
	require.Nil(t, updatedProps)

	createdProps.Name = "non-existent-queue"
	updatedProps, err = adminClient.UpdateQueue(context.Background(), createdProps)
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

	require.EqualValues(t, queueName, props.Name)

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

func TestAdminClient_getQueuePage(t *testing.T) {
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

		for _, props := range page {
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
	require.EqualValues(t, 1, len(pager.PageResponse()))
	require.EqualValues(t, allNames[len(allNames)-2], pager.PageResponse()[0].Name)
	require.NoError(t, pager.Err())

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, allNames[len(allNames)-1], pager.PageResponse()[0].Name)
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

		for _, props := range page {
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
	require.EqualValues(t, 1, len(pager.PageResponse()))
	require.EqualValues(t, allNames[len(allNames)-2], pager.PageResponse()[0].Name)
	require.NoError(t, pager.Err())

	require.True(t, pager.NextPage(context.Background()))
	require.EqualValues(t, allNames[len(allNames)-1], pager.PageResponse()[0].Name)
	require.False(t, pager.NextPage(context.Background()))
}

func TestAdminClient_deserializeATOMQueueEnvelope(t *testing.T) {
	reader, err := os.Open("testdata/queue_create_response.xml")
	require.NoError(t, err)
	envelope, err := deserializeQueueEnvelope(reader)
	require.NoError(t, err)

	queueProps, err := newQueueProperties("myqueuename", &envelope.Content.QueueDescription)
	require.NoError(t, err)
	require.NotNil(t, queueProps)

	es := EntityStatusReceiveDisabled

	require.EqualValues(t, &QueueProperties{
		AutoDeleteOnIdle:                    toDurationPtr(10 * time.Minute),
		DeadLetteringOnMessageExpiration:    to.BoolPtr(true),
		DefaultMessageTimeToLive:            toDurationPtr(time.Hour * 6),
		DuplicateDetectionHistoryTimeWindow: toDurationPtr(4 * time.Hour),
		EnableBatchedOperations:             to.BoolPtr(false),
		EnablePartitioning:                  to.BoolPtr(true),
		LockDuration:                        toDurationPtr(45 * time.Second),
		MaxDeliveryCount:                    to.Int32Ptr(100),
		MaxSizeInMegabytes:                  to.Int32Ptr(16 * 4096),
		Name:                                "myqueuename",
		RequiresDuplicateDetection:          to.BoolPtr(true),
		RequiresSession:                     to.BoolPtr(true),
		Status:                              &es,
	}, queueProps)
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

func toDurationPtr(d time.Duration) *time.Duration {
	return &d
}

func deleteQueue(t *testing.T, ac *AdminClient, queueName string) {
	_, err := ac.DeleteQueue(context.Background(), queueName)
	require.NoError(t, err)
}
