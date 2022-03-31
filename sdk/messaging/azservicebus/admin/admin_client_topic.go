// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
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
	DefaultMessageTimeToLive *string

	// DuplicateDetectionHistoryTimeWindow is the duration of duplicate detection history.
	// Default value is 10 minutes.
	DuplicateDetectionHistoryTimeWindow *string

	// EnableBatchedOperations indicates whether server-side batched operations are enabled.
	EnableBatchedOperations *bool

	// Status is the current status of the topic.
	Status *EntityStatus

	// AutoDeleteOnIdle is the idle interval after which the topic is automatically deleted.
	AutoDeleteOnIdle *string

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

// CreateTopicResponse contains response fields for Client.CreateTopic
type CreateTopicResponse struct {
	TopicProperties
}

// CreateTopicOptions contains optional parameters for Client.CreateTopic
type CreateTopicOptions struct {
	// For future expansion
}

// CreateTopic creates a topic using defaults for all options.
func (ac *Client) CreateTopic(ctx context.Context, topicName string, properties *TopicProperties, options *CreateTopicOptions) (CreateTopicResponse, error) {
	newProps, _, err := ac.createOrUpdateTopicImpl(ctx, topicName, properties, true)

	if err != nil {
		return CreateTopicResponse{}, err
	}

	return CreateTopicResponse{
		TopicProperties: *newProps,
	}, nil
}

// GetTopicResponse contains response fields for Client.GetTopic
type GetTopicResponse struct {
	TopicProperties
}

// GetTopicOptions contains optional parameters for Client.GetTopic
type GetTopicOptions struct {
	// For future expansion
}

// GetTopic gets a topic by name.
// If the entity does not exist this function will return a nil GetTopicResponse and a nil error.
func (ac *Client) GetTopic(ctx context.Context, topicName string, options *GetTopicOptions) (*GetTopicResponse, error) {
	var atomResp *atom.TopicEnvelope
	_, err := ac.em.Get(ctx, "/"+topicName, &atomResp)

	if err != nil {
		if errors.Is(err, atom.ErrFeedEmpty) {
			return nil, nil
		}

		return nil, err
	}

	props, err := newTopicProperties(&atomResp.Content.TopicDescription)

	if err != nil {
		return nil, err
	}

	return &GetTopicResponse{
		TopicProperties: *props,
	}, nil
}

// GetTopicRuntimePropertiesResponse contains the result for Client.GetTopicRuntimeProperties
type GetTopicRuntimePropertiesResponse struct {
	// Value is the result of the request.
	TopicRuntimeProperties
}

// GetTopicRuntimePropertiesOptions contains optional parameters for Client.GetTopicRuntimeProperties
type GetTopicRuntimePropertiesOptions struct {
	// For future expansion
}

// GetTopicRuntimeProperties gets runtime properties of a topic, like the SizeInBytes, or SubscriptionCount.
// If the entity does not exist this function will return a nil GetTopicRuntimePropertiesResponse and a nil error.
func (ac *Client) GetTopicRuntimeProperties(ctx context.Context, topicName string, options *GetTopicRuntimePropertiesOptions) (*GetTopicRuntimePropertiesResponse, error) {
	var atomResp *atom.TopicEnvelope
	_, err := ac.em.Get(ctx, "/"+topicName, &atomResp)

	if err != nil {
		if errors.Is(err, atom.ErrFeedEmpty) {
			return nil, nil
		}

		return nil, err
	}

	props, err := newTopicRuntimeProperties(&atomResp.Content.TopicDescription)

	if err != nil {
		return nil, err
	}

	return &GetTopicRuntimePropertiesResponse{
		TopicRuntimeProperties: *props,
	}, nil
}

// TopicItem is the data returned by the Client.ListTopics pager
type TopicItem struct {
	TopicProperties

	TopicName string
}

// ListTopicsResponse contains response fields for the Client.PageResponse method
type ListTopicsResponse struct {
	// Items is the result of the request.
	Items []*TopicItem
}

// ListTopicsOptions can be used to configure the ListTopics method.
type ListTopicsOptions struct {
	// MaxPageSize is the maximum size of each page of results.
	MaxPageSize int32
}

// ListTopicsPager provides iteration over ListTopics pages.
type ListTopicsPager struct {
	innerPager pagerFunc
	done       bool
}

func (p *ListTopicsPager) More() bool {
	return !p.done
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *ListTopicsPager) NextPage(ctx context.Context) (ListTopicsResponse, error) {
	resp, err := p.getNextPage(ctx)

	if resp == nil {
		p.done = true
		return ListTopicsResponse{}, nil
	}

	return *resp, err
}

func (p *ListTopicsPager) getNextPage(ctx context.Context) (*ListTopicsResponse, error) {
	var feed *atom.TopicFeed
	_, err := p.innerPager(ctx, &feed)

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
		Items: all,
	}, nil
}

// ListTopics lists topics.
func (ac *Client) ListTopics(options *ListTopicsOptions) *ListTopicsPager {
	var pageSize int32

	if options != nil {
		pageSize = options.MaxPageSize
	}

	pagerFunc := ac.newPagerFunc("/$Resources/Topics", pageSize, topicFeedLen)

	return &ListTopicsPager{
		innerPager: pagerFunc,
	}
}

