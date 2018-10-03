package servicebus

import (
	"context"
	"github.com/Azure/azure-service-bus-go/internal/test"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func (suite *serviceBusSuite) TestQueueSendReceiveWithLock() {
	tests := map[string]func(context.Context, *testing.T, *Queue){
		"SimpleSendReceiveWithLock": testQueueSendAndReceiveWithRenewLock,
	}

	ns := suite.getNewSasInstance()
	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			queueName := suite.randEntityName()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
			defer cancel()
			cleanup := makeQueue(ctx, t, ns, queueName)
			q, err := ns.NewQueue(queueName)
			suite.NoError(err)
			defer func() {
				cleanup()
			}()
			testFunc(ctx, t, q)
			q.Close(ctx)
			if !t.Failed() {
				checkZeroQueueMessages(ctx, t, ns, queueName)
			}
		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func testQueueSendAndReceiveWithRenewLock(ctx context.Context, t *testing.T, queue *Queue) {
	ttl := 5 * time.Minute
	numMessages := rand.Intn(100) + 20
	activeMessages := make([]*Message, 0, numMessages)
	expected := make(map[string]int, numMessages)
	seen := make(map[string]int, numMessages)
	errs := make(chan error, 1)

	renewEvery := time.Second * 35
	processingTime := time.Second * 240
	buffer := time.Second * 15

	t.Logf("Sending/receiving %d messages", numMessages)

	// Receiving Loop
	go func() {
		inner, cancel := context.WithCancel(ctx)
		numSeen := 0
		errs <- queue.Receive(inner, HandlerFunc(func(ctx context.Context, msg *Message) DispositionAction {
			numSeen++
			seen[string(msg.Data)]++

			activeMessages = append(activeMessages, msg)
			if numSeen >= numMessages {
				cancel()
			}
			return func(ctx context.Context) {
				//Do nothing as we want the message to remain in an uncompleted state.
			}
		}))
	}()

	// Sending Loop
	for i := 0; i < numMessages; i++ {
		payload := test.RandomString("hello", 10)
		expected[payload]++
		msg := NewMessageFromString(payload)
		msg.TTL = &ttl
		assert.NoError(t, queue.Send(ctx, msg))
	}

	// Renewal Loop
	go func() {
		assert.NoError(t, queue.RenewLocks(ctx, activeMessages))
		time.Sleep(renewEvery)
	}()

	// Wait for the all the messages to be send and recieved.
	// The renewal loop should keep locks live on the messages during this wait
	time.Sleep(processingTime + buffer)

	// Then finally accept all the messages we're holding locks on
	for _, msg := range activeMessages {
		msg.Complete()
	}

	//Check for any errors and check no messages left on the queue
	// (If there are this would mean that a lock wasn't held and the message was requeued)
	assert.EqualError(t, <-errs, context.Canceled.Error())

	assert.Equal(t, len(expected), len(seen))
	for k, v := range seen {
		assert.Equal(t, expected[k], v)
	}
}
