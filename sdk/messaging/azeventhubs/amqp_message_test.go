// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAMQPAnnotatedMessageUnitTest(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		msg := &AMQPAnnotatedMessage{}
		amqpMessage := msg.toAMQPMessage()

		// we duplicate/inflate these since we modify them
		// in various parts of the API.
		require.NotNil(t, amqpMessage.Properties)
		require.NotNil(t, amqpMessage.Annotations)
	})
}
