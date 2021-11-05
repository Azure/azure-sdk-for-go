// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/devigned/tab"
)

// Client provides methods to create Sender and Receiver
// instances to send and receive messages from Service Bus.
type Client struct {
	config    clientConfig
	namespace interface {
		// used internally by `Client`
		internal.NamespaceWithNewAMQPLinks
		// for child clients
		internal.NamespaceForAMQPLinks
		internal.NamespaceForMgmtClient
	}
	linksMu     *sync.Mutex
	linkCounter uint64
	links       map[uint64]internal.Closeable
}

// ClientOptions contains options for the `NewClient` and `NewClientFromConnectionString`
// functions.
type ClientOptions struct {
	// TLSConfig configures a client with a custom *tls.Config.
	TLSConfig *tls.Config
}

// NewClient creates a new Client for a Service Bus namespace, using a TokenCredential.
// A Client allows you create receivers (for queues or subscriptions) and senders (for queues and topics).
// fullyQualifiedNamespace is the Service Bus namespace name (ex: myservicebus.servicebus.windows.net)
// credential is one of the credentials in the `github.com/Azure/azure-sdk-for-go/sdk/azidentity` package.
func NewClient(fullyQualifiedNamespace string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if fullyQualifiedNamespace == "" {
		return nil, errors.New("fullyQualifiedNamespace must not be empty")
	}

	if credential == nil {
		return nil, errors.New("credential was nil")
	}

	return newClientImpl(clientConfig{
		credential:              credential,
		fullyQualifiedNamespace: fullyQualifiedNamespace,
	}, options)
}

// NewClientFromConnectionString creates a new Client for a Service Bus namespace using a connection string.
// A Client allows you create receivers (for queues or subscriptions) and senders (for queues and topics).
// connectionString is a Service Bus connection string for the namespace or for an entity.
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	if connectionString == "" {
		return nil, errors.New("connectionString must not be empty")
	}

	return newClientImpl(clientConfig{
		connectionString: connectionString,
	}, options)
}

type clientConfig struct {
	connectionString string
	credential       azcore.TokenCredential
	// the Service Bus namespace name (ex: myservicebus.servicebus.windows.net)
	fullyQualifiedNamespace string
	tlsConfig               *tls.Config
}

func applyClientOptions(client *Client, options *ClientOptions) error {
	if options == nil {
		return nil
	}

	client.config.tlsConfig = options.TLSConfig
	return nil
}

func newClientImpl(config clientConfig, options *ClientOptions) (*Client, error) {
	client := &Client{
		linksMu: &sync.Mutex{},
		config:  config,
		links:   map[uint64]internal.Closeable{},
	}

	if err := applyClientOptions(client, options); err != nil {
		return nil, err
	}

	var err error
	var nsOptions []internal.NamespaceOption

	if client.config.connectionString != "" {
		nsOptions = append(nsOptions, internal.NamespaceWithConnectionString(client.config.connectionString))
	} else if client.config.credential != nil {
		option := internal.NamespacesWithTokenCredential(
			client.config.fullyQualifiedNamespace,
			client.config.credential)

		nsOptions = append(nsOptions, option)
	}

	if client.config.tlsConfig != nil {
		nsOptions = append(nsOptions, internal.NamespaceWithTLSConfig(client.config.tlsConfig))
	}

	client.namespace, err = internal.NewNamespace(nsOptions...)
	return client, err
}

// NewReceiver creates a Receiver for a queue. A receiver allows you to receive messages.
func (client *Client) NewReceiverForQueue(queue string, options *ReceiverOptions) (*Receiver, error) {
	id, cleanupOnClose := client.getCleanupForCloseable()
	receiver, err := newReceiver(client.namespace, &entity{Queue: queue}, cleanupOnClose, options, nil)

	if err != nil {
		return nil, err
	}

	client.addCloseable(id, receiver)
	return receiver, nil
}

// NewReceiver creates a Receiver for a subscription. A receiver allows you to receive messages.
func (client *Client) NewReceiverForSubscription(topic string, subscription string, options *ReceiverOptions) (*Receiver, error) {
	id, cleanupOnClose := client.getCleanupForCloseable()
	receiver, err := newReceiver(client.namespace, &entity{Topic: topic, Subscription: subscription}, cleanupOnClose, options, nil)

	if err != nil {
		return nil, err
	}

	client.addCloseable(id, receiver)
	return receiver, nil
}

