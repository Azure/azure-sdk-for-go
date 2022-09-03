// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestEventData_Annotations(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		re := newReceivedEventData(&amqp.Message{})

		require.Empty(t, re.Body)
		require.Nil(t, re.EnqueuedTime)
		require.Equal(t, int64(0), re.SequenceNumber)
		require.Nil(t, re.Offset)
		require.Nil(t, re.PartitionKey)
	})

	t.Run("invalid types", func(t *testing.T) {
		// invalid types for properties doesn't crash us
		re := newReceivedEventData(&amqp.Message{
			Annotations: amqp.Annotations{
				"x-opt-partition-key":   99,
				"x-opt-sequence-number": "101",
				"x-opt-offset":          102,
				"x-opt-enqueued-time":   "now",
			},
		})

		require.Empty(t, re.Body)
		require.Nil(t, re.EnqueuedTime)
		require.Equal(t, int64(0), re.SequenceNumber)
		require.Nil(t, re.Offset)
		require.Nil(t, re.PartitionKey)
	})
}

func TestEventData_newReceivedEventData(t *testing.T) {
	now := time.Now().UTC()

	origAMQPMessage := &amqp.Message{
		Properties: &amqp.MessageProperties{
			ContentType:   to.Ptr("content type"),
			MessageID:     "message id",
			CorrelationID: to.Ptr("correlation id"),
		},
		Data: [][]byte{[]byte("hello world")},
		Annotations: map[any]any{
			"hello":                 "world",
			5:                       "ignored",
			"x-opt-partition-key":   "partition key",
			"x-opt-sequence-number": int64(101),
			"x-opt-offset":          "102",
			"x-opt-enqueued-time":   now,
		},
		ApplicationProperties: map[string]any{
			"application property 1": "application prioperty value 1",
		},
	}

	re := newReceivedEventData(origAMQPMessage)

	expectedBody := [][]byte{
		[]byte("hello world"),
	}

	expectedAppProperties := map[string]any{
		"application property 1": "application prioperty value 1",
	}

	require.Equal(t, &ReceivedEventData{
		EventData: EventData{
			Body:          expectedBody[0],
			ContentType:   to.Ptr("content type"),
			CorrelationID: to.Ptr("correlation id"),
			MessageID:     to.Ptr("message id"),
			Properties:    expectedAppProperties,
		},
		EnqueuedTime:   &now,
		SequenceNumber: 101,
		SystemProperties: map[string]any{
			"hello": "world",
		},
		Offset:       to.Ptr[int64](102),
		PartitionKey: to.Ptr("partition key"),
	}, re)

	require.Equal(t, &amqp.Message{
		Properties: &amqp.MessageProperties{
			ContentType:   to.Ptr("content type"),
			MessageID:     "message id",
			CorrelationID: to.Ptr("correlation id"),
		},
		Data:        [][]byte{[]byte("hello world")},
		Annotations: map[any]any{},
		ApplicationProperties: map[string]any{
			"application property 1": "application prioperty value 1",
		},
	}, re.toAMQPMessage())
}
