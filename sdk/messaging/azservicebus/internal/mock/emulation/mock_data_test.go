// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation_test

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/mock/emulation"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestNewConnection(t *testing.T) {
	md := emulation.NewMockData(t, nil)
	defer md.Close()

	client, err := md.NewConnection(context.Background())
	require.NoError(t, err)

	sess, err := client.NewSession(context.Background(), nil)
	require.NoError(t, err)

	rcvr, err := sess.NewReceiver(context.Background(), "testqueue", &amqp.ReceiverOptions{
		Credit: -1,
	})
	require.NoError(t, err)

	sender, err := sess.NewSender(context.Background(), "testqueue", nil)
	require.NoError(t, err)

	err = sender.Send(context.Background(), &amqp.Message{
		Value: []byte("hello"),
	}, nil)
	require.NoError(t, err)

	err = rcvr.IssueCredit(1)
	require.NoError(t, err)

	msg, err := rcvr.Receive(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, []byte("hello"), msg.Value)
}

func TestClosingConnectionClosesChildren(t *testing.T) {
	md := emulation.NewMockData(t, nil)
	defer md.Close()

	client, err := md.NewConnection(context.Background())
	require.NoError(t, err)

	sess, err := client.NewSession(context.Background(), nil)
	require.NoError(t, err)

	rcvr, err := sess.NewReceiver(context.Background(), "testqueue", nil)
	require.NoError(t, err)

	sender, err := sess.NewSender(context.Background(), "testqueue", nil)
	require.NoError(t, err)

	err = client.Close()
	require.NoError(t, err)

	// TODO: non-deterministic
	time.Sleep(time.Second)

	msg, err := rcvr.Receive(context.Background(), nil)
	require.Nil(t, msg)

	var connErr *amqp.ConnError
	require.ErrorAs(t, err, &connErr)

	err = sender.Send(context.Background(), &amqp.Message{}, nil)
	require.ErrorAs(t, err, &connErr)
}
