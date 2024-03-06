// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

// internalBatchSender is an interface for an *azservicebus.Sender
type internalBatchSender interface {
	SendMessageBatch(ctx context.Context, batch internalBatch) error
	NewMessageBatch(ctx context.Context, options *azservicebus.MessageBatchOptions) (internalBatch, error)
}

type internalBatch interface {
	AddMessage(m *azservicebus.Message, options *azservicebus.AddMessageOptions) error
	NumMessages() int32
}

type senderWrapper struct {
	inner *TrackingSender
}

func (sw *senderWrapper) SendMessageBatch(ctx context.Context, batch internalBatch) error {
	return sw.inner.SendMessageBatch(ctx, batch.(*azservicebus.MessageBatch), nil)
}

func (sw *senderWrapper) NewMessageBatch(ctx context.Context, options *azservicebus.MessageBatchOptions) (internalBatch, error) {
	return sw.inner.NewMessageBatch(ctx, options)
}

func NewStreamingMessageBatch(ctx context.Context, sender internalBatchSender, expectedTotal int) (*StreamingMessageBatch, error) {
	batch, err := sender.NewMessageBatch(ctx, nil)

	if err != nil {
		return nil, err
	}

	return &StreamingMessageBatch{
		sender:        sender,
		currentBatch:  batch,
		expectedTotal: expectedTotal,
	}, nil
}

type StreamingMessageBatch struct {
	sender        internalBatchSender
	currentBatch  internalBatch
	expectedTotal int
	total         int
}

// Add appends to the current batch. If it's full it'll send it, allocate a new one.
func (sb *StreamingMessageBatch) Add(ctx context.Context, msg *azservicebus.Message, options *azservicebus.AddMessageOptions) error {
	err := sb.currentBatch.AddMessage(msg, options)

	if err == nil {
		return nil
	}

	if err != azservicebus.ErrMessageTooLarge {
		// must be a fatal error, just give up.
		return err
	}

	log.Printf("Sending message batch with %d messages. %d/%d messages sent so far.", sb.currentBatch.NumMessages(), sb.total, sb.expectedTotal)
	if err := sb.sender.SendMessageBatch(ctx, sb.currentBatch); err != nil {
		return err
	}

	sb.total += int(sb.currentBatch.NumMessages())

	// throttle a teeny bit.
	time.Sleep(time.Second)

	batch, err := sb.sender.NewMessageBatch(ctx, nil)

	if err != nil {
		return err
	}

	if err := batch.AddMessage(msg, nil); err != nil {
		// if we can't add this message here (ie, by itself) into the batch then
		// we'll just error out.
		return err
	}

	sb.currentBatch = batch
	return nil
}

// Close sends any messages currently held in our batch.
func (sb *StreamingMessageBatch) Close(ctx context.Context) error {
	if sb.currentBatch.NumMessages() == 0 {
		return nil
	}

	log.Printf("Sending final message batch")
	if err := sb.sender.SendMessageBatch(ctx, sb.currentBatch); err != nil {
		return err
	}

	return nil
}
