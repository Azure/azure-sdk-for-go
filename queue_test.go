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
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2015-08-01/servicebus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Azure/azure-service-bus-go/atom"
	"github.com/Azure/azure-service-bus-go/internal/test"
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

func (suite *serviceBusSuite) TestQueueManagementPopulatedQueue() {
	tests := map[string]func(context.Context, *testing.T, *QueueManager, string, *Queue){
		"TestCountDetails": testCountDetails,
	}

	ns := suite.getNewSasInstance()
	qm := ns.NewQueueManager()

	for name, testFunc := range tests {
		suite.T().Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()

			// setup queue for population
			queueName := suite.randEntityName()
			_, err := qm.Put(ctx, queueName)
			if err != nil {
				suite.T().Fatal(err)
			}

			q, err := ns.NewQueue(queueName)
			if err != nil {
				suite.T().Fatal(err)
			}

			testFunc(ctx, t, qm, queueName, q)
			defer suite.cleanupQueue(queueName)
		})
	}
}

func getCountDetailsResults(ctx context.Context, qm *QueueManager, t *testing.T, queueName string, wg *sync.WaitGroup) {
	for {
		qq, err := qm.Get(ctx, queueName)
		assert.NoError(t, err)
		if *qq.CountDetails.ActiveMessageCount == 0 {
			// sleep...   evil....
			time.Sleep(1 * time.Second)
		} else {
			// if not 0.... then assert on what we got. Should be 1
			assert.Equal(t, *qq.CountDetails.ActiveMessageCount, int32(1))
			wg.Done()
			break
		}
	}
}

func testCountDetails(ctx context.Context, t *testing.T, qm *QueueManager, queueName string, q *Queue) {
	if assert.NoError(t, q.Send(ctx, NewMessageFromString("Hello World!"))) {
		var wg sync.WaitGroup
		wg.Add(1)
		go getCountDetailsResults(ctx, qm, t, queueName, &wg)
		end, _ := ctx.Deadline()
		waitUntil(t, &wg, time.Until(end))
	}
}

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

func (suite *serviceBusSuite) TestQueueManager_NotFound() {
	ns := suite.getNewSasInstance()
	qm := ns.NewQueueManager()
	entity, err := qm.Get(context.Background(), "somethingNotThere")
	suite.Nil(entity)
	suite.Require().NotNil(err)
	suite.True(IsErrNotFound(err))
	suite.Equal("entity at /somethingNotThere not found", err.Error())
}

