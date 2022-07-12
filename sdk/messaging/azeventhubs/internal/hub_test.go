// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

// Delete all of the hubs under a given namespace
// az eventhubs eventhub delete -g ehtest --namespace-name {ns} --ids $(az eventhubs eventhub list -g ehtest --namespace-name {ns} --query "[].id" -o tsv)

import (
	"context"
	"encoding/xml"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/aad"
	"github.com/Azure/azure-amqp-common-go/v3/auth"
	"github.com/Azure/azure-amqp-common-go/v3/sas"
	"github.com/Azure/azure-amqp-common-go/v3/uuid"
	"github.com/Azure/go-amqp"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/test"
)

type (
	// eventHubSuite encapsulates a end to end test of Event Hubs with build up and tear down of all EH resources
	eventHubSuite struct {
		test.BaseSuite
	}
)

var (
	defaultTimeout = 30 * time.Second
)

const (
	connStr = "Endpoint=sb://namespace.servicebus.windows.net/;SharedAccessKeyName=keyName;SharedAccessKey=secret;EntityPath=hubName"

	hubDescription = `
        <EventHubDescription xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect" 
            xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
            <MessageRetentionInDays>7</MessageRetentionInDays>
            <AuthorizationRules></AuthorizationRules>
            <Status>Active</Status>
            <CreatedAt>2018-07-18T17:04:00.67Z</CreatedAt>
            <UpdatedAt>2018-07-18T17:04:00.95Z</UpdatedAt>
            <PartitionCount>4</PartitionCount>
            <PartitionIds xmlns:d2p1="http://schemas.microsoft.com/2003/10/Serialization/Arrays">
                <d2p1:string>0</d2p1:string>
                <d2p1:string>1</d2p1:string>
                <d2p1:string>2</d2p1:string>
                <d2p1:string>3</d2p1:string>
            </PartitionIds>
        </EventHubDescription>
`

	hubEntry1 = `
<entry xmlns="http://www.w3.org/2005/Atom">
    <id>https://foo.servicebus.windows.net/goehcqqrjf-tag3lxnf?api-version=2017-04</id>
    <title type="text">goehcqqrjf-tag3lxnf</title>
    <published>2018-07-18T17:04:00Z</published>
    <updated>2018-07-18T17:04:00Z</updated>
    <author>
        <name>foo</name>
    </author>
    <link rel="self" href="https://foo.servicebus.windows.net/goehcqqrjf-tag3lxnf?api-version=2017-04"/>
    <content type="application/xml">` + hubDescription +
		`</content>
</entry>`

	feedOfHubDescriptions = `
<feed xmlns="http://www.w3.org/2005/Atom">
			<title type="text">Queues</title>
			<id>https://foo.servicebus.windows.net/$Resources/EventHubs</id>
			<updated>2018-05-03T00:21:15Z</updated>
			<link rel="self" href="https://sbdjtest.servicebus.windows.net/$Resources/EventHubs"/>` + hubEntry1 + `</feed>`
)

func TestEH(t *testing.T) {
	suite.Run(t, new(eventHubSuite))
}

func (suite *eventHubSuite) TestNewHubWithNameAndEnvironment() {
	revert := suite.captureEnv()
	defer revert()
	os.Clearenv()
	require.NoError(suite.T(), os.Setenv("EVENTHUB_CONNECTION_STRING", connStr))
	_, err := NewHubWithNamespaceNameAndEnvironment("hello", "world")
	require.NoError(suite.T(), err)
}

