// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestSender_UserFacingError(t *testing.T) {
	fakeAMQPLinks := &internal.FakeAMQPLinks{}

	sender, err := newSender(newSenderArgs{
		ns: &internal.FakeNS{
			AMQPLinks: fakeAMQPLinks,
		},
		queueOrTopic:   "queue",
		cleanupOnClose: func() {},
		retryOptions: RetryOptions{
			MaxRetries:    0,
			RetryDelay:    0,
			MaxRetryDelay: 0,
		},
	})
	require.NoError(t, err)

	fakeAMQPLinks.Err = amqp.ErrConnClosed

	var asSBError *Error

	batch, err := sender.NewMessageBatch(context.Background(), nil)
	require.Nil(t, batch)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	err = sender.CancelScheduledMessages(context.Background(), []int64{1}, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	seqNumbers, err := sender.ScheduleMessages(context.Background(), []*Message{}, time.Now(), nil)
	require.Empty(t, seqNumbers)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	err = sender.SendMessage(context.Background(), &Message{}, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)

	err = sender.SendMessageBatch(context.Background(), nil, nil)
	require.ErrorAs(t, err, &asSBError)
	require.Equal(t, CodeConnectionLost, asSBError.Code)
}
