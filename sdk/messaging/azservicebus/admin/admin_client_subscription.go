// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
)

// SubscriptionProperties represents the static properties of the subscription.
type SubscriptionProperties struct {
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

	// AutoDeleteOnIdle is the idle interval after which the subscription is automatically deleted.
	AutoDeleteOnIdle *time.Duration

	// ForwardTo is the name of the recipient entity to which all the messages sent to the queue
	// are forwarded to.
	ForwardTo *string

	// ForwardDeadLetteredMessagesTo is the absolute URI of the entity to forward dead letter messages
	ForwardDeadLetteredMessagesTo *string

	// EnableBatchedOperations indicates whether server-side batched operations are enabled.
	EnableBatchedOperations *bool

	// UserMetadata is custom metadata that user can associate with the subscription.
	UserMetadata *string
}

// SubscriptionRuntimeProperties represent dynamic properties of a subscription, such as the ActiveMessageCount.
type SubscriptionRuntimeProperties struct {
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

type CreateSubscriptionResult struct {
	SubscriptionProperties
}

type CreateSubscriptionResponse struct {
	// Value is the result of the request.
	CreateSubscriptionResult
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type CreateSubscriptionOptions struct {
	// For future expansion
}

// CreateSubscription creates a subscription to a topic with configurable properties
func (ac *Client) CreateSubscription(ctx context.Context, topicName string, subscriptionName string, properties *SubscriptionProperties, options *CreateSubscriptionOptions) (*CreateSubscriptionResponse, error) {
	newProps, resp, err := ac.createOrUpdateSubscriptionImpl(ctx, topicName, subscriptionName, properties, true)

	if err != nil {
		return nil, err
	}

	return &CreateSubscriptionResponse{
		RawResponse: resp,
		CreateSubscriptionResult: CreateSubscriptionResult{
			SubscriptionProperties: *newProps,
		},
	}, nil
}

type GetSubscriptionResult struct {
	SubscriptionProperties
}

type GetSubscriptionResponse struct {
	GetSubscriptionResult

	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type GetSubscriptionOptions struct {
	// For future expansion
}

// GetSubscription gets a subscription by name.
func (ac *Client) GetSubscription(ctx context.Context, topicName string, subscriptionName string, options *GetSubscriptionOptions) (*GetSubscriptionResponse, error) {
	var atomResp *atom.SubscriptionEnvelope
	resp, err := ac.em.Get(ctx, fmt.Sprintf("/%s/Subscriptions/%s", topicName, subscriptionName), &atomResp)

	if err != nil {
		return nil, err
	}

	props, err := newSubscriptionProperties(&atomResp.Content.SubscriptionDescription)

	if err != nil {
		return nil, err
	}

	return &GetSubscriptionResponse{
		RawResponse: resp,
		GetSubscriptionResult: GetSubscriptionResult{
			SubscriptionProperties: *props,
		},
	}, nil
}

type GetSubscriptionRuntimePropertiesResult struct {
	SubscriptionRuntimeProperties
}

type GetSubscriptionRuntimePropertiesResponse struct {
	GetSubscriptionRuntimePropertiesResult
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

// GetSubscriptionRuntimeProperties gets runtime properties of a subscription, like the SizeInBytes, or SubscriptionCount.
func (ac *Client) GetSubscriptionRuntimeProperties(ctx context.Context, topicName string, subscriptionName string) (*GetSubscriptionRuntimePropertiesResponse, error) {
	var atomResp *atom.SubscriptionEnvelope
	rawResp, err := ac.em.Get(ctx, fmt.Sprintf("/%s/Subscriptions/%s", topicName, subscriptionName), &atomResp)

	if err != nil {
		return nil, err
	}

	return &GetSubscriptionRuntimePropertiesResponse{
		RawResponse: rawResp,
		GetSubscriptionRuntimePropertiesResult: GetSubscriptionRuntimePropertiesResult{
			SubscriptionRuntimeProperties: *newSubscriptionRuntimeProperties(&atomResp.Content.SubscriptionDescription),
		},
	}, nil
}

// ListSubscriptionsOptions can be used to configure the ListSusbscriptions method.
type ListSubscriptionsOptions struct {
	// MaxPageSize is the maximum size of each page of results.
	MaxPageSize int32
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

type SubscriptionPropertiesItem struct {
	SubscriptionProperties

	TopicName        string
	SubscriptionName string
}

type ListSubscriptionsResponse struct {
	// Value is the result of the request.
	Items []*SubscriptionPropertiesItem
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

// ListSubscriptions lists subscriptions for a topic.
func (ac *Client) ListSubscriptions(topicName string, options *ListSubscriptionsOptions) SubscriptionPropertiesPager {
	var pageSize int32

	if options != nil {
		pageSize = options.MaxPageSize
	}

	return &subscriptionPropertiesPager{
		topicName:  topicName,
		innerPager: ac.newPagerFunc(fmt.Sprintf("/%s/Subscriptions?", topicName), pageSize, subFeedLen),
	}
}

func subFeedLen(v interface{}) int {
	feed := v.(**atom.SubscriptionFeed)
	return len((*feed).Entries)
}

// ListSubscriptionsRuntimePropertiesOptions can be used to configure the ListSubscriptionsRuntimeProperties method.
type ListSubscriptionsRuntimePropertiesOptions struct {
	// MaxPageSize is the maximum size of each page of results.
	MaxPageSize int32
}

type SubscriptionRuntimePropertiesItem struct {
	SubscriptionRuntimeProperties

	TopicName        string
	SubscriptionName string
}

type ListSubscriptionsRuntimePropertiesResponse struct {
	// Value is the result of the request.
	Items []*SubscriptionRuntimePropertiesItem
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
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
func (ac *Client) ListSubscriptionsRuntimeProperties(topicName string, options *ListSubscriptionsRuntimePropertiesOptions) SubscriptionRuntimePropertiesPager {
	var pageSize int32

	if options != nil {
		pageSize = options.MaxPageSize
	}

	return &subscriptionRuntimePropertiesPager{
		innerPager: ac.newPagerFunc(fmt.Sprintf("/%s/Subscriptions?", topicName), pageSize, subFeedLen),
	}
}

type UpdateSubscriptionResult struct {
	SubscriptionProperties
}

type UpdateSubscriptionResponse struct {
	UpdateSubscriptionResult

	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type UpdateSubscriptionOptions struct {
	// For future expansion
}

// UpdateSubscription updates an existing subscription.
func (ac *Client) UpdateSubscription(ctx context.Context, topicName string, subscriptionName string, properties SubscriptionProperties, options *UpdateSubscriptionOptions) (*UpdateSubscriptionResponse, error) {
	newProps, resp, err := ac.createOrUpdateSubscriptionImpl(ctx, topicName, subscriptionName, &properties, false)

	if err != nil {
		return nil, err
	}

	return &UpdateSubscriptionResponse{
		RawResponse: resp,
		UpdateSubscriptionResult: UpdateSubscriptionResult{
			SubscriptionProperties: *newProps,
		},
	}, nil
}

type DeleteSubscriptionOptions struct {
	// For future expansion
}

type DeleteSubscriptionResponse struct {
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

// DeleteSubscription deletes a subscription.
func (ac *Client) DeleteSubscription(ctx context.Context, topicName string, subscriptionName string, options *DeleteSubscriptionOptions) (*DeleteSubscriptionResponse, error) {
	resp, err := ac.em.Delete(ctx, fmt.Sprintf("/%s/Subscriptions/%s", topicName, subscriptionName))
	defer atom.CloseRes(ctx, resp)
	return &DeleteSubscriptionResponse{
		RawResponse: resp,
	}, err
}

func (ac *Client) createOrUpdateSubscriptionImpl(ctx context.Context, topicName string, subscriptionName string, props *SubscriptionProperties, creating bool) (*SubscriptionProperties, *http.Response, error) {
	if props == nil {
		props = &SubscriptionProperties{}
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
	resp, err := ac.em.Put(ctx, fmt.Sprintf("/%s/Subscriptions/%s", topicName, subscriptionName), env, &atomResp, mw...)

	if err != nil {
		return nil, nil, err
	}

	newProps, err := newSubscriptionProperties(&atomResp.Content.SubscriptionDescription)

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
		AutoDeleteOnIdle:                          utils.DurationToStringPtr(props.AutoDeleteOnIdle),
		// TODO: when we get rule serialization in place.
		// DefaultRuleDescription:                    props.DefaultRuleDescription,
		// are these attributes just not valid anymore?
	}

	return atom.WrapWithSubscriptionEnvelope(desc)
}

func newSubscriptionProperties(desc *atom.SubscriptionDescription) (*SubscriptionProperties, error) {
	defaultMessageTimeToLive, err := utils.ISO8601StringToDuration(desc.DefaultMessageTimeToLive)

	if err != nil {
		return nil, err
	}

	lockDuration, err := utils.ISO8601StringToDuration(desc.LockDuration)

	if err != nil {
		return nil, err
	}

	autoDeleteOnIdle, err := utils.ISO8601StringToDuration(desc.AutoDeleteOnIdle)

	if err != nil {
		return nil, err
	}

	return &SubscriptionProperties{
		RequiresSession:                                 desc.RequiresSession,
		DeadLetteringOnMessageExpiration:                desc.DeadLetteringOnMessageExpiration,
		EnableDeadLetteringOnFilterEvaluationExceptions: desc.DeadLetteringOnFilterEvaluationExceptions,
		MaxDeliveryCount:                                desc.MaxDeliveryCount,
		ForwardTo:                                       desc.ForwardTo,
		ForwardDeadLetteredMessagesTo:                   desc.ForwardDeadLetteredMessagesTo,
		UserMetadata:                                    desc.UserMetadata,
		LockDuration:                                    lockDuration,
		DefaultMessageTimeToLive:                        defaultMessageTimeToLive,
		EnableBatchedOperations:                         desc.EnableBatchedOperations,
		Status:                                          (*EntityStatus)(desc.Status),
		AutoDeleteOnIdle:                                autoDeleteOnIdle,
	}, nil
}

func newSubscriptionRuntimeProperties(desc *atom.SubscriptionDescription) *SubscriptionRuntimeProperties {
	return &SubscriptionRuntimeProperties{
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
	innerPager pagerFunc

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
	var feed *atom.SubscriptionFeed
	resp, err := p.innerPager(ctx, &feed)

	if err != nil || feed == nil {
		return nil, err
	}

	var all []*SubscriptionPropertiesItem

	for _, env := range feed.Entries {
		props, err := newSubscriptionProperties(&env.Content.SubscriptionDescription)

		if err != nil {
			return nil, err
		}

		all = append(all, &SubscriptionPropertiesItem{
			SubscriptionName:       env.Title,
			SubscriptionProperties: *props,
		})
	}

	return &ListSubscriptionsResponse{
		RawResponse: resp,
		Items:       all,
	}, nil
}

// subscriptionRuntimePropertiesPager provides iteration over SubscriptionRuntimeProperties pages.
type subscriptionRuntimePropertiesPager struct {
	topicName  string
	innerPager pagerFunc

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
	var feed *atom.SubscriptionFeed
	resp, err := p.innerPager(ctx, &feed)

	if err != nil || feed == nil {
		return nil, err
	}

	var all []*SubscriptionRuntimePropertiesItem

	for _, entry := range feed.Entries {
		all = append(all, &SubscriptionRuntimePropertiesItem{
			TopicName:                     p.topicName,
			SubscriptionName:              entry.Title,
			SubscriptionRuntimeProperties: *newSubscriptionRuntimeProperties(&entry.Content.SubscriptionDescription),
		})
	}

	return &ListSubscriptionsRuntimePropertiesResponse{
		RawResponse: resp,
		Items:       all,
	}, nil
}
