package servicebus

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func (suite *serviceBusSuite) TestMessageIterator() {
	tests := map[string]func(ctx context.Context, t *testing.T, queue *Queue){
		"MultiplePages": testMessageIteratorMultiplePages,
		"NoMessages":    testMessageIteratorNoMessages,
		"Continue":      testMessageIteratorContinue,
		"StartHalfway":  testMessageIteratorStartHalfway,
		"LargePages":    testMessageIteratorLargePageSize,
		"PeekOne":       testMessageIteratorPeekOne,
	}

	ns := suite.getNewSasInstance()

	window := 30 * time.Second

	for name, testFunc := range tests {
		suite.T().Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()

			queueName := suite.randEntityName()
			cleanup := makeQueue(
				ctx,
				t,
				ns,
				queueName,
				QueueEntityWithDuplicateDetection(&window))
			defer cleanup()

			q, err := ns.NewQueue(queueName)
			defer func() {
				q.Close(context.Background())
			}()
			suite.NoError(err)

			testFunc(ctx, t, q)
		})
	}
}

func testMessageIteratorMultiplePages(ctx context.Context, t *testing.T, queue *Queue) {
	numMessages := rand.Intn(50) + 75

	expectedMessages := make([]*Message, numMessages)

	for i := 0; i < numMessages; i++ {
		msg := NewMessageFromString(fmt.Sprintf("message payload %d", i))
		require.NoError(t, queue.Send(ctx, msg))
		expectedMessages[i] = msg
	}

	subject, err := queue.Peek(ctx)
	require.NoError(t, err)

	for i := 0; i < numMessages; i++ {
		cursor, err := subject.Next(ctx)
		require.NoError(t, err)

		assert.Equal(t, string(expectedMessages[i].Data), string(cursor.Data))
	}
}

func testMessageIteratorNoMessages(ctx context.Context, t *testing.T, queue *Queue) {
	subject, err := queue.Peek(ctx)
	require.NoError(t, err)

	cursor, err := subject.Next(ctx)

	assert.EqualError(t, err, ErrNoMessages{}.Error())
	assert.Nil(t, cursor)
}

func testMessageIteratorContinue(ctx context.Context, t *testing.T, queue *Queue) {
	subject, err := queue.Peek(ctx)
	require.NoError(t, err)

	cursor, err := subject.Next(ctx)
	assert.EqualError(t, err, ErrNoMessages{}.Error())
	assert.Nil(t, cursor)

	numMessages := uint(rand.Intn(50) + 75)
	expectedMessages := make([]*Message, numMessages)

	for i := uint(0); i < numMessages; i++ {
		msg := NewMessageFromString(fmt.Sprintf("message payload 0x%x", i))
		require.NoError(t, queue.Send(ctx, msg))
		expectedMessages[i] = msg
	}

	matchCount := uint(0)
	for i := uint(0); i < numMessages; i++ {
		cursor, err := subject.Next(ctx)
		if _, ok := err.(ErrNoMessages); ok {
			i--
			t.Error(err)
			continue
		} else if err != nil {
			t.Error(err)
			return
		}

		want, got := string(expectedMessages[i].Data), string(cursor.Data)
		assert.Equal(t, want, got)
		if got == want {
			matchCount++
		}
	}

	logMessageMatches(t, matchCount, numMessages)
}

func testMessageIteratorStartHalfway(ctx context.Context, t *testing.T, queue *Queue) {
	numMessages := rand.Intn(50) + 75

	expectedMessages := make([]*Message, numMessages)

	for i := 0; i < numMessages; i++ {
		msg := NewMessage([]byte{byte(i)})
		require.NoError(t, queue.Send(ctx, msg))
		expectedMessages[i] = msg
	}

	startingPoint := rand.Intn(numMessages/10) + numMessages/2
	subject, err := queue.Peek(ctx, PeekFromSequenceNumber(int64(startingPoint)))
	require.NoError(t, err)

	matchCount := uint(0)
	for i := startingPoint; i < numMessages; i++ {
		cursor, err := subject.Next(ctx)
		require.NoError(t, err)

		got, want := int(cursor.Data[0]), int(expectedMessages[i].Data[0])
		assert.Equal(t, want, got)
		if got == want {
			matchCount++
		}
	}

	numExpected := uint(numMessages) - uint(startingPoint)
	logMessageMatches(t, matchCount, numExpected)
}

func testMessageIteratorLargePageSize(ctx context.Context, t *testing.T, queue *Queue) {
	const pageSize = 500
	const deciPageSize = pageSize / 10

	subject, err := queue.Peek(ctx, PeekWithPageSize(pageSize))
	require.NoError(t, err)

	newPayload := func(n int) string {
		return "abc123-" + strconv.Itoa(n)
	}

	for i := 0; i < deciPageSize; i++ {
		msg := NewMessageFromString(newPayload(i))
		require.NoError(t, queue.Send(ctx, msg))
	}

	matchCount := uint(0)
	for i := 0; i < deciPageSize; i++ {
		msg, err := subject.Next(ctx)
		require.NoError(t, err)
		got, want := string(msg.Data), newPayload(i)
		assert.Equal(t, want, got)
		if got == want {
			matchCount++
		}
	}

	logMessageMatches(t, matchCount, deciPageSize)

	_, err = subject.Next(ctx)
	assert.EqualError(t, err, ErrNoMessages{}.Error())

	for i := 0; i < pageSize; i++ {
		select {
		case <-ctx.Done():
			t.Error(ctx.Err())
			return
		default:
			// Intentionally Left Blank
		}
		msg := NewMessageFromString(newPayload(i))
		require.NoError(t, queue.Send(ctx, msg))
	}

	require.NoError(t, queue.Send(ctx, NewMessageFromString(newPayload(pageSize+1))))

	got := uint32(0)
	for {
		_, err := subject.Next(ctx)
		if _, ok := err.(ErrNoMessages); ok {
			break
		}

		got++
	}

	const want = pageSize + 1
	if got != want {
		t.Logf("got: %d want: %d", got, want)
		t.Fail()
	}
}

func testMessageIteratorPeekOne(ctx context.Context, t *testing.T, queue *Queue) {
	createPayload := func(x int) string {
		return fmt.Sprintf("payload-%d", x)
	}

	for i := 0; i < 5; i++ {
		msg := NewMessageFromString(createPayload(i))
		require.NoError(t, queue.Send(ctx, msg))
	}

	msg, err := queue.PeekOne(ctx)
	if err == nil {
		assert.Equal(t, createPayload(0), string(msg.Data))
	} else {
		t.Error(err)
	}

	const seqNum = 2
	msg, err = queue.PeekOne(ctx, PeekFromSequenceNumber(seqNum))
	if err == nil {
		assert.Equal(t, createPayload(seqNum), string(msg.Data))
	} else {
		t.Error(err)
	}
}

func logMessageMatches(t *testing.T, matches, total uint) {
	if testing.Verbose() || t.Failed() {
		t.Logf(
			"%d/%d (%0.2f%%) messages matched",
			matches,
			total,
			float32(matches)/float32(total)*100)
	}
}