func (suite *eventHubSuite) TestUnmarshalHubEntry() {
	var entry hubEntry
	err := xml.Unmarshal([]byte(hubEntry1), &entry)
	suite.Nil(err)
	suite.Require().NotNil(entry)
	suite.Equal("https://foo.servicebus.windows.net/goehcqqrjf-tag3lxnf?api-version=2017-04", entry.ID)
	suite.Equal("goehcqqrjf-tag3lxnf", entry.Title)
	suite.Require().NotNil(entry.Author)
	suite.Equal("foo", *entry.Author.Name)
	suite.Require().NotNil(entry.Link)
	suite.Equal("https://foo.servicebus.windows.net/goehcqqrjf-tag3lxnf?api-version=2017-04", entry.Link.HREF)
	suite.Require().NotNil(entry.Content)
	suite.Require().NotNil(entry.Content.HubDescription)
	suite.Require().NotNil(entry.Content.HubDescription.PartitionCount)
	suite.Equal(int32(4), *entry.Content.HubDescription.PartitionCount)
	suite.Require().NotNil(entry.Content.HubDescription.PartitionIDs)
	suite.Equal([]string{"0", "1", "2", "3"}, *entry.Content.HubDescription.PartitionIDs)
	suite.Require().NotNil(entry.Content.HubDescription.MessageRetentionInDays)
	suite.Equal(int32(7), *entry.Content.HubDescription.MessageRetentionInDays)
}

func (suite *eventHubSuite) TestUnmarshalHubList() {
	var feed hubFeed
	suite.Require().NoError(xml.Unmarshal([]byte(feedOfHubDescriptions), &feed))
	suite.Require().NotNil(feed)
	suite.Require().NotNil(feed.Entries)
	suite.Len(feed.Entries, 1)
}

func (suite *eventHubSuite) TestHubManagementWrites() {
	tests := map[string]func(context.Context, *testing.T, *HubManager, string){
		"TestPutDefaultHub": testPutHub,
	}

	hm, err := NewHubManagerFromConnectionString(os.Getenv("EVENTHUB_CONNECTION_STRING"))
	suite.Require().NoError(err)

	for name, testFunc := range tests {
		suite.T().Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			name := suite.randEntityName()
			testFunc(ctx, t, hm, name)
			defer suite.DeleteEventHub(name)
		})
	}
}

func testPutHub(ctx context.Context, t *testing.T, hm *HubManager, name string) {
	hd, err := hm.Put(ctx, name)
	require.NoError(t, err)
	require.NotNil(t, hd.PartitionCount)
	assert.Equal(t, int32(4), *hd.PartitionCount)
	require.NotNil(t, hd.PartitionIDs)
	assert.Equal(t, []string{"0", "1", "2", "3"}, *hd.PartitionIDs)
	require.NotNil(t, hd.MessageRetentionInDays)
	assert.Equal(t, int32(7), *hd.MessageRetentionInDays)
}

func (suite *eventHubSuite) TestHubManagementReads() {
	tests := map[string]func(context.Context, *testing.T, *HubManager, []string){
		"TestGetHub":   testGetHub,
		"TestListHubs": testListHubs,
	}

	hm, err := NewHubManagerFromConnectionString(os.Getenv("EVENTHUB_CONNECTION_STRING"))
	suite.Require().NoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	names := []string{suite.randEntityName(), suite.randEntityName()}
	for _, name := range names {
		if _, err := hm.Put(ctx, name); err != nil {
			suite.Require().NoError(err)
		}
	}

	for name, testFunc := range tests {
		suite.T().Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			testFunc(ctx, t, hm, names)
		})
	}

	for _, name := range names {
		suite.DeleteEventHub(name)
	}
}

func testGetHub(ctx context.Context, t *testing.T, hm *HubManager, names []string) {
	q, err := hm.Get(ctx, names[0])
	require.NoError(t, err)
	require.NotNil(t, q)
	assert.Equal(t, q.Name, names[0])
}

func testListHubs(ctx context.Context, t *testing.T, hm *HubManager, names []string) {
	hubs, err := hm.List(ctx)
	require.NoError(t, err)
	require.NotNil(t, hubs)
	hubNames := make([]string, len(hubs))
	for idx, q := range hubs {
		hubNames[idx] = q.Name
	}

	for _, name := range names {
		assert.Contains(t, hubNames, name)
	}
}

