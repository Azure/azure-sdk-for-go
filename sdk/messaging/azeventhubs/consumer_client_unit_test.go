// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestUnitNewConsumerClient(t *testing.T) {
	t.Run("ConnectionStringNoEntityPath", func(t *testing.T) {
		connectionStringNoEntityPath := "Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>"

		client, err := NewConsumerClientFromConnectionString(connectionStringNoEntityPath, "eventHubName", DefaultConsumerGroup, nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, "eventHubName", client.eventHub)

		client, err = NewConsumerClientFromConnectionString(connectionStringNoEntityPath, "", DefaultConsumerGroup, nil)
		require.EqualError(t, err, "connection string does not contain an EntityPath. eventHub cannot be an empty string")
		require.Nil(t, client)
	})

	t.Run("ConnectionStringWithEntityPath", func(t *testing.T) {
		connectionStringWithEntityPath := "Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>;EntityPath=eventHubName"

		client, err := NewConsumerClientFromConnectionString(connectionStringWithEntityPath, "", DefaultConsumerGroup, nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, "eventHubName", client.eventHub)

		client, err = NewConsumerClientFromConnectionString(connectionStringWithEntityPath, "eventHubName", DefaultConsumerGroup, nil)
		require.EqualError(t, err, "connection string contains an EntityPath. eventHub must be an empty string")
		require.Nil(t, client)
	})

	t.Run("TokenCredential", func(t *testing.T) {
		tokenCredential := fakeTokenCredential{}
		client, err := NewConsumerClient("ripark.servicebus.windows.net", "eventHubName", DefaultConsumerGroup, tokenCredential, nil)
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, "eventHubName", client.eventHub)
	})
}

func TestUnit_getOffsetExpression(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		expr, err := getStartExpression(StartPosition{})
		require.NoError(t, err)
		require.Equal(t, "amqp.annotation.x-opt-offset > '@latest'", expr)

		expr, err = getStartExpression(StartPosition{Earliest: to.Ptr(true)})
		require.NoError(t, err)
		require.Equal(t, "amqp.annotation.x-opt-offset > '-1'", expr)

		expr, err = getStartExpression(StartPosition{Latest: to.Ptr(true)})
		require.NoError(t, err)
		require.Equal(t, "amqp.annotation.x-opt-offset > '@latest'", expr)

		expr, err = getStartExpression(StartPosition{Latest: to.Ptr(true), Inclusive: true})
		require.NoError(t, err)
		require.Equal(t, "amqp.annotation.x-opt-offset >= '@latest'", expr)

		expr, err = getStartExpression(StartPosition{Offset: to.Ptr(int64(101))})
		require.NoError(t, err)
		require.Equal(t, "amqp.annotation.x-opt-offset > '101'", expr)

		expr, err = getStartExpression(StartPosition{Offset: to.Ptr(int64(101)), Inclusive: true})
		require.NoError(t, err)
		require.Equal(t, "amqp.annotation.x-opt-offset >= '101'", expr)

		expr, err = getStartExpression(StartPosition{SequenceNumber: to.Ptr(int64(202))})
		require.NoError(t, err)
		require.Equal(t, "amqp.annotation.x-opt-sequence-number > '202'", expr)

		expr, err = getStartExpression(StartPosition{SequenceNumber: to.Ptr(int64(202)), Inclusive: true})
		require.NoError(t, err)
		require.Equal(t, "amqp.annotation.x-opt-sequence-number >= '202'", expr)

		enqueueTime, err := time.Parse(time.RFC3339, "2020-01-01T01:02:03Z")
		require.NoError(t, err)

		expr, err = getStartExpression(StartPosition{EnqueuedTime: &enqueueTime})
		require.NoError(t, err)
		require.Equal(t, "amqp.annotation.x-opt-enqueued-time > '1577840523000'", expr)

		expr, err = getStartExpression(StartPosition{EnqueuedTime: &enqueueTime, Inclusive: true})
		require.NoError(t, err)
		require.Equal(t, "amqp.annotation.x-opt-enqueued-time >= '1577840523000'", expr)
	})

	t.Run("Invalid", func(t *testing.T) {
		enqueueTime, err := time.Parse(time.RFC3339, "2020-01-01T01:02:03Z")
		require.NoError(t, err)

		expr, err := getStartExpression(StartPosition{
			EnqueuedTime: &enqueueTime,
			Offset:       to.Ptr[int64](101),
		})
		require.EqualError(t, err, "only a single start point can be set: Earliest, EnqueuedTime, Latest, Offset, or SequenceNumber")
		require.Empty(t, expr)

		expr, err = getStartExpression(StartPosition{
			Offset: to.Ptr[int64](202),
			Latest: to.Ptr(true),
		})
		require.EqualError(t, err, "only a single start point can be set: Earliest, EnqueuedTime, Latest, Offset, or SequenceNumber")
		require.Empty(t, expr)

		expr, err = getStartExpression(StartPosition{
			Latest:         to.Ptr(true),
			SequenceNumber: to.Ptr[int64](202),
		})
		require.EqualError(t, err, "only a single start point can be set: Earliest, EnqueuedTime, Latest, Offset, or SequenceNumber")
		require.Empty(t, expr)

		expr, err = getStartExpression(StartPosition{
			SequenceNumber: to.Ptr[int64](202),
			Earliest:       to.Ptr(true),
		})
		require.EqualError(t, err, "only a single start point can be set: Earliest, EnqueuedTime, Latest, Offset, or SequenceNumber")
		require.Empty(t, expr)
	})
}

type fakeTokenCredential struct {
	azcore.TokenCredential
}
