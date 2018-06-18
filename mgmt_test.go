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
	"encoding/xml"
	"time"

	"github.com/Azure/azure-service-bus-go/atom"
	"github.com/Azure/go-autorest/autorest/date"
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
