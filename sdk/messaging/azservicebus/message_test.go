// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestMessageUnitTest(t *testing.T) {
	t.Run("toAMQPMessage", func(t *testing.T) {
		message := &Message{}

		// basic thing - it's totally fine to send a message nothing in it.
		amqpMessage := message.toAMQPMessage()
		require.Empty(t, amqpMessage.Annotations)
		require.Nil(t, amqpMessage.Properties.MessageID)

		scheduledEnqueuedTime := time.Now()

		message = &Message{
			MessageID:               to.StringPtr("message id"),
			Body:                    []byte("the body"),
			PartitionKey:            to.StringPtr("partition key"),
			TransactionPartitionKey: to.StringPtr("via partition key"),
			SessionID:               to.StringPtr("session id"),
			ScheduledEnqueueTime:    &scheduledEnqueuedTime,
		}

		amqpMessage = message.toAMQPMessage()

		require.EqualValues(t, "message id", amqpMessage.Properties.MessageID)
		require.EqualValues(t, "session id", *amqpMessage.Properties.GroupID)

		require.EqualValues(t, "the body", string(amqpMessage.Data[0]))
		require.EqualValues(t, 1, len(amqpMessage.Data))

		require.EqualValues(t, map[interface{}]interface{}{
			partitionKeyAnnotation:          "partition key",
			viaPartitionKeyAnnotation:       "via partition key",
			scheduledEnqueuedTimeAnnotation: scheduledEnqueuedTime,
		}, amqpMessage.Annotations)
	})
}

func TestAMQPMessageToReceivedMessage(t *testing.T) {
	t.Run("empty_message", func(t *testing.T) {
		// nothing should blow up.
		rm := newReceivedMessage(&amqp.Message{})
		require.NotNil(t, rm)
	})

	t.Run("annotations", func(t *testing.T) {
		lockedUntil := time.Now().Add(time.Hour)
		enqueuedTime := time.Now().Add(2 * time.Hour)
		scheduledEnqueueTime := time.Now().Add(3 * time.Hour)

		amqpMessage := &amqp.Message{
			Annotations: map[interface{}]interface{}{
				"x-opt-locked-until":            lockedUntil,
				"x-opt-sequence-number":         int64(101),
				"x-opt-partition-key":           "partitionKey1",
				"x-opt-enqueued-time":           enqueuedTime,
				"x-opt-scheduled-enqueue-time":  scheduledEnqueueTime,
				"x-opt-enqueue-sequence-number": int64(202),
			},
		}

		receivedMessage := newReceivedMessage(amqpMessage)

		require.EqualValues(t, lockedUntil, *receivedMessage.LockedUntil)
		require.EqualValues(t, int64(101), *receivedMessage.SequenceNumber)
		require.EqualValues(t, "partitionKey1", *receivedMessage.PartitionKey)
		require.EqualValues(t, enqueuedTime, *receivedMessage.EnqueuedTime)
		require.EqualValues(t, scheduledEnqueueTime, *receivedMessage.ScheduledEnqueueTime)
		require.EqualValues(t, int64(202), *receivedMessage.EnqueuedSequenceNumber)
	})
}

func TestAMQPMessageToMessage(t *testing.T) {
	d := 30 * time.Second
	until := time.Now().Add(d)
	pID := int16(12)

	// ServiceBus encoded the lock token in .Net's serialisation format but requries it to submitted in
	// amqps (RFC 4122) format. These are both the same GUID encoded in both formats and are used to
	// test the conversion occurs correctly.
	dotNetEncodedLockTokenGUID := []byte{205, 89, 49, 187, 254, 253, 77, 205, 162, 38, 172, 76, 45, 235, 91, 225}

	groupSequence := uint32(1)

	amqpMsg := &amqp.Message{
		DeliveryTag: dotNetEncodedLockTokenGUID,
		Properties: &amqp.MessageProperties{
			MessageID:          "messageID",
			To:                 to.StringPtr("to"),
			Subject:            to.StringPtr("subject"),
			ReplyTo:            to.StringPtr("replyTo"),
			ReplyToGroupID:     to.StringPtr("replyToGroupID"),
			CorrelationID:      "correlationID",
			ContentType:        to.StringPtr("contentType"),
			ContentEncoding:    to.StringPtr("contentEncoding"),
			AbsoluteExpiryTime: &until,
			CreationTime:       &until,
			GroupID:            to.StringPtr("groupID"),
			GroupSequence:      &groupSequence,
		},
		Annotations: amqp.Annotations{
			"x-opt-locked-until":            until,
			"x-opt-sequence-number":         int64(1),
			"x-opt-partition-id":            pID,
			"x-opt-partition-key":           "key",
			"x-opt-enqueued-time":           until,
			"x-opt-deadletter-source":       "deadLetterSource",
			"x-opt-scheduled-enqueue-time":  until,
			"x-opt-enqueue-sequence-number": int64(1),
			"x-opt-via-partition-key":       "via",
			"custom-annotation":             "value",
		},
		ApplicationProperties: map[string]interface{}{
			"test": "foo",
		},
		Header: &amqp.MessageHeader{
			TTL: d,
		},
		Data: [][]byte{[]byte("foo")},
	}

	msg := newReceivedMessage(amqpMsg)

	require.EqualValues(t, msg.MessageID, amqpMsg.Properties.MessageID, "messageID")
	require.EqualValues(t, msg.SessionID, amqpMsg.Properties.GroupID, "groupID")
	require.EqualValues(t, msg.ContentType, amqpMsg.Properties.ContentType, "contentType")
	require.EqualValues(t, *msg.CorrelationID, amqpMsg.Properties.CorrelationID, "correlation")
	require.EqualValues(t, msg.ReplyToSessionID, amqpMsg.Properties.ReplyToGroupID, "replyToGroupID")
	require.EqualValues(t, msg.ReplyTo, amqpMsg.Properties.ReplyTo, "replyTo")
	require.EqualValues(t, *msg.TimeToLive, amqpMsg.Header.TTL, "ttl")
	require.EqualValues(t, msg.Subject, amqpMsg.Properties.Subject, "subject")
	require.EqualValues(t, msg.To, amqpMsg.Properties.To, "to")

	body, err := msg.Body()
	require.NoError(t, err)
	require.EqualValues(t, body, amqpMsg.Data[0], "data")

	expectedAMQPEncodedLockTokenGUID := [16]byte{187, 49, 89, 205, 253, 254, 205, 77, 162, 38, 172, 76, 45, 235, 91, 225}

	require.EqualValues(t, msg.LockToken, expectedAMQPEncodedLockTokenGUID, "locktoken")

	require.EqualValues(t, map[string]interface{}{
		"test": "foo",
	}, msg.ApplicationProperties)
}