// TopicRuntimePropertiesItem contains fields for the Client.ListTopicsRuntimeProperties method
type TopicRuntimePropertiesItem struct {
	TopicRuntimeProperties

	TopicName string
}

// ListTopicsRuntimePropertiesResponse contains response fields for TopicRuntimePropertiesPager.PageResponse
type ListTopicsRuntimePropertiesResponse struct {
	// Items is the result of the request.
	Items []*TopicRuntimePropertiesItem
}

// ListTopicsRuntimePropertiesOptions can be used to configure the ListTopicsRuntimeProperties method.
type ListTopicsRuntimePropertiesOptions struct {
	// MaxPageSize is the maximum size of each page of results.
	MaxPageSize int32
}

// ListTopicsRuntimePropertiesPager provides iteration over ListTopicsRuntimeProperties pages.
type ListTopicsRuntimePropertiesPager struct {
	innerPager pagerFunc
	done       bool
}

// More returns true if there are more pages.
func (p *ListTopicsRuntimePropertiesPager) More() bool {
	return !p.done
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *ListTopicsRuntimePropertiesPager) NextPage(ctx context.Context) (ListTopicsRuntimePropertiesResponse, error) {
	resp, err := p.getNextPage(ctx)

	if resp == nil {
		p.done = true
		return ListTopicsRuntimePropertiesResponse{}, nil
	}

	return *resp, err
}

func (p *ListTopicsRuntimePropertiesPager) getNextPage(ctx context.Context) (*ListTopicsRuntimePropertiesResponse, error) {
	var feed *atom.TopicFeed
	_, err := p.innerPager(ctx, &feed)

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
		Items: all,
	}, nil
}

// ListTopicsRuntimeProperties lists runtime properties for topics.
func (ac *Client) ListTopicsRuntimeProperties(options *ListTopicsRuntimePropertiesOptions) *ListTopicsRuntimePropertiesPager {
	var pageSize int32

	if options != nil {
		pageSize = options.MaxPageSize
	}

	pagerFunc := ac.newPagerFunc("/$Resources/Topics", pageSize, topicFeedLen)

	return &ListTopicsRuntimePropertiesPager{
		innerPager: pagerFunc,
	}
}

// UpdateTopicResponse contains response fields for Client.UpdateTopic
type UpdateTopicResponse struct {
	TopicProperties
}

// UpdateTopicOptions contains optional parameters for Client.UpdateTopic
type UpdateTopicOptions struct {
	// For future expansion
}

// UpdateTopic updates an existing topic.
func (ac *Client) UpdateTopic(ctx context.Context, topicName string, properties TopicProperties, options *UpdateTopicOptions) (UpdateTopicResponse, error) {
	newProps, _, err := ac.createOrUpdateTopicImpl(ctx, topicName, &properties, false)

	if err != nil {
		return UpdateTopicResponse{}, err
	}

	return UpdateTopicResponse{
		TopicProperties: *newProps,
	}, nil
}

// DeleteTopicResponse contains the response fields for Client.DeleteTopic
type DeleteTopicResponse struct {
	// Value is the result of the request.
	Value *TopicProperties
}

// DeleteTopicOptions contains optional parameters for Client.DeleteTopic
type DeleteTopicOptions struct {
	// For future expansion
}

// DeleteTopic deletes a topic.
func (ac *Client) DeleteTopic(ctx context.Context, topicName string, options *DeleteTopicOptions) (DeleteTopicResponse, error) {
	resp, err := ac.em.Delete(ctx, "/"+topicName)
	defer atom.CloseRes(ctx, resp)
	return DeleteTopicResponse{}, err
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
		DefaultMessageTimeToLive:            props.DefaultMessageTimeToLive,
		MaxSizeInMegabytes:                  props.MaxSizeInMegabytes,
		RequiresDuplicateDetection:          props.RequiresDuplicateDetection,
		DuplicateDetectionHistoryTimeWindow: props.DuplicateDetectionHistoryTimeWindow,
		EnableBatchedOperations:             props.EnableBatchedOperations,

		Status:             (*atom.EntityStatus)(props.Status),
		UserMetadata:       props.UserMetadata,
		SupportOrdering:    props.SupportOrdering,
		AutoDeleteOnIdle:   props.AutoDeleteOnIdle,
		EnablePartitioning: props.EnablePartitioning,
	}

	return atom.WrapWithTopicEnvelope(desc)
}

func newTopicProperties(td *atom.TopicDescription) (*TopicProperties, error) {
	return &TopicProperties{
		MaxSizeInMegabytes:                  td.MaxSizeInMegabytes,
		RequiresDuplicateDetection:          td.RequiresDuplicateDetection,
		DefaultMessageTimeToLive:            td.DefaultMessageTimeToLive,
		DuplicateDetectionHistoryTimeWindow: td.DuplicateDetectionHistoryTimeWindow,
		EnableBatchedOperations:             td.EnableBatchedOperations,
		Status:                              (*EntityStatus)(td.Status),
		UserMetadata:                        td.UserMetadata,
		AutoDeleteOnIdle:                    td.AutoDeleteOnIdle,
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
