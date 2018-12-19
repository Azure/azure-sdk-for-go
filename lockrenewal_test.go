package servicebus

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Azure/azure-service-bus-go/internal/test"
	"github.com/stretchr/testify/assert"
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
			lockDuration := 30 * time.Second

			cleanup := makeQueue(ctx, t, ns, queueName, QueueEntityWithLockDuration(&lockDuration))
			q, err := ns.NewQueue(queueName)
			suite.NoError(err)
			defer cleanup()
			testFunc(ctx, t, q)
			suite.NoError(q.Close(ctx))
			if !t.Failed() {
				// If there are message on the queue this would mean that a lock wasn't held and the message was requeued.
				assertZeroQueueMessages(ctx, t, ns, queueName)
			}
		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func testQueueSendAndReceiveWithRenewLock(ctx context.Context, t *testing.T, queue *Queue) {
	ttl := 5 * time.Minute
	numMessages := rand.Intn(40) + 20
	activeMessages := make([]*Message, 0, numMessages)
	expected := make(map[string]int, numMessages)
	seen := make(map[string]int, numMessages)
	errs := make(chan error, 1)

	renewEvery := time.Second * 20
	processingTime := time.Second * 100

	t.Logf("Sending/receiving %d messages", numMessages)

	// Receiving Loop
	go func() {
		inner, cancel := context.WithCancel(ctx)
		numSeen := 0
		errs <- queue.Receive(inner, HandlerFunc(func(ctx context.Context, msg *Message) error {
			numSeen++
			seen[string(msg.Data)]++

			activeMessages = append(activeMessages, msg)
			if numSeen >= numMessages {
				cancel()
			}

			// Do nothing as we want the message to remain in an uncompleted state.
			return nil
		}))
	}()

	// Renewal Loop
	runRenewal := true
	go func() {
		for runRenewal {
			time.Sleep(renewEvery)
			err := queue.RenewLocks(ctx, activeMessages...)
			if err != nil {
				fmt.Println(err.Error())
			}
			// If a renewal is taking place when the test ends
			// it will fail and cause a panic without this check
			if err != nil && runRenewal {
				t.Error(err)
			}
		}
	}()

	// Sending Loop
	for i := 0; i < numMessages; i++ {
		payload := test.RandomString("hello", 10)
		expected[payload]++
		msg := NewMessageFromString(payload)
		msg.TTL = &ttl
		assert.NoError(t, queue.Send(ctx, msg))
	}

	// Wait for the all the messages to be send and recieved.
	// The renewal loop should keep locks live on the messages during this wait
	time.Sleep(processingTime)
	runRenewal = false

	// Then finally accept all the messages we're holding locks on
	for _, msg := range activeMessages {
		assert.NoError(t, msg.Complete(ctx))
	}

	//Check for any errors
	assert.EqualError(t, <-errs, context.Canceled.Error())

	assert.Equal(t, len(expected), len(seen))
	for k, v := range seen {
		assert.Equal(t, expected[k], v)
	}
}
