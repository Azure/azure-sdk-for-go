// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"sync"
	"testing"

	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

func TestMessageBatchConstants(t *testing.T) {
	smallBytes := [255]byte{0} // 'vbin8'
	largeBytes := [256]byte{0} // 'vbin32'

	require.Greater(t, calcActualSizeForPayload(largeBytes[:]), calcActualSizeForPayload(smallBytes[:]))

	require.EqualValues(t, calcActualSizeForPayload(smallBytes[:]), mustEncode(t, &amqp.Message{Data: [][]byte{smallBytes[:]}}))
	require.EqualValues(t, calcActualSizeForPayload(smallBytes[:])*2, mustEncode(t, &amqp.Message{Data: [][]byte{smallBytes[:], smallBytes[:]}}))

	require.EqualValues(t, calcActualSizeForPayload(largeBytes[:]), mustEncode(t, &amqp.Message{Data: [][]byte{largeBytes[:]}}))
	require.EqualValues(t, calcActualSizeForPayload(largeBytes[:])*2, mustEncode(t, &amqp.Message{Data: [][]byte{largeBytes[:], largeBytes[:]}}))

	require.EqualValues(t, calcActualSizeForPayload(largeBytes[:])+calcActualSizeForPayload(smallBytes[:]), mustEncode(t, &amqp.Message{Data: [][]byte{smallBytes[:], largeBytes[:]}}))
}

func TestMessageBatchUnitTests(t *testing.T) {
	as2k := [2048]byte{'A'}

	t.Run("sizeCalculationsAreCorrectVBin8", func(t *testing.T) {
		mb := newMessageBatch(8000)

		err := mb.AddMessage(&Message{
			Body: []byte("small body"),
			ApplicationProperties: map[string]any{
				"small": "value",
			},
		}, nil)

		require.NoError(t, err)
		require.EqualValues(t, 1, mb.NumMessages())
		require.EqualValues(t, 196, mb.NumBytes())

		actualBytes, err := mb.toAMQPMessage().MarshalBinary()
		require.NoError(t, err)

		require.Equal(t, 196, len(actualBytes))
	})

	t.Run("sizeCalculationsAreCorrectVBin32", func(t *testing.T) {
		mb := newMessageBatch(8000)

		err := mb.AddMessage(&Message{
			Body: []byte("small body"),
			ApplicationProperties: map[string]any{
				"hello":      "world",
				"anInt":      100,
				"aFLoat":     100.1,
				"lotsOfData": string(as2k[:]),
			},
		}, nil)

		require.NoError(t, err)
		require.EqualValues(t, 1, mb.NumMessages())
		require.EqualValues(t, 4381, mb.NumBytes())

		actualBytes, err := mb.toAMQPMessage().MarshalBinary()
		require.NoError(t, err)

		require.Equal(t, 4381, len(actualBytes))
	})

	// the first message gets special treatment since it gets used as the actual
	// batch message's envelope.
	t.Run("firstMessageTooLarge", func(t *testing.T) {
		mb := newMessageBatch(1)

		err := mb.AddMessage(&Message{
			Body: []byte("hello world"),
		}, nil)

		require.EqualError(t, err, ErrMessageTooLarge.Error())

		require.EqualValues(t, 0, mb.NumBytes())
		require.EqualValues(t, 0, len(mb.marshaledMessages))
	})

	t.Run("addTooManyMessages", func(t *testing.T) {
		mb := newMessageBatch(200)

		require.EqualValues(t, 0, mb.currentSize)
		err := mb.AddMessage(&Message{
			Body: []byte("hello world"),
		}, nil)
		require.NoError(t, err)
		require.EqualValues(t, 145, mb.currentSize)

		sizeBefore := mb.NumBytes()
		countBefore := mb.NumMessages()

		err = mb.AddMessage(&Message{
			Body: as2k[:],
		}, nil)
		require.EqualError(t, err, ErrMessageTooLarge.Error())

		require.Equal(t, sizeBefore, mb.NumBytes(), "size is unchanged when a message fails to get added")
		require.Equal(t, countBefore, mb.NumMessages(), "count is unchanged when a message fails to get added")
	})

	t.Run("addConcurrently", func(t *testing.T) {
		mb := newMessageBatch(10000)

		wg := sync.WaitGroup{}

		for i := byte(0); i < 100; i++ {
			wg.Add(1)
			go func(i byte) {
				defer wg.Done()

				err := mb.AddMessage(&Message{
					Body: []byte{i},
				}, nil)

				require.NoError(t, err)
			}(i)
		}

		wg.Wait()
		require.EqualValues(t, 100, mb.NumMessages())
	})
}

func mustEncode(t *testing.T, msg *amqp.Message) int {
	bytes, err := msg.MarshalBinary()
	require.NoError(t, err)
	return len(bytes)
}
