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

const defaultPrefetchSize = uint32(300)
const defaultMaxCreditSize = uint32(2048)

// StartPosition indicates the position to start receiving events within a partition.
// The default position is Latest.
//
// You can set this in the options for ConsumerClient.
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
//
// This type is instantiated from the [ConsumerClient] type, using [ConsumerClient.NewPartitionClient].
type PartitionClient struct {
	retryOptions  RetryOptions
	eventHub      string
	consumerGroup string
	partitionID   string
	ownerLevel    *int64
	prefetch      int32

	offsetExpression string

	links internal.LinksForPartitionClient[amqpwrap.AMQPReceiverCloser]
}

// ReceiveEventsOptions contains optional parameters for the ReceiveEvents function
type ReceiveEventsOptions struct {
	// For future expansion
}

// ReceiveEvents receives events until 'count' events have been received or the context has
// expired or been cancelled.
//
// If your ReceiveEvents call appears to be stuck there are some common causes:
//
//  1. The PartitionClientOptions.StartPosition defaults to "Latest" when the client is created. The connection
//     is lazily initialized, so it's possible the link was initialized to a position after events you've sent.
//     To make this deterministic, you can choose an explicit start point using sequence number, offset or a
//     timestamp. See the [PartitionClientOptions.StartPosition] field for more details.
//
//  2. You might have sent the events to a different partition than intended. By default, batches that are
//     created using [ProducerClient.NewEventDataBatch] do not target a specific partition. When a partition
//     is not specified, Azure Event Hubs service will choose the partition the events will be sent to.
//
//     To fix this, you can specify a PartitionID as part of your [EventDataBatchOptions.PartitionID] options or
//     open multiple [PartitionClient] instances, one for each partition. You can get the full list of partitions
//     at runtime using [ConsumerClient.GetEventHubProperties]. See the "example_consuming_events_test.go" for
//     an example of this pattern.
//
//  3. Network issues can cause internal retries. To see log messages related to this use the instructions in
//     the example function "Example_enableLogging".
func (pc *PartitionClient) ReceiveEvents(ctx context.Context, count int, options *ReceiveEventsOptions) ([]*ReceivedEventData, error) {
	var events []*ReceivedEventData

	prefetchDisabled := pc.prefetch < 0

	err := pc.links.Retry(ctx, EventConsumer, "ReceiveEvents", pc.partitionID, pc.retryOptions, func(ctx context.Context, lwid internal.LinkWithID[amqpwrap.AMQPReceiverCloser]) error {
		events = nil

		if prefetchDisabled {
			remainingCredits := lwid.Link.Credits()

			if count > int(remainingCredits) {
				newCredits := uint32(count) - remainingCredits

				log.Writef(EventConsumer, "(%s) Have %d outstanding credit, only issuing %d credits", lwid.String(), remainingCredits, newCredits)

				if err := lwid.Link.IssueCredit(newCredits); err != nil {
					log.Writef(EventConsumer, "(%s) Error when issuing credits: %s", lwid.String(), err)
					return err
				}
			}
		}

		for {
			amqpMessage, err := lwid.Link.Receive(ctx)

			if internal.IsOwnershipLostError(err) {
				log.Writef(EventConsumer, "(%s) Error, link ownership lost: %s", lwid.String(), err)
				events = nil
				return err
			}

			if err != nil {
				prefetched := getAllPrefetched(lwid.Link, count-len(events))

				for _, amqpMsg := range prefetched {
					re, err := newReceivedEventData(amqpMsg)

					if err != nil {
						log.Writef(EventConsumer, "(%s) Failed converting AMQP message to EventData: %s", lwid.String(), err)
						return err
					}

					events = append(events, re)

					if len(events) == count {
						return nil
					}
				}

				// this lets cancel errors just return
				return err
			}

			receivedEvent, err := newReceivedEventData(amqpMessage)

			if err != nil {
				log.Writef(EventConsumer, "(%s) Failed converting AMQP message to EventData: %s", lwid.String(), err)
				return err
			}

			events = append(events, receivedEvent)

			if len(events) == count {
				return nil
			}
		}
	})

	if err != nil && len(events) == 0 {
		transformedErr := internal.TransformError(err)
		log.Writef(EventConsumer, "No events received, returning error %s", transformedErr.Error())
		return nil, transformedErr
	}

	numEvents := len(events)
	lastSequenceNumber := events[numEvents-1].SequenceNumber

	pc.offsetExpression = formatOffsetExpressionForSequence(">", lastSequenceNumber)
	log.Writef(EventConsumer, "%d Events received, moving sequence to %d", numEvents, lastSequenceNumber)
	return events, nil
}

// Close releases resources for this client.
func (pc *PartitionClient) Close(ctx context.Context) error {
	if pc.links != nil {
		return pc.links.Close(ctx)
	}

	return nil
}

func (pc *PartitionClient) getEntityPath(partitionID string) string {
	return fmt.Sprintf("%s/ConsumerGroups/%s/Partitions/%s", pc.eventHub, pc.consumerGroup, partitionID)
}

func (pc *PartitionClient) newEventHubConsumerLink(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (internal.AMQPReceiverCloser, error) {
	var props map[string]interface{}

	if pc.ownerLevel != nil {
		props = map[string]interface{}{
			"com.microsoft:epoch": *pc.ownerLevel,
		}
	}

	receiverOptions := &amqp.ReceiverOptions{
		SettlementMode: to.Ptr(amqp.ModeFirst),
		Filters: []amqp.LinkFilter{
			amqp.LinkFilterSelector(pc.offsetExpression),
		},
		Properties: props,
	}

	if pc.prefetch > 0 {
		log.Writef(EventConsumer, "Enabling prefetch with %d credits", pc.prefetch)
		receiverOptions.Credit = uint32(pc.prefetch)
	} else if pc.prefetch == 0 {
		log.Writef(EventConsumer, "Enabling prefetch with %d credits", defaultPrefetchSize)
		receiverOptions.Credit = defaultPrefetchSize
	} else {
		// prefetch is disabled, enable manual credits and enable
		// a reasonable default max for the buffer.
		log.Writef(EventConsumer, "Disabling prefetch")
		receiverOptions.ManualCredits = true
		receiverOptions.Credit = defaultMaxCreditSize
	}

	receiver, err := session.NewReceiver(ctx, entityPath, receiverOptions)

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
	namespace internal.NamespaceForAMQPLinks

	eventHub    string
	partitionID string

	consumerGroup string

	retryOptions RetryOptions
}

func newPartitionClient(args partitionClientArgs, options *PartitionClientOptions) (*PartitionClient, error) {
	if options == nil {
		options = &PartitionClientOptions{}
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
		prefetch:         options.Prefetch,
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
