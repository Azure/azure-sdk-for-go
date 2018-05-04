package servicebus

import (
	"context"
	"encoding/xml"
	"log"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2015-08-01/servicebus"
	"github.com/stretchr/testify/assert"
)

const (
	topicDescription1 = `
		<TopicDescription 
            xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect" 
            xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
            <DefaultMessageTimeToLive>P10675199DT2H48M5.4775807S</DefaultMessageTimeToLive>
            <MaxSizeInMegabytes>1024</MaxSizeInMegabytes>
            <RequiresDuplicateDetection>false</RequiresDuplicateDetection>
            <DuplicateDetectionHistoryTimeWindow>PT10M</DuplicateDetectionHistoryTimeWindow>
            <EnableBatchedOperations>true</EnableBatchedOperations>
            <SizeInBytes>0</SizeInBytes>
            <FilteringMessagesBeforePublishing>false</FilteringMessagesBeforePublishing>
            <IsAnonymousAccessible>false</IsAnonymousAccessible>
            <AuthorizationRules></AuthorizationRules>
            <Status>Active</Status>
            <CreatedAt>2018-05-04T20:59:02.86Z</CreatedAt>
            <UpdatedAt>2018-05-04T20:59:03Z</UpdatedAt>
            <SupportOrdering>true</SupportOrdering>
            <AutoDeleteOnIdle>P10675199DT2H48M5.4775807S</AutoDeleteOnIdle>
            <EnablePartitioning>false</EnablePartitioning>
            <IsExpress>false</IsExpress>
            <EntityAvailabilityStatus>Available</EntityAvailabilityStatus>
            <EnableSubscriptionPartitioning>false</EnableSubscriptionPartitioning>
            <EnableExpress>false</EnableExpress>
        </TopicDescription>`

	topicEntry1 = `
		<entry xmlns="http://www.w3.org/2005/Atom">
			<id>https://sbdjtest.servicebus.windows.net/foo</id>
			<title type="text">foo</title>
			<published>2018-05-02T20:54:59Z</published>
			<updated>2018-05-02T20:54:59Z</updated>
			<author>
				<name>sbdjtest</name>
			</author>
			<link rel="self" href="https://sbdjtest.servicebus.windows.net/foo"/>
			<content type="application/xml">` + topicDescription1 +
		`</content>
		</entry>`
)

func (suite *serviceBusSuite) TestTopicEntryUnmarshal() {
	var entry TopicEntry
	err := xml.Unmarshal([]byte(topicEntry1), &entry)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "https://sbdjtest.servicebus.windows.net/foo", entry.ID)
	assert.Equal(suite.T(), "foo", entry.Title)
	assert.Equal(suite.T(), "sbdjtest", *entry.Author.Name)
	assert.Equal(suite.T(), "https://sbdjtest.servicebus.windows.net/foo", entry.Link.HREF)
	assert.Equal(suite.T(), "P10675199DT2H48M5.4775807S", *entry.Content.TopicDescription.DefaultMessageTimeToLive)
	assert.NotNil(suite.T(), entry.Content)
}

func (suite *serviceBusSuite) TestTopicUnmarshal() {
	var entry Entry
	err := xml.Unmarshal([]byte(topicEntry1), &entry)
	assert.Nil(suite.T(), err)

	var td TopicDescription
	err = xml.Unmarshal([]byte(entry.Content.Body), &td)
	t := suite.T()
	assert.Nil(t, err)
	assert.Equal(t, int32(1024), *td.MaxSizeInMegabytes)
	assert.Equal(t, false, *td.RequiresDuplicateDetection)
	assert.Equal(t, "P10675199DT2H48M5.4775807S", *td.DefaultMessageTimeToLive)
	assert.Equal(t, "PT10M", *td.DuplicateDetectionHistoryTimeWindow)
	assert.Equal(t, true, *td.EnableBatchedOperations)
	assert.Equal(t, false, *td.FilteringMessagesBeforePublishing)
	assert.Equal(t, false, *td.EnableExpress)
	assert.Equal(t, int64(0), *td.SizeInBytes)
	assert.EqualValues(t, servicebus.EntityStatusActive, *td.Status)
}

func (suite *serviceBusSuite) TestTopicManagementWrites() {
	tests := map[string]func(context.Context, *testing.T, *TopicManager, string){
		"TestPutDefaultTopic": testPutTopic,
	}

	ns := suite.getNewSasInstance()
	tm := ns.NewTopicManager()
	for name, testFunc := range tests {
		suite.T().Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			name := suite.RandomName("gosb", 6)
			testFunc(ctx, t, tm, name)
			defer suite.cleanupQueue(name)
		})
	}
}

func (suite *serviceBusSuite) TestTopicManagementReads() {
	tests := map[string]func(context.Context, *testing.T, *TopicManager, []string){
		"TestGetTopic":   testGetTopic,
		"TestListTopics": testListTopics,
	}

	ns := suite.getNewSasInstance()
	tm := ns.NewTopicManager()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	names := []string{suite.randEntityName(), suite.randEntityName()}
	for _, name := range names {
		if _, err := tm.Put(ctx, name); err != nil {
			suite.T().Fatal(err)
		}
	}

	for name, testFunc := range tests {
		suite.T().Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			testFunc(ctx, t, tm, names)
		})
	}

	for _, name := range names {
		suite.cleanupTopic(name)
	}
}

func testGetTopic(ctx context.Context, t *testing.T, tm *TopicManager, names []string) {
	topic, err := tm.Get(ctx, names[0])
	assert.Nil(t, err)
	assert.NotNil(t, t)
	assert.Equal(t, topic.Entry.Title, names[0])
}

