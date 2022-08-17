// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
)

// DefaultConsumerGroup is the name of the default consumer group in the Event Hubs service.
const DefaultConsumerGroup = "$Default"

// ConsumerClientOptions contains options for the `NewConsumerClient` and `NewConsumerClientFromConnectionString`
// functions.
type ConsumerClientOptions struct {
	// TLSConfig configures a client with a custom *tls.Config.
	TLSConfig *tls.Config

	// Application ID that will be passed to the namespace.
	ApplicationID string

	// NewWebSocketConn is a function that can create a net.Conn for use with websockets.
	// For an example, see ExampleNewClient_usingWebsockets() function in example_client_test.go.
	NewWebSocketConn func(ctx context.Context, args NewWebSocketConnArgs) (net.Conn, error)

	// RetryOptions controls how often operations are retried from this client and any
	// Receivers and Senders created from this client.
	RetryOptions RetryOptions

	// StartPosition is the position we will start receiving events from,
	// either an offset (inclusive) with Offset, or receiving events received
	// after a specific time using EnqueuedTime.
	StartPosition StartPosition

	// OwnerLevel is the priority for this consumer, also known as the 'epoch' level.
	// When used, a consumer with a higher OwnerLevel will take ownership of a partition
	// from consumers with a lower OwnerLevel.
	// Default is off.
	OwnerLevel *int64
}

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

// ConsumerClient is used to receive events from an Event Hub partition.
type ConsumerClient struct {
	retryOptions  RetryOptions
	namespace     *internal.Namespace
	eventHub      string
	consumerGroup string
	partitionID   string
	ownerLevel    *int64

	offsetExpression string

	links *internal.Links[amqpwrap.AMQPReceiverCloser]
}

// NewConsumerClient creates a ConsumerClient which uses an azcore.TokenCredential for authentication.
// The consumerGroup is the consumer group for this consumer.
// The fullyQualifiedNamespace is the Event Hubs namespace name (ex: myeventhub.servicebus.windows.net)
// The credential is one of the credentials in the `github.com/Azure/azure-sdk-for-go/sdk/azidentity` package.
func NewConsumerClient(fullyQualifiedNamespace string, eventHub string, partitionID string, consumerGroup string, credential azcore.TokenCredential, options *ConsumerClientOptions) (*ConsumerClient, error) {
	return newConsumerClientImpl(consumerClientArgs{
		fullyQualifiedNamespace: fullyQualifiedNamespace,
		credential:              credential,
		eventHub:                eventHub,
		partitionID:             partitionID,
		consumerGroup:           consumerGroup,
	}, options)
}

// NewConsumerClientFromConnectionString creates a ConsumerClient from a connection string.
// The consumerGroup is the consumer group for this consumer.
//
// connectionString can be one of the following formats:
//
// Connection string, no EntityPath. In this case eventHub cannot be empty.
// ex: Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>
//
// Connection string, has EntityPath. In this case eventHub must be empty.
// ex: Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>;EntityPath=<entity path>
func NewConsumerClientFromConnectionString(connectionString string, eventHub string, partitionID string, consumerGroup string, options *ConsumerClientOptions) (*ConsumerClient, error) {
	parsedConn, err := parseConn(connectionString, eventHub)

	if err != nil {
		return nil, err
	}

	return newConsumerClientImpl(consumerClientArgs{
		connectionString: connectionString,
		eventHub:         parsedConn.HubName,
		partitionID:      partitionID,
		consumerGroup:    consumerGroup,
	}, options)
}

// ReceiveEventsOptions contains optional parameters for the ReceiveEvents function
type ReceiveEventsOptions struct {
	// For future expansion
}

// ReceiveEvents receives events until the context has expired or been cancelled.
func (cc *ConsumerClient) ReceiveEvents(ctx context.Context, count int, options *ReceiveEventsOptions) ([]*ReceivedEventData, error) {
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
					events = append(events, newReceivedEventData(amqpMsg))
				}

				// this lets cancel errors just return
				return err
			}

			receivedEvent := newReceivedEventData(amqpMessage)
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

// GetEventHubProperties gets event hub properties, like the available partition IDs and when the Event Hub was created.
func (cc *ConsumerClient) GetEventHubProperties(ctx context.Context, options *GetEventHubPropertiesOptions) (EventHubProperties, error) {
	rpcLink, err := cc.links.GetManagementLink(ctx)

	if err != nil {
		return EventHubProperties{}, err
	}

	return getEventHubProperties(ctx, cc.namespace, rpcLink.Link, cc.eventHub, options)
}

