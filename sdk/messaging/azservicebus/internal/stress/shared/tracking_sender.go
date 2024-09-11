// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

func NewTrackingSender(tc *TelemetryClientWrapper, client *azservicebus.Client, queueOrTopic string, options *azservicebus.NewSenderOptions) (*TrackingSender, error) {
	tmpSender, err := client.NewSender(queueOrTopic, options)

	if err != nil {
		return nil, err
	}

	return &TrackingSender{tc, tmpSender}, nil
}

// TrackingSender reports metrics and errors automatically for its methods.
type TrackingSender struct {
	tc *TelemetryClientWrapper
	s  *azservicebus.Sender
}

func (ts *TrackingSender) NewMessageBatch(ctx context.Context, options *azservicebus.MessageBatchOptions) (*azservicebus.MessageBatch, error) {
	batch, err := ts.s.NewMessageBatch(ctx, options)

	if err != nil {
		TrackError(ctx, ts.tc, fmt.Errorf("error during NewMessageBatch: %w", err))
		return batch, err
	}

	return batch, err
}

func (ts *TrackingSender) SendMessage(ctx context.Context, message *azservicebus.Message, options *azservicebus.SendMessageOptions) error {
	end := TrackDuration(ctx, ts.tc, MetricMessagesSent)
	defer end(map[string]string{
		AttrMessageCount: "1",
	})

	if err := ts.s.SendMessage(ctx, message, options); err != nil {
		TrackError(ctx, ts.tc, fmt.Errorf("error during SendMessage: %w", err))
		return err
	}

	return nil
}

func (ts *TrackingSender) SendMessageBatch(ctx context.Context, batch *azservicebus.MessageBatch, options *azservicebus.SendMessageBatchOptions) error {
	end := TrackDuration(ctx, ts.tc, MetricMessagesSent)
	defer end(map[string]string{
		AttrMessageCount: fmt.Sprintf("%d", batch.NumMessages()),
	})

	if err := ts.s.SendMessageBatch(ctx, batch, options); err != nil {
		TrackError(ctx, ts.tc, fmt.Errorf("error during SendMessageBatch: %w", err))
		return err
	}

	return nil
}

func (ts *TrackingSender) Close(ctx context.Context) error {
	end := TrackDuration(ctx, ts.tc, MetricCloseDuration)
	defer end(nil)

	if err := ts.s.Close(ctx); err != nil {
		TrackError(ctx, ts.tc, fmt.Errorf("error during Close: %w", err))
		return err
	}

	return nil
}
