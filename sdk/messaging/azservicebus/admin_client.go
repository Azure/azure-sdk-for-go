// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/Azure/azure-amqp-common-go/v3/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/devigned/tab"
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

// CreateQueue creates a queue using defaults for all options.
func (ac *AdminClient) CreateQueue(ctx context.Context, queueName string) (*QueueProperties, error) {
	return ac.CreateQueueWithProperties(ctx, &QueueProperties{
		Name: queueName,
	})
}

// CreateQueue creates a queue with configurable properties.
func (ac *AdminClient) CreateQueueWithProperties(ctx context.Context, props *QueueProperties) (*QueueProperties, error) {
	if props == nil {
		return nil, errors.New("properties are required and cannot be nil")
	}

	reqBytes, mw, err := serializeQueueProperties(props, ac.em.TokenProvider())

	if err != nil {
		return nil, err
	}

	resp, err := ac.em.Put(ctx, "/"+props.Name, reqBytes, mw...)
	defer atom.CloseRes(ctx, resp)

	if err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	return deserializeQueueEnvelope(props.Name, b)
}

func serializeQueueProperties(props *QueueProperties, tokenProvider auth.TokenProvider) ([]byte, []atom.MiddlewareFunc, error) {
	qpr := &atom.QueuePutRequest{
		LockDuration:                        utils.DurationToStringPtr(props.LockDuration),
		MaxSizeInMegabytes:                  props.MaxSizeInMegabytes,
		RequiresDuplicateDetection:          props.RequiresDuplicateDetection,
		RequiresSession:                     props.RequiresSession,
		DefaultMessageTimeToLive:            utils.DurationToStringPtr(props.DefaultMessageTimeToLive),
		DeadLetteringOnMessageExpiration:    props.DeadLetteringOnMessageExpiration,
		DuplicateDetectionHistoryTimeWindow: utils.DurationToStringPtr(props.DuplicateDetectionHistoryTimeWindow),
		MaxDeliveryCount:                    props.MaxDeliveryCount,
		EnableBatchedOperations:             props.EnableBatchedOperations,
		Status:                              (*atom.EntityStatus)(props.Status),
		AutoDeleteOnIdle:                    utils.DurationToStringPtr(props.AutoDeleteOnIdle),
		EnablePartitioning:                  props.EnablePartitioning,
		ForwardTo:                           props.ForwardTo,
		ForwardDeadLetteredMessagesTo:       props.ForwardDeadLetteredMessagesTo,
	}

	atomRequest, mw := atom.ConvertToQueueRequest(qpr, tokenProvider)

	bytes, err := xml.MarshalIndent(atomRequest, "", "  ")

	if err != nil {
		return nil, nil, err
	}

	return bytes, mw, nil
}

func deserializeQueueEnvelope(name string, b []byte) (*QueueProperties, error) {
	var atomResp atom.QueueEnvelope

	if err := xml.Unmarshal(b, &atomResp); err != nil {
		return nil, atom.FormatManagementError(b)
	}

	respQD := atomResp.Content.QueueDescription

	lockDuration, err := utils.ISO8601StringToDuration(respQD.LockDuration)

	if err != nil {
		return nil, err
	}

	defaultMessageTimeToLive, err := utils.ISO8601StringToDuration(respQD.DefaultMessageTimeToLive)

	if err != nil {
		return nil, err
	}

	duplicateDetectionHistoryTimeWindow, err := utils.ISO8601StringToDuration(respQD.DuplicateDetectionHistoryTimeWindow)

	if err != nil {
		return nil, err
	}

	autoDeleteOnIdle, err := utils.ISO8601StringToDuration(respQD.AutoDeleteOnIdle)

	if err != nil {
		return nil, err
	}

	queuePropsResult := &QueueProperties{
		Name:                                name,
		LockDuration:                        lockDuration,
		MaxSizeInMegabytes:                  respQD.MaxSizeInMegabytes,
		RequiresDuplicateDetection:          respQD.RequiresDuplicateDetection,
		RequiresSession:                     respQD.RequiresSession,
		DefaultMessageTimeToLive:            defaultMessageTimeToLive,
		DeadLetteringOnMessageExpiration:    respQD.DeadLetteringOnMessageExpiration,
		DuplicateDetectionHistoryTimeWindow: duplicateDetectionHistoryTimeWindow,
		MaxDeliveryCount:                    respQD.MaxDeliveryCount,
		EnableBatchedOperations:             respQD.EnableBatchedOperations,
		Status:                              (*EntityStatus)(respQD.Status),
		AutoDeleteOnIdle:                    autoDeleteOnIdle,
		EnablePartitioning:                  respQD.EnablePartitioning,
		ForwardTo:                           respQD.ForwardTo,
		ForwardDeadLetteredMessagesTo:       respQD.ForwardDeadLetteredMessagesTo,
	}

	return queuePropsResult, nil
}

func (ac *AdminClient) GetQueue(ctx context.Context, queueName string) (*QueueProperties, error) {
	ctx, span := tab.StartSpan(ctx, tracing.SpanGetEntity)
	defer span.End()

	resp, err := ac.em.Get(ctx, "/"+queueName)

	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return deserializeQueueEnvelope(queueName, b)
}

// func (ac *AdminClient) GetQueueRuntimeProperties() (*QueueRuntimeProperties, error) {
// 	return nil, nil
// }

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

// func (ac *AdminClient) UpdateQueue(properties *QueueProperties) (*QueueProperties, error) {
// 	return nil, nil
// }

func (ac *AdminClient) DeleteQueue(ctx context.Context, queueName string) (*http.Response, error) {
	resp, err := ac.em.Delete(ctx, "/"+queueName)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// func (ac *AdminClient) ListQueues()                  {}
// func (ac *AdminClient) ListQueuesRuntimeProperties() {}

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
