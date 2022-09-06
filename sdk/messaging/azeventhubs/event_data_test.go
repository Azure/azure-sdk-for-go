// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import (
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestEventData_Annotations(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		re, err := newReceivedEventData(&amqp.Message{})
		require.NoError(t, err)

		require.Empty(t, re.Body)
		require.Nil(t, re.EnqueuedTime)
		require.Equal(t, int64(0), re.SequenceNumber)
		require.Nil(t, re.Offset)
		require.Nil(t, re.PartitionKey)
	})

	type badAnnotationValue struct {
		Name  string
		Value any
		Error string
	}

	badAnnotationValues := []badAnnotationValue{
		{Name: "x-opt-partition-key", Value: 99, Error: "partition key cannot be converted to a string"},
		{Name: "x-opt-sequence-number", Value: "101", Error: "sequence number cannot be converted to an int64"},
		{Name: "x-opt-enqueued-time", Value: "now", Error: "enqueued time cannot be converted to a time.Time"},
		{Name: "x-opt-offset", Value: 102, Error: "offset cannot be converted to an int64"},
	}

	for _, bav := range badAnnotationValues {
		t.Run(fmt.Sprintf("invalid types (%s)", bav.Name), func(t *testing.T) {
			// invalid types for properties doesn't crash us
			re, err := newReceivedEventData(&amqp.Message{
				Annotations: amqp.Annotations{
					bav.Name: bav.Value,
				},
			})

			require.Nil(t, re)
			require.EqualError(t, err, bav.Error)
		})
	}
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
			"application property 1": "application property value 1",
		},
	}

	re, err := newReceivedEventData(origAMQPMessage)
	require.NoError(t, err)

	expectedBody := [][]byte{
		[]byte("hello world"),
	}

	expectedAppProperties := map[string]any{
		"application property 1": "application property value 1",
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
		Data: [][]byte{[]byte("hello world")},
		ApplicationProperties: map[string]any{
			"application property 1": "application property value 1",
		},
	}, re.toAMQPMessage())
}

func TestEventData_newReceivedEventData_sequenceNumberPromotion(t *testing.T) {
	intValues := []any{
		int(101), int8(101), int16(101), int32(101), int64(101),
	}

	for _, iv := range intValues {
		t.Run(fmt.Sprintf("%T", iv), func(t *testing.T) {
			re, err := newReceivedEventData(&amqp.Message{
				Annotations: map[any]any{
					"x-opt-sequence-number": iv,
				},
			})

			require.NoError(t, err)
			require.Equal(t, int64(101), re.SequenceNumber)
		})
	}
}
