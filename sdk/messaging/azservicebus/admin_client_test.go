// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/conn"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/admin"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/stretchr/testify/require"
)

func TestAdminClient_Queue_Forwarding(t *testing.T) {
	cs := test.GetConnectionString(t)
	adminClient, err := admin.NewClientFromConnectionString(cs, nil)
	require.NoError(t, err)

	queueName := fmt.Sprintf("queue-%X", time.Now().UnixNano())
	forwardToQueueName := fmt.Sprintf("queue-fwd-%X", time.Now().UnixNano())

	_, err = adminClient.CreateQueue(context.Background(), forwardToQueueName, nil, nil)
	require.NoError(t, err)

	defer func() {
		_, err := adminClient.DeleteQueue(context.Background(), forwardToQueueName, nil)
		require.NoError(t, err)
	}()

	parsed, err := conn.ParsedConnectionFromStr(cs)
	require.NoError(t, err)

	formatted := fmt.Sprintf("%s%s", fmt.Sprintf("https://%s.%s/", parsed.Namespace, parsed.Suffix), forwardToQueueName)

	createResp, err := adminClient.CreateQueue(context.Background(), queueName, &admin.QueueProperties{
		ForwardTo:                     &formatted,
		ForwardDeadLetteredMessagesTo: &formatted,
	}, nil)

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

func TestAdminClient_GetQueueRuntimeProperties(t *testing.T) {
	adminClient, err := admin.NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	client, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
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

	props, err := adminClient.GetQueueRuntimeProperties(context.Background(), queueName, nil)
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

func TestAdminClient_TopicAndSubscriptionRuntimeProperties(t *testing.T) {
	adminClient, err := admin.NewClientFromConnectionString(test.GetConnectionString(t), nil)
	require.NoError(t, err)

	client, err := NewClientFromConnectionString(test.GetConnectionString(t), nil)
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
	getRuntimeResp, err := adminClient.GetTopicRuntimeProperties(context.Background(), topicName, nil)
	require.NoError(t, err)

	require.EqualValues(t, 1, getRuntimeResp.SubscriptionCount)
	require.NotEqual(t, time.Time{}, getRuntimeResp.CreatedAt)
	require.NotEqual(t, time.Time{}, getRuntimeResp.UpdatedAt)
	require.NotEqual(t, time.Time{}, getRuntimeResp.AccessedAt)

	require.Greater(t, getRuntimeResp.SizeInBytes, int64(0))
	require.EqualValues(t, int32(1), getRuntimeResp.ScheduledMessageCount)

	// validate subscription runtime properties
	getSubResp, err := adminClient.GetSubscriptionRuntimeProperties(context.Background(), topicName, subscriptionName, nil)
	require.NoError(t, err)

	require.EqualValues(t, 0, getSubResp.ActiveMessageCount)
	require.NotEqual(t, time.Time{}, getSubResp.CreatedAt)
	require.NotEqual(t, time.Time{}, getSubResp.UpdatedAt)
	require.NotEqual(t, time.Time{}, getSubResp.AccessedAt)
}
