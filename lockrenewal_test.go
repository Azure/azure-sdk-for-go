package servicebus

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

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Azure/azure-service-bus-go/internal/test"
)

func (suite *serviceBusSuite) TestQueueSendReceiveWithLock() {
	tests := map[string]func(context.Context, *testing.T, *Queue, int){
		"SimpleSendReceiveWithLock": testQueueSendAndReceiveWithRenewLock,
	}

	ns := suite.getNewSasInstance()
	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			queueName := suite.randEntityName()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
			defer cancel()
			lockDuration := 10 * time.Second

			cleanup := makeQueue(ctx, t, ns, queueName, QueueEntityWithLockDuration(&lockDuration))
			numMessages := rand.Intn(40) + 10
			q, err := ns.NewQueue(queueName, QueueWithPrefetchCount(uint32(numMessages)))
			suite.NoError(err)
			defer cleanup()
			testFunc(ctx, t, q, numMessages)
			suite.NoError(q.Close(ctx))
			if !t.Failed() {
				// If there are message on the queue this would mean that a lock wasn't held and the message was requeued.
				assertZeroQueueMessages(ctx, t, ns, queueName)
			}
		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func testQueueSendAndReceiveWithRenewLock(ctx context.Context, t *testing.T, queue *Queue, numMessages int) {
	ttl := 5 * time.Minute
	activeMessages := make([]*Message, 0, numMessages)
	expected := make(map[string]int, numMessages)
	seen := make(map[string]int, numMessages)
	errs := make(chan error, 1)

	renewEvery := time.Second * 6
	// all are held in memory. we wait a bit longer than the lock expiry set to 10s
	processingTime := 15 * time.Second
	t.Logf("processing time : %s", processingTime)
	t.Logf("processing time : %s \n", processingTime)
	t.Logf("Sending/receiving %d messages", numMessages)

	// Receiving Loop
	go func() {
		inner, cancel := context.WithCancel(ctx)
		numSeen := 0
		errs <- queue.Receive(inner, HandlerFunc(func(ctx context.Context, msg *Message) error {
			numSeen++
			seen[string(msg.Data)]++
			t.Logf("handling message %d - %s \n", numSeen, msg.LockToken)
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
			// If a renewal is taking place when the test ends
			// it will fail and cause a panic without this check
			if err != nil && runRenewal {
				t.Error(err)
			}
			t.Logf("renewed locks successfuly for %d messages\n", len(activeMessages))
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
		t.Logf("completing %d messages", len(activeMessages))
		assert.NoError(t, msg.Complete(ctx))
	}

	//Check for any errors
	assert.EqualError(t, <-errs, context.Canceled.Error())

	assert.Equal(t, len(expected), len(seen))
	for k, v := range seen {
		assert.Equal(t, expected[k], v)
	}
}
