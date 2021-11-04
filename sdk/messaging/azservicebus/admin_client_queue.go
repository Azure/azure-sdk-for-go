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
	"github.com/Azure/go-autorest/autorest/date"
)

// QueueProperties represents the static properties of the queue.
type QueueProperties struct {
	// Name of the queue relative to the namespace base address.
	Name string

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
}

// QueueRuntimeProperties represent dynamic properties of a queue, such as the ActiveMessageCount.
type QueueRuntimeProperties struct {
	// Name is the name of the queue.
	Name string

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

type AddQueueResponse struct {
	// Value is the result of the request.
	Value *QueueProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type AddQueueWithPropertiesResponse struct {
	// Value is the result of the request.
	Value *QueueProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type GetQueueResponse struct {
	// Value is the result of the request.
	Value *QueueProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type UpdateQueueResponse struct {
	// Value is the result of the request.
	Value *QueueProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type GetQueueRuntimePropertiesResponse struct {
	// Value is the result of the request.
	Value *QueueRuntimeProperties
	// RawResponse is the *http.Response for the request.
	RawResponse *http.Response
}

type DeleteQueueResponse struct {
	RawResponse *http.Response
}

// AddQueue creates a queue using defaults for all options.
func (ac *AdminClient) AddQueue(ctx context.Context, queueName string) (*AddQueueResponse, error) {
	return ac.AddQueueWithProperties(ctx, &QueueProperties{
		Name: queueName,
	})
}

// CreateQueue creates a queue with configurable properties.
func (ac *AdminClient) AddQueueWithProperties(ctx context.Context, properties *QueueProperties) (*AddQueueResponse, error) {
	props, resp, err := ac.createOrUpdateQueueImpl(ctx, properties, true)

	return &AddQueueResponse{
		RawResponse: resp,
		Value:       props,
	}, err
}

// GetQueue gets a queue by name.
func (ac *AdminClient) GetQueue(ctx context.Context, queueName string) (*GetQueueResponse, error) {
	var atomResp *atom.QueueEnvelope
	resp, err := ac.em.Get(ctx, "/"+queueName, &atomResp)

	if err != nil {
		return nil, err
	}

	props, err := newQueueProperties(atomResp.Title, &atomResp.Content.QueueDescription)

	if err != nil {
		return nil, err
	}

	return &GetQueueResponse{
		RawResponse: resp,
		Value:       props,
	}, nil
}

// GetQueueRuntimeProperties gets runtime properties of a queue, like the SizeInBytes, or ActiveMessageCount.
func (ac *AdminClient) GetQueueRuntimeProperties(ctx context.Context, queueName string) (*GetQueueRuntimePropertiesResponse, error) {
	var atomResp *atom.QueueEnvelope
	resp, err := ac.em.Get(ctx, "/"+queueName, &atomResp)

	if err != nil {
		return nil, err
	}

	props, err := newQueueRuntimeProperties(atomResp.Title, &atomResp.Content.QueueDescription), nil

	if err != nil {
		return nil, err
	}

	return &GetQueueRuntimePropertiesResponse{
		RawResponse: resp,
		Value:       props,
	}, nil
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

	if atom.NotFound(err) {
		return false, nil
	}

	return false, err
}

// UpdateQueue updates an existing queue.
func (ac *AdminClient) UpdateQueue(ctx context.Context, properties *QueueProperties) (*UpdateQueueResponse, error) {
	newProps, resp, err := ac.createOrUpdateQueueImpl(ctx, properties, false)

	if err != nil {
		return nil, err
	}

	return &UpdateQueueResponse{
		RawResponse: resp,
		Value:       newProps,
	}, err
}

// DeleteQueue deletes a queue.
func (ac *AdminClient) DeleteQueue(ctx context.Context, queueName string) (*DeleteQueueResponse, error) {
	resp, err := ac.em.Delete(ctx, "/"+queueName)
	defer atom.CloseRes(ctx, resp)
	return &DeleteQueueResponse{RawResponse: resp}, err
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
	PageResponse() *ListQueuesResponse

	// Err returns the last error encountered while paging.
	Err() error
}

type ListQueuesResponse struct {
	RawResponse *http.Response
	Value       []*QueueProperties
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
		innerPager: ac.getQueuePager(pageSize, skip),
	}
}

// ListQueuesRuntimePropertiesOptions can be used to configure the ListQueuesRuntimeProperties method.
type ListQueuesRuntimePropertiesOptions struct {
	// Top is the maximum size of each page of results.
	Top int
	// Skip is the starting index for the paging operation.
	Skip int
}

type ListQueuesRuntimePropertiesResponse struct {
	Value       []*QueueRuntimeProperties
	RawResponse *http.Response
}

// QueueRuntimePropertiesPager provides iteration over ListQueueRuntimeProperties pages.
type QueueRuntimePropertiesPager interface {
	// NextPage returns true if the pager advanced to the next page.
	// Returns false if there are no more pages or an error occurred.
	NextPage(context.Context) bool

	// PageResponse returns the current QueueRuntimeProperties.
	PageResponse() *ListQueuesRuntimePropertiesResponse

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
		innerPager: ac.getQueuePager(pageSize, skip),
	}
}

func (ac *AdminClient) createOrUpdateQueueImpl(ctx context.Context, props *QueueProperties, creating bool) (*QueueProperties, *http.Response, error) {
	if props == nil {
		return nil, nil, errors.New("properties are required and cannot be nil")
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
	resp, err := ac.em.Put(ctx, "/"+props.Name, env, &atomResp, mw...)

	if err != nil {
		return nil, nil, err
	}

	newProps, err := newQueueProperties(props.Name, &atomResp.Content.QueueDescription)

	if err != nil {
		return nil, nil, err
	}

	return newProps, resp, nil
}

// queuePropertiesPager provides iteration over QueueProperties pages.
type queuePropertiesPager struct {
	innerPager queueFeedPagerFunc

	lastErr      error
	lastResponse *ListQueuesResponse
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *queuePropertiesPager) NextPage(ctx context.Context) bool {
	p.lastResponse, p.lastErr = p.getNextPage(ctx)
	return p.lastResponse != nil
}

// PageResponse returns the current page.
func (p *queuePropertiesPager) PageResponse() *ListQueuesResponse {
	return p.lastResponse
}

// Err returns the last error encountered while paging.
func (p *queuePropertiesPager) Err() error {
	return p.lastErr
}

func (p *queuePropertiesPager) getNextPage(ctx context.Context) (*ListQueuesResponse, error) {
	feed, resp, err := p.innerPager(ctx)

	if err != nil {
		return nil, err
	}

	if feed == nil {
		return nil, nil
	}

	var all []*QueueProperties

	for _, env := range feed.Entries {
		props, err := newQueueProperties(env.Title, &env.Content.QueueDescription)

		if err != nil {
			return nil, err
		}

		all = append(all, props)
	}

	return &ListQueuesResponse{
		RawResponse: resp,
		Value:       all,
	}, nil
}

// queueRuntimePropertiesPager provides iteration over QueueRuntimeProperties pages.
type queueRuntimePropertiesPager struct {
	innerPager   queueFeedPagerFunc
	lastErr      error
	lastResponse *ListQueuesRuntimePropertiesResponse
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *queueRuntimePropertiesPager) NextPage(ctx context.Context) bool {
	p.lastResponse, p.lastErr = p.getNextPage(ctx)
	return p.lastResponse != nil
}

func (p *queueRuntimePropertiesPager) getNextPage(ctx context.Context) (*ListQueuesRuntimePropertiesResponse, error) {
	feed, resp, err := p.innerPager(ctx)

	if err != nil || feed == nil {
		return nil, err
	}

	var all []*QueueRuntimeProperties

	for _, entry := range feed.Entries {
		all = append(all, newQueueRuntimeProperties(entry.Title, &entry.Content.QueueDescription))
	}

	return &ListQueuesRuntimePropertiesResponse{
		RawResponse: resp,
		Value:       all,
	}, nil
}

// PageResponse returns the current page.
func (p *queueRuntimePropertiesPager) PageResponse() *ListQueuesRuntimePropertiesResponse {
	return p.lastResponse
}

// Err returns the last error encountered while paging.
func (p *queueRuntimePropertiesPager) Err() error {
	return p.lastErr
}

type queueFeedPagerFunc func(ctx context.Context) (*atom.QueueFeed, *http.Response, error)

func (ac *AdminClient) getQueuePager(top int, skip int) queueFeedPagerFunc {
	return func(ctx context.Context) (*atom.QueueFeed, *http.Response, error) {
		url := "/$Resources/Queues?"
		if top > 0 {
			url += fmt.Sprintf("&$top=%d", top)
		}

		if skip > 0 {
			url += fmt.Sprintf("&$skip=%d", skip)
		}

		var atomResp *atom.QueueFeed
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
	}

	return atom.WrapWithQueueEnvelope(qpr, tokenProvider)
}

func newQueueProperties(name string, desc *atom.QueueDescription) (*QueueProperties, error) {
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
		Name:                                name,
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
	}

	return queuePropsResult, nil
}

func newQueueRuntimeProperties(name string, desc *atom.QueueDescription) *QueueRuntimeProperties {
	return &QueueRuntimeProperties{
		Name:                           name,
		SizeInBytes:                    int64OrZero(desc.SizeInBytes),
		CreatedAt:                      dateTimeToTime(desc.CreatedAt),
		UpdatedAt:                      dateTimeToTime(desc.UpdatedAt),
		AccessedAt:                     dateTimeToTime(desc.AccessedAt),
		TotalMessageCount:              int64OrZero(desc.MessageCount),
		ActiveMessageCount:             int32OrZero(desc.CountDetails.ActiveMessageCount),
		DeadLetterMessageCount:         int32OrZero(desc.CountDetails.DeadLetterMessageCount),
		ScheduledMessageCount:          int32OrZero(desc.CountDetails.ScheduledMessageCount),
		TransferDeadLetterMessageCount: int32OrZero(desc.CountDetails.TransferDeadLetterMessageCount),
		TransferMessageCount:           int32OrZero(desc.CountDetails.TransferMessageCount),
	}
}

func dateTimeToTime(t *date.Time) time.Time {
	if t == nil {
		return time.Time{}
	}

	return t.Time
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
