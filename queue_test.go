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
	"encoding/xml"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2015-08-01/servicebus"
	"github.com/Azure/azure-service-bus-go/atom"
	"github.com/Azure/azure-service-bus-go/internal/test"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

const (
	queueDescription1 = `
		<QueueDescription 
            xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect" 
            xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
            <LockDuration>PT1M</LockDuration>
            <MaxSizeInMegabytes>1024</MaxSizeInMegabytes>
            <RequiresDuplicateDetection>false</RequiresDuplicateDetection>
            <RequiresSession>false</RequiresSession>
            <DefaultMessageTimeToLive>P14D</DefaultMessageTimeToLive>
            <DeadLetteringOnMessageExpiration>false</DeadLetteringOnMessageExpiration>
            <DuplicateDetectionHistoryTimeWindow>PT10M</DuplicateDetectionHistoryTimeWindow>
            <MaxDeliveryCount>10</MaxDeliveryCount>
            <EnableBatchedOperations>true</EnableBatchedOperations>
            <SizeInBytes>0</SizeInBytes>
            <MessageCount>0</MessageCount>
            <IsAnonymousAccessible>false</IsAnonymousAccessible>
            <Status>Active</Status>
            <CreatedAt>2018-05-04T16:38:27.913Z</CreatedAt>
            <UpdatedAt>2018-05-04T16:38:41.897Z</UpdatedAt>
            <SupportOrdering>true</SupportOrdering>
            <AutoDeleteOnIdle>P14D</AutoDeleteOnIdle>
            <EnablePartitioning>false</EnablePartitioning>
            <EntityAvailabilityStatus>Available</EntityAvailabilityStatus>
            <EnableExpress>false</EnableExpress>
        </QueueDescription>`

	queueDescription2 = `
		<QueueDescription 
            xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect" 
            xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
            <LockDuration>PT2M</LockDuration>
            <MaxSizeInMegabytes>2048</MaxSizeInMegabytes>
            <RequiresDuplicateDetection>false</RequiresDuplicateDetection>
            <RequiresSession>false</RequiresSession>
            <DefaultMessageTimeToLive>P14D</DefaultMessageTimeToLive>
            <DeadLetteringOnMessageExpiration>true</DeadLetteringOnMessageExpiration>
            <DuplicateDetectionHistoryTimeWindow>PT20M</DuplicateDetectionHistoryTimeWindow>
            <MaxDeliveryCount>100</MaxDeliveryCount>
            <EnableBatchedOperations>true</EnableBatchedOperations>
            <SizeInBytes>256</SizeInBytes>
            <MessageCount>23</MessageCount>
            <IsAnonymousAccessible>false</IsAnonymousAccessible>
            <Status>Active</Status>
            <CreatedAt>2018-05-04T16:38:27.913Z</CreatedAt>
            <UpdatedAt>2018-05-04T16:38:41.897Z</UpdatedAt>
            <SupportOrdering>true</SupportOrdering>
            <AutoDeleteOnIdle>P14D</AutoDeleteOnIdle>
            <EnablePartitioning>true</EnablePartitioning>
            <EntityAvailabilityStatus>Available</EntityAvailabilityStatus>
            <EnableExpress>false</EnableExpress>
        </QueueDescription>`

	queueEntry1 = `
		<entry xmlns="http://www.w3.org/2005/Atom">
			<id>https://sbdjtest.servicebus.windows.net/foo</id>
			<title type="text">foo</title>
			<published>2018-05-02T20:54:59Z</published>
			<updated>2018-05-02T20:54:59Z</updated>
			<author>
				<name>sbdjtest</name>
			</author>
			<link rel="self" href="https://sbdjtest.servicebus.windows.net/foo"/>
			<content type="application/xml">` + queueDescription1 +
		`</content>
		</entry>`

	queueEntry2 = `
		<entry xmlns="http://www.w3.org/2005/Atom">
			<id>https://sbdjtest.servicebus.windows.net/bar</id>
			<title type="text">bar</title>
			<published>2018-05-02T20:54:59Z</published>
			<updated>2018-05-02T20:54:59Z</updated>
			<author>
				<name>sbdjtest</name>
			</author>
			<link rel="self" href="https://sbdjtest.servicebus.windows.net/bar"/>
			<content type="application/xml">` + queueDescription2 +
		`</content>
		</entry>`

	feedOfQueues = `
		<feed xmlns="http://www.w3.org/2005/Atom">
			<title type="text">Queues</title>
			<id>https://sbdjtest.servicebus.windows.net/$Resources/Queues</id>
			<updated>2018-05-03T00:21:15Z</updated>
			<link rel="self" href="https://sbdjtest.servicebus.windows.net/$Resources/Queues"/>` + queueEntry1 + queueEntry2 +
		`</feed>`
)

