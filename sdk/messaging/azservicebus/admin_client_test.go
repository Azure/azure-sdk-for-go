// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/conn"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/stretchr/testify/require"
)

func TestAdminClient_Queue_Forwarding(t *testing.T) {
	cs := test.GetConnectionString(t)
	adminClient, err := admin.NewClientFromConnectionString(cs, nil)
	require.NoError(t, err)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	forwardToQueueName := fmt.Sprintf("queue-fwd-%X", time.Now().UnixNano())

	_, err = adminClient.CreateQueue(context.Background(), forwardToQueueName, nil)
	require.NoError(t, err)

	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), forwardToQueueName, nil)
		require.NoError(t, err)
	}()

	parsed, err := conn.ParseConnectionString(cs)
	require.NoError(t, err)

	formatted := fmt.Sprintf("%s%s", fmt.Sprintf("https://%s/", parsed.FullyQualifiedNamespace), forwardToQueueName)

	createResp, err := adminClient.CreateQueue(context.Background(), queueName, &admin.CreateQueueOptions{
		Properties: &admin.QueueProperties{
			ForwardTo:                     &formatted,
			ForwardDeadLetteredMessagesTo: &formatted,
		},
	})

	require.NoError(t, err)
	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), queueName, nil)
		require.NoError(t, err)
	}()

	require.EqualValues(t, formatted, *createResp.ForwardTo)
	require.EqualValues(t, formatted, *createResp.ForwardDeadLetteredMessagesTo)

	getResp, err := adminClient.GetQueue(context.Background(), queueName, nil)

	require.NoError(t, err)
	require.EqualValues(t, createResp.QueueProperties, getResp.QueueProperties)

	client, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{
		Body: []byte("this message will be auto-forwarded"),
	}, nil)
	require.NoError(t, err)

	receiver, err := client.NewReceiverForQueue(forwardToQueueName, nil)
	require.NoError(t, err)

	forwardedMessages, err := receiver.ReceiveMessages(context.Background(), 1, nil)
	require.NoError(t, err)

	require.EqualValues(t, "this message will be auto-forwarded", string(forwardedMessages[0].Body))
}

func TestAdminClient_GetQueueRuntimeProperties(t *testing.T) {
	adminClient, err := admin.NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	client, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)
	defer client.Close(context.Background())

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	_, err = adminClient.CreateQueue(context.Background(), queueName, nil)
	require.NoError(t, err)

	defer deleteQueue(t, adminClient, queueName)

	sender, err := client.NewSender(queueName, nil)
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		err = sender.SendMessage(context.Background(), &Message{
			Body: []byte("hello"),
		}, nil)
		require.NoError(t, err)
	}

	sequenceNumbers, err := sender.ScheduleMessages(context.Background(), []*Message{
		{Body: []byte("hello")},
	}, time.Now().Add(2*time.Hour), nil)
	require.NoError(t, err)
	require.NotEmpty(t, sequenceNumbers)

	receiver, err := client.NewReceiverForQueue(queueName, nil)
	require.NoError(t, err)

	messages, err := receiver.ReceiveMessages(context.Background(), 2, nil)
	require.NoError(t, err)

	require.NoError(t, receiver.DeadLetterMessage(context.Background(), messages[0], nil))

	props, err := adminClient.GetQueueRuntimeProperties(context.Background(), queueName, nil)
	require.NoError(t, err)

	require.EqualValues(t, 4, props.TotalMessageCount)

	require.EqualValues(t, 2, props.ActiveMessageCount)
	require.EqualValues(t, 1, props.DeadLetterMessageCount)
	require.EqualValues(t, 1, props.ScheduledMessageCount)
	require.EqualValues(t, 0, props.TransferDeadLetterMessageCount)
	require.EqualValues(t, 0, props.TransferMessageCount)

	require.Greater(t, props.SizeInBytes, int64(0))

	require.False(t, props.CreatedAt.IsZero())
	require.False(t, props.UpdatedAt.IsZero())
	require.False(t, props.AccessedAt.IsZero())
}

func TestAdminClient_TopicAndSubscriptionRuntimeProperties(t *testing.T) {
	adminClient, err := admin.NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	client, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	topicName := fmt.Sprintf("topic-%X", time.Now().UnixNano())
	subscriptionName := fmt.Sprintf("sub-%X", time.Now().UnixNano())

	_, err = adminClient.CreateTopic(context.Background(), topicName, nil)
	require.NoError(t, err)

	addSubResp, err := adminClient.CreateSubscription(context.Background(), topicName, subscriptionName, nil)
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
	}, time.Now().Add(2*time.Hour), nil)
	require.NoError(t, err)

	// validate the topic runtime properties
	getRuntimeResp, err := adminClient.GetTopicRuntimeProperties(context.Background(), topicName, nil)
	require.NoError(t, err)

	require.EqualValues(t, 1, getRuntimeResp.SubscriptionCount)
	require.False(t, getRuntimeResp.CreatedAt.IsZero())
	require.False(t, getRuntimeResp.UpdatedAt.IsZero())
	require.False(t, getRuntimeResp.AccessedAt.IsZero())

	require.Greater(t, getRuntimeResp.SizeInBytes, int64(0))
	require.EqualValues(t, int32(1), getRuntimeResp.ScheduledMessageCount)

	// validate subscription runtime properties
	getSubResp, err := adminClient.GetSubscriptionRuntimeProperties(context.Background(), topicName, subscriptionName, nil)
	require.NoError(t, err)

	require.EqualValues(t, 0, getSubResp.ActiveMessageCount)
	require.False(t, getSubResp.CreatedAt.IsZero())
	require.False(t, getSubResp.UpdatedAt.IsZero())
	require.False(t, getSubResp.AccessedAt.IsZero())
}

func TestAdminClient_StringToTime(t *testing.T) {
	tm, err := atom.StringToTime("2021-11-22T23:07:33.08708Z")
	require.False(t, tm.IsZero())
	require.NoError(t, err)

	// You'll see this uninitialized timestamp when you look at the response from a PUT request.
	// It's the reason we can't use the much simpler method of just declaring a field as  time.Time in the various
	// <Entity>Description structs.
	// We don't even return AccessedTime in in that context so any value will be fine.
	tm, err = atom.StringToTime("0001-01-01T00:00:00")
	require.True(t, tm.IsZero())
	require.Nil(t, err)

	// and if it's just some ill-f0rmed timestamp we'll just fallback to giving them the zero time.
	tm, err = atom.StringToTime("Not a timestamp")
	require.Error(t, err)
	require.True(t, tm.IsZero())

	tm, err = atom.StringToTime("")
	require.Error(t, err)
	require.True(t, tm.IsZero())
}
