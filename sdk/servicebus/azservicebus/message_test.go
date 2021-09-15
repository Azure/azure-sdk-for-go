package azservicebus

import (
	"testing"

	"github.com/Azure/azure-amqp-common-go/v3/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/stretchr/testify/require"
)

func TestMessageUnitTest(t *testing.T) {
	t.Run("toAMQPMessage", func(t *testing.T) {
		message := &Message{}

		// basic thing - it's totally fine to send a message nothing in it.
		amqpMessage, err := message.toAMQPMessage()
		require.NoError(t, err)
		require.Empty(t, amqpMessage.Annotations)
		require.NotEmpty(t, amqpMessage.Properties.MessageID, "MessageID is (currently) automatically filled out if you don't specify one")

		message = &Message{
			ID:                      "message id",
			Body:                    []byte("the body"),
			PartitionKey:            to.StringPtr("partition key"),
			TransactionPartitionKey: to.StringPtr("via partition key"),
			SessionID:               "session id",
		}

		amqpMessage, err = message.toAMQPMessage()
		require.NoError(t, err)

		require.EqualValues(t, "message id", amqpMessage.Properties.MessageID)
		require.EqualValues(t, "session id", amqpMessage.Properties.GroupID)

		require.EqualValues(t, "the body", string(amqpMessage.Data[0]))
		require.EqualValues(t, 1, len(amqpMessage.Data))

		require.EqualValues(t, map[interface{}]interface{}{
			annotationPartitionKey:    "partition key",
			annotationViaPartitionKey: "via partition key",
		}, amqpMessage.Annotations)
	})

	t.Run("convertLegacyMessage", func(t *testing.T) {
		receivedMessage := convertToReceivedMessage(&internal.Message{})

		require.EqualValues(t, "", receivedMessage.LockToken)
		require.Nil(t, receivedMessage.SequenceNumber)
		require.Nil(t, receivedMessage.PartitionKey)
		require.Nil(t, receivedMessage.TransactionPartitionKey) // `ViaPartitionKey`

		uuid, err := uuid.NewV4()
		require.NoError(t, err)

		receivedMessage = convertToReceivedMessage(&internal.Message{
			LockToken: &uuid,
			SystemProperties: &internal.SystemProperties{
				SequenceNumber:  to.Int64Ptr(111),
				PartitionID:     toInt16Ptr(101),
				ViaPartitionKey: to.StringPtr("via partition key"),
				PartitionKey:    to.StringPtr("partition key"),
			},
		})

		require.EqualValues(t, uuid.String(), receivedMessage.LockToken)
		require.EqualValues(t, 111, *receivedMessage.SequenceNumber)
		require.EqualValues(t, "partition key", *receivedMessage.PartitionKey)
		require.EqualValues(t, "via partition key", *receivedMessage.TransactionPartitionKey) // `ViaPartitionKey`
	})
}

func toInt16Ptr(i int16) *int16 {
	return &i
}