// NewSender creates a Sender, which allows you to send messages or schedule messages.
func (client *Client) NewSender(queueOrTopic string) (*Sender, error) {
	id, cleanupOnClose := client.getCleanupForCloseable()
	sender, err := newSender(client.namespace, queueOrTopic, cleanupOnClose)

	if err != nil {
		return nil, err
	}

	client.addCloseable(id, sender)
	return sender, nil
}

// AcceptSessionForQueue accepts a session from a queue with a specific session ID.
// NOTE: this receiver is initialized immediately, not lazily.
func (client *Client) AcceptSessionForQueue(ctx context.Context, queue string, sessionID string, options *SessionReceiverOptions) (*SessionReceiver, error) {
	id, cleanupOnClose := client.getCleanupForCloseable()
	sessionReceiver, err := newSessionReceiver(
		ctx,
		&sessionID,
		client.namespace,
		&entity{Queue: queue},
		cleanupOnClose,
		toReceiverOptions(options))

	if err != nil {
		return nil, err
	}

	if err := sessionReceiver.init(ctx); err != nil {
		return nil, err
	}

	client.addCloseable(id, sessionReceiver)
	return sessionReceiver, nil
}

// AcceptSessionForSubscription accepts a session from a subscription with a specific session ID.
// NOTE: this receiver is initialized immediately, not lazily.
func (client *Client) AcceptSessionForSubscription(ctx context.Context, topic string, subscription string, sessionID string, options *SessionReceiverOptions) (*SessionReceiver, error) {
	id, cleanupOnClose := client.getCleanupForCloseable()
	sessionReceiver, err := newSessionReceiver(
		ctx,
		&sessionID,
		client.namespace,
		&entity{Topic: topic, Subscription: subscription},
		cleanupOnClose,
		toReceiverOptions(options))

	if err != nil {
		return nil, err
	}

	if err := sessionReceiver.init(ctx); err != nil {
		return nil, err
	}

	client.addCloseable(id, sessionReceiver)
	return sessionReceiver, nil
}

// AcceptNextSessionForQueue accepts the next available session from a queue.
// NOTE: this receiver is initialized immediately, not lazily.
func (client *Client) AcceptNextSessionForQueue(ctx context.Context, queue string, options *SessionReceiverOptions) (*SessionReceiver, error) {
	id, cleanupOnClose := client.getCleanupForCloseable()
	sessionReceiver, err := newSessionReceiver(
		ctx,
		nil,
		client.namespace,
		&entity{Queue: queue},
		cleanupOnClose,
		toReceiverOptions(options))

	if err != nil {
		return nil, err
	}

	if err := sessionReceiver.init(ctx); err != nil {
		return nil, err
	}

	client.addCloseable(id, sessionReceiver)
	return sessionReceiver, nil
}

// AcceptNextSessionForSubscription accepts the next available session from a subscription.
// NOTE: this receiver is initialized immediately, not lazily.
func (client *Client) AcceptNextSessionForSubscription(ctx context.Context, topic string, subscription string, options *SessionReceiverOptions) (*SessionReceiver, error) {
	id, cleanupOnClose := client.getCleanupForCloseable()
	sessionReceiver, err := newSessionReceiver(
		ctx,
		nil,
		client.namespace,
		&entity{Topic: topic, Subscription: subscription},
		cleanupOnClose,
		toReceiverOptions(options))

	if err != nil {
		return nil, err
	}

	if err := sessionReceiver.init(ctx); err != nil {
		return nil, err
	}

	client.addCloseable(id, sessionReceiver)
	return sessionReceiver, nil
}

// Close closes the current connection Service Bus as well as any Senders or Receivers created
// using this client.
func (client *Client) Close(ctx context.Context) error {
	var lastError error

	var links []internal.Closeable

	client.linksMu.Lock()

	for _, link := range client.links {
		links = append(links, link)
	}

	client.linksMu.Unlock()

	for _, link := range links {
		if err := link.Close(ctx); err != nil {
			tab.For(ctx).Error(err)
			lastError = err
		}
	}

	if lastError != nil {
		return fmt.Errorf("errors while closing links: %w", lastError)
	}
	return nil
}

func (client *Client) addCloseable(id uint64, closeable internal.Closeable) {
	client.linksMu.Lock()
	client.links[id] = closeable
	client.linksMu.Unlock()
}

func (client *Client) getCleanupForCloseable() (uint64, func()) {
	id := atomic.AddUint64(&client.linkCounter, 1)

	return id, func() {
		client.linksMu.Lock()
		delete(client.links, id)
		client.linksMu.Unlock()
	}
}
