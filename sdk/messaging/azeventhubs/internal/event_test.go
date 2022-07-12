// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"testing"

	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestMessageConversion(t *testing.T) {
	subject := "subject"
	contentEncoding := "utf-75"
	contentType := "application/octet-stream"

	amqpMsg := &amqp.Message{
		Properties: &amqp.MessageProperties{
			MessageID:       "messageID",
			UserID:          []byte("userID"),
			CorrelationID:   "correlationID",
			Subject:         &subject,
			ContentEncoding: &contentEncoding,
			ContentType:     &contentType,
		},
		Annotations: amqp.Annotations{
			"annotation1": "annotation1Value",
			"dt-subject":  "dt-subject-value",
		},
		DeliveryAnnotations: amqp.Annotations{
			"deliveryAnnotation1": "deliveryAnnotation1Value",
		},
		ApplicationProperties: map[string]interface{}{
			"applicationProperty1": "applicationProperty1Value",
		},
		Data: [][]byte{
			[]byte("hello world"),
		},
	}

	event, err := eventFromMsg(amqpMsg)
	require.NoError(t, err)

	require.EqualValues(t, "hello world", string(event.Data))

	// AMQPMessage.Properties -> event.RawAMQPMessage (subset)
	require.EqualValues(t, "userID", string(event.RawAMQPMessage.Properties.UserID))
	require.EqualValues(t, "correlationID", event.RawAMQPMessage.Properties.CorrelationID)
	require.EqualValues(t, "subject", event.RawAMQPMessage.Properties.Subject)
	require.EqualValues(t, "utf-75", event.RawAMQPMessage.Properties.ContentEncoding)
	require.EqualValues(t, "application/octet-stream", event.RawAMQPMessage.Properties.ContentType)

	// AMQPMessage.ApplicationProperties -> Event.Properties
	require.EqualValues(t, "applicationProperty1Value", event.Properties["applicationProperty1"])

	// AMQPMessage.Annotations -> Event.SystemProperties.Annotations
	require.EqualValues(t, "annotation1Value", event.SystemProperties.Annotations["annotation1"])
	require.EqualValues(t, "dt-subject-value", event.SystemProperties.Annotations["dt-subject"])
}
