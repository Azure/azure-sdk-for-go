// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/mock"
	"github.com/Azure/go-amqp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUnitEventDataBatchConstants(t *testing.T) {
	smallBytes := [255]byte{0} // 'vbin8'
	largeBytes := [256]byte{0} // 'vbin32'

	require.Greater(t, calcActualSizeForPayload(largeBytes[:]), calcActualSizeForPayload(smallBytes[:]))

	require.EqualValues(t, calcActualSizeForPayload(smallBytes[:]), mustEncode(t, &amqp.Message{Data: [][]byte{smallBytes[:]}}))
	require.EqualValues(t, calcActualSizeForPayload(smallBytes[:])*2, mustEncode(t, &amqp.Message{Data: [][]byte{smallBytes[:], smallBytes[:]}}))

	require.EqualValues(t, calcActualSizeForPayload(largeBytes[:]), mustEncode(t, &amqp.Message{Data: [][]byte{largeBytes[:]}}))
	require.EqualValues(t, calcActualSizeForPayload(largeBytes[:])*2, mustEncode(t, &amqp.Message{Data: [][]byte{largeBytes[:], largeBytes[:]}}))

	require.EqualValues(t, calcActualSizeForPayload(largeBytes[:])+calcActualSizeForPayload(smallBytes[:]), mustEncode(t, &amqp.Message{Data: [][]byte{smallBytes[:], largeBytes[:]}}))
}

type eventBatchLinkForTest struct {
	amqpwrap.AMQPSenderCloser
	maxMessageSize uint64
}

func (l eventBatchLinkForTest) MaxMessageSize() uint64 {
	return l.maxMessageSize
}

func TestUnitEventDataBatchUnitTests(t *testing.T) {
	link := eventBatchLinkForTest{maxMessageSize: 10000}

	t.Run("default: uses link size", func(t *testing.T) {
		batch, err := newEventDataBatch(link, &EventDataBatchOptions{})
		require.NoError(t, err)
		require.NotNil(t, batch)
		require.Equal(t, link.MaxMessageSize(), batch.maxBytes)
		require.Nil(t, batch.partitionID)
		require.Nil(t, batch.partitionKey)

		batch, err = newEventDataBatch(link, nil)
		require.NoError(t, err)
		require.NotNil(t, batch)
		require.Equal(t, link.MaxMessageSize(), batch.maxBytes)
		require.Nil(t, batch.partitionID)
		require.Nil(t, batch.partitionKey)
	})

	t.Run("custom size", func(t *testing.T) {
		batch, err := newEventDataBatch(link, &EventDataBatchOptions{
			MaxBytes: 9,
		})
		require.NoError(t, err)
		require.NotNil(t, batch)
		require.Equal(t, uint64(9), batch.maxBytes)
	})

	t.Run("requested size is bigger than allowed size", func(t *testing.T) {
		batch, err := newEventDataBatch(link, &EventDataBatchOptions{MaxBytes: link.maxMessageSize + 1})
		require.EqualError(t, err, fmt.Sprintf("maximum message size for batch was set to %d bytes, which is larger than the maximum size allowed by link (%d)", link.maxMessageSize+1, link.MaxMessageSize()))
		require.Nil(t, batch)
	})

	t.Run("partition key", func(t *testing.T) {
		batch, err := newEventDataBatch(link, &EventDataBatchOptions{
			PartitionKey: to.Ptr("hello-partition-key"),
		})
		require.NoError(t, err)
		require.NotNil(t, batch)
		require.Equal(t, link.MaxMessageSize(), batch.maxBytes)
		require.Equal(t, "hello-partition-key", *batch.partitionKey)
		require.Nil(t, batch.partitionID)
	})

	t.Run("partition ID", func(t *testing.T) {
		batch, err := newEventDataBatch(link, &EventDataBatchOptions{
			PartitionID: to.Ptr("101"),
		})
		require.NoError(t, err)
		require.NotNil(t, batch)
		require.Equal(t, link.MaxMessageSize(), batch.maxBytes)
		require.Equal(t, "101", *batch.partitionID)
		require.Nil(t, batch.partitionKey)
	})

	as2k := [2048]byte{'A'}

	t.Run("sizeCalculationsAreCorrectVBin8", func(t *testing.T) {
		mb, err := newEventDataBatch(link, &EventDataBatchOptions{MaxBytes: 8000})
		require.NoError(t, err)

		err = mb.AddEventData(&EventData{
			Body: []byte("small body"),
			Properties: map[string]any{
				"small": "value",
			},
		}, nil)

		require.NoError(t, err)
		require.EqualValues(t, 1, mb.NumEvents())
		require.EqualValues(t, 172, mb.NumBytes())

		msg, err := mb.toAMQPMessage()
		require.NoError(t, err)
		actualBytes, err := msg.MarshalBinary()
		require.NoError(t, err)

		require.Equal(t, 172, len(actualBytes))
	})

	t.Run("sizeCalculationsAreCorrectVBin32", func(t *testing.T) {
		mb, err := newEventDataBatch(link, &EventDataBatchOptions{MaxBytes: 8000})
		require.NoError(t, err)

		err = mb.AddEventData(&EventData{
			Body: []byte("small body"),
			Properties: map[string]any{
				"hello":      "world",
				"anInt":      100,
				"aFLoat":     100.1,
				"lotsOfData": string(as2k[:]),
			},
		}, nil)

		require.NoError(t, err)
		require.EqualValues(t, 1, mb.NumEvents())
		require.EqualValues(t, 4357, mb.NumBytes())

		msg, err := mb.toAMQPMessage()
		require.NoError(t, err)

		actualBytes, err := msg.MarshalBinary()
		require.NoError(t, err)

		require.Equal(t, 4357, len(actualBytes))
	})

	// the first message gets special treatment since it gets used as the actual
	// batch message's envelope.
	t.Run("firstMessageTooLarge", func(t *testing.T) {
		mb, err := newEventDataBatch(link, &EventDataBatchOptions{MaxBytes: 1})
		require.NoError(t, err)

		err = mb.AddEventData(&EventData{
			Body: []byte("hello world"),
		}, nil)

		require.EqualError(t, err, ErrEventDataTooLarge.Error())

		require.EqualValues(t, 0, mb.NumBytes())
		require.EqualValues(t, 0, len(mb.marshaledMessages))
	})

	t.Run("addTooManyMessages", func(t *testing.T) {
		mb, err := newEventDataBatch(link, &EventDataBatchOptions{MaxBytes: 200})
		require.NoError(t, err)

		require.EqualValues(t, 0, mb.currentSize)
		err = mb.AddEventData(&EventData{
			Body: []byte("hello world"),
		}, nil)
		require.NoError(t, err)
		require.EqualValues(t, 121, mb.currentSize)

		sizeBefore := mb.NumBytes()
		countBefore := mb.NumEvents()

		err = mb.AddEventData(&EventData{
			Body: as2k[:],
		}, nil)
		require.EqualError(t, err, ErrEventDataTooLarge.Error())

		require.Equal(t, sizeBefore, mb.NumBytes(), "size is unchanged when a message fails to get added")
		require.Equal(t, countBefore, mb.NumEvents(), "count is unchanged when a message fails to get added")
	})

	t.Run("addConcurrently", func(t *testing.T) {
		mb, err := newEventDataBatch(link, &EventDataBatchOptions{MaxBytes: 10000})
		require.NoError(t, err)

		wg := sync.WaitGroup{}

		for i := byte(0); i < 100; i++ {
			wg.Add(1)
			go func(i byte) {
				defer wg.Done()

				err := mb.AddEventData(&EventData{
					Body: []byte{i},
				}, nil)

				require.NoError(t, err)
			}(i)
		}

		wg.Wait()
		require.EqualValues(t, 100, mb.NumEvents())
	})
}

func TestUnitEventDataBatchDontReuseOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mock.NewMockAMQPSenderCloser(ctrl)
	sender.EXPECT().MaxMessageSize().Return(uint64(200)).AnyTimes()

	t.Run("partitionID", func(t *testing.T) {
		pid := "6"
		batchForPartition, err := newEventDataBatch(sender, &EventDataBatchOptions{
			PartitionID: &pid,
		})
		require.NoError(t, err)

		require.Equal(t, "6", *batchForPartition.partitionID)
		pid = "7"
		require.Equal(t, "6", *batchForPartition.partitionID)
	})

	t.Run("partitionKey", func(t *testing.T) {
		pkey := "hello"

		batchForPartitionKey, err := newEventDataBatch(sender, &EventDataBatchOptions{
			PartitionKey: &pkey,
		})
		require.NoError(t, err)

		require.Equal(t, "hello", *batchForPartitionKey.partitionKey)
		pkey = "world"
		require.Equal(t, "hello", *batchForPartitionKey.partitionKey)
	})
}

func TestUnitEventDataBatchAlwaysHasProperties(t *testing.T) {
	ctrl := gomock.NewController(t)
	sender := mock.NewMockAMQPSenderCloser(ctrl)
	sender.EXPECT().MaxMessageSize().Return(uint64(200)).AnyTimes()

	batch, err := newEventDataBatch(sender, nil)
	require.NoError(t, err)

	t.Run("empty", func(t *testing.T) {
		amqpMsg, err := batch.toAMQPMessage()
		require.Nil(t, amqpMsg)
		require.EqualError(t, err, "batch is nil or empty")
	})

	t.Run("annotated message, empty", func(t *testing.T) {
		err = batch.AddAMQPAnnotatedMessage(&AMQPAnnotatedMessage{}, nil)
		require.NoError(t, err)

		msg, err := batch.toAMQPMessage()
		require.NoError(t, err)
		require.NotEmpty(t, msg, msg.Properties.MessageID)
	})

	t.Run("regular event, empty", func(t *testing.T) {
		err = batch.AddEventData(&EventData{}, nil)
		require.NoError(t, err)

		msg, err := batch.toAMQPMessage()
		require.NoError(t, err)
		require.NotEmpty(t, msg, msg.Properties.MessageID)
	})
}

func mustEncode(t *testing.T, msg *amqp.Message) int {
	bytes, err := msg.MarshalBinary()
	require.NoError(t, err)
	return len(bytes)
}
