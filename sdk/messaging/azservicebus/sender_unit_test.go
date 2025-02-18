// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"testing"
	"time"

	testtracing "github.com/Azure/azure-sdk-for-go/sdk/internal/test/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock/emulation"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/Azure/go-amqp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSender_UserFacingError(t *testing.T) {
	_, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source != "$cbs" {
				mr.EXPECT().Receive(gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
					return nil, &amqp.ConnError{}
				}).AnyTimes()
			}

			return nil
		},
		PreSenderMock: func(ms *emulation.MockSender, ctx context.Context) error {
			if ms.Target != "$cbs" {
				ms.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, m *amqp.Message, o *amqp.SendOptions) error {
					return &amqp.ConnError{}
				}).AnyTimes()
			}

			return nil
		},
	}, &ClientOptions{
		RetryOptions: noRetriesNeeded,
	})

	defer cleanup()

	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)

	var asSBError *Error

	err = sender.SendMessage(context.Background(), &Message{}, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	msgID := "testID"
	err = sender.SendMessage(context.Background(), &Message{MessageID: &msgID}, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	err = sender.CancelScheduledMessages(context.Background(), []int64{1}, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	seqNumbers, err := sender.ScheduleMessages(context.Background(), []*Message{}, time.Now(), nil)
	require.Empty(t, seqNumbers)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	// link is already initialized, so this will work.
	batch, err := sender.NewMessageBatch(context.Background(), nil)
	require.NoError(t, err)

	err = batch.AddMessage(&Message{
		Body: []byte("hello"),
	}, nil)
	require.NoError(t, err)

	err = sender.SendMessageBatch(context.Background(), batch, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)
}

func TestSender_TracingUserFacingError(t *testing.T) {
	_, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source != "$cbs" {
				mr.EXPECT().Receive(gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
					return nil, &amqp.ConnError{}
				}).AnyTimes()
			}

			return nil
		},
		PreSenderMock: func(ms *emulation.MockSender, ctx context.Context) error {
			if ms.Target != "$cbs" {
				ms.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, m *amqp.Message, o *amqp.SendOptions) error {
					return &amqp.ConnError{}
				}).AnyTimes()
			}

			return nil
		},
	}, &ClientOptions{
		RetryOptions: noRetriesNeeded,
	})

	defer cleanup()

	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)

	var asSBError *Error

	fakeClientCreds := clientCreds{fullyQualifiedNamespace: "fake.something"}

	sender.tracer = newTracer(testtracing.NewSpanValidator(t, testtracing.SpanMatcher{
		Name:   "send queue",
		Kind:   tracing.SpanKindProducer,
		Status: tracing.SpanStatusError,
		Attributes: []tracing.Attribute{
			{Key: tracing.MessagingSystem, Value: "servicebus"},
			{Key: tracing.ServerAddress, Value: "fake.something"},
			{Key: tracing.DestinationName, Value: "queue"},
			{Key: tracing.OperationName, Value: "send"},
			{Key: tracing.OperationType, Value: "send"},
		},
	}), fakeClientCreds, "queue", "")
	err = sender.SendMessage(context.Background(), &Message{}, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	msgID := "testID"
	sender.tracer = newTracer(testtracing.NewSpanValidator(t, testtracing.SpanMatcher{
		Name:   "send queue",
		Kind:   tracing.SpanKindProducer,
		Status: tracing.SpanStatusError,
		Attributes: []tracing.Attribute{
			{Key: tracing.MessagingSystem, Value: "servicebus"},
			{Key: tracing.ServerAddress, Value: "fake.something"},
			{Key: tracing.DestinationName, Value: "queue"},
			{Key: tracing.OperationName, Value: "send"},
			{Key: tracing.OperationType, Value: "send"},
			{Key: tracing.MessageID, Value: msgID},
		},
	}), fakeClientCreds, "queue", "")
	err = sender.SendMessage(context.Background(), &Message{MessageID: &msgID}, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	sender.tracer = newTracer(testtracing.NewSpanValidator(t, testtracing.SpanMatcher{
		Name:   "cancel_scheduled queue",
		Kind:   tracing.SpanKindClient,
		Status: tracing.SpanStatusError,
		Attributes: []tracing.Attribute{
			{Key: tracing.MessagingSystem, Value: "servicebus"},
			{Key: tracing.ServerAddress, Value: "fake.something"},
			{Key: tracing.DestinationName, Value: "queue"},
			{Key: tracing.OperationName, Value: "cancel_scheduled"},
			{Key: tracing.OperationType, Value: "send"},
			{Key: tracing.BatchMessageCount, Value: int64(1)},
		},
	}), fakeClientCreds, "queue", "")
	err = sender.CancelScheduledMessages(context.Background(), []int64{1}, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	createTracer := newTracer(testtracing.NewSpanValidator(t, testtracing.SpanMatcher{
		Name:   "create queue",
		Kind:   tracing.SpanKindProducer,
		Status: tracing.SpanStatusUnset,
		Attributes: []tracing.Attribute{
			{Key: tracing.MessagingSystem, Value: "servicebus"},
			{Key: tracing.ServerAddress, Value: "fake.something"},
			{Key: tracing.DestinationName, Value: "queue"},
			{Key: tracing.OperationName, Value: "create"},
			{Key: tracing.OperationType, Value: "create"},
			{Key: tracing.MessageID, Value: msgID},
		},
	}), fakeClientCreds, "queue", "")
	msg := &Message{MessageID: &msgID}
	createMessageSpan(context.Background(), createTracer, tracing.Span{}, msg.toAMQPMessage())

	sender.tracer = newTracer(testtracing.NewSpanValidator(t, testtracing.SpanMatcher{
		Name:   "schedule queue",
		Kind:   tracing.SpanKindClient,
		Status: tracing.SpanStatusError,
		Attributes: []tracing.Attribute{
			{Key: tracing.MessagingSystem, Value: "servicebus"},
			{Key: tracing.ServerAddress, Value: "fake.something"},
			{Key: tracing.DestinationName, Value: "queue"},
			{Key: tracing.OperationName, Value: "schedule"},
			{Key: tracing.OperationType, Value: "send"},
			{Key: tracing.BatchMessageCount, Value: int64(1)},
		},
		Links: []tracing.Link{{Attributes: []tracing.Attribute{{Key: tracing.MessageID, Value: msgID}}}},
	}), fakeClientCreds, "queue", "")
	seqNumbers, err := sender.ScheduleMessages(context.Background(), []*Message{msg}, time.Now(), nil)
	require.Empty(t, seqNumbers)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	// link is already initialized, so this will work.
	batch, err := sender.NewMessageBatch(context.Background(), nil)
	require.NoError(t, err)

	err = batch.AddMessage(&Message{
		MessageID: &msgID,
		Body:      []byte("hello"),
	}, nil)
	require.NoError(t, err)

	sender.tracer = newTracer(testtracing.NewSpanValidator(t, testtracing.SpanMatcher{
		Name:   "send queue",
		Kind:   tracing.SpanKindClient,
		Status: tracing.SpanStatusError,
		Attributes: []tracing.Attribute{
			{Key: tracing.MessagingSystem, Value: "servicebus"},
			{Key: tracing.ServerAddress, Value: "fake.something"},
			{Key: tracing.DestinationName, Value: "queue"},
			{Key: tracing.OperationName, Value: "send"},
			{Key: tracing.OperationType, Value: "send"},
			{Key: tracing.BatchMessageCount, Value: int64(1)},
		},
		Links: []tracing.Link{{Attributes: []tracing.Attribute{{Key: tracing.MessageID, Value: msgID}}}},
	}), fakeClientCreds, "queue", "")
	err = sender.SendMessageBatch(context.Background(), batch, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)
}

func TestSenderNewMessageBatch_ConnectionClosed(t *testing.T) {
	_, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source != "$cbs" {
				mr.EXPECT().Receive(gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
					return nil, &amqp.ConnError{}
				}).AnyTimes()
			}

			return nil
		},
		PreSenderMock: func(ms *emulation.MockSender, ctx context.Context) error {
			if ms.Target != "$cbs" {
				return &amqp.ConnError{}
			}

			return nil
		},
	}, &ClientOptions{
		RetryOptions: noRetriesNeeded,
	})

	defer cleanup()

	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)

	batch, err := sender.NewMessageBatch(context.Background(), nil)
	var asSBError *Error
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)
	require.Nil(t, batch)
}
