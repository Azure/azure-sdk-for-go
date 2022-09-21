// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventhubs

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/conn"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/go-amqp"
)

// NewWebSocketConnArgs are passed to your web socket creation function (ClientOptions.NewWebSocketConn)
type NewWebSocketConnArgs = exported.NewWebSocketConnArgs

// RetryOptions represent the options for retries.
type RetryOptions = exported.RetryOptions

// ProducerClientOptions contains options for the `NewProducerClient` and `NewProducerClientFromConnectionString`
// functions.
type ProducerClientOptions struct {
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
}

// ProducerClient can be used to send events to an Event Hub.
type ProducerClient struct {
	retryOptions RetryOptions
	namespace    internal.NamespaceForProducerOrConsumer
	eventHub     string

	links *internal.Links[amqpwrap.AMQPSenderCloser]
}

// anyPartitionID is what we target if we want to send a message and let Event Hubs pick a partition
// or if we're doing an operation that isn't partition specific, such as querying the management link
// to get event hub properties or partition properties.
const anyPartitionID = ""

// NewProducerClient creates a ProducerClient which uses an azcore.TokenCredential for authentication.
// The fullyQualifiedNamespace is the Event Hubs namespace name (ex: myeventhub.servicebus.windows.net)
// The credential is one of the credentials in the `github.com/Azure/azure-sdk-for-go/sdk/azidentity` package.
func NewProducerClient(fullyQualifiedNamespace string, eventHub string, credential azcore.TokenCredential, options *ProducerClientOptions) (*ProducerClient, error) {
	return newProducerClientImpl(producerClientCreds{
		fullyQualifiedNamespace: fullyQualifiedNamespace,
		credential:              credential,
		eventHub:                eventHub,
	}, options)
}

// NewProducerClientFromConnectionString creates a ProducerClient from a connection string.
//
// connectionString can be one of the following formats:
//
// Connection string, no EntityPath. In this case eventHub cannot be empty.
// ex: Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>
//
// Connection string, has EntityPath. In this case eventHub must be empty.
// ex: Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>;EntityPath=<entity path>
func NewProducerClientFromConnectionString(connectionString string, eventHub string, options *ProducerClientOptions) (*ProducerClient, error) {
	parsedConn, err := parseConn(connectionString, eventHub)

	if err != nil {
		return nil, err
	}

	return newProducerClientImpl(producerClientCreds{
		connectionString: connectionString,
		eventHub:         parsedConn.HubName,
	}, options)
}

// NewEventDataBatchOptions contains optional parameters for the NewEventDataBatch function
type NewEventDataBatchOptions struct {
	// MaxBytes overrides the max size (in bytes) for a batch.
	// By default NewMessageBatch will use the max message size provided by the service.
	MaxBytes uint64

	// PartitionKey is hashed to calculate the partition assignment. Messages and message
	// batches with the same PartitionKey are guaranteed to end up in the same partition.
	// Note that if you use this option then PartitionID cannot be set.
	PartitionKey *string

	// PartitionID is the ID of the partition to send these messages to.
	// Note that if you use this option then PartitionKey cannot be set.
	PartitionID *string
}

// NewEventDataBatch can be used to create a batch that contain multiple events.
// If the operation fails it can return an *azeventhubs.Error type if the failure is actionable.
func (pc *ProducerClient) NewEventDataBatch(ctx context.Context, options *NewEventDataBatchOptions) (*EventDataBatch, error) {
	var batch *EventDataBatch

	partitionID := anyPartitionID

	if options != nil && options.PartitionID != nil {
		partitionID = *options.PartitionID
	}

	err := pc.links.Retry(ctx, exported.EventProducer, "NewEventDataBatch", partitionID, pc.retryOptions, func(ctx context.Context, lwid internal.LinkWithID[amqpwrap.AMQPSenderCloser]) error {
		tmpBatch, err := newEventDataBatch(lwid.Link, options)

		if err != nil {
			return err
		}

		batch = tmpBatch
		return nil
	})

	if err != nil {
		return nil, internal.TransformError(err)
	}

	return batch, nil
}

// SendEventBatchOptions contains optional parameters for the SendEventBatch function
type SendEventBatchOptions struct {
	// For future expansion
}

// SendEventBatch sends an event data batch to Event Hubs.
func (pc *ProducerClient) SendEventBatch(ctx context.Context, batch *EventDataBatch, options *SendEventBatchOptions) error {
	err := pc.links.Retry(ctx, exported.EventProducer, "SendEventBatch", getPartitionID(batch.partitionID), pc.retryOptions, func(ctx context.Context, lwid internal.LinkWithID[amqpwrap.AMQPSenderCloser]) error {
		return lwid.Link.Send(ctx, batch.toAMQPMessage())
	})
	return internal.TransformError(err)
}