// GetPartitionProperties gets properties for a specific partition. This includes data like the last enqueued sequence number, the first sequence
// number and when an event was last enqueued to the partition.
func (cc *ConsumerClient) GetPartitionProperties(ctx context.Context, partitionID string, options *GetPartitionPropertiesOptions) (PartitionProperties, error) {
	rpcLink, err := cc.links.GetManagementLink(ctx)

	if err != nil {
		return PartitionProperties{}, err
	}

	return getPartitionProperties(ctx, cc.namespace, rpcLink.Link, cc.eventHub, partitionID, options)
}

// Close closes the consumer's link and the underlying AMQP connection.
func (cc *ConsumerClient) Close(ctx context.Context) error {
	if err := cc.links.Close(ctx); err != nil {
		log.Writef(EventConsumer, "Failed to close link (error might be cached): %s", err.Error())
	}
	return cc.namespace.Close(ctx, true)
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

func (cc *ConsumerClient) getEntityPath(partitionID string) string {
	return fmt.Sprintf("%s/ConsumerGroups/%s/Partitions/%s", cc.eventHub, cc.consumerGroup, partitionID)
}

const defaultLinkRxBuffer = 2048

func (cc *ConsumerClient) newEventHubConsumerLink(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (internal.AMQPReceiverCloser, error) {
	var receiverProps map[string]interface{}

	if cc.ownerLevel != nil {
		receiverProps = map[string]interface{}{
			"com.microsoft:epoch": *cc.ownerLevel,
		}
	}

	receiver, err := session.NewReceiver(ctx, entityPath, &amqp.ReceiverOptions{
		SettlementMode: to.Ptr(amqp.ModeFirst),
		ManualCredits:  true,
		Credit:         defaultLinkRxBuffer,
		Filters: []amqp.LinkFilter{
			amqp.LinkFilterSelector(cc.offsetExpression),
		},
		Properties: receiverProps,
	})

	if err != nil {
		return nil, err
	}

	return receiver, nil
}

type consumerClientArgs struct {
	connectionString string

	// the Event Hubs namespace name (ex: myservicebus.servicebus.windows.net)
	fullyQualifiedNamespace string
	credential              azcore.TokenCredential

	eventHub    string
	partitionID string

	consumerGroup string
}

func newConsumerClientImpl(args consumerClientArgs, options *ConsumerClientOptions) (*ConsumerClient, error) {
	if options == nil {
		options = &ConsumerClientOptions{}
	}

	offsetExpr, err := getOffsetExpression(options.StartPosition)

	if err != nil {
		return nil, err
	}

	client := &ConsumerClient{
		eventHub:         args.eventHub,
		partitionID:      args.partitionID,
		ownerLevel:       options.OwnerLevel,
		consumerGroup:    args.consumerGroup,
		offsetExpression: offsetExpr,
	}

	var nsOptions []internal.NamespaceOption

	if args.connectionString != "" {
		nsOptions = append(nsOptions, internal.NamespaceWithConnectionString(args.connectionString))
	} else if args.credential != nil {
		option := internal.NamespaceWithTokenCredential(
			args.fullyQualifiedNamespace,
			args.credential)

		nsOptions = append(nsOptions, option)
	}

	client.retryOptions = options.RetryOptions

	if options.TLSConfig != nil {
		nsOptions = append(nsOptions, internal.NamespaceWithTLSConfig(options.TLSConfig))
	}

	if options.NewWebSocketConn != nil {
		nsOptions = append(nsOptions, internal.NamespaceWithWebSocket(options.NewWebSocketConn))
	}

	if options.ApplicationID != "" {
		nsOptions = append(nsOptions, internal.NamespaceWithUserAgent(options.ApplicationID))
	}

	nsOptions = append(nsOptions, internal.NamespaceWithRetryOptions(options.RetryOptions))

	tempNS, err := internal.NewNamespace(nsOptions...)

	if err != nil {
		return nil, err
	}

	client.namespace = tempNS
	client.links = internal.NewLinks(tempNS, fmt.Sprintf("%s/$management", client.eventHub), client.getEntityPath, client.newEventHubConsumerLink)

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
