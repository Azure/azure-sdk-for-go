// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package mock

import (
	context "context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
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
				receiver.EXPECT().AcceptMessage(gomock.Any(), gomock.Any()).Return(nil)

				// let the caller fill in the blanks of whatever needs to happen here.
				handler(sentMessage, response)
				return response, nil
			}
		})
	}
}

// Cancelled matches context.Context instances that are cancelled.
var Cancelled gomock.Matcher = ContextCancelledMatcher{true}

// NotCancelled matches context.Context instances that are not cancelled.
var NotCancelled gomock.Matcher = ContextCancelledMatcher{false}

// NotCancelledAndHasTimeout matches context.Context instances that are not cancelled
// AND were also created from NewContextForTest.
var NotCancelledAndHasTimeout gomock.Matcher = gomock.All(ContextCancelledMatcher{false}, ContextHasTestValueMatcher{})

// CancelledAndHasTimeout matches context.Context instances that are cancelled
// AND were also created from NewContextForTest.
var CancelledAndHasTimeout gomock.Matcher = gomock.All(ContextCancelledMatcher{true}, ContextHasTestValueMatcher{})

type ContextCancelledMatcher struct {
	// WantCancelled should be set if we expect the context should
	// be cancelled. If true, we check if Err() != nil, if false we check
	// that it's nil.
	WantCancelled bool
}

// Matches returns whether x is a match.
func (m ContextCancelledMatcher) Matches(x any) bool {
	ctx := x.(context.Context)

	if m.WantCancelled {
		return ctx.Err() != nil
	} else {
		return ctx.Err() == nil
	}
}

// String describes what the matcher matches.
func (m ContextCancelledMatcher) String() string {
	return fmt.Sprintf("want cancelled:%v", m.WantCancelled)
}

type ContextHasTestValueMatcher struct{}

func (m ContextHasTestValueMatcher) Matches(x any) bool {
	ctx := x.(context.Context)
	return ctx.Value(testContextKey(0)) == "correctContextWasUsed"
}

func (m ContextHasTestValueMatcher) String() string {
	return "has test context value"
}

type testContextKey int

// NewContextWithTimeoutForTests creates a context with a lower timeout than requested just to keep
// unit test times reasonable.
//
// It validates that the passed in timeout is the actual defaultCloseTimeout and also
// adds in a testContextKey(0) as a value, which can be used to verify that the context
// has been properly propagated.
func NewContextWithTimeoutForTests(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	// (we're in the wrong package to share the value, but this is meant to match defaultCloseTimeout)
	if timeout != time.Minute {
		// panic'ing instead of require.Equal() otherwise I would need to take a 't' and not be signature
		// compatible with context.WithTimeout.
		panic(fmt.Sprintf("Incorrect close timeout: expected %s, actual %s", time.Minute, timeout))
	}

	parentWithValue := context.WithValue(parent, testContextKey(0), "correctContextWasUsed")

	// NOTE: if you're debugging then you might need to bump up this
	// value so you can single step.
	return context.WithTimeout(parentWithValue, time.Second)
}
