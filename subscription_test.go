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
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-amqp-common-go/uuid"
	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2015-08-01/servicebus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Azure/azure-service-bus-go/internal/test"
)

const (
	ruleDescription = `
		<RuleDescription xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect"
			xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
			<Filter i:type="TrueFilter">
				<SqlExpression>1=1</SqlExpression>
				<CompatibilityLevel>20</CompatibilityLevel>
			</Filter>
			<Action i:type="EmptyRuleAction"/>
			<CreatedAt>2018-12-19T19:37:23.9128676Z</CreatedAt>
			<Name>$Default</Name>
		</RuleDescription>
`
	ruleEntryContent = `
			<entry xml:base="https://azuresbtests-bmjlfhyx.servicebus.windows.net/goqsgk5tt-tagtdd84njqj7/subscriptions/goqe6ci2n-tagtdd84njqj7/rules?api-version=2017-04">
				<id>https://azuresbtests-bmjlfhyx.servicebus.windows.net/goqsgk5tt-tagtdd84njqj7/subscriptions/goqe6ci2n-tagtdd84njqj7/rules/$Default?api-version=2017-04</id>
				<title type="text">$Default</title>
				<published>2018-12-19T19:37:23Z</published>
				<updated>2018-12-19T19:37:23Z</updated>
				<link rel="self" href="rules/$Default?api-version=2017-04"/>
				<content type="application/xml"> ` + ruleDescription + `
				</content>
			</entry>
`

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

func (suite *serviceBusSuite) TestSubscriptionRuleEntry_Unmarshal() {
	var entry ruleEntry
	err := xml.Unmarshal([]byte(ruleEntryContent), &entry)
	suite.NoError(err)
	suite.Equal("https://azuresbtests-bmjlfhyx.servicebus.windows.net/goqsgk5tt-tagtdd84njqj7/subscriptions/goqe6ci2n-tagtdd84njqj7/rules/$Default?api-version=2017-04", entry.ID)
	suite.Equal("$Default", entry.Title)
	suite.Equal("rules/$Default?api-version=2017-04", entry.Link.HREF)
	suite.Require().NotNil(entry.Content)
	suite.Require().NotNil(entry.Content.RuleDescription)
	//suite.Equal("TrueFilter", entry.Content.RuleDescription.Filter.Attributes)
	suite.Equal("TrueFilter", entry.Content.RuleDescription.Filter.Type)
}

func (suite *serviceBusSuite) TestSubscriptionEntry_Unmarshal() {
	var entry subscriptionEntry
	err := xml.Unmarshal([]byte(subscriptionEntryContent), &entry)
	suite.NoError(err)
	suite.Equal("https://sbdjtest.servicebus.windows.net/gosbh6of3g-tagz3cfzrp93m/subscriptions/gosbwg424p-tagz3cfzrp93m?api-version=2017-04", entry.ID)
	suite.Equal("gosbwg424p-tagz3cfzrp93m", entry.Title)
	suite.Equal("https://sbdjtest.servicebus.windows.net/gosbh6of3g-tagz3cfzrp93m/subscriptions/gosbwg424p-tagz3cfzrp93m?api-version=2017-04", entry.Link.HREF)
	suite.Require().NotNil(entry.Content)
	suite.Require().NotNil(entry.Content.SubscriptionDescription)
	suite.Equal("PT1M", *entry.Content.SubscriptionDescription.LockDuration)
}

func (suite *serviceBusSuite) TestSubscriptionEntity_Unmarshal() {
	var entry subscriptionEntry
	err := xml.Unmarshal([]byte(subscriptionEntryContent), &entry)
	suite.NoError(err)
	s := entry.Content.SubscriptionDescription
	suite.Equal("PT1M", *s.LockDuration)
	suite.Equal(false, *s.RequiresSession)
	suite.Equal("P10675199DT2H48M5.4775807S", *s.DefaultMessageTimeToLive)
	suite.Equal(false, *s.DeadLetteringOnMessageExpiration)
	suite.Equal(int32(10), *s.MaxDeliveryCount)
	suite.Equal(true, *s.EnableBatchedOperations)
	suite.Equal(int64(0), *s.MessageCount)
	suite.EqualValues(servicebus.EntityStatusActive, *s.Status)
}

func (suite *serviceBusSuite) TestSubscriptionManager_NotFound() {
	ns := suite.getNewSasInstance()
	sm, err := ns.NewSubscriptionManager("foo")
	suite.Require().NoError(err)
	subEntity, err := sm.Get(context.Background(), "bar")
	suite.Nil(subEntity)
	suite.Require().NotNil(err)
	suite.True(IsErrNotFound(err))
	suite.Equal("entity at /foo/subscriptions/bar not found", err.Error())
}

func (suite *serviceBusSuite) TestSubscriptionManagement_Writes() {
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

func (suite *serviceBusSuite) TestSubscriptionManager_SubscriptionWithForwarding() {
	tests := map[string]func(context.Context, *testing.T, *SubscriptionManager, string, string){
		"TestSubscriptionWithAutoForward":         testSubscriptionWithAutoForward,
		"TestSubscriptionWithForwardDeadLetterTo": testSubscriptionWithForwardDeadLetterTo,
	}

	suite.testSubscriptionManager(tests)
}

func testSubscriptionWithForwardDeadLetterTo(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, subName string) {
	targetName := "target-" + subName
	target := buildSubscription(ctx, t, sm, targetName)
	defer func() {
		assert.NoError(t, sm.Delete(ctx, targetName))
	}()

	src := buildSubscription(ctx, t, sm, subName, SubscriptionWithForwardDeadLetteredMessagesTo(target))
	assert.Equal(t, target.TargetURI(), *src.ForwardDeadLetteredMessagesTo)
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
	tests := map[string]func(ctx context.Context, t *testing.T, sm *SubscriptionManager, topicName, name string){
		"TestSubscriptionDefaultSettings":                      testDefaultSubscription,
		"TestSubscriptionWithAutoDeleteOnIdle":                 testSubscriptionWithAutoDeleteOnIdle,
		"TestSubscriptionWithRequiredSessions":                 testSubscriptionWithRequiredSessions,
		"TestSubscriptionWithDeadLetteringOnMessageExpiration": testSubscriptionWithDeadLetteringOnMessageExpiration,
		"TestSubscriptionWithMessageTimeToLive":                testSubscriptionWithMessageTimeToLive,
		"TestSubscriptionWithLockDuration":                     testSubscriptionWithLockDuration,
		"TestSubscriptionWithBatchedOperations":                testSubscriptionWithBatchedOperations,
		"TestSubscriptionWithDefaultRule":                      testSubscriptionWithDefaultRule,
		"TestSubscriptionWithFalseRule":                        testSubscriptionWithFalseRule,
		"TestSubscriptionWithSQLFilterRule":                    testSubscriptionWithSQLFilterRule,
		"TestSubscriptionWithCorrelationFilterRule":            testSubscriptionWithCorrelationFilterRule,
		"TestSubscriptionWithDeleteRule":                       testSubscriptionWithDeleteRule,
		"TestSubscriptionWithActionRule":                       testSubscriptionWithActionRule,
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

func testSubscriptionWithDefaultRule(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, name string) {
	s := buildSubscription(ctx, t, sm, name)
	rules, err := sm.ListRules(ctx, s.Name)
	require.NoError(t, err)
	require.Len(t, rules, 1)
	rule := rules[0]
	assert.Equal(t, rule.Filter.Type, "TrueFilter")
	assert.Equal(t, *rule.Filter.SQLExpression, "1=1")
}

func testSubscriptionWithFalseRule(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, name string) {
	s := buildSubscription(ctx, t, sm, name)
	_, err := sm.PutRule(ctx, s.Name, "falseRule", FalseFilter{})
	require.NoError(t, err)
	rules, err := sm.ListRules(ctx, s.Name)
	require.Len(t, rules, 2)
	rule := rules[1]
	assert.Equal(t, "falseRule", rule.Name)
	assert.Equal(t, "FalseFilter", rule.Filter.Type)
	assert.Equal(t, "1=0", *rule.Filter.SQLExpression)
}

func testSubscriptionWithSQLFilterRule(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, name string) {
	s := buildSubscription(ctx, t, sm, name)
	_, err := sm.PutRule(ctx, s.Name, "sqlRuleNotNullLabel", SQLFilter{Expression: "label IS NOT NULL"})
	require.NoError(t, err)
	rules, err := sm.ListRules(ctx, s.Name)
	require.Len(t, rules, 2)
	rule := rules[1]
	assert.Equal(t, "sqlRuleNotNullLabel", rule.Name)
	assert.Equal(t, "SqlFilter", rule.Filter.Type)
	assert.Equal(t, "label IS NOT NULL", *rule.Filter.SQLExpression)
}

func testSubscriptionWithCorrelationFilterRule(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, name string) {
	s := buildSubscription(ctx, t, sm, name)

	// filter messages that only contain label=="foo" and CorrelationID=="bazz"
	filter := CorrelationFilter{
		Label:         ptrString("foo"),
		CorrelationID: ptrString("bazz"),
	}
	_, err := sm.PutRule(ctx, s.Name, "correlationRule", filter)
	require.NoError(t, err)
	rules, err := sm.ListRules(ctx, s.Name)
	require.Len(t, rules, 2)
	rule := rules[1]
	assert.Equal(t, "correlationRule", rule.Name)
	assert.Equal(t, "CorrelationFilter", rule.Filter.Type)
	assert.Equal(t, "bazz", *rule.Filter.CorrelationID)
	assert.Equal(t, "foo", *rule.Filter.Label)
}

func testSubscriptionWithDeleteRule(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, name string) {
	s := buildSubscription(ctx, t, sm, name)
	_, err := sm.PutRule(ctx, s.Name, "falseRule", FalseFilter{})
	require.NoError(t, err)

	rules, err := sm.ListRules(ctx, s.Name)
	require.NoError(t, err)
	require.Len(t, rules, 2)

	assert.NoError(t, sm.DeleteRule(ctx, s.Name, "falseRule"))

	rules, err = sm.ListRules(ctx, s.Name)
	assert.NoError(t, err)
	assert.Len(t, rules, 1)
}

func testSubscriptionWithActionRule(ctx context.Context, t *testing.T, sm *SubscriptionManager, _, name string) {
	s := buildSubscription(ctx, t, sm, name)
	action := SQLAction{Expression: `set label = "1"`}
	_, err := sm.PutRuleWithAction(ctx, s.Name, "actionOnAll", TrueFilter{}, action)
	require.NoError(t, err)

	rules, err := sm.ListRules(ctx, s.Name)
	require.NoError(t, err)
	require.Len(t, rules, 2)

	rule := rules[1]
	assert.Equal(t, action.Expression, rule.Action.SQLExpression)
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

func (suite *serviceBusSuite) TestSubscription_TwoWithReceiveAndDelete() {
	ns := suite.getNewSasInstance()
	topicName := suite.randEntityName()
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	topicCleanup := makeTopic(ctx, suite.T(), ns, topicName)
	defer topicCleanup()
	topic, err := ns.NewTopic(topicName)
	if suite.NoError(err) {
		subName1 := suite.randEntityName()
		subCleanup1 := makeSubscription(ctx, suite.T(), topic, subName1)
		defer subCleanup1()
		subName2 := suite.randEntityName()
		subCleanup2 := makeSubscription(ctx, suite.T(), topic, subName2)
		defer subCleanup2()

		sub1, err := topic.NewSubscription(subName1, SubscriptionWithReceiveAndDelete())
		suite.Require().NoError(err)

		sub2, err := topic.NewSubscription(subName2, SubscriptionWithReceiveAndDelete())
		suite.Require().NoError(err)

		wg := sync.WaitGroup{}
		wg.Add(2)

		go func(){
			suite.Require().NoError(sub1.ReceiveOne(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
				wg.Done()
				return nil
			})))
		}()

		suite.Require().NoError(topic.Send(ctx, NewMessageFromString("foo")))

		go func(){
			suite.Require().NoError(sub2.ReceiveOne(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
				wg.Done()
				return nil
			})))
		}()

		suite.Require().NoError(topic.Send(ctx, NewMessageFromString("foo")))

		wg.Wait()
		suite.NoError(sub1.Close(ctx))
		suite.NoError(sub2.Close(ctx))
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
		"Defer":         testSubscriptionDeferMessage,
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

func testSubscriptionDeferMessage(ctx context.Context, t *testing.T, topic *Topic, sub *Subscription) {
	rmsg := test.RandomString("foo", 10)
	require.NoError(t, topic.Send(ctx, NewMessageFromString(fmt.Sprintf("hello %s!", rmsg))))

	var sequenceNumber *int64
	err := sub.ReceiveOne(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
		sequenceNumber = msg.SystemProperties.SequenceNumber
		return msg.Defer(ctx)
	}))
	require.NoError(t, err)
	require.NotNil(t, sequenceNumber)

	handled := false
	err = sub.ReceiveDeferred(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
		handled = true
		return msg.Complete(ctx)
	}), *sequenceNumber)

	assert.True(t, handled, "expected message handler to be called")
	assert.NoError(t, err)
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

func (suite *serviceBusSuite) TestSubscriptionWithPrefetch() {
	tests := map[string]func(context.Context, *testing.T, *Topic, *Subscription){
		"SendAndReceive": testSubscriptionSendAndReceive,
	}
	suite.subscriptionMessageTestWithOptions(tests, SubscriptionWithPrefetchCount(10))
}

func testSubscriptionSendAndReceive(ctx context.Context, t *testing.T, topic *Topic, s *Subscription) {
	messages := []string{"foo", "bar", "bazz", "buzz"}
	for _, msg := range messages {
		require.NoError(t, topic.Send(ctx, NewMessageFromString(msg)))
	}

	count := 0
	for idx := range messages {
		err := s.ReceiveOne(ctx, HandlerFunc(func(ctx context.Context, msg *Message) error {
			assert.Equal(t, messages[idx], string(msg.Data))
			count++
			return msg.Complete(ctx)
		}))
		assert.NoError(t, err)
	}
	assert.Len(t, messages, count)
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
	if assert.Error(t, err) && !IsErrNotFound(err) {
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