var (
	timeout = 60 * time.Second
)

func (suite *serviceBusSuite) TestQueueEntryUnmarshal() {
	var entry queueEntry
	err := xml.Unmarshal([]byte(queueEntry1), &entry)
	suite.Nil(err)
	suite.Equal("https://sbdjtest.servicebus.windows.net/foo", entry.ID)
	suite.Equal("foo", entry.Title)
	suite.Equal("sbdjtest", *entry.Author.Name)
	suite.Equal("https://sbdjtest.servicebus.windows.net/foo", entry.Link.HREF)
	suite.Equal("PT1M", *entry.Content.QueueDescription.LockDuration)
	suite.NotNil(entry.Content)
}

func (suite *serviceBusSuite) TestQueueUnmarshal() {
	var entry atom.Entry
	err := xml.Unmarshal([]byte(queueEntry1), &entry)
	assert.Nil(suite.T(), err)

	var q QueueDescription
	err = xml.Unmarshal([]byte(entry.Content.Body), &q)
	suite.Nil(err)
	suite.Equal("PT1M", *q.LockDuration)
	suite.Equal(int32(1024), *q.MaxSizeInMegabytes)
	suite.Equal(false, *q.RequiresDuplicateDetection)
	suite.Equal(false, *q.RequiresSession)
	suite.Equal("P14D", *q.DefaultMessageTimeToLive)
	suite.Equal(false, *q.DeadLetteringOnMessageExpiration)
	suite.Equal("PT10M", *q.DuplicateDetectionHistoryTimeWindow)
	suite.Equal(int32(10), *q.MaxDeliveryCount)
	suite.Equal(true, *q.EnableBatchedOperations)
	suite.Equal(int64(0), *q.SizeInBytes)
	suite.Equal(int64(0), *q.MessageCount)
	suite.EqualValues(servicebus.EntityStatusActive, *q.Status)
}

func (suite *serviceBusSuite) TestQueueManagementWrites() {
	tests := map[string]func(context.Context, *testing.T, *QueueManager, string){
		"TestPutDefaultQueue": testPutQueue,
	}

	ns := suite.getNewSasInstance()
	qm := ns.NewQueueManager()
	for name, testFunc := range tests {
		suite.T().Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			name := suite.RandomName("gosb", 6)
			testFunc(ctx, t, qm, name)
			defer suite.cleanupQueue(name)
		})
	}
}

func testPutQueue(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	q, err := qm.Put(ctx, name)
	assert.NoError(t, err)
	if assert.NotNil(t, q) {
		assert.Equal(t, name, q.Name)
	}
}

func (suite *serviceBusSuite) TestQueueManagementReads() {
	tests := map[string]func(context.Context, *testing.T, *QueueManager, []string){
		"TestGetQueue":   testGetQueue,
		"TestListQueues": testListQueues,
	}

	ns := suite.getNewSasInstance()
	qm := ns.NewQueueManager()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	names := []string{suite.randEntityName(), suite.randEntityName()}
	for _, name := range names {
		if _, err := qm.Put(ctx, name); err != nil {
			suite.T().Fatal(err)
		}
	}

	for name, testFunc := range tests {
		suite.T().Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			testFunc(ctx, t, qm, names)
		})
	}

	for _, name := range names {
		suite.cleanupQueue(name)
	}
}

func testGetQueue(ctx context.Context, t *testing.T, qm *QueueManager, names []string) {
	q, err := qm.Get(ctx, names[0])
	assert.Nil(t, err)
	assert.NotNil(t, q)
	assert.Equal(t, q.Name, names[0])
}

func testListQueues(ctx context.Context, t *testing.T, qm *QueueManager, names []string) {
	qs, err := qm.List(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, qs)
	queueNames := make([]string, len(qs))
	for idx, q := range qs {
		queueNames[idx] = q.Name
	}

	for _, name := range names {
		assert.Contains(t, queueNames, name)
	}
}

func (suite *serviceBusSuite) randEntityName() string {
	return suite.RandomName("goq", 6)
}

