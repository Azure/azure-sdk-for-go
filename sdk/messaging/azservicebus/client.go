// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/devigned/tab"
)

// Client provides methods to create Sender, Receiver and Processor
// instances to send and receive messages from Service Bus.
type Client struct {
	config    clientConfig
	namespace *internal.Namespace
	linksMu   *sync.Mutex
	links     []interface {
		Close(ctx context.Context) error
	}
}

type clientConfig struct {
	connectionString string
	tokenCredential  azcore.TokenCredential
	// the Service Bus namespace name (ex: myservicebus.servicebus.windows.net)
	fullyQualifiedNamespace string
}

// ClientOption is the type for an option that can configure Client.
// For an example option, see `WithConnectionString`
type ClientOption func(client *Client) error

// NewClient creates a new Client for a Service Bus namespace, using a TokenCredential.
// A Client allows you create receivers (for queues or subscriptions) and senders (for queues and topics).
// fullyQualifiedNamespace is the Service Bus namespace name (ex: myservicebus.servicebus.windows.net)
// tokenCredential is one of the credentials in the `github.com/Azure/azure-sdk-for-go/sdk/azidentity` package.
func NewClient(fullyQualifiedNamespace string, tokenCredential azcore.TokenCredential, options ...ClientOption) (*Client, error) {
	return newClientImpl(clientConfig{
		tokenCredential:         tokenCredential,
		fullyQualifiedNamespace: fullyQualifiedNamespace,
	}, options...)
}

// NewClient creates a new Client for a Service Bus namespace, using a TokenCredential.
// A Client allows you create receivers (for queues or subscriptions) and senders (for queues and topics).
// connectionString is a Service Bus connection string for the namespace or for an entity.
func NewClientWithConnectionString(connectionString string, options ...ClientOption) (*Client, error) {
	return newClientImpl(clientConfig{
		connectionString: connectionString,
	}, options...)
}

func newClientImpl(config clientConfig, options ...ClientOption) (*Client, error) {
	client := &Client{
		linksMu: &sync.Mutex{},
		config:  config,
	}

	for _, opt := range options {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	var err error
	var nsOptions []internal.NamespaceOption

	if client.config.connectionString != "" {
		nsOptions = append(nsOptions, internal.NamespaceWithConnectionString(client.config.connectionString))
	} else if client.config.tokenCredential != nil {
		option := internal.NamespacesWithTokenCredential(
			client.config.fullyQualifiedNamespace,
			client.config.tokenCredential)

		nsOptions = append(nsOptions, option)
	} else {
		return nil, errors.New("credentials not specified - use `WithTokenCredential` or `WithConnectionString` to pass in credentials")
	}

	client.namespace, err = internal.NewNamespace(nsOptions...)
	return client, err
}

// NewProcessor creates a Processor, which allows you to receive messages from ServiceBus.
func (client *Client) NewProcessor(options ...ProcessorOption) (*Processor, error) {
	processor, err := newProcessor(client.namespace, options...)

	if err != nil {
		return nil, err
	}

	// TODO: clean up these links
	client.linksMu.Lock()
	client.links = append(client.links, processor)
	client.linksMu.Unlock()

	return processor, nil
}

// NewSender creates a Sender, which allows you to send messages or schedule messages.
func (client *Client) NewSender(queueOrTopic string) (*Sender, error) {
	sender, err := newSender(client.namespace, queueOrTopic)

	if err != nil {
		return nil, err
	}

	// TODO: clean up these links
	client.linksMu.Lock()
	client.links = append(client.links, sender)
	client.linksMu.Unlock()

	return sender, nil
}

// NewReceiver creates a Receiver, which allows you to receive messages.
func (client *Client) NewReceiver(options ...ReceiverOption) (*Receiver, error) {
	receiver, err := newReceiver(client.namespace, options...)

	if err != nil {
		return nil, err
	}

	// TODO: clean up these links
	client.linksMu.Lock()
	client.links = append(client.links, receiver)
	client.linksMu.Unlock()

	return receiver, nil
}

// Close closes the current connection Service Bus as well as any Sender, Receiver or Processors created
// using this client.
func (client *Client) Close(ctx context.Context) error {
	var lastError error

	client.linksMu.Lock()

	for _, link := range client.links {
		if err := link.Close(ctx); err != nil {
			tab.For(ctx).Error(err)
			lastError = err
		}
	}

	client.linksMu.Unlock()

	if lastError != nil {
		return fmt.Errorf("errors while closing links: %w", lastError)
	}
	return nil
}
