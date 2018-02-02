package servicebus

import (
	"context"
	"fmt"
	mgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func (suite *ServiceBusSuite) TestTopicManagement() {
	tests := map[string]func(*testing.T, SenderReceiverManager, string){
		"DefaultTopicCreation":        testDefaultTopic,
		"TopicWithPartitioning":       testPartitionedTopic,
		"TopicWithOrdering":           testSupportOrdering,
		"TopicWithDuplicateDetection": testTopicWithDuplicateDetection,
		"TopicWithAutoDeleteOnIdle":   testTopicWithAutoDeleteOnIdle,
		"TopicWithTimeToLive":         testTopicWithMessageTimeToLive,
		"TopicWithBatchOperations":    testTopicWithBatchedOperations,
		"TopicWithExpress":            testTopicWithExpress,
		"TopicWithMaxSizeInMegabytes": testTopicWithMaxSizeInMegabytes,
	}

	sb := suite.getNewInstance()
	defer func() {
		sb.Close()
	}()

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			entityName := randomName("gosbtest", 10)
			defer func(name string) {
				err := sb.DeleteTopic(context.Background(), name)
				if err != nil {
					log.Fatalln(err)
				}
			}(entityName)
			testFunc(t, sb, entityName)

		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func testDefaultTopic(t *testing.T, sb SenderReceiverManager, name string) {
	topic := buildTopic(t, sb, name)
	assert.False(t, *topic.EnableExpress, "should not have Express enabled")
	assert.False(t, *topic.EnableBatchedOperations, "should not have batching enabled")
	assert.False(t, *topic.EnablePartitioning, "should not have partitioning enabled")
	assert.False(t, *topic.SupportOrdering, "should not support ordering")
	assert.False(t, *topic.RequiresDuplicateDetection, "should not require dup detection")
	assert.Equal(t, "P10675199DT2H48M5.4775807S", *topic.AutoDeleteOnIdle, "auto delete is not 10 minutes")
	assert.Equal(t, "PT10M", *topic.DuplicateDetectionHistoryTimeWindow, "dup detection is not 10 minutes")
	assert.Equal(t, mgmt.Active, topic.Status, "topic status")
}

func testPartitionedTopic(t *testing.T, sb SenderReceiverManager, name string) {
	topic := buildTopic(t, sb, name, TopicWithPartitioning())
	assert.True(t, *topic.EnablePartitioning)
}

func testSupportOrdering(t *testing.T, sb SenderReceiverManager, name string) {
	topic := buildTopic(t, sb, name, TopicWithOrdering())
	assert.True(t, *topic.SupportOrdering)
}

func testTopicWithDuplicateDetection(t *testing.T, sb SenderReceiverManager, name string) {
	window := time.Duration(20 * time.Minute)
	topic := buildTopic(t, sb, name, TopicWithDuplicateDetection(&window))
	assert.True(t, *topic.RequiresDuplicateDetection)
	assert.Equal(t, "PT20M", *topic.DuplicateDetectionHistoryTimeWindow)
}

func testTopicWithAutoDeleteOnIdle(t *testing.T, sb SenderReceiverManager, name string) {
	window := time.Duration(20 * time.Minute)
	topic := buildTopic(t, sb, name, TopicWithAutoDeleteOnIdle(&window))
	assert.Equal(t, "PT20M", *topic.AutoDeleteOnIdle)
}

func testTopicWithBatchedOperations(t *testing.T, sb SenderReceiverManager, name string) {
	topic := buildTopic(t, sb, name, TopicWithBatchedOperations())
	assert.True(t, *topic.EnableBatchedOperations)
}

func testTopicWithExpress(t *testing.T, sb SenderReceiverManager, name string) {
	topic := buildTopic(t, sb, name, TopicWithExpress())
	assert.True(t, *topic.EnableExpress)
}

func testTopicWithMessageTimeToLive(t *testing.T, sb SenderReceiverManager, name string) {
	window := time.Duration(20 * time.Minute)
	topic := buildTopic(t, sb, name, TopicWithMessageTimeToLive(&window))
	assert.Equal(t, "PT20M", *topic.DefaultMessageTimeToLive)
}

func testTopicWithMaxSizeInMegabytes(t *testing.T, sb SenderReceiverManager, name string) {
	size := 2 * Megabytes
	topic := buildTopic(t, sb, name, TopicWithMaxSizeInMegabytes(size))
	assert.Equal(t, int32(size), *topic.MaxSizeInMegabytes)
}

func buildTopic(t *testing.T, sb SenderReceiverManager, name string, opts ...TopicOption) *mgmt.SBTopic {
	topic, err := sb.EnsureTopic(context.Background(), name, opts...)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("%v", err))
	}
	return topic
}

func (suite *ServiceBusSuite) TestTopicSend() {
	tests := map[string]func(*testing.T, SenderReceiverManager, string){}

	sb := suite.getNewInstance()
	defer func() {
		sb.Close()
	}()

	for name, testFunc := range tests {
		entityName := randomName("gosbtest", 10)
		sb.EnsureTopic(context.Background(), entityName)
		suite.T().Run(name, func(t *testing.T) { testFunc(t, sb, entityName) })
		err := sb.DeleteTopic(context.Background(), entityName)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
