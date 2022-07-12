// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	eventhub "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
)

func TestNewEventBatch(t *testing.T) {
	eb := eventhub.NewEventBatch("eventId", nil)
	assert.Equal(t, eventhub.DefaultMaxMessageSizeInBytes, eb.MaxSize)
}

func TestEventBatch_AddOneMessage(t *testing.T) {
	eb := eventhub.NewEventBatch("eventId", nil)
	event := eventhub.NewEventFromString("Foo")
	ok, err := eb.Add(event)
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestEventBatch_AddManyMessages(t *testing.T) {
	eb := eventhub.NewEventBatch("eventId", nil)
	wrapperSize := eb.Size()
	event := eventhub.NewEventFromString("Foo")
	ok, err := eb.Add(event)
	assert.True(t, ok)
	assert.NoError(t, err)

	msgSize := eb.Size() - wrapperSize

	limit := ((int(eb.MaxSize) - 100) / msgSize) - 1
	for i := 0; i < limit; i++ {
		ok, err := eb.Add(event)
		assert.True(t, ok)
		assert.NoError(t, err)
	}

	ok, err = eb.Add(event)
	assert.False(t, ok)
	assert.NoError(t, err)
}

func TestEventBatch_Clear(t *testing.T) {
	eb := eventhub.NewEventBatch("eventId", nil)
	ok, err := eb.Add(eventhub.NewEventFromString("Foo"))
	assert.True(t, ok)
	assert.NoError(t, err)
	assert.Equal(t, 163, eb.Size())

	eb.Clear()
	assert.Equal(t, 100, eb.Size())
}

func TestHugeBatches(t *testing.T) {
	data := make([]byte, 500)
	events := make([]*eventhub.Event, 0)

	for i := 0; i < 100; i++ {
		// 100 / 4 * 50000 = 1250000 bytes per partition
		partitionKey := strconv.Itoa(i % 4)
		evt := &eventhub.Event{
			Data:         data,
			PartitionKey: &partitionKey,
		}

		events = append(events, evt)
	}

	opts := &eventhub.BatchOptions{
		MaxSize: 10000,
	}
	iter := eventhub.NewEventBatchIterator(events...)
	iterCount := 0

	for !iter.Done() {
		_, err := iter.Next("batchId", opts)
		assert.NoError(t, err)

		iterCount++

		if iterCount > 101 {
			assert.Fail(t, "Too much iteration")
		}
	}

	assert.Greater(t, iterCount, 5)
}

func TestOneHugeEvent(t *testing.T) {
	data := make([]byte, 1100)
	events := []*eventhub.Event{
		{
			Data: data,
		},
	}
	opts := &eventhub.BatchOptions{
		MaxSize: 1000,
	}
	iter := eventhub.NewEventBatchIterator(events...)

	for !iter.Done() {
		_, err := iter.Next("batchId", opts)
		assert.Equal(t, err, eventhub.ErrMessageIsTooBig)
	}
}
