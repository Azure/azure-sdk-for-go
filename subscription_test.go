package servicebus

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2015-08-01/servicebus"
	"github.com/stretchr/testify/assert"
)

const (
	subscriptionDescription = `
	<SubscriptionDescription
      xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect"
      xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
      <LockDuration>PT1M</LockDuration>
      <RequiresSession>false</RequiresSession>
      <DefaultMessageTimeToLive>P10675199DT2H48M5.4775807S</DefaultMessageTimeToLive>
      <DeadLetteringOnMessageExpiration>false</DeadLetteringOnMessageExpiration>
      <DeadLetteringOnFilterEvaluationExceptions>true</DeadLetteringOnFilterEvaluationExceptions>
      <MessageCount>0</MessageCount>
      <MaxDeliveryCount>10</MaxDeliveryCount>
      <EnableBatchedOperations>true</EnableBatchedOperations>
      <Status>Active</Status>
      <CreatedAt>2018-05-04T22:41:54.183101Z</CreatedAt>
      <UpdatedAt>2018-05-04T22:41:54.183101Z</UpdatedAt>
      <AccessedAt>0001-01-01T00:00:00</AccessedAt>
      <AutoDeleteOnIdle>P10675199DT2H48M5.4775807S</AutoDeleteOnIdle>
      <EntityAvailabilityStatus>Available</EntityAvailabilityStatus>
  </SubscriptionDescription>`

	subscriptionEntry = `
	<entry xmlns="http://www.w3.org/2005/Atom">
		<id>https://sbdjtest.servicebus.windows.net/gosbh6of3g-tagz3cfzrp93m/subscriptions/gosbwg424p-tagz3cfzrp93m?api-version=2017-04</id>
		<title type="text">gosbwg424p-tagz3cfzrp93m</title>
		<published>2018-05-02T20:54:59Z</published>
		<updated>2018-05-02T20:54:59Z</updated>
		<link rel="self" href="https://sbdjtest.servicebus.windows.net/gosbh6of3g-tagz3cfzrp93m/subscriptions/gosbwg424p-tagz3cfzrp93m?api-version=2017-04"/>
		<content type="application/xml">` + subscriptionDescription +
		`</content>
	</entry>`
)

func (suite *serviceBusSuite) TestSubscriptionEntryUnmarshal() {
	var entry SubscriptionEntry
	err := xml.Unmarshal([]byte(subscriptionEntry), &entry)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "https://sbdjtest.servicebus.windows.net/gosbh6of3g-tagz3cfzrp93m/subscriptions/gosbwg424p-tagz3cfzrp93m?api-version=2017-04", entry.ID)
	assert.Equal(suite.T(), "gosbwg424p-tagz3cfzrp93m", entry.Title)
	assert.Equal(suite.T(), "https://sbdjtest.servicebus.windows.net/gosbh6of3g-tagz3cfzrp93m/subscriptions/gosbwg424p-tagz3cfzrp93m?api-version=2017-04", entry.Link.HREF)
	assert.Equal(suite.T(), "PT1M", *entry.Content.SubscriptionDescription.LockDuration)
	assert.NotNil(suite.T(), entry.Content)
}

func (suite *serviceBusSuite) TestSubscriptionUnmarshal() {
	var entry SubscriptionEntry
	err := xml.Unmarshal([]byte(subscriptionEntry), &entry)
	assert.Nil(suite.T(), err)
	t := suite.T()
	s := entry.Content.SubscriptionDescription
	assert.Equal(t, "PT1M", *s.LockDuration)
	assert.Equal(t, false, *s.RequiresSession)
	assert.Equal(t, "P10675199DT2H48M5.4775807S", *s.DefaultMessageTimeToLive)
	assert.Equal(t, false, *s.DeadLetteringOnMessageExpiration)
	assert.Equal(t, int32(10), *s.MaxDeliveryCount)
	assert.Equal(t, true, *s.EnableBatchedOperations)
	assert.Equal(t, int64(0), *s.MessageCount)
	assert.EqualValues(t, servicebus.EntityStatusActive, *s.Status)
}

func (suite *serviceBusSuite) TestSubscriptionManagementWrites() {
	tests := map[string]func(context.Context, *testing.T, *SubscriptionManager, string){
		"TestPutDefaultSubscription": testPutSubscription,
	}

	ns := suite.getNewSasInstance()
	tm := ns.NewTopicManager()

	outerCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	topicName := suite.RandomName("gosb", 6)
	_, err := tm.Put(outerCtx, topicName)
	if err != nil {
		suite.T().Fatal(err)
	}
	topic := ns.NewTopic(topicName)
	sm := topic.NewSubscriptionManager()
	for name, testFunc := range tests {
		suite.T().Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			name := suite.RandomName("gosb", 6)
			testFunc(ctx, t, sm, name)
			defer suite.cleanupSubscription(topicName, name)
		})
	}
}

func testPutSubscription(ctx context.Context, t *testing.T, sm *SubscriptionManager, name string) {
	topic, err := sm.Put(ctx, name)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	if assert.NotNil(t, topic) {
		assert.Equal(t, name, topic.Title)
	}
}

