// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"encoding/xml"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/go-autorest/autorest/date"
)

const (
	queueDescription1 = `
		<QueueDescription 
            xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect" 
            xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
            <LockDuration>PT1M</LockDuration>
            <MaxSizeInMegabytes>1024</MaxSizeInMegabytes>
            <RequiresDuplicateDetection>false</RequiresDuplicateDetection>
            <RequiresSession>false</RequiresSession>
            <DefaultMessageTimeToLive>P14D</DefaultMessageTimeToLive>
            <DeadLetteringOnMessageExpiration>false</DeadLetteringOnMessageExpiration>
            <DuplicateDetectionHistoryTimeWindow>PT10M</DuplicateDetectionHistoryTimeWindow>
            <MaxDeliveryCount>10</MaxDeliveryCount>
            <EnableBatchedOperations>true</EnableBatchedOperations>
            <SizeInBytes>0</SizeInBytes>
            <MessageCount>0</MessageCount>
            <IsAnonymousAccessible>false</IsAnonymousAccessible>
            <Status>Active</Status>
            <CreatedAt>2018-05-04T16:38:27.913Z</CreatedAt>
            <UpdatedAt>2018-05-04T16:38:41.897Z</UpdatedAt>
            <SupportOrdering>true</SupportOrdering>
            <AutoDeleteOnIdle>P14D</AutoDeleteOnIdle>
            <EnablePartitioning>false</EnablePartitioning>
            <EntityAvailabilityStatus>Available</EntityAvailabilityStatus>
            <EnableExpress>false</EnableExpress>
        </QueueDescription>`

	queueDescription2 = `
		<QueueDescription 
            xmlns="http://schemas.microsoft.com/netservices/2010/10/servicebus/connect" 
            xmlns:i="http://www.w3.org/2001/XMLSchema-instance">
            <LockDuration>PT2M</LockDuration>
            <MaxSizeInMegabytes>2048</MaxSizeInMegabytes>
            <RequiresDuplicateDetection>false</RequiresDuplicateDetection>
            <RequiresSession>false</RequiresSession>
            <DefaultMessageTimeToLive>P14D</DefaultMessageTimeToLive>
            <DeadLetteringOnMessageExpiration>true</DeadLetteringOnMessageExpiration>
            <DuplicateDetectionHistoryTimeWindow>PT20M</DuplicateDetectionHistoryTimeWindow>
            <MaxDeliveryCount>100</MaxDeliveryCount>
            <EnableBatchedOperations>true</EnableBatchedOperations>
            <SizeInBytes>256</SizeInBytes>
            <MessageCount>23</MessageCount>
            <IsAnonymousAccessible>false</IsAnonymousAccessible>
            <Status>Active</Status>
            <CreatedAt>2018-05-04T16:38:27.913Z</CreatedAt>
            <UpdatedAt>2018-05-04T16:38:41.897Z</UpdatedAt>
            <SupportOrdering>true</SupportOrdering>
            <AutoDeleteOnIdle>P14D</AutoDeleteOnIdle>
            <EnablePartitioning>true</EnablePartitioning>
            <EntityAvailabilityStatus>Available</EntityAvailabilityStatus>
            <EnableExpress>false</EnableExpress>
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

func (suite *serviceBusSuite) TestFeedUnmarshal() {
	var feed atom.Feed
	err := xml.Unmarshal([]byte(feedOfQueues), &feed)
	suite.Nil(err)
	updated, err := date.ParseTime(time.RFC3339, "2018-05-03T00:21:15Z")
	suite.Nil(err)
	suite.Equal("https://sbdjtest.servicebus.windows.net/$Resources/Queues", feed.ID)
	suite.Equal("Queues", feed.Title)
	suite.WithinDuration(updated, feed.Updated.ToTime(), 100*time.Millisecond)
	suite.Equal(updated, (*feed.Updated).ToTime())
	if suite.Len(feed.Entries, 2) {
		suite.NotNil(feed.Entries[0].Content)
	}

}

func (suite *serviceBusSuite) TestEntryUnmarshal() {
	var entry atom.Entry
	err := xml.Unmarshal([]byte(queueEntry1), &entry)
	suite.Nil(err)
	suite.Equal("https://sbdjtest.servicebus.windows.net/foo", entry.ID)
	suite.Equal("foo", entry.Title)
	suite.Equal("sbdjtest", *entry.Author.Name)
	suite.Equal("https://sbdjtest.servicebus.windows.net/foo", entry.Link.HREF)
	for _, item := range []string{
		`<QueueDescription`,
		"<LockDuration>PT1M</LockDuration>",
		"<MessageCount>0</MessageCount>",
	} {
		suite.Contains(entry.Content.Body, item)
	}
}
