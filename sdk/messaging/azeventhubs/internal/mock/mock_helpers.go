// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package mock

import (
	context "context"

	"github.com/Azure/go-amqp"
	gomock "github.com/golang/mock/gomock"
)

func SetupRPC(sender *MockAMQPSenderCloser, receiver *MockAMQPReceiverCloser, expectedCount int, handler func(sent *amqp.Message, response *amqp.Message)) {
	// this is an RPC pattern - when we send a message we give it a message ID, and the
	// response comes back with a correlation ID filled out, so you can match requests
	// to responses.
	ch := make(chan *amqp.Message, 1000)

	for i := 0; i < expectedCount; i++ {
		sender.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Nil()).Do(func(ctx context.Context, msg *amqp.Message, o *amqp.SendOptions) error {
			ch <- msg
			return nil
		})
	}

	// RPC loops forever. We get one extra Receive() call here (the one that waits on the ctx.Done())
	for i := 0; i < expectedCount+1; i++ {
		receiver.EXPECT().Receive(gomock.Any(), gomock.Nil()).DoAndReturn(func(ctx context.Context, o *amqp.ReceiveOptions) (*amqp.Message, error) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case sentMessage := <-ch:
				response := &amqp.Message{
					// this is how RPC responses are correlated with their
					// sent messages.
					Properties: &amqp.MessageProperties{
						CorrelationID: sentMessage.Properties.MessageID,
					},
				}
				// let the caller fill in the blanks of whatever needs to happen here.
				handler(sentMessage, response)
				return response, nil
			}
		})
	}
}
