// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
)

// SubscriptionProperties represents the static properties of the subscription.
type SubscriptionProperties struct {
	// Name of the subscription relative to the namespace base address.
	Name string

	// LockDuration is the duration a message is locked when using the PeekLock receive mode.
	// Default is 1 minute.
	LockDuration *time.Duration

	// RequiresSession indicates whether the queue supports the concept of sessions.
	// Sessionful-messages follow FIFO ordering.
	// Default is false.
	RequiresSession *bool

	// DefaultMessageTimeToLive is the duration after which the message expires, starting from when
	// the message is sent to Service Bus. This is the default value used when TimeToLive is not
	// set on a message itself.
	DefaultMessageTimeToLive *time.Duration

	// DeadLetteringOnMessageExpiration indicates whether this queue has dead letter
	// support when a message expires.
	DeadLetteringOnMessageExpiration *bool

	// EnableDeadLetteringOnFilterEvaluationExceptions indicates whether messages need to be
	// forwarded to dead-letter sub queue when subscription rule evaluation fails.
	EnableDeadLetteringOnFilterEvaluationExceptions *bool

	// MaxDeliveryCount is the maximum amount of times a message can be delivered before it is automatically
	// sent to the dead letter queue.
	// Default value is 10.
	MaxDeliveryCount *int32

	// Status is the current status of the queue.
	Status *EntityStatus

	// ForwardTo is the name of the recipient entity to which all the messages sent to the queue
	// are forwarded to.
	ForwardTo *string

	// ForwardDeadLetteredMessagesTo is the absolute URI of the entity to forward dead letter messages
	ForwardDeadLetteredMessagesTo *string

	// EnableBatchedOperations indicates whether server-side batched operations are enabled.
	EnableBatchedOperations *bool

	// UserMetadata is custom metadata that user can associate with the description.
	UserMetadata *string
}

// SubscriptionRuntimeProperties represent dynamic properties of a subscription, such as the ActiveMessageCount.
type SubscriptionRuntimeProperties struct {
	// Name is the name of the subscription.
	Name string

	// TotalMessageCount is the number of messages in the subscription.
	TotalMessageCount int64

	// ActiveMessageCount is the number of active messages in the entity.
	ActiveMessageCount int32

	// DeadLetterMessageCount is the number of dead-lettered messages in the entity.
	DeadLetterMessageCount int32

	// TransferMessageCount is the number of messages which are yet to be transferred/forwarded to destination entity.
	TransferMessageCount int32

	// TransferDeadLetterMessageCount is the number of messages transfer-messages which are dead-lettered
	// into transfer-dead-letter subqueue.
	TransferDeadLetterMessageCount int32

	// AccessedAt is when the entity was last updated.
	AccessedAt time.Time

	// CreatedAt is when the entity was created.
	CreatedAt time.Time

	// UpdatedAt is when the entity was last updated.
	UpdatedAt time.Time
}