func (suite *eventHubSuite) randEntityName() string {
	return suite.RandomName("goeh", 6)
}

func (suite *eventHubSuite) TestSasToken() {
	tests := map[string]func(context.Context, *testing.T, *Hub, []string, string){
		"TestMultiSendAndReceive":            testMultiSendAndReceive,
		"TestHubRuntimeInformation":          testHubRuntimeInformation,
		"TestHubPartitionRuntimeInformation": testHubPartitionRuntimeInformation,
	}

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			provider, err := sas.NewTokenProvider(sas.TokenProviderWithEnvironmentVars())
			suite.Require().NoError(err)
			hub, cleanup := suite.RandomHub()
			defer cleanup()
			client, closer := suite.newClientWithProvider(t, *hub.Name, provider)
			defer closer()
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			testFunc(ctx, t, client, *hub.PartitionIds, *hub.Name)
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func (suite *eventHubSuite) TestPartitioned() {
	tests := map[string]func(context.Context, *testing.T, *Hub, string){
		"TestSend":                testBasicSend,
		"TestSendTooBig":          testSendTooBig,
		"TestSendAndReceive":      testBasicSendAndReceive,
		"TestBatchSendAndReceive": testBatchSendAndReceive,
		"TestBatchSendTooLarge":   testBatchSendTooLarge,
	}

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			hub, cleanup := suite.RandomHub()
			defer cleanup()
			partitionID := (*hub.PartitionIds)[0]
			client, closer := suite.newClient(t, *hub.Name, HubWithPartitionedSender(partitionID))
			defer closer()
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			testFunc(ctx, t, client, partitionID)
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func (suite *eventHubSuite) TestWebSocket() {
	tests := map[string]func(context.Context, *testing.T, *Hub, string){
		"TestSend":                testBasicSend,
		"TestSendTooBig":          testSendTooBig,
		"TestSendAndReceive":      testBasicSendAndReceive,
		"TestBatchSendAndReceive": testBatchSendAndReceive,
		"TestBatchSendTooLarge":   testBatchSendTooLarge,
	}

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			hub, cleanup := suite.RandomHub()
			defer cleanup()
			partitionID := (*hub.PartitionIds)[0]
			client, closer := suite.newClient(t, *hub.Name, HubWithPartitionedSender(partitionID), HubWithWebSocketConnection())
			defer closer()
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			testFunc(ctx, t, client, partitionID)
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func (suite *eventHubSuite) TestSenderRetryOptionsThroughHub() {
	tests := map[string]func(ctx context.Context, t *testing.T, suite *eventHubSuite, hubName string){
		"TestWithCustomSenderMaxRetryCount": func(ctx context.Context, t *testing.T, suite *eventHubSuite, hubName string) {
			client, closer := suite.newClient(t, hubName, HubWithSenderMaxRetryCount(3))
			defer closer()

			err := client.Send(ctx, &Event{
				Data: []byte("hello world"),
			})

			assert.NoError(t, err)
			assert.EqualValues(t, 3, client.sender.retryOptions.maxRetries)
			assert.EqualValues(t, client.sender.retryOptions.recoveryBackoff, newSenderRetryOptions().recoveryBackoff)
		},
		"TestWithDefaultSenderMaxRetryCount": func(ctx context.Context, t *testing.T, suite *eventHubSuite, hubName string) {
			client, closer := suite.newClient(t, hubName)
			defer closer()

			err := client.Send(ctx, &Event{
				Data: []byte("hello world"),
			})

			assert.NoError(t, err)
			assert.EqualValues(t, -1, client.sender.retryOptions.maxRetries)
			assert.EqualValues(t, client.sender.retryOptions.recoveryBackoff, newSenderRetryOptions().recoveryBackoff)
		},
	}

	hub, cleanup := suite.RandomHub()
	defer cleanup()

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()

			testFunc(ctx, t, suite, *hub.Name)
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func testBasicSend(ctx context.Context, t *testing.T, client *Hub, _ string) {
	err := client.Send(ctx, NewEventFromString("Hello!"))
	assert.NoError(t, err)
}

func testSendTooBig(ctx context.Context, t *testing.T, client *Hub, _ string) {
	data := make([]byte, 2600*1024)
	_, _ = rand.Read(data)
	event := NewEvent(data)
	err := client.Send(ctx, event)
	assert.Error(t, err, "encoded message size exceeds max of 1048576")
}

func testBatchSendAndReceive(ctx context.Context, t *testing.T, client *Hub, partitionID string) {
	messages := []string{"hello", "world", "foo", "bar", "baz", "buzz"}
	var wg sync.WaitGroup
	wg.Add(len(messages))

	events := make([]*Event, len(messages))
	for idx, msg := range messages {
		events[idx] = NewEventFromString(msg)
	}

	ebi := NewEventBatchIterator(events...)
	if assert.NoError(t, client.SendBatch(ctx, ebi)) {
		count := 0
		_, err := client.Receive(context.Background(), partitionID, func(ctx context.Context, event *Event) error {
			assert.Equal(t, messages[count], string(event.Data))
			count++
			wg.Done()
			return nil
		}, ReceiveWithPrefetchCount(100))
		if !assert.NoError(t, err) {
			end, _ := ctx.Deadline()
			waitUntil(t, &wg, time.Until(end))
		}
	}
}

func testBatchSendTooLarge(ctx context.Context, t *testing.T, client *Hub, _ string) {
	events := make([]*Event, 200000)

	var wg sync.WaitGroup
	wg.Add(len(events))

	for idx := range events {
		events[idx] = NewEventFromString(test.RandomString("foo", 10))
	}

	ebi := NewEventBatchIterator(events...)
	assert.EqualError(t, client.SendBatch(ctx, ebi, BatchWithMaxSizeInBytes(10000000)), "encoded message size exceeds max of 1048576")
}

func testBasicSendAndReceive(ctx context.Context, t *testing.T, client *Hub, partitionID string) {
	numMessages := rand.Intn(100) + 20
	var wg sync.WaitGroup
	wg.Add(numMessages)

	messages := make([]string, numMessages)
	for i := 0; i < numMessages; i++ {
		messages[i] = test.RandomString("hello", 10)
	}

	for idx, message := range messages {
		event := NewEventFromString(message)
		event.SystemProperties = &SystemProperties{
			Annotations: map[string]interface{}{
				"x-opt-custom-annotation": "custom-value",
			},
		}
		if !assert.NoError(t, client.Send(ctx, event, SendWithMessageID(fmt.Sprintf("%d", idx)))) {
			assert.FailNow(t, "unable to send event")
		}
	}

	count := 0
	_, err := client.Receive(ctx, partitionID, func(ctx context.Context, event *Event) error {
		assert.Equal(t, messages[count], string(event.Data))
		require.NotNil(t, event.SystemProperties)
		assert.NotNil(t, event.SystemProperties.EnqueuedTime)
		assert.NotNil(t, event.SystemProperties.Offset)
		assert.NotNil(t, event.SystemProperties.SequenceNumber)
		assert.Equal(t, int64(count), *event.SystemProperties.SequenceNumber)
		require.NotNil(t, event.SystemProperties.Annotations)
		assert.Equal(t, *event.SystemProperties.EnqueuedTime, event.SystemProperties.Annotations["x-opt-enqueued-time"].(time.Time))
		assert.Equal(t, strconv.FormatInt(*event.SystemProperties.Offset, 10), event.SystemProperties.Annotations["x-opt-offset"].(string))
		assert.Equal(t, *event.SystemProperties.SequenceNumber, event.SystemProperties.Annotations["x-opt-sequence-number"].(int64))
		assert.Equal(t, "custom-value", event.SystemProperties.Annotations["x-opt-custom-annotation"].(string))
		count++
		wg.Done()
		return nil
	}, ReceiveWithPrefetchCount(100))
	if !assert.NoError(t, err) {
		end, _ := ctx.Deadline()
		waitUntil(t, &wg, time.Until(end))
	}
}

func (suite *eventHubSuite) TestEpochReceivers() {
	tests := map[string]func(context.Context, *testing.T, *Hub, []string, string){
		"TestEpochGreaterThenLess": testEpochGreaterThenLess,
		"TestEpochLessThenGreater": testEpochLessThenGreater,
	}

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			hub, cleanup := suite.RandomHub()
			require.NotNil(t, hub)
			defer cleanup()
			require.Len(t, *hub.PartitionIds, 4)
			partitionID := (*hub.PartitionIds)[0]
			client, closer := suite.newClient(t, *hub.Name, HubWithPartitionedSender(partitionID))
			defer closer()
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			testFunc(ctx, t, client, *hub.PartitionIds, *hub.Name)
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func testEpochGreaterThenLess(ctx context.Context, t *testing.T, client *Hub, partitionIDs []string, _ string) {
	partitionID := partitionIDs[0]
	r1, err := client.Receive(ctx, partitionID, func(c context.Context, event *Event) error { return nil }, ReceiveWithEpoch(4))
	if !assert.NoError(t, err) {
		assert.FailNow(t, "error receiving with epoch of 4")
	}

	_, err = client.Receive(ctx, partitionID, func(c context.Context, event *Event) error { return nil }, ReceiveWithEpoch(1))

	assert.Error(t, err, "receiving with epoch of 1 should have failed")
	assert.NoError(t, r1.Err(), "r1 should still be running with the higher epoch")
}

func testEpochLessThenGreater(ctx context.Context, t *testing.T, client *Hub, partitionIDs []string, _ string) {
	partitionID := partitionIDs[0]
	r1, err := client.Receive(ctx, partitionID, func(c context.Context, event *Event) error { return nil }, ReceiveWithEpoch(1))
	if !assert.NoError(t, err) {
		assert.FailNow(t, "error receiving with epoch of 1")
	}

	r2, err := client.Receive(ctx, partitionID, func(c context.Context, event *Event) error { return nil }, ReceiveWithEpoch(4))
	if !assert.NoError(t, err) {
		assert.FailNow(t, "error receiving with epoch of 4")
	}

	select {
	case <-r1.Done():
		break
	case <-ctx.Done():
		assert.FailNow(t, "r1 didn't finish in time")
	}

	assert.Error(t, r1.Err(), "r1 should have died with error since it has a lower epoch value")
	assert.NoError(t, r2.Err(), "r2 should not have an error and should still be processing")
}

func (suite *eventHubSuite) TestMultiPartition() {
	tests := map[string]func(context.Context, *testing.T, *Hub, []string, string){
		"TestMultiSendAndReceive":            testMultiSendAndReceive,
		"TestSendWithPartitionKeyAndReceive": testSendWithPartitionKeyAndReceive,
	}

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			hub, cleanup := suite.RandomHub()
			defer cleanup()
			client, closer := suite.newClient(t, *hub.Name)
			defer closer()
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			testFunc(ctx, t, client, *hub.PartitionIds, *hub.Name)
		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func testMultiSendAndReceive(ctx context.Context, t *testing.T, client *Hub, partitionIDs []string, _ string) {
	numMessages := rand.Intn(100) + 20
	var wg sync.WaitGroup
	wg.Add(numMessages)

	messages := make([]string, numMessages)
	for i := 0; i < numMessages; i++ {
		messages[i] = test.RandomString("hello", 10)
	}

	for idx, message := range messages {
		require.NoError(t, client.Send(ctx, NewEventFromString(message), SendWithMessageID(fmt.Sprintf("%d", idx))))
	}

	for _, partitionID := range partitionIDs {
		_, err := client.Receive(ctx, partitionID, func(ctx context.Context, event *Event) error {
			wg.Done()
			return nil
		}, ReceiveWithPrefetchCount(100))
		if !assert.NoError(t, err) {
			assert.FailNow(t, "unable to setup receiver")
		}
	}
	end, _ := ctx.Deadline()
	waitUntil(t, &wg, time.Until(end))
}

func testSendWithPartitionKeyAndReceive(ctx context.Context, t *testing.T, client *Hub, partitionIDs []string, _ string) {
	numMessages := rand.Intn(100) + 20
	var wg sync.WaitGroup
	wg.Add(numMessages)

	sentEvents := make(map[string]*Event)
	for i := 0; i < numMessages; i++ {
		content := test.RandomString("hello", 10)
		event := NewEventFromString(content)
		event.PartitionKey = &content
		id, err := uuid.NewV4()
		if !assert.NoError(t, err) {
			assert.FailNow(t, "error generating uuid")
		}
		event.ID = id.String()
		sentEvents[event.ID] = event
	}

	for _, event := range sentEvents {
		if !assert.NoError(t, client.Send(ctx, event)) {
			assert.FailNow(t, "failed to send message to hub")
		}
	}

	received := make(map[string][]*Event)
	for _, pID := range partitionIDs {
		received[pID] = []*Event{}
	}
	for _, partitionID := range partitionIDs {
		_, err := client.Receive(ctx, partitionID, func(ctx context.Context, event *Event) error {
			defer wg.Done()
			received[partitionID] = append(received[partitionID], event)
			return nil
		}, ReceiveWithPrefetchCount(100))
		if !assert.NoError(t, err) {
			assert.FailNow(t, "failed to receive from partition")
		}
	}

	// wait for events to arrive
	end, _ := ctx.Deadline()
	if waitUntil(t, &wg, time.Until(end)) {
		// collect all of the partitioned events
		receivedEventsByID := make(map[string]*Event)
		for _, pID := range partitionIDs {
			for _, event := range received[pID] {
				receivedEventsByID[event.ID] = event
			}
		}

		// verify the sent events have the same partition keys as the received events
		for key, event := range sentEvents {
			assert.Equal(t, event.PartitionKey, receivedEventsByID[key].PartitionKey)
		}
	}
}

func (suite *eventHubSuite) TestHubManagement() {
	tests := map[string]func(context.Context, *testing.T, *Hub, []string, string){
		"TestHubRuntimeInformation":          testHubRuntimeInformation,
		"TestHubPartitionRuntimeInformation": testHubPartitionRuntimeInformation,
	}

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			hub, cleanup := suite.RandomHub()
			defer cleanup()
			client, closer := suite.newClient(t, *hub.Name)
			defer closer()
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			testFunc(ctx, t, client, *hub.PartitionIds, *hub.Name)
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func testHubRuntimeInformation(ctx context.Context, t *testing.T, client *Hub, partitionIDs []string, hubName string) {
	info, err := client.GetRuntimeInformation(ctx)
	if assert.NoError(t, err) {
		assert.Equal(t, len(partitionIDs), info.PartitionCount)
		assert.Equal(t, hubName, info.Path)
	}
}

func testHubPartitionRuntimeInformation(ctx context.Context, t *testing.T, client *Hub, partitionIDs []string, hubName string) {
	info, err := client.GetPartitionInformation(ctx, partitionIDs[0])
	if assert.NoError(t, err) {
		assert.Equal(t, hubName, info.HubPath)
		assert.Equal(t, partitionIDs[0], info.PartitionID)
		assert.Equal(t, "-1", info.LastEnqueuedOffset) // brand new, so should be very last
	}
}

func (suite *eventHubSuite) newClient(t *testing.T, hubName string, opts ...HubOption) (*Hub, func()) {
	provider, err := aad.NewJWTProvider(
		aad.JWTProviderWithEnvironmentVars(),
		aad.JWTProviderWithAzureEnvironment(&suite.Env),
	)
	if !suite.NoError(err) {
		suite.FailNow("unable to make a new JWT provider")
	}
	return suite.newClientWithProvider(t, hubName, provider, opts...)
}

func (suite *eventHubSuite) newClientWithProvider(t *testing.T, hubName string, provider auth.TokenProvider, opts ...HubOption) (*Hub, func()) {
	opts = append(opts, HubWithEnvironment(suite.Env))
	client, err := NewHub(suite.Namespace, hubName, provider, opts...)
	if !suite.NoError(err) {
		suite.FailNow("unable to make a new Hub")
	}
	return client, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = client.Close(ctx)
	}
}

func waitUntil(t *testing.T, wg *sync.WaitGroup, d time.Duration) bool {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return true
	case <-time.After(d):
		assert.Fail(t, "took longer than "+fmtDuration(d))
		return false
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second) / time.Second
	return fmt.Sprintf("%d seconds", d)
}

func restoreEnv(capture map[string]string) error {
	os.Clearenv()
	for key, value := range capture {
		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (suite *eventHubSuite) captureEnv() func() {
	capture := make(map[string]string)
	for _, pair := range os.Environ() {
		keyValue := strings.Split(pair, "=")
		capture[keyValue[0]] = strings.Join(keyValue[1:], "=")
	}
	return func() {
		suite.NoError(restoreEnv(capture))
	}
}

func TestNewHub_withAzureEnvironmentVariable(t *testing.T) {
	_ = os.Setenv("AZURE_ENVIRONMENT", "AZURECHINACLOUD")
	h, err := NewHub("test", "test", &aad.TokenProvider{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(h.namespace.host, azure.ChinaCloud.ServiceBusEndpointSuffix) {
		t.Fatalf("did not set appropriate endpoint suffix. Expected: %v, Received: %v", azure.ChinaCloud.ServiceBusEndpointSuffix, h.namespace.host)
	}
	_ = os.Setenv("AZURE_ENVIRONMENT", "AZUREGERMANCLOUD")
	h, err = NewHub("test", "test", &aad.TokenProvider{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(h.namespace.host, azure.GermanCloud.ServiceBusEndpointSuffix) {
		t.Fatalf("did not set appropriate endpoint suffix. Expected: %v, Received: %v", azure.GermanCloud.ServiceBusEndpointSuffix, h.namespace.host)
	}
	_ = os.Setenv("AZURE_ENVIRONMENT", "AZUREUSGOVERNMENTCLOUD")
	h, err = NewHub("test", "test", &aad.TokenProvider{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(h.namespace.host, azure.USGovernmentCloud.ServiceBusEndpointSuffix) {
		t.Fatalf("did not set appropriate endpoint suffix. Expected: %v, Received: %v", azure.USGovernmentCloud.ServiceBusEndpointSuffix, h.namespace.host)
	}
	_ = os.Setenv("AZURE_ENVIRONMENT", "AZUREPUBLICCLOUD")
	h, err = NewHub("test", "test", &aad.TokenProvider{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(h.namespace.host, azure.PublicCloud.ServiceBusEndpointSuffix) {
		t.Fatalf("did not set appropriate endpoint suffix. Expected: %v, Received: %v", azure.PublicCloud.ServiceBusEndpointSuffix, h.namespace.host)
	}
	_ = os.Unsetenv("AZURE_ENVIRONMENT")
	h, err = NewHub("test", "test", &aad.TokenProvider{})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(h.namespace.host, azure.PublicCloud.ServiceBusEndpointSuffix) {
		t.Fatalf("did not set appropriate endpoint suffix. Expected: %v, Received: %v", azure.PublicCloud.ServiceBusEndpointSuffix, h.namespace.host)
	}
}

func TestIsRecoverableCloseError(t *testing.T) {
	require.True(t, isRecoverableCloseError(&amqp.DetachError{}))

	// if the caller closes a link we shouldn't reopen or create a new one to replace it
	require.False(t, isRecoverableCloseError(amqp.ErrLinkClosed))
}
