package servicebus

import (
	"context"
	"fmt"
	"testing"
	"time"

	mgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func (suite *ServiceBusSuite) TestQueueManagement() {
	tests := map[string]func(*testing.T, SenderReceiverManager, string){
		"TestQueueDefaultSettings":                      testDefaultQueue,
		"TestQueueWithAutoDeleteOnIdle":                 testQueueWithAutoDeleteOnIdle,
		"TestQueueWithRequiredSessions":                 testQueueWithRequiredSessions,
		"TestQueueWithDeadLetteringOnMessageExpiration": testQueueWithDeadLetteringOnMessageExpiration,
		"TestQueueWithPartitioning":                     testQueueWithPartitioning,
		"TestQueueWithMaxSizeInMegabytes":               testQueueWithMaxSizeInMegabytes,
		"TestQueueWithDuplicateDetection":               testQueueWithDuplicateDetection,
		"TestQueueWithMessageTimeToLive":                testQueueWithMessageTimeToLive,
		"TestQueueWithLockDuration":                     testQueueWithLockDuration,
	}

	sb := suite.getNewInstance()
	defer func() {
		sb.Close()
	}()

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			entityName := randomName("gosbtest", 10)
			defer func(name string) {
				err := sb.DeleteQueue(context.Background(), name)
				if err != nil {
					log.Fatalln(err)
				}
			}(entityName)
			testFunc(t, sb, entityName)

		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func testDefaultQueue(t *testing.T, sb SenderReceiverManager, name string) {
	q := buildQueue(t, sb, name)
	assert.False(t, *q.EnableExpress, "should not have Express enabled")
	assert.False(t, *q.EnablePartitioning, "should not have partitioning enabled")
	assert.False(t, *q.DeadLetteringOnMessageExpiration, "should not have dead lettering on expiration")
	assert.False(t, *q.RequiresDuplicateDetection, "should not require dup detection")
	assert.False(t, *q.RequiresSession, "should not require session")
	assert.Equal(t, "P10675199DT2H48M5.4775807S", *q.AutoDeleteOnIdle, "auto delete is not 10 minutes")
	assert.Equal(t, "PT10M", *q.DuplicateDetectionHistoryTimeWindow, "dup detection is not 10 minutes")
	assert.Equal(t, int32(5*Megabytes), *q.MaxSizeInMegabytes, "max size in MBs")
	assert.Equal(t, "P10675199DT2H48M5.4775807S", *q.DefaultMessageTimeToLive, "default TTL")
	assert.Equal(t, mgmt.Active, q.Status, "queue status")
	assert.Equal(t, "PT1M", *q.LockDuration, "lock duration")
}

func testQueueWithAutoDeleteOnIdle(t *testing.T, sb SenderReceiverManager, name string) {
	window := time.Duration(20 * time.Minute)
	q := buildQueue(t, sb, name, QueueWithAutoDeleteOnIdle(&window))
	assert.Equal(t, "PT20M", *q.AutoDeleteOnIdle)
}

func testQueueWithRequiredSessions(t *testing.T, sb SenderReceiverManager, name string) {
	q := buildQueue(t, sb, name, QueueWithRequiredSessions())
	assert.True(t, *q.RequiresSession)
}

func testQueueWithDeadLetteringOnMessageExpiration(t *testing.T, sb SenderReceiverManager, name string) {
	q := buildQueue(t, sb, name, QueueWithDeadLetteringOnMessageExpiration())
	assert.True(t, *q.DeadLetteringOnMessageExpiration)
}

func testQueueWithPartitioning(t *testing.T, sb SenderReceiverManager, name string) {
	q := buildQueue(t, sb, name, QueueWithPartitioning())
	assert.True(t, *q.EnablePartitioning)
}

func testQueueWithMaxSizeInMegabytes(t *testing.T, sb SenderReceiverManager, name string) {
	size := 3 * Megabytes
	q := buildQueue(t, sb, name, QueueWithMaxSizeInMegabytes(size))
	assert.Equal(t, int32(size), *q.MaxSizeInMegabytes)
}

func testQueueWithDuplicateDetection(t *testing.T, sb SenderReceiverManager, name string) {
	window := time.Duration(20 * time.Minute)
	q := buildQueue(t, sb, name, QueueWithDuplicateDetection(&window))
	assert.Equal(t, "PT20M", *q.DuplicateDetectionHistoryTimeWindow)
}

func testQueueWithMessageTimeToLive(t *testing.T, sb SenderReceiverManager, name string) {
	window := time.Duration(10 * 24 * 60 * time.Minute)
	q := buildQueue(t, sb, name, QueueWithMessageTimeToLive(&window))
	assert.Equal(t, "P10D", *q.DefaultMessageTimeToLive)
}

func testQueueWithLockDuration(t *testing.T, sb SenderReceiverManager, name string) {
	window := time.Duration(3 * time.Minute)
	q := buildQueue(t, sb, name, QueueWithLockDuration(&window))
	assert.Equal(t, "PT3M", *q.LockDuration)
}

func buildQueue(t *testing.T, sb SenderReceiverManager, name string, opts ...QueueOption) *mgmt.SBQueue {
	q, err := sb.EnsureQueue(context.Background(), name, opts...)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("%v", err))
	}
	return q
}