func (suite *serviceBusSuite) cleanupSubscription(topicName, subscriptionName string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ns := suite.getNewSasInstance()
	topic := ns.NewTopic(topicName)
	sm := topic.NewSubscriptionManager()
	err := sm.Delete(ctx, subscriptionName)
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *serviceBusSuite) TestSubscriptionManagement() {
	tests := map[string]func(context.Context, *testing.T, *SubscriptionManager, string, string){
		"TestSubscriptionDefaultSettings":                      testDefaultSubscription,
		"TestSubscriptionWithAutoDeleteOnIdle":                 testSubscriptionWithAutoDeleteOnIdle,
		"TestSubscriptionWithRequiredSessions":                 testSubscriptionWithRequiredSessions,
		"TestSubscriptionWithDeadLetteringOnMessageExpiration": testSubscriptionWithDeadLetteringOnMessageExpiration,
		"TestSubscriptionWithMessageTimeToLive":                testSubscriptionWithMessageTimeToLive,
		"TestSubscriptionWithLockDuration":                     testSubscriptionWithLockDuration,
		"TestSubscriptionWithBatchedOperations":                testSubscriptionWithBatchedOperations,
	}

	ns := suite.getNewSasInstance()
	tm := ns.NewTopicManager()

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			topicName := suite.randEntityName()
			subName := suite.randEntityName()

			_, err := tm.Put(ctx, topicName)
			if err != nil {
				t.Fatal(err)
			}
			topic := ns.NewTopic(topicName)
			sm := topic.NewSubscriptionManager()
			_, err = sm.Put(ctx, topicName)
			if err != nil {
				t.Fatal(err)
			}

			defer func(tName, sName string) {
				innerCtx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				err = sm.Delete(innerCtx, sName)
				err2 := tm.Delete(innerCtx, tName)
				if err != nil {
					suite.T().Fatal(err)
				}
				if err2 != nil {
					suite.T().Fatal(err2)
				}
			}(topicName, subName)

			testFunc(ctx, t, sm, topicName, subName)
		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func testDefaultSubscription(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, name string) {
	s := buildSubscription(ctx, t, sm, name)
	assert.False(t, *s.DeadLetteringOnMessageExpiration, "should not have dead lettering on expiration")
	assert.False(t, *s.RequiresSession, "should not require session")
	assert.Equal(t, "P10675199DT2H48M5.4775807S", *s.AutoDeleteOnIdle, "auto delete is not 10 minutes")
	assert.Equal(t, "P10675199DT2H48M5.4775807S", *s.DefaultMessageTimeToLive, "default TTL")
	assert.EqualValues(t, servicebus.EntityStatusActive, *s.Status)
	assert.Equal(t, "PT1M", *s.LockDuration, "lock duration")
}

func testSubscriptionWithBatchedOperations(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, name string) {
	s := buildSubscription(ctx, t, sm, name)
	assert.True(t, *s.EnableBatchedOperations)
}

func testSubscriptionWithAutoDeleteOnIdle(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, name string) {
	window := time.Duration(20 * time.Minute)
	s := buildSubscription(ctx, t, sm, name, SubscriptionWithAutoDeleteOnIdle(&window))
	assert.Equal(t, "PT20M", *s.AutoDeleteOnIdle)
}

func testSubscriptionWithRequiredSessions(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, name string) {
	s := buildSubscription(ctx, t, sm, name, SubscriptionWithRequiredSessions())
	assert.True(t, *s.RequiresSession)
}

func testSubscriptionWithDeadLetteringOnMessageExpiration(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, name string) {
	s := buildSubscription(ctx, t, sm, name, SubscriptionWithDeadLetteringOnMessageExpiration())
	assert.True(t, *s.DeadLetteringOnMessageExpiration)
}

func testSubscriptionWithMessageTimeToLive(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, name string) {
	window := time.Duration(10 * 24 * 60 * time.Minute)
	s := buildSubscription(ctx, t, sm, name, SubscriptionWithMessageTimeToLive(&window))
	assert.Equal(t, "P10D", *s.DefaultMessageTimeToLive)
}

func testSubscriptionWithLockDuration(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, name string) {
	window := time.Duration(3 * time.Minute)
	s := buildSubscription(ctx, t, sm, name, SubscriptionWithLockDuration(&window))
	assert.Equal(t, "PT3M", *s.LockDuration)
}

func buildSubscription(ctx context.Context, t *testing.T, sm *SubscriptionManager, name string, opts ...SubscriptionOption) *SubscriptionDescription {
	s, err := sm.Put(ctx, name, opts...)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("%v", err))
	}
	return &s.Content.SubscriptionDescription
}

func (suite *serviceBusSuite) TestSubscription() {
	tests := map[string]func(context.Context, *testing.T, *Topic, *Subscription){
		"SimpleReceive": testSubscriptionReceive,
	}

	ns := suite.getNewSasInstance()
	tm := ns.NewTopicManager()
	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			topicName := suite.randEntityName()
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			_, err := tm.Put(ctx, topicName)
			if err != nil {
				log.Fatalln(err)
			}

			topic := ns.NewTopic(topicName)
			sm := topic.NewSubscriptionManager()
			subName := suite.randEntityName()
			sm.Put(ctx, subName)
			subscription := topic.NewSubscription(subName)
			defer func() {
				closeCtx, closeCancel := context.WithTimeout(context.Background(), timeout)
				defer closeCancel()
				topic.Close(closeCtx)
				suite.cleanupTopic(topicName)
				subscription.Close(closeCtx)
				suite.cleanupSubscription(topicName, subName)

			}()
			testFunc(ctx, t, topic, subscription)
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func testSubscriptionReceive(ctx context.Context, t *testing.T, topic *Topic, sub *Subscription) {
	err := topic.Send(ctx, NewEventFromString("hello!"))
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(1)
	_, err = sub.Receive(ctx, func(eventCtx context.Context, evt *Event) error {
		wg.Done()
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	end, _ := ctx.Deadline()
	waitUntil(t, &wg, time.Until(end))
}