func (suite *serviceBusSuite) TestQueueManagement() {
	tests := map[string]func(context.Context, *testing.T, *QueueManager, string){
		"TestQueueDefaultSettings":                      testDefaultQueue,
		"TestQueueWithRequiredSessions":                 testQueueWithRequiredSessions,
		"TestQueueWithDeadLetteringOnMessageExpiration": testQueueWithDeadLetteringOnMessageExpiration,
		"TestQueueWithMaxSizeInMegabytes":               testQueueWithMaxSizeInMegabytes,
		"TestQueueWithDuplicateDetection":               testQueueWithDuplicateDetection,
		"TestQueueWithMessageTimeToLive":                testQueueWithMessageTimeToLive,
		"TestQueueWithLockDuration":                     testQueueWithLockDuration,
		"TestQueueWithAutoDeleteOnIdle":                 testQueueWithAutoDeleteOnIdle,
		"TestQueueWithPartitioning":                     testQueueWithPartitioning,
	}

	ns := suite.getNewSasInstance()
	qm := ns.NewQueueManager()
	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			name := suite.randEntityName()
			testFunc(ctx, t, qm, name)
			defer suite.cleanupQueue(name)

		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func testDefaultQueue(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	q := buildQueue(ctx, t, qm, name)
	assert.False(t, *q.EnableExpress, "should not have Express enabled")
	assert.False(t, *q.EnablePartitioning, "should not have partitioning enabled")
	assert.False(t, *q.DeadLetteringOnMessageExpiration, "should not have dead lettering on expiration")
	assert.False(t, *q.RequiresDuplicateDetection, "should not require dup detection")
	assert.False(t, *q.RequiresSession, "should not require session")
	assert.Equal(t, "P10675199DT2H48M5.4775807S", *q.AutoDeleteOnIdle, "auto delete is not 10 minutes")
	assert.Equal(t, "PT10M", *q.DuplicateDetectionHistoryTimeWindow, "dup detection is not 10 minutes")
	assert.Equal(t, int32(1*Megabytes), *q.MaxSizeInMegabytes, "max size in MBs")
	assert.Equal(t, "P10675199DT2H48M5.4775807S", *q.DefaultMessageTimeToLive, "default TTL")
	assert.Equal(t, "PT1M", *q.LockDuration, "lock duration")
}

func testQueueWithAutoDeleteOnIdle(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	window := time.Duration(20 * time.Minute)
	q := buildQueue(ctx, t, qm, name, QueueEntityWithAutoDeleteOnIdle(&window))
	assert.Equal(t, "PT20M", *q.AutoDeleteOnIdle)
}

func testQueueWithRequiredSessions(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	q := buildQueue(ctx, t, qm, name, QueueEntityWithRequiredSessions())
	assert.True(t, *q.RequiresSession)
}

func testQueueWithDeadLetteringOnMessageExpiration(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	q := buildQueue(ctx, t, qm, name, QueueEntityWithDeadLetteringOnMessageExpiration())
	assert.True(t, *q.DeadLetteringOnMessageExpiration)
}

func testQueueWithPartitioning(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	q := buildQueue(ctx, t, qm, name, QueueEntityWithPartitioning())
	assert.True(t, *q.EnablePartitioning)
}

func testQueueWithMaxSizeInMegabytes(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	size := 3 * Megabytes
	q := buildQueue(ctx, t, qm, name, QueueEntityWithMaxSizeInMegabytes(size))
	assert.Equal(t, int32(size), *q.MaxSizeInMegabytes)
}

func testQueueWithDuplicateDetection(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	window := time.Duration(20 * time.Minute)
	q := buildQueue(ctx, t, qm, name, QueueEntityWithDuplicateDetection(&window))
	assert.Equal(t, "PT20M", *q.DuplicateDetectionHistoryTimeWindow)
}

func testQueueWithMessageTimeToLive(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	window := time.Duration(10 * 24 * 60 * time.Minute)
	q := buildQueue(ctx, t, qm, name, QueueEntityWithMessageTimeToLive(&window))
	assert.Equal(t, "P10D", *q.DefaultMessageTimeToLive)
}

func testQueueWithLockDuration(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	window := time.Duration(3 * time.Minute)
	q := buildQueue(ctx, t, qm, name, QueueEntityWithLockDuration(&window))
	assert.Equal(t, "PT3M", *q.LockDuration)
}

func buildQueue(ctx context.Context, t *testing.T, qm *QueueManager, name string, opts ...QueueManagementOption) *QueueEntity {
	_, err := qm.Put(ctx, name, opts...)
	if !assert.NoError(t, err) {
		assert.FailNow(t, fmt.Sprintf("%v", err))
	}

	q, err := qm.Get(ctx, name)
	if !assert.NoError(t, err) {
		assert.FailNow(t, fmt.Sprintf("%v", err))
	}

	return q
}

func (suite *serviceBusSuite) TestQueueClient() {
	tests := map[string]func(context.Context, *testing.T, *Queue){
		"SimpleSend":              testQueueSend,
		"SendAndReceiveInOrder":   testQueueSendAndReceiveInOrder,
		"DuplicateDetection":      testDuplicateDetection,
		"MessageProperties":       testMessageProperties,
		"Retry":                   testRequeueOnFail,
		"ReceiveOne":              testReceiveOne,
		"SendAndReceiveScheduled": testQueueSendAndReceiveScheduled,
	}

	ns := suite.getNewSasInstance()
	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			queueName := suite.randEntityName()
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			window := time.Duration(30 * time.Second)
			cleanup := makeQueue(ctx, t, ns, queueName,
				QueueEntityWithPartitioning(),
				QueueEntityWithDuplicateDetection(&window))
			q, err := ns.NewQueue(ctx, queueName)
			suite.NoError(err)
			defer func() {
				q.Close(ctx)
				cleanup()
			}()
			testFunc(ctx, t, q)
			if !t.Failed() && name != "SimpleSend" {
				checkZeroQueueMessages(ctx, t, ns, queueName)
			}
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func testReceiveOne(ctx context.Context, t *testing.T, q *Queue) {
	if assert.NoError(t, q.Send(ctx, NewMessageFromString("Hello World!"))) {
		messageWithContext, err := q.ReceiveOne(ctx)
		if assert.NoError(t, err) {
			span, _ := opentracing.StartSpanFromContext(messageWithContext.Ctx, "continue-message-span")
			defer span.Finish()

			messageWithContext.Complete()
		}
	}
}

func testRequeueOnFail(ctx context.Context, t *testing.T, q *Queue) {
	if assert.NoError(t, q.Send(ctx, NewMessageFromString("Hello World!"))) {
		var wg sync.WaitGroup
		wg.Add(2)
		var receivedMsg *Message
		fail := true
		listenHandle, err := q.Receive(context.Background(),
			func(ctx context.Context, msg *Message) DispositionAction {
				receivedMsg = msg
				defer func() {
					wg.Done()
				}()
				if fail {
					fail = false
					return msg.Abandon()
				}
				return msg.Complete()
			})
		if assert.NoError(t, err) {
			defer listenHandle.Close(ctx)
			end, _ := ctx.Deadline()
			waitUntil(t, &wg, time.Until(end))

			if assert.NoError(t, err) {
				assert.EqualValues(t, 2, receivedMsg.DeliveryCount)
			}
		}
	}
}

func testMessageProperties(ctx context.Context, t *testing.T, q *Queue) {
	if assert.NoError(t, q.Send(ctx, NewMessageFromString("Hello World!"))) {
		var wg sync.WaitGroup
		wg.Add(1)
		var receivedMsg *Message
		listenHandle, err := q.Receive(context.Background(),
			func(ctx context.Context, msg *Message) DispositionAction {
				receivedMsg = msg
				defer func() {
					wg.Done()
				}()
				return msg.Complete()
			})
		if assert.NoError(t, err) {
			defer listenHandle.Close(ctx)
			end, _ := ctx.Deadline()
			waitUntil(t, &wg, time.Until(end))

			if assert.NoError(t, err) {
				sp := receivedMsg.SystemProperties
				assert.NotNil(t, sp.LockedUntil, "LockedUntil")
				assert.NotNil(t, sp.EnqueuedSequenceNumber, "EnqueuedSequenceNumber")
				assert.NotNil(t, sp.EnqueuedTime, "EnqueuedTime")
				assert.NotNil(t, sp.SequenceNumber, "SequenceNumber")
				assert.NotNil(t, sp.PartitionID, "PartitionID")
				assert.NotNil(t, sp.PartitionKey, "PartitionKey")
			}
		}
	}
}

func testQueueSend(ctx context.Context, t *testing.T, queue *Queue) {
	err := queue.Send(ctx, NewMessageFromString("hello!"))
	assert.Nil(t, err)
}

func testQueueSendAndReceiveInOrder(ctx context.Context, t *testing.T, queue *Queue) {
	numMessages := rand.Intn(100) + 20
	messages := make([]string, numMessages)
	for i := 0; i < numMessages; i++ {
		messages[i] = test.RandomString("hello", 10)
	}

	for _, message := range messages {
		err := queue.Send(ctx, NewMessageFromString(message))
		if err != nil {
			t.Fatal(err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(numMessages)
	// ensure in-order processing of messages from the queue
	count := 0
	listener, err := queue.Receive(ctx, func(ctx context.Context, event *Message) DispositionAction {
		assert.Equal(t, messages[count], string(event.Data))
		count++
		wg.Done()
		return event.Complete()
	})
	if assert.NoError(t, err) {
		defer listener.Close(ctx)
		end, _ := ctx.Deadline()
		waitUntil(t, &wg, time.Until(end))
	}
}

func testQueueSendAndReceiveScheduled(ctx context.Context, t *testing.T, queue *Queue) {
	msg := NewMessageFromString("to the future!!")
	futureTime := time.Now().UTC().Add(15 * time.Second)
	msg.SystemProperties = &SystemProperties{
		ScheduledEnqueueTime: &futureTime,
	}
	if assert.NoError(t, queue.Send(ctx, msg)) {
		var wg sync.WaitGroup
		wg.Add(1)
		listener, err := queue.Receive(ctx, func(ctx context.Context, received *Message) DispositionAction {
			defer wg.Done()
			arrivalTime := time.Now().UTC()
			assert.True(t, arrivalTime.After(futureTime))
			return received.Complete()
		})
		if assert.NoError(t, err) {
			defer listener.Close(ctx)
			end, _ := ctx.Deadline()
			waitUntil(t, &wg, time.Until(end))
		}
	}
}

func testDuplicateDetection(ctx context.Context, t *testing.T, queue *Queue) {
	messages := []*Message{
		NewMessageFromString("hello, "),
		NewMessageFromString("world!"),
	}

	messages[0].ID = "foo"
	messages[1].ID = "bar"

	for _, msg := range messages {
		if !assert.NoError(t, queue.Send(ctx, msg)) {
			t.FailNow()
		}
	}

	// send dup
	if !assert.NoError(t, queue.Send(ctx, messages[0])) {
		t.FailNow()
	}

	var wg sync.WaitGroup
	wg.Add(2)
	received := make(map[interface{}]string)
	var all []*Message
	queue.Receive(ctx, func(ctx context.Context, message *Message) DispositionAction {
		all = append(all, message)
		if _, ok := received[message.ID]; !ok {
			// caught a new one
			defer wg.Done()
			received[message.ID] = string(message.Data)
		} else {
			// caught a dup
			assert.Fail(t, "received a duplicate message")
			for _, item := range all {
				t.Logf("mID: %q, gID: %q, gSeq: %q, lockT: %q", item.ID, *item.GroupID, *item.GroupSequence, *item.LockToken)
			}
		}
		return message.Complete()
	})
	end, _ := ctx.Deadline()
	waitUntil(t, &wg, time.Until(end))
}

func (suite *serviceBusSuite) TestQueueWithReceiveAndDelete() {
	tests := map[string]func(context.Context, *testing.T, *Queue){
		"SimpleSendAndReceive": testQueueSendAndReceiveWithReceiveAndDelete,
	}

	ns := suite.getNewSasInstance()
	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			queueName := suite.randEntityName()
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			cleanup := makeQueue(ctx, t, ns, queueName)
			q, err := ns.NewQueue(ctx, queueName, QueueWithReceiveAndDelete())
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

func testQueueSendAndReceiveWithReceiveAndDelete(ctx context.Context, t *testing.T, queue *Queue) {
	numMessages := rand.Intn(100) + 20
	messages := make([]string, numMessages)
	for i := 0; i < numMessages; i++ {
		messages[i] = test.RandomString("hello", 10)
	}

	for _, message := range messages {
		if !assert.NoError(t, queue.Send(ctx, NewMessageFromString(message))) {
			assert.FailNow(t, "failed to send message")
		}
	}

	var wg sync.WaitGroup
	wg.Add(numMessages)
	count := 0
	queue.Receive(ctx, func(ctx context.Context, msg *Message) DispositionAction {
		assert.Equal(t, messages[count], string(msg.Data))
		count++
		wg.Done()
		return nil
	})
	end, _ := ctx.Deadline()
	waitUntil(t, &wg, time.Until(end))
}

//func (suite *serviceBusSuite) TestQueueWithRequiredSessions() {
//	tests := map[string]func(context.Context, *testing.T, *Queue){
//		"TestSendAndReceiveSession": testQueueWithRequiredSessionSendAndReceive,
//	}
//
//	for name, testFunc := range tests {
//		setupTestTeardown := func(t *testing.T) {
//			ns := suite.getNewSasInstance()
//			queueName := suite.randEntityName()
//			sessionID := mustUUID(t).String()
//			ctx, cancel := context.WithTimeout(context.Background(), timeout)
//			cleanup := makeQueue(ctx, t, ns, queueName,
//				QueueEntityWithPartitioning(),
//				QueueEntityWithRequiredSessions())
//			q, err := ns.NewQueue(ctx, queueName, QueueWithRequiredSession(sessionID))
//			if suite.NoError(err) {
//				testFunc(ctx, t, q)
//				if !t.Failed() {
//					checkZeroQueueMessages(ctx, t, ns, queueName)
//				}
//			}
//			defer func() {
//				if q != nil {
//					q.Close(ctx)
//				}
//				cancel()
//				cleanup()
//			}()
//		}
//
//		suite.T().Run(name, setupTestTeardown)
//	}
//}
//
//func testQueueWithRequiredSessionSendAndReceive(ctx context.Context, t *testing.T, queue *Queue) {
//	numMessages := rand.Intn(100) + 20
//	messages := make([]string, numMessages)
//	for i := 0; i < numMessages; i++ {
//		messages[i] = test.RandomString("hello", 10)
//	}
//
//	for _, message := range messages {
//		err := queue.Send(ctx, NewMessageFromString(message))
//		if err != nil {
//			t.Fatal(err)
//		}
//	}
//
//	var wg sync.WaitGroup
//	wg.Add(numMessages)
//	// ensure in-order processing of messages from the queue
//	count := 0
//	handler := func(ctx context.Context, event *Message) DispositionAction {
//		if !assert.Equal(t, messages[count], string(event.Data)) {
//			assert.FailNow(t, fmt.Sprintf("message %d %q didn't match %q", count, messages[count], string(event.Data)))
//		}
//		count++
//		wg.Done()
//		return event.Complete()
//	}
//	listenHandle, err := queue.Receive(ctx, handler)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer func() {
//		ctx, cancel := context.WithTimeout(context.Background(), timeout)
//		defer cancel()
//		listenHandle.Close(ctx)
//	}()
//
//	end, _ := ctx.Deadline()
//	waitUntil(t, &wg, time.Until(end))
//}
//
//func mustUUID(t *testing.T) uuid.UUID {
//	id, err := uuid.NewV4()
//	if err != nil {
//		t.Fatal(err)
//	}
//	return id
//}

func makeQueue(ctx context.Context, t *testing.T, ns *Namespace, name string, opts ...QueueManagementOption) func() {
	qm := ns.NewQueueManager()
	entity, err := qm.Get(ctx, name)
	if !assert.NoError(t, err) {
		assert.FailNow(t, "could not GET a subscription")
	}

	if entity == nil {
		entity, err = qm.Put(ctx, name, opts...)
		if !assert.NoError(t, err) {
			assert.FailNow(t, "could not PUT a subscription")
		}
	}
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		_ = qm.Delete(ctx, entity.Name)
	}
}

func checkZeroQueueMessages(ctx context.Context, t *testing.T, ns *Namespace, name string) {
	qm := ns.NewQueueManager()
	maxTries := 10
	for i := 0; i < maxTries; i++ {
		q, err := qm.Get(ctx, name)
		if !assert.NoError(t, err) {
			return
		}
		if *q.MessageCount == 0 {
			return
		}
		t.Logf("try %d out of %d, message count was %d, not 0", i+1, maxTries, *q.MessageCount)
		time.Sleep(1 * time.Second)
	}

	assert.Fail(t, "message count never reached zero")
}

func waitUntil(t *testing.T, wg *sync.WaitGroup, d time.Duration) {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(d):
		t.Fatal("took longer than " + fmtDuration(d))
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second) / time.Second
	return fmt.Sprintf("%d seconds", d)
}

func (suite *serviceBusSuite) cleanupQueue(name string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ns := suite.getNewSasInstance()
	qm := ns.NewQueueManager()
	err := qm.Delete(ctx, name)
	if err != nil {
		suite.T().Fatal(err)
	}
}
