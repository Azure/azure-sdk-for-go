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

// TopicProperties represents the static properties of the topic.
type TopicProperties struct {
	// Name of the topic relative to the namespace base address.
	Name string

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

	// SupportsOrdering defines whether ordering needs to be maintained. If true, messages
	// sent to topic will be forwarded to the subscription, in order.
	SupportsOrdering *bool
}

// TopicRuntimeProperties represent dynamic properties of a topic, such as the ActiveMessageCount.
type TopicRuntimeProperties struct {
	// Name is the name of the topic.
	Name string

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

type AddTopicResponse struct {
	// Value is the result of the request.
	Value *TopicProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type GetTopicResponse struct {
	// Value is the result of the request.
	Value *TopicProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type GetTopicRuntimePropertiesResponse struct {
	// Value is the result of the request.
	Value *TopicRuntimeProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type ListTopicsResponse struct {
	// Value is the result of the request.
	Value []*TopicProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type ListTopicsRuntimePropertiesResponse struct {
	// Value is the result of the request.
	Value []*TopicRuntimeProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type UpdateTopicResponse struct {
	// Value is the result of the request.
	Value *TopicProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type DeleteTopicResponse struct {
	// Value is the result of the request.
	Value *TopicProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

// AddTopic creates a topic using defaults for all options.
func (ac *AdminClient) AddTopic(ctx context.Context, topicName string) (*AddTopicResponse, error) {
	return ac.AddTopicWithProperties(ctx, &TopicProperties{
		Name: topicName,
	})
}

// AddTopicWithProperties creates a topic with configurable properties
func (ac *AdminClient) AddTopicWithProperties(ctx context.Context, properties *TopicProperties) (*AddTopicResponse, error) {
	props, resp, err := ac.createOrUpdateTopicImpl(ctx, properties, true)

	if err != nil {
		return nil, err
	}

	return &AddTopicResponse{
		RawResponse: resp,
		Value:       props,
	}, nil
}

// GetTopic gets a topic by name.
func (ac *AdminClient) GetTopic(ctx context.Context, topicName string) (*GetTopicResponse, error) {
	var atomResp *atom.TopicEnvelope
	resp, err := ac.em.Get(ctx, "/"+topicName, &atomResp)

	if err != nil {
		return nil, err
	}

	props, err := newTopicProperties(atomResp.Title, &atomResp.Content.TopicDescription)

	if err != nil {
		return nil, err
	}

	return &GetTopicResponse{
		RawResponse: resp,
		Value:       props,
	}, nil
}

// GetTopicRuntimeProperties gets runtime properties of a topic, like the SizeInBytes, or SubscriptionCount.
func (ac *AdminClient) GetTopicRuntimeProperties(ctx context.Context, topicName string) (*GetTopicRuntimePropertiesResponse, error) {
	var atomResp *atom.TopicEnvelope
	resp, err := ac.em.Get(ctx, "/"+topicName, &atomResp)

	if err != nil {
		return nil, err
	}

	return &GetTopicRuntimePropertiesResponse{
		RawResponse: resp,
		Value:       newTopicRuntimeProperties(atomResp.Title, &atomResp.Content.TopicDescription),
	}, nil
}

// ListTopicsOptions can be used to configure the ListTopics method.
type ListTopicsOptions struct {
	// Top is the maximum size of each page of results.
	Top int
	// Skip is the starting index for the paging operation.
	Skip int
}

// TopicPropertiesPager provides iteration over ListTopicProperties pages.
type TopicPropertiesPager interface {
	// NextPage returns true if the pager advanced to the next page.
	// Returns false if there are no more pages or an error occurred.
	NextPage(context.Context) bool

	// PageResponse returns the current TopicProperties.
	PageResponse() *ListTopicsResponse

	// Err returns the last error encountered while paging.
	Err() error
}

// ListTopics lists topics.
func (ac *AdminClient) ListTopics(options *ListTopicsOptions) TopicPropertiesPager {
	var pageSize int
	var skip int

	if options != nil {
		skip = options.Skip
		pageSize = options.Top
	}

	return &topicPropertiesPager{
		innerPager: ac.getTopicPager(pageSize, skip),
	}
}

// ListTopicsRuntimePropertiesOptions can be used to configure the ListTopicsRuntimeProperties method.
type ListTopicsRuntimePropertiesOptions struct {
	// Top is the maximum size of each page of results.
	Top int
	// Skip is the starting index for the paging operation.
	Skip int
}

// TopicRuntimePropertiesPager provides iteration over ListTopicRuntimeProperties pages.
type TopicRuntimePropertiesPager interface {
	// NextPage returns true if the pager advanced to the next page.
	// Returns false if there are no more pages or an error occurred.
	NextPage(context.Context) bool

	// PageResponse returns the current TopicRuntimeProperties.
	PageResponse() *ListTopicsRuntimePropertiesResponse

	// Err returns the last error encountered while paging.
	Err() error
}

// ListTopicsRuntimeProperties lists runtime properties for topics.
func (ac *AdminClient) ListTopicsRuntimeProperties(options *ListTopicsRuntimePropertiesOptions) TopicRuntimePropertiesPager {
	var pageSize int
	var skip int

	if options != nil {
		skip = options.Skip
		pageSize = options.Top
	}

	return &topicRuntimePropertiesPager{
		innerPager: ac.getTopicPager(pageSize, skip),
	}
}

// TopicExists checks if a topic exists.
// Returns true if the topic is found
// (false, nil) if the topic is not found
// (false, err) if an error occurred while trying to check if the topic exists.
func (ac *AdminClient) TopicExists(ctx context.Context, topicName string) (bool, error) {
	_, err := ac.GetTopic(ctx, topicName)

	if err == nil {
		return true, nil
	}

	if atom.NotFound(err) {
		return false, nil
	}

	return false, err
}

// UpdateTopic updates an existing topic.
func (ac *AdminClient) UpdateTopic(ctx context.Context, properties *TopicProperties) (*UpdateTopicResponse, error) {
	newProps, resp, err := ac.createOrUpdateTopicImpl(ctx, properties, false)

	if err != nil {
		return nil, err
	}

	return &UpdateTopicResponse{
		RawResponse: resp,
		Value:       newProps,
	}, nil
}

// DeleteTopic deletes a topic.
func (ac *AdminClient) DeleteTopic(ctx context.Context, topicName string) (*DeleteTopicResponse, error) {
	resp, err := ac.em.Delete(ctx, "/"+topicName)
	defer atom.CloseRes(ctx, resp)
	return &DeleteTopicResponse{
		RawResponse: resp,
	}, err
}

func (ac *AdminClient) getTopicPager(top int, skip int) topicFeedPagerFunc {
	return func(ctx context.Context) (*atom.TopicFeed, *http.Response, error) {
		url := "/$Resources/Topics?"
		if top > 0 {
			url += fmt.Sprintf("&$top=%d", top)
		}

		if skip > 0 {
			url += fmt.Sprintf("&$skip=%d", skip)
		}

		var atomResp *atom.TopicFeed
		resp, err := ac.em.Get(ctx, url, &atomResp)

		if err != nil {
			return nil, nil, err
		}

		if len(atomResp.Entries) == 0 {
			return nil, nil, nil
		}

		skip += len(atomResp.Entries)
		return atomResp, resp, nil
	}
}

func (ac *AdminClient) createOrUpdateTopicImpl(ctx context.Context, props *TopicProperties, creating bool) (*TopicProperties, *http.Response, error) {
	if props == nil {
		return nil, nil, errors.New("properties are required and cannot be nil")
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
	resp, err := ac.em.Put(ctx, "/"+props.Name, env, &atomResp, mw...)

	if err != nil {
		return nil, nil, err
	}

	topicProps, err := newTopicProperties(props.Name, &atomResp.Content.TopicDescription)

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
		SupportOrdering:    props.SupportsOrdering,
		AutoDeleteOnIdle:   utils.DurationToStringPtr(props.AutoDeleteOnIdle),
		EnablePartitioning: props.EnablePartitioning,
	}

	return atom.WrapWithTopicEnvelope(desc)
}

func newTopicProperties(name string, td *atom.TopicDescription) (*TopicProperties, error) {
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
		Name:                                name,
		MaxSizeInMegabytes:                  td.MaxSizeInMegabytes,
		RequiresDuplicateDetection:          td.RequiresDuplicateDetection,
		DefaultMessageTimeToLive:            defaultMessageTimeToLive,
		DuplicateDetectionHistoryTimeWindow: duplicateDetectionHistoryTimeWindow,
		EnableBatchedOperations:             td.EnableBatchedOperations,
		Status:                              (*EntityStatus)(td.Status),
		AutoDeleteOnIdle:                    autoDeleteOnIdle,
		EnablePartitioning:                  td.EnablePartitioning,
		SupportsOrdering:                    td.SupportOrdering,
	}, nil
}

func newTopicRuntimeProperties(name string, desc *atom.TopicDescription) *TopicRuntimeProperties {
	return &TopicRuntimeProperties{
		Name:                  name,
		SizeInBytes:           int64OrZero(desc.SizeInBytes),
		CreatedAt:             dateTimeToTime(desc.CreatedAt),
		UpdatedAt:             dateTimeToTime(desc.UpdatedAt),
		AccessedAt:            dateTimeToTime(desc.AccessedAt),
		ScheduledMessageCount: int32OrZero(desc.CountDetails.ScheduledMessageCount),
		SubscriptionCount:     int32OrZero(desc.SubscriptionCount),
	}
}

// topicPropertiesPager provides iteration over TopicProperties pages.
type topicPropertiesPager struct {
	innerPager topicFeedPagerFunc

	lastErr      error
	lastResponse *ListTopicsResponse
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *topicPropertiesPager) NextPage(ctx context.Context) bool {
	p.lastResponse, p.lastErr = p.getNextPage(ctx)
	return p.lastResponse != nil
}

// PageResponse returns the current page.
func (p *topicPropertiesPager) PageResponse() *ListTopicsResponse {
	return p.lastResponse
}

// Err returns the last error encountered while paging.
func (p *topicPropertiesPager) Err() error {
	return p.lastErr
}

func (p *topicPropertiesPager) getNextPage(ctx context.Context) (*ListTopicsResponse, error) {
	feed, resp, err := p.innerPager(ctx)

	if err != nil || feed == nil {
		return nil, err
	}

	var all []*TopicProperties

	for _, env := range feed.Entries {
		props, err := newTopicProperties(env.Title, &env.Content.TopicDescription)

		if err != nil {
			return nil, err
		}

		all = append(all, props)
	}

	return &ListTopicsResponse{
		RawResponse: resp,
		Value:       all,
	}, nil
}

type topicFeedPagerFunc func(ctx context.Context) (*atom.TopicFeed, *http.Response, error)

// topicRuntimePropertiesPager provides iteration over TopicRuntimeProperties pages.
type topicRuntimePropertiesPager struct {
	innerPager   topicFeedPagerFunc
	lastErr      error
	lastResponse *ListTopicsRuntimePropertiesResponse
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *topicRuntimePropertiesPager) NextPage(ctx context.Context) bool {
	p.lastResponse, p.lastErr = p.getNextPage(ctx)
	return p.lastResponse != nil
}

func (p *topicRuntimePropertiesPager) getNextPage(ctx context.Context) (*ListTopicsRuntimePropertiesResponse, error) {
	feed, resp, err := p.innerPager(ctx)

	if err != nil || feed == nil {
		return nil, err
	}

	var all []*TopicRuntimeProperties

	for _, entry := range feed.Entries {
		all = append(all, newTopicRuntimeProperties(entry.Title, &entry.Content.TopicDescription))
	}

	return &ListTopicsRuntimePropertiesResponse{
		RawResponse: resp,
		Value:       all,
	}, nil
}

// PageResponse returns the current page.
func (p *topicRuntimePropertiesPager) PageResponse() *ListTopicsRuntimePropertiesResponse {
	return p.lastResponse
}

// Err returns the last error encountered while paging.
func (p *topicRuntimePropertiesPager) Err() error {
	return p.lastErr
}
