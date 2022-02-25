// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/internal/auth"
)

// TopicProperties represents the static properties of the topic.
type TopicProperties struct {
	// MaxSizeInMegabytes - The maximum size of the topic in megabytes, which is the size of memory
	// allocated for the topic.
	// Default is 1024.
	MaxSizeInMegabytes *int32

	// RequiresDuplicateDetection indicates if this topic requires duplicate detection.
	RequiresDuplicateDetection *bool

	// DefaultMessageTimeToLive is the duration after which the message expires, starting from when
	// the message is sent to Service Bus. This is the default value used when TimeToLive is not
	// set on a message itself.
	DefaultMessageTimeToLive *time.Duration

	// DuplicateDetectionHistoryTimeWindow is the duration of duplicate detection history.
	// Default value is 10 minutes.
	DuplicateDetectionHistoryTimeWindow *time.Duration

	// EnableBatchedOperations indicates whether server-side batched operations are enabled.
	EnableBatchedOperations *bool

	// Status is the current status of the topic.
	Status *EntityStatus

	// AutoDeleteOnIdle is the idle interval after which the topic is automatically deleted.
	AutoDeleteOnIdle *time.Duration

	// EnablePartitioning indicates whether the topic is to be partitioned across multiple message brokers.
	EnablePartitioning *bool

	// SupportOrdering defines whether ordering needs to be maintained. If true, messages
	// sent to topic will be forwarded to the subscription, in order.
	SupportOrdering *bool

	// UserMetadata is custom metadata that user can associate with the topic.
	UserMetadata *string
}

// TopicRuntimeProperties represent dynamic properties of a topic, such as the ActiveMessageCount.
type TopicRuntimeProperties struct {
	// SizeInBytes - The size of the topic, in bytes.
	SizeInBytes int64

	// CreatedAt is when the entity was created.
	CreatedAt time.Time

	// UpdatedAt is when the entity was last updated.
	UpdatedAt time.Time

	// AccessedAt is when the entity was last updated.
	AccessedAt time.Time

	// SubscriptionCount is the number of subscriptions to the topic.
	SubscriptionCount int32

	// ScheduledMessageCount is the number of messages that are scheduled to be entopicd.
	ScheduledMessageCount int32
}

type CreateTopicResult struct {
	TopicProperties
}

