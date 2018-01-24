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
	}

	spToken := suite.servicePrincipalToken()
	sb, err := NewWithSPToken(spToken, suite.SubscriptionID, ResourceGroupName, suite.Namespace, RootRuleName, suite.Environment)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		log.Debug("before close")
		sb.Close()
		log.Debug("after close")
	}()

	for name, testFunc := range tests {
		entityName := randomName("gosbtest", 10)
		suite.T().Run(name, func(t *testing.T) { testFunc(t, sb, entityName) })
		err = sb.DeleteTopic(context.Background(), entityName)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func testDefaultTopic(t *testing.T, sb SenderReceiverManager, name string) {
	topic := buildOrFail(t, sb, name)
	assert.False(t, *topic.EnableExpress, "should not have Express enabled")
	assert.False(t, *topic.EnableBatchedOperations, "should not have batching enabled")
	assert.False(t, *topic.EnablePartitioning, "should not have partitioning enabled")
	assert.False(t, *topic.SupportOrdering, "should not support ordering")
	assert.Equal(t, "P10675199DT2H48M5.4775807S", *topic.AutoDeleteOnIdle, "auto delete is not 10 minutes")
	assert.False(t, *topic.RequiresDuplicateDetection, "should not require dup detection")
	assert.Equal(t, "PT10M", *topic.DuplicateDetectionHistoryTimeWindow, "dup detection is not 10 minutes")
}

func testPartitionedTopic(t *testing.T, sb SenderReceiverManager, name string) {
	topic := buildOrFail(t, sb, name, TopicWithPartitioning())
	assert.True(t, *topic.EnablePartitioning)
}

func testSupportOrdering(t *testing.T, sb SenderReceiverManager, name string) {
	topic := buildOrFail(t, sb, name, TopicWithOrdering())
	assert.True(t, *topic.SupportOrdering)
}

func testTopicWithDuplicateDetection(t *testing.T, sb SenderReceiverManager, name string) {
	window := time.Duration(20 * time.Minute)
	topic := buildOrFail(t, sb, name, TopicWithDuplicateDetection(&window))
	assert.True(t, *topic.RequiresDuplicateDetection)
	assert.Equal(t, "PT20M", *topic.DuplicateDetectionHistoryTimeWindow)
}

func testTopicWithAutoDeleteOnIdle(t *testing.T, sb SenderReceiverManager, name string) {
	window := time.Duration(20 * time.Minute)
	topic := buildOrFail(t, sb, name, TopicWithAutoDeleteOnIdle(&window))
	assert.Equal(t, "PT20M", *topic.AutoDeleteOnIdle)
}

func testTopicWithBatchedOperations(t *testing.T, sb SenderReceiverManager, name string) {
	topic := buildOrFail(t, sb, name, TopicWithBatchedOperations())
	assert.True(t, *topic.EnableBatchedOperations)
}

func testTopicWithExpress(t *testing.T, sb SenderReceiverManager, name string) {
	topic := buildOrFail(t, sb, name, TopicWithExpress())
	assert.True(t, *topic.EnableExpress)
}

func testTopicWithMessageTimeToLive(t *testing.T, sb SenderReceiverManager, name string) {
	window := time.Duration(20 * time.Minute)
	topic := buildOrFail(t, sb, name, TopicWithMessageTimeToLive(&window))
	assert.Equal(t, "PT20M", *topic.DefaultMessageTimeToLive)
}

func buildOrFail(t *testing.T, sb SenderReceiverManager, name string, opts ...TopicOption) *mgmt.SBTopic {
	topic, err := sb.EnsureTopic(context.Background(), name, opts...)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("%v", err))
	}
	return topic
}
