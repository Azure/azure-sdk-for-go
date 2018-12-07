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
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-amqp-common-go/uuid"
	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2015-08-01/servicebus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	subscriptionDescriptionContent = `
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

	subscriptionEntryContent = `
	<entry xmlns="http://www.w3.org/2005/Atom">
		<id>https://sbdjtest.servicebus.windows.net/gosbh6of3g-tagz3cfzrp93m/subscriptions/gosbwg424p-tagz3cfzrp93m?api-version=2017-04</id>
		<title type="text">gosbwg424p-tagz3cfzrp93m</title>
		<published>2018-05-02T20:54:59Z</published>
		<updated>2018-05-02T20:54:59Z</updated>
		<link rel="self" href="https://sbdjtest.servicebus.windows.net/gosbh6of3g-tagz3cfzrp93m/subscriptions/gosbwg424p-tagz3cfzrp93m?api-version=2017-04"/>
		<content type="application/xml">` + subscriptionDescriptionContent +
		`</content>
	</entry>`
)

func (suite *serviceBusSuite) TestSubscriptionEntryUnmarshal() {
	var entry subscriptionEntry
	err := xml.Unmarshal([]byte(subscriptionEntryContent), &entry)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "https://sbdjtest.servicebus.windows.net/gosbh6of3g-tagz3cfzrp93m/subscriptions/gosbwg424p-tagz3cfzrp93m?api-version=2017-04", entry.ID)
	assert.Equal(suite.T(), "gosbwg424p-tagz3cfzrp93m", entry.Title)
	assert.Equal(suite.T(), "https://sbdjtest.servicebus.windows.net/gosbh6of3g-tagz3cfzrp93m/subscriptions/gosbwg424p-tagz3cfzrp93m?api-version=2017-04", entry.Link.HREF)
	assert.Equal(suite.T(), "PT1M", *entry.Content.SubscriptionDescription.LockDuration)
	assert.NotNil(suite.T(), entry.Content)
}

func (suite *serviceBusSuite) TestSubscriptionUnmarshal() {
	var entry subscriptionEntry
	err := xml.Unmarshal([]byte(subscriptionEntryContent), &entry)
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
	outerCtx, outerCancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer outerCancel()
	topicName := suite.RandomName("gosb", 6)
	cleanupTopic := makeTopic(outerCtx, suite.T(), ns, topicName)
	topic, err := ns.NewTopic(topicName)
	if suite.NoError(err) {
		sm := topic.NewSubscriptionManager()
		for name, testFunc := range tests {
			suite.T().Run(name, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
				defer cancel()
				name := suite.RandomName("gosb", 6)
				testFunc(ctx, t, sm, name)
				defer suite.cleanupSubscription(topic, name)
			})
		}
	}
	defer cleanupTopic()
}

func testPutSubscription(ctx context.Context, t *testing.T, sm *SubscriptionManager, name string) {
	sub, err := sm.Put(ctx, name)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	if assert.NotNil(t, sub) {
		assert.Equal(t, name, sub.Name)
	}
}

func (suite *serviceBusSuite) cleanupSubscription(topic *Topic, subscriptionName string) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	sm := topic.NewSubscriptionManager()
	err := sm.Delete(ctx, subscriptionName)
	if err != nil {
		suite.T().Fatal(err)
	}
}

func (suite *serviceBusSuite) TestSubscriptionManager_SubscriptionWithAutoForward() {
	tests := map[string]func(context.Context, *testing.T, *SubscriptionManager, string, string){
		"TestSubscriptionWithAutoForward": testSubscriptionWithAutoForward,
	}

	suite.testSubscriptionManager(tests)
}

func testSubscriptionWithAutoForward(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, subName string) {
	targetName := "target-" + subName
	target := buildSubscription(ctx, t, sm, targetName)
	defer func() {
		assert.NoError(t, sm.Delete(ctx, targetName))
	}()

	src := buildSubscription(ctx, t, sm, subName, SubscriptionWithAutoForward(target))
	assert.Equal(t, target.TargetURI(), *src.ForwardTo)
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

	suite.testSubscriptionManager(tests)
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

func buildSubscription(ctx context.Context, t *testing.T, sm *SubscriptionManager, name string, opts ...SubscriptionManagementOption) *SubscriptionEntity {
	_, err := sm.Put(ctx, name, opts...)
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("%v", err))
	}

	s, err := sm.Get(ctx, name)
	if !assert.NoError(t, err) {
		assert.FailNow(t, fmt.Sprintf("%v", err))
	}
	return s
}

func (suite *serviceBusSuite) testSubscriptionManager(tests map[string]func(context.Context, *testing.T, *SubscriptionManager, string, string)) {
	ns := suite.getNewSasInstance()
	for name, testFunc := range tests {
		topicName := suite.randEntityName()
		subName := suite.randEntityName()

		setupTestTeardown := func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			cleanupTopic := makeTopic(ctx, t, ns, topicName)
			topic, err := ns.NewTopic(topicName)
			if suite.NoError(err) {
				sm := topic.NewSubscriptionManager()
				if suite.NoError(err) {
					defer func(sName string) {
						ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
						defer cancel()
						if !suite.NoError(sm.Delete(ctx, sName)) {
							suite.Fail(err.Error())
						}
					}(subName)
					testFunc(ctx, t, sm, topicName, subName)
				}
			}
			defer cleanupTopic()
		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func (suite *serviceBusSuite) TestSubscription_WithReceiveAndDelete() {
	tests := map[string]func(context.Context, *testing.T, *Topic, *Subscription){
		"ReceiveOne": testSubscriptionReceiveOneNoComplete,
	}

	suite.subscriptionMessageTestWithOptions(tests, SubscriptionWithReceiveAndDelete())
}

func testSubscriptionReceiveOneNoComplete(ctx context.Context, t *testing.T, topic *Topic, sub *Subscription) {
	if assert.NoError(t, topic.Send(ctx, NewMessageFromString("hello!"))) {
		err := sub.ReceiveOne(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
			return nil
		}))
		assert.NoError(t, err)
	}
}

func (suite *serviceBusSuite) TestSubscriptionClient() {
	tests := map[string]func(context.Context, *testing.T, *Topic, *Subscription){
		"SimpleReceive": testSubscriptionReceive,
		"ReceiveOne":    testSubscriptionReceiveOne,
	}

	suite.subscriptionMessageTestWithMgmtOptions(tests)
}

func testSubscriptionReceive(ctx context.Context, t *testing.T, topic *Topic, sub *Subscription) {
	if assert.NoError(t, topic.Send(ctx, NewMessageFromString("hello!"))) {
		inner, cancel := context.WithCancel(ctx)
		err := sub.Receive(inner, HandlerFunc(func(eventCtx context.Context, msg *Message) error {
			defer cancel()
			return msg.Complete(ctx)
		}))
		assert.EqualError(t, err, context.Canceled.Error())
	}
}

func testSubscriptionReceiveOne(ctx context.Context, t *testing.T, topic *Topic, sub *Subscription) {
	if assert.NoError(t, topic.Send(ctx, NewMessageFromString("hello!"))) {
		err := sub.ReceiveOne(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
			return msg.Complete(ctx)
		}))
		assert.NoError(t, err)
	}
}

func (suite *serviceBusSuite) TestSubscription_NewDeadLetter() {
	tests := map[string]func(context.Context, *testing.T, *Topic, *Subscription){
		"ReceiveOneFromDeadLetter": testSubscriptionReceiveOneFromDeadLetter,
	}

	suite.subscriptionMessageTestWithOptions(tests)
}

func testSubscriptionReceiveOneFromDeadLetter(ctx context.Context, t *testing.T, topic *Topic, sub *Subscription) {
	require.NoError(t, topic.Send(ctx, NewMessageFromString("foo")))

	for i := 0; i < 10; i++ {
		err := sub.ReceiveOne(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
			return msg.Abandon(ctx)
		}))
		require.NoError(t, err)
	}

	dl := sub.NewDeadLetter()
	called := false
	err := dl.ReceiveOne(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
		called = true
		assert.Equal(t, "foo", string(msg.Data))
		return msg.Complete(ctx)
	}))
	assert.True(t, called)
	assert.NoError(t, err)
}

func (suite *serviceBusSuite) TestSubscriptionSessionClient() {
	tests := map[string]func(context.Context, *testing.T, *TopicSession, *SubscriptionSession){
		"ReceiveOne": testSubscriptionSessionReceiveOne,
	}

	ns := suite.getNewSasInstance()
	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			topicName := suite.randEntityName()
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()

			topicCleanup := makeTopic(ctx, t, ns, topicName)
			topic, err := ns.NewTopic(topicName)
			if suite.NoError(err) {
				subName := suite.randEntityName()
				subCleanup := makeSubscription(ctx, t, topic, subName, SubscriptionWithRequiredSessions())
				subscription, err := topic.NewSubscription(subName)
				id, err := uuid.NewV4()
				suite.Require().NoError(err)
				sessionID := id.String()

				ts := topic.NewSession(&sessionID)
				ss := subscription.NewSession(&sessionID)
				defer func() {
					closeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
					defer cancel()
					suite.NoError(ts.Close(closeCtx))
					suite.NoError(ss.Close(closeCtx))
					subCleanup()
					topicCleanup()
				}()

				if suite.NoError(err) {
					testFunc(ctx, t, ts, ss)
					if !t.Failed() {
						checkZeroSubscriptionMessages(ctx, t, topic, subName)
					}
				}
			}
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func testSubscriptionSessionReceiveOne(ctx context.Context, t *testing.T, topic *TopicSession, sub *SubscriptionSession) {
	const want = "hello!"
	require.NoError(t, topic.Send(ctx, NewMessageFromString(want)))

	closerCtx, cancel := context.WithCancel(ctx)
	err := sub.ReceiveOne(closerCtx, NewSessionHandler(
		HandlerFunc(func(ctx context.Context, msg *Message) error {
			assert.Equal(t, string(msg.Data), want)
			defer cancel()
			return msg.Complete(ctx)
		}),
		func(ms *MessageSession) error {
			return nil
		},
		func() {}))
	assert.Error(t, err, "context canceled")
}

func (suite *serviceBusSuite) subscriptionMessageTestWithOptions(tests map[string]func(context.Context, *testing.T, *Topic, *Subscription), opts ...SubscriptionOption) {
	suite.subscriptionMessageTest(tests, opts, []SubscriptionManagementOption{})
}

func (suite *serviceBusSuite) subscriptionMessageTestWithMgmtOptions(tests map[string]func(context.Context, *testing.T, *Topic, *Subscription), mgmtOpts ...SubscriptionManagementOption) {
	suite.subscriptionMessageTest(tests, []SubscriptionOption{}, mgmtOpts)
}

func (suite *serviceBusSuite) subscriptionMessageTest(tests map[string]func(context.Context, *testing.T, *Topic, *Subscription), subOpts []SubscriptionOption, mgmtOpts []SubscriptionManagementOption) {
	ns := suite.getNewSasInstance()
	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			topicName := suite.randEntityName()
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()

			topicCleanup := makeTopic(ctx, t, ns, topicName)
			topic, err := ns.NewTopic(topicName)
			if suite.NoError(err) {
				subName := suite.randEntityName()
				subCleanup := makeSubscription(ctx, t, topic, subName, mgmtOpts...)
				subscription, err := topic.NewSubscription(subName, subOpts...)

				if suite.NoError(err) {
					defer subCleanup()
					testFunc(ctx, t, topic, subscription)
					if !t.Failed() && !strings.HasSuffix(name, "NoZeroCheck") {
						checkZeroSubscriptionMessages(ctx, t, topic, subName)
					}
				}
				defer topicCleanup()
			}
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func makeSubscription(ctx context.Context, t *testing.T, topic *Topic, name string, opts ...SubscriptionManagementOption) func() {
	sm := topic.NewSubscriptionManager()
	entity, err := sm.Get(ctx, name)
	if !assert.NoError(t, err) {
		assert.FailNow(t, "could not GET a subscription")
	}

	if entity == nil {
		entity, err = sm.Put(ctx, name, opts...)
		if !assert.NoError(t, err) {
			assert.FailNow(t, "could not PUT a subscription")
		}
	}
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		_ = sm.Delete(ctx, entity.Name)
	}
}

func checkZeroSubscriptionMessages(ctx context.Context, t *testing.T, topic *Topic, name string) {
	sm := topic.NewSubscriptionManager()
	maxTries := 10
	for i := 0; i < maxTries; i++ {
		s, err := sm.Get(ctx, name)
		if !assert.NoError(t, err) {
			return
		}
		if *s.MessageCount == 0 {
			return
		}
		t.Logf("try %d out of %d, message count was %d, not 0", i+1, maxTries, *s.MessageCount)
		time.Sleep(1 * time.Second)
	}

	assert.Fail(t, "message count never reached zero")
}