func (suite *serviceBusSuite) TestQueueManagement_Writes() {
	tests := map[string]func(context.Context, *testing.T, *QueueManager, string){
		"TestPutDefaultQueue": testPutQueue,
	}

	ns := suite.getNewSasInstance()
	qm := ns.NewQueueManager()
	for name, testFunc := range tests {
		suite.T().Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
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

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	names := []string{suite.randEntityName(), suite.randEntityName()}
	for _, name := range names {
		if _, err := qm.Put(ctx, name); err != nil {
			suite.T().Fatal(err)
		}
	}

	for name, testFunc := range tests {
		suite.T().Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
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

func (suite *serviceBusSuite) TestQueueManager_QueueWithForwarding() {
	tests := map[string]func(context.Context, *testing.T, *QueueManager, string){
		"TestQueueWithAutoForward":         testQueueWithAutoForward,
		"TestQueueWithForwardDeadLetterTo": testQueueWithForwardDeadLetterTo,
	}

	suite.testQueueMgmt(tests)
}

func testQueueWithForwardDeadLetterTo(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	targetQueueName := "target-" + name
	target := buildQueue(ctx, t, qm, targetQueueName)
	defer func() {
		assert.NoError(t, qm.Delete(ctx, targetQueueName))
	}()
	src := buildQueue(ctx, t, qm, name, QueueEntityWithForwardDeadLetteredMessagesTo(target))
	assert.Equal(t, target.TargetURI(), *src.ForwardDeadLetteredMessagesTo)
}

func testQueueWithAutoForward(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	targetQueueName := "target-" + name
	target := buildQueue(ctx, t, qm, targetQueueName)
	defer func() {
		assert.NoError(t, qm.Delete(ctx, targetQueueName))
	}()
	src := buildQueue(ctx, t, qm, name, QueueEntityWithAutoForward(target))
	assert.Equal(t, target.TargetURI(), *src.ForwardTo)
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
		"TestQueueWithAutoForward":                      testQueueWithAutoForward,
	}

	suite.testQueueMgmt(tests)
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
	assert.Equal(t, int32(1024), *q.MaxSizeInMegabytes, "max size in MBs")
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
	size := 3 * 1024
	q := buildQueue(ctx, t, qm, name, QueueEntityWithMaxSizeInMegabytes(size))
	assert.Equal(t, int32(size), *q.MaxSizeInMegabytes)
}

func testQueueWithDuplicateDetection(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	window := time.Duration(20 * time.Minute)
	q := buildQueue(ctx, t, qm, name, QueueEntityWithDuplicateDetection(&window))
	assert.Equal(t, "PT20M", *q.DuplicateDetectionHistoryTimeWindow)
}

func testQueueWithMessageTimeToLive(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	window := time.Duration(10 * 24 * time.Hour)
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
	require.NoError(t, err)

	q, err := qm.Get(ctx, name)
	require.NoError(t, err)

	return q
}

func (suite *serviceBusSuite) testQueueMgmt(tests map[string]func(context.Context, *testing.T, *QueueManager, string)) {
	ns := suite.getNewSasInstance()
	qm := ns.NewQueueManager()
	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			name := suite.randEntityName()
			testFunc(ctx, t, qm, name)
			defer suite.cleanupQueue(name)

		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func (suite *serviceBusSuite) TestQueueClient() {
	tests := map[string]func(context.Context, *testing.T, *Queue){
		"SimpleSend_NoZeroCheck": testQueueSend,
		"DuplicateDetection":     testDuplicateDetection,
		"MessageProperties":      testMessageProperties,
		"Retry":                  testRequeueOnFail,
		"Defer":                  testDeferMessage,
	}

	window := time.Duration(30 * time.Second)
	mgmtOpts := []QueueManagementOption{
		QueueEntityWithPartitioning(),
		QueueEntityWithDuplicateDetection(&window),
	}
	suite.queueMessageTestWithMgmtOptions(tests, mgmtOpts...)
}

func testRequeueOnFail(ctx context.Context, t *testing.T, q *Queue) {
	const payload = "Hello World!!!"

	if assert.NoError(t, q.Send(ctx, NewMessageFromString(payload))) {
		inner, cancel := context.WithCancel(ctx)
		errs := make(chan error, 1)

		go func() {
			fail := true

			errs <- q.Receive(inner,
				HandlerFunc(func(ctx context.Context, msg *Message) error {
					assert.EqualValues(t, payload, string(msg.Data))
					if fail {
						fail = false
						assert.EqualValues(t, 1, msg.DeliveryCount)
						return msg.Abandon(ctx)
					}
					assert.EqualValues(t, 2, msg.DeliveryCount)
					cancel()
					return msg.Complete(ctx)
				}))
		}()

		select {
		case <-ctx.Done():
			t.Error(ctx.Err())
			return
		case err := <-errs:
			assert.EqualError(t, err, context.Canceled.Error())
		}
	}
}

func testMessageProperties(ctx context.Context, t *testing.T, q *Queue) {
	if assert.NoError(t, q.Send(ctx, NewMessageFromString("Hello World!"))) {
		err := q.ReceiveOne(context.Background(),
			HandlerFunc(func(ctx context.Context, msg *Message) error {
				sp := msg.SystemProperties
				assert.NotNil(t, sp.LockedUntil, "LockedUntil")
				assert.NotNil(t, sp.EnqueuedSequenceNumber, "EnqueuedSequenceNumber")
				assert.NotNil(t, sp.EnqueuedTime, "EnqueuedTime")
				assert.NotNil(t, sp.SequenceNumber, "SequenceNumber")
				assert.NotNil(t, sp.PartitionID, "PartitionID")
				assert.NotNil(t, sp.PartitionKey, "PartitionKey")
				return msg.Complete(ctx)
			}))

		assert.NoError(t, err)
	}
}

func testQueueSend(ctx context.Context, t *testing.T, queue *Queue) {
	rmsg := test.RandomString("foo", 10)
	assert.NoError(t, queue.Send(ctx, NewMessageFromString(fmt.Sprintf("hello %s!", rmsg))))
}

func testDeferMessage(ctx context.Context, t *testing.T, queue *Queue) {
	rmsg := test.RandomString("foo", 10)
	require.NoError(t, queue.Send(ctx, NewMessageFromString(fmt.Sprintf("hello %s!", rmsg))))

	var sequenceNumber *int64
	err := queue.ReceiveOne(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
		sequenceNumber = msg.SystemProperties.SequenceNumber
		return msg.Defer(ctx)
	}))
	require.NoError(t, err)
	require.NotNil(t, sequenceNumber)

	handled := false
	err = queue.ReceiveDeferred(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
		handled = true
		return msg.Complete(ctx)
	}), *sequenceNumber)

	assert.True(t, handled, "expected message handler to be called")
	assert.NoError(t, err)
}