func testListTopics(ctx context.Context, t *testing.T, tm *TopicManager, names []string) {
	feed, err := tm.List(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, feed)
	queueNames := make([]string, len(feed.Entries))
	for idx, entry := range feed.Entries {
		queueNames[idx] = entry.Title
	}

	for _, name := range names {
		assert.Contains(t, queueNames, name)
	}
}

func testPutTopic(ctx context.Context, t *testing.T, tm *TopicManager, name string) {
	topic, err := tm.Put(ctx, name)
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	if assert.NotNil(t, topic) {
		assert.Equal(t, name, topic.Title)
	}
}

func (suite *serviceBusSuite) TestTopicManagement() {
	tests := map[string]func(context.Context, *testing.T, *TopicManager, string){
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

	ns := suite.getNewSasInstance()
	tm := ns.NewTopicManager()
	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			entityName := suite.randEntityName()
			defer suite.cleanupTopic(entityName)
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			testFunc(ctx, t, tm, entityName)

		}
		suite.T().Run(name, setupTestTeardown)
	}
}

func testDefaultTopic(ctx context.Context, t *testing.T, tm *TopicManager, name string) {
	topic := buildTopic(ctx, t, tm, name)
	assert.False(t, *topic.EnableExpress, "should not have Express enabled")
	assert.True(t, *topic.EnableBatchedOperations, "should not have batching enabled")
	assert.False(t, *topic.EnablePartitioning, "should not have partitioning enabled")
	assert.True(t, *topic.SupportOrdering, "should not support ordering")
	assert.False(t, *topic.RequiresDuplicateDetection, "should not require dup detection")
	assert.Equal(t, "P10675199DT2H48M5.4775807S", *topic.AutoDeleteOnIdle, "auto delete is not 10 minutes")
	assert.Equal(t, "PT10M", *topic.DuplicateDetectionHistoryTimeWindow, "dup detection is not 10 minutes")
	assert.EqualValues(t, servicebus.EntityStatusActive, *topic.Status, "topic status")
}

func testPartitionedTopic(ctx context.Context, t *testing.T, tm *TopicManager, name string) {
	topic := buildTopic(ctx, t, tm, name, TopicWithPartitioning())
	assert.True(t, *topic.EnablePartitioning)
}

func testSupportOrdering(ctx context.Context, t *testing.T, tm *TopicManager, name string) {
	topic := buildTopic(ctx, t, tm, name, TopicWithOrdering())
	assert.True(t, *topic.SupportOrdering)
}

func testTopicWithDuplicateDetection(ctx context.Context, t *testing.T, tm *TopicManager, name string) {
	window := time.Duration(20 * time.Minute)
	topic := buildTopic(ctx, t, tm, name, TopicWithDuplicateDetection(&window))
	assert.True(t, *topic.RequiresDuplicateDetection)
	assert.Equal(t, "PT20M", *topic.DuplicateDetectionHistoryTimeWindow)
}

func testTopicWithAutoDeleteOnIdle(ctx context.Context, t *testing.T, tm *TopicManager, name string) {
	window := time.Duration(20 * time.Minute)
	topic := buildTopic(ctx, t, tm, name, TopicWithAutoDeleteOnIdle(&window))
	assert.Equal(t, "PT20M", *topic.AutoDeleteOnIdle)
}

func testTopicWithBatchedOperations(ctx context.Context, t *testing.T, tm *TopicManager, name string) {
	topic := buildTopic(ctx, t, tm, name, TopicWithBatchedOperations())
	assert.True(t, *topic.EnableBatchedOperations)
}

func testTopicWithExpress(ctx context.Context, t *testing.T, tm *TopicManager, name string) {
	topic := buildTopic(ctx, t, tm, name, TopicWithExpress())
	assert.True(t, *topic.EnableExpress)
}

func testTopicWithMessageTimeToLive(ctx context.Context, t *testing.T, tm *TopicManager, name string) {
	window := time.Duration(20 * time.Minute)
	topic := buildTopic(ctx, t, tm, name, TopicWithMessageTimeToLive(&window))
	assert.Equal(t, "PT20M", *topic.DefaultMessageTimeToLive)
}

func testTopicWithMaxSizeInMegabytes(ctx context.Context, t *testing.T, tm *TopicManager, name string) {
	size := 2 * Megabytes
	topic := buildTopic(ctx, t, tm, name, TopicWithMaxSizeInMegabytes(size))
	assert.Equal(t, int32(size), *topic.MaxSizeInMegabytes)
}

func buildTopic(ctx context.Context, t *testing.T, tm *TopicManager, name string, opts ...TopicOption) *TopicDescription {
	te, err := tm.Put(ctx, name, opts...)
	if err != nil {
		t.Fatal(err)
	}
	return &te.Content.TopicDescription
}

func (suite *serviceBusSuite) TestTopic() {
	tests := map[string]func(context.Context, *testing.T, *Topic){
		"SimpleSend": testTopicSend,
	}

	ns := suite.getNewSasInstance()
	tm := ns.NewTopicManager()
	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			name := suite.randEntityName()
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			_, err := tm.Put(ctx, name)
			if err != nil {
				log.Fatalln(err)
			}

			topic := ns.NewTopic(name)
			defer func() {
				topic.Close(ctx)
				suite.cleanupTopic(name)
			}()
			testFunc(ctx, t, topic)
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func testTopicSend(ctx context.Context, t *testing.T, topic *Topic) {
	err := topic.Send(ctx, NewEventFromString("hello!"))
	assert.Nil(t, err)
}

func (suite *serviceBusSuite) cleanupTopic(name string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ns := suite.getNewSasInstance()
	tm := ns.NewTopicManager()
	err := tm.Delete(ctx, name)
	if err != nil {
		suite.T().Fatal(err)
	}
}
