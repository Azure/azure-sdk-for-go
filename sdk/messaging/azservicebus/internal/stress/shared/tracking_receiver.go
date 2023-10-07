// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
	"github.com/microsoft/ApplicationInsights-Go/appinsights"
)

func NewTrackingReceiverForQueue(tc appinsights.TelemetryClient, client *azservicebus.Client, queueName string, options *azservicebus.ReceiverOptions) (*TrackingReceiver, error) {
	tmpReceiver, err := client.NewReceiverForQueue(queueName, options)

	if err != nil {
		return nil, err
	}

	return &TrackingReceiver{r: tmpReceiver, tc: tc}, nil
}

func NewTrackingReceiverForSubscription(tc appinsights.TelemetryClient, client *azservicebus.Client, topicName string, subscriptionName string, options *azservicebus.ReceiverOptions) (*TrackingReceiver, error) {
	tmpReceiver, err := client.NewReceiverForSubscription(topicName, subscriptionName, options)

	if err != nil {
		return nil, err
	}

	return &TrackingReceiver{r: tmpReceiver, tc: tc}, nil
}

// TrackingReceiver reports metrics and errors automatically for its methods.
type TrackingReceiver struct {
	r  *azservicebus.Receiver
	tc appinsights.TelemetryClient
}

func (tr *TrackingReceiver) ReceiveMessages(ctx context.Context, maxMessages int, options *azservicebus.ReceiveMessagesOptions) ([]*azservicebus.ReceivedMessage, error) {
	// TODO: there's no "receiver duration" equivalent metric?
	// TODO: it's a little tricky doing this above the library since you don't know how much time was taken just waiting
	//       for a message to arrive vs actually processing/receiving the messages.
	messages, err := tr.r.ReceiveMessages(ctx, maxMessages, options)

	if err != nil {
		TrackError(ctx, tr.tc, fmt.Errorf("error receiving events: %w", err))
		return messages, err
	}

	TrackMetric(ctx, tr.tc, MetricMessageReceived, float64(len(messages)), nil)
	return messages, err
}

func (tr *TrackingReceiver) CompleteMessage(ctx context.Context, message *azservicebus.ReceivedMessage, options *azservicebus.CompleteMessageOptions) error {
	end := TrackDuration(ctx, tr.tc, MetricSettlementRequestDuration)
	defer end(map[string]string{AttrAMQPDeliveryState: "accepted"})

	if err := tr.r.CompleteMessage(ctx, message, options); err != nil {
		tr.tc.TrackException(fmt.Errorf("error completing message: %w", err))
		return err
	}

	return nil
}

func (tr *TrackingReceiver) AbandonMessage(ctx context.Context, message *azservicebus.ReceivedMessage, options *azservicebus.AbandonMessageOptions) error {
	end := TrackDuration(ctx, tr.tc, MetricSettlementRequestDuration)
	defer end(map[string]string{AttrAMQPDeliveryState: "rejected"})

	if err := tr.r.AbandonMessage(ctx, message, options); err != nil {
		tr.tc.TrackException(fmt.Errorf("error completing message: %w", err))
		return err
	}

	return nil
}

func (tr *TrackingReceiver) PeekMessages(ctx context.Context, maxMessageCount int, options *azservicebus.PeekMessagesOptions) ([]*azservicebus.ReceivedMessage, error) {
	peeked, err := tr.r.PeekMessages(ctx, maxMessageCount, options)

	if err != nil {
		tr.tc.TrackException(fmt.Errorf("error peeking messages: %w", err))
		return peeked, err
	}

	TrackMetric(ctx, tr.tc, MetricMessagePeeked, float64(len(peeked)), nil)
	return peeked, err
}

func (tr *TrackingReceiver) RenewMessageLock(ctx context.Context, msg *azservicebus.ReceivedMessage, options *azservicebus.RenewMessageLockOptions) error {
	end := TrackDuration(ctx, tr.tc, MetricLockRenew)
	defer end(nil) // TODO: does Liudmila insert a message ID here as baggage?

	if err := tr.r.RenewMessageLock(ctx, msg, options); err != nil {
		TrackError(ctx, tr.tc, err)
		return err
	}

	return nil
}

func (tr *TrackingReceiver) Close(ctx context.Context) error {
	end := TrackDuration(ctx, tr.tc, MetricCloseDuration)
	defer end(nil)

	if err := tr.r.Close(ctx); err != nil {
		TrackError(ctx, tr.tc, fmt.Errorf("error during close: %w", err))
		return err
	}

	return nil
}

func isCancelError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
