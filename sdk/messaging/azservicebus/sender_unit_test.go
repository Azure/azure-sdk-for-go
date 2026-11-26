// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock/emulation"
	"github.com/Azure/go-amqp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSender_SendMessage_RecoversFromConnectionScopedNotAllowed(t *testing.T) {
	var sendAttempts int
	var md *emulation.MockData
	md, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreSenderMock: func(ms *emulation.MockSender, ctx context.Context) error {
			if ms.Target == "queue" {
				ms.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, msg *amqp.Message, o *amqp.SendOptions) error {
					sendAttempts++
					if sendAttempts == 1 {
						return fmt.Errorf("wrapped: %w", &amqp.ConnError{RemoteErr: &amqp.Error{Condition: amqp.ErrCondNotAllowed}})
					}

					return md.AllQueues()[ms.Target].Send(ctx, msg, ms.LinkEvent(), ms.Status)
				}).AnyTimes()
			}

			return nil
		},
	}, &ClientOptions{
		RetryOptions: exported.RetryOptions{
			MaxRetries:    1,
			RetryDelay:    time.Millisecond,
			MaxRetryDelay: time.Millisecond,
		},
	})
	defer cleanup()

	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)
	senderLinks := sender.links

	err = sender.SendMessage(context.Background(), &Message{Body: []byte("hello")}, nil)
	require.NoError(t, err)
	require.Equal(t, 2, sendAttempts)
	require.Equal(t, 2, countEmulationEvents(md.Events.All(), emulation.EventTypeConnOpen))
	require.Equal(t, 3, len(md.Events.GetOpenLinks()))
	require.Same(t, senderLinks, sender.links)
}

func TestSender_SendMessage_RecoveryExhaustionLeavesSenderUsable(t *testing.T) {
	failSends := true
	var sendAttempts int
	var md *emulation.MockData
	md, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreSenderMock: func(ms *emulation.MockSender, ctx context.Context) error {
			if ms.Target == "queue" {
				ms.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, msg *amqp.Message, o *amqp.SendOptions) error {
					sendAttempts++
					if failSends {
						return fmt.Errorf("wrapped: %w", &amqp.ConnError{RemoteErr: &amqp.Error{Condition: amqp.ErrCondNotAllowed}})
					}

					return md.AllQueues()[ms.Target].Send(ctx, msg, ms.LinkEvent(), ms.Status)
				}).AnyTimes()
			}

			return nil
		},
	}, &ClientOptions{
		RetryOptions: exported.RetryOptions{
			MaxRetries:    1,
			RetryDelay:    time.Millisecond,
			MaxRetryDelay: time.Millisecond,
		},
	})
	defer cleanup()

	sender, err := client.NewSender("queue", nil)
	require.NoError(t, err)

	err = sender.SendMessage(context.Background(), &Message{Body: []byte("first")}, nil)
	require.Error(t, err)
	require.Equal(t, 2, sendAttempts)

	failSends = false
	err = sender.SendMessage(context.Background(), &Message{Body: []byte("second")}, nil)
	require.NoError(t, err)
	require.Equal(t, 3, sendAttempts)
}

func countEmulationEvents(events []emulation.Event, eventType emulation.EventType) int {
	count := 0
	for _, event := range events {
		if event.Type == eventType {
			count++
		}
	}
	return count
}

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

func TestSenderNewMessageBatch_VendorPropertyOverridesMaxMessageSize(t *testing.T) {
	_, client, cleanup := newClientWithMockedConn(t, &emulation.MockDataOptions{
		PreReceiverMock: func(mr *emulation.MockReceiver, ctx context.Context) error {
			return nil
		},
		PreSenderMock: func(ms *emulation.MockSender, ctx context.Context) error {
			if ms.Target != "$cbs" {
				// Set MaxMessageSize to 100 MB and the vendor batch-size property to 1 MB
				// so this test can verify the vendor property is used as the batch limit.
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
	// A 1 MiB body plus AMQP envelope overhead exceeds the 1 MiB batch limit.
	largeBody := make([]byte, 1048576)
	err = batch.AddMessage(&Message{Body: largeBody}, nil)
	require.ErrorIs(t, err, ErrMessageTooLarge, "A 1 MiB message should exceed the vendor batch limit minus overhead")
}

func TestSenderNewMessageBatch_FallsBackWhenVendorPropertyAbsent(t *testing.T) {
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

func TestSenderNewMessageBatch_UserMaxBytesOverridesVendorProperty(t *testing.T) {
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
