package servicebus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMessageBatch(t *testing.T) {
	mb := NewMessageBatch(StandardMaxMessageSizeInBytes, "messageID", &BatchOptions{})
	assert.Equal(t, StandardMaxMessageSizeInBytes, mb.MaxSize)
}

func TestMessageBatch_AddOneMessage(t *testing.T) {
	mb := NewMessageBatch(StandardMaxMessageSizeInBytes, "messageID", &BatchOptions{})
	msg := NewMessageFromString("Foo")
	ok, err := mb.Add(msg)
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestMessageBatch_AddManyMessages(t *testing.T) {
	bOpts := &BatchOptions{
		SessionID: ptrString("foobarbazzbuzz"),
	}
	mb := NewMessageBatch(StandardMaxMessageSizeInBytes, "messageID", bOpts)
	wrapperSize := mb.Size()
	msg := NewMessageFromString("Foo")
	ok, err := mb.Add(msg)
	assert.True(t, ok)
	assert.NoError(t, err)
	msgSize := mb.Size() - wrapperSize

	limit := ((int(mb.MaxSize) - 100) / msgSize) - 1
	for i := 0; i < limit; i++ {
		ok, err := mb.Add(msg)
		assert.True(t, ok)
		assert.NoError(t, err)
	}

	ok, err = mb.Add(msg)
	assert.False(t, ok)
	assert.NoError(t, err)
}

func TestMessageBatch_Clear(t *testing.T) {
	mb := NewMessageBatch(StandardMaxMessageSizeInBytes, "messageID", &BatchOptions{})
	ok, err := mb.Add(NewMessageFromString("foo"))
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, 175, mb.Size())

	mb.Clear()
	assert.Equal(t, 100, mb.Size())
}
