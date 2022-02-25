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

// QueueProperties represents the static properties of the queue.
type QueueProperties struct {
	// LockDuration is the duration a message is locked when using the PeekLock receive mode.
	// Default is 1 minute.
	LockDuration *time.Duration

	// MaxSizeInMegabytes - The maximum size of the queue in megabytes, which is the size of memory
	// allocated for the queue.
	// Default is 1024.
	MaxSizeInMegabytes *int32

	// RequiresDuplicateDetection indicates if this queue requires duplicate detection.
	RequiresDuplicateDetection *bool

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

	// DuplicateDetectionHistoryTimeWindow is the duration of duplicate detection history.
	// Default value is 10 minutes.
	DuplicateDetectionHistoryTimeWindow *time.Duration

	// MaxDeliveryCount is the maximum amount of times a message can be delivered before it is automatically
	// sent to the dead letter queue.
	// Default value is 10.
	MaxDeliveryCount *int32

	// EnableBatchedOperations indicates whether server-side batched operations are enabled.
	EnableBatchedOperations *bool

	// Status is the current status of the queue.
	Status *EntityStatus

	// AutoDeleteOnIdle is the idle interval after which the queue is automatically deleted.
	AutoDeleteOnIdle *time.Duration

	// EnablePartitioning indicates whether the queue is to be partitioned across multiple message brokers.
	EnablePartitioning *bool

	// ForwardTo is the name of the recipient entity to which all the messages sent to the queue
	// are forwarded to.
	ForwardTo *string

	// ForwardDeadLetteredMessagesTo is the absolute URI of the entity to forward dead letter messages
	ForwardDeadLetteredMessagesTo *string

	// UserMetadata is custom metadata that user can associate with the queue.
	UserMetadata *string
}

// QueueRuntimeProperties represent dynamic properties of a queue, such as the ActiveMessageCount.
type QueueRuntimeProperties struct {
	// SizeInBytes - The size of the queue, in bytes.
	SizeInBytes int64

	// CreatedAt is when the entity was created.
	CreatedAt time.Time

	// UpdatedAt is when the entity was last updated.
	UpdatedAt time.Time

	// AccessedAt is when the entity was last updated.
	AccessedAt time.Time

	// TotalMessageCount is the number of messages in the queue.
	TotalMessageCount int64

	// ActiveMessageCount is the number of active messages in the entity.
	ActiveMessageCount int32

	// DeadLetterMessageCount is the number of dead-lettered messages in the entity.
	DeadLetterMessageCount int32

	// ScheduledMessageCount is the number of messages that are scheduled to be enqueued.
	ScheduledMessageCount int32

	// TransferDeadLetterMessageCount is the number of messages transfer-messages which are dead-lettered
	// into transfer-dead-letter subqueue.
	TransferDeadLetterMessageCount int32

	// TransferMessageCount is the number of messages which are yet to be transferred/forwarded to destination entity.
	TransferMessageCount int32
}

type CreateQueueOptions struct {
	// for future expansion
}

type CreateQueueResult struct {
	QueueProperties
}

