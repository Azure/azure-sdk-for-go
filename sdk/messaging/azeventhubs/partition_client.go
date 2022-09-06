// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
)

// DefaultConsumerGroup is the name of the default consumer group in the Event Hubs service.
const DefaultConsumerGroup = "$Default"

// StartPosition indicates the position to start receiving events within a partition.
// The default position is Latest.
type StartPosition struct {
	// Offset will start the consumer after the specified offset. Can be exclusive
	// or inclusive, based on the Inclusive property.
	// NOTE: offsets are not stable values, and might refer to different events over time
	// as the Event Hub events reach their age limit and are discarded.
	Offset *int64

	// SequenceNumber will start the consumer after the specified sequence number. Can be exclusive
	// or inclusive, based on the Inclusive property.
	SequenceNumber *int64

	// EnqueuedTime will start the consumer before events that were enqueued on or after EnqueuedTime.
	// Can be exclusive or inclusive, based on the Inclusive property.
	EnqueuedTime *time.Time

	// Inclusive configures whether the events directly at Offset, SequenceNumber or EnqueuedTime will be included (true)
	// or excluded (false).
	Inclusive bool

	// Earliest will start the consumer at the earliest event.
	Earliest *bool

	// Latest will start the consumer after the last event.
	Latest *bool
}

// PartitionClient is used to receive events from an Event Hub partition.
type PartitionClient struct {
	retryOptions  RetryOptions
	eventHub      string
	consumerGroup string
	partitionID   string
	ownerLevel    *int64

	offsetExpression string

	links internal.LinksForPartitionClient[amqpwrap.AMQPReceiverCloser]
}

// ReceiveEventsOptions contains optional parameters for the ReceiveEvents function
type ReceiveEventsOptions struct {
	// For future expansion
}

// ReceiveEvents receives events until 'count' events have been received or the context has
// expired or been cancelled.
func (cc *PartitionClient) ReceiveEvents(ctx context.Context, count int, options *ReceiveEventsOptions) ([]*ReceivedEventData, error) {
	var events []*ReceivedEventData

	err := cc.links.Retry(ctx, EventConsumer, "ReceiveEvents", cc.partitionID, cc.retryOptions, func(ctx context.Context, lwid internal.LinkWithID[amqpwrap.AMQPReceiverCloser]) error {
		events = nil

		outstandingCredits := lwid.Link.Credits()

		if count > int(outstandingCredits) {
			newCredits := uint32(count) - outstandingCredits

			log.Writef(EventConsumer, "Have %d outstanding credit, only issuing %d credits", outstandingCredits, newCredits)

			if err := lwid.Link.IssueCredit(newCredits); err != nil {
				return err
			}
		}

		for {
			amqpMessage, err := lwid.Link.Receive(ctx)

			if err != nil {
				prefetched := getAllPrefetched(lwid.Link, count-len(events))

				for _, amqpMsg := range prefetched {
					re, err := newReceivedEventData(amqpMsg)

					if err != nil {
						return err
					}

					events = append(events, re)
				}

				// this lets cancel errors just return
				return err
			}

			receivedEvent, err := newReceivedEventData(amqpMessage)

			if err != nil {
				return err
			}

			events = append(events, receivedEvent)

			if len(events) == count {
				return nil
			}
		}
	})

	if err != nil && len(events) == 0 {
		// TODO: if we get a "partition ownership lost" we need to think about whether that's retryable.
		return nil, internal.TransformError(err)
	}

	cc.offsetExpression = formatOffsetExpressionForSequence(">", events[len(events)-1].SequenceNumber)
	return events, nil
}

// Close closes the consumer's link and the underlying AMQP connection.
func (cc *PartitionClient) Close(ctx context.Context) error {
	if cc.links != nil {
		return cc.links.Close(ctx)
	}

	return nil
}

func (s *PartitionClient) getEntityPath(partitionID string) string {
	return fmt.Sprintf("%s/ConsumerGroups/%s/Partitions/%s", s.eventHub, s.consumerGroup, partitionID)
}

const defaultLinkRxBuffer = 2048

