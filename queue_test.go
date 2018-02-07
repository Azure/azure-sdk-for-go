package servicebus

import (
	"context"

	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-amqp-common-go/uuid"
	"github.com/Azure/azure-service-bus-go/internal/test"
	"github.com/stretchr/testify/assert"
)

const (
	queueDescription1 = `
		<QueueDescription xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect" xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
			<LockDuration>PT30S</LockDuration>
			<MaxSizeInMegabytes>16384</MaxSizeInMegabytes>
			<RequiresDuplicateDetection>false</RequiresDuplicateDetection>
			<RequiresSession>false</RequiresSession>
			<DefaultMessageTimeToLive>P14D</DefaultMessageTimeToLive>
			<DeadLetteringOnMessageExpiration>false</DeadLetteringOnMessageExpiration>
			<DuplicateDetectionHistoryTimeWindow>PT10M</DuplicateDetectionHistoryTimeWindow>
			<MaxDeliveryCount>10</MaxDeliveryCount>
			<EnableBatchedOperations>true</EnableBatchedOperations>
			<SizeInBytes>0</SizeInBytes>
			<MessageCount>0</MessageCount>
		</QueueDescription>`

	queueDescription2 = `
		<QueueDescription xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect" xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
			<LockDuration>PT15S</LockDuration>
			<MaxSizeInMegabytes>1024</MaxSizeInMegabytes>
			<RequiresDuplicateDetection>true</RequiresDuplicateDetection>
			<RequiresSession>true</RequiresSession>
			<DefaultMessageTimeToLive>P14D</DefaultMessageTimeToLive>
			<DeadLetteringOnMessageExpiration>false</DeadLetteringOnMessageExpiration>
			<DuplicateDetectionHistoryTimeWindow>PT10M</DuplicateDetectionHistoryTimeWindow>
			<MaxDeliveryCount>100</MaxDeliveryCount>
			<EnableBatchedOperations>true</EnableBatchedOperations>
			<SizeInBytes>10</SizeInBytes>
			<MessageCount>10</MessageCount>
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
	timeout = 20 * time.Second
)

func (suite *serviceBusSuite) TestQueueEntryUnmarshal() {
	var entry QueueEntry
	err := xml.Unmarshal([]byte(queueEntry1), &entry)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "https://sbdjtest.servicebus.windows.net/foo", entry.ID)
	assert.Equal(suite.T(), "foo", entry.Title)
	assert.Equal(suite.T(), "sbdjtest", *entry.Author.Name)
	assert.Equal(suite.T(), "https://sbdjtest.servicebus.windows.net/foo", entry.Link.HREF)
	assert.Equal(suite.T(), "PT30S", *entry.Content.QueueDescription.LockDuration)
	assert.NotNil(suite.T(), entry.Content)
}

func (suite *serviceBusSuite) TestQueueUnmarshal() {
	var entry Entry
	err := xml.Unmarshal([]byte(queueEntry1), &entry)
	assert.Nil(suite.T(), err)

	var q QueueDescription
	err = xml.Unmarshal([]byte(entry.Content.Body), &q)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "PT30S", *q.LockDuration)
	assert.Equal(suite.T(), int32(16384), *q.MaxSizeInMegabytes)
	assert.Equal(suite.T(), false, *q.RequiresDuplicateDetection)
	assert.Equal(suite.T(), false, *q.RequiresSession)
	assert.Equal(suite.T(), "P14D", *q.DefaultMessageTimeToLive)
	assert.Equal(suite.T(), false, *q.DeadLetteringOnMessageExpiration)
	assert.Equal(suite.T(), "PT10M", *q.DuplicateDetectionHistoryTimeWindow)
	assert.Equal(suite.T(), int32(10), *q.MaxDeliveryCount)
	assert.Equal(suite.T(), true, *q.EnableBatchedOperations)
	assert.Equal(suite.T(), int64(0), *q.SizeInBytes)
	assert.Equal(suite.T(), int64(0), *q.MessageCount)
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

			err := qm.Delete(ctx, name)
			assert.Nil(t, err)
		})
	}
}

func testPutQueue(ctx context.Context, t *testing.T, qm *QueueManager, name string) {
	q, err := qm.Put(ctx, name)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	if assert.NotNil(t, q) {
		assert.Equal(t, name, q.Title)
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

	names := []string{suite.randQueueName(), suite.randQueueName()}
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
		if err := qm.Delete(ctx, name); err != nil {
			suite.T().Fatal(err)
		}
	}
}

func testGetQueue(ctx context.Context, t *testing.T, qm *QueueManager, names []string) {
	q, err := qm.Get(ctx, names[0])
	assert.Nil(t, err)
	assert.NotNil(t, q)
	assert.Equal(t, q.Entry.Title, names[0])
}

func testListQueues(ctx context.Context, t *testing.T, qm *QueueManager, names []string) {
	q, err := qm.List(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, q)
	queueNames := make([]string, len(q.Entries))
	for idx, entry := range q.Entries {
		queueNames[idx] = entry.Title
	}

	for _, name := range names {
		assert.Contains(t, queueNames, name)
	}
}

func (suite *serviceBusSuite) randQueueName() string {
	return suite.RandomName("goq", 6)
}

func (suite *serviceBusSuite) TestQueueManagement() {
	tests := map[string]func(*testing.T, *QueueManager, string){
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
			name := suite.randQueueName()
			defer func(n string) {
				err := qm.Delete(context.Background(), n)
				if err != nil {
					t.Fatal(err)
				}
			}(name)
			testFunc(t, qm, name)

		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func testDefaultQueue(t *testing.T, qm *QueueManager, name string) {
	q := buildQueue(t, qm, name)
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

func testQueueWithAutoDeleteOnIdle(t *testing.T, qm *QueueManager, name string) {
	window := time.Duration(20 * time.Minute)
	q := buildQueue(t, qm, name, QueueWithAutoDeleteOnIdle(&window))
	assert.Equal(t, "PT20M", *q.AutoDeleteOnIdle)
}

func testQueueWithRequiredSessions(t *testing.T, qm *QueueManager, name string) {
	q := buildQueue(t, qm, name, QueueWithRequiredSessions())
	assert.True(t, *q.RequiresSession)
}

func testQueueWithDeadLetteringOnMessageExpiration(t *testing.T, qm *QueueManager, name string) {
	q := buildQueue(t, qm, name, QueueWithDeadLetteringOnMessageExpiration())
	assert.True(t, *q.DeadLetteringOnMessageExpiration)
}

func testQueueWithPartitioning(t *testing.T, qm *QueueManager, name string) {
	q := buildQueue(t, qm, name, QueueWithPartitioning())
	assert.True(t, *q.EnablePartitioning)
}

func testQueueWithMaxSizeInMegabytes(t *testing.T, qm *QueueManager, name string) {
	size := 3 * Megabytes
	q := buildQueue(t, qm, name, QueueWithMaxSizeInMegabytes(size))
	assert.Equal(t, int32(size), *q.MaxSizeInMegabytes)
}

func testQueueWithDuplicateDetection(t *testing.T, qm *QueueManager, name string) {
	window := time.Duration(20 * time.Minute)
	q := buildQueue(t, qm, name, QueueWithDuplicateDetection(&window))
	assert.Equal(t, "PT20M", *q.DuplicateDetectionHistoryTimeWindow)
}

func testQueueWithMessageTimeToLive(t *testing.T, qm *QueueManager, name string) {
	window := time.Duration(10 * 24 * 60 * time.Minute)
	q := buildQueue(t, qm, name, QueueWithMessageTimeToLive(&window))
	assert.Equal(t, "P10D", *q.DefaultMessageTimeToLive)
}

func testQueueWithLockDuration(t *testing.T, qm *QueueManager, name string) {
	window := time.Duration(3 * time.Minute)
	q := buildQueue(t, qm, name, QueueWithLockDuration(&window))
	assert.Equal(t, "PT3M", *q.LockDuration)
}

func buildQueue(t *testing.T, qm *QueueManager, name string, opts ...QueueOption) QueueDescription {
	q, err := qm.Put(context.Background(), name, opts...)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("%v", err))
	}
	return q.Content.QueueDescription
}

func (suite *serviceBusSuite) TestQueue() {
	tests := map[string]func(context.Context, *testing.T, *Queue){
		"SimpleSend":            testQueueSend,
		"SendAndReceiveInOrder": testQueueSendAndReceiveInOrder,
		"DuplicateDetection":    testDuplicateDetection,
	}

	ns := suite.getNewSasInstance()
	qm := ns.NewQueueManager()
	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			queueName := suite.randQueueName()
			defer qm.Delete(context.Background(), queueName)

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			window := 3 * time.Minute
			_, err := qm.Put(
				ctx,
				queueName,
				QueueWithPartitioning(),
				QueueWithDuplicateDetection(nil),
				QueueWithLockDuration(&window))
			if err != nil {
				log.Fatalln(err)
			}

			q := ns.NewQueue(queueName)
			defer q.Close(ctx)
			testFunc(ctx, t, q)
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func testQueueSend(ctx context.Context, t *testing.T, queue *Queue) {
	err := queue.Send(ctx, NewEventFromString("hello!"))
	assert.Nil(t, err)
}

func testQueueSendAndReceiveInOrder(ctx context.Context, t *testing.T, queue *Queue) {
	numMessages := rand.Intn(100) + 20
	messages := make([]string, numMessages)
	for i := 0; i < numMessages; i++ {
		messages[i] = test.RandomString("hello", 10)
	}

	for _, message := range messages {
		err := queue.Send(ctx, NewEventFromString(message))
		if err != nil {
			t.Fatal(err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(numMessages)
	// ensure in-order processing of messages from the queue
	count := 0
	queue.Receive(ctx, func(ctx context.Context, event *Event) error {
		assert.Equal(t, messages[count], string(event.Data))
		count++
		wg.Done()
		return nil
	})
	end, _ := ctx.Deadline()
	waitUntil(t, &wg, time.Until(end))
}

func testDuplicateDetection(ctx context.Context, t *testing.T, queue *Queue) {
	dupID := mustUUID(t)
	messages := []struct {
		ID   string
		Data string
	}{
		{
			ID:   dupID.String(),
			Data: "hello 1!",
		},
		{
			ID:   dupID.String(),
			Data: "hello 1!",
		},
		{
			ID:   mustUUID(t).String(),
			Data: "hello 2!",
		},
	}

	for _, msg := range messages {
		err := queue.Send(ctx, NewEventFromString(msg.Data), SendWithMessageID(msg.ID))
		if err != nil {
			t.Fatal(err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(2)
	received := make(map[interface{}]string)
	queue.Receive(ctx, func(ctx context.Context, event *Event) error {
		// we should get 2 messages discarding the duplicate ID
		received[event.ID] = string(event.Data)
		wg.Done()
		return nil
	})
	end, _ := ctx.Deadline()
	waitUntil(t, &wg, time.Until(end))
	assert.Equal(t, 2, len(received), "should not have more than 2 messages", received)
}

func (suite *serviceBusSuite) TestQueueWithRequiredSessions() {
	tests := map[string]func(context.Context, *testing.T, *Queue){
		"TestSendAndReceiveSession": testQueueWithRequiredSessionSendAndReceive,
	}

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			ns := suite.getNewSasInstance()
			qm := ns.NewQueueManager()

			queueName := suite.randQueueName()
			window := 3 * time.Minute
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			_, err := qm.Put(
				ctx,
				queueName,
				QueueWithPartitioning(),
				QueueWithDuplicateDetection(nil),
				QueueWithLockDuration(&window),
				QueueWithRequiredSessions())
			if err != nil {
				t.Fatal(err)
			}

			q := ns.NewQueue(queueName)
			defer q.Close(ctx)
			testFunc(ctx, t, q)

			qd, err := qm.Get(ctx, queueName)
			if assert.NoError(t, err) {
				assert.Zero(t, *qd.Content.QueueDescription.MessageCount, "message count for queue should be zero")
			}
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func testQueueWithRequiredSessionSendAndReceive(ctx context.Context, t *testing.T, queue *Queue) {
	sessionID := mustUUID(t).String()
	numMessages := rand.Intn(100) + 20
	messages := make([]string, numMessages)
	for i := 0; i < numMessages; i++ {
		messages[i] = test.RandomString("hello", 10)
	}

	for idx, message := range messages {
		err := queue.Send(ctx, NewEventFromString(message), SendWithSession(sessionID, uint32(idx)))
		if err != nil {
			t.Fatal(err)
		}
	}

	var wg sync.WaitGroup
	wg.Add(numMessages)
	// ensure in-order processing of messages from the queue
	count := 0
	handler := func(ctx context.Context, event *Event) error {
		assert.Equal(t, messages[count], string(event.Data))
		count++
		wg.Done()
		return nil
	}
	queue.Receive(ctx, handler, ReceiverWithSession(sessionID))
	end, _ := ctx.Deadline()
	waitUntil(t, &wg, time.Until(end))
}

func mustUUID(t *testing.T) uuid.UUID {
	id, err := uuid.NewV4()
	if err != nil {
		t.Fatal(err)
	}
	return id
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
		t.Error("took longer than " + fmtDuration(d))
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second) / time.Second
	return fmt.Sprintf("%d seconds", d)
}
