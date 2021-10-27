// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
)

// AdminClient allows you to administer resources in a Service Bus Namespace.
// For example, you can create queues, enabling capabilities like partitioning, duplicate detection, etc..
// NOTE: For sending and receiving messages you'll need to use the `Client` type instead.
type AdminClient struct {
	em *atom.EntityManager
}

type AdminClientOptions struct {
	// for future expansion
}

// NewAdminClient creates an AdminClient authenticating using a connection string.
func NewAdminClientWithConnectionString(connectionString string, options *AdminClientOptions) (*AdminClient, error) {
	em, err := atom.NewEntityManagerWithConnectionString(connectionString, internal.Version)

	if err != nil {
		return nil, err
	}

	return &AdminClient{em: em}, nil
}

// NewAdminClient creates an AdminClient authenticating using a TokenCredential.
func NewAdminClient(fullyQualifiedNamespace string, tokenCredential azcore.TokenCredential, options *AdminClientOptions) (*AdminClient, error) {
	em, err := atom.NewEntityManager(fullyQualifiedNamespace, tokenCredential, internal.Version)

	if err != nil {
		return nil, err
	}

	return &AdminClient{em: em}, nil
}

// AddQueue creates a queue using defaults for all options.
func (ac *AdminClient) AddQueue(ctx context.Context, queueName string) (*QueueProperties, error) {
	return ac.AddQueueWithProperties(ctx, &QueueProperties{
		Name: queueName,
	})
}

// CreateQueue creates a queue with configurable properties.
func (ac *AdminClient) AddQueueWithProperties(ctx context.Context, properties *QueueProperties) (*QueueProperties, error) {
	return ac.createOrUpdateQueueImpl(ctx, properties, true)
}

// GetQueue gets a queue by name.
func (ac *AdminClient) GetQueue(ctx context.Context, queueName string) (*QueueProperties, error) {
	name, desc, err := ac.getQueueImpl(ctx, queueName)

	if err != nil {
		return nil, err
	}

	return newQueueProperties(name, desc)
}

// GetQueueRuntimeProperties gets runtime properties of a queue, like the SizeInBytes, or ActiveMessageCount.
func (ac *AdminClient) GetQueueRuntimeProperties(ctx context.Context, queueName string) (*QueueRuntimeProperties, error) {
	name, desc, err := ac.getQueueImpl(ctx, queueName)

	if err != nil {
		return nil, err
	}

	return newQueueRuntimeProperties(name, desc), nil
}

// QueueExists checks if a queue exists.
// Returns true if the queue is found
// (false, nil) if the queue is not found
// (false, err) if an error occurred while trying to check if the queue exists.
func (ac *AdminClient) QueueExists(ctx context.Context, queueName string) (bool, error) {
	_, err := ac.GetQueue(ctx, queueName)

	if err == nil {
		return true, nil
	}

	var httpResponse azcore.HTTPResponse

	if errors.As(err, &httpResponse) && httpResponse.RawResponse().StatusCode == 404 {
		return false, nil
	}

	return false, err
}

// UpdateQueue updates an existing queue.
func (ac *AdminClient) UpdateQueue(ctx context.Context, properties *QueueProperties) (*QueueProperties, error) {
	return ac.createOrUpdateQueueImpl(ctx, properties, false)
}

// DeleteQueue deletes a queue.
func (ac *AdminClient) DeleteQueue(ctx context.Context, queueName string) (*http.Response, error) {
	resp, err := ac.em.Delete(ctx, "/"+queueName)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ListQueuesOptions can be used to configure the ListQueues method.
type ListQueuesOptions struct {
	// Top is the maximum size of each page of results.
	Top int
	// Skip is the starting index for the paging operation.
	Skip int
}

// QueuePropertiesPager provides iteration over ListQueueProperties pages.
type QueuePropertiesPager interface {
	// NextPage returns true if the pager advanced to the next page.
	// Returns false if there are no more pages or an error occurred.
	NextPage(context.Context) bool

	// PageResponse returns the current QueueProperties.
	PageResponse() []*QueueProperties

	// Err returns the last error encountered while paging.
	Err() error
}

// ListQueues lists queues.
func (ac *AdminClient) ListQueues(options *ListQueuesOptions) QueuePropertiesPager {
	var pageSize int
	var skip int

	if options != nil {
		skip = options.Skip
		pageSize = options.Top
	}

	return &queuePropertiesPager{
		adminClient: ac,
		pageSize:    pageSize,
		skip:        skip,
	}
}

// ListQueuesRuntimePropertiesOptions can be used to configure the ListQueuesRuntimeProperties method.
type ListQueuesRuntimePropertiesOptions struct {
	// Top is the maximum size of each page of results.
	Top int
	// Skip is the starting index for the paging operation.
	Skip int
}

// QueueRuntimePropertiesPager provides iteration over ListQueueRuntimeProperties pages.
type QueueRuntimePropertiesPager interface {
	// NextPage returns true if the pager advanced to the next page.
	// Returns false if there are no more pages or an error occurred.
	NextPage(context.Context) bool

	// PageResponse returns the current QueueRuntimeProperties.
	PageResponse() []*QueueRuntimeProperties

	// Err returns the last error encountered while paging.
	Err() error
}

// ListQueuesRuntimeProperties lists runtime properties for queues.
func (ac *AdminClient) ListQueuesRuntimeProperties(options *ListQueuesRuntimePropertiesOptions) QueueRuntimePropertiesPager {
	var pageSize int
	var skip int

	if options != nil {
		skip = options.Skip
		pageSize = options.Top
	}

	return &queueRuntimePropertiesPager{
		adminClient: ac,
		pageSize:    pageSize,
		skip:        skip,
	}
}

// func (ac *AdminClient) GetNamespaceProperties() {}

// func (ac *AdminClient) CreateTopic()                        {}
// func (ac *AdminClient) CreateSubscription()                 {}
// func (ac *AdminClient) CreateRule()                         {}
// func (ac *AdminClient) DeleteTopic()                        {}
// func (ac *AdminClient) DeleteSubscription()                 {}
// func (ac *AdminClient) DeleteRule()                         {}
// func (ac *AdminClient) GetRule()                            {}
// func (ac *AdminClient) GetSubscription()                    {}
// func (ac *AdminClient) GetSubscriptionRuntimeProperties()   {}
// func (ac *AdminClient) GetTopic()                           {}
// func (ac *AdminClient) GetTopicRuntimeProperties()          {}
// func (ac *AdminClient) ListRules()                          {}
// func (ac *AdminClient) ListTopics()                         {}
// func (ac *AdminClient) ListTopicsRuntimeProperties()        {}
// func (ac *AdminClient) ListSubscriptions()                  {}
// func (ac *AdminClient) ListSubscriptionsRuntimeProperties() {}

// func (ac *AdminClient) TopicExists()               {}
// func (ac *AdminClient) SubscriptionExists()        {}
// func (ac *AdminClient) RuleExists()                {}

// func (ac *AdminClient) UpdateTopic()        {}
// func (ac *AdminClient) UpdateSubscription() {}
// func (ac *AdminClient) UpdateRule()         {}
