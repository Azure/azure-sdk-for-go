// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/stretchr/testify/require"
)

type fakeBatch struct {
	max      int
	messages []*azservicebus.Message
}

func (fb *fakeBatch) AddMessage(msg *azservicebus.Message) error {
	if len(fb.messages) == fb.max {
		return azservicebus.ErrMessageTooLarge
	}

	fb.messages = append(fb.messages, msg)
	return nil
}

func (fb *fakeBatch) NumMessages() int32 {
	return int32(len(fb.messages))
}

type fakeSender struct {
	sent         []*fakeBatch
	maxBatchSize int
}

func (s *fakeSender) SendMessageBatch(ctx context.Context, batch internalBatch) error {
	s.sent = append(s.sent, batch.(*fakeBatch))
	return nil
}

func (s *fakeSender) NewMessageBatch(ctx context.Context, options *azservicebus.MessageBatchOptions) (internalBatch, error) {
	return &fakeBatch{
		max: s.maxBatchSize,
	}, nil
}

func TestStreamingBatch(t *testing.T) {
	ctx := context.Background()

	sender := &fakeSender{
		maxBatchSize: 2,
	}
	stats := NewStats("")

	streamingBatch, err := NewStreamingMessageBatch(ctx, sender, stats)
	require.NoError(t, err)
	require.NotNil(t, streamingBatch.currentBatch, "first batch is auto-created right when we create the first batch")

	currentBatch := func() *fakeBatch {
		return streamingBatch.currentBatch.(*fakeBatch)
	}

	for i := 0; i < currentBatch().max; i++ {
		err = streamingBatch.Add(ctx, &azservicebus.Message{
			Body: []byte(fmt.Sprintf("%d", i)),
		})
		require.NoError(t, err)
		require.Empty(t, sender.sent, "Nothing will be sent yet, since the batch is not yet full")
	}

	// this next message won't fix into the batch, so now we'll:
	// 1. Send the current batch
	// 2. Create a new batch
	// 3. Add this new message to the new batch
	err = streamingBatch.Add(ctx, &azservicebus.Message{
		Body: []byte("last"),
	})
	require.NoError(t, err)

	// check what's been sent.
	require.EqualValues(t, 1, len(sender.sent))
	require.EqualValues(t, []string{"0", "1"}, getSortedBodies(sender.sent[0].messages))

	// flush the streaming batch and it'll send whatever remains
	sender.sent = nil
	require.NoError(t, streamingBatch.Close(context.Background()))
	require.EqualValues(t, 1, len(sender.sent))
	require.EqualValues(t, []string{"last"}, getSortedBodies(sender.sent[0].messages))
}

func TestStreamingBatchTooLargeToFitByItself(t *testing.T) {
	ctx := context.Background()

	sender := &fakeSender{
		// purposefully make it impossible to put any messages into a batch
		maxBatchSize: 0,
	}
	stats := NewStats("")

	streamingBatch, err := NewStreamingMessageBatch(ctx, sender, stats)
	require.NoError(t, err)

	// this next message won't fix into the batch, so now we'll:
	// 1. Send the current batch
	// 2. Create a new batch
	// 3. Add this new message to the new batch
	err = streamingBatch.Add(ctx, &azservicebus.Message{
		Body: []byte("last"),
	})

	// This is the only fatal case for the streaming batch - when we have a message that literally can't be sent.
	// There's no good strategy here, the user has to make a decision about whether to keep the message, etc...
	require.EqualError(t, err, azservicebus.ErrMessageTooLarge.Error())

	// and _nothing_ is sent at this point (they could still find another smaller message to add into the batch)
	require.Empty(t, 0, len(sender.sent))
}

func getSortedBodies(messages []*azservicebus.Message) []string {
	var bodies []string

	for _, m := range messages {
		bodies = append(bodies, string(m.Body))
	}

	sort.Strings(bodies)
	return bodies
}
