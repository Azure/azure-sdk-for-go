// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock/emulation"
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

func TestNewMessageBatch_VendorPropertyOverridesMaxMessageSize(t *testing.T) {
	_, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			return nil
		},
		PreSenderMock: func(ms *emulation.MockSender, ctx context.Context) error {
			if ms.Target != "$cbs" {
				// Override: set MaxMessageSize to 100 MB (Premium large-message)
				// and vendor property to 1 MB (correct batch limit).
				// These are registered before the default expectations in mock_data_sender.go,
				// but with AnyTimes() gomock uses the last registered matching call,
				// so we need the defaults to NOT be registered or to use Times(0).
				// Actually, in gomock, when multiple .AnyTimes() expectations match,
				// the FIRST registered one is used. So pre-mock wins.
				ms.EXPECT().MaxMessageSize().Return(uint64(100 * 1024 * 1024)).AnyTimes()
				ms.EXPECT().Properties().Return(map[string]any{
					"com.microsoft:max-message-batch-size": uint64(1048576),
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

	batch, err := sender.NewMessageBatch(context.Background(), nil)
	require.NoError(t, err)
	require.NotNil(t, batch)

	// The batch should use the vendor property (1 MB), not MaxMessageSize (100 MB).
	// We verify indirectly: a message just over 1 MB should be rejected.
	largeBody := make([]byte, 1048576)
	err = batch.AddMessage(&Message{Body: largeBody}, nil)
	require.ErrorIs(t, err, ErrMessageTooLarge, "A 1 MB message should exceed the vendor batch limit minus overhead")
}

func TestNewMessageBatch_FallsBackWhenVendorPropertyAbsent(t *testing.T) {
	_, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			return nil
		},
		PreSenderMock: func(ms *emulation.MockSender, ctx context.Context) error {
			if ms.Target != "$cbs" {
				ms.EXPECT().MaxMessageSize().Return(uint64(262144)).AnyTimes() // 256 KB
				ms.EXPECT().Properties().Return(map[string]any(nil)).AnyTimes()
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
	require.NoError(t, err)
	require.NotNil(t, batch)

	// Should fall back to MaxMessageSize (256 KB). A 256 KB body should be rejected.
	body := make([]byte, 262144)
	err = batch.AddMessage(&Message{Body: body}, nil)
	require.ErrorIs(t, err, ErrMessageTooLarge, "A 256 KB message should exceed the link limit minus overhead")
}

func TestNewMessageBatch_UserMaxBytesOverridesVendorProperty(t *testing.T) {
	_, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			return nil
		},
		PreSenderMock: func(ms *emulation.MockSender, ctx context.Context) error {
			if ms.Target != "$cbs" {
				ms.EXPECT().MaxMessageSize().Return(uint64(100 * 1024 * 1024)).AnyTimes()
				ms.EXPECT().Properties().Return(map[string]any{
					"com.microsoft:max-message-batch-size": uint64(1048576),
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

	batch, err := sender.NewMessageBatch(context.Background(), &MessageBatchOptions{
		MaxBytes: 512,
	})
	require.NoError(t, err)
	require.NotNil(t, batch)

	// User override of 512 bytes — a small message should still be rejected
	body := make([]byte, 512)
	err = batch.AddMessage(&Message{Body: body}, nil)
	require.ErrorIs(t, err, ErrMessageTooLarge, "A 512-byte message should exceed the user-specified 512-byte limit minus overhead")
}
