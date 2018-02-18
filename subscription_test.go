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

func (suite *ServiceBusSuite) TestSubscriptionManagement() {
	tests := map[string]func(*testing.T, SenderReceiverManager, string, string){
		"TestSubscriptionDefaultSettings":                      testDefaultSubscription,
		"TestSubscriptionWithAutoDeleteOnIdle":                 testSubscriptionWithAutoDeleteOnIdle,
		"TestSubscriptionWithRequiredSessions":                 testSubscriptionWithRequiredSessions,
		"TestSubscriptionWithDeadLetteringOnMessageExpiration": testSubscriptionWithDeadLetteringOnMessageExpiration,
		"TestSubscriptionWithMessageTimeToLive":                testSubscriptionWithMessageTimeToLive,
		"TestSubscriptionWithLockDuration":                     testSubscriptionWithLockDuration,
		"TestSubscriptionWithBatchedOperations":                testSubscriptionWithBatchedOperations,
	}

	sb := suite.getNewInstance()
	defer func() {
		sb.Close()
	}()

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			ctx := context.Background()
			topicName := randomName("gosbtest", 10)
			subName := randomName("gosbtest", 10)
			_, err := sb.EnsureTopic(ctx, topicName)
			if err != nil {
				log.Fatalln(err)
			}

			defer func(tName, sName string) {
				err = sb.DeleteSubscription(ctx, tName, sName)
				err2 := sb.DeleteTopic(ctx, tName)
				if err != nil {
					log.Fatalln(err)
				}

				if err2 != nil {
					log.Fatalln(err2)
				}
			}(topicName, subName)

			testFunc(t, sb, topicName, subName)
		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func testDefaultSubscription(t *testing.T, sb SenderReceiverManager, topicName, name string) {
	s := buildSubscription(t, sb, topicName, name)
	assert.False(t, *s.DeadLetteringOnMessageExpiration, "should not have dead lettering on expiration")
	assert.False(t, *s.RequiresSession, "should not require session")
	assert.Equal(t, "P10675199DT2H48M5.4775807S", *s.AutoDeleteOnIdle, "auto delete is not 10 minutes")
	assert.Nil(t, s.DuplicateDetectionHistoryTimeWindow, "dup detection is nil")
	assert.Equal(t, "P10675199DT2H48M5.4775807S", *s.DefaultMessageTimeToLive, "default TTL")
	assert.Equal(t, mgmt.Active, s.Status, "subscription status")
	assert.Equal(t, "PT1M", *s.LockDuration, "lock duration")
}

func testSubscriptionWithBatchedOperations(t *testing.T, sb SenderReceiverManager, topicName, name string) {
	s := buildSubscription(t, sb, topicName, name, SubscriptionWithBatchedOperations())
	assert.True(t, *s.EnableBatchedOperations)
}

func testSubscriptionWithAutoDeleteOnIdle(t *testing.T, sb SenderReceiverManager, topicName, name string) {
	window := time.Duration(20 * time.Minute)
	s := buildSubscription(t, sb, topicName, name, SubscriptionWithAutoDeleteOnIdle(&window))
	assert.Equal(t, "PT20M", *s.AutoDeleteOnIdle)
}

func testSubscriptionWithRequiredSessions(t *testing.T, sb SenderReceiverManager, topicName, name string) {
	s := buildSubscription(t, sb, topicName, name, SubscriptionWithRequiredSessions())
	assert.True(t, *s.RequiresSession)
}

func testSubscriptionWithDeadLetteringOnMessageExpiration(t *testing.T, sb SenderReceiverManager, topicName, name string) {
	s := buildSubscription(t, sb, topicName, name, SubscriptionWithDeadLetteringOnMessageExpiration())
	assert.True(t, *s.DeadLetteringOnMessageExpiration)
}

func testSubscriptionWithMessageTimeToLive(t *testing.T, sb SenderReceiverManager, topicName, name string) {
	window := time.Duration(10 * 24 * 60 * time.Minute)
	s := buildSubscription(t, sb, topicName, name, SubscriptionWithMessageTimeToLive(&window))
	assert.Equal(t, "P10D", *s.DefaultMessageTimeToLive)
}

func testSubscriptionWithLockDuration(t *testing.T, sb SenderReceiverManager, topicName, name string) {
	window := time.Duration(3 * time.Minute)
	s := buildSubscription(t, sb, topicName, name, SubscriptionWithLockDuration(&window))
	assert.Equal(t, "PT3M", *s.LockDuration)
}

func buildSubscription(t *testing.T, sb SenderReceiverManager, topicName, name string, opts ...SubscriptionOption) *mgmt.SBSubscription {
	s, err := sb.EnsureSubscription(context.Background(), topicName, name, opts...)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("%v", err))
	}
	return s
}