// GetPartitionProperties gets properties for a specific partition. This includes data like the last enqueued sequence number, the first sequence
// number and when an event was last enqueued to the partition.
func (pc *ProducerClient) GetPartitionProperties(ctx context.Context, partitionID string, options *GetPartitionPropertiesOptions) (PartitionProperties, error) {
	rpcLink, err := pc.links.GetManagementLink(ctx)

	if err != nil {
		return PartitionProperties{}, err
	}

	return getPartitionProperties(ctx, pc.namespace, rpcLink.Link, pc.eventHub, partitionID, options)
}

// GetEventHubProperties gets event hub properties, like the available partition IDs and when the Event Hub was created.
func (pc *ProducerClient) GetEventHubProperties(ctx context.Context, options *GetEventHubPropertiesOptions) (EventHubProperties, error) {
	rpcLink, err := pc.links.GetManagementLink(ctx)

	if err != nil {
		return EventHubProperties{}, err
	}

	return getEventHubProperties(ctx, pc.namespace, rpcLink.Link, pc.eventHub, options)
}

// Close closes the producer's AMQP links and the underlying AMQP connection.
func (pc *ProducerClient) Close(ctx context.Context) error {
	if err := pc.links.Close(ctx); err != nil {
		azlog.Writef(EventProducer, "Failed when closing links while shutting down producer client: %s", err.Error())
	}
	return pc.namespace.Close(ctx, true)
}

func (pc *ProducerClient) getEntityPath(partitionID string) string {
	if partitionID != anyPartitionID {
		return fmt.Sprintf("%s/Partitions/%s", pc.eventHub, partitionID)
	} else {
		// this is the "let Event Hubs" decide link - any sends that occur here will
		// end up getting distributed to different partitions on the service side, rather
		// then being specified in the client.
		return pc.eventHub
	}
}

func (pc *ProducerClient) newEventHubProducerLink(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (amqpwrap.AMQPSenderCloser, error) {
	sender, err := session.NewSender(ctx, entityPath, &amqp.SenderOptions{
		SettlementMode:              to.Ptr(amqp.ModeMixed),
		RequestedReceiverSettleMode: to.Ptr(amqp.ModeFirst),
		IgnoreDispositionErrors:     true,
	})

	if err != nil {
		return nil, err
	}

	return sender, nil
}

type producerClientCreds struct {
	connectionString string

	// the Event Hubs namespace name (ex: myservicebus.servicebus.windows.net)
	fullyQualifiedNamespace string
	credential              azcore.TokenCredential

	eventHub string
}

func newProducerClientImpl(creds producerClientCreds, options *ProducerClientOptions) (*ProducerClient, error) {
	client := &ProducerClient{
		eventHub: creds.eventHub,
	}

	var err error
	var nsOptions []internal.NamespaceOption

	if creds.connectionString != "" {
		nsOptions = append(nsOptions, internal.NamespaceWithConnectionString(creds.connectionString))
	} else if creds.credential != nil {
		option := internal.NamespaceWithTokenCredential(
			creds.fullyQualifiedNamespace,
			creds.credential)

		nsOptions = append(nsOptions, option)
	}

	if options != nil {
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
	}

	tmpNS, err := internal.NewNamespace(nsOptions...)

	if err != nil {
		return nil, err
	}

	client.namespace = tmpNS

	client.links = internal.NewLinks(tmpNS, fmt.Sprintf("%s/$management", client.eventHub), client.getEntityPath, client.newEventHubProducerLink)

	return client, err
}

func parseConn(connectionString string, eventHub string) (*conn.ParsedConn, error) {
	parsedConn, err := conn.ParsedConnectionFromStr(connectionString)

	if err != nil {
		return nil, err
	}

	if parsedConn.HubName == "" {
		if eventHub == "" {
			return nil, errors.New("connection string does not contain an EntityPath. eventHub cannot be an empty string")
		}
		parsedConn.HubName = eventHub
	} else if parsedConn.HubName != "" {
		if eventHub != "" {
			return nil, errors.New("connection string contains an EntityPath. eventHub must be an empty string")
		}
	}

	return parsedConn, nil
}

func getPartitionID(partitionID *string) string {
	if partitionID != nil {
		return *partitionID
	}

	return anyPartitionID
}
