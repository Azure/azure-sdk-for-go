package servicebus

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func (suite *serviceBusSuite) TestMessageIterator() {
	tests := map[string]func(ctx context.Context, t *testing.T, queue *Queue){
		"MultiplePages": testMessageIteratorMultiplePages,
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
