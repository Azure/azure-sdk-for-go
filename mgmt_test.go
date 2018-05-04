package servicebus

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
