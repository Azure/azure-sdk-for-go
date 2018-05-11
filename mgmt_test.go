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

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/stretchr/testify/assert"
)

func (suite *serviceBusSuite) TestFeedUnmarshal() {
	var feed Feed
	err := xml.Unmarshal([]byte(feedOfQueues), &feed)
	assert.Nil(suite.T(), err)
	updated, err := date.ParseTime(time.RFC3339, "2018-05-03T00:21:15Z")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "https://sbdjtest.servicebus.windows.net/$Resources/Queues", feed.ID)
	assert.Equal(suite.T(), "Queues", feed.Title)
	assert.WithinDuration(suite.T(), updated, feed.Updated.ToTime(), 100*time.Millisecond)
	assert.Equal(suite.T(), updated, (*feed.Updated).ToTime())
	if assert.Len(suite.T(), feed.Entries, 2) {
		assert.NotNil(suite.T(), feed.Entries[0].Content)
	}

}

func (suite *serviceBusSuite) TestEntryUnmarshal() {
	var entry Entry
	err := xml.Unmarshal([]byte(queueEntry1), &entry)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), "https://sbdjtest.servicebus.windows.net/foo", entry.ID)
	assert.Equal(suite.T(), "foo", entry.Title)
	assert.Equal(suite.T(), "sbdjtest", *entry.Author.Name)
	assert.Equal(suite.T(), "https://sbdjtest.servicebus.windows.net/foo", entry.Link.HREF)
	for _, item := range []string{
		`<QueueDescription`,
		"<LockDuration>PT1M</LockDuration>",
		"<MessageCount>0</MessageCount>",
	} {
		assert.Contains(suite.T(), entry.Content.Body, item)
	}
}
