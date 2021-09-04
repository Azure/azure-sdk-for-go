// Package azservicebus provides clients for sending and receiving messages with Azure ServiceBus
// as well as modifying resources like Queues, Topics and Subscriptions.
package azservicebus

import (
	"context"
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/servicebus/azservicebus/internal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/devigned/tab"
)

// ServiceBusClient provides methods to create Sender, Receiver and Processor
// instances to send and receive messages from Service Bus.
type ServiceBusClient struct {
	namespace *internal.Namespace
	linksMu   *sync.Mutex
	links     []interface {
		Close(ctx context.Context) error
	}
}

// ServiceBusClientOption is the type for an option that can configure ServiceBusClient.
// For an example option, see `ServiceBusWithConnectionString`
type ServiceBusClientOption func(client *ServiceBusClient) error

// ServiceBusWithConnectionString configures a namespace with the information provided in a Service Bus connection string
func ServiceBusWithConnectionString(connStr string) ServiceBusClientOption {
	return func(client *ServiceBusClient) error {
		fn := internal.NamespaceWithConnectionString(connStr)
		return fn(client.namespace)
	}
}

// NewServiceBusClient creates a new ServiceBusClient.
// ServiceBusClient allows you create receivers (for queues or subscriptions) and
// senders (for queues and topics).
// For creating/deleting/updating queues, topics and subscriptions look at
// ServiceBusAdministrationClient.
func NewServiceBusClient(options ...ServiceBusClientOption) (*ServiceBusClient, error) {
	client := &ServiceBusClient{
		namespace: &internal.Namespace{
			Environment: azure.PublicCloud,
		},
		linksMu: &sync.Mutex{},
	}

	for _, opt := range options {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// NewProcessor creates a Processor, which allows you to receive messages from ServiceBus.
func (client *ServiceBusClient) NewProcessor(options ...ProcessorOption) (*Processor, error) {
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
func (client *ServiceBusClient) NewSender(queueOrTopic string, options ...SenderOption) (*Sender, error) {
	sender, err := newSender(client.namespace, queueOrTopic, options...)

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
func (client *ServiceBusClient) NewReceiver(options ...ReceiverOption) (*Receiver, error) {
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
func (client *ServiceBusClient) Close(ctx context.Context) error {
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
		return fmt.Errorf("errors while closing links: %W", lastError)
	}
	return nil
}