func (suite *serviceBusSuite) TestQueueWithoutDuplicateDetection() {
	tests := map[string]func(context.Context, *testing.T, *Queue){
		"SendBatch_NoZeroCheck": testSendBatch,
	}

	mgmtOpts := []QueueManagementOption{
		QueueEntityWithPartitioning(),
	}
	suite.queueMessageTestWithMgmtOptions(tests, mgmtOpts...)
}

func testSendBatch(ctx context.Context, t *testing.T, q *Queue) {
	msgCount := 10000
	messages := make([]*Message, msgCount)
	for i := 0; i < msgCount; i++ {
		rmsg := test.RandomString("foo", 10)
		messages[i] = NewMessageFromString(fmt.Sprintf("hello %s!", rmsg))
	}

	iterator := NewMessageBatchIterator(StandardMaxMessageSizeInBytes, messages...)
	assert.NoError(t, q.SendBatch(ctx, iterator))
	assertMessageCount(ctx, t, q.namespace, q.Name, int64(msgCount))
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

	received := make(map[interface{}]string)
	inner, cancel := context.WithCancel(ctx)

	var all []*Message
	err := queue.Receive(inner, HandlerFunc(func(ctx context.Context, message *Message) error {
		all = append(all, message)
		if _, ok := received[message.ID]; !ok {
			// caught a new one
			received[message.ID] = string(message.Data)
		} else {
			// caught a dup
			assert.Fail(t, "received a duplicate message")
			for _, item := range all {
				t.Logf("mID: %q, gID: %q, gSeq: %q, lockT: %q", item.ID, *item.SessionID, *item.GroupSequence, *item.LockToken)
			}
		}
		if len(all) == len(messages) {
			cancel()
		}
		return message.Complete(ctx)
	}))
	assert.EqualError(t, err, context.Canceled.Error())
}

func (suite *serviceBusSuite) TestQueueWithReceiveAndDelete() {
	tests := map[string]func(context.Context, *testing.T, *Queue){
		"SimpleSendAndReceive": testQueueSendAndReceiveWithReceiveAndDelete,
	}
	suite.queueMessageTestWithQueueOptions(tests, QueueWithReceiveAndDelete())
}

func testQueueSendAndReceiveWithReceiveAndDelete(ctx context.Context, t *testing.T, queue *Queue) {
	ttl := 5 * time.Minute
	numMessages := rand.Intn(100) + 20
	expected := make(map[string]int, numMessages)
	seen := make(map[string]int, numMessages)
	errs := make(chan error, 1)

	t.Logf("Sending/receiving %d messages", numMessages)

	go func() {
		inner, cancel := context.WithCancel(ctx)
		numSeen := 0
		errs <- queue.Receive(inner, HandlerFunc(func(ctx context.Context, msg *Message) error {
			numSeen++
			seen[string(msg.Data)]++
			if numSeen >= numMessages {
				cancel()
			}
			return nil
		}))
	}()

	for i := 0; i < numMessages; i++ {
		payload := test.RandomString("hello", 10)
		expected[payload]++
		msg := NewMessageFromString(payload)
		msg.TTL = &ttl
		assert.NoError(t, queue.Send(ctx, msg))
	}

	assert.EqualError(t, <-errs, context.Canceled.Error())

	assert.Equal(t, len(expected), len(seen))
	for k, v := range seen {
		assert.Equal(t, expected[k], v)
	}
}

func (suite *serviceBusSuite) TestQueueWithPrefetch() {
	tests := map[string]func(context.Context, *testing.T, *Queue){
		"SendAndReceive": testQueueSendAndReceive,
	}
	suite.queueMessageTestWithQueueOptions(tests, QueueWithPrefetchCount(10))
}

func testQueueSendAndReceive(ctx context.Context, t *testing.T, q *Queue) {
	messages := []string{"foo", "bar", "bazz", "buzz"}
	for _, msg := range messages {
		require.NoError(t, q.Send(ctx, NewMessageFromString(msg)))
	}

	count := 0
	for idx := range messages {
		err := q.ReceiveOne(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
			assert.Equal(t, messages[idx], string(msg.Data))
			count++
			return msg.Complete(ctx)
		}))
		assert.NoError(t, err)
	}
	assert.Len(t, messages, count)
}

func (suite *serviceBusSuite) TestIssue73QueueClient() {
	tests := map[string]func(context.Context, *testing.T, *Queue){
		"SimpleSend200_NoZeroCheck": func(ctx context.Context, t *testing.T, queue *Queue) {
			const count = 50
			for i := 0; i < count; i++ {
				testQueueSend(ctx, t, queue)
			}
			assertMessageCount(ctx, t, queue.namespace, queue.Name, count)
		},
	}
	window := time.Duration(20 * time.Second)
	ttl := time.Duration(14 * 24 * time.Hour)
	suite.queueMessageTestWithMgmtOptions(
		tests,
		QueueEntityWithMessageTimeToLive(&ttl),
		QueueEntityWithDuplicateDetection(&window),
		QueueEntityWithMaxDeliveryCount(10),
		QueueEntityWithMaxSizeInMegabytes(1024),
	)
}

