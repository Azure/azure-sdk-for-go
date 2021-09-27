// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/devigned/tab"
)

// Client provides methods to create Sender, Receiver and Processor
// instances to send and receive messages from Service Bus.
type Client struct {
	config struct {
		connectionString string
	}
	namespace *internal.Namespace
	linksMu   *sync.Mutex
	links     []interface {
		Close(ctx context.Context) error
	}
}

// ClientOption is the type for an option that can configure Client.
// For an example option, see `WithConnectionString`
type ClientOption func(client *Client) error

// WithConnectionString configures a namespace with the information provided in a Service Bus connection string
func WithConnectionString(connStr string) ClientOption {
	return func(client *Client) error {
		client.config.connectionString = connStr
		return nil
	}
}

// NewClient creates a new Client.
// Client allows you create receivers (for queues or subscriptions) and
// senders (for queues and topics).
// For creating/deleting/updating queues, topics and subscriptions look at
// `AdministrationClient`.
func NewClient(options ...ClientOption) (*Client, error) {
	client := &Client{
		linksMu: &sync.Mutex{},
	}

	for _, opt := range options {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	var err error
	client.namespace, err = internal.NewNamespace(internal.NamespaceWithConnectionString(client.config.connectionString))
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