type AddSubscriptionResponse struct {
	// Value is the result of the request.
	Value *SubscriptionProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type GetSubscriptionResponse struct {
	// Value is the result of the request.
	Value *SubscriptionProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type GetSubscriptionRuntimePropertiesResponse struct {
	// Value is the result of the request.
	Value *SubscriptionRuntimeProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type ListSubscriptionsResponse struct {
	// Value is the result of the request.
	Value []*SubscriptionProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type ListSubscriptionsRuntimePropertiesResponse struct {
	// Value is the result of the request.
	Value []*SubscriptionRuntimeProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type UpdateSubscriptionResponse struct {
	// Value is the result of the request.
	Value *SubscriptionProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type DeleteSubscriptionResponse struct {
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

// AddSubscription creates a subscription to a topic using defaults for all options.
func (ac *AdminClient) AddSubscription(ctx context.Context, topicName string, subscriptionName string) (*AddSubscriptionResponse, error) {
	return ac.AddSubscriptionWithProperties(ctx, topicName, &SubscriptionProperties{
		Name: subscriptionName,
	})
}

// AddSubscriptionWithProperties creates a subscription to a topic with configurable properties
func (ac *AdminClient) AddSubscriptionWithProperties(ctx context.Context, topicName string, properties *SubscriptionProperties) (*AddSubscriptionResponse, error) {
	newProps, resp, err := ac.createOrUpdateSubscriptionImpl(ctx, topicName, properties, true)

	if err != nil {
		return nil, err
	}

	return &AddSubscriptionResponse{
		RawResponse: resp,
		Value:       newProps,
	}, nil
}

// GetSubscription gets a subscription by name.
func (ac *AdminClient) GetSubscription(ctx context.Context, topicName string, subscriptionName string) (*GetSubscriptionResponse, error) {
	var atomResp *atom.SubscriptionEnvelope
	resp, err := ac.em.Get(ctx, fmt.Sprintf("/%s/Subscriptions/%s", topicName, subscriptionName), &atomResp)

	if err != nil {
		return nil, err
	}

	props, err := newSubscriptionProperties(topicName, atomResp.Title, &atomResp.Content.SubscriptionDescription)

	if err != nil {
		return nil, err
	}

	return &GetSubscriptionResponse{
		RawResponse: resp,
		Value:       props,
	}, nil
}

// GetSubscriptionRuntimeProperties gets runtime properties of a subscription, like the SizeInBytes, or SubscriptionCount.
func (ac *AdminClient) GetSubscriptionRuntimeProperties(ctx context.Context, topicName string, subscriptionName string) (*GetSubscriptionRuntimePropertiesResponse, error) {
	var atomResp *atom.SubscriptionEnvelope
	rawResp, err := ac.em.Get(ctx, fmt.Sprintf("/%s/Subscriptions/%s", topicName, subscriptionName), &atomResp)

	if err != nil {
		return nil, err
	}

	return &GetSubscriptionRuntimePropertiesResponse{
		RawResponse: rawResp,
		Value:       newSubscriptionRuntimeProperties(topicName, atomResp.Title, &atomResp.Content.SubscriptionDescription),
	}, nil
}

// ListSubscriptionsOptions can be used to configure the ListSusbscriptions method.
type ListSubscriptionsOptions struct {
	// Top is the maximum size of each page of results.
	Top int
	// Skip is the starting index for the paging operation.
	Skip int
}

// SubscriptionPropertiesPager provides iteration over ListSubscriptionProperties pages.
type SubscriptionPropertiesPager interface {
	// NextPage returns true if the pager advanced to the next page.
	// Returns false if there are no more pages or an error occurred.
	NextPage(context.Context) bool

	// PageResponse returns the current SubscriptionProperties.
	PageResponse() *ListSubscriptionsResponse

	// Err returns the last error encountered while paging.
	Err() error
}

// ListSubscriptions lists subscriptions for a topic.
func (ac *AdminClient) ListSubscriptions(topicName string, options *ListSubscriptionsOptions) SubscriptionPropertiesPager {
	var pageSize int
	var skip int

	if options != nil {
		skip = options.Skip
		pageSize = options.Top
	}

	return &subscriptionPropertiesPager{
		topicName:  topicName,
		innerPager: ac.getSubscriptionPager(topicName, pageSize, skip),
	}
}

// ListSubscriptionsRuntimePropertiesOptions can be used to configure the ListSubscriptionsRuntimeProperties method.
type ListSubscriptionsRuntimePropertiesOptions struct {
	// Top is the maximum size of each page of results.
	Top int
	// Skip is the starting index for the paging operation.
	Skip int
}

// SubscriptionRuntimePropertiesPager provides iteration over ListTopicRuntimeProperties pages.
type SubscriptionRuntimePropertiesPager interface {
	// NextPage returns true if the pager advanced to the next page.
	// Returns false if there are no more pages or an error occurred.
	NextPage(context.Context) bool

	// PageResponse returns the current SubscriptionRuntimeProperties.
	PageResponse() *ListSubscriptionsRuntimePropertiesResponse

	// Err returns the last error encountered while paging.
	Err() error
}

// ListSubscriptionsRuntimeProperties lists runtime properties for subscriptions for a topic.
func (ac *AdminClient) ListSubscriptionsRuntimeProperties(topicName string, options *ListSubscriptionsRuntimePropertiesOptions) SubscriptionRuntimePropertiesPager {
	var pageSize int
	var skip int

	if options != nil {
		skip = options.Skip
		pageSize = options.Top
	}

	return &subscriptionRuntimePropertiesPager{
		innerPager: ac.getSubscriptionPager(topicName, pageSize, skip),
	}
}

// SubscriptionExists checks if a subscription exists.
// Returns true if the subscription is found
// (false, nil) if the subscription is not found
// (false, err) if an error occurred while trying to check if the subscription exists.
func (ac *AdminClient) SubscriptionExists(ctx context.Context, topicName string, subscriptionName string) (bool, error) {
	_, err := ac.GetSubscription(ctx, topicName, subscriptionName)

	if err == nil {
		return true, nil
	}

	if atom.NotFound(err) {
		return false, nil
	}

	return false, err
}

// UpdateSubscription updates an existing subscription.
func (ac *AdminClient) UpdateSubscription(ctx context.Context, topicName string, properties *SubscriptionProperties) (*UpdateSubscriptionResponse, error) {
	newProps, resp, err := ac.createOrUpdateSubscriptionImpl(ctx, topicName, properties, false)

	if err != nil {
		return nil, err
	}

	return &UpdateSubscriptionResponse{
		RawResponse: resp,
		Value:       newProps,
	}, nil
}

// DeleteSubscription deletes a subscription.
func (ac *AdminClient) DeleteSubscription(ctx context.Context, topicName string, subscriptionName string) (*DeleteSubscriptionResponse, error) {
	resp, err := ac.em.Delete(ctx, fmt.Sprintf("/%s/Subscriptions/%s", topicName, subscriptionName))
	defer atom.CloseRes(ctx, resp)
	return &DeleteSubscriptionResponse{
		RawResponse: resp,
	}, err
}

func (ac *AdminClient) createOrUpdateSubscriptionImpl(ctx context.Context, topicName string, props *SubscriptionProperties, creating bool) (*SubscriptionProperties, *http.Response, error) {
	if props == nil {
		return nil, nil, errors.New("properties are required and cannot be nil")
	}

	env := newSubscriptionEnvelope(props, ac.em.TokenProvider())
	var mw []atom.MiddlewareFunc

	if !creating {
		// an update requires the entity to already exist.
		mw = append(mw, func(next atom.RestHandler) atom.RestHandler {
			return func(ctx context.Context, req *http.Request) (*http.Response, error) {
				req.Header.Set("If-Match", "*")
				return next(ctx, req)
			}
		})
	}

	var atomResp *atom.SubscriptionEnvelope
	resp, err := ac.em.Put(ctx, fmt.Sprintf("/%s/Subscriptions/%s", topicName, props.Name), env, &atomResp, mw...)

	if err != nil {
		return nil, nil, err
	}

	newProps, err := newSubscriptionProperties(topicName, props.Name, &atomResp.Content.SubscriptionDescription)

	if err != nil {
		return nil, nil, err
	}

	return newProps, resp, nil
}

func newSubscriptionEnvelope(props *SubscriptionProperties, tokenProvider auth.TokenProvider) *atom.SubscriptionEnvelope {
	desc := &atom.SubscriptionDescription{
		DefaultMessageTimeToLive:                  utils.DurationToStringPtr(props.DefaultMessageTimeToLive),
		LockDuration:                              utils.DurationToStringPtr(props.LockDuration),
		RequiresSession:                           props.RequiresSession,
		DeadLetteringOnMessageExpiration:          props.DeadLetteringOnMessageExpiration,
		DeadLetteringOnFilterEvaluationExceptions: props.EnableDeadLetteringOnFilterEvaluationExceptions,
		MaxDeliveryCount:                          props.MaxDeliveryCount,
		ForwardTo:                                 props.ForwardTo,
		ForwardDeadLetteredMessagesTo:             props.ForwardDeadLetteredMessagesTo,
		UserMetadata:                              props.UserMetadata,
		EnableBatchedOperations:                   props.EnableBatchedOperations,
		// TODO: when we get rule serialization in place.
		// DefaultRuleDescription:                    props.DefaultRuleDescription,
		// are these attributes just not valid anymore?
	}

	return atom.WrapWithSubscriptionEnvelope(desc)
}

func newSubscriptionProperties(topicName string, subscriptionName string, desc *atom.SubscriptionDescription) (*SubscriptionProperties, error) {
	defaultMessageTimeToLive, err := utils.ISO8601StringToDuration(desc.DefaultMessageTimeToLive)

	if err != nil {
		return nil, err
	}

	lockDuration, err := utils.ISO8601StringToDuration(desc.LockDuration)

	if err != nil {
		return nil, err
	}

	return &SubscriptionProperties{
		Name:                             subscriptionName,
		RequiresSession:                  desc.RequiresSession,
		DeadLetteringOnMessageExpiration: desc.DeadLetteringOnMessageExpiration,
		EnableDeadLetteringOnFilterEvaluationExceptions: desc.DeadLetteringOnFilterEvaluationExceptions,
		MaxDeliveryCount:              desc.MaxDeliveryCount,
		ForwardTo:                     desc.ForwardTo,
		ForwardDeadLetteredMessagesTo: desc.ForwardDeadLetteredMessagesTo,
		UserMetadata:                  desc.UserMetadata,
		LockDuration:                  lockDuration,
		DefaultMessageTimeToLive:      defaultMessageTimeToLive,
		EnableBatchedOperations:       desc.EnableBatchedOperations,
		Status:                        (*EntityStatus)(desc.Status),
	}, nil
}

func newSubscriptionRuntimeProperties(topicName string, subscriptionName string, desc *atom.SubscriptionDescription) *SubscriptionRuntimeProperties {
	return &SubscriptionRuntimeProperties{
		Name:                           subscriptionName,
		TotalMessageCount:              *desc.MessageCount,
		ActiveMessageCount:             *desc.CountDetails.ActiveMessageCount,
		DeadLetterMessageCount:         *desc.CountDetails.DeadLetterMessageCount,
		TransferMessageCount:           *desc.CountDetails.TransferMessageCount,
		TransferDeadLetterMessageCount: *desc.CountDetails.TransferDeadLetterMessageCount,
		CreatedAt:                      dateTimeToTime(desc.CreatedAt),
		UpdatedAt:                      dateTimeToTime(desc.UpdatedAt),
		AccessedAt:                     dateTimeToTime(desc.AccessedAt),
	}
}

// subscriptionPropertiesPager provides iteration over SubscriptionProperties pages.
type subscriptionPropertiesPager struct {
	topicName  string
	innerPager subscriptionFeedPagerFunc

	lastErr      error
	lastResponse *ListSubscriptionsResponse
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *subscriptionPropertiesPager) NextPage(ctx context.Context) bool {
	p.lastResponse, p.lastErr = p.getNext(ctx)
	return p.lastResponse != nil
}

// PageResponse returns the current page.
func (p *subscriptionPropertiesPager) PageResponse() *ListSubscriptionsResponse {
	return p.lastResponse
}

// Err returns the last error encountered while paging.
func (p *subscriptionPropertiesPager) Err() error {
	return p.lastErr
}

func (p *subscriptionPropertiesPager) getNext(ctx context.Context) (*ListSubscriptionsResponse, error) {
	feed, resp, err := p.innerPager(ctx)

	if err != nil || len(feed.Entries) == 0 {
		return nil, err
	}

	var all []*SubscriptionProperties

	for _, env := range feed.Entries {
		props, err := newSubscriptionProperties(p.topicName, env.Title, &env.Content.SubscriptionDescription)

		if err != nil {
			return nil, err
		}

		all = append(all, props)
	}

	return &ListSubscriptionsResponse{
		RawResponse: resp,
		Value:       all,
	}, nil
}

// subscriptionRuntimePropertiesPager provides iteration over SubscriptionRuntimeProperties pages.
type subscriptionRuntimePropertiesPager struct {
	topicName  string
	innerPager subscriptionFeedPagerFunc

	lastErr      error
	lastResponse *ListSubscriptionsRuntimePropertiesResponse
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *subscriptionRuntimePropertiesPager) NextPage(ctx context.Context) bool {
	p.lastResponse, p.lastErr = p.getNextPage(ctx)
	return p.lastResponse != nil
}

// PageResponse returns the current page.
func (p *subscriptionRuntimePropertiesPager) PageResponse() *ListSubscriptionsRuntimePropertiesResponse {
	return p.lastResponse
}

// Err returns the last error encountered while paging.
func (p *subscriptionRuntimePropertiesPager) Err() error {
	return p.lastErr
}

func (p *subscriptionRuntimePropertiesPager) getNextPage(ctx context.Context) (*ListSubscriptionsRuntimePropertiesResponse, error) {
	feed, resp, err := p.innerPager(ctx)

	if err != nil || len(feed.Entries) == 0 {
		return nil, err
	}

	var all []*SubscriptionRuntimeProperties

	for _, entry := range feed.Entries {
		all = append(all, newSubscriptionRuntimeProperties(p.topicName, entry.Title, &entry.Content.SubscriptionDescription))
	}

	return &ListSubscriptionsRuntimePropertiesResponse{
		RawResponse: resp,
		Value:       all,
	}, nil
}

type subscriptionFeedPagerFunc func(ctx context.Context) (*atom.SubscriptionFeed, *http.Response, error)

func (ac *AdminClient) getSubscriptionPager(topicName string, top int, skip int) subscriptionFeedPagerFunc {
	return func(ctx context.Context) (*atom.SubscriptionFeed, *http.Response, error) {
		url := fmt.Sprintf("%s/Subscriptions?", topicName)
		if top > 0 {
			url += fmt.Sprintf("&$top=%d", top)
		}

		if skip > 0 {
			url += fmt.Sprintf("&$skip=%d", skip)
		}

		var atomResp *atom.SubscriptionFeed
		resp, err := ac.em.Get(ctx, url, &atomResp)

		if err != nil {
			return nil, nil, err
		}

		skip += len(atomResp.Entries)
		return atomResp, resp, nil
	}
}
