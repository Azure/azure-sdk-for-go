// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azeventhubs

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"
)

// Error represents an Event Hub specific error.
// NOTE: the Code is considered part of the published API but the message that
// comes back from Error(), as well as the underlying wrapped error, are NOT and
// are subject to change.
type Error = exported.Error

// Code is an error code, usable by consuming code to work with
// programatically.
type Code = exported.Code

const (
	// CodeConnectionLost means our connection was lost and all retry attempts failed.
	// This typically reflects an extended outage or connection disruption and may
	// require manual intervention.
	CodeConnectionLost = exported.CodeConnectionLost

	// CodeOwnershipLost means that a partition that you were reading from was opened
	// by another link with a higher epoch/owner level.
	CodeOwnershipLost = exported.CodeOwnershipLost
)

// ConsumerClientOptions configures optional parameters for a ConsumerClient.
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
	OwnerLevel *uint64
}

// ConsumerClient can create PartitionClient instances, which can read events from
// a partition.
type ConsumerClient struct {
	consumerGroup string
	eventHub      string
	retryOptions  RetryOptions
	namespace     *internal.Namespace
	links         *internal.Links[amqpwrap.AMQPReceiverCloser]

	clientID string
}

// NewConsumerClient creates a ConsumerClient which uses an azcore.TokenCredential for authentication.
// The fullyQualifiedNamespace is the Event Hubs namespace name (ex: myeventhub.servicebus.windows.net)
// The credential is one of the credentials in the `github.com/Azure/azure-sdk-for-go/sdk/azidentity` package.
func NewConsumerClient(fullyQualifiedNamespace string, eventHub string, consumerGroup string, credential azcore.TokenCredential, options *ConsumerClientOptions) (*ConsumerClient, error) {
	return newConsumerClient(consumerClientArgs{
		consumerGroup:           consumerGroup,
		fullyQualifiedNamespace: fullyQualifiedNamespace,
		eventHub:                eventHub,
		credential:              credential,
	}, options)
}

// NewConsumerClientFromConnectionString creates a ConsumerClient from a connection string.
//
// connectionString can be one of the following formats:
//
// Connection string, no EntityPath. In this case eventHub cannot be empty.
// ex: Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>
//
// Connection string, has EntityPath. In this case eventHub must be empty.
// ex: Endpoint=sb://<your-namespace>.servicebus.windows.net/;SharedAccessKeyName=<key-name>;SharedAccessKey=<key>;EntityPath=<entity path>
func NewConsumerClientFromConnectionString(connectionString string, eventHub string, consumerGroup string, options *ConsumerClientOptions) (*ConsumerClient, error) {
	parsedConn, err := parseConn(connectionString, eventHub)

	if err != nil {
		return nil, err
	}

	return newConsumerClient(consumerClientArgs{
		consumerGroup:    consumerGroup,
		connectionString: connectionString,
		eventHub:         parsedConn.HubName,
	}, options)
}

// NewPartitionClientOptions provides options for the Subscribe function.
type NewPartitionClientOptions struct {
	// StartPosition is the position we will start receiving events from,
	// either an offset (inclusive) with Offset, or receiving events received
	// after a specific time using EnqueuedTime.
	StartPosition StartPosition

	// OwnerLevel is the priority for this partition client, also known as the 'epoch' level.
	// When used, a partition client with a higher OwnerLevel will take ownership of a partition
	// from partition clients with a lower OwnerLevel.
	// Default is off.
	OwnerLevel *int64
}

// NewPartitionClient creates a client that can receive events from a partition.
func (cc *ConsumerClient) NewPartitionClient(partitionID string, options *NewPartitionClientOptions) (*PartitionClient, error) {
	return newPartitionClient(partitionClientArgs{
		namespace:     cc.namespace,
		eventHub:      cc.eventHub,
		partitionID:   partitionID,
		consumerGroup: cc.consumerGroup,
		retryOptions:  cc.retryOptions,
	}, options)
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

type consumerClientDetails struct {
	FullyQualifiedNamespace string
	ConsumerGroup           string
	EventHubName            string
	ClientID                string
}

func (cc *ConsumerClient) getDetails() consumerClientDetails {
	return consumerClientDetails{
		FullyQualifiedNamespace: cc.namespace.FQDN,
		ConsumerGroup:           cc.consumerGroup,
		EventHubName:            cc.eventHub,
		ClientID:                cc.clientID,
	}
}

// Close closes the connection for this client.
func (cc *ConsumerClient) Close(ctx context.Context) error {
	return cc.namespace.Close(ctx, true)
}

type consumerClientArgs struct {
	connectionString string

	// the Event Hubs namespace name (ex: myservicebus.servicebus.windows.net)
	fullyQualifiedNamespace string
	credential              azcore.TokenCredential

	consumerGroup string
	eventHub      string
}

func newConsumerClient(args consumerClientArgs, options *ConsumerClientOptions) (*ConsumerClient, error) {
	if options == nil {
		options = &ConsumerClientOptions{}
	}

	clientUUID, err := uuid.New()

	if err != nil {
		return nil, err
	}

	client := &ConsumerClient{
		consumerGroup: args.consumerGroup,
		eventHub:      args.eventHub,
		clientID:      clientUUID.String(),
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
	client.links = internal.NewLinks[amqpwrap.AMQPReceiverCloser](tempNS, fmt.Sprintf("%s/$management", client.eventHub), nil, nil)

	return client, nil
}