type CreateQueueResponse struct {
	CreateQueueResult

	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

// CreateQueue creates a queue with configurable properties.
func (ac *Client) CreateQueue(ctx context.Context, queueName string, properties *QueueProperties, options *CreateQueueOptions) (*CreateQueueResponse, error) {
	newProps, resp, err := ac.createOrUpdateQueueImpl(ctx, queueName, properties, true)

	if err != nil {
		return nil, err
	}

	return &CreateQueueResponse{
		RawResponse: resp,
		CreateQueueResult: CreateQueueResult{
			QueueProperties: *newProps,
		},
	}, nil
}

type UpdateQueueResult struct {
	QueueProperties
}

type UpdateQueueResponse struct {
	UpdateQueueResult

	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type UpdateQueueOptions struct {
	// for future expansion
}

// UpdateQueue updates an existing queue.
func (ac *Client) UpdateQueue(ctx context.Context, queueName string, properties QueueProperties, options *UpdateQueueOptions) (*UpdateQueueResponse, error) {
	newProps, resp, err := ac.createOrUpdateQueueImpl(ctx, queueName, &properties, false)

	if err != nil {
		return nil, err
	}

	return &UpdateQueueResponse{
		RawResponse: resp,
		UpdateQueueResult: UpdateQueueResult{
			QueueProperties: *newProps,
		},
	}, err
}

type GetQueueResult struct {
	QueueProperties
}

type GetQueueResponse struct {
	GetQueueResult

	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type GetQueueOptions struct {
	// For future expansion
}

// GetQueue gets a queue by name.
func (ac *Client) GetQueue(ctx context.Context, queueName string, options *GetQueueOptions) (*GetQueueResponse, error) {
	var atomResp *atom.QueueEnvelope
	resp, err := ac.em.Get(ctx, "/"+queueName, &atomResp)

	if err != nil {
		return nil, err
	}

	props, err := newQueueProperties(&atomResp.Content.QueueDescription)

	if err != nil {
		return nil, err
	}

	return &GetQueueResponse{
		RawResponse: resp,
		GetQueueResult: GetQueueResult{
			QueueProperties: *props,
		},
	}, nil
}

type GetQueueRuntimePropertiesResult struct {
	QueueRuntimeProperties
}

type GetQueueRuntimePropertiesResponse struct {
	GetQueueRuntimePropertiesResult

	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type GetQueueRuntimePropertiesOptions struct {
	// For future expansion
}

// GetQueueRuntimeProperties gets runtime properties of a queue, like the SizeInBytes, or ActiveMessageCount.
func (ac *Client) GetQueueRuntimeProperties(ctx context.Context, queueName string, options *GetQueueRuntimePropertiesOptions) (*GetQueueRuntimePropertiesResponse, error) {
	var atomResp *atom.QueueEnvelope
	resp, err := ac.em.Get(ctx, "/"+queueName, &atomResp)

	if err != nil {
		return nil, err
	}

	props, err := newQueueRuntimeProperties(&atomResp.Content.QueueDescription)

	if err != nil {
		return nil, err
	}

	return &GetQueueRuntimePropertiesResponse{
		RawResponse: resp,
		GetQueueRuntimePropertiesResult: GetQueueRuntimePropertiesResult{
			QueueRuntimeProperties: *props,
		},
	}, nil
}

type DeleteQueueResponse struct {
	RawResponse *http.Response
}

type DeleteQueueOptions struct {
	// for future expansion
}

// DeleteQueue deletes a queue.
func (ac *Client) DeleteQueue(ctx context.Context, queueName string, options *DeleteQueueOptions) (*DeleteQueueResponse, error) {
	resp, err := ac.em.Delete(ctx, "/"+queueName)
	defer atom.CloseRes(ctx, resp)
	return &DeleteQueueResponse{RawResponse: resp}, err
}

// ListQueuesOptions can be used to configure the ListQueues method.
type ListQueuesOptions struct {
	// MaxPageSize is the maximum size of each page of results.
	MaxPageSize int32
}

type ListQueuesResult struct {
	Items []*QueueItem
}

type ListQueuesResponse struct {
	ListQueuesResult
	RawResponse *http.Response
}

type QueueItem struct {
	QueueName string
	QueueProperties
}

// ListQueues lists queues.
func (ac *Client) ListQueues(options *ListQueuesOptions) *QueuePager {
	var pageSize int32

	if options != nil {
		pageSize = options.MaxPageSize
	}

	return &QueuePager{
		innerPager: ac.newPagerFunc("/$Resources/Queues", pageSize, queueFeedLen),
	}
}

// ListQueuesRuntimePropertiesOptions can be used to configure the ListQueuesRuntimeProperties method.
type ListQueuesRuntimePropertiesOptions struct {
	// MaxPageSize is the maximum size of each page of results.
	MaxPageSize int32
}

type ListQueuesRuntimePropertiesResponse struct {
	Items       []*QueueRuntimePropertiesItem
	RawResponse *http.Response
}

type QueueRuntimePropertiesItem struct {
	QueueName string
	QueueRuntimeProperties
}

// ListQueuesRuntimeProperties lists runtime properties for queues.
func (ac *Client) ListQueuesRuntimeProperties(options *ListQueuesRuntimePropertiesOptions) *QueueRuntimePropertiesPager {
	var pageSize int32

	if options != nil {
		pageSize = options.MaxPageSize
	}

	return &QueueRuntimePropertiesPager{
		innerPager: ac.newPagerFunc("/$Resources/Queues", pageSize, queueFeedLen),
	}
}

func (ac *Client) createOrUpdateQueueImpl(ctx context.Context, queueName string, props *QueueProperties, creating bool) (*QueueProperties, *http.Response, error) {
	if props == nil {
		props = &QueueProperties{}
	}

	env, mw := newQueueEnvelope(props, ac.em.TokenProvider())

	if !creating {
		// an update requires the entity to already exist.
		mw = append(mw, func(next atom.RestHandler) atom.RestHandler {
			return func(ctx context.Context, req *http.Request) (*http.Response, error) {
				req.Header.Set("If-Match", "*")
				return next(ctx, req)
			}
		})
	}

	var atomResp *atom.QueueEnvelope
	resp, err := ac.em.Put(ctx, "/"+queueName, env, &atomResp, mw...)

	if err != nil {
		return nil, nil, err
	}

	newProps, err := newQueueProperties(&atomResp.Content.QueueDescription)

	if err != nil {
		return nil, nil, err
	}

	return newProps, resp, nil
}

// QueuePager provides iteration over ListQueues pages.
type QueuePager struct {
	innerPager pagerFunc

	lastErr      error
	lastResponse *ListQueuesResponse
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *QueuePager) NextPage(ctx context.Context) bool {
	p.lastResponse, p.lastErr = p.getNextPage(ctx)
	return p.lastResponse != nil
}

// PageResponse returns the current page.
func (p *QueuePager) PageResponse() *ListQueuesResponse {
	return p.lastResponse
}

// Err returns the last error encountered while paging.
func (p *QueuePager) Err() error {
	return p.lastErr
}

func (p *QueuePager) getNextPage(ctx context.Context) (*ListQueuesResponse, error) {
	var feed *atom.QueueFeed
	resp, err := p.innerPager(ctx, &feed)

	if err != nil || feed == nil {
		return nil, err
	}

	var all []*QueueItem

	for _, env := range feed.Entries {
		queueName := env.Title
		props, err := newQueueProperties(&env.Content.QueueDescription)

		if err != nil {
			return nil, err
		}

		all = append(all, &QueueItem{
			QueueName:       queueName,
			QueueProperties: *props,
		})
	}

	return &ListQueuesResponse{
		RawResponse: resp,
		ListQueuesResult: ListQueuesResult{
			Items: all,
		},
	}, nil
}

// QueueRuntimePropertiesPager provides iteration over ListQueueRuntimeProperties pages.
type QueueRuntimePropertiesPager struct {
	innerPager   pagerFunc
	lastErr      error
	lastResponse *ListQueuesRuntimePropertiesResponse
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *QueueRuntimePropertiesPager) NextPage(ctx context.Context) bool {
	p.lastResponse, p.lastErr = p.getNextPage(ctx)
	return p.lastResponse != nil
}

// PageResponse returns the current page.
func (p *QueueRuntimePropertiesPager) PageResponse() *ListQueuesRuntimePropertiesResponse {
	return p.lastResponse
}

// Err returns the last error encountered while paging.
func (p *QueueRuntimePropertiesPager) Err() error {
	return p.lastErr
}

func (p *QueueRuntimePropertiesPager) getNextPage(ctx context.Context) (*ListQueuesRuntimePropertiesResponse, error) {
	var feed *atom.QueueFeed
	resp, err := p.innerPager(ctx, &feed)

	if err != nil || feed == nil {
		return nil, err
	}

	var all []*QueueRuntimePropertiesItem

	for _, entry := range feed.Entries {
		rt, err := newQueueRuntimeProperties(&entry.Content.QueueDescription)

		if err != nil {
			return nil, err
		}

		all = append(all, &QueueRuntimePropertiesItem{
			QueueName:              entry.Title,
			QueueRuntimeProperties: *rt,
		})
	}

	return &ListQueuesRuntimePropertiesResponse{
		RawResponse: resp,
		Items:       all,
	}, nil
}

func newQueueEnvelope(props *QueueProperties, tokenProvider auth.TokenProvider) (*atom.QueueEnvelope, []atom.MiddlewareFunc) {
	qpr := &atom.QueueDescription{
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
		UserMetadata:                        props.UserMetadata,
	}

	return atom.WrapWithQueueEnvelope(qpr, tokenProvider)
}

func newQueueProperties(desc *atom.QueueDescription) (*QueueProperties, error) {
	lockDuration, err := utils.ISO8601StringToDuration(desc.LockDuration)

	if err != nil {
		return nil, err
	}

	defaultMessageTimeToLive, err := utils.ISO8601StringToDuration(desc.DefaultMessageTimeToLive)

	if err != nil {
		return nil, err
	}

	duplicateDetectionHistoryTimeWindow, err := utils.ISO8601StringToDuration(desc.DuplicateDetectionHistoryTimeWindow)

	if err != nil {
		return nil, err
	}

	autoDeleteOnIdle, err := utils.ISO8601StringToDuration(desc.AutoDeleteOnIdle)

	if err != nil {
		return nil, err
	}

	queuePropsResult := &QueueProperties{
		LockDuration:                        lockDuration,
		MaxSizeInMegabytes:                  desc.MaxSizeInMegabytes,
		RequiresDuplicateDetection:          desc.RequiresDuplicateDetection,
		RequiresSession:                     desc.RequiresSession,
		DefaultMessageTimeToLive:            defaultMessageTimeToLive,
		DeadLetteringOnMessageExpiration:    desc.DeadLetteringOnMessageExpiration,
		DuplicateDetectionHistoryTimeWindow: duplicateDetectionHistoryTimeWindow,
		MaxDeliveryCount:                    desc.MaxDeliveryCount,
		EnableBatchedOperations:             desc.EnableBatchedOperations,
		Status:                              (*EntityStatus)(desc.Status),
		AutoDeleteOnIdle:                    autoDeleteOnIdle,
		EnablePartitioning:                  desc.EnablePartitioning,
		ForwardTo:                           desc.ForwardTo,
		ForwardDeadLetteredMessagesTo:       desc.ForwardDeadLetteredMessagesTo,
		UserMetadata:                        desc.UserMetadata,
	}

	return queuePropsResult, nil
}

func newQueueRuntimeProperties(desc *atom.QueueDescription) (*QueueRuntimeProperties, error) {
	qrt := &QueueRuntimeProperties{
		SizeInBytes:                    int64OrZero(desc.SizeInBytes),
		TotalMessageCount:              int64OrZero(desc.MessageCount),
		ActiveMessageCount:             int32OrZero(desc.CountDetails.ActiveMessageCount),
		DeadLetterMessageCount:         int32OrZero(desc.CountDetails.DeadLetterMessageCount),
		ScheduledMessageCount:          int32OrZero(desc.CountDetails.ScheduledMessageCount),
		TransferDeadLetterMessageCount: int32OrZero(desc.CountDetails.TransferDeadLetterMessageCount),
		TransferMessageCount:           int32OrZero(desc.CountDetails.TransferMessageCount),
	}

	var err error

	if qrt.CreatedAt, err = atom.StringToTime(desc.CreatedAt); err != nil {
		return nil, err
	}

	if qrt.UpdatedAt, err = atom.StringToTime(desc.UpdatedAt); err != nil {
		return nil, err
	}

	if qrt.AccessedAt, err = atom.StringToTime(desc.AccessedAt); err != nil {
		return nil, err
	}

	return qrt, nil
}

func int32OrZero(i *int32) int32 {
	if i == nil {
		return 0
	}

	return *i
}

func int64OrZero(i *int64) int64 {
	if i == nil {
		return 0
	}

	return *i
}

func queueFeedLen(v interface{}) int {
	feed := v.(**atom.QueueFeed)
	return len((*feed).Entries)
}
