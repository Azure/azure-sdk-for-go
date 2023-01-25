// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation_test

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock/emulation"
	"github.com/stretchr/testify/require"
)

func TestMockQueue(t *testing.T) {
	events := &emulation.Events{}
	mq := emulation.NewQueue("entity", events)
	defer mq.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// no messages exist yet
	msg, err := mq.Receive(ctx, emulation.LinkEvent{})
	require.Nil(t, msg)
	require.ErrorIs(t, err, context.DeadlineExceeded)

	err = mq.Send(context.Background(), &amqp.Message{
		Value: []byte("first message"),
	}, emulation.LinkEvent{})
	require.NoError(t, err)

	// messages exist, but no credits have been added yet.
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msg, err = mq.Receive(ctx, emulation.LinkEvent{})
	require.Nil(t, msg)
	require.ErrorIs(t, err, context.DeadlineExceeded)

	// now issue credit
	err = mq.IssueCredit(1, emulation.LinkEvent{})
	require.NoError(t, err)

	msg, err = mq.Receive(context.Background(), emulation.LinkEvent{})
	require.NoError(t, err)
	require.Equal(t, []byte("first message"), msg.Value)

	// and the one message has been consumed.
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msg, err = mq.Receive(ctx, emulation.LinkEvent{})
	require.Nil(t, msg)
	require.ErrorIs(t, err, context.DeadlineExceeded)

	// reissue credit
	err = mq.IssueCredit(1, emulation.LinkEvent{})
	require.NoError(t, err)

	err = mq.Send(context.Background(), &amqp.Message{Value: []byte("second message")}, emulation.LinkEvent{})
	require.NoError(t, err)

	msg, err = mq.Receive(context.Background(), emulation.LinkEvent{})
	require.NoError(t, err)
	require.Equal(t, []byte("second message"), msg.Value)
}
