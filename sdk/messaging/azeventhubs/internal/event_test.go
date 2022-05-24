package eventhub

import (
	"testing"

	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

//	MIT License
//
//	Copyright (c) Microsoft Corporation. All rights reserved.
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE

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