type CreateTopicResponse struct {
	CreateTopicResult

	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type CreateTopicOptions struct {
	// For future expansion
}

// CreateTopic creates a topic using defaults for all options.
func (ac *Client) CreateTopic(ctx context.Context, topicName string, properties *TopicProperties, options *CreateTopicOptions) (*CreateTopicResponse, error) {
	newProps, resp, err := ac.createOrUpdateTopicImpl(ctx, topicName, properties, true)

	if err != nil {
		return nil, err
	}

	return &CreateTopicResponse{
		RawResponse: resp,
		CreateTopicResult: CreateTopicResult{
			TopicProperties: *newProps,
		},
	}, nil
}

type GetTopicResult struct {
	TopicProperties
}

type GetTopicResponse struct {
	GetTopicResult

	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type GetTopicOptions struct {
	// For future expansion
}

// GetTopic gets a topic by name.
func (ac *Client) GetTopic(ctx context.Context, topicName string, options *GetTopicOptions) (*GetTopicResponse, error) {
	var atomResp *atom.TopicEnvelope
	resp, err := ac.em.Get(ctx, "/"+topicName, &atomResp)

	if err != nil {
		return nil, err
	}

	props, err := newTopicProperties(&atomResp.Content.TopicDescription)

	if err != nil {
		return nil, err
	}

	return &GetTopicResponse{
		RawResponse: resp,
		GetTopicResult: GetTopicResult{
			TopicProperties: *props,
		},
	}, nil
}

type GetTopicRuntimePropertiesResult struct {
	// Value is the result of the request.
	TopicRuntimeProperties
}

type GetTopicRuntimePropertiesResponse struct {
	GetTopicRuntimePropertiesResult

	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type GetTopicRuntimePropertiesOptions struct {
	// For future expansion
}

// GetTopicRuntimeProperties gets runtime properties of a topic, like the SizeInBytes, or SubscriptionCount.
func (ac *Client) GetTopicRuntimeProperties(ctx context.Context, topicName string, options *GetTopicRuntimePropertiesOptions) (*GetTopicRuntimePropertiesResponse, error) {
	var atomResp *atom.TopicEnvelope
	resp, err := ac.em.Get(ctx, "/"+topicName, &atomResp)

	if err != nil {
		return nil, err
	}

	props, err := newTopicRuntimeProperties(&atomResp.Content.TopicDescription)

	if err != nil {
		return nil, err
	}

	return &GetTopicRuntimePropertiesResponse{
		RawResponse: resp,
		GetTopicRuntimePropertiesResult: GetTopicRuntimePropertiesResult{
			TopicRuntimeProperties: *props,
		},
	}, nil
}

type TopicItem struct {
	TopicProperties

	TopicName string
}

type ListTopicsResponse struct {
	// Items is the result of the request.
	Items []*TopicItem
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

// ListTopicsOptions can be used to configure the ListTopics method.
type ListTopicsOptions struct {
	// MaxPageSize is the maximum size of each page of results.
	MaxPageSize int32
}

// TopicsPager provides iteration over TopicProperties pages.
type TopicsPager struct {
	innerPager pagerFunc

	lastErr      error
	lastResponse *ListTopicsResponse
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *TopicsPager) NextPage(ctx context.Context) bool {
	p.lastResponse, p.lastErr = p.getNextPage(ctx)
	return p.lastResponse != nil
}

// PageResponse returns the current page.
func (p *TopicsPager) PageResponse() *ListTopicsResponse {
	return p.lastResponse
}

// Err returns the last error encountered while paging.
func (p *TopicsPager) Err() error {
	return p.lastErr
}

func (p *TopicsPager) getNextPage(ctx context.Context) (*ListTopicsResponse, error) {
	var feed *atom.TopicFeed
	resp, err := p.innerPager(ctx, &feed)

	if err != nil || feed == nil {
		return nil, err
	}

	var all []*TopicItem

	for _, env := range feed.Entries {
		props, err := newTopicProperties(&env.Content.TopicDescription)

		if err != nil {
			return nil, err
		}

		all = append(all, &TopicItem{
			TopicProperties: *props,
			TopicName:       env.Title,
		})
	}

	return &ListTopicsResponse{
		RawResponse: resp,
		Items:       all,
	}, nil
}

// ListTopics lists topics.
func (ac *Client) ListTopics(options *ListTopicsOptions) *TopicsPager {
	var pageSize int32

	if options != nil {
		pageSize = options.MaxPageSize
	}

	pagerFunc := ac.newPagerFunc("/$Resources/Topics", pageSize, topicFeedLen)

	return &TopicsPager{
		innerPager: pagerFunc,
	}
}

type TopicRuntimePropertiesItem struct {
	TopicRuntimeProperties

	TopicName string
}

type ListTopicsRuntimePropertiesResponse struct {
	// Items is the result of the request.
	Items []*TopicRuntimePropertiesItem
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

// ListTopicsRuntimePropertiesOptions can be used to configure the ListTopicsRuntimeProperties method.
type ListTopicsRuntimePropertiesOptions struct {
	// MaxPageSize is the maximum size of each page of results.
	MaxPageSize int32
}

// TopicRuntimePropertiesPager provides iteration over TopicRuntimeProperties pages.
type TopicRuntimePropertiesPager struct {
	innerPager   pagerFunc
	lastErr      error
	lastResponse *ListTopicsRuntimePropertiesResponse
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *TopicRuntimePropertiesPager) NextPage(ctx context.Context) bool {
	p.lastResponse, p.lastErr = p.getNextPage(ctx)
	return p.lastResponse != nil
}

// PageResponse returns the current page.
func (p *TopicRuntimePropertiesPager) PageResponse() *ListTopicsRuntimePropertiesResponse {
	return p.lastResponse
}

// Err returns the last error encountered while paging.
func (p *TopicRuntimePropertiesPager) Err() error {
	return p.lastErr
}

func (p *TopicRuntimePropertiesPager) getNextPage(ctx context.Context) (*ListTopicsRuntimePropertiesResponse, error) {
	var feed *atom.TopicFeed
	resp, err := p.innerPager(ctx, &feed)

	if err != nil || feed == nil {
		return nil, err
	}

	var all []*TopicRuntimePropertiesItem

	for _, entry := range feed.Entries {
		props, err := newTopicRuntimeProperties(&entry.Content.TopicDescription)

		if err != nil {
			return nil, err
		}

		all = append(all, &TopicRuntimePropertiesItem{
			TopicName:              entry.Title,
			TopicRuntimeProperties: *props,
		})
	}

	return &ListTopicsRuntimePropertiesResponse{
		RawResponse: resp,
		Items:       all,
	}, nil
}

// ListTopicsRuntimeProperties lists runtime properties for topics.
func (ac *Client) ListTopicsRuntimeProperties(options *ListTopicsRuntimePropertiesOptions) *TopicRuntimePropertiesPager {
	var pageSize int32

	if options != nil {
		pageSize = options.MaxPageSize
	}

	pagerFunc := ac.newPagerFunc("/$Resources/Topics", pageSize, topicFeedLen)

	return &TopicRuntimePropertiesPager{
		innerPager: pagerFunc,
	}
}

type UpdateTopicResult struct {
	TopicProperties
}

type UpdateTopicResponse struct {
	UpdateTopicResult

	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type UpdateTopicOptions struct {
	// For future expansion
}

// UpdateTopic updates an existing topic.
func (ac *Client) UpdateTopic(ctx context.Context, topicName string, properties TopicProperties, options *UpdateTopicOptions) (*UpdateTopicResponse, error) {
	newProps, resp, err := ac.createOrUpdateTopicImpl(ctx, topicName, &properties, false)

	if err != nil {
		return nil, err
	}

	return &UpdateTopicResponse{
		RawResponse: resp,
		UpdateTopicResult: UpdateTopicResult{
			TopicProperties: *newProps,
		},
	}, nil
}

type DeleteTopicResponse struct {
	// Value is the result of the request.
	Value *TopicProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type DeleteTopicOptions struct {
	// For future expansion
}

// DeleteTopic deletes a topic.
func (ac *Client) DeleteTopic(ctx context.Context, topicName string, options *DeleteTopicOptions) (*DeleteTopicResponse, error) {
	resp, err := ac.em.Delete(ctx, "/"+topicName)
	defer atom.CloseRes(ctx, resp)
	return &DeleteTopicResponse{
		RawResponse: resp,
	}, err
}

func (ac *Client) createOrUpdateTopicImpl(ctx context.Context, topicName string, props *TopicProperties, creating bool) (*TopicProperties, *http.Response, error) {
	if props == nil {
		props = &TopicProperties{}
	}

	env := newTopicEnvelope(props, ac.em.TokenProvider())

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

	var atomResp *atom.TopicEnvelope
	resp, err := ac.em.Put(ctx, "/"+topicName, env, &atomResp, mw...)

	if err != nil {
		return nil, nil, err
	}

	topicProps, err := newTopicProperties(&atomResp.Content.TopicDescription)

	if err != nil {
		return nil, nil, err
	}

	return topicProps, resp, nil
}

func newTopicEnvelope(props *TopicProperties, tokenProvider auth.TokenProvider) *atom.TopicEnvelope {
	desc := &atom.TopicDescription{
		DefaultMessageTimeToLive:            utils.DurationToStringPtr(props.DefaultMessageTimeToLive),
		MaxSizeInMegabytes:                  props.MaxSizeInMegabytes,
		RequiresDuplicateDetection:          props.RequiresDuplicateDetection,
		DuplicateDetectionHistoryTimeWindow: utils.DurationToStringPtr(props.DuplicateDetectionHistoryTimeWindow),
		EnableBatchedOperations:             props.EnableBatchedOperations,

		Status:             (*atom.EntityStatus)(props.Status),
		UserMetadata:       props.UserMetadata,
		SupportOrdering:    props.SupportOrdering,
		AutoDeleteOnIdle:   utils.DurationToStringPtr(props.AutoDeleteOnIdle),
		EnablePartitioning: props.EnablePartitioning,
	}

	return atom.WrapWithTopicEnvelope(desc)
}

func newTopicProperties(td *atom.TopicDescription) (*TopicProperties, error) {
	defaultMessageTimeToLive, err := utils.ISO8601StringToDuration(td.DefaultMessageTimeToLive)

	if err != nil {
		return nil, err
	}

	duplicateDetectionHistoryTimeWindow, err := utils.ISO8601StringToDuration(td.DuplicateDetectionHistoryTimeWindow)

	if err != nil {
		return nil, err
	}

	autoDeleteOnIdle, err := utils.ISO8601StringToDuration(td.AutoDeleteOnIdle)

	if err != nil {
		return nil, err
	}

	return &TopicProperties{
		MaxSizeInMegabytes:                  td.MaxSizeInMegabytes,
		RequiresDuplicateDetection:          td.RequiresDuplicateDetection,
		DefaultMessageTimeToLive:            defaultMessageTimeToLive,
		DuplicateDetectionHistoryTimeWindow: duplicateDetectionHistoryTimeWindow,
		EnableBatchedOperations:             td.EnableBatchedOperations,
		Status:                              (*EntityStatus)(td.Status),
		UserMetadata:                        td.UserMetadata,
		AutoDeleteOnIdle:                    autoDeleteOnIdle,
		EnablePartitioning:                  td.EnablePartitioning,
		SupportOrdering:                     td.SupportOrdering,
	}, nil
}

func newTopicRuntimeProperties(desc *atom.TopicDescription) (*TopicRuntimeProperties, error) {
	props := &TopicRuntimeProperties{
		SizeInBytes:           int64OrZero(desc.SizeInBytes),
		ScheduledMessageCount: int32OrZero(desc.CountDetails.ScheduledMessageCount),
		SubscriptionCount:     int32OrZero(desc.SubscriptionCount),
	}

	var err error

	if props.CreatedAt, err = atom.StringToTime(desc.CreatedAt); err != nil {
		return nil, err
	}

	if props.UpdatedAt, err = atom.StringToTime(desc.UpdatedAt); err != nil {
		return nil, err
	}

	if props.AccessedAt, err = atom.StringToTime(desc.AccessedAt); err != nil {
		return nil, err
	}

	return props, nil
}

func topicFeedLen(pv interface{}) int {
	topicFeed := pv.(**atom.TopicFeed)
	return len((*topicFeed).Entries)
}