func (s *PartitionClient) newEventHubConsumerLink(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (internal.AMQPReceiverCloser, error) {
	var props map[string]interface{}

	if s.ownerLevel != nil {
		props = map[string]interface{}{
			"com.microsoft:epoch": *s.ownerLevel,
		}
	}

	receiver, err := session.NewReceiver(ctx, entityPath, &amqp.ReceiverOptions{
		SettlementMode: to.Ptr(amqp.ModeFirst),
		ManualCredits:  true,
		Credit:         defaultLinkRxBuffer,
		Filters: []amqp.LinkFilter{
			amqp.LinkFilterSelector(s.offsetExpression),
		},
		Properties: props,
	})

	if err != nil {
		return nil, err
	}

	return receiver, nil
}

func (pc *PartitionClient) init(ctx context.Context) error {
	return pc.links.Retry(ctx, EventConsumer, "Init", pc.partitionID, pc.retryOptions, func(ctx context.Context, lwid internal.LinkWithID[amqpwrap.AMQPReceiverCloser]) error {
		return nil
	})
}

type partitionClientArgs struct {
	namespace *internal.Namespace

	eventHub    string
	partitionID string

	consumerGroup string

	retryOptions RetryOptions
}

func newPartitionClient(args partitionClientArgs, options *NewPartitionClientOptions) (*PartitionClient, error) {
	if options == nil {
		options = &NewPartitionClientOptions{}
	}

	offsetExpr, err := getOffsetExpression(options.StartPosition)

	if err != nil {
		return nil, err
	}

	client := &PartitionClient{
		eventHub:         args.eventHub,
		partitionID:      args.partitionID,
		ownerLevel:       options.OwnerLevel,
		consumerGroup:    args.consumerGroup,
		offsetExpression: offsetExpr,
		retryOptions:     args.retryOptions,
	}

	client.links = internal.NewLinks(args.namespace, fmt.Sprintf("%s/$management", client.eventHub), client.getEntityPath, client.newEventHubConsumerLink)

	return client, nil
}

func getAllPrefetched(receiver amqpwrap.AMQPReceiver, max int) []*amqp.Message {
	var messages []*amqp.Message

	for i := 0; i < max; i++ {
		msg := receiver.Prefetched()

		if msg == nil {
			break
		}

		messages = append(messages, msg)
	}

	return messages
}

func getOffsetExpression(startPosition StartPosition) (string, error) {
	lt := ">"

	if startPosition.Inclusive {
		lt = ">="
	}

	var errMultipleFieldsSet = errors.New("only a single start point can be set: Earliest, EnqueuedTime, Latest, Offset, or SequenceNumber")

	offsetExpr := ""

	if startPosition.EnqueuedTime != nil {
		// time-based, non-inclusive
		offsetExpr = fmt.Sprintf("amqp.annotation.x-opt-enqueued-time %s '%d'", lt, startPosition.EnqueuedTime.UnixMilli())
	}

	if startPosition.Offset != nil {
		// offset-based, non-inclusive
		// ex: amqp.annotation.x-opt-enqueued-time %s '165805323000'
		if offsetExpr != "" {
			return "", errMultipleFieldsSet
		}

		offsetExpr = fmt.Sprintf("amqp.annotation.x-opt-offset %s '%d'", lt, *startPosition.Offset)
	}

	if startPosition.Latest != nil && *startPosition.Latest {
		if offsetExpr != "" {
			return "", errMultipleFieldsSet
		}

		offsetExpr = "amqp.annotation.x-opt-offset > '@latest'"
	}

	if startPosition.SequenceNumber != nil {
		if offsetExpr != "" {
			return "", errMultipleFieldsSet
		}

		offsetExpr = formatOffsetExpressionForSequence(lt, *startPosition.SequenceNumber)
	}

	if startPosition.Earliest != nil && *startPosition.Earliest {
		if offsetExpr != "" {
			return "", errMultipleFieldsSet
		}

		return "amqp.annotation.x-opt-offset > '-1'", nil
	}

	if offsetExpr != "" {
		return offsetExpr, nil
	}

	// default to the start
	return "amqp.annotation.x-opt-offset > '@latest'", nil
}

func formatOffsetExpressionForSequence(op string, sequenceNumber int64) string {
	return fmt.Sprintf("amqp.annotation.x-opt-sequence-number %s '%d'", op, sequenceNumber)
}