func (suite *serviceBusSuite) TestQueue_NewSession() {
	ns := suite.getNewSasInstance()
	q, err := ns.NewQueue("foo")
	suite.NoError(err)
	sessionID := "123"
	qs := NewQueueSession(q, &sessionID)
	suite.Equal(sessionID, *qs.SessionID())
}

func (suite *serviceBusSuite) TestQueue_NewDeadLetter() {
	tests := map[string]func(context.Context, *testing.T, *Queue){
		"ReceiveOneFromDeadLetter": testReceiveOneFromDeadLetter,
	}
	suite.queueMessageTestWithMgmtOptions(tests, QueueEntityWithMaxDeliveryCount(10))
}

func testReceiveOneFromDeadLetter(ctx context.Context, t *testing.T, q *Queue) {
	qdl := q.NewDeadLetter()
	require.NotNil(t, qdl)
	err := q.Send(ctx, NewMessageFromString("foo"))
	require.NoError(t, err)

	// abandon multiple 10 times until the message is dead lettered
	for i := 0; i < 10; i++ {
		err := q.ReceiveOne(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
			return msg.Abandon(ctx)
		}))
		require.NoError(t, err)
	}

	dl := q.NewDeadLetter()
	called := false
	err = dl.ReceiveOne(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
		assert.Equal(t, "foo", string(msg.Data))
		called = true
		return msg.Complete(ctx)
	}))
	assert.True(t, called)
	assert.NoError(t, err)
}

func (suite *serviceBusSuite) queueMessageTestWithQueueOptions(
	tests map[string]func(context.Context, *testing.T, *Queue),
	queueOpts ...QueueOption) {
	suite.queueMessageTest(tests, queueOpts, []QueueManagementOption{})
}

func (suite *serviceBusSuite) queueMessageTestWithMgmtOptions(
	tests map[string]func(context.Context, *testing.T, *Queue),
	mgmtOptions ...QueueManagementOption) {
	suite.queueMessageTest(tests, []QueueOption{}, mgmtOptions)
}

func (suite *serviceBusSuite) queueMessageTest(
	tests map[string]func(context.Context, *testing.T, *Queue),
	queueOpts []QueueOption,
	mgmtOptions []QueueManagementOption) {

	ns := suite.getNewSasInstance()
	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			queueName := suite.randEntityName()
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			cleanup := makeQueue(ctx, t, ns, queueName, mgmtOptions...)
			q, err := ns.NewQueue(queueName, queueOpts...)
			suite.NoError(err)
			defer func() {
				cleanup()
			}()
			testFunc(ctx, t, q)
			suite.NoError(q.Close(ctx))
			if !t.Failed() && !strings.HasSuffix(name, "NoZeroCheck") {
				assertZeroQueueMessages(ctx, t, ns, queueName)
			}
		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func makeQueue(ctx context.Context, t *testing.T, ns *Namespace, name string, opts ...QueueManagementOption) func() {
	qm := ns.NewQueueManager()
	entity, err := qm.Get(ctx, name)
	if err != nil && !IsErrNotFound(err) {
		assert.FailNow(t, "could not GET a queue entity")
	}

	if entity == nil {
		entity, err = qm.Put(ctx, name, opts...)
		if !assert.NoError(t, err) {
			assert.FailNow(t, "could not PUT a queue entity")
		}
	}
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		_ = qm.Delete(ctx, entity.Name)
	}
}

func assertZeroQueueMessages(ctx context.Context, t *testing.T, ns *Namespace, name string) {
	assertMessageCount(ctx, t, ns, name, 0)
}

func assertMessageCount(ctx context.Context, t *testing.T, ns *Namespace, name string, count int64) {
	qm := ns.NewQueueManager()
	maxTries := 10
	for i := 0; i < maxTries; i++ {
		q, err := qm.Get(ctx, name)
		if !assert.NoError(t, err) {
			return
		}
		if *q.MessageCount == count {
			return
		}
		t.Logf("try %d out of %d, message count was %d, not %d", i+1, maxTries, *q.MessageCount, count)
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
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	ns := suite.getNewSasInstance()
	qm := ns.NewQueueManager()
	err := qm.Delete(ctx, name)
	if err != nil {
		suite.T().Fatal(err)
	}
}
