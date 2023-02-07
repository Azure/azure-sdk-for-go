// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock/emulation"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSender_UserFacingError(t *testing.T) {
	_, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source != "$cbs" {
				mr.EXPECT().Receive(gomock.Any()).DoAndReturn(func(ctx context.Context) (*amqp.Message, error) {
					return nil, &amqp.ConnectionError{}
				}).AnyTimes()
			}

			return nil
		},
		PreSenderMock: func(ms *emulation.MockSender, ctx context.Context) error {
			if ms.Target != "$cbs" {
				ms.EXPECT().Send(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, m *amqp.Message) error {
					return &amqp.ConnectionError{}
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

func TestSenderNewMessageBatch_ConnectionClosed(t *testing.T) {
	_, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			if mr.Source != "$cbs" {
				mr.EXPECT().Receive(gomock.Any()).DoAndReturn(func(ctx context.Context) (*amqp.Message, error) {
					return nil, &amqp.ConnectionError{}
				}).AnyTimes()
			}

			return nil
		},
		PreSenderMock: func(ms *emulation.MockSender, ctx context.Context) error {
			if ms.Target != "$cbs" {
				return &amqp.ConnectionError{}
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
