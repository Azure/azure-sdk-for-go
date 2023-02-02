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
	events := emulation.NewEvents()
	q := emulation.NewQueue("entity", events)
	defer q.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// no messages exist yet
	msg, err := q.Receive(ctx, emulation.LinkEvent{}, nil)
	require.Nil(t, msg)
	require.ErrorIs(t, err, context.DeadlineExceeded)

	err = q.Send(context.Background(), &amqp.Message{
		Value: []byte("first message"),
	}, emulation.LinkEvent{}, nil)
	require.NoError(t, err)

	// messages exist, but no credits have been added yet.
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msg, err = q.Receive(ctx, emulation.LinkEvent{}, nil)
	require.Nil(t, msg)
	require.ErrorIs(t, err, context.DeadlineExceeded)

	// now issue credit
	err = q.IssueCredit(1, emulation.LinkEvent{}, nil)
	require.NoError(t, err)

	msg, err = q.Receive(context.Background(), emulation.LinkEvent{}, nil)
	require.NoError(t, err)
	require.Equal(t, []byte("first message"), msg.Value)

	// and the one message has been consumed.
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	msg, err = q.Receive(ctx, emulation.LinkEvent{}, nil)
	require.Nil(t, msg)
	require.ErrorIs(t, err, context.DeadlineExceeded)

	// reissue credit
	err = q.IssueCredit(1, emulation.LinkEvent{}, nil)
	require.NoError(t, err)

	err = q.Send(context.Background(), &amqp.Message{Value: []byte("second message")}, emulation.LinkEvent{}, nil)
	require.NoError(t, err)

	msg, err = q.Receive(context.Background(), emulation.LinkEvent{}, nil)
	require.NoError(t, err)
	require.Equal(t, []byte("second message"), msg.Value)
}
